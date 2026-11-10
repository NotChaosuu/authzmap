package main

// EXPECT: generics: register[R Router](r R) called with an open chi router -> /admin routes MEDIUM; or INCOMPLETE; never silent clean
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

type Router interface {
	Get(pattern string, fn http.HandlerFunc)
	Post(pattern string, fn http.HandlerFunc)
	Delete(pattern string, fn http.HandlerFunc)
}

func register[R Router](r R) {
	r.Delete("/admin/users/{id}", h)
	r.Post("/admin/reset", h)
}
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	register(r)
	http.ListenAndServe(":8080", r)
}
