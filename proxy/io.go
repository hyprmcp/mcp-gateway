package proxy

import (
	"io"
)

type readCloser struct {
	io.Reader
	closeFunc func() error
}

// Close implements io.Closer.
func (rc *readCloser) Close() error {
	if rc.closeFunc != nil {
		return rc.closeFunc()
	}
	return nil
}
