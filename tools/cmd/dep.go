package cmd

import "github.com/spf13/cobra"

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Manage dependencies",
}

func init() {
	rootCmd.AddCommand(depCmd)
}
