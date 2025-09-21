package main

import (
	"fmt"
	"os"

	"github.com/becojo/semsearch/pkg/cli"
	"github.com/becojo/semsearch/pkg/rule"
)

func main() {
	args := os.Args[1:]

	// Handle completion generation before normal CLI parsing
	if len(args) == 1 && args[0] == "--bash-completion" {
		fmt.Print(cli.GetBashCompletion())
		return
	}

	if showHelp(args) {
		fmt.Println(cli.Help())
		return
	}

	state, err := cli.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err.Error())
		os.Exit(1)
		return
	}

	if command := os.Getenv("SEMSEARCH_COMMAND"); command != "" {
		state.Command(command)
	}

	runner := rule.NewRunner(state)

	if err := runner.Prepare(); err != nil {
		fmt.Fprintln(os.Stderr, "error preparing runner:", err.Error())
		os.Exit(1)
		return
	}

	err = runner.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running semsearch:", err.Error())
	}

	err = runner.Cleanup()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error cleaning up:", err.Error())
	}

	if err != nil {
		os.Exit(1)
	}
}

func showHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}

	for _, arg := range args {
		if arg == "help" || arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}
