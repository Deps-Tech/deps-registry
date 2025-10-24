package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all manifests in the registry",
	Run:   runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	hasErrors := false

	for _, itemType := range []string{"deps", "scripts"} {
		basePath := filepath.Join("..", itemType)
		items, err := os.ReadDir(basePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("Error reading %s: %v\n", itemType, err)
			hasErrors = true
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
				if err := validateManifest(versionPath); err != nil {
					fmt.Printf("❌ %s/%s: %v\n", item.Name(), version.Name(), err)
					hasErrors = true
				} else {
					fmt.Printf("✓ %s/%s\n", item.Name(), version.Name())
				}
			}
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("\nAll manifests are valid")
}

func validateManifest(path string) error {
	m, err := manifest.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load manifest: %w", err)
	}

	if m.ID == "" {
		return fmt.Errorf("missing id")
	}

	if m.Version == "" {
		return fmt.Errorf("missing version")
	}

	if len(m.Files) == 0 {
		return fmt.Errorf("no files listed")
	}

	for file := range m.Files {
		filePath := filepath.Join(path, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file %s not found", file)
		}
	}

	diskFiles, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range diskFiles {
		if f.IsDir() || f.Name() == "dep.json" {
			continue
		}

		if _, ok := m.Files[f.Name()]; !ok {
			return fmt.Errorf("file %s not in manifest", f.Name())
		}
	}

	return nil
}
