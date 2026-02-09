package git

import "time"

// Branch represents a Git branch
type Branch struct {
	Name         string
	IsRemote     bool
	IsHead       bool
	LastCommit   *Commit // Last commit on this branch
}

// Commit represents a Git commit
type Commit struct {
	Hash      string
	ShortHash string
	Message   string
	Author    string
	Date      time.Time
}

// FileChange represents a changed file
type FileChange struct {
	Path       string
	ChangeType string // "added", "modified", "deleted", "renamed"
	OldPath    string // for renames
}
