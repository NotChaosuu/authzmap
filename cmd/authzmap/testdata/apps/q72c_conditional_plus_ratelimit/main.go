package main

// EXPECT: config-gated auth plus a 429 rate-limit stop -> conditional-auth MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var AuthToken string

func allow(r *http.Request) bool { return true }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allow(r) {
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		if AuthToken != "" && r.Header.Get("Authorization") != "Bearer "+AuthToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
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
