package proxy

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/log"
	"github.com/jetski-sh/mcp-proxy/oauth"
	"github.com/jetski-sh/mcp-proxy/webhook"
	"github.com/opencontainers/go-digest"
)

type mcpAwareTransport struct {
	Transport http.RoundTripper
	config    *config.Webhook
}

func (c *mcpAwareTransport) getTransport() http.RoundTripper {
	if c.Transport != nil {
		return c.Transport
	}
	return http.DefaultTransport
}

// RoundTrip implements http.RoundTripper.
func (c *mcpAwareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.config == nil || req.Method != http.MethodPost {
		return c.getTransport().RoundTrip(req)
	}

	log := log.Get(req.Context())
	webhookPayload := webhookPayloadFromReq(req)

	var reqGetter requestGetterWriter
	var respGetter responseGetterWriter
	var wg sync.WaitGroup
	wg.Add(2)

	if req.Body != nil {
		reqGetter = &rpcBufferGetter{}
		req.Body = &readCloser{
			Reader: io.TeeReader(req.Body, reqGetter),
			closeFunc: sync.OnceValue(func() error {
				wg.Done()
				return req.Body.Close()
			}),
		}
	}

	resp, err := c.getTransport().RoundTrip(req)

	if err != nil {
		webhookPayload.HttpError = err.Error()
	} else {
		webhookPayload.HttpStatusCode = resp.StatusCode

		switch resp.Header.Get("Content-Type") {
		case "application/json", "application/json; charset=utf-8":
			respGetter = &rpcBufferGetter{}
		case "text/event-stream":
			respGetter = newRPCEventStreamGetter()
		default:
			log.Info("unknown response content type",
				"contentType", resp.Header.Get("Content-Type"))
		}

		if respGetter != nil {
			resp.Body = &readCloser{
				Reader: io.TeeReader(resp.Body, respGetter),
				closeFunc: sync.OnceValue(func() error {
					wg.Done()
					return resp.Body.Close()
				}),
			}
		} else {
			wg.Done()
		}
	}

	go func() {
		wg.Wait()

		webhookPayload.Duration = time.Since(webhookPayload.StartedAt)

		if reqGetter == nil {
			log.Error(err, "request getter is not set")
		} else if r, err := reqGetter.GetJSONRPCRequest(); err != nil {
			log.Error(err, "could not get JSON-RPC request")
		} else {
			webhookPayload.MCPRequest = r
		}

		if respGetter == nil {
			log.Error(err, "response getter is not set")
		} else if r, err := respGetter.GetJSONRPCResponse(); err != nil {
			log.Error(err, "could not get JSON-RPC response")
		} else {
			webhookPayload.MCPResponse = r
		}

		log.Info("webhook payload assembled", "payload", webhookPayload)

		if err := webhook.Send(
			context.Background(),
			c.config.Method,
			c.config.Url.String(),
			webhookPayload,
		); err != nil {
			log.Error(err, "webhook error")
		}
	}()

	return resp, err
}

func webhookPayloadFromReq(req *http.Request) webhook.WebhookPayload {
	webhookPayload := webhook.WebhookPayload{
		MCPSessionID: req.Header.Get("Mcp-Session-Id"),
		StartedAt:    time.Now(),
		UserAgent:    req.UserAgent(),
	}

	if rawToken := oauth.GetRawToken(req.Context()); rawToken != "" {
		webhookPayload.AuthTokenDigest = digest.FromString(rawToken)
	}

	if token := oauth.GetToken(req.Context()); token != nil {
		webhookPayload.Subject, _ = token.Subject()
		var email string
		if err := token.Get("email", &email); err != nil {
			log.Get(req.Context()).Error(err, "could not get email claim from token")
		} else {
			webhookPayload.SubjectEmail = email
		}
	}

	return webhookPayload
}
