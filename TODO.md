# TODO

See `ITERATIONS.md` for full iteration plan.

## Current: Iteration 9 - Meal Plan Foundation

- [ ] Create `internal/plan/plan.go` (MealPlan, PlanEntry types)
- [ ] Create `internal/plan/service.go` (Service with CRUD methods)
- [ ] Create `internal/plan/service_test.go` (unit tests)
- [ ] Create `cmd/plan.go` (command group)
- [ ] Create `cmd/plan_create.go` (`plan create [--days N]`)
- [ ] Create `cmd/plan_show.go` (`plan show`)
- [ ] Create `cmd/plan_clear.go` (`plan clear`)
- [ ] Add i18n keys for all new strings (EN + DE)
- [ ] Update `cmd/interactive.go` (tab completion for plan commands)

## Testing Notes

- API mock server in `e2e/e2e_test.go` enables testing `recipe generate` without real API calls
- Set `ANTHROPIC_API_URL` env var to override API endpoint for testing
