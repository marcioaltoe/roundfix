package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CauseGraph is the historical Task Graph projection used by cause reports.
// Unlike execution loading, it does not require an active PRD or Task status.
type CauseGraph struct {
	QATaskID string
	Tasks    map[string]string // Task ID to its graph-authored filename
	Dir      string
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

// CauseTaskText reads only the title and Overview of a graph node, bounded in
// bytes. A missing Task file is absent evidence; other read errors propagate.
func (g CauseGraph) CauseTaskText(id string, limit int) (string, error) {
	file, ok := g.Tasks[id]
	if !ok {
		return "", nil
	}
	data, err := os.ReadFile(filepath.Join(g.Dir, file))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read cause Task %q: %w", id, err)
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
