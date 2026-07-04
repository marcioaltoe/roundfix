// Package spec owns the on-disk Spec contract: discovering active Specs
// under docs/specs/, loading and validating a Task Graph into deterministic
// topological order, parsing task files, rewriting task status, and reading
// the QA Report verdict. Nothing outside this package touches spec markdown.
package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Task status values; the task file is the sole owner of its status.
const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

const (
	manifestSchema  = "spec-tasks/v1"
	prdStatusActive = "active"
	archivedDirName = "_archived"
)

// Status is the lifecycle state a task file carries in its frontmatter.
type Status string

// AllowedStatus reports whether the status is one of the four task lifecycle
// states.
func AllowedStatus(status Status) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// Spec identifies one Spec directory under docs/specs/.
type Spec struct {
	Slug string
	Dir  string
}

// Task is one Task Graph node joined with its parsed task file. File is the
// task file path relative to the repository root, so it can double as a git
// pathspec and be re-resolved with ReloadTask.
type Task struct {
	ID           string
	File         string
	Title        string
	Needs        []string
	Status       Status
	Type         string
	Verification []string
}

// Graph is a validated Task Graph with Tasks in deterministic topological
// order (Kahn's algorithm, manifest node order as tiebreak).
type Graph struct {
	Spec  Spec
	Tasks []Task
}

type manifestNode struct {
	ID    string   `yaml:"id"`
	File  string   `yaml:"file"`
	Needs []string `yaml:"needs"`
}

type manifestFrontmatter struct {
	Schema string `yaml:"schema"`
	Graph  struct {
		Nodes []manifestNode `yaml:"nodes"`
	} `yaml:"graph"`
}

// ListActive discovers the Specs eligible for the Implement Command:
// directories under docs/specs/ (excluding _archived/) whose _prd.md
// frontmatter carries status active. Directories without a readable active
// PRD are skipped; Load names the exact problem when a slug is requested
// explicitly. The result is sorted by slug.
func ListActive(gitRoot string) ([]Spec, error) {
	specsDir := specsRoot(gitRoot)
	entries, err := os.ReadDir(specsDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Spec root %q: %w", specsDir, err)
	}
	var specs []Spec
	// os.ReadDir returns entries sorted by filename, so the result is
	// already sorted by slug.
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == archivedDirName {
			continue
		}
		dir := filepath.Join(specsDir, entry.Name())
		content, err := os.ReadFile(filepath.Join(dir, "_prd.md"))
		if err != nil {
			continue
		}
		status, err := prdStatus(content)
		if err != nil || status != prdStatusActive {
			continue
		}
		specs = append(specs, Spec{Slug: entry.Name(), Dir: dir})
	}
	return specs, nil
}

// Load parses and validates one Spec's Task Graph and task files, returning
// the Tasks in deterministic topological order. Every validation failure is a
// typed error naming the offending Task or check.
func Load(gitRoot string, slug string) (*Graph, error) {
	dir := filepath.Join(specsRoot(gitRoot), slug)
	if err := requireActive(slug, dir); err != nil {
		return nil, err
	}

	manifestPath := filepath.Join(dir, "_tasks.md")
	nodes, err := loadManifestNodes(manifestPath)
	if err != nil {
		return nil, err
	}

	order, err := topologicalOrder(nodes)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(nodes))
	for _, index := range order {
		task, err := loadTask(dir, slug, nodes[index])
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &Graph{Spec: Spec{Slug: slug, Dir: dir}, Tasks: tasks}, nil
}

func specsRoot(gitRoot string) string {
	return filepath.Join(gitRoot, "docs", "specs")
}

func requireActive(slug string, dir string) error {
	prdPath := filepath.Join(dir, "_prd.md")
	content, err := os.ReadFile(prdPath)
	if errors.Is(err, os.ErrNotExist) {
		return SpecNotFoundError{Slug: slug, Dir: dir}
	}
	if err != nil {
		return fmt.Errorf("read Spec PRD %q: %w", prdPath, err)
	}
	status, err := prdStatus(content)
	if err != nil {
		return fmt.Errorf("parse Spec PRD %q: %w", prdPath, err)
	}
	if status != prdStatusActive {
		return InactiveSpecError{Slug: slug, Status: status}
	}
	return nil
}

func prdStatus(content []byte) (string, error) {
	frontmatterBytes, _, err := splitFrontmatter(content)
	if err != nil {
		return "", err
	}
	var prd struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal(frontmatterBytes, &prd); err != nil {
		return "", err
	}
	return prd.Status, nil
}

func loadManifestNodes(manifestPath string) ([]manifestNode, error) {
	content, err := os.ReadFile(manifestPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ManifestError{Path: manifestPath, Reason: "file does not exist; run the write-tasks workflow to create the Task Graph"}
	}
	if err != nil {
		return nil, fmt.Errorf("read Task Graph manifest %q: %w", manifestPath, err)
	}
	frontmatterBytes, _, err := splitFrontmatter(content)
	if err != nil {
		return nil, ManifestError{Path: manifestPath, Reason: "invalid frontmatter", Err: err}
	}
	var manifest manifestFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &manifest); err != nil {
		return nil, ManifestError{Path: manifestPath, Reason: "invalid frontmatter", Err: err}
	}
	if manifest.Schema != manifestSchema {
		return nil, ManifestSchemaError{Path: manifestPath, Schema: manifest.Schema}
	}
	nodes := manifest.Graph.Nodes
	if len(nodes) == 0 {
		return nil, ManifestError{Path: manifestPath, Reason: "graph has no nodes"}
	}

	known := make(map[string]bool, len(nodes))
	for index, node := range nodes {
		if node.ID == "" || node.File == "" {
			return nil, ManifestError{Path: manifestPath, Reason: fmt.Sprintf("graph node %d is missing id or file", index+1)}
		}
		if known[node.ID] {
			return nil, ManifestError{Path: manifestPath, Reason: fmt.Sprintf("duplicate Task id %q", node.ID)}
		}
		known[node.ID] = true
	}
	for _, node := range nodes {
		for _, need := range node.Needs {
			if !known[need] {
				return nil, UnknownNeedError{TaskID: node.ID, Need: need}
			}
		}
	}
	return nodes, nil
}

// topologicalOrder runs Kahn's algorithm over the manifest nodes, always
// picking the first ready node in manifest order so the same manifest yields
// the same order on every run.
func topologicalOrder(nodes []manifestNode) ([]int, error) {
	indexByID := make(map[string]int, len(nodes))
	for index, node := range nodes {
		indexByID[node.ID] = index
	}
	indegree := make([]int, len(nodes))
	dependents := make([][]int, len(nodes))
	for index, node := range nodes {
		for _, need := range node.Needs {
			needIndex := indexByID[need]
			dependents[needIndex] = append(dependents[needIndex], index)
			indegree[index]++
		}
	}

	placed := make([]bool, len(nodes))
	order := make([]int, 0, len(nodes))
	for len(order) < len(nodes) {
		next := -1
		for index := range nodes {
			if !placed[index] && indegree[index] == 0 {
				next = index
				break
			}
		}
		if next < 0 {
			var remaining []string
			for index, node := range nodes {
				if !placed[index] {
					remaining = append(remaining, node.ID)
				}
			}
			return nil, CycleError{TaskIDs: remaining}
		}
		placed[next] = true
		order = append(order, next)
		for _, dependent := range dependents[next] {
			indegree[dependent]--
		}
	}
	return order, nil
}

func loadTask(dir string, slug string, node manifestNode) (Task, error) {
	path := filepath.Join(dir, node.File)
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Task{}, MissingTaskFileError{TaskID: node.ID, Path: path}
	}
	if err != nil {
		return Task{}, fmt.Errorf("read Task %q file %q: %w", node.ID, path, err)
	}
	document, err := parseTaskDocument(content)
	if err != nil {
		return Task{}, TaskFileError{TaskID: node.ID, Path: path, Err: err}
	}
	if len(document.Verification) == 0 {
		return Task{}, MissingVerificationError{TaskID: node.ID, Path: path}
	}
	return Task{
		ID:           node.ID,
		File:         filepath.Join("docs", "specs", slug, node.File),
		Title:        document.Title,
		Needs:        append([]string(nil), node.Needs...),
		Status:       Status(document.Frontmatter.Status),
		Type:         document.Frontmatter.Type,
		Verification: document.Verification,
	}, nil
}

// splitFrontmatter mirrors the Round artifact frontmatter contract: an
// opening --- line, YAML, and a closing --- line ahead of the markdown body.
func splitFrontmatter(content []byte) ([]byte, []byte, error) {
	text := string(content)
	const opening = "---\n"
	if !strings.HasPrefix(text, opening) {
		return nil, nil, errors.New("missing YAML frontmatter opening marker")
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return nil, nil, errors.New("missing YAML frontmatter closing marker")
	}
	bodyStart := end + len("\n---")
	if strings.HasPrefix(rest[bodyStart:], "\n") {
		bodyStart++
	}
	return []byte(rest[:end]), []byte(rest[bodyStart:]), nil
}
