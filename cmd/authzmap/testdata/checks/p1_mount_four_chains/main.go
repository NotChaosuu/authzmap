package main

// EXPECT: one api router mounted 4 ways: auth, auth via Group, logger-only, and an enforcing in-tree gate -> only /c/* flagged
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

func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Key") != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	api := chi.NewRouter()
	api.Delete("/items/{id}", h)
	r := chi.NewRouter()
	r.With(RequireAuth).Mount("/a", api)
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		g.Mount("/b", api)
	})
	r.With(Logger).Mount("/c", api)
	r.With(Gate).Mount("/d", api)
	http.ListenAndServe(":8080", r)
}
