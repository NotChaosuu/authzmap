//go:build windows

package main

import "github.com/go-chi/chi/v5"

func adminRoutes(r chi.Router) {
	r.With(RequireAuth).Delete("/admin/users/{id}", h)
}
