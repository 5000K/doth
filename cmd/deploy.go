package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/5000K/doth/model"
	"github.com/5000K/doth/util"
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

		configPaths, err := cmd.Flags().GetStringSlice("config")
		if err != nil {
			fmt.Println("failed to get config paths: " + err.Error())
			return
		}

		// 0 configs are valid - not every user needs configs.
		var configs model.ConfigMap = nil

		if len(configPaths) > 0 {
			configs, err = model.LoadConfigFiles(configPaths)
			if err != nil {
				fmt.Println("failed to load config files: " + err.Error())
				return
			}
		}

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

		autoconfirm, err := cmd.Flags().GetBool("autoconfirm")
		if err != nil {
			fmt.Printf("Error getting autoconfirm flag: %v\n", err)
			return
		}

		_ = len(configs)

		planDeploy(configs).Run(
			dry, verbose, model.CreatePipelineConfig(
				true, autoconfirm,
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringSliceP("config", "c", make([]string, 0), "Path to the config.yaml to apply to the doth project when deploying it")
	deployCmd.Flags().BoolP("verbose", "v", false, "Print verbose output")
	deployCmd.Flags().BoolP("dry", "d", false, "Print the actions that would be taken without actually performing them")
	deployCmd.Flags().BoolP("autoconfirm", "y", false, "Automatically confirm all prompts with 'yes'")
}

func readModules() []model.Module {
	doth, err := model.LoadDothFileFromCwd()
	if err != nil {
		panic("failed to load doth file: " + err.Error())
	}

	moduleRoot := util.CleanPath(doth.ModulePath)
	patterns, err := filepath.Glob(filepath.Join(moduleRoot, "*", "module.y*ml"))
	if err != nil {
		panic("failed to glob for module files: " + err.Error())
	}

	modules := make([]model.Module, 0)

	for _, moduleConfig := range patterns {
		mod, err := model.LoadConfigFileFromPath[model.Module](moduleConfig)
		if err != nil {
			fmt.Printf("failed to load module config from %s: %v\n", moduleConfig, err)
			continue
		}

		mod.BasePath = filepath.Dir(moduleConfig)
		modules = append(modules, *mod)
	}
	return modules
}

func planDeploy(configs model.ConfigMap) *model.Pipeline {
	pipeline := model.NewPipeline()
	modules := readModules()

	for _, mod := range modules {
		for _, file := range mod.Files {
			sourcePath := filepath.Join(mod.BasePath, file.Name)
			targetPath := util.CleanPath(filepath.Join(mod.Target, file.Name))

			switch file.Strategy {
			case model.StrategyCopy:
				pipeline.AddModule(model.NewCopyFileStep(sourcePath, targetPath))
			case model.StrategyLink:
				pipeline.AddModule(model.NewCreateSymlinkStep(sourcePath, targetPath))
			case model.StrategyRender:
				pipeline.AddModule(model.NewRenderFileStep(sourcePath, configs, targetPath))
			default:
				fmt.Printf("unknown strategy %s for file %s in module %s\n", file.Strategy, file.Name, mod.BasePath)
			}
		}
	}

	return pipeline
}
