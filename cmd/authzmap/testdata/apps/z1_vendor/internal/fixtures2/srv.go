package fixtures2

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Fake() {
	r := chi.NewRouter()
	r.Delete("/admin/fake/{id}", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":4", r)
}
