package main

// EXPECT: sub mounted inside an authed group first, then open at /legacy -> /legacy/* flagged, /secure/* quiet
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
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
	sub := chi.NewRouter()
	sub.Delete("/users/{id}", h)
	sub.Post("/purge", h)
	r := chi.NewRouter()
	r.Route("/secure", func(s chi.Router) { s.Use(RequireAuth); s.Mount("/", sub) })
	r.Mount("/legacy", sub)
	http.ListenAndServe(":8080", r)
}
