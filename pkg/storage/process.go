package storage

import (
	"fmt"
	"sync"

	"awesomeProject/pkg/models"
)

type ProcessInfo struct {
	SessionID string
	Command   string
	Status    string // "running" or "completed"
	ExitCode  int
	Stdout    string
	Stderr    string
	mu        sync.Mutex
}

var (
	sessionsMu sync.RWMutex
	sessions   = make(map[string]*ProcessInfo)
)

func CreateSession(sessionID, command string) *ProcessInfo {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	info := &ProcessInfo{
		SessionID: sessionID,
		Command:   command,
		Status:    "running",
		ExitCode:  0,
	}
	sessions[sessionID] = info
	return info
}

func GetSession(sessionID string) (*ProcessInfo, error) {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()

	info, ok := sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	return info, nil
}

func ListSessions() []models.ProcessSession {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()

	var result []models.ProcessSession
	for _, info := range sessions {
		result = append(result, models.ProcessSession{
			SessionID: info.SessionID,
			Command:   info.Command,
			Status:    info.Status,
			ExitCode:  info.ExitCode,
			Stdout:    info.Stdout,
			Stderr:    info.Stderr,
		})
	}
	return result
}

func (pi *ProcessInfo) SetCompleted(exitCode int) {
	pi.mu.Lock()
	defer pi.mu.Unlock()
	pi.Status = "completed"
	pi.ExitCode = exitCode
}

func (pi *ProcessInfo) AppendStdout(data string) {
	pi.mu.Lock()
	defer pi.mu.Unlock()
	pi.Stdout += data
}

func (pi *ProcessInfo) AppendStderr(data string) {
	pi.mu.Lock()
	defer pi.mu.Unlock()
	pi.Stderr += data
}

func (pi *ProcessInfo) GetSnapshot() models.ProcessSession {
	pi.mu.Lock()
	defer pi.mu.Unlock()
	return models.ProcessSession{
		SessionID: pi.SessionID,
		Command:   pi.Command,
		Status:    pi.Status,
		ExitCode:  pi.ExitCode,
		Stdout:    pi.Stdout,
		Stderr:    pi.Stderr,
	}
}

// ClearSessions clears all sessions (used for testing)
func ClearSessions() {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	sessions = make(map[string]*ProcessInfo)
}
