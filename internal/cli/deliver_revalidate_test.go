// Suite: Delivery Revalidation command workflow.
// Invariant: strict findings and changed production premises are derived from the item worktree and real prior merge commits.
// Boundary IN: real disposable Git repositories, Spec loading, Spec Consistency Check, and deliver status output.
// Boundary OUT: no network or live Run Database is used.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

func TestRevalidateReportsAnUnresolvedDeclarationOnlyAfterTheMergeRemovedIt(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryRevalidationGitFixture(t)

	before, err := fixture.workflow(fixture.before).Revalidate(t.Context(), fixture.before, "clean", nil)
	if err != nil {
		t.Fatalf("revalidate before merge: %v", err)
	}
	after, err := fixture.workflow(fixture.after).Revalidate(t.Context(), fixture.after, "clean", []string{fixture.mergeCommit})
	if err != nil {
		t.Fatalf("revalidate after merge: %v", err)
	}

	if containsRevalidationString(before.Findings, "SC-REF-UNRESOLVED") {
		findings, _ := strictSpecFindings(filepath.Join(fixture.before, "docs", "specs"), fixture.before, "clean")
		t.Fatalf("findings before merge = %v (%+v), want no unresolved declaration", before.Findings, findings)
	}
	if !containsRevalidationString(after.Findings, "SC-REF-UNRESOLVED") {
		t.Fatalf("findings after merge = %v, want SC-REF-UNRESOLVED", after.Findings)
	}
}

func TestRevalidateNamesADeclaredProductionFileAndTheMergeThatChangedIt(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryRevalidationGitFixture(t)

	result, err := fixture.workflow(fixture.after).Revalidate(
		t.Context(),
		fixture.after,
		"clean",
		[]string{fixture.mergeCommit},
	)
	if err != nil {
		t.Fatalf("revalidate after merge: %v", err)
	}

	if got, want := result.ChangedPremises, []string{"internal/changed.go", "internal/removed.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("changed premises = %v, want %v", got, want)
	}
	if got, want := result.ChangedBy, []string{fixture.mergeCommit}; !reflect.DeepEqual(got, want) {
		t.Fatalf("changed-by merges = %v, want %v", got, want)
	}
}

func TestRevalidateIgnoresChangedTestsGuidesAndUndeclaredFiles(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryRevalidationGitFixture(t)

	result, err := fixture.workflow(fixture.after).Revalidate(
		t.Context(),
		fixture.after,
		"clean",
		[]string{fixture.mergeCommit},
	)
	if err != nil {
		t.Fatalf("revalidate after merge: %v", err)
	}

	for _, ignored := range []string{"internal/changed_test.go", "docs/guide.md", "internal/undeclared.go"} {
		if containsRevalidationString(result.ChangedPremises, ignored) {
			t.Fatalf("changed premises = %v, want %q ignored", result.ChangedPremises, ignored)
		}
	}
	if got, want := result.ChangedPremises, []string{"internal/changed.go", "internal/removed.go"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("changed premises = %v, want declared production paths %v", got, want)
	}
}

func TestRevalidateRefusesAnUnreadableMergeCommit(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryRevalidationGitFixture(t)

	_, err := fixture.workflow(fixture.after).Revalidate(
		t.Context(),
		fixture.after,
		"clean",
		[]string{"unreadable-merge"},
	)
	if err == nil || !strings.Contains(err.Error(), `read prior merge "unreadable-merge"`) {
		t.Fatalf("unreadable merge error = %v, want named merge", err)
	}
}

func TestDeliverStatusPrintsAnItemWarning(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{implementTestSlug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Warning = "premise-changed: internal/delivery/engine.go (merge abc123)"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repoDir, item); err != nil {
		t.Fatalf("record item warning: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)

	want := implementTestSlug + "\tqueued\t-\t-\n" +
		"Warning: " + implementTestSlug + " premise-changed: internal/delivery/engine.go (merge abc123)\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("deliver status exit=%d stdout=%q stderr=%q, want stdout %q", code, stdout.String(), stderr.String(), want)
	}
}

func TestDeliverStatusPrintsNoWarningLineWithoutAWarning(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	if _, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{implementTestSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)

	want := implementTestSlug + "\tqueued\t-\t-\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("deliver status exit=%d stdout=%q stderr=%q, want unchanged stdout %q", code, stdout.String(), stderr.String(), want)
	}
}

type deliveryRevalidationGitFixture struct {
	before      string
	after       string
	mergeCommit string
}

func newDeliveryRevalidationGitFixture(t *testing.T) deliveryRevalidationGitFixture {
	t.Helper()
	_, source := newSpecCheckWorkspace(t, "clean")
	specDir := filepath.Join(source, "docs", "specs", "clean")
	mustWrite(t, filepath.Join(specDir, "_tasks.md"), `---
schema: spec-tasks/v1
spec: clean
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | backend | Revalidation fixture |
`)
	mustWrite(t, filepath.Join(specDir, "task_01.md"), `---
task: task_01
spec: clean
status: pending
type: backend
complexity: low
---

# Task 01: Revalidation fixture

## Context

- interface: `+"`internal/removed.go`"+`
- interface: `+"`internal/changed.go`"+`
- interface: `+"`internal/changed_test.go`"+`
- interface: `+"`docs/guide.md`"+`

## Verification

- `+"`test -f internal/changed.go`"+` — expected: pass.
`)
	if err := os.MkdirAll(filepath.Join(source, "internal"), 0o755); err != nil {
		t.Fatalf("create fixture production directory: %v", err)
	}
	mustWrite(t, filepath.Join(source, "internal", "removed.go"), "package internal\n")
	mustWrite(t, filepath.Join(source, "internal", "changed.go"), "package internal\n")
	mustWrite(t, filepath.Join(source, "internal", "changed_test.go"), "package internal\n")
	mustWrite(t, filepath.Join(source, "docs", "guide.md"), "# Guide\n")
	mustWrite(t, filepath.Join(source, "internal", "undeclared.go"), "package internal\n")
	gittest.Run(t, source, "add", "-A")
	gittest.Run(t, source, "commit", "-m", "test: seed Delivery Revalidation fixture")
	gittest.Run(t, source, "switch", "main")
	gittest.Run(t, source, "merge", "--ff-only", "ma/spec-check")

	origin := filepath.Join(t.TempDir(), "origin.git")
	gittest.InitRepo(t, origin, "--bare", "--initial-branch=main")
	gittest.Run(t, source, "remote", "add", "origin", origin)
	gittest.Run(t, source, "push", "-u", "origin", "main")

	before := filepath.Join(t.TempDir(), "before")
	gittest.Run(t, "", "clone", origin, before)
	gittest.Harden(t, before)

	if err := os.Remove(filepath.Join(source, "internal", "removed.go")); err != nil {
		t.Fatalf("remove declared production file: %v", err)
	}
	mustWrite(t, filepath.Join(source, "internal", "changed.go"), "package internal\n\nconst changed = true\n")
	mustWrite(t, filepath.Join(source, "internal", "changed_test.go"), "package internal\n\nconst changedTest = true\n")
	mustWrite(t, filepath.Join(source, "docs", "guide.md"), "# Changed guide\n")
	mustWrite(t, filepath.Join(source, "internal", "undeclared.go"), "package internal\n\nconst undeclaredChanged = true\n")
	gittest.Run(t, source, "add", "-A")
	gittest.Run(t, source, "commit", "-m", "feat: merge earlier queue item")
	mergeCommit := strings.TrimSpace(gittest.Run(t, source, "rev-parse", "HEAD"))
	gittest.Run(t, source, "push", "origin", "main")

	after := filepath.Join(t.TempDir(), "after")
	gittest.Run(t, "", "clone", origin, after)
	gittest.Harden(t, after)
	return deliveryRevalidationGitFixture{before: before, after: after, mergeCommit: mergeCommit}
}

func (fixture deliveryRevalidationGitFixture) workflow(workDir string) *commandDeliveryWorkflow {
	return &commandDeliveryWorkflow{
		loaded: roundconfig.Loaded{
			GitRoot: workDir,
			Config:  roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}},
		},
		git: preflight.ExecGitRunner{},
	}
}

func containsRevalidationString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
