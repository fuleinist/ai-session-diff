package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fuleinist/ai-session-diff/internal/diff"
	"github.com/fuleinist/ai-session-diff/internal/session"
)

func Show(sessionID string) {
	s, err := session.Load(sessionID)
	if err != nil {
		fmt.Printf("✗ Session not found: %s\n", sessionID)
		fmt.Printf("  Run 'ai-session-diff list' to see available sessions\n")
		os.Exit(1)
	}

	beforeDir := filepath.Join(".ai-sessions", sessionID, "before")
	afterDir := filepath.Join(".ai-sessions", sessionID, "after")

	renderer := diff.TerminalRenderer{}
	output, err := renderer.RenderDiff(beforeDir, afterDir, s.FilesChanged)
	if err != nil {
		fmt.Printf("✗ Failed to generate diff: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(output)
}

func Report(sessionID string) {
	s, err := session.Load(sessionID)
	if err != nil {
		fmt.Printf("✗ Session not found: %s\n", sessionID)
		os.Exit(1)
	}

	beforeDir := filepath.Join(".ai-sessions", sessionID, "before")
	afterDir := filepath.Join(".ai-sessions", sessionID, "after")

	renderer := diff.HTMLRenderer{}
	html, err := renderer.RenderReport(beforeDir, afterDir, s)
	if err != nil {
		fmt.Printf("✗ Failed to generate report: %v\n", err)
		os.Exit(1)
	}

	reportPath := filepath.Join(".ai-sessions", sessionID, "session.html")
	if err := os.WriteFile(reportPath, []byte(html), 0644); err != nil {
		fmt.Printf("✗ Failed to write report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Report generated: %s\n", reportPath)
}

func List(aiOnly bool) {
	sessions, err := session.ListAll()
	if err != nil {
		fmt.Printf("✗ Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions recorded yet.")
		fmt.Println("Commit changes to start recording sessions!")
		return
	}

	fmt.Println("⚡ AI Session Diff — Session History")
	fmt.Println("─────────────────────────────────────────────────")

	for _, s := range sessions {
		if aiOnly && !s.AIDetected {
			continue
		}
		aiTag := ""
		if s.AIDetected {
			aiTag = fmt.Sprintf(" [%s]", s.AIProvider)
		}
		date := s.StartedAt
		if len(date) > 10 {
			date = date[:10]
		}
		fmt.Printf("  %s  %s  +%d -%d  %s%s\n",
			s.ID[:8],
			date,
			s.LinesAdded,
			s.LinesRemoved,
			s.CommitSha[:7],
			aiTag,
		)
	}
}

func Status() {
	prePath := filepath.Join(".git", "hooks", "pre-commit")
	postPath := filepath.Join(".git", "hooks", "post-commit")

	preInstalled := fileExists(prePath)
	postInstalled := fileExists(postPath)

	fmt.Println("⚡ AI Session Diff — Status")
	fmt.Println("─────────────────────────────")

	if preInstalled && postInstalled {
		fmt.Println("✓ Hooks installed")
		fmt.Printf("  Pre-commit:  %s\n", prePath)
		fmt.Printf("  Post-commit: %s\n", postPath)
	} else {
		fmt.Println("✗ Hooks not installed")
		fmt.Println("  Run 'ai-session-diff install' to enable")
	}

	sessions, _ := session.ListAll()
	fmt.Printf("\n  Sessions recorded: %d\n", len(sessions))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}