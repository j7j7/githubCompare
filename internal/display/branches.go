package display

import (
	"fmt"
	"strings"

	"github.com/githubCompare/internal/git"
)

// PrintBranches displays branches in a formatted way
func PrintBranches(branches []git.Branch) {
	if len(branches) == 0 {
		PrintWarning("No branches found")
		return
	}

	PrintSection("Available Branches")
	
	localBranches := []git.Branch{}
	remoteBranches := []git.Branch{}
	
	for _, branch := range branches {
		if branch.IsRemote {
			remoteBranches = append(remoteBranches, branch)
		} else {
			localBranches = append(localBranches, branch)
		}
	}
	
	if len(localBranches) > 0 {
		fmt.Println("\n  Local branches:")
		for i, branch := range localBranches {
			marker := "  "
			if branch.IsHead {
				marker = "→ "
			}
			
			fmt.Printf("    %d. %s", i+1, marker)
			Branch.Printf("%s", branch.Name)
			
			if branch.IsHead {
				fmt.Printf(" (current HEAD)")
			}
			
			if branch.LastCommit != nil {
				timeAgo := timeAgoString(branch.LastCommit.Date)
				fmt.Printf("\n       ")
				Commit.Printf("%s", branch.LastCommit.ShortHash)
				fmt.Printf(" - %s - %s", timeAgo, branch.LastCommit.Author)
				fmt.Printf("\n       %s\n", branch.LastCommit.Message)
			} else {
				fmt.Println()
			}
		}
	}
	
	if len(remoteBranches) > 0 {
		fmt.Println("\n  Remote branches:")
		for i, branch := range remoteBranches {
			fmt.Printf("    %d. ", i+1)
			Branch.Printf("%s", branch.Name)
			fmt.Printf(" (remote)")
			
			if branch.LastCommit != nil {
				timeAgo := timeAgoString(branch.LastCommit.Date)
				fmt.Printf("\n       ")
				Commit.Printf("%s", branch.LastCommit.ShortHash)
				fmt.Printf(" - %s - %s", timeAgo, branch.LastCommit.Author)
				fmt.Printf("\n       %s\n", branch.LastCommit.Message)
			} else {
				fmt.Println()
			}
		}
	}
	
	fmt.Println()
}

// FormatBranchOption formats a branch for selection
func FormatBranchOption(branch git.Branch, index int) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("%d.", index+1))
	
	if branch.IsHead {
		parts = append(parts, "→")
	}
	
	parts = append(parts, branch.Name)
	
	if branch.IsRemote {
		parts = append(parts, "(remote)")
	}
	
	// Add last commit info if available
	if branch.LastCommit != nil {
		timeAgo := timeAgoString(branch.LastCommit.Date)
		message := branch.LastCommit.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}
		parts = append(parts, fmt.Sprintf("- %s (%s): %s", branch.LastCommit.ShortHash, timeAgo, message))
	}
	
	return strings.Join(parts, " ")
}
