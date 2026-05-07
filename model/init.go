package model

import (
	"github.com/5000K/doth/util"
)

type InitConfig struct {
	PipelineConfig
	ModuleDir string `yaml:"moduleDir"`
	Force     bool   `yaml:"force"`
}

func (c *InitConfig) ForceOperations() bool {
	return c.Force
}

func Init(config *InitConfig, dry bool, verbose bool) {
	config.ModuleDir = util.CleanPath(config.ModuleDir)

	pipeline := NewPipeline[*InitConfig]().
		AddModule(newCreateDirStep[*InitConfig](config.ModuleDir, "modules directory")).
		AddModule(newCreateDirStep[*InitConfig](localStateDir, "local state directory")).
		AddModule(newCreateFileStep[*InitConfig](dothFileLocation, []byte(dothFileTemplate), "default configuration file")).
		AddModule(newCreateFileStep[*InitConfig](gitignoreFileLocation, []byte(gitignoreFileTemplate), "default .gitignore file"))

	if dry {
		err := pipeline.ApplyDry(config)
		if err != nil {
			panic(err)
		}
	} else if verbose {
		err := pipeline.ApplyVerbose(config)
		if err != nil {
			panic(err)
		}
	} else {
		err := pipeline.Apply(config)
		if err != nil {
			panic(err)
		}
	}
}
