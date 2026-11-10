package main

// EXPECT: one helper registered on two routers guarded by different auth -> no gaps
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(t string) bool { return t == "s" }
func SessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Cookie")) {
			http.Error(w, "x", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "x", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func files(r chi.Router)                       { r.Delete("/files/{id}", h); r.Post("/files", h) }
func main() {
	web := chi.NewRouter()
	web.Use(SessionAuth)
	files(web)
	api := chi.NewRouter()
	api.Use(TokenAuth)
	files(api)
	go http.ListenAndServe(":8081", api)
	http.ListenAndServe(":8080", web)
}
