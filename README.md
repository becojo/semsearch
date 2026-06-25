# semsearch

CLI interface to create and run Opengrep rules that are more complex that what the original Semgrep CLI can handle.

<!-- help start -->
```
Usage: semsearch [options]

Pattern options:
  -l    --language <language>               Add a language to the rule (default: generic)
  -p    --pattern <pattern>                 Pattern to match
  -pi   --pattern-inside <pattern>          Pattern to match inside the matched pattern
  -pni  --pattern-not-inside <pattern>      Pattern to match not inside the matched pattern
  -pr   --pattern-regex <pattern>           Pattern to match using a regex
  -pnr  --pattern-not-regex <pattern>       Pattern to match not using a regex
  -mr   --metavariable-regex <name=regex>   Metavariable to match using a regex
  -fm   --focus-metavariable <name>         Metavariable name to focus on

Pattern group options:
  -ps   --patterns [...]                    Start a pattern group where all patterns must match
  -pe   --pattern-either [...]              Start a pattern group where any pattern may match
  -mp   --metavariable-pattern <name> [...] Start a pattern group to match a metavariable
  -psk  --pattern-sinks [...]               Set the pattern sinks for the current rule
  -pso  --pattern-sources [...]             Set the pattern sources for the current rule
  ^     --pop                               Exit the current pattern group

Search options:
  -i    --path <path>                       Add the path to the search
  -e    --eval <string>                     Evaluate the rule on the given string

Rule options:
  -m    --message <message>                 Message to display
  -fx   --fix <pattern>                     Fix pattern
  -fr   --fix-regex <regex>                 Fix pattern using a regex
  -af   --autofix                           Automatically write fixes
  --id  <id>                                Rule ID
  --metadata <key=value>                    Add metadata to the rule
  --severity <severity>                     Set the severity of the rule
  --option <key=value>                      Set an option for the rule (such as: generic_ellipsis_max_span)
  --path-include <path>                     Limit the search to the specified path
  --path-exclude <path>                     Exclude the specified path from the search
  --rule                                    Start a new rule

Run options:
  -f    --format <format>                   Output format (json, text, sarif, vim, emacs)
  -c    --config <config>                   Add additional rules
  --debug                                   Output semsearch debug information
  --verbose                                 Enable Opengrep verbose mode
  --export                                  Output the rule instead of running Opengrep

Shell completion:
  --bash-completion                         Output bash completion script
```
<!-- help end -->

## Examples

Listing official GitHub Actions used in this repository:

```sh
semsearch -l yaml -p 'uses: "$USES"' -mr 'USES=actions' --path-include .github/workflows
```

```
    .github/workflows/ci.yml
    ❯❱ rule-1
           20┆ - uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0 # v7
            ⋮┆----------------------------------------
           23┆ - uses: actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16 # v6
```

Output the YAML rule used above instead of running it:
```sh
semsearch -l yaml -p 'uses: "$USES"' -mr 'USES=actions' --export
```

```yaml
rules:
- id: rule-1
  severity: WARNING
  message: ""
  languages:
  - yaml
  paths:
    include:
    - .github/workflows
  patterns:
  - pattern: 'uses: "$USES"'
  - metavariable-regex:
      metavariable: $USES
      regex: actions
```

## Installation

Download the [latest release](https://github.com/becojo/semsearch/releases) or install with `go install`:

```sh
go install github.com/becojo/semsearch/cmd/semsearch@latest
```
