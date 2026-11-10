package main

// EXPECT: RequireRole("admin") never rejects -> /admin flagged; RequireAuth enforces for /api
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func roleOf(r *http.Request) string { return r.Header.Get("X-Role") }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if roleOf(r) != role {
				log.Printf("role mismatch want %s", role)
			}
			next.ServeHTTP(w, r)
		})
	}
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.With(RequireRole("admin")).Delete("/admin/users/{id}", h)
	r.With(RequireRole("admin")).Post("/admin/reset", h)
	r.With(RequireAuth).Get("/api/me", h)
	r.With(RequireAuth).Post("/api/notes", h)
	http.ListenAndServe(":8080", r)
}
