package htmlresponse

import (
	_ "embed"
	"html/template"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/hyprmcp/mcp-gateway/config"
)

var (
	//go:embed template.html
	ts  string
	tpl *template.Template
)

const (
	ContentTypeTextHTML = "text/html"
)

func init() {
	if t, err := template.New("").Parse(ts); err != nil {
		panic(err)
	} else {
		tpl = t
	}
}

type handler struct {
	config         *config.Config
	alwaysCallNext bool
}

func NewHandler(config *config.Config, alwaysCallNext bool) *handler {
	return &handler{config: config, alwaysCallNext: alwaysCallNext}
}

func (h *handler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAcceptHTML := strings.Contains(r.Header.Get("Accept"), ContentTypeTextHTML)
		if h.alwaysCallNext || !isAcceptHTML {
			next.ServeHTTP(w, r)
		}

		if isAcceptHTML {
			h.handleHtml(w, r)
		}
	})
}

func (h *handler) handleHtml(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name string
		Url  string
	}

	u, _ := url.Parse(h.config.Host.String())
	u.Path = r.URL.Path
	u.RawQuery = r.URL.RawQuery
	data.Url = u.String()

	ps := strings.Split(r.URL.Path, "/")
	if nameIdx := slices.IndexFunc(ps, func(s string) bool { return s != "" && s != "mcp" }); nameIdx >= 0 {
		data.Name = ps[nameIdx]
	}

	w.Header().Set("Content-Type", ContentTypeTextHTML)

	if !h.alwaysCallNext {
		w.WriteHeader(http.StatusNotAcceptable)
	}

	_ = tpl.Execute(w, data)
}
