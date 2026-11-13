package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	const specSlug = "asking-spec"
	repo := newAuthorizationGitRepository(t)
	recordPath := writeAuthorizationRecordAt(t, repo, authorizationDocument(
		"approved", "2026-09-09", specSlug,
		"\noperations:\n  - implement\n  - push\n", "",
	))
	revision := commitAuthorizationFixture(t, repo, "docs: record authorization")
	specsRoot := filepath.Join(repo, "docs", "specs")
	for _, absent := range []bool{false, true} {
		name := "existing root"
		if absent {
			name = "absent root"
			if err := os.RemoveAll(filepath.Join(repo, "docs")); err != nil {
				t.Fatal(err)
			}
		}
		t.Run(name, func(t *testing.T) {
			resolved := ReadSpecAuthorization(t.Context(), repo, specsRoot, specSlug, revision)
			if resolved.Outcome != AuthorizationGranted || !resolved.Permits(AuthorizationOperationImplement) || !resolved.Permits(AuthorizationOperationPush) {
				t.Fatalf("authorization=%+v, want granted implement and push", resolved)
			}
			if resolved.Record.Source.Path != recordPath || resolved.Record.Source.Revision != revision {
				t.Fatalf("source=%+v, want path %q at %s", resolved.Record.Source, recordPath, revision)
			}
			if absent {
				if _, err := os.Lstat(filepath.Join(repo, "docs")); !os.IsNotExist(err) {
					t.Fatalf("authorization read recreated docs: %v", err)
				}
			}
		})
	}
}
