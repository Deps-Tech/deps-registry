package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func copyFiles(src, dest string) ([]string, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return nil, err
	}

	var copiedFileNames []string

	if !srcInfo.IsDir() {
		// Source is a single file
		fileName := srcInfo.Name()
		destPath := filepath.Join(dest, fileName)
		input, err := ioutil.ReadFile(src)
		if err != nil {
			return nil, fmt.Errorf("could not read source file '%s': %w", src, err)
		}
		if err = ioutil.WriteFile(destPath, input, srcInfo.Mode()); err != nil {
			return nil, fmt.Errorf("could not write destination file '%s': %w", destPath, err)
		}
		copiedFileNames = append(copiedFileNames, fileName)
		return copiedFileNames, nil
	}

	// Source is a directory, walk it
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("could not determine relative path for '%s': %w", path, err)
		}

		if relPath == "." {
			return nil // Skip the root directory itself
		}

		destPath := filepath.Join(dest, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read source file '%s': %w", path, err)
		}

		if err := ioutil.WriteFile(destPath, input, info.Mode()); err != nil {
			return fmt.Errorf("could not write destination file '%s': %w", destPath, err)
		}

		copiedFileNames = append(copiedFileNames, relPath)
		return nil
	})

	return copiedFileNames, err
}

func getLatestVersionInfo(itemPath string) (string, string, error) {
	versions, err := ioutil.ReadDir(itemPath)
	if err != nil {
		return "", "", err
	}

	var versionNames []string
	for _, v := range versions {
		if v.IsDir() {
			versionNames = append(versionNames, v.Name())
		}
	}

	if len(versionNames) == 0 {
		return "", "", fmt.Errorf("no versions found for this item")
	}

	sort.Strings(versionNames)
	latestVersion := versionNames[len(versionNames)-1]

	manifestPath := filepath.Join(itemPath, latestVersion, "dep.json")
	manifestFile, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return "", "", fmt.Errorf("could not read manifest for version %s: %w", latestVersion, err)
	}

	var manifest DepManifest
	err = json.Unmarshal(manifestFile, &manifest)
	if err != nil {
		return "", "", fmt.Errorf("could not parse manifest for version %s: %w", latestVersion, err)
	}

	return latestVersion, manifest.SourceURL, nil
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}
	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func DetectCycle(depID string, allDeps map[string]bool, visited map[string]bool) error {
	if visited[depID] {
		return fmt.Errorf("cyclical dependency detected: %s", depID)
	}
	visited[depID] = true

	depPath := filepath.Join("..", "deps", depID)
	versions, err := ioutil.ReadDir(depPath)
	if err != nil {
		return nil // Should be handled by the main check
	}

	for _, version := range versions {
		if !version.IsDir() {
			continue
		}
		versionStr := version.Name()
		manifestPath := filepath.Join(depPath, versionStr, "dep.json")
		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			continue // Should be handled by the main check
		}

		var manifest DepManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue // Should be handled by the main check
		}

		for subDep := range manifest.Dependencies {
			subDepRoot := getDepRootDir(subDep)
			if allDeps[subDepRoot] {
				if err := DetectCycle(subDepRoot, allDeps, visited); err != nil {
					return fmt.Errorf("%s -> %s", depID, err.Error())
				}
			}
		}
	}
	delete(visited, depID)
	return nil
}
