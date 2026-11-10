package aghhttp

import "net/http"

type Registrar interface {
	Register(method, path string, h http.HandlerFunc)
}

type WrapFunc func(method string, h http.HandlerFunc) http.Handler

type DefaultRegistrar struct {
	mux    *http.ServeMux
	wrapFn WrapFunc
}

func NewDefaultRegistrar(mux *http.ServeMux, wrap WrapFunc) *DefaultRegistrar {
	return &DefaultRegistrar{mux: mux, wrapFn: wrap}
}

func (r *DefaultRegistrar) Register(method, path string, h http.HandlerFunc) {
	r.mux.Handle(path, r.wrapFn(method, h))
}
