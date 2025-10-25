package parser

import (
	"regexp"
	"strings"
)

var (
	reRequireStatic = regexp.MustCompile(`require\s*\(?\s*["']([\w\.-]+)["']\s*\)?`)
	reNetworkAccess = regexp.MustCompile(`(http\.request|socket\.tcp|socket\.connect)`)
	reFFI           = regexp.MustCompile(`\brequire\s*\(?\s*["']ffi["']\s*\)?`)
	reFileAccess    = regexp.MustCompile(`io\.open\s*\(\s*["']([^"']+)["']`)
)

type RegexResult struct {
	RawModules  []string
	UsesNetwork bool
	UsesFFI     bool
	FilePaths   []string
}

func ParseWithRegex(source string) *RegexResult {
	result := &RegexResult{
		RawModules: []string{},
		FilePaths:  []string{},
	}

	matches := reRequireStatic.FindAllStringSubmatch(source, -1)
	for _, match := range matches {
		if len(match) > 1 {
			module := strings.TrimPrefix(match[1], "lib.")
			result.RawModules = append(result.RawModules, module)
		}
	}

	result.UsesNetwork = reNetworkAccess.MatchString(source)
	result.UsesFFI = reFFI.MatchString(source)

	fileMatches := reFileAccess.FindAllStringSubmatch(source, -1)
	for _, match := range fileMatches {
		if len(match) > 1 {
			result.FilePaths = append(result.FilePaths, match[1])
		}
	}

	return result
}
