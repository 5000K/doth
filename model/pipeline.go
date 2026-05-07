package model

import "github.com/5000K/doth/util"

type PipelineConfig interface {
	ForceOperations() bool
}

type PipelineModule[T PipelineConfig] interface {
	Apply(input T) error
	ApplyDry(input T) (string, error)
}

type Pipeline[T PipelineConfig] struct {
	modules []PipelineModule[T]
}

func NewPipeline[T PipelineConfig]() *Pipeline[T] {
	return &Pipeline[T]{
		modules: make([]PipelineModule[T], 0),
	}
}

func (p *Pipeline[T]) AddModule(module PipelineModule[T]) *Pipeline[T] {
	p.modules = append(p.modules, module)
	return p
}

func (p *Pipeline[T]) Apply(input T) error {
	for _, module := range p.modules {
		err := module.Apply(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline[T]) ApplyDry(input T) error {
	for _, module := range p.modules {
		dryResult, err := module.ApplyDry(input)
		if err != nil {
			return err
		}

		if len(dryResult) > 0 {
			println(dryResult)
		}
	}
	return nil
}

func (p *Pipeline[T]) ApplyVerbose(input T) error {
	for _, module := range p.modules {
		dryResult, err := module.ApplyDry(input)
		if err != nil {
			return err
		}

		if len(dryResult) > 0 {
			println(dryResult)
		}

		err = module.Apply(input)
		if err != nil {
			return err
		}
	}
	return nil
}

type CreateDirStep[T PipelineConfig] struct {
	path string
	role string
}

func newCreateDirStep[T PipelineConfig](path string, role string) PipelineModule[T] {
	path = util.CleanPath(path)
	return &CreateDirStep[T]{
		path: path,
		role: role,
	}
}

func (s *CreateDirStep[T]) Apply(config T) error {
	if exists, _ := util.Exists(s.path); exists {
		if !config.ForceOperations() {
			return nil
		}

		if err := util.Delete(s.path); err != nil {
			return err
		}
	}

	return util.CreateDir(s.path)
}

func (s *CreateDirStep[T]) ApplyDry(config T) (string, error) {
	if exists, _ := util.Exists(s.path); !exists {
		return "Create " + s.role + " at " + s.path, nil
	}

	if config.ForceOperations() {
		return "Force recreate " + s.role + " at " + s.path + " (delete and create)", nil
	}

	return "Skip creation of " + s.role + " at " + s.path + " (already exists)", nil
}

type CreateFileStep[T PipelineConfig] struct {
	path    string
	content []byte
	role    string
}

func newCreateFileStep[T PipelineConfig](path string, content []byte, role string) PipelineModule[T] {
	path = util.CleanPath(path)
	return &CreateFileStep[T]{
		path:    path,
		content: content,
		role:    role,
	}
}

func (s *CreateFileStep[T]) Apply(config T) error {
	if exists, _ := util.Exists(s.path); exists {
		if !config.ForceOperations() {
			return nil
		}

		if err := util.Delete(s.path); err != nil {
			return err
		}
	}

	return util.WriteConfigFile(s.path, s.content)
}

func (s *CreateFileStep[T]) ApplyDry(config T) (string, error) {
	if exists, _ := util.Exists(s.path); !exists {
		return "Create " + s.role + " at " + s.path, nil
	}

	if config.ForceOperations() {
		return "Force recreate " + s.role + " at " + s.path + " (delete and create)", nil
	}

	return "Skip creation of " + s.role + " at " + s.path + " (already exists)", nil
}
