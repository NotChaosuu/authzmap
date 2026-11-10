package main

// EXPECT: every route below except the GET /static, GET /docs, GET /api/public, POST /login, GET /health is NOT public and should be reported (all unauthed, siblings authed)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/me", h)
		r.Post("/api/notes", h)
		r.Delete("/api/notes/{id}", h)
		r.Put("/api/settings", h)
	})
	// genuinely public
	r.Get("/static/app.js", h)
	r.Get("/docs/index.html", h)
	r.Get("/api/public/stats", h)
	r.Post("/login", h)
	r.Get("/health", h)
	// should NOT be public
	r.Put("/api/public/settings", h)
	r.Patch("/api/public/users/{id}", h)
	r.Delete("/api/public/users/{id}", h)
	r.Post("/static/{file}", h)
	r.Put("/static/{file}", h)
	r.Delete("/docs/{id}", h)
	r.Post("/health/reset", h)
	r.Delete("/health", h)
	r.Get("/Admin/users", h)
	r.Get("/ADMIN/export", h)
	r.Get("/admin/", h)
	r.Get("/api/public/../admin/users", h)
	r.Get("/static/../admin/dump", h)
	r.Get("/api/%61dmin/users", h)
	r.Get("/login/../admin/users", h)
	r.Post("/api/v1/users/{id}/password/reset", h)
	r.Post("/api/v1/admin/users/register", h)
	r.Post("/api/v1/login/admin/promote", h)
	r.Get("/metrics", h)
	r.Get("/debug/pprof/*", h)
	r.Get("/api/users/{id}/token", h)
	r.Post("/oauth/token/revoke-all", h)
	r.Post("/webhooks/github/../../admin/x", h)
	http.ListenAndServe(":8080", r)
}
