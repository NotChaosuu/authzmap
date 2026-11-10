package main

import (
	"github.com/go-chi/chi/v5"
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
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	pub := chi.NewRouter()
	pub.Get("/", h)
	adm := chi.NewRouter()
	adm.Delete("/admin/users/{id}", h)
	adm.Post("/admin/reset", h)
	wrapped := RequireAuth(pub) // auth wraps the wrong router
	_ = wrapped
	srv := &http.Server{Addr: ":8081", Handler: RequireAuth(pub)}
	go srv.ListenAndServe()
	http.ListenAndServe(":9090", adm)
}
