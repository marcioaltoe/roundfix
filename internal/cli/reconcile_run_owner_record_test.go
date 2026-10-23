package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

// Each PID has fixed proof evidence; recording terminations catches any attempt
// to act on the unrelated process that now holds the recorded PID.
type runOwnerRecordProcesses struct {
	proofErrors           map[int]error
	inspected, terminated []int
}

func (p *runOwnerRecordProcesses) InspectTree(_ context.Context, pid int, _ string) ([]int, error) {
	p.inspected = append(p.inspected, pid)
	return nil, p.proofErrors[pid]
}

func (p *runOwnerRecordProcesses) TerminateTreeAndWait(_ context.Context, pid int, _ string) ([]store.TerminationOutcome, error) {
	p.terminated = append(p.terminated, pid)
	return nil, errors.New("owner record must never be terminated")
}

func newRunOwnerRecordFixture(t *testing.T, identity string, proofErr error) (string, store.Run, *runOwnerRecordProcesses) {
	t.Helper()
	home, repo, _ := newReconcileWorkspace(t)
	run := createReconcileMetadataRun(t, home, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: repo, LocalBranch: "ma/widget-flow",
		HeadSHA:  strings.TrimSpace(gitImplementOutput(t, repo, "rev-parse", "HEAD")),
		SpecSlug: "owner-record-spec", Agent: "codex", OwnerPID: 4242, OwnerIdentity: identity,
	}, store.StateFailed)
	processes := &runOwnerRecordProcesses{proofErrors: map[int]error{4242: proofErr}}
	withReconcileProcessController(t, processes)
	return home, run, processes
}

func runOwnerRecordReconcile(t *testing.T, args ...string) (reconcileReport, string) {
	t.Helper()
	var out, diagnostics bytes.Buffer
	code := runCLIContext(t, t.Context(), append([]string{"reconcile"}, args...), &out, &diagnostics)
	if code != exitOK || diagnostics.Len() != 0 {
		t.Fatalf("reconcile = %d, stderr=%q, stdout=%q", code, diagnostics.String(), out.String())
	}
	var report reconcileReport
	if strings.Contains(strings.Join(args, " "), "json") {
		if err := json.Unmarshal(out.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
	}
	return report, out.String()
}

func TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal(t *testing.T) {
	t.Parallel()
	home, run, processes := newRunOwnerRecordFixture(t, "test:old", fmt.Errorf("proof: %w", store.ErrOwnerProcessIdentityMismatch))
	before := readReconcileBytes(t, store.DatabasePath(home))
	_, text := runOwnerRecordReconcile(t)
	proof := fmt.Sprintf("terminal Run %q outcome %q records owner PID 4242, which now has a different start identity: the recorded owner ended", run.ID, run.State)
	for _, want := range []string{
		"Run owner record candidate: Run " + run.ID + " owner PID 4242\n  process-ids: []\n  proof: " + proof + "\n  action: would release the Run owner record with --apply\n  refusal-reason: -\n",
		" processes-applied=0 run-owner-records-released=0 run-branches-applied=0 ",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("text missing %q: %s", want, text)
		}
	}
	dry, raw := runOwnerRecordReconcile(t, "--format=json")
	if len(dry.ProcessCandidates) != 1 || len(dry.PreservedCandidates) != 0 || dry.DebrisSummary.ProcessCandidates != 1 || dry.DebrisSummary.RunOwnerRecordsReleased != 0 {
		t.Fatalf("dry report = %+v", dry)
	}
	candidate := dry.ProcessCandidates[0]
	if candidate.Kind != "runOwnerRecord" || candidate.RunID != run.ID || candidate.Outcome != run.State || candidate.OwnerPID != 4242 || candidate.Proof != proof || candidate.Action != "would release the Run owner record with --apply" || candidate.RefusalReason != "" ||
		!strings.Contains(raw, `"processIds":[]`) || !strings.Contains(raw, `"runOwnerRecordsReleased":0`) {
		t.Fatalf("candidate = %+v; JSON=%s", candidate, raw)
	}
	if !bytes.Equal(before, readReconcileBytes(t, store.DatabasePath(home))) {
		t.Fatal("dry run wrote the Run Database")
	}
	applied, _ := runOwnerRecordReconcile(t, "--apply", "--format=json")
	if applied.DebrisSummary.RunOwnerRecordsReleased != 1 || applied.DebrisSummary.ProcessesApplied != 0 || applied.DebrisSummary.Preserved != 0 || applied.Summary.OperationalFailures != 0 || len(applied.ProcessCandidates) != 1 || applied.ProcessCandidates[0].Action != "released Run owner record; no signal sent" {
		t.Fatalf("apply report = %+v", applied)
	}
	db, err := store.OpenReader(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	got, found, err := db.Run(t.Context(), run.ID)
	if err != nil || !found {
		t.Fatalf("read released Run: %v", err)
	}
	want := run
	want.OwnerPID = nil
	want.OwnerIdentity = ""
	want.OwnerIdentityUnproven = false
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Run after release = %+v, want %+v", got, want)
	}
	events, err := db.RunEventsAfter(t.Context(), run.ID, 0, 100)
	if err != nil || len(events) != 1 || !strings.Contains(string(events[0].Event.Payload), `"reason":"reconcile proved owner PID 4242 now has a different start identity"`) {
		t.Fatalf("release journal = %+v, %v", events, err)
	}
	inspected := len(processes.inspected)
	second, _ := runOwnerRecordReconcile(t, "--apply", "--format=json")
	if len(second.ProcessCandidates) != 0 || len(second.PreservedCandidates) != 0 || second.DebrisSummary.RunOwnerRecordsReleased != 0 || len(processes.inspected) != inspected || len(processes.terminated) != 0 {
		t.Fatalf("second apply = %+v; processes = %+v", second, processes)
	}
	table := &doctorResidueProcessTable{read: func(context.Context, ProcessLineage) ([]store.OwnedProcess, error) {
		t.Fatal("released Run still read a lineage")
		return nil, nil
	}}
	residue, err := Residue(t.Context(), db, table)
	if err != nil || len(residue) != 0 || table.calls != 0 {
		t.Fatalf("residue = %+v, %v; calls=%d", residue, err, table.calls)
	}
	for _, tt := range []struct {
		name                                 string
		missingRun, changedIdentity, badHome bool
	}{
		{name: "changed identity", changedIdentity: true},
		{name: "missing Run", missingRun: true},
		{name: "database open failure", badHome: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			candidate := reconcileDebrisResult{Kind: "runOwnerRecord", RunID: run.ID, OwnerPID: 4242, run: run}
			// Use the original inspected evidence against the already-released row.
			if tt.changedIdentity {
				candidate.run.OwnerIdentity = "test:new"
			}
			if tt.missingRun {
				candidate.RunID = "absent"
			}
			applyHome := home
			if tt.badHome {
				applyHome = store.DatabasePath(home)
			}
			report := reconcileReport{ProcessCandidates: []reconcileDebrisResult{candidate}}
			applyReconcileProcesses(t.Context(), applyHome, &report)
			got := report.ProcessCandidates[0]
			if tt.missingRun || tt.badHome {
				if got.Action != "preserve" || !strings.HasPrefix(got.RefusalReason, "release Run owner record: ") || report.DebrisSummary.Preserved != 1 || report.Summary.OperationalFailures != 1 {
					t.Fatalf("failed release = %+v", report)
				}
			} else if got.Action != "none; Run owner record already released" || report.DebrisSummary.RunOwnerRecordsReleased != 0 || report.Summary.OperationalFailures != 0 {
				t.Fatalf("stale apply = %+v", report)
			}
		})
	}

}

func TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, identity string
		proofErr       error
	}{
		{"unreadable read failure", "test:old", fmt.Errorf("%w: read live owner identity: read denied", store.ErrOwnerIdentityUnreadable)},
		{"uncomparable token", "foreign:old", fmt.Errorf("%w: recorded owner identity is not comparable", store.ErrOwnerIdentityUnreadable)},
		{"bare unproven", "test:old", store.ErrOwnerProcessIdentityUnproven},
		{"unsupported", "test:old", store.ErrOwnerProcessUnsupported},
		{"without identity", "", store.ErrOwnerProcessIdentityMismatch},
	} {
		t.Run(tt.name, func(t *testing.T) {
			home, run, processes := newRunOwnerRecordFixture(t, tt.identity, tt.proofErr)
			before := readReconcileBytes(t, store.DatabasePath(home))
			reason := "owned process tree cannot be proven: " + tt.proofErr.Error()
			if tt.identity == "" {
				reason = "recorded owner process identity is unproven"
			}
			for _, args := range [][]string{{"--format=json"}, {"--apply", "--format=json"}} {
				report, _ := runOwnerRecordReconcile(t, args...)
				if len(report.ProcessCandidates) != 0 || len(report.PreservedCandidates) != 1 || report.DebrisSummary.Preserved != 1 || report.DebrisSummary.RunOwnerRecordsReleased != 0 {
					t.Fatalf("report = %+v", report)
				}
				candidate := report.PreservedCandidates[0]
				if candidate.Kind != "process" || candidate.RunID != run.ID || candidate.OwnerPID != 4242 || candidate.Action != "preserve" || candidate.RefusalReason != reason {
					t.Fatalf("preserved = %+v, want reason %q", candidate, reason)
				}
				if !bytes.Equal(before, readReconcileBytes(t, store.DatabasePath(home))) {
					t.Fatal("preserved Run Database changed")
				}
			}
			if len(processes.terminated) != 0 {
				t.Fatalf("termination requested: %v", processes.terminated)
			}
			if tt.identity == "" && len(processes.inspected) != 0 {
				t.Fatal("identity-less owner was inspected")
			}
		})
	}
}
