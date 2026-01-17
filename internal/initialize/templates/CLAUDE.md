# Claude Code Guide for todo-rank

This document helps Claude Code effectively use the `todo-rank` CLI tool.

## What is todo-rank?

A CLI tool that scans codebases for TODO/FIXME/@next comments and ranks them by priority score, helping you decide what to fix first.

## Installation Check

Before using, verify installation:

```bash
which todo-rank || echo "Not installed"
```

If not installed:

```bash
go install github.com/simulacre7/todo-rank/cmd/todo-rank@latest
export PATH="$HOME/go/bin:$PATH"
```

## Recommended Workflows

### 1. Start of Session - Identify Priorities

When starting work on a codebase, run:

```bash
todo-rank --min-score 80
```

This shows P0 (Now) and P1 (Soon) items that need immediate attention.

### 2. Before Implementing Features

Check for related TODOs in the area you're working on:

```bash
todo-rank --root ./path/to/module --min-score 40
```

### 3. Generate Task List for User

When user asks "what should I work on?" or "show me TODOs":

```bash
todo-rank --format md
```

### 4. Focus on Critical Code Paths

For production-critical code:

```bash
todo-rank --root ./cmd --min-score 120
```

### 5. Technical Debt Review

For comprehensive review:

```bash
todo-rank --min-score 0
```

## Command Reference

```bash
todo-rank [options]

Options:
  --root <path>      Scan directory (default: .)
  --ignore <csv>     Skip directories (default: .git,node_modules,dist)
  --format <text|md> Output format (default: text)
  --out <path>       Save to file (default: stdout)
  --min-score <n>    Filter by minimum score (default: 0)
  --tags <csv>       Tags to find (default: TODO,FIXME,@next)
```

## Understanding Scores

| Score Range | Level | Meaning |
|-------------|-------|---------|
| >= 120 | P0 (Now) | Critical, fix immediately |
| >= 80 | P1 (Soon) | Important, fix soon |
| >= 40 | P2 (Later) | Can wait |
| < 40 | P3 (Cleanup) | Low priority cleanup |

### Score Calculation

- **Tags**: FIXME (+100), TODO (+50), @next (+30)
- **Priority markers**: [P0] (+100), [P1] (+70), [P2] (+40), [P3] (+10)
- **Path bonus**: `cmd/` or `main.go` (+20)
- **Test penalty**: `*_test.*` files (-20)

## Best Practices

1. **Always check high-priority TODOs first** before starting new work
2. **Use `--format md`** when presenting results to users for better readability
3. **Filter by path** (`--root`) when working on specific modules
4. **Respect developer intent**: [P0] markers mean the original author considered it critical
5. **Don't ignore test TODOs entirely**: they may indicate missing test coverage

## Example Output Interpretation

```
P0 (Now)
(200) cmd/server/main.go:42
  FIXME[P0]: data race when shutting down
```

This means:
- Score 200 = FIXME(100) + [P0](100)
- Location: cmd/server/main.go line 42
- Developer marked as critical ([P0]) data race issue

## Integration Tips

When writing code that includes TODOs:

```go
// TODO: description                    // Basic, score 50
// TODO[P1]: description                // With priority, score 120
// FIXME[P0]: critical issue            // Urgent fix, score 200
// @next: upcoming refactor             // Next sprint item, score 30
```

Use priority markers to communicate urgency to future developers and tools.
