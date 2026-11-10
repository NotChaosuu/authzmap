package main

// EXPECT: (*AuthMW).Handler never rejects -> flagged
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type AuthMW struct{ secret string }

func (m *AuthMW) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != m.secret {
			log.Println("bad auth")
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	auth := &AuthMW{secret: "x"}
	r := chi.NewRouter()
	r.Use(auth.Handler)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
