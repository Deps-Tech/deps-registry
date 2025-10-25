package parser

import (
	"regexp"
	"strings"
)

type WarningType int

const (
	WarningDynamicRequire WarningType = iota
	WarningVariableRequire
	WarningTableRequire
	WarningConcatRequire
)

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityError
)

type Warning struct {
	Type     WarningType
	Line     int
	Module   string
	Severity Severity
	Message  string
}

var dynamicPatterns = []struct {
	pattern  *regexp.Regexp
	warnType WarningType
	severity Severity
	message  string
}{
	{
		pattern:  regexp.MustCompile(`require\s*\(\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\)`),
		warnType: WarningVariableRequire,
		severity: SeverityWarning,
		message:  "Dynamic require detected with variable",
	},
	{
		pattern:  regexp.MustCompile(`require\s*\(\s*.+\[.+\]\s*\)`),
		warnType: WarningTableRequire,
		severity: SeverityWarning,
		message:  "Dynamic require detected with table index",
	},
	{
		pattern:  regexp.MustCompile(`require\s*\(\s*.+\.\..+\s*\)`),
		warnType: WarningConcatRequire,
		severity: SeverityWarning,
		message:  "Dynamic require detected with concatenation",
	},
}

func DetectDynamicRequires(source string) []Warning {
	warnings := []Warning{}
	lines := strings.Split(source, "\n")

	for lineNum, line := range lines {
		for _, pattern := range dynamicPatterns {
			if pattern.pattern.MatchString(line) {
				match := pattern.pattern.FindStringSubmatch(line)
				module := ""
				if len(match) > 1 {
					module = match[1]
				}

				warnings = append(warnings, Warning{
					Type:     pattern.warnType,
					Line:     lineNum + 1,
					Module:   module,
					Severity: pattern.severity,
					Message:  pattern.message,
				})
			}
		}
	}

	return warnings
}
