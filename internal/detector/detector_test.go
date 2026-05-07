package detector

import (
	"os"
	"testing"
)

func TestDetectUnknown(t *testing.T) {
	// Clean environment
	os.Unsetenv("CLAUDE_CODE")
	os.Unsetenv("ANTHROPIC_MODEL")

	d := New()
	provider, detected := d.Detect()

	if detected {
		t.Error("Expected no AI provider detected in clean environment")
	}
	if provider != ProviderUnknown {
		t.Errorf("Expected unknown provider, got %s", provider)
	}
}

func TestName(t *testing.T) {
	d := New()

	tests := []struct {
		provider Provider
		expected string
	}{
		{ProviderCursor, "Cursor"},
		{ProviderClaudeCode, "Claude Code"},
		{ProviderGitHubCopilot, "GitHub Copilot"},
		{ProviderUnknown, "Unknown"},
	}

	for _, tc := range tests {
		name := d.Name(tc.provider)
		if name != tc.expected {
			t.Errorf("Name(%s) = %s, want %s", tc.provider, name, tc.expected)
		}
	}
}

func TestDetectClaudeCode(t *testing.T) {
	os.Setenv("CLAUDE_CODE", "true")
	defer os.Unsetenv("CLAUDE_CODE")

	d := New()
	provider, detected := d.Detect()

	if !detected {
		t.Error("Expected Claude Code to be detected")
	}
	if provider != ProviderClaudeCode {
		t.Errorf("Expected claude-code provider, got %s", provider)
	}
}

func TestDetectAnthropicModel(t *testing.T) {
	os.Setenv("ANTHROPIC_MODEL", "claude-sonnet-4-20250514")
	defer os.Unsetenv("ANTHROPIC_MODEL")

	d := New()
	provider, detected := d.Detect()

	if !detected {
		t.Error("Expected Claude Code to be detected via ANTHROPIC_MODEL")
	}
	if provider != ProviderClaudeCode {
		t.Errorf("Expected claude-code provider, got %s", provider)
	}
}