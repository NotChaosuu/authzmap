package main

// EXPECT: net/http middleware named RequireAuth-ish (AuthMiddleware) whose only stop is a 429 rate limit -> never rejects on auth -> MEDIUM
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}

type limiter struct{}

func (limiter) Allow(k string) bool { return true }

var lim limiter

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lim.Allow(r.RemoteAddr) {
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	r.Use(AuthMiddleware)
	r.Delete("/admin/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
