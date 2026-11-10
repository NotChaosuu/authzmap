package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// public by design
	r.Post("/login", login)

	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/me", me)
		r.Get("/orders", orders)
	})

	// /stats is missing the auth its sibling has.
	r.Route("/admin", func(r chi.Router) {
		r.Get("/stats", stats)
		r.With(RequireAuth).Post("/users", createUser)
	})

	http.ListenAndServe(":8080", r)
}
