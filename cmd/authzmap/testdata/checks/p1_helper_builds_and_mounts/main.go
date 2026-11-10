package main

// EXPECT: a helper creates the sub-router itself and mounts it; called on an authed parent (/v1) and an open one (/v2) -> /v2 flagged, /v1 quiet
import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func mountAPI(parent chi.Router, prefix string) {
	api := chi.NewRouter()
	api.Delete("/users/{id}", h)
	parent.Mount(prefix, api)
}
func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		mountAPI(g, "/v1")
	})
	mountAPI(r, "/v2")
	http.ListenAndServe(":8080", r)
}
