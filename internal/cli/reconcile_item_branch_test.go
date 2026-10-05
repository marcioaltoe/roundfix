package cli

// Suite: full-scan item branch reporting and explicit apply through reconcile.
// Invariant: full scans expose item proof; Run-ID scans omit the collection;
// live Queue items and unproven or dirty worktrees stay intact.
// Boundary IN: public command runner, real local Git and Delivery Queue Store.
// Boundary OUT: network and Branch Disposition persistence.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func newReconcileItemFixture(t *testing.T) (reconcileMergedFixture, runworktree.ItemRef) {
	t.Helper()
	f := newReconcileMergeEvidenceFixture(t, false)
	location := filepath.Join(f.homeDir, ".roundfix", "worktrees")
	branch := "roundfix/deliver-" + reconcileMergedSpecSlug + "-0123456789abcdef"
	ref, err := runworktree.ItemRefFor(f.repoDir, location, branch)
	if err != nil {
		t.Fatal(err)
	}
	gitImplement(t, f.repoDir, "worktree", "add", "-b", branch, ref.Path, f.taskCommit)
	return f, ref
}

func itemReconcileOutput(t *testing.T, args ...string) []byte {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), append([]string{"reconcile"}, args...), &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("reconcile exit=%d stderr=%s stdout=%s", code, &stderr, &stdout)
	}
	return stdout.Bytes()
}

func itemReconcileReport(t *testing.T, args ...string) reconcileReport {
	t.Helper()
	var report reconcileReport
	if err := json.Unmarshal(itemReconcileOutput(t, append(args, "--format=json")...), &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec(t *testing.T) {
	t.Parallel()
	f, ref := newReconcileItemFixture(t)
	text := string(itemReconcileOutput(t))
	for _, want := range []string{
		"Item branch candidate: " + ref.Branch + " (Spec " + reconcileMergedSpecSlug + ")",
		"  head: " + f.taskCommit,
		"  worktree: ",
		"  proof: item branch is superseded by the delivery of Spec",
		"  action: would reclaim with --apply",
		"  refusal-reason: -",
		"item-branch-candidates=1 item-branches-applied=0\n",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
	report := itemReconcileReport(t)
	if report.ItemBranchCandidates == nil || len(*report.ItemBranchCandidates) != 1 {
		t.Fatalf("report = %+v", report)
	}
	candidate := (*report.ItemBranchCandidates)[0]
	if candidate.ItemBranch != ref.Branch || candidate.Head != f.taskCommit || candidate.Kind != "itemBranch" || candidate.SpecSlug != reconcileMergedSpecSlug || candidate.Action != "would reclaim with --apply" {
		t.Fatalf("candidate = %+v", candidate)
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, f.repoDir, ref.Branch, true)
	// Existing debris entries must not acquire the item-only fields.
	raw := itemReconcileOutput(t, "--format=json")
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"runBranchCandidates", "processCandidates", "stagingCandidates"} {
		var entries []map[string]json.RawMessage
		if err := json.Unmarshal(envelope[field], &entries); err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			for _, key := range []string{"itemBranch", "head"} {
				if _, exists := entry[key]; exists {
					t.Fatalf("%s acquired %s", field, key)
				}
			}
		}
	}
	applied := itemReconcileReport(t, "--apply")
	if applied.DebrisSummary.ItemBranchesApplied != 1 || (*applied.ItemBranchCandidates)[0].Action != "released" {
		t.Fatalf("apply = %+v", applied)
	}
	assertReconcilePathState(t, ref.Path, false)
	assertReconcileBranchState(t, f.repoDir, ref.Branch, false)
	again := itemReconcileReport(t)
	if again.ItemBranchCandidates == nil || len(*again.ItemBranchCandidates) != 0 {
		t.Fatalf("empty full scan = %+v", again)
	}
}

func seedLiveReconcileItem(t *testing.T, f reconcileMergedFixture, ref runworktree.ItemRef, stage store.DeliveryStage) {
	t.Helper()
	ctx := context.Background()
	reader, err := store.Open(ctx, f.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = reader.Close() }()
	queue, err := reader.CreateDeliveryQueue(ctx, f.repoDir, []string{reconcileMergedSpecSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Branch, item.Worktree, item.Stage = ref.Branch, ref.Path, stage
	if err := reader.UpdateDeliveryQueueItem(ctx, f.repoDir, item); err != nil {
		t.Fatal(err)
	}
}

func TestReconcileKeepsTheItemBranchOfALiveQueueItem(t *testing.T) {
	t.Parallel()
	for _, stage := range []store.DeliveryStage{store.DeliveryStageQueued, store.DeliveryStageRunning, store.DeliveryStageParked} {
		t.Run(string(stage), func(t *testing.T) {
			f, ref := newReconcileItemFixture(t)
			seedLiveReconcileItem(t, f, ref, stage)
			report := itemReconcileReport(t, "--apply")
			if report.ItemBranchCandidates == nil || len(*report.ItemBranchCandidates) != 0 {
				t.Fatalf("live candidates = %+v", report)
			}
			found := false
			for _, candidate := range report.PreservedCandidates {
				if candidate.Kind == "itemBranch" && candidate.ItemBranch == ref.Branch && candidate.RefusalReason == fmt.Sprintf("item branch belongs to live Delivery Queue item %q", reconcileMergedSpecSlug) {
					found = true
				}
			}
			if !found {
				t.Fatalf("live refusal missing: %+v", report.PreservedCandidates)
			}
			assertReconcilePathState(t, ref.Path, true)
			assertReconcileBranchState(t, f.repoDir, ref.Branch, true)
		})
	}
}

func TestReconcileWithARunIDOmitsItemBranchCandidates(t *testing.T) {
	t.Parallel()
	f, ref := newReconcileItemFixture(t)
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(itemReconcileOutput(t, f.run.ID, "--format=json"), &envelope); err != nil {
		t.Fatal(err)
	}
	if _, exists := envelope["itemBranchCandidates"]; exists {
		t.Fatal("Run-ID scan exposed item branches")
	}
	itemReconcileOutput(t, f.run.ID, "--apply")
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, f.repoDir, ref.Branch, true)
}

func TestReconcileItemBranchSafety(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"dirty", "unmerged", "wrong-path", "live-after-proof", "merged-item", "moved-head"} {
		t.Run(change, func(t *testing.T) {
			f, ref := newReconcileItemFixture(t)
			if change == "live-after-proof" || change == "moved-head" {
				selection, err := loadReconcileRuns(context.Background(), f.homeDir, f.repoDir, f.repoDir, "")
				if err != nil {
					t.Fatal(err)
				}
				selection.worktreeLocation = filepath.Join(f.homeDir, ".roundfix", "worktrees")
				report := inspectReconcileRuns(context.Background(), f.repoDir, reconcileOptions{}, selection)
				if change == "live-after-proof" {
					seedLiveReconcileItem(t, f, ref, store.DeliveryStageParked)
				} else {
					gitImplement(t, ref.Path, "commit", "--allow-empty", "-m", "fix: move item head")
				}
				applyReconcileItemBranches(context.Background(), f.homeDir, &report)
				if report.DebrisSummary.ItemBranchesApplied != 0 || report.Summary.OperationalFailures != 1 {
					t.Fatalf("stale apply = %+v", report)
				}
			} else {
				switch change {
				case "dirty":
					mustWrite(t, filepath.Join(ref.Path, "untracked.txt"), "dirty\n")
				case "unmerged":
					gitImplement(t, f.repoDir, "rm", "--", "docs/history/specs/"+reconcileMergedSpecSlug+"/_prd.md")
					gitImplement(t, f.repoDir, "commit", "-m", "docs: remove archive proof")
				case "wrong-path":
					moved := filepath.Join(t.TempDir(), "moved")
					gitImplement(t, f.repoDir, "worktree", "move", ref.Path, moved)
					ref.Path = moved
				case "merged-item":
					seedLiveReconcileItem(t, f, ref, store.DeliveryStageMerged)
				}
				report := itemReconcileReport(t, "--apply")
				if change == "merged-item" {
					if report.DebrisSummary.ItemBranchesApplied != 1 {
						t.Fatalf("merged item = %+v", report)
					}
					assertReconcileBranchState(t, f.repoDir, ref.Branch, false)
					return
				}
				found := false
				for _, candidate := range report.PreservedCandidates {
					if candidate.Kind == "itemBranch" && candidate.ItemBranch == ref.Branch && candidate.RefusalReason != "" {
						found = true
					}
				}
				if !found {
					t.Fatalf("missing refusal = %+v", report)
				}
			}
			assertReconcilePathState(t, ref.Path, true)
			assertReconcileBranchState(t, f.repoDir, ref.Branch, true)
		})
	}
}

func TestReconcileEmptyFullScanIncludesItemBranches(t *testing.T) {
	t.Parallel()
	newReconcileWorkspace(t)
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(itemReconcileOutput(t, "--format=json"), &envelope); err != nil {
		t.Fatal(err)
	}
	if string(envelope["itemBranchCandidates"]) != "[]" {
		t.Fatalf("empty collection = %s", envelope["itemBranchCandidates"])
	}
}
