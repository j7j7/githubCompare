package utils

import (
	"testing"
)

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		check   func(*RepoInfo) bool
	}{
		{
			name:    "HTTPS GitHub URL",
			url:     "https://github.com/owner/repo",
			wantErr: false,
			check: func(info *RepoInfo) bool {
				return info.Owner == "owner" && info.Name == "repo" && info.Protocol == "https"
			},
		},
		{
			name:    "HTTPS GitHub URL with .git",
			url:     "https://github.com/owner/repo.git",
			wantErr: false,
			check: func(info *RepoInfo) bool {
				return info.Owner == "owner" && info.Name == "repo"
			},
		},
		{
			name:    "SSH GitHub URL",
			url:     "git@github.com:owner/repo.git",
			wantErr: false,
			check: func(info *RepoInfo) bool {
				return info.Owner == "owner" && info.Name == "repo" && info.Protocol == "ssh"
			},
		},
		{
			name:    "SSH URL without .git",
			url:     "git@github.com:owner/repo",
			wantErr: false,
			check: func(info *RepoInfo) bool {
				return info.Owner == "owner" && info.Name == "repo"
			},
		},
		{
			name:    "Invalid URL",
			url:     "not-a-url",
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseRepoURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRepoURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && info != nil && tt.check != nil {
				if !tt.check(info) {
					t.Errorf("ParseRepoURL() returned unexpected result")
				}
			}
		})
	}
}
