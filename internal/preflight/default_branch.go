package preflight

import (
	"context"
	"strings"
)

const (
	originHeadRef    = "refs/remotes/origin/HEAD"
	originHeadPrefix = "refs/remotes/origin/"
)

// DefaultBranchSource identifies how the repository default branch was
// determined, so a Preflight Validation veto message can explain the
// detection and a false positive stays diagnosable from the message alone.
type DefaultBranchSource string

const (
	// DefaultBranchFromOriginHead means origin/HEAD named the default branch.
	DefaultBranchFromOriginHead DefaultBranchSource = "origin/HEAD"
	// DefaultBranchFromNameMatch means origin/HEAD was unset and the current
	// branch matched a conventional default branch name (main or master).
	DefaultBranchFromNameMatch DefaultBranchSource = "name-match"
	// DefaultBranchFromOriginBranch means exactly one fetched main/master ref names the default.
	DefaultBranchFromOriginBranch DefaultBranchSource = "origin-branch"
	// DefaultBranchUndetermined means no default branch could be determined.
	DefaultBranchUndetermined DefaultBranchSource = "undetermined"
)

// DefaultBranch is the result of repository default-branch detection.
type DefaultBranch struct {
	// Name is the detected default branch name. Empty when Source is
	// DefaultBranchUndetermined.
	Name   string
	Source DefaultBranchSource
}

// IsDefault reports whether branch is the repository default branch under
// this detection. An undetermined detection never matches: without
// origin/HEAD, a name match or one fetched main/master ref names the default.
func (detection DefaultBranch) IsDefault(branch string) bool {
	return detection.Source != DefaultBranchUndetermined && branch == detection.Name
}

// DetectDefaultBranch resolves the repository default branch through the git
// runner. origin/HEAD is authoritative when set, regardless of the branch
// name it points at. When origin/HEAD is unset or unreadable (fresh clones,
// no remote), a currentBranch named main or master is treated as the default
// by name match. Otherwise exactly one fetched origin/main or origin/master
// names the default; both or neither leave it undetermined.
func DetectDefaultBranch(ctx context.Context, workDir string, currentBranch string, runner GitRunner) DefaultBranch {
	if runner == nil {
		runner = ExecGitRunner{}
	}
	target, err := runner.RunGit(ctx, workDir, "symbolic-ref", originHeadRef)
	if err == nil {
		branch, ok := strings.CutPrefix(strings.TrimSpace(target), originHeadPrefix)
		if ok && branch != "" {
			return DefaultBranch{Name: branch, Source: DefaultBranchFromOriginHead}
		}
	}
	if currentBranch == "main" || currentBranch == "master" {
		return DefaultBranch{Name: currentBranch, Source: DefaultBranchFromNameMatch}
	}
	refs, err := runner.RunGit(ctx, workDir, "for-each-ref", "--format=%(refname)", "refs/remotes/origin/main", "refs/remotes/origin/master")
	if err == nil {
		name := ""
		for _, ref := range strings.Split(strings.TrimSpace(refs), "\n") {
			if ref != "refs/remotes/origin/main" && ref != "refs/remotes/origin/master" {
				continue
			}
			if name != "" {
				return DefaultBranch{Source: DefaultBranchUndetermined}
			}
			name = strings.TrimPrefix(ref, originHeadPrefix)
		}
		if name != "" {
			return DefaultBranch{Name: name, Source: DefaultBranchFromOriginBranch}
		}
	}
	return DefaultBranch{Source: DefaultBranchUndetermined}
}
