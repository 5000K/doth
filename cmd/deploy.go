package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the modules to their target locations.",
	Long: `Deploys the modules to their target locations.

Takes in one or more config.yaml files that will be applied when deploying the modules. The config files will be merged together, with the values from the later files taking precedence over the earlier ones.

This command will copy, symlink, or render the files from the modules to their respective target locations based on the defined strategies.
If a module has dependencies, they will be installed first using the appropriate strategies.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringSliceP("config", "c", make([]string, 0), "Path to the config.yaml to apply to the doth project when deploying it.")
}
