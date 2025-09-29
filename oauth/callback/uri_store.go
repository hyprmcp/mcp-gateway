package callback

import (
	"errors"
	"net/url"
	"sync"
)

type URIStore interface {
	Get(state string) (*url.URL, error)
	Set(state string, url string)
}

type store struct {
	data map[string]string
	lock *sync.RWMutex
}

func NewURIStore() URIStore {
	return &store{
		data: make(map[string]string),
		lock: new(sync.RWMutex),
	}
}

func (s *store) Get(state string) (*url.URL, error) {
	if state == "" {
		return nil, errors.New("empty state")
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	if uri, ok := s.data[state]; ok {
		return url.Parse(uri)
	}

	return nil, errors.New("no redirect_uri found")
}

func (s *store) Set(state string, url string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[state] = url
}
