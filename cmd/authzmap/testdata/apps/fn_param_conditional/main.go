package main

// EXPECT: middleware hands on through its own parameter fn (mailpit's shape); basic auth only when credentials are configured -> conditional MEDIUM
import "net/http"

type creds struct{ user, pass string }

var uiCredentials *creds

func middleWareFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "" {
			http.Error(w, "bad host", http.StatusForbidden)
			return
		}
		if uiCredentials != nil {
			user, pass, ok := r.BasicAuth()
			if !ok || user != uiCredentials.user || pass != uiCredentials.pass {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		fn(w, r)
	}
}
func h(w http.ResponseWriter, r *http.Request) {}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/messages/{id}", middleWareFunc(h))
	mux.HandleFunc("POST /api/admin/reset", middleWareFunc(h))
	http.ListenAndServe(":8080", mux)
}
