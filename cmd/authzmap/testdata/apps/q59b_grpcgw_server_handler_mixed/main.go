package main

// EXPECT: authed chi REST server plus an open grpc-gateway served by a second http.Server{Handler: gw} -> unanalysed-api MEDIUM or INCOMPLETE; never clean
import (
	"context"
	pb "example.com/q59b/gen/adminpb"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func valid(s string) bool { return s == "x" }
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
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	go (&http.Server{Addr: ":8081", Handler: gw}).ListenAndServe()
	r := chi.NewRouter()
	r.Use(RequireAuth)
	r.Delete("/api/users/{id}", h)
	http.ListenAndServe(":8080", r)
}
