// Suite: pre-PR review final-message verdicts
// Invariant: only the final non-blank Agent message is classified, while the answer file keeps the full transcript.
// Boundary IN: review command classification and answer persistence through its fake Agent runner.
// Boundary OUT: ACP stream parsing, provider selection, and real reviewers.

package cli

import (
	"os"
	"testing"

	"roundfix/internal/agent"
)

func TestReviewClassifiesTheFinalMessageAfterProgressText(t *testing.T) {
	const progress = "I am checking the candidate."
	const answer = "Findings:\n- review.txt:2: preserve the final message"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Message:    progress + "\n\n" + answer,
				Messages:   []string{progress, answer},
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed {
		t.Fatalf("review exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr)
	}
	if record.Outcome != reviewOutcomeFindings || record.Findings != "- review.txt:2: preserve the final message" {
		t.Fatalf("review record = %+v, want one findings verdict", record)
	}
	if len(record.FindingItems) != 1 || record.FindingItems[0].ID != "F1" || record.FindingItems[0].Text != "review.txt:2: preserve the final message" {
		t.Fatalf("finding items = %+v, want F1", record.FindingItems)
	}
}

func TestReviewAcceptsNoFindingsAsTheFinalMessage(t *testing.T) {
	const progress = "I am checking the candidate."
	const answer = "No findings."
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Message:    progress + "\n\n" + answer,
				Messages:   []string{progress, answer},
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review exit=%d record=%+v stderr=%q, want reviewed", code, record, stderr)
	}
}

func TestReviewStillBlocksBothVerdictsInTheFinalMessage(t *testing.T) {
	const progress = "I am checking the candidate."
	const answer = "No findings.\nFindings:\n- internal/cli/review.go:42: conflicting verdict"
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Message:    progress + "\n\n" + answer,
				Messages:   []string{progress, answer},
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	code, record, stderr := fixture.run(t)

	assertBlockedReviewCommand(t, code, record, stderr, "both no-findings and findings verdicts are present")
}

func TestReviewAnswerFileKeepsEveryMessage(t *testing.T) {
	const progress = "I am checking the candidate."
	const answer = "No findings."
	const transcript = progress + "\n\n" + answer
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{
				Message:    transcript,
				Messages:   []string{progress, answer},
				StopReason: "end_turn",
			},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	_, record, _ := fixture.run(t)

	written, err := os.ReadFile(record.AnswerPath)
	if err != nil {
		t.Fatalf("read review answer: %v", err)
	}
	if got := string(written); got != transcript {
		t.Fatalf("review answer = %q, want %q", got, transcript)
	}
}
