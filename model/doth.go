package model

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
