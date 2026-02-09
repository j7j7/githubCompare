# githubCompare

A command-line tool to compare changes between two Git commits or branches and export only the changed files as a ZIP archive.

## Features

- ✅ Cross-platform support (Linux, macOS, Windows)
- ✅ Works with public and private repositories
- ✅ **Beautiful interactive interface** with colors and intuitive prompts
- ✅ **Command-line mode** for automation (non-interactive)
- ✅ Interactive branch and commit selection with formatted display
- ✅ Supports both HTTPS and SSH URLs
- ✅ Creates ZIP archives with only changed files
- ✅ Preserves directory structure
- ✅ Automatic cleanup of temporary files
- ✅ Color-coded change types (added/modified/deleted/renamed)

## Installation

### Prerequisites

- Go 1.21 or later
- Git (already installed on your system)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd githubCompare

# Download dependencies
go mod download

# Build for your platform
make build

# Or build for all platforms
make build-all
```

The binaries will be in the `dist/` directory.

## Usage

### Interactive Mode (Recommended)

```bash
githubCompare --repo https://github.com/owner/repo
```

This will:
1. **Display a beautiful header** with repository information
2. **Clone the repository** with progress indication
3. **Show formatted branch list** and prompt you to select one
4. **Display commits** with author, date, and messages in an easy-to-read format
5. **Show a summary** of changes with color-coded file types
6. **Create a ZIP file** with only the changed files

The interface uses colors to make it easy to understand:
- 🟢 **Green** for added files and success messages
- 🟡 **Yellow** for modified files and commits
- 🔴 **Red** for deleted files and errors
- 🟣 **Magenta** for renamed files and branches
- 🔵 **Blue** for information and sections

### Command-Line Mode (Non-Interactive)

For automation or when you know the exact refs:

```bash
# Specify start and end commits directly (skips all prompts)
githubCompare --repo https://github.com/owner/repo \
  --start abc1234 \
  --end def5678 \
  --output changes.zip

# Compare branches
githubCompare --repo https://github.com/owner/repo \
  --start main \
  --end feature-branch

# Mix commits and branches
githubCompare --repo https://github.com/owner/repo \
  --start main \
  --end abc1234
```

### Authentication

```bash
# Use with private repository (HTTPS)
githubCompare --repo https://github.com/owner/private-repo \
  --auth-token YOUR_GITHUB_TOKEN

# Use SSH URL (uses your SSH keys automatically)
githubCompare --repo git@github.com:owner/repo.git
```

### Other Options

```bash
# Keep temporary directory for inspection
githubCompare --repo https://github.com/owner/repo --no-cleanup

# Specify custom output path
githubCompare --repo https://github.com/owner/repo \
  --start main --end feature \
  --output /path/to/custom-output.zip
```

### Command Line Options

- `--repo, -r` - Repository URL (required)
- `--output, -o` - Output ZIP file path (optional, auto-generated if not provided)
- `--start, -s` - Start commit/branch (optional, will prompt if not provided)
- `--end, -e` - End commit/branch (optional, will prompt if not provided)
- `--auth-token` - Authentication token for private repos (HTTPS)
- `--no-cleanup` - Keep temporary directory after execution

## Examples

### Compare two commits

```bash
githubCompare --repo https://github.com/microsoft/vscode \
  --start v1.80.0 \
  --end v1.81.0 \
  --output vscode-changes.zip
```

### Interactive mode

```bash
githubCompare --repo https://github.com/owner/repo
# Follow the prompts to select branch and commits
```

## Development

### Running Tests

```bash
make test
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make verify
```

## How It Works

1. **Clone**: The tool clones the repository to a temporary directory
2. **List**: It fetches all branches and commits
3. **Select**: You interactively select the start and end points
4. **Compare**: It uses Git's diff functionality to find changed files
5. **Archive**: It creates a ZIP file containing only the changed files
6. **Cleanup**: Temporary files are automatically removed (unless `--no-cleanup` is used)

## Authentication

### Public Repositories
No authentication needed for public repositories.

### Private Repositories

**HTTPS**: Use the `--auth-token` flag with a GitHub Personal Access Token:
```bash
githubCompare --repo https://github.com/owner/private-repo --auth-token ghp_xxxxx
```

**SSH**: The tool automatically uses your SSH keys from `~/.ssh/`:
- `~/.ssh/id_ed25519` (preferred)
- `~/.ssh/id_rsa` (fallback)

## Output

The ZIP file contains:
- Only files that changed between the two commits
- Preserved directory structure
- Files are stored with forward slashes (works on all platforms)

The output filename format is:
```
{repo-name}_{start-ref}_to_{end-ref}_{timestamp}.zip
```

Example: `vscode_abc1234_to_def5678_20260109_143022.zip`

## Troubleshooting

### Authentication Errors
- For HTTPS: Ensure your token has `repo` scope
- For SSH: Ensure your SSH key is added to your GitHub account

### Clone Failures
- Check your internet connection
- Verify the repository URL is correct
- For private repos, ensure authentication is set up correctly

### No Changes Found
- Verify that the start commit is an ancestor of the end commit
- Check that you're comparing the correct branches/commits

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
