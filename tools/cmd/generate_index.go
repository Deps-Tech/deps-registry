package cmd

import (
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
	URL    string `json:"url"`
	Sha256 string `json:"sha256"`
}

type DependencyInfo struct {
	Latest   string                 `json:"latest"`
	Versions map[string]VersionInfo `json:"versions"`
}

type Index struct {
	LastUpdated  string                    `json:"last_updated"`
	Dependencies map[string]DependencyInfo `json:"dependencies"`
}

func runGenerateIndex(cmd *cobra.Command, args []string) {
	distPath := "dist"
	depsPath := filepath.Join(distPath, "deps")

	if _, err := os.Stat(depsPath); os.IsNotExist(err) {
		fmt.Println("No deps directory found, skipping index generation.")
		return
	}

	index := Index{
		LastUpdated:  time.Now().UTC().Format(time.RFC3339),
		Dependencies: make(map[string]DependencyInfo),
	}

	files, err := os.ReadDir(depsPath)
	if err != nil {
		fmt.Printf("Error reading deps directory: %v\n", err)
		os.Exit(1)
	}

	depVersions := make(map[string][]string)
	re := regexp.MustCompile(`^(.+)-(\d+\.\d+\.\d+)\.zip$`)

	for _, file := range files {
		if !file.IsDir() {
			matches := re.FindStringSubmatch(file.Name())
			if len(matches) == 3 {
				depName := matches[1]
				version := matches[2]
				depVersions[depName] = append(depVersions[depName], version)
			}
		}
	}

	for depName, versions := range depVersions {
		sort.Sort(sort.Reverse(sort.StringSlice(versions)))
		latestVersion := versions[0]

		depInfo := DependencyInfo{
			Latest:   latestVersion,
			Versions: make(map[string]VersionInfo),
		}

		for _, version := range versions {
			fileName := fmt.Sprintf("%s-%s.zip", depName, version)
			filePath := filepath.Join(depsPath, fileName)
			sha256, err := hashFile(filePath)
			if err != nil {
				fmt.Printf("Error hashing file %s: %v\n", fileName, err)
				os.Exit(1)
			}
			url := fmt.Sprintf("%s/deps/%s", strings.TrimSuffix(cdnURL, "/"), fileName)
			depInfo.Versions[version] = VersionInfo{
				URL:    url,
				Sha256: sha256,
			}
		}
		index.Dependencies[depName] = depInfo
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
