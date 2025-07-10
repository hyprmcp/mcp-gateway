package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxyHandler(url *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
		},
	}
}
