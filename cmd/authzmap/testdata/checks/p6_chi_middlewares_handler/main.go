package main

// EXPECT: chi.Chain(RequireAuth).HandlerFunc(h) on one route -> quiet; chi.Chain(Logger).HandlerFunc(h) on /b -> MEDIUM
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

func main() {
	r := chi.NewRouter()
	r.Method("DELETE", "/a/users/{id}", chi.Chain(RequireAuth).HandlerFunc(h))
	r.Method("DELETE", "/b/users/{id}", chi.Chain(Logger).HandlerFunc(h))
	http.ListenAndServe(":8080", r)
}
