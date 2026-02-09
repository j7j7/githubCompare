package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/githubCompare/internal/display"
	"github.com/githubCompare/internal/git"
)

// SelectBranch prompts the user to select a branch
func SelectBranch(branches []git.Branch) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("no branches available")
	}

	// Display branches nicely
	display.PrintBranches(branches)

	options := make([]string, len(branches))
	for i, branch := range branches {
		options[i] = display.FormatBranchOption(branch, i)
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select a branch:",
		Options: options,
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", fmt.Errorf("failed to select branch: %w", err)
	}

	// Extract branch name (get the branch name from the formatted option)
	parts := strings.Fields(selected)
	if len(parts) >= 2 {
		// Find the branch name (skip index, arrow if present)
		for i, part := range parts {
			if i > 0 && !strings.HasSuffix(part, ".") && part != "→" && part != "(remote)" && part != "(current" && part != "HEAD)" {
				return part, nil
			}
		}
	}

	return "", fmt.Errorf("failed to parse selected branch")
}

// SelectCommit prompts the user to select a commit
func SelectCommit(commits []git.Commit, promptText string) (string, error) {
	if len(commits) == 0 {
		return "", fmt.Errorf("no commits available")
	}

	// Display commits nicely
	display.PrintCommits(commits, 20)

	options := make([]string, len(commits))
	for i, commit := range commits {
		options[i] = display.FormatCommitOption(commit, i)
	}

	var selected string
	prompt := &survey.Select{
		Message: promptText,
		Options: options,
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	// Extract commit hash (first field after the number)
	parts := strings.Fields(selected)
	if len(parts) >= 2 {
		hash := parts[1] // Skip the index number
		// Find full hash
		for _, commit := range commits {
			if commit.ShortHash == hash {
				return commit.Hash, nil
			}
		}
		return hash, nil
	}

	return "", fmt.Errorf("failed to parse selected commit")
}
