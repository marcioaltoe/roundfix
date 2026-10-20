package spec

import (
	"fmt"
	"path/filepath"
	"strings"

	"roundfix/internal/config"
)

// RepositoryPath returns path relative to repositoryRoot, slash-separated. Both
// arguments must be absolute. Configured Spec paths and Git's checkout
// root can spell the same location differently across a symbolic link. Compare
// lexically first to preserve existing paths, then resolve both spellings before
// handing paths to Git or writing them into archive files. Neither the main
// checkout nor the process work directory supplies a replacement root.
func RepositoryPath(repositoryRoot, path string) (string, error) {
	relative := func(root, target string) (string, bool) {
		rel, err := filepath.Rel(root, target)
		if err != nil || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", false
		}
		return filepath.ToSlash(rel), true
	}
	if rel, ok := relative(repositoryRoot, path); ok {
		return rel, nil
	}
	if rel, ok := relative(config.CanonicalizeResolvedPath(repositoryRoot), config.CanonicalizeResolvedPath(path)); ok {
		return rel, nil
	}
	return "", fmt.Errorf("path %q is outside the repository %q", path, repositoryRoot)
}
