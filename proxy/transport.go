package proxy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/jsonrpc"
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

	var wg sync.WaitGroup

	isToolsListRequest := false

	if req.Body != nil {
		defer req.Body.Close()
		if data, err := io.ReadAll(req.Body); err != nil {
			// TODO think about error handling some more
			return nil, err
		} else if rpcMsg, err := jsonrpc.ParseMessage(data); err == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(data))
			log.Error(err, "body parse error")
		} else if rpcReq, ok := rpcMsg.(*jsonrpc.Request); !ok {
			req.Body = io.NopCloser(bytes.NewBuffer(data))
		} else {
			webhookPayload.MCPRequest = rpcReq
			isToolsListRequest = rpcReq.Method == "tools/list"

			if rpcReq.Method == "tools/call" {
				newReq := &jsonrpc.Request{
					ID:     rpcReq.ID,
					Method: rpcReq.Method,
					// TODO: Parse the request params as MCP tool call request, remove the telemetry arguments and set req.Body to a new
					// buffer of the modified request
					Params:      rpcReq.Params,
					Notif:       rpcReq.Notif,
					Meta:        rpcReq.Meta,
					ExtraFields: rpcReq.ExtraFields,
				}
				if newData, err := json.Marshal(newReq); err != nil {
					log.Error(err, "failed to marshal rpc request")
				} else {
					data = newData
				}
			}

			req.Body = io.NopCloser(bytes.NewBuffer(data))
		}
	}

	resp, err := c.getTransport().RoundTrip(req)

	if err != nil {
		webhookPayload.HttpError = err.Error()
	} else {
		webhookPayload.HttpStatusCode = resp.StatusCode

		switch resp.Header.Get("Content-Type") {
		case "application/json", "application/json; charset=utf-8":
			defer resp.Body.Close()

			if data, err := io.ReadAll(resp.Body); err != nil {
				// TODO think about error handling some more
				return nil, err
			} else if rpcMsg, err := jsonrpc.ParseMessage(data); err != nil {
				log.Error(nil, "failed to parse JSONRPC message")
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
			} else if rpcResp, ok := rpcMsg.(*jsonrpc.Response); !ok {
				log.Error(err, "not a JSONRPC response")
				resp.Body = io.NopCloser(bytes.NewBuffer(data))
			} else {
				webhookPayload.MCPResponse = rpcResp

				if isToolsListRequest {
					newResp := &jsonrpc.Response{
						ID: rpcResp.ID,
						// TODO: parse the rpcResponse as MCP tools list response and add the telemetry arguments to all tools
						Result: rpcResp.Result,
						Error:  rpcResp.Error,
						Meta:   rpcResp.Meta,
					}

					if newData, err := json.Marshal(newResp); err != nil {
						log.Error(err, "failed to serialize modified JSONRPC response")
					} else {
						data = newData
					}
				}

				resp.Body = io.NopCloser(bytes.NewBuffer(data))
			}
		case "text/event-stream":
			wg.Add(1)

			resp.Body = &eventStreamReader{
				s: bufio.NewScanner(resp.Body),
				mutateFunc: func(e Event) Event {
					// TODO: parse the event data as JSONRPC response and add the telemetry arguments to all tools
					if rpcMsg, err := jsonrpc.ParseMessage([]byte(e.Data)); err != nil {
						log.Error(err, "failed to parse JSONRPC message")
					} else if rpcResp, ok := rpcMsg.(*jsonrpc.Response); !ok {
						log.Error(nil, "not a JSONRPC response")
					} else {
						webhookPayload.MCPResponse = rpcResp

						if isToolsListRequest {
							newResp := &jsonrpc.Response{
								ID: rpcResp.ID,
								// TODO: parse the rpcResponse as MCP tools list response and add the telemetry arguments to all tools
								Result: rpcResp.Result,
								Error:  rpcResp.Error,
								Meta:   rpcResp.Meta,
							}

							if newData, err := json.Marshal(newResp); err == nil {
								e.Data = string(newData)
							} else {
								log.Error(err, "failed to serialize modified JSONRPC response")
							}
						}
					}

					return e
				},
				closeFunc: func() error {
					wg.Done()
					return resp.Body.Close()
				},
			}
		default:
			log.Info("unknown response content type",
				"contentType", resp.Header.Get("Content-Type"))
		}
	}

	go func() {
		wg.Wait()

		webhookPayload.Duration = time.Since(webhookPayload.StartedAt)

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
