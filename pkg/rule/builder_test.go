package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const expected = `rules:
- id: rule-1
  severity: INFO
  message: msg
  languages:
  - go
  - generic
  options:
    generic_ellipsis_max_span: "5"
  patterns:
  - focus-metavariable: $PKG
  - pattern-either:
    - pattern: |-
        package $PKG
        ...
        var $VAR = $VAL
    - metavariable-pattern:
        metavariable: $PKG
        patterns:
        - pattern: pkgname
  - patterns:
    - pattern-not: not this
    - pattern-not-inside: not inside this
  fix-regex:
    regex: match
    replacement: replacement
  metadata:
    key: value
- id: extra-taint-rule
  severity: INFO
  message: ""
  languages:
  - go
  - generic
  paths:
    include:
    - path/to/include
    - path/to/includetoo
    exclude:
    - path/to/exclude
  mode: taint
  pattern-sources:
  - pattern: source
  pattern-sinks:
  - pattern: sink
`

func TestBuilder(t *testing.T) {
	state := Builder().
		Rule().
		Language("go").
		Language("generic").
		Severity("INFO").
		Message("msg").
		FocusMetavariable("PKG").
		PatternEither().
		Pattern("package $PKG\n...\nvar $VAR = $VAL").
		MetavariablePattern("PKG").
		Pattern(`pkgname`).
		Pop().
		Pop().
		Patterns().
		PatternNot("not this").
		PatternNotInside("not inside this").
		Fix("replacement").
		FixRegex("match").
		Option("generic_ellipsis_max_span", "5").
		Metadata("key", "value").
		Rule().
		ID("extra-taint-rule").
		PatternSources().
		Pattern("source").
		PatternSinks().
		Pattern("sink").
		PathInclude("path/to/include").
		PathInclude("path/to/includetoo").
		PathExclude("path/to/exclude")

	assert.Equal(t, expected, string(state.MarshalRules()))
}
