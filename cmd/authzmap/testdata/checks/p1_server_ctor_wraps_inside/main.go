package main

// EXPECT: the only server is built by a ctor that wraps its handler in RequireAuth -> quiet
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

func newServer(hd http.Handler) *http.Server {
	return &http.Server{Addr: ":8080", Handler: RequireAuth(hd)}
}
func main() {
	r := chi.NewRouter()
	r.Delete("/admin/users/{id}", h)
	srv := newServer(r)
	srv.ListenAndServe()
}
