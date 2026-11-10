package main

// EXPECT: grpc-gateway assigned to an http.Handler var then wrapped (h = RequireAuth(h)) before mux.Handle -> no gaps
import (
	"context"
	pb "example.com/q57/gen/adminpb"
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
func main() {
	gw := runtime.NewServeMux()
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	var h http.Handler = gw
	h = RequireAuth(h)
	mux := http.NewServeMux()
	mux.Handle("/v1/", h)
	http.ListenAndServe(":8080", mux)
}
