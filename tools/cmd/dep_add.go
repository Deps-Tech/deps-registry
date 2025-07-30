package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
}

func addDependency(cmd *cobra.Command, args []string) {
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
			Prompt:   &survey.Input{Message: "What is the local path to the dependency's source files?"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		ID         string
		Version    string
		SourceURL  string
		SourcePath string
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
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

	files, err := copyFiles(answers.SourcePath, versionPath)
	if err != nil {
		fmt.Printf("Error copying files: %v\n", err)
		return
	}

	manifest := DepManifest{
		ID:        answers.ID,
		Version:   answers.Version,
		SourceURL: answers.SourceURL,
		Files:     files,
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
