package main

// EXPECT: receiver cache keyed by the presented token: hit -> next (already authenticated), miss -> verify or 401 -> enforcing, quiet
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

type Auth struct{ cache map[string]bool }

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		if a.cache[tok] {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		a.cache[tok] = true
		next.ServeHTTP(w, r)
	})
}
func main() {
	a := &Auth{cache: map[string]bool{}}
	r := chi.NewRouter()
	r.Use(a.Require)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
