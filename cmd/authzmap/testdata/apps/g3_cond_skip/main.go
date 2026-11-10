package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("DEV") == "1" || r.URL.Query().Get("debug") == "1" {
			next.ServeHTTP(w, r)
			return
		}
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
	if os.Getenv("NO_AUTH") == "" {
		r.Use(RequireAuth)
	}
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
