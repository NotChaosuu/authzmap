package main

// EXPECT: non-constant path -> untraced/partial noted; scan should not be complete with 0 routes or should flag
import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := chi.NewRouter()
	prefix := os.Getenv("ADMIN_PREFIX")
	r.Delete(prefix+"/users/{id}", h)
	r.Post(prefix+"/reset", h)
	http.ListenAndServe(":8080", r)
}
