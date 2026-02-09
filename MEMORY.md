# githubCompare - Project Memory

## Project Overview
Command-line application to compare changes between two Git commits/branches and export only the changed files as a zip archive.

## Language & Technology Stack
- **Language**: Go (Golang)
- **Reason**: Cross-platform compilation, single binary output, excellent Git libraries, no runtime dependencies
- **Key Libraries**:
  - `go-git/go-git` - Pure Go implementation for Git operations
  - `archive/zip` - Standard library for ZIP creation
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/AlecAivazis/survey/v2` - Interactive prompts

## Key Requirements
1. Cross-platform support (OSX, Linux, Windows)
2. Works with public and private repositories
3. Lists branches and commits
4. Interactive selection of start and end points
5. Clones repo to temporary folder
6. Compares changes between two points
7. Creates ZIP with only changed files (preserving structure)
8. Git must be pre-installed on system

## Current Status
- ✅ Project structure created
- ✅ All core functionality implemented
- ✅ Unit tests created and passing
- ✅ Integration tests created
- ✅ Application builds successfully
- ✅ Binary verified and working
- ✅ Cross-platform ready (Go compiles to all platforms)

## Key Decisions
- Using Go-git library for Git operations (pure Go, no external git dependency)
- Temp folder handling with automatic cleanup
- ZIP archive will preserve directory structure
- Support for authentication (SSH keys, tokens for private repos)
