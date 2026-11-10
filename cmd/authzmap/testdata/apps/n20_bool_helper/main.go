package main

// EXPECT: enforcing (authorize() writes 401 and returns false; mw returns) -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func authorize(w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Authorization") == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authorize(w, r) {
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
