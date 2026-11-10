package main

// EXPECT: connect handler with an auth interceptor (connect.WithInterceptors(NewAuthInterceptor())) -> at most LOW/advisory; not a confident MEDIUM gap
import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"example.com/q51/gen/adminv1/adminv1connect"
	"net/http"
)

type srv struct{}

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Header().Get("Authorization") == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no token"))
			}
			return next(ctx, req)
		}
	}
}
func main() {
	mux := http.NewServeMux()
	mux.Handle(adminv1connect.NewAdminServiceHandler(srv{}, connect.WithInterceptors(NewAuthInterceptor())))
	http.ListenAndServe(":8080", mux)
}
