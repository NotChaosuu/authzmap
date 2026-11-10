package main

// EXPECT: routes appended in a helper (not a literal) -> either resolved or untraced/INCOMPLETE
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type route struct {
	method, path string
	h            http.HandlerFunc
}

func h(w http.ResponseWriter, r *http.Request) {}
func adminRoutes() []route {
	var rs []route
	rs = append(rs, route{"DELETE", "/admin/users/{id}", h})
	rs = append(rs, route{"POST", "/admin/reset", h})
	return rs
}
func main() {
	r := chi.NewRouter()
	for _, rt := range adminRoutes() {
		r.MethodFunc(rt.method, rt.path, rt.h)
	}
	http.ListenAndServe(":8080", r)
}
