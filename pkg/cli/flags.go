package cli

import (
	"strings"

	"github.com/becojo/semsearch/pkg/rule"
)

var shortcuts = map[string]string{
	// keep-sorted start
	"af":  "autofix",
	"c":   "config",
	"e":   "eval",
	"f":   "format",
	"fm":  "focus-metavariable",
	"fr":  "fix-regex",
	"fx":  "fix",
	"i":   "path",
	"l":   "language",
	"m":   "message",
	"mp":  "metavariable-pattern",
	"mr":  "metavariable-regex",
	"p":   "pattern",
	"pe":  "pattern-either",
	"pi":  "pattern-inside",
	"pn":  "pattern-not",
	"pni": "pattern-not-inside",
	"pnr": "pattern-not-regex",
	"pr":  "pattern-regex",
	"ps":  "patterns",
	"psk": "pattern-sinks",
	"pso": "pattern-sources",
	"sv":  "severity",
	// keep-sorted end
}

// Flags not expecting a value
var flags0 = map[string]func(*rule.State){
	// keep-sorted start block=yes
	"autofix":         func(s *rule.State) { s.Autofix() },
	"debug":           func(s *rule.State) { s.Debug() },
	"export":          func(s *rule.State) { s.Export() },
	"pattern-either":  func(s *rule.State) { s.PatternEither() },
	"pattern-sinks":   func(s *rule.State) { s.PatternSinks() },
	"pattern-sources": func(s *rule.State) { s.PatternSources() },
	"patterns":        func(s *rule.State) { s.Patterns() },
	"pop":             func(s *rule.State) { s.Pop() },
	"rule":            func(s *rule.State) { s.Rule() },
	"semgrep":         func(s *rule.State) { s.Command("semgrep") },
	"verbose":         func(s *rule.State) { s.Verbose() },
	// keep-sorted end
}

// Flags expecting a value
var flags1 = map[string]func(*rule.State, string){
	// keep-sorted start block=yes
	"config":               func(s *rule.State, v string) { s.Config(v) },
	"eval":                 func(s *rule.State, v string) { s.Eval(v) },
	"fix":                  func(s *rule.State, v string) { s.Fix(v) },
	"fix-regex":            func(s *rule.State, v string) { s.FixRegex(v) },
	"focus-metavariable":   func(s *rule.State, v string) { s.FocusMetavariable(v) },
	"format":               func(s *rule.State, v string) { s.Format(v) },
	"id":                   func(s *rule.State, v string) { s.ID(v) },
	"language":             func(s *rule.State, v string) { s.Language(v) },
	"message":              func(s *rule.State, v string) { s.Message(v) },
	"metadata":             kv(func(s *rule.State, k string, v string) { s.Metadata(k, v) }),
	"metavariable-pattern": func(s *rule.State, v string) { s.MetavariablePattern(v) },
	"metavariable-regex":   kv(func(s *rule.State, k string, v string) { s.MetavariableRegex(k, v) }),
	"option":               kv(func(s *rule.State, k string, v string) { s.Option(k, v) }),
	"path":                 func(s *rule.State, v string) { s.Path(v) },
	"path-exclude":         func(s *rule.State, v string) { s.PathExclude(v) },
	"path-include":         func(s *rule.State, v string) { s.PathInclude(v) },
	"pattern":              func(s *rule.State, v string) { s.Pattern(v) },
	"pattern-inside":       func(s *rule.State, v string) { s.PatternInside(v) },
	"pattern-not":          func(s *rule.State, v string) { s.PatternNot(v) },
	"pattern-not-inside":   func(s *rule.State, v string) { s.PatternNotInside(v) },
	"pattern-not-regex":    func(s *rule.State, v string) { s.PatternNotRegex(v) },
	"pattern-regex":        func(s *rule.State, v string) { s.PatternRegex(v) },
	"severity":             func(s *rule.State, v string) { s.Severity(v) },

	// keep-sorted end
}

func init() {
	flags0["^"] = flags0["pop"]

	for format := range rule.Formats {
		flags0[format] = func(s *rule.State) { s.Format(format) }
	}
}

func kv(f func(s *rule.State, k string, v string)) func(*rule.State, string) {
	return func(s *rule.State, v string) {
		k, v, _ := strings.Cut(v, "=")
		f(s, k, v)
	}
}
