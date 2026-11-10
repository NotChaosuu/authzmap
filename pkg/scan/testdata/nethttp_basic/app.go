package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", login)
	mux.Handle("GET /me", RequireAuth(http.HandlerFunc(me)))

	// /admin/stats is passed bare; /admin/users is wrapped.
	mux.HandleFunc("GET /admin/stats", stats)
	mux.Handle("POST /admin/users", RequireAuth(http.HandlerFunc(createUser)))

	http.ListenAndServe(":8080", mux)
}
