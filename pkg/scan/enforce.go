package scan

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

// Reading a middleware's own body, when it's declared in the tree, settles
// what its name can't, in both directions. miniflux's webSessionMiddleware
// redirects anonymous users to login, so it enforces although "session"
// alone doesn't say so. A RequireAuth that logs a bad token and calls next
// anyway never rejects anything, whatever its name says.
//
// Evidence of rejection is a 401 (status constant, literal, or a framework's
// ErrUnauthorized), a redirect to a login URL, or a branch on a missing
// identity, and only on a branch that then stops the request: returns
// without handing it on, or aborts. A 401 written before falling through to
// next doesn't reject. It must also come before the first unconditional
// hand-off: error-mapping middleware writes 401s after calling next.

type verdict int

const (
	vUnknown     verdict = iota
	vEnforces            // rejects some requests on an authentication check
	vPasses              // every path hands the request on; nothing stops it
	vConditional         // stops only under a flag, for some methods, or when a credential is sent
	vComplex             // too many branches to read: unverified, whatever the name
)

func (w *walker) enforces(x ast.Expr) bool { return w.mwVerdict(x) == vEnforces }

// mwVerdict reads x as middleware. The role matters for gin and hertz: a
// middleware that returns without aborting lets the chain run on, while a
// handler that returns has answered the request.
func (w *walker) mwVerdict(x ast.Expr) verdict {
	saved := w.r.asMW
	w.r.asMW = true
	defer func() { w.r.asMW = saved }()
	return w.verdictOf(x, 0)
}

// handlerChecks reads a route's handler the same way: many apps check the
// caller inside the handler (a token lookup that answers 401) rather than in
// middleware. Such routes are still reported, at low confidence.
func (w *walker) handlerChecks(h ast.Expr) bool {
	if h == nil {
		return false
	}
	saved := w.r.asMW
	w.r.asMW = false
	defer func() { w.r.asMW = saved }()
	return w.verdictOf(h, 0) == vEnforces
}

func (w *walker) verdictOf(x ast.Expr, d int) verdict {
	if x == nil || d > 4 {
		return vUnknown
	}
	switch e := unparen(x).(type) {
	case *ast.FuncLit:
		return w.bodyVerdict(e.Body, w.fc, 0, nil, e.Type)
	case *ast.CallExpr:
		if fi := w.resolveCall(e); fi != nil {
			return w.funcVerdict(fi)
		}
		// mw.New(cfg).Handler: judged by the factory
		if sel, ok := unparen(e.Fun).(*ast.SelectorExpr); ok {
			if c, ok := unparen(sel.X).(*ast.CallExpr); ok {
				return w.verdictOf(c, d+1)
			}
		}
	case *ast.Ident:
		if v, ok := w.env.enforce[e.Name]; ok {
			return v
		}
		if o, ok := w.env.origins[e.Name]; ok {
			return w.verdictOf(o, d+1)
		}
		if w.env.locals[e.Name] {
			return vUnknown
		}
		if fi := w.fc.pkg.funcs[e.Name]; fi != nil {
			return w.funcVerdict(fi)
		}
		if g := w.fc.pkg.globals[e.Name]; g != nil && g.init != nil {
			return w.with(g.fc).verdictOf(g.init, d+1)
		}
	case *ast.SelectorExpr:
		if id, ok := e.X.(*ast.Ident); ok && !w.env.locals[id.Name] {
			if _, isImport := w.fc.imports[id.Name]; isImport {
				if fi := w.ix.lookupFunc(w.fc, id.Name, e.Sel.Name); fi != nil {
					return w.funcVerdict(fi)
				}
				if p := w.ix.byImp[w.fc.imports[id.Name]]; p != nil {
					if g := p.globals[e.Sel.Name]; g != nil && g.init != nil {
						return w.with(g.fc).verdictOf(g.init, d+1)
					}
				}
				return vUnknown
			}
		}
		if info := w.fc.pkg.info; info != nil {
			if f, ok := info.Uses[e.Sel].(*types.Func); ok {
				if fi := w.ix.byObj[f]; fi != nil {
					return w.funcVerdict(fi)
				}
			}
		}
		if tr, ok := w.typeOf(e.X); ok {
			if fi := w.methodOf(tr, e.Sel.Name, 0); fi != nil {
				return w.funcVerdict(fi)
			}
		}
		// (&AuthMiddleware{}).Wrap, NewAuth(cfg).Handler
		switch unparen(e.X).(type) {
		case *ast.CompositeLit, *ast.UnaryExpr, *ast.CallExpr:
			if tr, ok := w.typeOfValue(unparen(e.X)); ok {
				if fi := w.methodOf(tr, e.Sel.Name, 0); fi != nil {
					return w.funcVerdict(fi)
				}
			}
		}
		if id, ok := e.X.(*ast.Ident); ok {
			if o, ok := w.env.origins[id.Name]; ok {
				if tr, ok := w.typeOfValue(o); ok {
					if fi := w.methodOf(tr, e.Sel.Name, 0); fi != nil {
						return w.funcVerdict(fi)
					}
				}
			}
		}
	}
	return vUnknown
}

func (w *walker) funcVerdict(fi *funcInfo) verdict {
	k := fnKey{fi, w.r.asMW}
	if w.r.enforceSet[k] {
		w.r.lastWhy = w.r.whyOf[fi.decl.Body]
		return w.r.enforce[k]
	}
	w.r.enforceSet[k] = true // a recursive middleware reads as unknown until done
	v := vUnknown
	if fi.decl.Body != nil {
		v = w.bodyVerdict(fi.decl.Body, fi.fc, 0, fi, fi.decl.Type)
	}
	w.r.lastWhy = w.r.whyOf[fi.decl.Body]
	w.r.enforce[k] = v
	return v
}

type bodyKey struct {
	b     *ast.BlockStmt
	depth int
	mw    bool
}

type fnKey struct {
	fi *funcInfo
	mw bool
}

// configParams collects parameters that configure a wrapper rather than
// carry the request, the next handler or an identity lookup.
func configParams(ft *ast.FuncType, identity, into map[string]bool) {
	if ft == nil || ft.Params == nil {
		return
	}
	for _, p := range ft.Params.List {
		t := strings.ToLower(typeBaseName(p.Type))
		switch {
		case strings.Contains(t, "handler"), strings.Contains(t, "responsewriter"), t == "request",
			t == "context", t == "ctx", t == "":
			continue
		}
		for _, n := range p.Names {
			if !identity[n.Name] {
				into[n.Name] = true
			}
		}
	}
}

func usesAny(x ast.Expr, names map[string]bool) bool {
	if len(names) == 0 {
		return false
	}
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && names[id.Name] {
			hit = true
		}
		return !hit
	})
	return hit
}

// identityParams: parameters that look the caller up, by type or name:
// `a authFunc`, `getUser func(*http.Request) *User`.
func identityParams(ft *ast.FuncType, into map[string]bool) {
	if ft == nil || ft.Params == nil {
		return
	}
	for _, p := range ft.Params.List {
		// by type the bar is higher: UserLevelFunc configures, authFunc looks up
		t := strings.ToLower(typeBaseName(p.Type))
		typeHit := false
		for _, w := range []string{"auth", "session", "token", "claims", "principal", "identity", "credential"} {
			if strings.Contains(t, w) && !strings.Contains(t, "author") {
				typeHit = true
			}
		}
		for _, n := range p.Names {
			if typeHit || identityCall(n.Name) {
				into[n.Name] = true
			}
		}
	}
}

// externParams marks parameters whose type comes from outside the tree
// (*gin.Context, http.ResponseWriter): a method call on them is never one of
// ours, however its name matches.
func (w *walker) externParams(ft *ast.FuncType) {
	if ft == nil || ft.Params == nil {
		return
	}
	for _, p := range ft.Params.List {
		if !w.externType(p.Type) {
			continue
		}
		for _, n := range p.Names {
			w.env.extern[n.Name] = true
		}
	}
}

func (w *walker) externType(t ast.Expr) bool {
	for {
		switch e := t.(type) {
		case *ast.StarExpr:
			t = e.X
			continue
		case *ast.SelectorExpr:
			if id, ok := e.X.(*ast.Ident); ok {
				ip, isImport := w.fc.imports[id.Name]
				return isImport && w.ix.byImp[ip] == nil
			}
		}
		return false
	}
}

func (w *walker) bodyVerdict(body *ast.BlockStmt, fc *fileCtx, depth int, owner *funcInfo, ft *ast.FuncType) verdict {
	if body == nil || depth > 2 {
		return vUnknown
	}
	key := bodyKey{body, depth, w.r.asMW}
	if v, ok := w.r.bodyCache[key]; ok {
		w.r.lastWhy = w.r.whyOf[body]
		return v
	}
	w.r.bodyCache[key] = vUnknown // cycles read as unknown
	v := w.bodyVerdictUncached(body, fc, depth, owner, ft)
	w.r.bodyCache[key] = v
	w.r.lastWhy = w.r.whyOf[body]
	return v
}

// helperRejects: a helper called from a handler lends its 401 to the caller
// only if it writes one unconditionally (func unauthorized(w)). A shared
// error writer that maps some error to 401 (c2's writeErr) doesn't.
func (w *walker) helperRejects(fi *funcInfo, depth int) bool {
	if fi.decl.Body == nil || depth > 2 {
		return false
	}
	for _, st := range fi.decl.Body.List {
		found := false
		ast.Inspect(st, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncLit, *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SelectStmt:
				_ = x
				return false
			}
			if rejectionMark(n) {
				found = true
			}
			return !found
		})
		if found {
			return true
		}
	}
	return false
}

// rejectionMark: a node that by itself says "unauthorized".
func rejectionMark(n ast.Node) bool {
	switch x := n.(type) {
	case *ast.Ident:
		return isUnauthorizedName(x.Name)
	case *ast.SelectorExpr:
		return isUnauthorizedName(x.Sel.Name)
	case *ast.CallExpr:
		for _, a := range x.Args {
			if lit, ok := unparen(a).(*ast.BasicLit); ok && lit.Kind == token.INT && lit.Value == "401" {
				return true
			}
		}
		return strings.Contains(exprName(x.Fun), "Redirect") && mentionsLogin(x)
	}
	return false
}

func (w *walker) bodyVerdictUncached(body *ast.BlockStmt, fc *fileCtx, depth int, owner *funcInfo, ft *ast.FuncType) verdict {
	saved := nextParams
	nextParams = handlerParams(ft)
	for k := range saved {
		if _, ok := nextParams[k]; !ok && depth > 0 {
			nextParams[k] = true
		}
	}
	defer func() { nextParams = saved }()
	enforces := false
	resolver := w.with(fc)
	params := map[string]bool{}
	// configuration parameters: a level, flag or role the caller passes in.
	// A rejection that depends on one can't be credited without the value
	// (writefreely's handler.Web(fn, UserLevelReader) vs UserLevelUser).
	cfgParams := map[string]bool{}
	identityParams(ft, params)
	configParams(ft, params, cfgParams)
	resolver.externParams(ft)
	if owner != nil && owner.recvType != "" && len(owner.decl.Recv.List[0].Names) > 0 {
		rn := owner.decl.Recv.List[0].Names[0].Name
		resolver.env.locals[rn] = true
		resolver.env.vtypes[rn] = typeRef{pkg: fc.pkg, name: owner.recvType}
	}
	// handler scopes are read path by path; the evidence scan below only
	// covers what they don't (callbacks, closures run by helpers, bodies
	// that never hand the request on)
	pv, why, covered, pathOK := w.pathVerdict(body, fc, depth, owner, ft, resolver, params, cfgParams)
	// handlers are the closures a middleware returns (or the body itself);
	// a closure passed to some other call is a callback, where only a 401
	// or a login redirect counts, not a bare identity branch
	returned := map[*ast.FuncLit]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		if rs, ok := n.(*ast.ReturnStmt); ok {
			for _, r := range rs.Results {
				if fl := returnedLit(r); fl != nil {
					returned[fl] = true
				}
			}
		}
		return true
	})

	var scope func(b *ast.BlockStmt, callback, handler bool)
	scope = func(root *ast.BlockStmt, callback, handler bool) {
		skipEvidence := pathOK && covered[root]
		// The hand-off that counts is the first unconditional one: a
		// top-level statement. `if static { next.ServeHTTP(); return }`
		// early in the body doesn't end the checks after it.
		next := token.NoPos
		for _, st := range root.List {
			var c *ast.CallExpr
			switch s := st.(type) {
			case *ast.ExprStmt:
				c, _ = unparen(s.X).(*ast.CallExpr)
			case *ast.ReturnStmt:
				if len(s.Results) == 1 {
					c, _ = unparen(s.Results[0]).(*ast.CallExpr)
				}
			}
			if c != nil && isHandOff(c) {
				next = c.Pos()
				break
			}
		}
		// variables holding the caller's identity: s := Auth(c, ...),
		// u, err := getUserSession(r)
		ids := map[string]bool{}
		ast.Inspect(root, func(n ast.Node) bool {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			as, ok := n.(*ast.AssignStmt)
			if !ok || len(as.Rhs) != 1 {
				return true
			}
			call, ok := unparen(as.Rhs[0]).(*ast.CallExpr)
			if !ok {
				return true
			}
			fn, isID := unparen(call.Fun).(*ast.Ident)
			if !identityCall(exprName(call.Fun)) && !(isID && params[fn.Name]) {
				return true
			}
			for _, l := range as.Lhs {
				if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
					ids[id.Name] = true
				}
			}
			return true
		})
		var evidence []token.Pos
		mark := func(p token.Pos) { evidence = append(evidence, p) }

		// marks reports rejection evidence inside one statement's own
		// expressions (nested blocks are walked as blocks).
		marks := func(st ast.Node, present bool) bool {
			if present || st == nil {
				return false
			}
			hit := false
			ast.Inspect(st, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncLit, *ast.BlockStmt:
					return false
				case *ast.CallExpr:
					if isHandOff(x) {
						return true
					}
					if depth < 2 {
						if fi := resolver.resolveCall(x); fi != nil && w.helperRejects(fi, depth+1) {
							hit = true
						}
					}
				}
				if rejectionMark(n) {
					hit = true
				}
				return !hit
			})
			return hit
		}

		// visitBlock walks a statement list. present is set inside a branch
		// that only runs when the caller already has an identity (`if err
		// == nil`, `if user != nil`): a 401 there rejects bad credentials but
		// lets anonymous requests through, so it doesn't count.
		var visitBlock func(list []ast.Stmt, present bool)
		var visitStmt func(st ast.Stmt, list []ast.Stmt, i int, present bool)
		visitBlock = func(list []ast.Stmt, present bool) {
			for i, st := range list {
				visitStmt(st, list, i, present)
			}
		}
		visitStmt = func(st ast.Stmt, list []ast.Stmt, i int, present bool) {
			// closures inside the statement are their own scopes; deferred
			// and concurrent ones can't stop the request
			ast.Inspect(st, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.DeferStmt, *ast.GoStmt:
					return false
				}
				if fl, ok := n.(*ast.FuncLit); ok {
					identityParams(fl.Type, params)
					resolver.externParams(fl.Type)
					// a zero-argument closure run by a helper is the handler's
					// own logic deferred (writefreely's handleHTTPError)
					cb := !returned[fl] && fl.Type.Params != nil && len(fl.Type.Params.List) > 0
					scope(fl.Body, cb, returned[fl])
					return false
				}
				if _, ok := n.(*ast.BlockStmt); ok && n != st {
					return false
				}
				return true
			})
			switch s := st.(type) {
			case *ast.IfStmt:
				if s.Init != nil {
					visitStmt(s.Init, nil, 0, present)
				}
				if usesAny(s.Cond, cfgParams) {
					return // depends on the caller's argument
				}
				if !present && !callback && identityCheck(s, ids) {
					mark(s.Pos())
				}
				visitBlock(s.Body.List, present || hasIdentity(s.Cond, ids))
				if s.Else != nil {
					visitStmt(s.Else, nil, 0, present || negative(s.Cond) && mentionsIdentity(s.Cond, ids))
				}
			case *ast.BlockStmt:
				visitBlock(s.List, present)
			case *ast.ForStmt:
				visitBlock(s.Body.List, present)
			case *ast.RangeStmt:
				visitBlock(s.Body.List, present)
			case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
				var body *ast.BlockStmt
				switch b := s.(type) {
				case *ast.SwitchStmt:
					body = b.Body
				case *ast.TypeSwitchStmt:
					body = b.Body
				case *ast.SelectStmt:
					body = b.Body
				}
				for _, c := range body.List {
					switch cc := c.(type) {
					case *ast.CaseClause:
						visitBlock(cc.Body, present)
					case *ast.CommClause:
						visitBlock(cc.Body, present)
					}
				}
			case *ast.LabeledStmt:
				visitStmt(s.Stmt, list, i, present)
			default:
				// a rejection mark counts when this branch then stops: the
				// statement is itself a return/abort, or one follows it here
				if marks(st, present) && stopsFrom(list, i) {
					mark(st.Pos())
				}
			}
		}
		visitBlock(root.List, false)
		for _, p := range evidence {
			if !skipEvidence && (next == token.NoPos || p < next) {
				enforces = true
			}
		}
	}
	// the body itself is a handler scope when it hands off directly (gin's
	// func(c *gin.Context) { ...; c.Next() })
	scope(body, false, true)
	if pathOK {
		switch pv {
		case vEnforces:
			return vEnforces
		case vConditional:
			w.r.whyOf[body] = why
			return vConditional
		case vComplex:
			return vComplex
		case vPasses:
			// a callback handed to another call may still reject; a local
			// closure whose result is ignored doesn't
			if callbackRejects(body, covered) || depth > 0 || defersRejection(body) {
				return vUnknown
			}
			return vPasses
		}
		if enforces {
			return vEnforces
		}
		return vUnknown
	}
	if enforces {
		return vEnforces
	}
	// func RequireAuth(next) http.Handler { return Auth{next: next} }: the
	// type's ServeHTTP is the middleware
	if v, ok := w.servedBy(body, resolver); ok {
		return v
	}
	return vUnknown
}

// servedBy finds `return T{...}` or `return &T{...}` of an in-tree type
// with a ServeHTTP method and reads that method.
func (w *walker) servedBy(body *ast.BlockStmt, resolver *walker) (verdict, bool) {
	for _, st := range body.List {
		rs, ok := st.(*ast.ReturnStmt)
		if !ok || len(rs.Results) != 1 || compositeOf(rs.Results[0]) == nil {
			continue
		}
		tr, ok := resolver.typeOfValue(rs.Results[0])
		if !ok {
			continue
		}
		if fi := resolver.methodOf(tr, "ServeHTTP", 0); fi != nil {
			v := w.funcVerdict(fi)
			if v == vConditional {
				w.r.whyOf[body] = w.r.whyOf[fi.decl.Body]
			}
			return v, true
		}
	}
	return vUnknown, false
}

// callbackRejects: a function literal passed straight to a call (an error
// handler in a library's options) writes a rejection.
func callbackRejects(body *ast.BlockStmt, covered map[*ast.BlockStmt]bool) bool {
	hit := false
	ast.Inspect(body, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok || hit {
			return !hit
		}
		for _, a := range c.Args {
			fl, ok := unparen(a).(*ast.FuncLit)
			if !ok || covered[fl.Body] {
				continue // the handler itself, already read path by path
			}
			ast.Inspect(fl.Body, func(m ast.Node) bool {
				if rejectionMark(m) {
					hit = true
				}
				return !hit
			})
		}
		return !hit
	})
	return hit
}

// defersRejection: the middleware stores a rejection for a later check
// instead of answering it (listmonk sets an echo.HTTPError 403 in the
// context and its per-route wrappers return it). Such a middleware hands
// every request on, but it isn't a no-op; leave the verdict to its name.
func defersRejection(body *ast.BlockStmt) bool {
	hit := false
	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || hit {
			return !hit
		}
		name := exprName(call.Fun)
		if i := strings.LastIndexByte(name, '.'); i >= 0 {
			name = name[i+1:]
		}
		switch name {
		case "Set", "SetValue", "WithValue", "Store", "Put":
		default:
			return true
		}
		for _, a := range call.Args {
			ast.Inspect(a, func(m ast.Node) bool {
				if rejectionMark(m) || forbiddenMark(m) {
					hit = true
				}
				return !hit
			})
		}
		return !hit
	})
	return hit
}

func forbiddenMark(n ast.Node) bool {
	switch x := n.(type) {
	case *ast.SelectorExpr:
		return x.Sel.Name == "StatusForbidden" || strings.HasPrefix(x.Sel.Name, "ErrForbidden")
	case *ast.BasicLit:
		return x.Kind == token.INT && x.Value == "403"
	}
	return false
}

// stopsFrom: list[i] is a return or abort, or a later statement in the same
// list is one (`http.Error(w, "", 401); return`).
func stopsFrom(list []ast.Stmt, i int) bool {
	if list == nil {
		return false
	}
	for _, st := range list[i:] {
		if isStop(st) {
			return true
		}
	}
	return false
}

func isStop(st ast.Stmt) bool {
	switch s := st.(type) {
	case *ast.ReturnStmt:
		if len(s.Results) == 1 {
			if c, ok := unparen(s.Results[0]).(*ast.CallExpr); ok && isHandOff(c) {
				return false
			}
		}
		return true
	case *ast.ExprStmt:
		if c, ok := unparen(s.X).(*ast.CallExpr); ok {
			name := exprName(c.Fun)
			if i := strings.LastIndexByte(name, '.'); i >= 0 {
				name = name[i+1:]
			}
			return strings.HasPrefix(name, "Abort") || name == "panic"
		}
	}
	return false
}

func returnsHandOff(s *ast.ReturnStmt) bool {
	for _, r := range s.Results {
		if c, ok := unparen(r).(*ast.CallExpr); ok && isHandOff(c) {
			return true
		}
	}
	return false
}

// identityCheck: an if on the caller's identity where one side stops the
// request: `if user == nil { return ctx.NotFound() }` or
// `if user != nil { return next(ctx) }` followed by a rejection.
func identityCheck(s *ast.IfStmt, ids map[string]bool) bool {
	if !authnCond(s.Cond, ids) {
		return false
	}
	stops, hands := false, false
	ast.Inspect(s.Body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt:
			if len(x.Results) == 1 {
				if c, ok := unparen(x.Results[0]).(*ast.CallExpr); ok && isHandOff(c) {
					hands = true
					return true
				}
				if id, ok := unparen(x.Results[0]).(*ast.Ident); ok && id.Name == "nil" {
					return true
				}
			}
			stops = true
		case *ast.CallExpr:
			if isHandOff(x) {
				hands = true
			} else if isAbort(x) {
				stops = true
			}
		}
		return true
	})
	// a branch that hands off only when the identity is present implies the
	// fallthrough rejects; one that stops when it's absent rejects directly.
	// `if user == nil { return next(c) }` lets anonymous requests through.
	switch {
	case stops && !hands:
		// only a stop on absence rejects: `if user != nil { ...return }`
		// is the present-identity path bailing out on some other error
		return negative(s.Cond)
	case hands && !stops:
		return !negative(s.Cond)
	}
	return false
}

// negative: the condition tests for absence (x == nil, x == "", !ok).
func negative(x ast.Expr) bool {
	switch e := unparen(x).(type) {
	case *ast.CallExpr:
		// if s.Abort(c), if isInvalid(tok)
		name := strings.ToLower(exprName(e.Fun))
		for _, w := range []string{"abort", "deny", "denied", "invalid", "unauthorized", "forbidden", "fail", "missing", "expired", "reject"} {
			if strings.Contains(name, w) {
				return true
			}
		}
		return false
	case *ast.UnaryExpr:
		return e.Op == token.NOT
	case *ast.BinaryExpr:
		switch e.Op {
		case token.EQL:
			return true
		case token.NEQ:
			// err != nil after an identity lookup means "no identity"
			return isErrIdent(e.X) || isErrIdent(e.Y)
		case token.LOR, token.LAND:
			return negative(e.X) || negative(e.Y)
		}
	}
	return false
}

var identityWords = []string{"user", "session", "token", "auth", "claims", "principal", "login", "logged",
	"admin", "identity", "account", "role", "permission", "authenticated", "member", "owner"}

// hasIdentity: the condition holds only when the caller has an identity:
// err == nil after a lookup, user != nil, ok.
func hasIdentity(x ast.Expr, ids map[string]bool) bool {
	switch e := unparen(x).(type) {
	case *ast.Ident:
		return e.Name == "ok" || ids[e.Name]
	case *ast.BinaryExpr:
		switch e.Op {
		case token.LAND:
			return hasIdentity(e.X, ids) || hasIdentity(e.Y, ids)
		case token.EQL:
			return (isErrIdent(e.X) || isErrIdent(e.Y)) && (isNil(e.X) || isNil(e.Y))
		case token.NEQ:
			return !isErrIdent(e.X) && !isErrIdent(e.Y) && (isNil(e.X) || isNil(e.Y)) && mentionsIdentity(e, ids)
		}
	}
	return false
}

func isNil(x ast.Expr) bool {
	id, ok := unparen(x).(*ast.Ident)
	return ok && id.Name == "nil"
}

func isErrIdent(x ast.Expr) bool {
	id, ok := unparen(x).(*ast.Ident)
	return ok && (id.Name == "err" || strings.HasSuffix(id.Name, "Err"))
}

// identityCall: a call that looks up or checks the caller, by its name.
func identityCall(name string) bool {
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		name = name[i+1:]
	}
	low := strings.ToLower(name)
	// UpdateUser, CreateSession: work on an identity, don't look one up
	for _, v := range []string{"update", "create", "delete", "save", "set", "insert", "remove", "list",
		"count", "add", "put", "patch", "send", "notify", "log", "store", "write", "invalidate", "revoke"} {
		if strings.HasPrefix(low, v) {
			return false
		}
	}
	for _, w := range identityWords {
		if strings.Contains(low, w) && !strings.Contains(low, "author") {
			return true
		}
	}
	return false
}

// authnCond: the condition asks whether there is a caller at all, as
// opposed to what the caller may do: user == nil, err != nil after a
// lookup, !ok, !IsAuthenticated(r), s.Abort(c). `!user.HasPerm(x)` on an
// identity nobody checked for is authorization only.
func authnCond(x ast.Expr, ids map[string]bool) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.BinaryExpr:
			if e.Op == token.EQL || e.Op == token.NEQ {
				for _, pair := range [][2]ast.Expr{{e.X, e.Y}, {e.Y, e.X}} {
					if isNil(pair[1]) || isEmptyString(pair[1]) {
						if id := rootIdent(pair[0]); id != nil && (ids[id.Name] || mentionsIdentity(pair[0], ids)) {
							hit = true
						}
					}
				}
			}
		case *ast.UnaryExpr:
			if id, ok := unparen(e.X).(*ast.Ident); ok && e.Op == token.NOT && ids[id.Name] {
				hit = true // !ok after a lookup
			}
		case *ast.CallExpr:
			name := strings.ToLower(exprName(e.Fun))
			for _, w := range []string{"authenticated", "loggedin", "signedin", "abort", "deny", "anonymous", "guest", "isvalid",
				"invalid", "expired", "unauthorized"} {
				if strings.Contains(name, w) {
					hit = true
				}
			}
		case *ast.SelectorExpr:
			low := strings.ToLower(e.Sel.Name)
			if strings.Contains(low, "authenticated") || strings.Contains(low, "loggedin") || strings.Contains(low, "anonymous") {
				hit = true
			}
		}
		return !hit
	})
	return hit
}

func isEmptyString(x ast.Expr) bool {
	lit, ok := unparen(x).(*ast.BasicLit)
	return ok && lit.Kind == token.STRING && (lit.Value == `""` || lit.Value == "``")
}

func mentionsIdentity(x ast.Expr, ids map[string]bool) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if ids[id.Name] {
			hit = true
			return false
		}
		low := strings.ToLower(id.Name)
		for _, w := range identityWords {
			if strings.Contains(low, w) && !strings.Contains(low, "author") {
				hit = true
			}
		}
		return !hit
	})
	return hit
}

func isAbort(c *ast.CallExpr) bool {
	name := exprName(c.Fun)
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		name = name[i+1:]
	}
	return strings.HasPrefix(name, "Abort") || name == "Error" || strings.Contains(name, "Redirect") || name == "WriteHeader"
}

// returnedLit finds the func literal in `return func(...)`, also through a
// conversion like http.HandlerFunc(func(...)).
func returnedLit(x ast.Expr) *ast.FuncLit {
	for i := 0; i < 3; i++ {
		switch e := unparen(x).(type) {
		case *ast.FuncLit:
			return e
		case *ast.CallExpr:
			if len(e.Args) != 1 {
				return nil
			}
			x = e.Args[0]
		default:
			return nil
		}
	}
	return nil
}

func isUnauthorizedName(n string) bool {
	switch {
	case n == "StatusUnauthorized", n == "StatusProxyAuthRequired", n == "Unauthenticated":
		return true
	case strings.HasPrefix(n, "ErrUnauthorized"), strings.HasPrefix(n, "ErrUnauthenticated"),
		strings.HasPrefix(n, "ErrNotAuthenticated"), strings.HasPrefix(n, "ErrJWT"),
		strings.HasPrefix(n, "ErrNotLoggedIn"), strings.HasPrefix(n, "ErrMissingToken"):
		return true
	}
	return false
}

// isHandOff: c.Next(), next.ServeHTTP(w, r), next(c), h.ServeHTTP(...),
// or a call of the middleware's own handler parameter, whatever it's named
// (mailpit's middleWareFunc(fn http.HandlerFunc) calling fn(w, r)).
func isHandOff(c *ast.CallExpr) bool {
	switch f := unparen(c.Fun).(type) {
	case *ast.SelectorExpr:
		return f.Sel.Name == "Next" || f.Sel.Name == "ServeHTTP"
	case *ast.Ident:
		return nextName(f.Name)
	}
	return false
}

// nextParams holds the handler parameters of the middleware being read.
// The walk is single-threaded; bodyVerdict sets and restores it.
var nextParams map[string]bool

func nextName(n string) bool {
	return n == "next" || n == "h" || n == "handler" || nextParams[n]
}

// handlerParams: parameters typed as a handler (http.HandlerFunc,
// echo.HandlerFunc, func(http.ResponseWriter, *http.Request)).
func handlerParams(ft *ast.FuncType) map[string]bool {
	out := map[string]bool{}
	if ft == nil || ft.Params == nil {
		return out
	}
	for _, f := range ft.Params.List {
		isH := false
		switch t := f.Type.(type) {
		case *ast.FuncType:
			for _, p := range t.Params.List {
				if n := typeBaseName(p.Type); n == "ResponseWriter" || n == "Context" || n == "RequestContext" || n == "Ctx" {
					isH = true
				}
			}
		default:
			n := typeBaseName(f.Type)
			isH = strings.HasSuffix(n, "Handler") || strings.HasSuffix(n, "HandlerFunc")
		}
		if isH {
			for _, nm := range f.Names {
				out[nm.Name] = true
			}
		}
	}
	return out
}

func mentionsLogin(x ast.Node) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		var s string
		switch e := n.(type) {
		case *ast.BasicLit:
			if e.Kind == token.STRING {
				s = e.Value
			}
		case *ast.Ident:
			s = e.Name
		}
		s = strings.ToLower(s)
		if strings.Contains(s, "login") || strings.Contains(s, "signin") || strings.Contains(s, "sign-in") || strings.Contains(s, "sign_in") {
			hit = true
		}
		return !hit
	})
	return hit
}
