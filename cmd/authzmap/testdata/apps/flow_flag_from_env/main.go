package main

// EXPECT: skip set from an env var (a setting) in a branch, then `if !skip && !valid` -> rejects only when a setting is on: conditional, MEDIUM
import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

var _ = os.Getenv

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}

func RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		skip := false
		if os.Getenv("DEV") != "" {
			skip = true
		}
		if !skip && !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(RequireToken)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
