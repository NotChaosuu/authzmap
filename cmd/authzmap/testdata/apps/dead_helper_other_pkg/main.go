package main

// EXPECT: an uncalled unexported helper in another package builds the routes open; main serves them authed -> quiet
import (
	"net/http"

	"example.com/dead_helper_other_pkg/routes"
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
	r := chi.NewRouter()
	r.Group(func(g chi.Router) { g.Use(RequireAuth); routes.Register(g) })
	http.ListenAndServe(":8080", r)
}
