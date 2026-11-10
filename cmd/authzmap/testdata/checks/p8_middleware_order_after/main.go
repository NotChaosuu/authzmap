package main

// EXPECT: chi: r.Get registered then r.Use(RequireAuth) panics in chi at runtime; but on a Group, auth wrapping via With after-the-fact doesn't apply: r.With(Logger).Delete(...) under an authed parent Group -> quiet
import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		g.With(Logger).Delete("/admin/users/{id}", h)
		g.Group(func(g2 chi.Router) {
			g2.Use(Logger)
			g2.Post("/admin/reset", h)
		})
	})
	http.ListenAndServe(":8080", r)
}
