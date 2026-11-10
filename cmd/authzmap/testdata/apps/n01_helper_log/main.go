package main

// EXPECT: RequireAuth never rejects (helper only logs) -> flagged
import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func checkToken(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" {
		log.Println("missing token")
	}
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkToken(w, r)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) { r.Use(RequireAuth); r.Delete("/api/users/{id}", h); r.Post("/admin/reset", h) })
	http.ListenAndServe(":8080", r)
}
