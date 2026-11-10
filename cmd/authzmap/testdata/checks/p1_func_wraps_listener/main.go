package main

// EXPECT: serve(hd) wraps its parameter: http.ListenAndServe(addr, RequireAuth(hd)); the only listener -> quiet
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

func serve(addr string, hd http.Handler) error {
	return http.ListenAndServe(addr, RequireAuth(hd))
}
func main() {
	r := chi.NewRouter()
	r.Delete("/admin/users/{id}", h)
	serve(":8080", r)
}
