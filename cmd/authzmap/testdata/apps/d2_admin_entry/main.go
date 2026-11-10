package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler { return next }
func h(w http.ResponseWriter, r *http.Request)   {}

func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	r.Post("/admin/impersonate/login", h)
	r.Post("/admin/jobs/callback", h)
	r.Post("/api/admin/sync/callback", h)
	r.Post("/internal/users/signin", h)
	r.Post("/api/v1/register", h)
	r.Post("/api/devices/register", h)
	r.Get("/admin/static/*", h)
	r.Post("/admin/assets/*", h)
	r.Get("/api/admin/status", h)
	http.ListenAndServe(":8080", r)
}
