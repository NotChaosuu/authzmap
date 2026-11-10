package scan

import (
	"go/ast"
	"go/token"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/NotChaosuu/authzmap/pkg/model"
)

// registrationVerbs is deliberately broad: it only decides what counts as
// "looks like a route" when reporting calls the walk never attributed.
var registrationVerbs = set("Get", "Post", "Put", "Delete", "Patch", "Head", "Options",
	"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "Any", "All",
	"Handle", "HandleFunc", "HandlePath", "Method", "MethodFunc", "Match", "Add")

// upperVerbs count as registrations even when the path isn't a constant:
// hardly anything but a router names a method GET.
var upperVerbs = set("GET", "POST", "PUT", "DELETE", "PATCH")

// unseenCalls finds route-like code that produced no route: a verb method
// whose path argument is a "/path" or "GET /path" constant, or hand-rolled
// dispatch on r.Method and r.URL.Path. It names the imports those files use
// instead of a supported router. Any of these make the scan incomplete: a
// project on an unsupported router must not read as "no gaps".
func unseenCalls(ix *index, routes []model.Route, untraced, mounts []model.Site) ([]model.Site, string) {
	seen := map[string]bool{}
	for _, s := range mounts {
		seen[s.File+":"+strconv.Itoa(s.Line)] = true
	}
	for _, r := range routes {
		seen[r.File+":"+strconv.Itoa(r.Line)] = true
	}
	for _, s := range untraced {
		seen[s.File+":"+strconv.Itoa(s.Line)] = true
	}
	var out []model.Site
	imports := map[string]int{}
	r := newRun(ix)
	for _, p := range ix.pkgs {
		for _, fc := range p.files {
			var here []model.Site
			w := &walker{r: r, ix: ix, fc: fc, env: newEnv(), consumed: map[*ast.FuncLit]bool{}}
			add := func(n ast.Node, call string) {
				line := ix.fset.Position(n.Pos()).Line
				if !seen[fc.rel+":"+strconv.Itoa(line)] {
					seen[fc.rel+":"+strconv.Itoa(line)] = true
					here = append(here, model.Site{File: fc.rel, Line: line, Call: call})
				}
			}
			ast.Inspect(fc.file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					if len(x.Args) < 2 {
						return true
					}
					sel, ok := x.Fun.(*ast.SelectorExpr)
					if !ok || !registrationVerbs[sel.Sel.Name] {
						return true
					}
					// the path is the first argument, or the second after a
					// method (Handle("POST", "/x", h)); constants fold. A
					// path built at run time counts after a method token
					// (grpc-gateway's Handle("DELETE", pattern_X, h)), from a
					// "/..." piece (baseURL + "/admin"), or on GET/POST/...
					for i := 0; i < 2 && i < len(x.Args)-1; i++ {
						v, ok := w.constString(x.Args[i])
						if ok && routeLiteral(v) {
							add(x, exprText(x.Fun))
							break
						}
						if !ok && (i == 1 || pathShaped(x.Args[i]) || upperVerbs[sel.Sel.Name]) {
							add(x, exprText(x.Fun))
							break
						}
						if !ok || !isMethodToken(v) {
							break
						}
					}
				case *ast.IfStmt:
					// hand-rolled dispatch: if r.Method == "DELETE" && r.URL.Path == "/admin/users" { s.handler(w, r) }
					if b, ok := unparen(x.Cond).(*ast.BinaryExpr); ok && b.Op == token.LAND &&
						comparesField(b, "Method") && w.comparesPath(b) && callsSomething(x.Body.List) {
						add(x.Cond, "dispatch on "+exprText(x.Cond))
					}
				case *ast.SwitchStmt:
					// switch r.URL.Path { case "/a": ... case "/b": ... }
					if sel, ok := unparen(x.Tag).(*ast.SelectorExpr); ok && sel.Sel.Name == "Path" && isURLField(sel.X) {
						for _, c := range x.Body.List {
							cc := c.(*ast.CaseClause)
							// a case that only returns true/false is a path guard, not a route
							if !callsSomething(cc.Body) {
								continue
							}
							for _, e := range cc.List {
								if v, ok := w.constString(e); ok && strings.HasPrefix(v, "/") {
									add(e, "dispatch on "+exprText(x.Tag)+" == "+v)
								}
							}
						}
					}
				}
				return true
			})
			// a file that imports a router authzmap doesn't model registers
			// routes it can't see, whatever the call shapes look like
			for _, imp := range fc.file.Imports {
				ip := strings.Trim(imp.Path.Value, `"`)
				if pkg := unmodelledRouter(ip); pkg != "" {
					add(imp, "imports "+pkg+": its routes are not analysed")
				}
			}
			if len(here) == 0 {
				continue
			}
			out = append(out, here...)
			for local, ip := range fc.imports {
				if fc.fws[local] == nil && strings.Contains(ip, ".") && !underFramework(ip) {
					imports[ip] += len(here)
				}
			}
		}
	}
	if len(out) == 0 {
		return []model.Site{}, ""
	}
	type kv struct {
		k string
		v int
	}
	var top []kv
	for k, v := range imports {
		top = append(top, kv{k, v})
	}
	routerish := func(p string) bool {
		for _, w := range []string{"route", "mux", "http", "web", "rest", "api"} {
			if strings.Contains(path.Base(p), w) {
				return true
			}
		}
		return false
	}
	sort.Slice(top, func(i, j int) bool {
		if ri, rj := routerish(top[i].k), routerish(top[j].k); ri != rj {
			return ri
		}
		if top[i].v != top[j].v {
			return top[i].v > top[j].v
		}
		return top[i].k < top[j].k
	})
	note := ""
	if len(top) > 0 {
		note = "the unchecked route-like calls are in files that import " + top[0].k
		if len(top) > 1 {
			note += ", " + top[1].k
		}
	}
	return out, note
}

// unmodelledRouterPkgs are HTTP routers and frameworks authzmap knows of
// but doesn't model.
var unmodelledRouterPkgs = []string{
	"github.com/emicklei/go-restful", "goa.design/goa", "github.com/goadesign/goa", "github.com/gorilla/rpc",
	"github.com/beego/beego", "github.com/astaxie/beego", "github.com/revel/revel", "github.com/gobuffalo/buffalo",
	"github.com/kataras/iris", "github.com/julienschmidt/httprouter", "github.com/uptrace/bunrouter",
	"github.com/go-zoo/bone", "github.com/bmizerany/pat", "github.com/gorilla/pat", "github.com/dimfeld/httptreemux",
	"github.com/gocraft/web", "github.com/fasthttp/router", "github.com/buaazp/fasthttprouter",
	"github.com/go-martini/martini", "github.com/zenazn/goji", "goji.io", "github.com/go-macaron/macaron",
	"gopkg.in/macaron.v1", "github.com/gogf/gf", "github.com/zeromicro/go-zero/rest", "github.com/twitchtv/twirp",
	"github.com/go-kratos/kratos/v2/transport/http", "github.com/danielgtaylor/huma", "github.com/go-fuego/fuego",
	"github.com/swaggest/rest", "github.com/savsgio/atreugo", "github.com/aofei/air", "github.com/go-pkgz/routegroup",
}

// underFramework reports whether ip is a subpackage of a supported router,
// such as chi/v5/middleware. Those never get named as the unmodelled router.
func underFramework(ip string) bool {
	for p := ip; ; {
		i := strings.LastIndexByte(p, '/')
		if i < 0 {
			return false
		}
		p = p[:i]
		if frameworkForImport(p) != nil {
			return true
		}
	}
}

func unmodelledRouter(ip string) string {
	for _, p := range unmodelledRouterPkgs {
		if ip == p || strings.HasPrefix(ip, p+"/") {
			return p
		}
	}
	return ""
}

// callsSomething: the branch does work (calls a handler, writes a response)
// rather than only returning a value.
func callsSomething(list []ast.Stmt) bool {
	hit := false
	for _, st := range list {
		ast.Inspect(st, func(n ast.Node) bool {
			if _, ok := n.(*ast.CallExpr); ok {
				hit = true
			}
			return !hit
		})
	}
	return hit
}

// comparesField: an == comparison on X.<field> somewhere in a && chain.
func comparesField(x ast.Expr, field string) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		if b, ok := n.(*ast.BinaryExpr); ok && b.Op == token.EQL {
			for _, side := range []ast.Expr{b.X, b.Y} {
				if sel, ok := unparen(side).(*ast.SelectorExpr); ok && sel.Sel.Name == field {
					hit = true
				}
			}
		}
		return !hit
	})
	return hit
}

// comparesPath: r.URL.Path == "/x" (or a constant) in a && chain.
func (w *walker) comparesPath(x ast.Expr) bool {
	hit := false
	ast.Inspect(x, func(n ast.Node) bool {
		b, ok := n.(*ast.BinaryExpr)
		if !ok || b.Op != token.EQL {
			return true
		}
		for _, pair := range [][2]ast.Expr{{b.X, b.Y}, {b.Y, b.X}} {
			sel, ok := unparen(pair[0]).(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Path" || !isURLField(sel.X) {
				continue
			}
			if v, ok := w.constString(pair[1]); ok && strings.HasPrefix(v, "/") {
				hit = true
			}
		}
		return !hit
	})
	return hit
}

func isURLField(x ast.Expr) bool {
	sel, ok := unparen(x).(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "URL"
}

// routeLiteral: "/users", "GET /users", "POST example.com/x".
func routeLiteral(s string) bool {
	if strings.HasPrefix(s, "/") {
		return true
	}
	if i := strings.IndexByte(s, ' '); i > 0 && isMethodToken(s[:i]) {
		return strings.Contains(s[i+1:], "/")
	}
	return false
}
