package main

// EXPECT: grpc-gateway served by http.Server{Handler: RequireAuth(gw)} -> quiet
import (
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	gw := runtime.NewServeMux()
	srv := &http.Server{Addr: ":8080", Handler: RequireAuth(gw)}
	srv.ListenAndServe()
}
