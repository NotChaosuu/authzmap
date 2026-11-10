package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func h(w http.ResponseWriter, r *http.Request) {}

func Register(r chi.Router) { r.Delete("/admin/users/{id}", h) }
