package display

import (
	"fmt"
	"time"

	"github.com/githubCompare/internal/git"
)

// PrintCommits displays commits in a formatted way
func PrintCommits(commits []git.Commit, limit int) {
	if len(commits) == 0 {
		PrintWarning("No commits found")
		return
	}

	displayCount := len(commits)
	if limit > 0 && limit < displayCount {
		displayCount = limit
	}

	PrintSection(fmt.Sprintf("Recent Commits (showing %d of %d)", displayCount, len(commits)))
	
	for i := 0; i < displayCount; i++ {
		commit := commits[i]
		timeAgo := timeAgoString(commit.Date)
		dateStr := commit.Date.Format("2006-01-02 15:04")
		
		fmt.Printf("  %d. ", i+1)
		Commit.Printf("%s", commit.ShortHash)
		fmt.Printf(" - %s (%s)", timeAgo, dateStr)
		fmt.Printf(" - %s", commit.Author)
		
		// Show commit message more prominently
		fmt.Printf("\n     ")
		Info.Printf("%s\n", commit.Message)
		if i < displayCount-1 {
			fmt.Println() // Only add blank line between items, not after last
		}
	}
}

// FormatCommitOption formats a commit for selection
func FormatCommitOption(commit git.Commit, index int) string {
	timeAgo := timeAgoString(commit.Date)
	message := commit.Message
	if len(message) > 55 {
		message = message[:52] + "..."
	}
	return fmt.Sprintf("%d. %s - %s - %s - %s", 
		index+1, commit.ShortHash, timeAgo, commit.Author, message)
}

// timeAgoString returns a human-readable time ago string
func timeAgoString(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (7 * 24))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else {
		return t.Format("2006-01-02")
	}
}
