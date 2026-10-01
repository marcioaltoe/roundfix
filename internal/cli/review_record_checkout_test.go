// Suite: checkout-local pre-PR review records.
// Invariant: each checkout reads and writes only its own review record and answer.
// Boundary IN: public review commands, real Git worktrees, and Artifact Directory files.
// Boundary OUT: real reviewer execution and publication workflows.
package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
)

func TestReviewRecordsLiveInTheCheckoutsOwnDirectory(t *testing.T) {
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "No findings.", StopReason: "end_turn"}},
		{result: agent.ExecuteResult{Message: "No findings.", StopReason: "end_turn"}},
	}}
	fixture := newReviewCheckoutFixture(t, runner)

	for _, checkout := range []string{fixture.primary, fixture.secondary} {
		digest := sha256.Sum256([]byte(filepath.Clean(checkout)))
		wantDir := filepath.Join(fixture.artifactDir, "pre-pr-review", hex.EncodeToString(digest[:])[:16])
		if got := reviewCheckoutDir(fixture.artifactDir, checkout); got != wantDir {
			t.Fatalf("review checkout directory = %q, want %q", got, wantDir)
		}
		code, record, stderr := runReviewAtCheckout(t, fixture, checkout)
		if code != exitOK || stderr != "" || record.Outcome != reviewOutcomeReviewed {
			t.Fatalf("review in %q exit=%d record=%+v stderr=%q", checkout, code, record, stderr)
		}
		checkoutDir := reviewCheckoutDir(fixture.artifactDir, checkout)
		if record.AnswerPath != filepath.Join(checkoutDir, reviewAnswerFileName) {
			t.Fatalf("review answer path = %q, want checkout-local path", record.AnswerPath)
		}
		for _, name := range []string{reviewRecordFileName, reviewAnswerFileName} {
			if _, err := os.Stat(filepath.Join(checkoutDir, name)); err != nil {
				t.Fatalf("stat checkout-local %s: %v", name, err)
			}
		}
		info, err := os.Stat(checkoutDir)
		if err != nil {
			t.Fatalf("stat checkout review directory: %v", err)
		}
		if info.Mode().Perm() != 0o755 {
			t.Fatalf("checkout review directory mode = %o, want 755", info.Mode().Perm())
		}
	}

	if reviewCheckoutDir(fixture.artifactDir, fixture.primary) == reviewCheckoutDir(fixture.artifactDir, fixture.secondary) {
		t.Fatal("linked checkouts resolved to the same review directory")
	}
	for _, name := range []string{reviewRecordFileName, reviewAnswerFileName} {
		if _, err := os.Stat(filepath.Join(fixture.artifactDir, name)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("old shared %s stat error = %v, want not exist", name, err)
		}
	}
}

func TestReviewInOneCheckoutLeavesAnotherCheckoutsRecord(t *testing.T) {
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "Findings:\n- review.txt:2: checkout-local finding", StopReason: "end_turn"}},
		{result: agent.ExecuteResult{Message: "No findings.", StopReason: "end_turn"}},
	}}
	fixture := newReviewCheckoutFixture(t, runner)

	code, first, stderr := runReviewAtCheckout(t, fixture, fixture.primary)
	if code != exitRunFailed || stderr != "" || first.Outcome != reviewOutcomeFindings {
		t.Fatalf("first checkout review exit=%d record=%+v stderr=%q", code, first, stderr)
	}
	firstDir := reviewCheckoutDir(fixture.artifactDir, fixture.primary)
	firstRecordBefore := mustReadReviewCheckoutFile(t, firstDir, reviewRecordFileName)
	firstAnswerBefore := mustReadReviewCheckoutFile(t, firstDir, reviewAnswerFileName)

	code, second, stderr := runReviewAtCheckout(t, fixture, fixture.secondary)
	if code != exitOK || stderr != "" || second.Outcome != reviewOutcomeReviewed {
		t.Fatalf("second checkout review exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	if got := mustReadReviewCheckoutFile(t, firstDir, reviewRecordFileName); !bytes.Equal(got, firstRecordBefore) {
		t.Fatalf("second checkout changed first checkout record: before=%q after=%q", firstRecordBefore, got)
	}
	if got := mustReadReviewCheckoutFile(t, firstDir, reviewAnswerFileName); !bytes.Equal(got, firstAnswerBefore) {
		t.Fatalf("second checkout changed first checkout answer: before=%q after=%q", firstAnswerBefore, got)
	}

	code, reused, _ := runReviewAtCheckout(t, fixture, fixture.primary)
	if code != exitRunFailed || !reused.Reused || reused.Outcome != reviewOutcomeFindings {
		t.Fatalf("first checkout reuse exit=%d record=%+v", code, reused)
	}
	if runner.preparedCalls != 2 {
		t.Fatalf("reviewer prompt calls = %d, want 2 after checkout-local reuse", runner.preparedCalls)
	}

	code, stdout, stderr := runReviewDisposeAtCheckout(t, fixture, fixture.primary, "F1", "--dismiss", "--evidence", "covered by the supported path")
	if code != exitOK || stdout == "" || stderr != "" {
		t.Fatalf("first checkout dispose exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestReviewDisposeReadsOnlyItsCheckoutsRecord(t *testing.T) {
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "Findings:\n- review.txt:2: checkout-local finding", StopReason: "end_turn"},
	}}}
	fixture := newReviewCheckoutFixture(t, runner)

	code, record, stderr := runReviewAtCheckout(t, fixture, fixture.primary)
	if code != exitRunFailed || stderr != "" || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("first checkout review exit=%d record=%+v stderr=%q", code, record, stderr)
	}

	code, stdout, stderr := runReviewDisposeAtCheckout(t, fixture, fixture.secondary, "F1", "--dismiss", "--evidence", "evidence")
	if code != exitPreflight || stdout != "" || !strings.Contains(stderr, "pre-pr review record does not exist") {
		t.Fatalf("second checkout dispose exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("disposition ledger stat error = %v, want not exist", err)
	}
}

func TestReviewIgnoresARecordAtTheSharedLocation(t *testing.T) {
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "No findings.", StopReason: "end_turn"},
	}}}
	fixture := newReviewCheckoutFixture(t, runner)
	if err := os.MkdirAll(fixture.artifactDir, 0o755); err != nil {
		t.Fatalf("create Artifact Directory: %v", err)
	}

	sharedRecordPath := filepath.Join(fixture.artifactDir, reviewRecordFileName)
	sharedAnswerPath := filepath.Join(fixture.artifactDir, reviewAnswerFileName)
	shared := newReviewRecord(
		fixture.primary,
		fixture.baseCommit,
		fixture.headCommit,
		roundconfig.PrePRReview{Provider: "codex", Source: "project"},
		reviewOutcomeFindings,
	)
	shared.Findings = "- internal/cli/review.go:1: legacy shared finding"
	shared.FindingItems = splitReviewFindings(shared.Findings)
	shared.AnswerPath = sharedAnswerPath
	writeDispositionRecordFile(t, sharedRecordPath, shared)
	mustWrite(t, sharedAnswerPath, "legacy shared answer")
	sharedRecordBefore := mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewRecordFileName)
	sharedAnswerBefore := mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewAnswerFileName)

	code, fresh, stderr := runReviewAtCheckout(t, fixture, fixture.primary)
	if code != exitOK || stderr != "" || fresh.Outcome != reviewOutcomeReviewed || fresh.Reused {
		t.Fatalf("fresh review exit=%d record=%+v stderr=%q", code, fresh, stderr)
	}
	if runner.preparedCalls != 1 {
		t.Fatalf("reviewer prompt calls = %d, want 1 when shared record is ignored", runner.preparedCalls)
	}
	if got := mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewRecordFileName); !bytes.Equal(got, sharedRecordBefore) {
		t.Fatalf("review changed shared record: before=%q after=%q", sharedRecordBefore, got)
	}
	if got := mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewAnswerFileName); !bytes.Equal(got, sharedAnswerBefore) {
		t.Fatalf("review changed shared answer: before=%q after=%q", sharedAnswerBefore, got)
	}
}

type reviewCheckoutFixture struct {
	primary     string
	secondary   string
	artifactDir string
	baseCommit  string
	headCommit  string
}

func newReviewCheckoutFixture(t *testing.T, runner *reviewCommandRunner) reviewCheckoutFixture {
	t.Helper()
	homeDir, primary, baseCommit, headCommit := newReviewCommandRepository(t)
	secondary := filepath.Join(t.TempDir(), "linked-checkout")
	gittest.Run(t, primary, "worktree", "add", "-b", "feature/review-linked", secondary, headCommit)
	resolvedSecondary, err := filepath.EvalSymlinks(secondary)
	if err != nil {
		t.Fatalf("resolve linked checkout: %v", err)
	}
	secondary = resolvedSecondary

	artifactDir := filepath.Join(t.TempDir(), "artifacts")
	writeReviewCommandConfig(t, primary, "codex", artifactDir)
	writeReviewCommandConfig(t, secondary, "codex", artifactDir)
	setCommandEnvironmentForTest(t, homeDir, primary)
	withAgentRunner(t, runner)
	return reviewCheckoutFixture{
		primary:     primary,
		secondary:   secondary,
		artifactDir: artifactDir,
		baseCommit:  baseCommit,
		headCommit:  headCommit,
	}
}

func runReviewAtCheckout(t *testing.T, fixture reviewCheckoutFixture, checkout string) (int, reviewRecord, string) {
	t.Helper()
	setCommandWorkDirForTest(t, checkout)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode review record from stdout %q: %v", stdout.String(), err)
	}
	return code, record, stderr.String()
}

func runReviewDisposeAtCheckout(t *testing.T, fixture reviewCheckoutFixture, checkout string, args ...string) (int, string, string) {
	t.Helper()
	setCommandWorkDirForTest(t, checkout)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	commandArgs := append([]string{"review", "dispose"}, args...)
	code := runCLI(t, commandArgs, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func mustReadReviewCheckoutFile(t *testing.T, directory string, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(directory, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return content
}
