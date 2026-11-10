package main

// EXPECT: grpc-gateway passed straight as http.Server{Handler: gw} with no auth -> unanalysed-api MEDIUM (or INCOMPLETE), not "no routes"
import (
	"context"
	pb "example.com/q59/gen/adminpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func main() {
	gw := runtime.NewServeMux()
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	srv := &http.Server{Addr: ":8080", Handler: gw}
	srv.ListenAndServe()
}
