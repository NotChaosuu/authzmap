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
	r.With(RequireAuth) // result discarded: no effect
	r.Delete("/api/users/{id}", h)
	r.Route("/admin", func(ar chi.Router) {
		ar.Use(RequireAuth)
		ar.Get("/stats", h)
		r.Delete("/admin/purge/{id}", h) // outer router: no auth
	})
	g := r.Group(nil)
	g.Use(RequireAuth)
	g.Get("/api/me", h)
	sub := chi.NewRouter()
	sub.Post("/settings/{id}", h)
	r.With(RequireAuth).Get("/x", h)
	r.Mount("/api/admin", sub)
	http.ListenAndServe(":8080", r)
}
