package main

// EXPECT: receiver field Lax (no switch word) skips auth -> conditional MEDIUM
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

type guard struct{ Lax bool }

func (g guard) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.Lax {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	g := guard{Lax: true}
	r := chi.NewRouter()
	r.Use(g.Wrap)
	r.Delete("/api/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
