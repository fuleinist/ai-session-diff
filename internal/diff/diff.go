package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fuleinist/ai-session-diff/internal/session"
)

type DiffResult struct {
	File        string
	Additions   int
	Deletions   int
	Hunks       []Hunk
}

type Hunk struct {
	OldStart, OldLines int
	NewStart, NewLines int
	Lines              []string
}

type TerminalRenderer struct{}

func (r TerminalRenderer) RenderDiff(beforeDir, afterDir string, files []string) (string, error) {
	var out strings.Builder

	for _, file := range files {
		beforePath := filepath.Join(beforeDir, file)
		afterPath := filepath.Join(afterDir, file)

		beforeData, _ := os.ReadFile(beforePath)
		afterData, _ := os.ReadFile(afterPath)

		beforeLines := strings.Split(string(beforeData), "\n")
		afterLines := strings.Split(string(afterData), "\n")

		diffLines := computeUnifiedDiff(file, beforeLines, afterLines)
		out.WriteString(strings.Join(diffLines, "\n"))
		out.WriteString("\n\n")
	}

	return out.String(), nil
}

func computeUnifiedDiff(file string, oldLines, newLines []string) []string {
	var result []string

	result = append(result, fmt.Sprintf("--- a/%s", file))
	result = append(result, fmt.Sprintf("+++ b/%s", file))

	oldContent := strings.Join(oldLines, "\n")
	newContent := strings.Join(newLines, "\n")

	oldWords := strings.Fields(oldContent)
	newWords := strings.Fields(newContent)

	lcs := longestCommonSubsequence(oldWords, newWords)

	additions := len(newWords) - len(lcs)
	deletions := len(oldWords) - len(lcs)

	result = append(result, fmt.Sprintf("@@ -1,%d +1,%d @@", len(oldLines), len(newLines)))

	for _, w := range newWords {
		result = append(result, "\x1b[32m+"+w+"\x1b[0m")
	}
	for _, w := range oldWords {
		if !contains(lcs, w) {
			result = append(result, "\x1b[31m-"+w+"\x1b[0m")
		}
	}

	_ = additions
	_ = deletions

	return result
}

func longestCommonSubsequence(a, b []string) []string {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] > dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	var lcs []string
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			lcs = append([]string{a[i-1]}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return lcs
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type HTMLRenderer struct{}

func (r HTMLRenderer) RenderReport(beforeDir, afterDir string, s *session.Session) (string, error) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Session Diff - ` + s.ID + `</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Segoe UI', system-ui, sans-serif; background: #0d1117; color: #c9d1d9; min-height: 100vh; }
        .header { background: linear-gradient(135deg, #238636, #2ea44f); padding: 24px; color: white; }
        .header h1 { font-size: 24px; margin-bottom: 8px; }
        .meta { display: flex; gap: 24px; font-size: 14px; opacity: 0.9; }
        .meta span { display: flex; align-items: center; gap: 6px; }
        .content { padding: 24px; }
        .file-list { display: flex; flex-direction: column; gap: 16px; }
        .file-card { background: #161b22; border-radius: 8px; overflow: hidden; border: 1px solid #30363d; }
        .file-header { padding: 12px 16px; background: #21262d; display: flex; justify-content: space-between; align-items: center; cursor: pointer; }
        .file-name { font-family: 'Consolas', monospace; font-size: 14px; color: #58a6ff; }
        .file-stats { font-size: 12px; color: #8b949e; }
        .stats-add { color: #3fb950; }
        .stats-del { color: #f85149; }
        .diff-content { padding: 16px; font-family: 'Consolas', 'Monaco', monospace; font-size: 13px; line-height: 1.5; overflow-x: auto; }
        .line { white-space: pre; }
        .line-add { background: rgba(63, 185, 80, 0.15); color: #3fb950; }
        .line-del { background: rgba(248, 81, 73, 0.15); color: #f85149; }
        .line-context { color: #8b949e; }
        .ai-badge { background: #7c3aed; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-left: 8px; }
        .toggle-theme { position: fixed; top: 20px; right: 20px; background: #21262d; border: 1px solid #30363d; color: #c9d1d9; padding: 8px 16px; border-radius: 6px; cursor: pointer; }
        .toggle-theme:hover { background: #30363d; }
        .light body { background: #ffffff; color: #24292f; }
        .light .file-card { background: #f6f8fa; border-color: #d0d7de; }
        .light .file-header { background: #f6f8fa; }
        .light .line-context { color: #57606a; }
    </style>
</head>
<body>
    <button class="toggle-theme" onclick="document.body.classList.toggle('light')">🌓 Toggle Theme</button>
    <div class="header">
        <h1>⚡ AI Session Diff</h1>
        <div class="meta">
            <span>📅 ` + s.StartedAt + `</span>
            <span>🌿 ` + s.Branch + `</span>
            <span>🔖 ` + s.CommitSha[:7] + `</span>
            <span>+` + fmt.Sprintf("%d", s.LinesAdded) + `</span>
            <span>-` + fmt.Sprintf("%d", s.LinesRemoved) + `</span>
            ` + func() string {
		if s.AIDetected {
			return `<span class="ai-badge">🤖 ` + s.AIProvider + `</span>`
		}
		return ""
	}() + `
        </div>
    </div>
    <div class="content">
        <div class="file-list">
            <div class="file-card">
                <div class="file-header">
                    <span class="file-name">session-overview</span>
                    <span class="file-stats">
                        <span class="stats-add">+` + fmt.Sprintf("%d", s.LinesAdded) + `</span>
                        <span class="stats-del">-` + fmt.Sprintf("%d", s.LinesRemoved) + `</span>
                    </span>
                </div>
                <div class="diff-content">
                    <div class="line line-context">Files changed: ` + fmt.Sprintf("%d", len(s.FilesChanged)) + `</div>
                    <div class="line line-context">Session: ` + s.ID + `</div>
                    <div class="line line-context">Message: ` + s.Message + `</div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`
	return html, nil
}