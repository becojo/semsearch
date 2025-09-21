package rule

import (
	"go.yaml.in/yaml/v2"
)

const (
	SEVERITY_WARNING = "WARNING"
	SEVERITY_ERROR   = "ERROR"
	SEVERITY_INFO    = "INFO"

	MODE_TAINT = "taint"
)

var severities = map[string]bool{
	SEVERITY_WARNING: true,
	SEVERITY_ERROR:   true,
	SEVERITY_INFO:    true,
}

type Rule struct {
	Id        string         `yaml:"id"`
	Severity  string         `yaml:"severity"`
	Message   string         `yaml:"message,omitempty"`
	Languages []string       `yaml:"languages"`
	Fix       string         `yaml:"fix,omitempty"`
	FixRegex  string         `yaml:"fix-regex,omitempty"`
	Options   map[string]any `yaml:"options,omitempty"`
	Metadata  map[string]any `yaml:"metadata,omitempty"`
	Paths     *RulePaths     `yaml:"paths,omitempty"`

	Patterns       *[]Pattern `yaml:"patterns,omitempty"`
	PatternSources *[]Pattern `yaml:"pattern-sources,omitempty"`
	PatternSinks   *[]Pattern `yaml:"pattern-sinks,omitempty"`
}

func (r Rule) MarshalYAML() (any, error) {
	languages := r.Languages
	if len(languages) == 0 {
		languages = []string{"generic"}
	}

	items := yaml.MapSlice{
		yaml.MapItem{Key: "id", Value: r.Id},
		yaml.MapItem{Key: "severity", Value: r.Severity},
		yaml.MapItem{Key: "message", Value: r.Message},
		yaml.MapItem{Key: "languages", Value: languages},
	}

	if r.Paths != nil {
		if len(r.Paths.Exclude) > 0 {
			items = append(items, yaml.MapItem{Key: "paths", Value: r.Paths})
		} else if len(r.Paths.Include) > 0 {
			items = append(items, yaml.MapItem{Key: "paths", Value: r.Paths})
		}
	}

	if len(r.Options) > 0 {
		items = append(items, yaml.MapItem{Key: "options", Value: r.Options})
	}

	if r.PatternSources != nil && len(*r.PatternSources) > 0 {
		items = append(items,
			yaml.MapItem{Key: "mode", Value: MODE_TAINT},
			yaml.MapItem{Key: "pattern-sources", Value: r.PatternSources},
			yaml.MapItem{Key: "pattern-sinks", Value: r.PatternSinks})
	} else if r.Patterns != nil && len(*r.Patterns) > 0 {
		items = append(items, yaml.MapItem{Key: "patterns", Value: r.Patterns})
	}

	if r.FixRegex != "" {
		items = append(items, yaml.MapItem{Key: "fix-regex", Value: []yaml.MapItem{
			{Key: "regex", Value: r.FixRegex},
			{Key: "replacement", Value: r.Fix},
		}})
	} else if r.Fix != "" {
		items = append(items, yaml.MapItem{Key: "fix", Value: r.Fix})
	}

	if len(r.Metadata) > 0 {
		items = append(items, yaml.MapItem{Key: "metadata", Value: r.Metadata})
	}

	return items, nil
}

type Pattern struct {
	Pattern             string               `yaml:"pattern,omitempty"`
	PatternNot          string               `yaml:"pattern-not,omitempty"`
	PatternInside       string               `yaml:"pattern-inside,omitempty"`
	PatternNotInside    string               `yaml:"pattern-not-inside,omitempty"`
	PatternRegex        string               `yaml:"pattern-regex,omitempty"`
	PatternNotRegex     string               `yaml:"pattern-not-regex,omitempty"`
	FocusMetavariable   string               `yaml:"focus-metavariable,omitempty"`
	MetavariableRegex   *MetavariableRegex   `yaml:"metavariable-regex,omitempty"`
	MetavariablePattern *MetavariablePattern `yaml:"metavariable-pattern,omitempty"`

	Patterns      *[]Pattern `yaml:"patterns,omitempty"`
	PatternEither *[]Pattern `yaml:"pattern-either,omitempty"`
}

type MetavariableRegex struct {
	Metavariable string `yaml:"metavariable,omitempty"`
	Regex        string `yaml:"regex,omitempty"`
}

type MetavariablePattern struct {
	Metavariable string     `yaml:"metavariable,omitempty"`
	Patterns     *[]Pattern `yaml:"patterns,omitempty"`
}

type RulePaths struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}
