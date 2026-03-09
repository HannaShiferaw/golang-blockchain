package memstore

import (
	"context"
	"sync"
)

type Store struct {
	mu sync.RWMutex
	m  map[string][]byte
}

func New() *Store {
	return &Store{m: map[string][]byte{}}
}

func (s *Store) Get(_ context.Context, key string) ([]byte, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	if !ok {
		return nil, false, nil
	}
	cp := make([]byte, len(v))
	copy(cp, v)
	return cp, true, nil
}

func (s *Store) Put(_ context.Context, key string, raw []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]byte, len(raw))
	copy(cp, raw)
	s.m[key] = cp
	return nil
}

