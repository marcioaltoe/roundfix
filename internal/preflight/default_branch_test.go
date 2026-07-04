package preflight

import (
	"context"
	"errors"
	"testing"
)

func TestDetectDefaultBranch(t *testing.T) {
	symbolicRefKey := gitKey("symbolic-ref", "refs/remotes/origin/HEAD")
	originHeadUnset := errors.New("fatal: ref refs/remotes/origin/HEAD is not a symbolic ref")

	tests := []struct {
		name          string
		originHead    string
		originHeadErr error
		currentBranch string
		want          DefaultBranch
	}{
		{
			name:          "origin/HEAD names the default regardless of branch name",
			originHead:    "refs/remotes/origin/trunk",
			currentBranch: "feature/review",
			want:          DefaultBranch{Name: "trunk", Source: DefaultBranchFromOriginHead},
		},
		{
			name:          "origin/HEAD wins over a main current branch",
			originHead:    "refs/remotes/origin/trunk",
			currentBranch: "main",
			want:          DefaultBranch{Name: "trunk", Source: DefaultBranchFromOriginHead},
		},
		{
			name:          "origin/HEAD keeps slashes in the default branch name",
			originHead:    "refs/remotes/origin/release/2024",
			currentBranch: "feature/review",
			want:          DefaultBranch{Name: "release/2024", Source: DefaultBranchFromOriginHead},
		},
		{
			name:          "unset origin/HEAD matches main by name",
			originHeadErr: originHeadUnset,
			currentBranch: "main",
			want:          DefaultBranch{Name: "main", Source: DefaultBranchFromNameMatch},
		},
		{
			name:          "unset origin/HEAD matches master by name",
			originHeadErr: originHeadUnset,
			currentBranch: "master",
			want:          DefaultBranch{Name: "master", Source: DefaultBranchFromNameMatch},
		},
		{
			name:          "unset origin/HEAD leaves a feature branch undetermined",
			originHeadErr: originHeadUnset,
			currentBranch: "feature/review",
			want:          DefaultBranch{Source: DefaultBranchUndetermined},
		},
		{
			name:          "malformed origin/HEAD target falls back to the name match",
			originHead:    "refs/heads/main",
			currentBranch: "main",
			want:          DefaultBranch{Name: "main", Source: DefaultBranchFromNameMatch},
		},
		{
			name:          "malformed origin/HEAD target on a feature branch is undetermined",
			originHead:    "refs/heads/main",
			currentBranch: "feature/review",
			want:          DefaultBranch{Source: DefaultBranchUndetermined},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := fakeGitRunner{
				outputs: map[string]string{},
				errors:  map[string]error{},
			}
			if tt.originHeadErr != nil {
				runner.errors[symbolicRefKey] = tt.originHeadErr
			} else {
				runner.outputs[symbolicRefKey] = tt.originHead
			}

			got := DetectDefaultBranch(context.Background(), "/repo", tt.currentBranch, runner)

			if got != tt.want {
				t.Fatalf("DetectDefaultBranch(currentBranch=%q) = %#v, want %#v", tt.currentBranch, got, tt.want)
			}
		})
	}
}

func TestDefaultBranchIsDefault(t *testing.T) {
	tests := []struct {
		name      string
		detection DefaultBranch
		branch    string
		want      bool
	}{
		{
			name:      "origin/HEAD default matches its branch",
			detection: DefaultBranch{Name: "trunk", Source: DefaultBranchFromOriginHead},
			branch:    "trunk",
			want:      true,
		},
		{
			name:      "origin/HEAD default rejects a main branch by name alone",
			detection: DefaultBranch{Name: "trunk", Source: DefaultBranchFromOriginHead},
			branch:    "main",
			want:      false,
		},
		{
			name:      "name-match default matches the conventional branch",
			detection: DefaultBranch{Name: "master", Source: DefaultBranchFromNameMatch},
			branch:    "master",
			want:      true,
		},
		{
			name:      "undetermined default matches no branch",
			detection: DefaultBranch{Source: DefaultBranchUndetermined},
			branch:    "main",
			want:      false,
		},
		{
			name:      "undetermined default rejects an empty branch despite the empty name",
			detection: DefaultBranch{Source: DefaultBranchUndetermined},
			branch:    "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.detection.IsDefault(tt.branch)
			if got != tt.want {
				t.Fatalf("(%#v).IsDefault(%q) = %v, want %v", tt.detection, tt.branch, got, tt.want)
			}
		})
	}
}
