package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// ArchiveRequest asks the Spec package to retire one completed Spec.
type ArchiveRequest struct {
	GitRoot    string
	Slug       string
	ArchivedAt time.Time
}

// ArchiveResult reports the filesystem paths touched by Archive.
type ArchiveResult struct {
	SourceDir   string
	ArchivedDir string
	ArchivedOn  string
}

// Archive verifies completion and QA evidence, stamps archive metadata in the
// PRD frontmatter, and moves the Spec under docs/specs/_archived/.
func Archive(req ArchiveRequest) (ArchiveResult, error) {
	graph, err := Load(req.GitRoot, req.Slug)
	if err != nil {
		return ArchiveResult{}, err
	}
	for _, task := range graph.Tasks {
		if task.Status != StatusCompleted {
			return ArchiveResult{}, fmt.Errorf("Task %q is %q; archive requires every Task to be %q", task.ID, task.Status, StatusCompleted)
		}
	}
	verdict, err := QAVerdict(graph.Spec.Dir)
	if err != nil {
		if errors.Is(err, ErrNoQAReport) {
			return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", err)
		}
		return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", err)
	}
	if verdict != VerdictPass {
		return ArchiveResult{}, fmt.Errorf("no passing QA verdict: newest QA Report verdict is %q; expected %q", verdict, VerdictPass)
	}

	archiveRoot := filepath.Join(specsRoot(req.GitRoot), archivedDirName)
	archivedDir := filepath.Join(archiveRoot, req.Slug)
	if _, err := os.Stat(archivedDir); err == nil {
		return ArchiveResult{}, fmt.Errorf("archived Spec destination %q already exists", archivedDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ArchiveResult{}, fmt.Errorf("stat archived Spec destination %q: %w", archivedDir, err)
	}

	archivedOn := archiveDate(req.ArchivedAt)
	prdPath := filepath.Join(graph.Spec.Dir, "_prd.md")
	if err := stampArchiveMetadata(prdPath, req.Slug, archivedOn); err != nil {
		return ArchiveResult{}, err
	}
	if err := os.MkdirAll(archiveRoot, 0o755); err != nil {
		return ArchiveResult{}, fmt.Errorf("create archived Spec root %q: %w", archiveRoot, err)
	}
	if err := os.Rename(graph.Spec.Dir, archivedDir); err != nil {
		return ArchiveResult{}, fmt.Errorf("move Spec %q to %q: %w", graph.Spec.Dir, archivedDir, err)
	}
	return ArchiveResult{
		SourceDir:   graph.Spec.Dir,
		ArchivedDir: archivedDir,
		ArchivedOn:  archivedOn,
	}, nil
}

func archiveDate(value time.Time) string {
	if value.IsZero() {
		value = time.Now()
	}
	return value.Format("2006-01-02")
}

func stampArchiveMetadata(prdPath string, slug string, archivedOn string) error {
	content, err := os.ReadFile(prdPath)
	if err != nil {
		return fmt.Errorf("read Spec PRD %q: %w", prdPath, err)
	}
	frontmatterBytes, body, err := splitFrontmatter(content)
	if err != nil {
		return fmt.Errorf("parse Spec PRD %q: %w", prdPath, err)
	}
	var frontmatter yaml.Node
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return fmt.Errorf("parse Spec PRD %q frontmatter: %w", prdPath, err)
	}
	mapping := archiveFrontmatterMapping(&frontmatter)
	if mapping == nil {
		return fmt.Errorf("parse Spec PRD %q frontmatter: expected a YAML mapping", prdPath)
	}
	setArchiveFrontmatterValue(mapping, "status", "archived")
	setArchiveFrontmatterValue(mapping, "archived", archivedOn)
	setArchiveFrontmatterValue(mapping, "source_slug", slug)

	var encoded bytes.Buffer
	encoder := yaml.NewEncoder(&encoded)
	if err := encoder.Encode(&frontmatter); err != nil {
		return fmt.Errorf("encode Spec PRD %q frontmatter: %w", prdPath, err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("encode Spec PRD %q frontmatter: %w", prdPath, err)
	}
	next := append([]byte("---\n"), encoded.Bytes()...)
	next = append(next, []byte("---\n\n")...)
	next = append(next, body...)
	if err := os.WriteFile(prdPath, next, 0o644); err != nil {
		return fmt.Errorf("write Spec PRD %q: %w", prdPath, err)
	}
	return nil
}

func archiveFrontmatterMapping(document *yaml.Node) *yaml.Node {
	if document.Kind == yaml.DocumentNode && len(document.Content) == 1 {
		document = document.Content[0]
	}
	if document.Kind != yaml.MappingNode {
		return nil
	}
	return document
}

func setArchiveFrontmatterValue(mapping *yaml.Node, key string, value string) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			mapping.Content[index+1] = archiveScalarNode(value)
			return
		}
	}
	mapping.Content = append(mapping.Content, archiveScalarNode(key), archiveScalarNode(value))
}

func archiveScalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}
