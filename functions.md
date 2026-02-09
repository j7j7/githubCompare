# githubCompare - Functions & Components

## Core Packages Structure

### Package: `main`
**Purpose**: Entry point and CLI setup

#### Functions:
- `main()` - Application entry point, initializes Cobra CLI

---

### Package: `git`
**Purpose**: Git operations (clone, list, diff)

#### Functions:
- `CloneRepository(url, authToken string) (repoPath string, err error)`
  - Clones repository to temp directory
  - Returns path to cloned repo
  
- `ListBranches(repoPath string) ([]Branch, error)`
  - Lists all branches in repository
  - Returns slice of Branch structs
  
- `ListCommits(repoPath, ref string, limit int) ([]Commit, error)`
  - Lists commits for given branch/ref
  - Returns slice of Commit structs with hash, message, author, date
  
- `GetChangedFiles(repoPath, startRef, endRef string) ([]FileChange, error)`
  - Compares two refs and returns changed files
  - Returns FileChange structs (path, status: added/modified/deleted)
  
- `ValidateRefs(repoPath, startRef, endRef string) error`
  - Validates that both refs exist
  - Checks that start is ancestor of end

#### Types:
```go
type Branch struct {
    Name      string
    IsRemote  bool
    IsHead    bool
}

type Commit struct {
    Hash      string
    ShortHash string
    Message   string
    Author    string
    Date      time.Time
}

type FileChange struct {
    Path       string
    ChangeType string // "added", "modified", "deleted", "renamed"
    OldPath    string // for renames
}
```

---

### Package: `archive`
**Purpose**: ZIP file creation

#### Functions:
- `CreateZipFromChanges(repoPath string, changes []FileChange, outputPath string) error`
  - Creates ZIP archive with changed files
  - Preserves directory structure
  - Handles errors for missing files

- `GenerateOutputName(repoName, startRef, endRef string) string`
  - Generates meaningful ZIP filename
  - Format: `{repoName}_{startRef}_to_{endRef}_{timestamp}.zip`

---

### Package: `interactive`
**Purpose**: User interaction and prompts

#### Functions:
- `SelectBranch(branches []Branch) (string, error)`
  - Interactive branch selector
  - Returns selected branch name

- `SelectCommit(commits []Commit, prompt string) (string, error)`
  - Interactive commit selector
  - Returns selected commit hash

- `ConfirmAction(message string) (bool, error)`
  - Yes/No confirmation prompt

---

### Package: `utils`
**Purpose**: Helper utilities

#### Functions:
- `CreateTempDir(prefix string) (string, error)`
  - Creates temporary directory
  - Returns path

- `CleanupTemp(path string) error`
  - Removes temporary directory

- `ParseRepoURL(url string) (RepoInfo, error)`
  - Parses Git URL (HTTPS/SSH)
  - Extracts owner, repo name

- `ValidateAuth(url, token string) error`
  - Tests authentication credentials

#### Types:
```go
type RepoInfo struct {
    URL      string
    Name     string
    Owner    string
    Protocol string // "https" or "ssh"
}
```

---

### Package: `cmd`
**Purpose**: CLI commands (using Cobra)

#### Commands:
- `rootCmd` - Main command
- `compareCmd` - Compare command with all flags
- `versionCmd` - Version information

---

## Flow Diagram

```
User Input (repo URL)
    ↓
Parse & Validate URL
    ↓
Clone Repository → Temp Directory
    ↓
List Branches → User Selection
    ↓
List Commits (for selected branch) → User Selects Start
    ↓
List Commits → User Selects End
    ↓
Get Changed Files (diff start..end)
    ↓
Create ZIP Archive (changed files only)
    ↓
Cleanup Temp Directory
    ↓
Output Success Message
```

## Dependencies

### External Libraries
- `github.com/go-git/go-git/v5` - Git operations
- `github.com/spf13/cobra` - CLI framework
- `github.com/AlecAivazis/survey/v2` - Interactive prompts
- `github.com/schollz/progressbar/v3` - Progress bars (optional)

### Standard Library
- `archive/zip` - ZIP creation
- `os` - File operations
- `path/filepath` - Path handling
- `io/ioutil` - Temp directories
- `fmt`, `log` - Output and logging
