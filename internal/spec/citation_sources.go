package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func citationSourceExcluded(path string) bool {
	if archivePathUnder(path, ".agents/skills") || archivePathUnder(path, "skills") {
		return true
	}
	for _, directory := range strings.Split(filepath.ToSlash(filepath.Dir(path)), "/") {
		if directory == "testdata" {
			return true
		}
	}
	return false
}

func citationMarkdown(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown", ".mdx":
		return true
	default:
		return false
	}
}

// CitationSources reads only this repository's files. Skill trees, testdata,
// symlinks (including index symlinks on Windows) and binary files are excluded
// before citation scanning. Untracked files are opt-in; ignored files stay out.
func CitationSources(ctx context.Context, root string, withUntracked bool) ([]string, error) {
	args := []string{"ls-files", "-z", "--stage", "--cached"}
	if withUntracked {
		args = append(args, "--others", "--exclude-standard")
	}
	output, err := archiveGitRead(ctx, root, args...)
	if err != nil {
		return nil, fmt.Errorf("list citation sources: %w", err)
	}
	seen := map[string]bool{}
	var paths []string
	for _, entry := range strings.Split(string(output), "\x00") {
		path := entry
		if metadata, name, staged := strings.Cut(entry, "\t"); staged && len(strings.Fields(metadata)) == 3 {
			if !strings.HasPrefix(metadata, "100") {
				continue
			}
			path = name
		}
		if path == "" || seen[path] || citationSourceExcluded(path) {
			continue
		}
		seen[path] = true
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if historyTrackedSymlink(root, path) {
			continue
		}
		if err := historyRewriteRegularFile(root, path); err != nil {
			if errors.Is(err, os.ErrNotExist) && archivePathUnder(path, "docs/history") {
				continue
			}
			return nil, err
		}
		file, err := os.Open(filepath.Join(root, path))
		if err != nil {
			return nil, fmt.Errorf("open citation source %q: %w", path, err)
		}
		var prefix [8000]byte
		n, readErr := io.ReadFull(file, prefix[:])
		closeErr := file.Close()
		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			readErr = nil
		}
		if err := errors.Join(readErr, closeErr); err != nil {
			return nil, fmt.Errorf("inspect citation source %q: %w", path, err)
		}
		if bytes.IndexByte(prefix[:n], 0) >= 0 {
			continue
		}
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

// SourceHistoryCitations reports tracked text citations in source, line, target
// order, preserving the caller's target spelling. It never reads excluded files
// for citations and never sends repository content outside this machine.
func SourceHistoryCitations(ctx context.Context, root string, removed []string) ([]HistoryCitation, error) {
	if len(removed) == 0 {
		return nil, nil
	}
	paths, err := CitationSources(ctx, root, false)
	if err != nil {
		return nil, err
	}
	requested := append([]string(nil), removed...)
	sort.Strings(requested)
	original := map[string]string{}
	var targets []string
	for _, target := range requested {
		canonical := strings.TrimSuffix(target, "/")
		if _, exists := original[canonical]; !exists {
			original[canonical] = target
			targets = append(targets, canonical)
		}
	}
	var citations []HistoryCitation
	for _, path := range paths {
		if historyPathRemoved(path, targets) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		content, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			return nil, fmt.Errorf("read history citation source %q: %w", path, err)
		}
		found := map[HistoryCitation]bool{}
		if citationMarkdown(path) {
			for _, citation := range archiveContentCitations(content, path, targets, "") {
				found[citation] = true
			}
		}
		// Probe with the writer's bounded path pass, discarding its output.
		// Consuming matched tokens also avoids counting their internal dots
		// as the start of a second, unrelated relative path.
		rw := CitationRewrite{}
		for i, line := range strings.Split(string(content), "\n") {
			rw.Locate = func(candidate string) (GitLocator, bool, error) {
				candidate = citationProseTarget(candidate, path)
				matched := false
				for _, target := range targets {
					if candidate == target || strings.HasPrefix(candidate, target+"/") {
						found[HistoryCitation{Path: path, Line: i + 1, Target: target}] = true
						matched = true
					}
				}
				return GitLocator{Path: candidate}, matched, nil
			}
			if _, _, err := rw.ApplyText([]byte(line), path); err != nil {
				return nil, err
			}
		}
		var sorted []HistoryCitation
		for citation := range found {
			sorted = append(sorted, citation)
		}
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Line != sorted[j].Line {
				return sorted[i].Line < sorted[j].Line
			}
			return sorted[i].Target < sorted[j].Target
		})
		for _, citation := range sorted {
			citation.Target = original[citation.Target]
			citations = append(citations, citation)
		}
	}
	return citations, nil
}
