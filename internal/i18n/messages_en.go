package i18n

var messagesEN = map[string]string{
	// General
	MsgError: "Error",

	// Config
	MsgConfigSet:      "Configuration updated: %s = %s",
	MsgConfigNotFound: "Key not found: %s",

	// Recipe commands
	MsgRecipeNotFound:    "Recipe not found: %s",
	MsgRecipeDeleted:     "Recipe deleted: %s",
	MsgRecipeAdded:       "Recipe added: %s (%s)",
	MsgRecipeSaved:       "Recipe saved to collection: %s (%s)",
	MsgRecipeListEmpty:   "No recipes found. Use 'recipe_box recipe add' to add one.",
	MsgRecipeListFound:   "Found %d recipe(s):",
	MsgRecipeGenerating:  "Generating recipe...",
	MsgRecipeSaveHint:    "Use 'recipe_box recipe save' to add this recipe to your collection.",
	MsgRecipeNoGenerated: "No generated recipe to save. Use 'recipe_box recipe generate' first.",

	// Shopping list
	MsgShoppingEmpty:          "Shopping list is empty.",
	MsgShoppingEmptyHint:      "Use 'recipe_box shopping generate <recipe-id>' to add items.",
	MsgShoppingCleared:        "Shopping list cleared.",
	MsgShoppingAdded:          "Added %d ingredients from '%s' to shopping list",
	MsgShoppingShowHint:       "Use 'recipe_box shopping show' to view your shopping list.",
	MsgShoppingListHeader:     "Shopping List (%d items)",
	MsgShoppingGenerateNoArgs: "Please provide a recipe ID or use --plan to generate from meal plan",
	MsgShoppingPlanNoRecipes:  "No recipes found in the meal plan for the specified days.",

	// Categories
	MsgCategoryProduce: "Produce",
	MsgCategoryDairy:   "Dairy",
	MsgCategoryMeat:    "Meat",
	MsgCategoryPantry:  "Pantry",
	MsgCategoryFrozen:  "Frozen",
	MsgCategoryOther:   "Other",

	// Units
	MsgServings: "servings",
	MsgMinutes:  "min",
	MsgItems:    "items",

	// Recipe view labels
	MsgLabelServings:    "Servings",
	MsgLabelPrepTime:    "Prep time",
	MsgLabelCookTime:    "Cook time",
	MsgLabelTotalTime:   "Total time",
	MsgLabelTags:        "Tags",
	MsgLabelIngredients: "Ingredients",
	MsgLabelSteps:       "Steps",

	// Recipe add
	MsgRecipeAddHeader:      "Add a new recipe",
	MsgRecipeAddIngredients: "Ingredients (enter empty line when done)",
	MsgRecipeAddSteps:       "Steps (enter empty line when done)",

	// Prompts
	MsgPromptTitle:       "Title",
	MsgPromptDescription: "Description (optional)",
	MsgPromptServings:    "Servings",
	MsgPromptPrepTime:    "Prep time (minutes)",
	MsgPromptCookTime:    "Cook time (minutes)",
	MsgPromptIngredient:  "Ingredient (e.g., '2 cups flour')",
	MsgPromptStep:        "Step %d",
	MsgPromptTags:        "Tags (comma-separated, optional)",

	// Interactive generate prompts
	MsgGenerateHeader:      "Generate a Recipe",
	MsgGeneratePrompt:      "What kind of recipe would you like?",
	MsgGenerateIngredients: "Ingredients to use (comma-separated, optional)",
	MsgGenerateCuisine:     "Cuisine type (e.g., italian, mexican, optional)",
	MsgGenerateVegetarian:  "Vegetarian? (y/n)",
	MsgGenerateQuick:       "Quick recipe under 30 min? (y/n)",
	MsgGenerateSavePrompt:  "Save this recipe to your collection? (y/n)",
	MsgYes:                 "y",
	MsgNo:                  "n",

	// Scaling
	MsgScaledFrom: "(scaled from %d servings)",

	// Generate servings prompt
	MsgGenerateServings: "Number of servings (default 4)",

	// Interactive mode
	MsgInteractiveMode: "Interactive Mode",
	MsgInteractiveHint: "Type 'help' for commands, 'exit' to quit",
	MsgGoodbye:         "Goodbye!",
	MsgAvailableCmd:    "Available Commands:",
	MsgOther:           "Other",

	// Command descriptions
	MsgCmdRootShort:     "A CLI for managing recipes and shopping lists",
	MsgCmdRootLong:      "Recipe Box is a command-line tool for managing your personal recipe collection.\n\nYou can add recipes manually, generate them with AI, and create shopping lists.",
	MsgCmdRecipeShort:   "Manage recipes",
	MsgCmdRecipeLong:    "Add, list, view, and delete recipes from your collection.",
	MsgCmdShoppingShort: "Manage shopping lists",
	MsgCmdShoppingLong:  "Generate, view, and clear shopping lists from your recipes.",
	MsgCmdConfigShort:   "Manage configuration",
	MsgCmdConfigLong:    "View and modify Recipe Box configuration settings.",
	MsgCmdHelpShort:     "Show available commands",
	MsgCmdExitShort:     "Exit interactive mode",

	// Recipe subcommand descriptions
	MsgCmdRecipeListShort:     "List all recipes",
	MsgCmdRecipeListLong:      "Display all recipes in your collection.",
	MsgCmdRecipeViewShort:     "View a recipe",
	MsgCmdRecipeViewLong:      "Display the full details of a recipe.\n\nUse --servings to scale ingredient quantities.",
	MsgCmdRecipeAddShort:      "Add a new recipe",
	MsgCmdRecipeAddLong:       "Interactively add a new recipe to your collection.",
	MsgCmdRecipeGenerateShort: "Generate a recipe with AI",
	MsgCmdRecipeGenerateLong:  "Generate a recipe using Claude AI.\n\nExamples:\n  recipe_box recipe generate --prompt \"quick pasta dish\"\n  recipe_box recipe generate --ingredients \"chicken, lemon, garlic\"\n  recipe_box recipe generate --cuisine italian --vegetarian\n  recipe_box recipe generate --quick --cuisine mexican",
	MsgCmdRecipeSaveShort:     "Save the last generated recipe",
	MsgCmdRecipeSaveLong:      "Save the most recently generated recipe to your collection.\n\nUse this after 'recipe generate' to permanently save a recipe you like.",
	MsgCmdRecipeDeleteShort:   "Delete a recipe",
	MsgCmdRecipeDeleteLong:    "Remove a recipe from your collection.",

	// Shopping subcommand descriptions
	MsgCmdShoppingGenerateShort: "Generate shopping list from recipe or meal plan",
	MsgCmdShoppingGenerateLong:  "Add ingredients from recipes to your shopping list.\n\nUse --plan to generate from your meal plan (aggregates all recipes).\nUse --days with --plan to limit to the first N days.\nUse --servings to scale ingredient quantities (for single recipes).\n\nExamples:\n  recipe_box shopping generate abc123\n  recipe_box shopping generate abc123 def456\n  recipe_box shopping generate abc123 --servings 8\n  recipe_box shopping generate --plan\n  recipe_box shopping generate --plan --days 3",
	MsgCmdShoppingShowShort:     "Show current shopping list",
	MsgCmdShoppingShowLong:      "Show all items in your shopping list, grouped by category.",
	MsgCmdShoppingClearShort:    "Clear the shopping list",
	MsgCmdShoppingClearLong:     "Remove all items from your shopping list.",

	// Config subcommand descriptions
	MsgCmdConfigGetShort: "Get a configuration value",
	MsgCmdConfigSetShort: "Set a configuration value",

	// Plan command descriptions
	MsgCmdPlanShort:         "Manage meal plans",
	MsgCmdPlanLong:          "Create, view, and manage your weekly meal plans.",
	MsgCmdPlanCreateShort:   "Create a new meal plan",
	MsgCmdPlanCreateLong:    "Create a new meal plan starting from today.\n\nExamples:\n  recipe_box plan create\n  recipe_box plan create --days 5",
	MsgCmdPlanShowShort:     "Show current meal plan",
	MsgCmdPlanShowLong:      "Display the current meal plan with all scheduled meals.",
	MsgCmdPlanClearShort:    "Clear the meal plan",
	MsgCmdPlanClearLong:     "Remove the current meal plan.",
	MsgCmdPlanAddShort:      "Add a recipe to a day",
	MsgCmdPlanAddLong:       "Add a recipe to a specific day in the meal plan.\n\nExamples:\n  recipe_box plan add monday abc123\n  recipe_box plan add tuesday abc123 --servings 4\n  recipe_box plan add wednesday abc123 --servings 6 --days 2",
	MsgCmdPlanRemoveShort:   "Remove a recipe from a day",
	MsgCmdPlanRemoveLong:    "Remove a recipe from a specific day in the meal plan.\n\nExamples:\n  recipe_box plan remove monday\n  recipe_box plan remove 2024-01-15",
	MsgCmdPlanGenerateShort: "Generate a meal plan with AI",
	MsgCmdPlanGenerateLong:  "Generate a balanced meal plan using Claude AI.\n\nExamples:\n  recipe_box plan generate\n  recipe_box plan generate --days 5\n  recipe_box plan generate --vegetarian\n  recipe_box plan generate --cuisine italian --quick",

	// Plan messages
	MsgPlanCreated:       "Meal plan created for %d days (starting %s)",
	MsgPlanCleared:       "Meal plan cleared.",
	MsgPlanEmpty:         "Meal plan is empty.",
	MsgPlanEmptyHint:     "Use 'plan add <day> <recipe-id>' to add meals.",
	MsgPlanNoPlan:        "No meal plan found.",
	MsgPlanNoPlanHint:    "Use 'plan create' to start a new meal plan.",
	MsgPlanHeader:        "Meal Plan (%d days)",
	MsgPlanLeftovers:     "leftovers from %s",
	MsgPlanEmptyDay:      "(empty)",
	MsgPlanDayFormat:     "%s (%s)",
	MsgPlanEntryAdded:    "Added '%s' to %s",
	MsgPlanEntryRemoved:  "Removed recipe from %s",
	MsgPlanDateNotInPlan: "Date '%s' is not in the current meal plan",
	MsgPlanCoversDays:    "covers %d days",

	// Plan generate messages
	MsgPlanGenerating:          "Generating meal plan...",
	MsgPlanGeneratedHeader:     "Suggested Meal Plan (%d days)",
	MsgPlanSuggestionDay:       "%s: %s",
	MsgPlanApprovePrompt:       "Apply this meal plan? (y/n)",
	MsgPlanApproved:            "Meal plan applied!",
	MsgPlanDiscarded:           "Meal plan discarded.",
	MsgPlanCreateRecipesPrompt: "Create full recipes for these meals? (y/n)",
	MsgPlanCreatingRecipes:     "Creating recipe: %s...",
	MsgPlanRecipeCreated:       "Recipe created: %s (%s)",
	MsgPlanAllRecipesCreated:   "All recipes created and added to plan.",
	MsgPlanServingsInfo:        "%d servings",
	MsgPlanVarietyNote:         "balanced variety with different cuisines and cooking styles",

	// Config key descriptions
	MsgConfigKeyAPIKey:   "Claude API key",
	MsgConfigKeyLanguage: "Language (en/de)",

	// Error messages
	MsgErrNoAPIKey: "API key not configured. Run: recipe_box config set api_key <your-key>",
}
