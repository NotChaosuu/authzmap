package main

// EXPECT: two http.Server values: TLS one authed, plain one with the raw router -> flagged
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
	r.Put("/settings", h)
	secure := &http.Server{Addr: ":443", Handler: RequireAuth(r)}
	plain := http.Server{Addr: ":80"}
	plain.Handler = r
	go secure.ListenAndServeTLS("c.pem", "k.pem")
	plain.ListenAndServe()
}
