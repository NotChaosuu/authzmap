package main

// EXPECT: optional auth (no Authorization header -> passes) is a bypass -> conditional-auth MEDIUM (checks credentials only when sent)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		if len(tok) > 0 {
			if !valid(tok) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
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
