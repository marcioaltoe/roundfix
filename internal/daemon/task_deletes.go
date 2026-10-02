package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"roundfix/internal/spec"
)

// undeletedTaskPaths lists deletion declarations still present in the worktree.
// Lstat counts directories and dangling symlinks as remaining paths too.
func undeletedTaskPaths(workDir string, task spec.Task) ([]string, error) {
	var remaining []string
	for _, ref := range task.Context {
		if ref.Kind != spec.ContextKindDeletes {
			continue
		}
		_, err := os.Lstat(filepath.Join(workDir, filepath.FromSlash(ref.Path)))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("inspect deletes path %q: %w", ref.Path, err)
		}
		remaining = append(remaining, ref.Path)
	}
	return remaining, nil
}
