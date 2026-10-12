package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// StaleHistoryCitations reports absent files held by the History Full Tag.
// held contains file paths from git ls-tree -r, never directory entries.
func StaleHistoryCitations(ctx context.Context, repositoryRoot, historyRoot string, held map[string]bool) ([]HistoryCitation, error) {
	output, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "ls-files", "-z", "--", "*.md").Output()
	if err != nil {
		return nil, fmt.Errorf("list stale history citation sources: %w", err)
	}
	citations := []HistoryCitation{}
	historyRoot = strings.TrimSuffix(filepath.ToSlash(historyRoot), "/")
	for _, source := range strings.Split(string(output), "\x00") {
		if source == "" || archivePathUnder(source, historyRoot) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("scan stale history citations: %w", err)
		}
		content, err := os.ReadFile(filepath.Join(repositoryRoot, source))
		if err != nil {
			return nil, fmt.Errorf("read stale history citation source %q: %w", source, err)
		}
		for i, line := range strings.Split(string(content), "\n") {
			targets := map[string]bool{}
			// Resolve relative destinations separately; raw repository paths in prose,
			// code spans and fences use the existing detector's boundary characters.
			for _, span := range scanArchiveMarkdown([]byte(line)) {
				target, _, relative := archiveLinkTarget(line[span.start:span.end], filepath.Dir(source))
				if relative {
					targets[filepath.ToSlash(target)] = true
				}
			}
			for _, token := range strings.FieldsFunc(line, func(r rune) bool {
				return strings.ContainsRune(" \t\r`\"'()[]<>#,;:", r)
			}) {
				targets[token] = true
				// A prose sentence may end just after the path.
				targets[strings.TrimRight(token, ".")] = true
			}
			found := []string{}
			for target := range targets {
				if !archivePathUnder(target, historyRoot) || !held[target] {
					continue
				}
				_, err := os.Stat(filepath.Join(repositoryRoot, target))
				if errors.Is(err, os.ErrNotExist) {
					found = append(found, target)
				} else if err != nil {
					return nil, fmt.Errorf("inspect stale history citation target %q: %w", target, err)
				}
			}
			sort.Strings(found)
			for _, target := range found {
				citations = append(citations, HistoryCitation{Path: source, Line: i + 1, Target: target})
			}
		}
	}
	return citations, nil
}
