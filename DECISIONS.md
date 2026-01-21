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

## ADR-007: Interactive Mode as Default

**Date**: 2026-01-21
**Status**: Accepted

### Context
Currently, each command requires running `recipe_box <command>` from the shell. Users want a more fluid experience where they can stay in the application and issue multiple commands without retyping the binary name.

### Decision
Implement an **interactive REPL mode** as the default when running `recipe_box` with no arguments.

- Use `go-prompt` library for REPL functionality
- Emoji-based prompt: `🍳 > `
- Tab completion for commands and recipe IDs
- Command history via arrow keys
- Exit via `exit` command
- Traditional CLI mode (`recipe_box recipe list`) remains available for scripting

### Consequences
- **Pros**:
  - More fluid user experience
  - Autocomplete reduces typing errors
  - History allows quick command repetition
  - Feels like a proper application
- **Cons**:
  - Additional dependency (go-prompt)
  - Need to handle two entry points (interactive vs direct)
  - Slightly more complex main.go

---

## ADR-008: Terminal Styling

**Date**: 2026-01-21
**Status**: Accepted

### Context
Plain text output is functional but lacks visual hierarchy. Users want colored/styled output to improve readability.

### Decision
Add **rich but subtle terminal styling** using `fatih/color` or `lipgloss`.

Style guide:
- **Headers/Titles**: Bold, subtle color (cyan or blue)
- **Categories**: Muted colors to group visually
- **Success messages**: Green
- **Errors**: Red
- **Prompts/Labels**: Dim/gray
- **Values/Data**: Default color (white/terminal default)

Respect `NO_COLOR` environment variable for accessibility.

### Consequences
- **Pros**:
  - Better visual hierarchy
  - Easier to scan output
  - More professional feel
- **Cons**:
  - Additional dependency
  - Need to test in different terminals
  - Must handle non-color terminals gracefully

---

## ADR-009: Servings Scaling

**Date**: 2026-01-21
**Status**: Accepted

### Context
Recipes are generated/stored with a fixed serving count. Users want to scale recipes up or down without modifying the original.

### Decision
Add `--servings` flag to relevant commands:

- `recipe view <id> --servings N` - Display with scaled ingredients
- `recipe generate --servings N` - Generate for specific serving count
- `shopping generate <id> --servings N` - Generate list with scaled quantities

Scaling is linear: `new_quantity = original_quantity * (new_servings / original_servings)`

### Consequences
- **Pros**:
  - Flexible for different household sizes
  - No need to duplicate recipes
  - Shopping list accuracy for actual needs
- **Cons**:
  - Some ingredients don't scale linearly (e.g., salt, spices)
  - Fractional quantities may look odd (0.33 onions)

---

## ADR-010: Localization Strategy

**Date**: 2026-01-21
**Status**: Accepted

### Context
MVP localization covered main user-facing messages but some strings remain hardcoded. Need a systematic approach.

### Decision
Implement **comprehensive localization**:

1. **All user-visible strings** go through `i18n.T()`
2. **Cobra command descriptions** (Short, Long) localized
3. **Unit names** localized (tbsp→EL, cups→Tassen, etc.)
4. **Error messages** localized where user-actionable
5. **Audit process**: grep for `fmt.Print` and `fmt.Errorf` to find hardcoded strings

**Keep in English (static)**:
- Command names (`recipe`, `shopping`, `config`, `help`)
- Flag names (`--servings`, `--prompt`, `--quick`)
- Technical error details (for debugging)
- Log messages
- JSON field names

### Consequences
- **Pros**:
  - Consistent language experience
  - Professional feel for non-English users
  - Scalable to more languages later
- **Cons**:
  - More translation work
  - Longer message files
  - Need to maintain parity between languages

---

## ADR-011: Interactive Recipe Generation Prompts

**Date**: 2026-01-21
**Status**: Accepted

### Context
In REPL mode, running `recipe generate` with no flags requires the user to know all available options. A more guided experience would help users discover and use generation features.

### Decision
Add **interactive prompts** when `recipe generate` is run without flags in REPL mode:

1. Prompt for recipe description (main request)
2. Prompt for ingredients to use (optional)
3. Prompt for cuisine type (optional)
4. Prompt for vegetarian preference (y/n)
5. Prompt for quick recipe preference (y/n)
6. After generation, prompt to save to collection (y/n)

Use `prompt.Input()` from go-prompt library (not bufio.Reader) to avoid stdin conflicts with the REPL.

### Consequences
- **Pros**:
  - Guided experience for new users
  - Discoverable options without reading help
  - Seamless flow from generate to save
- **Cons**:
  - Slightly slower than flags for power users
  - Flags still available for scripting/quick use

---

## ADR-012: ASCII Art for Recipes

**Date**: 2026-01-21
**Status**: Accepted

### Context
Terminal output is text-heavy. A visual element would make the experience more engaging and help users quickly identify dishes.

### Decision
Add **ASCII art** to AI-generated recipes:

1. Add `ascii_art` field to Recipe struct (optional, `omitempty`)
2. Request ASCII art in AI prompt (8-12 lines, max 40 chars wide)
3. Display at top of recipe output (before title)
4. Art generated inline with recipe (no external APIs/images)

### Consequences
- **Pros**:
  - More engaging visual experience
  - Differentiates AI-generated recipes
  - Minimal token cost (~200-400 extra tokens per recipe)
  - Art cached with recipe (only generated once)
- **Cons**:
  - Slightly larger API responses
  - Quality varies by dish complexity
  - Manual recipes won't have art (acceptable)

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
