package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"testing"
)

func TestX(t *testing.T) {
	r := chi.NewRouter()
	r.Delete("/admin/test/{id}", h)
	http.ListenAndServe(":3", r)
}
