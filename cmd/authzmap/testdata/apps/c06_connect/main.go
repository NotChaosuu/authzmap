package main

// EXPECT: connect-rpc handler with dynamic path and no auth interceptor -> reported, or partial, not silent
import (
	"example.com/c06_connect/gen/adminv1/adminv1connect"
	"net/http"
)

type srv struct{}

func main() {
	mux := http.NewServeMux()
	path, handler := adminv1connect.NewAdminServiceHandler(srv{})
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", mux)
}
