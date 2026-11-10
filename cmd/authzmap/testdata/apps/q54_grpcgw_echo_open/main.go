package main

// EXPECT: grpc-gateway mounted in echo via e.Any + echo.WrapHandler with no auth -> unanalysed-api MEDIUM
import (
	"context"
	pb "example.com/q54/gen/adminpb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
)

func main() {
	gw := runtime.NewServeMux()
	_ = pb.RegisterAdminServiceHandlerFromEndpoint(context.Background(), gw, "localhost:9090", nil)
	e := echo.New()
	e.Any("/v1/*", echo.WrapHandler(gw))
	e.Start(":8080")
}
