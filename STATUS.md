# Status

## Current Focus

**Iteration 4: Shopping List** - Generate shopping lists from recipes.

## Blockers

None.

## Handoff Notes

**Last session (2026-01-21)**:
- Completed Iteration 3: AI Generation
- Switched from Claude API to Ollama (local AI, no cost)
- Created Ollama client (`internal/ai/ollama.go`)
- Implemented `recipe generate` with flags (--prompt, --ingredients, --cuisine, --vegetarian, --quick)
- Implemented `recipe save` to save last generated recipe
- Default model: llama3.2 (configurable via `config set model`)
- Updated ADR-005 to document Ollama choice
- 9 unit tests for AI client (24 total tests passing)

**To use**: Install Ollama, run `ollama pull llama3.2`, then `ollama serve`

**Next**: Start Iteration 4 - Shopping list commands.
