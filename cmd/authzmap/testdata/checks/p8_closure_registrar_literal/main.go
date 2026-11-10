package main

// EXPECT: local closure reg := func(p string) { r.Delete(p, h) }; reg("/admin/keys") on an open router -> flagged with a real path
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
	r := chi.NewRouter()
	r.With(RequireAuth).Delete("/admin/users/{id}", h)
	reg := func(p string) { r.Delete(p, h) }
	reg("/admin/keys/{id}")
	http.ListenAndServe(":8080", r)
}
