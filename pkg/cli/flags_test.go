package cli

import (
	"testing"
)

var allCommands = map[string]bool{}

func init() {
	for cmd := range flags0 {
		allCommands[cmd] = true
	}
	for cmd := range flags1 {
		allCommands[cmd] = true
	}
}

func TestShortcutsMapToCommands(t *testing.T) {
	for shortcut, cmd := range shortcuts {
		if _, ok := allCommands[cmd]; !ok {
			t.Errorf("shortcut '%s' points to unknown command '%s'", shortcut, cmd)
		}
	}
}

func TestAmbiguousShortcutsCommands(t *testing.T) {
	for shortcut := range shortcuts {
		if _, ok := allCommands[shortcut]; ok {
			t.Errorf("shortcut '%s' is ambiguous with a command", shortcut)
		}
	}
}
