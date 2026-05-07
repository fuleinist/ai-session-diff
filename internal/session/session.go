package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Session struct {
	ID            string   `json:"id"`
	StartedAt     string   `json:"startedAt"`
	EndedAt       string   `json:"endedAt"`
	CommitSha     string   `json:"commitSha"`
	Branch        string   `json:"branch"`
	FilesChanged  []string `json:"filesChanged"`
	AIProvider    string   `json:"aiProvider"`
	AIDetected    bool     `json:"aiDetected"`
	LinesAdded    int      `json:"linesAdded"`
	LinesRemoved  int      `json:"linesRemoved"`
	Message       string   `json:"message"`
}

func Load(sessionID string) (*Session, error) {
	path := filepath.Join(".ai-sessions", sessionID, "session.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid session file: %w", err)
	}
	return &s, nil
}

func ListAll() ([]Session, error) {
	sessionDir := ".ai-sessions"
	if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
		return []Session{}, nil
	}

	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		sessPath := filepath.Join(sessionDir, entry.Name(), "session.json")
		if data, err := os.ReadFile(sessPath); err == nil {
			var s Session
			if json.Unmarshal(data, &s) == nil {
				sessions = append(sessions, s)
			}
		}
	}

	return sessions, nil
}

func Save(s *Session) error {
	os.MkdirAll(filepath.Join(".ai-sessions", s.ID), 0755)
	path := filepath.Join(".ai-sessions", s.ID, "session.json")
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Session) Summary() string {
	date := s.StartedAt
	if len(date) > 10 {
		date = date[:10]
	}
	aiTag := ""
	if s.AIDetected {
		aiTag = fmt.Sprintf(" [%s]", s.AIProvider)
	}
	return fmt.Sprintf("%s %s +%d -%d%s",
		s.ID[:8], date, s.LinesAdded, s.LinesRemoved, aiTag)
}