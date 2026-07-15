package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// FallbackCandidateSet carries one ACP Runtime's ordered fallback vocabulary.
type FallbackCandidateSet struct {
	Models  []string
	Efforts []string
}

// FallbackSelection is a proven Agent Model and its highest accepted
// reasoning effort. An empty ReasoningEffort means the model manages reasoning.
type FallbackSelection struct {
	Model           string
	ReasoningEffort string
}

// ProbeFallback proves the first non-failed candidate with a disposable Agent
// Session. Candidate and effort order are caller-owned and preserved.
func (runner ACPXRunner) ProbeFallback(ctx context.Context, runtime RuntimeSpec, candidates FallbackCandidateSet) (FallbackSelection, bool, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return FallbackSelection{}, false, fmt.Errorf("resolve fallback probe working directory: %w", err)
	}
	failedModel := strings.TrimSpace(runtime.Model)
	for _, value := range candidates.Models {
		if err := ctx.Err(); err != nil {
			return FallbackSelection{}, false, err
		}
		model := strings.TrimSpace(value)
		if model == "" || model == failedModel {
			continue
		}
		selection, functional, err := runner.probeFallbackCandidate(ctx, runtime, model, candidates.Efforts, workDir)
		if err != nil {
			return FallbackSelection{}, false, fmt.Errorf("probe fallback Agent Model %q: %w", model, err)
		}
		if functional {
			return selection, true, nil
		}
	}
	return FallbackSelection{}, false, nil
}

func (runner ACPXRunner) probeFallbackCandidate(ctx context.Context, runtime RuntimeSpec, model string, efforts []string, workDir string) (FallbackSelection, bool, error) {
	candidate := runtime
	candidate.Model = model
	candidate.ReasoningEffort = ""
	if err := validateRuntimeSelection(candidate); err != nil {
		return FallbackSelection{}, false, err
	}
	sessionName, err := disposablePreflightSessionName()
	if err != nil {
		return FallbackSelection{}, false, err
	}

	setupCtx, setupCancel := context.WithTimeout(ctx, acpxPreflightSetupTimeout)
	codexEnv, err := runner.codexEnvForSession(setupCtx, candidate, sessionName)
	if err != nil {
		setupCancel()
		return FallbackSelection{}, false, err
	}
	selection := FallbackSelection{Model: model}
	functional := true
	probeErr := runner.applyDisposableSelection(setupCtx, candidate, sessionName, workDir, codexEnv)
	if probeErr != nil {
		functional = false
	} else {
		for _, value := range efforts {
			effort := strings.TrimSpace(value)
			if effort == "" {
				continue
			}
			candidate.ReasoningEffort = effort
			effortErr := runner.applyDisposableEffort(setupCtx, candidate, sessionName, workDir, codexEnv)
			if effortErr == nil {
				selection.ReasoningEffort = effort
				break
			}
			if !isSelectionAssignmentRejection(effortErr) {
				probeErr = effortErr
				functional = false
				break
			}
		}
	}
	setupCancel()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.WithoutCancel(ctx), acpxPreflightCleanupTimeout)
	defer cleanupCancel()
	cleanupErr := runner.closeDisposableSession(cleanupCtx, candidate, sessionName, workDir)
	if cleanupErr != nil {
		if probeErr != nil {
			return FallbackSelection{}, false, errors.Join(probeErr, cleanupErr)
		}
		return FallbackSelection{}, false, cleanupErr
	}
	if probeErr != nil && !isSelectionAssignmentRejection(probeErr) {
		return FallbackSelection{}, false, probeErr
	}
	if !functional {
		return FallbackSelection{}, false, nil
	}
	return selection, true, nil
}

func (runner ACPXRunner) applyDisposableEffort(ctx context.Context, runtime RuntimeSpec, sessionName string, workDir string, codexEnv []string) error {
	key, err := acpxReasoningEffortConfigKey(runtime)
	if err != nil {
		return err
	}
	args, err := acpxSetConfigArgs(runtime, key, runtime.ReasoningEffort, sessionName, workDir)
	if err != nil {
		return err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return selectionPreflightError(runtime, "set "+key, fmt.Errorf("set disposable acpx Agent Session %s %q: %w", key, runtime.ReasoningEffort, classifyModelNotAdvertised(runtime, err)))
	}
	return nil
}

func isSelectionAssignmentRejection(err error) bool {
	var selectionErr *SelectionPreflightError
	if !errors.As(err, &selectionErr) {
		return false
	}
	var commandErr *InfrastructureError
	return errors.As(selectionErr.Err, &commandErr)
}
