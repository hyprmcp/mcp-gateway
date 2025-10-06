package proxyutil

import (
	"net/http"
	"net/http/httputil"
)

func ModifyResponseChain(fns ...func(*http.Response) error) func(*http.Response) error {
	return func(resp *http.Response) error {
		for _, fn := range fns {
			if fn != nil {
				if err := fn(resp); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func RewriteChain(fns ...func(*httputil.ProxyRequest)) func(*httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		for _, fn := range fns {
			if fn != nil {
				fn(r)
			}
		}
	}
}
