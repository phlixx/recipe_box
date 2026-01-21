# Iteration Plan

## Overview

MVP (v0.1) broken into **5 iterations**. Each iteration delivers working functionality.

```
Iter 1: Foundation     → Config & storage infrastructure
Iter 2: Recipe CRUD    → Manual recipe management
Iter 3: AI Generation  → Claude integration
Iter 4: Shopping List  → Generate from recipe
Iter 5: Localization   → DE/EN support
```

---

## Iteration 1: Foundation

**Goal**: Infrastructure for config and data storage

**Deliverables**:
- [ ] Config system (`~/.recipe_box/config.json`)
- [ ] `config set <key> <value>` command
- [ ] `config get <key>` command
- [ ] Storage layer for JSON persistence
- [ ] Recipe data model (`internal/recipe/`)

**Exit Criteria**:
```bash
recipe_box config set api_key "sk-xxx"
recipe_box config get api_key  # → sk-xxx
```

**Files to create**:
- `internal/config/config.go`
- `internal/storage/storage.go`
- `internal/recipe/recipe.go` (model only)
- `cmd/config.go`

---

## Iteration 2: Recipe CRUD (Manual)

**Goal**: Manage recipes without AI

**Deliverables**:
- [ ] `recipe add` - Interactive prompts to add recipe
- [ ] `recipe list` - Show all saved recipes
- [ ] `recipe view <id>` - Show recipe details
- [ ] `recipe delete <id>` - Remove recipe

**Exit Criteria**:
```bash
recipe_box recipe add          # prompts for title, ingredients, steps
recipe_box recipe list         # shows saved recipes
recipe_box recipe view abc123  # shows full recipe
recipe_box recipe delete abc123
```

**Files to create**:
- `cmd/recipe.go` (command group)
- `cmd/recipe_add.go`
- `cmd/recipe_list.go`
- `cmd/recipe_view.go`
- `cmd/recipe_delete.go`
- `internal/recipe/service.go`

---

## Iteration 3: AI Generation

**Goal**: Generate recipes with Claude

**Deliverables**:
- [ ] Claude API client (`internal/ai/`)
- [ ] `recipe generate` - Generate recipe from prompt/flags
- [ ] `recipe save` - Save last generated recipe to collection
- [ ] Structured output parsing (title, ingredients, steps)

**Exit Criteria**:
```bash
recipe_box recipe generate --prompt "quick pasta dish"
# → Shows generated recipe

recipe_box recipe generate --ingredients "chicken, lemon"
# → Shows recipe using those ingredients

recipe_box recipe save
# → Saves last generated to collection
```

**Files to create**:
- `internal/ai/claude.go`
- `cmd/recipe_generate.go`
- `cmd/recipe_save.go`

**Dependencies**:
- Requires `ANTHROPIC_API_KEY` in config (from Iter 1)

---

## Iteration 4: Shopping List

**Goal**: Generate shopping list from recipe

**Deliverables**:
- [ ] `shopping generate <recipe-id>` - Create list from recipe
- [ ] `shopping show` - Display current shopping list
- [ ] `shopping clear` - Clear the list
- [ ] Shopping list storage

**Exit Criteria**:
```bash
recipe_box shopping generate abc123
# → Creates shopping list from recipe ingredients

recipe_box shopping show
# → Shows grouped shopping list

recipe_box shopping clear
```

**Files to create**:
- `internal/shopping/shopping.go`
- `cmd/shopping.go`
- `cmd/shopping_generate.go`
- `cmd/shopping_show.go`
- `cmd/shopping_clear.go`

---

## Iteration 5: Localization

**Goal**: Support German and English

**Deliverables**:
- [ ] `config set language de|en`
- [ ] CLI messages in selected language
- [ ] AI generates recipes in selected language
- [ ] Shopping list output localized

**Exit Criteria**:
```bash
recipe_box config set language de
recipe_box recipe generate --prompt "schnelles Nudelgericht"
# → Recipe in German

recipe_box shopping show
# → Categories in German (Gemüse, Milchprodukte, etc.)
```

**Files to create**:
- `internal/i18n/i18n.go`
- `internal/i18n/messages_en.go`
- `internal/i18n/messages_de.go`

---

## Summary

| Iter | Focus | Key Commands |
|------|-------|--------------|
| 1 | Foundation | `config set/get` |
| 2 | Recipe CRUD | `recipe add/list/view/delete` |
| 3 | AI Generation | `recipe generate/save` |
| 4 | Shopping | `shopping generate/show/clear` |
| 5 | Localization | Language support |

**Total: 5 iterations for MVP**

---

## Notes

- Each iteration should be a working increment (can use CLI after each)
- Iterations 1-2 work fully offline (no API needed)
- Iteration 3 requires Claude API key
- Order matters: 1 → 2 → 3 → 4 → 5 (dependencies)
