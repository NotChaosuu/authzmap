package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct{ deps int }

func h(w http.ResponseWriter, r *http.Request) {}

type Registrar[T any] struct{ r chi.Router }

func (g Registrar[T]) Add(p string, fn http.HandlerFunc) { g.r.Post(p, fn) }

func (s *Server) Private() http.Handler {
	r := chi.NewRouter()
	r.Get("/users", h)
	r.Delete("/users/{id}", h)
	return r
}

func (s *Server) Public(r chi.Router) {
	r.Get("/feed", h)
	Registrar[int]{r: r}.Add("/users/{id}/ban", h)
}
