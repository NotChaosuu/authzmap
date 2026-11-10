package main

// EXPECT: enforcing negroni-style middleware (Handler delegates to m.ServeHTTP(w, r, next)) -> no gaps; must NOT be never-rejects
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AuthenticateMiddleware struct{}

func (m *AuthenticateMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { m.ServeHTTP(w, r, h.ServeHTTP) })
}
func (m *AuthenticateMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	next(w, r)
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	auth := &AuthenticateMiddleware{}
	r := chi.NewRouter()
	r.Use(auth.Handler)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
