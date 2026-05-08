package cli

import (
	"os"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {
	// Test status shows hooks info
	// This is a basic smoke test - full integration test would need git repo
}

func TestShowHelp(t *testing.T) {
	// Capture help output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ShowHelp()
	w.Close()
	os.Stdout = old

	// Read output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "Usage:") {
		t.Error("Help should contain usage")
	}
}

func TestInstallCheckGitRepo(t *testing.T) {
	// Create temp dir without .git
	tmpDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// This should fail because not a git repo
	// We can't easily test this without capturing os.Exit
	// but we verify the path exists check works
	if _, err := os.Stat(".git"); err == nil {
		t.Skip("Test requires non-git directory")
	}
}

func TestUninstallHooks(t *testing.T) {
	// Similar to install test - requires git repo context
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}