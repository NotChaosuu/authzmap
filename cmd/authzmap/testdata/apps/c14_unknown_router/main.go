package main

// EXPECT: routes on an unmodelled router (go-pkgz/routegroup, julienschmidt.httprouter-like custom) -> INCOMPLETE
import (
	"github.com/go-pkgz/routegroup"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	mux := http.NewServeMux()
	g := routegroup.New(mux)
	g.HandleFunc("DELETE /admin/users/{id}", h)
	api := g.Mount("/api")
	api.HandleFunc("POST /admin/reset", h)
	http.ListenAndServe(":8080", mux)
}
