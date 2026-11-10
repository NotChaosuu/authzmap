package main

// EXPECT: an exported method called only via reflect.MethodByName builds an open router -> exported so never dropped -> reported (or INCOMPLETE)
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "k" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func registerAdmin(r chi.Router) { r.Delete("/admin/users/{id}", h) }
func authedRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(RequireAuth)
	registerAdmin(r)
	return r
}

type Srv struct{}

func (Srv) OpenRouter() chi.Router {
	r := chi.NewRouter()
	registerAdmin(r)
	return r
}
func main() {
	go http.ListenAndServe(":8443", authedRouter())
	m := reflect.ValueOf(Srv{}).MethodByName("OpenRouter")
	out := m.Call(nil)
	http.ListenAndServe(":8080", out[0].Interface().(http.Handler))
}
