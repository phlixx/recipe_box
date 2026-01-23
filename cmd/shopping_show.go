package cmd

import (
	"fmt"
	"sort"
	"strings"

	"recipe_box/internal/i18n"
	"recipe_box/internal/shopping"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var shoppingShowCmd = &cobra.Command{
	Use:  "show",
	RunE: runShoppingShow,
}

func init() {
	shoppingShowCmd.Short = i18n.T(i18n.MsgCmdShoppingShowShort)
	shoppingShowCmd.Long = i18n.T(i18n.MsgCmdShoppingShowLong)
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
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgShoppingEmpty))
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgShoppingEmptyHint))
		return nil
	}

	groups := list.GroupByCategory()

	// Sort categories for consistent output
	categories := make([]string, 0, len(groups))
	for cat := range groups {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	ui.TitlePrintf(i18n.T(i18n.MsgShoppingListHeader)+"\n", len(list.Items))
	fmt.Println(strings.Repeat("=", 30))

	for _, cat := range categories {
		items := groups[cat]
		fmt.Println()
		ui.CategoryPrintf("%s:\n", i18n.Category(cat))
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
	qtyStr := formatQuantity(item.Quantity)
	if item.Unit == "" {
		return fmt.Sprintf("%s %s", qtyStr, item.Name)
	}
	return fmt.Sprintf("%s %s %s", qtyStr, i18n.Unit(item.Unit), item.Name)
}

// formatQuantity formats a quantity with nice rounding
func formatQuantity(qty float64) string {
	// Round to 2 decimal places to avoid floating point display issues
	rounded := float64(int(qty*100+0.5)) / 100

	// If it's a whole number, display without decimals
	if rounded == float64(int(rounded)) {
		return fmt.Sprintf("%d", int(rounded))
	}

	// Otherwise display with minimal decimals
	return fmt.Sprintf("%g", rounded)
}
