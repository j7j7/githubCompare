package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/githubCompare/internal/archive"
	"github.com/githubCompare/internal/display"
	"github.com/githubCompare/internal/git"
	"github.com/githubCompare/internal/interactive"
	"github.com/githubCompare/internal/utils"
)

func runCompare(cmd *cobra.Command, args []string) {
	// Parse repository URL
	repoInfo, err := utils.ParseRepoURL(repoURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing repository URL: %v\n", err)
		os.Exit(1)
	}

	// Create temp directory
	tempDir, err := utils.CreateTempDir("githubCompare-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp directory: %v\n", err)
		os.Exit(1)
	}

	// Cleanup temp directory unless --no-cleanup is set
	if !noCleanup {
		defer func() {
			if err := utils.CleanupTemp(tempDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to cleanup temp directory: %v\n", err)
			}
		}()
	}

	// Display header
	display.PrintHeader("GitHub Compare")
	fmt.Printf("\n")
	display.Info.Printf("Repository: %s\n", repoURL)
	if repoInfo.Name != "" {
		display.Info.Printf("Project: %s/%s\n", repoInfo.Owner, repoInfo.Name)
	}
	fmt.Println()

	display.PrintSection("Cloning Repository")
	fmt.Printf("  Cloning %s...\n", repoURL)

	// Clone repository
	cloneOpts := git.CloneOptions{
		URL:       repoURL,
		AuthToken: authToken,
		TempDir:   tempDir,
	}

	repoPath, err := git.CloneRepository(cloneOpts)
	if err != nil {
		display.PrintError(fmt.Sprintf("Failed to clone repository: %v", err))
		os.Exit(1)
	}

	display.PrintSuccess("Repository cloned successfully")

	// List branches
	display.PrintSection("Fetching Branches")
	branches, err := git.ListBranches(repoPath)
	if err != nil {
		display.PrintError(fmt.Sprintf("Failed to list branches: %v", err))
		os.Exit(1)
	}

	// If both start and end are provided, skip interactive selection
	var startCommit, endCommit string
	if startRef != "" && endRef != "" {
		startCommit = startRef
		endCommit = endRef
		display.Info.Printf("Using start reference: %s\n", startRef)
		display.Info.Printf("Using end reference: %s\n", endRef)
		fmt.Println()
	} else {
		// Select branch if not provided
		selectedBranch := endRef
		if selectedBranch == "" {
			selectedBranch, err = interactive.SelectBranch(branches)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error selecting branch: %v\n", err)
				os.Exit(1)
			}
		}

		// List commits for selected branch
		display.PrintSection(fmt.Sprintf("Fetching Commits for Branch '%s'", selectedBranch))
		commits, err := git.ListCommits(repoPath, selectedBranch, 100)
		if err != nil {
			display.PrintError(fmt.Sprintf("Failed to list commits: %v", err))
			os.Exit(1)
		}

		if len(commits) == 0 {
			display.PrintWarning(fmt.Sprintf("No commits found for branch %s", selectedBranch))
			os.Exit(1)
		}

		// Select start commit
		if startRef != "" {
			startCommit = startRef
			display.Info.Printf("Using start reference: %s\n", startRef)
		} else {
			startCommit, err = interactive.SelectCommit(commits, "Select START commit (older commit):")
			if err != nil {
				display.PrintError(fmt.Sprintf("Failed to select start commit: %v", err))
				os.Exit(1)
			}
		}

		// Select end commit
		if endRef != "" {
			endCommit = endRef
			display.Info.Printf("Using end reference: %s\n", endRef)
		} else {
			endCommit, err = interactive.SelectCommit(commits, "Select END commit (newer commit):")
			if err != nil {
				display.PrintError(fmt.Sprintf("Failed to select end commit: %v", err))
				os.Exit(1)
			}
		}
	}

	// Validate refs
	display.PrintSection("Validating References")
	if err := git.ValidateRefs(repoPath, startCommit, endCommit); err != nil {
		display.PrintError(fmt.Sprintf("Reference validation failed: %v", err))
		os.Exit(1)
	}
	display.PrintSuccess("References validated")

	// Get changed files
	display.PrintSection("Comparing Changes")
	fileChanges, err := git.GetChangedFiles(repoPath, startCommit, endCommit)
	if err != nil {
		display.PrintError(fmt.Sprintf("Failed to compare changes: %v", err))
		os.Exit(1)
	}

	if len(fileChanges) == 0 {
		display.PrintWarning("No files changed between the selected commits.")
		os.Exit(0)
	}

	// Display changes summary
	startShort := startCommit
	endShort := endCommit
	if len(startCommit) > 7 {
		startShort = startCommit[:7]
	}
	if len(endCommit) > 7 {
		endShort = endCommit[:7]
	}
	display.PrintSummary(startShort, endShort, len(fileChanges))
	display.PrintChanges(fileChanges)

	// Generate output path if not provided
	if outputPath == "" {
		outputPath = archive.GenerateOutputName(repoInfo.Name, startShort, endShort)
	}

	// Ensure output directory exists
	if err := archive.EnsureOutputDir(outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Convert git.FileChange to archive.FileChange
	archiveChanges := make([]archive.FileChange, len(fileChanges))
	for i, fc := range fileChanges {
		archiveChanges[i] = archive.FileChange{
			Path:       fc.Path,
			ChangeType: fc.ChangeType,
			OldPath:    fc.OldPath,
		}
	}

	// Create ZIP archive
	display.PrintSection("Creating Archive")
	fmt.Printf("  Output: %s\n", outputPath)
	
	if err := archive.CreateZipFromChanges(repoPath, archiveChanges, outputPath); err != nil {
		display.PrintError(fmt.Sprintf("Failed to create ZIP archive: %v", err))
		os.Exit(1)
	}

	// Get absolute path for display
	absPath, _ := filepath.Abs(outputPath)
	
	display.PrintHeader("Complete!")
	display.PrintSuccess(fmt.Sprintf("Archive created: %s", absPath))
	display.Count.Printf("  Changed files: %d\n", len(fileChanges))
	
	if noCleanup {
		display.Info.Printf("  Temp directory kept: %s\n", repoPath)
	}
	fmt.Println()
}
