package debugsrv

import (
	"net/http"

	"example.com/p5_init_listener_other_pkg/routes"
	"github.com/go-chi/chi/v5"
)

func start() {
	r := chi.NewRouter()
	routes.Register(r)
	go http.ListenAndServe(":6060", r)
}

func init() { start() }
