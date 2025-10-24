package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/packager"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package all dependencies and scripts into ZIP files",
	Run:   runPackage,
}

func init() {
	rootCmd.AddCommand(packageCmd)
}

func runPackage(cmd *cobra.Command, args []string) {
	distPath := "dist"
	if err := os.MkdirAll(distPath, 0755); err != nil {
		fmt.Printf("Failed to create dist directory: %v\n", err)
		os.Exit(1)
	}

	for _, itemType := range []string{"deps", "scripts"} {
		if err := packageItems(itemType, distPath); err != nil {
			fmt.Printf("Failed to package %s: %v\n", itemType, err)
			os.Exit(1)
		}
	}

	fmt.Println("Packaging complete")
}

func packageItems(itemType, distPath string) error {
	basePath := filepath.Join("..", itemType)
	targetPath := filepath.Join(distPath, itemType)

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	items, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		itemPath := filepath.Join(basePath, item.Name())
		versions, err := os.ReadDir(itemPath)
		if err != nil {
			continue
		}

		for _, version := range versions {
			if !version.IsDir() {
				continue
			}

			versionPath := filepath.Join(itemPath, version.Name())
			zipName := fmt.Sprintf("%s-%s.zip", item.Name(), version.Name())
			zipPath := filepath.Join(targetPath, zipName)

			if err := packager.ZipDirectory(versionPath, zipPath); err != nil {
				return fmt.Errorf("failed to zip %s: %w", versionPath, err)
			}

			fmt.Printf("Packaged %s/%s\n", item.Name(), version.Name())
		}
	}

	return nil
}
