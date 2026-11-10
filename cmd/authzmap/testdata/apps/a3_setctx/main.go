package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ctxKey struct{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKey{}, r.Header.Get("X-User"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Use(AuthMiddleware)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
