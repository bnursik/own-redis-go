package storage

import (
	"sync"
	"time"
)

type Entry struct {
	value      string
	expiration time.Time
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Entry
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]Entry),
	}
}

func (s *Store) Set(key string, value string, expireMs int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiration time.Time
	if expireMs > 0 {
		expiration = time.Now().Add(time.Duration(expireMs) * time.Millisecond)
	}


	s.items[key] = Entry{
		value:      value,
		expiration: expiration,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.items[key]
	if !exists {
		return "", false
	}

	if !entry.expiration.IsZero() && time.Now().After(entry.expiration) {
		delete(s.items, key)
		return "", false
	}

	return entry.value, true
}

func (s *Store) Ping() string {
	return "PONG"
}
