package main

// EXPECT: check runs in a goroutine that can't stop the request -> never rejects
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(t string) bool { return t == "s" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		go func() {
			if !valid(tok) {
				return
			}
		}()
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
