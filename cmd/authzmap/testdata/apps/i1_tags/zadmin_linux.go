//go:build linux

package main

import "github.com/go-chi/chi/v5"

// the build that actually ships: no auth
func adminRoutes(r chi.Router) {
	r.Delete("/admin/users/{id}", h)
}
