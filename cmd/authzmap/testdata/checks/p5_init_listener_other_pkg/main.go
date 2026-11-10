package main

// EXPECT: another package starts its own open listener from init() via an unexported func; blank-imported -> reported
import (
	_ "example.com/p5_init_listener_other_pkg/debugsrv"
	"example.com/p5_init_listener_other_pkg/routes"
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
	r.Use(RequireAuth)
	routes.Register(r)
	http.ListenAndServe(":8080", r)
}
