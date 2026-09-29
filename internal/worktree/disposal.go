package worktree

import "context"

// RemoveRegisteredWorktree removes a registered Git worktree without forcing
// dirty or locked work through the repository administration lock.
func RemoveRegisteredWorktree(ctx context.Context, gitRoot, path string) error {
	_, err := runWorktreeCommand(ctx, execGitRunner{}, gitRoot, "worktree", "remove", path)
	return err
}

// AddDetachedWorktree adds a detached Git worktree through the repository
// administration lock.
func AddDetachedWorktree(ctx context.Context, repository, path, revision string) error {
	_, err := runWorktreeCommand(
		ctx,
		execGitRunner{},
		repository,
		"worktree", "add", "--detach", path, revision,
	)
	return err
}

// RemoveWorktreeForce force-removes a registered Git worktree through the
// repository administration lock.
func RemoveWorktreeForce(ctx context.Context, repository, path string) error {
	_, err := runWorktreeCommand(
		ctx,
		execGitRunner{},
		repository,
		"worktree", "remove", "--force", path,
	)
	return err
}
