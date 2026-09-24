// Suite: pre-PR review records and reviewer input.
// Invariant: each review is bound to one candidate, preserves its answer, and reports one substantive verdict.
// Boundary IN: candidate diff collection, prompt execution, verdict classification, and artifact persistence.
// Boundary OUT: policy resolution, top-level command routing, and real ACP adapters.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
)

func TestReviewRecordRoundTripsEachOutcome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		outcome  reviewOutcome
		findings string
		reason   string
	}{
		{
			name:    "reviewed",
			outcome: reviewOutcomeReviewed,
		},
		{
			name:     "findings",
			outcome:  reviewOutcomeFindings,
			findings: "internal/cli/review.go: head validation is missing",
		},
		{
			name:    "blocked",
			outcome: reviewOutcomeBlocked,
			reason:  "reviewer output was empty",
		},
		{
			name:    "omitted",
			outcome: reviewOutcomeOmitted,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			policy := roundconfig.PrePRReview{
				Provider: "codex",
				Source:   "project",
			}
			want := newReviewRecord(
				"github.com/example/project",
				"1111111111111111111111111111111111111111",
				"2222222222222222222222222222222222222222",
				policy,
				test.outcome,
			)
			want.Findings = test.findings
			want.Reason = test.reason

			var output bytes.Buffer
			if err := writeReviewRecord(&output, want); err != nil {
				t.Fatalf("write review record: %v", err)
			}

			var got reviewRecord
			if err := json.Unmarshal(output.Bytes(), &got); err != nil {
				t.Fatalf("decode review record: %v", err)
			}
			if got != want {
				t.Fatalf("review record = %+v, want %+v", got, want)
			}
		})
	}
}

func TestReviewRecordWithoutHeadIsRefusedBeforeWriting(t *testing.T) {
	t.Parallel()

	record := newReviewRecord(
		"github.com/example/project",
		"1111111111111111111111111111111111111111",
		"",
		roundconfig.PrePRReview{Provider: "codex", Source: "project"},
		reviewOutcomeReviewed,
	)
	var output bytes.Buffer

	err := writeReviewRecord(&output, record)

	if err == nil || !strings.Contains(err.Error(), "head commit is required") {
		t.Fatalf("write review record error = %v, want missing-head error", err)
	}
	if output.Len() != 0 {
		t.Fatalf("invalid review record wrote %q", output.String())
	}
}

func TestReviewRecordRefusesInvalidOutcomeDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		outcome  reviewOutcome
		findings string
		reason   string
		want     string
	}{
		{
			name:    "unknown outcome",
			outcome: "unknown",
			want:    "outcome \"unknown\" is invalid",
		},
		{
			name:    "findings outcome without findings",
			outcome: reviewOutcomeFindings,
			want:    "findings are required",
		},
		{
			name:    "blocked outcome without reason",
			outcome: reviewOutcomeBlocked,
			want:    "reason is required",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			record := newReviewRecord(
				"github.com/example/project",
				"1111111111111111111111111111111111111111",
				"2222222222222222222222222222222222222222",
				roundconfig.PrePRReview{Provider: "codex", Source: "project"},
				test.outcome,
			)
			record.Findings = test.findings
			record.Reason = test.reason

			var output bytes.Buffer
			err := writeReviewRecord(&output, record)

			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("write review record error = %v, want %q", err, test.want)
			}
			if output.Len() != 0 {
				t.Fatalf("invalid review record wrote %q", output.String())
			}
		})
	}
}

func TestReviewPromptCarriesTheCandidateDiff(t *testing.T) {
	t.Parallel()

	repository, baseCommit, headCommit := reviewCandidateFixture(t)
	runner := &recordingReviewRunner{}

	if _, err := runReviewSession(
		t.Context(),
		agent.ExecuteRequest{GitRoot: repository},
		baseCommit,
		headCommit,
		preflight.ExecGitRunner{},
		runner,
		runevent.Discard,
	); err != nil {
		t.Fatalf("run review session: %v", err)
	}

	for _, want := range []string{
		baseCommit,
		headCommit,
		"diff --git a/review.txt b/review.txt",
		"+the reviewer receives this changed line",
	} {
		if !strings.Contains(runner.request.Prompt, want) {
			t.Fatalf("review prompt does not contain %q:\n%s", want, runner.request.Prompt)
		}
	}
}

func TestReviewSessionReadsWithoutWriting(t *testing.T) {
	t.Parallel()

	repository, baseCommit, headCommit := reviewCandidateFixture(t)
	runner := &recordingReviewRunner{}

	if _, err := runReviewSession(
		t.Context(),
		agent.ExecuteRequest{GitRoot: repository},
		baseCommit,
		headCommit,
		preflight.ExecGitRunner{},
		runner,
		runevent.Discard,
	); err != nil {
		t.Fatalf("run review session: %v", err)
	}

	if !runner.request.Access.CanRead() {
		t.Fatal("review session cannot read files")
	}
	if runner.request.Access.CanWrite() {
		t.Fatal("review session can write")
	}
}

func TestReviewSessionRefIsUniquePerInvocation(t *testing.T) {
	t.Parallel()

	const headCommit = "2222222222222222222222222222222222222222"
	first := reviewSessionRef(headCommit, "/tmp/repository", 0)
	second := reviewSessionRef(headCommit, "/tmp/repository", 0)

	if first == second {
		t.Fatalf("review session references = %+v and %+v, want different names for separate invocations", first, second)
	}
}

func TestReviewCommandExitsZeroOnExplicitClean(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Output:     `{"jsonrpc":"2.0","method":"session/update"}`,
				Message:    "No findings",
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitOK {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitOK, stderr)
	}
	if record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review outcome = %q, want %q", record.Outcome, reviewOutcomeReviewed)
	}
	fixture.assertCandidate(t, record)
}

func TestReviewCommandDefaultsBaseToMain(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode review record: %v", err)
	}
	if record.BaseCommit != fixture.baseCommit {
		t.Fatalf("default base = %q, want main commit %q", record.BaseCommit, fixture.baseCommit)
	}
}

func TestReviewCommandExitsOneAndRecordsFindings(t *testing.T) {
	const findings = "internal/cli/review.go:42: blocked result can be lost"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "Findings:\n" + findings, StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr)
	}
	if record.Outcome != reviewOutcomeFindings || record.Findings != findings {
		t.Fatalf("review record = %+v, want findings %q", record, findings)
	}
	fixture.assertCandidate(t, record)
}

func TestReviewClassifiesVerdictVariants(t *testing.T) {
	tests := []struct {
		name         string
		answer       string
		wantCode     int
		wantOutcome  reviewOutcome
		wantFindings string
	}{
		{
			name:        "punctuated clean verdict",
			answer:      "No findings.",
			wantCode:    exitOK,
			wantOutcome: reviewOutcomeReviewed,
		},
		{
			name:        "emphasized clean verdict",
			answer:      "**No findings.**",
			wantCode:    exitOK,
			wantOutcome: reviewOutcomeReviewed,
		},
		{
			name:        "lowercase clean verdict",
			answer:      "no findings",
			wantCode:    exitOK,
			wantOutcome: reviewOutcomeReviewed,
		},
		{
			name:         "preamble before emphasized lowercase findings verdict",
			answer:       "I reviewed the candidate.\n**findings:**\ninternal/cli/review.go:42: first finding\ninternal/cli/review_test.go:42: second finding",
			wantCode:     exitRunFailed,
			wantOutcome:  reviewOutcomeFindings,
			wantFindings: "internal/cli/review.go:42: first finding\ninternal/cli/review_test.go:42: second finding",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{
				results: []reviewCommandRunResult{{
					result: agent.ExecuteResult{Message: test.answer, StopReason: "end_turn"},
				}},
			}
			fixture := newReviewCommandFixture(t, "codex", runner)

			code, record, stderr := fixture.run(t)

			if code != test.wantCode {
				t.Fatalf("review exit = %d, want %d; stderr=%q", code, test.wantCode, stderr)
			}
			if record.Outcome != test.wantOutcome || record.Findings != test.wantFindings {
				t.Fatalf("review record = %+v, want outcome %q and findings %q", record, test.wantOutcome, test.wantFindings)
			}
		})
	}
}

func TestReviewBlocksAmbiguousVerdict(t *testing.T) {
	tests := []struct {
		name       string
		answer     string
		wantReason string
	}{
		{
			name:       "both verdicts",
			answer:     "No findings.\nFindings:\ninternal/cli/review.go:42: contradictory verdict",
			wantReason: "both",
		},
		{
			name:       "neither verdict",
			answer:     "Review complete.",
			wantReason: "neither",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{
				results: []reviewCommandRunResult{{
					result: agent.ExecuteResult{Message: test.answer, StopReason: "end_turn"},
				}},
			}
			fixture := newReviewCommandFixture(t, "codex", runner)

			code, record, stderr := fixture.run(t)

			assertBlockedReviewCommand(t, code, record, stderr, test.wantReason)
		})
	}
}

func TestReviewKeepsTheRawAnswer(t *testing.T) {
	tests := []struct {
		name        string
		answer      string
		wantOutcome reviewOutcome
	}{
		{name: "reviewed", answer: "  **No findings.**\n", wantOutcome: reviewOutcomeReviewed},
		{name: "findings", answer: "Findings:\ninternal/cli/review.go:42: keep this finding\n", wantOutcome: reviewOutcomeFindings},
		{name: "blocked", answer: "Review complete, but no verdict was stated.\n", wantOutcome: reviewOutcomeBlocked},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{
				results: []reviewCommandRunResult{{
					result: agent.ExecuteResult{Message: test.answer, StopReason: "end_turn"},
				}},
			}
			fixture := newReviewCommandFixture(t, "codex", runner)

			_, record, _ := fixture.run(t)

			if record.Outcome != test.wantOutcome {
				t.Fatalf("review outcome = %q, want %q", record.Outcome, test.wantOutcome)
			}
			wantPath := filepath.Join(fixture.artifactDir, reviewAnswerFileName)
			if record.AnswerPath != wantPath {
				t.Fatalf("review answer path = %q, want %q", record.AnswerPath, wantPath)
			}
			answer, err := os.ReadFile(record.AnswerPath)
			if err != nil {
				t.Fatalf("read raw review answer: %v", err)
			}
			if string(answer) != test.answer {
				t.Fatalf("raw review answer = %q, want %q", string(answer), test.answer)
			}
		})
	}
}

func TestReviewCommandBlocksOnRuntimeFailure(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{err: errors.New("runtime unavailable")}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "runtime failure")
	if runner.preparedCalls != 1 {
		t.Fatalf("review prompt calls = %d, want 1 with no fallback", runner.preparedCalls)
	}
}

func TestReviewCommandDoesNotFallbackAfterPromptFailure(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			err: &agent.SelectionFailureError{Runtime: "codex", Reason: "adapter ended after prompt"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "runtime failure")
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 {
		t.Fatalf("post-prompt failure activated fallback: prepares=%d prompts=%d", runner.prepareCalls, runner.preparedCalls)
	}
}

func TestReviewCommandKeepsATimeoutBlocked(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings"},
			err:    agent.StopError{Err: context.DeadlineExceeded},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "timeout")
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 {
		t.Fatalf("timeout activated fallback: prepares=%d prompts=%d", runner.prepareCalls, runner.preparedCalls)
	}
}

func TestReviewCommandBlocksOnTransportAnomaly(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Message:          "No findings",
				StopReason:       "end_turn",
				TransportAnomaly: "acpx exited with exit code 1 after parsed session/prompt result",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "transport anomaly")
}

func TestReviewCommandBlocksOnEmptyAgentOutput(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{result: agent.ExecuteResult{Output: "raw protocol only", StopReason: "end_turn"}}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "empty agent output")
}

func TestReviewCommandClassifiesAgentMessageNotProtocolStream(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Output:     `{"jsonrpc":"2.0","result":{"text":"No findings"}}`,
				Message:    "review complete",
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "unclassifiable agent output")
}

func TestReviewCommandNoneRecordsOmissionWithoutAgentActivity(t *testing.T) {
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "none", runner)

	code, record, stderr := fixture.run(t)

	if code != exitOK {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitOK, stderr)
	}
	if record.Outcome != reviewOutcomeOmitted {
		t.Fatalf("review outcome = %q, want %q", record.Outcome, reviewOutcomeOmitted)
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("none policy used Agent runtime: probes=%d prepares=%d prompts=%d", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
	}
}

func TestReviewCommandRefusesUnimplementedProvider(t *testing.T) {
	for _, provider := range []string{"claude", "coderabbit"} {
		t.Run(provider, func(t *testing.T) {
			runner := &reviewCommandRunner{}
			fixture := newReviewCommandFixture(t, provider, runner)

			code, record, stderr := fixture.run(t)

			assertBlockedReviewCommand(t, code, record, stderr, provider)
			if record.Outcome == reviewOutcomeOmitted {
				t.Fatalf("provider %q was recorded as omitted", provider)
			}
			if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
				t.Fatalf("provider %q used Agent runtime: probes=%d prepares=%d prompts=%d", provider, runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
			}
		})
	}
}

func TestReviewCommandRefusesProviderProfileMismatch(t *testing.T) {
	tests := []struct {
		name              string
		preferredRuntime  string
		fallbackRuntime   string
		mismatchedRuntime string
	}{
		{
			name:              "preferred runtime",
			preferredRuntime:  "claude",
			fallbackRuntime:   "codex",
			mismatchedRuntime: "claude",
		},
		{
			name:              "fallback runtime",
			preferredRuntime:  "codex",
			fallbackRuntime:   "opencode",
			mismatchedRuntime: "opencode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{}
			fixture := newReviewCommandFixture(t, "codex", runner)
			writeReviewCommandProfileConfig(
				t,
				fixture.repository,
				fixture.provider,
				fixture.artifactDir,
				test.preferredRuntime,
				test.fallbackRuntime,
			)

			code, record, stderr := fixture.run(t)

			assertBlockedReviewCommand(t, code, record, stderr, "configuration error")
			for _, want := range []string{"codex", test.mismatchedRuntime} {
				if !strings.Contains(record.Reason, want) {
					t.Fatalf("review refusal = %q, want selected provider and mismatched runtime named", record.Reason)
				}
			}
			if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
				t.Fatalf("provider/profile mismatch used Agent runtime: probes=%d prepares=%d prompts=%d", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
			}
		})
	}
}

func TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt(t *testing.T) {
	runner := &reviewCommandRunner{
		prepareErrors: []error{
			&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"},
			nil,
		},
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("fallback review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if runner.prepareCalls != 2 || runner.preparedCalls != 1 {
		t.Fatalf("fallback calls: prepares=%d prompts=%d, want 2 and 1", runner.prepareCalls, runner.preparedCalls)
	}
	if !strings.Contains(stderr, "activating fallback 1") {
		t.Fatalf("fallback notification missing from stderr: %q", stderr)
	}
}

func TestReviewCommandProvesArtifactDirectoryWritableBeforeAgentActivity(t *testing.T) {
	homeDir, repository, baseCommit, headCommit := newReviewCommandRepository(t)
	artifactPath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(artifactPath, []byte("file"), 0o644); err != nil {
		t.Fatalf("write Artifact Directory obstruction: %v", err)
	}
	writeReviewCommandConfig(t, repository, "codex", artifactPath)
	setCommandEnvironmentForTest(t, homeDir, repository)
	runner := &reviewCommandRunner{}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review", "--base", baseCommit}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitPreflight, stderr.String())
	}
	if !strings.Contains(stderr.String(), "not a directory") {
		t.Fatalf("Artifact Directory refusal missing: %q", stderr.String())
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("Agent activity preceded Artifact Directory refusal: probes=%d prepares=%d prompts=%d", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
	}
	if headCommit == "" {
		t.Fatal("fixture head commit is empty")
	}
}

func TestReviewCommandRefusesUnknownFlags(t *testing.T) {
	homeDir, repository, _, _ := newReviewCommandRepository(t)
	setCommandEnvironmentForTest(t, homeDir, repository)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review", "--unknown"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("review exit = %d, want %d", code, exitPreflight)
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("unknown flag diagnostic missing: %q", stderr.String())
	}
}

func reviewCandidateFixture(t *testing.T) (repository string, baseCommit string, headCommit string) {
	t.Helper()

	repository = t.TempDir()
	gittest.InitRepo(t, repository, "--initial-branch=main")
	gittest.PersistIdentity(t, repository)
	path := filepath.Join(repository, "review.txt")
	if err := os.WriteFile(path, []byte("before review\n"), 0o644); err != nil {
		t.Fatalf("write base candidate: %v", err)
	}
	gittest.Run(t, repository, "add", "review.txt")
	gittest.Run(t, repository, "commit", "-m", "base")
	baseCommit = strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))

	if err := os.WriteFile(path, []byte("before review\nthe reviewer receives this changed line\n"), 0o644); err != nil {
		t.Fatalf("write head candidate: %v", err)
	}
	gittest.Run(t, repository, "add", "review.txt")
	gittest.Run(t, repository, "commit", "-m", "candidate")
	headCommit = strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
	return repository, baseCommit, headCommit
}

type recordingReviewRunner struct {
	request agent.ExecuteRequest
}

func (*recordingReviewRunner) Probe(context.Context, agent.ProbeRequest) error {
	return nil
}

func (runner *recordingReviewRunner) Run(
	_ context.Context,
	request agent.ExecuteRequest,
	_ runevent.Sink,
) (agent.ExecuteResult, error) {
	runner.request = request
	return agent.ExecuteResult{Message: "No findings"}, nil
}

func (*recordingReviewRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type reviewCommandRunResult struct {
	result agent.ExecuteResult
	err    error
}

type reviewCommandRunner struct {
	probeCalls    int
	prepareCalls  int
	preparedCalls int
	endCalls      int
	prepareErrors []error
	results       []reviewCommandRunResult
}

func (runner *reviewCommandRunner) Probe(context.Context, agent.ProbeRequest) error {
	runner.probeCalls++
	return nil
}

func (runner *reviewCommandRunner) PrepareSession(context.Context, agent.ExecuteRequest, runevent.Sink) error {
	index := runner.prepareCalls
	runner.prepareCalls++
	if index < len(runner.prepareErrors) {
		return runner.prepareErrors[index]
	}
	return nil
}

func (runner *reviewCommandRunner) RunPrepared(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	index := runner.preparedCalls
	runner.preparedCalls++
	if index < len(runner.results) {
		return runner.results[index].result, runner.results[index].err
	}
	return agent.ExecuteResult{}, nil
}

func (runner *reviewCommandRunner) Run(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	return agent.ExecuteResult{}, errors.New("review command bypassed the prepared-session boundary")
}

func (runner *reviewCommandRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	runner.endCalls++
	return nil
}

type reviewCommandFixture struct {
	homeDir     string
	repository  string
	artifactDir string
	baseCommit  string
	headCommit  string
	provider    string
}

func newReviewCommandFixture(t *testing.T, provider string, runner agent.Runner) reviewCommandFixture {
	t.Helper()
	homeDir, repository, baseCommit, headCommit := newReviewCommandRepository(t)
	artifactDir := filepath.Join(t.TempDir(), "artifacts")
	writeReviewCommandConfig(t, repository, provider, artifactDir)
	setCommandEnvironmentForTest(t, homeDir, repository)
	withAgentRunner(t, runner)
	return reviewCommandFixture{
		homeDir:     homeDir,
		repository:  repository,
		artifactDir: artifactDir,
		baseCommit:  baseCommit,
		headCommit:  headCommit,
		provider:    provider,
	}
}

func newReviewCommandRepository(t *testing.T) (homeDir, repository, baseCommit, headCommit string) {
	t.Helper()
	homeDir = t.TempDir()
	repository, baseCommit, headCommit = reviewCandidateFixture(t)
	gittest.Run(t, repository, "checkout", "-b", "feature/review")
	gittest.Run(t, repository, "branch", "-f", "main", baseCommit)
	resolved, err := filepath.EvalSymlinks(repository)
	if err != nil {
		t.Fatalf("resolve review repository: %v", err)
	}
	repository = resolved
	return homeDir, repository, baseCommit, headCommit
}

func writeReviewCommandConfig(t *testing.T, repository string, provider string, artifactDir string) {
	t.Helper()
	mustWrite(t, filepath.Join(repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: %s
`, artifactDir, provider))
}

func writeReviewCommandProfileConfig(
	t *testing.T,
	repository string,
	provider string,
	artifactDir string,
	preferredRuntime string,
	fallbackRuntime string,
) {
	t.Helper()
	mustWrite(t, filepath.Join(repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: %s
profiles:
  review:
    preferred:
      runtime: %s
      model: preferred-model
      reasoning_effort: high
    fallbacks:
      - runtime: %s
        model: fallback-model
        reasoning_effort: high
`, artifactDir, provider, preferredRuntime, fallbackRuntime))
}

func (fixture reviewCommandFixture) run(t *testing.T) (int, reviewRecord, string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)

	var stdoutRecord reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &stdoutRecord); err != nil {
		t.Fatalf("decode stdout review record %q: %v", stdout.String(), err)
	}
	fileBytes, err := os.ReadFile(filepath.Join(fixture.artifactDir, "pre-pr-review.json"))
	if err != nil {
		t.Fatalf("read persisted review record: %v", err)
	}
	var fileRecord reviewRecord
	if err := json.Unmarshal(fileBytes, &fileRecord); err != nil {
		t.Fatalf("decode persisted review record: %v", err)
	}
	if fileRecord != stdoutRecord {
		t.Fatalf("persisted record = %+v, stdout record = %+v", fileRecord, stdoutRecord)
	}
	return code, fileRecord, stderr.String()
}

func (fixture reviewCommandFixture) assertCandidate(t *testing.T, record reviewRecord) {
	t.Helper()
	if record.Repository != fixture.repository || record.BaseCommit != fixture.baseCommit || record.HeadCommit != fixture.headCommit {
		t.Fatalf("review candidate = %+v, want repository=%q base=%q head=%q", record, fixture.repository, fixture.baseCommit, fixture.headCommit)
	}
	if record.Provider != fixture.provider || record.Source != "project" {
		t.Fatalf("review policy = provider %q source %q, want %q/project", record.Provider, record.Source, fixture.provider)
	}
}

func assertBlockedReviewCommand(t *testing.T, code int, record reviewRecord, stderr string, reason string) {
	t.Helper()
	if code != exitPreflight {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitPreflight, stderr)
	}
	if record.Outcome != reviewOutcomeBlocked || !strings.Contains(record.Reason, reason) {
		t.Fatalf("review record = %+v, want blocked reason containing %q", record, reason)
	}
	if !strings.Contains(stderr, reason) {
		t.Fatalf("stderr = %q, want reason containing %q", stderr, reason)
	}
}
