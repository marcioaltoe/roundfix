// Suite: delivery recovery evidence at real local Git boundaries.
// Boundary IN: isolated repositories, Run Database, QA reader, and legacy archive fixtures.
// Boundary OUT: no GitHub calls or external Agent execution.
package cli

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch(t *testing.T) {
	t.Parallel()
	const prePRRow = "\n## Results\n\n| ID | Provenance | Status |\n| --- | --- | --- |\n| PR | Pull Request row | blocked (environment: no open Pull Request) |\n"
	tests := []struct {
		name, content string
		want          bool
	}{
		{"environment", "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 1\n---\nQA environment row\n", true},
		{"finding", "---\nverdict: partial\nrows_blocked_finding: 1\nrows_blocked_environment: 1\n---\nQA finding row\n", false},
		{"missing", "", false},
		{"unreadable", "---\nverdict: partial\nrows_blocked_environment: invalid\n---\n", false},
		{"pre-PR only", "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 1\n---\n" + prePRRow, false},
		{"environment beyond pre-PR", "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 2\n---\n" + prePRRow, true},
		{"pass", "---\nverdict: pass\nrows_blocked_finding: 0\nrows_blocked_environment: 1\n---\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
			itemHead := itemRecoveryHead(t, repo)
			branch := store.RunBranchPrefix + "qa-partial"
			gittest.Run(t, repo, "checkout", "-b", branch)
			qaDir := filepath.Join(repo, "docs", "specs", implementTestSlug, "qa")
			mustMkdir(t, qaDir)
			if tt.content != "" {
				mustWrite(t, filepath.Join(qaDir, "qa-report-2026-09-29-01.md"), "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 1\n---\nold environment evidence\n")
				mustWrite(t, filepath.Join(qaDir, "qa-report-2026-09-30-02.md"), tt.content)
				gittest.Run(t, repo, "add", "docs/specs")
				gittest.Run(t, repo, "commit", "-m", "docs: Run QA report")
			}
			gittest.Run(t, repo, "checkout", "--detach", itemHead)
			// Exercise the RunSpec outcome boundary without launching an Agent executor.
			result, err := workflow.runResult(t.Context(), repo, implementTestSlug, roundfixCommandResult{exitCode: exitRunFailed}, "", nil, &store.Run{ID: "qa-partial", State: store.StateUnresolved})
			if err != nil {
				t.Fatal(err)
			}
			if result.Outcome != delivery.RunOutcomeUnresolved || result.QAEnvironmentPartial != tt.want {
				t.Fatalf("result=%+v want environment partial %v", result, tt.want)
			}
		})
	}
}

func TestInspectItemReadsTheQAOverrideOfTheArchivedSpec(t *testing.T) {
	t.Parallel()
	for _, override := range []bool{true, false} {
		t.Run(map[bool]string{true: "override", false: "normal archive"}[override], func(t *testing.T) {
			repo, _, _ := commitLinkRewritingArchive(t, nil, nil)
			workflow := newItemRecoveryWorkflowForRepository(t, t.TempDir(), repo)
			_, destination, err := workflow.archivePaths(repo, implementTestSlug)
			if err != nil {
				t.Fatal(err)
			}
			if override {
				path := filepath.Join(repo, filepath.FromSlash(destination), "_prd.md")
				content := mustRead(t, path)
				mustWrite(t, path, strings.Replace(content, "---\n", "---\nqa_override: true\n", 1))
			}
			state, err := workflow.InspectItem(t.Context(), repo, implementTestSlug)
			if err != nil {
				t.Fatal(err)
			}
			if !state.Archived || state.QAOverride != override {
				t.Fatalf("state=%+v", state)
			}
		})
	}
}

func TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit(t *testing.T) {
	t.Parallel()
	for _, postArchive := range []bool{false, true} {
		t.Run(map[bool]string{false: "archive head", true: "operator commit after archive"}[postArchive], func(t *testing.T) {
			repo, _, _ := commitLinkRewritingArchive(t, nil, nil)
			workflow := newItemRecoveryWorkflowForRepository(t, t.TempDir(), repo)
			// Install an explicit delivery grant in a new active -> archive transition.
			source, destination, err := workflow.archivePaths(repo, implementTestSlug)
			if err != nil {
				t.Fatal(err)
			}
			gittest.Run(t, repo, "mv", destination, source)
			mustWrite(t, filepath.Join(repo, filepath.FromSlash(source), "_authorization.md"), implementFixtureAuthorization(implementTestSlug, "implement", "commit", "push", "pull_request", "merge"))
			gittest.Run(t, repo, "add", "-A")
			gittest.Run(t, repo, "commit", "-m", "docs: active delivery grant")
			gittest.Run(t, repo, "mv", source, destination)
			gittest.Run(t, repo, "commit", "-m", "docs: archive delivery grant")
			if postArchive {
				// A later revoked archived grant must not replace the pre-archive authority.
				mustWrite(t, filepath.Join(repo, filepath.FromSlash(destination), "_authorization.md"), implementFixtureAuthorization(implementTestSlug, "implement"))
				gittest.Run(t, repo, "add", "-A")
				gittest.Run(t, repo, "commit", "-m", "docs: post archive operator change")
			}
			authorization, err := workflow.Authorization(t.Context(), repo, implementTestSlug)
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"push", "pull_request", "merge"}; !reflect.DeepEqual(authorization.Operations, want) {
				t.Fatalf("operations=%v want %v", authorization.Operations, want)
			}
		})
	}
}

func TestArchiveReportsASpecAlreadyArchivedAtTheReviewedHead(t *testing.T) {
	t.Parallel()
	repo, _, head := commitRealArchive(t, nil)
	workflow := newItemRecoveryWorkflowForRepository(t, t.TempDir(), repo)
	result, err := workflow.Archive(t.Context(), repo, implementTestSlug, head)
	if err != nil {
		t.Fatal(err)
	}
	if !result.AlreadyArchived || result.Head != head || itemRecoveryHead(t, repo) != head {
		t.Fatalf("archive=%+v current=%s", result, itemRecoveryHead(t, repo))
	}
	if dirty := gittest.Run(t, repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty archive=%q", dirty)
	}
}

func TestOperatorArchiveHistoryReadsTheRunStartAndAncestry(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
	start := itemRecoveryHead(t, repo)
	run, err := workflow.store.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: "main", HeadSHA: start, SpecSlug: implementTestSlug, Agent: "codex", ArtifactDir: t.TempDir(), WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	got, err := workflow.RunStart(t.Context(), repo, run.ID)
	if err != nil || got != start {
		t.Fatalf("start=%q err=%v", got, err)
	}
	mustWrite(t, filepath.Join(repo, "operator.txt"), "operator\n")
	gittest.Run(t, repo, "add", "operator.txt")
	gittest.Run(t, repo, "commit", "-m", "docs: operator update")
	head := itemRecoveryHead(t, repo)
	if ok, err := workflow.Descends(t.Context(), repo, start, head); err != nil || !ok {
		t.Fatalf("descends=%v err=%v", ok, err)
	}
	if ok, err := workflow.Descends(t.Context(), repo, head, start); err != nil || ok {
		t.Fatalf("reverse descends=%v err=%v", ok, err)
	}
	gittest.Run(t, repo, "checkout", "--orphan", "test/unrelated-history")
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "test: unrelated root")
	unrelated := itemRecoveryHead(t, repo)
	if ok, err := workflow.Descends(t.Context(), repo, start, unrelated); err != nil || ok {
		t.Fatalf("unrelated descends=%v err=%v", ok, err)
	}
	if _, err := workflow.Descends(t.Context(), repo, "missing-revision", unrelated); err == nil {
		t.Fatal("missing ancestor must report a Git error")
	}

}
