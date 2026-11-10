package main

// EXPECT: &Auth{next} with a pointer-receiver ServeHTTP that rejects -> enforcing
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Auth struct{ next http.Handler }

func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	a.next.ServeHTTP(w, r)
}
func RequireAuth(next http.Handler) http.Handler { return &Auth{next: next} }
func h(w http.ResponseWriter, r *http.Request)   {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
