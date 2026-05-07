package detector

import (
	"os"
)

type Provider string

const (
	ProviderCursor       Provider = "cursor"
	ProviderClaudeCode    Provider = "claude-code"
	ProviderGitHubCopilot Provider = "github-copilot"
	ProviderUnknown       Provider = "unknown"
)

type Detector struct{}

func New() *Detector {
	return &Detector{}
}

func (d *Detector) Detect() (Provider, bool) {
	// Check for Cursor
	if _, err := os.Stat(".cursor"); err == nil {
		return ProviderCursor, true
	}

	// Check for Claude Code
	if os.Getenv("CLAUDE_CODE") != "" {
		return ProviderClaudeCode, true
	}
	if os.Getenv("ANTHROPIC_MODEL") != "" {
		return ProviderClaudeCode, true
	}

	// Check for GitHub Copilot
	if _, err := os.Stat(".github/copilot-innb"); err == nil {
		return ProviderGitHubCopilot, true
	}

	return ProviderUnknown, false
}

func (d *Detector) Name(p Provider) string {
	switch p {
	case ProviderCursor:
		return "Cursor"
	case ProviderClaudeCode:
		return "Claude Code"
	case ProviderGitHubCopilot:
		return "GitHub Copilot"
	default:
		return "Unknown"
	}
}