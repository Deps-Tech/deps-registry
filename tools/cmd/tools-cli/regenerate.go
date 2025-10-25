package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/filesystem"
	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
	"github.com/Deps-Tech/deps-registry/tools/internal/parser"
	"github.com/Deps-Tech/deps-registry/tools/internal/registry"
	"github.com/Deps-Tech/deps-registry/tools/internal/validator"
	"github.com/Deps-Tech/deps-registry/tools/internal/versioning"
	"github.com/spf13/cobra"
)

var (
	dryRun      bool
	skipValidation bool
)

var regenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate all manifests with improved dependency parser",
	Long:  "Scans all existing packages, re-analyzes dependencies with the new parser, and updates dep.json files",
	Run:   runRegenerate,
}

func init() {
	regenerateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be changed without writing")
	regenerateCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip validation checks")
	rootCmd.AddCommand(regenerateCmd)
}

func runRegenerate(cmd *cobra.Command, args []string) {
	fmt.Println("Starting manifest regeneration...")

	basePaths := []string{
		filepath.Join("..", "deps"),
		filepath.Join("..", "scripts"),
	}

	allManifests := make(map[string]*manifest.Manifest)
	itemPaths := make(map[string]string)

	for _, basePath := range basePaths {
		items, err := os.ReadDir(basePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("Error reading %s: %v\n", basePath, err)
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
				m, err := manifest.Load(versionPath)
				if err != nil {
					fmt.Printf("⚠️  %s/%s: failed to load manifest: %v\n", item.Name(), version.Name(), err)
					continue
				}

				allManifests[m.ID] = m
				itemPaths[m.ID] = versionPath
			}
		}
	}

	fmt.Printf("\nLoaded %d packages\n", len(allManifests))

	if !skipValidation {
		fmt.Println("\nChecking for issues...")

		cycles := validator.DetectCycles(allManifests)
		if len(cycles) > 0 {
			fmt.Printf("\n⚠️  Found %d circular dependencies:\n", len(cycles))
			for _, cycle := range cycles {
				fmt.Printf("   %v\n", cycle.Cycle)
			}
		}

		duplicates := validator.DetectDuplicates(allManifests)
		if len(duplicates) > 0 {
			fmt.Printf("\n⚠️  Found %d sets of duplicate packages:\n", len(duplicates))
			for _, dup := range duplicates {
				fmt.Printf("   %v\n", dup.Packages)
			}
		}
	}

	reg, err := parser.LoadRegistryFromManifests(basePaths)
	if err != nil {
		fmt.Printf("Error loading registry: %v\n", err)
		os.Exit(1)
	}

	for pkgID, aliases := range registry.WellKnownAliases {
		if pkg := reg.GetPackage(pkgID); pkg != nil {
			pkg.Provides = aliases
		}
	}

	updated := 0
	skipped := 0
	errors := 0

	for id, m := range allManifests {
		versionPath := itemPaths[id]

		ctx, err := parser.NewContext(id, versionPath, reg)
		if err != nil {
			fmt.Printf("❌ %s: failed to create context: %v\n", id, err)
			errors++
			continue
		}

		analysis, err := parser.AnalyzeWithContext(ctx, versionPath)
		if err != nil {
			fmt.Printf("❌ %s: analysis failed: %v\n", id, err)
			errors++
			continue
		}

		newDeps := make(map[string]string)
		for _, dep := range analysis.Dependencies {
			version := getLatestVersionForDep(dep, basePaths)
			if version != "" {
				newDeps[dep] = version
			}
		}

		providesAliases := registry.GetAliases(id)

		changed := false
		if !mapsEqual(m.Dependencies, newDeps) {
			changed = true
		}
		if !slicesEqual(m.Provides, providesAliases) {
			changed = true
		}

		if changed {
			m.Dependencies = newDeps
			m.Provides = providesAliases

			fileMap := make(map[string]manifest.FileInfo)
			for fileName := range m.Files {
				filePath := filepath.Join(versionPath, fileName)
				hash, _ := filesystem.SHA256File(filePath)
				info, _ := os.Stat(filePath)
				fileMap[fileName] = manifest.FileInfo{
					SHA256: hash,
					Size:   info.Size(),
				}
			}
			m.Files = fileMap

			if !dryRun {
				if err := manifest.Save(versionPath, m); err != nil {
					fmt.Printf("❌ %s: failed to save: %v\n", id, err)
					errors++
					continue
				}
			}

			fmt.Printf("✓ %s\n", id)
			if len(analysis.Warnings) > 0 {
				fmt.Printf("  ⚠️  %d warnings (dynamic requires detected)\n", len(analysis.Warnings))
			}
			updated++
		} else {
			skipped++
		}
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Updated: %d\n", updated)
	fmt.Printf("Skipped: %d\n", skipped)
	if errors > 0 {
		fmt.Printf("Errors:  %d\n", errors)
	}
	if dryRun {
		fmt.Println("\n(Dry run - no files were modified)")
	}
}

func getLatestVersionForDep(depID string, basePaths []string) string {
	for _, basePath := range basePaths {
		depPath := filepath.Join(basePath, depID)
		versions, err := os.ReadDir(depPath)
		if err != nil {
			continue
		}

		versionList := []string{}
		for _, v := range versions {
			if v.IsDir() {
				versionList = append(versionList, v.Name())
			}
		}

		if len(versionList) > 0 {
			return versioning.GetLatest(versionList)
		}
	}
	return ""
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, v := range a {
		aMap[v] = true
	}
	for _, v := range b {
		if !aMap[v] {
			return false
		}
	}
	return true
}

