package main

// EXPECT: beego routes unmodelled -> INCOMPLETE
import (
	"github.com/beego/beego/v2/server/web"
)

type AdminController struct{ web.Controller }

func (c *AdminController) Delete() {}
func main() {
	web.Router("/admin/users/:id", &AdminController{}, "delete:Delete")
	web.Router("/admin/reset", &AdminController{}, "post:Delete")
	web.Run()
}
