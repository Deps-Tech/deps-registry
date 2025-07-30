package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate all dependencies and scripts in the registry",
	Run:   runValidation,
}

var path string

func init() {
	validateCmd.Flags().StringVar(&path, "path", "", "Optional: path to a specific item type to validate (e.g., 'deps' or 'scripts')")
	rootCmd.AddCommand(validateCmd)
}

func runValidation(cmd *cobra.Command, args []string) {
	fmt.Println("Running validation...")
	hasErrors := false

	if path != "" {
		if err := validateRegistry(path); err != nil {
			fmt.Printf("Validation failed for %s: %v\n", path, err)
			hasErrors = true
		}
	} else {
		if err := validateRegistry("deps"); err != nil {
			fmt.Printf("Validation failed for dependencies: %v\n", err)
			hasErrors = true
		}
		if err := validateRegistry("scripts"); err != nil {
			fmt.Printf("Validation failed for scripts: %v\n", err)
			hasErrors = true
		}
	}

	if hasErrors {
		fmt.Println("Validation finished with errors.")
		os.Exit(1)
	}

	fmt.Println("Validation successful. All items are consistent.")
}

func validateRegistry(itemType string) error {
	basePath := filepath.Join("..", itemType)
	items, err := ioutil.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Directory not found: %s. Skipping.\n", basePath)
			return nil
		}
		return fmt.Errorf("could not read %s directory: %w", itemType, err)
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		itemID := item.Name()
		itemPath := filepath.Join(basePath, itemID)
		versions, err := ioutil.ReadDir(itemPath)
		if err != nil {
			return fmt.Errorf("could not read versions for %s: %w", itemID, err)
		}

		for _, version := range versions {
			if !version.IsDir() {
				continue
			}
			versionID := version.Name()
			versionPath := filepath.Join(itemPath, versionID)
			if err := validateVersion(versionPath); err != nil {
				return fmt.Errorf("validation failed for %s/%s@%s: %w", itemType, itemID, versionID, err)
			}
		}
	}
	return nil
}

func validateVersion(versionPath string) error {
	manifestPath := filepath.Join(versionPath, "dep.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return fmt.Errorf("manifest 'dep.json' not found")
	}

	manifestData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("could not read manifest: %w", err)
	}

	var manifest DepManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("could not parse manifest: %w", err)
	}

	diskFiles, err := ioutil.ReadDir(versionPath)
	if err != nil {
		return fmt.Errorf("could not read version directory: %w", err)
	}

	diskFileMap := make(map[string]bool)
	for _, f := range diskFiles {
		if f.Name() != "dep.json" {
			diskFileMap[f.Name()] = true
		}
	}

	manifestFileMap := make(map[string]bool)
	for _, f := range manifest.Files {
		manifestFileMap[f] = true
		if _, exists := diskFileMap[f]; !exists {
			return fmt.Errorf("file '%s' listed in manifest but not found on disk", f)
		}
	}

	for f := range diskFileMap {
		if _, exists := manifestFileMap[f]; !exists {
			return fmt.Errorf("file '%s' found on disk but not listed in manifest", f)
		}
	}

	return nil
}
