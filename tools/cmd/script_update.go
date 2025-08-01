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

var scriptUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing script with a new version",
	Run:   updateScript,
}

func init() {
	scriptCmd.AddCommand(scriptUpdateCmd)
	scriptUpdateCmd.Flags().String("id", "", "Script ID")
	scriptUpdateCmd.Flags().String("name", "", "Script Name")
	scriptUpdateCmd.Flags().String("version", "", "New script version")
	scriptUpdateCmd.Flags().String("source-url", "", "New source URL")
	scriptUpdateCmd.Flags().String("source-path", "", "Local path to new source files")
	scriptUpdateCmd.Flags().String("tags", "", "New comma-separated tags")
	scriptUpdateCmd.Flags().String("deps", "", "Dependencies in format: dep1:ver1,dep2:ver2")
}

func updateScript(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()
	idAnswer, _ := flags.GetString("id")

	if idAnswer == "" {
		idQuestion := &survey.Question{
			Name:     "id",
			Prompt:   &survey.Input{Message: "What is the ID of the script to update?"},
			Validate: survey.Required,
		}
		err := survey.Ask([]*survey.Question{idQuestion}, &idAnswer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
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

	answers := struct {
		Name       string
		Version    string
		SourceURL  string
		SourcePath string
		Tags       string
	}{}

	version, _ := flags.GetString("version")
	if version == "" {
		questions := []*survey.Question{
			{
				Name:     "version",
				Prompt:   &survey.Input{Message: "What is the new version to add?"},
				Validate: survey.Required,
			},
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "What is the new script's name?"},
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
			{
				Name:   "tags",
				Prompt: &survey.Input{Message: "Enter up to 10 comma-separated tags (e.g., 'tag1,tag2'):"},
			},
		}
		err = survey.Ask(questions, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.Version = version
		answers.Name, _ = flags.GetString("name")
		answers.SourceURL, _ = flags.GetString("source-url")
		answers.SourcePath, _ = flags.GetString("source-path")
		answers.Tags, _ = flags.GetString("tags")
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

	sourcePath := strings.Trim(answers.SourcePath, `"`)
	analysis, err := analyzeScript(sourcePath)
	if err != nil {
		fmt.Printf("Error analyzing script: %v\n", err)
		return
	}
	deps := analysis.Dependencies

	depVersions := make(map[string]string)
	depsFlag, _ := flags.GetString("deps")
	if depsFlag != "" {
		depPairs := strings.Split(depsFlag, ",")
		for _, pair := range depPairs {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				depVersions[parts[0]] = parts[1]
			}
		}
	} else if len(deps) > 0 {
		fmt.Println("Found dependencies. Please select the versions to use:")
		for _, dep := range deps {
			versions, err := getAvailableVersions(dep)
			if err != nil {
				fmt.Printf("Error getting versions for dependency '%s': %v\n", dep, err)
				return
			}
			if len(versions) == 0 {
				fmt.Printf("No versions found for dependency '%s'. Please add it first.\n", dep)
				return
			}

			selectedVersion := ""
			prompt := &survey.Select{
				Message: fmt.Sprintf("Select version for '%s':", dep),
				Options: versions,
			}
			survey.AskOne(prompt, &selectedVersion, survey.WithValidator(survey.Required))
			depVersions[dep] = selectedVersion
		}
	}

	if len(depVersions) > 0 {
		resolvedDeps, err := resolveDependencies(depVersions)
		if err != nil {
			fmt.Printf("Error resolving dependencies: %v\n", err)
			return
		}
		depVersions = resolvedDeps
	}

	if analysis.HasNetworkAccess {
		fmt.Println("Warning: This script appears to make network requests.")
	}

	if len(analysis.TouchedFiles) > 0 {
		fmt.Println("Found file operations. Please confirm to add them to the manifest:")
		for _, file := range analysis.TouchedFiles {
			fmt.Printf("- %s\n", file)
		}
		confirm := false
		prompt := &survey.Confirm{
			Message: "Add these files to the manifest?",
			Default: true,
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			analysis.TouchedFiles = nil
		}
	}

	files, err := copyFiles(answers.SourcePath, versionPath)
	if err != nil {
		fmt.Printf("Error copying files: %v\n", err)
		return
	}

	tags := parseTags(answers.Tags)
	if len(tags) > 10 {
		fmt.Println("Error: You can only specify up to 10 tags.")
		return
	}

	manifest := DepManifest{
		ID:               idAnswer,
		Name:             answers.Name,
		Version:          answers.Version,
		SourceURL:        answers.SourceURL,
		Files:            files,
		Dependencies:     depVersions,
		Tags:             tags,
		TouchedFiles:     analysis.TouchedFiles,
		HasNetworkAccess: analysis.HasNetworkAccess,
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
