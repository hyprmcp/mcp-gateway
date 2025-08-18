package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/jetski-sh/mcp-gateway/jsonrpc"
)

type requestGetterWriter interface {
	io.Writer

	GetJSONRPCRequest() (*jsonrpc.Request, error)
}

type responseGetterWriter interface {
	io.Writer

	GetJSONRPCResponse() (*jsonrpc.Response, error)
}

type rpcBufferGetter struct {
	bytes.Buffer
}

func (buf *rpcBufferGetter) GetJSONRPCRequest() (*jsonrpc.Request, error) {
	if msg, err := jsonrpc.ParseMessage(buf.Bytes()); err != nil {
		return nil, err
	} else if req, ok := msg.(*jsonrpc.Request); !ok {
		return nil, fmt.Errorf("not a JSON-RPC request")
	} else {
		return req, nil
	}
}

func (buf *rpcBufferGetter) GetJSONRPCResponse() (*jsonrpc.Response, error) {
	if msg, err := jsonrpc.ParseMessage(buf.Bytes()); err != nil {
		return nil, err
	} else if resp, ok := msg.(*jsonrpc.Response); !ok {
		return nil, fmt.Errorf("not a JSON-RPC response")
	} else {
		return resp, nil
	}
}

type rpcEventStreamGetter struct {
	EventStreamWriter
	events []Event
}

func newRPCEventStreamGetter() *rpcEventStreamGetter {
	var obj rpcEventStreamGetter
	obj.handler = func(e Event) {
		obj.events = append(obj.events, e)
	}
	return &obj
}

func (obj *rpcEventStreamGetter) GetJSONRPCResponse() (*jsonrpc.Response, error) {
	var errs []error

	for _, e := range obj.events {
		if e.Data == "" {
			continue
		}

		if msg, err := jsonrpc.ParseMessage([]byte(e.Data)); err != nil {
			errs = append(errs, fmt.Errorf("could not decode SSE data as JSON-RPC message: %w", err))
		} else if r, ok := msg.(*jsonrpc.Response); ok {
			return r, nil
		}
	}

	if len(errs) == 0 {
		return nil, errors.New("no JSON-RPC message found in SSE events")
	} else {
		return nil, errors.Join(errs...)
	}
}
