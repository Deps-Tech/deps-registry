package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Deps-Tech/deps-registry/tools/internal/filesystem"
	"github.com/Deps-Tech/deps-registry/tools/internal/manifest"
	"github.com/Deps-Tech/deps-registry/tools/internal/parser"
	"github.com/Deps-Tech/deps-registry/tools/internal/registry"
	"github.com/spf13/cobra"
)

var (
	sourcePath string
	tags       string
)

var addCmd = &cobra.Command{
	Use:   "add [script|dep]",
	Short: "Add a new script or dependency to the registry",
}

var addScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Add a new script to the registry",
	Run:   runAddScript,
}

var addDepCmd = &cobra.Command{
	Use:   "dep",
	Short: "Add a new dependency to the registry",
	Run:   runAddDep,
}

func init() {
	addScriptCmd.Flags().StringVar(&sourcePath, "source", "", "Path to script file or directory")
	addScriptCmd.Flags().StringVar(&tags, "tags", "", "Comma-separated tags")
	addScriptCmd.MarkFlagRequired("source")

	addDepCmd.Flags().StringVar(&sourcePath, "source", "", "Path to dependency file or directory")
	addDepCmd.MarkFlagRequired("source")

	addCmd.AddCommand(addScriptCmd, addDepCmd)
	rootCmd.AddCommand(addCmd)
}

func runAddScript(cmd *cobra.Command, args []string) {
	if err := addItem("scripts", sourcePath, tags); err != nil {
		fmt.Printf("Error adding script: %v\n", err)
		os.Exit(1)
	}
}

func runAddDep(cmd *cobra.Command, args []string) {
	if err := addItem("deps", sourcePath, ""); err != nil {
		fmt.Printf("Error adding dependency: %v\n", err)
		os.Exit(1)
	}
}

func addItem(itemType, source, tagList string) error {
	client := registry.NewClient("")

	metadata, err := extractMetadata(source)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	fmt.Printf("\nExtracted metadata:\n")
	fmt.Printf("  ID: %s\n", metadata.ID)
	fmt.Printf("  Name: %s\n", metadata.Name)
	fmt.Printf("  Version: %s\n", metadata.Version)
	if metadata.Author != "" {
		fmt.Printf("  Author: %s\n", metadata.Author)
	}

	dupInfo, err := client.CheckDuplicate(itemType, metadata.ID, metadata.Version)
	if err != nil && client.IsAvailable() {
		return fmt.Errorf("failed to check for duplicates: %w", err)
	}

	if dupInfo != nil && dupInfo.ExactMatch {
		fmt.Printf("\n❌ Error: %s '%s' v%s already exists\n", itemType, metadata.ID, metadata.Version)
		if dupInfo.PackageURL != "" {
			fmt.Printf("   URL: %s\n", dupInfo.PackageURL)
		}
		fmt.Printf("\n   To add a new version, use a different version number.\n")
		return fmt.Errorf("package already exists")
	}

	if dupInfo != nil && dupInfo.Exists && !dupInfo.ExactMatch {
		fmt.Printf("\n⚠️  Warning: %s '%s' already exists with versions: %v\n",
			itemType, metadata.ID, dupInfo.AllVersions)
		fmt.Printf("   You are adding a new version: %s\n", metadata.Version)
		fmt.Printf("\n   Continue? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("\n   Aborted.")
			return fmt.Errorf("user cancelled")
		}
	}

	availableDeps := getAvailableDeps()
	analysis, err := parser.AnalyzeLua(source, metadata.ID, availableDeps)
	if err != nil {
		return fmt.Errorf("failed to analyze: %w", err)
	}

	var depVersions map[string]string
	if len(analysis.Dependencies) > 0 {
		fmt.Printf("\nFound dependencies:\n")
		_, depVersions = resolveDependencies(client, analysis.Dependencies)
		for _, dep := range analysis.Dependencies {
			version := depVersions[dep]
			if version == "*" {
				fmt.Printf("  - %s (*) ⚠️  not in registry\n", dep)
			} else {
				fmt.Printf("  - %s (%s) ✓\n", dep, version)
			}
		}
	}

	if analysis.UsesNetwork {
		fmt.Println("\n⚠️  Uses network access")
	}
	if analysis.UsesFFI {
		fmt.Println("\n⚠️  Uses FFI")
	}
	if len(analysis.Warnings) > 0 {
		fmt.Printf("\n⚠️  %d warnings (dynamic requires detected):\n", len(analysis.Warnings))
		for _, w := range analysis.Warnings {
			fmt.Printf("   Line %d: %s\n", w.Line, w.Message)
		}
	}

	targetPath := filepath.Join("..", itemType, metadata.ID, metadata.Version)
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	files, err := copySourceFiles(source, targetPath)
	if err != nil {
		return err
	}

	fileMap := make(map[string]manifest.FileInfo)
	for _, file := range files {
		filePath := filepath.Join(targetPath, file)
		hash, _ := filesystem.SHA256File(filePath)
		info, _ := os.Stat(filePath)
		fileMap[file] = manifest.FileInfo{
			SHA256: hash,
			Size:   info.Size(),
		}
	}

	deps := make(map[string]string)
	for _, dep := range analysis.Dependencies {
		if version, ok := depVersions[dep]; ok {
			deps[dep] = version
		} else {
			deps[dep] = "*"
		}
	}

	tagSlice := []string{}
	if tagList != "" {
		for _, tag := range strings.Split(tagList, ",") {
			tagSlice = append(tagSlice, strings.TrimSpace(tag))
		}
	}

	m := &manifest.Manifest{
		ManifestVersion: "1.0",
		ID:              metadata.ID,
		Name:            metadata.Name,
		Version:         metadata.Version,
		Files:           fileMap,
		Dependencies:    deps,
		Security: manifest.Security{
			NetworkAccess: analysis.UsesNetwork,
			FileAccess:    analysis.FilePaths,
			UsesFFI:       analysis.UsesFFI,
		},
		Metadata: manifest.Metadata{
			Tags: tagSlice,
		},
	}

	if err := manifest.Save(targetPath, m); err != nil {
		return err
	}

	fmt.Printf("\n✓ Created %s/%s/%s\n", itemType, metadata.ID, metadata.Version)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review generated dep.json")
	fmt.Println("  2. git add " + filepath.Join(itemType, metadata.ID))
	fmt.Printf("  3. git commit -m \"feat(%s): add %s v%s\"\n", itemType, metadata.ID, metadata.Version)
	fmt.Println("  4. git push and create Pull Request")

	return nil
}

type Metadata struct {
	ID      string
	Name    string
	Version string
	Author  string
}

func extractMetadata(source string) (*Metadata, error) {
	luaFiles, err := findLuaFiles(source)
	if err != nil {
		return nil, err
	}

	if len(luaFiles) == 0 {
		return nil, fmt.Errorf("no Lua files found")
	}

	content, err := os.ReadFile(luaFiles[0])
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	name := extractField(contentStr, `script_name\s*\(\s*["'](.+?)["']\s*\)`)
	version := extractField(contentStr, `script_version\s*\(\s*["'](.+?)["']\s*\)`)
	author := extractField(contentStr, `script_author\s*\(\s*["'](.+?)["']\s*\)`)

	if name == "" {
		name = filepath.Base(source)
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	if version == "" {
		version = "1.0.0"
	}

	id := strings.ToLower(name)
	id = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(id, "-")
	id = strings.Trim(id, "-")

	return &Metadata{
		ID:      id,
		Name:    name,
		Version: version,
		Author:  author,
	}, nil
}

func extractField(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func findLuaFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if filepath.Ext(path) == ".lua" {
			return []string{path}, nil
		}
		return []string{}, nil
	}

	var files []string
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(p) == ".lua" {
			files = append(files, p)
		}
		return nil
	})

	return files, err
}

func copySourceFiles(source, target string) ([]string, error) {
	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	var files []string

	if !info.IsDir() {
		fileName := filepath.Base(source)
		targetFile := filepath.Join(target, fileName)
		content, err := os.ReadFile(source)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(targetFile, content, 0644); err != nil {
			return nil, err
		}
		return []string{fileName}, nil
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(target, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.WriteFile(targetPath, content, info.Mode()); err != nil {
			return err
		}

		files = append(files, relPath)
		return nil
	})

	return files, err
}

func getAvailableDeps() map[string]bool {
	deps := make(map[string]bool)
	depsPath := filepath.Join("..", "deps")

	items, err := os.ReadDir(depsPath)
	if err != nil {
		return deps
	}

	for _, item := range items {
		if item.IsDir() {
			deps[strings.ToLower(item.Name())] = true
		}
	}

	return deps
}

func resolveDependencies(client *registry.Client, deps []string) ([]string, map[string]string) {
	cdnAvailable := client.IsAvailable()
	if !cdnAvailable {
		fmt.Printf("\n⚠️  Warning: Cannot reach registry CDN\n")
		fmt.Printf("   Continuing with unknown dependency versions (*)\n")
	}

	versions := make(map[string]string)
	uniqueDeps := make(map[string]bool)
	result := []string{}

	for _, dep := range deps {
		if uniqueDeps[dep] {
			continue
		}
		uniqueDeps[dep] = true
		result = append(result, dep)

		if cdnAvailable {
			version, err := client.GetLatestVersion("deps", dep)
			if err == nil {
				versions[dep] = version
			} else {
				versions[dep] = "*"
			}
		} else {
			versions[dep] = "*"
		}
	}

	return result, versions
}
