package main

import "net/http"

type server struct{}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/healthz" {
		w.Write([]byte("ok"))
	} else if r.Method == http.MethodDelete && r.URL.Path == "/admin/users" {
		s.deleteAllUsers(w, r)
	} else if r.Method == http.MethodGet && r.URL.Path == "/v1/me" {
		s.withAuth(s.me)(w, r)
	}
}
func (s *server) deleteAllUsers(w http.ResponseWriter, r *http.Request) {}
func (s *server) me(w http.ResponseWriter, r *http.Request)             {}
func (s *server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "", 401)
			return
		}
		next(w, r)
	}
}

func main() {
	s := &server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	http.ListenAndServe(":8080", mux)
}
