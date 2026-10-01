// Suite: checkout-owned reviewer session lifecycle.
// Boundary IN: temporary repositories and public review/dispose commands.
// Boundary OUT: fake PrepareSession, RunPrepared and EndSession records.
package cli

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
)

type sessionReviewRunner struct {
	reviewCommandRunner
	prepares []agent.ExecuteRequest
	runs     []agent.ExecuteRequest
	ends     []agent.SessionRef
	events   []string
}

func (r *sessionReviewRunner) PrepareSession(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) error {
	r.prepares = append(r.prepares, req)
	r.events = append(r.events, "prepare:"+req.Session.Name)
	return r.reviewCommandRunner.PrepareSession(ctx, req, sink)
}
func (r *sessionReviewRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	r.runs = append(r.runs, req)
	return r.reviewCommandRunner.RunPrepared(ctx, req, sink)
}
func (r *sessionReviewRunner) EndSession(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
	r.ends = append(r.ends, session)
	r.events = append(r.events, "end:"+session.Name)
	return r.reviewCommandRunner.EndSession(ctx, runtime, session)
}
func newSessionReviewRunner(answers ...string) *sessionReviewRunner {
	r := &sessionReviewRunner{}
	for _, answer := range answers {
		r.results = append(r.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: answer, StopReason: "end_turn", ACPSessionID: "s-1"}})
	}
	return r
}

const sessionReviewFinding = "Findings:\n- review.txt:2 Failure: candidate defect"

func TestReviewRoundOneWithFindingsLeavesItsSessionOpen(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding)
	fixture := newReviewCommandFixture(t, "codex", r)
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if len(r.ends) != 0 || record.Lineage.Session != r.prepares[0].Session.Name || record.Lineage.Selection != 0 || !record.Lineage.SessionOpen {
		t.Fatalf("lineage=%+v ended=%+v", record.Lineage, r.ends)
	}
	persisted, err := readReviewRecord(reviewCheckoutDir(fixture.artifactDir, fixture.repository) + "/" + reviewRecordFileName)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(record.Lineage, persisted.Lineage) {
		t.Fatalf("persisted lineage=%+v want=%+v", persisted.Lineage, record.Lineage)
	}
}
func TestReviewRoundOneReviewedEndsItsSession(t *testing.T) {
	r := newSessionReviewRunner("No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	code, record, _ := runLineageReview(t, fixture)
	if code != exitOK || record.Lineage.SessionOpen || len(r.ends) != 1 || r.ends[0] != r.prepares[0].Session {
		t.Fatalf("exit=%d lineage=%+v ends=%+v", code, record.Lineage, r.ends)
	}
}
func TestReviewRoundOneBlockedEndsItsSession(t *testing.T) {
	r := newSessionReviewRunner("not a verdict")
	fixture := newReviewCommandFixture(t, "codex", r)
	code, record, _ := runLineageReview(t, fixture)
	if code != exitPreflight || record.Lineage.SessionOpen || len(r.ends) != 1 {
		t.Fatalf("exit=%d lineage=%+v ends=%+v", code, record.Lineage, r.ends)
	}
}
func TestReviewRoundOneValidationDismissedEndsItsSession(t *testing.T) {
	r := newSessionReviewRunner("Findings:\n- outside.txt:10 Failure: outside candidate")
	fixture := newReviewCommandFixture(t, "codex", r)
	code, record, _ := runLineageReview(t, fixture)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || record.Lineage.SessionOpen || len(r.ends) != 1 {
		t.Fatalf("exit=%d record=%+v ends=%+v", code, record, r.ends)
	}
}
func TestReviewRoundTwoContinuesTheRecordedSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	head := commitLineageFile(t, fixture, "delta.txt", "round two correction\n")
	disposeLineageFinding(t, "--fixed-by", head)
	code, second, stderr := runLineageReview(t, fixture)
	if code != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || !second.Lineage.Continued || second.Lineage.SessionOpen {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session != r.prepares[0].Session || r.prepares[1].Runtime != r.prepares[0].Runtime || len(r.ends) != 1 || r.ends[0] != r.prepares[0].Session {
		t.Fatalf("prepares=%+v ends=%+v", r.prepares, r.ends)
	}
	if !reflect.DeepEqual(second.Lineage.ACPSessionIDs, []string{"s-1", "s-1"}) || second.Lineage.PreviousHead != first.HeadCommit || len(second.Lineage.PreviousFindings) != 1 || len(second.Lineage.PreviousDispositions) != 1 {
		t.Fatalf("lineage=%+v", second.Lineage)
	}
	if !strings.Contains(r.runs[1].Prompt, "--- BEGIN ROUND 2 DELTA DIFF ---") || !strings.Contains(r.runs[1].Prompt, "fixed by "+head) || strings.Contains(r.runs[1].Prompt, "--- BEGIN CANDIDATE DIFF ---") {
		t.Fatalf("round 2 prompt=%q", r.runs[1].Prompt)
	}
}
func TestReviewRoundTwoDifferentIDsDoesNotReportContinued(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	r.results[1].result.ACPSessionID = "s-2"
	fixture := newReviewCommandFixture(t, "codex", r)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	_, record, _ := runLineageReview(t, fixture)
	if record.Lineage.Continued || !reflect.DeepEqual(record.Lineage.ACPSessionIDs, []string{"s-1", "s-2"}) {
		t.Fatalf("lineage=%+v", record.Lineage)
	}
}
func TestReviewRoundTwoEmptyIDsDoesNotReportContinued(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	r.results[0].result.ACPSessionID = ""
	r.results[1].result.ACPSessionID = ""
	fixture := newReviewCommandFixture(t, "codex", r)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	_, record, _ := runLineageReview(t, fixture)
	if record.Lineage.Continued || !reflect.DeepEqual(record.Lineage.ACPSessionIDs, []string{"", ""}) {
		t.Fatalf("lineage=%+v", record.Lineage)
	}
}
func TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	r.prepareErrors = []error{nil, errors.New("cannot resume"), nil}
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || record.Lineage.Continued || record.Lineage.SessionOpen {
		t.Fatalf("exit=%d lineage=%+v stderr=%q", code, record.Lineage, stderr)
	}
	if len(r.prepares) != 3 || r.prepares[1].Session.Name != first.Lineage.Session || r.prepares[2].Session.Name == first.Lineage.Session || len(r.ends) != 2 || r.ends[0].Name != first.Lineage.Session || r.ends[1] != r.prepares[2].Session || r.runs[1].Session != r.prepares[2].Session {
		t.Fatalf("prepares=%+v runs=%+v ends=%+v", r.prepares, r.runs, r.ends)
	}
}
func TestReviewLineageChangeEndsTheOpenSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	gittest.Run(t, fixture.repository, "checkout", "main")
	commitLineageFile(t, fixture, "upstream.txt", "upstream\n")
	gittest.Run(t, fixture.repository, "checkout", "feature/review")
	gittest.Run(t, fixture.repository, "rebase", "main")
	code, record, _ := runLineageReview(t, fixture)
	if code != exitOK || record.Lineage.Round != 1 || record.Lineage.Continued {
		t.Fatalf("exit=%d lineage=%+v", code, record.Lineage)
	}
	if len(r.ends) != 2 || r.ends[0].Name != first.Lineage.Session || r.events[1] != "end:"+first.Lineage.Session || r.events[2] != "prepare:"+record.Lineage.Session {
		t.Fatalf("events=%+v ends=%+v", r.events, r.ends)
	}
	if !reflect.DeepEqual(record.Lineage.ACPSessionIDs, []string{"s-1"}) {
		t.Fatalf("new lineage ids=%+v", record.Lineage.ACPSessionIDs)
	}
}
func TestReviewDismissedReuseEndsTheOpenSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding)
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	disposeLineageFinding(t, "--dismiss", "--evidence", "documented operator evidence")
	code, record, _ := runLineageReview(t, fixture)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || !record.Reused || record.Lineage.SessionOpen || len(r.ends) != 1 || r.ends[0].Name != first.Lineage.Session || len(r.prepares) != 1 || len(r.runs) != 1 {
		t.Fatalf("exit=%d record=%+v ends=%+v", code, record, r.ends)
	}
	runLineageReview(t, fixture)
	if len(r.ends) != 1 {
		t.Fatalf("closed reuse ended twice: %+v", r.ends)
	}
}
func TestReviewStandingReuseKeepsTheOpenSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding)
	fixture := newReviewCommandFixture(t, "codex", r)
	runLineageReview(t, fixture)
	code, record, _ := runLineageReview(t, fixture)
	if code != exitRunFailed || !record.Lineage.SessionOpen || len(r.ends) != 0 || len(r.runs) != 1 {
		t.Fatalf("exit=%d lineage=%+v ends=%+v", code, record.Lineage, r.ends)
	}
}
func TestReviewRoundTwoFindingsEndsItsSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, sessionReviewFinding)
	fixture := newReviewCommandFixture(t, "codex", r)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, record, _ := runLineageReview(t, fixture)
	if code != exitRunFailed || record.Lineage.SessionOpen || len(r.ends) != 1 {
		t.Fatalf("exit=%d lineage=%+v ends=%+v", code, record.Lineage, r.ends)
	}
}

func TestReviewRoundTwoContinuesTheRecordedFallbackSelection(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	r.prepareErrors = []error{&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, nil, nil}
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	if first.Lineage.Selection != 1 {
		t.Fatalf("first selection=%d", first.Lineage.Selection)
	}
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	_, second, _ := runLineageReview(t, fixture)
	if second.Lineage.Selection != 1 || !second.Lineage.Continued || r.prepares[2].Runtime != r.prepares[1].Runtime || r.prepares[2].Session != r.prepares[1].Session || len(r.ends) != 2 || r.ends[1] != r.prepares[1].Session {
		t.Fatalf("lineage=%+v prepares=%+v ends=%+v", second.Lineage, r.prepares, r.ends)
	}
}
func TestReviewRoundTwoPromptFailureEndsWithoutFallback(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding, "")
	r.results[1].err = &agent.SelectionFailureError{Runtime: "codex", Reason: "adapter stopped after prompt"}
	fixture := newReviewCommandFixture(t, "codex", r)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, record, _ := runLineageReview(t, fixture)
	if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Lineage.SessionOpen || len(r.prepares) != 2 || len(r.ends) != 1 || r.ends[0] != r.prepares[0].Session {
		t.Fatalf("exit=%d record=%+v prepares=%+v ends=%+v", code, record, r.prepares, r.ends)
	}
}
func TestReviewProviderOmissionEndsTheOpenSession(t *testing.T) {
	r := newSessionReviewRunner(sessionReviewFinding)
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _ := runLineageReview(t, fixture)
	writeReviewCommandConfig(t, fixture.repository, "none", fixture.artifactDir)
	code, record, _ := runLineageReview(t, fixture)
	if code != exitOK || record.Outcome != reviewOutcomeOmitted || len(r.runs) != 1 || len(r.ends) != 1 || r.ends[0].Name != first.Lineage.Session {
		t.Fatalf("exit=%d record=%+v ends=%+v", code, record, r.ends)
	}
}
