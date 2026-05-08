package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	hooksDir     = ".git/hooks"
	sessionDir   = ".ai-sessions"
	isWindows    = runtime.GOOS == "windows"
	preCommitHook = `#!/bin/sh
# AI Session Diff - Pre-commit hook
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
SESSION_DIR="$REPO_DIR/.ai-sessions"
[ ! -d "$SESSION_DIR" ] && mkdir -p "$SESSION_DIR"
SESSION_ID=$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid 2>/dev/null || echo "pre-$(date +%s)")
BEFORE_DIR="$SESSION_DIR/$SESSION_ID/before"
mkdir -p "$BEFORE_DIR"
git diff --cached --name-only | while IFS= read -r file; do
    [ -f "$REPO_DIR/$file" ] && cp "$REPO_DIR/$file" "$BEFORE_DIR/$file"
done
echo "$SESSION_ID" > "$SESSION_DIR/.current-session"
exit 0
`
	postCommitHook = `#!/bin/sh
# AI Session Diff - Post-commit hook
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"
SESSION_DIR="$REPO_DIR/.ai-sessions"
[ ! -f "$SESSION_DIR/.current-session" ] && exit 0
SESSION_ID=$(cat "$SESSION_DIR/.current-session")
rm -f "$SESSION_DIR/.current-session"
AFTER_DIR="$SESSION_DIR/$SESSION_ID/after"
mkdir -p "$AFTER_DIR"
git ls-files | while IFS= read -r file; do
    [ -f "$REPO_DIR/$file" ] && cp "$REPO_DIR/$file" "$AFTER_DIR/$file"
done
COMMIT_SHA=$(git rev-parse HEAD 2>/dev/null)
BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
COMMIT_MSG=$(git log -1 --format=%s 2>/dev/null || echo "")
COMMIT_TIME=$(git log -1 --format=%cI 2>/dev/null || echo "")
LINES_ADDED=$(git diff --cached --numstat | awk '{add += $1} END {print add+0}')
LINES_REMOVED=$(git diff --cached --numstat | awk '{sub += $2} END {print sub+0}')
FILES_CHANGED=$(git diff --cached --name-only | tr '\n' ' ')
AI_PROVIDER=$([ -d "$REPO_DIR/.cursor" ] && echo "cursor" || echo "unknown")
AI_DETECTED=$([ -d "$REPO_DIR/.cursor" ] && echo "true" || echo "false")
cat > "$SESSION_DIR/$SESSION_ID/session.json" << EOF
{"id":"$SESSION_ID","startedAt":"$COMMIT_TIME","endedAt":"$COMMIT_TIME","commitSha":"$COMMIT_SHA","branch":"$BRANCH","filesChanged":["${FILES_CHANGED% }"],"aiProvider":"$AI_PROVIDER","aiDetected":$AI_DETECTED,"linesAdded":$LINES_ADDED,"linesRemoved":$LINES_REMOVED,"message":"$COMMIT_MSG"}
EOF
exit 0
`
	preCommitHookWin = `@echo off
REM AI Session Diff - Pre-commit hook (Windows)
setlocal enabledelayedexpansion
set "SCRIPT_DIR=%~dp0"
set "REPO_DIR=%SCRIPT_DIR%.."
set "SESSION_DIR=%REPO_DIR%\.ai-sessions"
if not exist "%SESSION_DIR%" mkdir "%SESSION_DIR%"
for /f "delims=" %%i in ('powershell -NoProfile -Command "[guid]::NewGuid().ToString()"') do set SESSION_ID=%%i
set "BEFORE_DIR=%SESSION_DIR%\%SESSION_ID%\before"
if not exist "%BEFORE_DIR%" mkdir "%BEFORE_DIR%"
for /f "delims=" %%f in ('git diff --cached --name-only') do (
    set "FILE=%%f"
    set "SOURCE=%REPO_DIR%\!FILE!"
    if exist "!SOURCE!" (
        set "DEST=%BEFORE_DIR%\!FILE!"
        for %%d in ("!DEST!") do if not exist "%%~dpd" mkdir "%%~dpd"
        copy /Y "!SOURCE!" "!DEST!" >nul
    )
)
echo %SESSION_ID% > "%SESSION_DIR%\.current-session"
exit /b 0
`
	postCommitHookWin = `@echo off
REM AI Session Diff - Post-commit hook (Windows)
setlocal enabledelayedexpansion
set "SCRIPT_DIR=%~dp0"
set "REPO_DIR=%SCRIPT_DIR%.."
set "SESSION_DIR=%REPO_DIR%\.ai-sessions"
if not exist "%SESSION_DIR%\.current-session" exit /b 0
set /p SESSION_ID=<"%SESSION_DIR%\.current-session"
del "%SESSION_DIR%\.current-session" 2>nul
set "AFTER_DIR=%SESSION_DIR%\%SESSION_ID%\after"
if not exist "%AFTER_DIR%" mkdir "%AFTER_DIR%"
for /f "delims=" %%f in ('git ls-files') do (
    set "FILE=%%f"
    set "SOURCE=%REPO_DIR%\!FILE!"
    if exist "!SOURCE!" (
        set "DEST=%AFTER_DIR%\!FILE!"
        for %%d in ("!DEST!") do if not exist "%%~dpd" mkdir "%%~dpd"
        copy /Y "!SOURCE!" "!DEST!" >nul
    )
)
for /f "delims=" %%s in ('git rev-parse HEAD 2^>nul') do set COMMIT_SHA=%%s
for /f "delims=" %%b in ('git rev-parse --abbrev-ref HEAD 2^>nul') do set BRANCH=%%b
for /f "delims=" %%m in ('git log -1 --format=%%s 2^>nul') do set COMMIT_MSG=%%m
for /f "delims=" %%t in ('git log -1 --format=%%cI 2^>nul') do set COMMIT_TIME=%%t
git diff --cached --numstat > "%TEMP%\ai-diff-stats.txt"
set LINES_ADDED=0
set LINES_REMOVED=0
for /f "usebackq tokens=1,2" %%a in ("%TEMP%\ai-diff-stats.txt") do (
    set /a LINES_ADDED+=%%a
    set /a LINES_REMOVED+=%%b
)
del "%TEMP%\ai-diff-stats.txt" 2>nul
set "AI_PROVIDER=unknown"
set "AI_DETECTED=false"
if exist "%REPO_DIR%\.cursor" (
    set "AI_PROVIDER=cursor"
    set "AI_DETECTED=true"
)
set "FILES_JSON=["
for /f "delims=" %%f in ('git diff --cached --name-only') do set "FILES_JSON=!FILES_JSON!\"%%f\","
if defined FILES_JSON set FILES_JSON=[!FILES_JSON:~0,-1!]
powershell -NoProfile -Command "$json = @{id='%SESSION_ID%';startedAt='%COMMIT_TIME%';endedAt='%COMMIT_TIME%';commitSha='%COMMIT_SHA%';branch='%BRANCH%';filesChanged=%FILES_JSON%;aiProvider='%AI_PROVIDER%';aiDetected=%AI_DETECTED%;linesAdded=%LINES_ADDED%;linesRemoved=%LINES_REMOVED%;message='%COMMIT_MSG%'} | ConvertTo-Json -Compress; Set-Content -Path '%SESSION_DIR%\%SESSION_ID%\session.json' -Value $json"
exit /b 0
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
	hookContent := preCommitHook
	if isWindows {
		hookContent = preCommitHookWin
	}
	if err := os.WriteFile(prePath, []byte(hookContent), 0755); err != nil {
		fmt.Printf("✗ Failed to write pre-commit hook: %v\n", err)
		os.Exit(1)
	}

	// Write post-commit hook
	postPath := filepath.Join(hooksDir, "post-commit")
	hookContent = postCommitHook
	if isWindows {
		hookContent = postCommitHookWin
	}
	if err := os.WriteFile(postPath, []byte(hookContent), 0755); err != nil {
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
  export         Export all sessions as JSON to stdout
  export-csv     Export all sessions as CSV to stdout
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