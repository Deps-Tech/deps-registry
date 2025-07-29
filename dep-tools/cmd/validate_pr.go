package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	repoPath string
	prNumber int
)

var validatePrCmd = &cobra.Command{
	Use:   "validate-pr",
	Short: "Validate the structure of a pull request",
	Run:   validatePr,
}

func init() {
	rootCmd.AddCommand(validatePrCmd)
	validatePrCmd.Flags().StringVar(&repoPath, "path", "", "Repository path (e.g., 'owner/repo')")
	validatePrCmd.Flags().IntVar(&prNumber, "pr", 0, "Pull request number")
	validatePrCmd.MarkFlagRequired("path")
	validatePrCmd.MarkFlagRequired("pr")
}

func validatePr(cmd *cobra.Command, args []string) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Error: GITHUB_TOKEN environment variable not set.")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner, repo, err := parseRepoPath(repoPath)
	if err != nil {
		fmt.Printf("Error: Invalid repository path: %v\n", err)
		os.Exit(1)
	}

	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		fmt.Printf("Error fetching PR files: %v\n", err)
		os.Exit(1)
	}

	validatedManifests := make(map[string]bool)
	valid := true

	for _, file := range files {
		filename := file.GetFilename()
		fmt.Printf("Validating file: %s\n", filename)

		if !isValidFilePath(filename) {
			valid = false
			continue
		}

		if strings.HasPrefix(filename, "deps/") {
			dir := filepath.Dir(filename)
			if _, ok := validatedManifests[dir]; ok {
				continue
			}

			if !validateManifest(dir) {
				valid = false
			}
			validatedManifests[dir] = true
		}
	}

	if !valid {
		fmt.Println("PR validation failed.")
		os.Exit(1)
	}

	fmt.Println("PR validation successful.")
}

func parseRepoPath(path string) (string, string, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected 'owner/repo'")
	}
	return parts[0], parts[1], nil
}

func isValidFilePath(filename string) bool {
	allowedPrefixes := []string{
		".github/",
		"cmd/",
	}
	allowedFiles := map[string]bool{
		"README.md":  true,
		".gitignore": true,
		"go.mod":     true,
		"go.sum":     true,
		"main.go":    true,
	}

	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(filename, prefix) {
			return true
		}
	}

	if allowedFiles[filename] {
		return true
	}

	if !strings.HasPrefix(filename, "deps/") {
		fmt.Printf("Error: File '%s' is not in an allowed directory. All dependency files must be in 'deps/'.\n", filename)
		return false
	}

	parts := strings.Split(filename, "/")
	if len(parts) != 4 {
		fmt.Printf("Error: File path '%s' does not match the required 'deps/[id]/[version]/[filename]' structure.\n", filename)
		return false
	}

	return true
}

func validateManifest(dir string) bool {
	manifestPath := filepath.Join(dir, "dep.json")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		fmt.Printf("Error: Manifest 'dep.json' not found in '%s'. Each dependency version must have a manifest.\n", dir)
		return false
	}

	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		fmt.Printf("Error reading manifest file '%s': %v\n", manifestPath, err)
		return false
	}

	var manifest DepManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		fmt.Printf("Error parsing manifest JSON in '%s': %v\n", manifestPath, err)
		return false
	}

	parts := strings.Split(dir, string(os.PathSeparator))
	depID := parts[1]
	depVersion := parts[2]

	if manifest.ID == "" || manifest.Version == "" || manifest.SourceURL == "" || len(manifest.Files) == 0 {
		fmt.Printf("Error: Manifest '%s' is missing one or more required fields (id, version, sourceUrl, files).\n", manifestPath)
		return false
	}

	if manifest.ID != depID {
		fmt.Printf("Error: Manifest ID '%s' does not match directory ID '%s' in '%s'.\n", manifest.ID, depID, dir)
		return false
	}

	if manifest.Version != depVersion {
		fmt.Printf("Error: Manifest version '%s' does not match directory version '%s' in '%s'.\n", manifest.Version, depVersion, dir)
		return false
	}

	for _, file := range manifest.Files {
		filePath := filepath.Join(dir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' listed in manifest '%s' was not found in the directory.\n", file, manifestPath)
			return false
		}
	}

	dirFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory '%s': %v\n", dir, err)
		return false
	}

	manifestFileMap := make(map[string]bool)
	for _, f := range manifest.Files {
		manifestFileMap[f] = true
	}
	manifestFileMap["dep.json"] = true

	for _, dirFile := range dirFiles {
		if !dirFile.IsDir() {
			if _, ok := manifestFileMap[dirFile.Name()]; !ok {
				fmt.Printf("Error: File '%s' exists in directory '%s' but is not listed in the 'dep.json' manifest.\n", dirFile.Name(), dir)
				return false
			}
		}
	}

	return true
}
