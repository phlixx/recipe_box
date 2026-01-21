# TODO

See `ITERATIONS.md` for full iteration plan.

## Current: Iteration 6 - Interactive Mode

- [ ] Add `go-prompt` dependency
- [ ] Add `fatih/color` or `lipgloss` dependency
- [ ] Create `internal/ui/colors.go` - color definitions
- [ ] Create `cmd/interactive.go` - REPL logic
- [ ] Implement emoji prompt `🍳 > `
- [ ] Tab completion for commands
- [ ] Tab completion for recipe IDs
- [ ] Command history
- [ ] `exit` command
- [ ] Color recipe output (headers, categories)
- [ ] Color shopping list output
- [ ] Color success/error messages
- [ ] Update `main.go` to default to interactive mode
- [ ] Keep traditional CLI working

## Up Next

**Iteration 7**: Servings scaling (`--servings` flag)
**Iteration 8**: Localization hardening (all strings, units)
