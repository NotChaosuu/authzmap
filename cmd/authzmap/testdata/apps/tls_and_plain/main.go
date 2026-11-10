package main

// EXPECT: authed TLS listener, open plain listener for the same router -> flagged
import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := chi.NewRouter()
	r.Delete("/api/users/{id}", h)
	go http.ListenAndServeTLS(":443", "c.pem", "k.pem", RequireAuth(r))
	http.ListenAndServe(":80", r)
}
