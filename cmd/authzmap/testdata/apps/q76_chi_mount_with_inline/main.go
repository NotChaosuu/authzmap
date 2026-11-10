package main

// EXPECT: chi r.With(RequireAuth).Route("/admin", ...) -> admin routes authed; a sibling r.Route("/admin2", ...) open -> only /admin2/* flagged
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Route("/admin", func(g chi.Router) { g.Delete("/users/{id}", h) })
	r.Route("/admin2", func(g chi.Router) { g.Delete("/users/{id}", h) })
	http.ListenAndServe(":8080", r)
}
