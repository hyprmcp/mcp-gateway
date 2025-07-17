package webhook

import (
	"time"

	"github.com/jetski-sh/mcp-proxy/jsonrpc"
	"github.com/opencontainers/go-digest"
)

type WebhookPayload struct {
	Subject          string            `json:"subject"`
	SubjectEmail     string            `json:"subjectEmail"`
	MCPSessionID     string            `json:"mcpSessionId"`
	TokenDigest      digest.Digest     `json:"tokenDigest"`
	UserAgent        string            `json:"userAgent"`
	HttpError        string            `json:"httpError,omitempty"`
	HttpResponseCode int               `json:"httpResponseCode,omitempty"`
	RPCRequest       *jsonrpc.Request  `json:"rpcRequest,omitempty"`
	RPCResponse      *jsonrpc.Response `json:"rpcResponse,omitempty"`
	StartTime        time.Time         `json:"startTime"`
	Duration         time.Duration     `json:"duration"`
}
