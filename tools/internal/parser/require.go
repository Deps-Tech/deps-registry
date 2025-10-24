package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Analysis struct {
	Dependencies []string
	FilePaths    []string
	UsesNetwork  bool
	UsesFFI      bool
}

var (
	reRequire = regexp.MustCompile(`require\s*\(?\s*["']([\w\.-]+)["']\s*\)?`)
	reNetwork = regexp.MustCompile(`(http\.request|socket\.tcp|socket\.connect)`)
	reFFI     = regexp.MustCompile(`\brequire\s*\(?\s*["']ffi["']\s*\)?`)
)

func AnalyzeLua(sourcePath string, excludeID string, availableDeps map[string]bool) (*Analysis, error) {
	analysis := &Analysis{
		Dependencies: []string{},
		FilePaths:    []string{},
	}

	luaFiles, err := findLuaFiles(sourcePath)
	if err != nil {
		return nil, err
	}

	depSet := make(map[string]struct{})

	for _, file := range luaFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		contentStr := string(content)

		matches := reRequire.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				dep := strings.TrimPrefix(match[1], "lib.")
				rootDep := strings.Split(dep, ".")[0]

				if excludeID != "" && strings.EqualFold(rootDep, excludeID) {
					continue
				}

				if availableDeps[strings.ToLower(rootDep)] {
					depSet[dep] = struct{}{}
				}
			}
		}

		if reNetwork.MatchString(contentStr) {
			analysis.UsesNetwork = true
		}

		if reFFI.MatchString(contentStr) {
			analysis.UsesFFI = true
		}
	}

	for dep := range depSet {
		analysis.Dependencies = append(analysis.Dependencies, dep)
	}

	return analysis, nil
}

func findLuaFiles(path string) ([]string, error) {
	var files []string

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if filepath.Ext(path) == ".lua" {
			return []string{path}, nil
		}
		return []string{}, nil
	}

	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(p) == ".lua" {
			files = append(files, p)
		}
		return nil
	})

	return files, err
}
