# TODO

See `ITERATIONS.md` for full iteration plan.

## Next: Iteration 11 - Shopping from Plan

- [ ] Add `--plan` flag to `shopping generate` command
- [ ] Add `--days N` flag to limit shopping generation to first N days of plan
- [ ] Implement `GenerateFromPlan` method in shopping service
- [ ] Aggregate ingredients across all plan entries
- [ ] Scale quantities by servings specified in plan entries
- [ ] Track which recipe each item came from (for display)
- [ ] Add i18n keys for new strings (EN + DE)
- [ ] Update help text in interactive mode

## Testing Notes

- API mock server in `e2e/e2e_test.go` enables testing `recipe generate` without real API calls
- Set `ANTHROPIC_API_URL` env var to override API endpoint for testing
