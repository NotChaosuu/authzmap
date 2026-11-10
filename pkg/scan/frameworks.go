package scan

import "strings"

// framework describes how one router library registers routes and middleware.
type framework struct {
	name string
	// ctors are package-level functions that return a new router.
	ctors map[string]bool
	// types are the package's router types (pointer or not).
	types map[string]bool
	// verbs maps a per-method registration call to its HTTP method.
	verbs map[string]string
	// requestTime: Use() applies to every route on the router, including ones
	// registered before it (chi panics instead, gorilla and echo's root apply
	// middleware per request). Otherwise middleware is captured at registration.
	requestTime bool
}

var (
	fwChi = &framework{
		name:  "chi",
		ctors: set("NewRouter", "NewMux"),
		types: set("Mux", "Router", "Routes"),
		verbs: map[string]string{
			"Get": "GET", "Post": "POST", "Put": "PUT", "Delete": "DELETE",
			"Patch": "PATCH", "Head": "HEAD", "Options": "OPTIONS",
			"Connect": "CONNECT", "Trace": "TRACE",
			"Handle": "*", "HandleFunc": "*",
		},
		requestTime: true,
	}
	fwGin = &framework{
		name:  "gin",
		ctors: set("New", "Default"),
		types: set("Engine", "RouterGroup", "IRouter", "IRoutes"),
		verbs: map[string]string{
			"GET": "GET", "POST": "POST", "PUT": "PUT", "DELETE": "DELETE",
			"PATCH": "PATCH", "HEAD": "HEAD", "OPTIONS": "OPTIONS", "Any": "*",
		},
	}
	fwEcho = &framework{
		name:  "echo",
		ctors: set("New", "NewWithConfig"),
		types: set("Echo", "Group"),
		verbs: map[string]string{
			"GET": "GET", "POST": "POST", "PUT": "PUT", "DELETE": "DELETE",
			"PATCH": "PATCH", "HEAD": "HEAD", "OPTIONS": "OPTIONS",
			"CONNECT": "CONNECT", "TRACE": "TRACE", "Any": "*",
		},
	}
	fwFiber = &framework{
		name:  "fiber",
		ctors: set("New"),
		types: set("App", "Router", "Group"),
		verbs: map[string]string{
			"Get": "GET", "Post": "POST", "Put": "PUT", "Delete": "DELETE",
			"Patch": "PATCH", "Head": "HEAD", "Options": "OPTIONS",
			"Connect": "CONNECT", "Trace": "TRACE", "All": "*",
		},
	}
	fwGorilla = &framework{
		name:        "gorilla",
		ctors:       set("NewRouter"),
		types:       set("Router"),
		verbs:       map[string]string{"HandleFunc": "*", "Handle": "*"},
		requestTime: true,
	}
	// hertz groups capture middleware like gin's
	fwHertz = &framework{
		name:  "hertz",
		ctors: set("Default", "New"),
		types: set("Hertz", "RouterGroup", "Engine"),
		verbs: map[string]string{
			"GET": "GET", "POST": "POST", "PUT": "PUT", "DELETE": "DELETE",
			"PATCH": "PATCH", "HEAD": "HEAD", "OPTIONS": "OPTIONS", "Any": "*",
		},
	}
	fwNetHTTP = &framework{
		name:  "net/http",
		ctors: set("NewServeMux"),
		types: set("ServeMux"),
		verbs: map[string]string{"HandleFunc": "*", "Handle": "*"},
	}
)

var frameworks = []*framework{fwChi, fwGin, fwEcho, fwFiber, fwGorilla, fwHertz, fwNetHTTP}

// frameworkForImport maps an import path to the router library it provides.
// Subpackages (chi/v5/middleware, echo/v4/middleware) are not the router.
func frameworkForImport(p string) *framework {
	if p == "net/http" {
		return fwNetHTTP
	}
	if i := strings.LastIndexByte(p, '/'); i > 0 && isVersionSuffix(p[i+1:]) {
		p = p[:i]
	}
	switch p {
	case "github.com/go-chi/chi":
		return fwChi
	case "github.com/gin-gonic/gin":
		return fwGin
	case "github.com/labstack/echo":
		return fwEcho
	case "github.com/gofiber/fiber":
		return fwFiber
	case "github.com/gorilla/mux":
		return fwGorilla
	case "github.com/cloudwego/hertz/pkg/app/server", "github.com/cloudwego/hertz/pkg/route":
		return fwHertz
	}
	return nil
}

func isVersionSuffix(s string) bool {
	if len(s) < 2 || s[0] != 'v' {
		return false
	}
	for _, r := range s[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Echo's root router applies Use() per request; its groups capture
// middleware when a route is added. Gin and fiber always capture.
func requestTime(n *node) bool {
	if n.fw == fwEcho {
		return n.kind == kindRoot
	}
	return n.fw.requestTime
}

// handlerAdapters convert between handler types without adding behaviour, so
// they're unwrapped rather than listed as middleware.
var handlerAdapters = set(
	"http.HandlerFunc", "HandlerFunc", "http.Handler",
	"gin.WrapH", "gin.WrapF", "echo.WrapHandler", "echo.WrapMiddleware",
	"adaptor.HTTPHandler", "adaptor.HTTPHandlerFunc", "adaptor.HTTPMiddleware",
	"mux.MiddlewareFunc", "MiddlewareFunc",
	"adaptor.FiberApp", "adaptor.FiberHandler", "adaptor.FiberHandlerFunc", "negroni.Wrap", "negroni.WrapFunc",
)

// frameworkListen: methods with which a framework serves its own engine
// (r.Run(":8080"), e.Start, app.Listen, h.Spin): an open listener.
var frameworkListen = set("Run", "RunTLS", "RunUnix", "RunListener", "RunFd", "Start", "StartTLS", "StartServer",
	"StartAutoTLS", "StartH2CServer", "Listen", "ListenTLS", "ListenMutualTLS", "Listener", "Spin")

// serveSinks take the handler that serves requests as their last argument.
var serveSinks = set(
	"http.ListenAndServe", "http.ListenAndServeTLS", "http.Serve", "http.ServeTLS",
	"httptest.NewServer", "httptest.NewTLSServer", "httptest.NewUnstartedServer",
	"h2c.NewHandler",
)

// frameworkMethod reports a method the router library itself defines, as
// opposed to one on a wrapper type that embeds the router.
func frameworkMethod(fw *framework, m string) bool {
	if _, ok := fw.verbs[m]; ok {
		return true
	}
	switch m {
	case "Use", "Group", "Route", "Mount", "With", "Handle", "HandleFunc", "Method", "MethodFunc",
		"Match", "Add", "Pre", "Static", "StaticFS", "StaticFile", "File", "NoRoute", "NoMethod",
		"NotFound", "MethodNotAllowed", "Any", "All", "PathPrefix", "Path", "Methods", "Subrouter",
		"NewRoute", "Host", "Schemes", "Headers", "Queries", "Name", "ServeHTTP", "Routes", "Middlewares",
		"RouteNotFound", "Listen", "Start", "Run", "Handler", "HandlerFunc", "Spin":
		return true
	}
	return false
}

func set(xs ...string) map[string]bool {
	m := make(map[string]bool, len(xs))
	for _, x := range xs {
		m[x] = true
	}
	return m
}
