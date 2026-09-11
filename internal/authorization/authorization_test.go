// Suite: repository authorization record
// Invariant: the operative Spec 0119 grant permits its declared delivery actions
// Boundary IN: authorization reader and the tracked Spec-contained record
// Boundary OUT: daemon and CLI action dispatch
package authorization

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestAuthorizationParseCharacterization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		closing string
	}{
		// Task 02 changes this answer from a grant to a refusal.
		{name: "closing delimiter with extra characters grants today", closing: "---evil"},
		{name: "well-formed default-root record grants", closing: "---"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			const (
				slug       = "parse-characterization"
				recordPath = "docs/specs/parse-characterization/_authorization.md"
			)
			repoRoot := t.TempDir()
			recordFile := filepath.Join(repoRoot, filepath.FromSlash(recordPath))
			if err := os.MkdirAll(filepath.Dir(recordFile), 0o755); err != nil {
				t.Fatalf("create authorization directory: %v", err)
			}
			record := strings.Join([]string{
				"---",
				"status: approved",
				"granted: 2026-09-10",
				"action: characterize parser answers",
				"consuming: " + slug,
				"paths:",
				"  - Makefile",
				"operations:",
				"  - implement",
				tt.closing,
				"",
			}, "\n")
			if err := os.WriteFile(recordFile, []byte(record), 0o644); err != nil {
				t.Fatalf("write authorization record: %v", err)
			}

			resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot:   repoRoot,
				RecordPath: recordPath,
				Role:       AuthorizationRoleSpec,
				AskingSpec: slug,
			})
			if resolution.Outcome != AuthorizationGranted {
				t.Fatalf("authorization outcome = %q, want granted: %#v", resolution.Outcome, resolution.Reason)
			}
			if resolution.Record.Status != AuthorizationStatusApproved ||
				resolution.Record.Action != "characterize parser answers" ||
				!reflect.DeepEqual(resolution.Record.Consuming, []string{slug}) ||
				!reflect.DeepEqual(resolution.Record.Paths, []string{"Makefile"}) ||
				!reflect.DeepEqual(resolution.Record.Operations, []AuthorizationOperation{AuthorizationOperationImplement}) {
				t.Fatalf("authorization record = %#v, want exact characterization grant", resolution.Record)
			}
			for _, operation := range AllAuthorizationOperations() {
				if got, want := resolution.Permits(operation), operation == AuthorizationOperationImplement; got != want {
					t.Errorf("Permits(%q) = %t, want %t", operation, got, want)
				}
			}
		})
	}
}

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
