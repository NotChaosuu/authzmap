package main

// EXPECT: never rejects: the bad case writes 401 then falls through into the case that calls next -> auth-never-rejects MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func classify(s string) string {
	if valid(s) {
		return "good"
	}
	return "bad"
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch classify(r.Header.Get("Authorization")) {
		case "bad":
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			fallthrough
		case "good":
			next.ServeHTTP(w, r)
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
