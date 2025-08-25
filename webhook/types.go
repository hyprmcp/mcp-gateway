package webhook

import (
	"time"

	"github.com/hyprmcp/mcp-gateway/jsonrpc"
	"github.com/opencontainers/go-digest"
)

type WebhookPayload struct {
	Subject         string            `json:"subject"`
	SubjectEmail    string            `json:"subjectEmail"`
	MCPSessionID    string            `json:"mcpSessionId"`
	StartedAt       time.Time         `json:"startedAt"`
	Duration        time.Duration     `json:"duration"`
	AuthTokenDigest digest.Digest     `json:"authTokenDigest"`
	MCPRequest      *jsonrpc.Request  `json:"mcpRequest,omitempty"`
	MCPResponse     *jsonrpc.Response `json:"mcpResponse,omitempty"`
	UserAgent       string            `json:"userAgent"`
	HttpStatusCode  int               `json:"httpStatusCode,omitempty"`
	HttpError       string            `json:"httpError,omitempty"`
}
