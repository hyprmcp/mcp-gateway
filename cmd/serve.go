package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/cors"
	"github.com/go-logr/stdr"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/htmlresponse"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/oauth"
	"github.com/hyprmcp/mcp-gateway/proxy"
	"github.com/hyprmcp/mcp-gateway/proxy/proxyutil"
	"github.com/spf13/cobra"
)

type ServeOptions struct {
	Config        string
	Addr          string
	AuthProxyAddr string
	Verbosity     int
}

func BindServeOptions(cmd *cobra.Command, opts *ServeOptions) {
	cmd.Flags().StringVarP(&opts.Config, "config", "c", "config.yaml", "Path to the configuration file")
	cmd.Flags().StringVarP(&opts.Addr, "addr", "a", ":9000", "Address to listen on")
	cmd.Flags().StringVar(&opts.AuthProxyAddr, "auth-proxy-addr", "", "Address to listen on with the authentication server proxy (advanced feature)")
	cmd.Flags().IntVarP(&opts.Verbosity, "verbosity", "v", 0, "Set the logging verbosity; greater number means more logging")
}

func runServe(ctx context.Context, opts ServeOptions) error {
	done := make(chan error)

	stdr.SetVerbosity(opts.Verbosity)

	cfg, err := config.ParseFile(opts.Config)
	if err != nil {
		return err
	}

	log.Get(ctx).Info("Loaded configuration", "config", cfg)

	if opts.AuthProxyAddr != "" {
		go func() {
			log.Get(ctx).Info("starting auth proxy server", "addr", opts.AuthProxyAddr)
			authUrl, err := url.Parse(cfg.Authorization.Server)
			if err != nil {
				done <- fmt.Errorf("auth proxy serve failed: %w", err)
			} else if err := http.ListenAndServe(opts.AuthProxyAddr, &httputil.ReverseProxy{Rewrite: proxyutil.RewriteHostFunc(authUrl)}); !errors.Is(err, http.ErrServerClosed) {
				done <- fmt.Errorf("auth proxy serve failed: %w", err)
			} else {
				done <- nil
			}
		}()
	}

	handler := &delegateHandler{}

	routerCtx, routerCancel := context.WithCancel(ctx)
	defer func() { routerCancel() }()

	if h, err := newRouter(routerCtx, cfg); err != nil {
		return err
	} else {
		handler.delegate = h
	}

	go func() {
		err := WatchConfigChanges(
			opts.Config,
			func(c *config.Config) {
				newRouterCtx, newRouterCancel := context.WithCancel(ctx)
				log.Get(ctx).Info("Reconfiguring server after config change...")
				if h, err := newRouter(newRouterCtx, c); err != nil {
					newRouterCancel()
					log.Get(ctx).Error(err, "failed to reload server")
				} else {
					routerCancel()
					routerCancel = newRouterCancel
					routerCtx = newRouterCtx
					handler.delegate = h
				}
			},
		)
		if err != nil {
			log.Get(ctx).Error(err, "config watch failed")
		}
	}()

	go func() {
		log.Get(ctx).Info("Starting server", "addr", opts.Addr)
		if err := http.ListenAndServe(opts.Addr, cors.AllowAll().Handler(handler)); !errors.Is(err, http.ErrServerClosed) {
			done <- fmt.Errorf("serve failed: %w", err)
		} else {
			done <- nil
		}
	}()

	return <-done
}

func newRouter(ctx context.Context, config *config.Config) (http.Handler, error) {
	mux := http.NewServeMux()

	htmlHandler := htmlresponse.NewHandler(config, false)
	oauthManager, err := oauth.NewManager(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := oauthManager.Register(mux); err != nil {
		return nil, err
	}

	for _, proxyConfig := range config.Proxy {
		if proxyConfig.Http != nil && proxyConfig.Http.Url != nil {
			handler := proxy.NewProxyHandler(&proxyConfig, oauthManager.UpdateWWWAuthenticateHeader)
			handler = htmlHandler.Handler(handler)

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
