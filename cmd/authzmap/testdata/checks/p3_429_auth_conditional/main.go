package main

// EXPECT: 429-on-bad-credential gate with an Open receiver flag that skips it -> conditional MEDIUM
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

type gate struct{ Open bool }

func (g *gate) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if g.Open {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	g := &gate{}
	r := chi.NewRouter()
	r.Use(g.Check)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
