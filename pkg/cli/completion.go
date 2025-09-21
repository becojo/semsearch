package cli

import _ "embed"

//go:embed completion.bash
var bashCompletion string

// GetBashCompletion returns the bash completion script
func GetBashCompletion() string {
	return bashCompletion
}
