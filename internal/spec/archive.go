package spec

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
	SpecsRoot      string
	BuiltInRoot    bool
	Slug           string
	ArchivedAt     time.Time
	QAOverride     *QAArchiveOverride
	SourceRevision string
	Promote        []string
	RepositoryRoot string
}

// ArchiveResult reports the filesystem paths touched by Archive.
type ArchiveResult struct {
	SourceDir      string
	ArchivedDir    string
	ArchivedOn     string
	QAOverride     bool
	RewrittenLinks int
	RecordPath     string
	RemovedFiles   int
	RemovedBytes   int64
	Promoted       []string
}

// Archive keeps the existing completion and QA eligibility policy, then writes
// an Archive Record and removes the Spec folder. Its bytes remain at the
// caller's committed SourceRevision; no PRD or link is rewritten.
func Archive(req ArchiveRequest) (ArchiveResult, error) {
	qaOverride, err := validateQAArchiveOverride(req.QAOverride)
	if err != nil {
		return ArchiveResult{}, err
	}
	sourceDir := filepath.Join(filepath.Clean(req.SpecsRoot), req.Slug)
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
	if _, err := os.Lstat(archivedDir); err == nil {
		return ArchiveResult{}, fmt.Errorf("archived Spec destination %q already exists", archivedDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ArchiveResult{}, fmt.Errorf("stat archived Spec destination %q: %w", archivedDir, err)
	}

	if _, _, err := prepareArchiveLinks(sourceDir, archivedDir, req.Slug); err != nil {
		return ArchiveResult{}, err
	}
	if !regexp.MustCompile(`^[0-9a-fA-F]{40}$`).MatchString(req.SourceRevision) {
		return ArchiveResult{}, errors.New("archive requires a 40-hex SourceRevision")
	}
	recordPath := ArchiveRecordPath(archiveRoot, req.Slug)
	if _, err := os.Lstat(recordPath); err == nil {
		return ArchiveResult{}, fmt.Errorf("archived Spec destination %q already exists", recordPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return ArchiveResult{}, err
	}
	var promoted []string
	copies := make(map[string][]byte)
	for _, candidate := range req.Promote {
		clean := filepath.Clean(candidate)
		if filepath.IsAbs(candidate) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return ArchiveResult{}, fmt.Errorf("promoted path %q is outside the Spec", candidate)
		}
		base := filepath.Base(clean)
		if (filepath.Dir(clean) == "." && strings.HasPrefix(base, "_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md")) {
			return ArchiveResult{}, fmt.Errorf("cannot promote core artifact %q", candidate)
		}
		filePath := filepath.Join(sourceDir, clean)
		// Reject symlinks in every component, including directory components.
		current := sourceDir
		for _, component := range strings.Split(clean, string(filepath.Separator)) {
			current = filepath.Join(current, component)
			info, err := os.Lstat(current)
			if err != nil {
				return ArchiveResult{}, err
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return ArchiveResult{}, fmt.Errorf("promoted path %q is not a regular file", candidate)
			}
		}
		info, err := os.Lstat(filePath)
		if err != nil {
			return ArchiveResult{}, err
		}
		if !info.Mode().IsRegular() {
			return ArchiveResult{}, fmt.Errorf("promoted path %q is not a regular file", candidate)
		}
		if req.RepositoryRoot == "" {
			return ArchiveResult{}, errors.New("promotion requires RepositoryRoot")
		}
		destination := filepath.Join(req.RepositoryRoot, "docs", "references", base)
		if _, err := os.Lstat(destination); err == nil {
			return ArchiveResult{}, fmt.Errorf("promotion destination %q already exists", destination)
		} else if !errors.Is(err, os.ErrNotExist) {
			return ArchiveResult{}, err
		}
		if _, exists := copies[destination]; exists {
			return ArchiveResult{}, fmt.Errorf("duplicate promotion destination %q", destination)
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return ArchiveResult{}, err
		}
		copies[destination] = content
		promoted = append(promoted, "docs/references/"+base)
	}
	archivedOn := archiveDate(req.ArchivedAt)
	repositoryRoot := req.RepositoryRoot
	if repositoryRoot == "" {
		repositoryRoot = filepath.Dir(req.SpecsRoot)
		if req.BuiltInRoot {
			repositoryRoot = filepath.Dir(repositoryRoot)
		}
	}
	canonicalRoot, err := filepath.EvalSymlinks(repositoryRoot)
	if err != nil {
		return ArchiveResult{}, err
	}
	canonicalSource, err := filepath.EvalSymlinks(sourceDir)
	if err != nil {
		return ArchiveResult{}, err
	}
	source, err := filepath.Rel(canonicalRoot, canonicalSource)
	if err != nil {
		return ArchiveResult{}, err
	}
	var overrideRecord *QAArchiveOverrideRecord
	if qaOverride != nil {
		overrideRecord = &QAArchiveOverrideRecord{qaOverride.Approval, qaOverride.Reason, qaOverrideOutcome, qaOverrideQATaskStatus, qaOverride.Revision}
	}
	record, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: sourceDir, Slug: req.Slug, Source: filepath.ToSlash(source), SourceRevision: req.SourceRevision, Archived: archivedOn, Unproven: unproven, QAOverride: overrideRecord, Promoted: promoted})
	if err != nil {
		return ArchiveResult{}, err
	}
	content, err := RenderArchiveRecord(record)
	if err != nil {
		return ArchiveResult{}, err
	}
	result := ArchiveResult{SourceDir: sourceDir, ArchivedDir: archivedDir, ArchivedOn: archivedOn, QAOverride: qaOverride != nil, RecordPath: recordPath, Promoted: promoted}
	if err := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		result.RemovedFiles++
		result.RemovedBytes += info.Size()
		return nil
	}); err != nil {
		return ArchiveResult{}, err
	}
	var written []string
	rollback := func(cause error) error {
		for _, path := range written {
			cause = errors.Join(cause, os.Remove(path))
		}
		return cause
	}
	write := func(path string, data []byte) error {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		written = append(written, path)
		_, writeErr := file.Write(data)
		return errors.Join(writeErr, file.Close())
	}
	if err := write(recordPath, content); err != nil {
		return ArchiveResult{}, rollback(fmt.Errorf("write archive record: %w", err))
	}
	for _, destination := range promoted {
		path := filepath.Join(req.RepositoryRoot, filepath.FromSlash(destination))
		if err := write(path, copies[path]); err != nil {
			return ArchiveResult{}, rollback(fmt.Errorf("copy promotion: %w", err))
		}
	}
	if err := os.RemoveAll(sourceDir); err != nil {
		return ArchiveResult{}, fmt.Errorf("remove Spec; leftover folder %q: %w", sourceDir, err)
	}
	return result, nil
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

// markdownDestination indexes only a destination's bytes, preserving its wrapper and title.
type markdownDestination struct{ start, end int }
type archiveLinkRewrite struct {
	path         string
	original     []byte
	mode         fs.FileMode
	destinations map[string]string
}

var archiveLinkScheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

func archiveLinkTarget(destination, dir string) (target, suffix string, relative bool) {
	if destination == "" || strings.HasPrefix(destination, "/") || strings.HasPrefix(destination, "#") || archiveLinkScheme.MatchString(destination) {
		return "", "", false
	}
	path := destination
	if index := strings.IndexAny(path, "?#"); index >= 0 {
		path, suffix = path[:index], path[index:]
	}
	if path == "" {
		return "", "", false
	}
	decoded, err := url.PathUnescape(path)
	if err != nil {
		return "", "", false
	}
	return filepath.Clean(filepath.Join(dir, filepath.FromSlash(decoded))), suffix, true
}

func archiveLinkInside(target, specDir string) bool {
	rel, err := filepath.Rel(specDir, target)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

func prepareArchiveLinks(sourceDir, archivedDir, slug string) ([]archiveLinkRewrite, int, error) {
	var rewrites []archiveLinkRewrite
	var broken []string
	count := 0
	err := filepath.WalkDir(sourceDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.Type().IsRegular() || filepath.Ext(path) != ".md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read Markdown %q: %w", path, err)
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		activeDir := filepath.Dir(path)
		movedDir := filepath.Dir(filepath.Join(archivedDir, rel))
		replacements := make(map[string]string)
		for _, span := range scanArchiveMarkdown(content) {
			destination := string(content[span.start:span.end])
			target, suffix, relative := archiveLinkTarget(destination, activeDir)
			if !relative || archiveLinkInside(target, sourceDir) {
				continue
			}
			if _, err := os.Stat(target); err == nil {
				next, err := filepath.Rel(movedDir, target)
				if err != nil {
					return err
				}
				encoded := (&url.URL{Path: filepath.ToSlash(next)}).EscapedPath() + suffix
				if encoded != destination {
					replacements[destination] = encoded
					count++
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("stat link target %q: %w", target, err)
			} else {
				movedTarget, _, _ := archiveLinkTarget(destination, movedDir)
				if _, err := os.Stat(movedTarget); errors.Is(err, os.ErrNotExist) {
					line := bytes.Count(content[:span.start], []byte("\n")) + 1
					broken = append(broken, fmt.Sprintf("%s:%d %q", filepath.ToSlash(rel), line, destination))
				} else if err != nil {
					return fmt.Errorf("stat archived link target %q: %w", movedTarget, err)
				}
			}
		}
		if len(replacements) > 0 {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			rewrites = append(rewrites, archiveLinkRewrite{path: path, original: content, mode: info.Mode(), destinations: replacements})
		}
		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("scan Spec Markdown links: %w", err)
	}
	if len(broken) > 0 {
		return nil, 0, fmt.Errorf("Spec %q has relative links that leave the Spec and do not resolve: %s; fix or remove each link, then retry the archive", slug, strings.Join(broken, ", "))
	}
	return rewrites, count, nil
}

func rewriteArchiveLinks(content []byte, replacements map[string]string) []byte {
	var next bytes.Buffer
	end := 0
	for _, span := range scanArchiveMarkdown(content) {
		next.Write(content[end:span.start])
		destination := string(content[span.start:span.end])
		if replacement, ok := replacements[destination]; ok {
			next.WriteString(replacement)
		} else {
			next.WriteString(destination)
		}
		end = span.end
	}
	next.Write(content[end:])
	return next.Bytes()
}

func restoreArchiveLinks(rewrites []archiveLinkRewrite) error {
	var errs []error
	for _, rewrite := range rewrites {
		if err := os.WriteFile(rewrite.path, rewrite.original, rewrite.mode); err != nil {
			errs = append(errs, fmt.Errorf("restore link rewrite %q: %w", rewrite.path, err))
		}
	}
	return errors.Join(errs...)
}

// ArchiveLinksMatch reports whether only outward link destinations changed,
// reaching the same lexical target and preserving their query and fragment.
func ArchiveLinksMatch(active, archived []byte, activeDir, archivedDir, specDir string) bool {
	left, right := scanArchiveMarkdown(active), scanArchiveMarkdown(archived)
	if len(left) != len(right) {
		return false
	}
	aEnd, bEnd := 0, 0
	for i, a := range left {
		b := right[i]
		if !bytes.Equal(active[aEnd:a.start], archived[bEnd:b.start]) {
			return false
		}
		before, after := string(active[a.start:a.end]), string(archived[b.start:b.end])
		if before != after {
			t1, s1, r1 := archiveLinkTarget(before, activeDir)
			t2, s2, r2 := archiveLinkTarget(after, archivedDir)
			if !r1 || !r2 || archiveLinkInside(t1, specDir) || t1 != t2 || s1 != s2 {
				return false
			}
		}
		aEnd, bEnd = a.end, b.end
	}
	return bytes.Equal(active[aEnd:], archived[bEnd:])
}

func markdownSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }

// scanArchiveMarkdown masks code first, then reads inline and reference
// destinations without reformatting any surrounding Markdown.
func scanArchiveMarkdown(content []byte) []markdownDestination {
	masked := bytes.Clone(content)
	fence := byte(0)
	fenceLength := 0
	for start := 0; start < len(masked); {
		end := bytes.IndexByte(masked[start:], '\n')
		if end < 0 {
			end = len(masked)
		} else {
			end += start
		}
		line := masked[start:end]
		indent := 0
		for indent < len(line) && indent < 4 && line[indent] == ' ' {
			indent++
		}
		run := 0
		if indent < 4 && indent < len(line) && (line[indent] == '`' || line[indent] == '~') {
			for indent+run < len(line) && line[indent+run] == line[indent] {
				run++
			}
		}
		if fence != 0 {
			if run >= fenceLength && line[indent] == fence && len(bytes.TrimSpace(line[indent+run:])) == 0 {
				fence = 0
			}
			for i := start; i < end; i++ {
				masked[i] = ' '
			}
		} else if run >= 3 {
			fence, fenceLength = line[indent], run
			for i := start; i < end; i++ {
				masked[i] = ' '
			}
		} else {
			for i := start; i < end; i++ {
				if masked[i] == '\\' {
					i++
					continue
				}
				if masked[i] != '`' {
					continue
				}
				n := 1
				for i+n < end && masked[i+n] == '`' {
					n++
				}
				closeAt := -1
				for j := i + n; j < end; {
					if masked[j] != '`' {
						j++
						continue
					}
					k := j
					for k < end && masked[k] == '`' {
						k++
					}
					if k-j == n {
						closeAt = k
						break
					}
					j = k
				}
				if closeAt >= 0 {
					for j := i; j < closeAt; j++ {
						masked[j] = ' '
					}
					i = closeAt - 1
				} else {
					i += n - 1
				}
			}
		}
		start = end + 1
	}
	var spans []markdownDestination
	var labels []int
	for i := 0; i < len(masked); i++ {
		if masked[i] == '\\' {
			i++
			continue
		}
		if masked[i] == '[' {
			labels = append(labels, i)
			continue
		}
		if masked[i] != ']' || len(labels) == 0 {
			continue
		}
		labelStart := labels[len(labels)-1]
		labels = labels[:len(labels)-1]
		j := i
		if j+1 >= len(masked) {
			continue
		}
		reference := masked[j+1] == ':'
		if !reference && masked[j+1] != '(' {
			continue
		}
		if reference {
			lineStart := bytes.LastIndexByte(masked[:labelStart], '\n') + 1
			if labelStart-lineStart > 3 || len(bytes.TrimSpace(masked[lineStart:labelStart])) != 0 {
				continue
			}
		}
		k := j + 2
		for k < len(masked) && markdownSpace(masked[k]) {
			k++
		}
		begin := k
		angle := k < len(masked) && masked[k] == '<'
		if k < len(masked) && masked[k] == '<' {
			begin = k + 1
			k++
			for k < len(masked) && masked[k] != '>' && masked[k] != '\n' {
				if masked[k] == '\\' {
					k++
				}
				k++
			}
			if k >= len(masked) || masked[k] != '>' {
				continue
			}
		} else {
			parens := 0
			for k < len(masked) {
				if masked[k] == '\\' && k+1 < len(masked) {
					k += 2
					continue
				}
				if markdownSpace(masked[k]) {
					break
				}
				if masked[k] == '(' {
					parens++
				}
				if masked[k] == ')' {
					if parens == 0 {
						break
					}
					parens--
				}
				k++
			}
			if parens != 0 {
				continue
			}
		}
		end := k
		if angle {
			end++
		}
		linkEnd, valid := archiveMarkdownLinkEnd(masked, end, reference)
		if !valid {
			continue
		}
		if k > begin {
			spans = append(spans, markdownDestination{begin, k})
		}
		i = linkEnd
	}
	return spans
}

// A destination counts only when the surrounding link syntax is complete.
func archiveMarkdownLinkEnd(content []byte, end int, reference bool) (int, bool) {
	k := end
	for k < len(content) && markdownSpace(content[k]) {
		if reference && content[k] == '\n' {
			return k, true
		}
		k++
	}
	if reference && k == len(content) {
		return k, true
	}
	if !reference && k < len(content) && content[k] == ')' {
		return k, true
	}
	if k == end || k >= len(content) {
		return 0, false
	}
	quote := content[k]
	if quote != '\'' && quote != '"' && quote != '(' {
		return 0, false
	}
	if quote == '(' {
		quote = ')'
	}
	k++
	for k < len(content) && content[k] != quote {
		if content[k] == '\\' {
			k++
		}
		k++
	}
	if k >= len(content) {
		return 0, false
	}
	k++
	for k < len(content) && markdownSpace(content[k]) {
		if reference && content[k] == '\n' {
			return k, true
		}
		k++
	}
	if reference {
		return k, k == len(content)
	}
	return k, k < len(content) && content[k] == ')'
}
