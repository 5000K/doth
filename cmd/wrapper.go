/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/5000K/doth/model"
	"github.com/spf13/cobra"
)

// wrapperCmd represents the wrapper command
var wrapperCmd = &cobra.Command{
	Use:   "wrapper",
	Short: "Prints the current version of the wrapper into stdout.",
	Long:  `Prints the current version of the wrapper into stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(model.DothShWrapper)
	},
}

func init() {
	rootCmd.AddCommand(wrapperCmd)
}
