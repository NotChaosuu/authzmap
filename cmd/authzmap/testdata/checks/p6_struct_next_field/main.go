package main

// EXPECT: middleware as struct with a next field and ServeHTTP: authMW{next: mux} enforces -> quiet; RequireAuthMW{next} never rejects -> MEDIUM
import "net/http"

func valid(s string) bool                      { return s == "k" }
func h(w http.ResponseWriter, r *http.Request) {}

type authMW struct{ next http.Handler }

func (a authMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !valid(r.Header.Get("Authorization")) {
		w.WriteHeader(401)
		return
	}
	a.next.ServeHTTP(w, r)
}

type RequireAuthMW struct{ inner http.Handler }

func (m *RequireAuthMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println(r.URL.Path)
	m.inner.ServeHTTP(w, r)
}
func main() {
	a := http.NewServeMux()
	a.HandleFunc("DELETE /a/users/{id}", h)
	b := http.NewServeMux()
	b.HandleFunc("DELETE /b/users/{id}", h)
	go http.ListenAndServe(":8081", authMW{next: a})
	http.ListenAndServe(":8080", &RequireAuthMW{inner: b})
}
