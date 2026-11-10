package main

// EXPECT: alice: chain with RequireAuth via Then -> quiet; alice chain with Logger only via ThenFunc on an admin route -> MEDIUM; Append adds auth -> quiet
import (
	"github.com/justinas/alice"
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

func main() {
	base := alice.New(Logger)
	secure := base.Append(RequireAuth)
	mux := http.NewServeMux()
	mux.Handle("DELETE /a/users/{id}", secure.ThenFunc(h))
	mux.Handle("DELETE /b/users/{id}", base.ThenFunc(h))
	mux.Handle("DELETE /c/users/{id}", alice.New(Logger, RequireAuth).Then(http.HandlerFunc(h)))
	http.ListenAndServe(":8080", mux)
}
