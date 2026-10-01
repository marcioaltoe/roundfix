// Suite: events usage surface.
// Invariant: the CLI accepts usage and prints Surface Transcript 7 as JSONL.
// Boundary IN: command dispatch, temporary database reader, stdout/stderr.
// Boundary OUT: payload parser details, owned by internal/runevent/stream_usage_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestEventsHelpNamesTheUsageCategory(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"events", "--help"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 || !strings.Contains(stdout.String(), "task-status,batch,verification,outcome,agent-selection,usage") {
		t.Fatalf("help=%s stderr=%s exit=%d", &stdout, &stderr, code)
	}
}

func TestEventsUsageTranscript(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	run := createEventsRun(t, home, repo, store.StateClean)
	ctx := context.Background()
	writer, err := store.Open(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	tokens := int64(5639755)
	if err := writer.AppendTokenUsage(ctx, store.TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: "task_01", Session: "roundfix-run_1-task_01", Attempt: 1, Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high", Basis: "request-sum", TotalTokens: &tokens, Readings: 46, Time: time.Date(2026, 10, 1, 12, 4, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"events", run.ID, "--filter", "usage"}, {"events", run.ID}} {
		var stdout, stderr bytes.Buffer
		code := runCLIContext(t, ctx, args, &stdout, &stderr)
		if code != 0 || stderr.Len() != 0 {
			t.Fatalf("exit=%d stderr=%s", code, &stderr)
		}
		var record runevent.StreamRecord
		decoder := json.NewDecoder(&stdout)
		found := false
		for decoder.More() {
			if err := decoder.Decode(&record); err != nil {
				t.Fatal(err)
			}
			if record.Category != runevent.StreamCategoryUsage {
				if len(args) > 2 {
					t.Fatalf("usage filter emitted %+v", record)
				}
				continue
			}
			found = true
			if record.Schema != runevent.StreamSchema || record.RunID != run.ID || record.Time != "2026-10-01T12:04:00Z" || record.WorkItem != "task_01" || record.ScopeKind != "task" || record.ScopeID != "task_01" || record.Runtime != "codex" || record.Model != "gpt-6.1-sol" || record.ReasoningEffort != "high" || record.TokenBasis != "request-sum" || record.Tokens == nil || *record.Tokens != tokens || record.Summary != "task_01 used 5639755 tokens (request-sum)" {
				t.Fatalf("transcript=%+v", record)
			}
		}
		if !found {
			t.Fatal("usage record missing")
		}
	}
}
