package display

import (
	"strings"
	"github.com/fatih/color"
)

var (
	// Header colors
	Header    = color.New(color.FgCyan, color.Bold)
	Success   = color.New(color.FgGreen, color.Bold)
	Error     = color.New(color.FgRed, color.Bold)
	Warning   = color.New(color.FgYellow)
	Info      = color.New(color.FgBlue)
	
	// Content colors
	Branch    = color.New(color.FgMagenta)
	Commit    = color.New(color.FgYellow)
	File      = color.New(color.FgCyan)
	Count     = color.New(color.FgGreen, color.Bold)
	
	// Status colors
	Added     = color.New(color.FgGreen)
	Modified  = color.New(color.FgYellow)
	Deleted   = color.New(color.FgRed)
	Renamed   = color.New(color.FgMagenta)
)

// PrintHeader prints a formatted header
func PrintHeader(text string) {
	Header.Println("\n" + text)
	Header.Println(strings.Repeat("=", len(text)))
}

// PrintSection prints a section header
func PrintSection(text string) {
	Info.Printf("\n▶ %s\n", text)
}

// PrintSuccess prints a success message
func PrintSuccess(text string) {
	Success.Printf("✓ %s\n", text)
}

// PrintError prints an error message
func PrintError(text string) {
	Error.Printf("✗ %s\n", text)
}

// PrintWarning prints a warning message
func PrintWarning(text string) {
	Warning.Printf("⚠ %s\n", text)
}

// PrintInfo prints an info message
func PrintInfo(text string) {
	Info.Printf("ℹ %s\n", text)
}
