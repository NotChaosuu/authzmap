package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Get("/login", h)
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/users", h)
		r.Get("/api/users/{id}", h)
		r.Put("/api/users/{id}", h)
	})
	// all unauthenticated:
	r.Post("/api/projects/{id}/public", h)
	r.Delete("/api/public/users/{id}", h)
	r.Delete("/assets/{id}", h)
	r.Post("/static/upload", h)
	r.Delete("/docs/page", h)
	r.Delete("/webhooks/{id}", h)
	r.Post("/hooks/{id}/delete", h)
	r.Post("/password/change", h)
	r.Post("/api/nodes/register", h)
	r.Post("/api/users/register", h)
	r.Delete("/api/sessions", h)
	r.Get("/api/users/export/callback", h)
	r.Post("/api/jobs/{id}/callback", h)
	r.Get("/api/info", h)
	r.Get("/api/status", h)
	r.Post("/api/token", h)
	r.Post("/token", h)
	r.Post("/api/v1/users/{id}/files/{fileHash}", h)
	r.Delete("/api/repos/{code}", h)
	r.Delete("/api/countries/{countryCode}", h)
	r.Post("/oauth/apps", h)
	r.Delete("/oauth/{appID}", h)
	r.Post("/sso/{id}/delete", h)
	r.Delete("/dist/{id}", h)
	r.Post("/api/swagger/config", h)
	r.Post("/swagger/config", h)
	r.Post("/.well-known/acme-reset", h)
	r.Delete("/api/users/{id}/openapi.json", h)
	r.Delete("/forgot-everything", h)
	r.Post("/reset-password/admin-override", h)
	r.Handle("/api/*", h2)
	r.Post("/settings/token", h)
	r.Delete("/{user}/repos/{id}", h)
	r.Delete("/api/users/{id}/sessions/{sid}", h)
	http.ListenAndServe(":8080", r)
}

var h2 = http.HandlerFunc(h)
