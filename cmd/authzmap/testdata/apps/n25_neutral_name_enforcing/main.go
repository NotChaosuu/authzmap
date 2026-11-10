package main

// EXPECT: neutral-named middleware that DOES enforce 401 -> ideally treated as auth (no gap); at least not HIGH
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(t string) bool { return t == "s3cret" }
func gatekeeper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(gatekeeper)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
