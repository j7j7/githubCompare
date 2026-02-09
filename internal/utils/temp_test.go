package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTempDir(t *testing.T) {
	dir, err := CreateTempDir("test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("Temp directory was not created: %s", dir)
	}

	// Cleanup
	defer CleanupTemp(dir)
}

func TestCleanupTemp(t *testing.T) {
	dir, err := CreateTempDir("test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a file in the temp directory
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Cleanup
	if err := CleanupTemp(dir); err != nil {
		t.Fatalf("Failed to cleanup temp directory: %v", err)
	}

	// Verify directory is gone
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("Temp directory was not removed: %s", dir)
	}
}

func TestEnsureDir(t *testing.T) {
	dir, err := CreateTempDir("test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer CleanupTemp(dir)

	testDir := filepath.Join(dir, "nested", "path")
	if err := EnsureDir(testDir); err != nil {
		t.Fatalf("Failed to ensure directory: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", testDir)
	}
}
