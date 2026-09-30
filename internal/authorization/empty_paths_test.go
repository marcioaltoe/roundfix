// Suite: explicit empty authorization paths.
// Invariant: only a Spec-contained explicit empty sequence grants operations without bounding a path.
// Boundary IN: authorization frontmatter parsing and classification.
// Boundary OUT: CLI delivery and governed-path consumers.
package authorization

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const emptyPathsTestSlug = "empty-paths-grant"

func TestAnExplicitEmptyPathsListGrantsOperations(t *testing.T) {
	t.Parallel()

	repoRoot, recordPath := writeEmptyPathsAuthorization(t, AuthorizationRoleSpec, "paths: []\n")
	resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: recordPath,
		Role:       AuthorizationRoleSpec,
		AskingSpec: emptyPathsTestSlug,
	})

	if resolution.Outcome != AuthorizationGranted {
		t.Fatalf("authorization outcome = %q, want granted: %#v", resolution.Outcome, resolution.Reason)
	}
	if resolution.Record.Paths == nil || len(resolution.Record.Paths) != 0 {
		t.Fatalf("authorization Paths = %#v, want a non-nil empty slice", resolution.Record.Paths)
	}
	if !resolution.Permits(AuthorizationOperationImplement) {
		t.Fatal("explicit empty paths grant does not permit its listed implement operation")
	}
}

func TestANullPathsFieldStillRefuses(t *testing.T) {
	t.Parallel()

	assertEmptyPathsRefusal(t, AuthorizationRoleSpec, "paths:\n")
}

func TestAnAbsentPathsFieldStillRefuses(t *testing.T) {
	t.Parallel()

	assertEmptyPathsRefusal(t, AuthorizationRoleSpec, "")
}

func TestALegacyRecordWithEmptyPathsStillRefuses(t *testing.T) {
	t.Parallel()

	assertEmptyPathsRefusal(t, AuthorizationRoleLegacy, "paths: []\n")
}

func assertEmptyPathsRefusal(t *testing.T, role AuthorizationRole, pathsDeclaration string) {
	t.Helper()

	repoRoot, recordPath := writeEmptyPathsAuthorization(t, role, pathsDeclaration)
	resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: recordPath,
		Role:       role,
		AskingSpec: emptyPathsTestSlug,
	})

	if resolution.Outcome != AuthorizationRefused ||
		resolution.Reason.Code != AuthorizationReasonPaths ||
		resolution.Reason.Field != "paths" ||
		resolution.Reason.Detail != "paths must contain at least one exact repository-relative path" {
		t.Fatalf("authorization resolution = %#v, want the unchanged paths refusal", resolution)
	}
}

func writeEmptyPathsAuthorization(
	t *testing.T,
	role AuthorizationRole,
	pathsDeclaration string,
) (string, string) {
	t.Helper()

	recordPath := "docs/specs/" + emptyPathsTestSlug + "/_authorization.md"
	if role == AuthorizationRoleLegacy {
		recordPath = "docs/workflow/authorizations/2026-09-30-empty-paths.md"
	}
	repoRoot := t.TempDir()
	absolutePath := filepath.Join(repoRoot, filepath.FromSlash(recordPath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatalf("create authorization fixture directory: %v", err)
	}
	content := strings.Join([]string{
		"---",
		"status: approved",
		"granted: 2026-09-30",
		"action: grant operations without bounded paths",
		"consuming: " + emptyPathsTestSlug,
		strings.TrimSuffix(pathsDeclaration, "\n"),
		"operations:",
		"  - implement",
		"---",
		"",
	}, "\n")
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write authorization fixture: %v", err)
	}
	return repoRoot, recordPath
}
