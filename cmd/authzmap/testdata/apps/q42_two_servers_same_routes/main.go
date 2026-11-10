package main

// EXPECT: same route function mounted on an authed admin server and an open public server -> /api/users/{id} DELETE flagged once (weakest chain), not twice
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
func h(w http.ResponseWriter, r *http.Request) {}
func apiRoutes(r chi.Router) {
	r.Delete("/users/{id}", h)
	r.Post("/settings", h)
}
func main() {
	admin := chi.NewRouter()
	admin.Use(RequireAuth)
	admin.Route("/api", apiRoutes)
	public := chi.NewRouter()
	public.Route("/api", apiRoutes)
	go http.ListenAndServe(":9090", admin)
	http.ListenAndServe(":8080", public)
}
