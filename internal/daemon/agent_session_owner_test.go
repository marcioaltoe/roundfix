// Suite: Agent Session selection ownership.
// Invariant: a Fallback Selection remains eligible only until the first Agent output.
// Boundary IN: Session preparation, first-output events, selection failures, and fallback activation.
// Boundary OUT: acpx stream decoding, owned by internal/agent/acpx_runner_test.go.
package daemon

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
)

func TestFallbackEligibilitySurvivesASelectionFailure(t *testing.T) {
	t.Parallel()

	sink := &captureEventSink{}
	runner := &fallbackBoundaryRunner{
		runErrByModel: map[string]error{
			"preferred-model": &agent.SelectionFailure{Runtime: "codex", Reason: "usage limit exhausted"},
		},
		emitOutputByModel: map[string]bool{"fallback-model": true},
	}
	owner := fallbackBoundaryOwner(t, sink, runner)

	result, err := owner.Run(context.Background(), agent.ExecuteRequest{
		RunID:  "run-fallback-boundary",
		Prompt: "do the work",
	})

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Output != "fallback output" {
		t.Fatalf("Run() output = %q, want fallback output", result.Output)
	}
	if got := strings.Join(runner.runModels(), ","); got != "preferred-model,fallback-model" {
		t.Fatalf("run models = %q, want preferred-model,fallback-model", got)
	}
	if got := countAgentStatusEvents(sink, agent.AgentSelectionFailedStatus); got != 1 {
		t.Fatalf("selection-failed statuses = %d, want 1", got)
	}
	if got := countAgentStatusEvents(sink, agent.AgentWorkStartedStatus); got != 1 {
		t.Fatalf("work-started statuses = %d, want 1 from fallback output", got)
	}
	if events := eventsOfKind(sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 1 {
		t.Fatalf("fallback events = %+v, want one", events)
	}
}

func TestFallbackEligibilityEndsAfterAnyAgentOutput(t *testing.T) {
	t.Parallel()

	sink := &captureEventSink{}
	runner := &fallbackBoundaryRunner{
		runErrByModel: map[string]error{
			"preferred-model": &agent.SelectionFailure{Runtime: "codex", Reason: "adapter exited after output"},
		},
		emitOutputByModel: map[string]bool{"preferred-model": true},
	}
	owner := fallbackBoundaryOwner(t, sink, runner)

	_, err := owner.Run(context.Background(), agent.ExecuteRequest{
		RunID:  "run-fallback-boundary",
		Prompt: "do the work",
	})

	var selectionErr *agent.SelectionFailure
	if !errors.As(err, &selectionErr) {
		t.Fatalf("Run() error = %T %v, want original *agent.SelectionFailure", err, err)
	}
	if got := strings.Join(runner.runModels(), ","); got != "preferred-model" {
		t.Fatalf("run models = %q, want only preferred-model", got)
	}
	if got := countAgentStatusEvents(sink, agent.AgentWorkStartedStatus); got != 1 {
		t.Fatalf("work-started statuses = %d, want 1", got)
	}
	if events := eventsOfKind(sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 0 {
		t.Fatalf("fallback events = %+v, want none after Agent output", events)
	}
}

func TestFallbackEligibilitySurvivesAdapterStartFailure(t *testing.T) {
	t.Parallel()

	sink := &captureEventSink{}
	runner := &fallbackBoundaryRunner{
		prepareErrByModel: map[string]error{
			"preferred-model": agent.AdapterProbeError{
				Command:    "missing-adapter",
				Executable: "missing-adapter",
				Err:        errors.New("executable file not found"),
			},
		},
		emitOutputByModel: map[string]bool{"fallback-model": true},
	}
	owner := fallbackBoundaryOwner(t, sink, runner)

	_, err := owner.Run(context.Background(), agent.ExecuteRequest{
		RunID:  "run-fallback-boundary",
		Prompt: "do the work",
	})

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got := strings.Join(runner.preparedModels(), ","); got != "preferred-model,fallback-model" {
		t.Fatalf("prepared models = %q, want preferred-model,fallback-model", got)
	}
	if got := strings.Join(runner.runModels(), ","); got != "fallback-model" {
		t.Fatalf("run models = %q, want only fallback-model", got)
	}
	if got := countAgentStatusEvents(sink, agent.AgentSelectionFailedStatus); got != 1 {
		t.Fatalf("selection-failed statuses = %d, want 1", got)
	}
	if events := eventsOfKind(sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 1 {
		t.Fatalf("fallback events = %+v, want one", events)
	}
}

func fallbackBoundaryOwner(t *testing.T, sink *captureEventSink, runner agent.Runner) *agentSessionOwner {
	t.Helper()
	engine := &Engine{deps: Dependencies{
		Runner:   runner,
		Sink:     sink,
		Progress: &bytes.Buffer{},
		Now: func() time.Time {
			return time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC)
		},
	}}
	owner, err := engine.agentSessionOwner(agentSelectionOwnerConfig{
		Profiles: selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
			roundconfig.CategoryBackend: selectionProfileForTest(
				selectionForTest("codex", "preferred-model", "high"),
				selectionForTest("claude", "fallback-model", "high"),
			),
		}),
		RuntimeFactory: runtimeFactoryForLifecycleTest(nil),
	}, agentSessionScope{
		RunID:    "run-fallback-boundary",
		Kind:     "task",
		ID:       "task_02",
		Category: roundconfig.CategoryBackend,
		Session:  agent.SessionRef{Name: "roundfix-run-fallback-boundary-task_02", WorkDir: t.TempDir()},
	})
	if err != nil {
		t.Fatalf("agentSessionOwner() error = %v", err)
	}
	return owner
}

type fallbackBoundaryRunner struct {
	mu                sync.Mutex
	prepareErrByModel map[string]error
	runErrByModel     map[string]error
	emitOutputByModel map[string]bool
	prepared          []string
	ran               []string
}

func (*fallbackBoundaryRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *fallbackBoundaryRunner) PrepareSession(_ context.Context, req agent.ExecuteRequest, _ runevent.Sink) error {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	runner.prepared = append(runner.prepared, req.Runtime.Model)
	return runner.prepareErrByModel[req.Runtime.Model]
}

func (runner *fallbackBoundaryRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	return runner.RunPrepared(ctx, req, sink)
}

func (runner *fallbackBoundaryRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.mu.Lock()
	runner.ran = append(runner.ran, req.Runtime.Model)
	emitOutput := runner.emitOutputByModel[req.Runtime.Model]
	err := runner.runErrByModel[req.Runtime.Model]
	runner.mu.Unlock()
	result := agent.ExecuteResult{}
	if emitOutput {
		result.Output = "fallback output"
		if publishErr := sink.Publish(ctx, runevent.RunEvent{
			RunID:   req.RunID,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentMessage,
			Summary: "fallback output",
		}); publishErr != nil {
			return result, publishErr
		}
	}
	return result, err
}

func (*fallbackBoundaryRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func (runner *fallbackBoundaryRunner) runModels() []string {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return append([]string(nil), runner.ran...)
}

func (runner *fallbackBoundaryRunner) preparedModels() []string {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return append([]string(nil), runner.prepared...)
}
