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

func register(r chi.Router, depth int) {
	r.Delete("/items/{id}", h)
	if depth > 0 {
		r.Route("/sub", func(r chi.Router) { register(r, depth-1) })
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/healthz", h)
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		register(r, 2)
	})
	http.ListenAndServe(":8080", r)
}
