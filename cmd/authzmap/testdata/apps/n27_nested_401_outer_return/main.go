package main

// EXPECT: 401 in inner if/else, return in outer block (wakapi shape) -> enforcing -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var redirect = ""

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			if redirect == "" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, redirect, http.StatusFound)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
