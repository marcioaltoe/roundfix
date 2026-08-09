package runevent

import (
	"encoding/json"
	"testing"
	"time"
)

// Suite: applied Agent Selection receipt stream projection
// Invariant: an observed runtime-deferred effort remains visible in the agent-selection stream category.
// Boundary IN: Run Event receipt payload and public stream projection.
// Boundary OUT: Agent Session command execution and Run Event Journal persistence.

func TestProjectStreamEventProjectsSelectionReceiptInExistingCategory(t *testing.T) {
	t.Parallel()

	payload, err := json.Marshal(SelectionReceiptPayload{
		Event:                    SelectionReceiptEventApplied,
		Session:                  "roundfix-run-1-task_04",
		Runtime:                  "opencode",
		Model:                    "openrouter/deepseek/deepseek-v4-pro",
		RequestedReasoningEffort: "xhigh",
		ReasoningEffort:          "xhigh",
		Status:                   SelectionReceiptStatusApplied,
	})
	if err != nil {
		t.Fatalf("marshal selection receipt: %v", err)
	}
	event := RunEvent{
		RunID:   "run-1",
		Source:  SourceAgent,
		Kind:    KindAgentSelectionReceipt,
		Time:    time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC),
		Payload: payload,
	}

	record, ok, err := ProjectStreamEvent(7, event, nil)
	if err != nil {
		t.Fatalf("project selection receipt: %v", err)
	}
	if !ok {
		t.Fatal("expected selection receipt to project")
	}
	if record.Category != StreamCategorySelection || record.ReasoningEffort != "xhigh" {
		t.Fatalf("selection receipt stream record = %#v", record)
	}
	if record.ScopeKind != "agent_session" || record.ScopeID != "roundfix-run-1-task_04" || record.Status != SelectionReceiptStatusApplied {
		t.Fatalf("selection receipt scope/status = %#v", record)
	}
}
