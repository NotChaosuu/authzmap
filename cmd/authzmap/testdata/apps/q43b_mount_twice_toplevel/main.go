package main

// EXPECT: one sub-router mounted at /a and /b (both open) -> routes /a/users/{id}, /a/reset, /b/users/{id}, /b/reset (4 routes, no /b/a/ concatenation)
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
	adm := chi.NewRouter()
	adm.Delete("/users/{id}", h)
	adm.Post("/reset", h)
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/me", h)
	r.Mount("/a", adm)
	r.Mount("/b", adm)
	http.ListenAndServe(":8080", r)
}
