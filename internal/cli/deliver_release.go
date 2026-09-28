package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func releaseMergedSpecRuns(
	ctx context.Context,
	runStore *store.Store,
	repository string,
	merged runworktree.MergedHead,
) error {
	if runStore == nil {
		return errors.New("release merged Spec Runs: Run Database is unavailable")
	}
	repository = strings.TrimSpace(repository)
	if repository == "" {
		return errors.New("release merged Spec Runs: repository is required")
	}
	specSlug := strings.TrimSpace(merged.SpecSlug)
	if specSlug == "" {
		return errors.New("release merged Spec Runs: Spec slug is required")
	}
	repositoryRoot, err := roundconfig.RepositoryRoot(repository)
	if err != nil {
		return fmt.Errorf("release merged Spec %q Runs: resolve repository identity: %w", specSlug, err)
	}

	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{
		RepositoryRoot: repositoryRoot,
		States:         store.StatesAll,
	})
	if err != nil {
		return fmt.Errorf("release merged Spec %q Runs: list Runs: %w", specSlug, err)
	}

	var kept []error
	for _, run := range runs {
		if run.Kind != store.KindImplement || strings.TrimSpace(run.SpecSlug) != specSlug {
			continue
		}
		if !store.IsTerminalState(run.State) {
			kept = append(kept, fmt.Errorf("keep Run %q: Run is %s", run.ID, run.State))
			continue
		}

		run.GitRoot = repository
		inspected, inspectErr := runworktree.InspectTerminalRunMerged(ctx, run, []runworktree.MergedHead{merged})
		if inspectErr != nil {
			kept = append(kept, fmt.Errorf("keep Run %q: %w", run.ID, inspectErr))
			continue
		}
		switch inspected.State {
		case runworktree.ReconciliationReleased:
			continue
		case runworktree.ReconciliationSafe, runworktree.ReconciliationSuperseded:
			if err := runworktree.ApplyTerminalRun(ctx, runStore, inspected); err != nil {
				kept = append(kept, fmt.Errorf("keep Run %q: %w", run.ID, err))
			}
		default:
			reason := strings.TrimSpace(inspected.Reason)
			if reason == "" {
				reason = fmt.Sprintf("classification is %q", inspected.State)
			}
			kept = append(kept, fmt.Errorf("keep Run %q: %s", run.ID, reason))
		}
	}
	return errors.Join(kept...)
}
