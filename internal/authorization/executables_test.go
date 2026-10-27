package authorization

// Suite: authorized executables
// Invariant: only exact file paths can enter an executable grant.
// Boundary IN: authorization record reader.
// Boundary OUT: Task staging and Git mode recording.

import (
	"context"
	"slices"
	"strings"
	"testing"
)

func TestAuthorizationReadsExecutables(t *testing.T) {
	t.Parallel()
	for _, declaration := range []string{"executables: [scripts/launch.sh, scripts/check.sh]\n", "executables: null\n", ""} {
		t.Run(declaration, func(t *testing.T) {
			t.Parallel()
			root, record := writeEmptyPathsAuthorization(t, AuthorizationRoleSpec, "paths: []\n"+declaration)
			got := ReadAuthorization(context.Background(), AuthorizationReadRequest{RepoRoot: root, RecordPath: record, Role: AuthorizationRoleSpec, AskingSpec: emptyPathsTestSlug})
			if got.Outcome != AuthorizationGranted {
				t.Fatalf("resolution = %+v", got)
			}
			var want []string
			if strings.Contains(declaration, "scripts/") {
				want = []string{"scripts/launch.sh", "scripts/check.sh"}
			}
			if !slices.Equal(got.Record.Executables, want) {
				t.Fatalf("executables = %v, want %v", got.Record.Executables, want)
			}
		})
	}
}

func TestAuthorizationRefusesAMalformedExecutablesEntry(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, value string }{
		{"absolute path", "[/scripts/launch.sh]"},
		{"upward segment", "[scripts/../launch.sh]"},
		{"unclean path", "[scripts//launch.sh]"},
		{"glob", "['scripts/*.sh']"},
		{"trailing slash", "[scripts/]"},
		{"surrounding whitespace", "[' scripts/launch.sh ']"},
		{"duplicate", "[scripts/launch.sh, scripts/launch.sh]"},
		{"scalar value", "scripts/launch.sh"},
		{"empty path", "['']"},
		{"dot", "['.']"},
		{"parent", "['..']"},
		{"backslash", "['scripts\\launch.sh']"},
		{"non string", "[42]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root, record := writeEmptyPathsAuthorization(t, AuthorizationRoleSpec, "paths: []\nexecutables: "+test.value+"\n")
			got := ReadAuthorization(context.Background(), AuthorizationReadRequest{RepoRoot: root, RecordPath: record, Role: AuthorizationRoleSpec, AskingSpec: emptyPathsTestSlug})
			if got.Outcome == AuthorizationGranted || got.Reason.Code != AuthorizationReasonMalformedRecord || got.Reason.Field != "executables" || !strings.HasPrefix(got.Reason.Detail, "executables") {
				t.Fatalf("resolution = %+v", got)
			}
		})
	}
}
