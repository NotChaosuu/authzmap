package main

// EXPECT: enforcing (401 at function end, no explicit return) -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if valid(r.Header.Get("Authorization")) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
