// Suite: pre-PR review records.
// Invariant: each written record identifies one candidate and one valid outcome with its required detail.
// Boundary IN: review-record construction, validation, and JSON encoding.
// Boundary OUT: policy resolution, reviewer execution, command routing, and record storage.
package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
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
