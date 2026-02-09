package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ListBranches lists all branches in the repository
func ListBranches(repoPath string) ([]Branch, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	branches := []Branch{}
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get local branches
	branchIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branch := Branch{
				Name:     ref.Name().Short(),
				IsRemote: false,
				IsHead:   ref.Hash() == head.Hash(),
			}
			
			// Get last commit for this branch
			var commit *object.Commit
			commit, err = repo.CommitObject(ref.Hash())
			if err == nil {
				message := strings.TrimSpace(commit.Message)
				// Get first line of commit message
				if idx := strings.Index(message, "\n"); idx > 0 {
					message = message[:idx]
				}
				// Truncate if too long
				if len(message) > 60 {
					message = message[:57] + "..."
				}
				
				branch.LastCommit = &Commit{
					Hash:      commit.Hash.String(),
					ShortHash: commit.Hash.String()[:7],
					Message:   message,
					Author:    commit.Author.Name,
					Date:      commit.Author.When,
				}
			}
			
			branches = append(branches, branch)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process branches: %w", err)
	}

	// Get remote branches
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	for _, remote := range remotes {
		refs, err := remote.List(&git.ListOptions{})
		if err != nil {
			continue // Skip if we can't list remote refs
		}

		for _, ref := range refs {
			if ref.Name().IsBranch() {
				branchName := strings.TrimPrefix(ref.Name().String(), "refs/remotes/"+remote.Config().Name+"/")
				// Check if we already have this branch locally
				exists := false
				for _, b := range branches {
					if b.Name == branchName {
						exists = true
						break
					}
				}
				if !exists {
					branch := Branch{
						Name:     branchName,
						IsRemote: true,
						IsHead:   false,
					}
					
					// Get last commit for remote branch
					var commit *object.Commit
					commit, err = repo.CommitObject(ref.Hash())
					if err == nil {
						message := strings.TrimSpace(commit.Message)
						// Get first line of commit message
						if idx := strings.Index(message, "\n"); idx > 0 {
							message = message[:idx]
						}
						// Truncate if too long
						if len(message) > 60 {
							message = message[:57] + "..."
						}
						
						branch.LastCommit = &Commit{
							Hash:      commit.Hash.String(),
							ShortHash: commit.Hash.String()[:7],
							Message:   message,
							Author:    commit.Author.Name,
							Date:      commit.Author.When,
						}
					}
					
					branches = append(branches, branch)
				}
			}
		}
	}

	return branches, nil
}
