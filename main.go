package main

import (
	"os"

	"recipe_box/cmd"
)

func main() {
	// If no arguments, start interactive mode (ADR-007)
	// Otherwise, execute as traditional CLI
	if len(os.Args) == 1 {
		cmd.RunInteractive()
	} else {
		cmd.Execute()
	}
}
