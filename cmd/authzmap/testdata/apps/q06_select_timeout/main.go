package main

// EXPECT: enforcing: select on verify result vs timeout, both failure branches 401+return -> no gaps
import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func verify(s string) chan bool { c := make(chan bool, 1); c <- s == "x"; return c }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case ok := <-verify(r.Header.Get("Authorization")):
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		case <-time.After(time.Second):
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
