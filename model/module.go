package model

type Module struct {
	Skip     bool         `yaml:"skip"` // if true, the module will be fully ignored when deploying
	BasePath string       `yaml:"-"`    // module path on disk, dynamically set when loading
	Files    []ModuleFile `yaml:"files"`
	Target   string       `yaml:"target"`
}

type FileStrategy string

const (
	StrategyCopy   FileStrategy = "copy"
	StrategyLink   FileStrategy = "link"
	StrategyRender FileStrategy = "render"
)

type ModuleFile struct {
	// Name is the filename or relative path to the file within the module folder - may use glob patterns
	Name string `yaml:"name"`
	// Strategy is the strategy to use when deploying the file to it's target location. Defaults to "copy".
	Strategy FileStrategy `yaml:"strategy"`
}
