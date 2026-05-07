package model

type DothFile struct {
	// Relative path to the module directory, e.g. "./modules" or "."
	ModulePath string `yaml:"modulePath"`

	// Version of doth used. Incompable versions will cause an error on load.
	// There are automatic migrations from older->newer, but not the other way around.
	DothFormatVersion uint32 `yaml:"dothVersionDoNotEditManually"`
}
