package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var ErrForbidden = errors.New("forbidden")

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrForbidden):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(500)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := store(chi.URLParam(r, "id")); err != nil {
		writeErr(w, err)
		return
	}
}
func store(id string) error { return nil }

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/projects", h)
		r.Get("/api/projects/{id}", h)
		r.Put("/api/projects/{id}", h)
	})
	r.Delete("/api/projects/{id}", deleteProject)
	r.Delete("/admin/projects/{id}", deleteProject)
	http.ListenAndServe(":8080", r)
}
