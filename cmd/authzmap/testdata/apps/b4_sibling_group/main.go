package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler { return next }
func h(w http.ResponseWriter, r *http.Request)   {}

func main() {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/users", h)
	})
	r.Route("/api/v2", func(r chi.Router) {
		r.Get("/users", h)
		r.Post("/users", h)
		r.Get("/users/export", h)
	})
	http.ListenAndServe(":8080", r)
}
