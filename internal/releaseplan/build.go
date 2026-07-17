package releaseplan

import (
	"context"
	"errors"
	"fmt"
)

// Build resolves the requested committed range, classifies every committed
// change, and calculates the final read-only Release Plan.
func Build(ctx context.Context, request Request, source GitSource) (Plan, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if source == nil {
		return Plan{}, GitSourceError{
			Operation:  "build release plan",
			NextAction: "configure a Release Plan Git source",
			Err:        errors.New("missing Git source"),
		}
	}

	releaseRange, err := source.ResolveRange(ctx, request.From, request.To)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve release range: %w", err)
	}
	commits, err := source.Commits(ctx, releaseRange)
	if err != nil {
		return Plan{}, fmt.Errorf("load release commits: %w", err)
	}
	classification, err := ClassifyChanges(ClassifyRequest{
		Commits:      commits,
		ManualImpact: request.ManualImpact,
		ManualReason: request.ManualReason,
	})
	if err != nil {
		return Plan{}, fmt.Errorf("classify release changes: %w", err)
	}

	proposal, err := CalculateProposal(releaseRange.Base.Version, ProposalInput{
		Impact:                       classification.Classification.Impact,
		Breaking:                     classification.Classification.Breaking,
		ManualClassificationRequired: classification.Classification.ManualClassificationRequired,
	})
	if err != nil {
		return Plan{}, fmt.Errorf("calculate release proposal: %w", err)
	}
	return NewPlan(releaseRange.Base, releaseRange.Target, classification.Classification, proposal, classification.Changes), nil
}
