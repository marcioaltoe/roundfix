// Suite: QA mechanical audit across Spec Runs.
// Invariant: the authorization audit reads every governed Task commit and the grant from the Delivery Base.
// Boundary IN: disposable Git repositories, qaMechanicalRequest, and the real mechanical stage.
// Boundary OUT: QA report rendering and Agent execution, owned by their daemon suites.
package daemon

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

const qaEveryRunAuditSlug = "qa-every-run-audit"

type qaEveryRunAuditFixture struct {
	t                 *testing.T
	repoRoot          string
	specDir           string
	authorizationPath string
	initialHead       string
	plan              TaskPlan
}

func newQAEveryRunAuditFixture(t *testing.T, withDefaultBranch bool) *qaEveryRunAuditFixture {
	t.Helper()
	repoRoot := t.TempDir()
	gittest.InitRepo(t, repoRoot, "-b", "main")
	gittest.AppendConfig(t, repoRoot, "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n")

	specDir := filepath.Join(repoRoot, "docs", "specs", qaEveryRunAuditSlug)
	authorizationPath := "docs/workflow/authorizations/qa-every-run.md"
	for _, dir := range []string{
		filepath.Dir(filepath.Join(repoRoot, authorizationPath)),
		specDir,
		filepath.Join(repoRoot, "docs", "agents"),
		filepath.Join(repoRoot, "internal"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("create fixture directory %q: %v", dir, err)
		}
	}
	mustWriteForTest(t, filepath.Join(repoRoot, authorizationPath), qaEveryRunAuthorization("Makefile"))
	mustWriteForTest(t, filepath.Join(repoRoot, "docs", "agents", "agent-instructions.md"), "# Agent instructions\n")
	mustWriteForTest(t, filepath.Join(repoRoot, "docs", "agents", "domain.md"), "# Domain instructions\n")
	mustWriteForTest(t, filepath.Join(specDir, "_prd.md"), "# PRD\n\n## Project Constraints\n\n"+
		"- Identifier strategy: not applicable — the fixture creates no identifier. Source: `docs/agents/domain.md`.\n"+
		"- Authentication and HTTP: not applicable — the fixture opens no transport. Source: `docs/agents/agent-instructions.md`.\n"+
		"- Active ADR obligations: not applicable — the fixture cites no ADR. Source: `docs/agents/domain.md`.\n"+
		"- Tooling authority: applicable — recorded at `"+authorizationPath+"`; bounded files: `Makefile`. Source: `docs/agents/agent-instructions.md`.\n\n"+
		"## Success Metrics\n\nNone. This fixture measures QA mechanical authorization behavior.\n")
	mustWriteForTest(t, filepath.Join(specDir, "task_01.md"), "task one\n")
	mustWriteForTest(t, filepath.Join(specDir, "task_02.md"), "task two\n")
	mustWriteForTest(t, filepath.Join(repoRoot, "Makefile"), "verify:\n\t@true\n")
	mustWriteForTest(t, filepath.Join(repoRoot, "internal", "ordinary.go"), "package internal\n")
	runGitForTest(t, repoRoot, "add", "-A")
	runGitForTest(t, repoRoot, "commit", "-q", "-m", "initial")
	initialHead := strings.TrimSpace(runGitForTest(t, repoRoot, "rev-parse", "HEAD"))
	if withDefaultBranch {
		runGitForTest(t, repoRoot, "update-ref", "refs/remotes/origin/main", initialHead)
		runGitForTest(t, repoRoot, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	}
	runGitForTest(t, repoRoot, "switch", "-q", "-c", "feature/qa-every-run")

	return &qaEveryRunAuditFixture{
		t:                 t,
		repoRoot:          repoRoot,
		specDir:           specDir,
		authorizationPath: authorizationPath,
		initialHead:       initialHead,
		plan: TaskPlan{
			WorkDir:   repoRoot,
			HeadSHA:   initialHead,
			SpecsRoot: filepath.Join(repoRoot, "docs", "specs"),
			Spec:      spec.Spec{Slug: qaEveryRunAuditSlug, Dir: specDir},
			Tasks: []spec.Task{
				{ID: "task_01", File: filepath.Join(qaEveryRunAuditSlug, "task_01.md")},
				{ID: "task_02", File: filepath.Join(qaEveryRunAuditSlug, "task_02.md")},
			},
		},
	}
}

func qaEveryRunAuthorization(paths ...string) string {
	var content strings.Builder
	content.WriteString("---\nstatus: approved\ngranted: 2026-09-28\naction: audit bounded paths\nconsuming: " + qaEveryRunAuditSlug + "\npaths:\n")
	for _, path := range paths {
		content.WriteString("  - " + path + "\n")
	}
	content.WriteString("---\n")
	return content.String()
}

func (fixture *qaEveryRunAuditFixture) commitTask(taskID, subject, path, content string) string {
	fixture.t.Helper()
	mustWriteForTest(fixture.t, filepath.Join(fixture.repoRoot, path), content)
	runGitForTest(fixture.t, fixture.repoRoot, "add", "-A")
	runGitForTest(fixture.t, fixture.repoRoot, "commit", "-q", "-m", subject, "-m", "Roundfix-Spec: "+qaEveryRunAuditSlug+"\nRoundfix-Task: "+taskID)
	return strings.TrimSpace(runGitForTest(fixture.t, fixture.repoRoot, "rev-parse", "HEAD"))
}

func (fixture *qaEveryRunAuditFixture) commitGrant(subject string, paths ...string) string {
	fixture.t.Helper()
	mustWriteForTest(fixture.t, filepath.Join(fixture.repoRoot, fixture.authorizationPath), qaEveryRunAuthorization(paths...))
	runGitForTest(fixture.t, fixture.repoRoot, "add", "-A")
	runGitForTest(fixture.t, fixture.repoRoot, "commit", "-q", "-m", subject)
	return strings.TrimSpace(runGitForTest(fixture.t, fixture.repoRoot, "rev-parse", "HEAD"))
}

func (fixture *qaEveryRunAuditFixture) mechanicalRequest(runStart string) speccheck.MechanicalRequest {
	fixture.t.Helper()
	plan := fixture.plan
	plan.HeadSHA = runStart
	request, err := (&Engine{}).qaMechanicalRequest(context.Background(), plan, spec.Task{}, "")
	if err != nil {
		fixture.t.Fatalf("qaMechanicalRequest returned error: %v", err)
	}
	return request
}

func runEveryRunMechanicalStage(t *testing.T, request speccheck.MechanicalRequest) speccheck.MechanicalResult {
	t.Helper()
	result, err := speccheck.RunMechanicalStage(context.Background(), request)
	if err != nil {
		t.Fatalf("RunMechanicalStage returned error: %v", err)
	}
	return result
}

func resultHasAuthorizationFinding(result speccheck.MechanicalResult, commit string) bool {
	for _, finding := range result.Findings {
		if finding.Code == speccheck.CodeMechanicalAuthPaths && strings.Contains(finding.Detail, commit) {
			return true
		}
	}
	return false
}

func TestQADeliveryBaseIsTheMergeBaseWithTheDefaultBranch(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	fixture.commitTask("task_01", "feat: first work commit", filepath.Join("internal", "ordinary.go"), "package internal\n\nconst first = true\n")
	fixture.commitTask("task_02", "chore: second work commit", "Makefile", "verify:\n\t@true\nsecond:\n\t@true\n")

	base, resolved, err := qaDeliveryBase(context.Background(), fixture.plan)
	if err != nil {
		t.Fatalf("qaDeliveryBase returned error: %v", err)
	}
	if !resolved || base != fixture.initialHead {
		t.Fatalf("qaDeliveryBase = %q, %v, want %q, true", base, resolved, fixture.initialHead)
	}
}

func TestQADeliveryBasePrefersTheRemoteTrackingBranch(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "main")
	mustWriteForTest(t, filepath.Join(fixture.repoRoot, "default.txt"), "remote default\n")
	runGitForTest(t, fixture.repoRoot, "add", "default.txt")
	runGitForTest(t, fixture.repoRoot, "commit", "-q", "-m", "docs: advance remote default")
	remoteHead := strings.TrimSpace(runGitForTest(t, fixture.repoRoot, "rev-parse", "HEAD"))
	runGitForTest(t, fixture.repoRoot, "update-ref", "refs/remotes/origin/main", remoteHead)
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "feature/qa-every-run")
	runGitForTest(t, fixture.repoRoot, "merge", "-q", "--ff-only", "refs/remotes/origin/main")
	runGitForTest(t, fixture.repoRoot, "update-ref", "refs/heads/main", fixture.initialHead)
	fixture.commitTask("task_01", "feat: work after remote default", filepath.Join("internal", "ordinary.go"), "package internal\n\nconst remoteBase = true\n")

	base, resolved, err := qaDeliveryBase(context.Background(), fixture.plan)
	if err != nil {
		t.Fatalf("qaDeliveryBase returned error: %v", err)
	}
	if !resolved || base != remoteHead {
		t.Fatalf("qaDeliveryBase = %q, %v, want remote-tracking base %q, true", base, resolved, remoteHead)
	}
}

func TestQADeliveryBaseIsUnresolvedWithoutADefaultBranch(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, false)

	base, resolved, err := qaDeliveryBase(context.Background(), fixture.plan)
	if err != nil {
		t.Fatalf("qaDeliveryBase returned error: %v", err)
	}
	if resolved || base != "" {
		t.Fatalf("qaDeliveryBase = %q, %v, want empty, false", base, resolved)
	}
}

func TestQAMechanicalRequestAuditsATaskCommitFromAnEarlierRun(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	earlierCommit := fixture.commitTask("task_01", "chore: earlier ungranted tooling", ".golangci.yml", "linters: {}\n")
	fixture.commitTask("task_02", "chore: current authorized tooling", "Makefile", "verify:\n\t@true\ncurrent:\n\t@true\n")

	result := runEveryRunMechanicalStage(t, fixture.mechanicalRequest(earlierCommit))
	if !resultHasAuthorizationFinding(result, earlierCommit) {
		t.Fatalf("mechanical findings = %+v, want QA-AUTH-PATHS for earlier commit %s", result.Findings, earlierCommit)
	}
}

func TestQAMechanicalRequestAuditsEveryCommitOfATask(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	olderCommit := fixture.commitTask("task_01", "chore: older ungranted tooling", ".golangci.yml", "linters: {}\n")
	newerCommit := fixture.commitTask("task_01", "chore: newer authorized tooling", "Makefile", "verify:\n\t@true\nnewer:\n\t@true\n")

	result := runEveryRunMechanicalStage(t, fixture.mechanicalRequest(newerCommit))
	if !resultHasAuthorizationFinding(result, olderCommit) {
		t.Fatalf("mechanical findings = %+v, want QA-AUTH-PATHS for older commit %s", result.Findings, olderCommit)
	}
}

func TestQAMechanicalRequestLeavesDefaultBranchCommitsUnaudited(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "main")
	defaultCommit := fixture.commitTask("task_01", "chore: default branch tooling", ".golangci.yml", "linters: {}\n")
	runGitForTest(t, fixture.repoRoot, "update-ref", "refs/remotes/origin/main", defaultCommit)
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "feature/qa-every-run")
	runGitForTest(t, fixture.repoRoot, "merge", "-q", "--ff-only", "main")
	fixture.commitTask("task_02", "chore: work branch tooling", "Makefile", "verify:\n\t@true\nwork:\n\t@true\n")

	request := fixture.mechanicalRequest(fixture.initialHead)
	for _, commit := range request.TaskCommits {
		if commit.SHA == defaultCommit {
			t.Fatalf("TaskCommits = %+v, default-branch commit %s must stay unaudited", request.TaskCommits, defaultCommit)
		}
	}
	runEveryRunMechanicalStage(t, request)
}

func TestQAMechanicalRequestRefusesAGrantWidenedOnTheSpecBranch(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	fixture.commitGrant("chore: widen grant on work branch", "Makefile", ".golangci.yml")
	taskCommit := fixture.commitTask("task_02", "chore: use work-branch grant", ".golangci.yml", "linters: {}\n")

	result := runEveryRunMechanicalStage(t, fixture.mechanicalRequest(fixture.initialHead))
	if !resultHasAuthorizationFinding(result, taskCommit) {
		t.Fatalf("mechanical findings = %+v, want QA-AUTH-PATHS for work-branch grant consumer %s", result.Findings, taskCommit)
	}
}

func TestQAMechanicalRequestAcceptsAGrantLandedOnTheDefaultBranch(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, true)
	fixture.commitTask("task_01", "feat: work before default merge", filepath.Join("internal", "ordinary.go"), "package internal\n\nconst beforeMerge = true\n")
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "main")
	defaultGrant := fixture.commitGrant("chore: widen grant on default branch", "Makefile", ".golangci.yml")
	runGitForTest(t, fixture.repoRoot, "update-ref", "refs/remotes/origin/main", defaultGrant)
	runGitForTest(t, fixture.repoRoot, "switch", "-q", "feature/qa-every-run")
	runGitForTest(t, fixture.repoRoot, "merge", "-q", "--no-edit", "main")
	taskCommit := fixture.commitTask("task_02", "chore: use default-branch grant", ".golangci.yml", "linters: {}\n")

	result := runEveryRunMechanicalStage(t, fixture.mechanicalRequest(fixture.initialHead))
	if resultHasAuthorizationFinding(result, taskCommit) {
		t.Fatalf("mechanical findings = %+v, default-branch grant must authorize commit %s", result.Findings, taskCommit)
	}
}

func TestQAMechanicalRequestFallsBackToTheRunStartHeadWithASkip(t *testing.T) {
	t.Parallel()
	fixture := newQAEveryRunAuditFixture(t, false)
	runStart := fixture.commitTask("task_01", "chore: earlier ungranted tooling", ".golangci.yml", "linters: {}\n")
	currentCommit := fixture.commitTask("task_02", "chore: current authorized tooling", "Makefile", "verify:\n\t@true\ncurrent:\n\t@true\n")

	request := fixture.mechanicalRequest(runStart)
	if request.DeliveryTargetRevision != runStart {
		t.Fatalf("DeliveryTargetRevision = %q, want Run start head %s", request.DeliveryTargetRevision, runStart)
	}
	if len(request.TaskCommits) != 1 || request.TaskCommits[0].SHA != currentCommit {
		t.Fatalf("TaskCommits = %+v, want only current-Run commit %s", request.TaskCommits, currentCommit)
	}
	result := runEveryRunMechanicalStage(t, request)
	wantSkip := speccheck.MechanicalSkip{
		Detector:        speccheck.DetectorMechanicalAuthPaths,
		MissingArtifact: speccheck.MechanicalSkipEarlierRunTaskCommits,
	}
	for _, skip := range result.Skips {
		if skip == wantSkip {
			return
		}
	}
	t.Fatalf("mechanical skips = %+v, want %+v", result.Skips, wantSkip)
}
