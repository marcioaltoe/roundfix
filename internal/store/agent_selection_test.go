package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

func TestSchema9CreatesAgentSelectionTable(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != schemaVersion {
		t.Fatalf("expected schema version %d, got %d", schemaVersion, version)
	}
	assertAgentSelectionTableExists(t, ctx, runStore)
	assertAgentSelectionSchemaPrivacy(t, ctx, runStore)
}

func TestSchema8To9MigratesRunsEventsAndSelectionTable(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV8Fixture(t, ctx, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != schemaVersion {
		t.Fatalf("expected schema version %d, got %d", schemaVersion, version)
	}
	run, ok, err := runStore.Run(ctx, "run_v8_active")
	if err != nil || !ok {
		t.Fatalf("expected schema 8 Run to survive, ok=%v err=%v", ok, err)
	}
	if run.Agent != "codex" || run.Model != "gpt-5.6-sol" || run.ReasoningEffort != "" {
		t.Fatalf("expected compatibility summary to survive, got %#v", run)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read migrated Run Events: %v", err)
	}
	if len(events) != 1 || events[0].Event.Summary != "schema 8 event" {
		t.Fatalf("expected schema 8 Run Event to survive, got %#v", events)
	}
	assertAgentSelectionTableExists(t, ctx, runStore)
}

func TestAgentSelectionAttemptsRoundTripTaskQAReviewHistories(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	runStore.now = func() time.Time { return time.Date(2026, 7, 17, 8, 0, 0, 0, time.UTC) }

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	requests := []AgentSelectionAttemptRequest{
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
			ReasoningEffort: "", Status: AgentSelectionStatusActive,
			Time: time.Date(2026, 7, 17, 8, 1, 0, 0, time.UTC),
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
			Category: "backend", ProfileSource: "project", Attempt: 2,
			SelectionRole: AgentSelectionRoleFallback, FallbackIndex: 1,
			Runtime: "claude", Model: "claude-fable-5", ReasoningEffort: "xhigh",
			Status: AgentSelectionStatusFailed, ReasonCode: "model_not_advertised",
			Reason: "model not advertised by runtime",
			Time:   time.Date(2026, 7, 17, 8, 2, 0, 0, time.UTC),
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeQA, ScopeID: "qa:default",
			Category: "qa", ProfileSource: "built-in", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-terra",
			ReasoningEffort: "high", Status: AgentSelectionStatusClosed,
			Time: time.Date(2026, 7, 17, 8, 3, 0, 0, time.UTC),
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeReview, ScopeID: "review:run",
			Category: "review", ProfileSource: "user", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "opencode", Model: "custom-review-model",
			ReasoningEffort: "model-managed", Status: AgentSelectionStatusAttempting,
			Time: time.Date(2026, 7, 17, 8, 4, 0, 0, time.UTC),
		},
	}

	for _, req := range requests {
		if _, err := runStore.AppendAgentSelectionAttempt(ctx, req); err != nil {
			t.Fatalf("append %s %s attempt %d: %v", req.ScopeKind, req.ScopeID, req.Attempt, err)
		}
	}

	taskAttempts, err := runStore.AgentSelectionAttemptsForScope(ctx, run.ID, AgentSelectionScopeTask, "task_01")
	if err != nil {
		t.Fatalf("read task attempts: %v", err)
	}
	if len(taskAttempts) != 2 {
		t.Fatalf("expected two task attempts, got %#v", taskAttempts)
	}
	assertSelectionAttemptMatches(t, taskAttempts[0], requests[0])
	assertSelectionAttemptMatches(t, taskAttempts[1], requests[1])
	if taskAttempts[0].ReasoningEffort != "" {
		t.Fatalf("expected explicit empty reasoning to round-trip, got %q", taskAttempts[0].ReasoningEffort)
	}

	allAttempts, err := runStore.AgentSelectionAttempts(ctx, run.ID)
	if err != nil {
		t.Fatalf("read all attempts: %v", err)
	}
	if len(allAttempts) != len(requests) {
		t.Fatalf("expected %d attempts, got %#v", len(requests), allAttempts)
	}

	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read selection Run Events: %v", err)
	}
	if got, want := selectionEventKinds(events), []runevent.Kind{
		runevent.KindDaemonAgentSelectionActive,
		runevent.KindDaemonAgentSelectionFallback,
		runevent.KindDaemonAgentSelectionClosed,
		runevent.KindDaemonAgentSelectionAttempt,
	}; strings.Join(kindStrings(got), "|") != strings.Join(kindStrings(want), "|") {
		t.Fatalf("expected event kinds %v, got %v", want, got)
	}
	var payload map[string]any
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatalf("decode selection event payload: %v", err)
	}
	if payload["scope_kind"] != "task" || payload["scope_id"] != "task_01" ||
		payload["category"] != "backend" || payload["runtime"] != "codex" ||
		payload["model"] != "gpt-5.6-sol" || payload["reasoning_effort"] != "" {
		t.Fatalf("expected event payload to match persisted attempt, got %#v", payload)
	}
	assertNoSensitiveSelectionPayload(t, string(events[0].Event.Payload))
}

func TestSelectionAttemptOrderingRejectsInvalidOrOutOfOrderAppendsWithoutPartialRows(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	valid := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		Status: AgentSelectionStatusAttempting,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, valid); err != nil {
		t.Fatalf("append valid first attempt: %v", err)
	}

	tests := []struct {
		name string
		req  AgentSelectionAttemptRequest
	}{
		{name: "unknown scope", req: withAttemptOverride(valid, func(req *AgentSelectionAttemptRequest) { req.ScopeKind = "deploy"; req.ScopeID = "deploy-1" })},
		{name: "unknown role", req: withAttemptOverride(valid, func(req *AgentSelectionAttemptRequest) { req.Attempt = 2; req.SelectionRole = "backup" })},
		{name: "unknown status", req: withAttemptOverride(valid, func(req *AgentSelectionAttemptRequest) { req.Attempt = 2; req.Status = "ready" })},
		{name: "duplicate attempt", req: valid},
		{name: "gap attempt", req: withAttemptOverride(valid, func(req *AgentSelectionAttemptRequest) { req.Attempt = 3 })},
		{name: "missing Run", req: withAttemptOverride(valid, func(req *AgentSelectionAttemptRequest) { req.RunID = "run_missing"; req.Attempt = 1 })},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeRows := countAgentSelectionAttempts(t, ctx, runStore, run.ID)
			beforeEvents := countRunEvents(t, ctx, runStore, run.ID)
			if _, err := runStore.AppendAgentSelectionAttempt(ctx, tt.req); err == nil {
				t.Fatal("expected append to fail")
			}
			if got := countAgentSelectionAttempts(t, ctx, runStore, run.ID); got != beforeRows {
				t.Fatalf("expected row count to remain %d, got %d", beforeRows, got)
			}
			if got := countRunEvents(t, ctx, runStore, run.ID); got != beforeEvents {
				t.Fatalf("expected event count to remain %d, got %d", beforeEvents, got)
			}
		})
	}
}

func TestAgentSelectionAttemptLifecycleUpdatesSameAttempt(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	base := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		Status: AgentSelectionStatusAttempting,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, base); err != nil {
		t.Fatalf("append attempting: %v", err)
	}
	active := withAttemptOverride(base, func(req *AgentSelectionAttemptRequest) {
		req.Status = AgentSelectionStatusActive
	})
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, active); err != nil {
		t.Fatalf("append active update: %v", err)
	}
	closed := withAttemptOverride(base, func(req *AgentSelectionAttemptRequest) {
		req.Status = AgentSelectionStatusClosed
	})
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, closed); err != nil {
		t.Fatalf("append closed update: %v", err)
	}
	attempts, err := runStore.AgentSelectionAttemptsForScope(ctx, run.ID, AgentSelectionScopeTask, "task_01")
	if err != nil {
		t.Fatalf("read attempts: %v", err)
	}
	if len(attempts) != 1 || attempts[0].Status != AgentSelectionStatusClosed {
		t.Fatalf("expected one closed attempt row, got %#v", attempts)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if got, want := selectionEventKinds(events), []runevent.Kind{
		runevent.KindDaemonAgentSelectionAttempt,
		runevent.KindDaemonAgentSelectionActive,
		runevent.KindDaemonAgentSelectionClosed,
	}; strings.Join(kindStrings(got), "|") != strings.Join(kindStrings(want), "|") {
		t.Fatalf("expected lifecycle event kinds %v, got %v", want, got)
	}

	failed := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeQA, ScopeID: "qa:default",
		Category: "qa", ProfileSource: "built-in", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-terra",
		Status: AgentSelectionStatusAttempting,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, failed); err != nil {
		t.Fatalf("append qa attempting: %v", err)
	}
	failed = withAttemptOverride(failed, func(req *AgentSelectionAttemptRequest) {
		req.Status = AgentSelectionStatusFailed
		req.ReasonCode = "runtime_unavailable"
		req.Reason = "runtime failed before start"
	})
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, failed); err != nil {
		t.Fatalf("append failed update: %v", err)
	}
	qaAttempts, err := runStore.AgentSelectionAttemptsForScope(ctx, run.ID, AgentSelectionScopeQA, "qa:default")
	if err != nil {
		t.Fatalf("read QA attempts: %v", err)
	}
	if len(qaAttempts) != 1 {
		t.Fatalf("expected one persisted failed QA attempt, got %#v", qaAttempts)
	}
	if qaAttempts[0].Status != AgentSelectionStatusFailed ||
		qaAttempts[0].ReasonCode != "runtime_unavailable" ||
		qaAttempts[0].Reason != "runtime failed before start" {
		t.Fatalf("expected persisted failed QA attempt with reason, got %#v", qaAttempts[0])
	}
	events, err = runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events after failed QA attempt: %v", err)
	}
	if len(events) == 0 || events[len(events)-1].Event.Kind != runevent.KindDaemonAgentSelectionFallback {
		t.Fatalf("expected latest Run Event to be fallback, got %#v", events)
	}
}

func TestAgentSelectionActiveScopesReturnsLatestLifecycleInStableOrder(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	requests := []AgentSelectionAttemptRequest{
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeReview, ScopeID: "batch-002",
			Category: "review", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "opencode", Model: "review-model",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_02",
			Category: "frontend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "claude", Model: "task-model",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeQA, ScopeID: "qa",
			Category: "qa", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRoleFallback, FallbackIndex: 1,
			Runtime: "codex", Model: "qa-model", ReasoningEffort: "high",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "task-model",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_closed",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "closed-model",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_closed",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "closed-model",
			Status: AgentSelectionStatusClosed,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_superseded",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "old-model",
			Status: AgentSelectionStatusActive,
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_superseded",
			Category: "backend", ProfileSource: "project", Attempt: 2,
			SelectionRole: AgentSelectionRoleFallback, FallbackIndex: 1,
			Runtime: "codex", Model: "replacement-model",
			Status: AgentSelectionStatusFailed, ReasonCode: "runtime_unavailable",
		},
		{
			RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_failed",
			Category: "backend", ProfileSource: "project", Attempt: 1,
			SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "failed-model",
			Status: AgentSelectionStatusFailed, ReasonCode: "runtime_unavailable",
		},
	}
	for _, req := range requests {
		if _, err := runStore.AppendAgentSelectionAttempt(ctx, req); err != nil {
			t.Fatalf("append %s %s attempt %d status %s: %v", req.ScopeKind, req.ScopeID, req.Attempt, req.Status, err)
		}
	}

	active, err := runStore.ActiveAgentSelectionScopes(ctx, run.ID)
	if err != nil {
		t.Fatalf("list active Agent Selection scopes: %v", err)
	}
	gotScopes := make([]string, 0, len(active))
	for _, attempt := range active {
		gotScopes = append(gotScopes, string(attempt.ScopeKind)+":"+attempt.ScopeID)
	}
	wantScopes := []string{"task:task_01", "task:task_02", "qa:qa", "review:batch-002"}
	if !reflect.DeepEqual(gotScopes, wantScopes) {
		t.Fatalf("active Agent Selection scope order mismatch\nwant: %#v\ngot:  %#v", wantScopes, gotScopes)
	}
	if active[2].SelectionRole != AgentSelectionRoleFallback ||
		active[2].FallbackIndex != 1 ||
		active[2].Runtime != "codex" ||
		active[2].Model != "qa-model" ||
		active[2].ReasoningEffort != "high" {
		t.Fatalf("expected active QA fallback selection fields, got %#v", active[2])
	}
}

func TestAgentSelectionAttemptLifecycleRejectsImmutableFieldChanges(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	base := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		Status: AgentSelectionStatusAttempting,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, base); err != nil {
		t.Fatalf("append attempting: %v", err)
	}
	changedModel := withAttemptOverride(base, func(req *AgentSelectionAttemptRequest) {
		req.Status = AgentSelectionStatusActive
		req.Model = "different-model"
	})
	beforeRows := countAgentSelectionAttempts(t, ctx, runStore, run.ID)
	beforeEvents := countRunEvents(t, ctx, runStore, run.ID)

	if _, err := runStore.AppendAgentSelectionAttempt(ctx, changedModel); err == nil || !strings.Contains(err.Error(), "selection fields are immutable") {
		t.Fatalf("expected immutable field rejection, got %v", err)
	}
	if got := countAgentSelectionAttempts(t, ctx, runStore, run.ID); got != beforeRows {
		t.Fatalf("expected row count to remain %d, got %d", beforeRows, got)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != beforeEvents {
		t.Fatalf("expected event count to remain %d, got %d", beforeEvents, got)
	}
}

func TestAgentSelectionAttemptsStoreNoSensitiveFields(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	assertAgentSelectionSchemaPrivacy(t, ctx, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	_, err = runStore.AppendAgentSelectionAttempt(ctx, AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		Status: AgentSelectionStatusAttempting,
	})
	if err != nil {
		t.Fatalf("append selection attempt: %v", err)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one selection event, got %#v", events)
	}
	assertNoSensitiveSelectionPayload(t, string(events[0].Event.Payload))
}

func TestAgentSelectionAttemptsAppendExhaustedEventFromPersistedHistory(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	first := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		Status: AgentSelectionStatusFailed, ReasonCode: "runtime_unavailable", Reason: "runtime failed before start",
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, first); err != nil {
		t.Fatalf("append first attempt: %v", err)
	}
	beforeRows := countAgentSelectionAttempts(t, ctx, runStore, run.ID)
	cursor, err := runStore.AppendAgentSelectionExhausted(ctx, AgentSelectionExhaustedRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01",
		Category: "backend", RecoveryAction: "run profiles validate",
		Time: time.Date(2026, 7, 17, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("append exhausted event: %v", err)
	}
	if cursor != 2 {
		t.Fatalf("expected exhausted event cursor 2, got %d", cursor)
	}
	if got := countAgentSelectionAttempts(t, ctx, runStore, run.ID); got != beforeRows {
		t.Fatalf("expected exhausted event not to add attempt rows, got %d want %d", got, beforeRows)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if events[1].Event.Kind != runevent.KindDaemonAgentSelectionExhausted {
		t.Fatalf("expected exhausted event, got %#v", events[1].Event)
	}
	payload := string(events[1].Event.Payload)
	for _, want := range []string{"agent_selection_exhausted", "task_01", "gpt-5.6-sol", "runtime_unavailable", "run profiles validate"} {
		if !strings.Contains(payload, want) {
			t.Fatalf("expected exhausted payload to contain %q, got %s", want, payload)
		}
	}
	assertNoSensitiveSelectionPayload(t, payload)

	if _, err := runStore.AppendAgentSelectionExhausted(ctx, AgentSelectionExhaustedRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeQA, ScopeID: "qa:missing", Category: "qa",
	}); err == nil {
		t.Fatal("expected exhausted event without persisted attempts to fail")
	}
	if _, err := runStore.AppendAgentSelectionExhausted(ctx, AgentSelectionExhaustedRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01", Category: "frontend",
	}); err == nil || !strings.Contains(err.Error(), "does not match latest failed attempt category") {
		t.Fatalf("expected mismatched category to fail, got %v", err)
	}
	activeAttempt := AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeReview, ScopeID: "review:run",
		Category: "review", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "review-model",
		Status: AgentSelectionStatusActive,
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, activeAttempt); err != nil {
		t.Fatalf("append active review attempt: %v", err)
	}
	if _, err := runStore.AppendAgentSelectionExhausted(ctx, AgentSelectionExhaustedRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeReview, ScopeID: "review:run", Category: "review",
	}); err == nil || !strings.Contains(err.Error(), "must be failed") {
		t.Fatalf("expected latest active attempt to block exhaustion, got %v", err)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 3 {
		t.Fatalf("expected failed exhausted events to leave event count 3, got %d", got)
	}
}

func buildV8Fixture(t *testing.T, ctx context.Context, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			work_dir TEXT,
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			reasoning_effort TEXT NOT NULL DEFAULT '',
			owner_pid INTEGER,
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir, spec_slug,
			agent, model, reasoning_effort, owner_pid, stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v8_active', 'implement', 'Active', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', 'tmp/spec-worktree', '0035-agent-selection-profiles',
			'codex', 'gpt-5.6-sol', '', NULL, NULL, '2026-07-17T08:00:00Z', '2026-07-17T08:00:00Z', '')`,
		`INSERT INTO run_events (run_id, cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload)
		 VALUES ('run_v8_active', 1, 0, 'daemon', 'daemon.status', '', '', '', 'schema 8 event', '2026-07-17T08:00:00Z', '{"state":"Active"}')`,
		`PRAGMA user_version = 8`,
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("build v8 fixture: %v", err)
		}
	}
}

func assertAgentSelectionTableExists(t *testing.T, ctx context.Context, store *Store) {
	t.Helper()
	var name string
	if err := store.db.QueryRowContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'run_agent_selections'`).Scan(&name); err != nil {
		t.Fatalf("expected run_agent_selections table: %v", err)
	}
}

func assertAgentSelectionSchemaPrivacy(t *testing.T, ctx context.Context, store *Store) {
	t.Helper()
	rows, err := store.db.QueryContext(ctx, `PRAGMA table_info(run_agent_selections)`)
	if err != nil {
		t.Fatalf("inspect run_agent_selections: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue any
		var primaryKey int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatalf("scan run_agent_selections column: %v", err)
		}
		assertNoSensitiveSelectionPayload(t, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate run_agent_selections columns: %v", err)
	}
}

func assertNoSensitiveSelectionPayload(t *testing.T, value string) {
	t.Helper()
	lower := strings.ToLower(value)
	for _, forbidden := range []string{"prompt", "transcript", "credential", "token", "cookie", "secret", "runtime_owned_config"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("expected no sensitive Agent Selection field %q in %s", forbidden, value)
		}
	}
}

func assertSelectionAttemptMatches(t *testing.T, got AgentSelectionAttempt, want AgentSelectionAttemptRequest) {
	t.Helper()
	if got.RunID != strings.TrimSpace(want.RunID) ||
		got.ScopeKind != want.ScopeKind ||
		got.ScopeID != want.ScopeID ||
		got.Category != want.Category ||
		got.ProfileSource != want.ProfileSource ||
		got.Attempt != want.Attempt ||
		got.SelectionRole != want.SelectionRole ||
		got.FallbackIndex != want.FallbackIndex ||
		got.Runtime != want.Runtime ||
		got.Model != want.Model ||
		got.ReasoningEffort != want.ReasoningEffort ||
		got.Status != want.Status ||
		got.ReasonCode != want.ReasonCode ||
		got.Reason != want.Reason ||
		!got.CreatedAt.Equal(want.Time) {
		t.Fatalf("selection attempt mismatch\ngot  %#v\nwant %#v", got, want)
	}
}

func withAttemptOverride(req AgentSelectionAttemptRequest, apply func(*AgentSelectionAttemptRequest)) AgentSelectionAttemptRequest {
	apply(&req)
	return req
}

func countAgentSelectionAttempts(t *testing.T, ctx context.Context, store *Store, runID string) int {
	t.Helper()
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM run_agent_selections WHERE run_id = ?`, runID).Scan(&count); err != nil {
		t.Fatalf("count Agent Selection attempts for %s: %v", runID, err)
	}
	return count
}

func selectionEventKinds(events []JournalEvent) []runevent.Kind {
	kinds := make([]runevent.Kind, 0, len(events))
	for _, event := range events {
		kinds = append(kinds, event.Event.Kind)
	}
	return kinds
}

func kindStrings(kinds []runevent.Kind) []string {
	values := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		values = append(values, string(kind))
	}
	return values
}
