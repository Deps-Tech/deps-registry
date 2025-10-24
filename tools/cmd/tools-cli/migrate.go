package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/filesystem"
	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate all dep.json to new format with SHA256 hashes",
	Run:   runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) {
	for _, itemType := range []string{"deps", "scripts"} {
		basePath := filepath.Join("..", itemType)
		items, err := os.ReadDir(basePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("Error reading %s: %v\n", itemType, err)
			continue
		}

		for _, item := range items {
			if !item.IsDir() {
				continue
			}

			itemPath := filepath.Join(basePath, item.Name())
			versions, _ := os.ReadDir(itemPath)

			for _, version := range versions {
				if !version.IsDir() {
					continue
				}

				versionPath := filepath.Join(itemPath, version.Name())
				if err := migrateManifest(versionPath); err != nil {
					fmt.Printf("Error migrating %s/%s: %v\n", item.Name(), version.Name(), err)
				} else {
					fmt.Printf("Migrated %s/%s\n", item.Name(), version.Name())
				}
			}
		}
	}
}

func migrateManifest(path string) error {
	depPath := filepath.Join(path, "dep.json")
	data, err := os.ReadFile(depPath)
	if err != nil {
		return err
	}

	var old map[string]interface{}
	if err := json.Unmarshal(data, &old); err != nil {
		return err
	}

	if _, ok := old["manifestVersion"]; ok {
		return nil
	}

	files, err := filesystem.ListFiles(path)
	if err != nil {
		return err
	}

	fileMap := make(map[string]manifest.FileInfo)
	for _, file := range files {
		if file == "dep.json" {
			continue
		}

		filePath := filepath.Join(path, file)
		hash, err := filesystem.SHA256File(filePath)
		if err != nil {
			continue
		}

		info, _ := os.Stat(filePath)
		fileMap[file] = manifest.FileInfo{
			SHA256: hash,
			Size:   info.Size(),
		}
	}

	m := &manifest.Manifest{
		ManifestVersion: "1.0",
		ID:              getString(old, "id"),
		Name:            getString(old, "name"),
		Version:         getString(old, "version"),
		Files:           fileMap,
		Dependencies:    getStringMap(old, "dependencies"),
		Security: manifest.Security{
			NetworkAccess: getBool(old, "hasNetworkAccess"),
			FileAccess:    getStringSlice(old, "touchedFiles"),
		},
		Metadata: manifest.Metadata{
			SourceURL: getString(old, "sourceUrl"),
			Tags:      getStringSlice(old, "tags"),
		},
	}

	return manifest.Save(path, m)
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key].([]interface{}); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

func getStringMap(m map[string]interface{}, key string) map[string]string {
	if v, ok := m[key].(map[string]interface{}); ok {
		result := make(map[string]string)
		for k, val := range v {
			if s, ok := val.(string); ok {
				result[k] = s
			}
		}
		return result
	}
	return nil
}
