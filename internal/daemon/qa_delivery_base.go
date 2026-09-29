package daemon

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"roundfix/internal/preflight"
)

// qaDeliveryBase resolves the point where the audited head forked from the
// repository default branch. A repository without enough default-branch
// history leaves the base unresolved instead of refusing the QA gate.
func qaDeliveryBase(ctx context.Context, plan TaskPlan) (string, bool, error) {
	runner := preflight.ExecGitRunner{}
	currentBranch, err := runner.RunGit(ctx, plan.WorkDir, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		if !qaGitExited(err, 1) {
			return "", false, fmt.Errorf("read current branch for QA Delivery Base: %w", err)
		}
		currentBranch = ""
	}

	detected := preflight.DetectDefaultBranch(ctx, plan.WorkDir, strings.TrimSpace(currentBranch), runner)
	if detected.Source == preflight.DefaultBranchUndetermined {
		return "", false, nil
	}

	remoteRef := "refs/remotes/origin/" + detected.Name
	resolvedRef, found, err := qaResolveDeliveryRef(ctx, runner, plan.WorkDir, remoteRef)
	if err != nil {
		return "", false, fmt.Errorf("resolve remote-tracking default branch %q: %w", remoteRef, err)
	}
	if !found {
		localRef := "refs/heads/" + detected.Name
		resolvedRef, found, err = qaResolveDeliveryRef(ctx, runner, plan.WorkDir, localRef)
		if err != nil {
			return "", false, fmt.Errorf("resolve local default branch %q: %w", localRef, err)
		}
		if !found {
			return "", false, nil
		}
	}

	base, err := runner.RunGit(ctx, plan.WorkDir, "merge-base", resolvedRef, "HEAD")
	if err != nil {
		if qaGitExited(err, 1) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find QA Delivery Base with %q: %w", resolvedRef, err)
	}
	return strings.TrimSpace(base), true, nil
}

func qaResolveDeliveryRef(ctx context.Context, runner preflight.GitRunner, workDir, ref string) (string, bool, error) {
	if _, err := runner.RunGit(ctx, workDir, "rev-parse", "--verify", "--quiet", ref+"^{commit}"); err != nil {
		if qaGitExited(err, 1) {
			return "", false, nil
		}
		return "", false, err
	}
	return ref, true, nil
}

func qaGitExited(err error, code int) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == code
}
