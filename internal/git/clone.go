package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// CloneOptions contains options for cloning a repository
type CloneOptions struct {
	URL       string
	AuthToken string
	TempDir   string
}

// CloneRepository clones a Git repository to a temporary directory
func CloneRepository(opts CloneOptions) (string, error) {
	// Create temp directory for this clone
	clonePath := filepath.Join(opts.TempDir, "repo")
	if err := os.MkdirAll(clonePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Configure clone options - fetch all branches
	cloneOpts := &git.CloneOptions{
		URL:               opts.URL,
		Progress:         os.Stdout,
		SingleBranch:     false,
		Depth:            0, // Full clone to get all branches
		RecurseSubmodules: git.NoRecurseSubmodules,
	}

	// Set up authentication
	auth, err := getAuth(opts.URL, opts.AuthToken)
	if err != nil {
		return "", fmt.Errorf("failed to set up authentication: %w", err)
	}
	if auth != nil {
		cloneOpts.Auth = auth
	}

	// Clone the repository
	repo, err := git.PlainClone(clonePath, false, cloneOpts)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Fetch all remote branches
	remotes, err := repo.Remotes()
	if err == nil && len(remotes) > 0 {
		remote := remotes[0]
		err = remote.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{"refs/heads/*:refs/remotes/origin/*"},
		})
		// Ignore fetch errors - branches might already be fetched
	}

	return clonePath, nil
}

// getAuth returns the appropriate authentication method
func getAuth(url, authToken string) (transport.AuthMethod, error) {
	// SSH URL
	if isSSHURL(url) {
		return getSSHAuth()
	}

	// HTTPS URL with token
	if authToken != "" {
		return &http.BasicAuth{
			Username: "token", // GitHub requires non-empty username
			Password: authToken,
		}, nil
	}

	// No auth needed for public repos
	return nil, nil
}

// isSSHURL checks if the URL is an SSH URL
func isSSHURL(url string) bool {
	return len(url) > 4 && (url[:4] == "git@" || url[:6] == "ssh://")
}

// getSSHAuth returns SSH authentication
func getSSHAuth() (transport.AuthMethod, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshKeyPath := filepath.Join(homeDir, ".ssh", "id_ed25519")
	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		sshKeyPath = filepath.Join(homeDir, ".ssh", "id_rsa")
		if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("no SSH key found")
		}
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load SSH key: %w", err)
	}

	// Set up known hosts
	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	if _, err := os.Stat(knownHostsPath); err == nil {
		callback, err := knownhosts.New(knownHostsPath)
		if err == nil {
			publicKeys.HostKeyCallback = callback
		}
	}

	return publicKeys, nil
}
