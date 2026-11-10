package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func validate(t string) error { return nil }

// RequireAuth logs a bad token but never stops the request.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := validate(r.Header.Get("Authorization")); err != nil {
			log.Println("bad token", err)
		}
		next.ServeHTTP(w, r)
	})
}

func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
