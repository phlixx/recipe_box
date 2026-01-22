package recipe

type Recipe struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Servings    int          `json:"servings"`
	PrepTime    int          `json:"prep_time"`
	CookTime    int          `json:"cook_time"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
	Tags        []string     `json:"tags"`
	Source      Source       `json:"source"`
	Language    Language     `json:"language"`
	AsciiArt    string       `json:"ascii_art,omitempty"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Category string  `json:"category"`
}

type Source string

const (
	SourceAI       Source = "ai"
	SourceManual   Source = "manual"
	SourceImported Source = "imported"
)

type Language string

const (
	LangEN Language = "en"
	LangDE Language = "de"
)

type Category string

const (
	CategoryProduce Category = "produce"
	CategoryDairy   Category = "dairy"
	CategoryMeat    Category = "meat"
	CategoryPantry  Category = "pantry"
	CategoryFrozen  Category = "frozen"
	CategoryOther   Category = "other"
)

// Scale returns a copy of the recipe with ingredient quantities scaled to the target servings.
// The original recipe is not modified.
func (r *Recipe) Scale(targetServings int) *Recipe {
	if targetServings <= 0 || r.Servings <= 0 {
		return r
	}

	factor := float64(targetServings) / float64(r.Servings)

	// Create a copy of the recipe
	scaled := *r
	scaled.Servings = targetServings

	// Scale ingredients
	scaled.Ingredients = make([]Ingredient, len(r.Ingredients))
	for i, ing := range r.Ingredients {
		scaled.Ingredients[i] = Ingredient{
			Name:     ing.Name,
			Quantity: ing.Quantity * factor,
			Unit:     ing.Unit,
			Category: ing.Category,
		}
	}

	// Copy steps and tags (they don't need scaling)
	scaled.Steps = make([]string, len(r.Steps))
	copy(scaled.Steps, r.Steps)

	scaled.Tags = make([]string, len(r.Tags))
	copy(scaled.Tags, r.Tags)

	return &scaled
}
