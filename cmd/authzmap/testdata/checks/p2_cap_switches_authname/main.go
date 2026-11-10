package main

// EXPECT: auth-named middleware, 7 switches x 4 cases on credential headers (16384 paths); every switch has a pass-through case so a request with no headers passes -> MEDIUM (never-rejects or unverified), not trusted
import (
	"github.com/go-chi/chi/v5"
	"net/http"
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
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("X-S0") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S1") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S2") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S3") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S4") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S5") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		switch r.Header.Get("X-S6") {
		case "a":
			println("a")
		case "b":
			w.WriteHeader(401)
			return
		case "c":
			println("c")
		default:
			println("d")
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(AuthRequired)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
