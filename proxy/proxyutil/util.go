package proxyutil

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

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

// RemoveCORSHeaders removes all CORS related headers from the http.Response.
//
// This is necessary because the ReverseProxy appends these headers to the ones
// already existing on the downstream response, however the downstream response
// already contains CORS headers which get added by the CORS middleware.
func RemoveCORSHeaders(resp *http.Response) error {
	if resp == nil {
		return nil
	}
	resp.Header.Del("Access-Control-Allow-Origin")
	resp.Header.Del("Access-Control-Allow-Methods")
	resp.Header.Del("Access-Control-Allow-Headers")
	resp.Header.Del("Access-Control-Allow-Credentials")
	resp.Header.Del("Access-Control-Expose-Headers")
	resp.Header.Del("Access-Control-Max-Age")
	return nil
}
