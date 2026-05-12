package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigMap map[string]any

type PackageSource struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type PackageConfig struct {
	/**
	PackageSources is a list of package managers available on the target system.
	Use {package} as a placeholder for the package name (and version if needed) in the command.
	If multiple sources are valid for a package, the first one found will be used.

	Example configuration:

	packageSources:
	  - name: "pacman"
	    command: "sudo pacman -S {package}"
	  - name: "paru"
	    command: "paru -S {package}"
	  - name: "go"
	    command: "go install {package}"
	  - name: "npm"
	    command: "npm install -g {package}"

	*/
	PackageSources []PackageSource `yaml:"packageSources"`
}

func (configMap ConfigMap) GetPackageConfig() (*PackageConfig, error) {
	var config PackageConfig
	if sources, ok := configMap["packageSources"].([]any); ok {
		for _, source := range sources {
			if sourceMap, ok := source.(map[string]any); ok {
				name, _ := sourceMap["name"].(string)
				command, _ := sourceMap["command"].(string)
				config.PackageSources = append(config.PackageSources, PackageSource{
					Name:    name,
					Command: command,
				})
			}
		}
	}
	return &config, nil
}

func LoadConfigFiles(paths []string) (ConfigMap, error) {
	if len(paths) == 0 {
		return nil, errors.New("no config paths provided")
	}

	result, err := loadFile(paths[0])
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", paths[0], err)
	}

	for _, path := range paths[1:] {
		next, err := loadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", path, err)
		}
		result = mergeConfigs(result, next)
	}

	return result, nil
}

func loadFile(path string) (ConfigMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseConfig(path, data)
}

func parseConfig(path string, data []byte) (ConfigMap, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return parseYAML(data)
	case ".json":
		return parseJSON(data)
	default:
		return nil, fmt.Errorf("unsupported file extension %q", ext)
	}
}

func parseYAML(data []byte) (ConfigMap, error) {
	var result ConfigMap
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseJSON(data []byte) (ConfigMap, error) {
	var result ConfigMap
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func mergeConfigs(dst, src ConfigMap) ConfigMap {
	result := make(ConfigMap, len(dst))
	for k, v := range dst {
		result[k] = v
	}
	for k, srcVal := range src {
		dstVal, exists := result[k]
		if !exists {
			result[k] = srcVal
			continue
		}
		result[k] = mergeValues(dstVal, srcVal)
	}
	return result
}

func mergeValues(dst, src any) any {
	dstMap, dstIsMap := dst.(ConfigMap)
	srcMap, srcIsMap := src.(ConfigMap)
	if dstIsMap && srcIsMap {
		return mergeConfigs(dstMap, srcMap)
	}

	dstSlice, dstIsSlice := dst.([]any)
	srcSlice, srcIsSlice := src.([]any)
	if dstIsSlice && srcIsSlice {
		return mergeSlices(dstSlice, srcSlice)
	}

	return src
}

func mergeSlices(a, b []any) []any {
	seen := make(map[any]struct{}, len(a)+len(b))
	result := make([]any, 0, len(a)+len(b))

	for _, v := range append(a, b...) {
		if !isComparable(v) {
			result = append(result, v)
			continue
		}
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

func isComparable(v any) bool {
	switch v.(type) {
	case ConfigMap, []any:
		return false
	default:
		return true
	}
}
