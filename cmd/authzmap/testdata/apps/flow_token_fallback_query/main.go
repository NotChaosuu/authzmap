package main

// EXPECT: token falls back to ?token= when the header is empty; both are credentials -> enforcing, quiet
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
		if tok == "" {
			tok = r.URL.Query().Get("token")
		}
		if !valid(tok) {
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
