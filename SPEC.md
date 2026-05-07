# AI Session Diff — Specification

## 1. Concept & Vision

**AI Session Diff** is a git hook toolkit that records AI coding agent sessions (Cursor, Claude Code, Copilot, etc.) and generates beautiful visual diffs showing *what* changed and *when*. It answers the question: "What did the AI do in this session, and why?"

Think of it as a flight recorder for AI coding — not just a diff tool, but a session tracker with timeline visualization.

Target users: dev teams reviewing AI-generated code, solo devs auditing their own sessions, tech leads tracking AI velocity.

## 2. Design Language

- **Aesthetic**: Developer-focused, terminal-first, clean. Inspired by GitHub's diff UI.
- **Colors**: Terminal ANSI colors — green (+), red (-), cyan (info), yellow (warning).
- **Typography**: Monospace throughout (JetBrains Mono / Fira Code / system default).
- **Motion**: None in CLI; smooth scroll in HTML output.
- **Icons**: Unicode symbols — ✓ ✗ ⚡ 📁 ➕ ➖

## 3. Tech Stack

- **Language**: Go 1.21+
- **Hooks**: Shell scripts calling Go binary (git hooks)
- **Diff visualization**: Custom Go renderer (terminal + HTML)
- **Storage**: Local JSON session files in `.ai-sessions/`
- **Dependencies**: Standard library + `github.com/go-git/go-git/v5` for git ops

## 4. Features & MVP Scope

### Phase 1 — Core MVP (this build)

#### 4.1 Session Recording (git hooks)
- `pre-commit` hook: snapshot all staged files into `.ai-sessions/<uuid>/before/`
- `post-commit` hook: snapshot post-commit state into `.ai-sessions/<uuid>/after/`, record metadata JSON
- Hook activation: `ai-session-diff install` command installs hooks into `.git/hooks/`
- Session ID: UUID generated per commit chain (grouped by branch session)

#### 4.2 Session Metadata
Each session writes `session.json`:
```json
{
  "id": "uuid",
  "startedAt": "ISO8601",
  "endedAt": "ISO8601",
  "commitSha": "abc123",
  "branch": "main",
  "filesChanged": ["foo.go", "bar.ts"],
  "aiProvider": "cursor|claude|github-copilot|unknown",
  "aiDetected": true,
  "linesAdded": 42,
  "linesRemoved": 7
}
```

#### 4.3 AI Provider Detection
Detect AI tool from:
- `.cursor/` directory presence → "cursor"
- Environment var `CLAUDE_CODE=true` or `ANTHROPIC_MODEL` → "claude-code"
- `.github/copilot-innb` file → "github-copilot"
- Else "unknown"

#### 4.4 Terminal Diff Viewer
- `ai-session-diff show <session-id>`: prints unified diff in terminal
- Color-coded: additions (green), deletions (red), context (default)
- File headers with `+`/`-` counts
- `--stats` flag: just summary, no diff content

#### 4.5 HTML Diff Report
- `ai-session-diff report <session-id>`: generates `session.html`
- Side-by-side diff view with syntax highlighting
- File tree sidebar
- Session metadata header
- Dark/light theme toggle

#### 4.6 Session List
- `ai-session-diff list`: shows all recorded sessions with ID, date, file count, lines changed
- `ai-session-diff list --ai-only`: filter to sessions with AI detected

### CLI Commands

| Command | Description |
|---------|-------------|
| `ai-session-diff install` | Install git hooks in current repo |
| `ai-session-diff uninstall` | Remove hooks |
| `ai-session-diff show <id>` | Show terminal diff for session |
| `ai-session-diff report <id>` | Generate HTML report |
| `ai-session-diff list` | List all sessions |
| `ai-session-diff status` | Show hook installation status |

## 5. File Structure

```
ai-session-diff/
├── cmd/
│   └── main.go              # CLI entrypoint
├── internal/
│   ├── hook/                # Git hook scripts
│   │   ├── pre-commit.go    # Pre-commit hook logic
│   │   └── post-commit.go   # Post-commit hook logic
│   ├── diff/                # Diff generation
│   │   ├── diff.go          # Core diff algorithm
│   │   └── renderer.go     # Terminal + HTML renderers
│   ├── session/             # Session storage
│   │   └── session.go       # Session model + storage
│   ├── detector/            # AI provider detection
│   │   └── detector.go
│   └── cli/                 # CLI commands
│       ├── install.go
│       ├── show.go
│       ├── report.go
│       └── list.go
├── .ai-sessions/            # (created at runtime) Session data
├── go.mod
├── go.sum
├── SPEC.md
└── README.md
```

## 6. Acceptance Criteria

- [ ] `ai-session-diff install` installs hooks and reports success
- [ ] Hooks correctly snapshot before/after state on commit
- [ ] `ai-session-diff list` shows all sessions with metadata
- [ ] `ai-session-diff show <id>` renders colored terminal diff
- [ ] `ai-session-diff report <id>` generates valid HTML with diff
- [ ] AI provider detection correctly identifies cursor/claude-code
- [ ] Project builds: `go build -o ai-session-diff`
- [ ] Project tests: `go test ./...` (basic unit tests)
- [ ] README with install + usage instructions
