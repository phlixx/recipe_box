# Architecture Decisions

This document records significant technical decisions made during development.

---

## ADR-001: CLI Framework

**Date**: 2025-01-21
**Status**: Accepted

### Context
Need a CLI framework for building the recipe management tool. Options considered:
- Cobra (most popular Go CLI framework)
- urfave/cli
- Standard library flag package

### Decision
Use **Cobra** for CLI framework.

### Consequences
- **Pros**:
  - Rich feature set (subcommands, flags, completion)
  - Well documented
  - Large community
  - Used by kubectl, hugo, gh
- **Cons**:
  - Slightly heavier than alternatives
  - Learning curve for advanced features

---

## ADR-002: Project Structure

**Date**: 2025-01-21
**Status**: Accepted

### Context
Need to establish a maintainable project structure.

### Decision
Follow standard Go project layout:
- `cmd/` for CLI commands
- `internal/` for private application code
- `pkg/` for reusable packages (if needed)

### Consequences
- **Pros**:
  - Familiar to Go developers
  - Clear separation of concerns
  - Easy to navigate
- **Cons**:
  - May be overkill for small project initially

---

## ADR-003: Storage Format

**Date**: 2025-01-21
**Status**: Proposed

### Context
Need to persist recipes locally. Options:
- JSON files
- SQLite
- YAML files

### Decision
Start with **JSON files** for simplicity.

### Consequences
- **Pros**:
  - No dependencies
  - Human readable
  - Easy to debug
  - Simple backup/sync
- **Cons**:
  - No query capabilities
  - Entire file must be loaded
  - Concurrent access issues (acceptable for CLI)

---

## ADR-004: Storage Location

**Date**: 2025-01-21
**Status**: Accepted

### Context
Need to decide where to store recipes, config, and other data.

### Decision
Use **`~/.recipe_box/`** as the storage directory.

### Consequences
- **Pros**: Standard location, works across projects, easy to backup
- **Cons**: Hidden directory, need to handle ~ expansion

---

## ADR-005: API Key Management

**Date**: 2025-01-21
**Status**: Accepted

### Context
Need to store Claude API key for recipe generation.

### Decision
Store API key in **config file** (`~/.recipe_box/config.json`).

### Consequences
- **Pros**: Persistent, user-friendly, can store other settings
- **Cons**: Security consideration (file permissions), need config commands

---

## ADR-006: Units

**Date**: 2025-01-21
**Status**: Accepted

### Context
Recipes need measurements. Support metric, imperial, or both?

### Decision
**Metric only** for MVP. Can add imperial conversion later.

### Consequences
- **Pros**: Simpler implementation, consistent output
- **Cons**: Less convenient for US users

---

## Template for New Decisions

```markdown
## ADR-XXX: Title

**Date**: YYYY-MM-DD
**Status**: Proposed | Accepted | Deprecated | Superseded

### Context
What is the issue that we're seeing that is motivating this decision?

### Decision
What is the change that we're proposing and/or doing?

### Consequences
What becomes easier or more difficult to do because of this change?
```
