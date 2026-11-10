package main

// EXPECT: enforcing via panic, recovered to 401 -> no gaps
import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var errUnauthorized = errors.New("unauthorized")

func valid(t string) bool { return t == "s3cret" }
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e == errUnauthorized {
				http.Error(w, "unauthorized", 401)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			panic(errUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(Recover)
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
