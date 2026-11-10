package main

// EXPECT: connect path, handler pair registered with the handler wrapped in RequireAuth -> quiet
import (
	"example.com/connect_wrapped_pair/gen/adminv1connect"
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
func h(w http.ResponseWriter, r *http.Request) {}

func main() {
	mux := http.NewServeMux()
	path, svc := adminv1connect.NewAdminServiceHandler(nil)
	mux.Handle(path, RequireAuth(svc))
	http.ListenAndServe(":8080", mux)
}
