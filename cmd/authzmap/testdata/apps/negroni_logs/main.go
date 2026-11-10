package main

// EXPECT: negroni-style delegate whose ServeHTTP only logs -> never rejects
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthMiddleware struct{}

func (m *AuthMiddleware) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { m.ServeHTTP(w, r, h.ServeHTTP) })
}
func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Header.Get("Authorization") == "" {
		println("anonymous")
	}
	next(w, r)
}

var RequireAuth = (&AuthMiddleware{}).Wrap

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
