package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bunrouter"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request)            {}
func bh(w http.ResponseWriter, r bunrouter.Request) error { return nil }

func main() {
	r := chi.NewRouter()
	r.Get("/healthz", h)
	br := bunrouter.New()
	br.DELETE("/admin/users/:id", bh)
	br.POST("/admin/reset", bh)
	go http.ListenAndServe(":9000", br)
	http.ListenAndServe(":8080", r)
}
