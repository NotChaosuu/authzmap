package main

// EXPECT: receiver `logger` nil-check (logging only) before a real 401 -> enforcing, quiet (no false conditional)
import (
	"github.com/go-chi/chi/v5"
	"log"
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

type Auth struct{ logger *log.Logger }

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.logger != nil {
			a.logger.Println("auth", r.URL.Path)
		}
		if !valid(r.Header.Get("Authorization")) {
			if a.logger != nil {
				a.logger.Println("denied")
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	a := &Auth{}
	r := chi.NewRouter()
	r.Use(a.Require)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
