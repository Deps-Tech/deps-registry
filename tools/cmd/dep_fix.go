package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var depFixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix dependencies for all existing dependencies",
	Run:   fixDependencies,
}

func init() {
	depCmd.AddCommand(depFixCmd)
}

func fixDependencies(cmd *cobra.Command, args []string) {
	depsPath := filepath.Join("..", "deps")
	dirs, err := ioutil.ReadDir(depsPath)
	if err != nil {
		fmt.Printf("Error reading deps directory: %v\n", err)
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		depID := dir.Name()
		depPath := filepath.Join(depsPath, depID)
		versions, err := ioutil.ReadDir(depPath)
		if err != nil {
			fmt.Printf("Error reading versions for dependency '%s': %v\n", depID, err)
			continue
		}

		for _, version := range versions {
			if !version.IsDir() {
				continue
			}
			versionPath := filepath.Join(depPath, version.Name())
			manifestPath := filepath.Join(versionPath, "dep.json")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				continue
			}

			analysis, err := analyzeScript(versionPath)
			if err != nil {
				fmt.Printf("Error analyzing dependency '%s' version '%s': %v\n", depID, version.Name(), err)
				continue
			}

			if len(analysis.Dependencies) > 0 {
				data, err := ioutil.ReadFile(manifestPath)
				if err != nil {
					fmt.Printf("Error reading manifest for '%s' version '%s': %v\n", depID, version.Name(), err)
					continue
				}

				var manifest DepManifest
				if err := json.Unmarshal(data, &manifest); err != nil {
					fmt.Printf("Error parsing manifest for '%s' version '%s': %v\n", depID, version.Name(), err)
					continue
				}

				if manifest.Dependencies == nil {
					manifest.Dependencies = make(map[string]string)
				}
				for _, dep := range analysis.Dependencies {
					if _, exists := manifest.Dependencies[dep]; !exists {
						versions, err := getAvailableVersions(dep)
						if err != nil || len(versions) == 0 {
							fmt.Printf("Could not find versions for sub-dependency '%s' of '%s'. Skipping.\n", dep, depID)
							continue
						}
						manifest.Dependencies[dep] = versions[len(versions)-1] // Select latest
					}
				}

				manifestData, err := json.MarshalIndent(manifest, "", "  ")
				if err != nil {
					fmt.Printf("Error creating manifest for '%s' version '%s': %v\n", depID, version.Name(), err)
					continue
				}

				err = ioutil.WriteFile(manifestPath, manifestData, 0644)
				if err != nil {
					fmt.Printf("Error writing manifest for '%s' version '%s': %v\n", depID, version.Name(), err)
					continue
				}
				fmt.Printf("Updated manifest for '%s' version '%s'\n", depID, version.Name())
			}
		}
	}
	fmt.Println("Dependency fix process complete.")
}
