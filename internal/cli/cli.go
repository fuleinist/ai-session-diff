package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	hooksDir     = ".git/hooks"
	sessionDir   = ".ai-sessions"
	preCommitHook = `#!/bin/sh
# AI Session Diff - Pre-commit hook
# This hook snapshots file state before commit

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
SESSION_DIR="$REPO_DIR/.ai-sessions"

if [ ! -d "$SESSION_DIR" ]; then
    mkdir -p "$SESSION_DIR"
fi

# Generate session ID
SESSION_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "pre-$(date +%s)")
BEFORE_DIR="$SESSION_DIR/$SESSION_ID/before"
mkdir -p "$BEFORE_DIR"

# Snapshot all staged files
git diff --cached --name-only | while IFS= read -r file; do
    if [ -f "$REPO_DIR/$file" ]; then
        mkdir -p "$BEFORE_DIR/$(dirname "$file")"
        cp "$REPO_DIR/$file" "$BEFORE_DIR/$file"
    fi
done

# Write pre-commit marker with session ID
echo "$SESSION_ID" > "$SESSION_DIR/.current-session"
exit 0
`

	postCommitHook = `#!/bin/sh
# AI Session Diff - Post-commit hook
# This hook records session metadata after commit

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
SESSION_DIR="$REPO_DIR/.ai-sessions"

if [ ! -f "$SESSION_DIR/.current-session" ]; then
    exit 0
fi

SESSION_ID=$(cat "$SESSION_DIR/.current-session")
rm -f "$SESSION_DIR/.current-session"

AFTER_DIR="$SESSION_DIR/$SESSION_ID/after"
mkdir -p "$AFTER_DIR"

# Snapshot all committed files
git ls-files | while IFS= read -r file; do
    if [ -f "$REPO_DIR/$file" ]; then
        mkdir -p "$AFTER_DIR/$(dirname "$file")"
        cp "$REPO_DIR/$file" "$AFTER_DIR/$file"
    fi
done

# Get commit info
COMMIT_SHA=$(git rev-parse HEAD 2>/dev/null)
BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
COMMIT_MSG=$(git log -1 --format=%s 2>/dev/null || echo "")
COMMIT_TIME=$(git log -1 --format=%cI 2>/dev/null || echo "")

# Count changes
LINES_ADDED=$(git diff --cached --numstat | awk '{add += $1} END {print add+0}')
LINES_REMOVED=$(git diff --cached --numstat | awk '{sub += $2} END {print sub+0}')

# Files changed
FILES_CHANGED=$(git diff --cached --name-only | tr '\n' ' ')

# Write metadata
cat > "$SESSION_DIR/$SESSION_ID/session.json" << EOF
{
  "id": "$SESSION_ID",
  "startedAt": "$COMMIT_TIME",
  "endedAt": "$COMMIT_TIME",
  "commitSha": "$COMMIT_SHA",
  "branch": "$BRANCH",
  "filesChanged": [${FILES_CHANGED:0:-1}],
  "aiProvider": "$([ -d "$REPO_DIR/.cursor" ] && echo "cursor" || echo "unknown")",
  "aiDetected": $([ -d "$REPO_DIR/.cursor" ] && echo "true" || echo "false"),
  "linesAdded": $LINES_ADDED,
  "linesRemoved": $LINES_REMOVED,
  "message": "$COMMIT_MSG"
}
EOF

exit 0
`
)

func Install() {
	// Check we're in a git repo
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("✗ Not a git repository. Run from repo root.")
		os.Exit(1)
	}

	// Create hooks dir
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		fmt.Printf("✗ Failed to create hooks directory: %v\n", err)
		os.Exit(1)
	}

	// Write pre-commit hook
	prePath := filepath.Join(hooksDir, "pre-commit")
	if err := os.WriteFile(prePath, []byte(preCommitHook), 0755); err != nil {
		fmt.Printf("✗ Failed to write pre-commit hook: %v\n", err)
		os.Exit(1)
	}

	// Write post-commit hook
	postPath := filepath.Join(hooksDir, "post-commit")
	if err := os.WriteFile(postPath, []byte(postCommitHook), 0755); err != nil {
		fmt.Printf("✗ Failed to write post-commit hook: %v\n", err)
		os.Exit(1)
	}

	// Create session directory
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		fmt.Printf("✗ Failed to create session directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ AI Session Diff hooks installed successfully!")
	fmt.Println("  Pre-commit: snapshots staged files")
	fmt.Println("  Post-commit: records session metadata")
}

func Uninstall() {
	hooks := []string{"pre-commit", "post-commit"}
	for _, hook := range hooks {
		path := filepath.Join(hooksDir, hook)
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
			fmt.Printf("✓ Removed %s hook\n", hook)
		}
	}
	fmt.Print("✓ AI Session Diff hooks uninstalled.\n")
}

func ShowHelp() {
	fmt.Print(`AI Session Diff - Track and visualize AI coding sessions

Usage:
  ai-session-diff <command> [flags]

Commands:
  install       Install git hooks for session tracking
  uninstall     Remove installed hooks
  list          List all recorded sessions
  show <id>     Show terminal diff for a session
  report <id>   Generate HTML diff report
  status        Check hook installation status
  help          Show this help message

Flags:
  -h, --help    Show help

Examples:
  ai-session-diff install
  ai-session-diff list
  ai-session-diff list --ai-only
  ai-session-diff show abc123
  ai-session-diff report abc123
`)
}