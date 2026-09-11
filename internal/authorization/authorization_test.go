// Suite: repository authorization record
// Invariant: the operative Spec 0119 grant permits its declared delivery actions
// Boundary IN: authorization reader and the tracked Spec-contained record
// Boundary OUT: daemon and CLI action dispatch
package authorization

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestAuthorizationParseCharacterization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		closing     string
		wantOutcome AuthorizationOutcome
	}{
		{name: "closing delimiter with extra characters refuses", closing: "---evil", wantOutcome: AuthorizationRefused},
		{name: "well-formed default-root record grants", closing: "---", wantOutcome: AuthorizationGranted},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			const slug = "parse-characterization"
			repoRoot, recordPath := writeAuthorizationRecordForTest(
				t,
				slug,
				"characterize parser answers",
				"---",
				tt.closing,
				"\n",
			)

			resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot:   repoRoot,
				RecordPath: recordPath,
				Role:       AuthorizationRoleSpec,
				AskingSpec: slug,
			})
			if resolution.Outcome != tt.wantOutcome {
				t.Fatalf("authorization outcome = %q, want %q: %#v", resolution.Outcome, tt.wantOutcome, resolution.Reason)
			}
			if tt.wantOutcome == AuthorizationRefused {
				if resolution.Reason.Code != AuthorizationReasonMalformedRecord ||
					!strings.Contains(resolution.Reason.Detail, tt.closing) {
					t.Fatalf("authorization reason = %#v, want malformed record naming %q", resolution.Reason, tt.closing)
				}
				return
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

func TestDelimiterWithTrailingCarriageReturnGrants(t *testing.T) {
	t.Parallel()

	const slug = "carriage-return-delimiter"
	repoRoot, recordPath := writeAuthorizationRecordForTest(
		t,
		slug,
		"accept carriage return delimiter",
		"---",
		"---",
		"\r\n",
	)

	resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: recordPath,
		Role:       AuthorizationRoleSpec,
		AskingSpec: slug,
	})
	if resolution.Outcome != AuthorizationGranted {
		t.Fatalf("authorization outcome = %q, want granted: %#v", resolution.Outcome, resolution.Reason)
	}
}

func TestMissingAuthorizationEvidenceRemainsUnresolved(t *testing.T) {
	t.Parallel()

	t.Run("unreadable record", func(t *testing.T) {
		t.Parallel()

		resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
			RepoRoot:   t.TempDir(),
			RecordPath: "docs/specs/missing/_authorization.md",
			Role:       AuthorizationRoleSpec,
		})
		if resolution.Outcome != AuthorizationUnresolved || resolution.Reason.Code != AuthorizationReasonUnreadableRecord {
			t.Fatalf("authorization resolution = %#v, want unresolved unreadable record", resolution)
		}
	})

	t.Run("unavailable revision", func(t *testing.T) {
		t.Parallel()

		repoRoot := t.TempDir()
		command := exec.Command("git", "init", "--quiet", repoRoot)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("initialize repository: %v: %s", err, output)
		}
		resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
			RepoRoot:   repoRoot,
			RecordPath: "docs/specs/missing/_authorization.md",
			Revision:   "refs/heads/unavailable",
			Role:       AuthorizationRoleSpec,
		})
		if resolution.Outcome != AuthorizationUnresolved || resolution.Reason.Code != AuthorizationReasonUnavailableRevision {
			t.Fatalf("authorization resolution = %#v, want unresolved unavailable revision", resolution)
		}
	})
}

func TestMalformedDelimiterRefuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opening   string
		closing   string
		offending string
	}{
		{name: "opening delimiter has extra characters", opening: "---evil", closing: "---", offending: "---evil"},
		{name: "closing delimiter has extra characters", opening: "---", closing: "---evil", offending: "---evil"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			const slug = "malformed-delimiter"
			repoRoot, recordPath := writeAuthorizationRecordForTest(
				t,
				slug,
				"refuse malformed delimiter",
				tt.opening,
				tt.closing,
				"\n",
			)

			resolution := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot:   repoRoot,
				RecordPath: recordPath,
				Role:       AuthorizationRoleSpec,
				AskingSpec: slug,
			})
			if resolution.Outcome != AuthorizationRefused {
				t.Fatalf("authorization outcome = %q, want refused: %#v", resolution.Outcome, resolution)
			}
			if resolution.Reason.Code != AuthorizationReasonMalformedRecord ||
				resolution.Reason.Field != "frontmatter" ||
				!strings.Contains(resolution.Reason.Detail, tt.offending) {
				t.Fatalf("authorization reason = %#v, want malformed frontmatter naming %q", resolution.Reason, tt.offending)
			}
		})
	}
}

func writeAuthorizationRecordForTest(
	t *testing.T,
	slug string,
	action string,
	opening string,
	closing string,
	lineEnding string,
) (string, string) {
	t.Helper()

	recordPath := "docs/specs/" + slug + "/_authorization.md"
	repoRoot := t.TempDir()
	recordFile := filepath.Join(repoRoot, filepath.FromSlash(recordPath))
	if err := os.MkdirAll(filepath.Dir(recordFile), 0o755); err != nil {
		t.Fatalf("create authorization directory: %v", err)
	}
	record := strings.Join([]string{
		opening,
		"status: approved",
		"granted: 2026-09-10",
		"action: " + action,
		"consuming: " + slug,
		"paths:",
		"  - Makefile",
		"operations:",
		"  - implement",
		closing,
		"",
	}, lineEnding)
	if err := os.WriteFile(recordFile, []byte(record), 0o644); err != nil {
		t.Fatalf("write authorization record: %v", err)
	}
	return repoRoot, recordPath
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
