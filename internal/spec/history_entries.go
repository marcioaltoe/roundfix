package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// HistoryKindAction describes how a retired family leaves its full text in Git.
type HistoryKindAction string

// HistoryKindPlan contains pending files and the bytes that leave the tree.
type HistoryKindPlan struct {
	Kind                  ArchiveKind
	Action                HistoryKindAction
	Files                 []string
	BytesBefore           int64
	Lineage               []HistoryLineage
	Refused               []HistoryRefusal
	specRoot, archiveRoot string
}

// HistoryLineage is the exact line appended to an accepted successor.
type HistoryLineage struct{ Successor, Removed, Line string }

// HistoryRefusal leaves one superseded ADR for a maintainer to resolve.
type HistoryRefusal struct{ Path, Reason string }

const historyADRRefusal = "superseded ADR has no accepted successor; name the accepted ADR that replaces it in superseded_by, or set its status to deprecated"

var historyEntryProvenance = regexp.MustCompile("^Full text in Git at `([[:xdigit:]]{40}|[[:xdigit:]]{64})`: `([^`\\r\\n]+)`\\.$")
var historyEntryLink = regexp.MustCompile(`\]\(([^\s)]+)`)

// IsReducedHistoryEntry recognizes the provenance line at the end of an entry.
func IsReducedHistoryEntry(content []byte) bool {
	lines := strings.Split(string(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			return historyEntryProvenance.MatchString(strings.TrimSuffix(lines[i], "\r"))
		}
	}
	return false
}

// PlanHistoryKinds plans every retired family in a stable order.
func PlanHistoryKinds(ctx context.Context, repositoryRoot string) ([]HistoryKindPlan, error) {
	var plans []HistoryKindPlan
	for _, kind := range []ArchiveKind{ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindADR, ArchiveKindReview, ArchiveKindHandoff} {
		plan, err := PlanHistoryKind(ctx, repositoryRoot, kind)
		if err != nil {
			return nil, err
		}
		if len(plan.Files) > 0 || len(plan.Refused) > 0 {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

func historyKindRoots(kind ArchiveKind, specRoot string) []string {
	roots := []string{ArchiveDir(kind)}
	switch kind {
	case ArchiveKindFinding:
		roots = append(roots, "docs/findings/_archived", "_archived/findings")
	case ArchiveKindBacklog:
		roots = append(roots, "_archived/backlog")
	case ArchiveKindADR:
		roots = append(roots, "_archived/adr")
	case ArchiveKindReview:
		roots = append(roots, specRoot+"/_reviews", specRoot+"/reviews")
	}
	return roots
}

// PlanHistoryKind measures a family independently, including legacy locations.
// The optional Spec Root is repository-relative and defaults to docs/specs.
func PlanHistoryKind(ctx context.Context, repositoryRoot string, kind ArchiveKind, specRoots ...string) (HistoryKindPlan, error) {
	plan := HistoryKindPlan{Kind: kind, Action: "remove", specRoot: "docs/specs"}
	if len(specRoots) > 0 {
		plan.specRoot = specRoots[0]
	}
	if !historyRepositoryPath(plan.specRoot) {
		return plan, fmt.Errorf("invalid Spec Root %q", plan.specRoot)
	}
	switch kind {
	case ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindADR, ArchiveKindReview, ArchiveKindHandoff:
	default:
		return plan, fmt.Errorf("plan history kind %q: unsupported kind", kind)
	}
	seen := map[string]bool{}
	add := func(path string) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(repositoryRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if seen[rel] {
			return nil
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("Spec contains symlink or non-regular history entry %q", rel)
		}
		if kind == ArchiveKindADR {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if ClassifyADR(content).Reason == "superseded" {
				successor, next, err := historyADRSuccessor(repositoryRoot, content)
				if err != nil {
					return err
				}
				if successor == "" {
					seen[rel] = true
					plan.Refused = append(plan.Refused, HistoryRefusal{Path: rel, Reason: historyADRRefusal})
					return nil
				}
				id := historyADRIdentifier(filepath.Base(path))
				if !historyADRNamesPredecessor(next, filepath.Base(path)) {
					date, err := historyADRSupersessionDate(ctx, repositoryRoot, rel, content)
					if err != nil {
						return err
					}
					locator, err := ResolveGitLocator(ctx, repositoryRoot, rel, "", "")
					if err != nil {
						return err
					}
					plan.Lineage = append(plan.Lineage, HistoryLineage{Successor: successor, Removed: rel, Line: fmt.Sprintf("Supersedes %s (%s); Git keeps its text at %s.", id, date, locator.Text(false))})
				}
			}
		}
		seen[rel] = true
		plan.Files = append(plan.Files, rel)
		plan.BytesBefore += info.Size()
		return nil
	}
	walk := func(rel string) error {
		root := filepath.Join(repositoryRoot, rel)
		return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if errors.Is(err, os.ErrNotExist) && path == root {
				return nil
			}
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return ctx.Err()
			}
			return add(path)
		})
	}
	roots := historyKindRoots(kind, plan.specRoot)
	for _, root := range roots {
		if kind == ArchiveKindReview && root != ArchiveDir(kind) {
			entries, err := os.ReadDir(filepath.Join(repositoryRoot, root))
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			if err != nil {
				return plan, err
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				dir := filepath.Join(root, entry.Name())
				state, _, err := ClassifyReview(ctx, repositoryRoot, filepath.Join(repositoryRoot, dir))
				if err != nil {
					return plan, err
				}
				if state == ReviewFinished {
					if err := walk(dir); err != nil {
						return plan, err
					}
				}
			}
		} else if err := walk(root); err != nil {
			return plan, fmt.Errorf("plan history kind %s: %w", kind, err)
		}
	}
	active := ""
	switch kind {
	case ArchiveKindFinding:
		active = "docs/findings"
	case ArchiveKindBacklog:
		active = "docs/backlog"
	case ArchiveKindADR:
		active = "docs/adr"
	}
	if active != "" {
		entries, err := os.ReadDir(filepath.Join(repositoryRoot, active))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return plan, err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			path := filepath.Join(repositoryRoot, active, entry.Name())
			if !entry.Type().IsRegular() {
				return plan, fmt.Errorf("active history entry %s is not a regular file", path)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return plan, err
			}
			status, _ := retirementStatus(content)
			retired := false
			switch kind {
			case ArchiveKindFinding:
				retired = status == "deferred" || status == "done" || status == "deprecated" || status == "superseded" || status == "closed" || status == "cancelled"
			case ArchiveKindBacklog:
				retired = status == "declined" || status == "deferred" || status == "done" || status == "deprecated" || status == "superseded" || status == "closed" || status == "cancelled"
			case ArchiveKindADR:
				retired = ClassifyADR(content).Retired
			}
			if retired {
				if err := add(path); err != nil {
					return plan, err
				}
			}
		}
	}
	sort.Strings(plan.Files)
	sort.Slice(plan.Lineage, func(i, j int) bool { return plan.Lineage[i].Removed < plan.Lineage[j].Removed })
	sort.Slice(plan.Refused, func(i, j int) bool { return plan.Refused[i].Path < plan.Refused[j].Path })
	return plan, nil
}

func historyADRSuccessor(root string, content []byte) (string, []byte, error) {
	front, _, err := splitFrontmatter(content)
	if err != nil {
		return "", nil, nil
	}
	var metadata struct {
		Successor yaml.Node `yaml:"superseded_by"`
	}
	if err := yaml.Unmarshal(front, &metadata); err != nil {
		return "", nil, fmt.Errorf("read superseded ADR: %w", err)
	}
	if metadata.Successor.Kind != yaml.ScalarNode {
		return "", nil, nil
	}
	entries, err := os.ReadDir(filepath.Join(root, "docs/adr"))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, fmt.Errorf("read successor directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !historyADRReferenceNames(metadata.Successor.Value, entry.Name()) {
			continue
		}
		path := "docs/adr/" + entry.Name()
		if err := historyRewriteRegularFile(root, path); err != nil {
			return "", nil, err
		}
		next, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return "", nil, fmt.Errorf("read successor %s: %w", path, err)
		}
		if status, _ := retirementStatus(next); status == "accepted" {
			return path, next, nil
		}
	}
	return "", nil, nil
}

// historyADRIdentifier reserves numeric identifiers for numbered, not date-named, ADRs.
func historyADRIdentifier(name string) string {
	stem := strings.TrimSuffix(filepath.Base(name), ".md")
	number, slug, ok := strings.Cut(stem, "-")
	if ok && number != "" && slug != "" && strings.Trim(number, "0123456789") == "" && !historyADRDateName.MatchString(stem) {
		return "ADR-" + number
	}
	return stem
}

var historyADRDateName = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-`)

func historyADRSupersessionDate(ctx context.Context, root, rel string, content []byte) (string, error) {
	front, _, err := splitFrontmatter(content)
	if err == nil {
		var metadata struct {
			Updated yaml.Node `yaml:"updated_at"`
		}
		if yaml.Unmarshal(front, &metadata) == nil && metadata.Updated.Kind == yaml.ScalarNode && metadata.Updated.Tag != "!!null" {
			var updated time.Time
			if err := metadata.Updated.Decode(&updated); err == nil {
				return updated.UTC().Format(time.DateOnly), nil
			}
			if updated, err := time.Parse(time.DateOnly, metadata.Updated.Value); err == nil {
				return updated.Format(time.DateOnly), nil
			}
		}
	}
	date, err := archiveGitRead(ctx, root, "log", "-1", "--format=%cI", "HEAD", "--", rel)
	if err != nil {
		return "", fmt.Errorf("read supersession date for %s: %w", rel, err)
	}
	if strings.TrimSpace(string(date)) == "" {
		date, err = archiveGitRead(ctx, root, "show", "-s", "--format=%cI", "HEAD")
		if err != nil {
			return "", fmt.Errorf("read lineage date: %w", err)
		}
	}
	committed, err := time.Parse(time.RFC3339, strings.TrimSpace(string(date)))
	if err != nil {
		return "", fmt.Errorf("parse lineage date: %w", err)
	}
	return committed.UTC().Format(time.DateOnly), nil
}

// historyADRReferenceNames matches numeric references only for numbered ADRs,
// and exact stems, filenames and paths for every ADR filename.
func historyADRReferenceNames(reference, name string) bool {
	if !strings.HasSuffix(name, ".md") {
		return false
	}
	stem := strings.TrimSuffix(name, ".md")
	if reference == stem || reference == name || strings.HasSuffix(reference, "/"+name) {
		return true
	}
	id := historyADRIdentifier(name)
	return id != stem && (reference == id || reference == strings.TrimPrefix(id, "ADR-"))
}

func historyADRNamesPredecessor(content []byte, name string) bool {
	id := historyADRIdentifier(name)
	if historyEntryNamesPath(string(content), id) || historyEntryNamesPath(string(content), name) {
		return true
	}
	front, _, err := splitFrontmatter(content)
	if err != nil {
		return false
	}
	var metadata struct {
		Supersedes yaml.Node `yaml:"supersedes"`
	}
	if err := yaml.Unmarshal(front, &metadata); err != nil {
		return false
	}
	references := []*yaml.Node{&metadata.Supersedes}
	if metadata.Supersedes.Kind == yaml.SequenceNode {
		references = metadata.Supersedes.Content
	}
	for _, reference := range references {
		if reference.Kind == yaml.ScalarNode && historyADRReferenceNames(reference.Value, name) {
			return true
		}
	}
	return false
}

func appendHistoryLineage(content []byte, line string) []byte {
	for _, existing := range strings.Split(string(content), "\n") {
		if existing == line {
			return content
		}
	}
	next := append([]byte(nil), content...)
	if len(next) > 0 && next[len(next)-1] != '\n' {
		next = append(next, '\n')
	}
	next = append(next, '\n')
	next = append(next, line...)
	return append(next, '\n')
}

// LineageCitations lists the exact lines apply will append, in removal order.
func (p HistoryKindPlan) LineageCitations(root string) ([]HistoryCitation, error) {
	contents := map[string][]byte{}
	var citations []HistoryCitation
	for _, lineage := range p.Lineage {
		content, ok := contents[lineage.Successor]
		if !ok {
			var err error
			content, err = os.ReadFile(filepath.Join(root, lineage.Successor))
			if err != nil {
				return nil, fmt.Errorf("read successor %s: %w", lineage.Successor, err)
			}
		}
		next := appendHistoryLineage(content, lineage.Line)
		if len(next) != len(content) {
			line := strings.Count(string(next[:len(next)-len(lineage.Line)-1]), "\n") + 1
			citations = append(citations, HistoryCitation{Path: lineage.Successor, Line: line, Target: lineage.Removed, Rewrite: lineage.Line})
		}
		contents[lineage.Successor] = next
	}
	return citations, nil
}

// ApplyHistoryKind validates every path before removing the first file.
func ApplyHistoryKind(repositoryRoot string, plan HistoryKindPlan) error {
	repositoryRoot = filepath.Clean(repositoryRoot)
	if plan.Action != "remove" {
		return fmt.Errorf("apply history kind %s: invalid action %q", plan.Kind, plan.Action)
	}
	switch plan.Kind {
	case ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindADR, ArchiveKindReview, ArchiveKindHandoff, ArchiveKindSpec:
	default:
		return fmt.Errorf("apply history kind %q: unsupported kind", plan.Kind)
	}
	specRoot := plan.specRoot
	if specRoot == "" {
		specRoot = "docs/specs"
	}
	if !historyRepositoryPath(specRoot) || (plan.archiveRoot != "" && !historyRepositoryPath(plan.archiveRoot)) {
		return fmt.Errorf("invalid history plan root")
	}
	roots := historyKindRoots(plan.Kind, specRoot)
	if plan.Kind == ArchiveKindSpec && plan.archiveRoot != "" {
		roots = []string{plan.archiveRoot}
	}
	for _, path := range plan.Files {
		allowed := false
		for _, root := range roots {
			allowed = allowed || archivePathUnder(path, root)
		}
		for _, active := range []struct {
			kind ArchiveKind
			root string
		}{{ArchiveKindFinding, "docs/findings"}, {ArchiveKindBacklog, "docs/backlog"}, {ArchiveKindADR, "docs/adr"}} {
			allowed = allowed || (plan.Kind == active.kind && filepath.ToSlash(filepath.Dir(path)) == active.root)
		}
		if !allowed || filepath.IsAbs(path) || filepath.ToSlash(filepath.Clean(path)) != path || strings.Contains(path, "\\") {
			return fmt.Errorf("apply history kind %s: path outside kind: %q", plan.Kind, path)
		}
		if plan.Kind == ArchiveKindSpec && (filepath.ToSlash(filepath.Dir(path)) != roots[0] || !strings.HasSuffix(path, ".md")) {
			return fmt.Errorf("apply history records: invalid path %q", path)
		}
		for current := filepath.Join(repositoryRoot, path); current != repositoryRoot; current = filepath.Dir(current) {
			info, err := os.Lstat(current)
			if err != nil {
				return fmt.Errorf("inspect history entry %q: %w", current, err)
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("inspect history entry %q: symbolic link", current)
			}
			if current == filepath.Join(repositoryRoot, path) && !info.Mode().IsRegular() {
				return fmt.Errorf("inspect history entry %q: not a regular file", current)
			}
		}
	}
	for _, lineage := range plan.Lineage {
		if plan.Kind != ArchiveKindADR || filepath.ToSlash(filepath.Dir(lineage.Successor)) != "docs/adr" {
			return fmt.Errorf("invalid lineage successor %q", lineage.Successor)
		}
		if err := historyRewriteRegularFile(repositoryRoot, lineage.Successor); err != nil {
			return err
		}
	}
	protected := map[string]bool{"docs/adr": true, "docs/backlog": true, "docs/findings": true, specRoot: true, "docs/history": true, "docs": true, ".": true}
	for _, path := range plan.Files {
		if err := os.Remove(filepath.Join(repositoryRoot, path)); err != nil {
			return fmt.Errorf("remove history entry %q: %w", path, err)
		}
		for dir := filepath.ToSlash(filepath.Dir(path)); !protected[dir]; dir = filepath.ToSlash(filepath.Dir(dir)) {
			entries, err := os.ReadDir(filepath.Join(repositoryRoot, dir))
			if err != nil {
				return fmt.Errorf("read history directory %q: %w", dir, err)
			}
			if len(entries) != 0 {
				break
			}
			if err := os.Remove(filepath.Join(repositoryRoot, dir)); err != nil {
				return fmt.Errorf("prune history directory %q: %w", dir, err)
			}
		}
	}
	for _, lineage := range plan.Lineage {
		path := filepath.Join(repositoryRoot, lineage.Successor)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read lineage successor %s: %w", lineage.Successor, err)
		}
		next := appendHistoryLineage(content, lineage.Line)
		if len(next) != len(content) {
			if err := os.WriteFile(path, next, 0o644); err != nil {
				return fmt.Errorf("write lineage successor %s: %w", lineage.Successor, err)
			}
		}
	}

	return nil
}

// HistoryCitation locates a tracked source line naming a removed path.
type HistoryCitation struct {
	Path    string
	Line    int
	Target  string
	Rewrite string
}

// HistoryCitations reads tracked Markdown, including Archive Records, and skips
// sources removed by the same plan.
func HistoryCitations(ctx context.Context, repositoryRoot string, removed []string) ([]HistoryCitation, error) {
	output, err := archiveGitRead(ctx, repositoryRoot, "ls-files", "-z", "--", "*.md")
	if err != nil {
		return nil, fmt.Errorf("list history citation sources: %w", err)
	}
	if len(removed) == 0 {
		return nil, nil
	}
	// Match canonical paths, but preserve the caller's target spelling in the
	// report. Consumers use a trailing slash to identify a removed folder.
	requested := append([]string(nil), removed...)
	sort.Strings(requested)
	originalTargets := make(map[string]string, len(requested))
	var targets []string
	for _, target := range requested {
		canonical := strings.TrimSuffix(target, "/")
		if _, exists := originalTargets[canonical]; !exists {
			originalTargets[canonical] = target
			targets = append(targets, canonical)
		}
	}
	var citations []HistoryCitation
	for _, path := range strings.Split(string(output), "\x00") {
		if path == "" || historyPathRemoved(path, targets) {
			continue
		}
		if historyTrackedSymlink(repositoryRoot, path) {
			continue
		}
		if err := historyRewriteRegularFile(repositoryRoot, path); err != nil {
			// Skip only a missing entry; a path behind a file and other errors must fail.
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}
		content, err := os.ReadFile(filepath.Join(repositoryRoot, path))
		// Skip only a missing entry; a path behind a file and other errors must fail.
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read history citation source %q: %w", path, err)
		}
		for _, citation := range archiveContentCitations(content, path, targets, "") {
			citation.Target = originalTargets[citation.Target]
			citations = append(citations, citation)
		}
	}
	return citations, nil
}

func historyPathRemoved(path string, removed []string) bool {
	for _, target := range removed {
		if path == target || archivePathUnder(path, strings.TrimSuffix(target, "/")) {
			return true
		}
	}
	return false
}

func historyEntryCitesPath(line, source, target string) bool {
	if historyEntryNamesPath(line, target) {
		return true
	}
	for _, link := range historyEntryLink.FindAllStringSubmatch(line, -1) {
		path, _, _ := strings.Cut(link[1], "#")
		resolved := filepath.ToSlash(filepath.Clean(filepath.Join(filepath.Dir(source), filepath.FromSlash(path))))
		if resolved == target || strings.HasPrefix(resolved, target+"/") {
			return true
		}
	}
	return false
}

func historyEntryNamesPath(line, target string) bool {
	if target == "" {
		return false
	}
	for start := 0; start < len(line); {
		i := strings.Index(line[start:], target)
		if i < 0 {
			return false
		}
		i += start
		end := i + len(target)
		boundary := func(c byte) bool { return citationPathBoundary(c) }
		suffix := strings.TrimLeft(line[end:], "`")
		if (end < len(line) && line[end] == '@') || strings.HasPrefix(suffix, " (not reachable from any branch or tag)") || strings.HasPrefix(suffix, " (pull request #") {
			start = end
			continue
		}
		if (i == 0 || boundary(line[i-1])) && (end == len(line) || line[end] == '/' || boundary(line[end])) {
			return true
		}
		start = end
	}
	return false
}

// PlanSupersededRecords selects only records no surviving tracked Markdown cites.
// A grant always keeps its record; keeping a record also keeps records it cites.
func PlanSupersededRecords(ctx context.Context, root, archiveRootRel string, removed []string) (HistoryKindPlan, error) {
	plan := HistoryKindPlan{Kind: ArchiveKindSpec, Action: "remove", archiveRoot: archiveRootRel, specRoot: "docs/specs"}
	if archiveRootRel != ArchiveDir(ArchiveKindSpec) {
		plan.specRoot = filepath.ToSlash(filepath.Dir(archiveRootRel))
	}
	entries, err := os.ReadDir(filepath.Join(root, archiveRootRel))
	if errors.Is(err, os.ErrNotExist) {
		return plan, nil
	}
	if err != nil {
		return plan, err
	}
	candidates := map[string]ArchiveRecord{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := archiveRootRel + "/" + entry.Name()
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return plan, err
		}
		if !IsArchiveRecord(entry.Name(), content) {
			continue
		}
		record, err := ParseArchiveRecord(content)
		if err != nil {
			return plan, err
		}
		if record.Disposition == ArchiveSuperseded && len(record.Regeneration) == 0 {
			candidates[path] = record
		}
	}
	if len(candidates) == 0 {
		return plan, nil
	}
	output, err := archiveGitRead(ctx, root, "ls-files", "-z", "--", "*.md")
	if err != nil {
		return plan, err
	}
	sources := map[string][]byte{}
	for _, path := range strings.Split(string(output), "\x00") {
		if path == "" || historyPathRemoved(path, removed) || strings.HasPrefix(path, ".agents/skills/") || strings.HasPrefix(path, "skills/") || strings.Contains("/"+path, "/testdata/") {
			continue
		}
		info, err := os.Lstat(filepath.Join(root, path))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return plan, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if err := historyRewriteRegularFile(root, path); err != nil {
			return plan, err
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return plan, err
		}
		sources[path] = content
	}
	for changed := true; changed; {
		changed = false
		for path, content := range sources {
			if _, removing := candidates[path]; removing {
				continue
			}
			for target, record := range candidates {
				if len(archiveContentCitations(content, path, []string{target}, "")) > 0 || historyEntryNamesPath(string(content), record.Spec) {
					delete(candidates, target)
					changed = true
				}
			}
		}
	}
	for path := range candidates {
		info, err := os.Lstat(filepath.Join(root, path))
		if err != nil {
			return plan, err
		}
		if !info.Mode().IsRegular() {
			return plan, fmt.Errorf("superseded record %s is not a regular file", path)
		}
		plan.Files = append(plan.Files, path)
		plan.BytesBefore += info.Size()
	}
	sort.Strings(plan.Files)
	return plan, nil
}

func historyTrackedSymlink(root, path string) bool {
	info, err := os.Lstat(filepath.Join(root, path))
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

// HistoryRewritePlan binds the citing files and destinations before any removal.
// Apply re-reads only these files, preserving edits made by the Baseline stage.
type HistoryRewritePlan struct {
	Citations []HistoryCitation
	Files     map[string][]byte
	Count     int
	rewrite   CitationRewrite
}

// PlanHistoryRewrite plans the shared Citation Rewrite for removed files and
// folder moves. It only reads the repository and never rewrites a removed source.
func PlanHistoryRewrite(ctx context.Context, root, specRootRel, archiveRootRel string, removed []string, moves map[string]string) (HistoryRewritePlan, error) {
	p := HistoryRewritePlan{Files: map[string][]byte{}}
	locators := map[string]GitLocator{}
	for source, destination := range moves {
		if !historyRepositoryPath(source) || !historyRepositoryPath(destination) {
			return p, fmt.Errorf("invalid citation move %q to %q", source, destination)
		}
	}
	for _, path := range removed {
		if !historyRepositoryPath(path) {
			return p, fmt.Errorf("invalid removal path %q", path)
		}
		if _, moved := moves[path]; moved {
			continue
		}
		locator, err := ResolveGitLocator(ctx, root, path, "history-full", "")
		if err != nil {
			return p, err
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return p, err
		}
		if IsReducedHistoryEntry(content) {
			lines := strings.Split(strings.TrimSpace(string(content)), "\n")
			match := historyEntryProvenance.FindStringSubmatch(lines[len(lines)-1])
			original, err := ResolveGitLocator(ctx, root, match[2], match[1], "")
			if err != nil {
				return p, err
			}
			// Only the named provenance commit can select its original active path.
			if original.Commit == match[1] {
				locator = original
			}
		}
		locators[path] = locator
	}
	p.rewrite = CitationRewrite{Moves: moves, Relative: true, Locate: func(path string) (GitLocator, bool, error) { l, ok := locators[path]; return l, ok, nil }}
	targets := append([]string(nil), removed...)
	folders := map[string]bool{}
	for path := range moves {
		if !historyPathRemoved(path, targets) {
			targets = append(targets, path)
			folders[path] = true
		}
	}
	citations, err := SourceHistoryCitations(ctx, root, targets)
	if err != nil {
		return p, err
	}
	contents := map[string][]byte{}
	for _, c := range citations {
		broader := false
		for _, other := range citations {
			if c.Path == other.Path && c.Line == other.Line && c.Target != other.Target && strings.HasPrefix(other.Target, strings.TrimSuffix(c.Target, "/")+"/") {
				broader = true
				break
			}
		}
		if broader {
			continue
		}

		if historyPathRemoved(c.Path, removed) {
			continue
		}
		content, ok := contents[c.Path]
		if !ok {
			content, err = os.ReadFile(filepath.Join(root, c.Path))
			if err != nil {
				return p, err
			}
			contents[c.Path] = content
		}
		if err := historyRewriteRegularFile(root, c.Path); err != nil {
			return p, err
		}
		if destination, ok := moves[c.Target]; ok {
			c.Rewrite = destination
		} else if locator, ok := locators[c.Target]; ok {
			c.Rewrite = locator.Text(!citationMarkdown(c.Path))
		}
		if _, planned := p.Files[c.Path]; !planned {
			next, count, err := p.rewriteSource(content, c.Path)
			if err != nil {
				return p, err
			}
			if count > 0 {
				p.Files[c.Path] = next
				p.Count += count
			}
		}
		if folders[c.Target] {
			c.Target += "/"
		}
		p.Citations = append(p.Citations, c)
	}
	return p, nil
}

// Apply writes only the files listed by the plan and counts actual replacements.
func (p HistoryRewritePlan) Apply(ctx context.Context, root string) (citations, files int, err error) {
	paths := make([]string, 0, len(p.Files))
	for path := range p.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return citations, files, err
		}
		if err := historyRewriteRegularFile(root, path); err != nil {
			return citations, files, err
		}
		content, readErr := os.ReadFile(filepath.Join(root, path))
		if readErr != nil {
			return citations, files, fmt.Errorf("read planned citation %s: %w", path, readErr)
		}
		if citationSourceExcluded(path) || bytes.IndexByte(content[:min(len(content), 8000)], 0) >= 0 {
			return citations, files, fmt.Errorf("planned citation %s is no longer a Citation Source", path)
		}
		next, count, rewriteErr := p.rewriteSource(content, path)
		if rewriteErr != nil {
			return citations, files, rewriteErr
		}
		if count == 0 {
			continue
		}
		if writeErr := os.WriteFile(filepath.Join(root, path), next, 0o644); writeErr != nil {
			return citations, files, fmt.Errorf("rewrite citation %s: %w", path, writeErr)
		}
		citations += count
		files++
	}
	return citations, files, nil
}

func (p HistoryRewritePlan) rewriteSource(content []byte, path string) ([]byte, int, error) {
	if citationMarkdown(path) {
		return p.rewrite.Apply(content, path, path)
	}
	return p.rewrite.ApplyText(content, path)
}

// Refuse a symlink before a citation rewrite can read or write through it.
func historyRewriteRegularFile(root, path string) error {
	if !historyRepositoryPath(path) {
		return fmt.Errorf("invalid citing path %q", path)
	}
	root = filepath.Clean(root)
	for current := filepath.Join(root, path); current != root; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("inspect citing file %s: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || (current == filepath.Join(root, path) && !info.Mode().IsRegular()) {
			return fmt.Errorf("citing file %s is not a regular repository file", path)
		}
	}
	return nil
}

func historyRepositoryPath(path string) bool {
	return path != "" && path != "." && !filepath.IsAbs(path) && path != ".." && !strings.HasPrefix(path, "../") && !strings.Contains(path, "\\") && filepath.ToSlash(filepath.Clean(path)) == path
}
