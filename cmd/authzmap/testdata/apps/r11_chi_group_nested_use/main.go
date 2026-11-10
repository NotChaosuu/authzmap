package main

// EXPECT: auth Use in an inner Group does not cover outer siblings; outer DELETE reported
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
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) { r.Use(RequireAuth); r.Get("/me", h); r.Post("/notes", h) })
		r.Delete("/notes/{id}", h)
		r.Route("/admin", func(r chi.Router) { r.Post("/reset", h) })
	})
	http.ListenAndServe(":8080", r)
}
