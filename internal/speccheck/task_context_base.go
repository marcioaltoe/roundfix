package speccheck

import (
	"context"
	"path/filepath"
	"strings"

	"roundfix/internal/preflight"
)

// taskContextTrackedAtBase reports whether path is tracked at HEAD or at the
// Delivery Base of the Git work tree whose top level is repoRoot. An unreadable
// revision cannot resolve the path; a nested root reads neither revision.
func taskContextTrackedAtBase(repoRoot, path string) bool {
	ctx := context.Background()
	runner := preflight.ExecGitRunner{}
	top, err := runner.RunGit(ctx, repoRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		return false
	}
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return false
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return false
	}
	top, err = filepath.EvalSymlinks(strings.TrimSpace(top))
	if err != nil || root != top {
		return false
	}
	path = filepath.ToSlash(path)
	if _, err := runner.RunGit(ctx, repoRoot, "cat-file", "-e", "HEAD:"+path); err == nil {
		return true
	}
	branch, _ := runner.RunGit(ctx, repoRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
	head, found, err := preflight.ResolveDefaultBranchHead(ctx, repoRoot, strings.TrimSpace(branch), runner)
	if err != nil || !found {
		return false
	}
	base, err := runner.RunGit(ctx, repoRoot, "merge-base", head.Ref, "HEAD")
	if err != nil {
		return false
	}
	_, err = runner.RunGit(ctx, repoRoot, "cat-file", "-e", strings.TrimSpace(base)+":"+path)
	return err == nil
}

// taskContextUnresolvedFix is the fix of the Task Context SC-REF-UNRESOLVED.
func taskContextUnresolvedFix(path, taskDisplayPath string) string {
	return "Remove the Task Context entry for " + path + " from " + taskDisplayPath + ", or declare it under deletes: when the Task removes it; replace it only with a path the Task itself reads or changes."
}
