package main

// EXPECT: func newServer(h) *http.Server called with RequireAuth(r) and with bare r; both served -> routes flagged (open listener)
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

func newServer(addr string, hd http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: hd}
}
func main() {
	r := chi.NewRouter()
	r.Delete("/admin/users/{id}", h)
	pub := newServer(":443", RequireAuth(r))
	internal := newServer(":8081", r)
	go internal.ListenAndServe()
	pub.ListenAndServe()
}
