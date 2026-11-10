package main

// EXPECT: servers in a slice literal, one authed one open, started in a loop -> flagged
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
	r.Delete("/admin/users/{id}", h)
	servers := []*http.Server{
		{Addr: ":443", Handler: RequireAuth(r)},
		{Addr: ":8080", Handler: r},
	}
	for _, s := range servers {
		go s.ListenAndServe()
	}
	select {}
}
