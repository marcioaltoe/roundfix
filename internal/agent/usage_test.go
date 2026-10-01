package agent

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Suite: prompt usage.
// Invariant: counts retain the adapter report's scope without changing prompt outcomes.
// Boundary IN: counting and RunPrompt reading JSON-RPC through the existing fake acpx process.
// Boundary OUT: live adapters, persistence and user-facing totals (later Tasks).

func usageInt(value int64) *int64 { return &value }

func assertTurnUsage(t *testing.T, got, want TurnUsage) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TurnUsage = %+v, want %+v", got, want)
	}
}

func TestCountTurnUsageSumsReadingsForALastRequestLineage(t *testing.T) {
	t.Parallel()
	for _, reported := range []*promptUsage{nil, {TotalTokens: 150849, InputTokens: usageInt(150000), OutputTokens: usageInt(849)}, {TotalTokens: 1}} {
		assertTurnUsage(t, countTurnUsage(true, reported, []int64{26830, 36330, 150849}, nil), TurnUsage{Basis: UsageBasisRequestSum, TotalTokens: 214009, Readings: 3})
	}
}

func TestCountTurnUsageTakesTheReportedTurnForOtherLineages(t *testing.T) {
	t.Parallel()
	cost := &ReportedCost{Amount: 3.5526034999999996, Currency: "USD"}
	report := &promptUsage{TotalTokens: 4019995, InputTokens: usageInt(19000), OutputTokens: usageInt(995), CachedReadTokens: usageInt(3900000), CachedWriteTokens: usageInt(100000), ThoughtTokens: usageInt(0)}
	want := TurnUsage{Basis: UsageBasisTurn, TotalTokens: 4019995, InputTokens: usageInt(19000), OutputTokens: usageInt(995), CachedReadTokens: usageInt(3900000), CachedWriteTokens: usageInt(100000), ThoughtTokens: usageInt(0), Readings: 3, Cost: cost}
	assertTurnUsage(t, countTurnUsage(false, report, []int64{100000, 110000, 120000}, cost), want)
	// A non-Codex report remains authoritative even below its last context reading.
	assertTurnUsage(t, countTurnUsage(false, &promptUsage{TotalTokens: 100, OutputTokens: usageInt(1)}, []int64{500, 600}, cost), TurnUsage{Basis: UsageBasisTurn, TotalTokens: 100, OutputTokens: usageInt(1), Readings: 2, Cost: cost})
	// Zero is a reported total; absent split fields remain absent.
	assertTurnUsage(t, countTurnUsage(false, &promptUsage{}, nil, nil), TurnUsage{Basis: UsageBasisTurn})
}

func TestCountTurnUsageCountsAGrownLastRequestReportAsTurn(t *testing.T) {
	t.Parallel()
	assertTurnUsage(t, countTurnUsage(true, &promptUsage{TotalTokens: 200, OutputTokens: usageInt(10)}, []int64{40, 100}, nil), TurnUsage{Basis: UsageBasisTurn, TotalTokens: 200, OutputTokens: usageInt(10), Readings: 2})
	// With no readings, a last-request lineage still uses the report.
	assertTurnUsage(t, countTurnUsage(true, &promptUsage{TotalTokens: 200}, nil, nil), TurnUsage{Basis: UsageBasisTurn, TotalTokens: 200})
}

func TestCountTurnUsageLeavesANoReportPromptUnreported(t *testing.T) {
	t.Parallel()
	for _, lastRequestOnly := range []bool{false, true} {
		assertTurnUsage(t, countTurnUsage(lastRequestOnly, nil, nil, nil), TurnUsage{})
	}
	cost := &ReportedCost{Amount: 1, Currency: "USD"}
	assertTurnUsage(t, countTurnUsage(false, nil, []int64{100}, cost), TurnUsage{Readings: 1, Cost: cost})
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: acpxPromptResponseLine("end_turn")})
	if run.err != nil {
		t.Fatal(run.err)
	}
	assertTurnUsage(t, run.result.Usage, TurnUsage{})
}

func usageFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/usage/" + name + "-turn.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestRunPromptReturnsTheCodexFixtureUsage(t *testing.T) {
	t.Parallel()
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: usageFixture(t, "codex")})
	if run.err != nil {
		t.Fatal(run.err)
	}
	assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisRequestSum, TotalTokens: 214009, Readings: 3})
	if len(run.sink.Events()) != 0 || run.result.Message != "" || len(run.result.Messages) != 0 {
		t.Fatalf("usage published Agent output: %+v", run.result)
	}
}

func TestRunPromptReturnsTheClaudeFixtureUsageAndCost(t *testing.T) {
	t.Parallel()
	run := runFakeACPXPrompt(t, fakeACPXPrompt{runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP}, stdout: usageFixture(t, "claude")})
	if run.err != nil {
		t.Fatal(run.err)
	}
	assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisTurn, TotalTokens: 4019995, InputTokens: usageInt(19000), OutputTokens: usageInt(995), CachedReadTokens: usageInt(3900000), CachedWriteTokens: usageInt(100000), ThoughtTokens: usageInt(0), Readings: 3, Cost: &ReportedCost{Amount: 3.5526034999999996, Currency: "USD"}})
	if len(run.sink.Events()) != 0 {
		t.Fatal("usage published a Run Event")
	}
}

func TestRunPromptKeepsUsageWhenALaterResultHasNone(t *testing.T) {
	t.Parallel()
	stdout := `{"jsonrpc":"2.0","id":1,"result":{"stopReason":"end_turn","usage":{"totalTokens":42,"outputTokens":0}}}` + "\n"
	for _, later := range []string{acpxPromptResponseLine("end_turn"), `{"jsonrpc":"2.0","id":2,"result":{"stopReason":"end_turn","usage":null}}` + "\n"} {
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout + later})
		if run.err != nil {
			t.Fatal(run.err)
		}
		assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisTurn, TotalTokens: 42, OutputTokens: usageInt(0)})
	}
}

func TestRunPromptIgnoresAMalformedUsagePayload(t *testing.T) {
	t.Parallel()
	for _, payload := range []string{`"wrong"`, `{"totalTokens":"wrong"}`, `{"inputTokens":[]}`} {
		stdout := acpxUpdateLine(`{"sessionId":"sess_fixture","update":{"sessionUpdate":"usage_update","used":"wrong","cost":{"amount":1,"currency":"USD"}}}`)
		stdout += acpxUpdateLine(`{"sessionId":"sess_fixture","update":{"sessionUpdate":"usage_update","used":12,"cost":{"amount":"wrong"}}}`)
		stdout += `{"jsonrpc":"2.0","id":1,"result":{"stopReason":"end_turn","usage":` + payload + `}}` + "\n"
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
		if run.err != nil {
			t.Fatal(run.err)
		}
		assertTurnUsage(t, run.result.Usage, TurnUsage{})
		if run.sink.HasStatus(AgentWorkStartedStatus) {
			t.Fatal("malformed usage counted as Agent output")
		}
	}
	// A malformed later report also cannot discard a valid report or cost.
	stdout := usageFixture(t, "claude") + `{"jsonrpc":"2.0","id":2,"result":{"stopReason":"end_turn","usage":{"totalTokens":"wrong"}}}` + "\n"
	stdout += acpxUpdateLine(`{"sessionId":"sess_fixture","update":{"sessionUpdate":"usage_update","used":"wrong"}}`)
	run := runFakeACPXPrompt(t, fakeACPXPrompt{runtime: RuntimeSpec{ID: "claude", Protocol: ProtocolACP}, stdout: stdout})
	if run.err != nil {
		t.Fatal(run.err)
	}
	if run.result.Usage.TotalTokens != 4019995 || run.result.Usage.Readings != 3 || run.result.Usage.Cost == nil {
		t.Fatalf("malformed report erased usage: %+v", run.result.Usage)
	}
}

func TestRunPromptReturnsUsageForAStoppedPrompt(t *testing.T) {
	t.Parallel()
	stdout := acpxUpdateLine(`{"sessionId":"sess_fixture","update":{"sessionUpdate":"usage_update","used":42,"cost":{"amount":1.5,"currency":"USD"}}}`)
	for _, tail := range []string{"", acpxPromptResponseLine("cancelled")} {
		run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout + tail, exitCode: 130})
		if !IsStopError(run.err) {
			t.Fatalf("error = %v, want StopError", run.err)
		}
		assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisRequestSum, TotalTokens: 42, Readings: 1, Cost: &ReportedCost{Amount: 1.5, Currency: "USD"}})
	}
}

func TestOnlyTheCodexLineageReportsItsLastRequestOnly(t *testing.T) {
	t.Parallel()
	if !adapterLineageContracts["codex"].LastRequestOnly {
		t.Fatal("Codex must report its last request only")
	}
	for name, contract := range adapterLineageContracts {
		if name != "codex" && contract.LastRequestOnly {
			t.Fatalf("%s reports last request only", name)
		}
	}
	for _, runtime := range []string{"claude", "custom"} {
		run := runFakeACPXPrompt(t, fakeACPXPrompt{runtime: RuntimeSpec{ID: runtime, Protocol: ProtocolACP}, stdout: usageFixture(t, "codex")})
		if run.err != nil {
			t.Fatal(run.err)
		}
		assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisTurn, TotalTokens: 150849, InputTokens: usageInt(150000), OutputTokens: usageInt(849), Readings: 3})
	}
}

func TestRunPromptPreservesUsageAfterTransportAnomaly(t *testing.T) {
	t.Parallel()
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: usageFixture(t, "codex"), exitCode: 1})
	if run.err != nil || run.result.TransportAnomaly == "" {
		t.Fatalf("result = %+v, error = %v", run.result, run.err)
	}
	assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisRequestSum, TotalTokens: 214009, Readings: 3})
}

func TestRunPromptUsageDoesNotChangeNoOutputClassification(t *testing.T) {
	t.Parallel()
	stdout := strings.Split(usageFixture(t, "codex"), "\n")[0] + "\n"
	run := runFakeACPXPrompt(t, fakeACPXPrompt{stdout: stdout})
	var failure *SelectionFailureError
	if !errors.As(run.err, &failure) || run.sink.HasStatus(AgentWorkStartedStatus) {
		t.Fatalf("usage changed output classification: %v", run.err)
	}
	assertTurnUsage(t, run.result.Usage, TurnUsage{Basis: UsageBasisRequestSum, TotalTokens: 26830, Readings: 1})
}
