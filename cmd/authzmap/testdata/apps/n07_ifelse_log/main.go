package main

// EXPECT: enforcing (handler not reached on failure though no 401 written) -> must NOT be never-rejects
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if valid(r.Header.Get("Authorization")) {
			next.ServeHTTP(w, r)
		} else {
			log.Println("rejected")
		}
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
