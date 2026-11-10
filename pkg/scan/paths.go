package scan

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strings"
)

// A handler scope (a middleware's returned closure, or a body that hands the
// request on itself) is read path by path. Each path through the code ends
// either having handed the request on or not; one that doesn't is a stop,
// whether it returns, aborts, panics or just falls off the end
// (`if valid { next(); return }; http.Error(401)`). Conditions taken along a
// path are kept as literals, so a stop can be asked what it depends on:
// a stop that turns into a hand-off by flipping only a package-level flag,
// the request method, or whether a credential was sent at all is
// conditional auth, not auth.
//
// gin and hertz run the next handler when a middleware returns without
// aborting, so there only Abort* stops the chain.

type atomKind int

const (
	akRequest  atomKind = iota
	akConfig            // a package-level flag or a configuration field
	akMethod            // r.Method == "GET"
	akPresence          // tok == "" / err == nil after a lookup: whether a credential was sent
	akBypass            // ?debug=1: a request parameter the caller picks
)

type atom struct {
	key   string
	kind  atomKind
	pos   bool
	expr  ast.Expr
	ids   map[string]bool // identity variables in scope where the condition is
	inert bool            // from a branch that can't stop or hand on
}

type ppath struct {
	atoms     []atom
	handed    bool
	aborted   bool
	strong    bool   // passed a rejection: a 401, an unauthorized error, a login redirect
	forbid    bool   // wrote a 403 / forbidden error
	brk       bool   // left its block through break/continue/goto
	fall      bool   // ended its case with fallthrough
	label     string // goto target not reached yet
	stopped   bool
	unmatched bool                // fell out of a switch with no default
	delegated bool                // the request went to another middleware ...
	deleg     verdict             // ... whose own verdict this is
	opaque    bool                // went through code that can't be read
	vars      map[string]verdict  // handlers built here: h = checkToken(next)
	vals      map[string]ast.Expr // the constant a reassigned local holds on this path
}

func (p *ppath) clone() *ppath {
	q := *p
	q.atoms = append([]atom(nil), p.atoms...)
	if p.vars != nil {
		q.vars = make(map[string]verdict, len(p.vars))
		for k, v := range p.vars {
			q.vars[k] = v
		}
	}
	if p.vals != nil {
		q.vals = make(map[string]ast.Expr, len(p.vals))
		for k, v := range p.vals {
			q.vals[k] = v
		}
	}
	return &q
}

type pathResult struct {
	scopes   int
	paths    []*ppath // ended paths
	opaque   bool
	overflow bool // too many paths to read even after merging
	why      string
}

type pathWalker struct {
	w     *walker // resolver for the scope's file
	base  *walker
	depth int
	gin   bool
	ctx   string // the gin context parameter's name

	recv      string
	declared  map[string]bool
	cfg       map[string]bool
	ids       map[string]bool
	req       map[string]bool         // the scope's own parameters: w, r, c
	locals    map[string]*ast.FuncLit // deny := func() { ... }
	aborters  map[string]bool         // abort := c.AbortWithStatus
	okVars    map[string]bool         // s, ok := store[key]: ok says a credential was found
	derived   map[string]bool         // names computed from the request
	setting   map[string]bool         // locals assigned once from a setting
	bypassVar map[string]bool         // locals assigned once from a request-parameter test
	tracked   map[string]bool         // locals assigned more than once and tested: followed per path
	decls     map[string][]idDecl     // every assignment of a name, in order
	res       *pathResult
	count     int
}

const (
	maxPaths  = 4096 // paths created in one scope before it counts as too branchy
	compactAt = 24   // live paths before irrelevant branches are merged
)

// ownHandOff reports whether the scope hands the request on at its own level
// (not inside a nested closure).
func ownHandOff(root *ast.BlockStmt, recv string) bool {
	hit := false
	ast.Inspect(root, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.DeferStmt, *ast.GoStmt:
			// go func() { next.ServeHTTP(w, r) }(): the hand-off runs anyway
			var c *ast.CallExpr
			if d, ok := x.(*ast.DeferStmt); ok {
				c = d.Call
			} else {
				c = x.(*ast.GoStmt).Call
			}
			if handOff(c, recv) {
				hit = true
			} else if fl, ok := unparen(c.Fun).(*ast.FuncLit); ok && ownHandOff(fl.Body, recv) {
				hit = true
			}
			return false
		case *ast.CallExpr:
			if handOff(x, recv) {
				hit = true
			}
			// m.ServeHTTP(w, r, next): handing the next handler to a delegate
			for _, a := range x.Args {
				if handOffValue(a) {
					hit = true
				}
			}
		}
		return !hit
	})
	return hit
}

// handOff refines isHandOff: ServeHTTP hands on only with (w, r), and never
// on the scope's own receiver (that's the same middleware, not the next one).
func handOff(c *ast.CallExpr, recv string) bool {
	if sel, ok := unparen(c.Fun).(*ast.SelectorExpr); ok && sel.Sel.Name == "ServeHTTP" {
		if len(c.Args) != 2 {
			return false
		}
		if id, ok := unparen(sel.X).(*ast.Ident); ok && recv != "" && id.Name == recv {
			return false
		}
		return true
	}
	return isHandOff(c)
}

// handOffValue: an argument that passes the next handler along without
// calling it (negroni's m.ServeHTTP(w, r, next)).
func handOffValue(x ast.Expr) bool {
	switch e := unparen(x).(type) {
	case *ast.Ident:
		return nextName(e.Name)
	case *ast.SelectorExpr:
		return e.Sel.Name == "ServeHTTP" || e.Sel.Name == "Next"
	}
	return false
}

// ginScope: the scope's context runs the next handler unless aborted.
func (w *walker) ginScope(ft *ast.FuncType) bool {
	if ft == nil || ft.Params == nil {
		return false
	}
	for _, p := range ft.Params.List {
		t := p.Type
		if s, ok := t.(*ast.StarExpr); ok {
			t = s.X
		}
		sel, ok := t.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		id, ok := sel.X.(*ast.Ident)
		if !ok {
			continue
		}
		ip := w.fc.imports[id.Name]
		switch {
		case sel.Sel.Name == "Context" && strings.HasSuffix(ip, "gin-gonic/gin"),
			sel.Sel.Name == "RequestContext" && strings.HasSuffix(ip, "cloudwego/hertz/pkg/app"):
			return true
		}
	}
	return false
}

// ginParam names the scope's gin (or hertz) context parameter.
func (w *walker) ginParam(ft *ast.FuncType) string {
	if ft == nil || ft.Params == nil {
		return ""
	}
	for _, p := range ft.Params.List {
		if len(p.Names) == 0 {
			continue
		}
		if w.ginScope(&ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{p}}}) {
			return p.Names[0].Name
		}
	}
	return ""
}

// declaredNames collects every name a function declares: parameters,
// results, := and var, so a package-level flag isn't confused with a local.
func declaredNames(ft *ast.FuncType, body *ast.BlockStmt) map[string]bool {
	out := map[string]bool{}
	addFields := func(fl *ast.FieldList) {
		if fl == nil {
			return
		}
		for _, f := range fl.List {
			for _, n := range f.Names {
				out[n.Name] = true
			}
		}
	}
	if ft != nil {
		addFields(ft.Params)
		addFields(ft.Results)
	}
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			addFields(x.Type.Params)
			addFields(x.Type.Results)
		case *ast.AssignStmt:
			if x.Tok == token.DEFINE {
				for _, l := range x.Lhs {
					if id, ok := l.(*ast.Ident); ok {
						out[id.Name] = true
					}
				}
			}
		case *ast.ValueSpec:
			for _, id := range x.Names {
				out[id.Name] = true
			}
		case *ast.RangeStmt:
			for _, e := range []ast.Expr{x.Key, x.Value} {
				if id, ok := e.(*ast.Ident); ok {
					out[id.Name] = true
				}
			}
		}
		return true
	})
	return out
}

type idDecl struct {
	pos token.Pos
	id  bool // assigned from an identity lookup or auth check
}

// scopeDecls records each assignment of a local, so a name reused for
// something else (an ok from a context lookup, then an ok from BasicAuth)
// is judged by the assignment it refers to.
func scopeDecls(root *ast.BlockStmt, params map[string]bool) map[string][]idDecl {
	out := map[string][]idDecl{}
	ast.Inspect(root, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		isID := false
		if len(as.Rhs) == 1 {
			if call, ok := unparen(as.Rhs[0]).(*ast.CallExpr); ok {
				fn, isIdent := unparen(call.Fun).(*ast.Ident)
				name := exprName(call.Fun)
				isID = identityCall(name) || authCheckCall(name) || isIdent && params[fn.Name] || credentialRead(call)
			}
		}
		for _, l := range as.Lhs {
			if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
				out[id.Name] = append(out[id.Name], idDecl{as.Pos(), isID})
			}
		}
		return true
	})
	return out
}

// idsAt: the identity variables as seen at pos: each name by its latest
// assignment before pos.
func (pw *pathWalker) idsAt(pos token.Pos) map[string]bool {
	if pw.decls == nil || !pos.IsValid() {
		return pw.ids
	}
	out := map[string]bool{}
	for name, ds := range pw.decls {
		last := -1
		for i, d := range ds {
			if d.pos < pos {
				last = i
			}
		}
		if last >= 0 && ds[last].id {
			out[name] = true
		}
	}
	return out
}

func (pw *pathWalker) atomIDs(a atom) map[string]bool {
	if a.ids != nil {
		return a.ids
	}
	return pw.ids
}

// scopeIDs: variables holding the caller's identity or the result of an
// auth check: s := Auth(c, ...), u, err := getUserSession(r),
// err = auth.CheckRequest(c, ...).
func scopeIDs(root *ast.BlockStmt, params map[string]bool) map[string]bool {
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
		name := exprName(call.Fun)
		if !identityCall(name) && !authCheckCall(name) && !(isID && params[fn.Name]) && !credentialRead(call) {
			return true
		}
		for _, l := range as.Lhs {
			if id, ok := l.(*ast.Ident); ok && id.Name != "_" {
				ids[id.Name] = true
			}
		}
		return true
	})
	return ids
}

// authCheckCall: a call whose qualified name says it checks credentials:
// auth.CheckRequest, jwt.Parse, s.VerifyToken, signer.Validate.
func authCheckCall(name string) bool {
	var ws []string
	for _, part := range strings.Split(name, ".") {
		ws = append(ws, camelWords(part)...)
	}
	check, subject := false, false
	for _, w := range ws {
		switch w {
		case "authenticate", "authorize", "authorise":
			check, subject = true, true
		case "check", "verify", "validate", "parse", "decode":
			check = true
		case "auth", "jwt", "token", "sign", "signature", "signed", "credential", "credentials",
			"session", "apikey", "hmac", "bearer", "oauth", "oidc":
			subject = true
		}
	}
	return check && subject
}

func camelWords(s string) []string {
	var out []string
	start := 0
	for i := 1; i <= len(s); i++ {
		if i == len(s) || s[i] == '_' || (s[i] >= 'A' && s[i] <= 'Z' && !(s[i-1] >= 'A' && s[i-1] <= 'Z')) {
			if w := strings.Trim(strings.ToLower(s[start:i]), "_"); w != "" {
				out = append(out, w)
			}
			start = i
		}
	}
	return out
}

func (pw *pathWalker) run(root *ast.BlockStmt) {
	pw.res.scopes++
	live := pw.exec(root.List, []*ppath{{}})
	for _, p := range live {
		pw.end(p)
	}
}

func (pw *pathWalker) end(p *ppath) {
	if pw.gin {
		p.stopped = p.aborted
	} else {
		p.stopped = !p.handed || p.aborted
	}
	pw.res.paths = append(pw.res.paths, p)
}

func (pw *pathWalker) grow(n int) bool {
	pw.count += n
	if pw.count > maxPaths {
		pw.res.opaque = true
		pw.res.overflow = true
		return false
	}
	return true
}

// exec runs a statement list over the live paths and returns the ones that
// reach its end (or leave it through break).
func (pw *pathWalker) exec(list []ast.Stmt, in []*ppath) []*ppath {
	live := in
	pending := map[int][]*ppath{}
	for i, st := range list {
		live = append(live, pending[i]...)
		var next []*ppath
		for _, p := range live {
			if p.brk {
				next = append(next, p)
				continue
			}
			next = append(next, pw.stmt(st, p)...)
		}
		live = live[:0]
		for _, p := range next {
			// goto forward to a label in this block
			if p.brk && p.label != "" {
				if j := labelIndex(list, p.label); j > i {
					p.brk, p.label = false, ""
					pending[j] = append(pending[j], p)
					continue
				}
			}
			live = append(live, p)
		}
		if len(live) > compactAt {
			live = pw.compact(live)
		}
		if pw.res.opaque {
			return nil
		}
	}
	return live
}

func labelIndex(list []ast.Stmt, name string) int {
	for i, st := range list {
		if l, ok := st.(*ast.LabeledStmt); ok && l.Label.Name == name {
			return i
		}
	}
	return -1
}

// compact merges live paths that differ only in conditions that have
// nothing to do with auth (a log level, a header that isn't a credential):
// ten harmless ifs are 1024 paths, but one decision.
func (pw *pathWalker) compact(live []*ppath) []*ppath {
	seen := map[string]bool{}
	var out []*ppath
	for _, p := range live {
		var keep []atom
		var b strings.Builder
		for _, a := range p.atoms {
			if pw.relevant(a) {
				keep = append(keep, a)
				fmt.Fprintf(&b, "%s=%v;", a.key, a.pos)
			}
		}
		fmt.Fprintf(&b, "|%v%v%v%v%v%v%d%s", p.handed, p.aborted, p.strong, p.brk, p.fall, p.opaque, p.deleg, p.label)
		if p.delegated {
			b.WriteString("d")
		}
		vk := make([]string, 0, len(p.vars))
		for k, v := range p.vars {
			vk = append(vk, fmt.Sprintf("%s:%d", k, v))
		}
		for k, v := range p.vals {
			vk = append(vk, k+"="+types.ExprString(v))
		}
		sort.Strings(vk)
		b.WriteString(strings.Join(vk, ","))
		k := b.String()
		if seen[k] {
			continue
		}
		seen[k] = true
		p.atoms = keep
		out = append(out, p)
	}
	return out
}

// relevant: a condition that can matter to the verdict.
func (pw *pathWalker) relevant(a atom) bool {
	if a.inert {
		return false
	}
	if a.kind != akRequest || a.expr == nil {
		return true
	}
	e := litExpr(a)
	if authnCond(e, pw.atomIDs(a)) || credentialRead(a.expr) || mentionsOptions(a.expr) {
		return true
	}
	hit := false
	ast.Inspect(a.expr, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && (pw.atomIDs(a)[id.Name] || pw.okVars[id.Name]) {
			hit = true
		}
		return !hit
	})
	return hit
}

func (pw *pathWalker) stmt(st ast.Stmt, p *ppath) []*ppath {
	switch s := st.(type) {
	case *ast.ReturnStmt:
		for _, r := range s.Results {
			pw.effects(r, p)
		}
		pw.end(p)
		return nil
	case *ast.ExprStmt:
		if pw.effects(s.X, p) {
			pw.end(p)
			return nil
		}
		return []*ppath{p}
	case *ast.IfStmt:
		if s.Init != nil {
			out := pw.stmt(s.Init, p)
			if len(out) == 0 {
				return nil
			}
			p = out[0]
		}
		pw.effects(s.Cond, p)
		// `if r.URL.Query().Get("v") == "1" { log() }`: a branch that can't
		// stop or hand on can't make anything conditional, unless it sets a
		// local a later check reads (`ok = true`, `tok = "k"`)
		inert := s.Else == nil && !controls(s.Body) && !pw.setsTracked(s.Body)
		cond := pw.subst(s.Cond, p)
		var out []*ppath
		for _, conj := range pw.dnf(cond, true) {
			if inert {
				conj = quiet(conj)
			}
			q := p.clone()
			q.atoms = append(q.atoms, conj...)
			if !pw.grow(1) {
				return nil
			}
			out = append(out, pw.exec(s.Body.List, []*ppath{q})...)
		}
		for _, conj := range pw.dnf(cond, false) {
			if inert {
				conj = quiet(conj)
			}
			q := p.clone()
			q.atoms = append(q.atoms, conj...)
			if !pw.grow(1) {
				return nil
			}
			if s.Else != nil {
				out = append(out, pw.stmt(s.Else, q)...)
			} else {
				out = append(out, q)
			}
		}
		return out
	case *ast.BlockStmt:
		return pw.exec(s.List, []*ppath{p})
	case *ast.ForStmt, *ast.RangeStmt:
		var body *ast.BlockStmt
		if f, ok := s.(*ast.ForStmt); ok {
			if f.Init != nil {
				pw.effects(f.Init, p)
			}
			body = f.Body
		} else {
			pw.effects(s.(*ast.RangeStmt).X, p)
			body = s.(*ast.RangeStmt).Body
		}
		if !pw.grow(1) {
			return nil
		}
		once := pw.exec(body.List, []*ppath{p.clone()})
		for _, q := range once {
			if q.label == "" {
				q.brk = false
			}
		}
		return append([]*ppath{p}, once...)
	case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		var body *ast.BlockStmt
		switch b := s.(type) {
		case *ast.SwitchStmt:
			if b.Init != nil {
				pw.effects(b.Init, p)
			}
			if b.Tag != nil {
				pw.effects(b.Tag, p)
			}
			body = b.Body
		case *ast.TypeSwitchStmt:
			body = b.Body
		case *ast.SelectStmt:
			body = b.Body
		}
		var out, carry []*ppath
		hasDefault := false
		for i, c := range body.List {
			var stmts []ast.Stmt
			var conds []ast.Expr
			switch cc := c.(type) {
			case *ast.CaseClause:
				stmts, conds = cc.Body, cc.List
				hasDefault = hasDefault || cc.List == nil
			case *ast.CommClause:
				stmts = cc.Body
				hasDefault = hasDefault || cc.Comm == nil
			}
			q := p.clone()
			if len(conds) > 0 {
				key := "case@" + pw.w.site(c.Pos())
				q.atoms = append(q.atoms, atom{key: key, kind: pw.switchKind(s, conds), pos: true})
			}
			_ = i
			if !pw.grow(1) {
				return nil
			}
			ins := []*ppath{q}
			for _, f := range carry {
				f.fall, f.brk = false, false
				ins = append(ins, f)
			}
			carry = nil
			for _, r := range pw.exec(stmts, ins) {
				if r.fall {
					carry = append(carry, r) // into the next case's body
					continue
				}
				if r.label == "" {
					r.brk = false
				}
				out = append(out, r)
			}
		}
		out = append(out, carry...)
		if !hasDefault {
			p.unmatched = true
			out = append(out, p)
		}
		return out
	case *ast.LabeledStmt:
		return pw.stmt(s.Stmt, p)
	case *ast.BranchStmt:
		p.brk = true
		switch s.Tok {
		case token.FALLTHROUGH:
			p.fall = true
		case token.GOTO:
			if s.Label != nil {
				p.label = s.Label.Name
			}
		}
		return []*ppath{p}
	case *ast.DeferStmt, *ast.GoStmt:
		// runs after the handler or beside it: it can't stop the request,
		// but a hand-off there (go next.ServeHTTP(w, r)) always happens
		var c *ast.CallExpr
		if d, ok := s.(*ast.DeferStmt); ok {
			c = d.Call
		} else {
			c = s.(*ast.GoStmt).Call
		}
		if handOff(c, pw.recv) {
			p.handed = true
		} else if fl, ok := unparen(c.Fun).(*ast.FuncLit); ok && ownHandOff(fl.Body, pw.recv) {
			p.handed = true
		}
		return []*ppath{p}
	case *ast.AssignStmt:
		for i, l := range s.Lhs {
			if id, ok := l.(*ast.Ident); ok && pw.tracked[id.Name] {
				var v ast.Expr
				if len(s.Rhs) == len(s.Lhs) && (s.Tok == token.ASSIGN || s.Tok == token.DEFINE) {
					v = pw.constVal(s.Rhs[i])
				}
				p.setVal(id.Name, v)
			}
		}
		for i, r := range s.Rhs {
			// h = checkToken(next): a handler built here, called later
			if c, ok := unparen(r).(*ast.CallExpr); ok && i < len(s.Lhs) && pw.constructor(c) {
				if id, ok := s.Lhs[i].(*ast.Ident); ok && id.Name != "_" {
					if p.vars == nil {
						p.vars = map[string]verdict{}
					}
					p.vars[id.Name] = pw.delegVerdict(c)
				}
				continue
			}
			pw.effects(r, p)
		}
		return []*ppath{p}
	case *ast.DeclStmt:
		if g, ok := s.Decl.(*ast.GenDecl); ok && g.Tok == token.VAR {
			for _, sp := range g.Specs {
				vs := sp.(*ast.ValueSpec)
				for i, n := range vs.Names {
					if !pw.tracked[n.Name] {
						continue
					}
					var v ast.Expr
					if len(vs.Values) == len(vs.Names) {
						v = pw.constVal(vs.Values[i])
					} else if len(vs.Values) == 0 {
						v = zeroVal(vs.Type, n.Pos()) // var ok bool
					}
					p.setVal(n.Name, v)
				}
			}
		}
		pw.effects(s, p)
		return []*ppath{p}
	case *ast.IncDecStmt:
		if id, ok := s.X.(*ast.Ident); ok && pw.tracked[id.Name] {
			p.setVal(id.Name, nil)
		}
	}
	pw.effects(st, p)
	return []*ppath{p}
}

func (p *ppath) setVal(name string, v ast.Expr) {
	if v == nil {
		delete(p.vals, name)
		return
	}
	if p.vals == nil {
		p.vals = map[string]ast.Expr{}
	}
	p.vals[name] = v
}

// constVal returns x if a later condition can use it in place of the local:
// a constant it can fold or read as a setting (true, "k", nil, a package
// constant), or a credential read. nil means the value is unknown.
func (pw *pathWalker) constVal(x ast.Expr) ast.Expr {
	if credentialRead(x) {
		return x // tok := r.Header.Get("Authorization")
	}
	switch e := unparen(x).(type) {
	case *ast.BasicLit:
		return e
	case *ast.Ident:
		if e.Name == "true" || e.Name == "false" || e.Name == "nil" {
			return e
		}
		if _, isConst := pw.w.fc.pkg.consts[e.Name]; isConst && !pw.declared[e.Name] {
			return e
		}
	}
	return nil
}

func zeroVal(t ast.Expr, pos token.Pos) ast.Expr {
	id, ok := t.(*ast.Ident)
	if !ok {
		return nil
	}
	switch id.Name {
	case "bool":
		return &ast.Ident{NamePos: pos, Name: "false"}
	case "string":
		return &ast.BasicLit{ValuePos: pos, Kind: token.STRING, Value: `""`}
	}
	return nil
}

// setsTracked: the block assigns a local that's followed per path.
func (pw *pathWalker) setsTracked(b *ast.BlockStmt) bool {
	hit := false
	ast.Inspect(b, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.AssignStmt:
			for _, l := range x.Lhs {
				if id, ok := l.(*ast.Ident); ok && pw.tracked[id.Name] {
					hit = true
				}
			}
		}
		return !hit
	})
	return hit
}

// subst puts the path's known value in place of each tracked local the
// condition reads: with ok = true set on this path, `!ok && !valid(tok)`
// is false here; with tok = "k", `!valid(tok)` checks a constant.
func (pw *pathWalker) subst(x ast.Expr, p *ppath) ast.Expr {
	if len(p.vals) == 0 {
		return x
	}
	var walk func(ast.Expr) ast.Expr
	walk = func(x ast.Expr) ast.Expr {
		switch e := x.(type) {
		case *ast.Ident:
			if v, ok := p.vals[e.Name]; ok {
				return v
			}
		case *ast.ParenExpr:
			if y := walk(e.X); y != e.X {
				c := *e
				c.X = y
				return &c
			}
		case *ast.UnaryExpr:
			if y := walk(e.X); y != e.X {
				c := *e
				c.X = y
				return &c
			}
		case *ast.BinaryExpr:
			l, r := walk(e.X), walk(e.Y)
			if l != e.X || r != e.Y {
				c := *e
				c.X, c.Y = l, r
				return &c
			}
		case *ast.CallExpr:
			var args []ast.Expr
			for i, a := range e.Args {
				if y := walk(a); y != a {
					if args == nil {
						args = append([]ast.Expr(nil), e.Args...)
					}
					args[i] = y
				}
			}
			if args != nil {
				c := *e
				c.Args = args
				return &c
			}
		}
		return x
	}
	return walk(x)
}

// controls: the block can end the request or pass it on: a return, a
// branch, a hand-off, an abort, a panic.
func controls(b *ast.BlockStmt) bool {
	hit := false
	ast.Inspect(b, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.ReturnStmt, *ast.BranchStmt:
			hit = true
		case *ast.CallExpr:
			name := exprName(x.Fun)
			if isHandOff(x) || strings.Contains(name, "Abort") || name == "panic" || strings.Contains(name, "Fatal") {
				hit = true
			}
			for _, a := range x.Args {
				if handOffValue(a) {
					hit = true
				}
			}
		}
		return !hit
	})
	return hit
}

// quiet demotes a bypass literal to a plain request condition.
func quiet(conj []atom) []atom {
	out := append([]atom(nil), conj...)
	for i := range out {
		out[i].kind = akRequest
		out[i].inert = true
	}
	return out
}

// switchKind: a switch on the request method is a method condition.
func (pw *pathWalker) switchKind(s ast.Stmt, conds []ast.Expr) atomKind {
	sw, ok := s.(*ast.SwitchStmt)
	if !ok || sw.Tag == nil || !methodExpr(sw.Tag) {
		return akRequest
	}
	for _, c := range conds {
		if mentionsOptions(c) {
			return akRequest
		}
	}
	return akMethod
}

// effects records what evaluating x does to the path: hands off, aborts,
// passes a rejection mark, delegates to another middleware. It reports
// whether x ends the request outright (panic, os.Exit).
func (pw *pathWalker) effects(x ast.Node, p *ppath) bool {
	terminal := false
	ast.Inspect(x, func(n ast.Node) bool {
		switch c := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.CallExpr:
			if v, ok := pw.handVia(c, p); ok {
				if !p.handed {
					p.handed = true
					if v != vPasses {
						p.delegated, p.deleg = true, v
						p.opaque = p.opaque || v == vUnknown
					}
				}
				return false
			}
			if pw.constructor(c) {
				return false // built and dropped: it never sees the request
			}
			name := exprName(c.Fun)
			short := name
			if i := strings.LastIndexByte(short, '.'); i >= 0 {
				short = short[i+1:]
			}
			switch {
			case name == "panic", name == "os.Exit", name == "runtime.Goexit",
				strings.HasPrefix(name, "log.Fatal"), strings.HasPrefix(name, "log.Panic"):
				terminal = true
			case strings.HasPrefix(short, "Abort") && strings.Contains(name, ".") && !p.handed:
				// aborting after the handler ran (error mapping) stops nothing
				p.aborted = true
			}
			if id, ok := unparen(c.Fun).(*ast.Ident); ok {
				// abort := c.AbortWithStatus; abort(401)
				if pw.aborters[id.Name] && !p.handed {
					p.aborted = true
				}
				// deny := func() { http.Error(w, "", 401) }; deny()
				if fl := pw.locals[id.Name]; fl != nil {
					pw.closureEffects(fl, p)
				}
			}
			// gin: a helper handed the context can abort the chain
			// (gotify's a.evaluateOr401(ctx, ...))
			if pw.gin && !p.handed && pw.passesCtx(c) {
				if fi := pw.w.resolveCall(c); fi != nil {
					aborts, rejects := pw.base.abortsIn(fi, c, pw.w, 0)
					if aborts {
						p.aborted = true
						p.strong = p.strong || rejects
					}
				} else if !pw.ctxMethod(c) {
					p.opaque = true // an outside helper may abort
				}
			}
			if pw.depth < 2 {
				if fi := pw.w.resolveCall(c); fi != nil && pw.base.helperRejects(fi, pw.depth+1) {
					p.strong = true
				}
			}
		}
		if rejectionMark(n) {
			p.strong = true
		}
		if forbiddenMark(n) {
			p.forbid = true
		}
		return true
	})
	return terminal
}

// handVia reports a call that passes the request on, and how: plainly
// (vPasses), or through another middleware whose verdict it returns:
// h.ServeHTTP after h = checkToken(next), checkToken(next).ServeHTTP(w, r),
// checkToken(next)(c), or m.ServeHTTP(w, r, next).
func (pw *pathWalker) handVia(c *ast.CallExpr, p *ppath) (verdict, bool) {
	if handOff(c, pw.recv) {
		switch f := unparen(c.Fun).(type) {
		case *ast.Ident:
			if v, ok := p.vars[f.Name]; ok {
				return v, true
			}
		case *ast.SelectorExpr:
			if id, ok := unparen(f.X).(*ast.Ident); ok {
				if v, ok := p.vars[id.Name]; ok {
					return v, true
				}
			}
			if inner, ok := unparen(f.X).(*ast.CallExpr); ok && pw.constructor(inner) {
				return pw.delegVerdict(inner), true
			}
		}
		return vPasses, true
	}
	if id, ok := unparen(c.Fun).(*ast.Ident); ok {
		if v, ok := p.vars[id.Name]; ok {
			return v, true
		}
	}
	// checkToken(next)(c), or a factory's handler run on this request:
	// SignRequired(auth)(c)
	if inner, ok := unparen(c.Fun).(*ast.CallExpr); ok && (pw.constructor(inner) || pw.passesRequest(c)) {
		return pw.delegVerdict(inner), true
	}
	if pw.takesNext(c) && pw.passesRequest(c) {
		return pw.delegVerdict(c), true
	}
	return 0, false
}

func (pw *pathWalker) takesNext(c *ast.CallExpr) bool {
	for _, a := range c.Args {
		if handOffValue(a) {
			return true
		}
	}
	return false
}

// passesRequest: the call gets the scope's own request parameters, so it
// handles this request rather than building a handler.
func (pw *pathWalker) passesRequest(c *ast.CallExpr) bool {
	for _, a := range c.Args {
		if id, ok := unparen(a).(*ast.Ident); ok && pw.req[id.Name] {
			return true
		}
	}
	return false
}

// constructor reports a call like checkToken(next) that builds a handler
// around next rather than handling the request.
func (pw *pathWalker) constructor(c *ast.CallExpr) bool {
	return pw.takesNext(c) && !pw.passesRequest(c) && !handOff(c, pw.recv)
}

func (pw *pathWalker) delegVerdict(c *ast.CallExpr) verdict {
	fi := pw.w.resolveCall(c)
	if fi == nil {
		return vUnknown
	}
	return pw.base.funcVerdict(fi)
}

// closureEffects applies what a local closure does unconditionally: its
// top-level statements, not those under its own ifs.
func (pw *pathWalker) closureEffects(fl *ast.FuncLit, p *ppath) {
	for _, st := range fl.Body.List {
		switch st.(type) {
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SelectStmt,
			*ast.GoStmt, *ast.DeferStmt:
			continue
		}
		ast.Inspect(st, func(n ast.Node) bool {
			switch c := n.(type) {
			case *ast.FuncLit:
				return false
			case *ast.CallExpr:
				name := exprName(c.Fun)
				if i := strings.LastIndexByte(name, '.'); i >= 0 && strings.HasPrefix(name[i+1:], "Abort") && !p.handed {
					p.aborted = true
				}
				if handOff(c, pw.recv) {
					p.handed = true
				}
			}
			if rejectionMark(n) {
				p.strong = true
			}
			return true
		})
	}
}

// passesCtx: the call takes the gin context itself as an argument.
func (pw *pathWalker) passesCtx(c *ast.CallExpr) bool {
	if pw.ctx == "" {
		return false
	}
	for _, a := range c.Args {
		if id, ok := unparen(a).(*ast.Ident); ok && id.Name == pw.ctx {
			return true
		}
	}
	return false
}

// ctxMethod: c.JSON(...), c.Set(...): a call on the context, not a helper.
func (pw *pathWalker) ctxMethod(c *ast.CallExpr) bool {
	sel, ok := unparen(c.Fun).(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id := rootIdent(sel.X)
	return id != nil && id.Name == pw.ctx
}

// abortsIn reports whether a function (or what it calls, a few levels down)
// aborts a gin chain, and whether it does so with a rejection.
func (w *walker) abortsIn(fi *funcInfo, call *ast.CallExpr, caller *walker, depth int) (aborts, rejects bool) {
	if fi == nil || fi.decl.Body == nil || depth > 3 {
		return false, false
	}
	r := w.with(fi.fc)
	if fi.recvType != "" && len(fi.decl.Recv.List[0].Names) > 0 {
		rn := fi.decl.Recv.List[0].Names[0].Name
		r.env.locals[rn] = true
		r.env.vtypes[rn] = typeRef{pkg: fi.fc.pkg, name: fi.recvType}
	}
	// respond(c, 401, ...) into `if code >= 500 { c.Abort... }`: bind the
	// constant arguments so a branch the call can't take is skipped
	if call != nil && caller != nil {
		i := 0
		for _, fld := range fi.decl.Type.Params.List {
			for _, nm := range fld.Names {
				if i < len(call.Args) {
					if v, ok := caller.constInt(call.Args[i]); ok {
						r.env.ints[nm.Name] = v
					}
					if b, ok := caller.constBool(call.Args[i]); ok {
						r.env.bools[nm.Name] = b
					}
				}
				i++
			}
		}
	}
	var visit func(n ast.Node) bool
	visit = func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.IfStmt:
			if b, ok := r.constBool(x.Cond); ok {
				if x.Init != nil {
					ast.Inspect(x.Init, visit)
				}
				if b {
					ast.Inspect(x.Body, visit)
				} else if x.Else != nil {
					ast.Inspect(x.Else, visit)
				}
				return false
			}
		case *ast.CallExpr:
			name := exprName(x.Fun)
			if i := strings.LastIndexByte(name, '.'); i >= 0 && strings.HasPrefix(name[i+1:], "Abort") {
				aborts = true
				for _, a := range x.Args {
					ast.Inspect(a, func(m ast.Node) bool {
						if rejectionMark(m) || forbiddenMark(m) {
							rejects = true
						}
						return !rejects
					})
				}
				if rejectionMark(x) {
					rejects = true
				}
				return true
			}
			if callee := r.resolveCall(x); callee != nil && callee != fi {
				a, rj := w.abortsIn(callee, x, r, depth+1)
				aborts, rejects = aborts || a, rejects || rj
			}
		}
		return true
	}
	ast.Inspect(fi.decl.Body, visit)
	return aborts, rejects
}

// dnf turns a condition into the conjunctions of literals under which it
// is `want`; nil means it can't be. && and || short-circuit as in Go.
func (pw *pathWalker) dnf(x ast.Expr, want bool) [][]atom {
	if b, ok := pw.constOf(x); ok {
		if b == want {
			return [][]atom{{}}
		}
		return nil
	}
	switch e := unparen(x).(type) {
	case *ast.UnaryExpr:
		if e.Op == token.NOT {
			return pw.dnf(e.X, !want)
		}
	case *ast.BinaryExpr:
		var out [][]atom
		switch {
		case e.Op == token.LAND && want:
			out = cross(pw.dnf(e.X, true), pw.dnf(e.Y, true))
		case e.Op == token.LAND:
			out = append(pw.dnf(e.X, false), cross(pw.dnf(e.X, true), pw.dnf(e.Y, false))...)
		case e.Op == token.LOR && want:
			out = append(pw.dnf(e.X, true), cross(pw.dnf(e.X, false), pw.dnf(e.Y, true))...)
		case e.Op == token.LOR:
			out = cross(pw.dnf(e.X, false), pw.dnf(e.Y, false))
		default:
			return [][]atom{{pw.atomOf(e, want)}}
		}
		// a condition that expands to more than 16 cases is kept whole
		if len(out) > 16 {
			return [][]atom{{pw.atomOf(e, want)}}
		}
		return out
	}
	return [][]atom{{pw.atomOf(x, want)}}
}

func cross(a, b [][]atom) [][]atom {
	var out [][]atom
	for _, x := range a {
		for _, y := range b {
			c := append(append([]atom(nil), x...), y...)
			out = append(out, c)
		}
	}
	return out
}

// constOf folds a condition: literals, bound parameters, and package
// constants (const enforce = false).
func (pw *pathWalker) constOf(x ast.Expr) (bool, bool) {
	if b, ok := pw.w.constBool(x); ok {
		return b, true
	}
	if id, ok := unparen(x).(*ast.Ident); ok && !pw.declared[id.Name] {
		p := pw.w.fc.pkg
		if c, ok := p.consts[id.Name]; ok && p.globals[id.Name] == nil {
			if lit, ok := unparen(c.expr).(*ast.Ident); ok && (lit.Name == "true" || lit.Name == "false") {
				return lit.Name == "true", true
			}
		}
	}
	return false, false
}

// atomOf normalises a != b to the negation of a == b, so both sides of a
// comparison share a key.
func (pw *pathWalker) atomOf(x ast.Expr, pos bool) atom {
	x = unparen(x)
	norm := x
	if b, ok := x.(*ast.BinaryExpr); ok {
		// len(tok) > 0 is tok != ""
		if arg, cmp, ok := lenCompare(b); ok {
			norm = &ast.BinaryExpr{X: arg, Op: token.EQL, Y: &ast.BasicLit{Kind: token.STRING, Value: `""`}}
			if cmp {
				pos = !pos
			}
			return atom{key: pw.keyOf(norm, x.Pos()), kind: pw.kindOf(norm), pos: pos, expr: norm, ids: pw.idsAt(x.Pos())}
		}
		if b.Op == token.NEQ {
			norm = &ast.BinaryExpr{X: b.X, Op: token.EQL, Y: b.Y}
			pos = !pos
		}
	}
	return atom{key: pw.keyOf(norm, x.Pos()), kind: pw.kindOf(norm), pos: pos, expr: norm, ids: pw.idsAt(x.Pos())}
}

// keyOf names a condition by its text and, for each local it reads, the
// assignment it refers to: `!ok` after a context lookup and `!ok` after
// BasicAuth are different conditions.
func (pw *pathWalker) keyOf(x ast.Expr, at token.Pos) string {
	key := types.ExprString(x)
	ast.Inspect(x, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		ds := pw.decls[id.Name]
		last := token.NoPos
		for _, d := range ds {
			if d.pos < at {
				last = d.pos
			}
		}
		if last.IsValid() {
			key += fmt.Sprintf("#%s@%d", id.Name, last)
		}
		return true
	})
	return key
}

// lenCompare: len(x) > 0, len(x) != 0, len(x) >= 1 (nonEmpty true) or
// len(x) == 0, len(x) < 1 (nonEmpty false).
func lenCompare(b *ast.BinaryExpr) (arg ast.Expr, nonEmpty, ok bool) {
	c, isCall := unparen(b.X).(*ast.CallExpr)
	lit, isLit := unparen(b.Y).(*ast.BasicLit)
	if !isCall || !isLit || lit.Kind != token.INT || len(c.Args) != 1 {
		return nil, false, false
	}
	if id, ok := c.Fun.(*ast.Ident); !ok || id.Name != "len" {
		return nil, false, false
	}
	switch {
	case (b.Op == token.GTR || b.Op == token.NEQ) && lit.Value == "0", b.Op == token.GEQ && lit.Value == "1":
		return c.Args[0], true, true
	case b.Op == token.EQL && lit.Value == "0", b.Op == token.LSS && lit.Value == "1", b.Op == token.LEQ && lit.Value == "0":
		return c.Args[0], false, true
	}
	return nil, false, false
}

func (pw *pathWalker) kindOf(x ast.Expr) atomKind {
	if id, ok := unparen(x).(*ast.Ident); ok && pw.okVars[id.Name] {
		return akPresence // s, ok := sessions[key]
	}
	if requestParam(x) {
		return akBypass
	}
	if id, ok := unparen(x).(*ast.Ident); ok && pw.bypassVar[id.Name] {
		return akBypass // ok := r.URL.Query().Get("debug") == "1"
	}
	if b, ok := x.(*ast.BinaryExpr); ok && b.Op == token.EQL {
		if (methodExpr(b.X) || methodExpr(b.Y)) && !mentionsOptions(b) {
			return akMethod
		}
		for _, pair := range [][2]ast.Expr{{b.X, b.Y}, {b.Y, b.X}} {
			if !isNil(pair[1]) && !isEmptyString(pair[1]) {
				continue
			}
			if credentialRead(pair[0]) {
				return akPresence // r.Header.Get("Authorization") == ""
			}
			id := rootIdent(pair[0])
			if id == nil {
				continue
			}
			if bare, ok := unparen(pair[0]).(*ast.Ident); ok && pw.req[bare.Name] {
				continue // if c == nil: the request itself, not a credential
			}
			if isErrIdent(pair[0]) {
				if pw.idsAt(x.Pos())[id.Name] {
					return akPresence // err == nil after an identity lookup
				}
				continue
			}
			if !pw.cfgExpr(pair[0]) {
				return akPresence
			}
		}
	}
	if pw.cfgExpr(x) {
		return akConfig
	}
	return akRequest
}

// cfgExpr: the expression reads only configuration: package-level vars,
// configuration parameters, imported packages' settings, or a receiver
// field named like a switch (m.Strict, s.cfg.AuthEnabled). No calls.
func (pw *pathWalker) cfgExpr(x ast.Expr) bool {
	ok, any := true, false
	ast.Inspect(x, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.FuncLit:
			ok = false
			return false
		case *ast.CallExpr, *ast.IndexExpr:
			// os.Getenv("NOAUTH"), config.Get().AuthDisabled(), devMode():
			// a call that doesn't read the request is a setting
			if pw.readsRequest(e) {
				ok = false
			} else {
				any = true
			}
			return false
		case *ast.SelectorExpr:
			id := rootIdent(e)
			if id == nil {
				ok = false
				return false
			}
			switch {
			case pw.cfg[id.Name], pw.globalVar(id.Name):
				any = true
			case pw.isImport(id.Name):
				any = true
			case id.Name == pw.recv && pw.recv != "":
				// the middleware's own settings: m.AllowAnonymous, m.Public
				any = true
			default:
				ok = false
			}
			return false
		case *ast.Ident:
			switch {
			case e.Name == "true" || e.Name == "false" || e.Name == "nil":
			case pw.cfg[e.Name], pw.globalVar(e.Name), pw.setting[e.Name]:
				any = true
			default:
				if _, isConst := pw.w.fc.pkg.consts[e.Name]; !(isConst && !pw.declared[e.Name]) {
					ok = false
				}
			}
		}
		return ok
	})
	return ok && any
}

// readsRequest: the expression mentions the request, its parameters, or
// anything computed from them.
func (pw *pathWalker) readsRequest(x ast.Node) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && (pw.derived[id.Name] || pw.req[id.Name]) {
			hit = true
		}
		return !hit
	})
	return hit
}

func (pw *pathWalker) globalVar(name string) bool {
	return !pw.declared[name] && pw.w.fc.pkg.globals[name] != nil
}

func (pw *pathWalker) isImport(name string) bool {
	if pw.declared[name] {
		return false
	}
	_, ok := pw.w.fc.imports[name]
	return ok
}

// requestParam: a condition on a request parameter the caller chooses
// (?debug=1, X-Debug: on) compared with a fixed value, or its presence.
// The path and the method are routing, not this.
func requestParam(x ast.Expr) bool {
	read := func(e ast.Expr) bool {
		c, ok := unparen(e).(*ast.CallExpr)
		if !ok {
			return false
		}
		sel, ok := unparen(c.Fun).(*ast.SelectorExpr)
		if !ok {
			return false
		}
		switch sel.Sel.Name {
		case "Get", "Has", "Query", "QueryParam", "FormValue", "PostFormValue", "DefaultQuery", "GetHeader", "Param":
		default:
			return false
		}
		if len(c.Args) != 1 {
			return false
		}
		lit, ok := unparen(c.Args[0]).(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING || credentialName(lit.Value) {
			return false
		}
		return true
	}
	switch e := unparen(x).(type) {
	case *ast.CallExpr:
		if sel, ok := unparen(e.Fun).(*ast.SelectorExpr); ok && sel.Sel.Name == "Has" {
			return read(e)
		}
	case *ast.BinaryExpr:
		for _, pair := range [][2]ast.Expr{{e.X, e.Y}, {e.Y, e.X}} {
			if lit, ok := unparen(pair[1]).(*ast.BasicLit); ok && lit.Value != `""` && read(pair[0]) {
				return true
			}
		}
	}
	return false
}

func methodExpr(x ast.Expr) bool {
	switch e := unparen(x).(type) {
	case *ast.SelectorExpr:
		return e.Sel.Name == "Method"
	case *ast.CallExpr:
		if sel, ok := unparen(e.Fun).(*ast.SelectorExpr); ok && len(e.Args) == 0 {
			return sel.Sel.Name == "Method"
		}
	}
	return false
}

func mentionsOptions(x ast.Node) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.BasicLit:
			hit = hit || strings.Contains(strings.ToUpper(e.Value), "OPTIONS")
		case *ast.Ident:
			hit = hit || e.Name == "MethodOptions"
		}
		return !hit
	})
	return hit
}

// decide sorts the ended paths. Only stops that reject on an auth check
// decide between enforcing and conditional: a Host check or a rate limit
// that also stops requests says nothing about auth. A rejecting stop is
// conditional when a path that hands the request on differs from it only
// in removable conditions (a flag, the method, whether a credential was
// sent, a request parameter) and passes no check of its own.
func (pw *pathWalker) decide() (v verdict, why atomKind) {
	var authStops, condStops, hands []*ppath
	plainStops, opaque := 0, false
	for _, p := range pw.res.paths {
		switch {
		case p.delegated && p.deleg == vEnforces:
			authStops = append(authStops, p)
		case p.delegated && p.deleg == vConditional:
			condStops = append(condStops, p)
		case p.delegated && p.deleg == vPasses:
			hands = append(hands, p)
		case p.opaque:
			opaque = true
		case p.stopped && pw.evidence(p) && !gateOnly(p):
			authStops = append(authStops, p)
		case p.stopped && p.unmatched:
			// a value no case names: nothing the code shows rejects there
		case p.stopped:
			plainStops++
		default:
			hands = append(hands, p)
		}
	}
	// With no stop that rejects on an auth check, the other stops (an
	// OPTIONS preflight, a body limit, a rate limiter) don't make it auth:
	// it still never rejects a caller for lacking credentials.
	if len(authStops) == 0 {
		switch {
		case len(condStops) > 0:
			return vConditional, akConfig
		case opaque:
			return vUnknown, 0
		}
		_ = plainStops
		return vPasses, 0
	}
	// conditions some auth rejection depends on: passing one is a check
	guards := map[string]bool{}
	for _, p := range authStops {
		for _, a := range p.atoms {
			guards[a.key] = true
		}
	}
	kind := akRequest
	for _, p := range authStops {
		k, ok := akRequest, false
		for _, q := range hands {
			if k, ok = flipsTo(p, q, guards); ok && (k != akBypass || pw.separateCheck(p)) {
				break
			}
			ok = false
		}
		if !ok {
			return vEnforces, 0
		}
		kind = k
	}
	return vConditional, kind
}

// gateOnly: the stop depends on nothing but settings, the method or a
// request parameter (navidrome's requireQuickConnect answers 401 when the
// feature is off): a feature switch, not a check of the caller.
func gateOnly(p *ppath) bool {
	if len(p.atoms) == 0 {
		return false
	}
	for _, a := range p.atoms {
		switch a.kind {
		case akConfig, akMethod:
		default:
			// a request parameter (X-Admin: 1) is the check itself here
			return false
		}
	}
	return true
}

// separateCheck: besides the request parameter, the stop failed a check of
// a credential or identity. Without one the parameter is the check itself
// (X-Admin: 1), not a way around it.
func (pw *pathWalker) separateCheck(p *ppath) bool {
	for _, a := range p.atoms {
		if a.kind == akBypass || a.expr == nil {
			continue
		}
		if e := litExpr(a); credentialRead(a.expr) || authnCond(e, pw.atomIDs(a)) {
			return true
		}
	}
	return false
}

// evidence: the stop rejects on an auth check: a 401 or unauthorized
// error, a login redirect, a missing identity, a failed credential check.
func (pw *pathWalker) evidence(p *ppath) bool {
	return p.strong || pw.identityAbsent(p) || pw.checkedBefore(p) || p.forbid && pw.readsIdentity(p)
}

// readsIdentity: the path's conditions read the caller: a role, a user,
// a permission (role != "admin" before a 403).
func (pw *pathWalker) readsIdentity(p *ppath) bool {
	for _, a := range p.atoms {
		if a.expr != nil && a.kind == akRequest && mentionsIdentity(a.expr, pw.atomIDs(a)) {
			return true
		}
	}
	return false
}

func removable(a atom) bool {
	switch a.kind {
	case akConfig, akMethod, akBypass:
		return true
	case akPresence:
		// a stop that needs the credential to be there: tok != "", ok, or
		// err == nil after a lookup
		b, ok := a.expr.(*ast.BinaryExpr)
		if !ok {
			return a.pos
		}
		if isErrIdent(b.X) || isErrIdent(b.Y) {
			return a.pos
		}
		return !a.pos
	}
	return false
}

// litExpr rebuilds the condition a literal stands for.
func litExpr(a atom) ast.Expr {
	if a.pos {
		return a.expr
	}
	if b, ok := a.expr.(*ast.BinaryExpr); ok && b.Op == token.EQL {
		return &ast.BinaryExpr{X: b.X, Op: token.NEQ, Y: b.Y}
	}
	return &ast.UnaryExpr{Op: token.NOT, X: a.expr}
}

// flipsTo: q agrees with p on every literal p can't drop, and disagrees on
// at least one it can.
func flipsTo(p, q *ppath, guards map[string]bool) (atomKind, bool) {
	pa := map[string]atom{}
	for _, a := range p.atoms {
		pa[a.key] = a
	}
	flipped, kind := false, akRequest
	for _, a := range q.atoms {
		b, ok := pa[a.key]
		switch {
		case ok && b.pos == a.pos:
		case ok && removable(b):
			flipped, kind = true, b.kind
		case !ok && (a.kind == akConfig || a.kind == akMethod || a.kind == akBypass):
			// q tests a setting p never reached
		case !ok && !guards[a.key]:
			// q went on past where p stopped and tested something no auth
			// rejection depends on (a WebSocket upgrade, gzip)
		default:
			// q passes a check of its own (canWrite where p failed canRead)
			return 0, false
		}
	}
	return kind, flipped
}

// identityAbsent: the path went through "there is no caller" (user == nil,
// err != nil after a lookup or check, !ok).
func (pw *pathWalker) identityAbsent(p *ppath) bool {
	for _, a := range p.atoms {
		if a.expr == nil {
			continue
		}
		e := litExpr(a)
		if negative(e) && (authnCond(e, pw.atomIDs(a)) || credentialRead(a.expr)) {
			return true
		}
	}
	return false
}

// credentialRead: the expression reads a credential from the request: a
// header, cookie or parameter named like one (Authorization, X-Sign,
// api_key, session).
func credentialRead(x ast.Node) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return !hit
		}
		sel, ok := unparen(c.Fun).(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch sel.Sel.Name {
		case "Cookie":
			hit = true // a cookie read is a session read, whatever it's named
			return false
		case "Get", "GetHeader", "Header", "Query", "QueryParam", "FormValue", "Param", "Params", "Values":
		default:
			return true
		}
		for _, a := range c.Args {
			if lit, ok := unparen(a).(*ast.BasicLit); ok && lit.Kind == token.STRING && credentialName(lit.Value) {
				hit = true
			}
		}
		return !hit
	})
	return hit
}

func credentialName(s string) bool {
	low := strings.ToLower(s)
	for _, w := range []string{"authorization", "auth", "token", "sign", "key", "session", "cookie", "bearer", "jwt", "secret", "hmac", "credential"} {
		if strings.Contains(low, w) {
			return true
		}
	}
	return false
}

// pathVerdict reads the handler scopes of a body. ok is false when the
// body has none (no own-level hand-off) or the paths can't be read.
func (w *walker) pathVerdict(body *ast.BlockStmt, fc *fileCtx, depth int, owner *funcInfo, ft *ast.FuncType,
	resolver *walker, params, cfg map[string]bool) (v verdict, why string, covered map[*ast.BlockStmt]bool, ok bool) {
	recv := ""
	if owner != nil && owner.recvType != "" && len(owner.decl.Recv.List[0].Names) > 0 {
		recv = owner.decl.Recv.List[0].Names[0].Name
	}
	declared := declaredNames(ft, body)
	res := &pathResult{}
	var scopes []struct {
		b  *ast.BlockStmt
		ft *ast.FuncType
	}
	// a gin middleware runs the chain on even without c.Next()
	ginMW := func(b *ast.BlockStmt, ft *ast.FuncType) bool {
		return w.r.asMW && resolver.ginScope(ft)
	}
	if ownHandOff(body, recv) || ginMW(body, ft) {
		scopes = append(scopes, struct {
			b  *ast.BlockStmt
			ft *ast.FuncType
		}{body, ft})
	}
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.DeferStmt, *ast.GoStmt:
			return false
		case *ast.ReturnStmt:
			for _, r := range x.Results {
				if fl := returnedLit(r); fl != nil && (ownHandOff(fl.Body, recv) || ginMW(fl.Body, fl.Type)) {
					scopes = append(scopes, struct {
						b  *ast.BlockStmt
						ft *ast.FuncType
					}{fl.Body, fl.Type})
				}
			}
		}
		return true
	})
	if len(scopes) == 0 {
		return vUnknown, "", nil, false
	}
	covered = map[*ast.BlockStmt]bool{}
	for _, s := range scopes {
		covered[s.b] = true
	}
	var pws []*pathWalker
	for _, s := range scopes {
		pw := &pathWalker{w: resolver, base: w, depth: depth, gin: resolver.ginScope(s.ft), recv: recv,
			declared: declared, cfg: cfg, ids: scopeIDs(s.b, params), res: res, ctx: resolver.ginParam(s.ft),
			decls: scopeDecls(s.b, params)}
		pw.scan(s.b, s.ft)
		pw.run(s.b)
		pws = append(pws, pw)
		if res.overflow {
			// too many branches even after merging: say so, never trust the name
			return vComplex, "", covered, true
		}
		if res.opaque {
			return vUnknown, "", nil, false
		}
	}
	v, kind := pws[0].decide()
	if v == vConditional {
		switch kind {
		case akMethod:
			why = "rejects only some methods"
		case akPresence:
			why = "checks credentials only when sent"
		case akBypass:
			why = "a request parameter skips the check"
		default:
			why = "rejects only when a setting is on"
		}
	}
	return v, why, covered, true
}

// scan collects what the scope declares up front: its own parameters,
// local closures, abort aliases and comma-ok lookups.
func (pw *pathWalker) scan(b *ast.BlockStmt, ft *ast.FuncType) {
	pw.req, pw.locals, pw.aborters, pw.okVars = map[string]bool{}, map[string]*ast.FuncLit{}, map[string]bool{}, map[string]bool{}
	defer pw.dataflow(b)
	if ft != nil && ft.Params != nil {
		for _, f := range ft.Params.List {
			for _, n := range f.Names {
				pw.req[n.Name] = true
			}
		}
	}
	ast.Inspect(b, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			return true
		}
		if len(as.Lhs) == 2 && len(as.Rhs) == 1 {
			rhs := unparen(as.Rhs[0])
			_, isIndex := rhs.(*ast.IndexExpr)
			c, isCall := rhs.(*ast.CallExpr)
			if id, ok := as.Lhs[1].(*ast.Ident); ok && (isIndex || isCall && credentialRead(c)) {
				pw.okVars[id.Name] = true
			}
		}
		for i, l := range as.Lhs {
			id, ok := l.(*ast.Ident)
			if !ok || i >= len(as.Rhs) {
				continue
			}
			switch r := unparen(as.Rhs[i]).(type) {
			case *ast.FuncLit:
				pw.locals[id.Name] = r
			case *ast.SelectorExpr:
				if strings.HasPrefix(r.Sel.Name, "Abort") {
					pw.aborters[id.Name] = true
				}
			}
		}
		return true
	})
}

// dataflow sorts the scope's locals: computed from the request (directly
// or through other locals), or assigned exactly once from a setting or
// from a test of a request parameter. Locals assigned more than once and
// tested are tracked instead: their constant value is followed per path.
func (pw *pathWalker) dataflow(b *ast.BlockStmt) {
	pw.derived, pw.setting, pw.bypassVar = map[string]bool{}, map[string]bool{}, map[string]bool{}
	type asg struct {
		name string
		rhs  ast.Expr
	}
	var all []asg
	count := map[string]int{}
	ast.Inspect(b, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.AssignStmt:
			for i, l := range x.Lhs {
				id, ok := l.(*ast.Ident)
				if !ok || id.Name == "_" {
					continue
				}
				rhs := x.Rhs[0]
				if len(x.Rhs) == len(x.Lhs) {
					rhs = x.Rhs[i]
				}
				all = append(all, asg{id.Name, rhs})
				count[id.Name]++
			}
		case *ast.RangeStmt:
			for _, e := range []ast.Expr{x.Key, x.Value} {
				if id, ok := e.(*ast.Ident); ok {
					all = append(all, asg{id.Name, x.X})
					count[id.Name]++
				}
			}
		}
		return true
	})
	for changed := true; changed; {
		changed = false
		for _, a := range all {
			if !pw.derived[a.name] && pw.readsRequest(a.rhs) {
				pw.derived[a.name] = true
				changed = true
			}
		}
	}
	pw.tracked = map[string]bool{}
	tested := map[string]bool{}
	ast.Inspect(b, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.IfStmt:
			ast.Inspect(x.Cond, func(n ast.Node) bool {
				if id, ok := n.(*ast.Ident); ok {
					tested[id.Name] = true
				}
				return true
			})
		case *ast.DeclStmt:
			if g, ok := x.Decl.(*ast.GenDecl); ok && g.Tok == token.VAR {
				for _, sp := range g.Specs {
					for _, n := range sp.(*ast.ValueSpec).Names {
						count[n.Name]++
					}
				}
			}
		}
		return true
	})
	for name, c := range count {
		if c > 1 && tested[name] && !pw.req[name] {
			pw.tracked[name] = true
		}
	}
	for _, a := range all {
		if count[a.name] != 1 {
			continue
		}
		if !pw.derived[a.name] && pw.cfgExpr(a.rhs) {
			pw.setting[a.name] = true
		}
		if requestParam(a.rhs) {
			pw.bypassVar[a.name] = true
		}
	}
}

// checkedBefore: a stop after a failed credential or auth check, whatever
// status it writes (Cloudreve's SignRequired answers 200 with an error body).
func (pw *pathWalker) checkedBefore(p *ppath) bool {
	for _, a := range p.atoms {
		if a.expr == nil {
			continue
		}
		e := litExpr(a)
		if negative(e) {
			if credentialRead(a.expr) {
				return true
			}
			// if !authorize(w, r) { return }: a failed auth-check call
			if c, ok := unparen(a.expr).(*ast.CallExpr); ok && !a.pos {
				if n := exprName(c.Fun); authCheckCall(n) || identityCall(n) {
					return true
				}
			}
			if id := rootIdent(a.expr); id != nil && pw.atomIDs(a)[id.Name] {
				return true
			}
			if b, ok := a.expr.(*ast.BinaryExpr); ok {
				for _, s := range []ast.Expr{b.X, b.Y} {
					if id := rootIdent(s); id != nil && pw.atomIDs(a)[id.Name] {
						return true
					}
				}
			}
		}
	}
	return false
}
