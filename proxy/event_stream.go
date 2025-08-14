package proxy

import (
	"bytes"
	"fmt"
)

type Event struct {
	Event string
	Data  string
	ID    string
	Retry string
}

func (e *Event) Bytes() []byte {
	buf := new(bytes.Buffer)

	if e.ID != "" {
		fmt.Fprintln(buf, "id: "+e.ID)
	}

	if e.Event != "" {
		fmt.Fprintln(buf, "event: "+e.Event)
	}

	if e.Data != "" {
		fmt.Fprintln(buf, "data: "+e.Data)
	}

	if e.Retry != "" {
		fmt.Fprintln(buf, "retry: "+e.Retry)
	}

	fmt.Fprintln(buf)

	return buf.Bytes()
}
