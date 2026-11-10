package main

// EXPECT: another package registers an open router from an unexported func inside init(); main serves the authed one -> the init one is live (init runs on import) -> reported
import (
	"example.com/p5_init_other_pkg/plugins"
	"example.com/p5_init_other_pkg/routes"
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
	r.Group(func(g chi.Router) { g.Use(RequireAuth); routes.Register(g) })
	r.Mount("/", plugins.Root)
	http.ListenAndServe(":8080", r)
}
