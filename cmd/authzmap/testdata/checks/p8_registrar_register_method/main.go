package main

// EXPECT: AdGuardHome shape: NewDefaultRegistrar(mux, wrap) in another package, passed through a config struct, Register(method, path, h) per route -> routes found (partial) or INCOMPLETE; never complete with them missing
import (
	"example.com/p8_registrar_register_method/aghhttp"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request)              {}
func noAuth(method string, h http.HandlerFunc) http.Handler { return h }

type webConfig struct {
	mux     *http.ServeMux
	httpReg aghhttp.Registrar
}
type webAPI struct{ httpReg aghhttp.Registrar }

func (web *webAPI) registerControlHandlers() {
	web.httpReg.Register(http.MethodPost, "/control/update", h)
	web.httpReg.Register(http.MethodDelete, "/control/users", h)
}
func newWeb(conf *webConfig) *webAPI {
	w := &webAPI{httpReg: conf.httpReg}
	w.registerControlHandlers()
	return w
}
func main() {
	mux := http.NewServeMux()
	httpReg := aghhttp.NewDefaultRegistrar(mux, noAuth)
	newWeb(&webConfig{mux: mux, httpReg: httpReg})
	http.ListenAndServe(":8080", mux)
}
