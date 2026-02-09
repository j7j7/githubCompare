package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// ResolveRef tries to resolve a reference using multiple formats
func ResolveRef(repo *git.Repository, ref string) (*plumbing.Hash, error) {
	// Try different ref formats
	refFormats := []string{
		ref,                                    // As-is
		"refs/heads/" + ref,                    // Local branch
		"refs/remotes/origin/" + ref,           // Remote branch
		"origin/" + ref,                        // Remote branch shorthand
		"refs/tags/" + ref,                    // Tag
	}

	for _, refFormat := range refFormats {
		hash, err := repo.ResolveRevision(plumbing.Revision(refFormat))
		if err == nil {
			return hash, nil
		}
	}

	return nil, fmt.Errorf("reference not found: %s", ref)
}
