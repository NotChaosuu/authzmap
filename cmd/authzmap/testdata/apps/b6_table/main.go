package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler { return next }
func h(w http.ResponseWriter, r *http.Request)   {}

type route struct {
	method, path string
	fn           http.HandlerFunc
	auth         bool
}

var routes = []route{
	{"GET", "/api/users", h, true},
	{"DELETE", "/api/users/{id}", h, false},
	{"POST", "/admin/reset", h, false},
}

func main() {
	r := chi.NewRouter()
	for _, rt := range routes {
		if rt.auth {
			r.With(RequireAuth).Method(rt.method, rt.path, rt.fn)
		} else {
			r.Method(rt.method, rt.path, rt.fn)
		}
	}
	m := map[string]http.HandlerFunc{"/admin/drop": h, "/admin/users": h}
	for p, f := range m {
		r.Post(p, f)
	}
	http.ListenAndServe(":8080", r)
}
