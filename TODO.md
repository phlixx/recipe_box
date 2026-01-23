# TODO

See `ITERATIONS.md` for full iteration plan.

## Next: Iteration 12 - AI Meal Planning

- [ ] Add `plan generate` command
- [ ] Add `--days N` flag (default 7)
- [ ] Add dietary/cuisine constraint flags (reuse from recipe generate)
- [ ] Implement `GenerateMealPlan` method in AI service
- [ ] Add approval flow for generated suggestions
- [ ] Option to auto-create recipes or just show suggestions
- [ ] Consider balance/variety in AI prompt
- [ ] Add i18n keys for new strings (EN + DE)

## Testing Notes

- API mock server in `e2e/e2e_test.go` enables testing `recipe generate` without real API calls
- Set `ANTHROPIC_API_URL` env var to override API endpoint for testing
