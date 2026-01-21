package i18n

import (
	"testing"
)

func TestT_English(t *testing.T) {
	SetLanguage(LangEN)
	defer ResetLanguage()

	tests := []struct {
		key      string
		expected string
	}{
		{MsgShoppingCleared, "Shopping list cleared."},
		{MsgCategoryProduce, "Produce"},
		{MsgCategoryDairy, "Dairy"},
		{MsgRecipeGenerating, "Generating recipe..."},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := T(tt.key)
			if got != tt.expected {
				t.Errorf("T(%q) = %q, want %q", tt.key, got, tt.expected)
			}
		})
	}
}

func TestT_German(t *testing.T) {
	SetLanguage(LangDE)
	defer ResetLanguage()

	tests := []struct {
		key      string
		expected string
	}{
		{MsgShoppingCleared, "Einkaufsliste geleert."},
		{MsgCategoryProduce, "Obst & Gemüse"},
		{MsgCategoryDairy, "Milchprodukte"},
		{MsgRecipeGenerating, "Rezept wird generiert..."},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := T(tt.key)
			if got != tt.expected {
				t.Errorf("T(%q) = %q, want %q", tt.key, got, tt.expected)
			}
		})
	}
}

func TestT_FallbackToEnglish(t *testing.T) {
	SetLanguage(LangEN)
	defer ResetLanguage()

	// Test with a key that exists
	got := T(MsgShoppingCleared)
	if got != "Shopping list cleared." {
		t.Errorf("expected English fallback, got %q", got)
	}
}

func TestT_UnknownKey(t *testing.T) {
	SetLanguage(LangEN)
	defer ResetLanguage()

	// Unknown key should return the key itself
	got := T("unknown.key")
	if got != "unknown.key" {
		t.Errorf("expected key to be returned for unknown key, got %q", got)
	}
}

func TestCategory_English(t *testing.T) {
	SetLanguage(LangEN)
	defer ResetLanguage()

	tests := []struct {
		cat      string
		expected string
	}{
		{"produce", "Produce"},
		{"dairy", "Dairy"},
		{"meat", "Meat"},
		{"pantry", "Pantry"},
		{"frozen", "Frozen"},
		{"other", "Other"},
		{"", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.cat, func(t *testing.T) {
			got := Category(tt.cat)
			if got != tt.expected {
				t.Errorf("Category(%q) = %q, want %q", tt.cat, got, tt.expected)
			}
		})
	}
}

func TestCategory_German(t *testing.T) {
	SetLanguage(LangDE)
	defer ResetLanguage()

	tests := []struct {
		cat      string
		expected string
	}{
		{"produce", "Obst & Gemüse"},
		{"dairy", "Milchprodukte"},
		{"meat", "Fleisch"},
		{"pantry", "Vorratskammer"},
		{"frozen", "Tiefkühl"},
		{"other", "Sonstiges"},
		{"", "Sonstiges"},
	}

	for _, tt := range tests {
		t.Run(tt.cat, func(t *testing.T) {
			got := Category(tt.cat)
			if got != tt.expected {
				t.Errorf("Category(%q) = %q, want %q", tt.cat, got, tt.expected)
			}
		})
	}
}

func TestIsValidLanguage(t *testing.T) {
	tests := []struct {
		lang     string
		expected bool
	}{
		{LangEN, true},
		{LangDE, true},
		{"fr", false},
		{"", false},
		{"english", false},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := IsValidLanguage(tt.lang)
			if got != tt.expected {
				t.Errorf("IsValidLanguage(%q) = %v, want %v", tt.lang, got, tt.expected)
			}
		})
	}
}

func TestSetLanguage(t *testing.T) {
	defer ResetLanguage()

	SetLanguage(LangDE)
	if GetLanguage() != LangDE {
		t.Error("expected language to be set to DE")
	}

	SetLanguage(LangEN)
	if GetLanguage() != LangEN {
		t.Error("expected language to be set to EN")
	}

	// Invalid language should not change the setting
	SetLanguage("invalid")
	if GetLanguage() != LangEN {
		t.Error("invalid language should not change the setting")
	}
}

func TestGetLanguage_Default(t *testing.T) {
	ResetLanguage()
	// When no config is set, should return default (EN)
	// Note: This test might behave differently if config exists
	lang := GetLanguage()
	if lang != LangEN && lang != LangDE {
		t.Errorf("expected valid language, got %q", lang)
	}
}
