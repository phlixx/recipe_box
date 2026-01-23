<p align="center">
  <h1 align="center">Recipe Box</h1>
  <p align="center">
    <strong>Your AI-powered kitchen companion in the terminal</strong>
  </p>
  <p align="center">
    Generate recipes, scale servings, build shopping lists — all from the command line.
  </p>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#usage">Usage</a> •
  <a href="#localization">Deutsch/English</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.21-blue" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
  <img src="https://img.shields.io/badge/powered%20by-Claude%20AI-orange" alt="Powered by Claude">
</p>

---

```
$ recipe_box
Recipe Box - Interactive Mode
Type 'help' for commands, 'exit' to quit

🍳 > recipe generate --cuisine italian --quick

    🍝 Pasta Aglio e Olio
    ═══════════════════════════════════════

    A classic Roman pasta dish - simple, fast, and delicious.
    Ready in 20 minutes. Serves 4.

    Ingredients:
      - 400g spaghetti
      - 6 cloves garlic, thinly sliced
      - 1/2 cup extra virgin olive oil
      - 1 tsp red pepper flakes
      - Fresh parsley, chopped
      ...

🍳 > recipe save
✓ Recipe saved: a3f8c2d1

🍳 > shopping generate a3f8c2d1
✓ Added 5 items to shopping list
```

---

## Why Recipe Box?

**Tired of browser tabs full of recipe sites?** Recipe Box keeps your recipes in one place, generates new ones on demand, and builds shopping lists automatically.

| What you get | How it works |
|--------------|--------------|
| **AI-generated recipes** | Describe what you want, get a complete recipe |
| **Smart shopping lists** | One command turns recipes into categorized lists |
| **Flexible scaling** | Cooking for 2 or 20? Scale any recipe instantly |
| **Works offline** | Your saved recipes are always available locally |
| **Terminal-native** | Fast, scriptable, no GUI needed |

---

## Features

- **AI Recipe Generation** — Generate recipes by cuisine, ingredients, dietary needs, or freeform prompts
- **AI Meal Planning** — Generate balanced weekly meal plans with one command
- **Recipe Collection** — Save, organize, and manage your personal recipe library
- **Shopping Lists** — Auto-generate shopping lists grouped by category (produce, dairy, etc.)
- **Servings Scaling** — Scale any recipe up or down with automatic quantity adjustments
- **Interactive Mode** — REPL with tab completion, history, and colored output
- **Bilingual** — Full English and German support (UI, recipes, units)

---

## Installation

### Prerequisites

- Go 1.21+
- Claude API key → [console.anthropic.com](https://console.anthropic.com)

### Build from Source

```bash
git clone https://github.com/phlixx/recipe_box.git
cd recipe_box
make build
```

Binary: `./bin/recipe_box`

---

## Quick Start

```bash
# 1. Configure your API key
recipe_box config set api_key YOUR_API_KEY

# 2. Generate your first recipe
recipe_box recipe generate --prompt "quick healthy lunch"

# 3. Save it to your collection
recipe_box recipe save

# 4. Create a shopping list
recipe_box shopping generate <recipe-id>
recipe_box shopping show
```

Or just run `recipe_box` to enter interactive mode with tab completion.

---

## Usage

### Interactive Mode

Launch without arguments for the full experience:

```
$ recipe_box

🍳 > help
Available commands:
  recipe      Manage recipes
  shopping    Manage shopping list
  config      Settings
  help        Show this help
  exit        Quit

🍳 > recipe gen<TAB>     # Tab completion
🍳 > recipe view a3f<TAB> # Completes recipe IDs too
```

### Recipe Commands

```bash
# Generate with various options
recipe_box recipe generate --prompt "comfort food for winter"
recipe_box recipe generate --ingredients "salmon, asparagus, lemon"
recipe_box recipe generate --cuisine mexican --vegetarian
recipe_box recipe generate --quick --servings 6

# Manage your collection
recipe_box recipe list
recipe_box recipe view <id>
recipe_box recipe view <id> --servings 4  # Scale it
recipe_box recipe delete <id>
recipe_box recipe save                     # Save last generated
recipe_box recipe add                      # Add manually
```

### Meal Planning

```bash
# Create a meal plan
recipe_box plan create --days 7

# AI-generate a balanced week
recipe_box plan generate --days 7
recipe_box plan generate --vegetarian --quick
recipe_box plan generate --cuisine italian

# Manually add recipes to days
recipe_box plan add monday <recipe-id> --servings 4
recipe_box plan add tuesday <recipe-id> --days 2  # Covers leftovers

# View and manage
recipe_box plan show
recipe_box plan remove wednesday
recipe_box plan clear
```

### Shopping Lists

```bash
# Generate from one or more recipes
recipe_box shopping generate <id>
recipe_box shopping generate <id1> <id2> --servings 4

# Generate from meal plan
recipe_box shopping generate --plan           # Entire plan
recipe_box shopping generate --plan --days 3  # First 3 days only

# View and manage
recipe_box shopping show    # Grouped by category
recipe_box shopping clear
```

### Configuration

```bash
recipe_box config set api_key YOUR_KEY
recipe_box config set language de    # Switch to German
recipe_box config get language
```

---

## Localization

Recipe Box speaks English and German. Set your language:

```bash
recipe_box config set language de
```

| What changes | Example |
|--------------|---------|
| All CLI messages | "Rezept gespeichert" |
| AI-generated recipes | Recipes in German |
| Shopping categories | "Gemüse", "Milchprodukte" |
| Cooking units | tbsp → EL, cups → Tassen |

---

## Data Storage

Everything stays on your machine:

```
~/.recipe_box/
├── config.json      # API key, language
├── recipes/         # Your saved recipes
└── shopping.json    # Current shopping list
```

---

## Development

```bash
make build          # Build binary
make test           # Run tests
make check          # fmt + vet + test
make smoke          # Quick smoke test
make test-coverage  # Coverage report
```

<details>
<summary>Project Structure</summary>

```
recipe_box/
├── cmd/              # CLI commands (Cobra)
├── internal/
│   ├── ai/           # Claude API client
│   ├── config/       # Configuration
│   ├── i18n/         # Localization (EN/DE)
│   ├── recipe/       # Recipe domain
│   ├── shopping/     # Shopping lists
│   ├── storage/      # JSON persistence
│   └── ui/           # Terminal formatting
├── e2e/              # End-to-end tests
└── main.go
```

</details>

---

## License

MIT — use it, fork it, cook with it.

## Contributing

PRs welcome! Run `make check` before submitting.
