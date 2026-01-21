package i18n

var messagesDE = map[string]string{
	// General
	MsgError: "Fehler",

	// Config
	MsgConfigSet:      "Konfiguration aktualisiert: %s = %s",
	MsgConfigNotFound: "Schlüssel nicht gefunden: %s",

	// Recipe commands
	MsgRecipeNotFound:    "Rezept nicht gefunden: %s",
	MsgRecipeDeleted:     "Rezept gelöscht: %s",
	MsgRecipeAdded:       "Rezept hinzugefügt: %s (%s)",
	MsgRecipeSaved:       "Rezept in Sammlung gespeichert: %s (%s)",
	MsgRecipeListEmpty:   "Keine Rezepte gefunden. Nutze 'recipe_box recipe add' um eins hinzuzufügen.",
	MsgRecipeListFound:   "%d Rezept(e) gefunden:",
	MsgRecipeGenerating:  "Rezept wird generiert...",
	MsgRecipeSaveHint:    "Nutze 'recipe_box recipe save' um dieses Rezept zur Sammlung hinzuzufügen.",
	MsgRecipeNoGenerated: "Kein generiertes Rezept zum Speichern. Nutze zuerst 'recipe_box recipe generate'.",

	// Shopping list
	MsgShoppingEmpty:      "Einkaufsliste ist leer.",
	MsgShoppingEmptyHint:  "Nutze 'recipe_box shopping generate <rezept-id>' um Artikel hinzuzufügen.",
	MsgShoppingCleared:    "Einkaufsliste geleert.",
	MsgShoppingAdded:      "%d Zutaten von '%s' zur Einkaufsliste hinzugefügt",
	MsgShoppingShowHint:   "Nutze 'recipe_box shopping show' um deine Einkaufsliste anzuzeigen.",
	MsgShoppingListHeader: "Einkaufsliste (%d Artikel)",

	// Categories
	MsgCategoryProduce: "Obst & Gemüse",
	MsgCategoryDairy:   "Milchprodukte",
	MsgCategoryMeat:    "Fleisch",
	MsgCategoryPantry:  "Vorratskammer",
	MsgCategoryFrozen:  "Tiefkühl",
	MsgCategoryOther:   "Sonstiges",

	// Units
	MsgServings: "Portionen",
	MsgMinutes:  "Min",
	MsgItems:    "Artikel",

	// Recipe view labels
	MsgLabelServings:    "Portionen",
	MsgLabelPrepTime:    "Vorbereitungszeit",
	MsgLabelCookTime:    "Kochzeit",
	MsgLabelTotalTime:   "Gesamtzeit",
	MsgLabelTags:        "Tags",
	MsgLabelIngredients: "Zutaten",
	MsgLabelSteps:       "Zubereitung",

	// Recipe add
	MsgRecipeAddHeader:      "Neues Rezept hinzufügen",
	MsgRecipeAddIngredients: "Zutaten (leere Zeile zum Beenden)",
	MsgRecipeAddSteps:       "Schritte (leere Zeile zum Beenden)",

	// Prompts
	MsgPromptTitle:       "Titel",
	MsgPromptDescription: "Beschreibung (optional)",
	MsgPromptServings:    "Portionen",
	MsgPromptPrepTime:    "Vorbereitungszeit (Minuten)",
	MsgPromptCookTime:    "Kochzeit (Minuten)",
	MsgPromptIngredient:  "Zutat (z.B. '200 g Mehl')",
	MsgPromptStep:        "Schritt %d",
	MsgPromptTags:        "Tags (kommagetrennt, optional)",

	// Interactive generate prompts
	MsgGenerateHeader:      "Rezept generieren",
	MsgGeneratePrompt:      "Was für ein Rezept möchtest du?",
	MsgGenerateIngredients: "Zutaten verwenden (kommagetrennt, optional)",
	MsgGenerateCuisine:     "Küche (z.B. italienisch, mexikanisch, optional)",
	MsgGenerateVegetarian:  "Vegetarisch? (j/n)",
	MsgGenerateQuick:       "Schnelles Rezept unter 30 Min? (j/n)",
	MsgGenerateSavePrompt:  "Rezept in Sammlung speichern? (j/n)",
	MsgYes:                 "j",
	MsgNo:                  "n",
}
