package main

// EXPECT: registration through a method value / func var -> found or untraced
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	del := r.Delete
	del("/admin/users/{id}", h)
	reg := func(p string) { r.Post(p, h) }
	reg("/admin/reset")
	http.ListenAndServe(":8080", r)
}
