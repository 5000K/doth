package model

import (
	"os"
	"os/exec"

	"github.com/5000K/doth/template"
	"github.com/5000K/doth/util"
)

type PipelineConfig interface {
	ForceOperations() bool
	Autoconfirm() bool
}

func CreatePipelineConfig(force bool, autoconfirm bool) PipelineConfig {
	return &pipelineConfigImpl{
		force:       force,
		autoconfirm: autoconfirm,
	}
}

type pipelineConfigImpl struct {
	force       bool
	autoconfirm bool
}

func (c *pipelineConfigImpl) ForceOperations() bool {
	return c.force
}

func (c *pipelineConfigImpl) Autoconfirm() bool {
	return c.autoconfirm
}

// PipelineModule represents a single operation that can be applied, or printed out in a dry run.
type PipelineModule interface {
	Apply(input PipelineConfig) error
	ApplyDry(input PipelineConfig) (string, error)
}

// Pipeline represents a planned sequence of small operations. It can apply them, or print out what it would do (dry run), or print out what it is doing as it does it (verbose).
type Pipeline struct {
	modules []PipelineModule
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		modules: make([]PipelineModule, 0),
	}
}

func (p *Pipeline) AddModule(module PipelineModule) *Pipeline {
	p.modules = append(p.modules, module)
	return p
}

func (p *Pipeline) Run(dry bool, verbose bool, config PipelineConfig) error {
	if dry {
		return p.ApplyDry(config)
	} else if verbose {
		return p.ApplyVerbose(config)
	} else {
		return p.Apply(config)
	}
}

func (p *Pipeline) Apply(input PipelineConfig) error {
	for _, module := range p.modules {
		err := module.Apply(input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) ApplyDry(input PipelineConfig) error {
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

func (p *Pipeline) ApplyVerbose(input PipelineConfig) error {
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

type CreateDirStep struct {
	path string
	role string
}

func NewCreateDirStep(path string, role string) PipelineModule {
	path = util.CleanPath(path)
	return &CreateDirStep{
		path: path,
		role: role,
	}
}

func (s *CreateDirStep) Apply(config PipelineConfig) error {
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

func (s *CreateDirStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.path); !exists {
		return "Create " + s.role + " at " + s.path, nil
	}

	if config.ForceOperations() {
		return "Force recreate " + s.role + " at " + s.path + " (delete and create)", nil
	}

	return "Skip creation of " + s.role + " at " + s.path + " (already exists)", nil
}

type CreateDirIfNotExistsStep struct {
	path string
	role string
}

func NewCreateDirIfNotExistsStep(path string, role string) PipelineModule {
	path = util.CleanPath(path)
	return &CreateDirIfNotExistsStep{
		path: path,
		role: role,
	}
}

func (s *CreateDirIfNotExistsStep) Apply(config PipelineConfig) error {
	if exists, _ := util.Exists(s.path); exists {
		return nil
	}

	return util.CreateDir(s.path)
}

func (s *CreateDirIfNotExistsStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.path); exists {
		return "Skip creation of " + s.role + " at " + s.path + " (already exists)", nil
	}

	return "Create " + s.role + " at " + s.path, nil
}

type CreateFileStep struct {
	path        string
	content     []byte
	role        string
	permissions os.FileMode
}

func NewCreateFileStep(path string, content []byte, role string) PipelineModule {
	path = util.CleanPath(path)
	return &CreateFileStep{
		path:        path,
		content:     content,
		role:        role,
		permissions: 0644,
	}
}

func NewCreateFileStepWithPermissions(path string, content []byte, role string, permissions os.FileMode) PipelineModule {
	path = util.CleanPath(path)
	return &CreateFileStep{
		path:        path,
		content:     content,
		role:        role,
		permissions: permissions,
	}
}

func (s *CreateFileStep) Apply(config PipelineConfig) error {
	if exists, _ := util.Exists(s.path); exists {
		if !config.ForceOperations() {
			return nil
		}

		if err := util.Delete(s.path); err != nil {
			return err
		}
	}

	return os.WriteFile(s.path, s.content, s.permissions)
}

func (s *CreateFileStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.path); !exists {
		return "Create " + s.role + " at " + s.path, nil
	}

	if config.ForceOperations() {
		return "Force recreate " + s.role + " at " + s.path + " (delete and create)", nil
	}

	return "Skip creation of " + s.role + " at " + s.path + " (already exists)", nil
}

type CopyFileStep struct {
	src string
	dst string
}

func NewCopyFileStep(src string, dst string) PipelineModule {
	return &CopyFileStep{
		src: util.CleanPath(src),
		dst: util.CleanPath(dst),
	}
}

func (s *CopyFileStep) Apply(config PipelineConfig) error {
	if exists, _ := util.Exists(s.dst); exists {
		if err := util.Delete(s.dst); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(s.src)
	if err != nil {
		return err
	}
	return util.WriteConfigFile(s.dst, data)
}

func (s *CopyFileStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.dst); exists {
		return "Replace file at " + s.dst + " with copy of " + s.src, nil
	}
	return "Copy " + s.src + " to " + s.dst, nil
}

type CreateSymlinkStep struct {
	src string
	dst string
}

func NewCreateSymlinkStep(src string, dst string) PipelineModule {
	return &CreateSymlinkStep{
		src: util.CleanPath(src),
		dst: util.CleanPath(dst),
	}
}

func (s *CreateSymlinkStep) Apply(config PipelineConfig) error {
	if exists, _ := util.Exists(s.dst); exists {
		if err := util.Delete(s.dst); err != nil {
			return err
		}
	}
	return os.Symlink(s.src, s.dst)
}

func (s *CreateSymlinkStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.dst); exists {
		return "Replace " + s.dst + " with symlink to " + s.src, nil
	}
	return "Create symlink " + s.dst + " to " + s.src, nil
}

type RenderFileStep struct {
	templatePath string
	configMap    ConfigMap
	dst          string
}

func NewRenderFileStep(templatePath string, configMap ConfigMap, dst string) PipelineModule {
	return &RenderFileStep{
		templatePath: util.CleanPath(templatePath),
		configMap:    configMap,
		dst:          util.CleanPath(dst),
	}
}

func (s *RenderFileStep) Apply(config PipelineConfig) error {
	if exists, _ := util.Exists(s.dst); exists {
		if err := util.Delete(s.dst); err != nil {
			return err
		}
	}

	tmplData, err := os.ReadFile(s.templatePath)
	if err != nil {
		return err
	}

	rendered, err := template.RenderTemplate(string(tmplData), s.configMap)
	if err != nil {
		return err
	}

	return util.WriteConfigFile(s.dst, []byte(rendered))
}

func (s *RenderFileStep) ApplyDry(config PipelineConfig) (string, error) {
	if exists, _ := util.Exists(s.dst); exists {
		return "Replace " + s.dst + " with rendered output of template " + s.templatePath, nil
	}
	return "Render template " + s.templatePath + " to " + s.dst, nil
}

type ConfirmStep struct {
	message string
}

func NewConfirmStep(message string) PipelineModule {
	return &ConfirmStep{
		message: message,
	}
}

func (s *ConfirmStep) Apply(config PipelineConfig) error {
	if config.Autoconfirm() {
		return nil
	}

	return util.ConfirmAction(s.message)
}

func (s *ConfirmStep) ApplyDry(config PipelineConfig) (string, error) {
	if config.Autoconfirm() {
		return "Auto-confirm enabled:\n\t" + s.message + " (y/N) y", nil
	}
	return "Ask for confirmation: " + s.message, nil
}

type LogStep struct {
	message string
	dryOnly bool
}

func NewLogStep(message string, dryOnly bool) PipelineModule {
	return &LogStep{
		message: message,
		dryOnly: dryOnly,
	}
}

func (s *LogStep) Apply(config PipelineConfig) error {
	if s.dryOnly {
		return nil
	}
	println(s.message)
	return nil
}

func (s *LogStep) ApplyDry(config PipelineConfig) (string, error) {
	return s.message, nil
}

type ExecuteShellCommandStep struct {
	command string
}

func NewExecuteShellCommandStep(command string) PipelineModule {
	return &ExecuteShellCommandStep{
		command: command,
	}
}

func (s *ExecuteShellCommandStep) Apply(config PipelineConfig) error {
	return exec.Command("/bin/sh", "-c", s.command).Run()
}

func (s *ExecuteShellCommandStep) ApplyDry(config PipelineConfig) (string, error) {
	return "Run command: " + s.command, nil
}
