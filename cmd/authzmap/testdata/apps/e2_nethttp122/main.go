package main

import "net/http"

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

func init() {
	http.HandleFunc("DELETE /api/users/{id}", h)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /api/me", requireAuth(http.HandlerFunc(h)))
	mux.Handle("PUT /api/me", requireAuth(http.HandlerFunc(h)))
	mux.Handle("GET /api/users/{id}", requireAuth(http.HandlerFunc(h)))
	mux.HandleFunc("DELETE /api/me", h)
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("/"))))
	mux.Handle("/admin/", http.StripPrefix("/admin", adminMux()))
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}

func adminMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("POST /users/{id}/delete", h)
	return m
}
