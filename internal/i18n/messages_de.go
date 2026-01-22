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

	// Scaling
	MsgScaledFrom: "(skaliert von %d Portionen)",

	// Generate servings prompt
	MsgGenerateServings: "Anzahl Portionen (Standard 4)",

	// Interactive mode
	MsgInteractiveMode: "Interaktiver Modus",
	MsgInteractiveHint: "Gib 'help' für Befehle ein, 'exit' zum Beenden",
	MsgGoodbye:         "Auf Wiedersehen!",
	MsgAvailableCmd:    "Verfügbare Befehle:",
	MsgOther:           "Sonstiges",

	// Command descriptions
	MsgCmdRootShort:     "CLI zum Verwalten von Rezepten und Einkaufslisten",
	MsgCmdRootLong:      "Recipe Box ist ein Kommandozeilen-Tool für deine persönliche Rezeptsammlung.\n\nDu kannst Rezepte manuell hinzufügen, mit KI generieren und Einkaufslisten erstellen.",
	MsgCmdRecipeShort:   "Rezepte verwalten",
	MsgCmdRecipeLong:    "Rezepte hinzufügen, auflisten, anzeigen und löschen.",
	MsgCmdShoppingShort: "Einkaufslisten verwalten",
	MsgCmdShoppingLong:  "Einkaufslisten aus Rezepten erstellen, anzeigen und leeren.",
	MsgCmdConfigShort:   "Einstellungen verwalten",
	MsgCmdConfigLong:    "Recipe Box Einstellungen anzeigen und ändern.",
	MsgCmdHelpShort:     "Verfügbare Befehle anzeigen",
	MsgCmdExitShort:     "Interaktiven Modus beenden",

	// Recipe subcommand descriptions
	MsgCmdRecipeListShort:     "Alle Rezepte auflisten",
	MsgCmdRecipeListLong:      "Zeigt alle Rezepte deiner Sammlung an.",
	MsgCmdRecipeViewShort:     "Rezept anzeigen",
	MsgCmdRecipeViewLong:      "Zeigt alle Details eines Rezepts an.\n\nMit --servings kannst du die Zutatenmengen skalieren.",
	MsgCmdRecipeAddShort:      "Neues Rezept hinzufügen",
	MsgCmdRecipeAddLong:       "Fügt interaktiv ein neues Rezept zu deiner Sammlung hinzu.",
	MsgCmdRecipeGenerateShort: "Rezept mit KI generieren",
	MsgCmdRecipeGenerateLong:  "Generiert ein Rezept mit Claude AI.\n\nBeispiele:\n  recipe_box recipe generate --prompt \"schnelles Nudelgericht\"\n  recipe_box recipe generate --ingredients \"Hähnchen, Zitrone, Knoblauch\"\n  recipe_box recipe generate --cuisine italienisch --vegetarian\n  recipe_box recipe generate --quick --cuisine mexikanisch",
	MsgCmdRecipeSaveShort:     "Letztes generiertes Rezept speichern",
	MsgCmdRecipeSaveLong:      "Speichert das zuletzt generierte Rezept in deiner Sammlung.\n\nNutze dies nach 'recipe generate', um ein Rezept dauerhaft zu speichern.",
	MsgCmdRecipeDeleteShort:   "Rezept löschen",
	MsgCmdRecipeDeleteLong:    "Entfernt ein Rezept aus deiner Sammlung.",

	// Shopping subcommand descriptions
	MsgCmdShoppingGenerateShort: "Einkaufsliste aus Rezept erstellen",
	MsgCmdShoppingGenerateLong:  "Fügt Zutaten eines Rezepts zur Einkaufsliste hinzu.\n\nMit --servings kannst du die Mengen für eine andere Portionsanzahl skalieren.\n\nBeispiele:\n  recipe_box shopping generate abc123\n  recipe_box shopping generate abc123 def456\n  recipe_box shopping generate abc123 --servings 8",
	MsgCmdShoppingShowShort:     "Aktuelle Einkaufsliste anzeigen",
	MsgCmdShoppingShowLong:      "Zeigt alle Artikel der Einkaufsliste, nach Kategorien gruppiert.",
	MsgCmdShoppingClearShort:    "Einkaufsliste leeren",
	MsgCmdShoppingClearLong:     "Entfernt alle Artikel von der Einkaufsliste.",

	// Config subcommand descriptions
	MsgCmdConfigGetShort: "Einstellung abrufen",
	MsgCmdConfigSetShort: "Einstellung setzen",

	// Plan command descriptions
	MsgCmdPlanShort:       "Essenspläne verwalten",
	MsgCmdPlanLong:        "Wochenpläne erstellen, anzeigen und verwalten.",
	MsgCmdPlanCreateShort: "Neuen Essensplan erstellen",
	MsgCmdPlanCreateLong:  "Erstellt einen neuen Essensplan ab heute.\n\nBeispiele:\n  recipe_box plan create\n  recipe_box plan create --days 5",
	MsgCmdPlanShowShort:   "Aktuellen Essensplan anzeigen",
	MsgCmdPlanShowLong:    "Zeigt den aktuellen Essensplan mit allen geplanten Mahlzeiten.",
	MsgCmdPlanClearShort:  "Essensplan löschen",
	MsgCmdPlanClearLong:   "Entfernt den aktuellen Essensplan.",

	// Plan messages
	MsgPlanCreated:    "Essensplan für %d Tage erstellt (ab %s)",
	MsgPlanCleared:    "Essensplan gelöscht.",
	MsgPlanEmpty:      "Essensplan ist leer.",
	MsgPlanEmptyHint:  "Nutze 'plan add <tag> <rezept-id>' um Mahlzeiten hinzuzufügen.",
	MsgPlanNoPlan:     "Kein Essensplan gefunden.",
	MsgPlanNoPlanHint: "Nutze 'plan create' um einen neuen Essensplan zu erstellen.",
	MsgPlanHeader:     "Essensplan (%d Tage)",
	MsgPlanLeftovers:  "Reste von %s",
	MsgPlanEmptyDay:   "(leer)",
	MsgPlanDayFormat:  "%s (%s)",

	// Config key descriptions
	MsgConfigKeyAPIKey:   "Claude API-Schlüssel",
	MsgConfigKeyLanguage: "Sprache (en/de)",

	// Error messages
	MsgErrNoAPIKey: "API-Schlüssel nicht konfiguriert. Führe aus: recipe_box config set api_key <dein-schlüssel>",
}
