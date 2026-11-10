package main

// EXPECT: kubernetes-dashboard shape: go-restful (unmodelled) ws.Route(ws.DELETE(path).To(h)) mounted with http.Handle("/", container) next to one modelled net/http route -> must be INCOMPLETE (exit 3) or find the go-restful routes; never complete with them missing
import (
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func del(req *restful.Request, resp *restful.Response) {}
func main() {
	ws := new(restful.WebService)
	ws.Path("/api/v1")
	ws.Route(ws.DELETE("/users/{id}").To(del))
	ws.Route(
		ws.POST("/admin/reset").
			To(del))
	c := restful.NewContainer()
	c.Add(ws)
	http.Handle("/", c)
	http.Handle("/settings", RequireAuth(http.HandlerFunc(h)))
	http.ListenAndServe(":8080", nil)
}
