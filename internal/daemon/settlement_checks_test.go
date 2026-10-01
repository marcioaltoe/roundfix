package daemon

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

type fakeSettlementChecker struct {
	findings      []speccheck.Finding
	audit         func(speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error)
	calls         int
	refusingCalls int
}

func (c *fakeSettlementChecker) RefusingFindings(string, string, string) ([]speccheck.Finding, error) {
	c.refusingCalls++
	return c.findings, nil
}
func (c *fakeSettlementChecker) AuditCommit(_ context.Context, _ speccheck.MechanicalRequest, commit speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
	c.calls++
	if c.audit != nil {
		return c.audit(commit)
	}
	return nil, nil
}

func settlementAuditFixture(t *testing.T, checker *fakeSettlementChecker) (*taskCycleFixture, TaskPlan, *taskFakeRunner, *taskFakeVerifier, *engineFakeCommitter, *Engine) {
	t.Helper()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	commitTaskFixtureSource(t, fixture.gitRoot, "seed gated Spec")
	plan.HeadSHA = strings.TrimSpace(gittest.Run(t, fixture.gitRoot, "rev-parse", "HEAD"))
	plan.Authorization = spec.ReadSpecAuthorization(context.Background(), plan.WorkDir, plan.SpecsRoot, plan.Spec.Slug, plan.HeadSHA)
	fixture.worktree.snapshots = [][]string{nil, {"Makefile"}}
	runner.afterTask = func(string) {
		mustWriteForTest(t, filepath.Join(fixture.gitRoot, "Makefile"), "changed\n")
	}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
	engine.deps.SettlementChecker = checker
	return fixture, plan, runner, verifier, committer, engine
}

func auditFailure() []speccheck.MechanicalFinding {
	return []speccheck.MechanicalFinding{{Code: speccheck.CodeMechanicalAuthPaths, Detail: "outside grant: Makefile", File: "grant.md", Line: 1, Fix: "repair tree"}}
}

func TestSettlementAuditFindingReturnsFeedbackAndRepairSettles(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, plan, runner, _, committer, engine := settlementAuditFixture(t, checker)
	runner.afterTask = func(string) {
		content := "escaped\n"
		if len(runner.requests) > 1 {
			content = "repaired\n"
		}
		mustWriteForTest(t, filepath.Join(fixture.gitRoot, "Makefile"), content)
	}
	checker.audit = func(commit speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
		if commit.Parent != plan.HeadSHA || !strings.Contains(strings.Join(commit.Changed, ","), "Makefile") {
			t.Fatalf("prospective commit = %+v", commit)
		}
		content, err := os.ReadFile(filepath.Join(plan.WorkDir, "Makefile"))
		if err != nil {
			return nil, err
		}
		if string(content) == "escaped\n" {
			return auditFailure(), nil
		}
		return nil, nil
	}
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 || checker.calls != 2 || len(committer.messages) != 1 {
		t.Fatalf("result %+v, err %v, calls %d, commits %v", result, err, checker.calls, committer.messages)
	}
	if len(runner.requests) != 2 || runner.requests[0].Session != runner.requests[1].Session {
		t.Fatalf("requests = %+v", runner.requests)
	}
	if prompt := runner.requests[1].Prompt; !strings.Contains(prompt, "Verification Feedback") || !strings.Contains(prompt, "settlement check: authorization") || !strings.Contains(prompt, "Makefile") {
		t.Fatalf("Feedback = %s", prompt)
	}
}

func TestSettlementCheckFailureOnFinalAttemptFailsWithNamedReason(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{audit: func(speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
		return auditFailure(), nil
	}}
	fixture, plan, runner, _, committer, engine := settlementAuditFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || result.Completed != 0 || len(committer.messages) != 0 || len(runner.requests) != 2 {
		t.Fatalf("result %+v, err %v, commits %v", result, err, committer.messages)
	}
	outcome, _ := taskOutcomeByID(result.Outcomes, "task_01")
	if !strings.Contains(outcome.Reason, "Settlement check failed: settlement check: authorization:") || !strings.Contains(outcome.Reason, "QA-AUTH-PATHS") {
		t.Fatalf("reason = %s", outcome.Reason)
	}
	path := filepath.Join(plan.ArtifactDir, "runs", plan.RunID, "verification", "batch-001-attempt-2-settlement-authorization.log")
	content, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(content), "Makefile") {
		t.Fatalf("diagnostics %q, %v", content, err)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "failed" {
		t.Fatalf("status = %s", got)
	}
	if head := strings.TrimSpace(gittest.Run(t, plan.WorkDir, "rev-parse", "HEAD")); head != plan.HeadSHA {
		t.Fatalf("Task commit created: %s", head)
	}
}

func TestSettlementCheckerErrorFailsTheCheck(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{audit: func(speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
		return nil, errors.New("cannot read grant")
	}}
	_, plan, runner, _, committer, engine := settlementAuditFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || len(committer.messages) != 0 {
		t.Fatalf("result %+v, err %v", result, err)
	}
	if len(runner.requests) != 2 || !strings.Contains(runner.requests[1].Prompt, "cannot read grant") {
		t.Fatalf("Feedback requests = %+v", runner.requests)
	}
}

func TestSettlementChecksRunAfterAFailedCommand(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{audit: func(speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
		return auditFailure(), nil
	}}
	_, plan, runner, verifier, _, engine := settlementAuditFixture(t, checker)
	verifier.failOn = map[string]error{"declared check": errors.New("exit status 1")}
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || checker.calls != 2 {
		t.Fatalf("result %+v, err %v, calls %d", result, err, checker.calls)
	}
	prompt := runner.requests[1].Prompt
	if !strings.Contains(prompt, "Failed command: declared check") || !strings.Contains(prompt, "Failed command: settlement check: authorization") {
		t.Fatalf("Feedback = %s", prompt)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check", "declared check"}) {
		t.Fatalf("commands = %v", verifier.commands)
	}
}

func TestSettlementChecksPublishVerificationEvents(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, plan, _, _, _, engine := settlementAuditFixture(t, checker)
	checker.audit = func(speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
		if checker.calls == 1 {
			return auditFailure(), nil
		}
		return nil, nil
	}
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("result %+v, %v", result, err)
	}
	var phases []string
	for _, event := range fixture.sink.snapshot() {
		if event.Kind == runevent.KindDaemonVerification && eventPayloadString(t, event, "command") == "settlement check: authorization" {
			payload := eventPayloadMap(t, event)
			if payload["task"] != "task_01" {
				t.Fatalf("event payload = %+v", payload)
			}
			phases = append(phases, eventPayloadString(t, event, "phase"))
		}
	}
	if !reflect.DeepEqual(phases, []string{"started", "failed", "started", "command-passed"}) {
		t.Fatalf("phases = %v", phases)
	}
}

func TestGatelessGraphNeverCallsTheSettlementChecker(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, _, runner, verifier, committer, engine := settlementAuditFixture(t, checker)
	plan := fixture.plan()
	plan.Tasks = plan.Tasks[:1]
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 || checker.calls != 0 || checker.refusingCalls != 0 || len(runner.requests) != 1 || len(committer.messages) != 1 {
		t.Fatalf("result %+v, err %v, calls %d", result, err, checker.calls)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check"}) {
		t.Fatalf("commands = %v", verifier.commands)
	}
}

func TestDefaultSettlementCheckerAuditsAProspectiveCommit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	grant := "docs/specs/check/_authorization.md"
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, grant)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(root, grant), "---\nstatus: approved\ngranted: 2026-09-09\naction: test\nconsuming: check\npaths:\n  - Makefile\n---\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "grant")
	parent := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	findings, err := (SpecCheckSettlementChecker{}).AuditCommit(context.Background(), speccheck.MechanicalRequest{RepoRoot: root, AuthorizationPath: grant, ConsumingSpec: "check", DeliveryTargetRevision: parent}, speccheck.ProspectiveTaskCommit{TaskID: "task_02", Parent: parent, Changed: []string{".agents/skills/escaped/SKILL.md"}})
	if err != nil || len(findings) != 1 || !strings.Contains(findings[0].Detail, ".agents/skills/escaped/SKILL.md") {
		t.Fatalf("findings %+v, err %v", findings, err)
	}
}

func TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"check"}}})
	plan := fixture.qaPlan()
	prd := filepath.Join(plan.Spec.Dir, "_prd.md")
	content, err := os.ReadFile(prd)
	if err != nil {
		t.Fatal(err)
	}
	head, _, found := strings.Cut(string(content), "## Project Constraints")
	if !found {
		t.Fatal("fixture lacks Project Constraints")
	}
	mustWriteForTest(t, prd, head)
	gate, err := qaGatePrecondition(plan)
	if err != nil || !gate.Blocking {
		t.Fatalf("precondition %+v, %v", gate, err)
	}
	findings, err := (SpecCheckSettlementChecker{}).RefusingFindings(plan.SpecsRoot, plan.WorkDir, plan.Spec.Slug)
	if err != nil || !reflect.DeepEqual(findings, gate.Findings) {
		t.Fatalf("findings %+v, want %+v, err %v", findings, gate.Findings, err)
	}
}

type settlementVerifierFunc func(context.Context, VerifyRequest) (VerifyResult, error)

func (f settlementVerifierFunc) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	return f(ctx, req)
}

func TestSettlementChecksSkipAnUnobservedCommand(t *testing.T) {
	testSettlementChecksSkip(t, &VerificationUnknownError{Command: "check", Err: errors.New("unobserved")})
}

func TestSettlementChecksSkipAnInfrastructureError(t *testing.T) {
	testSettlementChecksSkip(t, errors.New("cannot start verifier"))
}

func TestSettlementChecksSkipAStopRequest(t *testing.T) {
	testSettlementChecksSkip(t, ErrStopRequested)
}

func testSettlementChecksSkip(t *testing.T, commandErr error) {
	t.Helper()
	checker := &fakeSettlementChecker{}
	_, _, _, _, _, engine := settlementAuditFixture(t, checker)
	engine.deps.Verifier = settlementVerifierFunc(func(context.Context, VerifyRequest) (VerifyResult, error) { return VerifyResult{}, commandErr })
	calls := 0
	_, err := engine.runVerificationAttempt(context.Background(), verificationAttemptRequest{
		RunID: "test", ArtifactDir: t.TempDir(), BatchNumber: 1, Attempt: 1,
		Commands: []string{"check"}, Checks: []verificationCheck{{Label: "settlement check: authorization", Run: func(context.Context, string) (string, error) { calls++; return "", nil }}},
		Publish: func(context.Context, string, map[string]any) error { return nil },
	})
	var unknown *VerificationUnknownError
	if !errors.As(commandErr, &unknown) && !errors.Is(err, commandErr) {
		t.Fatalf("err = %v, want %v", err, commandErr)
	}
	if calls != 0 {
		t.Fatalf("check calls = %d", calls)
	}
}

func TestSettlementChecksSkipGrantReadsForOrdinaryPaths(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, plan, _, _, _, engine := settlementAuditFixture(t, checker)
	fixture.worktree.snapshots = [][]string{{"src/ordinary.go"}}
	if err := os.MkdirAll(filepath.Join(plan.WorkDir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(plan.WorkDir, "src/ordinary.go"), "ordinary\n")
	mustWriteForTest(t, filepath.Join(plan.Spec.Dir, "_prd.md"), "malformed authorization reference\n")
	failure, err := engine.authorizationSettlementCheck(plan, plan.Tasks[0], nil).Run(context.Background(), "")
	if err != nil || failure != "" || checker.calls != 0 {
		t.Fatalf("failure %s, err %v, calls %d", failure, err, checker.calls)
	}
}

func TestSettlementChecksExcludePreexistingGovernedPaths(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, plan, _, _, _, engine := settlementAuditFixture(t, checker)
	fixture.worktree.snapshots = [][]string{{"Makefile"}}
	mustWriteForTest(t, filepath.Join(plan.WorkDir, "Makefile"), "preexisting\n")
	failure, err := engine.authorizationSettlementCheck(plan, plan.Tasks[0], []string{"Makefile"}).Run(context.Background(), "")
	if err != nil || failure != "" || checker.calls != 0 {
		t.Fatalf("failure %s, err %v, calls %d", failure, err, checker.calls)
	}
}

func TestSettlementChecksRunAgainOnTemporaryRetry(t *testing.T) {
	t.Parallel()
	checker := &fakeSettlementChecker{}
	fixture, plan, _, _, _, engine := settlementAuditFixture(t, checker)
	fixture.worktree.snapshots = [][]string{{"Makefile"}}
	plan.settlementChecks = true
	plan.verificationGate = newVerificationGate(1)
	attempts := 0
	engine.deps.Verifier = settlementVerifierFunc(func(_ context.Context, req VerifyRequest) (VerifyResult, error) {
		attempts++
		if attempts == 1 {
			return VerifyResult{}, &TemporaryVerificationFailureError{CommandFailure: &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 75")}}
		}
		return VerifyResult{OutputPath: req.OutputPath}, nil
	})
	mustWriteForTest(t, filepath.Join(plan.WorkDir, "Makefile"), "changed\n")
	retryUsed := false
	outcome, err := engine.verifyTask(context.Background(), plan, plan.Tasks[0], 1, 1, &retryUsed, nil)
	if err != nil || outcome.Failure != "" || !retryUsed || checker.calls != 2 {
		t.Fatalf("outcome %+v, err %v, retry %v, checks %d", outcome, err, retryUsed, checker.calls)
	}
}
