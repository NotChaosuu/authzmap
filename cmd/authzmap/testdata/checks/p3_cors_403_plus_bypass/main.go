package main

// EXPECT: a CORS origin 403 (not auth) first, then a receiver bypass, then the real 401 -> conditional MEDIUM (the CORS stop must not count as enforcing)
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

type srv struct{ NoAuth bool }

func (s *srv) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if o := r.Header.Get("Origin"); o != "" && o != "https://app.example" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if s.NoAuth {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	s := &srv{}
	r := chi.NewRouter()
	r.Use(s.Auth)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
