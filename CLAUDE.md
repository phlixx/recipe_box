# Recipe Box CLI

A Go CLI for generating recipes and shopping lists.

## Project Overview

**Goal**: Build a CLI tool for home cooks and meal preppers to:
- Generate recipes using Claude AI
- Plan meals for a week (with leftover tracking)
- Generate shopping lists from meal plans
- Manage a personal recipe collection
- Support German and English localization

**Tech Stack**: Go, Cobra (CLI framework), Claude API

**Key Docs**: See `PRODUCT.md` for use cases, `ITERATIONS.md` for roadmap.

---

## Session Continuity

### State Files (Lean Model)

| File | Purpose | Updates |
|------|---------|---------|
| `PRODUCT.md` | Use cases, data model, scope | When scope/requirements change |
| `ITERATIONS.md` | Iteration plan with deliverables | When roadmap changes |
| `STATUS.md` | Current focus, blockers, handoff notes | Every session |
| `TODO.md` | Pending tasks only (no completed section) | When tasks change |
| `DECISIONS.md` | Architecture decisions (append-only) | When decisions made |
| `git log` | Completed work (source of truth) | Automatic |

### Session Types

| Type | Focus | Key Updates |
|------|-------|-------------|
| **Implement** | Write code, fix bugs, add features | STATUS, TODO, README |
| **Steer** | Plan phases, design decisions, adjust scope | STATUS, PRODUCT, DECISIONS, ITERATIONS |

Most sessions are **Implement**. Use **Steer** when:
- Planning a new phase or major feature
- Making architecture/design decisions
- Changing project scope or requirements
- Reviewing or adjusting the roadmap

### Starting a Session

```
1. Run: git log --oneline -5
2. Read: STATUS.md
3. Read: TODO.md
4. Validate: Do the files align with git history? Flag if not.
5. Identify: Is this an Implement or Steer session?
```

**Validation check**: If STATUS.md mentions work that doesn't appear in recent commits, or TODO.md lists tasks as pending that were already done, fix the inconsistency before proceeding.

### Ending a Session

**Both session types** — Update STATUS.md:
- Current focus (what's being worked on)
- Blockers (if any)
- Handoff notes (context for next session)

**Implement sessions** — Also update:
- TODO.md (remove completed, add new)
- README.md (if user-facing features changed)

**Steer sessions** — Also update:
- PRODUCT.md (if scope or requirements changed)
- DECISIONS.md (append ADR if design decisions made)
- ITERATIONS.md (if roadmap changed)
- TODO.md (if planning surfaced concrete tasks)

**Do NOT duplicate git history in state files.**

---

## Development Loop

```
UNDERSTAND → PLAN → IMPLEMENT → VERIFY → DOCUMENT
```

### 1. Understand
- Read relevant existing code
- Check if similar patterns exist in codebase
- Clarify requirements if ambiguous

### 2. Plan
- Break into small, testable increments
- Identify files to create/modify
- Consider edge cases

### 3. Implement
- Write code in small chunks
- Follow existing patterns in codebase
- Keep changes focused and minimal

### 4. Verify

**Run after every change:**
```bash
make check
```

This runs: `fmt` → `vet` → `test`

For commits, also run:
```bash
make build && make smoke
```

### 5. Document
- Update STATUS.md (focus/blockers)
- Update TODO.md (remove done, add new)
- Add code comments only where non-obvious

---

## Code Conventions

### Project Structure

```
recipe_box/
├── main.go              # Entry point
├── cmd/                 # CLI commands (Cobra)
│   └── root.go          # Root command
├── internal/
│   ├── ai/              # Claude API client (Iter 3)
│   ├── config/          # Config management (Iter 1)
│   ├── i18n/            # Localization (Iter 5)
│   ├── recipe/          # Recipe domain logic
│   ├── shopping/        # Shopping list logic
│   └── storage/         # Data persistence
├── CLAUDE.md            # This file
├── DECISIONS.md         # Architecture decisions
├── ITERATIONS.md        # MVP roadmap
├── PRODUCT.md           # Use cases & data model
├── STATUS.md            # Current focus & handoff
└── TODO.md              # Pending tasks
```

### Go Conventions

- Use `internal/` for private packages
- Error handling: Return errors, don't panic
- Naming: `camelCase` for private, `PascalCase` for public
- Keep functions small and focused
- Use table-driven tests

---

## Verification Loop

**ALWAYS run verification after making changes.** This is non-negotiable.

```
┌─────────────────────────────────────────────────────┐
│                   IMPLEMENT                          │
│         Write code in small increments              │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                    VERIFY                            │
│              make check                              │
│         (fmt → vet → test)                          │
└──────────────────────┬──────────────────────────────┘
                       │
              ┌────────┴────────┐
              │                 │
              ▼                 ▼
         ✓ PASS            ✗ FAIL
              │                 │
              ▼                 ▼
         CONTINUE             FIX
      (next increment     (debug, resolve,
        or commit)         then re-verify)
```

### Checklist

Before considering any task complete:

**Implement sessions:**
- [ ] `make check` passes (fmt, vet, test)
- [ ] `make build` succeeds
- [ ] Smoke test if CLI behavior changed
- [ ] STATUS.md updated
- [ ] Completed tasks removed from TODO.md
- [ ] README.md updated if user-facing features changed

**Steer sessions:**
- [ ] PRODUCT.md updated if scope changed
- [ ] DECISIONS.md has ADR for each design decision
- [ ] ITERATIONS.md reflects current roadmap
- [ ] STATUS.md updated with handoff notes

> **If `make check` doesn't pass, the work isn't done.**

---

## Common Patterns

### Adding a New Command

1. Create `cmd/<name>.go`
2. Define command with Cobra
3. Register in `cmd/root.go`
4. Add business logic in `internal/`
5. Write tests

### Adding a New Feature

1. Define interface in `internal/<domain>/`
2. Implement concrete type
3. Write unit tests
4. Wire up to CLI command
5. Integration test

---

## Quick Reference

| Action | Command |
|--------|---------|
| **Full check** | `make check` |
| Build binary | `make build` |
| Run tests | `make test` |
| Run with args | `make run ARGS="config get api_key"` |
| Test coverage | `make test-coverage` |
| Smoke test | `make smoke` |
| Git history | `git log --oneline -5` |

### Running Specific Tests

```bash
# All tests
make test

# Specific package
go test -v ./internal/config/...

# Specific test
go test -v -run TestConfig_SetAndGet ./internal/config/
```
