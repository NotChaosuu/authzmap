package main

// EXPECT: oapi-codegen chi HandlerWithOptions with Middlewares: []MiddlewareFunc{RequireAuth} -> routes authed, no gaps. Second server without middlewares -> gaps
import (
	"example.com/c04_oapi_chi_mw/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type impl struct{}

func (impl) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	r := chi.NewRouter()
	api.HandlerWithOptions(impl{}, api.ChiServerOptions{BaseURL: "/v1", BaseRouter: r, Middlewares: []api.MiddlewareFunc{RequireAuth}})
	r2 := chi.NewRouter()
	api.HandlerWithOptions(impl{}, api.ChiServerOptions{BaseURL: "/open", BaseRouter: r2})
	r.Mount("/", r2)
	http.ListenAndServe(":8080", r)
}
