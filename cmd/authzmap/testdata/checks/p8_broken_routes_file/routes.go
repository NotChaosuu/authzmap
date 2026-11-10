package main

import "github.com/go-chi/chi/v5"

func register(r chi.Router) {
	r.Delete("/admin/users/{id}", h)
	r.Post("/public/ping", h
}
