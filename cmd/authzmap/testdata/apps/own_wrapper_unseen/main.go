package main

// EXPECT: an app's own router wrapper registers with a path it builds at run time (gitea's shape): its call sites stay unchecked -> INCOMPLETE, never a clean pass
import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	r      chi.Router
	prefix string
}

func (m *Router) Delete(p string, h http.HandlerFunc) { m.r.Delete(m.prefix+p, h) }
func h(w http.ResponseWriter, r *http.Request)        {}
func main() {
	m := &Router{r: chi.NewRouter(), prefix: os.Getenv("PREFIX")}
	m.Delete("/admin/users/{id}", h)
	m.Delete("/admin/repos/{id}", h)
	http.ListenAndServe(":8080", m.r)
}
