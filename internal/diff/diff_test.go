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

func TestComputeUnifiedDiff_Order(t *testing.T) {
	// Test that deletions appear before additions when line positions differ
	oldLines := []string{"line1", "line2", "line3", "line4", "line5"}
	newLines := []string{"line1", "line2", "modified_line3", "line4", "line5", "added_line6"}

	diffLines := computeUnifiedDiff("test.txt", oldLines, newLines)

	// Check that we have both deletions and additions
	hasDeletion := false
	hasAddition := false
	for _, line := range diffLines {
		if strings.HasPrefix(line, "\x1b[31m-") {
			hasDeletion = true
		}
		if strings.HasPrefix(line, "\x1b[32m+") {
			hasAddition = true
		}
	}

	if !hasDeletion {
		t.Error("Expected at least one deletion line")
	}
	if !hasAddition {
		t.Error("Expected at least one addition line")
	}

	// Find positions of first deletion and first addition
	firstDelIdx := -1
	firstAddIdx := -1
	for i, line := range diffLines {
		if firstDelIdx == -1 && strings.HasPrefix(line, "\x1b[31m-") {
			firstDelIdx = i
		}
		if firstAddIdx == -1 && strings.HasPrefix(line, "\x1b[32m+") {
			firstAddIdx = i
		}
	}

	// Deletion should come before addition (unified diff order)
	if firstDelIdx != -1 && firstAddIdx != -1 {
		if firstDelIdx > firstAddIdx {
			t.Errorf("Diff order wrong: first deletion at %d, first addition at %d (deletion should come first)", firstDelIdx, firstAddIdx)
		}
	}
}
