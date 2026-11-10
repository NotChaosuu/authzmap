package main

// EXPECT: enforcing: switch on a config mode string, every mode rejects -> no gaps (NOT conditional-auth)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

var AuthMode = "token"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := false
		if AuthMode == "token" {
			ok = valid(r.Header.Get("Authorization"))
		} else {
			c, err := r.Cookie("session")
			ok = err == nil && valid(c.Value)
		}
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
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
