// Suite: review verdicts after permission refusal and prompt overflow
// Invariant: a refusal is recorded without replacing a delivered verdict; prompt overflow names the final answer.
// Boundary IN: review command and disposition ledger through the fake Agent runner.
// Boundary OUT: ACP parsing, real providers, and prompt admission bounds.

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

	"roundfix/internal/agent"
)

func TestReviewCommandClassifiesFindingsAfterARefusedPermission(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{
		Message:    "Findings:\n- review.txt:2: Failure: first regression\n- review.txt:2: Failure: second regression",
		StopReason: "end_turn", PermissionRefused: true,
	}}}}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")

	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || !record.PermissionRefused {
		t.Fatalf("review exit=%d record=%+v stderr=%q, want findings", code, record, stderr)
	}
	fixture.assertCandidate(t, record)
	if len(record.FindingItems) != 2 || record.FindingItems[0].ID != "F1" || record.FindingItems[1].ID != "F2" {
		t.Fatalf("finding items = %+v, want F1 and F2", record.FindingItems)
	}
	for _, finding := range record.FindingItems {
		if finding.Anchor == nil || finding.Validation == nil || finding.Validation.Status != "stands" {
			t.Fatalf("finding did not pass anchor validation: %+v", finding)
		}
	}
	// fixture.run compares decoded stdout and persisted records. Check the wire key too.
	data, err := os.ReadFile(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["permissionRefused"]) != "true" {
		t.Fatalf("permissionRefused = %s, want true", fields["permissionRefused"])
	}

	const evidence = "the failure is excluded by the supported caller contract"
	disposeCode, stdout, disposeStderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", evidence)
	if disposeCode != exitOK || disposeStderr != "" {
		t.Fatalf("dispose exit=%d stderr=%q, want success", disposeCode, disposeStderr)
	}
	line := readSingleDispositionLine(t, fixture.artifactDir)
	entry := decodeDispositionLine(t, line)
	if stdout != line || entry.Finding != "F1" || entry.Disposition != "dismissed" || entry.Evidence != evidence || entry.HeadCommit != record.HeadCommit || entry.Text != record.FindingItems[0].Text {
		t.Fatalf("disposition = %+v, stdout=%q ledger=%q", entry, stdout, line)
	}
}

func TestReviewCommandPassesNoFindingsAfterARefusedPermission(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{
		Message:    "The requested command was refused.\n\nNo findings.",
		Messages:   []string{"The requested command was refused.", "No findings."},
		StopReason: "end_turn", PermissionRefused: true,
	}}}}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || record.Reason != "" || !record.PermissionRefused {
		t.Fatalf("review exit=%d record=%+v stderr=%q, want reviewed", code, record, stderr)
	}
}

func TestReviewCommandNamesTheRefusalWhenTheAnswerHasNoVerdict(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ name, answer, reason string }{
		{"no verdict", "I could not finish checking.", "unclassifiable agent output: neither a no-findings nor findings verdict is present"},
		{"no anchor", "Findings:\n- Failure: a regression without an anchor", "findings name no file and line"},
	} {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{
				Message: test.answer, StopReason: "end_turn", PermissionRefused: true,
			}}}}
			fixture := newReviewCommandFixture(t, "claude", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")
			code, record, stderr := fixture.run(t)
			want := test.reason + " (after the read-only session refused a permission request)"
			assertBlockedReviewCommand(t, code, record, stderr, want)
			if record.Reason != want || !record.PermissionRefused {
				t.Fatalf("reason = %q, want %q", record.Reason, want)
			}
		})
	}
}

func TestReviewCommandNamesAPromptThatIsTooLong(t *testing.T) {
	t.Parallel()
	const line = "Prompt is too long · the request is ~1294791 tokens (limit 1000000) · reduce the prompt"
	const prefix = "Prompt is too long "
	longPrefix := prefix + strings.Repeat("a", reviewPromptTooLongLineLimit-len(prefix)-1)
	for _, test := range []struct {
		name     string
		result   agent.ExecuteResult
		err      error
		wantLine string
	}{
		{"runtime failure", agent.ExecuteResult{Message: line}, &agent.BatchFailureError{Reason: "agent/protocol error"}, line},
		{"parsed result", agent.ExecuteResult{Message: line, StopReason: "end_turn"}, nil, line},
		{"transport anomaly", agent.ExecuteResult{Message: line, TransportAnomaly: "unexpected exit"}, nil, line},
		{"trimmed later line", agent.ExecuteResult{Message: "Review failed.\n \t" + line + " \nNo findings."}, nil, line},
		{"rune boundary", agent.ExecuteResult{Message: longPrefix + "界" + strings.Repeat("a", 20)}, nil, longPrefix},
		{"final message", agent.ExecuteResult{Message: "progress\n\n" + line, Messages: []string{"progress", line}}, nil, line},
		{"refused permission", agent.ExecuteResult{Message: line, PermissionRefused: true}, nil, line},
	} {
		t.Run(test.name, func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: test.result, err: test.err}}}
			fixture := newReviewCommandFixture(t, "claude", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")
			code, record, stderr := fixture.run(t)
			want := "review prompt too long: " + test.wantLine
			if test.result.PermissionRefused {
				want += reviewPermissionRefusedSuffix
			}
			assertBlockedReviewCommand(t, code, record, stderr, want)
			if record.Reason != want || !utf8.ValidString(record.Reason) || len(test.wantLine) > reviewPromptTooLongLineLimit {
				t.Fatalf("reason = %q, want %q within byte and rune bounds", record.Reason, want)
			}
		})
	}
}

func TestReviewPermissionFlagIsOptionalAndPromptOverflowUsesFinalAnswer(t *testing.T) {
	t.Parallel()
	const progress = "Prompt is too long in a previous attempt."
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{
		Message: progress + "\n\nNo findings.", Messages: []string{progress, "No findings."}, StopReason: "end_turn",
	}}}}
	fixture := newReviewCommandFixture(t, "claude", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, fixture.provider, fixture.artifactDir, "claude", "claude")
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || record.PermissionRefused {
		t.Fatalf("review exit=%d record=%+v stderr=%q, want clean final verdict", code, record, stderr)
	}
	data, err := os.ReadFile(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if _, present := fields["permissionRefused"]; present {
		t.Fatalf("false permissionRefused must be omitted: %s", data)
	}
}

func TestReviewPromptTooLongKeepsFailurePrecedence(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		err  error
		want string
	}{
		{"admission", reviewPrePromptError{err: errors.New("admission refused")}, "admission refused"},
		{"timeout", context.DeadlineExceeded, "review timeout: context deadline exceeded"},
		{"Spec read", reviewSpecReadError{err: errors.New("unreadable Spec")}, "Spec context read failure: unreadable Spec"},
	} {
		t.Run(test.name, func(t *testing.T) {
			record, code := classifyReviewCommandResult(reviewRecord{Outcome: reviewOutcomeBlocked}, agent.ExecuteResult{Message: "Prompt is too long"}, test.err)
			if code != exitPreflight || record.Reason != test.want {
				t.Fatalf("exit=%d reason=%q, want %q", code, record.Reason, test.want)
			}
		})
	}
}
