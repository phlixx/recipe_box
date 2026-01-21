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
