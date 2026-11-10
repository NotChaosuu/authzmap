package main

// EXPECT: enforcing: reads need canRead, writes need canWrite; every method is checked -> no gaps (NOT conditional-auth)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool           { return s == "x" }
func canRead(r *http.Request) bool  { return valid(r.Header.Get("Authorization")) }
func canWrite(r *http.Request) bool { return valid(r.Header.Get("X-Write-Token")) }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !canRead(r) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		if !canWrite(r) {
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
