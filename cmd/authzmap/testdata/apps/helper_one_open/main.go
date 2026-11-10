package main

// EXPECT: the same helper on an authed router and an open one -> the open copy is a gap
import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func valid(t string) bool { return t == "s" }
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
func files(r chi.Router)                       { r.Delete("/files/{id}", h) }
func main() {
	api := chi.NewRouter()
	api.Use(TokenAuth)
	files(api)
	api.Get("/me", h)
	internal := chi.NewRouter()
	files(internal)
	go http.ListenAndServe(":9000", internal)
	http.ListenAndServe(":8080", api)
}
