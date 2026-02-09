package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/githubCompare/internal/archive"
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

	fmt.Printf("Cloning repository %s...\n", repoURL)

	// Clone repository
	cloneOpts := git.CloneOptions{
		URL:       repoURL,
		AuthToken: authToken,
		TempDir:   tempDir,
	}

	repoPath, err := git.CloneRepository(cloneOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repository cloned successfully to %s\n", repoPath)

	// List branches
	fmt.Println("\nFetching branches...")
	branches, err := git.ListBranches(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing branches: %v\n", err)
		os.Exit(1)
	}

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
	fmt.Printf("\nFetching commits for branch '%s'...\n", selectedBranch)
	commits, err := git.ListCommits(repoPath, selectedBranch, 100)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing commits: %v\n", err)
		os.Exit(1)
	}

	if len(commits) == 0 {
		fmt.Fprintf(os.Stderr, "No commits found for branch %s\n", selectedBranch)
		os.Exit(1)
	}

	// Select start commit
	var startCommit string
	if startRef != "" {
		startCommit = startRef
	} else {
		startCommit, err = interactive.SelectCommit(commits, "Select start commit:")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting start commit: %v\n", err)
			os.Exit(1)
		}
	}

	// Select end commit
	var endCommit string
	if endRef != "" {
		endCommit = endRef
	} else {
		endCommit, err = interactive.SelectCommit(commits, "Select end commit:")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting end commit: %v\n", err)
			os.Exit(1)
		}
	}

	// Validate refs
	fmt.Println("\nValidating references...")
	if err := git.ValidateRefs(repoPath, startCommit, endCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Error validating references: %v\n", err)
		os.Exit(1)
	}

	// Get changed files
	fmt.Println("\nComparing changes...")
	fileChanges, err := git.GetChangedFiles(repoPath, startCommit, endCommit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting changed files: %v\n", err)
		os.Exit(1)
	}

	if len(fileChanges) == 0 {
		fmt.Println("No files changed between the selected commits.")
		os.Exit(0)
	}

	fmt.Printf("Found %d changed file(s)\n", len(fileChanges))

	// Generate output path if not provided
	if outputPath == "" {
		outputPath = archive.GenerateOutputName(repoInfo.Name, startCommit[:7], endCommit[:7])
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
	fmt.Printf("\nCreating ZIP archive: %s\n", outputPath)
	if err := archive.CreateZipFromChanges(repoPath, archiveChanges, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating ZIP archive: %v\n", err)
		os.Exit(1)
	}

	// Get absolute path for display
	absPath, _ := filepath.Abs(outputPath)
	fmt.Printf("\n✓ Successfully created archive: %s\n", absPath)
	fmt.Printf("  Changed files: %d\n", len(fileChanges))

	if noCleanup {
		fmt.Printf("  Temp directory kept: %s\n", repoPath)
	}
}
