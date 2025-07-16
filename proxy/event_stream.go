package proxy

import (
	"bufio"
	"bytes"
	"strings"
)

type Event struct {
	Event string
	Data  string
	ID    string
	Retry string
}

type EventStreamWriter struct {
	buf     bytes.Buffer
	handler func(Event)
}

func (w *EventStreamWriter) Write(p []byte) (n int, err error) {
	n, _ = w.buf.Write(p)
	s := bufio.NewScanner(&w.buf)
	var event Event

	for s.Scan() {
		line := s.Text()

		switch {
		case line == "":
			// double EOL -> end of event, call handler
			if event.Event != "" || event.Data != "" || event.ID != "" {
				w.handler(event)
			}
			event = Event{}
		case line[0] == ':':
			// skip comments
		case strings.HasPrefix(line, "event:"):
			event.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			event.Data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		case strings.HasPrefix(line, "id:"):
			event.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		case strings.HasPrefix(line, "retry:"):
			event.Retry = strings.TrimSpace(strings.TrimPrefix(line, "retry:"))
		}
	}

	if event.Event != "" || event.Data != "" || event.ID != "" {
		w.handler(event)
	}

	return n, s.Err()
}
