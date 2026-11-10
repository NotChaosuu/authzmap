package main

// EXPECT: reads need the read scope, writes the write scope -> quiet
import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
func h(w http.ResponseWriter, r *http.Request) {}

func hasScope(r *http.Request, s string) bool { return r.Header.Get("Authorization") == s }
func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !hasScope(r, "read") {
				w.WriteHeader(401)
				return
			}
		} else if !hasScope(r, "write") {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(Gate)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
