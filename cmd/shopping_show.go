package cmd

import (
	"fmt"
	"sort"
	"strings"

	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"

	"github.com/spf13/cobra"
)

var shoppingShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current shopping list",
	Long:  `Show all items in your shopping list, grouped by category.`,
	RunE:  runShoppingShow,
}

func init() {
	shoppingCmd.AddCommand(shoppingShowCmd)
}

func runShoppingShow(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	svc := shopping.NewService(store)
	list, err := svc.Get()
	if err != nil {
		return fmt.Errorf("failed to get shopping list: %w", err)
	}

	if len(list.Items) == 0 {
		fmt.Println("Shopping list is empty.")
		fmt.Println("Use 'recipe_box shopping generate <recipe-id>' to add items.")
		return nil
	}

	groups := list.GroupByCategory()

	// Sort categories for consistent output
	categories := make([]string, 0, len(groups))
	for cat := range groups {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	fmt.Printf("Shopping List (%d items)\n", len(list.Items))
	fmt.Println(strings.Repeat("=", 30))

	for _, cat := range categories {
		items := groups[cat]
		fmt.Printf("\n%s:\n", formatCategory(cat))
		for _, item := range items {
			fmt.Printf("  - %s\n", formatItem(item))
		}
	}

	return nil
}

func formatCategory(cat string) string {
	if cat == "" {
		return "Other"
	}
	return strings.ToUpper(cat[:1]) + cat[1:]
}

func formatItem(item shopping.Item) string {
	if item.Quantity == 0 && item.Unit == "" {
		return item.Name
	}
	if item.Unit == "" {
		return fmt.Sprintf("%g %s", item.Quantity, item.Name)
	}
	return fmt.Sprintf("%g %s %s", item.Quantity, item.Unit, item.Name)
}
