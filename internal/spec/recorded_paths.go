package spec

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RecordedPathsHeading names the Daemon-owned Task file section (ADR-0166).
const RecordedPathsHeading = "## Recorded paths"

const recordedPathsExplanation = "The Daemon recorded these paths, which this Task changed without declaring them in `## Context`."

// UndeclaredTaskPaths returns, sorted and unique, the committed paths that are
// not taskFile, not declared as interface: or creates: in task.Context, and
// not reported true by governed.
func UndeclaredTaskPaths(task Task, taskFile string, committed []string, governed func(string) bool) []string {
	excluded := map[string]bool{filepath.ToSlash(filepath.Clean(taskFile)): true}
	for _, ref := range task.Context {
		if ref.Kind == ContextKindInterface || ref.Kind == ContextKindCreates {
			excluded[filepath.ToSlash(filepath.Clean(ref.Path))] = true
		}
	}

	seen := make(map[string]bool, len(committed))
	paths := make([]string, 0, len(committed))
	for _, path := range committed {
		path = filepath.ToSlash(filepath.Clean(path))
		if excluded[path] || seen[path] || governed != nil && governed(path) {
			continue
		}
		seen[path] = true
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

// RecordTaskPaths writes paths as the Task file's recorded section, replacing
// an existing one and preserving every other byte. No paths removes an
// existing section and writes nothing else.
func RecordTaskPaths(taskPath string, paths []string) error {
	for _, path := range paths {
		if strings.ContainsAny(path, "`\r\n") {
			return fmt.Errorf("record Task path %q: path cannot be recorded because it contains a backtick or line break", path)
		}
	}

	info, err := os.Stat(taskPath)
	if err != nil {
		return fmt.Errorf("stat Task file %q before recording paths: %w", taskPath, err)
	}
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return fmt.Errorf("read Task file %q before recording paths: %w", taskPath, err)
	}

	base := content
	if offset, ok := recordedPathsSectionOffset(content); ok {
		base = content[:offset]
	}
	updated := append([]byte(nil), base...)
	if len(paths) > 0 {
		updated = append(updated, '\n')
		updated = append(updated, RecordedPathsHeading...)
		updated = append(updated, "\n\n"...)
		updated = append(updated, recordedPathsExplanation...)
		updated = append(updated, "\n\n"...)
		for _, path := range paths {
			updated = append(updated, "- `"...)
			updated = append(updated, path...)
			updated = append(updated, "`\n"...)
		}
	}
	if bytes.Equal(updated, content) {
		return nil
	}
	if err := replaceTaskFile(taskPath, updated, info.Mode().Perm()); err != nil {
		return fmt.Errorf("replace Task file %q after recording paths: %w", taskPath, err)
	}
	return nil
}

// RecordedTaskPaths reads the paths from the Daemon-owned recorded section.
func RecordedTaskPaths(content []byte) []string {
	offset, ok := recordedPathsSectionOffset(content)
	if !ok {
		return nil
	}
	section := content[offset:]
	if len(section) > 0 && section[0] == '\n' {
		section = section[1:]
	}
	lines := strings.Split(string(section), "\n")
	paths := make([]string, 0)
	for _, line := range lines[1:] {
		if !strings.HasPrefix(line, "- `") || !strings.HasSuffix(line, "`") {
			continue
		}
		paths = append(paths, strings.TrimSuffix(strings.TrimPrefix(line, "- `"), "`"))
	}
	return paths
}

func recordedPathsSectionOffset(content []byte) (int, bool) {
	heading := []byte(RecordedPathsHeading + "\n")
	if bytes.HasPrefix(content, heading) {
		return 0, true
	}
	marker := append([]byte{'\n'}, heading...)
	if offset := bytes.Index(content, marker); offset >= 0 {
		return offset, true
	}
	return 0, false
}
