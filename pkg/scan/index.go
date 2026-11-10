package scan

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"strconv"
	"strings"
)

// index is everything the walker needs to follow a router across functions,
// files and packages: declarations by name, struct fields of router type,
// package-level constants and router variables.
type index struct {
	fset  *token.FileSet
	root  string
	pkgs  []*pkgInfo
	byDir map[string]*pkgInfo
	byImp map[string]*pkgInfo
	// methodsByName holds every method in the tree, for calls whose receiver
	// type can't be determined syntactically.
	methodsByName map[string][]*funcInfo
	// byObj resolves a typed callee (-types) to its declaration.
	byObj map[types.Object]*funcInfo
	// typedTouch (-types) reports an expression of router type.
	typedTouch func(*funcInfo, ast.Expr) bool
}

type pkgInfo struct {
	dir        string
	importPath string
	name       string
	files      []*fileCtx
	funcs      map[string]*funcInfo // "F" or "T.M"
	structs    map[string]*structInfo
	// routerIfaces: in-tree interfaces with router verbs (oapi-codegen's
	// EchoRouter, an app's own Registrar); a function taking one registers
	// routes on whatever router its caller passes
	routerIfaces map[string]bool
	consts       map[string]constDef // package-level string constants and vars
	globals      map[string]*globalVar
	info         *types.Info // nil in AST mode or when type-checking failed
}

type fileCtx struct {
	pkg     *pkgInfo
	file    *ast.File
	rel     string                // path relative to the scan root, slash-separated
	fws     map[string]*framework // local import name -> router framework
	imports map[string]string     // local import name -> import path
}

type funcInfo struct {
	decl     *ast.FuncDecl
	fc       *fileCtx
	recvType string // "" for plain functions
	name     string
	// touches: the function can register routes or hand back a router, so a
	// call to it is worth following.
	touches bool
	callees []*funcInfo
	// assignsRouterState: assigns a router to a struct field or package var;
	// walked first so other functions find that router already built.
	assignsRouterState bool
	// variants: the same name declared again under other build tags
	variants []*funcInfo
}

type structInfo struct {
	fc       *fileCtx
	fields   map[string]ast.Expr // field name -> type
	embedded map[string]ast.Expr // embedded type's name -> type
	order    []string            // field names in declaration order, for positional literals
}

type constDef struct {
	expr ast.Expr
	fc   *fileCtx
}

type globalVar struct {
	name string
	typ  ast.Expr
	init ast.Expr
	fc   *fileCtx
}

func newIndex(fset *token.FileSet, root string) *index {
	return &index{
		fset: fset, root: root,
		byDir: map[string]*pkgInfo{}, byImp: map[string]*pkgInfo{},
		methodsByName: map[string][]*funcInfo{},
		byObj:         map[types.Object]*funcInfo{},
	}
}

func (ix *index) pkgFor(dir, importPath, name string) *pkgInfo {
	if p := ix.byDir[dir]; p != nil {
		return p
	}
	p := &pkgInfo{
		dir: dir, importPath: importPath, name: name,
		funcs: map[string]*funcInfo{}, structs: map[string]*structInfo{}, routerIfaces: map[string]bool{},
		consts: map[string]constDef{}, globals: map[string]*globalVar{},
	}
	ix.byDir[dir] = p
	if importPath != "" {
		ix.byImp[importPath] = p
	}
	ix.pkgs = append(ix.pkgs, p)
	return p
}

// addFile records a parsed file's imports and top-level declarations.
func (ix *index) addFile(p *pkgInfo, f *ast.File, rel string) *fileCtx {
	fc := &fileCtx{pkg: p, file: f, rel: rel, fws: map[string]*framework{}, imports: map[string]string{}}
	for _, imp := range f.Imports {
		ip, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		local := importName(ip)
		if imp.Name != nil {
			local = imp.Name.Name
		}
		if local == "_" || local == "." {
			continue
		}
		fc.imports[local] = ip
		if fw := frameworkForImport(ip); fw != nil {
			fc.fws[local] = fw
		}
	}
	p.files = append(p.files, fc)

	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			fi := &funcInfo{decl: d, fc: fc, name: d.Name.Name}
			key := d.Name.Name
			if d.Recv != nil && len(d.Recv.List) > 0 {
				fi.recvType = recvTypeName(d.Recv.List[0].Type)
				key = fi.recvType + "." + key
				ix.methodsByName[d.Name.Name] = append(ix.methodsByName[d.Name.Name], fi)
			}
			// the same function in build-tagged files (admin_linux.go and
			// admin_windows.go): keep every variant, each is walked
			if first, dup := p.funcs[key]; dup {
				first.variants = append(first.variants, fi)
			} else {
				p.funcs[key] = fi
			}
		case *ast.GenDecl:
			ix.addGenDecl(fc, d)
		}
	}
	return fc
}

func (ix *index) addGenDecl(fc *fileCtx, d *ast.GenDecl) {
	p := fc.pkg
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if it, ok := s.Type.(*ast.InterfaceType); ok {
				if routerMethods(it) {
					p.routerIfaces[s.Name.Name] = true
				}
				continue
			}
			st, ok := s.Type.(*ast.StructType)
			if !ok {
				continue
			}
			si := &structInfo{fc: fc, fields: map[string]ast.Expr{}, embedded: map[string]ast.Expr{}}
			for _, fld := range st.Fields.List {
				if len(fld.Names) == 0 {
					si.embedded[typeBaseName(fld.Type)] = fld.Type
					continue
				}
				for _, n := range fld.Names {
					si.fields[n.Name] = fld.Type
					si.order = append(si.order, n.Name)
				}
			}
			p.structs[s.Name.Name] = si
		case *ast.ValueSpec:
			for i, n := range s.Names {
				var init ast.Expr
				if i < len(s.Values) {
					init = s.Values[i]
				}
				if d.Tok == token.CONST {
					if init != nil {
						p.consts[n.Name] = constDef{expr: init, fc: fc}
					}
					continue
				}
				p.globals[n.Name] = &globalVar{name: n.Name, typ: s.Type, init: init, fc: fc}
				if init != nil {
					p.consts[n.Name] = constDef{expr: init, fc: fc}
				}
			}
		}
	}
}

// importName is the identifier an unaliased import usually binds: the last
// path segment, skipping a major-version suffix (.../echo/v4 -> echo) and
// a gopkg.in style ".vN".
func importName(ip string) string {
	seg := path.Base(ip)
	if isVersionSuffix(seg) {
		seg = path.Base(path.Dir(ip))
	}
	if i := strings.Index(seg, ".v"); i > 0 {
		seg = seg[:i]
	}
	seg = strings.TrimPrefix(seg, "go-")
	return strings.ReplaceAll(seg, "-", "")
}

func recvTypeName(x ast.Expr) string {
	return typeBaseName(x)
}

// typeBaseName strips pointers, generic instantiation and a package qualifier.
func typeBaseName(x ast.Expr) string {
	for {
		switch t := x.(type) {
		case *ast.StarExpr:
			x = t.X
		case *ast.ParenExpr:
			x = t.X
		case *ast.IndexExpr:
			x = t.X
		case *ast.IndexListExpr:
			x = t.X
		case *ast.Ident:
			return t.Name
		case *ast.SelectorExpr:
			return t.Sel.Name
		default:
			return ""
		}
	}
}

// routerType reports the framework of a type expression written in fc, e.g.
// `chi.Router`, `*gin.RouterGroup`, `*http.ServeMux`.
func (fc *fileCtx) routerType(x ast.Expr) *framework {
	for {
		switch t := x.(type) {
		case *ast.StarExpr:
			x = t.X
			continue
		case *ast.ParenExpr:
			x = t.X
			continue
		case *ast.SelectorExpr:
			id, ok := t.X.(*ast.Ident)
			if !ok {
				return nil
			}
			if fw := fc.fws[id.Name]; fw != nil && fw.types[t.Sel.Name] {
				return fw
			}
		}
		return nil
	}
}

// structOf resolves a type expression to a struct declared in the tree.
func (ix *index) structOf(fc *fileCtx, x ast.Expr) (*pkgInfo, string, *structInfo) {
	for {
		switch t := x.(type) {
		case *ast.StarExpr:
			x = t.X
			continue
		case *ast.ParenExpr:
			x = t.X
			continue
		case *ast.IndexExpr: // Registrar[T]
			x = t.X
			continue
		case *ast.IndexListExpr:
			x = t.X
			continue
		case *ast.Ident:
			if si := fc.pkg.structs[t.Name]; si != nil {
				return fc.pkg, t.Name, si
			}
		case *ast.SelectorExpr:
			if id, ok := t.X.(*ast.Ident); ok {
				if p := ix.byImp[fc.imports[id.Name]]; p != nil {
					if si := p.structs[t.Sel.Name]; si != nil {
						return p, t.Sel.Name, si
					}
				}
			}
		}
		return nil, "", nil
	}
}

// lookupFunc resolves a package-qualified or local function name.
func (ix *index) lookupFunc(fc *fileCtx, qual, name string) *funcInfo {
	if qual == "" {
		return fc.pkg.funcs[name]
	}
	ip, ok := fc.imports[qual]
	if !ok {
		return nil
	}
	if p := ix.byImp[ip]; p != nil {
		return p.funcs[name]
	}
	return nil
}

// computeTouches marks the functions worth following: those that mention a
// router type or framework package directly, plus everything that calls
// them. Calls into anything else are skipped, which keeps the walk linear-ish
// on large trees.
func (ix *index) computeTouches(resolve func(*funcInfo, *ast.CallExpr) *funcInfo, typed func(*funcInfo, ast.Expr) bool) {
	var all []*funcInfo
	for _, p := range ix.pkgs {
		for _, fi := range p.funcs {
			all = append(all, fi)
			all = append(all, fi.variants...)
		}
	}
	for _, fi := range all {
		fc := fi.fc
		seen := map[*funcInfo]bool{}
		direct := false
		if fi.recvType != "" {
			if si := fc.pkg.structs[fi.recvType]; si != nil && ix.structHasRouter(fc.pkg, si, 0) {
				direct = true
			}
		}
		for _, fld := range fi.decl.Type.Params.List {
			if ix.routerIface(fc, fld.Type) {
				direct = true
			}
		}
		ast.Inspect(fi.decl, func(n ast.Node) bool {
			if e, ok := n.(ast.Expr); ok && typed != nil && !direct && typed(fi, e) {
				direct = true
			}
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if id, ok := x.X.(*ast.Ident); ok {
					if fw := fc.fws[id.Name]; fw != nil && (fw.types[x.Sel.Name] || fw.ctors[x.Sel.Name] ||
						(fw == fwNetHTTP && (x.Sel.Name == "Handle" || x.Sel.Name == "HandleFunc"))) {
						direct = true
					}
				}
			case *ast.Ident:
				if g := fc.pkg.globals[x.Name]; g != nil && ix.globalIsRouter(g) {
					direct = true
				}
			case *ast.CallExpr:
				if registrationShaped(fc, x) || apiCtor(fc, x) || serveSinks[exprName(x.Fun)] {
					direct = true
				}
				if c := resolve(fi, x); c != nil && c != fi && !seen[c] {
					seen[c] = true
					fi.callees = append(fi.callees, c)
				}
			case *ast.CompositeLit:
				// newServer(h) { return &http.Server{Handler: RequireAuth(h)} }
				if exprName(x.Type) == "http.Server" {
					direct = true
				}
			case *ast.AssignStmt:
				for _, l := range x.Lhs {
					switch l := l.(type) {
					case *ast.SelectorExpr:
						fi.assignsRouterState = fi.assignsRouterState || mentionsRouter(fc, x.Rhs)
					case *ast.Ident:
						if fc.pkg.globals[l.Name] != nil && x.Tok == token.ASSIGN {
							fi.assignsRouterState = fi.assignsRouterState || mentionsRouter(fc, x.Rhs)
						}
					}
				}
			}
			return true
		})
		fi.touches = direct
	}
	for changed := true; changed; {
		changed = false
		for _, fi := range all {
			if fi.touches {
				continue
			}
			for _, c := range fi.callees {
				if c.touches {
					fi.touches = true
					changed = true
					break
				}
			}
		}
	}
}

// registrationShaped: x.GET("/path", ...) in a file that imports a router
// whose verbs include GET. Walked even when x can't be traced, so the route
// is reported as partial rather than never seen.
func registrationShaped(fc *fileCtx, c *ast.CallExpr) bool {
	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok || len(c.Args) < 2 {
		return false
	}
	if !pathShaped(c.Args[0]) {
		return false
	}
	for _, fw := range fc.fws {
		if _, ok := fw.verbs[sel.Sel.Name]; ok && fw != fwNetHTTP {
			return true
		}
	}
	return false
}

// apiCtor: a call that builds an API server served as one handler
// (grpc-gateway's NewServeMux, gqlgen's handler, a connect NewXHandler); a
// function that builds one is walked even if it touches no router.
func apiCtor(fc *fileCtx, c *ast.CallExpr) bool {
	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	ip, n := fc.imports[id.Name], sel.Sel.Name
	switch {
	case strings.Contains(ip, "grpc-gateway") && n == "NewServeMux",
		strings.Contains(ip, "gqlgen/graphql/handler"), strings.Contains(ip, "graph-gophers/graphql-go"),
		strings.HasSuffix(ip, "connect") && strings.HasPrefix(n, "New") && strings.HasSuffix(n, "Handler"):
		return true
	}
	return false
}

// pathShaped: "/x", or a path built from a prefix: baseURL + "/admin".
func pathShaped(x ast.Expr) bool {
	switch e := x.(type) {
	case *ast.BasicLit:
		return e.Kind == token.STRING && len(e.Value) >= 3 && e.Value[1] == '/'
	case *ast.BinaryExpr:
		return e.Op == token.ADD && (pathShaped(e.Y) || pathShaped(e.X))
	}
	return false
}

var routerVerbs = set("GET", "POST", "PUT", "DELETE", "PATCH", "Get", "Post", "Put", "Delete", "Patch",
	"Handle", "HandleFunc", "Any", "Add")

// routerMethods: at least two verb methods shaped like a registration,
// (path string, handler ...): GET(path string, h echo.HandlerFunc, ...).
// A cache's Get(key) / Delete(ctx, key) isn't one.
func routerMethods(it *ast.InterfaceType) bool {
	n := 0
	for _, m := range it.Methods.List {
		ft, ok := m.Type.(*ast.FuncType)
		if !ok || !registrationSig(ft) {
			continue
		}
		for _, nm := range m.Names {
			if routerVerbs[nm.Name] {
				n++
			}
		}
	}
	// one Handle/HandleFunc(pattern string, fn) is enough: that's a mux
	return n >= 2 || n == 1 && handleOnly(it)
}

func handleOnly(it *ast.InterfaceType) bool {
	for _, m := range it.Methods.List {
		for _, nm := range m.Names {
			if ft, ok := m.Type.(*ast.FuncType); ok && registrationSig(ft) && (nm.Name == "Handle" || nm.Name == "HandleFunc") {
				return true
			}
		}
	}
	return false
}

func registrationSig(ft *ast.FuncType) bool {
	var types []ast.Expr
	for _, f := range ft.Params.List {
		k := len(f.Names)
		if k == 0 {
			k = 1
		}
		for i := 0; i < k; i++ {
			types = append(types, f.Type)
		}
	}
	if len(types) < 2 {
		return false
	}
	if id, ok := types[0].(*ast.Ident); !ok || id.Name != "string" {
		return false
	}
	for _, t := range types[1:] {
		if e, ok := t.(*ast.Ellipsis); ok {
			t = e.Elt
		}
		if _, ok := t.(*ast.FuncType); ok {
			return true
		}
		if n := typeBaseName(t); strings.Contains(n, "Handler") || strings.Contains(n, "Handle") {
			return true
		}
	}
	return false
}

// routerIface: t names an in-tree router interface.
func (ix *index) routerIface(fc *fileCtx, t ast.Expr) bool {
	switch e := t.(type) {
	case *ast.StarExpr:
		return ix.routerIface(fc, e.X)
	case *ast.Ident:
		return fc.pkg.routerIfaces[e.Name]
	case *ast.SelectorExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			if p := ix.byImp[fc.imports[id.Name]]; p != nil {
				return p.routerIfaces[e.Sel.Name]
			}
		}
	}
	return false
}

// mentionsRouter: a cheap syntactic test for "this expression builds a router".
func mentionsRouter(fc *fileCtx, xs []ast.Expr) bool {
	found := false
	for _, x := range xs {
		ast.Inspect(x, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok {
					if fw := fc.fws[id.Name]; fw != nil && (fw.ctors[sel.Sel.Name] || fw.types[sel.Sel.Name]) {
						found = true
					}
				}
			}
			return !found
		})
	}
	return found
}

func (ix *index) globalIsRouter(g *globalVar) bool {
	if g.typ != nil && g.fc.routerType(g.typ) != nil {
		return true
	}
	return g.init != nil && mentionsRouter(g.fc, []ast.Expr{g.init})
}

func (ix *index) structHasRouter(p *pkgInfo, si *structInfo, depth int) bool {
	if depth > 3 {
		return false
	}
	for _, t := range si.fields {
		if si.fc.routerType(t) != nil {
			return true
		}
	}
	for _, t := range si.embedded {
		if si.fc.routerType(t) != nil {
			return true
		}
		if ep, _, es := ix.structOf(si.fc, t); es != nil && ix.structHasRouter(ep, es, depth+1) {
			return true
		}
	}
	return false
}
