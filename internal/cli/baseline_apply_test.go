// Suite: Public Baseline apply CLI
// Invariant: baseline apply keeps requested results on stdout, diagnostics on stderr, and maps exact apply outcomes to stable exits.
// Boundary IN: real flag parsing, strict plan-file handoff, command dispatch, text/JSON rendering, and real binary execution.
// Boundary OUT: transaction phase injection, which is owned by internal/baseline/transaction_test.go.

package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineApplyCommand(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)
	plan, planPath := baselineApplyTestPlan(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "apply",
		"--repo", repo,
		"--plan", planPath,
		"--confirm-plan", plan.PlanDigest,
		"--format=json",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("baseline apply exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	result, err := baseline.ParseResult(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse baseline apply result: %v\n%s", err, stdout.String())
	}
	if result.Operation != "apply" || result.State != "verified" ||
		result.PlanDigest != plan.PlanDigest ||
		len(result.VerifiedPostimages) != len(plan.Postimages) {
		t.Fatalf("baseline apply result = %+v", result)
	}
}

func TestBaselineApplyStdoutStderrAndExitCodes(t *testing.T) {
	t.Parallel()
	t.Run("confirmation refusal is actionable JSON", func(t *testing.T) {
		repo := newBaselineApplyTestRepository(t)
		plan, planPath := baselineApplyTestPlan(t, repo)
		before := baselinePlanTestTree(t, repo)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", "sha256:" + strings.Repeat("0", 64),
			"--format", "json",
		}, &stdout, &stderr)
		if code != exitUnverified || !strings.Contains(stderr.String(), "confirmed Plan Digest") {
			t.Fatalf("confirmation exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		result, err := baseline.ParseResult(stdout.Bytes())
		if err != nil {
			t.Fatalf("parse confirmation result: %v\n%s", err, stdout.String())
		}
		if result.State != "action_required" || result.Category != "approval" ||
			result.NextAction == "" || result.PlanDigest != plan.PlanDigest {
			t.Fatalf("confirmation result = %+v", result)
		}
		if after := baselinePlanTestTree(t, repo); after != before {
			t.Fatalf("incorrect confirmation changed repository: before=%s after=%s", before, after)
		}
	})

	t.Run("stale preimage is exit three", func(t *testing.T) {
		repo := newBaselineApplyTestRepository(t)
		plan, planPath := baselineApplyTestPlan(t, repo)
		writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@echo stale\n")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", plan.PlanDigest,
			"--format=json",
		}, &stdout, &stderr)
		if code != exitUnverified || !strings.Contains(stderr.String(), "stale") {
			t.Fatalf("stale exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		result, err := baseline.ParseResult(stdout.Bytes())
		if err != nil || result.Category != "stale" || result.State != "action_required" {
			t.Fatalf("stale result = %+v error=%v", result, err)
		}
	})

	t.Run("invalid plan is exit two", func(t *testing.T) {
		planPath := filepath.Join(t.TempDir(), "invalid.json")
		if err := os.WriteFile(planPath, []byte(`{"schemaVersion":"wrong"}`), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "apply",
			"--plan", planPath,
			"--confirm-plan", "sha256:" + strings.Repeat("0", 64),
			"--format=json",
		}, &stdout, &stderr)
		if code != exitPreflight || stderr.Len() == 0 {
			t.Fatalf("invalid plan exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		result, err := baseline.ParseResult(stdout.Bytes())
		if err != nil || result.Category != "invalid" || result.State != "failed" {
			t.Fatalf("invalid result = %+v error=%v", result, err)
		}
	})

	t.Run("incomplete rollback is exit one", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		err := &baseline.ApplyError{
			Kind:       baseline.ApplyErrorExecution,
			NextAction: "recover the transaction before retrying",
			Err: &baseline.IncompleteRollbackError{
				Cause:    errors.New("postimage verification failed"),
				Rollback: errors.New("restore blocked"),
			},
		}
		code := printBaselineApplyFailure(
			baseline.PlanDocument{},
			err,
			true,
			&stdout,
			&stderr,
		)
		if code != exitRunFailed || !strings.Contains(stderr.String(), "rollback is incomplete") {
			t.Fatalf("rollback exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		result, parseErr := baseline.ParseResult(stdout.Bytes())
		if parseErr != nil || result.Category != "execution" || result.State != "failed" {
			t.Fatalf("rollback result = %+v error=%v", result, parseErr)
		}
	})

	t.Run("cancellation is exit 130", func(t *testing.T) {
		repo := newBaselineApplyTestRepository(t)
		plan, planPath := baselineApplyTestPlan(t, repo)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(ctx, []string{
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", plan.PlanDigest,
			"--format=json",
		}, &stdout, &stderr)
		if code != exitSIGINT || !strings.Contains(stderr.String(), context.Canceled.Error()) {
			t.Fatalf("canceled exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		result, parseErr := baseline.ParseResult(stdout.Bytes())
		if parseErr != nil || result.Category != "execution" {
			t.Fatalf("canceled result = %+v error=%v", result, parseErr)
		}
	})

	t.Run("text success reports recommendations without running them", func(t *testing.T) {
		repo := newBaselineApplyTestRepository(t)
		plan, planPath := baselineApplyTestPlan(t, repo)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", plan.PlanDigest,
			"--format=text",
		}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 ||
			!strings.Contains(stdout.String(), "Baseline apply: verified") ||
			!strings.Contains(stdout.String(), "Recommendation (not run): make verify") {
			t.Fatalf("text apply exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
	})
}

func TestBaselineApplyCommandRealCLI(t *testing.T) {
	t.Parallel()
	repo := newBaselineApplyTestRepository(t)
	plan, planPath := baselineApplyTestPlan(t, repo)
	binary := buildBaselineReleaseBinary(t)
	command := exec.Command(
		binary,
		"baseline", "apply",
		"--repo", repo,
		"--plan", planPath,
		"--confirm-plan", plan.PlanDigest,
		"--format=json",
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		t.Fatalf("real baseline apply: %v exit=%v stdout=%s stderr=%s",
			err, errors.As(err, &exitErr), stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("real baseline apply stderr = %q", stderr.String())
	}
	result, err := baseline.ParseResult(stdout.Bytes())
	if err != nil || result.State != "verified" {
		t.Fatalf("real baseline apply result = %+v error=%v", result, err)
	}
}

func newBaselineApplyTestRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func baselineApplyTestPlan(t *testing.T, repo string) (baseline.PlanDocument, string) {
	t.Helper()
	args := []string{
		"baseline", "plan",
		"--repo", repo,
		"--profile", "go-cli-tui",
		"--decision", "preservation.mode=greenfield",
		"--decision", "language.generated=English",
		"--decision", "verification.gate=make verify",
		"--decision", "branch.prefix=ma/",
		"--decision", "spec.scaffold=true",
		"--decision", "domain.layout=single-context",
		"--decision", "triage.external=false",
		"--decision", "autonomous.enabled=true",
		"--decision", "runtime.backend=codex gpt-5.5 xhigh",
		"--decision", "runtime.design=claude opus xhigh",
		"--decision", "secondbrain.enabled=false",
		"--decision", "repository.extension.enabled=false",
		"--format=json",
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunContext(context.Background(), args, &stdout, &stderr); code != exitOK {
		t.Fatalf("build CLI plan exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	plan, err := baseline.ParsePlanDocument(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse CLI plan: %v\n%s", err, stdout.String())
	}
	planPath := filepath.Join(t.TempDir(), "baseline-plan.json")
	if err := os.WriteFile(planPath, stdout.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return plan, planPath
}
