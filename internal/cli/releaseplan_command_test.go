package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
)

// Suite: Release Plan Command
// Invariant: public CLI planning is deterministic, read-only, and stdout/stderr isolated.
// Boundary IN: roundfix release plan flags and temporary Git repositories.
// Boundary OUT: release mutation, Run creation, Roundfix configuration, network services.

func TestReleasePlanCommandMatchesPRDOutcomes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		baseTag              string
		commits              []releasePlanCommandCommit
		args                 []string
		wantExit             int
		wantState            releaseplan.State
		wantImpact           releaseplan.Impact
		wantProposed         string
		wantApprovalQuestion string
		wantBreaking         bool
		wantBlocking         int
		wantManualReason     string
	}{
		{
			name:     "fix only proposes patch without approval",
			baseTag:  "v0.4.0",
			commits:  []releasePlanCommandCommit{{subject: "fix: correct release output", paths: []string{"internal/cli/release.go"}}},
			wantExit: exitOK, wantState: releaseplan.StateReady, wantImpact: releaseplan.ImpactPatch, wantProposed: "v0.4.1",
		},
		{
			name:     "compatible feature requires minor approval",
			baseTag:  "v0.4.1",
			commits:  []releasePlanCommandCommit{{subject: "feat: expose release plan", paths: []string{"internal/cli/release.go"}}},
			wantExit: exitUnverified, wantState: releaseplan.StateApprovalRequired, wantImpact: releaseplan.ImpactMinor, wantProposed: "v0.5.0",
			wantApprovalQuestion: "Approve the minor increment to v0.5.0?",
		},
		{
			name:     "major breaking requires major approval",
			baseTag:  "v1.4.2",
			commits:  []releasePlanCommandCommit{{subject: "feat!: remove legacy release command", paths: []string{"internal/cli/release.go"}}},
			wantExit: exitUnverified, wantState: releaseplan.StateApprovalRequired, wantImpact: releaseplan.ImpactMajor, wantProposed: "v2.0.0",
			wantApprovalQuestion: "Approve the major increment to v2.0.0?", wantBreaking: true,
		},
		{
			name:     "version zero breaking maps to minor approval",
			baseTag:  "v0.7.3",
			commits:  []releasePlanCommandCommit{{subject: "fix!: change agent contract", paths: []string{"internal/cli/release.go"}}},
			wantExit: exitUnverified, wantState: releaseplan.StateApprovalRequired, wantImpact: releaseplan.ImpactMajor, wantProposed: "v0.8.0",
			wantApprovalQuestion: "Approve the minor increment to v0.8.0?", wantBreaking: true,
		},
		{
			name:    "maintenance only produces no release",
			baseTag: "v0.4.1",
			commits: []releasePlanCommandCommit{
				{subject: "docs: record release plan evidence", paths: []string{"docs/specs/0034-release-plan/task_04.md"}},
				{subject: "test: add release fixture", paths: []string{"internal/releaseplan/testdata/plan.json"}},
				{subject: "ci: verify release plan", paths: []string{".github/workflows/ci-conventions.yml"}},
			},
			wantExit: exitOK, wantState: releaseplan.StateNoRelease, wantImpact: releaseplan.ImpactNone,
		},
		{
			name:         "ambiguous shipped change requires manual classification",
			baseTag:      "v0.4.1",
			commits:      []releasePlanCommandCommit{{subject: "chore: tune release command", paths: []string{"internal/cli/release.go"}}},
			wantExit:     exitUnverified,
			wantState:    releaseplan.StateManualClassificationRequired,
			wantImpact:   releaseplan.ImpactNone,
			wantBlocking: 1,
		},
		{
			name:    "manual classification records reason and still gates minor approval",
			baseTag: "v0.4.1",
			commits: []releasePlanCommandCommit{
				{subject: "chore: tune release command", paths: []string{"internal/cli/release.go"}},
			},
			args:     []string{"--impact", "minor", "--reason", "public release command behavior changed"},
			wantExit: exitUnverified, wantState: releaseplan.StateApprovalRequired, wantImpact: releaseplan.ImpactMinor, wantProposed: "v0.5.0",
			wantApprovalQuestion: "Approve the minor increment to v0.5.0?", wantBlocking: 1, wantManualReason: "public release command behavior changed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := newReleasePlanCommandRepo(t, tt.baseTag, tt.commits...)
			args := append([]string{"--format", "json"}, tt.args...)
			code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, args...)

			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, tt.wantExit, stdout, stderr)
			}
			if stderr != "" {
				t.Fatalf("expected no stderr for valid plan, got %q", stderr)
			}
			plan := decodeReleasePlanJSON(t, stdout)
			if plan.SchemaVersion != releaseplan.SchemaVersion {
				t.Fatalf("schemaVersion = %q, want %q", plan.SchemaVersion, releaseplan.SchemaVersion)
			}
			if plan.State != tt.wantState || plan.Classification.Impact != tt.wantImpact {
				t.Fatalf("state/impact = %s/%s, want %s/%s", plan.State, plan.Classification.Impact, tt.wantState, tt.wantImpact)
			}
			if plan.ProposedVersion != tt.wantProposed {
				t.Fatalf("proposedVersion = %q, want %q", plan.ProposedVersion, tt.wantProposed)
			}
			if plan.Classification.Breaking != tt.wantBreaking {
				t.Fatalf("breaking = %v, want %v", plan.Classification.Breaking, tt.wantBreaking)
			}
			if plan.Approval.Question != tt.wantApprovalQuestion {
				t.Fatalf("approval question = %q, want %q", plan.Approval.Question, tt.wantApprovalQuestion)
			}
			if len(plan.Classification.BlockingCommits) != tt.wantBlocking {
				t.Fatalf("blocking commits = %v, want %d", plan.Classification.BlockingCommits, tt.wantBlocking)
			}
			if plan.Classification.ManualReason != tt.wantManualReason {
				t.Fatalf("manual reason = %q, want %q", plan.Classification.ManualReason, tt.wantManualReason)
			}
		})
	}
}

func TestReleasePlanCommandMixedOrderSelectsHighestImpact(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		commits []releasePlanCommandCommit
	}{
		{
			name: "fix then feature",
			commits: []releasePlanCommandCommit{
				{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
				{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}},
			},
		},
		{
			name: "feature then fix",
			commits: []releasePlanCommandCommit{
				{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}},
				{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := newReleasePlanCommandRepo(t, "v0.4.1", tt.commits...)
			code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, "--format", "json")
			if code != exitUnverified {
				t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout, stderr)
			}
			plan := decodeReleasePlanJSON(t, stdout)
			if plan.State != releaseplan.StateApprovalRequired || plan.ProposedVersion != "v0.5.0" {
				t.Fatalf("state/proposed = %s/%s, want approval_required/v0.5.0", plan.State, plan.ProposedVersion)
			}
			if plan.Classification.Impact != releaseplan.ImpactMinor {
				t.Fatalf("impact = %s, want minor", plan.Classification.Impact)
			}
		})
	}
}

func TestReleasePlanTextPrintsOnlyDeterminingOrBlockingCommits(t *testing.T) {
	t.Parallel()
	t.Run("determining commits", func(t *testing.T) {
		repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
			releasePlanCommandCommit{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
			releasePlanCommandCommit{subject: "docs: record planning evidence", paths: []string{"docs/specs/0034-release-plan/task_04.md"}},
			releasePlanCommandCommit{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}},
		)

		code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir)

		if code != exitUnverified {
			t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout, stderr)
		}
		assertReleasePlanNoStderr(t, stderr)
		if !strings.HasPrefix(stdout, "Decision: approval_required\n") {
			t.Fatalf("text output should lead with decision, got %q", stdout)
		}
		if !strings.Contains(stdout, "feat: expose release plan") {
			t.Fatalf("expected determining feature commit in text output, got %q", stdout)
		}
		for _, unexpected := range []string{"fix: correct release output", "docs: record planning evidence"} {
			if strings.Contains(stdout, unexpected) {
				t.Fatalf("text output included non-determining commit %q: %q", unexpected, stdout)
			}
		}
	})

	t.Run("blocking commits", func(t *testing.T) {
		repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
			releasePlanCommandCommit{subject: "docs: record planning evidence", paths: []string{"docs/specs/0034-release-plan/task_04.md"}},
			releasePlanCommandCommit{subject: "chore: tune release command", paths: []string{"internal/cli/release.go"}},
		)

		code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir)

		if code != exitUnverified {
			t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout, stderr)
		}
		assertReleasePlanNoStderr(t, stderr)
		if !strings.Contains(stdout, "Blocking commits:") || !strings.Contains(stdout, "chore: tune release command") {
			t.Fatalf("expected blocking commit in text output, got %q", stdout)
		}
		if strings.Contains(stdout, "docs: record planning evidence") {
			t.Fatalf("text output included non-blocking maintenance commit: %q", stdout)
		}
		if !strings.Contains(stdout, "rerun roundfix release plan --from v0.4.1 --to HEAD --impact <none|patch|minor|major> --reason <text>") {
			t.Fatalf("expected explicit manual rerun shape, got %q", stdout)
		}
	})
}

func TestReleasePlanJSONIncludesEveryCommitEvidence(t *testing.T) {
	t.Parallel()
	repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
		releasePlanCommandCommit{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
		releasePlanCommandCommit{subject: "docs: record planning evidence", paths: []string{"docs/specs/0034-release-plan/task_04.md"}},
		releasePlanCommandCommit{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}},
	)

	code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, "--format", "json")

	if code != exitUnverified {
		t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout, stderr)
	}
	assertReleasePlanNoStderr(t, stderr)
	plan := decodeReleasePlanJSON(t, stdout)
	if len(plan.Changes) != 3 {
		t.Fatalf("changes length = %d, want evidence for every commit: %+v", len(plan.Changes), plan.Changes)
	}
	for _, subject := range []string{"fix: correct release output", "docs: record planning evidence", "feat: expose release plan"} {
		if !releasePlanJSONHasSubject(plan, subject) {
			t.Fatalf("JSON changes missing subject %q: %+v", subject, plan.Changes)
		}
	}
	if strings.Contains(stdout, "roundfix:") || stderr != "" {
		t.Fatalf("JSON result leaked diagnostics stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestReleasePlanExitCodesAndInvalidInputIsolation(t *testing.T) {
	t.Parallel()
	t.Run("ready and no release exit zero", func(t *testing.T) {
		for _, tt := range []struct {
			name    string
			commits []releasePlanCommandCommit
		}{
			{name: "ready", commits: []releasePlanCommandCommit{{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}}}},
			{name: "no release", commits: []releasePlanCommandCommit{{subject: "docs: record planning evidence", paths: []string{"docs/specs/0034-release-plan/task_04.md"}}}},
		} {
			t.Run(tt.name, func(t *testing.T) {
				repoDir := newReleasePlanCommandRepo(t, "v0.4.1", tt.commits...)
				code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir)
				if code != exitOK {
					t.Fatalf("exit = %d, want 0 stdout=%q stderr=%q", code, stdout, stderr)
				}
				assertReleasePlanNoStderr(t, stderr)
				if !strings.HasPrefix(stdout, "Decision: ") {
					t.Fatalf("valid plan should write stdout decision, got %q", stdout)
				}
			})
		}
	})

	t.Run("decision gates exit three", func(t *testing.T) {
		for _, tt := range []struct {
			name    string
			commits []releasePlanCommandCommit
		}{
			{name: "approval", commits: []releasePlanCommandCommit{{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}}}},
			{name: "manual", commits: []releasePlanCommandCommit{{subject: "chore: tune release command", paths: []string{"internal/cli/release.go"}}}},
		} {
			t.Run(tt.name, func(t *testing.T) {
				repoDir := newReleasePlanCommandRepo(t, "v0.4.1", tt.commits...)
				code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir)
				if code != exitUnverified {
					t.Fatalf("exit = %d, want 3 stdout=%q stderr=%q", code, stdout, stderr)
				}
				assertReleasePlanNoStderr(t, stderr)
				if !strings.HasPrefix(stdout, "Decision: ") {
					t.Fatalf("gated plan should remain a stdout result, got %q", stdout)
				}
			})
		}
	})

	t.Run("invalid input exits two with no partial plan", func(t *testing.T) {
		repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
			releasePlanCommandCommit{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
		)
		code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, "--format", "xml")
		if code != exitPreflight {
			t.Fatalf("exit = %d, want 2 stdout=%q stderr=%q", code, stdout, stderr)
		}
		if stdout != "" {
			t.Fatalf("invalid input emitted partial stdout plan: %q", stdout)
		}
		assertReleasePlanOneDiagnostic(t, stderr, "unsupported --format")
	})

	t.Run("manual downgrade exits two with no partial plan", func(t *testing.T) {
		repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
			releasePlanCommandCommit{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}},
		)
		code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, "--impact", "patch", "--reason", "try to lower feature impact")
		if code != exitPreflight {
			t.Fatalf("exit = %d, want 2 stdout=%q stderr=%q", code, stdout, stderr)
		}
		if stdout != "" {
			t.Fatalf("manual downgrade emitted partial stdout plan: %q", stdout)
		}
		assertReleasePlanOneDiagnostic(t, stderr, "at least the automatic minimum")
	})
}

func TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary(t *testing.T) {
	t.Parallel()
	var rootStdout, rootStderr bytes.Buffer
	rootCode := runCLI(t, []string{"--help"}, &rootStdout, &rootStderr)
	if rootCode != exitOK {
		t.Fatalf("root help exit = %d, want 0 stderr=%q", rootCode, rootStderr.String())
	}
	if !strings.Contains(rootStdout.String(), "roundfix release plan") || !strings.Contains(rootStdout.String(), "release    Plan the next release version") {
		t.Fatalf("root help missing release plan command: %q", rootStdout.String())
	}

	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"release", "plan", "--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("release plan help exit = %d, want 0 stderr=%q", code, stderr.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected help stderr empty, got %q", stderr.String())
	}
	help := stdout.String()
	for _, want := range []string{
		"--from", "--to", "--impact", "--reason", "--format",
		"defaults to the latest reachable stable", "--to defaults to committed HEAD",
		"ready", "approval_required", "manual_classification_required", "no_release",
		"read-only", "creates no Run", "reads no Roundfix configuration", "never mutates",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("release plan help missing %q: %q", want, help)
		}
	}
}

func TestReleasePlanReadOnlyPreservesRepositoryForOutcomes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		commits  []releasePlanCommandCommit
		args     []string
		wantExit int
	}{
		{
			name:     "ready",
			commits:  []releasePlanCommandCommit{{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}}},
			wantExit: exitOK,
		},
		{
			name:     "approval",
			commits:  []releasePlanCommandCommit{{subject: "feat: expose release plan", paths: []string{"internal/cli/feature.go"}}},
			wantExit: exitUnverified,
		},
		{
			name:     "manual required",
			commits:  []releasePlanCommandCommit{{subject: "chore: tune release command", paths: []string{"internal/cli/release.go"}}},
			wantExit: exitUnverified,
		},
		{
			name:     "invalid range",
			commits:  []releasePlanCommandCommit{{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}}},
			args:     []string{"--from", "release-1"},
			wantExit: exitPreflight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := newReleasePlanCommandRepo(t, "v0.4.1", tt.commits...)
			code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, tt.args...)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d stdout=%q stderr=%q", code, tt.wantExit, stdout, stderr)
			}
			if code == exitPreflight && stdout != "" {
				t.Fatalf("invalid command emitted stdout plan: %q", stdout)
			}
		})
	}
}

func TestReleasePlanDirtyTreeBlocksWithActionableDiagnostic(t *testing.T) {
	t.Parallel()
	repoDir := newReleasePlanCommandRepo(t, "v0.4.1",
		releasePlanCommandCommit{subject: "fix: correct release output", paths: []string{"internal/cli/fix.go"}},
	)
	writeReleasePlanFile(t, repoDir, "internal/cli/fix.go", "dirty tracked change\n")
	writeReleasePlanFile(t, repoDir, "scratch.txt", "untracked change\n")

	code, stdout, stderr := runReleasePlanCommandInRepo(t, repoDir, "--format", "json")

	if code != exitPreflight {
		t.Fatalf("exit = %d, want 2 stdout=%q stderr=%q", code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("dirty tree emitted partial stdout plan: %q", stdout)
	}
	for _, want := range []string{"internal/cli/fix.go", "scratch.txt", "commit, stash, or remove"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("dirty diagnostic missing %q: %q", want, stderr)
		}
	}
	assertReleasePlanOneDiagnostic(t, stderr, "dirty worktree")
}

func TestReleasePlanResetTextAndJSONInventoryMatchThroughRunBoundary(t *testing.T) {
	t.Parallel()
	gitRunner := newResetPlanRecordingGitRunner()
	ghRunner := &resetPlanRecordingGHRunner{
		output: `[
			[
				{"id":30,"node_id":"RE_node_30","name":"Third release","tag_name":"v0.3.0","target_commitish":"main"},
				{"id":10,"node_id":"RE_node_10","name":"First release","tag_name":"v0.1.0","target_commitish":"main"}
			],
			[
				{"id":20,"node_id":"RE_node_20","name":"Second release","tag_name":"v0.2.0","target_commitish":"main"}
			]
		]`,
	}
	restore := setReleasePlanCommandRunnersForTest(t, gitRunner, ghRunner)
	t.Cleanup(restore)

	var textStdout, textStderr bytes.Buffer
	textExit := runCLIContext(t, context.Background(), []string{"release", "plan", "--reset-to", "v0.0.1"}, &textStdout, &textStderr)
	if textExit != exitUnverified {
		t.Fatalf("text exit = %d, want 3 stdout=%q stderr=%q", textExit, textStdout.String(), textStderr.String())
	}
	assertReleasePlanNoStderr(t, textStderr.String())

	var jsonStdout, jsonStderr bytes.Buffer
	jsonExit := runCLIContext(t, context.Background(), []string{"release", "plan", "--reset-to", "v0.0.1", "--format", "json"}, &jsonStdout, &jsonStderr)
	if jsonExit != exitUnverified {
		t.Fatalf("JSON exit = %d, want 3 stdout=%q stderr=%q", jsonExit, jsonStdout.String(), jsonStderr.String())
	}
	assertReleasePlanNoStderr(t, jsonStderr.String())
	plan := decodeReleaseResetPlanJSON(t, jsonStdout.String())

	if plan.State != releaseplan.StateApprovalRequired || !plan.Approval.Required {
		t.Fatalf("decision = state:%q approval:%+v, want approval_required and required", plan.State, plan.Approval)
	}
	if plan.TargetVersion != "v0.0.1" || plan.Target.Name != "HEAD" || plan.Target.CommitSHA != resetPlanTargetCommit {
		t.Fatalf("target = version:%q ref:%+v, want v0.0.1 HEAD at %s", plan.TargetVersion, plan.Target, resetPlanTargetCommit)
	}
	if plan.PlanDigest == "" || !strings.Contains(textStdout.String(), "Plan digest: "+plan.PlanDigest) {
		t.Fatalf("text and JSON digests differ: text=%q JSON=%q", textStdout.String(), plan.PlanDigest)
	}
	if len(plan.Tags) != 5 {
		t.Fatalf("tags = %+v, want every two local and three remote stable tags exactly once", plan.Tags)
	}
	if len(plan.Releases) != 3 {
		t.Fatalf("releases = %+v, want every release from both pages exactly once", plan.Releases)
	}
	if got := []int64{plan.Releases[0].ID, plan.Releases[1].ID, plan.Releases[2].ID}; !reflect.DeepEqual(got, []int64{10, 20, 30}) {
		t.Fatalf("release IDs = %v, want deterministic tag order 10, 20, 30", got)
	}
	for _, tag := range plan.Tags {
		for _, value := range []string{tag.ImmutableID, tag.TargetCommit} {
			if value == "" || !strings.Contains(textStdout.String(), value) {
				t.Fatalf("text output missing tag identity value %q from %+v: %q", value, tag, textStdout.String())
			}
		}
	}
	for _, release := range plan.Releases {
		for _, value := range []string{release.ImmutableID, release.NodeID, release.TargetCommit} {
			if value == "" || !strings.Contains(textStdout.String(), value) {
				t.Fatalf("text output missing release identity value %q from %+v: %q", value, release, textStdout.String())
			}
		}
	}
	if gitRunner.mutationCalls != 0 || ghRunner.mutationCalls != 0 {
		t.Fatalf("mutation calls = git:%d GitHub:%d, want zero", gitRunner.mutationCalls, ghRunner.mutationCalls)
	}
	if len(ghRunner.calls) != 2 {
		t.Fatalf("GitHub calls = %v, want one paginated read per text/JSON plan", ghRunner.calls)
	}
	for _, call := range ghRunner.calls {
		want := []string{"api", "--method", "GET", "--paginate", "--slurp", "repos/{owner}/{repo}/releases?per_page=100"}
		if !reflect.DeepEqual(call, want) {
			t.Fatalf("GitHub call = %v, want read-only exhaustive pagination %v", call, want)
		}
	}
}

func TestReleasePlanResetInventoriesTemporaryGitRemoteAndPaginatedGitHubReadOnly(t *testing.T) {
	t.Parallel()
	repoDir := newEmptyReleasePlanGitRepo(t)
	writeReleasePlanFile(t, repoDir, "README.md", "seed\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	commitReleasePlan(t, repoDir, "chore: seed")
	gitReleasePlan(t, repoDir, "tag", "v0.1.0")
	writeReleasePlanFile(t, repoDir, "release.txt", "second\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	commitReleasePlan(t, repoDir, "feat: second release")
	gitReleasePlan(t, repoDir, "tag", "v0.2.0")

	remoteDir := t.TempDir()
	gittest.InitRepo(t, remoteDir, "--bare")
	gitReleasePlan(t, repoDir, "remote", "add", "origin", remoteDir)
	gitReleasePlan(t, repoDir, "push", "origin", "main", "v0.1.0", "v0.2.0")
	gitReleasePlan(t, repoDir, "tag", "v0.3.0")
	gitReleasePlan(t, repoDir, "push", "origin", "v0.3.0")
	gitReleasePlan(t, repoDir, "tag", "-d", "v0.3.0")
	gitReleasePlan(t, repoDir, "tag", "v0.4.0")

	gitRunner := &resetPlanExecGitRunner{delegate: preflight.ExecGitRunner{}}
	ghRunner := &resetPlanRecordingGHRunner{
		output: `[
			[
				{"id":40,"node_id":"RE_node_40","name":"Fourth release","tag_name":"v0.4.0","target_commitish":"main"},
				{"id":10,"node_id":"RE_node_10","name":"First release","tag_name":"v0.1.0","target_commitish":"main"}
			],
			[
				{"id":20,"node_id":"RE_node_20","name":"Second release","tag_name":"v0.2.0","target_commitish":"main"},
				{"id":30,"node_id":"RE_node_30","name":"Third release","tag_name":"v0.3.0","target_commitish":"main"}
			]
		]`,
	}
	restore := setReleasePlanCommandRunnersForTest(t, gitRunner, ghRunner)
	t.Cleanup(restore)
	before := snapshotReleasePlanRepo(t, repoDir)
	setCommandWorkDirForTest(t, repoDir)

	var stdout, stderr bytes.Buffer
	code := runCLIContext(t,
		context.Background(),
		[]string{"release", "plan", "--reset-to", "v0.0.1", "--format", "json"},
		&stdout,
		&stderr)

	if code != exitUnverified {
		t.Fatalf("exit = %d, want 3 stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	assertReleasePlanNoStderr(t, stderr.String())
	plan := decodeReleaseResetPlanJSON(t, stdout.String())
	if plan.State != releaseplan.StateApprovalRequired || !plan.Approval.Required {
		t.Fatalf("decision = state:%q approval:%+v, want approval_required and required", plan.State, plan.Approval)
	}
	if len(plan.Tags) != 6 {
		t.Fatalf("tags = %+v, want three shared/local tags plus one remote-only and one local-only identity", plan.Tags)
	}
	tagInventory := map[string]bool{}
	for _, tag := range plan.Tags {
		tagInventory[string(tag.Source)+":"+tag.Name] = true
	}
	wantTags := []string{
		"local:v0.1.0",
		"local:v0.2.0",
		"local:v0.4.0",
		"remote:v0.1.0",
		"remote:v0.2.0",
		"remote:v0.3.0",
	}
	for _, want := range wantTags {
		if !tagInventory[want] {
			t.Fatalf("tag inventory missing %q: %+v", want, plan.Tags)
		}
	}
	if got := []int64{plan.Releases[0].ID, plan.Releases[1].ID, plan.Releases[2].ID, plan.Releases[3].ID}; !reflect.DeepEqual(got, []int64{10, 20, 30, 40}) {
		t.Fatalf("release IDs = %v, want every paginated release in deterministic tag order", got)
	}
	if gitRunner.mutationCalls != 0 || ghRunner.mutationCalls != 0 {
		t.Fatalf("mutation calls = git:%d GitHub:%d, want zero", gitRunner.mutationCalls, ghRunner.mutationCalls)
	}
	if len(ghRunner.calls) != 1 {
		t.Fatalf("GitHub calls = %v, want one exhaustive paginated read", ghRunner.calls)
	}
	assertReleasePlanRepoUnchanged(t, repoDir, before)
}

func TestReleasePlanResetRejectsConflictingOrMalformedFlagsBeforeInventory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "from", args: []string{"--reset-to", "v0.0.1", "--from", "v0.4.0"}, want: "--reset-to cannot be combined with --from"},
		{name: "to", args: []string{"--reset-to", "v0.0.1", "--to", "main"}, want: "--reset-to cannot be combined with --to"},
		{name: "impact", args: []string{"--reset-to", "v0.0.1", "--impact", "major"}, want: "--reset-to cannot be combined with --impact"},
		{name: "reason", args: []string{"--reset-to", "v0.0.1", "--reason", "reset"}, want: "--reset-to cannot be combined with --reason"},
		{name: "malformed target", args: []string{"--reset-to", "0.0.1"}, want: "expected a stable vMAJOR.MINOR.PATCH tag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRunner := newResetPlanRecordingGitRunner()
			ghRunner := &resetPlanRecordingGHRunner{output: `[[]]`}
			restore := setReleasePlanCommandRunnersForTest(t, gitRunner, ghRunner)
			t.Cleanup(restore)

			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, context.Background(), append([]string{"release", "plan"}, tt.args...), &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("exit = %d, want 2 stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if stdout.String() != "" {
				t.Fatalf("invalid reset flags emitted partial stdout plan: %q", stdout.String())
			}
			assertReleasePlanOneDiagnostic(t, stderr.String(), tt.want)
			if len(gitRunner.calls) != 0 || len(ghRunner.calls) != 0 {
				t.Fatalf("invalid flags reached inventory providers: git=%v GitHub=%v", gitRunner.calls, ghRunner.calls)
			}
		})
	}
}

func TestReleasePlanResetFailsClosedForDirtyOrIncompleteInventory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		configure func(*resetPlanRecordingGitRunner, *resetPlanRecordingGHRunner)
		want      string
	}{
		{
			name: "dirty target",
			configure: func(git *resetPlanRecordingGitRunner, _ *resetPlanRecordingGHRunner) {
				git.dirtyStatus = " M internal/cli/releaseplan_command.go\x00"
			},
			want: "dirty worktree",
		},
		{
			name: "remote inventory unavailable",
			configure: func(git *resetPlanRecordingGitRunner, _ *resetPlanRecordingGHRunner) {
				git.remoteErr = errors.New("remote unavailable")
			},
			want: "inventory remote stable tags",
		},
		{
			name: "GitHub inventory unavailable",
			configure: func(_ *resetPlanRecordingGitRunner, gh *resetPlanRecordingGHRunner) {
				gh.err = errors.New("GitHub unavailable")
			},
			want: "inventory GitHub Releases",
		},
		{
			name: "GitHub pagination output incomplete",
			configure: func(_ *resetPlanRecordingGitRunner, gh *resetPlanRecordingGHRunner) {
				gh.output = `[null]`
			},
			want: "complete paginated GitHub Release inventory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitRunner := newResetPlanRecordingGitRunner()
			ghRunner := &resetPlanRecordingGHRunner{output: `[[]]`}
			tt.configure(gitRunner, ghRunner)
			restore := setReleasePlanCommandRunnersForTest(t, gitRunner, ghRunner)
			t.Cleanup(restore)

			var stdout, stderr bytes.Buffer
			code := runCLIContext(t, context.Background(), []string{"release", "plan", "--reset-to", "v0.0.1", "--format", "json"}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("exit = %d, want 2 stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if stdout.String() != "" {
				t.Fatalf("incomplete reset inventory emitted partial stdout plan: %q", stdout.String())
			}
			assertReleasePlanOneDiagnostic(t, stderr.String(), tt.want)
			if gitRunner.mutationCalls != 0 || ghRunner.mutationCalls != 0 {
				t.Fatalf("mutation calls = git:%d GitHub:%d, want zero", gitRunner.mutationCalls, ghRunner.mutationCalls)
			}
		})
	}
}

type releasePlanCommandCommit struct {
	subject string
	body    string
	paths   []string
}

const resetPlanTargetCommit = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

type resetPlanRecordingGitRunner struct {
	calls         [][]string
	dirtyStatus   string
	remoteErr     error
	mutationCalls int
}

type resetPlanExecGitRunner struct {
	delegate      preflight.GitRunner
	calls         [][]string
	mutationCalls int
}

func (runner *resetPlanExecGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	runner.calls = append(runner.calls, append([]string(nil), args...))
	if resetPlanGitCallMutates(args) {
		runner.mutationCalls++
		return "", errors.New("mutation forbidden")
	}
	return runner.delegate.RunGit(ctx, workDir, args...)
}

func newResetPlanRecordingGitRunner() *resetPlanRecordingGitRunner {
	return &resetPlanRecordingGitRunner{}
}

func (runner *resetPlanRecordingGitRunner) RunGit(_ context.Context, _ string, args ...string) (string, error) {
	call := append([]string(nil), args...)
	runner.calls = append(runner.calls, call)
	if resetPlanGitCallMutates(args) {
		runner.mutationCalls++
		return "", errors.New("mutation forbidden")
	}
	switch strings.Join(args, "\x00") {
	case "rev-parse\x00--show-toplevel":
		return "/fixture/roundfix", nil
	case "--no-optional-locks\x00status\x00--porcelain=v1\x00-z":
		return runner.dirtyStatus, nil
	case "rev-parse\x00--verify\x00HEAD":
		return resetPlanTargetCommit, nil
	case "rev-parse\x00--verify\x00HEAD^{commit}":
		return resetPlanTargetCommit, nil
	case "for-each-ref\x00--format=%(refname)\x00refs/tags":
		return "refs/tags/v0.2.0\nrefs/tags/v0.1.0\nrefs/tags/v0.3.0-rc.1", nil
	case "rev-parse\x00--verify\x00refs/tags/v0.1.0^{commit}":
		return "1111111111111111111111111111111111111111", nil
	case "rev-parse\x00--verify\x00refs/tags/v0.2.0^{commit}":
		return "2222222222222222222222222222222222222222", nil
	case "remote":
		return "origin\nbackup", nil
	case "ls-remote\x00--tags\x00origin":
		if runner.remoteErr != nil {
			return "", runner.remoteErr
		}
		return strings.Join([]string{
			"1111111111111111111111111111111111111111\trefs/tags/v0.1.0",
			"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\trefs/tags/v0.3.0",
			"3333333333333333333333333333333333333333\trefs/tags/v0.3.0^{}",
			"cccccccccccccccccccccccccccccccccccccccc\trefs/tags/v0.4.0-rc.1",
		}, "\n"), nil
	case "ls-remote\x00--tags\x00backup":
		if runner.remoteErr != nil {
			return "", runner.remoteErr
		}
		return "2222222222222222222222222222222222222222\trefs/tags/v0.2.0", nil
	default:
		return "", fmt.Errorf("unexpected read-only git call: %v", args)
	}
}

type resetPlanRecordingGHRunner struct {
	calls         [][]string
	output        string
	err           error
	mutationCalls int
}

func (runner *resetPlanRecordingGHRunner) RunGH(_ context.Context, _ string, args ...string) (string, error) {
	call := append([]string(nil), args...)
	runner.calls = append(runner.calls, call)
	if resetPlanGHCallMutates(args) {
		runner.mutationCalls++
		return "", errors.New("mutation forbidden")
	}
	return runner.output, runner.err
}

func resetPlanGitCallMutates(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "add", "commit", "push", "tag", "update-ref", "branch", "checkout", "switch", "reset", "restore", "clean":
		return true
	default:
		return false
	}
}

func resetPlanGHCallMutates(args []string) bool {
	if len(args) == 0 {
		return false
	}
	if args[0] != "api" {
		return true
	}
	for index, arg := range args {
		if (arg == "--method" || arg == "-X") && index+1 < len(args) && args[index+1] != "GET" {
			return true
		}
	}
	return false
}

func setReleasePlanCommandRunnersForTest(t *testing.T, gitRunner preflight.GitRunner, ghRunner preflight.GHRunner) func() {
	t.Helper()
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.releasePlanGitRunner = gitRunner
		dependencies.releasePlanGHRunner = ghRunner
	})
	return func() {}
}

func newReleasePlanCommandRepo(t *testing.T, baseTag string, commits ...releasePlanCommandCommit) string {
	t.Helper()
	repoDir := newEmptyReleasePlanGitRepo(t)
	writeReleasePlanFile(t, repoDir, "seed.txt", "seed\n")
	gitReleasePlan(t, repoDir, "add", "-A")
	commitReleasePlan(t, repoDir, "chore: seed")
	gitReleasePlan(t, repoDir, "tag", baseTag)
	gitReleasePlan(t, repoDir, "remote", "add", "origin", "https://example.invalid/roundfix.git")

	for index, commit := range commits {
		for _, changedPath := range commit.paths {
			writeReleasePlanFile(t, repoDir, changedPath, fmt.Sprintf("%s\ncommit=%d\n", commit.subject, index))
		}
		gitReleasePlan(t, repoDir, "add", "-A")
		if commit.body == "" {
			commitReleasePlan(t, repoDir, commit.subject)
		} else {
			commitReleasePlan(t, repoDir, commit.subject, commit.body)
		}
	}
	return repoDir
}

func runReleasePlanCommandInRepo(t *testing.T, repoDir string, args ...string) (int, string, string) {
	t.Helper()
	before := snapshotReleasePlanRepo(t, repoDir)
	setCommandWorkDirForTest(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	commandArgs := append([]string{"release", "plan"}, args...)
	code := runCLIContext(t, context.Background(), commandArgs, &stdout, &stderr)
	assertReleasePlanRepoUnchanged(t, repoDir, before)
	return code, stdout.String(), stderr.String()
}

func decodeReleasePlanJSON(t *testing.T, output string) releasePlanJSON {
	t.Helper()
	var plan releasePlanJSON
	decoder := json.NewDecoder(strings.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
		t.Fatalf("decode release plan JSON %q: %v", output, err)
	}
	var trailing struct{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		t.Fatalf("expected exactly one release plan JSON object, got trailing content in %q", output)
	}
	return plan
}

func decodeReleaseResetPlanJSON(t *testing.T, output string) releaseResetPlanJSON {
	t.Helper()
	var plan releaseResetPlanJSON
	decoder := json.NewDecoder(strings.NewReader(output))
	if err := decoder.Decode(&plan); err != nil {
		t.Fatalf("decode release reset plan JSON %q: %v", output, err)
	}
	var trailing struct{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		t.Fatalf("expected exactly one release reset plan JSON object, got trailing content in %q", output)
	}
	return plan
}

func releasePlanJSONHasSubject(plan releasePlanJSON, subject string) bool {
	for _, change := range plan.Changes {
		if change.Subject == subject {
			return true
		}
	}
	return false
}

func assertReleasePlanNoStderr(t *testing.T, stderr string) {
	t.Helper()
	if stderr != "" {
		t.Fatalf("expected stderr empty, got %q", stderr)
	}
}

func assertReleasePlanOneDiagnostic(t *testing.T, stderr string, want string) {
	t.Helper()
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr missing %q: %q", want, stderr)
	}
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	if len(lines) != 1 || lines[0] == "" {
		t.Fatalf("expected one actionable stderr diagnostic, got %q", stderr)
	}
}
