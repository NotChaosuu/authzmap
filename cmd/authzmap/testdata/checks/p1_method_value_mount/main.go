package main

// EXPECT: m := r.Mount; m("/open", api) alongside an authed mount -> /open flagged
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

func main() {
	api := chi.NewRouter()
	api.Delete("/items/{id}", h)
	r := chi.NewRouter()
	r.With(RequireAuth).Mount("/secure", api)
	m := r.Mount
	m("/open", api)
	http.ListenAndServe(":8080", r)
}
