package main

// EXPECT: iris routes -> found or INCOMPLETE
import (
	"github.com/kataras/iris/v12"
)

func h(ctx iris.Context) {}
func main() {
	app := iris.New()
	app.Delete("/admin/users/{id}", h)
	app.Post("/admin/reset", h)
	app.Listen(":8080")
}
