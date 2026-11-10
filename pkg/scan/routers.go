package scan

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
)

// routerMethod applies a method call on a known router and returns the
// router the call evaluates to (a group, a With() view, the same router for
// gin's chained Use), or nil.
func (w *walker) routerMethod(n *node, m string, call *ast.CallExpr) *node {
	if frameworkListen[m] && n.fw != fwNetHTTP && n.fw != fwChi && n.fw != fwGorilla {
		w.listenOpen(n)
		return nil
	}
	// g.Handle(method, path, chain...) with chain built here: list it out
	if call.Ellipsis.IsValid() && len(call.Args) > 0 {
		last := call.Args[len(call.Args)-1]
		elems, ok := w.env.slices[keyOf(last)]
		if a, isCall := unparen(last).(*ast.CallExpr); isCall && exprName(a.Fun) == "append" && len(a.Args) > 0 && !a.Ellipsis.IsValid() {
			if base, found := w.env.slices[keyOf(a.Args[0])]; found {
				elems, ok = append(append([]ast.Expr(nil), base...), a.Args[1:]...), true
			}
		}
		if ok {
			c := *call
			c.Args = append(append([]ast.Expr(nil), call.Args[:len(call.Args)-1]...), elems...)
			c.Ellipsis = token.NoPos
			call = &c
		}
	}
	// mux.Handle(adminv1connect.NewAdminServiceHandler(srv)): connect's
	// constructor returns (path, handler), spread into Handle
	if (m == "Handle" || m == "Mount") && len(call.Args) == 1 {
		if c, ok := unparen(call.Args[0]).(*ast.CallExpr); ok {
			if api := w.apiKind(c, 0); api != "" {
				w.registerPath(n, "*", "<"+exprText(c.Fun)+">", true, c, nil, call)
				w.r.recs[len(w.r.recs)-1].api = api
				return nil
			}
		}
	}
	switch n.fw {
	case fwChi:
		return w.chiMethod(n, m, call)
	case fwGin, fwHertz:
		return w.ginMethod(n, m, call)
	case fwEcho:
		return w.echoMethod(n, m, call)
	case fwFiber:
		return w.fiberMethod(n, m, call)
	case fwNetHTTP:
		if m == "Handle" || m == "HandleFunc" {
			w.netHTTPRegister(n, call)
		}
		return nil
	case fwGorilla:
		// direct Use/HandleFunc; builder chains go through gorillaChain
		switch m {
		case "Use":
			n.mw = append(n.mw, w.mwNames(call.Args)...)
		case "HandleFunc", "Handle":
			if len(call.Args) >= 2 {
				w.register(n, "*", call.Args[0], call.Args[len(call.Args)-1], nil, call)
			}
		}
		return nil
	}
	return nil
}

func (w *walker) derive(parent *node, kind string, prefix string, mw []mwEntry, pos ast.Node) *node {
	c := w.r.newNode(parent.fw, kind, w.site(pos.Pos()))
	c.parent, c.parentLen, c.prefix, c.mw = parent, len(parent.mw), prefix, mw
	return c
}

func (w *walker) chiMethod(n *node, m string, call *ast.CallExpr) *node {
	args := call.Args
	switch m {
	case "Use":
		n.mw = append(n.mw, w.mwNames(args)...)
	case "With":
		return w.derive(n, kindInline, "", w.mwNames(args), call)
	case "Group":
		c := w.derive(n, kindGroup, "", nil, call)
		if len(args) > 0 {
			w.enterRouterFunc(args[0], c)
		}
		return c
	case "Route":
		if len(args) < 2 {
			return nil
		}
		p, _ := w.pathArg(args[0])
		c := w.derive(n, kindGroup, p, nil, call)
		w.enterRouterFunc(args[1], c)
		return c
	case "Mount":
		if len(args) == 2 {
			w.mountOrRoute(n, "*", args[0], args[1], call, true)
		}
	case "Handle", "HandleFunc":
		if len(args) == 2 {
			w.mountOrRoute(n, "*", args[0], args[1], call, false)
		}
	case "Method", "MethodFunc":
		if len(args) == 3 {
			meth, _ := w.constString(args[0])
			w.register(n, strings.ToUpper(meth), args[1], args[2], nil, call)
		}
	default:
		if verb, ok := fwChi.verbs[m]; ok && len(args) >= 2 {
			w.register(n, verb, args[0], args[len(args)-1], nil, call)
		}
	}
	return nil
}

func (w *walker) ginMethod(n *node, m string, call *ast.CallExpr) *node {
	args := call.Args
	switch m {
	case "Use":
		n.mw = append(n.mw, w.mwNames(args)...)
		return n
	case "Group":
		if len(args) == 0 {
			return nil
		}
		p, _ := w.pathArg(args[0])
		return w.derive(n, kindGroup, p, w.mwNames(args[1:]), call)
	case "Handle":
		if len(args) >= 3 {
			meth, _ := w.constString(args[0])
			w.register(n, strings.ToUpper(meth), args[1], args[len(args)-1], args[2:len(args)-1], call)
		}
	case "Match":
		if len(args) >= 3 {
			for _, meth := range w.stringList(args[0]) {
				w.register(n, strings.ToUpper(meth), args[1], args[len(args)-1], args[2:len(args)-1], call)
			}
		}
	default:
		// gin accepts a route with no handlers of its own when the group's
		// middleware answers it: g.OPTIONS("/*any")
		if verb, ok := n.fw.verbs[m]; ok && len(args) >= 1 {
			w.register(n, verb, args[0], lastOrNil(args[1:]), initOf(args[1:]), call)
		}
	}
	return nil
}

// echoMethod handles echo's registration shape: verb(path, handler, mw...),
// with route middleware after the handler.
func (w *walker) echoMethod(n *node, m string, call *ast.CallExpr) *node {
	args := call.Args
	switch m {
	case "Use", "Pre":
		n.mw = append(n.mw, w.mwNames(args)...)
	case "Group":
		if len(args) == 0 {
			return nil
		}
		p, _ := w.pathArg(args[0])
		return w.derive(n, kindGroup, p, w.mwNames(args[1:]), call)
	case "Add":
		if len(args) >= 3 {
			meth, _ := w.constString(args[0])
			w.register(n, strings.ToUpper(meth), args[1], args[2], args[3:], call)
		}
	case "Match":
		if len(args) >= 3 {
			for _, meth := range w.stringList(args[0]) {
				w.register(n, strings.ToUpper(meth), args[1], args[2], args[3:], call)
			}
		}
	default:
		if verb, ok := fwEcho.verbs[m]; ok && len(args) >= 2 {
			w.register(n, verb, args[0], args[1], args[2:], call)
		}
	}
	return nil
}

func (w *walker) fiberMethod(n *node, m string, call *ast.CallExpr) *node {
	args := call.Args
	switch m {
	case "Use":
		// Use([prefix...], handlers...); a sub-app instead of a handler mounts it.
		var prefixes []string
		i := 0
		for ; i < len(args); i++ {
			if s, ok := w.constString(args[i]); ok {
				prefixes = append(prefixes, s)
				continue
			}
			if l := w.stringList(args[i]); l != nil {
				prefixes = append(prefixes, l...)
				continue
			}
			break
		}
		rest := args[i:]
		if len(rest) == 1 {
			if wi := w.evalValue(rest[0]); wi.n != nil {
				p := ""
				if len(prefixes) > 0 {
					p = prefixes[0]
				}
				w.mount(n, p, wi, call, false, p)
				return nil
			}
		}
		for _, e := range w.mwNames(rest) {
			if len(prefixes) == 0 {
				n.mw = append(n.mw, e)
			}
			for _, p := range prefixes {
				n.mw = append(n.mw, mwEntry{name: e.name, prefix: p})
			}
		}
	case "Group":
		if len(args) == 0 {
			return nil
		}
		p, _ := w.pathArg(args[0])
		return w.derive(n, kindGroup, p, w.mwNames(args[1:]), call)
	case "Route":
		if len(args) < 2 {
			return nil
		}
		p, _ := w.pathArg(args[0])
		c := w.derive(n, kindGroup, p, nil, call)
		w.enterRouterFunc(args[1], c)
		return c
	case "Mount":
		if len(args) == 2 {
			if wi := w.evalValue(args[1]); wi.n != nil {
				p, _ := w.constString(args[0])
				w.mount(n, p, wi, call, false, p)
			}
		}
	case "Add":
		if len(args) >= 3 {
			meths := w.stringList(args[0])
			if meths == nil {
				if s, ok := w.constString(args[0]); ok {
					meths = []string{s}
				}
			}
			for _, meth := range meths {
				w.register(n, strings.ToUpper(meth), args[1], args[len(args)-1], args[2:len(args)-1], call)
			}
		}
	default:
		if verb, ok := fwFiber.verbs[m]; ok && len(args) >= 2 {
			w.register(n, verb, args[0], args[len(args)-1], args[1:len(args)-1], call)
		}
	}
	return nil
}

// netHTTPRegister handles mux.Handle(pattern, h). A ServeMux or other router
// as the handler is a mount: with http.StripPrefix its routes sit under the
// stripped prefix, without it they already carry the full path.
func (w *walker) netHTTPRegister(n *node, call *ast.CallExpr) {
	if len(call.Args) < 2 {
		return
	}
	pattern, dyn := w.pathArg(call.Args[0])
	method, p := "*", pattern
	if !dyn {
		method, p = parsePattern(pattern)
	}
	if wi := w.evalValue(call.Args[1]); wi.n != nil && wi.n.canon() != n.canon() {
		w.mount(n, wi.strip, wi, call, true, p)
		return
	}
	w.registerPath(n, method, p, dyn, call.Args[1], nil, call)
}

func parsePattern(pat string) (method, path string) {
	method = "*"
	rest := strings.TrimSpace(pat)
	if i := strings.IndexAny(rest, " \t"); i > 0 && isMethodToken(rest[:i]) {
		method = rest[:i]
		rest = strings.TrimSpace(rest[i+1:])
	}
	if !strings.HasPrefix(rest, "/") {
		if j := strings.IndexByte(rest, '/'); j >= 0 {
			rest = rest[j:] // drop a host: "GET example.com/x"
		}
	}
	return method, rest
}

func isMethodToken(s string) bool {
	switch s {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE":
		return true
	}
	return false
}

// mountOrRoute: chi's Mount/Handle take either a router (a mount) or a
// plain handler (a route that matches everything under the path).
func (w *walker) mountOrRoute(n *node, method string, pathX, h ast.Expr, call *ast.CallExpr, isMount bool) {
	if wi := w.evalValue(h); wi.n != nil && wi.n.canon() != n.canon() {
		if !isMount {
			// r.Handle("/v1/*", sub): chi doesn't strip, sub sees the full path
			pat, _ := w.constString(pathX)
			w.mount(n, "", wi, call, true, pat)
			return
		}
		p, _ := w.constString(pathX)
		p = strings.TrimSuffix(strings.TrimSuffix(p, "*"), "/")
		w.mount(n, p, wi, call, false, p)
		return
	}
	if isMount {
		p, dyn := w.pathArg(pathX)
		if !dyn {
			p = strings.TrimSuffix(p, "/") + "/*"
		}
		w.registerPath(n, method, p, dyn, h, nil, call)
		return
	}
	w.register(n, method, pathX, h, nil, call)
}

// mount attaches a router under the parent. With absolute set the router
// sees the request's full path (chi Handle, net/http, gorilla), so its own
// paths are not under the parent's prefix; otherwise (chi/fiber Mount) the
// prefix is stripped and prepended. Mounting makes it a group of the parent:
// siblings are judged within it, not against the parent's routes.
func (w *walker) mount(parent *node, prefix string, wi wrapInfo, call *ast.CallExpr, absolute bool, pattern string) {
	sub := wi.n.canon()
	for p := parent; p != nil; p = p.parent {
		if p.canon() == sub {
			return // a router mounted under itself
		}
	}
	// each mount is its own edge: a router mounted open at /internal and
	// behind auth at /admin has its routes under both, each with its chain.
	// parent stays the view the mount was made on (r.With(auth).Mount).
	e := edge{parent: parent, prefix: prefix, absolute: absolute, wrappers: append([]mwEntry(nil), wi.wrappers...)}
	for _, old := range sub.edges {
		if old.parent == e.parent && old.prefix == e.prefix && old.absolute == e.absolute && sameNames(old.wrappers, e.wrappers) {
			return
		}
	}
	sub.edges = append(sub.edges, e)
	if sub.kind == kindRoot {
		sub.kind = kindGroup
	}
	if call != nil {
		w.r.mounts = append(w.r.mounts, mountRec{sub: sub, parent: parent, pattern: pattern, file: w.fc.rel, line: w.line(call.Pos()), origin: w.origin})
	}
}

// enterRouterFunc walks a Group/Route callback with its router parameter
// bound to c. The callback may be a literal or a named function.
func (w *walker) enterRouterFunc(fx ast.Expr, c *node) {
	switch f := unparen(fx).(type) {
	case *ast.FuncLit:
		w.consumed[f] = true
		w.walkClosure(f, map[int]*node{0: c}, w.fc)
	case *ast.Ident, *ast.SelectorExpr:
		fake := &ast.CallExpr{Fun: fx, Args: []ast.Expr{&ast.Ident{Name: "\x00router"}}}
		fi := w.resolveCall(fake)
		if fi == nil {
			return
		}
		saved := w.env.routers["\x00router"]
		w.env.routers["\x00router"] = c
		var recv ast.Expr
		if sel, ok := fx.(*ast.SelectorExpr); ok {
			recv = sel.X
		}
		w.inline(fi, fake, recv)
		if saved != nil {
			w.env.routers["\x00router"] = saved
		} else {
			delete(w.env.routers, "\x00router")
		}
	}
}

// register records a route whose path is the first-class argument pathX.
func (w *walker) register(n *node, method string, pathX, h ast.Expr, mwArgs []ast.Expr, call *ast.CallExpr) {
	p, dyn := w.pathArg(pathX)
	w.registerPath(n, method, p, dyn, h, mwArgs, call)
}

func (w *walker) registerPath(n *node, method, p string, dyn bool, h ast.Expr, mwArgs []ast.Expr, call *ast.CallExpr) {
	if method == "" {
		method = "*"
	}
	extra := w.mwNames(mwArgs)
	ws, base := w.unwrapHandler(h)
	extra = append(extra, ws...)
	if k := keyOf(base); k != "" {
		extra = append(extra, w.env.hwraps[k]...)
	}
	// wrapper.DeleteUser on oapi-codegen's ServerInterfaceWrapper runs the
	// wrapper's HandlerMiddlewares around the handler
	if sel, ok := unparen(base).(*ast.SelectorExpr); ok {
		if rk := keyOf(sel.X); rk != "" {
			var ks []string
			for k := range w.env.chains {
				if strings.HasPrefix(k, rk+".") && strings.Contains(k[len(rk)+1:], "Middleware") {
					ks = append(ks, k)
				}
			}
			sort.Strings(ks)
			for _, k := range ks {
				extra = append(extra, w.env.chains[k]...)
			}
		}
	}
	if fl, ok := unparen(base).(*ast.FuncLit); ok {
		w.consumed[fl] = true
	}
	var anc map[*node]int
	if n.fw == fwFiber {
		// fiber runs handlers in registration order across the whole app: a
		// group created before app.Use(auth) still gets auth on routes added
		// after it
		anc = map[*node]int{}
		for a, i := n, 0; a != nil && i < 4096; a, i = a.parent, i+1 {
			anc[a] = len(a.mw)
		}
	}
	api := w.apiKind(base, 0)
	if api == "" && dyn && (strings.Contains(p, ".Method(") || strings.Contains(p, "reflect.")) {
		api = "routes built by reflection"
	}
	w.r.recs = append(w.r.recs, routeRec{
		anc: anc,
		n:   n, nmw: len(n.mw), method: method, path: p, dynamic: dyn,
		handler: exprName(base), extra: extra, checks: w.handlerChecks(base),
		file: w.fc.rel, line: w.line(call.Pos()), origin: w.origin, api: api,
	})
}

// serveAPI: an API server handed straight to a server (http.Server{Handler:
// gw}) is a route of its own: everything under "/".
func (w *walker) serveAPI(x ast.Expr) {
	api := w.apiKind(x, 0)
	if api == "" {
		return
	}
	n := w.r.newNode(fwNetHTTP, kindRoot, w.site(x.Pos()))
	w.registerPath(n, "*", "/", false, x, nil, &ast.CallExpr{Fun: x})
	w.r.recs[len(w.r.recs)-1].api = api
}

// apiKind names an API server handed to the router as a single handler,
// whose operations the walk can't see: grpc-gateway's ServeMux, gqlgen,
// connect-rpc, a generated OpenAPI server built outside the tree.
func (w *walker) apiKind(x ast.Expr, depth int) string {
	if x == nil || depth > 4 {
		return ""
	}
	switch e := unparen(x).(type) {
	case *ast.Ident:
		if o, ok := w.env.origins[e.Name]; ok {
			return w.apiKind(o, depth+1)
		}
	case *ast.CallExpr:
		if sel, ok := unparen(e.Fun).(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && !w.env.locals[id.Name] {
				ip, n := w.fc.imports[id.Name], sel.Sel.Name
				switch {
				case strings.Contains(ip, "grpc-gateway") && n == "NewServeMux":
					return "grpc-gateway"
				case strings.Contains(ip, "gqlgen/graphql/handler"), strings.Contains(ip, "graph-gophers/graphql-go"),
					strings.Contains(ip, "graphql-go/handler"):
					return "GraphQL (" + id.Name + "." + n + ")"
				case strings.HasSuffix(ip, "connect") && strings.HasPrefix(n, "New") && strings.HasSuffix(n, "Handler"):
					return "connect-rpc (" + id.Name + "." + n + ")"
				case n == "NewServer" && ip != "" && w.ix.byImp[ip] == nil && w.fc.fws[id.Name] == nil:
					return "a generated API server (" + id.Name + "." + n + ")"
				}
			}
		}
		// http.StripPrefix("/api", srv), cors(gql)
		for _, a := range e.Args {
			if k := w.apiKind(a, depth+1); k != "" {
				return k
			}
		}
	}
	return ""
}

func lastOrNil(xs []ast.Expr) ast.Expr {
	if len(xs) == 0 {
		return nil
	}
	return xs[len(xs)-1]
}

func initOf(xs []ast.Expr) []ast.Expr {
	if len(xs) == 0 {
		return nil
	}
	return xs[:len(xs)-1]
}

func (w *walker) stringList(x ast.Expr) []string {
	cl, ok := unparen(x).(*ast.CompositeLit)
	if !ok {
		return nil
	}
	var out []string
	for _, e := range cl.Elts {
		if s, ok := w.constString(e); ok {
			out = append(out, s)
		}
	}
	return out
}

// ---- gorilla/mux ----

var gorillaBuilder = set("Path", "PathPrefix", "Methods", "Host", "Schemes", "Headers",
	"HeadersRegexp", "Queries", "Name", "MatcherFunc", "Subrouter", "NewRoute",
	"Handler", "HandlerFunc", "Handle", "HandleFunc", "BuildOnly", "BuildVarsFunc", "SkipClean", "StrictSlash", "UseEncodedPath")

type gstep struct {
	name string
	call *ast.CallExpr
}

// gorillaSteps flattens r.PathPrefix("/a").Methods("GET").HandlerFunc(h)
// into its base receiver and the calls applied to it, innermost first.
func gorillaSteps(call *ast.CallExpr) (ast.Expr, []gstep) {
	var steps []gstep
	var x ast.Expr = call
	for {
		c, ok := unparen(x).(*ast.CallExpr)
		if !ok {
			break
		}
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok || !gorillaBuilder[sel.Sel.Name] {
			break
		}
		steps = append(steps, gstep{name: sel.Sel.Name, call: c})
		x = sel.X
	}
	for i, j := 0, len(steps)-1; i < j; i, j = i+1, j-1 {
		steps[i], steps[j] = steps[j], steps[i]
	}
	return x, steps
}

// gorillaRoot strips a call chain down to the receiver it starts from:
// r.PathPrefix("/a").Subrouter().HandleFunc(...) -> r.
func gorillaRoot(call *ast.CallExpr) ast.Expr {
	var x ast.Expr = call
	for {
		c, ok := unparen(x).(*ast.CallExpr)
		if !ok {
			return unparen(x)
		}
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok {
			return nil
		}
		x = sel.X
	}
}

// peek looks a receiver up without evaluating calls.
func (w *walker) peek(x ast.Expr) *node {
	switch unparen(x).(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return w.lookup(x)
	}
	return nil
}

func (w *walker) gorillaChain(call *ast.CallExpr) *node {
	base, steps := gorillaSteps(call)
	if len(steps) == 0 {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if n := w.eval(sel.X); n != nil {
				return w.routerMethod(n, sel.Sel.Name, call)
			}
		}
		return nil
	}
	n := w.eval(base)
	if n == nil || n.fw != fwGorilla {
		return nil
	}
	var (
		p        string
		dyn      bool
		isPrefix bool
		methods  []string
		handler  ast.Expr
		at       *ast.CallExpr
	)
	cur := n
	for _, s := range steps {
		args := s.call.Args
		switch s.name {
		case "Path", "PathPrefix":
			if len(args) == 1 {
				p, dyn = w.pathArg(args[0])
				isPrefix = s.name == "PathPrefix"
			}
		case "Methods":
			for i := range args {
				if m, ok := w.constString(args[i]); ok {
					methods = append(methods, strings.ToUpper(m))
				}
			}
		case "HandleFunc", "Handle":
			if len(args) >= 2 {
				p, dyn = w.pathArg(args[0])
				handler, at = args[len(args)-1], s.call
			}
		case "Handler", "HandlerFunc":
			if len(args) == 1 {
				handler, at = args[0], s.call
			}
		case "Subrouter":
			prefix := ""
			if isPrefix || p != "" {
				prefix = p
			}
			cur = w.derive(cur, kindGroup, prefix, nil, s.call)
			p, dyn, isPrefix, methods, handler = "", false, false, nil, nil
		}
	}
	if handler == nil {
		switch steps[len(steps)-1].name {
		case "Subrouter", "StrictSlash", "SkipClean", "UseEncodedPath":
			return cur
		}
		return nil
	}
	if wi := w.evalValue(handler); wi.n != nil && wi.n.canon() != cur.canon() {
		if isPrefix {
			// the sub-router sees the path StripPrefix left it
			w.mount(cur, wi.strip, wi, at, true, p)
			return nil
		}
	}
	if isPrefix && !dyn {
		p = strings.TrimSuffix(p, "/") + "/*"
	}
	if len(methods) == 0 {
		methods = []string{"*"}
	}
	for _, m := range methods {
		w.registerPath(cur, m, p, dyn, handler, nil, at)
	}
	return nil
}
