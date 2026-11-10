// Package chi is a minimal stub of go-chi/chi/v5's surface, enough for the typed
// scanner's fixture to type-check offline.
package chi

import "net/http"

type Router interface {
	Get(pattern string, h http.HandlerFunc)
	Post(pattern string, h http.HandlerFunc)
	Use(mw ...func(http.Handler) http.Handler)
}

type Mux struct{}

func NewRouter() *Mux { return &Mux{} }

func (m *Mux) Get(pattern string, h http.HandlerFunc)    {}
func (m *Mux) Post(pattern string, h http.HandlerFunc)   {}
func (m *Mux) Use(mw ...func(http.Handler) http.Handler) {}
