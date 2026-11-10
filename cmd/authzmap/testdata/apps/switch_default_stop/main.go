package main

// EXPECT: switch on scheme; only the bearer case hands on -> enforcing
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func check(t string) bool { return t == "s" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch scheme := r.Header.Get("X-Scheme"); scheme {
		case "bearer":
			if check(r.Header.Get("Authorization")) {
				next.ServeHTTP(w, r)
				return
			}
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
