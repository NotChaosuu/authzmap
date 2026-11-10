package main

// EXPECT: connect idiom mux.Handle(adminv1connect.NewAdminServiceHandler(srv{})) with no auth -> unanalysed-api MEDIUM
import (
	"example.com/q52/gen/adminv1/adminv1connect"
	"net/http"
)

type srv struct{}

func main() {
	mux := http.NewServeMux()
	mux.Handle(adminv1connect.NewAdminServiceHandler(srv{}))
	http.ListenAndServe(":8080", mux)
}
