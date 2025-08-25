package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/cors"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/oauth"
	"github.com/hyprmcp/mcp-gateway/proxy"
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

	if h, err := newRouter(ctx, cfg); err != nil {
		return err
	} else {
		handler.delegate = h
	}

	go func() {
		err := WatchConfigChanges(
			opts.Config,
			func(c *config.Config) {
				log.Get(ctx).Info("Reconfiguring server after config change...")
				if h, err := newRouter(ctx, c); err != nil {
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

func newRouter(ctx context.Context, config *config.Config) (http.Handler, error) {
	mux := http.NewServeMux()

	oauthManager, err := oauth.NewManager(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := oauthManager.Register(mux); err != nil {
		return nil, err
	}

	for _, proxyConfig := range config.Proxy {
		if proxyConfig.Http != nil && proxyConfig.Http.Url != nil {
			handler := proxy.NewProxyHandler(&proxyConfig)

			if proxyConfig.Authentication.Enabled {
				handler = oauthManager.Handler(handler)
			}

			mux.Handle(proxyConfig.Path, handler)
		}
	}

	return mux, nil
}

func WatchConfigChanges(path string, callback func(*config.Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func() { _ = watcher.Close() }()

	// We watch the parent directory of the config file, rather than just the file itself, because Kubernetes uses
	// symlinks when mounting ConfigMaps/Secrets and just watching the file doesn't work well in those cases.
	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	if err := watcher.Add(fileDir); err != nil {
		return fmt.Errorf("failed to watch directory %s: %w", fileDir, err)
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

			// Check if the event is for our config file
			if filepath.Base(event.Name) == fileName {
				log.Root().Info("starting config reload", "op", event.Op, "path", event.Name)

				if cfg, err := config.ParseFile(path); err != nil {
					log.Root().Error(err, "config reload error", "event", event)
				} else {
					callback(cfg)
				}
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
