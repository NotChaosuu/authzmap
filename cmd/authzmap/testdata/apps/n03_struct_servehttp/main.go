package main

// EXPECT: Auth.ServeHTTP never rejects -> flagged
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Auth struct{ next http.Handler }

func (a Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		log.Println("unauthenticated")
	}
	a.next.ServeHTTP(w, r)
}
func RequireAuth(next http.Handler) http.Handler { return Auth{next: next} }
func h(w http.ResponseWriter, r *http.Request)   {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
