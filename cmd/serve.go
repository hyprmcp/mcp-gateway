package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
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
	cfg, err := config.ParseFile(opts.Config)
	if err != nil {
		return err
	}

	log.Get(ctx).Info("Loaded configuration", "config", cfg)

	handler := &delegateHandler{}

	if h, err := newHandler(ctx, cfg); err != nil {
		return err
	} else {
		handler.delegate = h
	}

	go func() {
		err := WatchConfigChanges(
			opts.Config,
			func(c *config.Config) {
				log.Get(ctx).Info("Reconfiguring server after config change...")
				if h, err := newHandler(ctx, c); err != nil {
					log.Get(ctx).Error(err, "failed to reload server")
				} else {
					handler.delegate = h
				}
			},
		)
		if err != nil {
			log.Get(ctx).Error(err, "config watch failed")
		}
	}()

	log.Get(ctx).Info("Starting server", "addr", opts.Addr)

	if err := http.ListenAndServe(opts.Addr, cors.AllowAll().Handler(handler)); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func newHandler(ctx context.Context, config *config.Config) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, p := range config.Proxy {
		if p.Http != nil && p.Http.Url != nil {
			mux.Handle(p.Path, proxy.NewProxyHandler(&p))
		}
	}

	var handler http.Handler = mux

	if oauthHandler, err := oauth.NewOAuthMiddleware(ctx, config); err != nil {
		return nil, err
	} else {
		handler = oauthHandler(handler)
	}

	return handler, nil
}

func WatchConfigChanges(path string, callback func(*config.Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func() { _ = watcher.Close() }()

	if err := watcher.Add(path); err != nil {
		return err
	}

	for {
		select {
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return err
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op != fsnotify.Rename && event.Op != fsnotify.Write {
				continue
			}

			log.Root().Info("starting config reload", "op", event.Op, "path", path)

			if cfg, err := config.ParseFile(path); err != nil {
				log.Root().Error(err, "config reload error", "op", event.Op)
			} else {
				callback(cfg)
			}
		}
	}
}

type delegateHandler struct {
	delegate http.Handler
}

func (h *delegateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.delegate.ServeHTTP(w, r)
}
