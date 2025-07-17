package jsonrpc

import (
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
)

type Message any
type Request = jsonrpc2.Request
type Response = jsonrpc2.Response

// ParseMessage parses a JSON-RPC message, returning either a *Request or *Response.
func ParseMessage(data []byte) (Message, error) {
	var probe struct {
		Method string `json:"method"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	} else if probe.Method != "" {
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		return &req, nil
	} else {
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}
}
