// Package scan resolves the HTTP routes a Go service registers, each with its
// middleware chain. The default scan only parses source, so it needs no build.
package scan

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/model"
)

// Result is what a scan saw, not just what it found: a scan with Problems is
// incomplete and must not be read as clean.
type Result struct {
	Root     string          `json:"root"`
	Mode     string          `json:"mode"`
	Files    int             `json:"files"`
	Problems []model.Problem `json:"problems"`
	Routes   []model.Route   `json:"routes"`
	// Untraced lists registrations on a value authzmap couldn't trace to a
	// router. Their routes are included, marked partial.
	Untraced []model.Site `json:"untraced"`
	// Unseen lists route-like calls no supported router accounts for (a
	// router authzmap doesn't model). They are not analysed.
	Unseen []model.Site `json:"unseen"`
	Notes  []string     `json:"notes,omitempty"`
}

// Complete reports whether every file parsed (and, with -types, every
// package type-checked) and no route-like code went unanalysed.
func (r *Result) Complete() bool { return len(r.Problems) == 0 && len(r.Unseen) == 0 }

// ErrNoGoFiles is returned when the target holds no Go source at all.
var ErrNoGoFiles = errors.New("no Go files found")

// Scan parses every non-test .go file under root (a directory or a single
// file) and resolves the routes it registers.
func Scan(root string) (*Result, error) {
	fset := token.NewFileSet()
	res := &Result{Root: root, Mode: "ast"}
	files, err := goFiles(root)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("%w under %s", ErrNoGoFiles, root)
	}
	base := root
	if st, err := os.Stat(root); err == nil && !st.IsDir() {
		base = filepath.Dir(root)
	}
	// module directories are matched by prefix, so both sides must be absolute
	if abs, err := filepath.Abs(base); err == nil {
		base = abs
	}
	mods := findModules(base)
	ix := newIndex(fset, base)
	for _, f := range files {
		res.Files++
		if abs, err := filepath.Abs(f); err == nil {
			f = abs
		}
		rel := relPath(base, f)
		file, perr := parser.ParseFile(fset, f, nil, parser.SkipObjectResolution)
		if perr != nil {
			// the parser's message starts with the absolute path; File has it already
			msg := strings.TrimPrefix(firstLine(perr.Error()), f+":")
			res.Problems = append(res.Problems, model.Problem{File: rel, Message: msg})
			if file == nil {
				continue
			}
		}
		dir := filepath.Dir(f)
		p := ix.pkgFor(dir, mods.importPath(dir), file.Name.Name)
		ix.addFile(p, file, rel)
	}
	walkAll(ix, res)
	return res, nil
}

// goFiles lists the non-test .go files under root, skipping dependency and
// fixture trees the go tool would skip too.
func goFiles(root string) ([]string, error) {
	st, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		if !strings.HasSuffix(root, ".go") {
			return nil, fmt.Errorf("%s is not a Go file or directory", root)
		}
		return []string{root}, nil
	}
	var out []string
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p != root && skipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go") {
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

// skipDir skips dependency trees and test scaffolding: fixtures and mocks
// stand up fake servers whose routes aren't the service's.
func skipDir(name string) bool {
	switch name {
	case "vendor", "node_modules", "testdata", "fixtures", "fixture", "mocks", "mock", "testutil", "testutils", "testhelpers", "testhelper":
		return true
	}
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

func relPath(base, p string) string {
	if r, err := filepath.Rel(base, p); err == nil {
		return filepath.ToSlash(r)
	}
	return filepath.ToSlash(p)
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// modules maps directories to Go module paths so pkg.F(...) calls can be
// followed into the package that declares F.
type modules map[string]string

func findModules(base string) modules {
	m := modules{}
	// a go.mod above the scan root still names the packages inside it
	for d := base; ; {
		if mp := readModule(filepath.Join(d, "go.mod")); mp != "" {
			m[d] = mp
			break
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	_ = filepath.WalkDir(base, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if p != base && skipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() == "go.mod" {
			if mp := readModule(p); mp != "" {
				m[filepath.Dir(p)] = mp
			}
		}
		return nil
	})
	return m
}

func readModule(gomod string) string {
	f, err := os.Open(gomod)
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "module") {
			mp := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return strings.Trim(mp, `"`)
		}
	}
	return ""
}

func (m modules) importPath(dir string) string {
	best := ""
	for md := range m {
		if (dir == md || strings.HasPrefix(dir, md+string(filepath.Separator))) && len(md) > len(best) {
			best = md
		}
	}
	if best == "" {
		return ""
	}
	rel, err := filepath.Rel(best, dir)
	if err != nil {
		return ""
	}
	if rel == "." {
		return m[best]
	}
	return path.Join(m[best], filepath.ToSlash(rel))
}

// walkAll runs the route walk over an indexed tree and fills res.
//
// It walks twice. The first pass walks every function and records which
// ones are reached from a caller; the second starts clean and walks only
// the entry points, following calls into the rest. A function walked on its
// own and again from its caller would otherwise apply its Use() calls twice
// to a router shared through a struct field.
func walkAll(ix *index, res *Result) {
	first := newRun(ix)
	ix.computeTouches(func(fi *funcInfo, call *ast.CallExpr) *funcInfo {
		w := &walker{r: first, ix: ix, fc: fi.fc, env: newEnv(), consumed: map[*ast.FuncLit]bool{}}
		if c := w.resolveCall(call); c != nil {
			return c
		}
		// resolveCall only accepts a uniquely named method once it is known
		// to touch routers; that is what is being computed here, so accept
		// it regardless or the result depends on map order
		if sel, ok := unparen(call.Fun).(*ast.SelectorExpr); ok {
			if ms := ix.methodsByName[sel.Sel.Name]; len(ms) == 1 {
				return ms[0]
			}
		}
		return nil
	}, ix.typedTouch)
	fns := orderedFuncs(ix)
	walkPass(ix, first, fns, nil)
	entry := first.entryPoints(fns)

	r := newRun(ix)
	walkPass(ix, r, fns, entry)

	res.Routes = r.finish()
	res.Untraced = r.untraced
	if res.Untraced == nil {
		res.Untraced = []model.Site{}
	}
	if res.Problems == nil {
		res.Problems = []model.Problem{}
	}
	var note string
	var mounts []model.Site
	for _, m := range r.mounts {
		mounts = append(mounts, model.Site{File: m.file, Line: m.line})
	}
	// calls followed into a function that registered routes were traced
	mounts = append(mounts, r.viaSites...)
	res.Unseen, note = unseenCalls(ix, res.Routes, res.Untraced, mounts)
	if note != "" {
		res.Notes = append(res.Notes, note)
	}
	if r.budgetHit {
		res.Problems = append(res.Problems, model.Problem{Message: "too many nested calls to follow them all; some routes may be missing or lack their caller's middleware"})
	}
	for _, c := range r.cuts {
		res.Problems = append(res.Problems, model.Problem{File: c.File, Message: c.Call + " was cut off (too deep, or the walk budget ran out); routes it registers deeper down are missing"})
	}
}

// orderedFuncs lists the functions worth walking: init() and functions that
// build router state first, then the rest in source order.
func orderedFuncs(ix *index) []*funcInfo {
	var fns []*funcInfo
	for _, p := range ix.pkgs {
		for _, fi := range p.funcs {
			for _, v := range append([]*funcInfo{fi}, fi.variants...) {
				if v.touches && v.decl.Body != nil {
					fns = append(fns, v)
				}
			}
		}
	}
	rank := func(fi *funcInfo) int {
		switch {
		case fi.name == "init" && fi.recvType == "":
			return 0
		case fi.assignsRouterState:
			return 1
		}
		return 2
	}
	sort.SliceStable(fns, func(i, j int) bool {
		a, b := fns[i], fns[j]
		if rank(a) != rank(b) {
			return rank(a) < rank(b)
		}
		if a.fc.rel != b.fc.rel {
			return a.fc.rel < b.fc.rel
		}
		return a.decl.Pos() < b.decl.Pos()
	})
	return fns
}

// walkPass walks package-level initializers, then each function in fns
// (only those in only, when it's set).
func walkPass(ix *index, r *run, fns []*funcInfo, only map[*funcInfo]bool) {
	consumed := map[*ast.FuncLit]bool{}
	for _, p := range ix.pkgs {
		names := make([]string, 0, len(p.globals))
		for n := range p.globals {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, name := range names {
			g := p.globals[name]
			if g.init == nil {
				continue
			}
			// The initializer is its own origin, so a function it calls
			// (a cobra command's Run closure calling RegisterRoutes) counts
			// as reached from a caller.
			origin := &funcInfo{fc: g.fc, name: "var " + name, decl: &ast.FuncDecl{Name: ast.NewIdent(name)}}
			w := &walker{r: r, ix: ix, fc: g.fc, env: newEnv(), origin: origin, consumed: consumed}
			if fl, ok := unparen(g.init).(*ast.FuncLit); ok {
				consumed[fl] = true
				w.walkClosure(fl, nil, g.fc)
				continue
			}
			if n := w.eval(g.init); n != nil {
				r.globals[p.dir+":"+name] = n
			}
			w.sweep(g.init)
		}
	}
	for _, fi := range fns {
		if only != nil && !only[fi] {
			continue
		}
		e := newEnv()
		w := &walker{r: r, ix: ix, fc: fi.fc, env: e, origin: fi, stack: []*funcInfo{fi}, consumed: consumed}
		if fi.decl.Recv != nil && len(fi.decl.Recv.List) > 0 && len(fi.decl.Recv.List[0].Names) > 0 {
			rn := fi.decl.Recv.List[0].Names[0].Name
			e.locals[rn] = true
			e.vtypes[rn] = typeRef{pkg: fi.fc.pkg, name: fi.recvType}
		}
		w.bindParams(fi.decl.Type, nil, fi.name)
		w.walkStmts(fi.decl.Body.List)
	}
}
