package main

// EXPECT: deny closure only logs -> never rejects -> flagged
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deny := func(reason string) { log.Printf("auth denied: %s (status %d)", reason, http.StatusUnauthorized) }
		if !valid(r.Header.Get("Authorization")) {
			deny("bad token")
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
