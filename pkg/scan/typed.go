package scan

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/model"
	"golang.org/x/tools/go/packages"
)

// ScanTyped type-checks the module at root (go/types) and runs the same walk
// with type information: a router reached any way is recognized by its type,
// constants fold, and calls resolve exactly. The target's dependencies must
// resolve; packages that don't type-check are reported as problems and read
// from their syntax alone.
func ScanTyped(root string) (*Result, error) {
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("-types needs a module directory, not a file: %s", root)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports,
		Dir:   abs,
		Fset:  fset,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}
	res := &Result{Root: root, Mode: "types"}
	ix := newIndex(fset, abs)
	mods := findModules(abs)
	seen := map[string]bool{}
	sort.Slice(pkgs, func(i, j int) bool { return pkgs[i].PkgPath < pkgs[j].PkgPath })
	for _, lp := range pkgs {
		for _, e := range lp.Errors {
			res.Problems = append(res.Problems, model.Problem{File: errFile(abs, e.Pos), Message: firstLine(e.Msg)})
		}
		if len(lp.GoFiles) == 0 && len(lp.CompiledGoFiles) == 0 {
			continue
		}
		dir := ""
		if len(lp.GoFiles) > 0 {
			dir = filepath.Dir(lp.GoFiles[0])
		} else {
			dir = filepath.Dir(lp.CompiledGoFiles[0])
		}
		if hasSkippedDir(abs, dir) {
			continue
		}
		p := ix.pkgFor(dir, lp.PkgPath, lp.Name)
		if lp.TypesInfo != nil && len(lp.Syntax) > 0 && !lp.IllTyped {
			p.info = lp.TypesInfo
		} else if lp.TypesInfo != nil && len(lp.Syntax) > 0 {
			p.info = lp.TypesInfo // partial type info is still better than none
		}
		if len(lp.Syntax) > 0 {
			for i, f := range lp.Syntax {
				name := ""
				if i < len(lp.CompiledGoFiles) {
					name = lp.CompiledGoFiles[i]
				}
				if !strings.HasSuffix(name, ".go") || seen[name] {
					continue // cgo-generated or duplicate
				}
				seen[name] = true
				res.Files++
				ix.addFile(p, f, relPath(abs, name))
			}
			continue
		}
		// No syntax (the package failed to load): read it untyped.
		for _, name := range lp.GoFiles {
			if seen[name] || strings.HasSuffix(name, "_test.go") {
				continue
			}
			seen[name] = true
			res.Files++
			f, perr := parser.ParseFile(fset, name, nil, parser.SkipObjectResolution)
			if perr != nil && f == nil {
				continue
			}
			if p.importPath == "" {
				p.importPath = mods.importPath(dir)
			}
			ix.addFile(p, f, relPath(abs, name))
		}
	}
	if res.Files == 0 {
		return nil, fmt.Errorf("%w under %s", ErrNoGoFiles, root)
	}
	for _, p := range ix.pkgs {
		if p.info == nil {
			continue
		}
		for _, fc := range p.files {
			for _, d := range fc.file.Decls {
				fd, ok := d.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if obj := p.info.Defs[fd.Name]; obj != nil {
					key := fd.Name.Name
					if fd.Recv != nil && len(fd.Recv.List) > 0 {
						key = recvTypeName(fd.Recv.List[0].Type) + "." + key
					}
					if fi := p.funcs[key]; fi != nil && fi.decl == fd {
						ix.byObj[obj] = fi
					}
				}
			}
		}
	}
	walkAllTyped(ix, res)
	return res, nil
}

func walkAllTyped(ix *index, res *Result) {
	ix.typedTouch = func(fi *funcInfo, x ast.Expr) bool {
		info := fi.fc.pkg.info
		if info == nil {
			return false
		}
		t := info.TypeOf(x)
		return t != nil && routerFrameworkOf(t) != nil
	}
	walkAll(ix, res)
}

func hasSkippedDir(root, dir string) bool {
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return false
	}
	for _, seg := range strings.Split(filepath.ToSlash(rel), "/") {
		if seg != "." && skipDir(seg) {
			return true
		}
	}
	return false
}

func errFile(root, pos string) string {
	if pos == "" || pos == "-" {
		return ""
	}
	// pos is file:line:col; keep the file part relative to the root
	f := pos
	for i := 0; i < 2; i++ {
		if j := strings.LastIndexByte(f, ':'); j > 1 {
			f = f[:j]
		}
	}
	return relPath(root, f)
}
