package daemon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestVerificationAttemptCancelledAtTheImmediateGrantStartsNoCommand(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gate := newVerificationGate(1)
	verifier := &immediateGrantCountingVerifier{}
	engine := &Engine{deps: Dependencies{Verifier: verifier, Progress: io.Discard}}
	plan := TaskPlan{RunID: "run-immediate-cancellation", verificationGate: gate}
	task := spec.Task{ID: "task_01", Verification: []string{"verify-task_01"}}
	var phases []string
	request := verificationAttemptRequest{
		RunID:       plan.RunID,
		WorkDir:     t.TempDir(),
		ArtifactDir: t.TempDir(),
		BatchNumber: 1,
		WorkItem:    task.ID,
		Attempt:     1,
		Mode:        verificationShared,
		Capacity:    1,
		Commands:    task.Verification,
		Publish: func(_ context.Context, _ string, payload map[string]any) error {
			phase, ok := payload["phase"].(string)
			if !ok {
				t.Fatalf("expected a Verification phase, got %v", payload)
			}
			phases = append(phases, phase)
			if phase == string(runevent.VerificationPhaseWaiting) {
				cancel()
			}
			return nil
		},
	}

	_, err := engine.runTaskVerificationRequest(ctx, plan, task, request)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected cancelled capacity acquisition, got %v", err)
	} else if want := fmt.Sprintf("acquire Verification Capacity for run %q Task %s attempt 1 retry 0: %v", plan.RunID, task.ID, context.Canceled); err.Error() != want {
		t.Errorf("expected acquisition error %q, got %q", want, err)
	}
	if len(verifier.requests) != 0 {
		t.Errorf("expected no Verification command after cancellation, got %d requests: %+v", len(verifier.requests), verifier.requests)
	}
	if want := []string{string(runevent.VerificationPhaseWaiting)}; !slices.Equal(phases, want) {
		t.Errorf("expected only Waiting and no started phase, got %v", phases)
	}
	release, acquired := gate.TryAcquire(verificationShared)
	if !acquired {
		t.Fatal("expected capacity released after cancelled immediate grant")
	}
	release()
}

func TestPreWorkProbeCancelledAtTheImmediateGrantStartsNoProbe(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workDir := t.TempDir()
	runStore, err := store.Open(ctx, t.TempDir())
	if err != nil {
		t.Fatalf("open Run store: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run store: %v", err)
		}
	})
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     workDir,
		LocalBranch: "fix/immediate-grant-cancellation",
		SpecSlug:    taskCycleSlug,
	})
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	gate := newVerificationGate(1)
	verifier := &immediateGrantCountingVerifier{}
	engine := &Engine{deps: Dependencies{
		Verifier: verifier,
		Runs:     &cancelAfterRunStateUpdate{RunStateStore: runStore, cancel: cancel},
	}}
	plan := TaskPlan{
		RunID:                   run.ID,
		WorkDir:                 workDir,
		ArtifactDir:             t.TempDir(),
		VerificationConcurrency: 1,
		verificationGate:        gate,
	}
	task := spec.Task{ID: "task_01", Verification: []string{"verify-task_01"}}

	_, err = engine.verifyTaskPreWork(ctx, plan, task, 1)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected cancelled capacity acquisition, got %v", err)
	} else if want := fmt.Sprintf("acquire Verification Capacity for run %q Task %s attempt 1 retry 0: %v", plan.RunID, task.ID, context.Canceled); err.Error() != want {
		// ProbeCommands also rejects cancellation before calling the verifier.
		// The acquisition error proves rejection happened at the grant.
		t.Errorf("expected acquisition error %q, got %q", want, err)
	}
	if got := runStateForTest(runStore, run.ID); got != store.StateVerifying {
		t.Errorf("expected delegated state update before cancellation, got %q", got)
	}
	if len(verifier.requests) != 0 {
		t.Errorf("expected no pre-work probe after cancellation, got %d requests: %+v", len(verifier.requests), verifier.requests)
	}
	release, acquired := gate.TryAcquire(verificationShared)
	if !acquired {
		t.Fatal("expected capacity released after cancelled immediate grant")
	}
	release()
}

type cancelAfterRunStateUpdate struct {
	RunStateStore
	cancel context.CancelFunc
}

func (runs *cancelAfterRunStateUpdate) UpdateRunState(ctx context.Context, runID string, state string) error {
	if err := runs.RunStateStore.UpdateRunState(ctx, runID, state); err != nil {
		return err
	}
	runs.cancel()
	return nil
}

type immediateGrantCountingVerifier struct {
	requests []VerifyRequest
}

func (verifier *immediateGrantCountingVerifier) Verify(_ context.Context, request VerifyRequest) (VerifyResult, error) {
	verifier.requests = append(verifier.requests, request)
	return VerifyResult{OutputPath: request.OutputPath}, nil
}
