package model

import (
	"os"

	"github.com/5000K/doth/util"
	"gopkg.in/yaml.v3"
)

type Dependency struct {
	Name string `yaml:"name"`

	// a map of source -> package name.
	// sources are provided by the user
	Packages map[string]string `yaml:"packages"`
}

type DothFile struct {
	// Relative path to the module directory, e.g. "./modules" or "."
	ModulePath string `yaml:"modulePath"`

	Deps []Dependency `yaml:"deps"`

	RequireConfig bool `yaml:"requireConfig"`

	// Version of doth used. Incompable versions will cause an error on load.
	// There are automatic migrations from older->newer, but not the other way around.
	DothFormatVersion uint32 `yaml:"dothVersionDoNotEditManually"`
}

func LoadDothFileFromCwd() (*DothFile, error) {
	d, err := LoadConfigFileFromPath[DothFile](util.CleanPath(DothFileLocationAlt))
	if err != nil {
		return LoadConfigFileFromPath[DothFile](util.CleanPath(DothFileLocation))
	}
	return d, nil
}

func LoadConfigFileFromPath[T any](path string) (*T, error) {
	configFile := new(T)

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, configFile)
	if err != nil {
		return nil, err
	}

	return configFile, nil
}
