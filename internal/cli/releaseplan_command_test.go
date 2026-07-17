package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"roundfix/internal/releaseplan"
)

// Suite: Release Plan Command
// Invariant: public CLI planning is deterministic, read-only, and stdout/stderr isolated.
// Boundary IN: roundfix release plan flags and temporary Git repositories.
// Boundary OUT: release mutation, Run creation, Roundfix configuration, network services.

func TestReleasePlanCommandMatchesPRDOutcomes(t *testing.T) {
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
	var rootStdout, rootStderr bytes.Buffer
	rootCode := Run([]string{"--help"}, &rootStdout, &rootStderr)
	if rootCode != exitOK {
		t.Fatalf("root help exit = %d, want 0 stderr=%q", rootCode, rootStderr.String())
	}
	if !strings.Contains(rootStdout.String(), "roundfix release plan") || !strings.Contains(rootStdout.String(), "release    Plan the next release version") {
		t.Fatalf("root help missing release plan command: %q", rootStdout.String())
	}

	var stdout, stderr bytes.Buffer
	code := Run([]string{"release", "plan", "--help"}, &stdout, &stderr)
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

type releasePlanCommandCommit struct {
	subject string
	body    string
	paths   []string
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
	t.Chdir(repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	commandArgs := append([]string{"release", "plan"}, args...)
	code := RunContext(context.Background(), commandArgs, &stdout, &stderr)
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
