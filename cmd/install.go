/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/5000K/doth/model"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the dependencies in your main config and all modules.",
	Long:  `Installs the dependencies in your main config and all modules.`,
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

		configPaths, err := cmd.Flags().GetStringSlice("config")
		if err != nil {
			fmt.Println("failed to get config paths: " + err.Error())
			return
		}

		// 0 configs are valid - not every user needs configs.
		var configs model.ConfigMap = make(model.ConfigMap)

		if len(configPaths) > 0 {
			configs, err = model.LoadConfigFiles(configPaths)
			if err != nil {
				fmt.Println("failed to load config files: " + err.Error())
				return
			}
		}

		config, err := configs.GetPackageConfig()
		sources := make([]model.PackageSource, 0)

		if err == nil {
			sources = config.PackageSources
		}

		// check for source flags and add them to the sources list
		if cmd.Flags().Changed("apt") {
			sources = append(sources, model.PackageSource{
				Name:    "apt",
				Command: "sudo apt install {package}",
			})
		}
		if cmd.Flags().Changed("apt-get") {
			sources = append(sources, model.PackageSource{
				Name:    "apt-get",
				Command: "sudo apt-get install {package}",
			})
		}
		if cmd.Flags().Changed("dnf") {
			sources = append(sources, model.PackageSource{
				Name:    "dnf",
				Command: "sudo dnf install {package}",
			})
		}
		if cmd.Flags().Changed("pacman") {
			sources = append(sources, model.PackageSource{
				Name:    "pacman",
				Command: "sudo pacman -S {package}",
			})
		}
		if cmd.Flags().Changed("yay") {
			sources = append(sources, model.PackageSource{
				Name:    "yay",
				Command: "yay -S {package}",
			})
		}
		if cmd.Flags().Changed("paru") {
			sources = append(sources, model.PackageSource{
				Name:    "paru",
				Command: "paru -S {package}",
			})
		}
		if cmd.Flags().Changed("go") {
			sources = append(sources, model.PackageSource{
				Name:    "go",
				Command: "go install {package}",
			})
		}
		if cmd.Flags().Changed("npm") {
			sources = append(sources, model.PackageSource{
				Name:    "npm",
				Command: "npm install -g {package}",
			})
		}
		if cmd.Flags().Changed("brew") {
			sources = append(sources, model.PackageSource{
				Name:    "brew",
				Command: "brew install {package}",
			})
		}

		planInstall(sources).Run(
			dry, verbose, model.CreatePipelineConfig(
				true, false, // no autoconfirm for anything that executes shell code unescaped
			))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringSliceP("config", "c", []string{}, "Path(s) to config.yaml file(s) to use when installing dependencies. The config files will be merged together, with the values from the later files taking precedence over the earlier ones.")

	installCmd.Flags().BoolP("dry", "d", false, "Print the dependencies that would be installed without actually installing them.")
	installCmd.Flags().BoolP("verbose", "v", false, "Print verbose output.")

	// source flags
	installCmd.Flags().Bool("apt", false, "Use apt to install dependencies")
	installCmd.Flags().Bool("apt-get", false, "Use apt-get to install dependencies")
	installCmd.Flags().Bool("dnf", false, "Use dnf to install dependencies")
	installCmd.Flags().Bool("pacman", false, "Use pacman to install dependencies")
	installCmd.Flags().Bool("yay", false, "Use yay to install dependencies")
	installCmd.Flags().Bool("paru", false, "Use paru to install dependencies")
	installCmd.Flags().Bool("go", false, "Use go to install dependencies")
	installCmd.Flags().Bool("npm", false, "Use npm to install dependencies")
	installCmd.Flags().Bool("brew", false, "Use brew to install dependencies")
}

func collectDependencies() []model.Dependency {
	doth, err := model.LoadDothFileFromCwd()
	if err != nil {
		return make([]model.Dependency, 0)
	}

	if doth == nil {
		return make([]model.Dependency, 0)
	}

	dependencies := make([]model.Dependency, len(doth.Deps))
	dependencies = append(dependencies, doth.Deps...)

	modules := readModules()

	for _, module := range modules {
		dependencies = append(dependencies, module.Deps...)
	}

	// filter out duplicate dependencies by name, keeping the last one (which should be the most specific one for the module)
	seen := make(map[string]bool)
	uniqueDependencies := make([]model.Dependency, 0)

	for i := len(dependencies) - 1; i >= 0; i-- {
		dep := dependencies[i]
		if _, ok := seen[dep.Name]; !ok {
			uniqueDependencies = append(uniqueDependencies, dep)
			seen[dep.Name] = true
		}
	}

	// reverse the uniqueDependencies slice to restore the original order (with duplicates removed)
	for i, j := 0, len(uniqueDependencies)-1; i < j; i, j = i+1, j-1 {
		uniqueDependencies[i], uniqueDependencies[j] = uniqueDependencies[j], uniqueDependencies[i]
	}

	return uniqueDependencies
}

func planInstall(sources []model.PackageSource) *model.Pipeline {
	pipeline := model.NewPipeline().AddModule(model.NewConfirmStep("doth install gives this doth project shell execution rights. Only continue if you trust this project 10000% (=> you wrote it)."))

	deps := collectDependencies()

	for _, dep := range deps {
		// find the first source that has a matching package entry for this dependency
		for _, source := range sources {
			packageName, ok := dep.Packages[source.Name]
			if !ok {
				continue
			}

			command := strings.ReplaceAll(source.Command, "{package}", packageName)
			pipeline.AddModule(model.NewExecuteShellCommandStep(command))
			break
		}
	}

	return pipeline
}
