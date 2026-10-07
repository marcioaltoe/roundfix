package spec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// causeTaskReadLimit bounds how much of a Task file is read to find its
// title and Overview, whatever the evidence limit the caller applies.
const causeTaskReadLimit = 64 << 10

// CauseGraph is the historical Task Graph projection used by cause reports.
// Unlike execution loading, it does not require an active PRD or Task status.
type CauseGraph struct {
	QATaskID string
	Tasks    map[string]string // Task ID to its graph-authored filename
	Dir      string
	readTask func(string) ([]byte, error)
}

// ReadCauseGraph reads graph membership and QA identity without execution
// eligibility checks, so archived and legacy graphs remain readable.
func ReadCauseGraph(root, slug string) (CauseGraph, error) {
	if filepath.Base(slug) != slug || slug == "." || slug == ".." {
		return CauseGraph{}, fmt.Errorf("unsafe Spec slug %q", slug)
	}
	dir := filepath.Join(root, slug)
	path := filepath.Join(dir, "_tasks.md")
	if _, err := os.Stat(path); err != nil {
		return CauseGraph{}, fmt.Errorf("inspect cause Task Graph: %w", err)
	}
	nodes, _, _, qa, _, err := loadManifestNodes(path)
	if err != nil {
		return CauseGraph{}, err
	}
	graph := CauseGraph{QATaskID: qa.TaskID, Tasks: map[string]string{}, Dir: dir}
	for _, node := range nodes {
		graph.Tasks[node.ID] = node.File
	}
	return graph, nil
}

// ReadCauseGraphAt reads the historical graph and its Tasks through the same
// immutable tree reader. It shares manifest validation with the folder form.
func ReadCauseGraphAt(source, slug string, read func(string) ([]byte, error)) (CauseGraph, error) {
	manifestPath := filepath.Join(source, "_tasks.md")
	content, err := read(filepath.ToSlash(manifestPath))
	if err != nil {
		return CauseGraph{}, fmt.Errorf("read archived cause Task Graph %q: %w", slug, err)
	}
	nodes, _, _, qa, _, err := parseManifestNodes(manifestPath, content)
	if err != nil {
		return CauseGraph{}, err
	}
	graph := CauseGraph{QATaskID: qa.TaskID, Tasks: map[string]string{}, Dir: source, readTask: read}
	for _, node := range nodes {
		graph.Tasks[node.ID] = node.File
	}
	return graph, nil
}

// CauseTaskText reads only the title and Overview of a graph node, bounded in
// bytes. A missing Task file is absent evidence; other read errors propagate.
func (g CauseGraph) CauseTaskText(id string, limit int) (string, error) {
	file, ok := g.Tasks[id]
	if !ok {
		return "", nil
	}
	// Only a plain file directly in the Spec directory is evidence: a graph
	// that names a path, a symbolic link or a non-regular file gives none.
	if file == "" || filepath.Base(file) != file || file == "." || file == ".." {
		return "", nil
	}
	var data []byte
	if g.readTask != nil {
		var err error
		data, err = g.readTask(filepath.ToSlash(filepath.Join(g.Dir, file)))
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("read cause Task %q: %w", id, err)
		}
		if len(data) > causeTaskReadLimit {
			data = data[:causeTaskReadLimit]
		}
	} else {
		path := filepath.Join(g.Dir, file)
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return "", nil
		}
		if err != nil {
			return "", fmt.Errorf("inspect cause Task %q: %w", id, err)
		}
		if !info.Mode().IsRegular() {
			return "", nil
		}
		handle, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("read cause Task %q: %w", id, err)
		}
		defer handle.Close()
		var readErr error
		data, readErr = io.ReadAll(io.LimitReader(handle, causeTaskReadLimit))
		if readErr != nil {
			return "", fmt.Errorf("read cause Task %q: %w", id, readErr)
		}
	}
	var title string
	var overview []string
	inOverview := false
	for _, line := range strings.Split(string(data), "\n") {
		if title == "" && strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
		}
		if strings.HasPrefix(line, "## ") {
			if inOverview {
				break
			}
			inOverview = strings.TrimSpace(line) == "## Overview"
			continue
		}
		if inOverview {
			overview = append(overview, line)
		}
	}
	text := strings.TrimSpace(title + "\n" + strings.Join(overview, "\n"))
	if len(text) > limit {
		text = text[:limit]
	}
	return text, nil
}
