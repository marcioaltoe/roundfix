package runevent

import (
	"encoding/json"
	"reflect"
	"strings"
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

func TestProjectVacuousEventJournaledBeforeTheFixReadsProbedCommands(t *testing.T) {
	t.Parallel()
	event := RunEvent{
		RunID:       "run_legacy_vacuous",
		Batch:       5,
		Source:      SourceDaemon,
		Kind:        KindDaemonVerification,
		ReviewIssue: "task_05",
		Time:        time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{
			"attempt":1,
			"phase":"failed",
			"task":"task_05",
			"classification":"verification_vacuous",
			"probed_commands":[
				{"command":"test -f first","verdict":"passed","probe_log_path":"/tmp/first.log"},
				{"command":"test -f missing","verdict":"failed","probe_log_path":"/tmp/missing.log"},
				{"command":"test -f second","verdict":"passed","probe_log_path":"/tmp/second.log"}
			]
		}`),
	}

	record, ok, err := ProjectStreamEvent(4, event, AllStreamCategories())
	if err != nil {
		t.Fatalf("project legacy vacuous event: %v", err)
	}
	if !ok {
		t.Fatal("legacy vacuous event did not project")
	}
	want := []string{"test -f first", "test -f second"}
	if !reflect.DeepEqual(record.Commands, want) {
		t.Fatalf("commands = %q, want passed commands %q", record.Commands, want)
	}
}

func TestProjectVacuousEventWithoutCommandsIsMalformed(t *testing.T) {
	t.Parallel()
	event := RunEvent{
		RunID:       "run_malformed_vacuous",
		Batch:       5,
		Source:      SourceDaemon,
		Kind:        KindDaemonVerification,
		ReviewIssue: "task_05",
		Time:        time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC),
		Payload:     []byte(`{"attempt":1,"phase":"failed","task":"task_05","classification":"verification_vacuous"}`),
	}

	_, _, err := ProjectStreamEvent(4, event, AllStreamCategories())
	if err == nil {
		t.Fatal("expected vacuous event without commands to fail projection")
	}
	if !strings.Contains(err.Error(), `missing payload field "commands"`) {
		t.Fatalf("projection error = %q, want missing commands", err)
	}
}
