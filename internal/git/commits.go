package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ListCommits lists commits for a given reference
func ListCommits(repoPath, ref string, limit int) ([]Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve the reference
	hash, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve reference %s: %w", ref, err)
	}

	commits := []Commit{}
	count := 0

	// Iterate through commit history using Log
	commitIter, err := repo.Log(&git.LogOptions{From: *hash})
	if err != nil {
		return nil, fmt.Errorf("failed to create commit iterator: %w", err)
	}
	defer commitIter.Close()

	err = commitIter.ForEach(func(c *object.Commit) error {
		if limit > 0 && count >= limit {
			return fmt.Errorf("limit reached")
		}

		commit := Commit{
			Hash:      c.Hash.String(),
			ShortHash: c.Hash.String()[:7],
			Message:   strings.TrimSpace(c.Message),
			Author:    c.Author.Name,
			Date:      c.Author.When,
		}

		// Truncate message for display
		if len(commit.Message) > 80 {
			commit.Message = commit.Message[:77] + "..."
		}

		commits = append(commits, commit)
		count++
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}
