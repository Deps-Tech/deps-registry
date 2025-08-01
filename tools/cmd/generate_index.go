package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cdnURL string

var generateIndexCmd = &cobra.Command{
	Use:   "generate-index",
	Short: "Generate an index file for all packages",
	Run:   runGenerateIndex,
}

func init() {
	generateIndexCmd.Flags().StringVar(&cdnURL, "cdn-url", "", "The base URL of the CDN")
	generateIndexCmd.MarkFlagRequired("cdn-url")
	rootCmd.AddCommand(generateIndexCmd)
}

type VersionInfo struct {
	URL      string      `json:"url"`
	Sha256   string      `json:"sha256"`
	Manifest DepManifest `json:"manifest"`
}

type DependencyInfo struct {
	Latest   string                 `json:"latest"`
	Versions map[string]VersionInfo `json:"versions"`
}

type Index struct {
	LastUpdated  string                    `json:"last_updated"`
	Dependencies map[string]DependencyInfo `json:"dependencies"`
	Scripts      map[string]DependencyInfo `json:"scripts"`
}

func runGenerateIndex(cmd *cobra.Command, args []string) {
	distPath := "dist"

	deps, err := generateIndexForType("deps", distPath, cdnURL)
	if err != nil {
		fmt.Printf("Error generating index for dependencies: %v\n", err)
		os.Exit(1)
	}

	scripts, err := generateIndexForType("scripts", distPath, cdnURL)
	if err != nil {
		fmt.Printf("Error generating index for scripts: %v\n", err)
		os.Exit(1)
	}

	index := Index{
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
		Dependencies: deps,
		Scripts:      scripts,
	}

	indexPath := filepath.Join(distPath, "index.json")
	indexData, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling index.json: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		fmt.Printf("Error writing index.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated index.json")
}

func generateIndexForType(itemType, distPath, cdnURL string) (map[string]DependencyInfo, error) {
	itemsPath := filepath.Join(distPath, itemType)
	result := make(map[string]DependencyInfo)

	if _, err := os.Stat(itemsPath); os.IsNotExist(err) {
		return result, nil
	}

	files, err := os.ReadDir(itemsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s directory: %w", itemsPath, err)
	}

	itemVersions := make(map[string][]string)
	re := regexp.MustCompile(`^(.+)-(\d+\.\d+\.\d+)\.zip$`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matches := re.FindStringSubmatch(file.Name())
		if len(matches) == 3 {
			itemName := matches[1]
			version := matches[2]
			itemVersions[itemName] = append(itemVersions[itemName], version)
		}
	}

	for itemName, versions := range itemVersions {
		sort.Sort(sort.Reverse(sort.StringSlice(versions)))
		latestVersion := versions[0]

		itemInfo := DependencyInfo{
			Latest:   latestVersion,
			Versions: make(map[string]VersionInfo),
		}

		for _, version := range versions {
			fileName := fmt.Sprintf("%s-%s.zip", itemName, version)
			filePath := filepath.Join(itemsPath, fileName)
			sha256, err := hashFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("error hashing file %s: %w", fileName, err)
			}

			manifest, err := readManifestFromZip(filePath)
			if err != nil {
				return nil, fmt.Errorf("error reading manifest from %s: %w", fileName, err)
			}

			url := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(cdnURL, "/"), itemType, fileName)
			itemInfo.Versions[version] = VersionInfo{
				URL:      url,
				Sha256:   sha256,
				Manifest: *manifest,
			}
		}
		result[itemName] = itemInfo
	}

	return result, nil
}

func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func readManifestFromZip(zipPath string) (*DepManifest, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "dep.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var manifest DepManifest
			if err := json.Unmarshal(content, &manifest); err != nil {
				return nil, err
			}
			return &manifest, nil
		}
	}

	return nil, fmt.Errorf("dep.json not found in %s", zipPath)
}
