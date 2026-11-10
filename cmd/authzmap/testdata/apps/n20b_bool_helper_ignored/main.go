package main

// EXPECT: authorize() result ignored (no return) -> next runs -> flagged
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
		authorize(w, r)
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
