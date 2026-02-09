package display

import (
	"fmt"

	"github.com/githubCompare/internal/git"
)

// PrintChanges displays file changes in a formatted way
func PrintChanges(changes []git.FileChange) {
	if len(changes) == 0 {
		PrintWarning("No changes found")
		return
	}

	PrintSection(fmt.Sprintf("Changed Files (%d total)", len(changes)))
	
	// Group by change type
	added := []git.FileChange{}
	modified := []git.FileChange{}
	deleted := []git.FileChange{}
	renamed := []git.FileChange{}
	
	for _, change := range changes {
		switch change.ChangeType {
		case "added":
			added = append(added, change)
		case "modified":
			modified = append(modified, change)
		case "deleted":
			deleted = append(deleted, change)
		case "renamed":
			renamed = append(renamed, change)
		}
	}
	
	if len(added) > 0 {
		Added.Printf("\n  ➕ Added (%d):\n", len(added))
		for _, change := range added {
			File.Printf("      + %s\n", change.Path)
		}
	}
	
	if len(modified) > 0 {
		Modified.Printf("\n  ✏️  Modified (%d):\n", len(modified))
		for _, change := range modified {
			File.Printf("      ~ %s\n", change.Path)
		}
	}
	
	if len(renamed) > 0 {
		Renamed.Printf("\n  🔄 Renamed (%d):\n", len(renamed))
		for _, change := range renamed {
			File.Printf("      %s → %s\n", change.OldPath, change.Path)
		}
	}
	
	if len(deleted) > 0 {
		Deleted.Printf("\n  ➖ Deleted (%d):\n", len(deleted))
		for _, change := range deleted {
			File.Printf("      - %s\n", change.Path)
		}
	}
	
	fmt.Println()
}

// PrintSummary prints a summary of changes
func PrintSummary(startRef, endRef string, changeCount int) {
	PrintSection("Comparison Summary")
	fmt.Printf("  Start: ")
	Commit.Printf("%s\n", startRef)
	fmt.Printf("  End:   ")
	Commit.Printf("%s\n", endRef)
	fmt.Printf("  Files: ")
	Count.Printf("%d changed\n", changeCount)
	fmt.Println()
}
