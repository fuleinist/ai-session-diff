# AI Session Diff

**Git hooks that record AI coding sessions and generate visual diffs.**

Track what AI coding tools (Cursor, Claude Code, Copilot) changed — and why.

## Installation

```bash
# Clone the repo
git clone https://github.com/fuleinist/ai-session-diff.git
cd ai-session-diff

# Build
go build -o ai-session-diff .

# Install hooks in your project
./ai-session-diff install
```

## Usage

```bash
# List all recorded sessions
ai-session-diff list

# Show diff for a session (terminal)
ai-session-diff show <session-id>

# Generate HTML report
ai-session-diff report <session-id>

# Check hook status
ai-session-diff status

# Uninstall hooks
ai-session-diff uninstall
```

## How It Works

1. **Pre-commit hook** snapshots staged files
2. **Post-commit hook** records session metadata and creates the diff
3. Session data stored in `.ai-sessions/<uuid>/`

## AI Detection

Automatically detects:
- **Cursor** — via `.cursor/` directory
- **Claude Code** — via `CLAUDE_CODE` or `ANTHROPIC_MODEL` env vars
- **GitHub Copilot** — via `.github/copilot-innb` file

## Example Output

```
⚡ AI Session Diff — Session History
─────────────────────────────────────────────────
  abc12345  2026-05-07  +42 -7  a1b2c3d [cursor]
  def67890  2026-05-06  +15 -3  e4f5g6h
```

## License

MIT