package y

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Example() {
	r := chi.NewRouter()
	r.Delete("/admin/vendored/{id}", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":2", r)
}
