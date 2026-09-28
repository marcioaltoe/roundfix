// Suite: Profile proof timeout recovery
// Invariant: only a proof-owned setup deadline is retried, at most once, across every profile-readiness consumer.
// Boundary IN: profiles validate and operational Preflight with an injected Agent runner.
// Boundary OUT: the ACP Runtime adapter, whose disposable-session outcome the runner simulates.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func TestProfileProofRetriesATimedOutSetupOnce(t *testing.T) {
	t.Parallel()
	homeDir, _ := withCLIWorkspace(t)
	config := roundconfig.Builtin()
	preferred := config.Profiles[roundconfig.CategoryBackend].Profile.Preferred
	fallback := config.Profiles[roundconfig.CategoryBackend].Profile.Fallbacks[0]
	runner := newProfileProofRetryRunner(func(_ context.Context, selection roundconfig.AgentSelection, attempt int) error {
		if selection == preferred && attempt == 1 {
			return fmt.Errorf("open disposable session: %w", context.DeadlineExceeded)
		}
		return nil
	})
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("profiles validate exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if !response.OK {
		t.Fatalf("profiles validate response = %+v, want success", response)
	}
	if got := runner.attempts[preferred]; got != 2 {
		t.Fatalf("preferred proof attempts = %d, want 2", got)
	}
	if got := runner.attempts[fallback]; got != 1 {
		t.Fatalf("fallback proof attempts = %d, want 1", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("profiles validate stderr = %q, want empty", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofReportsASecondTimeoutAsTemporary(t *testing.T) {
	t.Parallel()
	homeDir, _ := withCLIWorkspace(t)
	config := roundconfig.Builtin()
	preferred := config.Profiles[roundconfig.CategoryBackend].Profile.Preferred
	runner := newProfileProofRetryRunner(func(_ context.Context, selection roundconfig.AgentSelection, _ int) error {
		if selection == preferred {
			return fmt.Errorf("open disposable session: %w", context.DeadlineExceeded)
		}
		return nil
	})
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("profiles validate exit = %d, want %d", code, exitPreflight)
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if response.OK || len(response.Proofs) == 0 {
		t.Fatalf("profiles validate response = %+v, want failed proof", response)
	}
	proof := response.Proofs[0]
	if proof.Classification != "temporary" {
		t.Fatalf("proof classification = %q, want temporary", proof.Classification)
	}
	for _, want := range []string{"rerun", "when load drops", "configured profile was not shown to be wrong"} {
		if !strings.Contains(proof.NextAction, want) {
			t.Fatalf("proof next action missing %q: %q", want, proof.NextAction)
		}
	}
	if strings.Contains(proof.NextAction, "roundfix profiles configure") {
		t.Fatalf("proof next action advises reconfiguration: %q", proof.NextAction)
	}
	if got := runner.attempts[preferred]; got != 2 {
		t.Fatalf("preferred proof attempts = %d, want 2", got)
	}
	if !strings.Contains(stderr.String(), "classification: temporary") || strings.Contains(stderr.String(), "roundfix profiles configure") {
		t.Fatalf("profiles validate stderr = %q, want temporary rerun advice without reconfiguration", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofDoesNotRetryARejectedSelection(t *testing.T) {
	t.Parallel()
	homeDir, _ := withCLIWorkspace(t)
	config := roundconfig.Builtin()
	preferred := config.Profiles[roundconfig.CategoryBackend].Profile.Preferred
	runner := newProfileProofRetryRunner(func(_ context.Context, selection roundconfig.AgentSelection, _ int) error {
		return &agent.SelectionRejectedError{
			Assignment: agent.SelectionAssignment{
				Runtime:         selection.Runtime,
				Model:           selection.Model,
				ReasoningEffort: selection.ReasoningEffort,
			},
			Operation: "set Agent Selection",
			Err:       errors.New("adapter rejected selection"),
		}
	})
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("profiles validate exit = %d, want %d", code, exitPreflight)
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if len(response.Proofs) == 0 || response.Proofs[0].Classification != agent.SelectionRejected {
		t.Fatalf("profiles validate response = %+v, want selection_rejected", response)
	}
	if got := runner.attempts[preferred]; got != 1 {
		t.Fatalf("rejected proof attempts = %d, want 1", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofDoesNotRetryACancelledCommand(t *testing.T) {
	t.Parallel()
	homeDir, _ := withCLIWorkspace(t)
	config := roundconfig.Builtin()
	preferred := config.Profiles[roundconfig.CategoryBackend].Profile.Preferred
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runner := newProfileProofRetryRunner(func(ctx context.Context, _ roundconfig.AgentSelection, _ int) error {
		cancel()
		return ctx.Err()
	})
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, ctx, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("profiles validate exit = %d, want %d", code, exitPreflight)
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if len(response.Proofs) == 0 || response.Proofs[0].Classification != "" {
		t.Fatalf("profiles validate response = %+v, want unclassified cancellation", response)
	}
	if !strings.Contains(response.Proofs[0].Error, context.Canceled.Error()) {
		t.Fatalf("proof error = %q, want context cancellation", response.Proofs[0].Error)
	}
	if got := runner.attempts[preferred]; got != 1 {
		t.Fatalf("cancelled proof attempts = %d, want 1", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce(t *testing.T) {
	t.Parallel()
	config := roundconfig.Builtin()
	profile := config.Profiles[roundconfig.CategoryBackend].Profile
	preferred := profile.Preferred
	fallback := profile.Fallbacks[0]
	runner := newProfileProofRetryRunner(func(_ context.Context, selection roundconfig.AgentSelection, attempt int) error {
		if selection == fallback && attempt == 1 {
			return context.DeadlineExceeded
		}
		return nil
	})

	result, err := runProfileOperationalPreflight(
		context.Background(),
		commandRequest{name: "implement"},
		config,
		[]roundconfig.WorkCategory{roundconfig.CategoryBackend},
		"/workspace",
		runner,
		io.Discard,
	)

	if err != nil {
		t.Fatalf("operational Preflight error = %v", err)
	}
	if result.Err != nil || len(result.Proofs) != 2 {
		t.Fatalf("operational Preflight result = %+v, want two passed proofs", result)
	}
	for _, proof := range result.Proofs {
		if proof.Status != "passed" {
			t.Fatalf("proof = %+v, want passed", proof)
		}
	}
	if got := runner.attempts[preferred]; got != 1 {
		t.Fatalf("preferred proof attempts = %d, want 1", got)
	}
	if got := runner.attempts[fallback]; got != 2 {
		t.Fatalf("fallback proof attempts = %d, want 2", got)
	}
}

type profileProofRetryRunner struct {
	fakeAgentRunner
	attempts map[roundconfig.AgentSelection]int
	prove    func(context.Context, roundconfig.AgentSelection, int) error
}

func newProfileProofRetryRunner(prove func(context.Context, roundconfig.AgentSelection, int) error) *profileProofRetryRunner {
	return &profileProofRetryRunner{
		attempts: make(map[roundconfig.AgentSelection]int),
		prove:    prove,
	}
}

func (runner *profileProofRetryRunner) ProveProfileSelection(ctx context.Context, req agent.ProbeRequest) (agent.SelectionProof, error) {
	selection := roundconfig.AgentSelection{
		Runtime:         req.Runtime.ID,
		Model:           req.Runtime.Model,
		ReasoningEffort: req.Runtime.ReasoningEffort,
	}
	runner.attempts[selection]++
	if runner.prove != nil {
		if err := runner.prove(ctx, selection, runner.attempts[selection]); err != nil {
			return agent.SelectionProof{}, err
		}
	}
	return agent.SelectionProof{
		Runtime:               selection.Runtime,
		Model:                 selection.Model,
		ReasoningEffort:       selection.ReasoningEffort,
		EffectiveAccessPolicy: req.Runtime.RequestedPolicy(),
		Status:                agent.SelectionProofStatusProven,
	}, nil
}
