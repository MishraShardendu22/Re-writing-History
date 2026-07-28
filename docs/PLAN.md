# Phase 1 Plan: Re-writing-History

## 1. Design Intent
Establish a dark-mode Git History & Commit Graph Telemetry UI for `Re-writing-History`. The UI will re-interpret the canonical design system for Git engineering tools: dark base (`#09090b`), raised card panels (`#121318`), indigo commit SHA pills, interactive commit log timeline, AI rewrite diff viewer, and responsive layout.

## 2. Primitives & Token System
* **Colors**:
  * Surface layers: `--bg-base` (`#09090b`), `--bg-raised` (`#121318`), `--bg-overlay` (`#1c1d24`).
  * Text: `--text-primary` (`#f4f4f5`), `--text-secondary` (`#a1a1aa`), `--text-muted` (`#71717a`).
  * Borders: `--border-default` (`#27272a`), `--border-focus` (`#6366f1`).
  * Status: Original (`#3b82f6`), Rewritten (`#10b981`), AI Generated (`#8b5cf6`).
* **Typography**: Inter font for UI + Monospace font for Git SHAs, timestamps, and commit diff previews.

## 3. Concrete Component Work
* **Commit Graph Visualizer**: Vertical commit log timeline displaying commit hash, author, date range filter, and AI rewrite status.
* **Diff & Rewrite Inspector**: Side-by-side comparison panel for original vs. AI-enhanced commit messages.
* **HTTP Web Server**: Add `util.StartWebDashboard(port)` helper in Go server.

## 4. Verification Plan
* `go build`: Successful Go compilation.
* `go test ./...`: Go tests pass cleanly.
