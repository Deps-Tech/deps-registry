package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Deps-Tech/deps-registry/tools/internal/indexer"
	"github.com/spf13/cobra"
)

var (
	cdnURL string
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate index.json for the registry",
	Run:   runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&cdnURL, "cdn-url", "", "Base URL of the CDN")
	indexCmd.MarkFlagRequired("cdn-url")
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) {
	distPath := "dist"

	idx, err := indexer.Generate(distPath, cdnURL)
	if err != nil {
		fmt.Printf("Failed to generate index: %v\n", err)
		os.Exit(1)
	}

	indexPath := filepath.Join(distPath, "index.json")
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal index: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		fmt.Printf("Failed to write index: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated index.json with %d dependencies and %d scripts\n",
		len(idx.Dependencies), len(idx.Scripts))
}

