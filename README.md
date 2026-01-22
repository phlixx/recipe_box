# Recipe Box CLI

A command-line tool for managing recipes and shopping lists, powered by Claude AI.

## Features

- **AI-Powered Recipe Generation** - Generate recipes using Claude AI with customizable criteria (cuisine, ingredients, dietary restrictions)
- **Recipe Management** - Add, list, view, and delete recipes from your personal collection
- **Shopping Lists** - Generate shopping lists from recipes with automatic ingredient categorization
- **Servings Scaling** - Scale recipes up or down, with automatic quantity adjustments
- **Localization** - Full support for English and German (DE/EN)
- **Interactive Mode** - REPL interface with tab completion and colored output

## Installation

### Prerequisites

- Go 1.21 or later
- Claude API key (get one at [console.anthropic.com](https://console.anthropic.com))

### Build from Source

```bash
git clone https://github.com/phlixx/recipe_box.git
cd recipe_box
make build
```

The binary will be at `./bin/recipe_box`.

## Quick Start

1. **Set your API key:**
   ```bash
   recipe_box config set api_key YOUR_API_KEY
   ```

2. **Set your language (optional, defaults to English):**
   ```bash
   recipe_box config set language de  # or "en"
   ```

3. **Generate a recipe:**
   ```bash
   recipe_box recipe generate --prompt "quick pasta dish"
   ```

4. **Save the generated recipe:**
   ```bash
   recipe_box recipe save
   ```

5. **Create a shopping list:**
   ```bash
   recipe_box shopping generate <recipe-id>
   recipe_box shopping show
   ```

## Usage

### Interactive Mode

Run `recipe_box` without arguments to enter interactive mode:

```
$ recipe_box
Recipe Box - Interactive Mode
Type 'help' for commands, 'exit' to quit

🍳 > recipe list
🍳 > recipe generate --cuisine italian
🍳 > exit
```

Features in interactive mode:
- Tab completion for commands and recipe IDs
- Command history (arrow keys)
- Colored output

### Recipe Commands

```bash
# Generate a recipe with AI
recipe_box recipe generate --prompt "quick weeknight dinner"
recipe_box recipe generate --ingredients "chicken, lemon, garlic"
recipe_box recipe generate --cuisine italian --vegetarian
recipe_box recipe generate --quick --servings 6

# Save the last generated recipe
recipe_box recipe save

# Manually add a recipe
recipe_box recipe add

# List all recipes
recipe_box recipe list

# View a recipe (with optional scaling)
recipe_box recipe view <id>
recipe_box recipe view <id> --servings 8

# Delete a recipe
recipe_box recipe delete <id>
```

### Shopping List Commands

```bash
# Add ingredients from a recipe to shopping list
recipe_box shopping generate <recipe-id>
recipe_box shopping generate <recipe-id> --servings 4

# Add from multiple recipes
recipe_box shopping generate <id1> <id2> <id3>

# View shopping list (grouped by category)
recipe_box shopping show

# Clear shopping list
recipe_box shopping clear
```

### Configuration

```bash
# Set API key
recipe_box config set api_key YOUR_KEY

# Set language (en or de)
recipe_box config set language de

# Get a config value
recipe_box config get language
```

## Localization

Recipe Box supports English and German. When set to German:

- CLI messages and help text appear in German
- AI generates recipes in German
- Shopping list categories are in German
- Cooking units are translated (tbsp → EL, cups → Tassen, tsp → TL)

```bash
recipe_box config set language de
```

## Data Storage

All data is stored locally in `~/.recipe_box/`:

- `config.json` - Configuration (API key, language)
- `recipes/` - Saved recipes as JSON files
- `shopping.json` - Current shopping list

## Development

### Prerequisites

- Go 1.21+
- Make

### Commands

```bash
make build        # Build the binary
make test         # Run tests
make check        # Run fmt, vet, and test
make test-coverage # Run tests with coverage
make smoke        # Quick smoke test
make run ARGS="recipe list"  # Run with arguments
```

### Project Structure

```
recipe_box/
├── cmd/           # CLI commands (Cobra)
├── internal/
│   ├── ai/        # Claude API client
│   ├── config/    # Configuration management
│   ├── i18n/      # Localization (EN/DE)
│   ├── recipe/    # Recipe domain logic
│   ├── shopping/  # Shopping list logic
│   ├── storage/   # JSON persistence
│   └── ui/        # Terminal colors/formatting
├── e2e/           # End-to-end tests
└── main.go        # Entry point
```

## License

MIT

## Contributing

Contributions are welcome! Please read the existing code patterns and run `make check` before submitting PRs.
