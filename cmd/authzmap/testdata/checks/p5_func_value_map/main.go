package main

// EXPECT: the open builder is only referenced as a value in a map (never called by name); it is live -> the open copy must be reported
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

func openRouter() chi.Router {
	r := chi.NewRouter()
	registerAdmin(r)
	return r
}

var builders = map[string]func() chi.Router{"secure": authedRouter, "legacy": openRouter}

func main() {
	go http.ListenAndServe(":8443", builders["secure"]())
	http.ListenAndServe(":8080", builders["legacy"]())
}
