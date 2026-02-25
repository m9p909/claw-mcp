package session

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// SessionStore manages MCP protocol sessions
type SessionStore struct {
	sessions map[string]*SessionMetadata
	mu       sync.RWMutex
}

// SessionMetadata tracks session information
type SessionMetadata struct {
	SessionID      string
	CreatedAt      time.Time
	LastActivityAt time.Time
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*SessionMetadata),
	}
}

// GenerateSessionID creates a cryptographically secure session ID
func GenerateSessionID() (string, error) {
	// Generate 16 bytes (128 bits) of random data
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Encode as hex to ensure visible ASCII characters (0x21-0x7E)
	return hex.EncodeToString(b), nil
}

// CreateSession creates a new session and returns its ID
func (s *SessionStore) CreateSession() (string, error) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[sessionID] = &SessionMetadata{
		SessionID:      sessionID,
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}

	return sessionID, nil
}

// ValidateSession checks if a session exists and updates last activity time
func (s *SessionStore) ValidateSession(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return false
	}

	// Update last activity time
	session.LastActivityAt = time.Now()
	return true
}

// TerminateSession removes a session from the store
func (s *SessionStore) TerminateSession(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionID]; exists {
		delete(s.sessions, sessionID)
		return true
	}

	return false
}

// GetSession returns session metadata (read-only)
func (s *SessionStore) GetSession(sessionID string) *SessionMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sessions[sessionID]
}
