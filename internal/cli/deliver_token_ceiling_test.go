// Suite: Delivery token ceiling command surface.
// Invariant: the ceiling is validated before writes and exposes the authored status and retry transcripts.
// Boundary IN: CLI dispatch, Delivery Engine and a temporary Roundfix Home with real SQLite records.
// Boundary OUT: detached owner launch; no ACP, provider, GitHub or live Home is reached.
package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func TestDeliverStartRecordsAndPrintsTheTokenCeiling(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repo, allDeliveryOperations...)
	started := 0
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = readyDeliveryReadiness
		deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "start", implementTestSlug, "--max-tokens", "5000000"}, &out, &diag)
	want := "Limits: deadline none, retries per item none, concurrency 1, tokens 5000000\n"
	if code != exitOK || diag.Len() != 0 || out.String() != want || started != 1 {
		t.Fatalf("exit=%d starts=%d stdout=%q stderr=%q", code, started, out.String(), diag.String())
	}
	if queue := openDeliveryQueueForCLI(t, home, repo); queue.Limits.MaxTokens != 5000000 {
		t.Fatalf("limits=%+v", queue.Limits)
	}
}

func TestDeliverStartRefusesAMaxTokensBelowOne(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "-1"} {
		t.Run(value, func(t *testing.T) {
			assertDeliverStartLimitsRefused(t, []string{"--max-tokens", value}, "max-tokens must be at least 1")
		})
	}
}

func TestDeliverStartRefusesANonIntegerMaxTokens(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"1.5", "not-a-number", "9223372036854775808"} {
		t.Run(value, func(t *testing.T) {
			assertDeliverStartLimitsRefused(t, []string{"--max-tokens", value}, "invalid value")
		})
	}
}

func seedTokenCeilingTranscript(t *testing.T) (string, string) {
	t.Helper()
	home, repo := withCLIWorkspace(t)
	if err := os.MkdirAll(filepath.Join(repo, "docs", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Error(err)
		}
	}()
	queue, err := writer.CreateDeliveryQueueWithLimits(t.Context(), repo, []string{"0300-example", "0301-example"}, store.DeliveryQueueLimits{MaxTokens: 5000000})
	if err != nil {
		t.Fatal(err)
	}
	run, err := writer.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: "feat/spent", HeadSHA: "spent-head", SpecSlug: "0300-example", Agent: "codex"})
	if err != nil {
		t.Fatal(err)
	}
	tokens := int64(5639755)
	if err := writer.AppendTokenUsage(t.Context(), store.TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: "task_01", Session: "spent", Basis: "request-sum", TotalTokens: &tokens}); err != nil {
		t.Fatal(err)
	}
	first := queue.Items[0]
	first.Stage, first.RunID = store.DeliveryStageMerged, run.ID
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, first); err != nil {
		t.Fatal(err)
	}
	next := queue.Items[1]
	next.Stage, next.Blocker = store.DeliveryStageParked, "queue-token-ceiling"
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, next); err != nil {
		t.Fatal(err)
	}
	return home, repo
}

func TestDeliverStatusPrintsTheTokenCeilingQuestion(t *testing.T) {
	t.Parallel()
	seedTokenCeilingTranscript(t)
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
	want := "0300-example\tmerged\t-\t-\n0301-example\tparked\tqueue-token-ceiling\t-\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens 5000000\n" +
		"Usage: 5639755 tokens from 1 of 1 prompt(s) across 1 Run(s); cost not reported\n" +
		"Pending question: 0301-example parked queue-token-ceiling\n" +
		"Answer: record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens\n"
	if code != exitOK || diag.Len() != 0 || out.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q, want stdout=%q", code, out.String(), diag.String(), want)
	}
}

func TestDeliverRetryIsRefusedAtTheTokenCeiling(t *testing.T) {
	t.Parallel()
	home, repo := seedTokenCeilingTranscript(t)
	before := openDeliveryQueueForCLI(t, home, repo)
	started := 0
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "retry", "0301-example"}, &out, &diag)
	want := "Retry refused\n\nReason:\n  retry Delivery Queue item \"0301-example\": queue token ceiling 5000000 was reached with 5639755 tokens; start a new queue with roundfix deliver start\n"
	if code != exitPreflight || out.Len() != 0 || !strings.HasPrefix(diag.String(), want) || strings.Contains(diag.String(), "Usage:") || started != 0 {
		t.Fatalf("exit=%d starts=%d stdout=%q stderr=%q", code, started, out.String(), diag.String())
	}
	if after := openDeliveryQueueForCLI(t, home, repo); !reflect.DeepEqual(before, after) {
		t.Fatalf("refused retry changed queue: before=%+v after=%+v", before, after)
	}
}

func TestDeliverHelpNamesTheTokenCeiling(t *testing.T) {
	t.Parallel()
	var out, diag bytes.Buffer
	code := Run([]string{"deliver", "--help"}, &out, &diag)
	if code != exitOK || diag.Len() != 0 || !strings.Contains(out.String(), "--max-tokens <n>") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diag.String())
	}
}
