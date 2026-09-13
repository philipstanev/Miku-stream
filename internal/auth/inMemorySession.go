package auth

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type InMemorySession struct {
	mu   sync.Mutex
	data map[string]SessionInfo
}

func NewInMemorySession() *InMemorySession {
	return &InMemorySession{
		data: make(map[string]SessionInfo),
	}
}

func (s *InMemorySession) GetInfo(sessionID string) (SessionInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.data[sessionID]
	if ok {
		return data, nil
	}
	return SessionInfo{}, errors.New("User not found")
}

func (s *InMemorySession) PutUser(info SessionInfo) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	SessionID := uuid.New().String()
	s.data[SessionID] = info
	return SessionID, nil

}

func (s *InMemorySession) RemoveUser(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, sessionID)
	return nil
}
