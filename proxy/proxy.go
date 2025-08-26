package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hyprmcp/mcp-gateway/config"
)

func NewProxyHandler(config *config.Proxy) http.Handler {
	url := (*url.URL)(config.Http.Url)

	return &httputil.ReverseProxy{
		Rewrite: RewriteFullFunc(url),
		Transport: &mcpAwareTransport{
			config: config,
		},
	}
}

func RewriteFullFunc(url *url.URL) func(r *httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		r.Out.URL.Scheme = url.Scheme
		r.Out.URL.Host = url.Host
		r.Out.URL.Path = url.Path
		r.Out.URL.RawPath = url.RawPath
		if r.Out.URL.RawQuery == "" || url.RawQuery == "" {
			r.Out.URL.RawQuery = r.Out.URL.RawQuery + url.RawQuery
		} else {
			r.Out.URL.RawQuery = url.RawQuery + "&" + r.Out.URL.RawQuery
		}
		r.Out.Host = ""
	}
}

func RewriteHostFunc(url *url.URL) func(r *httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		r.Out.URL.Scheme = url.Scheme
		r.Out.URL.Host = url.Host
		r.Out.Host = ""
	}
}
