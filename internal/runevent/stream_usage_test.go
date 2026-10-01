// Suite: usage stream projection.
// Invariant: usage keeps absent fields absent and uses the stable stream schema.
// Boundary IN: event payload projection and category selection.
// Boundary OUT: journal writes and CLI replay.
package runevent

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func usageStreamEvent() RunEvent {
	return RunEvent{RunID: "run_20261001T120000Z_0123456789abcdef", Source: SourceDaemon, Kind: KindDaemonTokenUsage, ReviewIssue: "task_01", Time: time.Date(2026, 10, 1, 12, 4, 0, 0, time.UTC), Payload: []byte(`{"scope_kind":"task","scope_id":"task_01","session":"roundfix-run_1-task_01","attempt":1,"runtime":"codex","model":"gpt-6.1-sol","reasoning_effort":"high","basis":"request-sum","total_tokens":5639755,"readings":46}`)}
}
func TestUsageEventProjectsToTheUsageCategory(t *testing.T) {
	event := usageStreamEvent()
	record, ok, err := ProjectStreamEvent(41, event, nil)
	if err != nil || !ok {
		t.Fatalf("projection=%+v %v", record, err)
	}
	raw, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"schema":"` + StreamSchema + `","run_id":"run_20261001T120000Z_0123456789abcdef","category":"usage","time":"2026-10-01T12:04:00Z","cursor":41,"work_item":"task_01","summary":"task_01 used 5639755 tokens (request-sum)","scope_kind":"task","scope_id":"task_01","runtime":"codex","model":"gpt-6.1-sol","reasoning_effort":"high","token_basis":"request-sum","tokens":5639755}`
	var gotMap, wantMap map[string]any
	if err := json.Unmarshal(raw, &gotMap); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &wantMap); err != nil {
		t.Fatal(err)
	}
	if len(gotMap) != len(wantMap) {
		t.Fatalf("record=%s", raw)
	}
	for key, value := range wantMap {
		if gotMap[key] != value {
			t.Fatalf("%s=%v want=%v record=%s", key, gotMap[key], value, raw)
		}
	}
	if IsDaemonKind(KindDaemonTokenUsage) {
		t.Fatal("usage changed the Live Run View daemon vocabulary")
	}
	event.Payload = []byte(`{"scope_kind":"qa","scope_id":"qa","basis":"unreported","cost":{"amount":0,"currency":"USD"}}`)
	record, ok, err = ProjectStreamEvent(42, event, nil)
	if err != nil || !ok || record.Tokens != nil || record.CostAmount == nil || *record.CostAmount != 0 || record.Summary != "qa reported no usage" {
		t.Fatalf("unreported=%+v err=%v", record, err)
	}
	event.Payload = []byte(`{"scope_kind":"task","scope_id":"task_01","basis":"turn","total_tokens":0,"input_tokens":0,"output_tokens":1,"cached_read_tokens":2,"cached_write_tokens":3,"thought_tokens":4,"cost":{"amount":1.5,"currency":"EUR"}}`)
	record, ok, err = ProjectStreamEvent(43, event, nil)
	if err != nil || !ok || record.InputTokens == nil || *record.InputTokens != 0 || *record.OutputTokens != 1 || *record.CachedReadTokens != 2 || *record.CachedWriteTokens != 3 || *record.ThoughtTokens != 4 || *record.CostAmount != 1.5 || record.CostCurrency != "EUR" {
		t.Fatalf("fields=%+v err=%v", record, err)
	}
	for _, payload := range []string{`{"scope_kind":"task","scope_id":"x","basis":"turn"}`, `{"scope_kind":"task","scope_id":"x","basis":"unreported","total_tokens":0}`, `{"scope_id":"x","basis":"unreported"}`, `{"scope_kind":"task","scope_id":"x","basis":"bad"}`, `{"scope_kind":"task","scope_id":"x","basis":"turn","total_tokens":"1"}`} {
		event.Payload = []byte(payload)
		if _, _, err := ProjectStreamEvent(44, event, nil); err == nil {
			t.Fatalf("accepted malformed usage %s", payload)
		}
	}
}
func TestDefaultFilterIncludesUsage(t *testing.T) {
	filter := AllStreamCategories()
	if !filter.Includes(StreamCategoryUsage) {
		t.Fatal("default omitted usage")
	}
	if _, ok, err := ProjectStreamEvent(41, usageStreamEvent(), filter); err != nil || !ok {
		t.Fatalf("default projection=%t %v", ok, err)
	}
}
func TestUsageFilterSelectsOnlyUsage(t *testing.T) {
	filter, err := ParseStreamCategoryFilter("usage")
	if err != nil {
		t.Fatal(err)
	}
	if len(filter) != 1 || !filter.Includes(StreamCategoryUsage) {
		t.Fatalf("filter=%v", filter)
	}
	event := usageStreamEvent()
	if _, ok, err := ProjectStreamEvent(41, event, filter); err != nil || !ok {
		t.Fatalf("usage=%t %v", ok, err)
	}
	event.Kind = KindDaemonOutcome
	event.Payload = []byte(`{"state":"Clean"}`)
	if _, ok, err := ProjectStreamEvent(42, event, filter); err != nil || ok {
		t.Fatalf("outcome selected=%t %v", ok, err)
	}
	if _, err := ParseStreamCategoryFilter("usage,typo"); err == nil || !strings.Contains(err.Error(), "typo") {
		t.Fatalf("invalid filter=%v", err)
	}
}
