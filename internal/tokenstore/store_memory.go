package tokenstore

import (
	"context"
	"sync"
	"time"
)

type memEntry struct {
	userID    int64
	expiresAt time.Time
}

type memoryStore struct {
	mu     sync.RWMutex
	data   map[string]memEntry        // hash → entry
	byUser map[int64]map[string]bool  // userID → set of hashes
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		data:   make(map[string]memEntry),
		byUser: make(map[int64]map[string]bool),
	}
}

func (s *memoryStore) Save(_ context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[tokenHash] = memEntry{userID: userID, expiresAt: expiresAt}
	if s.byUser[userID] == nil {
		s.byUser[userID] = make(map[string]bool)
	}
	s.byUser[userID][tokenHash] = true
	return nil
}

func (s *memoryStore) Exists(_ context.Context, tokenHash string) (bool, error) {
	s.mu.RLock()
	entry, ok := s.data[tokenHash]
	s.mu.RUnlock()
	if !ok {
		return false, nil
	}
	if time.Now().After(entry.expiresAt) {
		// Lazy eviction of expired entry
		s.mu.Lock()
		delete(s.data, tokenHash)
		if set := s.byUser[entry.userID]; set != nil {
			delete(set, tokenHash)
		}
		s.mu.Unlock()
		return false, nil
	}
	return true, nil
}

func (s *memoryStore) Delete(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.data[tokenHash]; ok {
		delete(s.data, tokenHash)
		if set := s.byUser[entry.userID]; set != nil {
			delete(set, tokenHash)
		}
	}
	return nil
}

func (s *memoryStore) DeleteByUser(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for hash := range s.byUser[userID] {
		delete(s.data, hash)
	}
	delete(s.byUser, userID)
	return nil
}
