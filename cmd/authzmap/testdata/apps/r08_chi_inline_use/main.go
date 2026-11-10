package main

// EXPECT: auth applied via r.With(RequireAuth) stored in a var, then routes on the var; plus a route on the original r -> only that one reported
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	authed := r.With(RequireAuth)
	authed.Get("/api/me", h)
	authed.Post("/api/notes", h)
	authed.Delete("/api/notes/{id}", h)
	r.Delete("/api/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
