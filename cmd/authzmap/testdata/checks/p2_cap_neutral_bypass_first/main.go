package main

// EXPECT: neutral-named Gate: bypass first + 13 non-mergeable credential branches -> MEDIUM (not auth / unverified / conditional), never quiet
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

type C struct{ Open bool }

var Cfg = &C{}

func Gate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("debug") == "1" {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-Token-0") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D0") != "" {
			w.Header().Set("X-Seen", "0")
		} else {
			println(0)
		}
		if r.Header.Get("X-E0") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-1") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D1") != "" {
			w.Header().Set("X-Seen", "1")
		} else {
			println(1)
		}
		if r.Header.Get("X-E1") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-2") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D2") != "" {
			w.Header().Set("X-Seen", "2")
		} else {
			println(2)
		}
		if r.Header.Get("X-E2") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-3") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D3") != "" {
			w.Header().Set("X-Seen", "3")
		} else {
			println(3)
		}
		if r.Header.Get("X-E3") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-4") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D4") != "" {
			w.Header().Set("X-Seen", "4")
		} else {
			println(4)
		}
		if r.Header.Get("X-E4") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-5") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D5") != "" {
			w.Header().Set("X-Seen", "5")
		} else {
			println(5)
		}
		if r.Header.Get("X-E5") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-6") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D6") != "" {
			w.Header().Set("X-Seen", "6")
		} else {
			println(6)
		}
		if r.Header.Get("X-E6") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-7") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D7") != "" {
			w.Header().Set("X-Seen", "7")
		} else {
			println(7)
		}
		if r.Header.Get("X-E7") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-8") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D8") != "" {
			w.Header().Set("X-Seen", "8")
		} else {
			println(8)
		}
		if r.Header.Get("X-E8") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-9") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D9") != "" {
			w.Header().Set("X-Seen", "9")
		} else {
			println(9)
		}
		if r.Header.Get("X-E9") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-10") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D10") != "" {
			w.Header().Set("X-Seen", "10")
		} else {
			println(10)
		}
		if r.Header.Get("X-E10") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-11") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D11") != "" {
			w.Header().Set("X-Seen", "11")
		} else {
			println(11)
		}
		if r.Header.Get("X-E11") == "1" {
			println("e")
		}
		if r.Header.Get("X-Token-12") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if r.Header.Get("X-D12") != "" {
			w.Header().Set("X-Seen", "12")
		} else {
			println(12)
		}
		if r.Header.Get("X-E12") == "1" {
			println("e")
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(Gate)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
