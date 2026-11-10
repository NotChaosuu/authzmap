package main

// EXPECT: conditional-auth: inside the middleware, auth wraps next only when a config flag is set, else next runs bare -> conditional-auth MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

var AuthEnabled bool

func checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := next
		if AuthEnabled {
			h = checkToken(next)
		}
		h.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
