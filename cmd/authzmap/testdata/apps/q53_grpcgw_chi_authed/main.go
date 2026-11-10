package main

// EXPECT: grpc-gateway mounted under chi with RequireAuth -> no gaps
import (
	"context"
	pb "example.com/q53/gen/adminpb"
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
func main() {
	gw := runtime.NewServeMux()
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	r := chi.NewRouter()
	r.With(RequireAuth).Mount("/v1", gw)
	http.ListenAndServe(":8080", r)
}
