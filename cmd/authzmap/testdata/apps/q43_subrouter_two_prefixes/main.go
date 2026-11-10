package main

// EXPECT: the same admin sub-router mounted at /admin (authed group) and at /internal (open) -> /internal/* flagged, /admin/* quiet
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
	r.Group(func(g chi.Router) { g.Use(RequireAuth); g.Mount("/admin", adm) })
	r.Mount("/internal", adm)
	http.ListenAndServe(":8080", r)
}
