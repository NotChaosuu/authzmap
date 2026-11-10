package main

// EXPECT: enforcing: config picks OIDC or basic auth; both reject -> no gaps (NOT conditional-auth)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

type Config struct{ UseOIDC bool }

var cfg Config

func oidcOK(r *http.Request) bool  { return valid(r.Header.Get("Authorization")) }
func basicOK(r *http.Request) bool { u, p, ok := r.BasicAuth(); return ok && u == "a" && p == "b" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.UseOIDC {
			if !oidcOK(r) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		if !basicOK(r) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
