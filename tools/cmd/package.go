package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package all dependencies and scripts for distribution",
	Run:   runPackaging,
}

func init() {
	rootCmd.AddCommand(packageCmd)
}

func runPackaging(cmd *cobra.Command, args []string) {
	fmt.Println("Running packaging...")
	distPath := filepath.Join("..", "dist")
	if err := os.MkdirAll(distPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating dist directory: %v\n", err)
		os.Exit(1)
	}

	if err := packageItems("deps", distPath); err != nil {
		fmt.Printf("Packaging failed for dependencies: %v\n", err)
		os.Exit(1)
	}

	if err := packageItems("scripts", distPath); err != nil {
		fmt.Printf("Packaging failed for scripts: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Packaging successful.")
}

func packageItems(itemType, distPath string) error {
	basePath := filepath.Join("..", itemType)
	items, err := ioutil.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nothing to package
		}
		return err
	}

	index := make(map[string][]string)

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		itemID := item.Name()
		itemPath := filepath.Join(basePath, itemID)
		versions, err := ioutil.ReadDir(itemPath)
		if err != nil {
			return err
		}

		var versionNames []string
		for _, version := range versions {
			if !version.IsDir() {
				continue
			}
			versionID := version.Name()
			versionNames = append(versionNames, versionID)
			versionPath := filepath.Join(itemPath, versionID)
			zipName := fmt.Sprintf("%s-%s.zip", itemID, versionID)
			zipPath := filepath.Join(distPath, itemType, zipName)

			if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
				return err
			}
			if err := zipDirectory(versionPath, zipPath); err != nil {
				return fmt.Errorf("failed to zip %s: %w", versionPath, err)
			}
		}
		index[itemID] = versionNames
	}

	indexPath := filepath.Join(distPath, fmt.Sprintf("%s_index.json", itemType))
	indexData, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(indexPath, indexData, 0644)
}

func zipDirectory(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
