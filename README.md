# semsearch

CLI interface to create and run Semgrep rules that are more complex that what the Semgrep CLI can handle.


```
Usage: semsearch [options]

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
  --debug                                   Debug mode
  --export                                  Output the rule instead of running Semgrep
```

## Examples


### Search functions related to `*State`
``` sh
semsearch -l go -fm F -pe -p 'func ($S *State) $F(...) {...}'  -p 'func $F(...) *State {...}'
```

``` sh
    cmd/semsearch.go
    ❯❱ id
          123┆ func NewState() *State {
             ┆----------------------------------------
          135┆ func (s *State) Args() []string {
             ┆----------------------------------------
          161┆ func (s *State) AddCondition(cond Condition) {
             ┆----------------------------------------
          166┆ func (s *State) Build(args []string) {
             ┆----------------------------------------
          276┆ func (s *State) Tempfile(name string) (*os.File, error) {
             ┆----------------------------------------
          285┆ func (s *State) Cleanup() {
             ┆----------------------------------------
          291┆ func (s *State) Prepare() {
             ┆----------------------------------------
          308┆ func (s *State) Exec() {
```


### Output the Semgrep rule instead of running it

``` sh
semsearch -l go -fm F -pe -p 'func ($S *State) $F(...) {...}'  -p 'func $F(...) *State {...}' --export
```

``` yaml
rules:
- id: id
  patterns:
  - focus-metavariable: $F
  - pattern-either:
    - pattern: func ($S *State) $F(...) {...}
    - pattern: func $F(...) *State {...}
  severity: WARNING
  message: ""
  languages:
  - go
```

