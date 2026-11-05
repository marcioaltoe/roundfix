package cli

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"roundfix/internal/preflight"
)

// reviewEffectiveBase returns the candidate commit's merge base with the review base.
func reviewEffectiveBase(ctx context.Context, git preflight.GitRunner, root, base, commit string) (string, error) {
	effectiveBase, err := git.RunGit(ctx, root, "merge-base", base, commit)
	if err != nil {
		return "", fmt.Errorf("find review effective base of %s and %s: %w", base, commit, err)
	}
	return strings.TrimSpace(effectiveBase), nil
}

// reviewMeasuredFrom carries from's own change onto to's Effective Base.
// A conflicting carry keeps the original measurement and returns false.
func reviewMeasuredFrom(ctx context.Context, git preflight.GitRunner, root, base, from, to string) (string, bool, error) {
	fromBase, err := reviewEffectiveBase(ctx, git, root, base, from)
	if err != nil {
		return "", false, fmt.Errorf("measure review diff from %s to %s: %w", from, to, err)
	}
	toBase, err := reviewEffectiveBase(ctx, git, root, base, to)
	if err != nil {
		return "", false, fmt.Errorf("measure review diff from %s to %s: %w", from, to, err)
	}
	if fromBase == toBase {
		return from, true, nil
	}
	if fromBase == from {
		return toBase, true, nil
	}
	tree, err := git.RunGit(ctx, root, "merge-tree", "--write-tree", from, toBase)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return from, false, nil
		}
		return "", false, fmt.Errorf("carry review diff from %s to %s onto %s: %w", from, to, toBase, err)
	}
	firstLine, _, _ := strings.Cut(tree, "\n")
	return strings.TrimSpace(firstLine), true, nil
}
