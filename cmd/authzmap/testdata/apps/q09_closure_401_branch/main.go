package main

// EXPECT: enforcing via a local deny closure (writes 401) followed by return -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deny := func(msg string) { http.Error(w, msg, http.StatusUnauthorized) }
		if !valid(r.Header.Get("Authorization")) {
			deny("bad token")
			return
		}
		next.ServeHTTP(w, r)
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
