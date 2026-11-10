package main

// EXPECT: one sub-router at /v1 and /v2 (authed) and /beta (open) -> only /beta/* flagged; no stacked prefixes
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(s string) bool { return s == "k" }
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
	api.Delete("/items/{id}", h)
	r := chi.NewRouter()
	r.With(RequireAuth).Mount("/v1", api)
	r.With(RequireAuth).Mount("/v2", api)
	r.Mount("/beta", api)
	http.ListenAndServe(":8080", r)
}
