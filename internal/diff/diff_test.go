package diff

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fuleinist/ai-session-diff/internal/session"
)

func TestTerminalRendererRenderDiff(t *testing.T) {
	// Create temp directories
	beforeDir, _ := os.MkdirTemp("", "before")
	afterDir, _ := os.MkdirTemp("", "after")
	defer os.RemoveAll(beforeDir)
	defer os.RemoveAll(afterDir)

	// Create test files
	beforeFile := filepath.Join(beforeDir, "test.go")
	afterFile := filepath.Join(afterDir, "test.go")
	os.WriteFile(beforeFile, []byte("line1\nline2\nline3\n"), 0644)
	os.WriteFile(afterFile, []byte("line1\nmodified\nline3\n"), 0644)

	renderer := TerminalRenderer{}
	output, err := renderer.RenderDiff(beforeDir, afterDir, []string{"test.go"})
	if err != nil {
		t.Fatalf("RenderDiff failed: %v", err)
	}
	if output == "" {
		t.Error("Expected non-empty output")
	}
}

func TestHTMLRendererRenderReport(t *testing.T) {
	beforeDir, _ := os.MkdirTemp("", "before")
	afterDir, _ := os.MkdirTemp("", "after")
	defer os.RemoveAll(beforeDir)
	defer os.RemoveAll(afterDir)

	renderer := HTMLRenderer{}
	html, err := renderer.RenderReport(beforeDir, afterDir, &session.Session{
		ID:           "test-123",
		StartedAt:    "2026-05-08T10:00:00Z",
		CommitSha:    "abc123def456",
		Branch:       "main",
		LinesAdded:   10,
		LinesRemoved: 5,
		FilesChanged: []string{"test.go"},
		Message:      "test commit",
	})
	if err != nil {
		t.Fatalf("RenderReport failed: %v", err)
	}
	if html == "" {
		t.Error("Expected non-empty HTML")
	}
	if !strings.Contains(html, "AI Session Diff") {
		t.Error("HTML should contain title")
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	b := []string{"b", "c", "e"}
	lcs := longestCommonSubsequence(a, b)
	if len(lcs) != 2 {
		t.Errorf("Expected LCS length 2, got %d", len(lcs))
	}
}
