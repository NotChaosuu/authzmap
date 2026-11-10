package main

// EXPECT: Authenticate answers CORS preflight (OPTIONS -> 204) and otherwise only logs the token and hands on -> never rejects on auth -> MEDIUM
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

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		println(r.Header.Get("Authorization"))
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(Authenticate)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
