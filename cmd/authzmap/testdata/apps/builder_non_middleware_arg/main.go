package main

// EXPECT: withRole(h, role) takes a value, not middleware, so it isn't a builder: the route keeps withRole as its own entry
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}

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

// Wrap applies mws to h, first one outermost.
func Wrap(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func withRole(h http.Handler, role string) http.Handler { return h }

func main() {
	r := chi.NewRouter()
	r.Handle("/api/users/{id}", withRole(http.HandlerFunc(h), "admin"))
	http.ListenAndServe(":8080", r)
}
