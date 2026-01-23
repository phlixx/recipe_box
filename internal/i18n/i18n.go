package i18n

import (
	"recipe_box/internal/config"
)

// Message keys for localization
const (
	// General
	MsgError = "error"

	// Config
	MsgConfigSet      = "config.set"
	MsgConfigNotFound = "config.not_found"

	// Recipe commands
	MsgRecipeNotFound    = "recipe.not_found"
	MsgRecipeDeleted     = "recipe.deleted"
	MsgRecipeAdded       = "recipe.added"
	MsgRecipeSaved       = "recipe.saved"
	MsgRecipeListEmpty   = "recipe.list_empty"
	MsgRecipeListFound   = "recipe.list_found"
	MsgRecipeGenerating  = "recipe.generating"
	MsgRecipeSaveHint    = "recipe.save_hint"
	MsgRecipeNoGenerated = "recipe.no_generated"

	// Shopping list
	MsgShoppingEmpty          = "shopping.empty"
	MsgShoppingEmptyHint      = "shopping.empty_hint"
	MsgShoppingCleared        = "shopping.cleared"
	MsgShoppingAdded          = "shopping.added"
	MsgShoppingShowHint       = "shopping.show_hint"
	MsgShoppingListHeader     = "shopping.list_header"
	MsgShoppingGenerateNoArgs = "shopping.generate_no_args"
	MsgShoppingPlanNoRecipes  = "shopping.plan_no_recipes"

	// Categories
	MsgCategoryProduce = "category.produce"
	MsgCategoryDairy   = "category.dairy"
	MsgCategoryMeat    = "category.meat"
	MsgCategoryPantry  = "category.pantry"
	MsgCategoryFrozen  = "category.frozen"
	MsgCategoryOther   = "category.other"

	// Units (display names)
	MsgServings = "unit.servings"
	MsgMinutes  = "unit.minutes"
	MsgItems    = "unit.items"

	// Recipe view labels
	MsgLabelServings    = "label.servings"
	MsgLabelPrepTime    = "label.prep_time"
	MsgLabelCookTime    = "label.cook_time"
	MsgLabelTotalTime   = "label.total_time"
	MsgLabelTags        = "label.tags"
	MsgLabelIngredients = "label.ingredients"
	MsgLabelSteps       = "label.steps"

	// Recipe add
	MsgRecipeAddHeader      = "recipe.add_header"
	MsgRecipeAddIngredients = "recipe.add_ingredients"
	MsgRecipeAddSteps       = "recipe.add_steps"

	// Prompts
	MsgPromptTitle       = "prompt.title"
	MsgPromptDescription = "prompt.description"
	MsgPromptServings    = "prompt.servings"
	MsgPromptPrepTime    = "prompt.prep_time"
	MsgPromptCookTime    = "prompt.cook_time"
	MsgPromptIngredient  = "prompt.ingredient"
	MsgPromptStep        = "prompt.step"
	MsgPromptTags        = "prompt.tags"

	// Interactive generate prompts
	MsgGenerateHeader      = "generate.header"
	MsgGeneratePrompt      = "generate.prompt"
	MsgGenerateIngredients = "generate.ingredients"
	MsgGenerateCuisine     = "generate.cuisine"
	MsgGenerateVegetarian  = "generate.vegetarian"
	MsgGenerateQuick       = "generate.quick"
	MsgGenerateSavePrompt  = "generate.save_prompt"
	MsgYes                 = "yes"
	MsgNo                  = "no"

	// Scaling
	MsgScaledFrom = "scaled.from"

	// Generate servings prompt
	MsgGenerateServings = "generate.servings"

	// Interactive mode
	MsgInteractiveMode = "interactive.mode"
	MsgInteractiveHint = "interactive.hint"
	MsgGoodbye         = "interactive.goodbye"
	MsgAvailableCmd    = "interactive.available_commands"
	MsgOther           = "interactive.other"

	// Command descriptions (for Cobra and tab completion)
	MsgCmdRootShort     = "cmd.root.short"
	MsgCmdRootLong      = "cmd.root.long"
	MsgCmdRecipeShort   = "cmd.recipe.short"
	MsgCmdRecipeLong    = "cmd.recipe.long"
	MsgCmdShoppingShort = "cmd.shopping.short"
	MsgCmdShoppingLong  = "cmd.shopping.long"
	MsgCmdConfigShort   = "cmd.config.short"
	MsgCmdConfigLong    = "cmd.config.long"
	MsgCmdHelpShort     = "cmd.help.short"
	MsgCmdExitShort     = "cmd.exit.short"

	// Recipe subcommand descriptions
	MsgCmdRecipeListShort     = "cmd.recipe.list.short"
	MsgCmdRecipeListLong      = "cmd.recipe.list.long"
	MsgCmdRecipeViewShort     = "cmd.recipe.view.short"
	MsgCmdRecipeViewLong      = "cmd.recipe.view.long"
	MsgCmdRecipeAddShort      = "cmd.recipe.add.short"
	MsgCmdRecipeAddLong       = "cmd.recipe.add.long"
	MsgCmdRecipeGenerateShort = "cmd.recipe.generate.short"
	MsgCmdRecipeGenerateLong  = "cmd.recipe.generate.long"
	MsgCmdRecipeSaveShort     = "cmd.recipe.save.short"
	MsgCmdRecipeSaveLong      = "cmd.recipe.save.long"
	MsgCmdRecipeDeleteShort   = "cmd.recipe.delete.short"
	MsgCmdRecipeDeleteLong    = "cmd.recipe.delete.long"

	// Shopping subcommand descriptions
	MsgCmdShoppingGenerateShort = "cmd.shopping.generate.short"
	MsgCmdShoppingGenerateLong  = "cmd.shopping.generate.long"
	MsgCmdShoppingShowShort     = "cmd.shopping.show.short"
	MsgCmdShoppingShowLong      = "cmd.shopping.show.long"
	MsgCmdShoppingClearShort    = "cmd.shopping.clear.short"
	MsgCmdShoppingClearLong     = "cmd.shopping.clear.long"

	// Config subcommand descriptions
	MsgCmdConfigGetShort = "cmd.config.get.short"
	MsgCmdConfigSetShort = "cmd.config.set.short"

	// Plan command descriptions
	MsgCmdPlanShort         = "cmd.plan.short"
	MsgCmdPlanLong          = "cmd.plan.long"
	MsgCmdPlanCreateShort   = "cmd.plan.create.short"
	MsgCmdPlanCreateLong    = "cmd.plan.create.long"
	MsgCmdPlanShowShort     = "cmd.plan.show.short"
	MsgCmdPlanShowLong      = "cmd.plan.show.long"
	MsgCmdPlanClearShort    = "cmd.plan.clear.short"
	MsgCmdPlanClearLong     = "cmd.plan.clear.long"
	MsgCmdPlanAddShort      = "cmd.plan.add.short"
	MsgCmdPlanAddLong       = "cmd.plan.add.long"
	MsgCmdPlanRemoveShort   = "cmd.plan.remove.short"
	MsgCmdPlanRemoveLong    = "cmd.plan.remove.long"
	MsgCmdPlanGenerateShort = "cmd.plan.generate.short"
	MsgCmdPlanGenerateLong  = "cmd.plan.generate.long"

	// Plan messages
	MsgPlanCreated       = "plan.created"
	MsgPlanCleared       = "plan.cleared"
	MsgPlanEmpty         = "plan.empty"
	MsgPlanEmptyHint     = "plan.empty_hint"
	MsgPlanNoPlan        = "plan.no_plan"
	MsgPlanNoPlanHint    = "plan.no_plan_hint"
	MsgPlanHeader        = "plan.header"
	MsgPlanLeftovers     = "plan.leftovers"
	MsgPlanEmptyDay      = "plan.empty_day"
	MsgPlanDayFormat     = "plan.day_format"
	MsgPlanEntryAdded    = "plan.entry_added"
	MsgPlanEntryRemoved  = "plan.entry_removed"
	MsgPlanDateNotInPlan = "plan.date_not_in_plan"
	MsgPlanCoversDays    = "plan.covers_days"

	// Plan generate messages
	MsgPlanGenerating          = "plan.generating"
	MsgPlanGeneratedHeader     = "plan.generated_header"
	MsgPlanSuggestionDay       = "plan.suggestion_day"
	MsgPlanApprovePrompt       = "plan.approve_prompt"
	MsgPlanApproved            = "plan.approved"
	MsgPlanDiscarded           = "plan.discarded"
	MsgPlanCreateRecipesPrompt = "plan.create_recipes_prompt"
	MsgPlanCreatingRecipes     = "plan.creating_recipes"
	MsgPlanRecipeCreated       = "plan.recipe_created"
	MsgPlanAllRecipesCreated   = "plan.all_recipes_created"
	MsgPlanServingsInfo        = "plan.servings_info"
	MsgPlanVarietyNote         = "plan.variety_note"

	// Config key descriptions
	MsgConfigKeyAPIKey   = "config.key.api_key"
	MsgConfigKeyLanguage = "config.key.language"

	// Error messages
	MsgErrNoAPIKey = "error.no_api_key"
)

// Supported languages
const (
	LangEN = "en"
	LangDE = "de"
)

// DefaultLang is the fallback language
const DefaultLang = LangEN

// translations holds all message translations by language
var translations = map[string]map[string]string{
	LangEN: messagesEN,
	LangDE: messagesDE,
}

// currentLang caches the current language setting
var currentLang string

// GetLanguage returns the configured language or default
func GetLanguage() string {
	if currentLang != "" {
		return currentLang
	}

	cfg, err := config.New()
	if err != nil {
		return DefaultLang
	}

	lang, err := cfg.Get("language")
	if err != nil || (lang != LangEN && lang != LangDE) {
		return DefaultLang
	}

	currentLang = lang
	return lang
}

// SetLanguage sets the current language (for testing)
func SetLanguage(lang string) {
	if lang == LangEN || lang == LangDE {
		currentLang = lang
	}
}

// ResetLanguage clears the cached language
func ResetLanguage() {
	currentLang = ""
}

// T translates a message key to the current language
func T(key string) string {
	lang := GetLanguage()

	if msgs, ok := translations[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}

	// Fallback to English
	if msgs, ok := translations[LangEN]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}

	// Return key if not found
	return key
}

// Category returns the localized category name
func Category(cat string) string {
	switch cat {
	case "produce":
		return T(MsgCategoryProduce)
	case "dairy":
		return T(MsgCategoryDairy)
	case "meat":
		return T(MsgCategoryMeat)
	case "pantry":
		return T(MsgCategoryPantry)
	case "frozen":
		return T(MsgCategoryFrozen)
	case "other", "":
		return T(MsgCategoryOther)
	default:
		// Capitalize first letter for unknown categories
		if len(cat) > 0 {
			return string(cat[0]-32) + cat[1:]
		}
		return T(MsgCategoryOther)
	}
}

// IsValidLanguage checks if a language code is supported
func IsValidLanguage(lang string) bool {
	return lang == LangEN || lang == LangDE
}

// unitTranslations maps English unit names to localized versions
var unitTranslations = map[string]map[string]string{
	LangDE: {
		// Volume
		"tbsp":        "EL",
		"tablespoon":  "EL",
		"tablespoons": "EL",
		"tsp":         "TL",
		"teaspoon":    "TL",
		"teaspoons":   "TL",
		"cup":         "Tasse",
		"cups":        "Tassen",
		// Weight
		"oz":     "oz",
		"ounce":  "oz",
		"ounces": "oz",
		"lb":     "Pfund",
		"lbs":    "Pfund",
		"pound":  "Pfund",
		"pounds": "Pfund",
		// Metric units stay the same
		"g":  "g",
		"kg": "kg",
		"ml": "ml",
		"l":  "l",
		// Other common units
		"clove":  "Zehe",
		"cloves": "Zehen",
		"bunch":  "Bund",
		"pinch":  "Prise",
		"slice":  "Scheibe",
		"slices": "Scheiben",
		"piece":  "Stück",
		"pieces": "Stück",
		"can":    "Dose",
		"cans":   "Dosen",
	},
	// English keeps original units
	LangEN: {},
}

// Unit translates a cooking unit to the current language
func Unit(unit string) string {
	lang := GetLanguage()

	// Check if there's a translation for this language
	if langUnits, ok := unitTranslations[lang]; ok {
		if translated, ok := langUnits[unit]; ok {
			return translated
		}
	}

	// Return original unit if no translation found
	return unit
}
