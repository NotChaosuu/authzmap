package main

// EXPECT: sub-router mounted open at /internal FIRST, then inside an authed group at /admin -> /internal/* must be flagged (silent miss if the authed mount wins)
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
	r.Mount("/internal", adm)
	r.Group(func(g chi.Router) { g.Use(RequireAuth); g.Mount("/admin", adm) })
	http.ListenAndServe(":8080", r)
}
