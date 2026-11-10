package main

// EXPECT: api router mounted at /api/v1 (authed via With) and at /api (open alias) -> /api/* routes flagged MEDIUM
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
func main() {
	api := chi.NewRouter()
	api.Delete("/users/{id}", h)
	api.Post("/admin/reset", h)
	r := chi.NewRouter()
	r.With(RequireAuth).Mount("/api/v1", api)
	r.Mount("/api", api)
	http.ListenAndServe(":8080", r)
}
