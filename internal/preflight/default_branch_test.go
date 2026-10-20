package preflight

import (
	"context"
	"errors"
	"testing"
)

func TestDetectDefaultBranch(t *testing.T) {
	t.Parallel()
	symbolicRefKey := gitKey("symbolic-ref", "refs/remotes/origin/HEAD")
	originHeadUnset := errors.New("fatal: ref refs/remotes/origin/HEAD is not a symbolic ref")

	tests := []struct {
		name          string
		originHead    string
		originHeadErr error
		originRefs    string
		originRefsErr error
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
		{name: "one fetched main", originRefs: "refs/remotes/origin/main\n", want: DefaultBranch{Name: "main", Source: DefaultBranchFromOriginBranch}},
		{name: "one fetched master", originRefs: "refs/remotes/origin/master\n", want: DefaultBranch{Name: "master", Source: DefaultBranchFromOriginBranch}},
		{name: "both fetched defaults", originRefs: "refs/remotes/origin/main\nrefs/remotes/origin/master\n", want: DefaultBranch{Source: DefaultBranchUndetermined}},
		{name: "no fetched default", want: DefaultBranch{Source: DefaultBranchUndetermined}},
		{name: "fetched refs fail", originRefs: "refs/remotes/origin/main\n", originRefsErr: errors.New("refs failed"), want: DefaultBranch{Source: DefaultBranchUndetermined}},
		{name: "only exact refs count", originRefs: "refs/remotes/origin/main/extra\nrefs/remotes/origin/master-old\n", want: DefaultBranch{Source: DefaultBranchUndetermined}},
		{name: "origin HEAD wins over fetched defaults", originHead: "refs/remotes/origin/trunk", originRefs: "refs/remotes/origin/main\n", want: DefaultBranch{Name: "trunk", Source: DefaultBranchFromOriginHead}},
		{name: "name match wins over fetched defaults", currentBranch: "master", originRefs: "refs/remotes/origin/main\n", want: DefaultBranch{Name: "master", Source: DefaultBranchFromNameMatch}},
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

			refsKey := gitKey("for-each-ref", "--format=%(refname)", "refs/remotes/origin/main", "refs/remotes/origin/master")
			runner.outputs[refsKey] = tt.originRefs
			runner.errors[refsKey] = tt.originRefsErr
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

func TestResolveDefaultBranchHeadPrefersTheFetchedBranch(t *testing.T) {
	t.Parallel()
	gitErr := errors.New("cannot read refs")
	branch := DefaultBranch{Name: "main", Source: DefaultBranchFromNameMatch}
	for _, tc := range []struct {
		name, current, refs string
		err                 error
		want                DefaultBranchHead
		found               bool
	}{
		{name: "remote wins over local", current: "main", refs: "refs/heads/main\x00local\nrefs/remotes/origin/main\x00remote\n", want: DefaultBranchHead{Branch: branch, Ref: "refs/remotes/origin/main", Commit: "remote"}, found: true},
		{name: "local alone", current: "main", refs: "refs/heads/main\x00local\n", want: DefaultBranchHead{Branch: branch, Ref: "refs/heads/main", Commit: "local"}, found: true},
		{name: "undetermined"},
		{name: "absent refs", current: "main"},
		{name: "exact refs only", current: "main", refs: "refs/heads/main/extra\x00other\nrefs/remotes/origin/main/extra\x00other\n"},
		{name: "git error", current: "main", err: gitErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			key := gitKey("for-each-ref", "--format=%(refname)%00%(objectname)", "refs/remotes/origin/main", "refs/heads/main")
			runner := fakeGitRunner{outputs: map[string]string{
				key: tc.refs,
				gitKey("symbolic-ref", "refs/remotes/origin/HEAD"):                                                      "",
				gitKey("for-each-ref", "--format=%(refname)", "refs/remotes/origin/main", "refs/remotes/origin/master"): "",
			}, errors: map[string]error{key: tc.err}}
			got, found, err := ResolveDefaultBranchHead(t.Context(), "/repo", tc.current, runner)
			if !errors.Is(err, tc.err) || found != tc.found || got != tc.want {
				t.Fatalf("head=%+v found=%v err=%v, want %+v found=%v err=%v", got, found, err, tc.want, tc.found, tc.err)
			}
		})
	}
}
