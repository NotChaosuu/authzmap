package main

// EXPECT: an unexported method reached only through an interface call registers the open copy -> live -> reported
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

func registerAdmin(r chi.Router) { r.Delete("/admin/users/{id}", h) }
func authedRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	registerAdmin(r)
	return r
}

type mounter interface{ mount() chi.Router }
type legacy struct{}

func (legacy) mount() chi.Router {
	r := chi.NewRouter()
	registerAdmin(r)
	return r
}
func run(m mounter) { http.ListenAndServe(":8080", m.mount()) }
func main() {
	go http.ListenAndServe(":8443", authedRouter())
	run(legacy{})
}
