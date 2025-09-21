package cli

import (
	"fmt"

	"github.com/becojo/semsearch/pkg/rule"
)

func Parse(args []string) (*rule.State, error) {
	var cmd string
	var value string
	state := rule.Builder().Rule()

	for i := 0; i < len(args); i++ {
		cmd = normalizeShortcut(args[i])
		if cmd == "" {
			return nil, fmt.Errorf("invalid command: %s", args[i])
		}

		if f, ok := flags0[cmd]; ok {
			f(state)
			continue
		}

		f, ok := flags1[cmd]
		if !ok {
			return nil, fmt.Errorf("unknown command %s", cmd)
		}

		i += 1
		if i < len(args) {
			value = args[i]
		}

		f(state, value)
	}

	return state, nil
}

func normalizeShortcut(arg string) (cmd string) {
	if len(arg) > 2 && arg[0:2] == "--" {
		cmd = arg[2:]
	} else if arg[0] == '-' {
		cmd = shortcuts[arg[1:]]
	}
	return
}
