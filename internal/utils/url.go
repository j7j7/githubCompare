package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ParseRepoURL parses a Git repository URL and extracts information
func ParseRepoURL(url string) (*RepoInfo, error) {
	info := &RepoInfo{
		URL: url,
	}

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Determine protocol
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://") {
		info.Protocol = "ssh"
	} else {
		info.Protocol = "https"
	}

	// Parse HTTPS URL: https://github.com/owner/repo
	if info.Protocol == "https" {
		re := regexp.MustCompile(`https?://(?:[^@]+@)?([^/]+)/([^/]+)/([^/]+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) == 4 {
			info.Owner = matches[2]
			info.Name = matches[3]
			return info, nil
		}
	}

	// Parse SSH URL: git@github.com:owner/repo
	if info.Protocol == "ssh" {
		re := regexp.MustCompile(`(?:ssh://)?(?:git@)?([^:]+):([^/]+)/([^/]+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) == 4 {
			info.Owner = matches[2]
			info.Name = matches[3]
			return info, nil
		}
	}

	return nil, fmt.Errorf("unable to parse repository URL: %s", url)
}

// ValidateAuth tests if authentication is needed and available
func ValidateAuth(url string) bool {
	// For SSH, check if SSH key exists
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://") {
		// Check for SSH key in common locations
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		sshKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa")
		if _, err := os.Stat(sshKeyPath); err == nil {
			return true
		}
		sshKeyPath = filepath.Join(homeDir, ".ssh", "id_ed25519")
		if _, err := os.Stat(sshKeyPath); err == nil {
			return true
		}
		return false
	}
	// HTTPS repos can work without auth for public repos
	return true
}
