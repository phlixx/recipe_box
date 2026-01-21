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

## Summary - MVP (v0.1) ✅

| Iter | Focus | Key Commands | Status |
|------|-------|--------------|--------|
| 1 | Foundation | `config set/get` | ✅ Done |
| 2 | Recipe CRUD | `recipe add/list/view/delete` | ✅ Done |
| 3 | AI Generation | `recipe generate/save` | ✅ Done |
| 4 | Shopping | `shopping generate/show/clear` | ✅ Done |
| 5 | Localization | Language support | ✅ Done |

---

# Post-MVP: Phase 1 - UX & Polish (v0.2)

```
Iter 6: Interactive Mode  → REPL with autocomplete
Iter 7: Servings Scaling  → Scale recipes up/down
Iter 8: Localization+     → Complete all strings, units
```

---

## Iteration 6: Interactive Mode ✅

**Goal**: REPL as default mode with colored output

**Deliverables**:
- [x] Interactive mode when running `recipe_box` (no args)
- [x] Emoji prompt: `🍳 > `
- [x] Tab completion for commands
- [x] Tab completion for recipe IDs
- [x] Command history (arrow keys)
- [x] `exit` command to quit
- [x] Colored output (headers, categories, success/error)
- [x] Traditional CLI still works (`recipe_box recipe list`)

**Exit Criteria**:
```bash
$ recipe_box
🍳 > recipe list
Found 2 recipe(s):          # colored header
  816758d9  Chicken Teriyaki
  ...

🍳 > recipe view 81<TAB>    # autocompletes to 816758d9
🍳 > exit
$
```

**Files to create/modify**:
- `cmd/interactive.go` - REPL logic
- `internal/ui/colors.go` - Color/style definitions
- `main.go` - Entry point routing

**Dependencies**:
- `github.com/c-bata/go-prompt`
- `github.com/fatih/color` or `github.com/charmbracelet/lipgloss`

---

## Iteration 7: Servings Scaling

**Goal**: Scale recipe ingredients for different serving sizes

**Deliverables**:
- [ ] `recipe view <id> --servings N` - View with scaled ingredients
- [ ] `recipe generate --servings N` - Generate for specific servings
- [ ] `shopping generate <id> --servings N` - Shopping list with scaled quantities
- [ ] Display shows both original and scaled servings

**Exit Criteria**:
```bash
🍳 > recipe view 816758d9 --servings 4
Chicken Teriyaki (scaled: 4 servings, original: 10)

Ingredients:
  - 1 kg chicken thighs (was 2.5 kg)
  - 80 ml soy sauce (was 200 ml)
  ...
```

**Files to modify**:
- `cmd/recipe_view.go` - Add --servings flag
- `cmd/recipe_generate.go` - Add --servings flag
- `cmd/shopping_generate.go` - Add --servings flag
- `internal/recipe/recipe.go` - Add Scale() method

---

## Iteration 8: Localization Hardening

**Goal**: Complete and consistent localization

**Deliverables**:
- [ ] Audit all `fmt.Print*` calls for hardcoded strings
- [ ] Localize Cobra command descriptions (Short, Long, Use)
- [ ] Localize unit names (tbsp→EL, cups→Tassen, tsp→TL)
- [ ] Localize error messages (user-actionable ones)
- [ ] Add unit tests verifying all keys exist in both languages
- [ ] Interactive mode messages localized

**Exit Criteria**:
```bash
$ recipe_box config set language de
$ recipe_box
🍳 > hilfe                    # or still 'help'?
Verfügbare Befehle:
  rezept     Rezepte verwalten
  einkauf    Einkaufsliste verwalten
  config     Einstellungen
  exit       Beenden

🍳 > rezept liste
2 Rezept(e) gefunden:
  ...
```

**Decision needed**: Localize command names too? (`recipe`→`rezept`, `shopping`→`einkauf`)

**Files to modify**:
- `internal/i18n/messages_*.go` - Add all missing keys
- `cmd/*.go` - Use i18n for Short/Long descriptions
- `internal/i18n/units.go` - Unit name translations

---

## Summary - Post-MVP Phase 1

| Iter | Focus | Key Features | Status |
|------|-------|--------------|--------|
| 6 | Interactive Mode | REPL, autocomplete, colors | ✅ Done |
| 7 | Servings | Scale ingredients up/down | |
| 8 | Localization+ | Complete translations, units | |

**Total: 3 iterations for v0.2**

---

## Future Phases (Ideas)

**Phase 2 - Meal Planning**:
- Weekly meal planner
- Leftover tracking
- Shopping from plan

**Phase 3 - Advanced Features**:
- Recipe import from URL
- Pantry tracking
- Export to PDF/Markdown
- Recipe sharing

---

## Notes

- Each iteration should be a working increment (can use CLI after each)
- Iterations 1-2 work fully offline (no API needed)
- Iteration 3+ requires Claude API key
- Order matters: 1 → 2 → 3 → 4 → 5 (dependencies)
- Post-MVP iterations are more independent
