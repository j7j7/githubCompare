package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateZipFromChanges creates a ZIP archive containing only the changed files
func CreateZipFromChanges(repoPath string, changes []FileChange, outputPath string) error {
	// Create the ZIP file
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Track added files to avoid duplicates
	addedFiles := make(map[string]bool)

	for _, change := range changes {
		// Skip deleted files (or handle them differently if needed)
		if change.ChangeType == "deleted" {
			continue
		}

		// Skip if already added (for renames)
		if addedFiles[change.Path] {
			continue
		}

		filePath := filepath.Join(repoPath, change.Path)

		// Check if file exists
		info, err := os.Stat(filePath)
		if err != nil {
			// File might not exist if it was deleted in the end commit
			continue
		}

		// Skip directories
		if info.IsDir() {
			continue
		}

		// Add file to ZIP
		if err := addFileToZip(zipWriter, filePath, change.Path); err != nil {
			return fmt.Errorf("failed to add file %s to ZIP: %w", change.Path, err)
		}

		addedFiles[change.Path] = true
	}

	return nil
}

// addFileToZip adds a single file to the ZIP archive
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create ZIP header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Set the name in the ZIP (use forward slashes)
	header.Name = strings.ReplaceAll(zipPath, "\\", "/")
	header.Method = zip.Deflate

	// Create writer for this file
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy file contents
	_, err = io.Copy(writer, file)
	return err
}

// FileChange represents a file change (imported from git package)
type FileChange struct {
	Path       string
	ChangeType string
	OldPath    string
}
