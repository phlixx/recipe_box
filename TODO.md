# TODO

See `ITERATIONS.md` for full iteration plan.

## Phase 2 Complete ✅

All Phase 2 (Meal Planning) iterations are complete:
- Iteration 9: Meal Plan Foundation ✅
- Iteration 10: Meal Plan Management ✅
- Iteration 11: Shopping from Plan ✅
- Iteration 12: AI Meal Planning ✅

## Future: Phase 3 - Advanced Features

Ideas for future development:
- Recipe import from URL
- Pantry tracking (track what you have)
- Export to PDF/Markdown
- Recipe sharing
- Nutritional information
- Recipe categories and filtering

## Testing Notes

- API mock server in `e2e/e2e_test.go` enables testing `recipe generate` and `plan generate` without real API calls
- Set `ANTHROPIC_API_URL` env var to override API endpoint for testing
