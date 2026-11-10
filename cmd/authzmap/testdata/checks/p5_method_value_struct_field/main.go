package main

// EXPECT: an unexported method stored as a struct field func and called through the field -> live -> reported
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

type app struct{ build func() chi.Router }
type srv struct{}

func (srv) openRouter() chi.Router {
	r := chi.NewRouter()
	registerAdmin(r)
	return r
}
func main() {
	s := srv{}
	a := app{build: s.openRouter}
	go http.ListenAndServe(":8443", authedRouter())
	http.ListenAndServe(":8080", a.build())
}
