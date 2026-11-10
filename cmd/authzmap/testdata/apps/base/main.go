package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "no", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Get("/healthz", h)
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/users", h)
		r.Delete("/api/users/{id}", h)
	})
	r.Delete("/api/projects/{id}", h)
	http.ListenAndServe(":8080", r)
}
