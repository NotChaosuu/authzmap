package main

// EXPECT: routes built by reflection -> unseen/untraced -> INCOMPLETE (exit 3), not complete:true
import (
	"net/http"
	"reflect"
	"strings"
)

type API struct{}

func (API) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {}
func (API) AdminReset(w http.ResponseWriter, r *http.Request)      {}
func main() {
	mux := http.NewServeMux()
	v := reflect.ValueOf(API{})
	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := v.Method(i).Interface().(func(http.ResponseWriter, *http.Request))
		mux.HandleFunc("/"+strings.ToLower(t.Method(i).Name), m)
	}
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":8080", mux)
}
