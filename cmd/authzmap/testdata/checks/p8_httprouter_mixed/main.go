package main

// EXPECT: httprouter is documented as unmodelled: a mixed app must be INCOMPLETE (exit 3), never clean
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func valid(s string) bool                                            { return s == "k" }
func hh(w http.ResponseWriter, r *http.Request, p httprouter.Params) {}
func Protect(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if !valid(r.Header.Get("Authorization")) {
			w.WriteHeader(401)
			return
		}
		next(w, r, p)
	}
}
func main() {
	router := httprouter.New()
	router.DELETE("/admin/users/:id", Protect(hh))
	router.DELETE("/admin/keys/:id", hh)
	http.ListenAndServe(":8080", router)
}
