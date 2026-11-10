package main

// EXPECT: two unrelated `ok` variables (a context lookup, then BasicAuth) and checks after the auth block; auth only when credentials are configured -> conditional MEDIUM
import (
	"context"
	"net/http"
	"strings"
)

type key int

var uiPass string

func wrap(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(key(1)).(int64); !ok {
			r = r.WithContext(context.Background())
		}
		if uiPass != "" {
			_, pass, ok := r.BasicAuth()
			if !ok || pass != uiPass {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "identity")
		fn(w, r)
	}
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/messages/{id}", wrap(h))
	http.ListenAndServe(":8080", mux)
}
