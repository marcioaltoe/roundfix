// Suite: durable prompt usage.
// Invariant: absent usage stays absent; rows and events commit together.
// Boundary IN: migrated SQLite, queue writes, and journal retention.
// Boundary OUT: adapter parsing and CLI rendering.
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

func usageInt(n int64) *int64       { return &n }
func usageFloat(n float64) *float64 { return &n }
func usageStore(t *testing.T) (*Store, Run) {
	t.Helper()
	s := openTestStore(t, context.Background(), t.TempDir())
	t.Cleanup(func() { closeStore(t, s) })
	run, err := s.CreateRun(context.Background(), sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatal(err)
	}
	return s, run
}
func usageRecord(runID, scope, session string, n *int64) TokenUsageRecord {
	basis := "unreported"
	if n != nil {
		basis = "turn"
	}
	return TokenUsageRecord{RunID: runID, ScopeKind: "task", ScopeID: scope, Session: session, Basis: basis, TotalTokens: n, Runtime: "claude", Model: "opus", ReasoningEffort: "high"}
}
func appendUsage(t *testing.T, s *Store, r TokenUsageRecord) {
	t.Helper()
	if err := s.AppendTokenUsage(context.Background(), r); err != nil {
		t.Fatal(err)
	}
}
func usageReport(t *testing.T, s *Store, id string) TokenUsageReport {
	t.Helper()
	r, err := s.RunTokenUsage(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestMigrationAddsTokenUsageTablesToThePreviousSchema(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	s := openTestStore(t, ctx, home)
	run, err := s.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatal(err)
	}
	queue, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"spec"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.ExecContext(ctx, `DROP TABLE run_token_usage; DROP TABLE delivery_queue_runs; ALTER TABLE delivery_queues DROP COLUMN max_tokens; PRAGMA user_version = `+fmt.Sprint(schemaVersion-1)); err != nil {
		t.Fatal(err)
	}
	closeStore(t, s)
	s = openTestStore(t, ctx, home)
	defer closeStore(t, s)
	version, err := s.MigrationVersion(ctx)
	if err != nil || version != schemaVersion {
		t.Fatalf("version=%d err=%v", version, err)
	}
	if _, found, err := s.Run(ctx, run.ID); err != nil || !found {
		t.Fatalf("old Run lost: %v", err)
	}
	got, found, err := s.DeliveryQueue(ctx, "/repo")
	if err != nil || !found || !reflect.DeepEqual(got, queue) {
		t.Fatalf("old queue changed: %+v %v", got, err)
	}
	for _, statement := range []string{
		`INSERT INTO run_token_usage (run_id,scope_kind,scope_id,session,basis,total_tokens,created_at) VALUES ('missing','task','x','s','turn',1,'now')`,
		`INSERT INTO run_token_usage (run_id,scope_kind,scope_id,session,basis,created_at) VALUES ('` + run.ID + `','bad','x','s','unreported','now')`,
		`INSERT INTO run_token_usage (run_id,scope_kind,scope_id,session,basis,created_at) VALUES ('` + run.ID + `','task','x','s','bad','now')`,
		`INSERT INTO run_token_usage (run_id,scope_kind,scope_id,session,basis,total_tokens,created_at) VALUES ('` + run.ID + `','task','x','s','unreported',0,'now')`,
		`INSERT INTO run_token_usage (run_id,scope_kind,scope_id,session,basis,created_at) VALUES ('` + run.ID + `','task','x','s','turn','now')`,
		`INSERT INTO delivery_queue_runs VALUES ('missing','spec','run')`,
	} {
		if _, err := s.db.ExecContext(ctx, statement); err == nil {
			t.Fatalf("constraint accepted %s", statement)
		}
	}
	appendUsage(t, s, usageRecord(run.ID, "task_01", "s", usageInt(7)))
	if _, err := s.db.ExecContext(ctx, `INSERT INTO delivery_queue_runs VALUES ('/repo','spec','run')`); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO delivery_queue_runs VALUES ('/repo','other','run')`); err == nil {
		t.Fatal("duplicate queue Run accepted")
	}
}

func TestAppendTokenUsageWritesTheRowAndItsRunEvent(t *testing.T) {
	s, run := usageStore(t)
	ctx := context.Background()
	r := usageRecord(run.ID, "task_01", "session", usageInt(12))
	r.Attempt = 2
	r.Readings = 3
	r.InputTokens = usageInt(8)
	r.OutputTokens = usageInt(4)
	r.CostAmount = usageFloat(1.25)
	r.CostCurrency = "USD"
	r.Time = time.Date(2026, 10, 1, 12, 4, 0, 0, time.UTC)
	appendUsage(t, s, r)
	events, err := s.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Event.Kind != runevent.KindDaemonTokenUsage || events[0].Event.Source != runevent.SourceDaemon || events[0].Event.ReviewIssue != r.ScopeID || !events[0].Event.Time.Equal(r.Time) {
		t.Fatalf("events: %+v", events)
	}
	var payload map[string]any
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["total_tokens"] != float64(12) || payload["attempt"] != float64(2) || payload["input_tokens"] != float64(8) || payload["session"] != "session" || payload["cost"].(map[string]any)["amount"] != 1.25 {
		t.Fatalf("payload=%v", payload)
	}
	var session string
	var attempt, readings int
	var tokens int64
	if err := s.db.QueryRowContext(ctx, `SELECT session,attempt,readings,total_tokens FROM run_token_usage WHERE run_id=?`, run.ID).Scan(&session, &attempt, &readings, &tokens); err != nil || session != "session" || attempt != 2 || readings != 3 || tokens != 12 {
		t.Fatalf("row %s %d %d %d: %v", session, attempt, readings, tokens, err)
	}
	if _, err := s.db.ExecContext(ctx, `CREATE TRIGGER reject_usage_event BEFORE INSERT ON run_events WHEN NEW.kind='daemon.token_usage' BEGIN SELECT RAISE(ABORT,'event rejected'); END`); err != nil {
		t.Fatal(err)
	}
	if err := s.AppendTokenUsage(ctx, r); err == nil {
		t.Fatal("event failure accepted")
	}
	if got := usageReport(t, s, run.ID).Total.Prompts; got != 1 {
		t.Fatalf("event failure left row: %d", got)
	}
	if _, err := s.db.ExecContext(ctx, `DROP TRIGGER reject_usage_event`); err != nil {
		t.Fatal(err)
	}
	appendUsage(t, s, usageRecord(run.ID, "qa", "qa", nil))
	events, err = s.RunEventsAfter(ctx, run.ID, 1, 10)
	if err != nil || len(events) != 1 {
		t.Fatalf("unreported event: %v %v", events, err)
	}
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	// Decode into a fresh map: omitted JSON members must stay absent.
	payload = map[string]any{}
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if _, ok := payload["total_tokens"]; ok || payload["basis"] != "unreported" {
		t.Fatalf("unreported=%v", payload)
	}
}

func TestRunTokenUsageSumsScopesAndTheRun(t *testing.T) {
	s, run := usageStore(t)
	r := usageRecord(run.ID, "second", "s1", usageInt(10))
	r.InputTokens = usageInt(7)
	r.OutputTokens = usageInt(3)
	appendUsage(t, s, r)
	r = usageRecord(run.ID, "first", "s2", usageInt(20))
	r.Basis = "request-sum"
	r.Runtime = "codex"
	appendUsage(t, s, r)
	r = usageRecord(run.ID, "second", "s1", usageInt(5))
	r.InputTokens = usageInt(4)
	r.OutputTokens = usageInt(1)
	appendUsage(t, s, r)
	report := usageReport(t, s, run.ID)
	if len(report.Scopes) != 2 || report.Scopes[0].ScopeID != "second" || report.Scopes[1].ScopeID != "first" {
		t.Fatalf("scope order=%+v", report.Scopes)
	}
	if report.Total.Tokens == nil || *report.Total.Tokens != 35 || report.Total.Prompts != 3 || report.Total.ReportedPrompts != 3 || report.Total.Sessions != 2 || report.Total.InputTokens != nil || !reflect.DeepEqual(report.Total.Bases, []string{"turn", "request-sum"}) {
		t.Fatalf("total=%+v", report.Total)
	}
	scope := report.Scopes[0]
	if *scope.Tokens != 15 || scope.InputTokens == nil || *scope.InputTokens != 11 || *scope.OutputTokens != 4 || len(scope.Selections) != 1 {
		t.Fatalf("scope=%+v", scope)
	}
}

func TestRunTokenUsageKeepsUnreportedPromptsOutOfTheTotal(t *testing.T) {
	s, run := usageStore(t)
	if r := usageReport(t, s, run.ID); r.Total.Tokens != nil || r.Total.Prompts != 0 {
		t.Fatalf("empty=%+v", r)
	}
	appendUsage(t, s, usageRecord(run.ID, "unknown", "s", nil))
	r := usageReport(t, s, run.ID)
	if r.Total.Tokens != nil || r.Total.Prompts != 1 || r.Total.ReportedPrompts != 0 || r.Scopes[0].Tokens != nil {
		t.Fatalf("unreported became zero: %+v", r)
	}
	zero := usageRecord(run.ID, "zero", "s2", usageInt(0))
	zero.InputTokens = usageInt(0)
	appendUsage(t, s, zero)
	r = usageReport(t, s, run.ID)
	if r.Total.Tokens == nil || *r.Total.Tokens != 0 || r.Total.ReportedPrompts != 1 || r.Total.Prompts != 2 || r.Total.InputTokens == nil {
		t.Fatalf("reported zero=%+v", r.Total)
	}
}

func TestSessionSpendSumsIncreasesAndRestartsOnADrop(t *testing.T) {
	s, run := usageStore(t)
	for _, entry := range []struct {
		session, currency string
		amount            float64
	}{{"s", "USD", 3}, {"s", "EUR", 1}, {"s", "USD", 5}, {"s", "USD", 2}, {"s", "USD", 4}, {"s", "EUR", 2}, {"s2", "USD", 1}, {"s2", "USD", 1}} {
		r := usageRecord(run.ID, "task_01", entry.session, nil)
		r.CostAmount = usageFloat(entry.amount)
		r.CostCurrency = entry.currency
		appendUsage(t, s, r)
	}
	r := usageReport(t, s, run.ID)
	if r.Total.Sessions != 2 || r.Total.CostSessions != 2 || r.Total.Tokens != nil || len(r.Total.Costs) != 2 || math.Abs(r.Total.Costs[0].Amount-10) > 1e-9 || r.Total.Costs[1].Amount != 2 {
		t.Fatalf("spend=%+v", r.Total)
	}
	if !reflect.DeepEqual(r.Total.Costs, r.Scopes[0].Costs) {
		t.Fatalf("scope cost=%v total=%v", r.Scopes[0].Costs, r.Total.Costs)
	}
}

func TestDeliveryQueueRunsLinkEveryRecordedRunID(t *testing.T) {
	s, _ := usageStore(t)
	ctx := context.Background()
	q, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"spec"})
	if err != nil {
		t.Fatal(err)
	}
	item := q.Items[0]
	item.RunID = "run1"
	item.Stage = DeliveryStageParked
	item.Blocker = "retry"
	for range 2 {
		if err := s.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
			t.Fatal(err)
		}
	}
	item.RunID = "run2"
	item.Stage = DeliveryStageRunning
	if _, _, err := s.RetryDeliveryQueueItem(ctx, "/repo", item, "retry"); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM delivery_queue_runs WHERE git_root='/repo'`).Scan(&count); err != nil || count != 2 {
		t.Fatalf("links=%d err=%v", count, err)
	}
	item.Stage = DeliveryStageParked
	item.RunID = ""
	if err := s.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"new-spec"}); err != nil {
		t.Fatal(err)
	}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM delivery_queue_runs`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("replaced queue retained links=%d err=%v", count, err)
	}
}

func TestDeliveryQueueTokenUsageSumsLinkedRuns(t *testing.T) {
	s, first := usageStore(t)
	ctx := context.Background()
	secondRequest := sampleImplementCreateRunRequest()
	secondRequest.SpecSlug = "second-spec"
	second, err := s.CreateRun(ctx, secondRequest)
	if err != nil {
		t.Fatal(err)
	}
	thirdRequest := sampleImplementCreateRunRequest()
	thirdRequest.SpecSlug = "third-spec"
	third, err := s.CreateRun(ctx, thirdRequest)
	if err != nil {
		t.Fatal(err)
	}
	appendUsage(t, s, usageRecord(first.ID, "task_01", "same-session", usageInt(10)))
	r := usageRecord(second.ID, "task_01", "same-session", usageInt(20))
	r.CostAmount = usageFloat(2)
	r.CostCurrency = "USD"
	appendUsage(t, s, r)
	r = usageRecord(first.ID, "task_01", "same-session", nil)
	r.CostAmount = usageFloat(3)
	r.CostCurrency = "USD"
	appendUsage(t, s, r)
	appendUsage(t, s, usageRecord(third.ID, "task_01", "s", usageInt(999)))
	q, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"spec"})
	if err != nil {
		t.Fatal(err)
	}
	item := q.Items[0]
	item.RunID = first.ID
	item.Stage = DeliveryStageParked
	item.Blocker = "retry"
	if err := s.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	item.RunID = second.ID
	item.Stage = DeliveryStageRunning
	if _, _, err := s.RetryDeliveryQueueItem(ctx, "/repo", item, "retry"); err != nil {
		t.Fatal(err)
	}
	// A linked Run without prompts still contributes to the Run count.
	if err := linkUsageTestRun(ctx, s, "/repo", "missing-usage"); err != nil {
		t.Fatal(err)
	}
	report, err := s.DeliveryQueueTokenUsage(ctx, "/repo")
	if err != nil {
		t.Fatal(err)
	}
	if report.Total.Tokens == nil || *report.Total.Tokens != 30 || report.Total.Runs != 3 || report.Total.Prompts != 3 || report.Total.ReportedPrompts != 2 || report.Total.Sessions != 2 || report.Total.Costs[0].Amount != 5 {
		t.Fatalf("queue=%+v", report)
	}
}
func linkUsageTestRun(ctx context.Context, s *Store, root, run string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO delivery_queue_runs VALUES (?,'spec',?)`, root, run)
	return err
}

func TestDeliveryQueueLimitsRoundTripMaxTokens(t *testing.T) {
	s, _ := usageStore(t)
	ctx := context.Background()
	want := DeliveryQueueLimits{MaxTokens: 5000000, MaxRetries: 2, Deadline: time.Date(2026, 10, 1, 20, 0, 0, 0, time.UTC)}
	if _, err := s.CreateDeliveryQueueWithLimits(ctx, "/repo", []string{"spec"}, want); err != nil {
		t.Fatal(err)
	}
	got, found, err := s.DeliveryQueue(ctx, "/repo")
	if err != nil || !found || got.Limits != want {
		t.Fatalf("limits=%+v err=%v", got.Limits, err)
	}
}

func TestRetentionLeavesTokenUsageTables(t *testing.T) {
	s, run := usageStore(t)
	ctx := context.Background()
	appendUsage(t, s, usageRecord(run.ID, "task_01", "s", usageInt(12)))
	q, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"spec"})
	if err != nil {
		t.Fatal(err)
	}
	item := q.Items[0]
	item.RunID = run.ID
	if err := s.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatal(err)
	}
	pruned, err := s.PruneTerminalRuns(ctx, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if countRunEvents(t, ctx, s, run.ID) != 0 {
		t.Fatalf("journal not pruned: %+v", pruned)
	}
	report, err := s.DeliveryQueueTokenUsage(ctx, "/repo")
	if err != nil || report.Total.Runs != 1 || report.Total.Prompts != 1 || *report.Total.Tokens != 12 {
		t.Fatalf("retention removed usage: %+v err=%v", report, err)
	}
}

func TestDeliveryQueueUsageLinkFailureRollsBackTheItem(t *testing.T) {
	s, _ := usageStore(t)
	ctx := context.Background()
	q, err := s.CreateDeliveryQueue(ctx, "/repo", []string{"spec"})
	if err != nil {
		t.Fatal(err)
	}
	item := q.Items[0]
	item.Stage = DeliveryStageParked
	item.Blocker = "retry"
	if err := s.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.ExecContext(ctx, `CREATE TRIGGER reject_queue_link BEFORE INSERT ON delivery_queue_runs BEGIN SELECT RAISE(ABORT,'link rejected'); END`); err != nil {
		t.Fatal(err)
	}
	proposed := item
	proposed.RunID = "new-run"
	proposed.Stage = DeliveryStageRunning
	if err := s.UpdateDeliveryQueueItem(ctx, "/repo", proposed); err == nil {
		t.Fatal("update accepted rejected link")
	}
	if _, _, err := s.RetryDeliveryQueueItem(ctx, "/repo", proposed, "retry"); err == nil {
		t.Fatal("retry accepted rejected link")
	}
	got, found, err := s.DeliveryQueue(ctx, "/repo")
	if err != nil || !found || !reflect.DeepEqual(got.Items[0], item) {
		t.Fatalf("rejected link changed item=%+v err=%v", got, err)
	}
}
