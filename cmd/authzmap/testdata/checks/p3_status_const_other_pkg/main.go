package main

// EXPECT: status code is a constant from another package (errs.Denied = 403) on a failed credential check, with a receiver bypass -> conditional MEDIUM
import (
	"example.com/p3_status_const_other_pkg/errs"
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

type mw struct{ Insecure bool }

func (m mw) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Insecure && !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(errs.Denied)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(mw{}.Auth)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
