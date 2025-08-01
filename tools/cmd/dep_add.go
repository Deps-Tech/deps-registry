package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new dependency",
	Run:   addDependency,
}

func init() {
	depCmd.AddCommand(addCmd)
	addCmd.Flags().String("id", "", "Dependency ID")
	addCmd.Flags().String("version", "", "Dependency version")
	addCmd.Flags().String("source-url", "", "Source URL")
	addCmd.Flags().StringSlice("source-path", []string{}, "Local path(s) to source files (can be specified multiple times)")
	addCmd.Flags().String("deps", "", "Dependencies in format: dep1:ver1,dep2:ver2")
}

func addDependency(cmd *cobra.Command, args []string) {
	var err error
	answers := struct {
		ID         string
		Version    string
		SourceURL  string
		SourcePath string // For interactive mode, we'll take a comma-separated string
		Deps       string
	}{}

	flags := cmd.Flags()
	id, _ := flags.GetString("id")

	if id == "" {
		questions := []*survey.Question{
			{
				Name:     "id",
				Prompt:   &survey.Input{Message: "What is the new dependency's ID (e.g., 'awesome-lib')?"},
				Validate: survey.Required,
			},
			{
				Name:     "version",
				Prompt:   &survey.Input{Message: "What is the version to add (e.g., '1.2.0')?"},
				Validate: survey.Required,
			},
			{
				Name:     "sourceUrl",
				Prompt:   &survey.Input{Message: "What is the source URL (link to forum, etc.)?"},
				Validate: survey.Required,
			},
			{
				Name:     "sourcePath",
				Prompt:   &survey.Input{Message: "What are the local paths to the dependency's source files (comma-separated)?"},
				Validate: survey.Required,
			},
			{
				Name:   "deps",
				Prompt: &survey.Input{Message: "Enter dependencies in format: dep1:ver1,dep2:ver2"},
			},
		}
		err = survey.Ask(questions, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.ID = id
		answers.Version, _ = flags.GetString("version")
		answers.SourceURL, _ = flags.GetString("source-url")
		// This flag is now a slice, but we handle it later to support both modes.
		// For simplicity, we don't assign it to the answers struct here.
		answers.Deps, _ = flags.GetString("deps")
	}

	depPath := filepath.Join("..", "deps", answers.ID)
	if _, err := os.Stat(depPath); !os.IsNotExist(err) {
		fmt.Printf("Error: Dependency '%s' already exists. Use the 'update' command to add a new version.\n", answers.ID)
		return
	}

	versionPath := filepath.Join(depPath, answers.Version)
	err = os.MkdirAll(versionPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory structure: %v\n", err)
		return
	}

	var sourcePaths []string
	pathsFromFlags, _ := flags.GetStringSlice("source-path")

	if len(pathsFromFlags) > 0 {
		sourcePaths = pathsFromFlags
	} else {
		// From interactive survey, split the comma-separated string
		sourcePaths = strings.Split(answers.SourcePath, ",")
	}

	var allFiles []string
	for _, p := range sourcePaths {
		sourcePath := strings.Trim(strings.TrimSpace(p), `"`)
		if sourcePath == "" {
			continue
		}
		files, err := copyFiles(sourcePath, versionPath)
		if err != nil {
			fmt.Printf("Error copying files from '%s': %v\n", sourcePath, err)
			return
		}
		allFiles = append(allFiles, files...)
	}

	depVersions := make(map[string]string)
	if answers.Deps != "" {
		depPairs := strings.Split(answers.Deps, ",")
		for _, pair := range depPairs {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				depVersions[parts[0]] = parts[1]
			}
		}
	}

	manifest := DepManifest{
		ID:           answers.ID,
		Version:      answers.Version,
		SourceURL:    answers.SourceURL,
		Files:        allFiles,
		Dependencies: depVersions,
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Printf("Error creating manifest: %v\n", err)
		return
	}

	manifestPath := filepath.Join(versionPath, "dep.json")
	err = ioutil.WriteFile(manifestPath, manifestData, 0644)
	if err != nil {
		fmt.Printf("Error writing manifest file: %v\n", err)
		return
	}

	fmt.Printf("Success: Added dependency '%s' version '%s'. The following files are now staged for commit. Please review and then run 'git add .', 'git commit', and 'git push'.\n", answers.ID, answers.Version)
}
