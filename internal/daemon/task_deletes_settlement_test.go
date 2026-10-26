package daemon

// Suite: declared deletions through gated Task settlement.
// Invariant: missing working-tree sources settle; broken links and unauthorized deletions remain findings.
// Boundary IN: synthetic local Git repositories and fake Agent deletion turns.
// Boundary OUT: TaskCycle outcomes and default settlement checker findings.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const settlementDeletedSource = "packages/backend/test/gone.test.ts"

func retiredHistoryDeletionFixture(t *testing.T) (*taskCycleFixture, TaskPlan, *taskFakeRunner, *taskFakeVerifier, *engineFakeCommitter) {
	t.Helper()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	files := map[string]string{
		"docs/agents/setup-context.json":          `{"managedArtifacts":[{"id":"guide.docs-layout","path":"docs/agents/docs-layout.md"}]}`,
		"docs/agents/docs-layout.md":              "<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->\n" + spec.RetirementRuleSentinel + "\n<!-- setup-context-driven:end id=guide.docs-layout -->\n",
		"docs/history/findings/2026-01-01-old.md": "# Retired finding\n",
		"docs/user-guide/index.md":                "# Guide\n\n[Guide](guide.md)\n",
		"docs/user-guide/guide.md":                "# User guide\n",
		settlementDeletedSource:                   "// Obsolete test.\n",
	}
	for path, content := range files {
		full := filepath.Join(fixture.gitRoot, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteForTest(t, full, content)
	}
	taskPath := taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01")
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, taskPath, string(content)+"\n## Context\n\n- deletes: `"+settlementDeletedSource+"`\n")
	commitTaskFixtureSource(t, fixture.gitRoot, "test: seed retired-history deletion")
	gittest.Run(t, fixture.gitRoot, "rm", "docs/history/findings/2026-01-01-old.md")
	gittest.Run(t, fixture.gitRoot, "commit", "-m", "test: retire old finding")
	fixture.reloadGraph()
	plan.Tasks = fixture.graph.Tasks
	for i := range plan.Tasks {
		if plan.Tasks[i].Type == spec.TaskTypeQA {
			plan.Tasks[i].Status = spec.StatusCompleted
		}
	}
	plan.HeadSHA = strings.TrimSpace(gittest.Run(t, fixture.gitRoot, "rev-parse", "HEAD"))
	plan.Authorization = spec.ReadSpecAuthorization(context.Background(), plan.WorkDir, plan.SpecsRoot, plan.Spec.Slug, plan.HeadSHA)
	fixture.worktree.snapshots = [][]string{nil, {settlementDeletedSource}}
	return fixture, plan, runner, verifier, committer
}

func TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"disk-only", "staged"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			fixture, plan, runner, verifier, committer := retiredHistoryDeletionFixture(t)
			runner.afterTask = func(string) {
				if mode == "staged" {
					cmd := exec.CommandContext(t.Context(), "git", "-C", fixture.gitRoot, "rm", "-q", "--ignore-unmatch", "--", settlementDeletedSource)
					cmd.Env = gittest.IsolatedEnv()
					if out, err := cmd.CombinedOutput(); err != nil {
						t.Errorf("stage deletion: %v\n%s", err, out)
					}
					return
				}
				if err := os.Remove(filepath.Join(fixture.gitRoot, filepath.FromSlash(settlementDeletedSource))); err != nil && !errors.Is(err, os.ErrNotExist) {
					t.Errorf("remove declared source: %v", err)
				}
			}
			result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
			if err != nil || result.Completed != 1 || result.Failed != 0 {
				t.Fatalf("TaskCycle = %+v, %v", result, err)
			}
		})
	}
}

func TestSettlementStillReportsALinkToADeletedFile(t *testing.T) {
	t.Parallel()
	fixture, plan, _, _, _ := retiredHistoryDeletionFixture(t)
	if err := os.Remove(filepath.Join(fixture.gitRoot, "docs/user-guide/guide.md")); err != nil {
		t.Fatal(err)
	}
	findings, err := (SpecCheckSettlementChecker{}).RefusingFindings(plan.SpecsRoot, plan.WorkDir, plan.Spec.Slug)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range findings {
		if finding.Code == speccheck.CodeHistoryRetained && strings.Contains(finding.Summary, "docs/user-guide/index.md") && strings.Contains(finding.Summary, "guide.md") {
			return
		}
	}
	t.Fatalf("missing broken-link finding: %+v", findings)
}

func TestSettlementRefusesAnUnauthorizedGovernedDeletion(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	const grant = "docs/specs/check/_authorization.md"
	const escaped = ".agents/skills/escaped/SKILL.md"
	for path, content := range map[string]string{
		grant:   "---\nstatus: approved\ngranted: 2026-09-09\naction: test\nconsuming: check\npaths:\n  - Makefile\n---\n",
		escaped: "# Tracked skill\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteForTest(t, full, content)
	}
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "test: seed authorization and tracked skill")
	parent := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(escaped))); err != nil {
		t.Fatal(err)
	}
	findings, err := (SpecCheckSettlementChecker{}).AuditCommit(context.Background(), speccheck.MechanicalRequest{RepoRoot: root, AuthorizationPath: grant, ConsumingSpec: "check", DeliveryTargetRevision: parent}, speccheck.ProspectiveTaskCommit{TaskID: "task_02", Parent: parent, Changed: []string{escaped}})
	if err != nil || len(findings) != 1 {
		t.Fatalf("findings %+v, err %v", findings, err)
	}
	if findings[0].Code != speccheck.CodeMechanicalAuthPaths || !strings.Contains(findings[0].Detail, escaped) {
		t.Fatalf("unauthorized deletion finding = %+v", findings[0])
	}
}
