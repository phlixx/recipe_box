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
	MsgShoppingEmpty      = "shopping.empty"
	MsgShoppingEmptyHint  = "shopping.empty_hint"
	MsgShoppingCleared    = "shopping.cleared"
	MsgShoppingAdded      = "shopping.added"
	MsgShoppingShowHint   = "shopping.show_hint"
	MsgShoppingListHeader = "shopping.list_header"

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
