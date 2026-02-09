package utils

// RepoInfo contains parsed repository information
type RepoInfo struct {
	URL      string
	Name     string
	Owner    string
	Protocol string // "https" or "ssh"
}
