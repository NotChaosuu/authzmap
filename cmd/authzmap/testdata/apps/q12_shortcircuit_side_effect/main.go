package main

// EXPECT: never rejects: deny() writes 401 inside || short-circuit, next always runs -> auth-never-rejects MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func deny(w http.ResponseWriter) bool {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
	return true
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = valid(r.Header.Get("Authorization")) || deny(w)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
