package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

var tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		// Verifier only parses; Authenticator is missing
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Delete("/api/users/{id}", h)
		r.Post("/admin/reset", h)
	})
	http.ListenAndServe(":8080", r)
}
