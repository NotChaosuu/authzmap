package main

import "net/http/httptest"

func newTestServer() *httptest.Server { return httptest.NewServer(routes()) }
