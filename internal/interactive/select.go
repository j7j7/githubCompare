package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/githubCompare/internal/git"
)

// SelectBranch prompts the user to select a branch
func SelectBranch(branches []git.Branch) (string, error) {
	if len(branches) == 0 {
		return "", fmt.Errorf("no branches available")
	}

	options := make([]string, len(branches))
	for i, branch := range branches {
		display := branch.Name
		if branch.IsHead {
			display += " (HEAD)"
		}
		if branch.IsRemote {
			display += " (remote)"
		}
		options[i] = display
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select a branch:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", fmt.Errorf("failed to select branch: %w", err)
	}

	// Extract branch name (remove display suffix)
	branchName := selected
	if idx := strings.Index(selected, " ("); idx != -1 {
		branchName = selected[:idx]
	}

	return branchName, nil
}

// SelectCommit prompts the user to select a commit
func SelectCommit(commits []git.Commit, promptText string) (string, error) {
	if len(commits) == 0 {
		return "", fmt.Errorf("no commits available")
	}

	options := make([]string, len(commits))
	for i, commit := range commits {
		dateStr := commit.Date.Format("2006-01-02 15:04")
		options[i] = fmt.Sprintf("%s - %s - %s", commit.ShortHash, dateStr, commit.Message)
	}

	var selected string
	prompt := &survey.Select{
		Message: promptText,
		Options: options,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return "", fmt.Errorf("failed to select commit: %w", err)
	}

	// Extract commit hash
	hash := strings.Split(selected, " - ")[0]

	// Find full hash
	for _, commit := range commits {
		if commit.ShortHash == hash {
			return commit.Hash, nil
		}
	}

	return hash, nil
}
