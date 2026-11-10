package main

// EXPECT: an uncalled helper named only in a comment also builds the authed routes open -> quiet
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

// legacyRouter used to be served on :9000.
func legacyRouter() chi.Router { r := chi.NewRouter(); routes(r); return r }
func routes(r chi.Router)      { r.Delete("/api/users/{id}", h) }
func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) { g.Use(RequireAuth); routes(g) })
	http.ListenAndServe(":8080", r)
}
