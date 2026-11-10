package main

// EXPECT: never rejects: local check closure writes 401 but its result is ignored, next always runs -> auth-never-rejects MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		check := func() bool {
			if !valid(r.Header.Get("Authorization")) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return false
			}
			return true
		}
		check()
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
