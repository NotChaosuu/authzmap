package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request)                       {}
func hr(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

const adminDelete = "/admin/users/:id"

func main() {
	r := chi.NewRouter()
	r.With(RequireAuth).Get("/api/me", h)
	adm := httprouter.New()
	adm.DELETE(adminDelete, hr)
	adm.Handle("POST", "/admin/reset", hr)
	r.Mount("/internal", adm)
	http.ListenAndServe(":8080", r)
}
