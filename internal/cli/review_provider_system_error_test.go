package cli

import (
	"testing"

	"roundfix/internal/agent"
)

func TestReviewRecordsAProviderSystemErrorAsARuntimeFailure(t *testing.T) {
	t.Parallel()
	const detail = "Selected model is at capacity. Please try a different model."
	failure := &agent.BatchFailureError{
		Reason: agent.ProviderSystemErrorReason,
		Protocol: &agent.ProtocolFailure{
			Step: "session/prompt", PromptSent: true, ProviderSystemError: true, Detail: detail,
		},
	}
	for _, tc := range []struct {
		name, answer string
	}{
		{"capacity text", detail},
		{"apparent clean verdict", "No findings."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			record := reviewRecord{Outcome: reviewOutcomeBlocked}
			result := agent.ExecuteResult{Message: tc.answer, Messages: []string{tc.answer}, StopReason: "end_turn"}
			got, code := classifyReviewCommandResult(record, result, failure)
			want := "review runtime failure: Agent Batch failed: provider system error at session/prompt: " + detail
			if code != exitPreflight || got.Outcome != reviewOutcomeBlocked || got.Reason != want || got.FailedStep != "session/prompt" {
				t.Fatalf("classification = %+v, code = %d, want runtime failure %q", got, code, want)
			}
			if got.Findings != "" || len(got.FindingItems) != 0 {
				t.Fatalf("provider failure parsed as review output: %+v", got)
			}
		})
	}
}
