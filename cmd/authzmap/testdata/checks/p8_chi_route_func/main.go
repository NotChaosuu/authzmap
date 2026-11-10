package main

// EXPECT: r.Route("/admin") with Use(RequireAuth) inside, r.Route("/ops") without -> /ops flagged
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.Use(RequireAuth)
		r.Delete("/users/{id}", h)
	})
	r.Route("/ops", func(r chi.Router) {
		r.Post("/restart", h)
	})
	http.ListenAndServe(":8080", r)
}
