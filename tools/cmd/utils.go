package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func copyFiles(src, dest string) ([]string, error) {
	var fileNames []string
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		sourcePath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		input, err := ioutil.ReadFile(sourcePath)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(destPath, input, 0644)
		if err != nil {
			return nil, err
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
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
