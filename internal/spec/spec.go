// Package spec owns the on-disk Spec contract: discovering active Specs
// under a Spec Root, loading and validating a Task Graph into deterministic
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

// Spec identifies one Spec directory under a Spec Root.
type Spec struct {
	Slug string
	Dir  string
}

// SkippedSpec reports a Spec directory that discovery intentionally did not
// offer to the Implement picker, with the reason a user can fix.
type SkippedSpec struct {
	Dir    string
	Reason string
}

// Task is one Task Graph node joined with its parsed task file. File is the
// task file path relative to the Spec Root.
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
// directories under the Spec Root (excluding _archived/) whose _prd.md
// frontmatter carries status active. Directories without a readable active
// PRD are skipped; Load names the exact problem when a slug is requested
// explicitly. The result is sorted by slug.
func ListActive(specsRoot string) ([]Spec, error) {
	specs, _, err := ListActiveDetailed(specsRoot)
	return specs, err
}

// ListActiveDetailed discovers active Specs and reports non-active Spec
// directories skipped because their PRD is missing, unreadable, or inactive.
// _archived/ is outside active discovery and is not reported as skipped.
func ListActiveDetailed(specsRoot string) ([]Spec, []SkippedSpec, error) {
	root := filepath.Clean(specsRoot)
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read Spec Root %q: %w", root, err)
	}
	var specs []Spec
	var skipped []SkippedSpec
	// os.ReadDir returns entries sorted by filename, so the result is
	// already sorted by slug.
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == archivedDirName {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		content, err := os.ReadFile(filepath.Join(dir, "_prd.md"))
		if err != nil {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: prdReadSkipReason(err),
			})
			continue
		}
		status, err := prdStatus(content)
		if err != nil {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: fmt.Sprintf("unreadable _prd.md frontmatter: %v", err),
			})
			continue
		}
		if status != prdStatusActive {
			skipped = append(skipped, SkippedSpec{
				Dir:    specDisplayDir(root, entry.Name()),
				Reason: fmt.Sprintf("status %q is not active", status),
			})
			continue
		}
		specs = append(specs, Spec{Slug: entry.Name(), Dir: dir})
	}
	return specs, skipped, nil
}

// Load parses and validates one Spec's Task Graph and task files from the
// Spec Root, returning the Tasks in deterministic topological order. Every
// validation failure is a typed error naming the offending Task or check.
func Load(specsRoot string, slug string) (*Graph, error) {
	root := filepath.Clean(specsRoot)
	dir := filepath.Join(root, slug)
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

func specDisplayDir(specsRoot string, slug string) string {
	if filepath.Base(specsRoot) == "specs" && filepath.Base(filepath.Dir(specsRoot)) == "docs" {
		return filepath.ToSlash(filepath.Join("docs", "specs", slug))
	}
	if filepath.IsAbs(specsRoot) {
		return filepath.ToSlash(filepath.Join(specsRoot, slug))
	}
	return filepath.ToSlash(filepath.Join("docs", "specs", slug))
}

func prdReadSkipReason(err error) string {
	if errors.Is(err, os.ErrNotExist) {
		return "missing _prd.md"
	}
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return fmt.Sprintf("read _prd.md: %v", pathErr.Err)
	}
	return fmt.Sprintf("read _prd.md: %v", err)
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
		File:         filepath.Join(slug, node.File),
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
