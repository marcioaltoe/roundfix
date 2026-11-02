package authorization

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Suite: Release Ceiling record keys.
// Invariant: the ceiling and exact release paths form one operative release grant.
// Boundary IN: public authorization reader, working-tree records.
// Boundary OUT: release planning and release execution.
func TestReleaseCeilingIsReadAndValidated(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, fields, ceiling string
		release               bool
		code                  AuthorizationReasonCode
		field                 string
	}{
		{name: "granted pair", fields: "release_ceiling: minor\nrelease_paths: [VERSION, CHANGELOG.md]\n", ceiling: "minor", release: true},
		{name: "patch", fields: "release_ceiling: patch\nrelease_paths: [VERSION, CHANGELOG.md]\n", ceiling: "patch", release: true},
		{name: "major", fields: "release_ceiling: major\nrelease_paths: [VERSION, CHANGELOG.md]\n", ceiling: "major", release: true},
		{name: "without release", fields: "release_ceiling: minor\nrelease_paths: [VERSION]\n", code: AuthorizationReasonContradictory, field: "release_ceiling"},
		{name: "unknown ceiling", fields: "release_ceiling: huge\nrelease_paths: [VERSION]\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_ceiling"},
		{name: "ceiling without paths", fields: "release_ceiling: minor\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_paths"},
		{name: "paths without ceiling", fields: "release_paths: [VERSION]\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_paths"},
		{name: "without keys"},
		{name: "release without keys", release: true},
		{name: "null ceiling", fields: "release_ceiling: null\nrelease_paths: [VERSION]\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_ceiling"},
		{name: "empty ceiling", fields: "release_ceiling: ''\nrelease_paths: [VERSION]\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_ceiling"},
		{name: "empty paths", fields: "release_ceiling: minor\nrelease_paths: []\n", release: true, code: AuthorizationReasonMalformedRecord, field: "release_paths"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root, record := writeEmptyPathsAuthorization(t, AuthorizationRoleSpec, "paths: []\n"+tt.fields)
			file := filepath.Join(root, filepath.FromSlash(record))
			if tt.release {
				content, err := os.ReadFile(file)
				if err != nil {
					t.Fatal(err)
				}
				content = []byte(strings.Replace(string(content), "  - implement", "  - implement\n  - release", 1))
				if err := os.WriteFile(file, content, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got := ReadAuthorization(context.Background(), AuthorizationReadRequest{RepoRoot: root, RecordPath: record, Role: AuthorizationRoleSpec, AskingSpec: emptyPathsTestSlug})
			if tt.code != "" {
				if got.Outcome != AuthorizationRefused || got.Reason.Code != tt.code || got.Reason.Field != tt.field {
					t.Fatalf("resolution = %+v", got)
				}
				return
			}
			if got.Outcome != AuthorizationGranted || got.Record.ReleaseCeiling != tt.ceiling || !got.Permits(AuthorizationOperationImplement) || got.Permits(AuthorizationOperationRelease) != tt.release {
				t.Fatalf("resolution = %+v", got)
			}
			var want []string
			if tt.ceiling != "" {
				want = []string{"VERSION", "CHANGELOG.md"}
			}
			if !reflect.DeepEqual(got.Record.ReleasePaths, want) {
				t.Fatalf("paths=%v want=%v", got.Record.ReleasePaths, want)
			}
		})
	}
}

func TestReleasePathsUseExactExecutablePathRules(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"[/VERSION]", "[x/../VERSION]", "[x//VERSION]", "['*.md']", "[x/]", "[' VERSION ']", "[VERSION, VERSION]", "VERSION", "['']", "['.']", "['..']", "['x\\VERSION']", "[42]"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			root, record := writeEmptyPathsAuthorization(t, AuthorizationRoleSpec, "paths: []\nrelease_paths: "+value+"\n")
			got := ReadAuthorization(context.Background(), AuthorizationReadRequest{RepoRoot: root, RecordPath: record, Role: AuthorizationRoleSpec})
			if got.Outcome != AuthorizationRefused || got.Reason.Code != AuthorizationReasonMalformedRecord || got.Reason.Field != "release_paths" {
				t.Fatalf("resolution=%+v", got)
			}
		})
	}
}
