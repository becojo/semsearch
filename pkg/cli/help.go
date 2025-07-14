package cli

import "strings"

var help string = `Usage: semsearch [options]

Pattern options:
  -l,   --language <language>               Add a language to the rule
  -p,   --pattern <pattern>                 Pattern to search for
  -pi,  --pattern-inside <pattern>          Pattern to search for inside the matched pattern
  -pni, --pattern-not-inside <pattern>      Pattern to search for not inside the matched pattern
  -pr,  --pattern-regex <pattern>           Pattern to search for using a regex
  -pnr, --pattern-not-regex <pattern>       Pattern to search for not using a regex
  -mr,  --metavariable-regex <name=regex>   Metavariable to search for using a regex
  -fm,  --focus-metavariable <name>         Metavariable to focus on

Pattern group options:
  -ps,  --patterns                         Start a pattern group where all patterns must match
  -pe,  --pattern-either <pattern>         Start a pattern group where any pattern may match
  -mp,  --metavariable-pattern <name>      Start a metavariable pattern group
  -psk, --pattern-sinks                    Set the pattern sinks for the current rule
  -pso, --pattern-sources                  Set the pattern sources for the current rule
  ^,    --pop                              Exit the current pattern group

Search options:
  -i,   --path <path>                       Add the path to the search
  -e,   --eval <string>                     Evaluate the rule on the given string

Rule options:
  --id  <id>                                Rule ID
  -m,   --message <message>                 Message to display
  -fx,  --fix <replacement>                 Fix pattern
  -af,  --autofix                           Write fixes
  --severity <severity>                     Set the severity of the rule
  --metadata <key=value>                    Add metadata to the rule
  --option <key=value>                      Set an option for the rule
  --path-include <path>                     Limit the search to the specified path
  --path-exclude <path>                     Exclude the specified path from the search
  --rule                                    Start a new rule

Run options:
  -f,   --format <format>                   Output format (json, text, sarif, vim)
  -c,   --config <config>                   Add additional rules
  --debug                                   Output semsearch debug information
  --verbose                                 Enable Opengrep verbose mode
  --export                                  Output the rule instead of running Opengrep
`

func Help() string {
	return strings.TrimSpace(help)
}
