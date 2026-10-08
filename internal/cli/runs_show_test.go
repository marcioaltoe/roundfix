// Suite: Run usage readers.
// Invariant: each recorded scope stays visible, absent usage stays absent, and reads never write.
// Boundary IN: CLI dispatch and temporary Roundfix Home and repository.
// Boundary OUT: adapter parsing and store aggregation, covered by their owning packages.
package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func seedUsageRun(t *testing.T, unreported bool) (string, string, store.Run) {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	run := createEventsRun(t, home, repo, store.StateClean)
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	a, b, cost := int64(5639755), int64(4019995), 3.55
	rows := []store.TokenUsageRecord{
		{ScopeKind: "task", ScopeID: "task_01", Session: "first", Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high", Basis: "request-sum", TotalTokens: &a},
		{ScopeKind: "task", ScopeID: "task_02", Session: "second", Runtime: "claude", Model: "opus", ReasoningEffort: "high", Basis: "turn", TotalTokens: &b, CostAmount: &cost, CostCurrency: "USD"},
		{ScopeKind: "qa", ScopeID: "qa", Session: "gate", Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high", Basis: "unreported"},
	}
	for _, row := range rows {
		row.RunID = run.ID
		if unreported {
			row.Basis = "unreported"
			row.TotalTokens = nil
			row.CostAmount = nil
			row.CostCurrency = ""
		}
		if err := writer.AppendTokenUsage(t.Context(), row); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return home, repo, run
}

func TestRunsShowPrintsEachScopeAndTheTotal(t *testing.T) {
	t.Parallel()
	_, _, run := seedUsageRun(t, false)
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"runs", "show", run.ID}, &out, &diag)
	want := fmt.Sprintf("Run %s: %s %s %s\n", run.ID, run.Kind, run.SpecSlug, run.State) +
		"task_01\tcodex/gpt-6.1-sol/high\t5639755 tokens (request-sum) from 1 of 1 prompt(s)\tcost not reported\n" +
		"task_02\tclaude/opus/high\t4019995 tokens (turn) from 1 of 1 prompt(s)\tcost 3.55 USD from 1 of 1 Agent Session(s)\n" +
		"qa\tcodex/gpt-6.1-sol/high\tnone reported by 1 prompt(s)\tcost not reported\n" +
		"Tokens: 9659750 from 2 of 3 prompt(s); cost 3.55 USD from 1 of 3 Agent Session(s)\n"
	if code != 0 || diag.Len() != 0 || out.String() != want {
		t.Fatalf("exit=%d stderr=%s\nstdout=%q\nwant=%q", code, &diag, out.String(), want)
	}
}

func TestRunsShowRefusesAnUnknownRun(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"runs", "show", "run_missing"}, &out, &diag)
	want := "roundfix: runs show failed: Run \"run_missing\" does not exist\nRun 'roundfix runs show --help' for usage.\n"
	if code != 2 || out.Len() != 0 || diag.String() != want {
		t.Fatalf("exit=%d out=%s stderr=%q", code, &out, diag.String())
	}
}

func TestRunsShowJSONIsSchemaV1(t *testing.T) {
	t.Parallel()
	_, _, run := seedUsageRun(t, false)
	for _, args := range [][]string{{"runs", "show", run.ID, "--json"}, {"runs", "show", "--json", run.ID}} {
		var out, diag bytes.Buffer
		code := runCLI(t, args, &out, &diag)
		var got struct {
			Schema   string `json:"schema"`
			RunID    string `json:"run_id"`
			Kind     string `json:"kind"`
			SpecSlug string `json:"spec_slug"`
			State    string `json:"state"`
			store.TokenUsageReport
		}
		if err := json.Unmarshal(out.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if code != 0 || diag.Len() != 0 || got.Schema != "roundfix/runs-show/v1" || got.RunID != run.ID || got.Kind != run.Kind || got.SpecSlug != run.SpecSlug || got.State != run.State || len(got.Scopes) != 3 || got.Total.Tokens == nil || *got.Total.Tokens != 9659750 || got.Total.Prompts != 3 || got.Total.ReportedPrompts != 2 || len(got.Total.Costs) != 1 || got.Total.Costs[0].Amount != 3.55 {
			t.Fatalf("exit=%d stderr=%s report=%+v", code, &diag, got)
		}
		if got.Scopes[2].Tokens != nil || got.Scopes[0].InputTokens != nil || got.Scopes[2].ScopeKind != "qa" {
			t.Fatalf("unknown counts fabricated: %+v", got.Scopes)
		}
		var raw map[string]any
		if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
			t.Fatal(err)
		}
		scopes := raw["scopes"].([]any)
		for _, key := range []string{"tokens", "input_tokens", "output_tokens", "cached_read_tokens", "cached_write_tokens", "thought_tokens"} {
			value, exists := scopes[2].(map[string]any)[key]
			if !exists || value != nil {
				t.Fatalf("unreported %s must be explicit null: %v", key, scopes[2])
			}
		}
	}
}

func TestRunsShowRefusesAMissingOrExtraArgument(t *testing.T) {
	t.Parallel()
	withCLIWorkspace(t)
	for _, args := range [][]string{{"runs", "show"}, {"runs", "show", "one", "two"}, {"runs", "show", "one", "--unknown"}, {"runs", "show", "--json"}} {
		var out, diag bytes.Buffer
		if code := runCLI(t, args, &out, &diag); code != 2 || out.Len() != 0 || !strings.Contains(diag.String(), "runs show failed") {
			t.Fatalf("args=%v exit=%d out=%s stderr=%s", args, code, &out, &diag)
		}
	}
}

func snapshotUsageHome(t *testing.T, home string) map[string]string {
	t.Helper()
	result := map[string]string{}
	err := filepath.WalkDir(home, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || strings.HasSuffix(path, "-shm") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		// A mode=ro live reader may create SQLite coordination sidecars.
		// An empty WAL contains no database writes; any appended frame fails.
		if strings.HasSuffix(path, "-wal") && len(data) == 0 {
			return nil
		}
		result[path] = fmt.Sprintf("%x:%o:%d", sha256.Sum256(data), info.Mode(), info.ModTime().UnixNano())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestRunsShowWritesNothing(t *testing.T) {
	t.Parallel()
	home, _, run := seedUsageRun(t, false)
	before := snapshotUsageHome(t, home)
	for _, args := range [][]string{{"runs", "show", run.ID}, {"runs", "show", run.ID, "--json"}, {"runs", "show", "run_missing"}} {
		var out, diag bytes.Buffer
		runCLI(t, args, &out, &diag)
		after := snapshotUsageHome(t, home)
		if len(after) != len(before) {
			t.Fatalf("files changed: %v", after)
		}
		for path, value := range before {
			if after[path] != value {
				t.Fatalf("reader wrote %s", path)
			}
		}
	}
}

func TestRunsShowPrintsNoPromptsRecorded(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	run := createEventsRun(t, home, repo, store.StateClean)
	var out, diag bytes.Buffer
	if code := runCLI(t, []string{"runs", "show", run.ID}, &out, &diag); code != 0 || diag.Len() != 0 || !strings.HasSuffix(out.String(), "Tokens: no prompts recorded; cost not reported\n") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
}

func TestTokenTotalsRenderReportedZeroAndCurrencies(t *testing.T) {
	t.Parallel()
	zero := int64(0)
	tokens, cost := formatTokenTotals(store.TokenTotals{Prompts: 2, ReportedPrompts: 1, Tokens: &zero, Costs: []store.TokenCost{{Currency: "USD", Amount: 3.556}, {Currency: "EUR", Amount: 1.2}}, Sessions: 2, CostSessions: 1})
	if tokens != "0 from 1 of 2 prompt(s)" || cost != "cost 3.56 USD + 1.20 EUR from 1 of 2 Agent Session(s)" {
		t.Fatalf("tokens=%s cost=%s", tokens, cost)
	}
}
