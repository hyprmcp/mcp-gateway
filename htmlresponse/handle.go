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

func init() {
	if t, err := template.New("").Parse(ts); err != nil {
		panic(err)
	} else {
		tpl = t
	}
}

type handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *handler {
	return &handler{config: config}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) error {
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

	return tpl.Execute(w, data)
}
