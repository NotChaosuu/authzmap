package main

// EXPECT: a route table filled by append in a helper with a var of the element type -> routes resolved
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type route struct {
	method, path string
	fn           http.HandlerFunc
}

func h(w http.ResponseWriter, r *http.Request) {}
func table() []route {
	var out []route
	out = append(out, route{method: "DELETE", path: "/admin/users/{id}", fn: h})
	out = append(out, route{method: "POST", path: "/admin/reset", fn: h})
	return out
}
func main() {
	r := chi.NewRouter()
	for _, rt := range table() {
		r.MethodFunc(rt.method, rt.path, rt.fn)
	}
	http.ListenAndServe(":8080", r)
}
