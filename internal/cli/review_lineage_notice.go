package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/preflight"
)

// reviewLineageNotice returns an advisory restart notice without changing review artifacts.
func reviewLineageNotice(ctx context.Context, git preflight.GitRunner, loaded roundconfig.Loaded, specSlug, worktree string) string {
	provider := loaded.Config.PrePRReview.Provider
	if (provider != "codex" && provider != "claude") || strings.TrimSpace(worktree) == "" {
		return ""
	}
	if git == nil {
		git = preflight.ExecGitRunner{}
	}
	state, err := preflight.InspectGit(ctx, worktree, git)
	if err != nil {
		return ""
	}
	artifactDir, err := roundconfig.ResolveArtifactDirectory(loaded.Config.Defaults.ArtifactDir, state.Root, loaded.HomeDir)
	if err != nil {
		return ""
	}
	prior, err := readReviewRecord(filepath.Join(reviewCheckoutDir(artifactDir, state.Root), reviewRecordFileName))
	if err != nil || prior.HeadCommit == state.HEAD {
		return ""
	}
	baseTip, err := resolveReviewBaseCommit(ctx, "", state, git)
	if err != nil {
		return ""
	}
	base, err := resolveReviewMergeBase(ctx, state.Root, baseTip, state.HEAD, git)
	if err != nil {
		return ""
	}
	candidate := newReviewRecord(state.Root, base, state.HEAD, loaded.Config.PrePRReview, reviewOutcomeBlocked)
	candidate.BaseTipCommit = baseTip
	plan, err := decideReviewLineage(ctx, &prior, candidate, git)
	if err != nil || plan.Restart == "" {
		return ""
	}
	return fmt.Sprintf("Review: %s head %s opens a new Reviewer Lineage at round 1, a full review of the candidate: %s", specSlug, state.HEAD, plan.Restart)
}
