package main

// EXPECT: else-if chain where only the valid branch hands on -> enforcing, no gaps
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(t string) bool   { return t == "s" }
func expired(t string) bool { return t == "old" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		if valid(tok) {
			next.ServeHTTP(w, r)
		} else if expired(tok) {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
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
