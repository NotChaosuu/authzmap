package main

// EXPECT: open builder called only via `go` statement -> live -> reported
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

func registerAdmin(r chi.Router) { r.Delete("/admin/users/{id}", h) }
func authedRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	registerAdmin(r)
	return r
}

func openListener() {
	r := chi.NewRouter()
	registerAdmin(r)
	http.ListenAndServe(":8080", r)
}
func main() {
	go openListener()
	http.ListenAndServe(":8443", authedRouter())
}
