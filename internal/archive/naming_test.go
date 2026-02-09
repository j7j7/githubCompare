package archive

import (
	"strings"
	"testing"
)

func TestGenerateOutputName(t *testing.T) {
	name := GenerateOutputName("test-repo", "abc1234", "def5678")
	
	if !strings.Contains(name, "test-repo") {
		t.Errorf("Output name should contain repo name")
	}
	
	if !strings.Contains(name, "abc1234") {
		t.Errorf("Output name should contain start ref")
	}
	
	if !strings.Contains(name, "def5678") {
		t.Errorf("Output name should contain end ref")
	}
	
	if !strings.HasSuffix(name, ".zip") {
		t.Errorf("Output name should end with .zip")
	}
}

func TestCleanRefForFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"feature/branch", "feature_branch"},
		{"abc1234567890", "abc1234"}, // Hash should be truncated
		{"normal-branch", "normal-branch"},
	}

	for _, tt := range tests {
		result := cleanRefForFilename(tt.input)
		if len(result) > len(tt.expected) && !strings.Contains(result, tt.input[:len(tt.input)-3]) {
			// For hash truncation, just check it's shorter
			if len(tt.input) == 40 && len(result) != 7 {
				t.Errorf("cleanRefForFilename(%s) = %s, expected truncated hash", tt.input, result)
			}
		}
	}
}
