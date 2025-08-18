package proxy

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type eventStreamReader struct {
	buf        bytes.Buffer
	s          *bufio.Scanner
	mutateFunc func(Event) Event
	closeFunc  func() error
	err        error
	events     []Event
}

func (rw *eventStreamReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// If the buffer is empty, and the scanner has not reached EOF, read from the scanner until a full SSE event has been
	// read. Then, call the mutate function on the event and write the mutated SSE event to the buffer.
	if rw.buf.Len() == 0 {
		if rw.err != nil {
			return 0, rw.err
		}

		// will be set to false if we break out of the loop early
		scannerHasEOF := true
		var event Event
	scanLoop:
		for rw.s.Scan() {
			line := rw.s.Text()

			switch {
			case line == "":
				// double EOL -> end of event
				scannerHasEOF = false
				break scanLoop
			case line[0] == ':':
				// skip comments
			case strings.HasPrefix(line, "event:"):
				event.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			case strings.HasPrefix(line, "data:"):
				// When the EventSource receives multiple consecutive lines that begin with data:, it concatenates them, inserting
				// a newline character between each one. Trailing newlines are removed.
				dataLine := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if event.Data == "" {
					event.Data = dataLine
				} else {
					event.Data += "\n" + dataLine
				}
			case strings.HasPrefix(line, "id:"):
				event.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
			case strings.HasPrefix(line, "retry:"):
				event.Retry = strings.TrimSpace(strings.TrimPrefix(line, "retry:"))
			}
		}

		if event.Event != "" || event.Data != "" || event.ID != "" {
			if rw.mutateFunc != nil {
				event = rw.mutateFunc(event)
			}

			// Write on bytes.Buffer never returns an error
			_, _ = rw.buf.Write(event.Bytes())
			rw.events = append(rw.events, event)
		}

		// bufio.Scanner does not expose EOF errors, so we need to remember whether the loop has ended due to EOF, another
		// error or because a full event has been read.
		if err := rw.s.Err(); err != nil {
			rw.err = err
		} else if scannerHasEOF {
			rw.err = io.EOF
		}
	}

	return rw.buf.Read(p)
}

func (rw *eventStreamReader) Close() error {
	if rw.closeFunc != nil {
		return rw.closeFunc()
	}

	return nil
}
