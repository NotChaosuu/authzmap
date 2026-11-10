package main

// EXPECT: deny closure writes 401 but caller does not return -> next runs -> flagged
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deny := func() { http.Error(w, "unauthorized", http.StatusUnauthorized) }
		if !valid(r.Header.Get("Authorization")) {
			deny()
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
