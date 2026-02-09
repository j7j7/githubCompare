package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetChangedFiles compares two references and returns all changed files
func GetChangedFiles(repoPath, startRef, endRef string) ([]FileChange, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve start reference
	startHash, err := repo.ResolveRevision(plumbing.Revision(startRef))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve start reference %s: %w", startRef, err)
	}

	// Resolve end reference
	endHash, err := repo.ResolveRevision(plumbing.Revision(endRef))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve end reference %s: %w", endRef, err)
	}

	// Get commit objects
	startCommit, err := repo.CommitObject(*startHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get start commit: %w", err)
	}

	endCommit, err := repo.CommitObject(*endHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get end commit: %w", err)
	}

	// Get trees
	startTree, err := startCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get start tree: %w", err)
	}

	endTree, err := endCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get end tree: %w", err)
	}

	// Get diff
	changes, err := object.DiffTree(startTree, endTree)
	if err != nil {
		return nil, fmt.Errorf("failed to diff trees: %w", err)
	}

	fileChanges := []FileChange{}

	for _, change := range changes {
		fileChange := FileChange{}

		// Determine change type
		if change.From.Name == "" {
			fileChange.ChangeType = "added"
			fileChange.Path = change.To.Name
		} else if change.To.Name == "" {
			fileChange.ChangeType = "deleted"
			fileChange.Path = change.From.Name
		} else if change.From.Name != change.To.Name {
			fileChange.ChangeType = "renamed"
			fileChange.Path = change.To.Name
			fileChange.OldPath = change.From.Name
		} else {
			fileChange.ChangeType = "modified"
			fileChange.Path = change.To.Name
		}

		fileChanges = append(fileChanges, fileChange)
	}

	return fileChanges, nil
}

// ValidateRefs validates that both references exist and start is ancestor of end
func ValidateRefs(repoPath, startRef, endRef string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve references
	startHash, err := repo.ResolveRevision(plumbing.Revision(startRef))
	if err != nil {
		return fmt.Errorf("start reference %s not found: %w", startRef, err)
	}

	endHash, err := repo.ResolveRevision(plumbing.Revision(endRef))
	if err != nil {
		return fmt.Errorf("end reference %s not found: %w", endRef, err)
	}

	// Check if start is ancestor of end
	startCommit, err := repo.CommitObject(*startHash)
	if err != nil {
		return fmt.Errorf("failed to get start commit: %w", err)
	}

	endCommit, err := repo.CommitObject(*endHash)
	if err != nil {
		return fmt.Errorf("failed to get end commit: %w", err)
	}

	// Check if start is ancestor of end
	isAncestor, err := endCommit.IsAncestor(startCommit)
	if err != nil {
		return fmt.Errorf("failed to check ancestry: %w", err)
	}

	if !isAncestor && startHash.String() != endHash.String() {
		return fmt.Errorf("start reference %s is not an ancestor of end reference %s", startRef, endRef)
	}

	return nil
}

// GetCommitHash returns the full hash for a reference
func GetCommitHash(repoPath, ref string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", fmt.Errorf("failed to resolve reference %s: %w", ref, err)
	}

	return hash.String(), nil
}
