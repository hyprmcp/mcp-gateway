package oauth

import "net/http"

type muxRegistration struct {
	pattern string
	handler http.Handler
}

func (reg *muxRegistration) Register(mux *http.ServeMux) {
	mux.Handle(reg.pattern, reg.handler)
}

type muxRegistrations []muxRegistration

func (reg *muxRegistrations) Add(pattern string, handler http.Handler) {
	*reg = append(*reg, muxRegistration{pattern: pattern, handler: handler})
}

func (reg muxRegistrations) Register(mux *http.ServeMux) {
	for _, r := range reg {
		r.Register(mux)
	}
}
