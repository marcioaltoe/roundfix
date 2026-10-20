// Suite: Profile proof timeout recovery
// Invariant: only a proof-owned setup deadline is retried, at most twice after backoff, across every profile-readiness consumer.
// Boundary IN: profiles validate and operational Preflight with an injected Agent runner.
// Boundary OUT: the ACP Runtime adapter, whose disposable-session outcome the runner simulates.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

func TestProfileProofRetriesATimedOutSetupOnce(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
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
	assertProfileProofWaits(t, *waits, profileProofRetryWaits[0])
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofReportsAThirdTimeoutAsTemporary(t *testing.T) {
	t.Parallel()
	for _, jsonMode := range []bool{true, false} {
		t.Run(fmt.Sprintf("json=%t", jsonMode), func(t *testing.T) {
			waits := recordProfileProofSleeps(t)
			homeDir, _ := withCLIWorkspace(t)
			config := roundconfig.Builtin()
			preferred := config.Profiles[roundconfig.CategoryBackend].Profile.Preferred
			runner := newProfileProofRetryRunner(func(_ context.Context, _ roundconfig.AgentSelection, _ int) error {
				return fmt.Errorf("open disposable session: %w", context.DeadlineExceeded)
			})
			withAgentRunner(t, runner)
			var stdout, stderr bytes.Buffer
			args := []string{"profiles", "validate", "--category", "backend"}
			if jsonMode {
				args = append(args, "--json")
			}
			code := runCLI(t, args, &stdout, &stderr)
			if code != exitPreflight {
				t.Fatalf("exit = %d, want 2", code)
			}
			wantError := fmt.Sprintf("profile proof setup timed out on 3 attempts of %s each: open disposable session: context deadline exceeded", config.Defaults.ProfileProofSetupTimeout)
			const wantNext = "rerun the command when load drops, or raise defaults.profile_proof_setup_timeout; the configured profile was not shown to be wrong"
			if jsonMode {
				response := decodeProfilesValidateResponse(t, stdout.String())
				if response.OK || len(response.Proofs) == 0 {
					t.Fatalf("response = %+v, want failed proof", response)
				}
				proof := response.Proofs[0]
				if proof.Classification != "temporary" || proof.Error != wantError || proof.NextAction != wantNext {
					t.Fatalf("proof = %+v, want exact temporary timeout and recovery action", proof)
				}
			} else {
				if stdout.Len() != 0 {
					t.Fatalf("stdout = %q, want empty", stdout.String())
				}
				wantStderr := fmt.Sprintf("roundfix: profiles failed: profile proof failed for runtime %q, model %q, reasoning_effort %q; affected categories: backend preferred source=built-in; classification: temporary; adapter error: %s; next: %s; fallback: Fallback Chains activate only after Run creation (ADR-0050); Preflight proves every configured tuple and substitutes none\nRun 'roundfix profiles --help' for usage.\n", preferred.Runtime, preferred.Model, preferred.ReasoningEffort, wantError, wantNext)
				if stderr.String() != wantStderr {
					t.Fatalf("stderr = %q, want %q", stderr.String(), wantStderr)
				}
			}
			if runner.attempts[preferred] != 3 {
				t.Fatalf("attempts = %d, want 3", runner.attempts[preferred])
			}
			assertProfileProofWaits(t, *waits, profileProofRetryWaits...)
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestProfileProofDoesNotRetryARejectedSelection(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
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
	assertProfileProofWaits(t, *waits)
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofDoesNotRetryACancelledCommand(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
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
	assertProfileProofWaits(t, *waits)
	assertNoRunDatabase(t, homeDir)
	t.Run("default sleeper honors cancellation", func(t *testing.T) {
		if err := defaultCommandDependencies().profileProofSleep(ctx, time.Hour); !errors.Is(err, context.Canceled) {
			t.Fatalf("sleep error = %v, want cancellation", err)
		}
	})
	t.Run("default sleeper timer expires", func(t *testing.T) {
		if err := defaultCommandDependencies().profileProofSleep(t.Context(), 0); err != nil {
			t.Fatalf("sleep error = %v, want nil", err)
		}
	})
}

func TestOperationalPreflightPassesAfterAFallbackProofTimesOutOnce(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
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

	var stderr bytes.Buffer
	result, err := runProfileOperationalPreflight(
		commandContextForTest(t, t.Context()),
		commandRequest{name: "implement"},
		config,
		[]roundconfig.WorkCategory{roundconfig.CategoryBackend},
		"/workspace",
		runner,
		&stderr,
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
	assertProfileProofWaits(t, *waits, profileProofRetryWaits[0])
	if stderr.Len() != 0 {
		t.Fatalf("operational Preflight stderr = %q, want empty", stderr.String())
	}
}

func TestProfileProofPassesOnTheThirdAttempt(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
	homeDir, _ := withCLIWorkspace(t)
	preferred := roundconfig.Builtin().Profiles[roundconfig.CategoryBackend].Profile.Preferred
	runner := newProfileProofRetryRunner(func(_ context.Context, selection roundconfig.AgentSelection, attempt int) error {
		if selection == preferred && attempt < 3 {
			return context.DeadlineExceeded
		}
		return nil
	})
	withAgentRunner(t, runner)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)
	if code != exitOK || !decodeProfilesValidateResponse(t, stdout.String()).OK {
		t.Fatalf("exit = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}
	if runner.attempts[preferred] != 3 || stderr.Len() != 0 {
		t.Fatalf("attempts = %d, stderr = %q", runner.attempts[preferred], stderr.String())
	}
	assertProfileProofWaits(t, *waits, profileProofRetryWaits...)
	assertNoRunDatabase(t, homeDir)
}

func TestProfileProofStopsWhenCancelledDuringTheBackoff(t *testing.T) {
	t.Parallel()
	waits := recordProfileProofSleeps(t)
	homeDir, _ := withCLIWorkspace(t)
	preferred := roundconfig.Builtin().Profiles[roundconfig.CategoryBackend].Profile.Preferred
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.profileProofSleep = func(ctx context.Context, wait time.Duration) error {
			*waits = append(*waits, wait)
			cancel()
			return ctx.Err()
		}
	})
	runner := newProfileProofRetryRunner(func(_ context.Context, _ roundconfig.AgentSelection, _ int) error {
		return context.DeadlineExceeded
	})
	withAgentRunner(t, runner)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, ctx, []string{"profiles", "validate", "--category", "backend", "--json"}, &stdout, &stderr)
	if code != exitPreflight {
		t.Fatalf("exit = %d, want 2", code)
	}
	response := decodeProfilesValidateResponse(t, stdout.String())
	if len(response.Proofs) == 0 || response.Proofs[0].Classification != "" || !strings.Contains(response.Proofs[0].Error, context.Canceled.Error()) {
		t.Fatalf("response = %+v, want unclassified cancellation", response)
	}
	if runner.attempts[preferred] != 1 {
		t.Fatalf("attempts = %d, want 1", runner.attempts[preferred])
	}
	assertProfileProofWaits(t, *waits, profileProofRetryWaits[0])
	assertNoRunDatabase(t, homeDir)
}

func recordProfileProofSleeps(t *testing.T) *[]time.Duration {
	t.Helper()
	var waits []time.Duration
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.profileProofSleep = func(_ context.Context, wait time.Duration) error {
			waits = append(waits, wait)
			return nil
		}
	})
	return &waits
}

func assertProfileProofWaits(t *testing.T, got []time.Duration, want ...time.Duration) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("proof waits = %v, want %v", got, want)
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
