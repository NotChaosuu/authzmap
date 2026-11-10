package main

// EXPECT: the 401 code is held in a local variable -> enforcing; quiet. Second middleware RequireToken with const status and a ?debug bypass -> conditional MEDIUM on /b
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

const denied = http.StatusUnauthorized

func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusUnauthorized
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(code)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("debug") == "1" {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(denied)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.With(Gate).Delete("/a/users/{id}", h)
	r.With(RequireToken).Delete("/b/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
