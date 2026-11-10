package main

// EXPECT: registrar interface embedding another (Reader{Get} + Writer{Post,Delete}); open router -> /admin routes MEDIUM
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

type Reader interface {
	Get(pattern string, fn http.HandlerFunc)
}
type Writer interface {
	Post(pattern string, fn http.HandlerFunc)
	Delete(pattern string, fn http.HandlerFunc)
}
type RW interface {
	Reader
	Writer
}

func adminRoutes(r RW) {
	r.Delete("/admin/users/{id}", h)
	r.Post("/admin/reset", h)
}
func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	adminRoutes(r)
	http.ListenAndServe(":8080", r)
}
