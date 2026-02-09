package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	repoURL    string
	outputPath string
	startRef   string
	endRef     string
	authToken  string
	noCleanup  bool
)

var rootCmd = &cobra.Command{
	Use:   "githubCompare",
	Short: "Compare Git repository changes between two commits/branches",
	Long: `githubCompare is a CLI tool that allows you to compare changes between
two Git commits or branches and export only the changed files as a ZIP archive.

It supports both public and private repositories, and can work with HTTPS or SSH URLs.`,
	Run: runCompare,
}

func init() {
	rootCmd.Flags().StringVarP(&repoURL, "repo", "r", "", "Repository URL (required)")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output ZIP file path (optional, auto-generated if not provided)")
	rootCmd.Flags().StringVarP(&startRef, "start", "s", "", "Start commit/branch (optional, will prompt if not provided)")
	rootCmd.Flags().StringVarP(&endRef, "end", "e", "", "End commit/branch (optional, will prompt if not provided)")
	rootCmd.Flags().StringVar(&authToken, "auth-token", "", "Authentication token for private repos (HTTPS)")
	rootCmd.Flags().BoolVar(&noCleanup, "no-cleanup", false, "Keep temporary directory after execution")

	rootCmd.MarkFlagRequired("repo")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
