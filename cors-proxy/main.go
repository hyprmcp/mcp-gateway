package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/cors"
)

func main() {
	upstream, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	cors := cors.AllowAll()
	if err := http.ListenAndServe(":8081", cors.Handler(proxy)); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
