package scan

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strconv"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/authsig"
	"github.com/NotChaosuu/authzmap/pkg/model"
	"golang.org/x/tools/go/types/typeutil"
)

// Node kinds. kindInline is a With()/block view of another router: same
// sibling group, extra middleware.
const (
	kindRoot   = model.ScopeRoot
	kindGroup  = model.ScopeGroup
	kindParam  = model.ScopeParam
	kindInline = "inline"
)

const (
	maxInlineDepth = 12        // nested calls followed; deeper ones are cut and reported
	walkBudget     = 3_000_000 // statements, across the whole run
)

type mwEntry struct {
	name     string
	prefix   string // fiber: app.Use("/api", mw) only covers /api
	enforces bool   // its body rejects unauthenticated requests (401 / login redirect)
	passes   bool   // its body always hands the request on: whatever its name, it isn't a check
	// cond: added inside an if/loop/switch, so it may not run; routes
	// registered after that block see it as conditional.
	cond bool
	// from is the block view it was copied up from; routes in that block
	// already have it.
	from       int
	unverified bool // its code has too many branches to read
}

// node is one router value: an engine, a group, a With() view, a subrouter.
// Routes keep a pointer to their node and resolve the full prefix and chain
// at the end, so a router mounted after its routes were registered still
// gets the mount prefix and the parent's middleware.
type node struct {
	id        int
	fw        *framework
	kind      string
	parent    *node
	parentLen int // parent's middleware count when derived; -1 = all of it
	prefix    string
	mw        []mwEntry
	outer     []mwEntry // wrappers around the whole router: auth(r) passed to a server
	partial   string
	site      string
	shadowOf  *node
	absolute  bool // mounted without stripping: its paths are full paths
	// edges: each place the router is mounted, with its own prefix and the
	// wrappers around that mount; the routes exist once per edge
	edges []edge
	// listeners: the wrappers around each server that serves the router
	// directly; an open listener beside an authed one leaves the routes open
	listeners [][]mwEntry
}

type edge struct {
	parent   *node
	prefix   string
	absolute bool
	wrappers []mwEntry
}

type mountRec struct {
	sub     *node
	parent  *node
	pattern string // the path it's mounted under, as written
	file    string
	line    int
	origin  *funcInfo
}

func (n *node) canon() *node {
	for n.shadowOf != nil {
		n = n.shadowOf
	}
	return n
}

type routeRec struct {
	n       *node
	nmw     int
	method  string
	path    string
	dynamic bool
	handler string
	extra   []mwEntry
	partial string // set on synthetic routes
	anc     map[*node]int
	absPath bool // path is already the full path
	checks  bool // the handler itself rejects anonymous callers
	api     string
	file    string
	line    int
	origin  *funcInfo
}

type wrapInfo struct {
	n        *node
	wrappers []mwEntry
	strip    string
}

type run struct {
	ix         *index
	nextID     int
	recs       []routeRec
	globals    map[string]*node
	objKeys    map[types.Object]string
	wraps      map[ast.Expr]wrapInfo
	inlined    map[*funcInfo]map[*funcInfo]bool
	steps      int
	budgetHit  bool
	untraced   []model.Site
	untracedAt map[string]bool
	defaultMux map[*pkgInfo]*node
	enforce    map[fnKey]verdict
	enforceSet map[fnKey]bool
	used       map[string]bool       // identifiers named anywhere outside their own declaration
	inClosure  map[*ast.FuncLit]bool // local closures being walked where they're called
	viaSites   []model.Site          // calls that led into registrations
	lastFields map[string]*node      // router fields of the struct the last inlined call returned
	asMW       bool                  // the verdict being computed is for middleware
	bodyCache  map[bodyKey]verdict
	whyOf      map[*ast.BlockStmt]string // why a conditional verdict is conditional
	lastWhy    string
	mounts     []mountRec
	cuts       []model.Site // recursion the walk refused to follow
	cutAt      map[string]bool
	anc        map[*node]int // fiber: each ancestor's middleware count when the route being resolved was added
}

type env struct {
	routers  map[string]*node
	values   map[string]string
	strs     map[string]string
	bools    map[string]bool  // parameters bound to true/false at the call site
	ints     map[string]int64 // and to integer constants (a depth argument)
	tables   map[string]*ast.CompositeLit
	parent   *env            // the enclosing block's env
	declared map[string]bool // names declared in this block (:= or var)
	vtypes   map[string]typeRef
	locals   map[string]bool
	chains   map[string][]mwEntry
	handlers map[string]wrapInfo
	origins  map[string]ast.Expr   // var -> the expression it was assigned from
	enforce  map[string]verdict    // param bound to middleware whose body was read
	methods  map[string]methodVal  // del := r.Delete
	negroni  map[string]bool       // chain variables that are negroni stacks
	slices   map[string][]ast.Expr // chain := append(chain, RequireAuth()): a handler slice's elements
	hwraps   map[string][]mwEntry  // h = RequireAuth(h) on a handler that isn't a router
	literal  map[string]bool       // params bound to a struct literal: an unset field is nil
	extern   map[string]bool       // vars whose type comes from outside the tree
}

// methodVal is a router's registration method taken as a value.
type methodVal struct {
	n    *node
	name string
}

type typeRef struct {
	pkg  *pkgInfo
	name string
}

func newEnv() *env {
	return &env{
		routers: map[string]*node{}, values: map[string]string{}, strs: map[string]string{}, bools: map[string]bool{}, ints: map[string]int64{}, tables: map[string]*ast.CompositeLit{},
		vtypes: map[string]typeRef{}, locals: map[string]bool{}, chains: map[string][]mwEntry{},
		handlers: map[string]wrapInfo{}, origins: map[string]ast.Expr{}, enforce: map[string]verdict{}, methods: map[string]methodVal{}, negroni: map[string]bool{}, slices: map[string][]ast.Expr{}, hwraps: map[string][]mwEntry{}, literal: map[string]bool{},
		extern: map[string]bool{}, declared: map[string]bool{},
	}
}

type walker struct {
	r         *run
	ix        *index
	fc        *fileCtx
	env       *env
	origin    *funcInfo
	stack     []*funcInfo
	returned  []*node
	retFields map[string]*node // router fields of a struct literal it returned
	consumed  map[*ast.FuncLit]bool
}

func newRun(ix *index) *run {
	return &run{
		ix: ix, globals: map[string]*node{}, objKeys: map[types.Object]string{},
		wraps: map[ast.Expr]wrapInfo{}, inlined: map[*funcInfo]map[*funcInfo]bool{},
		untracedAt: map[string]bool{}, defaultMux: map[*pkgInfo]*node{},
		enforce: map[fnKey]verdict{}, enforceSet: map[fnKey]bool{}, bodyCache: map[bodyKey]verdict{}, whyOf: map[*ast.BlockStmt]string{}, inClosure: map[*ast.FuncLit]bool{},
		cutAt: map[string]bool{},
	}
}

func (r *run) newNode(fw *framework, kind string, site string) *node {
	r.nextID++
	return &node{id: r.nextID, fw: fw, kind: kind, parentLen: -1, site: site}
}

func (w *walker) site(pos token.Pos) string {
	return w.fc.rel + ":" + strconv.Itoa(w.ix.fset.Position(pos).Line)
}

func (w *walker) line(pos token.Pos) int { return w.ix.fset.Position(pos).Line }

// with returns a walker for evaluating constants declared in another file.
func (w *walker) with(fc *fileCtx) *walker {
	c := *w
	c.fc = fc
	c.env = newEnv()
	return &c
}

func (w *walker) child(fc *fileCtx, e *env) *walker {
	return &walker{r: w.r, ix: w.ix, fc: fc, env: e, origin: w.origin, stack: w.stack, consumed: w.consumed}
}

// shadow copies env for a nested block. Router values become inline views,
// so a Use() inside an if doesn't leak onto routes after the block.
func (w *walker) shadow() *env {
	e := newEnv()
	e.parent = w.env
	for k, n := range w.env.routers {
		s := w.r.newNode(n.fw, kindInline, n.site)
		s.parent, s.parentLen, s.shadowOf = n, len(n.mw), n
		e.routers[k] = s
	}
	for k, v := range w.env.values {
		e.values[k] = v
	}
	for k, v := range w.env.strs {
		e.strs[k] = v
	}
	for k, v := range w.env.bools {
		e.bools[k] = v
	}
	for k, v := range w.env.ints {
		e.ints[k] = v
	}
	for k, v := range w.env.tables {
		e.tables[k] = v
	}
	for k, v := range w.env.vtypes {
		e.vtypes[k] = v
	}
	for k, v := range w.env.locals {
		e.locals[k] = v
	}
	for k, v := range w.env.chains {
		e.chains[k] = v
	}
	for k, v := range w.env.handlers {
		e.handlers[k] = v
	}
	for k, v := range w.env.origins {
		e.origins[k] = v
	}
	for k, v := range w.env.methods {
		e.methods[k] = v
	}
	for k, v := range w.env.negroni {
		e.negroni[k] = v
	}
	for k, v := range w.env.slices {
		e.slices[k] = v
	}
	for k, v := range w.env.hwraps {
		e.hwraps[k] = v
	}
	for k, v := range w.env.literal {
		e.literal[k] = v
	}
	for k, v := range w.env.enforce {
		e.enforce[k] = v
	}
	for k, v := range w.env.extern {
		e.extern[k] = v
	}
	return e
}

// tick counts work. Past the budget, calls stop being followed (each
// function is still walked on its own), so a huge tree degrades to partial
// chains rather than missing routes.
func (w *walker) tick() bool {
	w.r.steps++
	if w.r.steps > walkBudget {
		w.r.budgetHit = true
	}
	return true
}

// ---- statements ----

func (w *walker) walkBlock(b *ast.BlockStmt) { w.walkBlockCond(b, true) }

// walkBlockCond walks a nested block. cond is false for a bare { } block,
// which always runs, so its Use() calls carry over unconditionally.
func (w *walker) walkBlockCond(b *ast.BlockStmt, cond bool) {
	if b == nil {
		return
	}
	saved := w.env
	w.env = w.shadow()
	w.walkStmts(b.List)
	w.unshadowCond(saved, cond)
}

// unshadow leaves a nested block: middleware its routers gained there is
// carried to the outer routers marked conditional.
func (w *walker) unshadow(outer *env) { w.unshadowCond(outer, true) }

func (w *walker) unshadowCond(outer *env, cond bool) {
	for _, s := range w.env.routers {
		if s.shadowOf == nil || len(s.mw) == 0 {
			continue
		}
		for _, m := range s.mw {
			m.cond = m.cond || cond
			m.from = s.id
			s.shadowOf.mw = append(s.shadowOf.mw, m)
		}
	}
	w.env = outer
}

func (w *walker) walkStmts(list []ast.Stmt) {
	for _, st := range list {
		w.walkStmt(st)
	}
}

func (w *walker) walkStmt(st ast.Stmt) {
	if st == nil || !w.tick() {
		return
	}
	switch s := st.(type) {
	case *ast.ExprStmt:
		if call, ok := unparen(s.X).(*ast.CallExpr); ok {
			w.evalCall(call)
		}
		w.sweep(s.X)
	case *ast.AssignStmt:
		w.assign(s)
		for _, x := range s.Rhs {
			w.sweep(x)
		}
	case *ast.DeclStmt:
		w.declStmt(s)
	case *ast.ReturnStmt:
		for _, x := range s.Results {
			wi := w.sink(x)
			if wi.n != nil {
				w.returned = append(w.returned, wi.n)
			}
			// return &http.Server{Handler: RequireAuth(hd)}: that server
			// serves the wrapped handler
			if cl := compositeOf(x); cl != nil && hasHandlerField(cl) {
				w.compositeAssign("", cl)
			}
			// return &Router{group: g, auth: m}: a wrapper struct carrying a
			// router; the caller's variable gets its fields
			if cl := compositeOf(x); cl != nil && wi.n == nil {
				for _, f := range w.literalFields(cl) {
					if n := w.eval(f.value); n != nil {
						if w.retFields == nil {
							w.retFields = map[string]*node{}
						}
						w.retFields[f.name] = n
					}
				}
			}
			w.sweep(x)
		}
	case *ast.GoStmt:
		w.evalCall(s.Call)
		w.sweep(s.Call)
	case *ast.DeferStmt:
		w.evalCall(s.Call)
		w.sweep(s.Call)
	case *ast.BlockStmt:
		w.walkBlockCond(s, false)
	case *ast.IfStmt:
		// `if persistable { r.Put(...) }` with persistable bound to false at
		// the call site registers nothing
		if v, ok := w.constBool(s.Cond); ok && s.Init == nil {
			if v {
				w.walkBlockCond(s.Body, false)
			} else if s.Else != nil {
				w.walkStmt(s.Else)
			}
			return
		}
		saved := w.env
		w.env = w.shadow()
		w.walkStmt(s.Init)
		w.sweep(s.Cond)
		w.walkBlock(s.Body)
		w.walkStmt(s.Else)
		w.unshadow(saved)
	case *ast.ForStmt:
		saved := w.env
		w.env = w.shadow()
		w.walkStmt(s.Init)
		w.sweep(s.Cond)
		w.walkBlock(s.Body)
		w.walkStmt(s.Post)
		w.unshadow(saved)
	case *ast.RangeStmt:
		w.sweep(s.X)
		if !w.rangeTable(s) {
			w.walkBlock(s.Body)
		}
	case *ast.SwitchStmt:
		saved := w.env
		w.env = w.shadow()
		w.walkStmt(s.Init)
		w.sweep(s.Tag)
		w.walkClauses(s.Body)
		w.unshadow(saved)
	case *ast.TypeSwitchStmt:
		saved := w.env
		w.env = w.shadow()
		w.walkStmt(s.Init)
		w.walkClauses(s.Body)
		w.unshadow(saved)
	case *ast.SelectStmt:
		w.walkClauses(s.Body)
	case *ast.LabeledStmt:
		w.walkStmt(s.Stmt)
	case *ast.SendStmt:
		w.sweep(s.Value)
	}
}

func (w *walker) walkClauses(b *ast.BlockStmt) {
	for _, c := range b.List {
		switch cc := c.(type) {
		case *ast.CaseClause:
			for _, x := range cc.List {
				w.sweep(x)
			}
			w.walkBlock(&ast.BlockStmt{List: cc.Body})
		case *ast.CommClause:
			w.walkStmt(cc.Comm)
			w.walkBlock(&ast.BlockStmt{List: cc.Body})
		}
	}
}

// sweep walks function literals inside an expression that nothing else
// consumed: handlers passed to unknown helpers, goroutines, callbacks.
func (w *walker) sweep(x ast.Node) {
	if x == nil {
		return
	}
	ast.Inspect(x, func(n ast.Node) bool {
		fl, ok := n.(*ast.FuncLit)
		if !ok {
			return true
		}
		if !w.consumed[fl] {
			w.consumed[fl] = true
			w.walkClosure(fl, nil, w.fc)
		}
		return false
	})
}

// walkClosure walks a function literal with the enclosing env visible.
// bind maps parameter positions to routers (a Group() closure's `r`).
func (w *walker) walkClosure(fl *ast.FuncLit, bind map[int]*node, fc *fileCtx) {
	sub := w.child(fc, w.shadow())
	sub.bindParams(fl.Type, bind, "closure")
	sub.walkStmts(fl.Body.List)
}

// bindParams binds a function's parameters. Routers passed by the caller are
// in bind; any other router-typed parameter becomes a partial router, since
// what the caller applied to it isn't known.
func (w *walker) bindParams(ft *ast.FuncType, bind map[int]*node, fname string) {
	if ft == nil || ft.Params == nil {
		return
	}
	w.externParams(ft)
	i := 0
	for _, fld := range ft.Params.List {
		names := fld.Names
		if len(names) == 0 {
			i++
			continue
		}
		for _, nm := range names {
			w.env.locals[nm.Name] = true
			if n := bind[i]; n != nil {
				w.env.routers[nm.Name] = n
			} else if fw := w.paramRouterType(fld.Type); fw != nil && nm.Name != "_" {
				n := w.r.newNode(fw, kindParam, w.site(nm.Pos()))
				n.partial = fmt.Sprintf("router arrives as parameter %s of %s; middleware applied by the caller isn't traced", nm.Name, fname)
				w.env.routers[nm.Name] = n
			}
			if _, tn, si := w.ix.structOf(w.fc, fld.Type); si != nil {
				p, _, _ := w.ix.structOf(w.fc, fld.Type)
				w.env.vtypes[nm.Name] = typeRef{pkg: p, name: tn}
			}
			i++
		}
	}
}

func (w *walker) paramRouterType(t ast.Expr) *framework {
	if fw := w.fc.routerType(t); fw != nil {
		return fw
	}
	if info := w.fc.pkg.info; info != nil {
		if tv, ok := info.Types[t]; ok && tv.Type != nil {
			return routerFrameworkOf(tv.Type)
		}
	}
	return nil
}

func (w *walker) declStmt(s *ast.DeclStmt) {
	gd, ok := s.Decl.(*ast.GenDecl)
	if !ok {
		return
	}
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, nm := range vs.Names {
			w.env.locals[nm.Name] = true
			w.env.declared[nm.Name] = true
			if i < len(vs.Values) {
				w.assignOne(nm, vs.Values[i], token.DEFINE)
				w.sweep(vs.Values[i])
			} else if vs.Type != nil {
				if p, tn, si := w.ix.structOf(w.fc, vs.Type); si != nil {
					w.env.vtypes[nm.Name] = typeRef{pkg: p, name: tn}
				}
			}
		}
	}
}

func (w *walker) assign(s *ast.AssignStmt) {
	if s.Tok == token.DEFINE {
		for _, l := range s.Lhs {
			if id, ok := l.(*ast.Ident); ok {
				w.env.locals[id.Name] = true
				w.env.declared[id.Name] = true
			}
		}
	}
	if len(s.Lhs) == len(s.Rhs) {
		for i := range s.Lhs {
			w.assignOne(s.Lhs[i], s.Rhs[i], s.Tok)
		}
		return
	}
	if len(s.Rhs) == 1 {
		// r, err := newRouter()
		if n := w.eval(s.Rhs[0]); n != nil && len(s.Lhs) > 0 {
			w.bind(s.Lhs[0], n, s.Tok)
			return
		}
		// path, handler := adminv1connect.NewAdminServiceHandler(srv)
		for _, l := range s.Lhs {
			if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
				w.env.origins[id.Name] = s.Rhs[0]
			}
		}
	}
}

func (w *walker) assignOne(l, rhs ast.Expr, tok token.Token) {
	key := keyOf(l)
	if cl := compositeOf(rhs); cl != nil {
		if key != "" {
			w.env.tables[key] = cl
			if _, ok := cl.Type.(*ast.ArrayType); ok || handlerSliceType(cl.Type) {
				w.env.slices[key] = append([]ast.Expr(nil), cl.Elts...)
			}
		}
		w.compositeAssign(key, cl)
		return
	}
	// chain := make([]gin.HandlerFunc, 0); chain = append(chain, RequireAuth(), hs...)
	if c, ok := unparen(rhs).(*ast.CallExpr); ok && key != "" {
		switch exprName(c.Fun) {
		case "make":
			w.env.slices[key] = nil
		case "append":
			if len(c.Args) > 0 {
				base, ok := w.env.slices[keyOf(c.Args[0])]
				if cl := compositeOf(c.Args[0]); cl != nil {
					base, ok = cl.Elts, true // append([]gin.HandlerFunc{}, auth, h)
				}
				if ok || keyOf(c.Args[0]) == key {
					elems := append([]ast.Expr(nil), base...)
					for i, a := range c.Args[1:] {
						if c.Ellipsis.IsValid() && i == len(c.Args)-2 {
							elems = append(elems, w.env.slices[keyOf(a)]...)
							continue
						}
						elems = append(elems, a)
					}
					w.env.slices[key] = elems
				}
			}
		}
	}
	if sel, ok := l.(*ast.SelectorExpr); ok && sel.Sel.Name == "Handler" {
		// srv.Handler = auth(r)
		w.serve(rhs)
		return
	}
	w.r.lastFields = nil
	n := w.eval(rhs)
	if n != nil {
		w.bind(l, n, tok)
		w.wrapAssign(l, rhs, tok)
		return
	}
	if key == "" || key == "_" {
		return
	}
	// apiRouter := authz.NewRouter(api, m): its .group is the gin group
	if _, isCall := unparen(rhs).(*ast.CallExpr); isCall && len(w.r.lastFields) > 0 {
		for f, fn := range w.r.lastFields {
			w.env.routers[key+"."+f] = fn
		}
	}
	w.r.lastFields = nil
	delete(w.env.routers, key)
	delete(w.env.methods, key)
	// del := r.Delete: calling del registers on r
	if sel, ok := unparen(rhs).(*ast.SelectorExpr); ok {
		if recv := w.lookup(sel.X); recv != nil && frameworkMethod(recv.fw, sel.Sel.Name) {
			w.env.methods[key] = methodVal{n: recv, name: sel.Sel.Name}
		}
	}
	if s, ok := w.constString(rhs); ok {
		w.env.strs[key] = s
	} else {
		delete(w.env.strs, key)
	}
	if ch, ok := w.chainExpr(rhs); ok {
		w.env.chains[key] = ch
		if c, isCall := unparen(rhs).(*ast.CallExpr); isCall {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok && isNegroni(w.fc.imports[id.Name]) {
					w.env.negroni[key] = true
				}
			}
		}
	}
	if wi := w.wrapOf(rhs); wi.n != nil {
		w.env.handlers[key] = wi
		w.wrapAssign(l, rhs, tok)
	}
	// h = RequireAuth(h) on a handler that isn't a router (a grpc-gateway
	// mux): the wrapper travels with the variable
	selfWrap := false
	if c, ok := unparen(rhs).(*ast.CallExpr); ok && len(c.Args) == 1 && keyOf(c.Args[0]) != "" {
		if inner := keyOf(c.Args[0]); inner == key || w.env.hwraps[inner] != nil || w.env.origins[inner] != nil {
			if e := w.entryOf(c); e.name != "" && !handlerAdapters[exprName(c.Fun)] {
				w.env.hwraps[key] = append([]mwEntry{e}, w.env.hwraps[inner]...)
				selfWrap = inner == key
			}
		}
	}
	if fl, ok := unparen(rhs).(*ast.FuncLit); ok {
		w.env.origins[key] = fl // reg := func(p string) { ... }
	}
	if name := w.mwName(rhs); name != "" && name != key && !selfWrap {
		w.env.values[key] = name
		w.env.origins[key] = rhs
	}
	if tr, ok := w.typeOfValue(rhs); ok {
		w.env.vtypes[key] = tr
	}
}

// bind stores a router under an assignment target: a local, a package-level
// var, or a struct field (shared by every method of the type).
func (w *walker) bind(l ast.Expr, n *node, tok token.Token) {
	key := keyOf(l)
	if key == "" || key == "_" {
		return
	}
	delete(w.env.handlers, key)
	switch t := l.(type) {
	case *ast.Ident:
		if tok != token.DEFINE && !w.env.locals[t.Name] {
			if g := w.fc.pkg.globals[t.Name]; g != nil {
				w.r.globals[w.fc.pkg.dir+":"+t.Name] = n
			}
		}
	case *ast.SelectorExpr:
		if fk, _, _ := w.fieldKey(t); fk != "" {
			w.r.globals[fk] = n
		}
	}
	w.env.routers[key] = n
	// `x = build()` inside an if/for assigns the variable declared outside
	// it: carry the router up to the block that declared it, unless the
	// variable already held one there (`r := opts.Router; if r == nil { r =
	// chi.NewRouter() }` keeps the router it came with)
	if id, ok := l.(*ast.Ident); ok && tok != token.DEFINE && !w.env.declared[id.Name] {
		for p := w.env.parent; p != nil; p = p.parent {
			if p.routers[key] != nil {
				break
			}
			p.routers[key] = n
			if p.declared[id.Name] {
				break
			}
		}
	}
}

// wrapAssign: `handler = RequireAuth(handler)` makes the variable serve its
// router wrapped. Inside an if, the variable declared outside it is wrapped
// only when the condition holds: its wrappers are conditional there.
func (w *walker) wrapAssign(l, rhs ast.Expr, tok token.Token) {
	key := keyOf(l)
	call, ok := unparen(rhs).(*ast.CallExpr)
	if key == "" || key == "_" || !ok {
		return
	}
	wi := w.r.wraps[call]
	if wi.n == nil || len(wi.wrappers) == 0 {
		return
	}
	w.env.handlers[key] = wi
	id, ok := l.(*ast.Ident)
	if !ok || tok == token.DEFINE || w.env.declared[id.Name] {
		return
	}
	cw := wrapInfo{n: wi.n, strip: wi.strip}
	for _, m := range wi.wrappers {
		m.cond = true
		cw.wrappers = append(cw.wrappers, m)
	}
	for p := w.env.parent; p != nil; p = p.parent {
		p.handlers[key] = cw
		if p.declared[id.Name] {
			break
		}
	}
}

// bindFields binds a struct literal's fields under name.field in env e.
func (w *walker) bindFields(e *env, name string, cl *ast.CompositeLit) {
	e.literal[name] = true
	for _, f := range w.literalFields(cl) {
		k := name + "." + f.name
		if s, ok := w.constString(f.value); ok {
			e.strs[k] = s
		}
		if b, ok := w.constBool(f.value); ok {
			e.bools[k] = b
		}
		if ch, ok := w.chainExpr(f.value); ok {
			e.chains[k] = ch
		}
		if n := w.eval(f.value); n != nil {
			e.routers[k] = n
		}
		if sel, ok := unparen(f.value).(*ast.SelectorExpr); ok {
			if recv := w.lookup(sel.X); recv != nil && frameworkMethod(recv.fw, sel.Sel.Name) {
				e.methods[k] = methodVal{n: recv, name: sel.Sel.Name}
			}
		}
	}
}

type litField struct {
	name  string
	value ast.Expr
}

// literalFields names a struct literal's elements, positional ones (gr{m})
// by the struct's field order.
func (w *walker) literalFields(cl *ast.CompositeLit) []litField {
	var order []string
	if cl.Type != nil {
		if _, _, si := w.ix.structOf(w.fc, cl.Type); si != nil {
			order = si.order
		}
	}
	var out []litField
	for i, elt := range cl.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if id, ok := kv.Key.(*ast.Ident); ok {
				out = append(out, litField{id.Name, kv.Value})
			}
			continue
		}
		if i < len(order) {
			out = append(out, litField{order[i], elt})
		}
	}
	return out
}

func (w *walker) compositeAssign(key string, cl *ast.CompositeLit) {
	var tr typeRef
	if cl.Type != nil {
		if p, tn, si := w.ix.structOf(w.fc, cl.Type); si != nil {
			tr = typeRef{pkg: p, name: tn}
			if key != "" {
				w.env.vtypes[key] = tr
			}
		}
	}
	for _, elt := range cl.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		fname, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		if fname.Name == "Handler" {
			w.serve(kv.Value)
			continue
		}
		if ch, ok := w.chainExpr(kv.Value); ok && key != "" {
			w.env.chains[key+"."+fname.Name] = ch
		}
		if sel, ok := unparen(kv.Value).(*ast.SelectorExpr); ok && key != "" {
			if recv := w.lookup(sel.X); recv != nil && frameworkMethod(recv.fw, sel.Sel.Name) {
				w.env.methods[key+"."+fname.Name] = methodVal{n: recv, name: sel.Sel.Name}
			}
		}
		n := w.eval(kv.Value)
		if n == nil {
			continue
		}
		if tr.pkg != nil {
			w.r.globals["field:"+tr.pkg.dir+":"+tr.name+"."+fname.Name] = n
		}
		if key != "" {
			w.env.routers[key+"."+fname.Name] = n
		}
	}
}

// rangeTable walks a loop over a literal route table once per element,
// with the element's constant fields bound:
//
//	for _, rt := range routes { r.Method(rt.method, rt.path, rt.fn) }
//	for p, f := range map[string]http.HandlerFunc{"/admin/drop": h} { r.Post(p, f) }
func (w *walker) rangeTable(s *ast.RangeStmt) bool {
	cl := w.tableOf(s.X)
	if cl == nil || len(cl.Elts) == 0 || len(cl.Elts) > 512 {
		return false
	}
	keyName, valName := "", ""
	if id, ok := s.Key.(*ast.Ident); ok && id.Name != "_" {
		keyName = id.Name
	}
	if id, ok := s.Value.(*ast.Ident); ok && id.Name != "_" {
		valName = id.Name
	}
	var order []string
	if at, ok := cl.Type.(*ast.ArrayType); ok {
		if _, _, si := w.ix.structOf(w.fc, at.Elt); si != nil {
			order = si.order
		}
	}
	bound := false
	for _, elt := range cl.Elts {
		saved := w.env
		w.env = w.shadow()
		var val ast.Expr = elt
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if keyName != "" {
				if k, ok := w.constString(kv.Key); ok {
					w.env.strs[keyName] = k
				}
			}
			val = kv.Value
		}
		if valName != "" {
			if v, ok := w.constString(val); ok {
				w.env.strs[valName] = v
				bound = true
			} else if inner := compositeOf(val); inner != nil {
				// []Module{adminMod{}}: m.Routes(r) calls adminMod's method
				w.env.origins[valName] = val
				if tr, ok := w.typeOfValue(val); ok {
					w.env.vtypes[valName] = tr
				}
				for i, f := range inner.Elts {
					name, fx := "", f
					if kv, ok := f.(*ast.KeyValueExpr); ok {
						if id, ok := kv.Key.(*ast.Ident); ok {
							name, fx = id.Name, kv.Value
						}
					} else if i < len(order) {
						name = order[i]
					}
					if name == "" {
						continue
					}
					k := valName + "." + name
					if v, ok := w.constString(fx); ok {
						w.env.strs[k] = v
						bound = true
					}
					if b, ok := w.constBool(fx); ok {
						w.env.bools[k] = b
					}
				}
			} else {
				w.env.origins[valName] = val
			}
		}
		if keyName != "" {
			if _, ok := w.env.strs[keyName]; ok {
				bound = true
			}
		}
		w.walkStmts(s.Body.List)
		w.unshadowCond(saved, true)
	}
	_ = bound
	return true
}

// tableOf finds the literal a range loop iterates: a composite literal in
// place, or one assigned to a local or package variable.
func (w *walker) tableOf(x ast.Expr) *ast.CompositeLit {
	if cl := compositeOf(x); cl != nil {
		return cl
	}
	if call, ok := unparen(x).(*ast.CallExpr); ok {
		if fi := w.resolveCall(call); fi != nil && fi.fc.pkg == w.fc.pkg {
			return builtTable(fi)
		}
		return nil
	}
	id, ok := unparen(x).(*ast.Ident)
	if !ok {
		return nil
	}
	if cl := w.env.tables[id.Name]; cl != nil {
		return cl
	}
	if w.env.locals[id.Name] {
		return nil
	}
	if g := w.fc.pkg.globals[id.Name]; g != nil && g.init != nil {
		return compositeOf(g.init)
	}
	return nil
}

// builtTable reads the table a helper returns: a literal, or a slice it
// fills with append (rs = append(rs, route{"DELETE", "/admin/x", h})).
func builtTable(fi *funcInfo) *ast.CompositeLit {
	if fi.decl.Body == nil {
		return nil
	}
	var ret ast.Expr
	ast.Inspect(fi.decl.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if rs, ok := n.(*ast.ReturnStmt); ok && len(rs.Results) == 1 && ret == nil {
			ret = rs.Results[0]
		}
		return true
	})
	if cl := compositeOf(ret); cl != nil {
		return cl
	}
	id, ok := unparen(ret).(*ast.Ident)
	if !ok {
		return nil
	}
	out := &ast.CompositeLit{}
	ast.Inspect(fi.decl.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ValueSpec:
			for _, nm := range x.Names {
				if nm.Name == id.Name && x.Type != nil {
					out.Type = x.Type
				}
			}
		case *ast.AssignStmt:
			if len(x.Lhs) != 1 || len(x.Rhs) != 1 || keyOf(x.Lhs[0]) != id.Name {
				return true
			}
			if cl := compositeOf(x.Rhs[0]); cl != nil {
				out.Type = cl.Type
				out.Elts = append(out.Elts, cl.Elts...)
			}
			if c, ok := x.Rhs[0].(*ast.CallExpr); ok && exprName(c.Fun) == "append" && len(c.Args) > 1 && keyOf(c.Args[0]) == id.Name {
				out.Elts = append(out.Elts, c.Args[1:]...)
			}
		}
		return true
	})
	if len(out.Elts) == 0 {
		return nil
	}
	return out
}

func compositeOf(x ast.Expr) *ast.CompositeLit {
	switch e := x.(type) {
	case *ast.CompositeLit:
		return e
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			return compositeOf(e.X)
		}
	case *ast.ParenExpr:
		return compositeOf(e.X)
	}
	return nil
}

func unparen(x ast.Expr) ast.Expr {
	for {
		p, ok := x.(*ast.ParenExpr)
		if !ok {
			return x
		}
		x = p.X
	}
}

// ---- values ----

// eval returns the router an expression evaluates to, or nil.
func (w *walker) eval(x ast.Expr) *node {
	switch e := x.(type) {
	case *ast.ParenExpr:
		return w.eval(e.X)
	case *ast.StarExpr:
		return w.eval(e.X)
	case *ast.UnaryExpr:
		if e.Op == token.AND {
			if cl := compositeOf(e.X); cl != nil {
				return w.compositeRouter(cl)
			}
		}
		return w.eval(e.X)
	case *ast.CompositeLit:
		return w.compositeRouter(e)
	case *ast.Ident, *ast.SelectorExpr:
		return w.lookup(x)
	case *ast.CallExpr:
		return w.evalCall(e)
	case *ast.TypeAssertExpr:
		return w.eval(e.X)
	}
	return nil
}

// compositeRouter handles &http.ServeMux{} and fiber.App-like literals.
func (w *walker) compositeRouter(cl *ast.CompositeLit) *node {
	if cl.Type == nil {
		return nil
	}
	if fw := w.fc.routerType(cl.Type); fw != nil {
		return w.r.newNode(fw, kindRoot, w.site(cl.Pos()))
	}
	// &Router{Group: g} where Router embeds a router: the value is g
	_, _, si := w.ix.structOf(w.fc, cl.Type)
	if si == nil {
		return nil
	}
	for name, t := range si.embedded {
		if si.fc.routerType(t) == nil {
			continue
		}
		for _, elt := range cl.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if id, ok := kv.Key.(*ast.Ident); ok && id.Name == name {
					return w.eval(kv.Value)
				}
			}
		}
	}
	return nil
}

func (w *walker) lookup(x ast.Expr) *node {
	key := keyOf(x)
	if key != "" {
		if n := w.env.routers[key]; n != nil {
			return n
		}
	}
	switch e := x.(type) {
	case *ast.Ident:
		if !w.env.locals[e.Name] {
			if n := w.global(w.fc.pkg, e.Name); n != nil {
				return n
			}
		}
		if tr, ok := w.env.vtypes[e.Name]; ok {
			if n := w.embeddedRouter(tr, e.Name); n != nil {
				return n
			}
		}
	case *ast.SelectorExpr:
		if id, ok := e.X.(*ast.Ident); ok && !w.env.locals[id.Name] {
			if ip, ok := w.fc.imports[id.Name]; ok {
				if p := w.ix.byImp[ip]; p != nil {
					return w.global(p, e.Sel.Name)
				}
				if ip == "net/http" && e.Sel.Name == "DefaultServeMux" {
					return w.defaultMux()
				}
				return nil
			}
		}
		// ChiServerOptions{} passed in: options.BaseRouter is nil
		if id, ok := e.X.(*ast.Ident); ok && w.env.literal[id.Name] {
			return nil
		}
		if fk, ft, ffc := w.fieldKey(e); fk != "" {
			if n := w.r.globals[fk]; n != nil {
				return n
			}
			if fw := ffc.routerTypeOrTyped(ft, w); fw != nil {
				n := w.r.newNode(fw, kindParam, w.site(e.Pos()))
				n.partial = fmt.Sprintf("router in field %s is set up somewhere that wasn't traced", strings.TrimPrefix(fk[strings.LastIndexByte(fk, ':')+1:], "field:"))
				w.r.globals[fk] = n
				return n
			}
		}
	}
	// -types: anything whose static type is a router
	if info := w.fc.pkg.info; info != nil {
		if t := info.TypeOf(x); t != nil {
			if fw := routerFrameworkOf(t); fw != nil && key != "" {
				n := w.r.newNode(fw, kindParam, w.site(x.Pos()))
				n.partial = "router " + key + " comes from somewhere that wasn't traced"
				w.env.routers[key] = n
				return n
			}
		}
	}
	return nil
}

func (fc *fileCtx) routerTypeOrTyped(t ast.Expr, w *walker) *framework {
	if t == nil {
		return nil
	}
	if fw := fc.routerType(t); fw != nil {
		return fw
	}
	if info := fc.pkg.info; info != nil {
		if tv, ok := info.Types[t]; ok && tv.Type != nil {
			return routerFrameworkOf(tv.Type)
		}
	}
	return nil
}

// global returns a package-level router variable, creating a partial one if
// the var is declared with a router type but never assigned where we looked.
func (w *walker) global(p *pkgInfo, name string) *node {
	gk := p.dir + ":" + name
	if n := w.r.globals[gk]; n != nil {
		return n
	}
	g := p.globals[name]
	if g == nil || g.typ == nil {
		return nil
	}
	fw := g.fc.routerTypeOrTyped(g.typ, w)
	if fw == nil {
		return nil
	}
	n := w.r.newNode(fw, kindParam, g.fc.rel)
	n.partial = "package-level router " + name + " is set up somewhere that wasn't traced"
	w.r.globals[gk] = n
	return n
}

func (w *walker) defaultMux() *node {
	p := w.fc.pkg
	if n := w.r.defaultMux[p]; n != nil {
		return n
	}
	n := w.r.newNode(fwNetHTTP, kindRoot, "http.DefaultServeMux")
	w.r.defaultMux[p] = n
	return n
}

// fieldKey names a struct field holding a router so every method of the
// type shares it. It returns the field's declared type for lazy creation.
func (w *walker) fieldKey(sel *ast.SelectorExpr) (string, ast.Expr, *fileCtx) {
	if info := w.fc.pkg.info; info != nil {
		if s, ok := info.Selections[sel]; ok && s.Kind() == types.FieldVal {
			obj := s.Obj()
			k, ok := w.r.objKeys[obj]
			if !ok {
				k = fmt.Sprintf("field:%s.%s#%d", typeString(s.Recv()), obj.Name(), len(w.r.objKeys))
				w.r.objKeys[obj] = k
			}
			if n := w.r.globals[k]; n == nil {
				if fw := routerFrameworkOf(obj.Type()); fw != nil {
					n := w.r.newNode(fw, kindParam, w.site(sel.Pos()))
					n.partial = "router in field " + obj.Name() + " is set up somewhere that wasn't traced"
					w.r.globals[k] = n
				}
			}
			return k, nil, w.fc
		}
	}
	tr, ok := w.typeOf(sel.X)
	if !ok {
		return "", nil, nil
	}
	return w.fieldIn(tr, sel.Sel.Name, 0)
}

func (w *walker) fieldIn(tr typeRef, field string, depth int) (string, ast.Expr, *fileCtx) {
	if depth > 3 || tr.pkg == nil {
		return "", nil, nil
	}
	si := tr.pkg.structs[tr.name]
	if si == nil {
		return "", nil, nil
	}
	if t, ok := si.fields[field]; ok {
		return "field:" + tr.pkg.dir + ":" + tr.name + "." + field, t, si.fc
	}
	if t, ok := si.embedded[field]; ok {
		return "field:" + tr.pkg.dir + ":" + tr.name + "." + field, t, si.fc
	}
	for _, et := range si.embedded {
		if p, tn, es := w.ix.structOf(si.fc, et); es != nil {
			if k, t, fc := w.fieldIn(typeRef{pkg: p, name: tn}, field, depth+1); k != "" {
				return k, t, fc
			}
		}
	}
	return "", nil, nil
}

// embeddedRouter: `type Server struct{ *gin.Engine }` makes s.GET(...) a
// route on the embedded engine.
func (w *walker) embeddedRouter(tr typeRef, varName string) *node {
	si := tr.pkg.structs[tr.name]
	if si == nil {
		return nil
	}
	for name, t := range si.embedded {
		if fw := si.fc.routerType(t); fw != nil {
			k := "field:" + tr.pkg.dir + ":" + tr.name + "." + name
			if n := w.r.globals[k]; n != nil {
				return n
			}
			n := w.r.newNode(fw, kindParam, si.fc.rel)
			n.partial = "router embedded in " + tr.name + " is set up somewhere that wasn't traced"
			w.r.globals[k] = n
			return n
		}
	}
	return nil
}

// typeOf finds the struct type of a variable or field expression.
func (w *walker) typeOf(x ast.Expr) (typeRef, bool) {
	switch e := x.(type) {
	case *ast.Ident:
		tr, ok := w.env.vtypes[e.Name]
		return tr, ok
	case *ast.ParenExpr:
		return w.typeOf(e.X)
	case *ast.StarExpr:
		return w.typeOf(e.X)
	case *ast.SelectorExpr:
		tr, ok := w.typeOf(e.X)
		if !ok {
			return typeRef{}, false
		}
		_, t, fc := w.fieldIn(tr, e.Sel.Name, 0)
		if t == nil {
			return typeRef{}, false
		}
		if p, tn, si := w.ix.structOf(fc, t); si != nil {
			return typeRef{pkg: p, name: tn}, true
		}
	}
	return typeRef{}, false
}

// typeOfValue finds the struct type of &T{}, T{}, new(T) or NewT().
func (w *walker) typeOfValue(x ast.Expr) (typeRef, bool) {
	if cl := compositeOf(x); cl != nil && cl.Type != nil {
		if p, tn, si := w.ix.structOf(w.fc, cl.Type); si != nil {
			return typeRef{pkg: p, name: tn}, true
		}
	}
	call, ok := unparen(x).(*ast.CallExpr)
	if !ok {
		return w.typeOf(x)
	}
	if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "new" && len(call.Args) == 1 {
		if p, tn, si := w.ix.structOf(w.fc, call.Args[0]); si != nil {
			return typeRef{pkg: p, name: tn}, true
		}
	}
	if fi := w.resolveCall(call); fi != nil && fi.decl.Type.Results != nil && len(fi.decl.Type.Results.List) > 0 {
		if p, tn, si := w.ix.structOf(fi.fc, fi.decl.Type.Results.List[0].Type); si != nil {
			return typeRef{pkg: p, name: tn}, true
		}
	}
	return typeRef{}, false
}

// ---- calls ----

func (w *walker) evalCall(call *ast.CallExpr) *node {
	if !w.tick() {
		return nil
	}
	switch fun := unparen(call.Fun).(type) {
	case *ast.SelectorExpr:
		// (&http.Server{Handler: gw}).ListenAndServe()
		if cl := compositeOf(fun.X); cl != nil {
			w.compositeAssign("", cl)
		}
		if w.negroniCall(fun, call) {
			return nil
		}
		// g.post("/x", h) after g := reg{post: r.Post}
		if mv, ok := w.env.methods[keyOf(fun)]; ok {
			return w.routerMethod(mv.n, mv.name, call)
		}
		if id, ok := fun.X.(*ast.Ident); ok && !w.env.locals[id.Name] && w.env.routers[id.Name] == nil {
			if _, isImport := w.fc.imports[id.Name]; isImport {
				return w.pkgCall(call, id.Name, fun.Sel.Name)
			}
		}
		if root := gorillaRoot(call); root != nil {
			if g := w.peek(root); g != nil && g.fw == fwGorilla {
				return w.gorillaChain(call)
			}
		}
		if recv := w.eval(fun.X); recv != nil {
			if !frameworkMethod(recv.fw, fun.Sel.Name) {
				// a wrapper type's own method: r.SubGroup("/x") on a
				// type that embeds the router
				if fi := w.resolveCall(call); fi != nil {
					return w.inlineOnRouter(fi, call, recv)
				}
			}
			return w.routerMethod(recv, fun.Sel.Name, call)
		}
		if fi := w.resolveCall(call); fi != nil {
			return w.inline(fi, call, fun.X)
		}
		if n := w.untracedReceiver(fun, call); n != nil {
			return w.routerMethod(n, fun.Sel.Name, call)
		}
		return w.unknownCall(call)
	case *ast.IndexExpr, *ast.IndexListExpr:
		if fi := w.resolveCall(call); fi != nil {
			return w.inline(fi, call, nil)
		}
		return w.unknownCall(call)
	case *ast.Ident:
		if mv, ok := w.env.methods[fun.Name]; ok {
			return w.routerMethod(mv.n, mv.name, call)
		}
		if fl, ok := w.env.origins[fun.Name].(*ast.FuncLit); ok && w.env.locals[fun.Name] {
			return w.callClosure(fl, call)
		}
		if fi := w.resolveCall(call); fi != nil {
			return w.inline(fi, call, nil)
		}
		return w.unknownCall(call)
	case *ast.FuncLit:
		w.consumed[fun] = true
		bind := map[int]*node{}
		for i, a := range call.Args {
			if n := w.eval(a); n != nil {
				bind[i] = n
			}
		}
		sub := w.child(w.fc, w.shadow())
		sub.bindParams(fun.Type, bind, "func literal")
		sub.walkStmts(fun.Body.List)
		if len(sub.returned) > 0 {
			return sub.returned[0]
		}
		return nil
	}
	return w.unknownCall(call)
}

// pkgCall handles alias.F(...): router constructors, net/http's default mux,
// serving sinks, and functions declared elsewhere in the tree.
func (w *walker) pkgCall(call *ast.CallExpr, alias, name string) *node {
	if fw := w.fc.fws[alias]; fw != nil {
		if fw.ctors[name] {
			for _, a := range call.Args {
				w.sweep(a)
			}
			return w.r.newNode(fw, kindRoot, w.site(call.Pos()))
		}
		if fw == fwNetHTTP {
			switch name {
			case "Handle", "HandleFunc":
				w.netHTTPRegister(w.defaultMux(), call)
				return nil
			}
		}
	}
	if serveSinks[alias+"."+name] {
		for _, a := range call.Args {
			w.serve(a)
		}
		return nil
	}
	if fi := w.ix.lookupFunc(w.fc, alias, name); fi != nil {
		return w.inline(fi, call, nil)
	}
	if info := w.fc.pkg.info; info != nil {
		if fi := w.resolveCall(call); fi != nil {
			return w.inline(fi, call, nil)
		}
	}
	return w.unknownCall(call)
}

// unknownCall evaluates the arguments of a call we can't follow. When exactly
// one argument is a router, the call may be a wrapper (auth(r)); that is
// remembered and only applied if the result is served.
func (w *walker) unknownCall(call *ast.CallExpr) *node {
	var found []wrapInfo
	for _, a := range call.Args {
		if wi := w.evalValue(a); wi.n != nil {
			found = append(found, wi)
		}
	}
	if len(found) == 1 {
		name := exprName(call.Fun)
		wi := found[0]
		switch {
		case handlerAdapters[name]:
			w.r.wraps[call] = wi
		default:
			ents, ok := w.wrapperEntries(call)
			if !ok {
				ents = []mwEntry{w.entryOf(call)}
			}
			w.r.wraps[call] = wrapInfo{n: wi.n, wrappers: append(ents, wi.wrappers...), strip: wi.strip}
		}
	}
	return nil
}

// wrapperEntries reads an in-tree function that only wraps its handler
// parameter (secure(h) { return RequireAuth(h) }, or portainer's
// h = b.mwIsAdmin(h); h = b.mwAuthenticatedUser(h); return h) and returns
// the middleware it applies, outermost first.
func (w *walker) wrapperEntries(call *ast.CallExpr) ([]mwEntry, bool) {
	fi := w.resolveCall(call)
	if fi == nil {
		// bouncer.TeamLeaderAccess(h) on an interface: a method name that is
		// unique in the tree
		if sel, ok := unparen(call.Fun).(*ast.SelectorExpr); ok {
			if ms := w.ix.methodsByName[sel.Sel.Name]; len(ms) == 1 {
				fi = ms[0]
			}
		}
	}
	if fi == nil || fi.decl.Body == nil || len(fi.decl.Body.List) == 0 || len(fi.decl.Body.List) > 12 {
		return nil, false
	}
	// PublicAccess, OptionalAuth: the name says who may pass; keep it
	if optionalArg(fi.name) {
		return nil, false
	}
	params := map[string]bool{}
	for _, f := range fi.decl.Type.Params.List {
		for _, n := range f.Names {
			params[n.Name] = true
		}
	}
	r := w.with(fi.fc)
	if fi.recvType != "" && len(fi.decl.Recv.List[0].Names) > 0 {
		rn := fi.decl.Recv.List[0].Names[0].Name
		r.env.locals[rn] = true
		r.env.vtypes[rn] = typeRef{pkg: fi.fc.pkg, name: fi.recvType}
	}
	held := map[string][]mwEntry{} // the wrappers each variable now carries
	for p := range params {
		held[p] = []mwEntry{}
	}
	var unwrap func(x ast.Expr, depth int) ([]mwEntry, bool)
	unwrap = func(x ast.Expr, depth int) ([]mwEntry, bool) {
		if depth > 8 {
			return nil, false
		}
		switch e := unparen(x).(type) {
		case *ast.Ident:
			ents, ok := held[e.Name]
			return ents, ok
		case *ast.CallExpr:
			name := exprName(e.Fun)
			inner := r.innerHandlerArg(e.Args)
			if inner == nil {
				return nil, false
			}
			rest, ok := unwrap(inner, depth+1)
			if !ok {
				return nil, false
			}
			if handlerAdapters[name] {
				return rest, true
			}
			return append([]mwEntry{r.entryOf(e)}, rest...), true
		}
		return nil, false
	}
	for _, st := range fi.decl.Body.List {
		switch s := st.(type) {
		case *ast.AssignStmt:
			if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
				return nil, false
			}
			id, ok := s.Lhs[0].(*ast.Ident)
			if !ok {
				return nil, false
			}
			ents, ok := unwrap(s.Rhs[0], 0)
			if !ok {
				return nil, false
			}
			held[id.Name] = ents
		case *ast.ReturnStmt:
			if len(s.Results) != 1 {
				return nil, false
			}
			ents, ok := unwrap(s.Results[0], 0)
			if !ok || len(ents) == 0 {
				return nil, false
			}
			return ents, true
		default:
			return nil, false
		}
	}
	return nil, false
}

// evalValue evaluates an expression used as a handler: a router, possibly
// wrapped (auth(r), http.StripPrefix("/api", r), alice.New(a).Then(r)).
func (w *walker) evalValue(x ast.Expr) wrapInfo {
	switch e := unparen(x).(type) {
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok && len(e.Args) == 1 {
			switch sel.Sel.Name {
			case "Then", "ThenFunc", "Handler", "HandlerFunc":
				if ch, ok := w.chainExpr(sel.X); ok {
					wi := w.evalValue(e.Args[0])
					wi.wrappers = append(append([]mwEntry(nil), ch...), wi.wrappers...)
					return wi
				}
			}
		}
		if exprName(e.Fun) == "http.StripPrefix" && len(e.Args) == 2 {
			wi := w.evalValue(e.Args[1])
			if s, ok := w.constString(e.Args[0]); ok {
				wi.strip = s
			}
			return wi
		}
		if n := w.evalCall(e); n != nil {
			return wrapInfo{n: n}
		}
		return w.r.wraps[e]
	case *ast.Ident, *ast.SelectorExpr:
		// a variable holding a wrapped router serves it wrapped
		if k := keyOf(e); k != "" {
			if wi := w.env.handlers[k]; wi.n != nil && len(wi.wrappers) > 0 {
				return wi
			}
		}
		if n := w.lookup(e); n != nil {
			return wrapInfo{n: n}
		}
		if k := keyOf(e); k != "" {
			return w.env.handlers[k]
		}
	case *ast.FuncLit:
		return wrapInfo{}
	case *ast.CompositeLit, *ast.UnaryExpr:
		if wi := w.structMiddleware(e); wi.n != nil {
			return wi
		}
		if n := w.eval(x); n != nil {
			return wrapInfo{n: n}
		}
	default:
		if n := w.eval(x); n != nil {
			return wrapInfo{n: n}
		}
	}
	return wrapInfo{}
}

// structMiddleware: authMW{next: mux} or &RequireAuthMW{inner: b}, a type
// whose ServeHTTP is the middleware, holding the router it guards.
func (w *walker) structMiddleware(x ast.Expr) wrapInfo {
	cl := compositeOf(x)
	if cl == nil || cl.Type == nil {
		return wrapInfo{}
	}
	tr, ok := w.typeOfValue(x)
	if !ok {
		return wrapInfo{}
	}
	fi := w.methodOf(tr, "ServeHTTP", 0)
	if fi == nil {
		return wrapInfo{}
	}
	var inner *node
	for _, f := range w.literalFields(cl) {
		if n := w.eval(f.value); n != nil {
			inner = n
		}
	}
	if inner == nil {
		return wrapInfo{}
	}
	saved := w.r.asMW
	w.r.asMW = true
	v := w.funcVerdict(fi)
	w.r.asMW = saved
	e := mwEntry{name: exprName(cl.Type), enforces: v == vEnforces, passes: v == vPasses}
	if v == vConditional {
		e.name += " [rejects only under a condition]"
		e.cond = true
	}
	if v == vComplex {
		e.name += " [too many branches to verify]"
		e.unverified = true
	}
	return wrapInfo{n: inner, wrappers: []mwEntry{e}}
}

// wrapOf reads what evalValue already computed, without re-evaluating.
func (w *walker) wrapOf(x ast.Expr) wrapInfo {
	switch e := unparen(x).(type) {
	case *ast.CallExpr:
		return w.r.wraps[e]
	case *ast.Ident, *ast.SelectorExpr:
		if k := keyOf(e); k != "" {
			return w.env.handlers[k]
		}
	}
	return wrapInfo{}
}

// sink marks an expression as the handler that serves requests: whatever
// wraps a router there wraps every route on it.
func (w *walker) sink(x ast.Expr) wrapInfo {
	wi := w.evalValue(x)
	if wi.n != nil && len(wi.wrappers) > 0 {
		n := wi.n.canon()
		n.outer = appendUnique(n.outer, wi.wrappers...)
	}
	return wi
}

// handlerSliceType: gin.HandlersChain and other named handler slices.
func handlerSliceType(t ast.Expr) bool {
	n := typeBaseName(t)
	return strings.HasSuffix(n, "HandlersChain") || strings.HasSuffix(n, "Chain") || strings.HasSuffix(n, "Middlewares")
}

func hasHandlerField(cl *ast.CompositeLit) bool {
	for _, e := range cl.Elts {
		if kv, ok := e.(*ast.KeyValueExpr); ok {
			if id, ok := kv.Key.(*ast.Ident); ok && id.Name == "Handler" {
				return true
			}
		}
	}
	return false
}

// isNegroni reports whether ip is negroni, under urfave or its old codegangsta path.
func isNegroni(ip string) bool {
	return ip == "github.com/urfave/negroni" || ip == "github.com/codegangsta/negroni" || strings.HasPrefix(ip, "github.com/urfave/negroni/")
}

// negroniCall follows a negroni stack held in a variable: Use/UseFunc add
// middleware, UseHandler/UseHandlerFunc set what it serves.
func (w *walker) negroniCall(fun *ast.SelectorExpr, call *ast.CallExpr) bool {
	k := keyOf(fun.X)
	ch, isChain := w.env.chains[k]
	if k == "" || !isChain || !w.env.negroni[k] {
		return false
	}
	switch fun.Sel.Name {
	case "Use", "UseFunc":
		w.env.chains[k] = append(append([]mwEntry(nil), ch...), w.negroniEntries(call.Args)...)
	case "UseHandler", "UseHandlerFunc":
		if len(call.Args) == 1 {
			wi := w.evalValue(call.Args[0])
			if wi.n != nil {
				w.env.handlers[k] = wrapInfo{n: wi.n, wrappers: append(append([]mwEntry(nil), ch...), wi.wrappers...), strip: wi.strip}
			}
		}
	default:
		return false
	}
	return true
}

// negroniEntries: negroni.HandlerFunc(fn) is fn; negroni.Wrap(h) is a
// handler, not middleware.
func (w *walker) negroniEntries(args []ast.Expr) []mwEntry {
	var out []mwEntry
	for _, a := range args {
		x := a
		if c, ok := unparen(a).(*ast.CallExpr); ok && len(c.Args) == 1 {
			switch exprName(c.Fun) {
			case "negroni.HandlerFunc":
				x = c.Args[0]
			case "negroni.Wrap", "negroni.WrapFunc":
				continue
			}
		}
		if e := w.entryOf(x); e.name != "" {
			out = append(out, e)
		}
	}
	return out
}

// callClosure walks a local closure where it's called, with its string
// and router arguments bound, so reg("/admin/keys") registers that path.
func (w *walker) callClosure(fl *ast.FuncLit, call *ast.CallExpr) *node {
	// a closure that calls itself (walk := func(...) { walk(...) }) is
	// followed once
	if w.r.inClosure[fl] || !w.tick() {
		return nil
	}
	w.r.inClosure[fl] = true
	defer delete(w.r.inClosure, fl)
	bind := map[int]*node{}
	for i, a := range call.Args {
		if n := w.eval(a); n != nil {
			bind[i] = n
		}
	}
	sub := w.child(w.fc, w.shadow())
	sub.bindParams(fl.Type, bind, "local closure")
	i := 0
	for _, f := range fl.Type.Params.List {
		for _, nm := range f.Names {
			if i < len(call.Args) {
				if s, ok := w.constString(call.Args[i]); ok {
					sub.env.strs[nm.Name] = s
				}
			}
			i++
		}
	}
	sub.walkStmts(fl.Body.List)
	if len(sub.returned) > 0 {
		return sub.returned[0]
	}
	return nil
}

// listenOpen records a framework's own listener: no wrappers around it.
func (w *walker) listenOpen(n *node) {
	n = n.canon()
	for _, l := range n.listeners {
		if len(l) == 0 {
			return
		}
	}
	n.listeners = append(n.listeners, nil)
}

// serve marks x as what a server answers requests with. Each server is a
// separate way in: the same mux served open on one port and behind auth on
// another is open.
func (w *walker) serve(x ast.Expr) wrapInfo {
	wi := w.evalValue(x)
	if wi.n == nil {
		w.serveAPI(x)
		return wi
	}
	n := wi.n.canon()
	for _, l := range n.listeners {
		if sameNames(l, wi.wrappers) {
			return wi
		}
	}
	n.listeners = append(n.listeners, append([]mwEntry(nil), wi.wrappers...))
	return wi
}

func sameNames(a, b []mwEntry) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].name != b[i].name {
			return false
		}
	}
	return true
}

func appendUnique(xs []mwEntry, add ...mwEntry) []mwEntry {
	for _, a := range add {
		dup := false
		for _, x := range xs {
			if x.name == a.name {
				dup = true
				break
			}
		}
		if !dup {
			xs = append(xs, a)
		}
	}
	return xs
}

// resolveCall finds the declaration a call invokes, if it's in the tree.
func (w *walker) resolveCall(call *ast.CallExpr) *funcInfo {
	if info := w.fc.pkg.info; info != nil {
		f := typeutil.StaticCallee(info, call)
		if f == nil {
			// a method value passed as a callback: r.Route("/x", api.routes)
			var id *ast.Ident
			switch fun := unparen(call.Fun).(type) {
			case *ast.Ident:
				id = fun
			case *ast.SelectorExpr:
				id = fun.Sel
			}
			if id != nil {
				f, _ = info.Uses[id].(*types.Func)
			}
		}
		if f != nil {
			if fi := w.ix.byObj[f]; fi != nil {
				return fi
			}
			if f.Pkg() != nil {
				if p := w.ix.byImp[f.Pkg().Path()]; p != nil {
					if fi := p.funcs[funcKey(f)]; fi != nil {
						return fi
					}
				}
			}
		}
	}
	fn := unparen(call.Fun)
	switch g := fn.(type) {
	case *ast.IndexExpr:
		fn = g.X
	case *ast.IndexListExpr:
		fn = g.X
	}
	switch fun := fn.(type) {
	case *ast.Ident:
		if w.env.locals[fun.Name] {
			return nil
		}
		return w.fc.pkg.funcs[fun.Name]
	case *ast.SelectorExpr:
		if id, ok := fun.X.(*ast.Ident); ok && !w.env.locals[id.Name] {
			if _, ok := w.fc.imports[id.Name]; ok {
				return w.ix.lookupFunc(w.fc, id.Name, fun.Sel.Name)
			}
		}
		if tr, ok := w.typeOf(fun.X); ok {
			if fi := w.methodOf(tr, fun.Sel.Name, 0); fi != nil {
				return fi
			}
		}
		// adminMod{}.Routes(r), NewAdmin().Routes(g)
		switch unparen(fun.X).(type) {
		case *ast.CompositeLit, *ast.UnaryExpr, *ast.CallExpr:
			if tr, ok := w.typeOfValue(fun.X); ok {
				if fi := w.methodOf(tr, fun.Sel.Name, 0); fi != nil {
					return fi
				}
			}
		}
		// Unknown receiver type: accept the method only if its name is unique
		// in the tree and it touches routers, and the receiver isn't a value
		// of an outside type (gin's c.Set is not our apiClient.Set).
		if id := rootIdent(fun.X); id != nil && w.env.extern[id.Name] {
			return nil
		}
		if ms := w.ix.methodsByName[fun.Sel.Name]; len(ms) == 1 && ms[0].touches {
			return ms[0]
		}
	}
	return nil
}

func funcKey(f *types.Func) string {
	sig, ok := f.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return f.Name()
	}
	t := sig.Recv().Type()
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	if n, ok := t.(*types.Named); ok {
		return n.Obj().Name() + "." + f.Name()
	}
	return f.Name()
}

func (w *walker) methodOf(tr typeRef, name string, depth int) *funcInfo {
	if tr.pkg == nil || depth > 3 {
		return nil
	}
	if fi := tr.pkg.funcs[tr.name+"."+name]; fi != nil {
		return fi
	}
	if si := tr.pkg.structs[tr.name]; si != nil {
		for _, et := range si.embedded {
			if p, tn, es := w.ix.structOf(si.fc, et); es != nil {
				if fi := w.methodOf(typeRef{pkg: p, name: tn}, name, depth+1); fi != nil {
					return fi
				}
			}
		}
	}
	return nil
}

// inline walks a called function with its parameters bound to the caller's
// routers, strings and middleware, and returns the router it returns.
func (w *walker) inline(fi *funcInfo, call *ast.CallExpr, recv ast.Expr) *node {
	if len(fi.variants) == 0 {
		return w.inlineOne(fi, call, recv)
	}
	// build-tagged variants: walk each, so the weakest registration shows
	var res *node
	for _, v := range append([]*funcInfo{fi}, fi.variants...) {
		if n := w.inlineOne(v, call, recv); n != nil && res == nil {
			res = n
		}
	}
	return res
}

func (w *walker) inlineOne(fi *funcInfo, call *ast.CallExpr, recv ast.Expr) *node {
	if !fi.touches || fi.decl.Body == nil {
		return w.unknownCall(call)
	}
	depth := 0
	for _, s := range w.stack {
		if s == fi {
			depth++
		}
	}
	// A function calling itself is followed while an integer argument folds
	// to a constant (register(r, depth-1) terminates on `if depth > 0`);
	// otherwise, or past 8 levels, the recursion is cut and the scan
	// reported incomplete.
	if depth > 0 && (depth >= 8 || !w.intArg(call)) {
		w.cut(fi, call, "recursion through")
		return nil
	}
	if len(w.stack) >= maxInlineDepth || w.r.budgetHit {
		w.cut(fi, call, "the call chain through")
		return w.unknownCall(call)
	}
	e := newEnv()
	bind := map[int]*node{}
	i := 0
	for _, fld := range fi.decl.Type.Params.List {
		names := fld.Names
		if len(names) == 0 {
			names = []*ast.Ident{nil}
		}
		for _, nm := range names {
			if i >= len(call.Args) {
				break
			}
			arg := call.Args[i]
			n := w.eval(arg)
			if nm != nil {
				switch {
				case n != nil:
					bind[i] = n
				default:
					if s, ok := w.constString(arg); ok {
						e.strs[nm.Name] = s
					}
					if b, ok := w.constBool(arg); ok {
						e.bools[nm.Name] = b
					}
					if v, ok := w.constInt(arg); ok {
						e.ints[nm.Name] = v
					}
					if ch, ok := w.chainExpr(arg); ok {
						e.chains[nm.Name] = ch
					}
					if name := w.mwName(arg); name != "" {
						e.values[nm.Name] = name
						e.enforce[nm.Name] = w.mwVerdict(arg)
					}
					if wi := w.wrapOf(arg); wi.n != nil {
						e.handlers[nm.Name] = wi
					}
					// HandlerWithOptions(si, ChiServerOptions{BaseURL: "/v1",
					// BaseRouter: r, Middlewares: ...}): options.X in the callee
					if cl := compositeOf(arg); cl != nil {
						w.bindFields(e, nm.Name, cl)
					}
					if _, variadic := fld.Type.(*ast.Ellipsis); variadic {
						if call.Ellipsis.IsValid() {
							e.slices[nm.Name] = w.env.slices[keyOf(arg)]
						} else {
							e.slices[nm.Name] = append([]ast.Expr(nil), call.Args[i:]...)
						}
					}
				}
			}
			i++
		}
	}
	sub := &walker{r: w.r, ix: w.ix, fc: fi.fc, env: e, origin: w.origin,
		stack: append(append([]*funcInfo(nil), w.stack...), fi), consumed: w.consumed}
	sub.bindParams(fi.decl.Type, bind, fi.name)
	recvName := ""
	if fi.decl.Recv != nil && len(fi.decl.Recv.List) > 0 && len(fi.decl.Recv.List[0].Names) > 0 {
		recvName = fi.decl.Recv.List[0].Names[0].Name
		e.locals[recvName] = true
		e.vtypes[recvName] = typeRef{pkg: fi.fc.pkg, name: fi.recvType}
		if rk := keyOf(recv); rk != "" {
			for k, n := range w.env.routers {
				if strings.HasPrefix(k, rk+".") {
					e.routers[recvName+k[len(rk):]] = n
				}
			}
			if n := w.env.routers[rk]; n != nil {
				e.routers[recvName] = n
			}
		}
		if cl := compositeOf(recv); cl != nil {
			for _, f := range w.literalFields(cl) {
				if n := w.eval(f.value); n != nil {
					e.routers[recvName+"."+f.name] = n
				}
			}
		}
		// api.Group("/users").DELETE(...): the receiver is a call that
		// returned a wrapper struct; its router fields come with it
		if c, ok := unparen(recv).(*ast.CallExpr); ok {
			w.r.lastFields = nil
			if n := w.eval(c); n == nil {
				for f, fn := range w.r.lastFields {
					e.routers[recvName+"."+f] = fn
				}
			}
			w.r.lastFields = nil
		}
	}
	if w.origin != nil {
		m := w.r.inlined[w.origin]
		if m == nil {
			m = map[*funcInfo]bool{}
			w.r.inlined[w.origin] = m
		}
		m[fi] = true
	}
	before := len(w.r.recs)
	sub.walkStmts(fi.decl.Body.List)
	// the call counts as traced only if it produced a route with a known
	// path: gitea's m.Get(pattern, h) into its own Router collapses into one
	// dynamic route, and its call sites stay unchecked
	if call != nil && call.Pos().IsValid() {
		for _, rec := range w.r.recs[before:] {
			if !rec.dynamic {
				w.r.viaSites = append(w.r.viaSites, model.Site{File: w.fc.rel, Line: w.line(call.Pos())})
				break
			}
		}
	}
	if len(sub.returned) == 0 {
		if n := w.declaredRouter(fi, call); n != nil {
			sub.returned = append(sub.returned, n)
		}
	}
	if recvName != "" {
		if rk := keyOf(recv); rk != "" {
			for k, n := range e.routers {
				if strings.HasPrefix(k, recvName+".") {
					w.env.routers[rk+k[len(recvName):]] = n
				}
			}
		}
	}
	w.r.lastFields = sub.retFields
	if len(sub.returned) > 0 {
		return sub.returned[0]
	}
	return nil
}

// inlineOnRouter follows a method whose receiver evaluated to a router: the
// receiver and its embedded router field both name that router inside.
func (w *walker) inlineOnRouter(fi *funcInfo, call *ast.CallExpr, recv *node) *node {
	if fi.decl.Recv == nil || len(fi.decl.Recv.List) == 0 || len(fi.decl.Recv.List[0].Names) == 0 {
		return w.inline(fi, call, nil)
	}
	rn := fi.decl.Recv.List[0].Names[0].Name
	key := "\x00recv"
	w.env.routers[key] = recv
	if si := fi.fc.pkg.structs[fi.recvType]; si != nil {
		for name, t := range si.embedded {
			if si.fc.routerType(t) != nil {
				w.env.routers[key+"."+name] = recv
			}
		}
	}
	fake := &ast.CallExpr{Fun: &ast.SelectorExpr{X: ast.NewIdent(key), Sel: ast.NewIdent(fi.name)}, Args: call.Args, Lparen: call.Lparen}
	_ = rn
	n := w.inline(fi, fake, ast.NewIdent(key))
	for k := range w.env.routers {
		if k == key || strings.HasPrefix(k, key+".") {
			delete(w.env.routers, k)
		}
	}
	// r.SubGroup("/api/ipc", ipcAuth) returning a new group: the wrapper may
	// pass its middleware on through closures the walk can't follow, so the
	// middleware-looking arguments go on the group it returned
	if n != nil && n.canon() != recv.canon() {
		var mw []ast.Expr
		for _, a := range call.Args {
			if _, ok := w.constString(a); ok {
				continue
			}
			if _, ok := unparen(a).(*ast.BasicLit); ok {
				continue
			}
			mw = append(mw, a)
		}
		have := map[string]bool{}
		for _, m := range n.mw {
			have[m.name] = true
		}
		for _, e := range w.mwNames(mw) {
			if !have[e.name] {
				n.mw = append(n.mw, e)
			}
		}
	}
	return n
}

// constBool folds true, false, a parameter or table field bound to one,
// !x, && / ||, and comparisons of folded integers (depth > 0).
func (w *walker) constBool(x ast.Expr) (bool, bool) {
	switch e := unparen(x).(type) {
	case *ast.Ident:
		switch e.Name {
		case "true":
			return true, true
		case "false":
			return false, true
		}
		b, ok := w.env.bools[e.Name]
		return b, ok
	case *ast.SelectorExpr:
		if k := keyOf(e); k != "" {
			b, ok := w.env.bools[k]
			return b, ok
		}
	case *ast.UnaryExpr:
		if e.Op == token.NOT {
			if b, ok := w.constBool(e.X); ok {
				return !b, true
			}
		}
	case *ast.BinaryExpr:
		switch e.Op {
		case token.LAND, token.LOR:
			a, ok1 := w.constBool(e.X)
			b, ok2 := w.constBool(e.Y)
			if ok1 && ok2 {
				if e.Op == token.LAND {
					return a && b, true
				}
				return a || b, true
			}
		case token.GTR, token.GEQ, token.LSS, token.LEQ, token.EQL, token.NEQ:
			a, ok1 := w.constInt(e.X)
			b, ok2 := w.constInt(e.Y)
			if ok1 && ok2 {
				switch e.Op {
				case token.GTR:
					return a > b, true
				case token.GEQ:
					return a >= b, true
				case token.LSS:
					return a < b, true
				case token.LEQ:
					return a <= b, true
				case token.EQL:
					return a == b, true
				}
				return a != b, true
			}
		}
	}
	return false, false
}

// constInt folds integer literals, bound parameters and + / - of those.
func (w *walker) constInt(x ast.Expr) (int64, bool) {
	switch e := unparen(x).(type) {
	case *ast.BasicLit:
		if e.Kind == token.INT {
			v, err := strconv.ParseInt(e.Value, 0, 64)
			return v, err == nil
		}
	case *ast.Ident:
		v, ok := w.env.ints[e.Name]
		return v, ok
	case *ast.BinaryExpr:
		a, ok1 := w.constInt(e.X)
		b, ok2 := w.constInt(e.Y)
		if ok1 && ok2 {
			switch e.Op {
			case token.ADD:
				return a + b, true
			case token.SUB:
				return a - b, true
			}
		}
	}
	return 0, false
}

func (w *walker) intArg(call *ast.CallExpr) bool {
	for _, a := range call.Args {
		if _, ok := w.constInt(a); ok {
			return true
		}
	}
	return false
}

// cut records a call the walk refused to follow. It only matters, and is
// only reported, when a router is handed down: then routes registered
// deeper are lost. A recursive Serve or Shutdown that carries no router
// loses nothing.
func (w *walker) cut(fi *funcInfo, call *ast.CallExpr, what string) {
	carries := false
	for _, a := range call.Args {
		if w.peek(a) != nil {
			carries = true
		}
	}
	if !carries {
		return
	}
	s := w.site(call.Pos())
	if w.r.cutAt[s] {
		return
	}
	w.r.cutAt[s] = true
	w.r.cuts = append(w.r.cuts, model.Site{File: w.fc.rel, Line: w.line(call.Pos()), Call: what + " " + fi.name})
}

// declaredRouter: a function declared to return a router whose body didn't
// yield one we could follow (built by a call we don't model). Its routes
// still count, on a router marked partial.
func (w *walker) declaredRouter(fi *funcInfo, call *ast.CallExpr) *node {
	res := fi.decl.Type.Results
	if res == nil || len(res.List) == 0 {
		return nil
	}
	fw := fi.fc.routerTypeOrTyped(res.List[0].Type, w)
	if fw == nil {
		return nil
	}
	n := w.r.newNode(fw, kindRoot, w.site(call.Pos()))
	n.partial = "router built inside " + fi.name + " in a way that wasn't traced"
	return n
}

// untracedReceiver: a registration-shaped call (verb, "/path" literal) on a
// value we couldn't trace to a router. The route is recorded on a partial
// router rather than dropped, and the site is listed.
func (w *walker) untracedReceiver(fun *ast.SelectorExpr, call *ast.CallExpr) *node {
	key := keyOf(fun.X)
	if key == "" || len(call.Args) < 2 {
		return nil
	}
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}
	p, err := strconv.Unquote(lit.Value)
	if err != nil || !strings.HasPrefix(p, "/") {
		return nil
	}
	var fw *framework
	for _, f := range w.fc.fws {
		if f == fwNetHTTP {
			continue
		}
		if _, ok := f.verbs[fun.Sel.Name]; ok {
			fw = f
			break
		}
	}
	if fw == nil {
		return nil
	}
	s := w.site(call.Pos())
	if !w.r.untracedAt[s] {
		w.r.untracedAt[s] = true
		w.r.untraced = append(w.r.untraced, model.Site{File: w.fc.rel, Line: w.line(call.Pos()), Call: exprText(call.Fun)})
	}
	n := w.r.newNode(fw, kindParam, s)
	n.partial = "couldn't trace where router " + key + " comes from"
	w.env.routers[key] = n
	return n
}

// ---- middleware names ----

// mwName renders a middleware expression, appending what a variable was
// assigned from: authMW=jwtauth.Verifier.
func (w *walker) mwName(x ast.Expr) string {
	if _, ok := unparen(x).(*ast.FuncLit); ok {
		return ""
	}
	name := exprName(x)
	if name == "" {
		return ""
	}
	id := rootIdent(calleeOf(x))
	if id == nil {
		return name
	}
	origin := w.env.values[id.Name]
	if origin == "" && !w.env.locals[id.Name] {
		if g := w.fc.pkg.globals[id.Name]; g != nil && g.init != nil {
			origin = exprName(g.init)
		}
	}
	switch origin {
	case "", name, "make", "append", "new":
		return name
	}
	if i := strings.IndexByte(origin, '='); i >= 0 {
		origin = origin[i+1:]
	}
	return name + "=" + origin + strings.TrimPrefix(name, id.Name)
}

func calleeOf(x ast.Expr) ast.Expr {
	if c, ok := unparen(x).(*ast.CallExpr); ok {
		return c.Fun
	}
	return x
}

func (w *walker) callName(call *ast.CallExpr) string { return w.mwName(call) }

// mwNames renders middleware arguments; alice/chi chain values expand.
func (w *walker) mwNames(args []ast.Expr) []mwEntry {
	var out []mwEntry
	for _, a := range args {
		if ch, ok := w.chainExpr(a); ok {
			out = append(out, ch...)
			continue
		}
		if call, ok := unparen(a).(*ast.CallExpr); ok && call.Ellipsis.IsValid() {
			continue
		}
		if fl, ok := unparen(a).(*ast.FuncLit); ok {
			// an inline middleware literal only matters if it rejects requests
			saved := w.r.asMW
			w.r.asMW = true
			v := w.bodyVerdict(fl.Body, w.fc, 0, nil, fl.Type)
			w.r.asMW = saved
			if v == vEnforces {
				out = append(out, mwEntry{name: "func@" + w.site(fl.Pos()), enforces: true})
			}
			continue
		}
		if e := w.entryOf(a); e.name != "" {
			out = append(out, e)
		}
	}
	return out
}

func (w *walker) entryOf(x ast.Expr) mwEntry {
	w.r.lastWhy = ""
	v := w.mwVerdict(x)
	e := mwEntry{name: w.mwName(x), enforces: v == vEnforces, passes: v == vPasses}
	if v == vComplex {
		// too many branches to read: unverified, never trusted by name
		e.name += " [too many branches to verify]"
		e.unverified = true
	}
	if v == vConditional {
		why := w.r.lastWhy
		if why == "" {
			why = "rejects only under a condition"
		}
		e.name += " [" + why + "]"
		e.cond = true
	}
	// echojwt.WithConfig(echojwt.Config{Skipper: ...}): auth that some
	// requests skip is conditional auth
	if skipper(x) && authsig.ClassifyName(e.name) == authsig.Auth {
		e.name += "(Skipper)"
		e.cond = true
	}
	// handler.Web(fn, UserLevelNone): a level argument that says "nobody"
	// overrides whatever the wrapper can enforce
	if call, ok := unparen(x).(*ast.CallExpr); ok && !strings.HasPrefix(exprName(call.Fun), "authboss.Middleware") {
		for _, a := range call.Args {
			if id := rootIdent(a); id != nil && w.env.locals[id.Name] {
				continue // r.optional: a variable, value unknown; the body decides
			}
			if n := exprName(a); n != "" && optionalArg(n) {
				e.name += "(" + n + ")"
				e.enforces, e.passes = false, false
				break
			}
		}
	}
	return e
}

func skipper(x ast.Expr) bool {
	call, ok := unparen(x).(*ast.CallExpr)
	if !ok {
		return false
	}
	for _, a := range call.Args {
		if cl := compositeOf(a); cl != nil {
			for _, elt := range cl.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if id, ok := kv.Key.(*ast.Ident); ok && id.Name == "Skipper" && !isNil(kv.Value) {
						return true
					}
				}
			}
		}
	}
	return false
}

func optionalArg(name string) bool {
	low := strings.ToLower(name)
	if i := strings.LastIndexByte(low, '.'); i >= 0 {
		low = low[i+1:]
	}
	for _, w := range []string{"none", "anonymous", "guest", "public", "optional", "noauth"} {
		if strings.Contains(low, w) {
			return true
		}
	}
	return false
}

func names(es []mwEntry) []string {
	out := make([]string, 0, len(es))
	for _, e := range es {
		out = append(out, e.name)
	}
	return out
}

// chainExpr recognises middleware chain values: alice.New(a, b),
// chi.Chain(a, b), c.Append(x), and variables holding them.
func (w *walker) chainExpr(x ast.Expr) ([]mwEntry, bool) {
	switch e := unparen(x).(type) {
	case *ast.Ident:
		ch, ok := w.env.chains[e.Name]
		return ch, ok
	case *ast.SelectorExpr:
		// options.Middlewares bound from a struct literal
		if k := keyOf(e); k != "" {
			ch, ok := w.env.chains[k]
			return ch, ok
		}
	case *ast.CompositeLit:
		// []MiddlewareFunc{RequireAuth}
		if _, ok := e.Type.(*ast.ArrayType); ok {
			return w.mwNames(e.Elts), true
		}
	case *ast.CallExpr:
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil, false
		}
		if id, ok := sel.X.(*ast.Ident); ok {
			ip := w.fc.imports[id.Name]
			if (ip == "github.com/justinas/alice" && sel.Sel.Name == "New") ||
				(w.fc.fws[id.Name] == fwChi && sel.Sel.Name == "Chain") {
				return w.mwNames(e.Args), true
			}
			if isNegroni(ip) && (sel.Sel.Name == "New" || sel.Sel.Name == "Classic") {
				return w.negroniEntries(e.Args), true
			}
		}
		if sel.Sel.Name == "Append" || sel.Sel.Name == "Extend" {
			if base, ok := w.chainExpr(sel.X); ok {
				var add []mwEntry
				for _, a := range e.Args {
					if ch, ok := w.chainExpr(a); ok {
						add = append(add, ch...)
					} else if en := w.entryOf(a); en.name != "" {
						add = append(add, en)
					}
				}
				return append(append([]mwEntry(nil), base...), add...), true
			}
		}
	}
	return nil, false
}

// unwrapHandler peels middleware calls around a route's handler:
// RequireAuth(h) -> [RequireAuth], h. A call with no arguments is a handler
// factory (s.handleX()), not a wrapper.
func (w *walker) unwrapHandler(x ast.Expr) ([]mwEntry, ast.Expr) {
	var ws []mwEntry
	if x == nil {
		return nil, nil
	}
	for i := 0; i < 12; i++ {
		call, ok := unparen(x).(*ast.CallExpr)
		if !ok || len(call.Args) == 0 {
			return ws, x
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok && len(call.Args) == 1 {
			switch sel.Sel.Name {
			case "Then", "ThenFunc", "Handler", "HandlerFunc":
				if ch, ok := w.chainExpr(sel.X); ok {
					ws = append(ws, ch...)
					x = call.Args[0]
					continue
				}
			}
		}
		// chain.ToHandlerFunc(h, mw...), Wrap(h, auth, logger): a builder
		// whose other arguments are all middleware applies them to h
		if h, mws, ok := w.builderMiddleware(call); ok {
			ws = append(ws, w.entryOf(call))
			ws = append(ws, w.mwNames(mws)...)
			x = h
			continue
		}
		name := exprName(call.Fun)
		inner := w.innerHandlerArg(call.Args)
		if inner == nil {
			return ws, x
		}
		if !handlerAdapters[name] {
			// b.TeamLeaderAccess(h) that composes h = mw(h); return h
			if ents, ok := w.wrapperEntries(call); ok {
				ws = append(ws, ents...)
			} else {
				ws = append(ws, w.entryOf(call))
			}
		}
		x = inner
	}
	return ws, x
}

// builderMiddleware splits a builder call into the handler and the
// middleware applied to it. Every argument but one must be a middleware
// value; a spread slice counts only if every element is known and is one.
func (w *walker) builderMiddleware(call *ast.CallExpr) (h ast.Expr, mws []ast.Expr, ok bool) {
	if len(call.Args) < 2 {
		return nil, nil, false
	}
	for i, a := range call.Args {
		if call.Ellipsis.IsValid() && i == len(call.Args)-1 {
			elems, found := w.env.slices[keyOf(a)]
			if !found {
				return nil, nil, false
			}
			for _, e := range elems {
				if !w.middlewareValue(e) {
					return nil, nil, false
				}
			}
			mws = append(mws, elems...)
			continue
		}
		if w.middlewareValue(a) {
			mws = append(mws, a)
			continue
		}
		if h != nil {
			return nil, nil, false
		}
		h = a
	}
	return h, mws, h != nil && len(mws) > 0
}

// middlewareValue: a function or method value shaped like middleware,
// func(next X) X where X is a handler type.
func (w *walker) middlewareValue(x ast.Expr) bool {
	var ft *ast.FuncType
	if fl, ok := unparen(x).(*ast.FuncLit); ok {
		ft = fl.Type
	} else if fi := w.funcValue(x); fi != nil {
		ft = fi.decl.Type
	}
	if ft == nil || ft.Params == nil || ft.Results == nil || len(ft.Params.List) != 1 || len(ft.Results.List) != 1 {
		return false
	}
	in, out := ft.Params.List[0], ft.Results.List[0]
	if len(in.Names) > 1 || len(out.Names) > 1 {
		return false
	}
	return handlerType(in.Type) && handlerType(out.Type)
}

func handlerType(t ast.Expr) bool {
	if f, ok := t.(*ast.FuncType); ok {
		return takesRequest(f)
	}
	n := typeBaseName(t)
	return strings.HasSuffix(n, "Handler") || strings.HasSuffix(n, "HandlerFunc")
}

// innerHandlerArg picks the wrapped handler among a wrapper's arguments:
// the only one, else the most handler-like (a func literal or a function
// taking a ResponseWriter/Context beats a name like h or viewX, which beats
// any other value); ties go to the last.
func (w *walker) innerHandlerArg(args []ast.Expr) ast.Expr {
	if len(args) == 1 {
		return args[0]
	}
	best, bestScore := ast.Expr(nil), 0
	for _, a := range args {
		if s := w.handlerScore(a); s > 0 && s >= bestScore {
			best, bestScore = a, s
		}
	}
	return best
}

func (w *walker) handlerScore(x ast.Expr) int {
	switch a := unparen(x).(type) {
	case *ast.BasicLit, *ast.BinaryExpr, *ast.CompositeLit:
		return 0
	case *ast.FuncLit:
		return 3
	case *ast.Ident:
		if _, ok := w.constString(a); ok || a.Name == "nil" || a.Name == "true" || a.Name == "false" {
			return 0
		}
	}
	if fi := w.funcValue(x); fi != nil && takesRequest(fi.decl.Type) {
		return 3
	}
	name := strings.ToLower(exprName(x))
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		name = name[i+1:]
	}
	switch {
	case name == "h" || name == "next" || name == "fn" || name == "f" || name == "hdl":
		return 2
	case strings.Contains(name, "handle"), strings.Contains(name, "serve"), strings.HasPrefix(name, "view"),
		strings.Contains(name, "page"), strings.HasPrefix(name, "api"):
		return 2
	}
	return 1
}

// funcValue resolves a function or method value (not a call) to its
// declaration.
func (w *walker) funcValue(x ast.Expr) *funcInfo {
	switch e := unparen(x).(type) {
	case *ast.Ident:
		if w.env.locals[e.Name] {
			return nil
		}
		return w.fc.pkg.funcs[e.Name]
	case *ast.SelectorExpr:
		fake := &ast.CallExpr{Fun: e}
		return w.resolveCall(fake)
	}
	return nil
}

// takesRequest: the signature of an HTTP handler in any of the frameworks.
func takesRequest(ft *ast.FuncType) bool {
	if ft.Params == nil {
		return false
	}
	for _, p := range ft.Params.List {
		switch typeBaseName(p.Type) {
		case "ResponseWriter", "Request", "Context", "Ctx":
			return true
		}
	}
	return false
}

// ---- results ----

// finish resolves every recorded route against its router's final state.
func (r *run) finish() []model.Route {
	r.emptyMounts()
	seen := map[string]int{}
	var out []model.Route
	// A registration reached from a real caller makes the copy walked from
	// an uncalled function's router parameter redundant (oapi-codegen's
	// unused HandlerFromMux beside the Handler the server calls).
	res := make([][]resolved, len(r.recs))
	concrete := map[string]bool{}
	for i, rec := range r.recs {
		r.anc = rec.anc
		res[i] = r.resolveAll(rec.n, rec.nmw, rec.path, 0, 0, map[*node]bool{})
		r.anc = nil
		for j := range res[i] {
			if rec.partial != "" {
				res[i][j].partial = rec.partial
			}
			if !paramRooted(res[i][j].partial) && !r.dead(rec.origin) {
				concrete[rec.method+" "+rec.file+":"+strconv.Itoa(rec.line)] = true
			}
		}
	}
	// a closure body walked once unbound (r.Delete(<p>, h)) and again
	// where it's called with a literal: keep the literal paths
	literalAt := map[string]bool{}
	for _, rec := range r.recs {
		if !rec.dynamic {
			literalAt[rec.file+":"+strconv.Itoa(rec.line)] = true
		}
	}
	type way struct {
		rec routeRec
		res resolved
	}
	var ways []way
	for i, rec := range r.recs {
		if rec.dynamic && literalAt[rec.file+":"+strconv.Itoa(rec.line)] && strings.HasPrefix(rec.path, "<") && !strings.ContainsAny(rec.path[1:len(rec.path)-1], " +(\"") {
			continue
		}
		for _, x := range res[i] {
			ways = append(ways, way{rec, x})
		}
	}
	for _, wy := range ways {
		rec := wy.rec
		prefix, entries, partial := wy.res.prefix, append([]mwEntry(nil), wy.res.entries...), wy.res.partial
		if (paramRooted(partial) || r.dead(rec.origin)) && concrete[rec.method+" "+rec.file+":"+strconv.Itoa(rec.line)] {
			continue
		}
		entries = append(entries, rec.extra...)
		chain := names(entries)
		var enforcing, conditional, passing, unverified []string
		for _, e := range entries {
			if e.unverified {
				unverified = append(unverified, e.name)
			}
			if e.enforces {
				enforcing = append(enforcing, e.name)
			}
			if e.passes {
				passing = append(passing, e.name)
			}
			if e.cond {
				conditional = append(conditional, e.name)
			}
		}
		scope, kind := r.scopeOf(rec.n)
		path := rec.path
		if !rec.absPath && (!rec.dynamic || prefix != "") {
			path = joinPath(prefix, rec.path)
		}
		if rec.n.fw == fwNetHTTP {
			// a ServeMux is flat: siblings are routes sharing a first segment
			scope += ":" + firstSegment(path)
			kind = model.ScopePrefix
		}
		if chain == nil {
			chain = []string{}
		}
		rt := model.Route{
			Method: rec.method, Path: path, Handler: rec.handler, Framework: rec.n.fw.name,
			File: rec.file, Line: rec.line, Chain: chain, Scope: scope, ScopeKind: kind,
			Partial: partial, Dynamic: rec.dynamic, Enforcing: enforcing, Conditional: conditional,
			HandlerChecks: rec.checks, Passing: passing, API: rec.api, Unverified: unverified,
		}
		k := rt.Method + " " + rt.Path + " " + rt.File + ":" + strconv.Itoa(rt.Line)
		if i, ok := seen[k]; ok {
			out[i] = mergeRoute(out[i], rt)
			continue
		}
		seen[k] = len(out)
		out = append(out, rt)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		if out[i].Method != out[j].Method {
			return out[i].Method < out[j].Method
		}
		return out[i].Path < out[j].Path
	})
	// Scope ids carry an internal node number that depends on walk order;
	// renumber them by first appearance so output is stable run to run.
	renum := map[string]string{}
	for i := range out {
		s := out[i].Scope
		if _, ok := renum[s]; !ok {
			site, suffix := s, ""
			if j := strings.IndexByte(s, '#'); j >= 0 {
				site = s[:j]
				if k := strings.IndexByte(s[j:], ':'); k >= 0 {
					suffix = s[j+k:] // net/http ":/segment"
				}
			}
			renum[s] = site + "#" + strconv.Itoa(len(renum)+1) + suffix
		}
		out[i].Scope = renum[s]
	}
	return out
}

// dead: the walk started in an unexported function nothing in the tree
// names (an unused newBareRouter). Its copies of routes a live caller also
// registers are dropped.
func (r *run) dead(fi *funcInfo) bool {
	if fi == nil || fi.recvType != "" || fi.name == "main" || fi.name == "init" || ast.IsExported(fi.name) {
		return false
	}
	if r.used == nil {
		r.used = map[string]bool{}
		for _, p := range r.ix.pkgs {
			for _, fc := range p.files {
				ast.Inspect(fc.file, func(n ast.Node) bool {
					switch x := n.(type) {
					case *ast.FuncDecl:
						if x.Body != nil {
							ast.Inspect(x.Body, func(m ast.Node) bool {
								if id, ok := m.(*ast.Ident); ok {
									r.used[id.Name] = true
								}
								return true
							})
						}
						return false
					case *ast.Ident:
						r.used[x.Name] = true
					}
					return true
				})
			}
		}
	}
	return !r.used[fi.name]
}

func paramRooted(partial string) bool {
	return strings.HasPrefix(partial, "router arrives as parameter")
}

// emptyMounts turns a mounted router with no routes we could read (an
// OpenAPI framework's mux, a router built elsewhere) into a partial
// catch-all route at the mount point, instead of letting it vanish.
func (r *run) emptyMounts() {
	used := map[*node]bool{}
	for _, rec := range r.recs {
		// the bound only guards against a parent cycle
		for n, i := rec.n, 0; n != nil && i < 4096; n, i = n.parent, i+1 {
			used[n.canon()] = true
		}
	}
	done := map[*node]bool{}
	for _, m := range r.mounts {
		if used[m.sub] || done[m.sub] {
			continue
		}
		done[m.sub] = true
		pp, _, _ := r.resolve(m.parent, -1, "", 0, 0)
		p := concatPath(pp, m.pattern)
		if !strings.HasSuffix(p, "*") {
			p = strings.TrimSuffix(p, "/") + "/*"
		}
		r.recs = append(r.recs, routeRec{
			n: m.sub, nmw: len(m.sub.mw), method: "*", path: p, absPath: true, handler: "(mounted router)",
			partial: "a router is mounted here but none of its routes could be read",
			file:    m.file, line: m.line, origin: m.origin,
		})
	}
}

// mergeRoute folds a second walk of the same registration into the first.
// Each keeps its own chain: a helper registered on one server behind AuthA
// and on another behind AuthB is authed both ways, and the analysis judges
// the route by its weakest variant.
func mergeRoute(a, b model.Route) model.Route {
	if len(a.Variants) == 0 {
		a.Variants = []model.Variant{variantOf(a)}
	}
	vb := variantOf(b)
	for _, v := range a.Variants {
		if strings.Join(v.Chain, "\x00") == strings.Join(vb.Chain, "\x00") && v.Partial == vb.Partial {
			return a
		}
	}
	a.Variants = append(a.Variants, vb)
	if a.API == "" {
		a.API = b.API
	}
	return a
}

func variantOf(r model.Route) model.Variant {
	return model.Variant{Chain: r.Chain, Enforcing: r.Enforcing, Passing: r.Passing, Conditional: r.Conditional, Unverified: r.Unverified,
		HandlerChecks: r.HandlerChecks, Partial: r.Partial}
}

// entryPoints returns the functions not reached from any other walked
// function: the places a clean walk should start.
func (r *run) entryPoints(fns []*funcInfo) map[*funcInfo]bool {
	inlinedBy := map[*funcInfo]bool{}
	for o, m := range r.inlined {
		for f := range m {
			if f != o {
				inlinedBy[f] = true
			}
		}
	}
	entry := map[*funcInfo]bool{}
	for _, f := range fns {
		if !inlinedBy[f] {
			entry[f] = true
		}
	}
	// Functions only reachable from each other (a cycle with no outside
	// caller) would otherwise never be walked.
	for {
		covered := map[*funcInfo]bool{}
		for k := range entry {
			for f := range r.inlined[k] {
				covered[f] = true
			}
		}
		for o, m := range r.inlined {
			if o.decl != nil && o.decl.Body == nil { // package-level initializers
				for f := range m {
					covered[f] = true
				}
			}
		}
		var orphan *funcInfo
		for _, f := range fns {
			if !entry[f] && !covered[f] {
				orphan = f
				break
			}
		}
		if orphan == nil {
			return entry
		}
		entry[orphan] = true
	}
}

// resolved is one way a route is reached: its path prefix, the middleware
// in front of it, and why that may be incomplete.
type resolved struct {
	prefix  string
	entries []mwEntry
	partial string
}

// maxWays caps the distinct ways up from one router. Each mount and listener
// multiplies them; past this many, more copies of a route add nothing.
const maxWays = 32

// resolveAll walks every way up from n: its structural parent, each mount
// edge and each listener, one result per distinct way.
func (r *run) resolveAll(n *node, k int, rel string, depth, exclude int, onStack map[*node]bool) []resolved {
	if depth > 256 || onStack[n] {
		return []resolved{{partial: "router nesting too deep to resolve"}}
	}
	onStack[n] = true
	defer delete(onStack, n)
	if v, ok := r.anc[n]; ok {
		k = v
	} else if requestTime(n) || k < 0 || k > len(n.mw) {
		k = len(n.mw)
	}
	mine := append([]mwEntry(nil), n.outer...)
	for _, m := range n.mw[:k] {
		if exclude != 0 && m.from == exclude {
			continue
		}
		if m.prefix == "" || relHasPrefix(rel, m.prefix) {
			mine = append(mine, m)
		}
	}
	var out []resolved
	add := func(x resolved) {
		for _, o := range out {
			if o.prefix == x.prefix && sameNames(o.entries, x.entries) && o.partial == x.partial {
				return
			}
		}
		if len(out) < maxWays {
			out = append(out, x)
		}
	}
	if n.parent != nil {
		for _, p := range r.resolveAll(n.parent, n.parentLen, concatPath(n.prefix, rel), depth+1, n.id, onStack) {
			part := n.partial
			if part == "" {
				part = p.partial
			}
			pre := concatPath(p.prefix, n.prefix)
			if n.absolute {
				pre = n.prefix
			}
			add(resolved{pre, append(append([]mwEntry(nil), p.entries...), mine...), part})
		}
	}
	for _, e := range n.edges {
		for _, p := range r.resolveAll(e.parent, -1, concatPath(concatPath(e.prefix, n.prefix), rel), depth+1, n.id, onStack) {
			part := n.partial
			if part == "" {
				part = p.partial
			}
			pre := concatPath(p.prefix, concatPath(e.prefix, n.prefix))
			if e.absolute {
				// the sub-router sees the full path (net/http, chi Handle)
				pre = concatPath(e.prefix, n.prefix)
			}
			chain := append(append(append([]mwEntry(nil), p.entries...), e.wrappers...), mine...)
			add(resolved{pre, chain, part})
		}
	}
	if n.parent == nil && (len(n.edges) == 0 || len(n.listeners) > 0) {
		if len(n.listeners) == 0 {
			add(resolved{n.prefix, mine, n.partial})
		}
		for _, l := range n.listeners {
			add(resolved{n.prefix, append(append([]mwEntry(nil), l...), mine...), n.partial})
		}
	}
	return out
}

func (r *run) resolve(n *node, k int, rel string, depth, exclude int) (prefix string, chain []mwEntry, partial string) {
	if depth > 4096 {
		return "", nil, "router nesting too deep to resolve"
	}
	if v, ok := r.anc[n]; ok {
		k = v
	} else if requestTime(n) || k < 0 || k > len(n.mw) {
		k = len(n.mw)
	}
	mine := append([]mwEntry(nil), n.outer...)
	for _, m := range n.mw[:k] {
		if exclude != 0 && m.from == exclude {
			continue
		}
		if m.prefix == "" || relHasPrefix(rel, m.prefix) {
			mine = append(mine, m)
		}
	}
	if n.parent == nil {
		return n.prefix, mine, n.partial
	}
	pp, pc, ppart := r.resolve(n.parent, n.parentLen, concatPath(n.prefix, rel), depth+1, n.id)
	partial = n.partial
	if partial == "" {
		partial = ppart
	}
	if n.absolute {
		return n.prefix, append(pc, mine...), partial
	}
	return concatPath(pp, n.prefix), append(pc, mine...), partial
}

func concatPath(a, b string) string {
	switch {
	case b == "":
		return a
	case a == "":
		return b
	}
	return strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(b, "/")
}

func (r *run) scopeOf(n *node) (string, string) {
	for i := 0; i < 4096 && n.kind == kindInline && n.parent != nil; i++ {
		n = n.parent
	}
	return n.site + "#" + strconv.Itoa(n.id), n.kind
}

func typeString(t types.Type) string {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	if n, ok := t.(*types.Named); ok && n.Obj() != nil {
		if n.Obj().Pkg() != nil {
			return n.Obj().Pkg().Path() + "." + n.Obj().Name()
		}
		return n.Obj().Name()
	}
	return t.String()
}

// routerFrameworkOf maps a (pointer to a) named type to its router framework.
func routerFrameworkOf(t types.Type) *framework {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return nil
	}
	fw := frameworkForImport(named.Obj().Pkg().Path())
	if fw != nil && fw.types[named.Obj().Name()] {
		return fw
	}
	return nil
}
