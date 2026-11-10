package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func h(w http.ResponseWriter, r *http.Request) {}

// Register adds the admin routes.
func Register(r chi.Router) { r.Delete("/admin/users/{id}", h) }

func bare() chi.Router { r := chi.NewRouter(); Register(r); return r }
