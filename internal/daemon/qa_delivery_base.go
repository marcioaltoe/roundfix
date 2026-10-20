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

	head, found, err := preflight.ResolveDefaultBranchHead(ctx, plan.WorkDir, strings.TrimSpace(currentBranch), runner)
	if err != nil {
		return "", false, fmt.Errorf("resolve default branch for QA Delivery Base: %w", err)
	}
	if !found {
		return "", false, nil
	}

	base, err := runner.RunGit(ctx, plan.WorkDir, "merge-base", head.Ref, "HEAD")
	if err != nil {
		if qaGitExited(err, 1) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find QA Delivery Base with %q: %w", head.Ref, err)
	}
	return strings.TrimSpace(base), true, nil
}

func qaGitExited(err error, code int) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && exitErr.ExitCode() == code
}
