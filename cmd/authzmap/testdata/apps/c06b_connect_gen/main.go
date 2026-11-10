package main

// EXPECT: connect-rpc gen in tree, no auth interceptor -> reported or INCOMPLETE (dispatch)
import (
	"example.com/c06b_connect_gen/gen/adminv1/adminv1connect"
	"net/http"
)

type srv struct{}

func main() {
	mux := http.NewServeMux()
	path, handler := adminv1connect.NewAdminServiceHandler(srv{})
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", mux)
}
