# Product Specification

## Target Users

- **Home cooks** - Want recipe inspiration and organized shopping
- **Meal preppers** - Plan ahead, cook in batches, minimize waste

---

## Core Use Cases

### UC1: Generate Recipe with AI

**As a** user
**I want to** generate a recipe based on criteria (cuisine, ingredients, dietary restrictions)
**So that** I get inspiration without browsing multiple websites

```bash
recipe_box recipe generate --cuisine italian --vegetarian
recipe_box recipe generate --ingredients "chicken, rice, broccoli"
recipe_box recipe generate --quick  # under 30 min
```

**Acceptance Criteria:**
- Connects to Claude API
- Returns structured recipe (title, ingredients, steps, servings, prep time)
- Can save generated recipe to collection

---

### UC2: Manage Recipe Collection

**As a** user
**I want to** save, list, view, edit, and delete recipes
**So that** I have a personal recipe library

```bash
recipe_box recipe list
recipe_box recipe view <id>
recipe_box recipe add                    # manual entry
recipe_box recipe save <generated-id>    # save AI-generated
recipe_box recipe edit <id>
recipe_box recipe delete <id>
```

**Acceptance Criteria:**
- Recipes stored locally (JSON)
- Can edit ingredients (for substitutions/preferences)
- Search/filter by name, tags, cuisine

---

### UC3: Plan Meals for a Timeframe

**As a** meal prepper
**I want to** plan meals for a week (or custom period)
**So that** I know what to cook each day and can batch-cook efficiently

```bash
recipe_box plan create --days 7
recipe_box plan add monday <recipe-id> --servings 4 --days 2  # leftovers for 2 days
recipe_box plan show
recipe_box plan clear
```

**Acceptance Criteria:**
- Assign recipes to days
- Specify servings and how many days it covers (leftovers)
- View plan as calendar/list
- Warn if days are unplanned

---

### UC4: Generate Shopping List

**As a** user
**I want to** generate a shopping list from my meal plan
**So that** I buy exactly what I need

```bash
recipe_box shopping generate              # from current plan
recipe_box shopping generate --days 3     # partial plan
recipe_box shopping show
recipe_box shopping add "olive oil"       # add custom item
recipe_box shopping remove "salt"         # already have it
recipe_box shopping clear
```

**Acceptance Criteria:**
- Aggregates ingredients from all planned recipes
- Combines quantities (e.g., 2 recipes need onions → total onions)
- Allows manual additions/removals
- Groups by category (produce, dairy, pantry, etc.)

---

### UC5: Localization (German/English)

**As a** German/English speaker
**I want to** use the CLI in my language
**So that** recipes and shopping lists are in my preferred language

```bash
recipe_box config set language de
recipe_box config set language en
```

**Acceptance Criteria:**
- CLI output in selected language
- AI generates recipes in selected language
- Shopping list in selected language
- Ingredient names localized

---

## Data Model (Draft)

```
Recipe:
  - id: string
  - title: string
  - description: string
  - servings: int
  - prepTime: int (minutes)
  - cookTime: int (minutes)
  - ingredients: []Ingredient
  - steps: []string
  - tags: []string
  - source: "ai" | "manual" | "imported"
  - language: "en" | "de"

Ingredient:
  - name: string
  - quantity: float
  - unit: string
  - category: string (produce, dairy, meat, pantry, etc.)

MealPlan:
  - startDate: date
  - entries: []PlanEntry

PlanEntry:
  - date: date
  - recipeId: string
  - servings: int
  - coversDays: int (for leftovers)

ShoppingList:
  - items: []ShoppingItem
  - customItems: []string

ShoppingItem:
  - ingredient: Ingredient
  - totalQuantity: float
  - fromRecipes: []string (recipe IDs)
```

---

## MVP Scope (v0.1)

Focus on core loop first:

1. **Recipe generate** - AI generates recipe, display it
2. **Recipe save/list/view** - Basic CRUD
3. **Shopping generate** - From single recipe (no meal plan yet)
4. **Language config** - EN/DE support

### Out of Scope for MVP
- Meal planning (week view)
- Ingredient aggregation
- Recipe editing
- Categories/tags filtering

---

## Resolved Questions

1. **Storage location** → `~/.recipe_box/` (ADR-004)
2. **API key management** → Config file (ADR-005)
3. **Units** → Metric only for MVP (ADR-006)

## Open Questions

1. **Offline mode** - Cache AI responses? Allow offline recipe management?
