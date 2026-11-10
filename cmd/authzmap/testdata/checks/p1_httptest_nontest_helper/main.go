package main

// EXPECT: httptest.NewServer(routes()) in a non-test helper file (testutil.go) is a test server, not a listener -> quiet (or at most LOW)
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
