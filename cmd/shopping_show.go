package cmd

import (
	"fmt"
	"sort"
	"strings"

	"recipe_box/internal/i18n"
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
		fmt.Println(i18n.T(i18n.MsgShoppingEmpty))
		fmt.Println(i18n.T(i18n.MsgShoppingEmptyHint))
		return nil
	}

	groups := list.GroupByCategory()

	// Sort categories for consistent output
	categories := make([]string, 0, len(groups))
	for cat := range groups {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	fmt.Printf(i18n.T(i18n.MsgShoppingListHeader)+"\n", len(list.Items))
	fmt.Println(strings.Repeat("=", 30))

	for _, cat := range categories {
		items := groups[cat]
		fmt.Printf("\n%s:\n", i18n.Category(cat))
		for _, item := range items {
			fmt.Printf("  - %s\n", formatItem(item))
		}
	}

	return nil
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
