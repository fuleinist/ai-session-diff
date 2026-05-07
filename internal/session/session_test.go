package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNotFound(t *testing.T) {
	_, err := Load("nonexistent-session-id")
	if err == nil {
		t.Error("Expected error for nonexistent session")
	}
}

func TestListAllEmpty(t *testing.T) {
	// Create a fresh test directory
	os.MkdirAll(".ai-sessions", 0755)

	sessions, err := ListAll()
	if err != nil {
		t.Errorf("ListAll failed: %v", err)
	}
	if sessions == nil {
		sessions = []Session{}
	}
}

func TestSessionSummary(t *testing.T) {
	s := &Session{
		ID:           "abcd1234567890",
		StartedAt:    "2026-05-07T13:00:00Z",
		LinesAdded:   42,
		LinesRemoved: 7,
		AIProvider:   "cursor",
		AIDetected:   true,
	}

	summary := s.Summary()
	if len(summary) == 0 {
		t.Error("Summary should not be empty")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Clean up any existing test session
	testID := "test-save-load-123"
	os.RemoveAll(filepath.Join(".ai-sessions", testID))

	s := &Session{
		ID:           testID,
		CommitSha:    "abc123def456",
		Branch:       "main",
		LinesAdded:   10,
		LinesRemoved: 3,
		AIDetected:   true,
		AIProvider:   "cursor",
	}

	if err := Save(s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(testID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ID != s.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, s.ID)
	}
	if loaded.CommitSha != s.CommitSha {
		t.Errorf("CommitSha mismatch: got %s, want %s", loaded.CommitSha, s.CommitSha)
	}

	// Cleanup
	os.RemoveAll(filepath.Join(".ai-sessions", testID))
}