package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateOutputName generates a meaningful ZIP filename
func GenerateOutputName(repoName, startRef, endRef string) string {
	// Clean ref names for filename
	cleanStart := cleanRefForFilename(startRef)
	cleanEnd := cleanRefForFilename(endRef)

	// Truncate refs if too long
	if len(cleanStart) > 20 {
		cleanStart = cleanStart[:20]
	}
	if len(cleanEnd) > 20 {
		cleanEnd = cleanEnd[:20]
	}

	// Add timestamp
	timestamp := time.Now().Format("20060102_150405")

	filename := fmt.Sprintf("%s_%s_to_%s_%s.zip", repoName, cleanStart, cleanEnd, timestamp)
	return filename
}

// cleanRefForFilename removes invalid characters from ref name
func cleanRefForFilename(ref string) string {
	// Remove common invalid characters
	ref = strings.ReplaceAll(ref, "/", "_")
	ref = strings.ReplaceAll(ref, "\\", "_")
	ref = strings.ReplaceAll(ref, ":", "_")
	ref = strings.ReplaceAll(ref, "*", "_")
	ref = strings.ReplaceAll(ref, "?", "_")
	ref = strings.ReplaceAll(ref, "\"", "_")
	ref = strings.ReplaceAll(ref, "<", "_")
	ref = strings.ReplaceAll(ref, ">", "_")
	ref = strings.ReplaceAll(ref, "|", "_")

	// If it's a hash, use first 7 characters
	if len(ref) == 40 && isHexString(ref) {
		return ref[:7]
	}

	return ref
}

// isHexString checks if a string contains only hexadecimal characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// EnsureOutputDir ensures the output directory exists
func EnsureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if dir == "" || dir == "." {
		return nil // Current directory
	}
	return os.MkdirAll(dir, 0755)
}
