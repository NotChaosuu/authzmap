package main

// EXPECT: go-pkgz m.Auth = auth; m.Trace = optional (weak) -> /trace flagged; authboss LoadClientStateMiddleware not auth -> /ab flagged; authboss.Middleware2 and casbin Authorizer = auth
import (
	"github.com/casbin/casbin/v2"
	authz "github.com/casbin/chi-authz"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/auth"
	"github.com/volatiletech/authboss/v3"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	svc := auth.NewService(auth.Opts{})
	m := svc.Middleware()
	ab := authboss.New()
	e, _ := casbin.NewEnforcer("model.conf", "policy.csv")
	r := chi.NewRouter()
	r.Group(func(r chi.Router) { r.Use(m.Auth); r.Delete("/api/users/{id}", h); r.Post("/api/notes", h) })
	r.Group(func(r chi.Router) { r.Use(m.Trace); r.Delete("/trace/users/{id}", h); r.Post("/trace/admin/reset", h) })
	r.Group(func(r chi.Router) {
		r.Use(ab.LoadClientStateMiddleware)
		r.Delete("/ab/users/{id}", h)
		r.Post("/ab/admin/reset", h)
	})
	r.Group(func(r chi.Router) {
		r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
		r.Delete("/abok/users/{id}", h)
		r.Post("/abok/notes", h)
	})
	r.Group(func(r chi.Router) {
		r.Use(authz.Authorizer(e))
		r.Delete("/cas/users/{id}", h)
		r.Post("/cas/notes", h)
	})
	http.ListenAndServe(":8080", r)
}
