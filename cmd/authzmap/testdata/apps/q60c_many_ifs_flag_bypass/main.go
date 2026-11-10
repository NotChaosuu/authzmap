package main

// EXPECT: conditional-auth (package flag bypass) behind 10 harmless ifs -> conditional-auth MEDIUM; not silently trusted
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }

var SkipAuth bool

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-F0") != "" {
			w.Header().Set("X-Seen0", "1")
		}
		if r.Header.Get("X-F1") != "" {
			w.Header().Set("X-Seen1", "1")
		}
		if r.Header.Get("X-F2") != "" {
			w.Header().Set("X-Seen2", "1")
		}
		if r.Header.Get("X-F3") != "" {
			w.Header().Set("X-Seen3", "1")
		}
		if r.Header.Get("X-F4") != "" {
			w.Header().Set("X-Seen4", "1")
		}
		if r.Header.Get("X-F5") != "" {
			w.Header().Set("X-Seen5", "1")
		}
		if r.Header.Get("X-F6") != "" {
			w.Header().Set("X-Seen6", "1")
		}
		if r.Header.Get("X-F7") != "" {
			w.Header().Set("X-Seen7", "1")
		}
		if r.Header.Get("X-F8") != "" {
			w.Header().Set("X-Seen8", "1")
		}
		if r.Header.Get("X-F9") != "" {
			w.Header().Set("X-Seen9", "1")
		}
		if SkipAuth {
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
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	r.Post("/admin/reset", h)
	http.ListenAndServe(":8080", r)
}
