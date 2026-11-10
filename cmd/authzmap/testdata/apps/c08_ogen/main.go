package main

// EXPECT: ogen server is the whole handler; routes live in generated code outside tree -> INCOMPLETE or partial /* route
import (
	"example.com/c08_ogen/oas"
	"net/http"
)

type h struct{}

func main() {
	srv, err := oas.NewServer(h{})
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", srv)
}
