package main

// EXPECT: deleteUser on an authed /api route and on an open /legacy route -> /legacy DELETE flagged (sibling or sensitive) MEDIUM+
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request)          {}
func deleteUser(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Route("/api", func(g chi.Router) { g.Use(RequireAuth); g.Delete("/users/{id}", deleteUser); g.Get("/users", h) })
	r.Delete("/legacy/users/{id}", deleteUser)
	http.ListenAndServe(":8080", r)
}
