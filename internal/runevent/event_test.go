package runevent

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestSpecRunKindsUseDaemonNamespace(t *testing.T) {
	if KindDaemonTask != "daemon.task" {
		t.Fatalf("expected daemon.task kind, got %q", KindDaemonTask)
	}
	if KindDaemonQA != "daemon.qa" {
		t.Fatalf("expected daemon.qa kind, got %q", KindDaemonQA)
	}
}

func TestVerificationEventVocabulary(t *testing.T) {
	if VerificationPhaseStarted != "started" {
		t.Fatalf("expected started phase, got %q", VerificationPhaseStarted)
	}
	if VerificationPhaseCommandPassed != "command-passed" {
		t.Fatalf("expected command-passed phase, got %q", VerificationPhaseCommandPassed)
	}
	if VerificationPhaseFailed != "failed" {
		t.Fatalf("expected failed phase, got %q", VerificationPhaseFailed)
	}
	if VerificationPhaseVerdict != "verdict" {
		t.Fatalf("expected verdict phase, got %q", VerificationPhaseVerdict)
	}
	if VerificationVerdictPassed != "passed" {
		t.Fatalf("expected passed verdict, got %q", VerificationVerdictPassed)
	}
	if VerificationVerdictFailed != "failed" {
		t.Fatalf("expected failed verdict, got %q", VerificationVerdictFailed)
	}
}

func TestIsDaemonKindCoversSpecRunKindsAndSkipsUnknown(t *testing.T) {
	tests := []struct {
		name     string
		kind     Kind
		expected bool
	}{
		{name: "task kind is daemon vocabulary", kind: KindDaemonTask, expected: true},
		{name: "qa kind is daemon vocabulary", kind: KindDaemonQA, expected: true},
		{name: "selection attempt kind is daemon vocabulary", kind: KindDaemonAgentSelectionAttempt, expected: true},
		{name: "selection active kind is daemon vocabulary", kind: KindDaemonAgentSelectionActive, expected: true},
		{name: "selection fallback kind is daemon vocabulary", kind: KindDaemonAgentSelectionFallback, expected: true},
		{name: "selection exhausted kind is daemon vocabulary", kind: KindDaemonAgentSelectionExhausted, expected: true},
		{name: "selection closed kind is daemon vocabulary", kind: KindDaemonAgentSelectionClosed, expected: true},
		{name: "agent kind is not daemon vocabulary", kind: KindAgentMessage, expected: false},
		{name: "unknown daemon-prefixed kind stays skippable", kind: Kind("daemon.unknown"), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDaemonKind(tt.kind); got != tt.expected {
				t.Fatalf("IsDaemonKind(%q) = %v, want %v", tt.kind, got, tt.expected)
			}
		})
	}
}

func TestSourceFilterSinkDropsOnlyConfiguredSource(t *testing.T) {
	next := &recordingSink{}
	sink := NewSourceFilterSink(next, SourceAgent)
	events := []RunEvent{
		{Source: SourceDaemon, Kind: KindDaemonSelection, Summary: "selected"},
		{Source: SourceAgent, Kind: KindAgentRaw, Summary: "agent output"},
		{Source: SourceVerification, Kind: KindDaemonVerification, Summary: "verified"},
		{Source: SourceAgent, Kind: KindAgentStatus, Summary: "agent status"},
		{Source: SourceGit, Kind: KindDaemonCommit, Summary: "committed"},
	}

	for _, event := range events {
		if err := sink.Publish(context.Background(), event); err != nil {
			t.Fatalf("publish %s: %v", event.Summary, err)
		}
	}

	if len(next.events) != 3 {
		t.Fatalf("expected three forwarded events, got %#v", next.events)
	}
	for index, want := range []string{"selected", "verified", "committed"} {
		if next.events[index].Summary != want {
			t.Fatalf("expected forwarded event %d to be %q, got %#v", index, want, next.events)
		}
	}
}

func TestSourceFilterSinkPropagatesWrappedSinkErrors(t *testing.T) {
	next := &recordingSink{err: errors.New("writer failed")}
	sink := NewSourceFilterSink(next, SourceAgent)

	droppedErr := sink.Publish(context.Background(), RunEvent{Source: SourceAgent, Kind: KindAgentRaw})
	forwardedErr := sink.Publish(context.Background(), RunEvent{Source: SourceDaemon, Kind: KindDaemonOutcome})

	if droppedErr != nil {
		t.Fatalf("expected dropped event to skip wrapped sink error, got %v", droppedErr)
	}
	if !errors.Is(forwardedErr, next.err) {
		t.Fatalf("expected forwarded event to propagate wrapped sink error, got %v", forwardedErr)
	}
	if len(next.events) != 1 || next.events[0].Source != SourceDaemon {
		t.Fatalf("expected only the daemon event forwarded, got %#v", next.events)
	}
}

func TestProjectStreamEventCoversStableCategoriesAndRedactsPayload(t *testing.T) {
	at := time.Date(2026, 7, 10, 12, 0, 0, 123, time.UTC)
	events := []RunEvent{
		{
			RunID:       "run_123",
			Batch:       4,
			Source:      SourceDaemon,
			Kind:        KindDaemonTask,
			ReviewIssue: "task_03",
			Summary:     "Task task_03 settled completed.",
			Time:        at,
			Payload:     []byte(`{"task":"task_03","phase":"settled","status":"completed"}`),
		},
		{
			RunID:   "run_123",
			Batch:   4,
			Source:  SourceDaemon,
			Kind:    KindDaemonBatch,
			Summary: "Batch 004 completed.",
			Time:    at,
			Payload: []byte(`{"phase":"completed","batch":4}`),
		},
		{
			RunID:       "run_123",
			Batch:       4,
			Source:      SourceDaemon,
			Kind:        KindDaemonVerification,
			ReviewIssue: "task_03",
			Summary:     "Verification attempt 1 failed: make verify",
			Time:        at,
			Payload:     []byte(`{"attempt":1,"phase":"verdict","verdict":"failed","command":"make verify","diagnostic_path":"/tmp/log"}`),
		},
		{
			RunID:   "run_123",
			Source:  SourceDaemon,
			Kind:    KindDaemonOutcome,
			Summary: "Run reached Clean.",
			Time:    at,
			Payload: []byte(`{"state":"Clean","remaining":0}`),
		},
		{
			RunID:   "run_123",
			Source:  SourceAgent,
			Kind:    KindAgentRaw,
			Summary: "raw agent payload",
			Time:    at,
			Payload: []byte(`{"raw":"agent output"}`),
		},
	}
	filter := AllStreamCategories()
	categories := []string{}
	for index, event := range events {
		record, ok, err := ProjectStreamEvent(int64(index+1), event, filter)
		if err != nil {
			t.Fatalf("project event %d: %v", index, err)
		}
		if !ok {
			continue
		}
		categories = append(categories, string(record.Category))
		raw, err := json.Marshal(record)
		if err != nil {
			t.Fatalf("marshal projected record: %v", err)
		}
		encoded := string(raw)
		for _, forbidden := range []string{"daemon.", "agent output", "make verify", "diagnostic_path", "/tmp/log"} {
			if strings.Contains(encoded, forbidden) {
				t.Fatalf("expected projection to omit %q, got %s", forbidden, encoded)
			}
		}
		if record.Schema != StreamSchema || record.RunID != "run_123" || record.Cursor == 0 || record.Time == "" {
			t.Fatalf("expected stable envelope fields, got %#v", record)
		}
	}
	if got := strings.Join(categories, "|"); got != "task-status|batch|verification|outcome" {
		t.Fatalf("expected stable categories only, got %s", got)
	}
}

func TestProjectStreamEventNormalizesLegacyVerificationEvents(t *testing.T) {
	tests := []struct {
		name    string
		summary string
		payload string
		phase   string
	}{
		{
			name:    "started",
			summary: "Verification started: make verify",
			payload: `{"phase":"started","command":"make verify"}`,
			phase:   "started",
		},
		{
			name:    "failed",
			summary: "Verification failed: make verify",
			payload: `{"phase":"failed","command":"make verify","error":"exit status 1"}`,
			phase:   "failed",
		},
		{
			name:    "passed",
			summary: "Verification command passed: make verify",
			payload: `{"phase":"passed","command":"make verify"}`,
			phase:   "passed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, ok, err := ProjectStreamEvent(1, RunEvent{
				RunID:   "run_legacy",
				Batch:   1,
				Source:  SourceDaemon,
				Kind:    KindDaemonVerification,
				Summary: tt.summary,
				Time:    time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
				Payload: []byte(tt.payload),
			}, AllStreamCategories())
			if err != nil {
				t.Fatalf("ProjectStreamEvent: %v", err)
			}
			if !ok {
				t.Fatal("expected legacy verification event to be projected")
			}
			if record.Attempt != 1 || record.Phase != tt.phase {
				t.Fatalf("record = %#v, want attempt 1 phase %q", record, tt.phase)
			}
		})
	}
}

func TestProjectStreamEventRejectsNewVerificationPayloadWithoutAttempt(t *testing.T) {
	_, _, err := ProjectStreamEvent(1, RunEvent{
		RunID:   "run_123",
		Batch:   1,
		Source:  SourceDaemon,
		Kind:    KindDaemonVerification,
		Summary: "Verification attempt 1 for Batch 001 started: make verify",
		Time:    time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"phase":"started","command":"make verify"}`),
	}, AllStreamCategories())
	if err == nil {
		t.Fatal("expected new verification payload without attempt to fail")
	}
}

func TestProjectStreamEventRejectsMalformedRelevantDaemonPayload(t *testing.T) {
	_, _, err := ProjectStreamEvent(1, RunEvent{
		RunID:   "run_123",
		Batch:   1,
		Source:  SourceDaemon,
		Kind:    KindDaemonBatch,
		Summary: "Batch 001 completed.",
		Time:    time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"batch":1}`),
	}, AllStreamCategories())
	if err == nil {
		t.Fatal("expected missing phase in relevant daemon payload to fail")
	}
}

func TestParseStreamCategoryFilterValidatesAndDeduplicates(t *testing.T) {
	filter, err := ParseStreamCategoryFilter("verification,outcome,verification")
	if err != nil {
		t.Fatalf("parse filter: %v", err)
	}
	if !filter.Includes(StreamCategoryVerification) || !filter.Includes(StreamCategoryOutcome) {
		t.Fatalf("expected requested categories included, got %#v", filter)
	}
	if filter.Includes(StreamCategoryBatch) || filter.Includes(StreamCategoryTaskStatus) {
		t.Fatalf("expected unrequested categories excluded, got %#v", filter)
	}
	for _, raw := range []string{"", "verification,,outcome", "raw"} {
		if _, err := ParseStreamCategoryFilter(raw); err == nil {
			t.Fatalf("expected invalid filter %q to fail", raw)
		}
	}
}

func TestAgentSelectionEventProjectsStablePayloadAndKeepsReasoningEmpty(t *testing.T) {
	record, ok, err := ProjectStreamEvent(7, RunEvent{
		RunID:   "run_123",
		Batch:   2,
		Source:  SourceDaemon,
		Kind:    KindDaemonAgentSelectionAttempt,
		Summary: "Agent Selection attempt 1 for task task_01 (backend) attempting.",
		Time:    time.Date(2026, 7, 17, 8, 1, 0, 0, time.UTC),
		Payload: []byte(`{
			"event":"agent_selection_attempt",
			"scope_kind":"task",
			"scope_id":"task_01",
			"scope_identity":"task:task_01",
			"category":"backend",
			"profile_source":"project",
			"attempt":1,
			"selection_role":"preferred",
			"fallback_index":0,
			"runtime":"codex",
			"model":"gpt-5.6-sol",
			"reasoning_effort":"",
			"status":"attempting",
			"reason_code":"",
			"reason":""
		}`),
	}, AllStreamCategories())
	if err != nil {
		t.Fatalf("project selection event: %v", err)
	}
	if !ok {
		t.Fatal("expected selection event to project")
	}
	if record.Category != StreamCategorySelection ||
		record.ScopeKind != "task" ||
		record.ScopeID != "task_01" ||
		record.WorkCategory != "backend" ||
		record.Attempt != 1 ||
		record.SelectionRole != "preferred" ||
		record.Runtime != "codex" ||
		record.Model != "gpt-5.6-sol" ||
		record.ReasoningEffort != "" ||
		record.Status != "attempting" {
		t.Fatalf("expected projected selection fields, got %#v", record)
	}
	raw, err := json.Marshal(record)
	if err != nil {
		t.Fatalf("marshal projected selection record: %v", err)
	}
	for _, forbidden := range []string{"prompt", "transcript", "credential", "token", "cookie", "secret"} {
		if strings.Contains(strings.ToLower(string(raw)), forbidden) {
			t.Fatalf("expected projected selection record to omit %q, got %s", forbidden, raw)
		}
	}
}

func TestAgentSelectionEventProjectsFallbackNotificationPayload(t *testing.T) {
	record, ok, err := ProjectStreamEvent(8, RunEvent{
		RunID:   "run_123",
		Source:  SourceDaemon,
		Kind:    KindDaemonAgentSelectionFallback,
		Summary: "Agent Selection fallback for task task_01 (backend).",
		Time:    time.Date(2026, 7, 17, 8, 2, 0, 0, time.UTC),
		Payload: []byte(`{
			"event":"agent_selection_fallback",
			"category":"backend",
			"scope_kind":"task",
			"scope_id":"task_01",
			"scope_identity":"task:task_01",
			"failed_selection":{"runtime":"codex","model":"gpt-5.6-sol","reasoning_effort":""},
			"next_selection":{"runtime":"claude","model":"claude-fable-5","reasoning_effort":"xhigh"},
			"fallback_index":1,
			"reason_code":"model_not_advertised",
			"reason":"model was not advertised",
			"automatic":true
		}`),
	}, AllStreamCategories())
	if err != nil {
		t.Fatalf("project fallback notification: %v", err)
	}
	if !ok {
		t.Fatal("expected fallback notification to project")
	}
	if record.Category != StreamCategorySelection ||
		record.ScopeKind != "task" ||
		record.ScopeID != "task_01" ||
		record.WorkCategory != "backend" ||
		record.Status != "failed" ||
		record.FallbackIndex != 1 ||
		record.Runtime != "codex" ||
		record.Model != "gpt-5.6-sol" ||
		record.ReasoningEffort != "" ||
		record.ReasonCode != "model_not_advertised" {
		t.Fatalf("expected projected fallback fields, got %#v", record)
	}
}
