package model

type Dependency struct {
	Name string `yaml:"name"`

	// a map of source -> package name.
	// sources are provided by the user
	Packages map[string]string `yaml:"packages"`
}

type Module struct {
	BasePath string       // module path on disk, dynamically set when loading
	Files    []ModuleFile `yaml:"files"`
	Target   string       `yaml:"target"`
	Deps     []Dependency `yaml:"deps"`
}

type FileStrategy string

const (
	StrategyCopy    FileStrategy = "copy"
	StrategySymlink FileStrategy = "symlink"
	StrategyRender  FileStrategy = "render"
)

type ModuleFile struct {
	// Name is the filename or relative path to the file within the module folder - may use glob patterns
	Name string `yaml:"name"`
	// Strategy is the strategy to use when deploying the file to it's target location. Defaults to "copy".
	Strategy FileStrategy `yaml:"strategy"`
}
