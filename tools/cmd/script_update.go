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

var scriptUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing script with a new version",
	Run:   updateScript,
}

func init() {
	scriptCmd.AddCommand(scriptUpdateCmd)
}

func updateScript(cmd *cobra.Command, args []string) {
	idQuestion := &survey.Question{
		Name:     "id",
		Prompt:   &survey.Input{Message: "What is the ID of the script to update?"},
		Validate: survey.Required,
	}

	var idAnswer string
	err := survey.Ask([]*survey.Question{idQuestion}, &idAnswer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	scriptPath := filepath.Join("..", "scripts", idAnswer)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Error: Script '%s' not found. Use the 'add' command to create it first.\n", idAnswer)
		return
	}

	latestVersion, lastSourceURL, err := getLatestVersionInfo(scriptPath)
	if err != nil {
		fmt.Printf("Error getting latest version: %v\n", err)
		return
	}

	fmt.Printf("The latest known version is %s.\n", latestVersion)

	questions := []*survey.Question{
		{
			Name:     "version",
			Prompt:   &survey.Input{Message: "What is the new version to add?"},
			Validate: survey.Required,
		},
		{
			Name:     "sourceUrl",
			Prompt:   &survey.Input{Message: "What is the source URL?", Default: lastSourceURL},
			Validate: survey.Required,
		},
		{
			Name:     "sourcePath",
			Prompt:   &survey.Input{Message: "What is the local path to the new version's source files?"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Version    string
		SourceURL  string
		SourcePath string
	}{}

	err = survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	versionPath := filepath.Join(scriptPath, answers.Version)
	if _, err := os.Stat(versionPath); !os.IsNotExist(err) {
		fmt.Printf("Error: Version '%s' already exists for this script.\n", answers.Version)
		return
	}

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
		ID:        idAnswer,
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

	fmt.Printf("Success: Updated script '%s' with new version '%s'. Please review and commit the changes.\n", idAnswer, answers.Version)
}
