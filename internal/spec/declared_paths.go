package spec

import (
	"fmt"
	"path/filepath"
	"sort"
)

// DeclaredTaskPaths returns the creates: and interface: paths of one Task
// file and its RecordedTaskPaths, clean, sorted and deduplicated.
func DeclaredTaskPaths(content []byte, taskFile string) ([]string, error) {
	document, err := parseTaskDocument(content, taskFile)
	if err != nil {
		return nil, fmt.Errorf("parse Task %s: %w", taskFile, err)
	}
	seen := make(map[string]bool)
	for _, ref := range document.Context {
		if ref.Kind == ContextKindCreates || ref.Kind == ContextKindInterface {
			seen[ref.Path] = true
		}
	}
	for _, recorded := range RecordedTaskPaths(content) {
		clean, err := cleanTaskContextPath(filepath.ToSlash(filepath.Clean(recorded)))
		if err != nil {
			return nil, fmt.Errorf("read recorded path %q in Task %s: %w", recorded, taskFile, err)
		}
		seen[clean] = true
	}
	paths := make([]string, 0, len(seen))
	for name := range seen {
		paths = append(paths, name)
	}
	sort.Strings(paths)
	return paths, nil
}
