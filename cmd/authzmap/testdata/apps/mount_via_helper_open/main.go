package main

// EXPECT: a helper mounts the sub on whatever router it's given; called with an authed and an open router -> open copy flagged
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func attach(parent chi.Router, sub chi.Router) { parent.Mount("/files", sub) }
func main() {
	sub := chi.NewRouter()
	sub.Delete("/{id}", h)
	authed := chi.NewRouter()
	authed.Use(RequireAuth)
	attach(authed, sub)
	open := chi.NewRouter()
	attach(open, sub)
	go http.ListenAndServe(":9090", open)
	http.ListenAndServe(":8080", authed)
}
