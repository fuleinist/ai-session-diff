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

type LineClass int

const (
	LineContext LineClass = iota
	LineAdd
	LineDel
)

type DiffLine struct {
	Text  string
	Class LineClass
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

	lcs := longestCommonSubsequence(oldLines, newLines)

	// Build diff using LCS-based algorithm
	i, j := 0, 0
	lcsIdx := 0
	
	for lcsIdx < len(lcs) {
		// Skip deletions in oldLines
		for i < len(oldLines) && oldLines[i] != lcs[lcsIdx] {
			result = append(result, "\x1b[31m-"+oldLines[i]+"\x1b[0m")
			i++
		}
		// Skip additions in newLines
		for j < len(newLines) && newLines[j] != lcs[lcsIdx] {
			result = append(result, "\x1b[32m+"+newLines[j]+"\x1b[0m")
			j++
		}
		// Context line (in LCS)
		if i < len(oldLines) && j < len(newLines) {
			result = append(result, " "+oldLines[i])
			i++
			j++
			lcsIdx++
		}
	}
	
	// Remaining deletions in oldLines
	for i < len(oldLines) {
		result = append(result, "\x1b[31m-"+oldLines[i]+"\x1b[0m")
		i++
	}
	// Remaining additions in newLines
	for j < len(newLines) {
		result = append(result, "\x1b[32m+"+newLines[j]+"\x1b[0m")
		j++
	}

	return result
}

// fileDiffLines returns the diff lines for a single file, with class labels.
func fileDiffLines(file, beforeDir, afterDir string) []DiffLine {
	beforePath := filepath.Join(beforeDir, file)
	afterPath := filepath.Join(afterDir, file)

	beforeData, err := os.ReadFile(beforePath)
	if err != nil {
		return []DiffLine{{Text: "(could not read before)", Class: LineContext}}
	}
	afterData, err := os.ReadFile(afterPath)
	if err != nil {
		return []DiffLine{{Text: "(could not read after)", Class: LineContext}}
	}

	beforeLines := strings.Split(string(beforeData), "\n")
	afterLines := strings.Split(string(afterData), "\n")

	lcs := longestCommonSubsequence(beforeLines, afterLines)

	var lines []DiffLine
	for _, line := range beforeLines {
		if !contains(lcs, line) {
			lines = append(lines, DiffLine{Text: "-" + line, Class: LineDel})
		} else {
			lines = append(lines, DiffLine{Text: " " + line, Class: LineContext})
		}
	}
	for _, line := range afterLines {
		if !contains(lcs, line) {
			lines = append(lines, DiffLine{Text: "+" + line, Class: LineAdd})
		}
	}

	return lines
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
	// Build per-file diff blocks
	var fileBlocks strings.Builder
	for _, file := range s.FilesChanged {
		lines := fileDiffLines(file, beforeDir, afterDir)
		fileBlocks.WriteString(fmt.Sprintf(`            <div class="file-card">
                <div class="file-header">
                    <span class="file-name">%s</span>
                </div>
                <div class="diff-content">
`, file))
		for _, line := range lines {
			class := "line-context"
			if line.Class == LineAdd {
				class = "line line-add"
			} else if line.Class == LineDel {
				class = "line line-del"
			}
			// Escape HTML entities in line text
			escaped := strings.ReplaceAll(line.Text, "&", "&amp;")
			escaped = strings.ReplaceAll(escaped, "<", "&lt;")
			escaped = strings.ReplaceAll(escaped, ">", "&gt;")
			fileBlocks.WriteString(fmt.Sprintf(`                    <div class="%s">%s</div>
`, class, escaped))
		}
		fileBlocks.WriteString(`                </div>
            </div>
`)
	}

	aiBadge := ""
	if s.AIDetected {
		aiBadge = fmt.Sprintf(`<span class="ai-badge">🤖 %s</span>`, s.AIProvider)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Session Diff - %s</title>
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
            <span>📅 %s</span>
            <span>🌿 %s</span>
            <span>🔖 %s</span>
            <span>+%d</span>
            <span>-%d</span>
            %s
        </div>
    </div>
    <div class="content">
        <div class="file-list">
%s        </div>
    </div>
</body>
</html>`, s.ID, s.StartedAt, s.Branch, s.CommitSha[:7], s.LinesAdded, s.LinesRemoved, aiBadge, fileBlocks.String())

	return html, nil
}
