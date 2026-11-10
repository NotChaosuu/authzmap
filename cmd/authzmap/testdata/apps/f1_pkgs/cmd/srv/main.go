package main

import (
	"example.com/f1/internal/api"
	"example.com/f1/internal/mw"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	s := &api.Server{}
	r := chi.NewRouter()
	r.With(mw.RequireUser).Mount("/api", s.Private())
	r.Route("/pub", s.Public)
	http.ListenAndServe(":8080", r)
}
