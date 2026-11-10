package main

// EXPECT: JWTAuth's only stop is a 413 body-size limit; the token is never checked -> MEDIUM
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

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 1<<20 {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		_ = r.Header.Get("Authorization")
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(JWTAuth)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
