// Boundary: public review/dispose commands, real temporary Git history and records.
// Reviewer calls and sessions are captured by the existing review fake.
package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/delivery"
)

func fixedRoundTwoFixture(t *testing.T, answers ...string) (reviewCommandFixture, *sessionReviewRunner, reviewRecord, []string) {
	t.Helper()
	runner := newSessionReviewRunner(append([]string{"No findings.", "Findings:\n- review.txt:2 Failure: first defect\n- review.txt:1 Failure: second defect"}, answers...)...)
	fixture := newReviewCommandFixture(t, "codex", runner)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "round-two.txt", "round two content excluded from correction\n")
	code, prior, stderr := runLineageReview(t, fixture)
	if code != exitRunFailed || prior.Lineage.Round != 2 || len(prior.FindingItems) != 2 {
		t.Fatalf("round two code=%d record=%+v stderr=%q", code, prior, stderr)
	}
	var fixes []string
	for i, finding := range prior.FindingItems {
		fix := commitLineageFile(t, fixture, []string{"fix-one.txt", "fix-two.txt"}[i], []string{"first structural fix\n", "second structural fix\n"}[i])
		fixes = append(fixes, fix)
		disposeCorrectionFinding(t, finding.ID, "--fixed-by", fix)
	}
	return fixture, runner, prior, fixes
}

func disposeCorrectionFinding(t *testing.T, id string, args ...string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, append([]string{"review", "dispose", id}, args...), &stdout, &stderr); code != exitOK {
		t.Fatalf("dispose %s code=%d stdout=%q stderr=%q", id, code, stdout.String(), stderr.String())
	}
}

func TestCorrectionReviewReviewsTheFixOnce(t *testing.T) {
	t.Parallel()
	for _, answer := range []string{"No findings.", "Findings:\n- outside.txt:9 Failure: unanchored"} {
		t.Run(answer, func(t *testing.T) {
			fixture, runner, prior, fixes := fixedRoundTwoFixture(t, answer)
			calls := len(runner.runs)
			code, record, stderr := runLineageReview(t, fixture)
			wantStderr := ""
			if strings.HasPrefix(answer, "Findings:") {
				wantStderr = "roundfix: review finding F1 dismissed by validation (unanchored): anchor names no line of the candidate diff\n"
			}
			if code != exitOK || stderr != wantStderr || record.Lineage == nil || !record.Lineage.Correction || record.Lineage.Round != 2 || len(runner.runs) != calls+1 {
				t.Fatalf("correction code=%d record=%+v calls=%d stderr=%q", code, record, len(runner.runs)-calls, stderr)
			}
			if record.Lineage.PreviousHead != prior.HeadCommit || !reflect.DeepEqual(record.Lineage.PreviousFindings, prior.FindingItems) || len(record.Lineage.PreviousDispositions) != 2 || len(record.Segments) != 0 || record.Lineage.SessionOpen || record.Lineage.Continued {
				t.Fatalf("correction lineage=%+v segments=%+v", record.Lineage, record.Segments)
			}
			if runner.runs[calls].Session == runner.runs[calls-1].Session || runner.ends[len(runner.ends)-1] != runner.runs[calls].Session {
				t.Fatal("Correction Review did not open and end a fresh session")
			}
			prompt := runner.runs[calls].Prompt
			wantSentence := "This is the Correction Review of this candidate: round 2 reviewed " + prior.HeadCommit + " and raised the findings below, each fixed by the commit named. The diff below is the change since " + prior.HeadCommit + "; it is reviewed once and no further review of this candidate follows."
			for _, want := range []string{"This is round 2 of 2 of this review.", wantSentence, "--- BEGIN CORRECTION DELTA DIFF ---", "--- END CORRECTION DELTA DIFF ---", "first structural fix", "second structural fix"} {
				if !strings.Contains(prompt, want) {
					t.Errorf("prompt missing %q: %s", want, prompt)
				}
			}
			for i, finding := range prior.FindingItems {
				if !strings.Contains(prompt, "- "+finding.ID+" | fixed by "+fixes[i]+" | "+finding.Text) {
					t.Errorf("prompt missing finding %s and its fixing commit", finding.ID)
				}
			}
			for _, excluded := range []string{"round two content excluded from correction", "the reviewer receives this changed line", "BEGIN CANDIDATE DIFF", "BEGIN ROUND 2 DELTA DIFF"} {
				if strings.Contains(prompt, excluded) {
					t.Errorf("prompt includes %q", excluded)
				}
			}
			persisted, err := readReviewRecordAtHead(fixture.artifactDir, fixture.repository, record.HeadCommit)
			if err != nil || !reflect.DeepEqual(persisted, record) {
				t.Fatalf("head record=%+v err=%v", persisted, err)
			}
			encoded, err := json.Marshal(record)
			if err != nil || !bytes.Contains(encoded, []byte(`"correction":true`)) {
				t.Fatalf("JSON=%s err=%v", encoded, err)
			}
			result, err := deliveryReviewResult(record, record.HeadCommit)
			if err != nil || result.Outcome != delivery.ReviewOutcomeReviewed {
				t.Fatalf("delivery=%+v err=%v", result, err)
			}
			// Repeating the same head must not spend another reviewer call either.
			code, _, stderr = runLineageReview(t, fixture)
			if code != exitOK || stderr != wantStderr || len(runner.runs) != calls+1 {
				t.Fatalf("same head code=%d calls=%d stderr=%q", code, len(runner.runs)-calls, stderr)
			}
			commitLineageFile(t, fixture, "later.txt", "after correction\n")
			code, closed, stderr := runLineageReview(t, fixture)
			if code != exitOK || stderr != "" || closed.Outcome != reviewOutcomeCeilingClosed || closed.Lineage.ReviewedHead != record.HeadCommit || len(runner.runs) != calls+1 {
				t.Fatalf("ceiling code=%d record=%+v calls=%d stderr=%q", code, closed, len(runner.runs)-calls, stderr)
			}
		})
	}
}

func TestCorrectionReviewFindingsReachTheCeiling(t *testing.T) {
	t.Parallel()
	fixture, runner, _, _ := fixedRoundTwoFixture(t, "Findings:\n- fix-one.txt:1 Failure: correction still broken")
	calls := len(runner.runs)
	code, correction, stderr := runLineageReview(t, fixture)
	if code != exitRunFailed || correction.Outcome != reviewOutcomeFindings || !correction.Lineage.Correction || len(runner.runs) != calls+1 {
		t.Fatalf("correction code=%d record=%+v stderr=%q", code, correction, stderr)
	}
	result, err := deliveryReviewResult(correction, correction.HeadCommit)
	if err != nil || result.Outcome != delivery.ReviewOutcomeFindings {
		t.Fatalf("delivery=%+v err=%v", result, err)
	}
	fix := commitLineageFile(t, fixture, "final.txt", "last fix\n")
	code, blocked, _ := runLineageReview(t, fixture)
	if code != exitPreflight || blocked.Outcome != reviewOutcomeBlocked || len(runner.runs) != calls+1 {
		t.Fatalf("undisposed ceiling code=%d record=%+v calls=%d", code, blocked, len(runner.runs)-calls)
	}
	disposeCorrectionFinding(t, "F1", "--fixed-by", fix)
	code, closed, stderr := runLineageReview(t, fixture)
	if code != exitOK || stderr != "" || closed.Outcome != reviewOutcomeCeilingClosed || closed.Lineage.ReviewedHead != correction.HeadCommit || len(runner.runs) != calls+1 {
		t.Fatalf("ceiling code=%d record=%+v calls=%d stderr=%q", code, closed, len(runner.runs)-calls, stderr)
	}
}

func TestADismissedRoundTwoFindingKeepsTheCeiling(t *testing.T) {
	t.Parallel()
	fixture, runner, prior := roundTwoLineageFixture(t)
	disposeLineageFinding(t, "--dismiss", "--evidence", "checked against the contract")
	commitLineageFile(t, fixture, "dismissed.txt", "candidate moved\n")
	assertLineageCeilingClosed(t, fixture, runner, prior)
}

func TestBlockedCorrectionReviewRepeatsItsDelta(t *testing.T) {
	t.Parallel()
	fixture, runner, prior, _ := fixedRoundTwoFixture(t, "invalid verdict", "No findings.")
	code, blocked, _ := runLineageReview(t, fixture)
	if code != exitPreflight || blocked.Outcome != reviewOutcomeBlocked || !blocked.Lineage.Correction {
		t.Fatalf("blocked code=%d record=%+v", code, blocked)
	}
	commitLineageFile(t, fixture, "retry.txt", "retry correction\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || stderr != "" || !record.Lineage.Correction || record.Lineage.PreviousHead != prior.HeadCommit || !reflect.DeepEqual(record.Lineage.PreviousDispositions, blocked.Lineage.PreviousDispositions) || len(runner.runs) != 4 {
		t.Fatalf("retry code=%d record=%+v calls=%d stderr=%q", code, record, len(runner.runs), stderr)
	}
	for _, want := range []string{"first structural fix", "second structural fix", "retry correction", "BEGIN CORRECTION DELTA DIFF"} {
		if !strings.Contains(runner.runs[3].Prompt, want) {
			t.Errorf("retry missing %q", want)
		}
	}
	persisted, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil || !persisted.Lineage.Correction {
		t.Fatalf("persisted=%+v err=%v", persisted, err)
	}
}

func TestCorrectionReviewRequiresEveryFixedDisposition(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		change   func(*reviewCorrectionFixture)
		wantCode int
	}{
		{"missing", func(f *reviewCorrectionFixture) { f.dispositions = f.dispositions[:1] }, exitPreflight},
		{"duplicate", func(f *reviewCorrectionFixture) { f.dispositions = append(f.dispositions, f.dispositions[0]) }, exitPreflight},
		{"wrong text", func(f *reviewCorrectionFixture) { f.dispositions[1].Text = "another finding" }, exitPreflight},
		{"fix before round two", func(f *reviewCorrectionFixture) { f.dispositions[1].FixedBy = f.record.BaseCommit }, exitOK},
		{"unavailable fixing commit", func(f *reviewCorrectionFixture) { f.dispositions[1].FixedBy = "missing-fixing-commit" }, exitPreflight},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			f.record.Lineage = &reviewLineage{Round: 2}
			for i := range f.dispositions {
				f.dispositions[i].Disposition = "fixed"
				f.dispositions[i].FixedBy = f.head
				f.dispositions[i].Evidence = ""
			}
			tc.change(f)
			f.persist(t)
			runner := newSessionReviewRunner("No findings.")
			writeReviewCommandConfig(t, f.review.repository, "codex", f.review.artifactDir)
			setCommandEnvironmentForTest(t, f.review.homeDir, f.review.repository)
			withAgentRunner(t, runner)
			f.review.baseCommit = f.record.BaseCommit
			code, record, stderr := runLineageReview(t, f.review)
			if code != tc.wantCode || len(runner.runs) != 0 || record.Lineage.Correction || record.Lineage.ReviewedHead != f.candidate {
				t.Fatalf("ceiling code=%d record=%+v calls=%d stderr=%q", code, record, len(runner.runs), stderr)
			}
		})
	}
}

func TestCorrectionReviewBlocksAnOversizedDeltaBeforeTheReviewer(t *testing.T) {
	t.Parallel()
	for _, provider := range []string{"codex", "claude"} {
		t.Run(provider, func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			f.record.Lineage = &reviewLineage{Round: 2}
			f.record.Provider = provider
			f.head = commitReviewCorrectionFile(t, f.review, "oversized.txt", strings.Repeat("large correction line\n", reviewDiffBound/20+1))
			for i := range f.dispositions {
				f.dispositions[i].Disposition = "fixed"
				f.dispositions[i].FixedBy = f.head
				f.dispositions[i].Evidence = ""
			}
			f.persist(t)
			runner := newSessionReviewRunner("No findings.")
			writeReviewCommandProfileConfig(t, f.review.repository, provider, f.review.artifactDir, provider, provider)
			setCommandEnvironmentForTest(t, f.review.homeDir, f.review.repository)
			withAgentRunner(t, runner)
			f.review.baseCommit = f.record.BaseCommit
			code, record, stderr := runLineageReview(t, f.review)
			if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || !record.Lineage.Correction || !strings.Contains(record.Reason, "review diff too large:") || record.DiffBytes <= reviewDiffBound || len(record.Segments) != 0 || runner.probeCalls != 0 || runner.prepareCalls != 0 || len(runner.runs) != 0 {
				t.Fatalf("bound code=%d record=%+v calls=%d stderr=%q", code, record, len(runner.runs), stderr)
			}
		})
	}
}
