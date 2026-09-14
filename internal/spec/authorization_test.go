// Suite: authorization record reader.
// Invariant: only an operative record for the asking Spec grants its exact declared scope.
// Boundary IN: repository files, Git objects, and authorization record syntax.
// Boundary OUT: authoring diagnostics, changed-path audits, and regeneration output resolution.
package spec

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

const authorizationHistoryRevision = "6b8ea48725cbca13974eee0b400b3482202874f6"

func TestAuthorizationReaderClassifiesGrantState(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		want       AuthorizationOutcome
		wantReason AuthorizationReasonCode
		wantField  string
	}{
		{
			name: "approved record grants exact scope",
			content: authorizationDocument("approved", "2026-09-09", "asking-spec", `
operations:
  - implement
`, `
## Sanctioned regeneration

`+"```yaml\ncommand: make baseline-digests\n```\n"),
			want: AuthorizationGranted,
		},
		{
			name:       "proposal withholds on status",
			content:    authorizationDocument("proposed", "null", "asking-spec", "", ""),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonStatus,
			wantField:  "status",
		},
		{
			name:       "approved null date withholds on granted",
			content:    authorizationDocument("approved", "null", "asking-spec", "", ""),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonGranted,
			wantField:  "granted",
		},
		{
			name:       "unparseable date withholds on granted",
			content:    authorizationDocument("approved", "yesterday", "asking-spec", "", ""),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonGranted,
			wantField:  "granted",
		},
		{
			name:       "empty path list withholds on paths",
			content:    authorizationDocumentWithPaths(nil),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonPaths,
			wantField:  "paths",
		},
		{
			name: "prose mention cannot replace consuming field",
			content: authorizationDocument("approved", "2026-09-09", "another-spec", "", `
The asking-spec slug appears here only in prose.
`),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonConsuming,
			wantField:  "consuming",
		},
		{
			name:       "proposed record with grant date is contradictory",
			content:    authorizationDocument("proposed", "2026-09-09", "asking-spec", "", ""),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonContradictory,
			wantField:  "status,granted",
		},
		{
			name:       "unknown status is refused",
			content:    authorizationDocument("accepted", "2026-09-09", "asking-spec", "", ""),
			want:       AuthorizationRefused,
			wantReason: AuthorizationReasonStatus,
			wantField:  "status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoRoot, recordPath := writeAuthorizationRecord(t, tt.content)
			resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot:   repoRoot,
				RecordPath: recordPath,
				Role:       AuthorizationRoleSpec,
				AskingSpec: "asking-spec",
			})

			if resolved.Outcome != tt.want {
				t.Fatalf("ReadAuthorization() outcome = %q, want %q; reason = %#v", resolved.Outcome, tt.want, resolved.Reason)
			}
			if resolved.Reason.Code != tt.wantReason || resolved.Reason.Field != tt.wantField {
				t.Fatalf("ReadAuthorization() reason = %#v, want code %q and field %q", resolved.Reason, tt.wantReason, tt.wantField)
			}
			if resolved.Record.Status != AuthorizationStatus(strings.SplitN(tt.content, "\n", 3)[1][len("status: "):]) {
				t.Fatalf("ReadAuthorization() status = %q, want frontmatter status", resolved.Record.Status)
			}

			if tt.want != AuthorizationGranted {
				return
			}
			if got := resolved.Record.GrantedAt.Format("2006-01-02"); got != "2026-09-09" {
				t.Fatalf("ReadAuthorization() grant date = %q, want 2026-09-09", got)
			}
			if resolved.Record.Action != "implement the bounded change" {
				t.Fatalf("ReadAuthorization() action = %q", resolved.Record.Action)
			}
			if !reflect.DeepEqual(resolved.Record.Consuming, []string{"asking-spec"}) {
				t.Fatalf("ReadAuthorization() consuming = %q", resolved.Record.Consuming)
			}
			if !reflect.DeepEqual(resolved.Record.Paths, []string{"internal/spec/authorization.go"}) {
				t.Fatalf("ReadAuthorization() paths = %q", resolved.Record.Paths)
			}
			wantRegeneration := []AuthorizationRegeneration{{Command: "make baseline-digests"}}
			if !reflect.DeepEqual(resolved.Record.Regenerations, wantRegeneration) {
				t.Fatalf("ReadAuthorization() regenerations = %#v, want %#v", resolved.Record.Regenerations, wantRegeneration)
			}
		})
	}

	t.Run("malformed record is a typed refusal", func(t *testing.T) {
		repoRoot, recordPath := writeAuthorizationRecord(t, "not YAML frontmatter\n")
		resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
			RepoRoot: repoRoot, RecordPath: recordPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
		})
		if resolved.Outcome != AuthorizationRefused ||
			resolved.Reason.Code != AuthorizationReasonMalformedRecord ||
			resolved.Reason.Field != "frontmatter" {
			t.Fatalf("malformed record = %#v, want a frontmatter refusal", resolved)
		}
	})
}

func TestAuthorizationReaderTypesPermittedOperations(t *testing.T) {
	wantVocabulary := []AuthorizationOperation{
		AuthorizationOperationImplement,
		AuthorizationOperationCommit,
		AuthorizationOperationPush,
		AuthorizationOperationPullRequest,
		AuthorizationOperationMerge,
		AuthorizationOperationRelease,
	}
	if got := AllAuthorizationOperations(); !reflect.DeepEqual(got, wantVocabulary) {
		t.Fatalf("AllAuthorizationOperations() = %q, want exact closed vocabulary %q", got, wantVocabulary)
	}

	const allowedOperations = `
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
`
	repoRoot, recordPath := writeAuthorizationRecord(t, authorizationDocument(
		"approved",
		"2026-09-09",
		"asking-spec",
		allowedOperations,
		"",
	))
	resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot: repoRoot, RecordPath: recordPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
	})
	if resolved.Outcome != AuthorizationGranted {
		t.Fatalf("ReadAuthorization() outcome = %q, want granted; reason = %#v", resolved.Outcome, resolved.Reason)
	}
	for _, operation := range []AuthorizationOperation{
		AuthorizationOperationImplement,
		AuthorizationOperationCommit,
		AuthorizationOperationPush,
		AuthorizationOperationPullRequest,
		AuthorizationOperationMerge,
	} {
		if !resolved.Permits(operation) {
			t.Errorf("Permits(%q) = false, want true", operation)
		}
	}
	if resolved.Permits(AuthorizationOperationRelease) {
		t.Error("Permits(release) = true, want false")
	}
	if resolved.Permits(AuthorizationOperation("tag")) {
		t.Error("Permits(tag) = true, want false")
	}

	withoutOperationsRoot, withoutOperationsPath := writeAuthorizationRecord(t, authorizationDocument(
		"approved", "2026-09-09", "asking-spec", "", "",
	))
	withoutOperations := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot: withoutOperationsRoot, RecordPath: withoutOperationsPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
	})
	if withoutOperations.Outcome != AuthorizationGranted {
		t.Fatalf("record without operations outcome = %q, want granted; reason = %#v", withoutOperations.Outcome, withoutOperations.Reason)
	}
	for _, operation := range AllAuthorizationOperations() {
		if withoutOperations.Permits(operation) {
			t.Errorf("record without operations permits %q", operation)
		}
	}

	for _, token := range []string{"tag", "deploy"} {
		t.Run("rejects "+token, func(t *testing.T) {
			unknownRoot, unknownPath := writeAuthorizationRecord(t, authorizationDocument(
				"approved",
				"2026-09-09",
				"asking-spec",
				"\noperations:\n  - "+token+"\n",
				"",
			))
			unknown := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot: unknownRoot, RecordPath: unknownPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
			})
			if unknown.Outcome != AuthorizationRefused ||
				unknown.Reason.Code != AuthorizationReasonOperations ||
				unknown.Reason.Field != "operations" ||
				unknown.Reason.Value != token {
				t.Fatalf("unknown operation resolution = %#v, want refusal naming %q", unknown, token)
			}
		})
	}
}

func TestAuthorizationReaderRefusesEmptyPaths(t *testing.T) {
	t.Parallel()

	repoRoot, recordPath := writeAuthorizationRecord(t, authorizationDocumentWithPaths(nil))
	resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot: repoRoot, RecordPath: recordPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
	})
	if resolved.Outcome != AuthorizationRefused ||
		resolved.Reason.Code != AuthorizationReasonPaths ||
		resolved.Reason.Field != "paths" {
		t.Fatalf("empty paths resolution = %#v, want a paths refusal", resolved)
	}
}

func TestAuthorizationReaderRefusesEscapingPaths(t *testing.T) {
	symlinkRoot := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(symlinkRoot, "linked")); err != nil {
		t.Fatalf("create path-validation symlink: %v", err)
	}

	tests := []struct {
		name     string
		repoRoot string
		paths    []string
		wantPath string
	}{
		{name: "absolute", paths: []string{"/outside"}, wantPath: "/outside"},
		{name: "upward traversal", paths: []string{"docs/../outside"}, wantPath: "docs/../outside"},
		{name: "glob", paths: []string{"docs/*.md"}, wantPath: "docs/*.md"},
		{name: "symlink", repoRoot: symlinkRoot, paths: []string{"linked/file.md"}, wantPath: "linked/file.md"},
		{name: "duplicate", paths: []string{"docs/file.md", "docs/file.md"}, wantPath: "docs/file.md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoRoot := tt.repoRoot
			if repoRoot == "" {
				repoRoot = t.TempDir()
			}
			recordPath := writeAuthorizationRecordAt(t, repoRoot, authorizationDocumentWithPaths(tt.paths))
			resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
				RepoRoot: repoRoot, RecordPath: recordPath, Role: AuthorizationRoleSpec, AskingSpec: "asking-spec",
			})
			if resolved.Outcome != AuthorizationRefused ||
				resolved.Reason.Code != AuthorizationReasonPaths ||
				resolved.Reason.Field != "paths" ||
				resolved.Reason.Value != tt.wantPath {
				t.Fatalf("ReadAuthorization() = %#v, want path refusal naming %q", resolved, tt.wantPath)
			}
		})
	}
}

func TestAuthorizationReaderResolvesPreservedHistoricalRecords(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	records := authorizationHistoryPaths(t, repoRoot, authorizationHistoryRevision)
	if len(records) != 42 {
		t.Fatalf("historical authorization records = %d, want 42", len(records))
	}

	var multiSpec AuthorizationRecord
	for _, recordPath := range records {
		resolved := ReadAuthorization(context.Background(), AuthorizationReadRequest{
			RepoRoot:   repoRoot,
			RecordPath: recordPath,
			Revision:   authorizationHistoryRevision,
			Role:       AuthorizationRoleLegacy,
			AskingSpec: "0119-spec-contained-authorization",
		})
		if resolved.Outcome == AuthorizationUnresolved {
			t.Errorf("historical record %q is unresolved: %#v", recordPath, resolved.Reason)
			continue
		}
		if resolved.Outcome == AuthorizationGranted {
			t.Errorf("historical record %q grants authority to Spec 0119", recordPath)
		}
		if strings.HasSuffix(recordPath, "2026-08-12-the-authoring-and-baseline-corrections.md") {
			multiSpec = resolved.Record
		}
	}

	wantConsumers := []string{
		"0095-a-verification-that-ran-before-anyone-believed-it",
		"0096-a-failure-the-agent-can-read",
		"0098-a-hook-that-cannot-outrank-the-gate",
		"0104-a-gate-that-cannot-certify-its-own-cache",
		"0105-the-gates-own-economics",
		"0106-a-decision-that-reaches-every-artifact",
		"0107-the-authoring-rules-the-guides-do-not-carry",
		"0108-what-an-agent-loads-to-answer-one-question",
	}
	if !reflect.DeepEqual(multiSpec.Consuming, wantConsumers) {
		t.Fatalf("multi-Spec legacy consumers = %q, want %q", multiSpec.Consuming, wantConsumers)
	}
	if len(multiSpec.Paths) != 25 || multiSpec.Paths[0] != ".agents/skills/write-tasks/SKILL.md" || multiSpec.Paths[24] != "Makefile" {
		t.Fatalf("multi-Spec legacy paths = %q, want all 25 declared paths in order", multiSpec.Paths)
	}
	multiSpecGrant := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: "docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md",
		Revision:   authorizationHistoryRevision,
		Role:       AuthorizationRoleLegacy,
		AskingSpec: "0104-a-gate-that-cannot-certify-its-own-cache",
	})
	if multiSpecGrant.Outcome != AuthorizationGranted {
		t.Fatalf("multi-Spec legacy consumer = %#v, want granted", multiSpecGrant)
	}

	proseGrant := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: "docs/workflow/authorizations/2026-08-08-glossary-currency-clause.md",
		Revision:   authorizationHistoryRevision,
		Role:       AuthorizationRoleLegacy,
		AskingSpec: "0084-an-update-that-can-run",
	})
	if proseGrant.Outcome != AuthorizationGranted {
		t.Fatalf("prose legacy consumer = %#v, want granted", proseGrant)
	}
	if len(proseGrant.Record.Paths) != 3 {
		t.Fatalf("prose legacy paths = %q, want three exact declared paths", proseGrant.Record.Paths)
	}

	unreadable := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot: repoRoot, RecordPath: "docs/specs/missing/_authorization.md", Role: AuthorizationRoleSpec, AskingSpec: "missing",
	})
	if unreadable.Outcome != AuthorizationUnresolved || unreadable.Reason.Code != AuthorizationReasonUnreadableRecord {
		t.Fatalf("unreadable record = %#v, want unresolved unreadable-record reason", unreadable)
	}

	unavailable := ReadAuthorization(context.Background(), AuthorizationReadRequest{
		RepoRoot: repoRoot, RecordPath: records[0], Revision: strings.Repeat("0", 40), Role: AuthorizationRoleLegacy,
	})
	if unavailable.Outcome != AuthorizationUnresolved || unavailable.Reason.Code != AuthorizationReasonUnavailableRevision {
		t.Fatalf("unavailable revision = %#v, want unresolved unavailable-revision reason", unavailable)
	}
}

func TestOperationAuthorityResolvesExternalSpecRoot(t *testing.T) {
	const specSlug = "asking-spec"
	projectRoot := newAuthorizationGitRepository(t)
	projectRevision := commitAuthorizationFixture(t, projectRoot, "seed project")
	externalRepoRoot := newAuthorizationGitRepository(t)
	externalSpecsRoot := filepath.Join(externalRepoRoot, "specs")
	externalRecordPath := filepath.ToSlash(filepath.Join("specs", specSlug, "_authorization.md"))
	writeAuthorizationRecordAtPath(t, externalRepoRoot, externalRecordPath, authorizationDocument(
		"approved",
		"2026-09-09",
		specSlug,
		"\noperations:\n  - implement\n",
		"",
	))
	externalRevision := commitAuthorizationFixture(t, externalRepoRoot, "seed external Spec")
	writeAuthorizationRecordAtPath(t, externalRepoRoot, externalRecordPath, authorizationDocument(
		"proposed",
		"null",
		specSlug,
		"\noperations:\n  - implement\n",
		"",
	))

	resolved := ReadSpecAuthorization(context.Background(), projectRoot, externalSpecsRoot, specSlug, projectRevision)

	if resolved.Outcome != AuthorizationGranted || !resolved.Permits(AuthorizationOperationImplement) {
		t.Fatalf("external-root operation authority = %#v, want granted implement authority", resolved)
	}
	if resolved.Record.Source.Path != externalRecordPath {
		t.Fatalf("external-root record path = %q, want root-derived %q", resolved.Record.Source.Path, externalRecordPath)
	}
	if resolved.Record.Source.Revision != externalRevision {
		t.Fatalf("external-root record revision = %q, want external HEAD %q", resolved.Record.Source.Revision, externalRevision)
	}
	if _, err := os.Stat(filepath.Join(projectRoot, filepath.FromSlash(AuthorizationRecordPath(specSlug)))); !os.IsNotExist(err) {
		t.Fatalf("project repository contains duplicate authorization record, stat error = %v", err)
	}
}

func TestOperationAuthorityDefaultRootUnchanged(t *testing.T) {
	const specSlug = "asking-spec"
	projectRoot := newAuthorizationGitRepository(t)
	writeAuthorizationRecordAt(t, projectRoot, authorizationDocument(
		"approved",
		"2026-09-09",
		specSlug,
		"\noperations:\n  - implement\n",
		"",
	))
	projectRevision := commitAuthorizationFixture(t, projectRoot, "seed default Spec Root")

	resolved := ReadSpecAuthorization(
		context.Background(),
		projectRoot,
		filepath.Join(projectRoot, "docs", "specs"),
		specSlug,
		projectRevision,
	)

	if resolved.Outcome != AuthorizationGranted || !resolved.Permits(AuthorizationOperationImplement) {
		t.Fatalf("default-root operation authority = %#v, want granted implement authority", resolved)
	}
	if resolved.Record.Source.Path != AuthorizationRecordPath(specSlug) {
		t.Fatalf("default-root record path = %q, want unchanged %q", resolved.Record.Source.Path, AuthorizationRecordPath(specSlug))
	}
	if resolved.Record.Source.Revision != projectRevision {
		t.Fatalf("default-root record revision = %q, want unchanged %q", resolved.Record.Source.Revision, projectRevision)
	}
}

func TestSpecAuthorizationNamesUnavailableRevision(t *testing.T) {
	const (
		specSlug            = "asking-spec"
		unavailableRevision = "refs/heads/unavailable"
	)
	projectRoot := newAuthorizationGitRepository(t)
	writeAuthorizationRecordAt(t, projectRoot, authorizationDocument(
		"approved",
		"2026-09-09",
		specSlug,
		"\noperations:\n  - implement\n",
		"",
	))
	commitAuthorizationFixture(t, projectRoot, "seed default Spec Root")

	resolved := ReadSpecAuthorization(
		context.Background(),
		projectRoot,
		filepath.Join(projectRoot, "docs", "specs"),
		specSlug,
		unavailableRevision,
	)

	if resolved.Outcome != AuthorizationUnresolved {
		t.Fatalf("unavailable-revision authorization outcome = %q, want unresolved: %#v", resolved.Outcome, resolved)
	}
	if resolved.Reason.Code != AuthorizationReasonUnavailableRevision ||
		resolved.Reason.Field != "revision" ||
		resolved.Reason.Value != unavailableRevision {
		t.Fatalf("unavailable-revision authorization reason = %#v, want unavailable revision %q", resolved.Reason, unavailableRevision)
	}
	if resolved.Permits(AuthorizationOperationImplement) {
		t.Fatal("unavailable-revision authorization permits implement")
	}
}

func TestOperationAuthorityReportsUnresolvableSpecRoot(t *testing.T) {
	projectRoot := newAuthorizationGitRepository(t)
	projectRevision := commitAuthorizationFixture(t, projectRoot, "seed project")
	unresolvableRoot := filepath.Join(t.TempDir(), "not-a-git-repository")
	if err := os.MkdirAll(unresolvableRoot, 0o755); err != nil {
		t.Fatalf("create unresolvable Spec Root: %v", err)
	}

	resolved := ReadSpecAuthorization(context.Background(), projectRoot, unresolvableRoot, "asking-spec", projectRevision)

	if resolved.Outcome != AuthorizationUnresolved {
		t.Fatalf("unresolvable-root operation authority outcome = %q, want unresolved: %#v", resolved.Outcome, resolved)
	}
	if resolved.Reason.Code != AuthorizationReasonUnreadableRecord || resolved.Reason.Field != "spec_root" {
		t.Fatalf("unresolvable-root operation authority reason = %#v, want unreadable spec_root", resolved.Reason)
	}
}

func authorizationDocument(status, granted, consuming, extraFrontmatter, body string) string {
	return "---\n" +
		"status: " + status + "\n" +
		"granted: " + granted + "\n" +
		"action: implement the bounded change\n" +
		"consuming: " + consuming + "\n" +
		"paths:\n" +
		"  - internal/spec/authorization.go\n" +
		extraFrontmatter +
		"---\n" +
		body
}

func authorizationDocumentWithPaths(paths []string) string {
	var declared strings.Builder
	for _, path := range paths {
		declared.WriteString("  - ")
		declared.WriteString(path)
		declared.WriteByte('\n')
	}
	return "---\n" +
		"status: approved\n" +
		"granted: 2026-09-09\n" +
		"action: validate exact paths\n" +
		"consuming: asking-spec\n" +
		"paths:\n" + declared.String() +
		"---\n"
}

func writeAuthorizationRecord(t *testing.T, content string) (string, string) {
	t.Helper()
	repoRoot := t.TempDir()
	return repoRoot, writeAuthorizationRecordAt(t, repoRoot, content)
}

func writeAuthorizationRecordAt(t *testing.T, repoRoot, content string) string {
	t.Helper()
	const recordPath = "docs/specs/asking-spec/_authorization.md"
	return writeAuthorizationRecordAtPath(t, repoRoot, recordPath, content)
}

func writeAuthorizationRecordAtPath(t *testing.T, repoRoot, recordPath, content string) string {
	t.Helper()
	absPath := filepath.Join(repoRoot, filepath.FromSlash(recordPath))
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		t.Fatalf("create authorization directory: %v", err)
	}
	if err := os.WriteFile(absPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write authorization record: %v", err)
	}
	return recordPath
}

func authorizationHistoryPaths(t *testing.T, repoRoot, revision string) []string {
	t.Helper()
	command := exec.Command(
		"git", "-C", repoRoot, "-c", "core.fsmonitor=false",
		"ls-tree", "-r", "--name-only", revision, "docs/workflow/authorizations",
	)
	command.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("list historical authorization records at %s: %v: %s", revision, err, output)
	}
	return strings.Fields(string(output))
}

func newAuthorizationGitRepository(t *testing.T) string {
	t.Helper()
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "--initial-branch=main")
	return repoRoot
}

func commitAuthorizationFixture(t *testing.T, repoRoot, message string) string {
	t.Helper()
	gittest.Run(t, repoRoot, "add", "-A")
	gittest.Run(t, repoRoot, "commit", "--allow-empty", "-m", message)
	return strings.TrimSpace(gittest.Run(t, repoRoot, "rev-parse", "HEAD"))
}
