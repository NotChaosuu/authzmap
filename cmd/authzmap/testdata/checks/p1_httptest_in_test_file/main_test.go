package main

import (
	"net/http/httptest"
	"testing"
)

func TestX(t *testing.T) {
	srv := httptest.NewServer(routes())
	defer srv.Close()
}
