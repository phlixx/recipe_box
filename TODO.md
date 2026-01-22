# TODO

See `ITERATIONS.md` for full iteration plan.

## Next: Iteration 10 - Meal Plan Management

- [ ] Create `cmd/plan_add.go` (`plan add <day> <recipe-id> [--servings N] [--days N]`)
- [ ] Create `cmd/plan_remove.go` (`plan remove <day>`)
- [ ] Add day name parsing ("monday", "tuesday" -> dates in current plan)
- [ ] Update `plan show` to display leftover visualization
- [ ] Add `AddEntry` and `RemoveEntry` methods to plan service
- [ ] Add tab completion for day names in interactive mode
- [ ] Add tab completion for recipe IDs for `plan add` command
- [ ] Add i18n keys for new strings (EN + DE)

## Testing Notes

- API mock server in `e2e/e2e_test.go` enables testing `recipe generate` without real API calls
- Set `ANTHROPIC_API_URL` env var to override API endpoint for testing
