package baselineacp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/baseline"
)

const (
	PreferredModel          = "gpt-5.6-sol"
	FallbackModel           = "gpt-5.5"
	RequiredReasoningEffort = "xhigh"
)

// Runtime is the non-Run ACP boundary used by semantic Baseline analysis.
// Neither method accepts a Run, Run Event sink, checkout, or Agent-log path.
type Runtime interface {
	ProveProfileSelection(context.Context, agent.ProbeRequest) (agent.SelectionProof, error)
	RunSealedPrompt(context.Context, agent.SealedPromptRequest) (agent.SealedPromptResult, error)
}

// Analyzer supervises the fixed preferred and fallback Agent Selections.
type Analyzer struct {
	runtime Runtime
}

func NewAnalyzer(runtime Runtime) *Analyzer {
	return &Analyzer{runtime: runtime}
}

func NewDefaultAnalyzer() *Analyzer {
	return NewAnalyzer(agent.NewDefaultRunner())
}

// Classify eagerly proves both fixed selections, then admits the first
// complete proposal. Unavailable or invalid analysis returns the complete
// deterministic manual proposal.
func (analyzer *Analyzer) Classify(
	ctx context.Context,
	snapshot baseline.AnalysisSnapshot,
) (baseline.ClassificationProposal, error) {
	if ctx == nil {
		return baseline.ClassificationProposal{}, errors.New(
			"Baseline proposal analysis context is required",
		)
	}
	if err := ctx.Err(); err != nil {
		return baseline.ClassificationProposal{}, err
	}
	if analyzer == nil || analyzer.runtime == nil {
		return baseline.ClassificationProposal{}, errors.New(
			"Baseline proposal analysis runtime is required",
		)
	}
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		return baseline.ClassificationProposal{}, err
	}

	selections := []agent.RuntimeSpec{
		classificationRuntime(PreferredModel),
		classificationRuntime(FallbackModel),
	}
	proven := make([]bool, len(selections))
	for index, selection := range selections {
		proof, proofErr := analyzer.proveSelection(ctx, selection)
		if err := ctx.Err(); err != nil {
			return baseline.ClassificationProposal{}, err
		}
		if proofErr != nil {
			if cleanupUnproven(proofErr) {
				return baseline.ManualClassificationProposal(snapshot)
			}
			continue
		}
		proven[index] = exactProofMatches(selection, proof)
	}

	for index, selection := range selections {
		if !proven[index] {
			continue
		}
		proposal, attemptErr := analyzer.attempt(ctx, selection, canonical, snapshot)
		if attemptErr == nil {
			return proposal, nil
		}
		if err := ctx.Err(); err != nil {
			return baseline.ClassificationProposal{}, err
		}
		if cleanupUnproven(attemptErr) {
			break
		}
	}
	return baseline.ManualClassificationProposal(snapshot)
}

// Revise supervises the same fixed selections with the separate strict
// Revision Snapshot and Proposal contracts. Invalid or unavailable semantic
// translation preserves the maintainer's manual revision path.
func (analyzer *Analyzer) Revise(
	ctx context.Context,
	snapshot baseline.RevisionSnapshot,
) (baseline.RevisionProposal, error) {
	if ctx == nil {
		return baseline.RevisionProposal{}, errors.New("Baseline revision analysis context is required")
	}
	if err := ctx.Err(); err != nil {
		return baseline.RevisionProposal{}, err
	}
	if analyzer == nil || analyzer.runtime == nil {
		return baseline.RevisionProposal{}, errors.New("Baseline revision analysis runtime is required")
	}
	canonical, err := snapshot.CanonicalBytes()
	if err != nil {
		return baseline.RevisionProposal{}, err
	}
	selections := []agent.RuntimeSpec{
		classificationRuntime(PreferredModel),
		classificationRuntime(FallbackModel),
	}
	proven := make([]bool, len(selections))
	for index, selection := range selections {
		proof, proofErr := analyzer.proveSelection(ctx, selection)
		if err := ctx.Err(); err != nil {
			return baseline.RevisionProposal{}, err
		}
		if proofErr != nil {
			if cleanupUnproven(proofErr) {
				return baseline.ManualRevisionProposal(snapshot)
			}
			continue
		}
		proven[index] = exactProofMatches(selection, proof)
	}
	for index, selection := range selections {
		if !proven[index] {
			continue
		}
		proposal, attemptErr := analyzer.attemptRevision(ctx, selection, canonical, snapshot)
		if attemptErr == nil {
			return proposal, nil
		}
		if err := ctx.Err(); err != nil {
			return baseline.RevisionProposal{}, err
		}
		if cleanupUnproven(attemptErr) {
			break
		}
	}
	return baseline.ManualRevisionProposal(snapshot)
}

func cleanupUnproven(err error) bool {
	var sessionErr *agent.AgentSessionCleanupError
	var directoryErr *privateDirectoryCleanupError
	return errors.As(err, &sessionErr) || errors.As(err, &directoryErr)
}

func (analyzer *Analyzer) proveSelection(
	ctx context.Context,
	selection agent.RuntimeSpec,
) (agent.SelectionProof, error) {
	var proof agent.SelectionProof
	err := withPrivateDirectory("roundfix-baseline-proof-", func(workDir string) error {
		var proofErr error
		proof, proofErr = analyzer.runtime.ProveProfileSelection(ctx, agent.ProbeRequest{
			Runtime: selection,
			WorkDir: workDir,
		})
		return proofErr
	})
	return proof, err
}

func (analyzer *Analyzer) attempt(
	ctx context.Context,
	selection agent.RuntimeSpec,
	canonical []byte,
	snapshot baseline.AnalysisSnapshot,
) (baseline.ClassificationProposal, error) {
	var result agent.SealedPromptResult
	err := withPrivateDirectory("roundfix-baseline-analysis-", func(workDir string) error {
		var attemptErr error
		result, attemptErr = analyzer.runtime.RunSealedPrompt(ctx, agent.SealedPromptRequest{
			Runtime: selection,
			WorkDir: workDir,
			Input:   append([]byte(nil), canonical...),
		})
		return attemptErr
	})
	if err != nil {
		return baseline.ClassificationProposal{}, err
	}
	if result.ToolUsed {
		return baseline.ClassificationProposal{}, agent.ErrSealedToolUse
	}
	return baseline.ParseClassificationProposal(result.Output, snapshot)
}

func (analyzer *Analyzer) attemptRevision(
	ctx context.Context,
	selection agent.RuntimeSpec,
	canonical []byte,
	snapshot baseline.RevisionSnapshot,
) (baseline.RevisionProposal, error) {
	var result agent.SealedPromptResult
	err := withPrivateDirectory("roundfix-baseline-revision-", func(workDir string) error {
		var attemptErr error
		result, attemptErr = analyzer.runtime.RunSealedPrompt(ctx, agent.SealedPromptRequest{
			Runtime: selection,
			WorkDir: workDir,
			Input:   append([]byte(nil), canonical...),
		})
		return attemptErr
	})
	if err != nil {
		return baseline.RevisionProposal{}, err
	}
	if result.ToolUsed {
		return baseline.RevisionProposal{}, agent.ErrSealedToolUse
	}
	return baseline.ParseRevisionProposal(result.Output, snapshot)
}

func classificationRuntime(model string) agent.RuntimeSpec {
	return agent.RuntimeSpec{
		ID:              "codex",
		DisplayName:     "Codex",
		Protocol:        agent.ProtocolACP,
		Model:           model,
		ReasoningEffort: RequiredReasoningEffort,
	}
}

func exactProofMatches(selection agent.RuntimeSpec, proof agent.SelectionProof) bool {
	return proof.Status == agent.SelectionProofStatusProven &&
		strings.TrimSpace(proof.Runtime) == selection.ID &&
		strings.TrimSpace(proof.Model) == selection.Model &&
		strings.TrimSpace(proof.ReasoningEffort) == selection.ReasoningEffort
}

func withPrivateDirectory(prefix string, run func(string) error) error {
	workDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return fmt.Errorf("create private sealed Agent directory: %w", err)
	}
	if err := os.Chmod(workDir, 0o700); err != nil {
		_ = os.RemoveAll(workDir)
		return fmt.Errorf("make sealed Agent directory private: %w", err)
	}
	runErr := run(workDir)
	removeErr := os.RemoveAll(workDir)
	if removeErr != nil {
		cleanupErr := &privateDirectoryCleanupError{Path: workDir, Err: removeErr}
		if runErr != nil {
			return errors.Join(runErr, cleanupErr)
		}
		return cleanupErr
	}
	return runErr
}

type privateDirectoryCleanupError struct {
	Path string
	Err  error
}

func (err *privateDirectoryCleanupError) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("remove private sealed Agent directory %q: %v", err.Path, err.Err)
}

func (err *privateDirectoryCleanupError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}
