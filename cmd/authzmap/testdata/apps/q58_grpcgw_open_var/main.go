package main

// EXPECT: control for q57 without the wrap -> unanalysed-api MEDIUM
import (
	"context"
	pb "example.com/q58/gen/adminpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

func main() {
	gw := runtime.NewServeMux()
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	var h http.Handler = gw
	mux := http.NewServeMux()
	mux.Handle("/v1/", h)
	http.ListenAndServe(":8080", mux)
}
