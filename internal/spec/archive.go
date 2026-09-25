package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// QAArchiveOverride records the maintainer authority and revision that permit
// archiving a Spec whose newest QA evidence does not qualify for normal
// archive. It never changes the QA Task or QA Report.
type QAArchiveOverride struct {
	Approval string
	Reason   string
	Revision string
}

// ArchiveKind names a retired artifact family.
type ArchiveKind string

const (
	ArchiveKindSpec    ArchiveKind = "specs"
	ArchiveKindFinding ArchiveKind = "findings"
	ArchiveKindADR     ArchiveKind = "adr"
	ArchiveKindBacklog ArchiveKind = "backlog"
	ArchiveKindReview  ArchiveKind = "reviews"
	ArchiveKindHandoff ArchiveKind = "handoffs"
)

// ArchiveKinds returns every retired artifact family understood by the Spec
// package.
func ArchiveKinds() []ArchiveKind {
	return []ArchiveKind{
		ArchiveKindSpec,
		ArchiveKindFinding,
		ArchiveKindADR,
		ArchiveKindBacklog,
		ArchiveKindReview,
		ArchiveKindHandoff,
	}
}

// ArchiveDir returns the repository-relative directory holding retired
// artifacts of kind. Unknown kinds return an empty directory.
func ArchiveDir(kind ArchiveKind) string {
	switch kind {
	case ArchiveKindSpec:
		return "docs/history/specs"
	case ArchiveKindFinding:
		return "docs/history/findings"
	case ArchiveKindADR:
		return "docs/history/adr"
	case ArchiveKindBacklog:
		return "docs/history/backlog"
	case ArchiveKindReview:
		return "docs/history/reviews"
	case ArchiveKindHandoff:
		return "docs/history/handoffs"
	default:
		return ""
	}
}

// ArchiveRequest asks the Spec package to retire one completed Spec.
type ArchiveRequest struct {
	SpecsRoot   string
	BuiltInRoot bool
	Slug        string
	ArchivedAt  time.Time
	QAOverride  *QAArchiveOverride
}

// ArchiveResult reports the filesystem paths touched by Archive.
type ArchiveResult struct {
	SourceDir   string
	ArchivedDir string
	ArchivedOn  string
	QAOverride  bool
}

// Archive verifies either completion and QA evidence for a Spec with a Task
// Graph or a supersession record for a Spec without one, then moves the Spec
// under the resolved archived Spec root. A partial QA Report is eligible only
// when its blocked rows are declared unreachable. Superseded Specs move
// byte-identically because their amendment already records their disposition.
func Archive(req ArchiveRequest) (ArchiveResult, error) {
	qaOverride, err := validateQAArchiveOverride(req.QAOverride)
	if err != nil {
		return ArchiveResult{}, err
	}
	sourceDir := filepath.Join(filepath.Clean(req.SpecsRoot), req.Slug)
	stampMetadata := true
	var unproven []string
	qaOverrideOutcome := ""
	qaOverrideQATaskStatus := ""

	graph, err := Load(req.SpecsRoot, req.Slug)
	if err != nil {
		var manifestErr ManifestError
		if !errors.As(err, &manifestErr) || manifestErr.Reason != missingManifestReason || manifestErr.Err != nil {
			return ArchiveResult{}, err
		}
		if qaOverride != nil {
			return ArchiveResult{}, fmt.Errorf("QA archive override requires a Task Graph: %w", err)
		}
		if _, supersessionErr := ReadSupersession(sourceDir); supersessionErr != nil {
			if errors.Is(supersessionErr, ErrNoSupersession) {
				return ArchiveResult{}, err
			}
			return ArchiveResult{}, fmt.Errorf("invalid supersession proof: %w", supersessionErr)
		}
		stampMetadata = false
	} else {
		sourceDir = graph.Spec.Dir
		if qaOverride != nil {
			allTasksCompleted := true
			for _, task := range graph.Tasks {
				if task.Status != StatusCompleted {
					allTasksCompleted = false
				}
				if task.Type == TaskTypeQA && task.Status != StatusCompleted {
					qaOverrideQATaskStatus = string(task.Status)
				}
				if task.Type != TaskTypeQA && task.Status != StatusCompleted {
					return ArchiveResult{}, fmt.Errorf("Task %q is %q; QA archive override requires every non-QA Task to be %q", task.ID, task.Status, StatusCompleted)
				}
			}
			report, reportErr := ReadQAReport(graph.Spec.Dir)
			switch {
			case reportErr == nil:
				if allTasksCompleted {
					if _, eligibilityErr := archiveUnprovenActions(graph.Spec.Dir, report); eligibilityErr == nil {
						return ArchiveResult{}, errors.New("QA archive override is not allowed because the Spec already qualifies for normal archive")
					}
				}
				qaOverrideOutcome = report.Verdict
			case errors.Is(reportErr, ErrNoQAReport):
				qaOverrideOutcome = "missing"
			default:
				qaOverrideOutcome = qaArchiveOverrideErrorOutcome(graph.Spec.Dir, reportErr)
			}
		} else {
			for _, task := range graph.Tasks {
				if task.Status != StatusCompleted {
					return ArchiveResult{}, fmt.Errorf("Task %q is %q; archive requires every Task to be %q", task.ID, task.Status, StatusCompleted)
				}
			}
			report, reportErr := ReadQAReport(graph.Spec.Dir)
			if reportErr != nil {
				return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", reportErr)
			}
			unproven, reportErr = archiveUnprovenActions(graph.Spec.Dir, report)
			if reportErr != nil {
				return ArchiveResult{}, fmt.Errorf("no passing QA verdict: %w", reportErr)
			}
		}
	}

	archiveRoot := ArchiveSpecRoot(req.SpecsRoot, req.BuiltInRoot)
	archivedDir := filepath.Join(archiveRoot, req.Slug)
	if _, err := os.Stat(archivedDir); err == nil {
		return ArchiveResult{}, fmt.Errorf("archived Spec destination %q already exists", archivedDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ArchiveResult{}, fmt.Errorf("stat archived Spec destination %q: %w", archivedDir, err)
	}

	archivedOn := archiveDate(req.ArchivedAt)
	if stampMetadata {
		prdPath := filepath.Join(sourceDir, "_prd.md")
		if err := stampArchiveMetadata(prdPath, req.Slug, archivedOn, unproven, qaOverride, qaOverrideOutcome, qaOverrideQATaskStatus); err != nil {
			return ArchiveResult{}, err
		}
	}
	if err := os.MkdirAll(archiveRoot, 0o755); err != nil {
		return ArchiveResult{}, fmt.Errorf("create archived Spec root %q: %w", archiveRoot, err)
	}
	if err := os.Rename(sourceDir, archivedDir); err != nil {
		return ArchiveResult{}, fmt.Errorf("move Spec %q to %q: %w", sourceDir, archivedDir, err)
	}
	return ArchiveResult{
		SourceDir:   sourceDir,
		ArchivedDir: archivedDir,
		ArchivedOn:  archivedOn,
		QAOverride:  qaOverride != nil,
	}, nil
}

func validateQAArchiveOverride(override *QAArchiveOverride) (*QAArchiveOverride, error) {
	if override == nil {
		return nil, nil
	}
	normalized := &QAArchiveOverride{
		Approval: strings.TrimSpace(override.Approval),
		Reason:   strings.TrimSpace(override.Reason),
		Revision: strings.TrimSpace(override.Revision),
	}
	if normalized.Approval == "" {
		return nil, errors.New("QA archive override requires an approval source")
	}
	if normalized.Reason == "" {
		return nil, errors.New("QA archive override requires a reason")
	}
	if normalized.Revision == "" {
		return nil, errors.New("QA archive override requires an archived revision")
	}
	return normalized, nil
}

func qaArchiveOverrideErrorOutcome(specDir string, err error) string {
	var reportErr QAReportError
	if !errors.As(err, &reportErr) {
		return err.Error()
	}
	relativeErr := reportErr.Err
	var pathErr *os.PathError
	if errors.As(relativeErr, &pathErr) {
		relativeErr = &os.PathError{
			Op:   pathErr.Op,
			Path: qaArchiveOverrideRelativePath(specDir, pathErr.Path),
			Err:  pathErr.Err,
		}
	}
	return QAReportError{
		Path: qaArchiveOverrideRelativePath(specDir, reportErr.Path),
		Err:  relativeErr,
	}.Error()
}

func qaArchiveOverrideRelativePath(specDir string, path string) string {
	relativePath, err := filepath.Rel(specDir, path)
	if err != nil || filepath.IsAbs(relativePath) {
		return filepath.Join("qa", filepath.Base(path))
	}
	return relativePath
}

// ArchiveSpecRoot returns the filesystem directory holding retired Specs for
// one configured Spec Root. The repository's built-in Spec Root uses the
// documentation history layout (ArchiveDir). An external or configured
// non-default Spec Root keeps its archive beside the active root.
func ArchiveSpecRoot(specsRoot string, builtInRoot bool) string {
	cleanSpecsRoot := filepath.Clean(specsRoot)
	if !builtInRoot {
		// An external or configured non-default Spec Root owns its archive
		// beside its active Specs.
		return filepath.Join(cleanSpecsRoot, archivedDirName)
	}
	docsRoot := filepath.Dir(cleanSpecsRoot)
	return filepath.Join(filepath.Dir(docsRoot), filepath.FromSlash(ArchiveDir(ArchiveKindSpec)))
}

func archiveUnprovenActions(specDir string, report QAReport) ([]string, error) {
	if err := QAReportEligibility(specDir, report); err != nil {
		return nil, err
	}
	if report.Verdict != VerdictPartial {
		return nil, nil
	}

	declarations, err := Unreachable(specDir)
	if err != nil {
		return nil, fmt.Errorf("read unreachable acceptance declarations: %w", err)
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

func stampArchiveMetadata(prdPath string, slug string, archivedOn string, unproven []string, qaOverride *QAArchiveOverride, qaOverrideOutcome string, qaOverrideQATaskStatus string) error {
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
	if qaOverride != nil {
		setArchiveFrontmatterNode(mapping, "qa_override", archiveBoolNode(true))
		setArchiveFrontmatterValue(mapping, "qa_override_approval", qaOverride.Approval)
		setArchiveFrontmatterValue(mapping, "qa_override_reason", qaOverride.Reason)
		setArchiveFrontmatterValue(mapping, "qa_override_qa_outcome", qaOverrideOutcome)
		if qaOverrideQATaskStatus != "" {
			setArchiveFrontmatterValue(mapping, "qa_override_qa_task_status", qaOverrideQATaskStatus)
		}
		setArchiveFrontmatterValue(mapping, "qa_override_revision", qaOverride.Revision)
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

func archiveBoolNode(value bool) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: fmt.Sprintf("%t", value)}
}

func archiveSequenceNode(values []string) *yaml.Node {
	node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		node.Content = append(node.Content, archiveScalarNode(value))
	}
	return node
}
