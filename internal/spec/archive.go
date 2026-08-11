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

// ArchiveKind names a retired artifact family.
type ArchiveKind string

const (
	ArchiveKindSpec    ArchiveKind = "specs"
	ArchiveKindFinding ArchiveKind = "findings"
	ArchiveKindADR     ArchiveKind = "adr"
	ArchiveKindBacklog ArchiveKind = "backlog"
)

// ArchiveDir returns the repository-relative directory holding retired
// artifacts of kind. Unknown kinds return an empty directory.
func ArchiveDir(kind ArchiveKind) string {
	switch kind {
	case ArchiveKindSpec:
		return "_archived/specs"
	case ArchiveKindFinding:
		return "_archived/findings"
	case ArchiveKindADR:
		return "_archived/adr"
	case ArchiveKindBacklog:
		return "_archived/backlog"
	default:
		return ""
	}
}

// ArchiveRequest asks the Spec package to retire one completed Spec.
type ArchiveRequest struct {
	SpecsRoot  string
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
// PRD frontmatter, and moves the Spec under the resolved archived Spec root. A
// partial QA Report is eligible only when its blocked rows are declared
// unreachable.
func Archive(req ArchiveRequest) (ArchiveResult, error) {
	graph, err := Load(req.SpecsRoot, req.Slug)
	if err != nil {
		return ArchiveResult{}, err
	}
	for _, task := range graph.Tasks {
		if task.Status != StatusCompleted {
			return ArchiveResult{}, fmt.Errorf("Task %q is %q; archive requires every Task to be %q", task.ID, task.Status, StatusCompleted)
		}
	}
	report, err := ReadQAReport(graph.Spec.Dir)
	if err != nil {
		if errors.Is(err, ErrNoQAReport) {
			return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", err)
		}
		return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", err)
	}
	unproven, err := archiveUnprovenActions(graph.Spec.Dir, report)
	if err != nil {
		return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", err)
	}

	archiveRoot := ArchiveSpecRoot(req.SpecsRoot)
	archivedDir := filepath.Join(archiveRoot, req.Slug)
	if _, err := os.Stat(archivedDir); err == nil {
		return ArchiveResult{}, fmt.Errorf("archived Spec destination %q already exists", archivedDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ArchiveResult{}, fmt.Errorf("stat archived Spec destination %q: %w", archivedDir, err)
	}

	archivedOn := archiveDate(req.ArchivedAt)
	prdPath := filepath.Join(graph.Spec.Dir, "_prd.md")
	if err := stampArchiveMetadata(prdPath, req.Slug, archivedOn, unproven); err != nil {
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

// ArchiveSpecRoot returns the filesystem directory holding retired Specs for
// one configured Spec Root. The default repository layout uses ArchiveDir;
// external Spec Roots keep their archive beside the active root.
func ArchiveSpecRoot(specsRoot string) string {
	cleanSpecsRoot := filepath.Clean(specsRoot)
	docsRoot := filepath.Dir(cleanSpecsRoot)
	if filepath.Base(cleanSpecsRoot) != string(ArchiveKindSpec) || filepath.Base(docsRoot) != "docs" {
		// A configured non-default Spec Root owns its archive beside its active
		// Specs; this request does not carry a repository root for resolving it.
		return filepath.Join(cleanSpecsRoot, archivedDirName)
	}
	return filepath.Join(filepath.Dir(docsRoot), filepath.FromSlash(ArchiveDir(ArchiveKindSpec)))
}

func archiveUnprovenActions(specDir string, report QAReport) ([]string, error) {
	if report.Verdict == VerdictPass {
		return nil, nil
	}
	if report.Verdict != VerdictPartial {
		return nil, fmt.Errorf("newest QA Report verdict is %q; expected %q", report.Verdict, VerdictPass)
	}
	if report.RowsBlockedFinding > 0 {
		return nil, fmt.Errorf("rows_blocked_finding is %d; expected 0", report.RowsBlockedFinding)
	}
	if report.RowsBlockedEnvironment > 0 {
		return nil, fmt.Errorf("rows_blocked_environment is %d; expected 0", report.RowsBlockedEnvironment)
	}
	if report.RowsBlockedDeclared == 0 {
		return nil, fmt.Errorf("newest QA Report verdict is %q; expected %q", report.Verdict, VerdictPass)
	}

	declarations, err := Unreachable(specDir)
	if err != nil {
		return nil, fmt.Errorf("read unreachable acceptance declarations: %w", err)
	}
	if report.RowsBlockedDeclared > len(declarations) {
		plural := ""
		if len(declarations) != 1 {
			plural = "s"
		}
		return nil, fmt.Errorf(
			"rows_blocked_declared is %d, but Spec declares %d unreachable acceptance%s; shortfall is %d",
			report.RowsBlockedDeclared,
			len(declarations),
			plural,
			report.RowsBlockedDeclared-len(declarations),
		)
	}

	actions := make([]string, 0, len(declarations))
	for _, declaration := range declarations {
		actions = append(actions, declaration.SatisfiedBy)
	}
	return actions, nil
}

func archiveDate(value time.Time) string {
	if value.IsZero() {
		value = time.Now()
	}
	return value.Format("2006-01-02")
}

func stampArchiveMetadata(prdPath string, slug string, archivedOn string, unproven []string) error {
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
	if len(unproven) > 0 {
		setArchiveFrontmatterNode(mapping, "unproven", archiveSequenceNode(unproven))
	}

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
	setArchiveFrontmatterNode(mapping, key, archiveScalarNode(value))
}

func setArchiveFrontmatterNode(mapping *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			mapping.Content[index+1] = value
			return
		}
	}
	mapping.Content = append(mapping.Content, archiveScalarNode(key), value)
}

func archiveScalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func archiveSequenceNode(values []string) *yaml.Node {
	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		node.Content = append(node.Content, archiveScalarNode(value))
	}
	return node
}
