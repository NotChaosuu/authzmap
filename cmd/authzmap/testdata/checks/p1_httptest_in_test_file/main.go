package main

// EXPECT: httptest.NewServer(r) in a _test.go file is not a listener; production serves RequireAuth(r) -> quiet
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

func routes() chi.Router {
	r := chi.NewRouter()
	r.Delete("/admin/users/{id}", h)
	return r
}
func main() { http.ListenAndServe(":8080", RequireAuth(routes())) }
