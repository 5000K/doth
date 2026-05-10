/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/5000K/doth/model"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of doth.",
	Long:  `Print the version of doth.`,
	Run: func(cmd *cobra.Command, args []string) {
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			fmt.Printf("Error getting raw flag: %v\n", err)
			return
		}

		if raw {
			fmt.Printf("%s\n", model.Version)
		} else {
			fmt.Printf("doth dotfile manager, version %s\n", model.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("raw", "r", false, "Print only the version number")
}
