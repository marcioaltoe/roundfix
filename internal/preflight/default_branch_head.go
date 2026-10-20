package preflight

import (
	"context"
	"fmt"
	"strings"
)

// DefaultBranchHead names the fetched default head, or the local head when no
// remote-tracking branch exists.
type DefaultBranchHead struct {
	Branch DefaultBranch
	Ref    string
	Commit string
}

// ResolveDefaultBranchHead reads local refs without fetching or changing them.
// An undetermined default or absent refs returns found=false.
func ResolveDefaultBranchHead(ctx context.Context, workDir, currentBranch string, runner GitRunner) (DefaultBranchHead, bool, error) {
	if runner == nil {
		runner = ExecGitRunner{}
	}
	branch := DetectDefaultBranch(ctx, workDir, currentBranch, runner)
	if branch.Source == DefaultBranchUndetermined {
		return DefaultBranchHead{}, false, nil
	}
	remote, local := originHeadPrefix+branch.Name, "refs/heads/"+branch.Name
	refs, err := runner.RunGit(ctx, workDir, "for-each-ref", "--format=%(refname)%00%(objectname)", remote, local)
	if err != nil {
		return DefaultBranchHead{}, false, fmt.Errorf("read default branch head: %w", err)
	}
	var fallback DefaultBranchHead
	for _, line := range strings.Split(strings.TrimSpace(refs), "\n") {
		ref, commit, ok := strings.Cut(line, "\x00")
		if !ok {
			continue
		}
		head := DefaultBranchHead{Branch: branch, Ref: ref, Commit: commit}
		if ref == remote {
			return head, true, nil
		}
		if ref == local {
			fallback = head
		}
	}
	return fallback, fallback.Ref != "", nil
}
