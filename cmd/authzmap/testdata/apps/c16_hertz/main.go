package main

// EXPECT: hertz (cloudwego) routes unmodelled -> INCOMPLETE
import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func h(ctx context.Context, c *app.RequestContext) {}
func main() {
	s := server.Default()
	s.DELETE("/admin/users/:id", h)
	s.POST("/admin/reset", h)
	s.Spin()
}
