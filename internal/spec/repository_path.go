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
// resolved targets before preserving a lexical relative path. The configured
// Spec Root and Git top level select the sources and destinations for Archive
// Records, promoted copies and returned sources, which a repository commit must
// hold. Refuse an outside target before any write instead of following its link.
// Neither the main checkout nor the process work directory replaces the root.
func RepositoryPath(repositoryRoot, path string) (string, error) {
	relative := func(root, target string) (string, bool) {
		rel, err := filepath.Rel(root, target)
		if err != nil || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", false
		}
		return filepath.ToSlash(rel), true
	}
	if resolved, ok := relative(config.CanonicalizeResolvedPath(repositoryRoot), config.CanonicalizeResolvedPath(path)); ok {
		if lexical, ok := relative(repositoryRoot, path); ok {
			return lexical, nil
		}
		return resolved, nil
	}
	return "", fmt.Errorf("path %q is outside the repository %q", path, repositoryRoot)
}
