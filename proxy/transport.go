package proxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jetski-sh/mcp-proxy/log"
	"github.com/jetski-sh/mcp-proxy/oauth"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/opencontainers/go-digest"
	"go.uber.org/multierr"
)

type mcpAwareTransport struct {
	Transport http.RoundTripper
}

// RoundTrip implements http.RoundTripper.
func (c *mcpAwareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mcpSessionID := req.Header.Get("Mcp-Session-Id")
	log := log.Get(req.Context()).WithValues("mcpSessionId", mcpSessionID)

	token := oauth.TokenFromContext(req.Context())
	sub, _ := token.Subject()
	if subDecoded, err := base64.RawStdEncoding.DecodeString(sub); err == nil {
		sub = string(subDecoded)
	}
	log = log.WithValues("subject", sub)
	var email string
	if err := token.Get("email", &email); err != nil {
		log.Error(err, "could not get email claim")
	} else {
		log = log.WithValues("email", email)
	}
	if st, err := jwt.NewSerializer().Serialize(token); err != nil {
		log.Error(err, "could not serialize JWT")
	} else {
		log = log.WithValues("tokenDigest", digest.FromBytes(st))
	}

	transport := c.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var reqBody, respBody bytes.Buffer

	if req.Body != nil {
		// retain Close functionality of the original body
		req.Body = &readCloser{
			Reader: io.TeeReader(req.Body, &reqBody),
			closeFunc: sync.OnceValue(func() error {
				wg.Done()
				return req.Body.Close()
			}),
		}
	}

	timeBefore := time.Now()

	resp, err := transport.RoundTrip(req)

	var events []Event

	if resp != nil && resp.Body != nil {
		var w io.Writer
		switch resp.Header.Get("Content-Type") {
		case "application/json", "application/json; charset=utf-8":
			w = &respBody
		case "text/event-stream":
			w = &EventStreamWriter{
				handler: func(event Event) { events = append(events, event) },
			}
		}

		resp.Body = &readCloser{
			Reader: io.TeeReader(resp.Body, w),
			closeFunc: sync.OnceValue(func() error {
				wg.Done()
				return resp.Body.Close()
			}),
		}
	}

	go func() {
		wg.Wait()

		timeElapsed := time.Since(timeBefore)

		var rpcRequest JSONRPCMessage
		if err := json.Unmarshal(reqBody.Bytes(), &rpcRequest); err != nil {
			log.Error(err, "could not decode request body", "data", reqBody.String())
		} else if rpcRequest.Method == "tools/call" {
			var toolParams mcp.CallToolParams
			var toolResult mcp.CallToolResult

			if err := json.Unmarshal(rpcRequest.Params, &toolParams); err != nil {
				log.Error(err, "could not decode JSON-RPC params as MCP tool call")
			} else if rpcResponse, err := getJSONRPCResponse(respBody.Bytes(), events); err != nil {
				log.Error(err, "could not get JSON-RPC response")
			} else if rpcResponse.Error != nil {
				log.Info("got JSON-RPC error", "code", rpcResponse.Error.Code, "message", rpcResponse.Error.Message)
			} else if err := json.Unmarshal(rpcResponse.Result, &toolResult); err != nil {
				log.Error(err, "could not decode JSON-RPC result as tool result")
			} else {
				log.Info("proxy request done",
					"jsonRpcId", rpcRequest.ID,
					"userAgent", req.Header.Get("User-Agent"),
					"timeElapsed", timeElapsed,
					"toolParams", toolParams,
					"toolResult", toolResult,
				)
			}
		} else {
			log.Info("not a tool call", "method", rpcRequest.Method)
		}
	}()

	return resp, err
}

func getJSONRPCResponse(body []byte, events []Event) (*JSONRPCMessage, error) {
	var msg JSONRPCMessage

	if len(body) > 0 {
		if err := json.Unmarshal(body, &msg); err != nil {
			return nil, fmt.Errorf("could not decode body as JSON-RPC message: %w", err)
		}

		return &msg, nil
	} else {
		var aggErr error

		for _, e := range events {
			if e.Data == "" {
				continue
			}

			if err := json.Unmarshal([]byte(e.Data), &msg); err != nil {
				multierr.AppendInto(&aggErr, fmt.Errorf("could not decode SSE data as JSON-RPC message: %w", err))
			} else if msg.IsResponse() {
				return &msg, nil
			}
		}

		if aggErr == nil {
			aggErr = errors.New("no JSON-RPC message found in SSE events")
		}

		return nil, aggErr
	}
}
