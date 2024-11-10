# semsearch

CLI interface to create and run Semgrep rules that are more complex that what the Semgrep CLI can handle.


```
Usage: semsearch [options]

Pattern options:
  -l,   --language <language>              Add a language to the rule
  -p,   --pattern <pattern>                Pattern to search for
  -pi,  --pattern-inside <pattern>         Pattern to search for inside the matched pattern
  -pni, --pattern-not-inside <pattern>     Pattern to search for not inside the matched pattern
  -pr,  --pattern-regex <pattern>          Pattern to search for using a regex
  -pnr, --pattern-not-regex <pattern>      Pattern to search for not using a regex
  -mr,  --metavariable-regex <name=regex>  Metavariable to search for using a regex
  -fm,  --focus-metavariable <name>        Metavariable to focus on

Pattern group options:
  -ps,  --patterns                          Start a pattern group where all patterns must match
  -pe,  --pattern-either <pattern>          Start a pattern group where any pattern may match
  ],    --pop                               Exit the current pattern group

Search options:
  -i,   --path <path>                       Add the path to the search
  -e,   --eval <string>                     Evaluate the rule on the given string

Rule options:
  --id  <id>                                Rule ID
  -f,   --format <format>                   Output format (json, text, sarif, vim)
  -c,   --config <config>                   Add additional rules
  -m,   --message <message>                 Message to display
  -fx,  --fix <pattern>                     Fix pattern
  -af,  --autofix                           Write fixes


Other options:
  --debug                                  Debug mode
  --export                                 Output the rule instead of running Semgrep
```
