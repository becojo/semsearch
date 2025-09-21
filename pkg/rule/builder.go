package rule

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v2"
)

var Formats = map[string]bool{
	"json":           true,
	"vim":            true,
	"emacs":          true,
	"sarif":          true,
	"text":           true,
	"gitlab-sast":    true,
	"gitlab-secrets": true,
	"junit-xml":      true,
}

type State struct {
	// rules to run
	rules []*Rule
	// stack of patterns
	stack []*[]Pattern
	// paths to file or directories to scan
	paths []string
	// Output format
	format string
	// strings to evaluate
	evals []string
	// paths to additional rules
	configs []string
	// errors encountered during rule building
	warnings []string
	// autofix enabled
	autofix bool
	// debug mode
	debug bool
	// export mode
	export bool
	// command to run opengrep
	command string
	// opengrep verbose mode
	verbose bool
}

func Builder() *State {
	return &State{
		command: "opengrep",
		rules:   []*Rule{},
		format:  "text",
		stack:   []*[]Pattern{},
	}
}

func (s *State) pushPattern(p Pattern) {
	head := s.stack[len(s.stack)-1]
	if head == nil {
		s.warn("no pattern stack available to push pattern")
		return
	}
	*head = append(*head, p)
}

func (s *State) headRule() *Rule {
	if len(s.rules) == 0 {
		return &Rule{}
	}
	return s.rules[len(s.rules)-1]
}

func (s *State) warn(message string) {
	s.warnings = append(s.warnings, message)
}

// Set the output format of the findings.
func (s *State) Format(format string) *State {
	if !Formats[format] {
		s.warn(fmt.Sprintf("unknown output format '%s'", format))
	}
	s.format = format
	return s
}

// Enable export mode to output rules instead of running them.
func (s *State) Export() *State {
	s.export = true
	return s
}

// Serialize the rules to YAML format.
func (s *State) MarshalRules() []byte {
	y, err := yaml.Marshal(map[string]any{
		"rules": s.rules,
	})
	if err != nil {
		s.warn(fmt.Sprintf("failed to marshal rules: %v", err))
		return nil
	}
	return y
}

// Enable debug mode
func (s *State) Debug() *State {
	s.debug = true
	return s
}

// Create a new rule in the state.
func (s *State) Rule() *State {
	r := Rule{
		Id:       fmt.Sprintf("rule-%d", len(s.rules)+1),
		Patterns: &[]Pattern{},
		Severity: SEVERITY_WARNING,
		Metadata: map[string]any{},
		Options:  map[string]any{},
	}

	if len(s.rules) > 0 {
		h := s.headRule()
		r.Languages = append(r.Languages, h.Languages...)
		r.Severity = h.Severity
	}
	s.rules = append(s.rules, &r)
	s.stack = []*[]Pattern{r.Patterns}
	return s
}

// Run the rules on the provided path
func (s *State) Path(path string) *State {
	s.paths = append(s.paths, path)
	return s
}

// Set the pattern sources for the current rule.
func (s *State) PatternSources() *State {
	r := s.headRule()
	r.PatternSources = &[]Pattern{}
	s.stack = []*[]Pattern{r.PatternSources}
	return s
}

// Set the pattern sinks for the current rule.
func (s *State) PatternSinks() *State {
	r := s.headRule()
	r.PatternSinks = &[]Pattern{}
	s.stack = []*[]Pattern{r.PatternSinks}
	return s
}

// Set the language for the current rule.
func (s *State) Language(lang string) *State {
	r := s.headRule()
	r.Languages = append(r.Languages, lang)
	return s
}

// Add a pattern to the current rule.
func (s *State) Pattern(pattern string) *State {
	s.pushPattern(Pattern{Pattern: pattern})
	return s
}

// Add a pattern that matches either of the provided patterns.
func (s *State) PatternEither() *State {
	p := Pattern{PatternEither: &[]Pattern{}}
	s.pushPattern(p)
	s.stack = append(s.stack, p.PatternEither)
	return s
}

// Add a pattern that matches all of the provided patterns.
func (s *State) Patterns() *State {
	r := &[]Pattern{}
	s.pushPattern(Pattern{Patterns: r})
	s.stack = append(s.stack, r)
	return s
}

// Add a pattern that does not match the provided pattern.
func (s *State) PatternNot(pattern string) *State {
	s.pushPattern(Pattern{PatternNot: pattern})
	return s
}

// Add a pattern that matches inside the provided pattern.
func (s *State) PatternInside(pattern string) *State {
	s.pushPattern(Pattern{PatternInside: pattern})
	return s
}

// Add a pattern that does not match inside the provided pattern.
func (s *State) PatternNotInside(pattern string) *State {
	s.pushPattern(Pattern{PatternNotInside: pattern})
	return s
}

// Add a pattern that matches the provided regex.
func (s *State) PatternRegex(pattern string) *State {
	s.pushPattern(Pattern{PatternRegex: pattern})
	return s
}

// Add a pattern that does not match the provided regex.
func (s *State) PatternNotRegex(pattern string) *State {
	s.pushPattern(Pattern{PatternNotRegex: pattern})
	return s
}

// Set the focused metavariable for the current rule.
func (s *State) FocusMetavariable(metavariable string) *State {
	s.pushPattern(Pattern{
		FocusMetavariable: normalizeMetavariable(metavariable),
	})
	return s
}

// Add a metavariable regex pattern to the current rule.
func (s *State) MetavariableRegex(metavariable, regex string) *State {
	s.pushPattern(Pattern{
		MetavariableRegex: &MetavariableRegex{
			Metavariable: normalizeMetavariable(metavariable),
			Regex:        regex,
		},
	})
	return s
}

// Add a metavariable pattern to the current rule.
func (s *State) MetavariablePattern(metavariable string) *State {
	patterns := &[]Pattern{}
	s.pushPattern(Pattern{
		MetavariablePattern: &MetavariablePattern{
			Metavariable: normalizeMetavariable(metavariable),
			Patterns:     patterns,
		},
	})
	s.stack = append(s.stack, patterns)
	return s
}

// Exit the current pattern group
func (s *State) Pop() *State {
	if len(s.stack) > 1 {
		s.stack = s.stack[:len(s.stack)-1]
	} else {
		s.stack = []*[]Pattern{{}}
	}
	return s
}

// Add metadata to the current rule
func (s *State) Metadata(key string, value any) *State {
	r := s.headRule()
	r.Metadata[key] = value
	return s
}

// Set the current rule message
func (s *State) Message(message string) *State {
	r := s.headRule()
	r.Message = message
	return s
}

// Set the fix for the current rule
func (s *State) Fix(fix string) *State {
	r := s.headRule()
	r.Fix = fix
	return s
}

// Set the fix regex for the current rule
func (s *State) FixRegex(fixRegex string) *State {
	r := s.headRule()
	r.FixRegex = fixRegex
	return s
}

// Set the rule ID
func (s *State) ID(id string) *State {
	r := s.headRule()
	r.Id = id
	return s
}

// Set the severity for the current rule.
func (s *State) Severity(severity string) *State {
	r := s.headRule()
	severity = strings.ToUpper(severity)
	if _, ok := severities[severity]; !ok {
		severity = SEVERITY_WARNING
		s.warn(fmt.Sprintf("unknown severity '%s', using default '%s'", severity, SEVERITY_WARNING))
	}
	r.Severity = severity
	return s
}

// Enable autofix mode.
func (s *State) Autofix() *State {
	s.autofix = true
	return s
}

// Add a string to evaluate the rules against.
func (s *State) Eval(code string) *State {
	s.evals = append(s.evals, code)
	return s
}

// Add a rule or directory of rules to run.
func (s *State) Config(path string) *State {
	s.configs = append(s.configs, path)
	return s
}

// Set the command to invoke Opengrep.
func (s *State) Command(path string) *State {
	s.command = path
	return s
}

// Set an option for the current rule.
func (s *State) Option(name string, value string) *State {
	h := s.headRule()
	h.Options[name] = value
	return s
}

// Only include the specified path pattern in the search.
func (s *State) PathInclude(path string) *State {
	h := s.headRule()
	if h.Paths == nil {
		h.Paths = &RulePaths{}
	}
	h.Paths.Include = append(h.Paths.Include, path)
	return s
}

// Exclude the specified path pattern from the search.
func (s *State) PathExclude(path string) *State {
	h := s.headRule()
	if h.Paths == nil {
		h.Paths = &RulePaths{}
	}
	h.Paths.Exclude = append(h.Paths.Exclude, path)
	return s
}

// Enable Opengrep verbose mode.
func (s *State) Verbose() *State {
	s.verbose = true
	return s
}

func normalizeMetavariable(value string) string {
	if value[0] != '$' {
		return "$" + value
	}

	return value
}
