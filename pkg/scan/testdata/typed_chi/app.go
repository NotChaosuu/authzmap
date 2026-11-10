package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes hands its router to a helper; only --types follows it.
func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(Logger)
	r.Get("/health", health)
	registerAdmin(r)
	return r
}

func registerAdmin(r chi.Router) {
	r.Get("/admin/stats", stats)
}

func Logger(next http.Handler) http.Handler         { return next }
func health(w http.ResponseWriter, r *http.Request) {}
func stats(w http.ResponseWriter, r *http.Request)  {}
