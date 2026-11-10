package main

// EXPECT: 16 credential-header branches that only log: inert, so they merge -> never rejects, MEDIUM; never trusted by name
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

// RequireLogin is named like auth to make trust-by-name the risk
var _ = 0

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token-0") == "" {
			println(0)
		}
		if r.Header.Get("X-Token-1") == "" {
			println(1)
		}
		if r.Header.Get("X-Token-2") == "" {
			println(2)
		}
		if r.Header.Get("X-Token-3") == "" {
			println(3)
		}
		if r.Header.Get("X-Token-4") == "" {
			println(4)
		}
		if r.Header.Get("X-Token-5") == "" {
			println(5)
		}
		if r.Header.Get("X-Token-6") == "" {
			println(6)
		}
		if r.Header.Get("X-Token-7") == "" {
			println(7)
		}
		if r.Header.Get("X-Token-8") == "" {
			println(8)
		}
		if r.Header.Get("X-Token-9") == "" {
			println(9)
		}
		if r.Header.Get("X-Token-10") == "" {
			println(10)
		}
		if r.Header.Get("X-Token-11") == "" {
			println(11)
		}
		if r.Header.Get("X-Token-12") == "" {
			println(12)
		}
		if r.Header.Get("X-Token-13") == "" {
			println(13)
		}
		if r.Header.Get("X-Token-14") == "" {
			println(14)
		}
		if r.Header.Get("X-Token-15") == "" {
			println(15)
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(RequireLogin)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
