package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func parseLuaDependencies(sourcePath string) ([]string, error) {
	depSet := make(map[string]struct{})
	re := regexp.MustCompile(`require\s*\(?\s*["']([^"']+)["']\s*\)?`)

	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, err
	}

	var luaFiles []string
	if info.IsDir() {
		files, err := ioutil.ReadDir(sourcePath)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".lua" {
				luaFiles = append(luaFiles, filepath.Join(sourcePath, file.Name()))
			}
		}
	} else if filepath.Ext(sourcePath) == ".lua" {
		luaFiles = append(luaFiles, sourcePath)
	}

	for _, filePath := range luaFiles {
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			if len(match) > 1 {
				depSet[match[1]] = struct{}{}
			}
		}
	}

	var deps []string
	for dep := range depSet {
		deps = append(deps, dep)
	}

	return deps, nil
}

func getAvailableVersions(depID string) ([]string, error) {
	depPath := filepath.Join("..", "deps", depID)
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

	manifestPath := filepath.Join("..", "deps", depID, version, "dep.json")
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
