package cmd

import (
	"fmt"
	"strings"
	"time"

	"recipe_box/internal/i18n"
	"recipe_box/internal/plan"
	"recipe_box/internal/recipe"
	"recipe_box/internal/storage"
	"recipe_box/internal/ui"

	"github.com/spf13/cobra"
)

var planShowCmd = &cobra.Command{
	Use:  "show",
	RunE: runPlanShow,
}

func init() {
	planShowCmd.Short = i18n.T(i18n.MsgCmdPlanShowShort)
	planShowCmd.Long = i18n.T(i18n.MsgCmdPlanShowLong)
	planCmd.AddCommand(planShowCmd)
}

func runPlanShow(cmd *cobra.Command, args []string) error {
	store, err := storage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	planSvc := plan.NewService(store)
	p, err := planSvc.Get()
	if err != nil {
		return fmt.Errorf("failed to get meal plan: %w", err)
	}

	if p == nil {
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlan))
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanNoPlanHint))
		return nil
	}

	recipeSvc := recipe.NewService(store)

	dates := p.GetDates()
	hasEntries := len(p.Entries) > 0

	ui.TitlePrintf(i18n.T(i18n.MsgPlanHeader)+"\n", p.Days)
	fmt.Println(strings.Repeat("=", 30))

	for _, date := range dates {
		dayName := formatDayName(date)
		ui.SectionPrintf("\n%s:\n", fmt.Sprintf(i18n.T(i18n.MsgPlanDayFormat), dayName, date))

		// Check if this day has an entry
		entry := p.GetEntryForDate(date)
		if entry != nil {
			displayPlanEntry(recipeSvc, entry)
			continue
		}

		// Check if this day has leftovers
		isLeftover, sourceEntry := p.IsLeftoverDay(date)
		if isLeftover {
			displayLeftovers(recipeSvc, sourceEntry)
			continue
		}

		// Empty day
		ui.LabelPrintf("  %s\n", i18n.T(i18n.MsgPlanEmptyDay))
	}

	if !hasEntries {
		fmt.Println()
		ui.LabelPrintf("%s\n", i18n.T(i18n.MsgPlanEmptyHint))
	}

	return nil
}

func displayPlanEntry(recipeSvc *recipe.Service, entry *plan.PlanEntry) {
	r, err := recipeSvc.Get(entry.RecipeID)
	if err != nil {
		ui.ErrorPrintf("  %s (ID: %s)\n", i18n.T(i18n.MsgRecipeNotFound), entry.RecipeID)
		return
	}

	servingsInfo := ""
	if entry.Servings > 0 {
		servingsInfo = fmt.Sprintf(" (%d %s)", entry.Servings, i18n.T(i18n.MsgServings))
	}

	coversInfo := ""
	if entry.CoversDays > 1 {
		coversInfo = fmt.Sprintf(", covers %d days", entry.CoversDays)
	}

	fmt.Printf("  %s%s%s\n", r.Title, servingsInfo, coversInfo)
}

func displayLeftovers(recipeSvc *recipe.Service, sourceEntry *plan.PlanEntry) {
	r, err := recipeSvc.Get(sourceEntry.RecipeID)
	if err != nil {
		ui.LabelPrintf("  %s\n", i18n.T(i18n.MsgPlanLeftovers))
		return
	}

	sourceDayName := formatDayName(sourceEntry.Date)
	ui.LabelPrintf("  <- %s\n", fmt.Sprintf(i18n.T(i18n.MsgPlanLeftovers), sourceDayName+" ("+r.Title+")"))
}

func formatDayName(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Weekday().String()
}
