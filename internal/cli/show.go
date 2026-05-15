package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

func List(aiOnly bool, since string) {
	sessions, err := session.ListAll()
	if err != nil {
		fmt.Printf("✗ Failed to list sessions: %v\n", err)
		os.Exit(1)
	}

	if since != "" {
		sinceTime, err := parseDate(since)
		if err != nil {
			fmt.Printf("✗ Invalid --since date format: %s (use YYYY-MM-DD)\n", since)
			os.Exit(1)
		}
		sessions = filterSessionsSince(sessions, sinceTime)
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

func filterSessionsSince(sessions []session.Session, since int64) []session.Session {
	var filtered []session.Session
	for _, s := range sessions {
		if t, err := parseTimestamp(s.StartedAt); err == nil && t >= since {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func parseDate(s string) (int64, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func parseTimestamp(s string) (int64, error) {
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Unix(), nil
	}
	// Try date-only (YYYY-MM-DD)
	if len(s) >= 10 {
		return parseDate(s[:10])
	}
	return 0, fmt.Errorf("unparseable timestamp: %s", s)
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