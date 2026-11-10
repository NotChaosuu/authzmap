package main

// EXPECT: a package func authEnabled() reading os.Getenv gates the check -> conditional MEDIUM
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

func authEnabled() bool { return os.Getenv("AUTH") != "off" }
func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled() {
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
	r := chi.NewRouter()
	r.Use(Gate)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
