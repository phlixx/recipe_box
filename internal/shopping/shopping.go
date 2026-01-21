package shopping

import (
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
)

const shoppingListFile = "shopping_list.json"

// Item represents a shopping list item
type Item struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Category string  `json:"category"`
	RecipeID string  `json:"recipe_id"`
}

// List represents a shopping list
type List struct {
	Items []Item `json:"items"`
}

// Storage defines the interface for persistence
type Storage interface {
	Save(filename string, data any) error
	Load(filename string, dest any) error
	Exists(filename string) bool
	Delete(filename string) error
}

// Service handles shopping list operations
type Service struct {
	storage Storage
}

// NewService creates a new shopping service
func NewService(s *storage.Storage) *Service {
	return &Service{storage: s}
}

// GenerateFromRecipe creates shopping list items from a recipe
func (s *Service) GenerateFromRecipe(r *recipe.Recipe) error {
	list, err := s.Get()
	if err != nil {
		list = &List{Items: []Item{}}
	}

	for _, ing := range r.Ingredients {
		item := Item{
			Name:     ing.Name,
			Quantity: ing.Quantity,
			Unit:     ing.Unit,
			Category: ing.Category,
			RecipeID: r.ID,
		}
		list.Items = mergeItem(list.Items, item)
	}

	return s.storage.Save(shoppingListFile, list)
}

// Get returns the current shopping list
func (s *Service) Get() (*List, error) {
	var list List

	if !s.storage.Exists(shoppingListFile) {
		return &List{Items: []Item{}}, nil
	}

	err := s.storage.Load(shoppingListFile, &list)
	return &list, err
}

// Clear removes all items from the shopping list
func (s *Service) Clear() error {
	if !s.storage.Exists(shoppingListFile) {
		return nil
	}
	return s.storage.Delete(shoppingListFile)
}

// mergeItem adds an item to the list, combining quantities for same name+unit
func mergeItem(items []Item, newItem Item) []Item {
	for i, existing := range items {
		if existing.Name == newItem.Name && existing.Unit == newItem.Unit {
			items[i].Quantity += newItem.Quantity
			return items
		}
	}
	return append(items, newItem)
}

// GroupByCategory returns items grouped by category
func (l *List) GroupByCategory() map[string][]Item {
	groups := make(map[string][]Item)
	for _, item := range l.Items {
		cat := item.Category
		if cat == "" {
			cat = "other"
		}
		groups[cat] = append(groups[cat], item)
	}
	return groups
}
