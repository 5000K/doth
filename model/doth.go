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

	// Version of doth used. Incompable versions will cause an error on load.
	// There are automatic migrations from older->newer, but not the other way around.
	DothFormatVersion uint32 `yaml:"dothVersionDoNotEditManually"`
}

func LoadDothFileFromCwd() (*DothFile, error) {
	d, err := LoadDothFileFromPath(util.CleanPath(DothFileLocationAlt))
	if err != nil {
		return LoadDothFileFromPath(util.CleanPath(DothFileLocation))
	}
	return d, nil
}

func LoadDothFileFromPath(path string) (*DothFile, error) {
	dothFile := &DothFile{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, dothFile)
	if err != nil {
		return nil, err
	}

	return dothFile, nil
}
