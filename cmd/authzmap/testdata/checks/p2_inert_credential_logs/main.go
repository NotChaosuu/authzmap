package main

// EXPECT: 20 inert logging ifs that read the credential (`tok == "" || X-Ln`) then a real 401 check -> merged as inert, enforcing, quiet
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

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		if tok == "" || r.Header.Get("X-L0") != "" {
			println("missing", 0)
		}
		if tok == "" || r.Header.Get("X-L1") != "" {
			println("missing", 1)
		}
		if tok == "" || r.Header.Get("X-L2") != "" {
			println("missing", 2)
		}
		if tok == "" || r.Header.Get("X-L3") != "" {
			println("missing", 3)
		}
		if tok == "" || r.Header.Get("X-L4") != "" {
			println("missing", 4)
		}
		if tok == "" || r.Header.Get("X-L5") != "" {
			println("missing", 5)
		}
		if tok == "" || r.Header.Get("X-L6") != "" {
			println("missing", 6)
		}
		if tok == "" || r.Header.Get("X-L7") != "" {
			println("missing", 7)
		}
		if tok == "" || r.Header.Get("X-L8") != "" {
			println("missing", 8)
		}
		if tok == "" || r.Header.Get("X-L9") != "" {
			println("missing", 9)
		}
		if tok == "" || r.Header.Get("X-L10") != "" {
			println("missing", 10)
		}
		if tok == "" || r.Header.Get("X-L11") != "" {
			println("missing", 11)
		}
		if tok == "" || r.Header.Get("X-L12") != "" {
			println("missing", 12)
		}
		if tok == "" || r.Header.Get("X-L13") != "" {
			println("missing", 13)
		}
		if tok == "" || r.Header.Get("X-L14") != "" {
			println("missing", 14)
		}
		if tok == "" || r.Header.Get("X-L15") != "" {
			println("missing", 15)
		}
		if tok == "" || r.Header.Get("X-L16") != "" {
			println("missing", 16)
		}
		if tok == "" || r.Header.Get("X-L17") != "" {
			println("missing", 17)
		}
		if tok == "" || r.Header.Get("X-L18") != "" {
			println("missing", 18)
		}
		if tok == "" || r.Header.Get("X-L19") != "" {
			println("missing", 19)
		}
		if !valid(tok) {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(RequireLogin)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
