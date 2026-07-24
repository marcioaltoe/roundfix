// Suite: Human Baseline workflow
// Invariant: one linear terminal workflow produces the automation-equivalent Plan and mutates only after exact final approval.
// Boundary IN: adoption/update state detection, numbered prompts, consolidated review, no-TTY refusal, and exact apply wiring.
// Boundary OUT: catalog semantics, repository inventory internals, and transaction phase failure injection.

package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestHumanBaselineAdoption(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(humanBaselineAdoptionAnswers("1")),
			interactive: true,
		},
	)
	if code != exitOK {
		t.Fatalf("human adoption exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Baseline apply: verified") ||
		!strings.Contains(stdout.String(), "Consolidated Change Plan review") {
		t.Fatalf("human adoption output does not show review and verified apply:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Final confirmation for Plan Digest sha256:") {
		t.Fatalf("human adoption did not bind confirmation to digest:\n%s", stderr.String())
	}
	manifest := filepath.Join(repo, filepath.FromSlash(baselineSetupManifestPath))
	if _, err := os.Stat(manifest); err != nil {
		t.Fatalf("human adoption did not write verified Setup Manifest: %v", err)
	}
	after := baselinePlanTestTree(t, repo)
	if before == after {
		t.Fatal("approved human adoption did not change repository bytes")
	}
}

func TestHumanBaselineUpdate(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	applyHumanBaselineFixturePlan(t, repo)
	before := baselinePlanTestTree(t, repo)

	answers := []string{"1", "1"}
	for range 10 {
		answers = append(answers, "1")
	}
	answers = append(answers, "2")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(strings.Join(answers, "\n") + "\n"),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("human update decline exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Baseline workflow: update",
		"Current Baseline Profile: go-cli-tui",
		"Baseline Plan was declined; no repository bytes were written",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("human update output missing %q:\n%s", want, stdout.String())
		}
	}
	if !strings.Contains(stderr.String(), "Change Baseline Profile") {
		t.Fatalf("human update did not offer explicit profile change:\n%s", stderr.String())
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("declined human update changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}

	changeAnswers := []string{"1", "2", "2"}
	for range 10 {
		changeAnswers = append(changeAnswers, "1")
	}
	var changeReview bytes.Buffer
	var changePrompts bytes.Buffer
	humanPlan, err := driveHumanBaselinePlan(
		context.Background(),
		repo,
		&baselineHumanPrompt{
			reader: bufioReader(strings.Join(changeAnswers, "\n") + "\n"),
			writer: &changePrompts,
		},
		&changeReview,
	)
	if err != nil {
		t.Fatalf("build human profile-change Plan: %v\nreview=%s\nprompts=%s", err, changeReview.String(), changePrompts.String())
	}
	automation, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "rust-cli",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || automation.Plan == nil {
		t.Fatalf("build automation profile-change Plan: outcome=%+v error=%v", automation, err)
	}
	humanBytes, err := baseline.MarshalPlanDocument(humanPlan)
	if err != nil {
		t.Fatal(err)
	}
	automationBytes, err := baseline.MarshalPlanDocument(*automation.Plan)
	if err != nil {
		t.Fatal(err)
	}
	if humanPlan.Profile.ID != "rust-cli" || !bytes.Equal(humanBytes, automationBytes) {
		t.Fatalf(
			"profile-change parity failed: profile=%s human=%s automation=%s",
			humanPlan.Profile.ID,
			humanPlan.PlanDigest,
			automation.Plan.PlanDigest,
		)
	}
}

func TestConsolidatedReview(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "Preserve this repository-specific rule.\n")
	commitBaselinePlanTestRepository(t, repo)
	before := baselinePlanTestTree(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(humanBaselinePreservationAnswers()),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("preservation review exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Consolidated editable classification review:",
		"normative-clause -> repository-rules",
		"Consolidated Change Plan review",
		"Complete managed-entry ledger:",
		"Complete Upgrade Retention Contract ledger:",
		"Plan Digest: sha256:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("consolidated review missing %q:\n%s", want, stdout.String())
		}
	}
	fileChanges := strings.Index(stdout.String(), "File changes:")
	managedLedger := strings.Index(stdout.String(), "Complete managed-entry ledger:")
	retentionLedger := strings.Index(stdout.String(), "Complete Upgrade Retention Contract ledger:")
	if fileChanges < 0 || managedLedger < fileChanges || retentionLedger < managedLedger {
		t.Fatalf("consolidated review order is not fileChanges then ledgers:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Edit classifications and dispositions") {
		t.Fatalf("classification proposal was not editable:\n%s", stderr.String())
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("declined preservation review changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

func TestConsolidatedReviewEditsManagedClassification(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	writeBaselinePlanTestFile(
		t,
		repo,
		"AGENTS.md",
		"<!-- setup-context-driven:begin id=legacy.rule version=1 -->\nPreserve me.\n<!-- setup-context-driven:end id=legacy.rule -->\n",
	)
	commitBaselinePlanTestRepository(t, repo)

	answers := strings.TrimSuffix(humanBaselineAdoptionAnswers(""), "\n")
	answers = "2" + strings.TrimPrefix(answers, "1")
	answers += "\n2\n1\nmanaged setup evidence is not repository policy\n2\n"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(answers),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("edited classification exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Consolidated Change Plan review") ||
		strings.Contains(stdout.String(), "Category: decision") {
		t.Fatalf("edited managed classification did not produce a valid review:\n%s", stdout.String())
	}
}

func TestHumanAutomationPlanParity(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	var review bytes.Buffer
	var prompts bytes.Buffer
	humanPlan, err := driveHumanBaselinePlan(
		context.Background(),
		repo,
		&baselineHumanPrompt{
			reader: bufioReader(humanBaselineAdoptionAnswers("")),
			writer: &prompts,
		},
		&review,
	)
	if err != nil {
		t.Fatalf("build human Baseline Plan: %v\nreview=%s\nprompts=%s", err, review.String(), prompts.String())
	}
	automation, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || automation.Plan == nil {
		t.Fatalf("build automation Baseline Plan: outcome=%+v error=%v", automation, err)
	}
	humanBytes, err := baseline.MarshalPlanDocument(humanPlan)
	if err != nil {
		t.Fatal(err)
	}
	automationBytes, err := baseline.MarshalPlanDocument(*automation.Plan)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(humanBytes, automationBytes) {
		t.Fatalf(
			"human and automation plans differ\ndigests: human=%s automation=%s",
			humanPlan.PlanDigest,
			automation.Plan.PlanDigest,
		)
	}
}

func TestBaselineNoTTY(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo, "--format=json"},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader("1\n"),
			interactive: false,
		},
	)
	if code != exitUnverified {
		t.Fatalf("no-TTY baseline exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	result, err := baseline.ParseResult(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse no-TTY result: %v\n%s", err, stdout.String())
	}
	if result.Category != "interactive_input" ||
		!strings.Contains(result.NextAction, "roundfix baseline plan") ||
		!strings.Contains(result.NextAction, "roundfix baseline apply") {
		t.Fatalf("no-TTY result is not actionable: %+v", result)
	}
	if strings.Contains(stderr.String(), "Prompt ") {
		t.Fatalf("no-TTY baseline emitted a hidden prompt:\n%s", stderr.String())
	}
	if after := baselinePlanTestTree(t, repo); before != after {
		t.Fatalf("no-TTY baseline changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "roundfix")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = projectRoot
	build.Env = append(os.Environ(), "GOCACHE="+filepath.Join(t.TempDir(), "go-cache"))
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real roundfix CLI: %v\n%s", err, output)
	}
	command := exec.Command(binary, "baseline", "--repo", repo, "--format=json")
	command.Stdin = strings.NewReader("1\n1\n1\n")
	var realOut bytes.Buffer
	var realErr bytes.Buffer
	command.Stdout = &realOut
	command.Stderr = &realErr
	err = command.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != exitUnverified {
		t.Fatalf("redirected real baseline error=%v stdout=%s stderr=%s", err, realOut.String(), realErr.String())
	}
	if strings.Contains(realErr.String(), "Prompt ") {
		t.Fatalf("redirected real baseline emitted a prompt:\n%s", realErr.String())
	}
}

func newHumanBaselineRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func applyHumanBaselineFixturePlan(t *testing.T, repo string) {
	t.Helper()
	outcome, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("build update fixture plan: outcome=%+v error=%v", outcome, err)
	}
	if _, err := baseline.ApplyPlan(context.Background(), repo, *outcome.Plan, outcome.Plan.PlanDigest); err != nil {
		t.Fatalf("apply update fixture plan: %v", err)
	}
}

func humanBaselineFixtureDecisions() []baseline.DecisionValue {
	return []baseline.DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: true},
		{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		{ID: "runtime.design", Value: "claude opus xhigh"},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: false},
	}
}

func humanBaselineAdoptionAnswers(final string) string {
	answers := []string{
		"1",
		"1",
		"1",
		"make verify",
		"1",
		"1",
		"2",
		"1",
		"2",
		"2",
		"codex gpt-5.5 xhigh",
		"claude opus xhigh",
	}
	if final != "" {
		answers = append(answers, final)
	}
	return strings.Join(answers, "\n") + "\n"
}

func humanBaselinePreservationAnswers() string {
	answers := strings.TrimSuffix(humanBaselineAdoptionAnswers(""), "\n")
	answers = "2" + strings.TrimPrefix(answers, "1")
	return answers + "\n1\n2\n"
}

func bufioReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func decodeHumanBaselineResult(t *testing.T, data []byte) baseline.Result {
	t.Helper()
	var result baseline.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}
