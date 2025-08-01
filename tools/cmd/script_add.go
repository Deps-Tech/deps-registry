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

var scriptAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new script",
	Run:   addScript,
}

func init() {
	scriptCmd.AddCommand(scriptAddCmd)
	scriptAddCmd.Flags().String("id", "", "Script ID")
	scriptAddCmd.Flags().String("name", "", "Script Name")
	scriptAddCmd.Flags().String("version", "", "Script version")
	scriptAddCmd.Flags().String("source-url", "", "Source URL")
	scriptAddCmd.Flags().String("source-path", "", "Local path to source files")
	scriptAddCmd.Flags().String("tags", "", "Comma-separated tags")
	scriptAddCmd.Flags().String("deps", "", "Dependencies in format: dep1:ver1,dep2:ver2")
}

func addScript(cmd *cobra.Command, args []string) {
	answers := struct {
		ID         string
		Name       string
		Version    string
		SourceURL  string
		SourcePath string
		Tags       string
	}{}

	flags := cmd.Flags()
	id, _ := flags.GetString("id")
	if id == "" {
		questions := []*survey.Question{
			{
				Name:     "id",
				Prompt:   &survey.Input{Message: "What is the new script's ID (e.g., 'my-cool-script')?"},
				Validate: survey.Required,
			},
			{
				Name:     "name",
				Prompt:   &survey.Input{Message: "What is the new script's name?"},
				Validate: survey.Required,
			},
			{
				Name:     "version",
				Prompt:   &survey.Input{Message: "What is the version to add (e.g., '1.0.0')?"},
				Validate: survey.Required,
			},
			{
				Name:     "sourceUrl",
				Prompt:   &survey.Input{Message: "What is the source URL (link to forum, etc.)?"},
				Validate: survey.Required,
			},
			{
				Name:     "sourcePath",
				Prompt:   &survey.Input{Message: "What is the local path to the script's source files?"},
				Validate: survey.Required,
			},
			{
				Name:   "tags",
				Prompt: &survey.Input{Message: "Enter up to 10 comma-separated tags (e.g., 'tag1,tag2'):"},
			},
		}
		err := survey.Ask(questions, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		answers.ID = id
		answers.Name, _ = flags.GetString("name")
		answers.Version, _ = flags.GetString("version")
		answers.SourceURL, _ = flags.GetString("source-url")
		answers.SourcePath, _ = flags.GetString("source-path")
		answers.Tags, _ = flags.GetString("tags")
	}

	sourcePath := strings.Trim(answers.SourcePath, `"`)
	analysis, err := analyzeScript(sourcePath, answers.ID)
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

	scriptPath := filepath.Join("..", "scripts", answers.ID)
	if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
		fmt.Printf("Error: Script '%s' already exists. Use the 'update' command to add a new version.\n", answers.ID)
		return
	}

	versionPath := filepath.Join(scriptPath, answers.Version)
	err = os.MkdirAll(versionPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory structure: %v\n", err)
		return
	}

	files, err := copyFiles(sourcePath, versionPath)
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
		ID:               answers.ID,
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

	fmt.Printf("Success: Added script '%s' version '%s'. The following files are now staged for commit. Please review and then run 'git add .', 'git commit', and 'git push'.\n", answers.ID, answers.Version)
}
