package main

// EXPECT: enforcing: 401 after range loop over keys, next only on match -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

var keys = []string{"a", "b"}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		for _, k := range keys {
			if k == tok {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
