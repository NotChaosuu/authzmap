package main

// EXPECT: never rejects: break outer after 401 leaves the loop and falls into next -> auth-never-rejects MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	outer:
		for _, v := range r.Header.Values("Authorization") {
			switch {
			case !valid(v):
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				break outer
			}
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
