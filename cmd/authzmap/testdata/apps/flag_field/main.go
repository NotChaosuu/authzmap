package main

// EXPECT: a receiver field switch turns the 401 off -> conditional
import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

type authMW struct{ RequireLogin bool }

func valid(t string) bool { return t == "s" }
func (m authMW) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.RequireLogin {
			next.ServeHTTP(w, r)
			return
		}
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

var RequireAuth = authMW{RequireLogin: os.Getenv("OPEN") == ""}.Handler

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
