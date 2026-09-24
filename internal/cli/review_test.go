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
	"reflect"
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
			if !reflect.DeepEqual(got, want) {
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

func TestReviewPromptCarriesSpecDecisions(t *testing.T) {
	const slug = "0160-spec-aware-review"
	const specsRoot = "planning/specs"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), `# Spec-aware review

## Decisions

- Read the Spec from the candidate.

## Acceptance evidence

This content is outside the Decisions section.
`)
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), `# Spec-aware review

## Design

The reviewer receives the complete TechSpec.

## Rejected alternatives

An extra command-line argument was rejected.
`)
	gittest.Run(t, fixture.repository, "add", specsRoot)
	gittest.Run(t, fixture.repository, "commit", "-m", "add spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, fixture.provider, fixture.artifactDir, specsRoot)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("spec-aware review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if len(record.Specs) != 1 || record.Specs[0] != slug {
		t.Fatalf("review Specs = %v, want [%s]", record.Specs, slug)
	}
	for _, want := range []string{
		"--- BEGIN SPEC CONTEXT: " + slug + " ---",
		"## Decisions\n\n- Read the Spec from the candidate.",
		"The reviewer receives the complete TechSpec.",
		"An extra command-line argument was rejected.",
		"contradicts a recorded decision or adopts an alternative the Spec rejected",
	} {
		if !strings.Contains(runner.request.Prompt, want) {
			t.Fatalf("review prompt does not contain %q:\n%s", want, runner.request.Prompt)
		}
	}
	contextStart := strings.LastIndex(runner.request.Prompt, "--- BEGIN SPEC CONTEXT:")
	if contextStart < 0 {
		t.Fatalf("review prompt has no Spec context block:\n%s", runner.request.Prompt)
	}
	if strings.Contains(runner.request.Prompt[contextStart:], "This content is outside the Decisions section.") {
		t.Fatalf("Spec context carries PRD content outside Decisions:\n%s", runner.request.Prompt[contextStart:])
	}
}

func TestReviewSkipsASpecWithoutDecisions(t *testing.T) {
	const slug = "0159-no-decisions"
	const missingTechSpecSlug = "0158-missing-techspec"
	const specsRoot = "planning/specs"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Spec without decisions\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Technical design\n")
	missingTechSpecDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), missingTechSpecSlug)
	if err := os.MkdirAll(missingTechSpecDir, 0o755); err != nil {
		t.Fatalf("create incomplete Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(missingTechSpecDir, "_prd.md"), "# Incomplete Spec\n\n## Decisions\n\n- This PRD has no TechSpec.\n")
	gittest.Run(t, fixture.repository, "add", specsRoot)
	gittest.Run(t, fixture.repository, "commit", "-m", "add Spec without decisions")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, fixture.provider, fixture.artifactDir, specsRoot)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review with skipped Spec exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if len(record.Specs) != 0 {
		t.Fatalf("review Specs = %v, want empty", record.Specs)
	}
	if !reflect.DeepEqual(record.SkippedSpecs, []string{missingTechSpecSlug, slug}) {
		t.Fatalf("review skipped Specs = %v, want [%s %s]", record.SkippedSpecs, missingTechSpecSlug, slug)
	}
	if runner.preparedCalls != 1 {
		t.Fatalf("review prompt calls = %d, want 1", runner.preparedCalls)
	}
	if strings.Contains(runner.request.Prompt, "BEGIN SPEC CONTEXT") {
		t.Fatalf("review prompt carries skipped Spec context:\n%s", runner.request.Prompt)
	}
}

func TestReviewReadsAnArchivedSpecFromTheCandidate(t *testing.T) {
	const slug = "0158-archived-with-candidate"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	archiveRoot := filepath.FromSlash("docs/history/specs")
	specDir := filepath.Join(fixture.repository, archiveRoot, slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create archived Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Archived Spec\n\n## Decisions\n\n- Keep archived context visible.\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Archived technical design\n\nRead this file from HEAD.\n")
	gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join(archiveRoot, slug)))
	gittest.Run(t, fixture.repository, "commit", "-m", "archive Spec in candidate")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Worktree-only technical design\n")

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("archived Spec review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if !reflect.DeepEqual(record.Specs, []string{slug}) {
		t.Fatalf("review Specs = %v, want [%s]", record.Specs, slug)
	}
	for _, want := range []string{slug, "Keep archived context visible.", "Read this file from HEAD."} {
		if !strings.Contains(runner.request.Prompt, want) {
			t.Fatalf("archived Spec prompt does not contain %q:\n%s", want, runner.request.Prompt)
		}
	}
	if strings.Contains(runner.request.Prompt, "Worktree-only technical design") {
		t.Fatalf("archived Spec prompt read the worktree instead of HEAD:\n%s", runner.request.Prompt)
	}
}

func TestReviewBoundsSpecContext(t *testing.T) {
	const slug = "0160-oversized-context"
	const specsRoot = "planning/specs"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create Spec directory: %v", err)
	}
	oversized := strings.Repeat("candidate context ", reviewSpecContextPerSpecLimit)
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Oversized Spec\n\n## Decisions\n\n"+oversized+"\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Oversized technical design\n\n"+oversized+"\n")
	gittest.Run(t, fixture.repository, "add", specsRoot)
	gittest.Run(t, fixture.repository, "commit", "-m", "add oversized Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, fixture.provider, fixture.artifactDir, specsRoot)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("bounded Spec review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if !record.SpecContextTruncated {
		t.Fatalf("review record = %+v, want Spec context truncation recorded", record)
	}
	contextStart := strings.Index(runner.request.Prompt, "--- BEGIN SPEC CONTEXT:")
	if contextStart < 0 {
		t.Fatalf("review prompt has no Spec context block:\n%s", runner.request.Prompt)
	}
	carried := runner.request.Prompt[contextStart:]
	if !strings.Contains(carried, reviewSpecContextTruncationMarker) {
		t.Fatalf("bounded Spec context has no truncation marker:\n%s", carried)
	}
	if len(carried) > reviewSpecContextPerSpecLimit+1024 {
		t.Fatalf("bounded Spec context length = %d, want at most %d plus envelope", len(carried), reviewSpecContextPerSpecLimit)
	}

	contexts := make([]reviewSpecContext, 4)
	for index := range contexts {
		contexts[index] = reviewSpecContext{
			slug: fmt.Sprintf("oversized-%d", index),
			body: strings.Repeat("x", reviewSpecContextPerSpecLimit),
		}
	}
	bounded, dropped, truncated := boundReviewSpecContexts(contexts)
	if !truncated {
		t.Fatal("total-bound fixture was not recorded as truncated")
	}
	if len(dropped) == 0 {
		t.Fatal("total-bound fixture did not report a dropped Spec")
	}
	for _, context := range bounded {
		if !strings.Contains(context.body, reviewSpecContextTruncationMarker) {
			t.Fatalf("total-bound context %q has no visible marker", context.slug)
		}
	}
	totalContext := appendReviewSpecContexts("", reviewSpecContextResult{contexts: bounded, truncated: truncated})
	if len(totalContext) > reviewSpecContextTotalLimit {
		t.Fatalf("total Spec context length = %d, want at most %d", len(totalContext), reviewSpecContextTotalLimit)
	}
}

func TestReviewRecordsDroppedSpecs(t *testing.T) {
	const specsRoot = "planning/specs"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	for index := 1; index <= 4; index++ {
		slug := fmt.Sprintf("0160-bounded-context-%d", index)
		specDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), slug)
		if err := os.MkdirAll(specDir, 0o755); err != nil {
			t.Fatalf("create Spec directory: %v", err)
		}
		oversized := strings.Repeat(fmt.Sprintf("%d", index), reviewSpecContextPerSpecLimit)
		mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Spec\n\n## Decisions\n\n"+oversized+"\n")
		mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# TechSpec\n\n"+oversized+"\n")
	}
	gittest.Run(t, fixture.repository, "add", specsRoot)
	gittest.Run(t, fixture.repository, "commit", "-m", "add bounded Specs")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, fixture.provider, fixture.artifactDir, specsRoot)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("bounded Specs review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if !reflect.DeepEqual(record.SkippedSpecs, []string{"0160-bounded-context-3", "0160-bounded-context-4"}) {
		t.Fatalf("review skipped Specs = %v, want dropped Specs named", record.SkippedSpecs)
	}
	for _, slug := range record.SkippedSpecs {
		if strings.Contains(runner.request.Prompt, "BEGIN SPEC CONTEXT: "+slug) {
			t.Fatalf("review prompt carries dropped Spec %q", slug)
		}
	}
}

func TestReviewReportsSpecReadingFailureDistinctly(t *testing.T) {
	const slug = "0160-unreadable-context"
	const specsRoot = "planning/specs"
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, filepath.FromSlash(specsRoot), slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Spec\n\n## Decisions\n\n- Read candidate files.\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# Technical design\n")
	gittest.Run(t, fixture.repository, "add", specsRoot)
	gittest.Run(t, fixture.repository, "commit", "-m", "add unreadable Spec fixture")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, fixture.provider, fixture.artifactDir, specsRoot)
	withReviewSpecGitRunner(t, failingSpecFileGitRunner{delegate: preflight.ExecGitRunner{}})

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "Spec context read failure")
	if strings.Contains(record.Reason, "runtime failure") {
		t.Fatalf("Spec read reason = %q, want a distinct non-runtime reason", record.Reason)
	}
	if runner.preparedCalls != 0 {
		t.Fatalf("review prompt calls = %d, want 0 after Spec read failure", runner.preparedCalls)
	}
}

func TestReviewPromptWithoutSpecIsUnchanged(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	diff, err := reviewCandidateDiff(t.Context(), fixture.repository, fixture.baseCommit, fixture.headCommit, preflight.ExecGitRunner{})
	if err != nil {
		t.Fatalf("read candidate diff: %v", err)
	}
	wantPrompt := buildReviewPrompt(fixture.baseCommit, fixture.headCommit, diff)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if len(record.Specs) != 0 {
		t.Fatalf("review Specs = %v, want empty", record.Specs)
	}
	if runner.request.Prompt != wantPrompt {
		t.Fatalf("review prompt changed without a candidate Spec:\n--- got ---\n%s\n--- want ---\n%s", runner.request.Prompt, wantPrompt)
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
			name:        "plain preamble before clean verdict",
			answer:      "I reviewed the candidate.\nNo findings",
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

func TestReviewRecognisesEmphasizedFindingsHeader(t *testing.T) {
	tests := []struct {
		name         string
		answer       string
		wantFindings string
	}{
		{
			name:         "bold",
			answer:       "**Findings**:\ninternal/cli/review.go:42: bold header",
			wantFindings: "internal/cli/review.go:42: bold header",
		},
		{
			name:         "italic",
			answer:       "_Findings_:\ninternal/cli/review.go:42: italic header",
			wantFindings: "internal/cli/review.go:42: italic header",
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

			if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || record.Findings != test.wantFindings {
				t.Fatalf("emphasized findings review exit=%d record=%+v stderr=%q; want findings %q", code, record, stderr, test.wantFindings)
			}
		})
	}
}

func TestReviewBlocksEmphasizedFindingsBesideNoFindings(t *testing.T) {
	tests := []struct {
		name   string
		answer string
	}{
		{name: "bold", answer: "No findings.\n**Findings**:\ninternal/cli/review.go:42: contradictory verdict"},
		{name: "italic", answer: "No findings.\n_Findings_:\ninternal/cli/review.go:42: contradictory verdict"},
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

			assertBlockedReviewCommand(t, code, record, stderr, "both")
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

func TestReviewBlocksContentBesideANoFindingsVerdict(t *testing.T) {
	tests := []struct {
		name   string
		answer string
	}{
		{
			name:   "findings under a Findings heading",
			answer: "## Findings\n1. internal/cli/review.go:42: classifier accepts unrelated clean verdicts\n\n## Security\nNo findings",
		},
		{
			name:   "findings under a topic heading",
			answer: "## Correctness\ninternal/cli/review.go:42: classifier accepts topic findings\nNo findings",
		},
		{
			name:   "finding after the verdict",
			answer: "No findings\ninternal/cli/review.go:42: classifier ignores trailing findings",
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

			assertBlockedReviewCommand(t, code, record, stderr, "unclassifiable agent output")
		})
	}
}

func TestReviewReadsFindingsNoneAsNoFindings(t *testing.T) {
	tests := []struct {
		name        string
		answer      string
		wantCode    int
		wantOutcome reviewOutcome
	}{
		{name: "none", answer: "Findings: none", wantCode: exitOK, wantOutcome: reviewOutcomeReviewed},
		{name: "not applicable", answer: "Findings: n/a", wantCode: exitOK, wantOutcome: reviewOutcomeReviewed},
		{name: "no findings", answer: "Findings: no findings", wantCode: exitOK, wantOutcome: reviewOutcomeReviewed},
		{name: "plain preamble before none", answer: "I reviewed the candidate.\nFindings: none", wantCode: exitOK, wantOutcome: reviewOutcomeReviewed},
		{name: "structured content before none", answer: "## Correctness\n1. internal/cli/review.go:42: finding\nFindings: none", wantCode: exitPreflight, wantOutcome: reviewOutcomeBlocked},
		{name: "text follows none", answer: "Findings: none\ninternal/cli/review.go:42: finding follows", wantCode: exitRunFailed, wantOutcome: reviewOutcomeFindings},
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

			if code != test.wantCode || record.Outcome != test.wantOutcome {
				t.Fatalf("review exit=%d record=%+v stderr=%q, want exit=%d outcome=%q", code, record, stderr, test.wantCode, test.wantOutcome)
			}
		})
	}
}

func TestReviewIgnoresAQuotedVerdictInsideFindings(t *testing.T) {
	const answer = "Findings:\nNo findings\ninternal/cli/review.go:42: the findings body quoted the clean verdict"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: answer, StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("quoted verdict review exit=%d record=%+v stderr=%q, want findings", code, record, stderr)
	}
	if record.Findings != "No findings\ninternal/cli/review.go:42: the findings body quoted the clean verdict" {
		t.Fatalf("review findings = %q, want quoted verdict preserved", record.Findings)
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

func TestReviewKeepsNoAnswerWhenTheReviewerWasNotReached(t *testing.T) {
	runner := &reviewCommandRunner{
		prepareErrors: []error{errors.New("prepare review session")},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "runtime failure")
	if runner.preparedCalls != 0 {
		t.Fatalf("review prompt calls = %d, want 0", runner.preparedCalls)
	}
	if record.AnswerPath != "" {
		t.Fatalf("review answer path = %q, want empty", record.AnswerPath)
	}
	answerPath := filepath.Join(fixture.artifactDir, reviewAnswerFileName)
	if _, err := os.Stat(answerPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("review answer file stat error = %v, want not exist", err)
	}
}

func TestReviewRemovesAStaleAnswerFile(t *testing.T) {
	runner := &reviewCommandRunner{
		prepareErrors: []error{errors.New("prepare review session")},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	answerPath := filepath.Join(fixture.artifactDir, reviewAnswerFileName)
	if err := os.MkdirAll(fixture.artifactDir, 0o755); err != nil {
		t.Fatalf("create Artifact Directory: %v", err)
	}
	mustWrite(t, answerPath, "stale answer from an earlier review")

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "runtime failure")
	if _, err := os.Stat(answerPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stale review answer stat error = %v, want not exist", err)
	}
}

func TestReviewRecordsEmptySkippedSpecsAsAList(t *testing.T) {
	const slug = "0160-unreadable-skipped-specs"
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	specDir := filepath.Join(fixture.repository, "docs", "specs", slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create Spec directory: %v", err)
	}
	mustWrite(t, filepath.Join(specDir, "_prd.md"), "# Spec\n\n## Decisions\n\n- Read candidate files.\n")
	mustWrite(t, filepath.Join(specDir, "_techspec.md"), "# TechSpec\n")
	gittest.Run(t, fixture.repository, "add", "docs/specs")
	gittest.Run(t, fixture.repository, "commit", "-m", "add unreadable Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	withReviewSpecGitRunner(t, failingSpecFileGitRunner{delegate: preflight.ExecGitRunner{}})

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "Spec context read failure")
	if record.SkippedSpecs == nil {
		t.Fatal("review skipped Specs = nil, want an empty list")
	}
	recordBytes, err := os.ReadFile(filepath.Join(fixture.artifactDir, reviewRecordFileName))
	if err != nil {
		t.Fatalf("read review record: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(recordBytes, &fields); err != nil {
		t.Fatalf("decode review record fields: %v", err)
	}
	if string(fields["skippedSpecs"]) != "[]" {
		t.Fatalf("review skippedSpecs JSON = %s, want []", fields["skippedSpecs"])
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

func TestReviewRunsTheClaudeProvider(t *testing.T) {
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("Claude review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	fixture.assertCandidate(t, record)
	answer, err := os.ReadFile(record.AnswerPath)
	if err != nil {
		t.Fatalf("read Claude review answer: %v", err)
	}
	if string(answer) != "No findings" {
		t.Fatalf("Claude review answer = %q, want %q", string(answer), "No findings")
	}
	if !runner.request.Access.CanRead() || runner.request.Access.CanWrite() {
		t.Fatalf("Claude review access = %+v, want read-only", runner.request.Access)
	}
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 || runner.endCalls != 1 {
		t.Fatalf("Claude review calls: prepares=%d prompts=%d ends=%d, want 1 each", runner.prepareCalls, runner.preparedCalls, runner.endCalls)
	}
}

func TestReviewRefusesCodeRabbitNamingTheMissingSurface(t *testing.T) {
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "coderabbit", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "local CodeRabbit review surface")
	if record.Outcome == reviewOutcomeOmitted {
		t.Fatal("CodeRabbit refusal was recorded as omitted")
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("CodeRabbit refusal used Agent runtime: probes=%d prepares=%d prompts=%d", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
	}
}

func TestReviewCommandRefusesProviderProfileMismatch(t *testing.T) {
	tests := []struct {
		name              string
		provider          string
		preferredRuntime  string
		fallbackRuntime   string
		mismatchedRuntime string
	}{
		{
			name:              "preferred runtime",
			provider:          "codex",
			preferredRuntime:  "claude",
			fallbackRuntime:   "codex",
			mismatchedRuntime: "claude",
		},
		{
			name:              "fallback runtime",
			provider:          "codex",
			preferredRuntime:  "codex",
			fallbackRuntime:   "opencode",
			mismatchedRuntime: "opencode",
		},
		{
			name:              "Claude provider",
			provider:          "claude",
			preferredRuntime:  "codex",
			fallbackRuntime:   "claude",
			mismatchedRuntime: "codex",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{}
			fixture := newReviewCommandFixture(t, test.provider, runner)
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
			for _, want := range []string{test.provider, test.mismatchedRuntime} {
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

type failingSpecFileGitRunner struct {
	delegate preflight.GitRunner
}

func (runner failingSpecFileGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	if len(args) > 1 && args[0] == "show" && strings.Contains(args[1], "_prd.md") {
		return "", errors.New("candidate Spec object is unreadable")
	}
	return runner.delegate.RunGit(ctx, workDir, args...)
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
	request       agent.ExecuteRequest
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

func (runner *reviewCommandRunner) RunPrepared(_ context.Context, request agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	index := runner.preparedCalls
	runner.preparedCalls++
	runner.request = request
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

func writeReviewCommandConfigWithSpecsRoot(
	t *testing.T,
	repository string,
	provider string,
	artifactDir string,
	specsRoot string,
) {
	t.Helper()
	mustWrite(t, filepath.Join(repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: %s
specs:
  root: %q
`, artifactDir, provider, specsRoot))
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
	if !reflect.DeepEqual(fileRecord, stdoutRecord) {
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
