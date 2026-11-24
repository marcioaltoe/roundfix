package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
)

// Read lineage.model from JSON so these tests also compile against older records.

type modelNoticeReviewRunner struct {
	sessionReviewRunner
	opened map[string]string
}

func (r *modelNoticeReviewRunner) PrepareSession(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) error {
	if r.opened == nil {
		r.opened = map[string]string{}
	}
	if _, ok := r.opened[req.Session.Name]; !ok {
		r.opened[req.Session.Name] = req.Runtime.Model
	}
	return r.sessionReviewRunner.PrepareSession(ctx, req, sink)
}

func (r *modelNoticeReviewRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	result, err := r.sessionReviewRunner.RunPrepared(ctx, req, sink)
	if was := r.opened[req.Session.Name]; was != "" && was != req.Runtime.Model {
		result.Message = reportedModelNotice(was, req.Runtime.Model) + result.Message
	}
	return result, err
}

func reportedModelNotice(was, now string) string {
	return fmt.Sprintf("Warning: This session was recorded with model `%s` but is resuming with `%s`. Consider switching back to `%s` as it may affect Codex performance.\n\n", was, now, was)
}

func newModelNoticeRunner(answers ...string) *modelNoticeReviewRunner {
	r := &modelNoticeReviewRunner{}
	for _, answer := range answers {
		r.results = append(r.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: answer, StopReason: "end_turn", ACPSessionID: "s-1"}})
	}
	return r
}

func writeModelReviewProfile(t *testing.T, fixture reviewCommandFixture, model string) {
	t.Helper()
	mustWrite(t, filepath.Join(fixture.repository, ".roundfixrc.yml"), fmt.Sprintf(`defaults:
  artifact_dir: %q
pre_pr_review:
  provider: codex
profiles:
  review:
    preferred:
      runtime: codex
      model: %s
      reasoning_effort: high
    fallbacks:
      - runtime: codex
        model: fallback-model
        reasoning_effort: high
`, fixture.artifactDir, model))
}

// runModelReview checks persistence and returns the printed record and model.
func runModelReview(t *testing.T, fixture reviewCommandFixture) (int, reviewRecord, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode %q: %v (stderr %q)", stdout.String(), err, stderr.String())
	}
	var raw struct {
		Lineage map[string]any `json:"lineage"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	model, _ := raw.Lineage["model"].(string)
	content, err := os.ReadFile(modelReviewRecordPath(fixture))
	if err != nil {
		t.Fatal(err)
	}
	var persisted struct {
		Lineage map[string]any `json:"lineage"`
	}
	if err := json.Unmarshal(content, &persisted); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(raw.Lineage, persisted.Lineage) {
		t.Fatalf("printed lineage=%v persisted lineage=%v", raw.Lineage, persisted.Lineage)
	}
	return code, record, model, stderr.String()
}

func modelReviewRecordPath(fixture reviewCommandFixture) string {
	return filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName)
}

func TestReviewRoundTwoOpensANewSessionWhenTheConfiguredModelChanged(t *testing.T) {
	t.Parallel()
	r := newModelNoticeRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	writeModelReviewProfile(t, fixture, "gpt-6-luna")
	code1, first, model1, _ := runModelReview(t, fixture)
	if code1 != exitRunFailed || !first.Lineage.SessionOpen || model1 != "gpt-6-luna" {
		t.Fatalf("round 1 exit=%d lineage=%+v model=%q", code1, first.Lineage, model1)
	}
	head := commitLineageFile(t, fixture, "delta.txt", "round two correction\n")
	disposeLineageFinding(t, "--fixed-by", head)
	writeModelReviewProfile(t, fixture, "gpt-6.1-sol")
	code2, second, model2, stderr := runModelReview(t, fixture)
	if code2 != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || second.Lineage.Continued || second.Lineage.SessionOpen || stderr != "" || model2 != "gpt-6.1-sol" {
		t.Fatalf("round 2 exit=%d outcome=%s reason=%q lineage=%+v model=%q stderr=%q", code2, second.Outcome, second.Reason, second.Lineage, model2, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session.Name == r.prepares[0].Session.Name || r.prepares[0].Runtime.Model != "gpt-6-luna" || r.prepares[1].Runtime.Model != "gpt-6.1-sol" {
		t.Fatalf("prepares=%+v", r.prepares)
	}
	old, fresh := r.prepares[0].Session.Name, r.prepares[1].Session.Name
	if want := []string{"prepare:" + old, "end:" + old, "prepare:" + fresh, "end:" + fresh}; !reflect.DeepEqual(r.events, want) {
		t.Fatalf("events=%v want=%v", r.events, want)
	}
	if !strings.Contains(r.runs[1].Prompt, "--- BEGIN ROUND 2 DELTA DIFF ---") || !strings.Contains(r.runs[1].Prompt, "fixed by "+head) || strings.Contains(r.runs[1].Prompt, "--- BEGIN CANDIDATE DIFF ---") {
		t.Fatalf("round 2 prompt=%q", r.runs[1].Prompt)
	}
	t.Run("the_classifier_still_refuses_the_reported_answer", func(t *testing.T) {
		only := newSessionReviewRunner(reportedModelNotice("gpt-6-luna", "gpt-6.1-sol") + "No findings.")
		other := newReviewCommandFixture(t, "codex", only)
		code, record, _ := runLineageReview(t, other)
		if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Reason != "unclassifiable agent output: no-findings verdict does not account for other content" {
			t.Fatalf("exit=%d outcome=%s reason=%q", code, record.Outcome, record.Reason)
		}
	})
}

func TestReviewRoundTwoOpensANewSessionWhenTheRecordHoldsNoModel(t *testing.T) {
	t.Parallel()
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _, _ := runModelReview(t, fixture)
	path := modelReviewRecordPath(fixture)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err := json.Unmarshal(content, &record); err != nil {
		t.Fatal(err)
	}
	delete(record["lineage"].(map[string]any), "model")
	rewritten, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rewritten, 0o644); err != nil {
		t.Fatal(err)
	}
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, second, model, stderr := runModelReview(t, fixture)
	if code != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || second.Lineage.Continued || second.Lineage.SessionOpen || model != r.prepares[1].Runtime.Model || stderr != "" {
		t.Fatalf("exit=%d lineage=%+v model=%q stderr=%q", code, second.Lineage, model, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session.Name == first.Lineage.Session || len(r.ends) != 2 || r.ends[0].Name != first.Lineage.Session {
		t.Fatalf("prepares=%+v ends=%+v", r.prepares, r.ends)
	}
	if want := []string{"prepare:" + first.Lineage.Session, "end:" + first.Lineage.Session, "prepare:" + r.prepares[1].Session.Name, "end:" + r.prepares[1].Session.Name}; !reflect.DeepEqual(r.events, want) {
		t.Fatalf("events=%v want=%v", r.events, want)
	}
}

func TestReviewRecordsTheModelOfTheSessionOfEachRound(t *testing.T) {
	t.Parallel()
	t.Run("preferred_then_continued", func(t *testing.T) {
		r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
		fixture := newReviewCommandFixture(t, "codex", r)
		_, _, model1, _ := runModelReview(t, fixture)
		if model1 == "" || model1 != r.prepares[0].Runtime.Model {
			t.Fatalf("round 1 model=%q runtime=%q", model1, r.prepares[0].Runtime.Model)
		}
		content, err := os.ReadFile(modelReviewRecordPath(fixture))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), fmt.Sprintf("%q:%q", "model", model1)) && !strings.Contains(string(content), fmt.Sprintf("%q: %q", "model", model1)) {
			t.Fatalf("record file holds no model: %s", content)
		}
		commitLineageFile(t, fixture, "delta.txt", "fix\n")
		_, second, model2, _ := runModelReview(t, fixture)
		if !second.Lineage.Continued || model2 != model1 {
			t.Fatalf("round 2 lineage=%+v model=%q want %q", second.Lineage, model2, model1)
		}
	})
	t.Run("fallback", func(t *testing.T) {
		r := newSessionReviewRunner(sessionReviewFinding)
		r.prepareErrors = []error{&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, &agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, nil}
		fixture := newReviewCommandFixture(t, "codex", r)
		_, first, model, _ := runModelReview(t, fixture)
		if first.Lineage.Selection == nil || *first.Lineage.Selection != 1 || model == "" || model != r.prepares[2].Runtime.Model || model == r.prepares[0].Runtime.Model {
			t.Fatalf("lineage=%+v model=%q prepares=%q,%q", first.Lineage, model, r.prepares[0].Runtime.Model, r.prepares[2].Runtime.Model)
		}
	})
}

func TestReviewRoundTwoContinuesTheSessionWhenTheRecordedModelDiffersOnlyByWhitespace(t *testing.T) {
	t.Parallel()
	r := newSessionReviewRunner(sessionReviewFinding, "No findings.")
	fixture := newReviewCommandFixture(t, "codex", r)
	_, first, _, _ := runModelReview(t, fixture)
	path := modelReviewRecordPath(fixture)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err := json.Unmarshal(content, &record); err != nil {
		t.Fatal(err)
	}
	lineage := record["lineage"].(map[string]any)
	recorded, _ := lineage["model"].(string)
	if recorded == "" {
		t.Fatalf("round 1 recorded no model: %s", content)
	}
	lineage["model"] = "  " + recorded + "\n"
	rewritten, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, rewritten, 0o644); err != nil {
		t.Fatal(err)
	}
	commitLineageFile(t, fixture, "delta.txt", "fix\n")
	code, second, _, stderr := runModelReview(t, fixture)
	if code != exitOK || second.Outcome != reviewOutcomeReviewed || second.Lineage.Round != 2 || !second.Lineage.Continued || stderr != "" {
		t.Fatalf("exit=%d lineage=%+v stderr=%q", code, second.Lineage, stderr)
	}
	if len(r.prepares) != 2 || r.prepares[1].Session.Name != first.Lineage.Session {
		t.Fatalf("round 2 did not continue the round-1 session: prepares=%+v", r.prepares)
	}
}
