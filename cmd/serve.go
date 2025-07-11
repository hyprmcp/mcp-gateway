package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/cors"
	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/log"
	"github.com/jetski-sh/mcp-proxy/oauth"
	"github.com/jetski-sh/mcp-proxy/proxy"
	"github.com/spf13/cobra"
)

type ServeOptions struct {
	Config string
	Addr   string
}

func BindServeOptions(cmd *cobra.Command, opts *ServeOptions) {
	cmd.Flags().StringVarP(&opts.Config, "config", "c", "config.yaml", "Path to the configuration file")
	cmd.Flags().StringVarP(&opts.Addr, "addr", "a", ":9000", "Address to listen on")
}

func runServe(ctx context.Context, opts ServeOptions) error {
	config, err := config.ParseFile(opts.Config)
	if err != nil {
		return err
	}

	log.Get(ctx).Info("Loaded configuration", "config", config)

	mux := http.NewServeMux()

	for _, p := range config.Proxy {
		if p.Http != nil && p.Http.Url != nil {
			mux.Handle(p.Path, proxy.NewProxyHandler((*url.URL)(p.Http.Url)))
		}
	}

	var handler http.Handler = mux

	if len(config.AuthorizationServers) > 0 {
		handler = oauth.NewOAuthMiddleware(config.AuthorizationServers)(handler)
	}

	handler = cors.AllowAll().Handler(handler)

	log.Get(ctx).Info("Starting server", "addr", opts.Addr)

	if err := http.ListenAndServe(opts.Addr, handler); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}
