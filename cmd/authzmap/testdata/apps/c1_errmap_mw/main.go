package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var pages = map[int]string{http.StatusUnauthorized: "please sign in", http.StatusNotFound: "nope"}

func errorPages(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = map[int]string{http.StatusUnauthorized: "please sign in"}
		next.ServeHTTP(w, r)
	})
}

func RequireAuth(next http.Handler) http.Handler { return next }
func h(w http.ResponseWriter, r *http.Request)   {}

func main() {
	r := chi.NewRouter()
	r.Use(errorPages)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
