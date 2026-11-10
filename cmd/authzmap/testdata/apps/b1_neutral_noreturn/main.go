package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func lookup(r *http.Request) *string { return nil }

func gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if lookup(r) == nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Use(gate)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
