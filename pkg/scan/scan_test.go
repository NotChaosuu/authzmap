package scan

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NotChaosuu/authzmap/pkg/model"
)

// tree writes files into a temp dir and returns it.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, src := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func scanFiles(t *testing.T, files map[string]string) *Result {
	t.Helper()
	res, err := Scan(tree(t, files))
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	return res
}

// scanSrc scans one inline main.go.
func scanSrc(t *testing.T, src string) []model.Route {
	t.Helper()
	return scanFiles(t, map[string]string{"main.go": src}).Routes
}

func routeSet(routes []model.Route) map[string][]string {
	m := map[string][]string{}
	for _, r := range routes {
		m[r.Method+" "+r.Path] = r.Chain
	}
	return m
}

func byKey(routes []model.Route) map[string]model.Route {
	m := map[string]model.Route{}
	for _, r := range routes {
		m[r.Method+" "+r.Path] = r
	}
	return m
}

func chainIs(t *testing.T, got map[string][]string, key, want string) {
	t.Helper()
	c, ok := got[key]
	if !ok {
		t.Errorf("route %q missing; have %v", key, keys(got))
		return
	}
	if strings.Join(c, ",") != want {
		t.Errorf("%q chain = %v, want %s", key, c, want)
	}
}

func keys(m map[string][]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---- fixtures ----

func TestScanChiBasic(t *testing.T) {
	res, err := Scan("testdata/chi_basic")
	if err != nil {
		t.Fatal(err)
	}
	got := routeSet(res.Routes)
	if len(res.Routes) != 5 {
		t.Fatalf("got %d routes: %v", len(res.Routes), keys(got))
	}
	chainIs(t, got, "POST /login", "middleware.Logger")
	chainIs(t, got, "GET /me", "middleware.Logger,RequireAuth")
	chainIs(t, got, "GET /orders", "middleware.Logger,RequireAuth")
	chainIs(t, got, "GET /admin/stats", "middleware.Logger")
	chainIs(t, got, "POST /admin/users", "middleware.Logger,RequireAuth")
}

func TestScanGinBasic(t *testing.T) {
	res, err := Scan("testdata/gin_basic")
	if err != nil {
		t.Fatal(err)
	}
	got := routeSet(res.Routes)
	chainIs(t, got, "POST /login", "Logger")
	chainIs(t, got, "GET /api/me", "Logger,RequireAuth")
	chainIs(t, got, "GET /admin/stats", "Logger")
	chainIs(t, got, "POST /admin/users", "Logger,AdminAuth")
}

func TestScanNetHTTPBasic(t *testing.T) {
	res, err := Scan("testdata/nethttp_basic")
	if err != nil {
		t.Fatal(err)
	}
	got := routeSet(res.Routes)
	chainIs(t, got, "POST /login", "")
	chainIs(t, got, "GET /me", "RequireAuth")
	chainIs(t, got, "POST /admin/users", "RequireAuth")
}

// ---- what was not read is reported ----

// A syntax error in the file that registers every route must not read as
// "0 routes, no gaps": the file is a Problem, and the routes that parsed
// before the error are still found.
func TestParseErrorIsReported(t *testing.T) {
	res := scanFiles(t, map[string]string{
		"main.go": `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/admin/users", h)
	r.Get("/admin/stats", h)
}
func zzBroken( {
`,
	})
	if res.Complete() || len(res.Problems) != 1 || res.Problems[0].File != "main.go" {
		t.Fatalf("problems = %+v, want one for main.go", res.Problems)
	}
	if m := res.Problems[0].Message; strings.Contains(m, "main.go") || !strings.HasPrefix(m, "8:") {
		t.Errorf("message should be line:col: error, without the path again: %q", m)
	}
	if len(res.Routes) != 2 {
		t.Errorf("routes before the syntax error should still be read, got %d", len(res.Routes))
	}
}

func TestEmptyDirIsAnError(t *testing.T) {
	if _, err := Scan(t.TempDir()); err == nil {
		t.Fatal("a directory with no Go files must be an error, not an empty clean result")
	}
}

func TestSingleFileTarget(t *testing.T) {
	dir := tree(t, map[string]string{"a.go": "package a\nimport \"net/http\"\nfunc f() { http.HandleFunc(\"/x\", h) }\n"})
	res, err := Scan(filepath.Join(dir, "a.go"))
	if err != nil || len(res.Routes) != 1 {
		t.Fatalf("single-file scan: %v %+v", err, res)
	}
}

// ---- routers that cross function boundaries ----

// Routes registered in `func routes(r chi.Router)` are found.
func TestRouterParameterIsScanned(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func routes(r chi.Router) {
	r.With(RequireAuth).Get("/admin/real", h)
	r.Get("/admin/stats", h)
}`)
	got := byKey(routes)
	if len(got) != 2 {
		t.Fatalf("got %d routes, want 2", len(got))
	}
	if got["GET /admin/stats"].Partial == "" {
		t.Error("a router seen only as a parameter must be marked partial")
	}
}

// The caller's middleware reaches routes registered in a helper.
func TestHelperInheritsCallerChain(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	register(r, "/v1")
}
func register(r chi.Router, base string) {
	r.Get(base+"/users", h)
}`)
	got := byKey(routes)
	rt, ok := got["GET /v1/users"]
	if !ok {
		t.Fatalf("helper route with a folded path missing: %v", got)
	}
	if strings.Join(rt.Chain, ",") != "RequireAuth" || rt.Partial != "" {
		t.Errorf("chain=%v partial=%q, want [RequireAuth] and fully traced", rt.Chain, rt.Partial)
	}
	if len(routes) != 1 {
		t.Errorf("the helper's own standalone walk must not add a duplicate: %d routes", len(routes))
	}
}

// A router in a struct field, set up in a constructor, used in a method.
func TestStructFieldRouterAcrossMethods(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/gin-gonic/gin"
type Server struct{ engine *gin.Engine }
func New() *Server {
	s := &Server{engine: gin.New()}
	s.engine.Use(AuthRequired())
	s.routes()
	return s
}
func (s *Server) routes() {
	s.engine.GET("/items", h)
}`)
	rt, ok := byKey(routes)["GET /items"]
	if !ok || strings.Join(rt.Chain, ",") != "AuthRequired" {
		t.Fatalf("field router route = %+v (ok=%v)", rt, ok)
	}
	if len(routes) != 1 {
		t.Errorf("want 1 route, got %d (Use applied twice or duplicate walk?)", len(routes))
	}
}

// A factory returns a router that is mounted: prefix and parent chain apply.
func TestFactoryMount(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.Use(Logger)
	r.Mount("/admin", adminRouter())
}
func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(RequireAdmin)
	r.Get("/stats", h)
	return r
}`)
	rt, ok := byKey(routes)["GET /admin/stats"]
	if !ok {
		t.Fatalf("mounted route missing: %v", routeSet(routes))
	}
	if strings.Join(rt.Chain, ",") != "Logger,RequireAdmin" || rt.ScopeKind != model.ScopeGroup {
		t.Errorf("chain=%v kind=%s", rt.Chain, rt.ScopeKind)
	}
}

// A package-level router built in a var initializer and used from a func.
func TestPackageLevelRouter(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/labstack/echo/v4"
var e = echo.New()
func init() { e.Use(JWTAuth) }
func routes() { e.GET("/me", h) }`)
	chainIs(t, routeSet(routes), "GET /me", "JWTAuth")
}

// A cobra-style command whose Run closure lives in a package-level var.
func TestRoutesFromVarClosure(t *testing.T) {
	routes := scanFiles(t, map[string]string{
		"cmd.go": `package main
var webCmd = &Command{Run: func() { RegisterRoutes(NewEcho()) }}`,
		"routes.go": `package main
import "github.com/labstack/echo/v4"
func NewEcho() *echo.Echo {
	e := echo.NewWithConfig(echo.Config{})
	e.Use(RequireAuth)
	return e
}
func RegisterRoutes(e *echo.Echo) { e.GET("/tasks", h) }`,
	}).Routes
	rt, ok := byKey(routes)["GET /tasks"]
	if !ok || rt.Partial != "" || strings.Join(rt.Chain, ",") != "RequireAuth" {
		t.Fatalf("route = %+v (ok=%v), want traced with RequireAuth", rt, ok)
	}
}

// opengist's shape: a wrapper type embedding the router, with its own
// SubGroup method, and Use() inside a bare { } block.
func TestEmbeddingWrapperType(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/labstack/echo/v4"
type Router struct{ *echo.Group }
func NewRouter(g *echo.Group) *Router { return &Router{Group: g} }
func (r *Router) SubGroup(prefix string) *Router { return NewRouter(r.Group.Group(prefix)) }
func main() {
	e := echo.New()
	r := NewRouter(e.Group(""))
	r.POST("/upload", upload, requireLogin)
	s := r.SubGroup("/gists")
	{
		s.Use(requireLogin)
		s.GET("/:id", show)
	}
}`)
	got := routeSet(routes)
	chainIs(t, got, "POST /upload", "requireLogin")
	chainIs(t, got, "GET /gists/:id", "requireLogin")
}

// Cross-package call, resolved through go.mod's module path.
func TestCrossPackageCall(t *testing.T) {
	routes := scanFiles(t, map[string]string{
		"go.mod": "module example.com/app\n\ngo 1.22\n",
		"main.go": `package main
import (
	"github.com/go-chi/chi/v5"
	"example.com/app/api"
)
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	api.Register(r)
}`,
		"api/api.go": `package api
import "github.com/go-chi/chi/v5"
func Register(r chi.Router) { r.Delete("/users/{id}", h) }`,
	}).Routes
	rt, ok := byKey(routes)["DELETE /users/{id}"]
	if !ok || rt.Partial != "" || strings.Join(rt.Chain, ",") != "RequireAuth" {
		t.Fatalf("cross-package route = %+v (ok=%v)", rt, ok)
	}
}

// Relative and absolute roots give the same result; module mapping matches
// directories by prefix, so "." must be made absolute first.
func TestScanRootFormsAgree(t *testing.T) {
	dir := tree(t, map[string]string{
		"go.mod": "module example.com/app\n",
		"main.go": `package main
import ("github.com/go-chi/chi/v5"; "example.com/app/api")
func main() { r := chi.NewRouter(); r.Use(RequireAuth); api.Register(r) }`,
		"api/api.go": "package api\nimport \"github.com/go-chi/chi/v5\"\nfunc Register(r chi.Router) { r.Get(\"/x\", h) }\n",
	})
	abs, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	rel, err := Scan(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(abs.Routes) != 1 || len(rel.Routes) != 1 ||
		strings.Join(abs.Routes[0].Chain, ",") != strings.Join(rel.Routes[0].Chain, ",") || rel.Routes[0].Partial != "" {
		t.Fatalf("abs=%+v rel=%+v", abs.Routes, rel.Routes)
	}
}

// One helper reached from two servers is one route; only middleware on both
// paths counts.
func TestTwoCallPathsMerge(t *testing.T) {
	routes := scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func a() { r := chi.NewRouter(); r.Use(RequireAuth); reg(r) }
func b() { r := chi.NewRouter(); reg(r) }
func reg(r chi.Router) { r.Get("/x", h) }`)
	if len(routes) != 1 {
		t.Fatalf("got %d routes, want 1 merged", len(routes))
	}
	// both chains are kept; the analysis judges the route by the open one
	open := false
	for _, v := range routes[0].Variants {
		open = open || len(v.Chain) == 0
	}
	if len(routes[0].Variants) != 2 || !open {
		t.Errorf("auth on only one path must not count: variants=%+v", routes[0].Variants)
	}
}

// A registration on a value we can't trace is still a route, marked
// partial and listed.
func TestUntracedReceiver(t *testing.T) {
	res := scanFiles(t, map[string]string{"main.go": `package main
import "github.com/gin-gonic/gin"
func setup(app *App) { app.Router().GET("/x", h); deps.Router.POST("/y", h) }`})
	got := byKey(res.Routes)
	if rt, ok := got["POST /y"]; !ok || rt.Partial == "" {
		t.Errorf("untraced route = %+v (ok=%v)", rt, ok)
	}
	if len(res.Untraced) == 0 {
		t.Error("untraced registrations must be listed")
	}
}

// ---- middleware semantics per framework ----

// gin captures middleware at registration; echo's root applies Use() per
// request; a chi With() view doesn't leak.
func TestUseOrderSemantics(t *testing.T) {
	gin := routeSet(scanSrc(t, `package main
import "github.com/gin-gonic/gin"
func main() { r := gin.New(); r.GET("/before", h); r.Use(Auth()); r.GET("/after", h) }`))
	chainIs(t, gin, "GET /before", "")
	chainIs(t, gin, "GET /after", "Auth")

	echo := routeSet(scanSrc(t, `package main
import "github.com/labstack/echo/v4"
func main() { e := echo.New(); e.GET("/before", h); e.Use(Auth) }`))
	chainIs(t, echo, "GET /before", "Auth")

	chi := routeSet(scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	r.Use(Logger)
	r.With(RequireAuth).With(RateLimit).Get("/me", h)
	r.Get("/public", h)
}`))
	chainIs(t, chi, "GET /me", "Logger,RequireAuth,RateLimit")
	chainIs(t, chi, "GET /public", "Logger")
}

// Middleware added inside an if reaches later routes as conditional; routes
// inside the block have it unconditionally.
func TestConditionalUse(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import "github.com/labstack/echo/v4"
func main() {
	e := echo.New()
	g := e.Group("/debug")
	if auth := basicAuth(); auth != nil {
		g.Use(auth)
		g.GET("/inside", h)
	}
	g.GET("/after", h)
}`))
	if c := got["GET /debug/after"].Conditional; len(c) != 1 {
		t.Errorf("after-block route conditional = %v", c)
	}
	if c := got["GET /debug/inside"]; len(c.Conditional) != 0 || len(c.Chain) != 1 {
		t.Errorf("in-block route = %+v", c)
	}
}

func TestEchoFiberGorillaShapes(t *testing.T) {
	echo := routeSet(scanSrc(t, `package main
import "github.com/labstack/echo/v4"
func routes() {
	e := echo.New()
	api := e.Group("/api", RequireAuth)
	api.GET("/me", me)
	e.GET("/admin", stats, AdminOnly)
}`))
	chainIs(t, echo, "GET /api/me", "RequireAuth")
	chainIs(t, echo, "GET /admin", "AdminOnly")

	fiber := routeSet(scanSrc(t, `package main
import "github.com/gofiber/fiber/v2"
func routes() {
	app := fiber.New()
	app.Use("/api", RequireAuth)
	app.Get("/api/me", me)
	app.Get("/login", login)
}`))
	chainIs(t, fiber, "GET /api/me", "RequireAuth")
	chainIs(t, fiber, "GET /login", "")

	gorilla := routeSet(scanSrc(t, `package main
import "github.com/gorilla/mux"
func routes() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	api := r.PathPrefix("/api").Subrouter()
	api.Use(RequireAuth)
	api.Path("/me").Methods("GET").HandlerFunc(me)
	r.HandleFunc("/health", health)
}`))
	chainIs(t, gorilla, "POST /login", "")
	chainIs(t, gorilla, "GET /api/me", "RequireAuth")
	chainIs(t, gorilla, "* /health", "")
}

// net/http: a sub-mux mounted with StripPrefix, and auth wrapped around the
// whole mux where it's served.
func TestNetHTTPMountAndServeWrapper(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "net/http"
func main() {
	api := http.NewServeMux()
	api.HandleFunc("GET /users", users)
	root := http.NewServeMux()
	root.Handle("/api/", http.StripPrefix("/api", api))
	http.ListenAndServe(":8080", requireAuth(root))
}`))
	chainIs(t, got, "GET /api/users", "requireAuth")
}

// gin: g.OPTIONS("/*any") with no handler of its own is still a route.
func TestGinHandlerlessRoute(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "github.com/gin-gonic/gin"
func main() { g := gin.New(); g.Use(cors()); g.OPTIONS("/*any") }`))
	if _, ok := got["OPTIONS /*any"]; !ok {
		t.Fatalf("handler-less route missing: %v", keys(got))
	}
}

// Paths are normalized with a leading slash (gotify's g.GET("version")).
func TestPathNormalized(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "github.com/gin-gonic/gin"
func main() { g := gin.New(); g.GET("version", h) }`))
	if _, ok := got["GET /version"]; !ok {
		t.Fatalf("got %v", keys(got))
	}
}

// ---- reading middleware and handler bodies ----

// miniflux's shape: a "session" middleware that redirects anonymous users
// to login after an early pass-through for static assets.
func TestEnforcingMiddlewareBody(t *testing.T) {
	rt := scanSrc(t, `package main
import "net/http"
type sessionMW struct{}
func newSessionMW() *sessionMW { return &sessionMW{} }
func (m *sessionMW) handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isStatic(r) {
			next.ServeHTTP(w, r)
			return
		}
		if !request.IsAuthenticated(r) {
			http.Redirect(w, r, loginURL, http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Serve() http.Handler {
	s := newSessionMW()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /feeds", feeds)
	return s.handle(mux)
}`)
	if len(rt) != 1 || len(rt[0].Enforcing) != 1 {
		t.Fatalf("routes=%+v, want the session middleware read as enforcing", rt)
	}
}

// woodpecker's SetUser: anonymous requests pass; a 401 only for a present
// session that fails a check. Error mappers write 401 after next.
func TestNonEnforcingBodies(t *testing.T) {
	rt := scanSrc(t, `package main
import ("net/http"; "github.com/gin-gonic/gin")
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := token.ParseRequest(c.Request, func(t *Token) (string, error) {
			u, err := store.GetUser(t.ID)
			if err != nil { return "", err }
			return u.Hash, nil
		})
		if err == nil {
			if bad(t) { c.AbortWithStatus(http.StatusUnauthorized); return }
		}
		c.Next()
	}
}
func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if isAuthErr(c) { c.AbortWithStatus(401) }
	}
}
func main() { r := gin.New(); r.Use(SetUser(), Errors()); r.GET("/x", h) }`)
	if len(rt) != 1 || len(rt[0].Enforcing) != 0 {
		t.Fatalf("enforcing = %v, want none", rt[0].Enforcing)
	}
}

// A handler that looks the caller up and rejects absence checks auth itself;
// one that only checks permissions on an identity nobody checked doesn't.
func TestHandlerChecks(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import ("net/http"; "github.com/go-chi/chi/v5")
func checked(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if u == nil { http.Error(w, "no", http.StatusUnauthorized); return }
}
func permsOnly(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.HasPerm("x") { w.WriteHeader(403); return }
	if err := UpdateUser(u); err != nil { return }
}
func main() { r := chi.NewRouter(); r.Get("/a", checked); r.Get("/b", permsOnly) }`))
	if !got["GET /a"].HandlerChecks || got["GET /b"].HandlerChecks {
		t.Fatalf("a=%v b=%v", got["GET /a"].HandlerChecks, got["GET /b"].HandlerChecks)
	}
}

// writefreely's shape: handler.Download(fn, UserLevelUser) wraps the real
// handler; the wrapped handler, not the level, is what gets read.
func TestWrappedHandlerPicked(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import ("net/http"; "github.com/gorilla/mux")
func export(app *App, w http.ResponseWriter, r *http.Request) error {
	u := getUserSession(app, r)
	if u == nil { return ErrNotLoggedIn }
	return nil
}
func UserLevelUser(c *Config) int { return 1 }
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/export", handler.Download(export, UserLevelUser))
}`))
	rt := got["* /export"]
	if rt.Handler != "export" || !rt.HandlerChecks {
		t.Fatalf("handler=%q checks=%v", rt.Handler, rt.HandlerChecks)
	}
}

func TestResultIsDeterministic(t *testing.T) {
	dir := tree(t, map[string]string{"main.go": `package main
import "github.com/go-chi/chi/v5"
type S struct{ r chi.Router }
func (s *S) a() { s.r.Use(A); s.b() }
func (s *S) b() { s.r.Get("/x", h); s.r.Route("/y", func(r chi.Router) { r.Get("/z", h) }) }`})
	first, _ := Scan(dir)
	for i := 0; i < 5; i++ {
		again, _ := Scan(dir)
		if len(again.Routes) != len(first.Routes) {
			t.Fatalf("run %d: %d routes vs %d", i, len(again.Routes), len(first.Routes))
		}
		for j := range again.Routes {
			if strings.Join(again.Routes[j].Chain, ",") != strings.Join(first.Routes[j].Chain, ",") {
				t.Fatalf("run %d differs at %d", i, j)
			}
		}
	}
}

// chi Mount strips its prefix; Handle with a router doesn't (headscale
// mounted /api/v2 that way and the path came out doubled).
func TestChiMountVersusHandle(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func main() {
	r := chi.NewRouter()
	a := chi.NewRouter()
	a.Get("/users", h)
	b := chi.NewRouter()
	b.Post("/api/v2/token", h)
	r.Mount("/admin", a)
	r.Route("/api", func(r chi.Router) { r.Handle("/v2/*", b) })
}`))
	if _, ok := got["GET /admin/users"]; !ok {
		t.Errorf("Mount should prefix: %v", keys(got))
	}
	if _, ok := got["POST /api/v2/token"]; !ok {
		t.Errorf("Handle should keep the sub-router's full paths: %v", keys(got))
	}
}

// A router mounted with nothing we can read inside it shows up as a partial
// catch-all at the mount point instead of vanishing.
func TestEmptyMountIsVisible(t *testing.T) {
	rt := byKey(scanSrc(t, `package main
import ("net/http"; "github.com/go-chi/chi/v5")
func build(api *chi.Mux) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) { r.Handle("/v1/*", api) })
	return r
}
func main() { build(newAPI()); http.ListenAndServe(":1", nil) }
func newAPI() *chi.Mux { return chi.NewMux() }`))
	m, ok := rt["* /api/v1/*"]
	if !ok || m.Partial == "" {
		t.Fatalf("empty mount route = %+v (ok=%v), have %v", m, ok, rt)
	}
}

// Calls on a router authzmap doesn't model are reported, with the import.
func TestUnseenRouterIsNoted(t *testing.T) {
	res := scanFiles(t, map[string]string{"main.go": `package main
import ("net/http"; "github.com/go-pkgz/routegroup")
func main() {
	g := routegroup.New(http.NewServeMux())
	g.HandleFunc("GET /users", h)
	g.HandleFunc("DELETE /users/{id}", h)
}`})
	// the two calls, and the import of a router that isn't modelled
	if len(res.Unseen) != 3 || res.Complete() {
		t.Fatalf("unseen = %+v", res.Unseen)
	}
	if len(res.Notes) == 0 || !strings.Contains(res.Notes[0], "routegroup") {
		t.Errorf("notes = %v", res.Notes)
	}
}

// A supported router's own subpackage (chi's middleware) is never named as
// the import that might be the unmodelled router.
func TestUnseenNoteSkipsFrameworkSubpackages(t *testing.T) {
	res := scanFiles(t, map[string]string{"main.go": `package main
import ("net/http"; "github.com/go-chi/chi/v5"; "github.com/go-chi/chi/v5/middleware"; "github.com/julienschmidt/httprouter")
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	legacy := httprouter.New()
	legacy.GET("/legacy/export", nil)
	r.Mount("/old", legacy)
	http.ListenAndServe(":8080", r)
}`})
	if len(res.Notes) == 0 || strings.Contains(res.Notes[0], "chi/v5/middleware") || !strings.Contains(res.Notes[0], "httprouter") {
		t.Errorf("notes = %v", res.Notes)
	}
}

// Routes behind a flag the caller binds to false are never registered
// (navidrome's api.R(r, "/translation", ..., false)).
func TestConstantFlagSkipsBranch(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
func R(r chi.Router, p string, persistable bool) {
	r.Get(p, h)
	if persistable {
		r.Put(p, h)
	}
}
func main() { r := chi.NewRouter(); R(r, "/translation", false); R(r, "/user", true) }`))
	if _, ok := got["PUT /translation"]; ok {
		t.Error("route under a false flag must not be reported")
	}
	if _, ok := got["PUT /user"]; !ok {
		t.Errorf("route under a true flag missing: %v", keys(got))
	}
}

// opengist: r.SubGroup("/api/ipc", ipcAuth) hands its middleware to closures;
// the middleware still lands on the group it returns.
func TestWrapperGroupMiddlewareArgs(t *testing.T) {
	got := routeSet(scanSrc(t, `package main
import "github.com/labstack/echo/v4"
type Router struct{ *echo.Group }
func NewRouter(g *echo.Group) *Router { return &Router{Group: g} }
func (r *Router) SubGroup(prefix string, m ...Middleware) *Router {
	mws := make([]echo.MiddlewareFunc, len(m))
	for i := range m { mws[i] = func(next echo.HandlerFunc) echo.HandlerFunc { return next } }
	return NewRouter(r.Group.Group(prefix, mws...))
}
func main() {
	r := NewRouter(echo.New().Group(""))
	ipc := r.SubGroup("/api/ipc", ipcAuth)
	ipc.POST("/hook", hook)
}`))
	if c := got["POST /api/ipc/hook"]; !strings.Contains(strings.Join(c, ","), "ipcAuth") {
		t.Errorf("chain = %v", c)
	}
}

// owncast: a handler variable set inside an if is still the mux mounted
// later outside it.
func TestAssignInBlockReachesOuterVar(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import "net/http"
func requireAdmin(next http.HandlerFunc) http.HandlerFunc { return next }
func admin() http.Handler {
	m := http.NewServeMux()
	m.Handle("/api/admin/plugins/", requireAdmin(h))
	return m
}
func main() {
	var pluginAdmin http.Handler
	if enabled {
		pluginAdmin = admin()
	}
	root := http.NewServeMux()
	root.Handle("/api/admin/plugins/", pluginAdmin)
	http.ListenAndServe(":1", root)
}`))
	rt, ok := got["* /api/admin/plugins/"]
	if !ok || rt.Partial != "" || len(rt.Chain) == 0 {
		t.Fatalf("mounted admin mux route = %+v (ok=%v)", rt, ok)
	}
}

// oapi-codegen: r := opts.Router; if r == nil { r = chi.NewRouter() } keeps
// the router it was given rather than a fresh unguarded one.
func TestConditionalReassignKeepsRouter(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import "github.com/go-chi/chi/v5"
type Opts struct{ BaseRouter chi.Router }
func Handler(opts Opts) {
	r := opts.BaseRouter
	if r == nil {
		r = chi.NewRouter()
	}
	r.Post("/admin/config", h)
}`))
	if rt := got["POST /admin/config"]; rt.Partial == "" {
		t.Errorf("route on the given router should stay partial: %+v", rt)
	}
}

// A recursive function that carries no router loses no routes when cut.
func TestRecursionWithoutRouterIsComplete(t *testing.T) {
	res := scanFiles(t, map[string]string{"main.go": `package main
import "github.com/go-chi/chi/v5"
func serve(r chi.Router) { r.Get("/x", h); retry(3) }
func retry(n int) { if n > 0 { retry(n - 1) }; restart() }
func restart() { serveAgain() }
func serveAgain() { restart() }
func main() { serve(chi.NewRouter()) }`})
	if !res.Complete() {
		t.Errorf("problems = %+v", res.Problems)
	}
}

// A path switch that only answers true/false is a guard, not dispatch; a
// mount site is not an unseen route.
func TestCoverageIgnoresGuardsAndMounts(t *testing.T) {
	res := scanFiles(t, map[string]string{"main.go": `package main
import "net/http"
func exempt(r *http.Request) bool {
	switch r.URL.Path {
	case "/.well-known/webfinger", "/nodeinfo/2.0":
		return true
	}
	return false
}
func api() *http.ServeMux { m := http.NewServeMux(); m.HandleFunc("GET /v1/me", h); return m }
func main() {
	root := http.NewServeMux()
	root.Handle("/v1/", api())
	http.ListenAndServe(":1", root)
}`})
	if len(res.Unseen) != 0 || !res.Complete() {
		t.Errorf("unseen = %+v", res.Unseen)
	}
}

// ---- shapes from real apps ----

// gotify's shape: a gin middleware with no c.Next() that hands the context
// to a helper which aborts with a 401. gin runs the chain on unless
// aborted, so the helper's abort is the check.
func TestGinHelperAborts(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import ("errors"; "github.com/gin-gonic/gin")
type Auth struct{}
func (a *Auth) ok(ctx *gin.Context) bool { return ctx.GetHeader("X-Token") != "" }
func (a *Auth) abort401(ctx *gin.Context) { ctx.AbortWithError(401, errors.New("no")) }
func (a *Auth) evaluateOr401(ctx *gin.Context) { if !a.ok(ctx) { a.abort401(ctx) } }
func (a *Auth) RequireClient(ctx *gin.Context) { a.evaluateOr401(ctx) }
func (a *Auth) LogOnly(ctx *gin.Context) { println(ctx.FullPath()) }
func main() {
	a := &Auth{}
	r := gin.New()
	r.Group("/x", a.RequireClient).GET("/a", h)
	r.Group("/y", a.LogOnly).GET("/b", h)
}`))
	if rt := got["GET /x/a"]; len(rt.Passing) != 0 || len(rt.Enforcing) == 0 {
		t.Errorf("RequireClient: passing=%v enforcing=%v", rt.Passing, rt.Enforcing)
	}
	if rt := got["GET /y/b"]; len(rt.Passing) == 0 {
		t.Errorf("LogOnly hands every request on: passing=%v", rt.Passing)
	}
}

// A 401 at the end of the function, reached when the valid branch didn't
// hand on, is a stop; a gin c.JSON(401) followed by return without Abort
// is not.
func TestPathSensitiveStops(t *testing.T) {
	got := byKey(scanSrc(t, `package main
import ("net/http"; "github.com/go-chi/chi/v5"; "github.com/gin-gonic/gin")
func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if valid(r.Header.Get("Authorization")) { next.ServeHTTP(w, r); return }
		http.Error(w, "no", http.StatusUnauthorized)
	})
}
func GinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" { c.JSON(401, nil); return }
		c.Next()
	}
}
func main() {
	r := chi.NewRouter(); r.Use(RequireAuth); r.Get("/a", h)
	g := gin.New(); g.Use(GinAuth()); g.GET("/b", h)
}`))
	if rt := got["GET /a"]; len(rt.Passing) != 0 {
		t.Errorf("fall-through 401 must enforce: passing=%v", rt.Passing)
	}
	if rt := got["GET /b"]; len(rt.Passing) == 0 {
		t.Errorf("gin return without Abort must pass: enforcing=%v", rt.Enforcing)
	}
}

// A cache's Get(key)/Delete(key) interface is not a router: functions
// taking one aren't walked as registrars (Cloudreve hit the walk budget).
func TestRouterInterfaceShape(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", `package x
type KV interface { Get(key string) ([]byte, error); Delete(key string) error; Put(key string, v []byte) error }
type EchoRouter interface { GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route; POST(path string, h HandlerFunc, m ...MiddlewareFunc) *Route }`, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"KV": false, "EchoRouter": true}
	for _, d := range f.Decls {
		for _, sp := range d.(*ast.GenDecl).Specs {
			ts := sp.(*ast.TypeSpec)
			if got := routerMethods(ts.Type.(*ast.InterfaceType)); got != want[ts.Name.Name] {
				t.Errorf("%s: router=%v", ts.Name.Name, got)
			}
		}
	}
}

// An uncalled generated HandlerFromMux(r) walked from its parameter adds a
// prefix-less copy of every route the real caller already registered.
func TestParamCopyDropped(t *testing.T) {
	routes := scanSrc(t, `package main
import ("net/http"; "github.com/go-chi/chi/v5")
func HandlerFromMux(r chi.Router) http.Handler { return HandlerWithOptions(Options{BaseRouter: r}) }
type Options struct{ BaseRouter chi.Router }
func HandlerWithOptions(o Options) http.Handler {
	r := o.BaseRouter
	if r == nil { r = chi.NewRouter() }
	r.Delete("/users/{id}", h)
	return r
}
func main() {
	m := chi.NewRouter()
	m.Mount("/api", HandlerWithOptions(Options{}))
	http.ListenAndServe(":1", m)
}`)
	if len(routes) != 1 || routes[0].Path != "/api/users/{id}" || routes[0].Partial != "" {
		t.Fatalf("routes = %+v", routes)
	}
}
