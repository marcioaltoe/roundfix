// Suite: Delivery usage output.
// Invariant: limits reflect the stored ceiling; usage includes every linked Run, including retries.
// Boundary IN: CLI, temporary Roundfix Home/repository and persisted queue links.
// Boundary OUT: enforcement of the ceiling, owned by task_04.
package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func TestDeliverStartPrintsTokensNone(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repo, allDeliveryOperations...)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int { return exitOK }
	})
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag)
	if code != 0 || diag.Len() != 0 || out.String() != "Limits: deadline none, retries per item none, concurrency 1, tokens none\n" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
}

func TestDeliverStatusPrintsTheUsageLine(t *testing.T) {
	t.Parallel()
	home, repo, run := seedUsageRun(t, false)
	if err := os.MkdirAll(filepath.Join(repo, "docs", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := writer.CreateDeliveryQueueWithLimits(t.Context(), repo, []string{run.SpecSlug}, store.DeliveryQueueLimits{MaxTokens: 5000000})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.RunID = run.ID
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	// A retry replaces the current item Run but must not erase the old link.
	retry, err := writer.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: "feat/retry", HeadSHA: "retry-head", SpecSlug: run.SpecSlug, Agent: "codex"})
	if err != nil {
		t.Fatal(err)
	}
	tokens := int64(100)
	if err := writer.AppendTokenUsage(t.Context(), store.TokenUsageRecord{RunID: retry.ID, ScopeKind: "task", ScopeID: "task_01", Session: "retry", Basis: "turn", TotalTokens: &tokens}); err != nil {
		t.Fatal(err)
	}
	item.RunID = retry.ID
	if err := writer.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
	want := "Limits: deadline none, retries per item none, concurrency 1, tokens 5000000\nUsage: 9659850 tokens from 3 of 4 prompt(s) across 2 Run(s); cost 3.55 USD from 1 of 4 Agent Session(s)\n"
	if code != 0 || diag.Len() != 0 || !strings.Contains(out.String(), want) {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
}

func TestDeliverStatusPrintsNoRunsRecorded(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	writer, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
	want := implementTestSlug + "\tqueued\t-\t-\nLimits: deadline none, retries per item none, concurrency 1, tokens none\nUsage: no Runs recorded\n"
	if code != 0 || diag.Len() != 0 || out.String() != want {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
}
