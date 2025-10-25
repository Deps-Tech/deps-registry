package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
	"github.com/Deps-Tech/deps-registry/tools/internal/validator"
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
	allManifests := make(map[string]*manifest.Manifest)

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
				m, err := validateManifest(versionPath)
				if err != nil {
					fmt.Printf("❌ %s/%s: %v\n", item.Name(), version.Name(), err)
					hasErrors = true
				} else {
					fmt.Printf("✓ %s/%s\n", item.Name(), version.Name())
					if m != nil {
						allManifests[m.ID] = m
					}
				}
			}
		}
	}

	if !hasErrors {
		fmt.Println("\nRunning dependency graph validation...")

		cycles := validator.DetectCycles(allManifests)
		if len(cycles) > 0 {
			fmt.Printf("\n⚠️  Found %d circular dependencies:\n", len(cycles))
			for _, cycle := range cycles {
				fmt.Printf("   %v\n", cycle.Cycle)
			}
			hasErrors = true
		}

		duplicates := validator.DetectDuplicates(allManifests)
		if len(duplicates) > 0 {
			fmt.Printf("\n⚠️  Found %d sets of duplicate packages:\n", len(duplicates))
			for _, dup := range duplicates {
				fmt.Printf("   %v\n", dup.Packages)
			}
		}
	}

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("\nAll manifests are valid")
}

func validateManifest(path string) (*manifest.Manifest, error) {
	m, err := manifest.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	if m.ID == "" {
		return nil, fmt.Errorf("missing id")
	}

	if m.Version == "" {
		return nil, fmt.Errorf("missing version")
	}

	if len(m.Files) == 0 {
		return nil, fmt.Errorf("no files listed")
	}

	for file := range m.Files {
		filePath := filepath.Join(path, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("file %s not found", file)
		}
	}

	diskFiles, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range diskFiles {
		if f.IsDir() || f.Name() == "dep.json" {
			continue
		}

		if _, ok := m.Files[f.Name()]; !ok {
			return nil, fmt.Errorf("file %s not in manifest", f.Name())
		}
	}

	return m, nil
}
