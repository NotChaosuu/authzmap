package main

// EXPECT: hertz middleware writes 401 and returns without Abort -> chain continues -> auth-never-rejects MEDIUM
import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func valid(s string) bool { return s == "x" }
func RequireAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !valid(string(c.GetHeader("Authorization"))) {
			c.JSON(401, "unauthorized")
			return
		}
		c.Next(ctx)
	}
}
func h(ctx context.Context, c *app.RequestContext) {}
func main() {
	s := server.Default()
	g := s.Group("/api", RequireAuth())
	g.DELETE("/users/:id", h)
	g.POST("/admin/reset", h)
	s.Spin()
}
