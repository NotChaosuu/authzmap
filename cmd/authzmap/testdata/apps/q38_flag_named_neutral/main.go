package main

// EXPECT: conditional-auth: a receiver field named Public (not in the switch-word list) skips auth -> conditional-auth MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

type Auth struct{ Public bool }

func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.Public {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	a := &Auth{Public: true}
	r := chi.NewRouter()
	r.Use(a.Handler)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
