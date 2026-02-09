# githubCompare - Project Structure

```
githubCompare/
│
├── cmd/
│   └── root.go              # Cobra root command setup
│
├── internal/
│   ├── git/
│   │   ├── clone.go         # Repository cloning
│   │   ├── branches.go      # Branch listing
│   │   ├── commits.go       # Commit listing
│   │   ├── diff.go          # Change detection
│   │   └── types.go         # Git-related types
│   │
│   ├── archive/
│   │   ├── zip.go           # ZIP creation logic
│   │   └── naming.go        # Output file naming
│   │
│   ├── interactive/
│   │   ├── select.go        # Selection prompts
│   │   └── confirm.go       # Confirmation prompts
│   │
│   └── utils/
│       ├── temp.go          # Temp directory handling
│       ├── url.go           # URL parsing
│       └── auth.go          # Authentication helpers
│
├── main.go                  # Application entry point
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
│
├── README.md                # User documentation
├── MEMORY.md                # Project memory
├── plan.md                  # Project plan
├── functions.md             # Functions documentation
├── tree.md                  # This file
│
├── .gitignore               # Git ignore rules
├── .gitattributes           # Git attributes
│
├── Makefile                 # Build automation
│   # Targets:
│   # - make build-all      Build for all platforms
│   # - make build-linux
│   # - make build-windows
│   # - make build-darwin
│   # - make test
│   # - make clean
│
└── dist/                    # Build output (ignored by git)
    ├── githubCompare-linux-amd64
    ├── githubCompare-darwin-amd64
    ├── githubCompare-darwin-arm64
    └── githubCompare-windows-amd64.exe
```

## Notes

### Directory Naming Conventions
- `cmd/` - Command-line interface code
- `internal/` - Private application code (cannot be imported by other projects)
- `dist/` - Compiled binaries (not tracked in git)

### File Naming Conventions
- Snake case for Go files: `clone.go`, `branches.go`
- Descriptive names matching their primary function
- `types.go` for type definitions in each package

### Build Artifacts
The `dist/` directory will contain platform-specific binaries:
- Linux: `githubCompare-linux-amd64`
- macOS Intel: `githubCompare-darwin-amd64`
- macOS ARM: `githubCompare-darwin-arm64`
- Windows: `githubCompare-windows-amd64.exe`

### Temporary Files
During execution, the app will create temporary directories in the system temp location:
- Format: `/tmp/githubCompare-*` (Linux/macOS)
- Format: `%TEMP%\githubCompare-*` (Windows)
- Automatically cleaned up after execution
