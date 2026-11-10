package main

// EXPECT: gorilla route-builder chains found; /admin unauthed reported
import (
	"github.com/gorilla/mux"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	r := mux.NewRouter()
	r.NewRoute().Path("/admin/users/{id}").Methods("DELETE").HandlerFunc(h)
	r.Path("/admin/reset").Methods(http.MethodPost).Handler(http.HandlerFunc(h))
	r.Methods("PUT").Path("/admin/settings").HandlerFunc(h)
	http.ListenAndServe(":8080", r)
}
