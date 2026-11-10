package main

// EXPECT: router held in a package-level interface variable (var api Router = chi.NewRouter()) -> /admin routes MEDIUM, authed one quiet
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
	Delete(pattern string, fn http.HandlerFunc)
	Post(pattern string, fn http.HandlerFunc)
	With(m ...func(http.Handler) http.Handler) chi.Router
}

var api Router = chi.NewRouter()

func main() {
	api.With(RequireAuth).Delete("/users/{id}", h)
	api.Delete("/admin/users/{id}", h)
	api.Post("/admin/reset", h)
	http.ListenAndServe(":8080", api.(http.Handler))
}
