package cli

// Suite: digest-confirmed Dirty Discard through the public CLI.
// Invariant: confirmation binds the complete offered set; recording precedes removal.
// Boundary IN: temporary homes, Run Databases, and real Git repositories.
// Boundary OUT: live process discovery internals and provider services.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/store"
)

func newReconcileDirtyDiscardFixture(t *testing.T) reconcileMergedFixture {
	t.Helper()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{
		archivedSlug: reconcileMergedSpecSlug, taskStatus: "completed", includeQA: true,
	})
	archivePRD := filepath.Join("docs", "history", "specs", reconcileMergedSpecSlug, "_prd.md")
	mustWrite(t, filepath.Join(fixture.repoDir, archivePRD), "archived planning\n")
	gitImplement(t, fixture.repoDir, "add", filepath.ToSlash(archivePRD))
	gitImplement(t, fixture.repoDir, "commit", "-m", "docs: archive delivered planning")
	mustWrite(t, filepath.Join(fixture.ref.Path, "unrelated.txt"), "discard\n")
	return fixture
}

func TestReconcileDryRunListsTheDirtyDiscardPathsAndDigest(t *testing.T) {
	t.Parallel()
	fixture := newReconcileDirtyDiscardFixture(t)
	report := runReconcileMergedCommand(t, fixture)
	result := reconcileMergedResult(t, report, fixture.run.ID)
	digest := report.DiscardDirtyDigest
	wantCommand := "roundfix reconcile " + fixture.run.ID + " --discard-dirty " + digest
	if !slices.Equal(result.DiscardPaths, []string{"unrelated.txt"}) || !validDirtyDiscardDigest(digest) || report.DiscardDirtyCommand != wantCommand || result.Action != "would discard 1 uncommitted path(s) with --discard-dirty" {
		t.Fatalf("dry-run report = %+v", report)
	}
	if again := runReconcileMergedCommand(t, fixture); again.DiscardDirtyDigest != digest {
		t.Fatalf("digest changed: %q -> %q", digest, again.DiscardDirtyDigest)
	}
	// Assert the wire fields, not only a struct that would ignore missing JSON.
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--format=json"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("JSON exit=%d stderr=%q", code, stderr.String())
	}
	for _, field := range []string{`"discardPaths":["unrelated.txt"]`, `"discardDirtyDigest":"` + digest + `"`, `"discardDirtyCommand":"` + wantCommand + `"`} {
		if !strings.Contains(stdout.String(), field) {
			t.Fatalf("missing %q in JSON %s", field, stdout.String())
		}
	}
	stdout.Reset()
	code = runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("text exit=%d stderr=%q", code, stderr.String())
	}
	for _, line := range []string{"  discard-path: unrelated.txt\n", "Discard dirty leftovers with: " + wantCommand + "\n"} {
		if !strings.Contains(stdout.String(), line) {
			t.Fatalf("missing %q in %s", line, stdout.String())
		}
	}
	if !strings.Contains(stdout.String(), "Apply with: roundfix reconcile "+fixture.run.ID+" --apply\nDiscard dirty leftovers with:") {
		t.Fatalf("confirmation did not follow Apply with: %s", stdout.String())
	}
}

func TestReconcileDiscardDirtyRecordsThenRemovesTheConfirmedRuns(t *testing.T) {
	t.Parallel()
	fixture := newReconcileDirtyDiscardFixture(t)
	before := runReconcileMergedCommand(t, fixture)
	beforeDatabase := readReconcileBytes(t, store.DatabasePath(fixture.homeDir))
	report := runReconcileMergedCommand(t, fixture, "--discard-dirty="+before.DiscardDirtyDigest)
	result := reconcileMergedResult(t, report, fixture.run.ID)
	path := filepath.Join(builtinArtifactDirForRepo(t, fixture.repoDir), "runs", fixture.run.ID, "branch-disposition.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		SchemaVersion string `json:"schemaVersion"`
		daemon.BranchDisposition
	}
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != daemon.BranchDispositionRecordSchemaVersion || record.Superseded || record.DiscardDigest != before.DiscardDirtyDigest || !slices.Equal(record.DiscardedPaths, []string{"unrelated.txt"}) || record.RunHead != result.RunHead || record.Worktree != fixture.ref.Path || record.Reason != "dirty leftovers discarded by a confirmed Dirty Discard" || record.RefusalReason != before.Results[0].Evidence || len(record.Commits) != 0 || len(record.ChangedFiles) != 0 {
		t.Fatalf("discard record = %+v", record)
	}
	if report.Mode != "discard-dirty" || result.Action != "discarded" || result.RefusalReason != "" || report.Summary.Applied != 1 || !strings.HasSuffix(result.Evidence, "; branch record: "+path) {
		t.Fatalf("discard report = %+v", report)
	}
	if got := readReconcileBytes(t, store.DatabasePath(fixture.homeDir)); !bytes.Equal(got, beforeDatabase) {
		t.Fatal("Dirty Discard wrote to the Run Database")
	}
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
	after := runReconcileMergedCommand(t, fixture)
	if after.Results[0].Classification != "released" {
		t.Fatalf("next report = %+v", after)
	}
}

func TestReconcileDiscardDirtyRefusesAChangedDigestWithoutChangingAnything(t *testing.T) {
	t.Parallel()
	fixture := newReconcileDirtyDiscardFixture(t)
	before := runReconcileMergedCommand(t, fixture)
	mustWrite(t, filepath.Join(fixture.ref.Path, "another.txt"), "keep\n")
	current := runReconcileMergedCommand(t, fixture)
	message := fmt.Sprintf("the Dirty Discard set changed: confirmed digest %s, current digest %s; run roundfix reconcile and confirm the digest it prints", before.DiscardDirtyDigest, current.DiscardDirtyDigest)
	assertDirtyDiscardTranscript(t, []string{"reconcile", fixture.run.ID, "--discard-dirty", before.DiscardDirtyDigest}, message)
	assertNoDirtyDiscardRecord(t, fixture)
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReconcileDiscardDirtyRefusesWhenNoRunIsOffered(t *testing.T) {
	t.Parallel()
	fixture := newReconcileMergedFixture(t, reconcileMergedFixtureOptions{taskStatus: "completed"})
	mustWrite(t, filepath.Join(fixture.ref.Path, "unrelated.txt"), "keep\n")
	assertDirtyDiscardTranscript(t, []string{"reconcile", fixture.run.ID, "--discard-dirty", strings.Repeat("0", sha256.Size*2)}, "no ended Run has dirty leftovers a Dirty Discard can remove; run roundfix reconcile to see why each dirty Run is kept")
	assertNoDirtyDiscardRecord(t, fixture)
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestADirtyDiscardIsNotOfferedForARunWithALiveProcessTree(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"processCandidates", "preservedCandidates", "runOwnerRecord"} {
		t.Run(name, func(t *testing.T) {
			report := reconcileReport{Results: []reconcileResult{{RunID: "ended", DiscardPaths: []string{"unrelated.txt"}}}}
			candidate := reconcileDebrisResult{Kind: "process", RunID: "ended"}
			if name == "runOwnerRecord" {
				candidate.Kind = "runOwnerRecord"
			}
			if name == "preservedCandidates" {
				report.PreservedCandidates = []reconcileDebrisResult{candidate}
			} else {
				report.ProcessCandidates = []reconcileDebrisResult{candidate}
			}
			withholdDirtyDiscardFromLiveRuns(&report)
			if name == "runOwnerRecord" {
				if len(report.Results[0].DiscardPaths) == 0 || dirtyDiscardDigest(report.Results) == "" {
					t.Fatal("owner record withheld discard")
				}
			} else if len(report.Results[0].DiscardPaths) != 0 || report.Results[0].RefusalReason != "a live process tree is attributed to this Run, so its dirty leftovers are not offered for discard" || dirtyDiscardDigest(report.Results) != "" {
				t.Fatalf("live Run report = %+v", report)
			}
		})
	}
}

func TestReconcileDiscardDirtyNeedsItsDigestTranscript(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	assertDirtyDiscardTranscript(t, []string{"reconcile", "--discard-dirty"}, "flag needs an argument: --discard-dirty")
}

func TestReconcileDiscardDirtyRefusesAMalformedDigestTranscript(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	for _, digest := range []string{"abc", "", strings.Repeat("A", 64), strings.Repeat("g", 64), strings.Repeat("a", 63), strings.Repeat("a", 65)} {
		t.Run(fmt.Sprintf("digest-%q", digest), func(t *testing.T) {
			args := []string{"reconcile", "--discard-dirty", digest}
			if digest == "" {
				args = []string{"reconcile", "--discard-dirty="}
			}
			assertDirtyDiscardTranscript(t, args, "--discard-dirty needs the 64-character digest that roundfix reconcile prints")
		})
	}
}

func TestReconcileMutationSwitchesExcludeEachOtherTranscript(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	switches := [][]string{{"--apply"}, {"--discard-superseded"}, {"--discard-dirty", strings.Repeat("0", 64)}, {"--carry-forward"}}
	for i := range switches {
		for j := i + 1; j < len(switches); j++ {
			t.Run(switches[i][0]+"+"+switches[j][0], func(t *testing.T) {
				args := append([]string{"reconcile"}, switches[i]...)
				args = append(args, switches[j]...)
				assertDirtyDiscardTranscript(t, args, "--apply, --discard-superseded, --discard-dirty, and --carry-forward are mutually exclusive")
			})
		}
	}
}

func TestDirtyDiscardDigestBindsTheSortedWholeSet(t *testing.T) {
	t.Parallel()
	results := []reconcileResult{
		{RunID: "b", RunHead: "b-head", Worktree: "b-tree", DiscardPaths: []string{"z", "a"}},
		{RunID: "not-offered"},
		{RunID: "a", RunHead: "a-head", Worktree: "a-tree", DiscardPaths: []string{"x"}},
	}
	sum := sha256.Sum256([]byte("a\x00a-head\x00a-tree\x00x\x00\x00b\x00b-head\x00b-tree\x00a\x00z\x00\x00"))
	want := hex.EncodeToString(sum[:])
	if got := dirtyDiscardDigest(results); got != want {
		t.Fatalf("digest = %s, want %s", got, want)
	}
	if !slices.Equal(results[0].DiscardPaths, []string{"z", "a"}) || results[0].RunID != "b" {
		t.Fatal("digest reordered report")
	}
}

func TestReconcileDirtyDiscardPreservesWorkWhenTheRecordCannotBeWritten(t *testing.T) {
	t.Parallel()
	fixture := newReconcileDirtyDiscardFixture(t)
	report := runReconcileMergedCommand(t, fixture)
	artifactRoot := builtinArtifactDirForRepo(t, fixture.repoDir)
	if err := os.MkdirAll(artifactRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(artifactRoot, "runs"), "not a directory\n")
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"reconcile", fixture.run.ID, "--discard-dirty", report.DiscardDirtyDigest, "--format=json"}, &stdout, &stderr)
	if code != exitRunFailed {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	var failed reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &failed); err != nil {
		t.Fatal(err)
	}
	if failed.Summary.OperationalFailures != 1 || failed.Summary.Applied != 0 || failed.Results[0].Action != "preserve" || !strings.Contains(failed.Results[0].RefusalReason, "write branch disposition record directory") {
		t.Fatalf("record failure report = %+v", failed)
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func assertDirtyDiscardTranscript(t *testing.T, args []string, message string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
	want := "roundfix: reconcile failed: " + message + "\nRun 'roundfix reconcile --help' for usage.\n"
	if code != exitPreflight || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("transcript exit=%d stdout=%q stderr=%q; want exit=2 stdout=empty stderr=%q", code, stdout.String(), stderr.String(), want)
	}
}

func assertNoDirtyDiscardRecord(t *testing.T, fixture reconcileMergedFixture) {
	t.Helper()
	path := filepath.Join(builtinArtifactDirForRepo(t, fixture.repoDir), "runs", fixture.run.ID, "branch-disposition.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("record exists or cannot be inspected: %v", err)
	}
}

func TestReconcileWithoutADirtyDiscardKeepsItsWireBytes(t *testing.T) {
	t.Parallel()
	report := reconcileReport{
		SchemaVersion: reconcileSchemaVersion, Mode: "dry-run", Repository: "repo",
		ApplyCommand: "roundfix reconcile --apply",
		Results:      []reconcileResult{{RunID: "kept", Outcome: "Stopped", Classification: "dirty", Evidence: "kept reason", Action: "preserve", RefusalReason: "kept reason"}},
		Summary:      reconcileSummary{Total: 1, Dirty: 1, Preserved: 1},
	}
	report.Runs = report.Results
	const runJSON = `{"runId":"kept","outcome":"Stopped","classification":"dirty","runBranch":"","runHead":"","targetBranch":"","targetHead":"","worktree":"","supersedingReport":"","evidence":"kept reason","action":"preserve","refusalReason":"kept reason"}`
	const wantJSON = `{"schemaVersion":"roundfix-reconcile/v1","mode":"dry-run","repository":"repo","applyCommand":"roundfix reconcile --apply","results":[` + runJSON + `],"runs":[` + runJSON + `],"summary":{"total":1,"safe":0,"superseded":0,"unintegrated":0,"dirty":1,"unknown":0,"released":0,"applied":0,"preserved":1,"operationalFailures":0},"processCandidates":null,"runBranchCandidates":null,"stagingCandidates":null,"preservedCandidates":null,"debrisSummary":{"processCandidates":0,"runBranchCandidates":0,"stagingCandidates":0,"preserved":0,"processesApplied":0,"runOwnerRecordsReleased":0,"runBranchesApplied":0,"stagingApplied":0,"itemBranchCandidates":0,"itemBranchesApplied":0},"carryForwards":null}` + "\n"
	const wantText = `Repository: repo
Mode: dry-run
Run: kept
  outcome: Stopped
  classification: dirty
  run-branch: -
  run-head: -
  target-branch: -
  target-head: -
  worktree: -
  superseding-report: -
  evidence: kept reason
  action: preserve
  refusal-reason: kept reason
Summary: total=1 safe=0 superseded=0 unintegrated=0 dirty=1 unknown=0 released=0 applied=0 preserved=1 operational-failures=0
Debris summary: process-candidates=0 run-branch-candidates=0 staging-candidates=0 staging-applied=0 preserved=0 processes-applied=0 run-owner-records-released=0 run-branches-applied=0 item-branch-candidates=0 item-branches-applied=0
Apply with: roundfix reconcile --apply
`
	for _, tc := range []struct{ format, want string }{{"text", wantText}, {"json", wantJSON}} {
		t.Run(tc.format, func(t *testing.T) {
			var output bytes.Buffer
			if err := printReconcileReport(&output, tc.format, report); err != nil {
				t.Fatal(err)
			}
			if output.String() != tc.want {
				t.Fatalf("wire bytes changed:\ngot %q\nwant %q", output.String(), tc.want)
			}
		})
	}
}
