package scan

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path"
	"strconv"
	"strings"
)

// exprName renders a middleware or handler expression as a short dotted
// name. A call renders as its callee, so AuthRequired() becomes
// "AuthRequired" and mw.New(cfg).Handler becomes "mw.New.Handler".
func exprName(x ast.Expr) string {
	switch e := x.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		if base := exprName(e.X); base != "" {
			return base + "." + e.Sel.Name
		}
		return e.Sel.Name
	case *ast.CallExpr:
		return exprName(e.Fun)
	case *ast.ParenExpr:
		return exprName(e.X)
	case *ast.StarExpr:
		return exprName(e.X)
	case *ast.UnaryExpr:
		return exprName(e.X)
	case *ast.IndexExpr:
		return exprName(e.X)
	case *ast.IndexListExpr:
		return exprName(e.X)
	case *ast.CompositeLit:
		if e.Type != nil {
			return exprName(e.Type)
		}
	case *ast.FuncLit:
		return "func"
	}
	return ""
}

// rootIdent returns the identifier at the base of a selector chain.
func rootIdent(x ast.Expr) *ast.Ident {
	for {
		switch e := x.(type) {
		case *ast.Ident:
			return e
		case *ast.SelectorExpr:
			x = e.X
		case *ast.ParenExpr:
			x = e.X
		case *ast.StarExpr:
			x = e.X
		default:
			return nil
		}
	}
}

// keyOf is the env key for an addressable router expression: "r", "s.mux".
func keyOf(x ast.Expr) string {
	switch e := x.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		if k := keyOf(e.X); k != "" {
			return k + "." + e.Sel.Name
		}
	case *ast.ParenExpr:
		return keyOf(e.X)
	case *ast.StarExpr:
		return keyOf(e.X)
	}
	return ""
}

func exprText(x ast.Expr) string {
	s := types.ExprString(x)
	if len(s) > 60 {
		s = s[:57] + "..."
	}
	return s
}

// constString folds a string expression: literals, constants, string
// parameters bound at a call site, `a + "/x"` and path.Join of those.
func (w *walker) constString(x ast.Expr) (string, bool) {
	return w.foldString(x, 0)
}

func (w *walker) foldString(x ast.Expr, depth int) (string, bool) {
	if depth > 8 || x == nil {
		return "", false
	}
	if info := w.fc.pkg.info; info != nil {
		if tv, ok := info.Types[x]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
			return constant.StringVal(tv.Value), true
		}
	}
	if sel, ok := x.(*ast.SelectorExpr); ok {
		if k := keyOf(sel); k != "" {
			if s, ok := w.env.strs[k]; ok {
				return s, true
			}
		}
	}
	switch e := x.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			return s, err == nil
		}
	case *ast.ParenExpr:
		return w.foldString(e.X, depth+1)
	case *ast.Ident:
		if s, ok := w.env.strs[e.Name]; ok {
			return s, true
		}
		if w.env.locals[e.Name] {
			return "", false
		}
		if c, ok := w.fc.pkg.consts[e.Name]; ok {
			return w.with(c.fc).foldString(c.expr, depth+1)
		}
	case *ast.SelectorExpr:
		if id, ok := e.X.(*ast.Ident); ok && !w.env.locals[id.Name] {
			// http.MethodDelete
			if w.fc.imports[id.Name] == "net/http" && strings.HasPrefix(e.Sel.Name, "Method") && len(e.Sel.Name) > 6 {
				return strings.ToUpper(e.Sel.Name[6:]), true
			}
			if p := w.ix.byImp[w.fc.imports[id.Name]]; p != nil {
				if c, ok := p.consts[e.Sel.Name]; ok {
					return w.with(c.fc).foldString(c.expr, depth+1)
				}
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			a, ok1 := w.foldString(e.X, depth+1)
			b, ok2 := w.foldString(e.Y, depth+1)
			if ok1 && ok2 {
				return a + b, true
			}
		}
	case *ast.CallExpr:
		if name := exprName(e.Fun); name == "path.Join" || name == "filepath.Join" {
			var parts []string
			for _, a := range e.Args {
				s, ok := w.foldString(a, depth+1)
				if !ok {
					return "", false
				}
				parts = append(parts, s)
			}
			return path.Join(parts...), true
		}
	}
	return "", false
}

// pathArg returns a route path: its constant value, or the expression text in
// angle brackets with dynamic set.
func (w *walker) pathArg(x ast.Expr) (p string, dynamic bool) {
	if s, ok := w.constString(x); ok {
		return s, false
	}
	return "<" + exprText(x) + ">", true
}

// joinPath joins a router prefix and a registered path the way the routers
// do: one slash between them, a leading slash always, the registered path's
// trailing slash kept.
func joinPath(prefix, p string) string {
	switch {
	case p == "":
		if prefix == "" {
			return "/"
		}
		return prefix
	case prefix == "" || prefix == "/":
		if !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "<") {
			return "/" + p
		}
		return p
	}
	out := strings.TrimSuffix(prefix, "/") + "/" + strings.TrimPrefix(p, "/")
	// gin's r.Group("f") still serves /f
	if !strings.HasPrefix(out, "/") && !strings.HasPrefix(out, "<") {
		out = "/" + out
	}
	return out
}

// relHasPrefix reports whether a path relative to a router falls under a
// prefix-scoped middleware (fiber's app.Use("/api", mw)).
func relHasPrefix(rel, prefix string) bool {
	prefix = strings.TrimSuffix(prefix, "/")
	if prefix == "" {
		return true
	}
	return rel == prefix || strings.HasPrefix(rel, prefix+"/")
}

func firstSegment(p string) string {
	p = strings.TrimPrefix(p, "/")
	if i := strings.IndexByte(p, '/'); i >= 0 {
		return "/" + p[:i]
	}
	return "/" + p
}
