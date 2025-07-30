package cmd

import "github.com/spf13/cobra"

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Manage scripts",
}

func init() {
	rootCmd.AddCommand(scriptCmd)
}
