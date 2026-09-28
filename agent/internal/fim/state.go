package fim

import "sync"

type state struct {
	mu sync.RWMutex
	hashes map[string]string
}

func newState() *state {
	return &state{
		hashes: make(map[string]string),
	}
}

func (s *state) get(path string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	hash, ok := s.hashes[path]
	return hash, ok
}

func (s *state) set(path string, hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashes[path] = hash
}

func (s *state) delete(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.hashes, path)
}

	