package main

// EXPECT: a method on an app struct returns *http.Server with the router unwrapped; a second one wraps -> flagged
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

type App struct{ router chi.Router }

func (a *App) public() *http.Server { return &http.Server{Addr: ":80", Handler: a.router} }
func (a *App) private() *http.Server {
	return &http.Server{Addr: ":443", Handler: RequireAuth(a.router)}
}
func main() {
	r := chi.NewRouter()
	r.Post("/admin/reset", h)
	a := &App{router: r}
	go a.public().ListenAndServe()
	a.private().ListenAndServeTLS("c", "k")
}
