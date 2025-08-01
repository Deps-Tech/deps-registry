package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ScriptAnalysis struct {
	Dependencies     []string
	TouchedFiles     []string
	HasNetworkAccess bool
}

func analyzeScript(sourcePath string, id string, allDepIDs map[string]bool) (*ScriptAnalysis, error) {
	analysis := &ScriptAnalysis{
		Dependencies:     []string{},
		TouchedFiles:     []string{},
		HasNetworkAccess: false,
	}

	depSet := make(map[string]struct{})
	touchedFilesSet := make(map[string]struct{})

	reDep := regexp.MustCompile(`require\s*\(?\s*["']([\w\.-]+)["']\s*\)?`)
	reFileVar := regexp.MustCompile(`local\s+([\w_]+)\s*=\s*getWorkingDirectory\(\)\s*\.\.\s*["']([^"']+)["']`)
	reFileConcat := regexp.MustCompile(`local\s+([\w_]+)\s*=\s*([\w_]+)\s*\.\.\s*["']([^"']+)["']`)
	reFileUsage := regexp.MustCompile(`(io\.open|jsonSave|jsonRead|createDirectory|doesDirectoryExist)\s*\(\s*([\w_]+)`)
	reNet := regexp.MustCompile(`(http\.request|socket\.tcp)`)

	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, err
	}

	var luaFiles []string
	if info.IsDir() {
		err = filepath.Walk(sourcePath, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".lua" {
				luaFiles = append(luaFiles, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else if filepath.Ext(sourcePath) == ".lua" {
		luaFiles = append(luaFiles, sourcePath)
	}

	for _, filePath := range luaFiles {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		contentStr := string(content)
		pathVars := make(map[string]string)

		// Pass 1: Find base paths from getWorkingDirectory()
		matchesFileVar := reFileVar.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matchesFileVar {
			if len(match) == 3 {
				cleanPath := strings.ReplaceAll(match[2], `\\`, `\`)
				pathVars[match[1]] = "<working_dir>" + cleanPath
			}
		}

		// Pass 2: Find concatenated paths
		matchesFileConcat := reFileConcat.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matchesFileConcat {
			if len(match) == 4 {
				if basePath, ok := pathVars[match[2]]; ok {
					cleanPath := strings.ReplaceAll(match[3], `\\`, `\`)
					pathVars[match[1]] = basePath + cleanPath
				}
			}
		}

		// Pass 3: Find usage of path variables
		matchesFileUsage := reFileUsage.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matchesFileUsage {
			if len(match) == 3 {
				if resolvedPath, ok := pathVars[match[2]]; ok {
					touchedFilesSet[resolvedPath] = struct{}{}
				}
			}
		}

		// Dependencies
		matchesDep := reDep.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matchesDep {
			if len(match) > 1 {
				cleanDepID := strings.TrimPrefix(match[1], "lib.")
				depRootDir := getDepRootDir(cleanDepID)
				if id != "" && strings.EqualFold(depRootDir, id) {
					continue
				}

				if _, exists := allDepIDs[strings.ToLower(depRootDir)]; !exists {
					continue
				}

				depSet[cleanDepID] = struct{}{}
			}
		}

		// Network access
		if reNet.MatchString(contentStr) {
			analysis.HasNetworkAccess = true
		}
	}

	for dep := range depSet {
		analysis.Dependencies = append(analysis.Dependencies, dep)
	}
	for file := range touchedFilesSet {
		analysis.TouchedFiles = append(analysis.TouchedFiles, file)
	}

	// Filter out redundant parent directories
	paths := analysis.TouchedFiles
	toRemove := make(map[string]bool)
	for _, p1 := range paths {
		for _, p2 := range paths {
			if p1 != p2 && strings.HasPrefix(p2, p1+string(os.PathSeparator)) {
				toRemove[p1] = true
			}
		}
	}

	finalPaths := []string{}
	for _, p := range paths {
		if !toRemove[p] {
			finalPaths = append(finalPaths, p)
		}
	}
	analysis.TouchedFiles = finalPaths

	return analysis, nil
}

func getAvailableVersions(depID string) ([]string, error) {
	cleanDepID := strings.TrimPrefix(depID, "lib.")
	depPath := filepath.Join("..", "deps", getDepRootDir(cleanDepID))
	versions, err := ioutil.ReadDir(depPath)
	if err != nil {
		return nil, err
	}

	var versionNames []string
	for _, v := range versions {
		if v.IsDir() {
			versionNames = append(versionNames, v.Name())
		}
	}
	return versionNames, nil
}

func getDepRootDir(depID string) string {
	return strings.Split(depID, ".")[0]
}

func resolveDependencies(deps map[string]string) (map[string]string, error) {
	resolved := make(map[string]string)
	for dep, version := range deps {
		if err := resolve(dep, version, resolved); err != nil {
			return nil, err
		}
	}
	return resolved, nil
}

func resolve(depID, version string, resolved map[string]string) error {
	if _, ok := resolved[depID]; ok {
		return nil
	}

	resolved[depID] = version

	cleanDepID := strings.TrimPrefix(depID, "lib.")
	manifestPath := filepath.Join("..", "deps", getDepRootDir(cleanDepID), version, "dep.json")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		if strings.Contains(cleanDepID, ".") {
			return nil
		}
		return err
	}

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	var manifest DepManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	for dep, ver := range manifest.Dependencies {
		if err := resolve(dep, ver, resolved); err != nil {
			return err
		}
	}

	return nil
}
