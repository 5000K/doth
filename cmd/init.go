package cmd

import (
	"fmt"

	"github.com/5000K/doth/model"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a doth project in the current directory.",
	Long: `Initializes a doth project in the current directory.
	
This will create a doth.yaml file with the provided configuration, and a modules directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		dry, err := cmd.Flags().GetBool("dry")
		if err != nil {
			fmt.Printf("Error getting dry flag: %v\n", err)
			return
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Printf("Error getting verbose flag: %v\n", err)
			return
		}

		destructive, err := cmd.Flags().GetBool("destructive")
		if err != nil {
			fmt.Printf("Error getting destructive flag: %v\n", err)
			return
		}

		autoconfirm, err := cmd.Flags().GetBool("autoconfirm")
		if err != nil {
			fmt.Printf("Error getting autoconfirm flag: %v\n", err)
			return
		}

		modulesDir, err := cmd.Flags().GetString("modules")
		if err != nil {
			fmt.Printf("Error getting modules flag: %v\n", err)
			return
		}

		config := model.CreatePipelineConfig(destructive, autoconfirm)

		pipeline := model.NewPipeline()

		if destructive {
			pipeline.AddModule(model.NewConfirmStep("Force create a new project?"))
		}

		pipeline.
			AddModule(model.NewCreateDirStep(modulesDir, "modules directory")).
			AddModule(model.NewCreateDirStep(model.LocalStateDir, "local state directory")).
			AddModule(model.NewCreateFileStep(model.DothFileLocation, []byte(model.DothFileTemplate), "default configuration file")).
			AddModule(model.NewCreateFileStep(model.GitignoreFileLocation, []byte(model.GitignoreFileTemplate), "default .gitignore file")).
			AddModule(model.NewCreateFileStepWithPermissions(model.DothShWrapperLocation, []byte(model.DothShWrapper), "doth.sh wrapper script", 0744))

		err = pipeline.Run(dry, verbose, config)

		if err != nil {
			fmt.Printf("Error initializing project: %v\n", err)
		} else {
			fmt.Println("Project initialized successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("modules", "m", "./modules", "The directory the modules are located in. Relative to the current working directory")
	initCmd.Flags().BoolP("dry", "d", false, "Print the actions that would be taken without actually performing them")
	initCmd.Flags().BoolP("verbose", "v", false, "Print verbose output when running commands")
	initCmd.Flags().Bool("destructive", false, "Deletes and recreates the whole project. This should be used with caution")
	initCmd.Flags().BoolP("autoconfirm", "y", false, "Automatically confirm all prompts with 'yes'")
}
