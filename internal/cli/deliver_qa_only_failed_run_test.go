// Suite: Delivery Retry recovery of a QA-Only Failed Run.
// Invariant: only Delivery Retry accepts the completed Tasks; the QA gate never carries.
// Boundary IN: production carry-forward, disposable repositories and Run Database.
// Boundary OUT: Agent execution and Task settlement represented by committed evidence.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestAFailedImplementChildNamesTheRunItCreated(t *testing.T) {
	t.Parallel()
	for _, created := range []bool{true, false} {
		name := "no_run_created_returns_none"
		if created {
			name = "failed_run_is_returned_with_the_failure"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			before := &store.Run{ID: "run-earlier"}
			after := before
			want := ""
			if created {
				after = &store.Run{ID: "run-failed", State: store.StateFailed}
				want = after.ID
			}
			result, err := deliveryRunResult(roundfixCommandResult{exitCode: 1}, "", before, after)
			if result.RunID != want || err == nil || err.Error() != "roundfix implement failed with exit code 1: exit code 1" {
				t.Fatalf("result=%+v error=%v, want Run %q and unchanged failure", result, err, want)
			}
		})
	}
}

func TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"carries_the_settled_tasks", "another_task_not_completed_carries_nothing", "run_branch_gone_carries_nothing", "unresolved_run_carries_as_before", "preflight_reports_no_candidate", "reconcile_keeps_refusing_the_failed_run", "unloadable_committed_graph_carries_nothing", "completed_gate_carries_nothing", "moved_input_refuses_the_whole_set"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			seeds := []implementSeed{}
			ids := []string{}
			for i := 1; i <= 5; i++ {
				id := fmt.Sprintf("task_%02d", i)
				ids = append(ids, id)
				seeds = append(seeds, implementSeed{id: id})
			}
			graphSeeds := append(append([]implementSeed{}, seeds...), implementQAGateSeed("pending", ids...))
			home, repo := newImplementWorkspace(t, graphSeeds)
			workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
			state := store.StateFailed
			if name == "unresolved_run_carries_as_before" {
				state = store.StateUnresolved
			}
			settled := seeds
			if name == "another_task_not_completed_carries_nothing" {
				settled = seeds[:4]
			}
			fixture := createDeliveryRetryRun(t, workflow, repo, "ma/widget-flow", settled, state)
			if name == "carries_the_settled_tasks" {
				// Even stale settlement evidence cannot make the gate a candidate.
				_, err := workflow.store.AppendRunEvents(t.Context(), []runevent.RunEvent{
					{RunID: fixture.run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonVerification, ReviewIssue: "task_qa", Payload: []byte(`{"attempt":1,"phase":"verdict","task":"task_qa","verdict":"passed"}`)},
					{RunID: fixture.run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, ReviewIssue: "task_qa", Payload: []byte(`{"task":"task_qa","phase":"settled","status":"completed"}`)},
				})
				if err != nil {
					t.Fatal(err)
				}
			}
			if name == "completed_gate_carries_nothing" {
				if err := spec.SetStatus(implementTaskPath(fixture.ref.Path, "task_qa"), spec.StatusCompleted); err != nil {
					t.Fatal(err)
				}
				gitImplement(t, fixture.ref.Path, "add", ".")
				gitImplement(t, fixture.ref.Path, "commit", "-m", "complete gate")
			}
			if name == "unloadable_committed_graph_carries_nothing" {
				mustWrite(t, filepath.Join(fixture.ref.Path, "docs", "specs", implementTestSlug, "_tasks.md"), "invalid graph\n")
				gitImplement(t, fixture.ref.Path, "add", ".")
				gitImplement(t, fixture.ref.Path, "commit", "-m", "break graph")
			}
			if name == "moved_input_refuses_the_whole_set" {
				path := filepath.Join(repo, "docs", "specs", implementTestSlug, "_prd.md")
				mustWrite(t, path, mustRead(t, path)+"\nChanged binding input.\n")
				gitImplement(t, repo, "add", ".")
				gitImplement(t, repo, "commit", "-m", "amend binding input")
			}
			if name == "run_branch_gone_carries_nothing" {
				gitImplement(t, repo, "worktree", "remove", "--force", fixture.ref.Path)
				gitImplement(t, repo, "branch", "-D", fixture.ref.Branch)
			}
			if name == "preflight_reports_no_candidate" {
				candidates, err := inspectSpecCarryForwards(t.Context(), workflow.store, repo, roundconfig.SpecsRoot{Path: filepath.Join(repo, "docs", "specs")}, implementTestSlug)
				if err != nil || len(candidates) != 0 {
					t.Fatalf("Preflight candidates=%+v err=%v", candidates, err)
				}
				return
			}
			if name == "reconcile_keeps_refusing_the_failed_run" {
				var out, diagnostic bytes.Buffer
				code := runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--carry-forward"}, &out, &diagnostic)
				if code != exitPreflight || !strings.Contains(diagnostic.String(), "has outcome Failed; carry-forward accepts outcomes Stopped and Unresolved and BudgetExceeded") {
					t.Fatalf("reconcile code=%d stderr=%q", code, diagnostic.String())
				}
				return
			}
			// Uncommitted Run Worktree status must not classify a Failed Run.
			if name == "another_task_not_completed_carries_nothing" {
				if err := spec.SetStatus(implementTaskPath(fixture.ref.Path, "task_05"), spec.StatusCompleted); err != nil {
					t.Fatal(err)
				}
			}
			result, err := workflow.CarryForward(t.Context(), repo, implementTestSlug, "ma/widget-flow", fixture.run.ID)
			if name == "moved_input_refuses_the_whole_set" {
				var refusal deliveryCarryForwardRefusal
				if !errors.As(err, &refusal) || !strings.Contains(err.Error(), "declared input(s) moved") || strings.Contains(refusal.NextAction(), "roundfix reconcile") || !strings.Contains(refusal.NextAction(), "roundfix deliver retry") {
					t.Fatalf("refusal=%v", err)
				}
				for _, id := range ids {
					if _, err := os.Stat(filepath.Join(repo, "src", id+".txt")); !os.IsNotExist(err) {
						t.Fatalf("refused set carried %s: %v", id, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			carry := name == "carries_the_settled_tasks" || name == "unresolved_run_carries_as_before"
			for _, seed := range graphSeeds {
				content := mustRead(t, implementTaskPath(repo, seed.id))
				want := "pending"
				if carry && seed.id != "task_qa" {
					want = "completed"
				}
				if !strings.Contains(content, "status: "+want) {
					t.Fatalf("%s: want %s: %s", seed.id, want, content)
				}
				if carry && seed.id != "task_qa" && (!strings.Contains(content, fixture.run.ID) || !strings.Contains(content, fixture.commits[seed.id])) {
					t.Fatalf("missing provenance: %s", content)
				}
			}
			var output bytes.Buffer
			printDeliverRetryResult(&output, implementTestSlug, delivery.RetryResult{CarriedFrom: result})
			if carry {
				line := "Carried forward from Run " + fixture.run.ID + ": " + strings.Join(ids, ", ") + "\n"
				if strings.Count(output.String(), "Carried forward from Run ") != 1 || !strings.Contains(output.String(), line) {
					t.Fatalf("output=%q", output.String())
				}
			} else if len(result.Runs) != 0 {
				t.Fatalf("unexpected carry: %+v", result)
			}
		})
	}
}
