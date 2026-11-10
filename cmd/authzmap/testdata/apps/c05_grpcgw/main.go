package main

// EXPECT: grpc-gateway mux with generated handlers outside the tree -> partial/unseen warning, scan not silently complete
import (
	"context"
	pb "example.com/c05_grpcgw/gen/adminpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type srv struct{}

func main() {
	ctx := context.Background()
	gw := runtime.NewServeMux()
	pb.RegisterAdminServiceHandlerServer(ctx, gw, srv{})
	gw.HandlePath("DELETE", "/admin/users/{id}", func(w http.ResponseWriter, r *http.Request, p map[string]string) {})
	mux := http.NewServeMux()
	mux.Handle("/", gw)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":8080", mux)
}
