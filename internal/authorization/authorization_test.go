// Suite: repository authorization record
// Invariant: the operative Spec 0119 grant permits its declared delivery actions
// Boundary IN: authorization reader and the tracked Spec-contained record
// Boundary OUT: daemon and CLI action dispatch
package authorization

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCurrentRecordPermitsItsDeclaredOperations(t *testing.T) {
	t.Parallel()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate authorization test source")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", ".."))
	const recordPath = "docs/specs/0119-spec-contained-authorization/_authorization.md"
	resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: recordPath,
		Role:       AuthorizationRoleSpec,
		AskingSpec: "0119-spec-contained-authorization",
	})
	if resolution.Outcome != AuthorizationGranted {
		t.Fatalf("current authorization outcome = %q, want granted: %#v", resolution.Outcome, resolution.Reason)
	}
	for _, operation := range []AuthorizationOperation{
		AuthorizationOperationImplement,
		AuthorizationOperationCommit,
		AuthorizationOperationPush,
	} {
		if !resolution.Permits(operation) {
			t.Errorf("current authorization record %s does not permit %q", recordPath, operation)
		}
	}
}
