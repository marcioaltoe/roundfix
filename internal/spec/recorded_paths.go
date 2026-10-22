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
// not taskFile, not declared as interface:, creates:, or deletes: in task.Context, and
// not reported true by governed.
func UndeclaredTaskPaths(task Task, taskFile string, committed []string, governed func(string) bool) []string {
	excluded := map[string]bool{filepath.ToSlash(filepath.Clean(taskFile)): true}
	for _, ref := range task.Context {
		if ref.Kind == ContextKindInterface || ref.Kind == ContextKindCreates || ref.Kind == ContextKindDeletes {
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

	sections := recordedPathsSections(content)
	base := content
	if len(sections) > 0 {
		base = content[:sections[0].start]
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
	if len(sections) > 0 {
		cursor := sections[0].end
		for _, section := range sections[1:] {
			updated = append(updated, content[cursor:section.start]...)
			cursor = section.end
		}
		updated = append(updated, content[cursor:]...)
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
	var paths []string
	seen := make(map[string]bool)
	for _, bounds := range recordedPathsSections(content) {
		section := content[bounds.start:bounds.end]
		if len(section) > 0 && section[0] == '\n' {
			section = section[1:]
		}
		lines := strings.Split(string(section), "\n")
		for _, line := range lines[1:] {
			if !strings.HasPrefix(line, "- `") || !strings.HasSuffix(line, "`") {
				continue
			}
			path := strings.TrimSuffix(strings.TrimPrefix(line, "- `"), "`")
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths
}

type taskFileSection struct {
	heading    string
	start, end int
}

// recordedPathsSections returns, in file order, recorded sections followed
// only by Daemon-written settlement sections. Bounds include the newline
// before each heading and end at the next section's start or the file's end.
func recordedPathsSections(content []byte) []taskFileSection {
	var sections []taskFileSection
	inFence := false
	fence := ""
	for lineStart := 0; lineStart <= len(content); {
		lineEnd := bytes.IndexByte(content[lineStart:], '\n')
		if lineEnd < 0 {
			lineEnd = len(content)
		} else {
			lineEnd += lineStart
		}
		line := string(content[lineStart:lineEnd])
		trimmed := strings.TrimSpace(line)
		if marker, ok := markdownFenceMarker(trimmed); ok {
			if !inFence {
				inFence = true
				fence = marker
			} else if marker == fence {
				inFence = false
				fence = ""
			}
		} else if !inFence {
			if depth, _, ok := markdownHeading(trimmed); ok && depth == 2 {
				start := lineStart
				if start > 0 && content[start-1] == '\n' {
					start--
				}
				if len(sections) > 0 {
					sections[len(sections)-1].end = start
				}
				sections = append(sections, taskFileSection{heading: line, start: start, end: len(content)})
			}
		}
		if lineEnd == len(content) {
			break
		}
		lineStart = lineEnd + 1
	}
	first := len(sections)
settlementRun:
	for first > 0 {
		switch sections[first-1].heading {
		case RecordedPathsHeading, "## Result", "## Carry-forward provenance":
			first--
		default:
			break settlementRun
		}
	}
	var recorded []taskFileSection
	for _, section := range sections[first:] {
		if section.heading == RecordedPathsHeading {
			recorded = append(recorded, section)
		}
	}
	return recorded
}
