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
	MsgShoppingEmpty:      "Shopping list is empty.",
	MsgShoppingEmptyHint:  "Use 'recipe_box shopping generate <recipe-id>' to add items.",
	MsgShoppingCleared:    "Shopping list cleared.",
	MsgShoppingAdded:      "Added %d ingredients from '%s' to shopping list",
	MsgShoppingShowHint:   "Use 'recipe_box shopping show' to view your shopping list.",
	MsgShoppingListHeader: "Shopping List (%d items)",

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
}
