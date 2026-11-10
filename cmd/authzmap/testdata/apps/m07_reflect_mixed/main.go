package main

// EXPECT: reflect-registered admin routes beside authed routes -> visible at default or INCOMPLETE
import (
	"net/http"
	"reflect"
	"strings"
)

type API struct{}

func (API) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {}
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /app/me", RequireAuth(http.HandlerFunc(h)))
	mux.Handle("POST /app/notes", RequireAuth(http.HandlerFunc(h)))
	v := reflect.ValueOf(API{})
	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := v.Method(i).Interface().(func(http.ResponseWriter, *http.Request))
		mux.HandleFunc("/"+strings.ToLower(t.Method(i).Name), m)
	}
	http.ListenAndServe(":8080", mux)
}
