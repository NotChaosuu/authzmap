package main

// EXPECT: a debug header sets ok = true, but ok is reset to false before the check -> enforcing, quiet
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
		ok := false
		if r.Header.Get("X-Debug") == "1" {
			ok = true
		}
		ok = false
		if !ok && !valid(r.Header.Get("Authorization")) {
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
