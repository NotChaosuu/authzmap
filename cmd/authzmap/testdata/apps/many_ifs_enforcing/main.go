package main

// EXPECT: real 401 check behind 12 harmless ifs -> quiet (enforcing)
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

func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("p0") == "x" {
			println(0)
		}
		if r.URL.Query().Get("p1") == "x" {
			println(1)
		}
		if r.URL.Query().Get("p2") == "x" {
			println(2)
		}
		if r.URL.Query().Get("p3") == "x" {
			println(3)
		}
		if r.URL.Query().Get("p4") == "x" {
			println(4)
		}
		if r.URL.Query().Get("p5") == "x" {
			println(5)
		}
		if r.URL.Query().Get("p6") == "x" {
			println(6)
		}
		if r.URL.Query().Get("p7") == "x" {
			println(7)
		}
		if r.URL.Query().Get("p8") == "x" {
			println(8)
		}
		if r.URL.Query().Get("p9") == "x" {
			println(9)
		}
		if r.URL.Query().Get("p10") == "x" {
			println(10)
		}
		if r.URL.Query().Get("p11") == "x" {
			println(11)
		}
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "no", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(Gate)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
