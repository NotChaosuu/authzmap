package main

// EXPECT: 401 written without return behind 12 harmless switches -> never rejects MEDIUM
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
		switch r.Header.Get("X-Mode0") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode1") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode2") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode3") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode4") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode5") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode6") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode7") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode8") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode9") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode10") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		switch r.Header.Get("X-Mode11") {
		case "a":
			println(1)
		case "b":
			println(2)
		}
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
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
