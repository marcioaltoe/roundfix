// Suite: review selection retry policy.
// Boundary IN: public review commands in temporary repositories.
// Boundary OUT: persisted records, diagnostics, and fake session calls.
package cli

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
)

func retryProtocolFailure(sent bool) error {
	step := "session/set_model"
	if sent {
		step = "session/prompt"
	}
	return &agent.SelectionFailureError{Runtime: "codex", Reason: "agent/protocol error", Protocol: &agent.ProtocolFailure{Step: step, Message: "adapter refused the request", Code: -32603, PromptSent: sent}}
}

func retryReviewFixture(t *testing.T, runner agent.Runner) reviewCommandFixture {
	t.Helper()
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeReviewCommandProfileConfig(t, fixture.repository, "codex", fixture.artifactDir, "codex", "codex")

	return fixture
}

func assertReviewRetry(t *testing.T, record reviewRecord, stderr, step, message string, selection int) {
	t.Helper()
	want := []reviewSelectionRetry{{Selection: selection, Step: step, Message: message}}
	if !reflect.DeepEqual(record.SelectionRetries, want) {
		t.Fatalf("retries=%+v want=%+v", record.SelectionRetries, want)
	}
	notice := "roundfix: review Agent Selection failed before the prompt at " + step + " (" + message + "); retrying selection "
	if strings.Count(stderr, notice) != 1 {
		t.Fatalf("retry notice missing or repeated: %q", stderr)
	}
	if record.Provider != "codex" || record.Outcome == reviewOutcomeOmitted {
		t.Fatalf("failure changed policy: %+v", record)
	}
}

func TestReviewRetriesASelectionFailureBeforeThePromptOnce(t *testing.T) {
	r := newSessionReviewRunner("No findings.")
	r.prepareErrors = []error{&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, nil}
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(r.prepares) != 2 || len(r.runs) != 1 {
		t.Fatalf("exit=%d record=%+v prepares=%+v runs=%+v", code, record, r.prepares, r.runs)
	}
	if r.prepares[0].Runtime != r.prepares[1].Runtime || r.prepares[0].Session != r.prepares[1].Session || len(r.ends) != 2 || r.ends[0] != r.prepares[0].Session || r.events[1] != "end:"+r.prepares[0].Session.Name {
		t.Fatalf("retry did not close and prepare the same selection: prepares=%+v ends=%+v events=%+v", r.prepares, r.ends, r.events)
	}
	assertReviewRetry(t, record, stderr, "session setup", "", 0)
	if strings.Contains(stderr, "activating fallback") || record.FailedStep != "" || record.AdapterMessage != "" {
		t.Fatalf("unexpected failure metadata: %+v stderr=%q", record, stderr)
	}
}

func TestReviewRetriesAPromptProcessThatFailedBeforeSessionPrompt(t *testing.T) {
	r := newSessionReviewRunner("", "No findings.")
	r.results[0].err = retryProtocolFailure(false)
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(r.prepares) != 1 || len(r.runs) != 2 || r.runs[0].Session != r.runs[1].Session || r.runs[0].Runtime != r.runs[1].Runtime {
		t.Fatalf("exit=%d record=%+v prepares=%+v runs=%+v", code, record, r.prepares, r.runs)
	}
	assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
	if strings.Contains(stderr, "activating fallback") {
		t.Fatalf("unexpected fallback: %q", stderr)
	}
}

func TestReviewNeverRetriesAfterSessionPromptWasSent(t *testing.T) {
	r := &reviewCommandRunner{results: []reviewCommandRunResult{{err: retryProtocolFailure(true)}}}
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure:")
	if r.prepareCalls != 1 || r.preparedCalls != 1 || record.FailedStep != "session/prompt" || record.AdapterMessage != "adapter refused the request" || len(record.SelectionRetries) != 0 || strings.Contains(stderr, "retrying selection") || strings.Contains(stderr, "activating fallback") {
		t.Fatalf("record=%+v calls=%d/%d stderr=%q", record, r.prepareCalls, r.preparedCalls, stderr)
	}
}

func TestReviewActivatesTheFallbackAfterTheRetryFails(t *testing.T) {
	for _, phase := range []string{"preparation", "prompt"} {
		t.Run(phase, func(t *testing.T) {
			r := newSessionReviewRunner("No findings.")
			if phase == "preparation" {
				r.prepareErrors = []error{retryProtocolFailure(false), retryProtocolFailure(false), nil}
			} else {
				r.results = []reviewCommandRunResult{{err: retryProtocolFailure(false)}, {err: retryProtocolFailure(false)}, {result: agent.ExecuteResult{Message: "No findings."}}}
			}
			fixture := retryReviewFixture(t, r)
			code, record, stderr := fixture.run(t)
			if code != exitOK || record.Outcome != reviewOutcomeReviewed || record.Lineage.Selection == nil || *record.Lineage.Selection != 1 || r.request.Runtime.Model != "fallback-model" {
				t.Fatalf("exit=%d record=%+v request=%+v", code, record, r.request)
			}
			assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
			if strings.Count(stderr, "activating fallback 1") != 1 || strings.Index(stderr, "retrying selection") > strings.Index(stderr, "activating fallback") {
				t.Fatalf("stderr=%q", stderr)
			}
			if phase == "preparation" && (len(r.prepares) != 3 || len(r.runs) != 1) || phase == "prompt" && (len(r.prepares) != 2 || len(r.runs) != 3) {
				t.Fatalf("prepares=%+v runs=%+v", r.prepares, r.runs)
			}
		})
	}
}

func TestReviewBlocksAfterTheRetryFailsWithoutFallback(t *testing.T) {
	// Profiles require a fallback, so exercise the last selection with no fallback left.
	for _, phase := range []string{"preparation", "prompt"} {
		t.Run(phase, func(t *testing.T) {
			r := &reviewCommandRunner{}
			failures := []error{retryProtocolFailure(false), retryProtocolFailure(false), retryProtocolFailure(false), retryProtocolFailure(false)}
			if phase == "preparation" {
				r.prepareErrors = failures
			} else {
				for _, err := range failures {
					r.results = append(r.results, reviewCommandRunResult{err: err})
				}
			}
			fixture := retryReviewFixture(t, r)
			code, record, stderr := fixture.run(t)
			assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure:")
			want := []reviewSelectionRetry{
				{Selection: 0, Step: "session/set_model", Message: "adapter refused the request"},
				{Selection: 1, Step: "session/set_model", Message: "adapter refused the request"},
			}
			if !reflect.DeepEqual(record.SelectionRetries, want) || record.Provider != "codex" || !strings.HasSuffix(record.Reason, " (after one automatic retry before the prompt)") || record.FailedStep != "session/set_model" || record.AdapterMessage != "adapter refused the request" || strings.Count(stderr, "retrying selection") != 2 || strings.Count(stderr, "activating fallback") != 1 {
				t.Fatalf("record=%+v stderr=%q", record, stderr)
			}
			if phase == "preparation" && (r.prepareCalls != 4 || r.preparedCalls != 0) || phase == "prompt" && (r.prepareCalls != 2 || r.preparedCalls != 4) {
				t.Fatalf("calls=%d/%d", r.prepareCalls, r.preparedCalls)
			}
		})
	}
}

func TestReviewDoesNotRetryAGenericPreparationError(t *testing.T) {
	r := &reviewCommandRunner{prepareErrors: []error{errors.New("session filesystem unavailable")}}
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure: session filesystem unavailable")
	if r.prepareCalls != 1 || r.preparedCalls != 0 || len(record.SelectionRetries) != 0 || record.FailedStep != "" || strings.Contains(stderr, "retrying selection") || strings.Contains(stderr, "activating fallback") {
		t.Fatalf("record=%+v calls=%d/%d stderr=%q", record, r.prepareCalls, r.preparedCalls, stderr)
	}
}

func TestReviewRoundTwoRetriesAResumedPromptBeforeSessionPrompt(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "", "No findings.")
	r.results[1].err = retryProtocolFailure(false)
	fixture := retryReviewFixture(t, r)
	_, first, _ := runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || record.Lineage.Round != 2 || !record.Lineage.Continued || len(r.prepares) != 2 || len(r.runs) != 3 || r.runs[1].Session.Name != first.Lineage.Session || r.runs[1].Session != r.runs[2].Session || len(r.ends) != 1 {
		t.Fatalf("exit=%d record=%+v prepares=%+v runs=%+v ends=%+v", code, record, r.prepares, r.runs, r.ends)
	}
	assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
	if strings.Contains(stderr, "activating fallback") {
		t.Fatalf("stderr=%q", stderr)
	}
}

func TestReviewSelectionRetryStopsOnCancellationOrUnplacedPromptFailure(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
	}{
		{"canceled", errors.Join(retryProtocolFailure(false), context.Canceled)},
		{"deadline", errors.Join(retryProtocolFailure(false), context.DeadlineExceeded)},
		{"unplaced", &agent.SelectionFailureError{Runtime: "codex", Reason: "unknown prompt failure"}},
	} {
		for _, phase := range []string{"preparation", "prompt"} {
			if phase == "preparation" && test.name == "unplaced" {
				continue
			}
			t.Run(test.name+"/"+phase, func(t *testing.T) {
				r := &reviewCommandRunner{}
				if phase == "preparation" {
					r.prepareErrors = []error{test.err}
				} else {
					r.results = []reviewCommandRunResult{{err: test.err}}
				}
				fixture := retryReviewFixture(t, r)
				code, record, stderr := fixture.run(t)
				if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Provider != "codex" || len(record.SelectionRetries) != 0 || r.prepareCalls != 1 || r.preparedCalls > 1 || strings.Contains(stderr, "retrying selection") || strings.Contains(stderr, "activating fallback") {
					t.Fatalf("exit=%d record=%+v calls=%d/%d stderr=%q", code, record, r.prepareCalls, r.preparedCalls, stderr)
				}
			})
		}
	}
}

func TestReviewSelectionRetryBudgetSpansPreparationAndPrompt(t *testing.T) {
	r := newSessionReviewRunner("", "No findings.")
	r.prepareErrors = []error{retryProtocolFailure(false), nil}
	r.results[0].err = retryProtocolFailure(false)
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || len(r.prepares) != 3 || len(r.runs) != 2 || r.runs[0].Runtime.Model != "preferred-model" || r.runs[1].Runtime.Model != "fallback-model" {
		t.Fatalf("exit=%d record=%+v prepares=%+v runs=%+v", code, record, r.prepares, r.runs)
	}
	assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
}

func TestReviewFailureAfterRetryDoesNotActivateFallbackAfterPrompt(t *testing.T) {
	r := &reviewCommandRunner{results: []reviewCommandRunResult{{err: retryProtocolFailure(false)}, {err: retryProtocolFailure(true)}}}
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure:")
	assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
	if r.prepareCalls != 1 || r.preparedCalls != 2 || record.FailedStep != "session/prompt" || !strings.HasSuffix(record.Reason, " (after one automatic retry before the prompt)") || strings.Contains(stderr, "activating fallback") {
		t.Fatalf("record=%+v calls=%d/%d stderr=%q", record, r.prepareCalls, r.preparedCalls, stderr)
	}
}

func TestReviewPreviousSelectionRetryDoesNotSuffixUnretriedFallbackFailure(t *testing.T) {
	r := &reviewCommandRunner{prepareErrors: []error{retryProtocolFailure(false), retryProtocolFailure(false)}, results: []reviewCommandRunResult{{err: retryProtocolFailure(true)}}}
	fixture := retryReviewFixture(t, r)
	code, record, stderr := fixture.run(t)
	assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure:")
	assertReviewRetry(t, record, stderr, "session/set_model", "adapter refused the request", 0)
	if r.prepareCalls != 3 || r.preparedCalls != 1 || record.FailedStep != "session/prompt" || strings.Contains(record.Reason, "after one automatic retry") {
		t.Fatalf("record=%+v calls=%d/%d", record, r.prepareCalls, r.preparedCalls)
	}
}
