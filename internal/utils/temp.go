package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateTempDir creates a temporary directory with the given prefix
func CreateTempDir(prefix string) (string, error) {
	return ioutil.TempDir("", prefix)
}

// CleanupTemp removes a temporary directory
func CleanupTemp(path string) error {
	if path == "" {
		return nil
	}
	return os.RemoveAll(path)
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetTempBase returns the base temporary directory path
func GetTempBase() string {
	return filepath.Join(os.TempDir(), "githubCompare")
}
