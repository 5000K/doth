/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a module to the project.",
	Long: `Adds a module to the project.

This will create a directory and a doth.yaml file for the module.

If the target directory already exists, all existing files will be added with the strategy "copy".
This can be skipped with the --skip-existing flag (-s).
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "", "Internal name of the module")

	addCmd.Flags().StringP("path", "p", "", "Target path of the module")
	addCmd.MarkFlagRequired("path")

	addCmd.Flags().BoolP("skip-existing", "s", false, "Skip existing files when adding a module")
}
