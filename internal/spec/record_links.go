package spec

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RetirementRuleSentinel enables deletion checks only after the guide adopts them.
const RetirementRuleSentinel = "A retirement deletes the file"

// MissingLink identifies a broken relative destination, or a link to its own file.
type MissingLink struct {
	Path                string
	Line                int
	Destination, Target string
	Self                bool
}

// MissingRelativeLinks uses the archive scanner, which excludes code and URLs.
func MissingRelativeLinks(repositoryRoot, path string, content []byte) []MissingLink {
	var links []MissingLink
	for _, span := range scanArchiveMarkdown(content) {
		destination := string(content[span.start:span.end])
		target, _, relative := archiveLinkTarget(destination, filepath.Dir(path))
		target = filepath.ToSlash(target)
		if !relative || !historyRepositoryPath(target) {
			continue
		}
		self := target == path
		_, err := os.Stat(filepath.Join(repositoryRoot, target))
		if self || os.IsNotExist(err) {
			links = append(links, MissingLink{Path: path, Line: bytes.Count(content[:span.start], []byte("\n")) + 1, Destination: destination, Target: target, Self: self})
		}
	}
	return links
}

// RecordLinkPlan freezes destinations before removals, retaining rewritten bytes
// for inspection. Apply re-reads records to preserve other planned citation edits.
type RecordLinkPlan struct {
	Files     map[string][]byte
	Links     int
	Citations []HistoryCitation
	rewrites  map[string]map[string]GitLocator
}

// ResolveRecordLink tries the record-relative target before its original source.
func ResolveRecordLink(ctx context.Context, root string, link MissingLink, record ArchiveRecord) (GitLocator, error) {
	target := link.Target
	if link.Self {
		target = record.Source
	}
	locator, err := ResolveGitLocator(ctx, root, target, record.SourceRevision, record.PullRequest)
	if err != nil || link.Self || locator.Commit != "" {
		return locator, err
	}
	sourceTarget, _, relative := archiveLinkTarget(link.Destination, record.Source)
	sourceTarget = filepath.ToSlash(sourceTarget)
	if relative && historyRepositoryPath(sourceTarget) && sourceTarget != target {
		original, err := ResolveGitLocator(ctx, root, sourceTarget, record.SourceRevision, record.PullRequest)
		if err != nil {
			return locator, err
		}
		if original.Commit != "" {
			return original, nil
		}
	}
	return locator, nil
}

// PlanRecordLinks plans only surviving Archive Records, never legacy folders.
func PlanRecordLinks(ctx context.Context, root, archiveRootRel string, skip map[string]bool) (RecordLinkPlan, error) {
	p := RecordLinkPlan{Files: map[string][]byte{}, rewrites: map[string]map[string]GitLocator{}}
	entries, err := os.ReadDir(filepath.Join(root, archiveRootRel))
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return p, fmt.Errorf("list record links: %w", err)
	}
	for _, entry := range entries {
		path := filepath.ToSlash(filepath.Join(archiveRootRel, entry.Name()))
		if entry.IsDir() || skip[path] || !strings.HasSuffix(path, ".md") {
			continue
		}
		if err := historyRewriteRegularFile(root, path); err != nil {
			return p, err
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return p, fmt.Errorf("read record links %s: %w", path, err)
		}
		if !IsArchiveRecord(entry.Name(), content) {
			continue
		}
		record, err := ParseArchiveRecord(content)
		if err != nil {
			return p, fmt.Errorf("parse record links %s: %w", path, err)
		}
		locators := map[string]GitLocator{}
		for _, link := range MissingRelativeLinks(root, path, content) {
			locator, err := ResolveRecordLink(ctx, root, link, record)
			if err != nil {
				return p, err
			}
			locators[link.Target] = locator
			p.Citations = append(p.Citations, HistoryCitation{Path: path, Line: link.Line, Target: link.Target, Rewrite: locator.Text(false)})
		}
		if len(locators) == 0 {
			continue
		}
		next, count, err := rewriteRecordLinks(content, path, locators)
		if err != nil {
			return p, err
		}
		p.Files[path], p.rewrites[path] = next, locators
		p.Links += count
	}
	return p, nil
}

// Apply rewrites the frozen records and returns actual link and file counts.
func (p RecordLinkPlan) Apply(ctx context.Context, root string) (links, files int, err error) {
	paths := make([]string, 0, len(p.Files))
	for path := range p.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if err := ctx.Err(); err != nil {
			return links, files, err
		}
		if err := historyRewriteRegularFile(root, path); err != nil {
			return links, files, err
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return links, files, fmt.Errorf("read record links %s: %w", path, err)
		}
		next, count, err := rewriteRecordLinks(content, path, p.rewrites[path])
		if err != nil {
			return links, files, err
		}
		if count == 0 {
			continue
		}
		if err := os.WriteFile(filepath.Join(root, path), next, 0o644); err != nil {
			return links, files, fmt.Errorf("rewrite record links %s: %w", path, err)
		}
		links += count
		files++
	}
	return links, files, nil
}

// A repair owns only link spans. The general Citation Rewrite also rewrites
// prose paths, which here could alter an Archive Record's source metadata.
func rewriteRecordLinks(content []byte, path string, locators map[string]GitLocator) ([]byte, int, error) {
	spans := scanArchiveMarkdown(content)
	sort.SliceStable(spans, func(i, j int) bool {
		if spans[i].linkStart != spans[j].linkStart {
			return spans[i].linkStart < spans[j].linkStart
		}
		return spans[i].linkEnd > spans[j].linkEnd
	})
	rw := CitationRewrite{Locate: func(target string) (GitLocator, bool, error) { l, ok := locators[target]; return l, ok, nil }}
	var out bytes.Buffer
	end, count := 0, 0
	for _, span := range spans {
		if span.linkStart < end {
			continue
		}
		target, _, relative := archiveLinkTarget(string(content[span.start:span.end]), filepath.Dir(path))
		if _, ok := locators[filepath.ToSlash(target)]; !relative || !ok {
			continue
		}
		next, n, err := rw.Apply(content[span.linkStart:span.linkEnd], path, path)
		if err != nil {
			return nil, 0, err
		}
		out.Write(content[end:span.linkStart])
		out.Write(next)
		end = span.linkEnd
		count += n
	}
	out.Write(content[end:])
	return out.Bytes(), count, nil
}
