package main

// EXPECT: grpc-gateway behind StripPrefix on a mux, no auth -> unanalysed-api MEDIUM
import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	gw := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/rpc/", http.StripPrefix("/rpc", gw))
	http.ListenAndServe(":8080", mux)
}
