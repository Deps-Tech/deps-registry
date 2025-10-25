package parser

import (
	"os"
	"path/filepath"
)

type Analysis struct {
	Dependencies []string
	FilePaths    []string
	UsesNetwork  bool
	UsesFFI      bool
	Warnings     []Warning
	HasDynamic   bool
}

func AnalyzeLua(sourcePath string, excludeID string, availableDeps map[string]bool) (*Analysis, error) {
	registry := NewRegistry()
	for depID := range availableDeps {
		registry.AddPackage(&PackageInfo{
			ID: depID,
		})
	}

	ctx, err := NewContext(excludeID, sourcePath, registry)
	if err != nil {
		return nil, err
	}

	return AnalyzeWithContext(ctx, sourcePath)
}

func AnalyzeWithContext(ctx *Context, sourcePath string) (*Analysis, error) {
	analysis := &Analysis{
		Dependencies: []string{},
		FilePaths:    []string{},
		Warnings:     []Warning{},
	}

	luaFiles, err := findLuaFiles(sourcePath)
	if err != nil {
		return nil, err
	}

	allRawModules := []string{}

	for _, file := range luaFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		contentStr := string(content)

		regexResult := ParseWithRegex(contentStr)
		allRawModules = append(allRawModules, regexResult.RawModules...)
		analysis.FilePaths = append(analysis.FilePaths, regexResult.FilePaths...)

		if regexResult.UsesNetwork {
			analysis.UsesNetwork = true
		}
		if regexResult.UsesFFI {
			analysis.UsesFFI = true
		}

		warnings := DetectDynamicRequires(contentStr)
		analysis.Warnings = append(analysis.Warnings, warnings...)
		if len(warnings) > 0 {
			analysis.HasDynamic = true
		}
	}

	resolved := ResolveDependencies(ctx, allRawModules)

	for pkgID := range resolved {
		analysis.Dependencies = append(analysis.Dependencies, pkgID)
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
