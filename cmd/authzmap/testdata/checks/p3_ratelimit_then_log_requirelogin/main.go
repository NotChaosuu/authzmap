package main

// EXPECT: RequireLogin's only stop is a per-IP 429 limiter -> never rejects on auth -> MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
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

var _ = os.Getenv

type L struct{}

func (L) Allow(k string) bool { return true }

var lim L

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lim.Allow(r.RemoteAddr) {
			w.WriteHeader(429)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(RequireLogin)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
