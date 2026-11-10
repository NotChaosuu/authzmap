package main

// EXPECT: `var skip bool` set when ?internal=1, then `if !skip && !valid(...)` -> the parameter skips the check: conditional, MEDIUM
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
		var skip bool
		if r.URL.Query().Get("internal") == "1" {
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
