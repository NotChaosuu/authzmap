package main

// EXPECT: X-Legacy-Client: true swaps the token for a known-good constant -> bypass: conditional, MEDIUM
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
		tok := r.Header.Get("Authorization")
		if r.Header.Get("X-Legacy-Client") == "true" {
			tok = "k"
		}
		if !valid(tok) {
			w.WriteHeader(http.StatusUnauthorized)
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
