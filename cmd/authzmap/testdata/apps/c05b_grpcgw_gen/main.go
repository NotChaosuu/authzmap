package main

// EXPECT: grpc-gateway with generated .pb.gw.go in-tree; DELETE /admin/users/{id} unauthed -> found+reported or INCOMPLETE
import (
	"context"
	pb "example.com/c05b_grpcgw_gen/gen/adminpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type srv struct{}

func main() {
	ctx := context.Background()
	gw := runtime.NewServeMux()
	pb.RegisterAdminServiceHandlerServer(ctx, gw, srv{})
	mux := http.NewServeMux()
	mux.Handle("/", gw)
	http.ListenAndServe(":8080", mux)
}
