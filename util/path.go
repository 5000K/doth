package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CleanPath tries to turn any path into a usable path (cleans up, expands ~ to home dir, and tries to turn it into an absolute path)
func CleanPath(relativePath string) string {
	relativePath = filepath.Clean(relativePath)

	expandedPath := expandHomeDir(relativePath)

	// try to get the absolute path, if it fails, return the cleaned relative path
	absPath, err := filepath.Abs(expandedPath)

	if err != nil {
		return expandedPath
	}

	return absPath
}

// ~username syntax is not supported.
func expandHomeDir(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	return filepath.Join(home, path[1:])
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func WriteConfigFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func Delete(path string) error {
	return os.RemoveAll(path)
}
