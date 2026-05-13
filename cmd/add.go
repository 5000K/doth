package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/5000K/doth/model"
	"github.com/5000K/doth/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		if util.ConfirmRunIfRoot() == false {
			return
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Printf("Error getting name flag: %v\n", err)
			return
		}

		target, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Printf("Error getting target flag: %v\n", err)
			return
		}

		skipExisting, err := cmd.Flags().GetBool("skip-existing")
		if err != nil {
			fmt.Printf("Error getting skip-existing flag: %v\n", err)
			return
		}

		destructive, err := cmd.Flags().GetBool("destructive")
		if err != nil {
			fmt.Printf("Error getting destructive flag: %v\n", err)
			return
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Printf("Error getting verbose flag: %v\n", err)
			return
		}

		dry, err := cmd.Flags().GetBool("dry")
		if err != nil {
			fmt.Printf("Error getting dry flag: %v\n", err)
			return
		}

		autoconfirm, err := cmd.Flags().GetBool("autoconfirm")
		if err != nil {
			fmt.Printf("Error getting autoconfirm flag: %v\n", err)
			return
		}

		glob, err := cmd.Flags().GetString("glob")
		if err != nil {
			fmt.Printf("Error getting glob flag: %v\n", err)
			return
		}

		pipeline := planAddPipeline(name, target, skipExisting, destructive, glob)
		config := model.CreatePipelineConfig(destructive, autoconfirm)
		err = pipeline.Run(dry, verbose, config)

		if err != nil {
			fmt.Printf("Error adding module: %v\n", err)
		} else {
			fmt.Printf("Module '%s' added successfully.\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("name", "n", "", "Internal name of the module")
	addCmd.MarkFlagRequired("name")

	addCmd.Flags().StringP("target", "t", "", "Target path of the module")
	addCmd.MarkFlagRequired("target")

	addCmd.Flags().StringP("glob", "g", "**/*", "A glob pattern to specify which files to include. Relative to the target path")

	addCmd.Flags().BoolP("skip-existing", "s", false, "Skip copying existing files in the target folder")
	addCmd.Flags().BoolP("destructive", "f", false, "Deletes and recreates the module if it already exists")
	addCmd.Flags().BoolP("verbose", "v", false, "Print verbose output")
	addCmd.Flags().BoolP("dry", "d", false, "Print the actions that would be taken without actually performing them")
	addCmd.Flags().BoolP("autoconfirm", "y", false, "Automatically confirm all prompts with 'yes'")
}

func planAddPipeline(name string, path string, skipExisting bool, destructive bool, glob string) *model.Pipeline {
	pipeline := model.NewPipeline()

	doth, err := model.LoadDothFileFromCwd()
	if err != nil {
		panic("couldn't find a valid doth.y[a]ml in the current directory: " + err.Error())
	}

	if !util.IsValidFolderName(name) {
		panic(fmt.Errorf("invalid module name '%s'. Module names must be valid folder names.", name))
	}

	moduleRoot := util.CleanPath(doth.ModulePath)
	targetModulePath := filepath.Join(moduleRoot, name)
	targetSourcePath := util.CleanPath(path)

	// if target doesn't yet exist: set skipExisting to true, since there's nothing to skip
	if _, err := os.Stat(targetSourcePath); err != nil {
		pipeline.AddModule(model.NewLogStep("Couldn't find target directory - skipping file copy", true))
		skipExisting = true
	}

	if destructive {
		pipeline.AddModule(model.NewConfirmStep(fmt.Sprintf("Module '%s' already exists. Do you want to delete and recreate it?", name)))
	}

	pipeline.AddModule(model.NewCreateDirStep(targetModulePath, "existing module directory"))

	module := model.Module{
		Skip:   false,
		Files:  make([]model.ModuleFile, 0),
		Deps:   make([]model.Dependency, 0),
		Target: path,
	}

	// find all files in target path, add a copy file step for each of them, unless skipExisting is true
	if !skipExisting {
		existingFiles, err := filepath.Glob(filepath.Join(targetSourcePath, glob))
		if err != nil {
			panic(fmt.Errorf("error listing files in the target source path: %w", err))
		} else {
			for _, file := range existingFiles {
				// is folder - create dir step, is file - copy file step.
				info, err := os.Stat(file)
				if err != nil {
					panic(fmt.Errorf("error stating file '%s': %w", file, err))
				}
				if info.IsDir() {
					pipeline.AddModule(model.NewCreateDirStep(filepath.Join(targetModulePath, file), "existing directory"))
					continue
				}

				relativePath, _ := filepath.Rel(targetSourcePath, file)
				targetPath := filepath.Join(targetModulePath, relativePath)
				pipeline.AddModule(model.NewCopyFileStep(file, targetPath))

				module.Files = append(module.Files, model.ModuleFile{
					Name:     relativePath,
					Strategy: model.StrategyCopy,
				})
			}
		}
	}

	moduleFileData, err := yaml.Marshal(module)
	if err != nil {
		panic(fmt.Errorf("error creating module config: %w", err))
	}

	pipeline.AddModule(model.NewCreateFileStep(filepath.Join(targetModulePath, model.ModuleFileName), moduleFileData, "module configuration file"))

	return pipeline
}
