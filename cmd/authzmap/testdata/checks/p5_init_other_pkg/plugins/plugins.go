package plugins

import (
	"example.com/p5_init_other_pkg/routes"
	"github.com/go-chi/chi/v5"
)

var Root = chi.NewRouter()

func setup() { routes.Register(Root) }

func init() { setup() }
