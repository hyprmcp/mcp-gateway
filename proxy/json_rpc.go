package proxy

import "bytes"

// TODO: Maybe switch to github.com/modelcontextprotocol/go-sdk/jsonrpc once it properly exposes all needed functions
type JSONRPCMessage struct {
	ID any `json:"id"`

	// JSONRPC request values

	Method string  `json:"method"`
	Params RawJSON `json:"params"`

	// JSONRPC response values

	Result RawJSON       `json:"result"`
	Error  *JSONRPCError `json:"error"`
}

func (msg *JSONRPCMessage) IsRequest() bool {
	return msg.Method != "" && msg.ID != nil
}

func (msg *JSONRPCMessage) IsNotification() bool {
	return msg.Method != "" && msg.ID == nil
}

func (msg *JSONRPCMessage) IsResponse() bool {
	return msg.Result != nil || msg.Error != nil
}

type JSONRPCError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    RawJSON `json:"data"`
}

type RawJSON []byte

func (obj *RawJSON) UnmarshalJSON(data []byte) error {
	*obj = bytes.Clone(data)
	return nil
}

func (obj *RawJSON) MarshalJSON() ([]byte, error) {
	return *obj, nil
}
