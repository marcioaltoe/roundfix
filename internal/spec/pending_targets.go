package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PendingTaskCreates reads only active Specs under specsRoot. Unparseable Task
// files contribute nothing; directory and file read failures remain errors.
func PendingTaskCreates(specsRoot string) (map[string]bool, error) {
	targets := map[string]bool{}
	dirs, err := os.ReadDir(specsRoot)
	if errors.Is(err, os.ErrNotExist) {
		return targets, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pending Specs: %w", err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		full := filepath.Join(specsRoot, dir.Name())
		content, err := os.ReadFile(filepath.Join(full, "_prd.md"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read pending Spec PRD: %w", err)
		}
		status, err := prdStatus(content)
		if err != nil || status != prdStatusActive {
			continue
		}
		files, err := os.ReadDir(full)
		if err != nil {
			return nil, fmt.Errorf("read pending Tasks: %w", err)
		}
		for _, file := range files {
			if !file.Type().IsRegular() || !strings.HasPrefix(file.Name(), "task_") || !strings.HasSuffix(file.Name(), ".md") {
				continue
			}
			path := filepath.Join(full, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("read pending Task: %w", err)
			}
			task, err := parseTaskDocument(content, path)
			if err != nil || task.Frontmatter.Status == string(StatusCompleted) {
				continue
			}
			for _, ref := range task.Context {
				if ref.Kind == ContextKindCreates {
					targets[ref.Path] = true
				}
			}
		}
	}
	return targets, nil
}

// RemovedHistoryPaths reads HEAD's tree and deletion history through Git only,
// without rename detection and limited to the three retired history families.
// Present working-tree paths are excluded. Outside Git or without HEAD it is empty.
func RemovedHistoryPaths(ctx context.Context, root string) ([]string, error) {
	if _, err := archiveGitRead(ctx, root, "rev-parse", "--is-inside-work-tree"); err != nil {
		if ctx.Err() == nil && strings.Contains(err.Error(), "not a git repository") {
			return nil, nil
		}
		return nil, err
	}
	head, err := gitLocatorCommit(ctx, root, "HEAD")
	if err != nil || head == "" {
		return nil, err
	}
	families := []string{"docs/history/findings", "docs/history/backlog", "docs/history/adr"}
	seen := map[string]bool{}
	for _, args := range [][]string{
		{"ls-tree", "-r", "--name-only", "-z", "HEAD", "--"},
		{"log", "--format=", "--name-only", "-z", "--diff-filter=D", "--no-renames", "HEAD", "--"},
	} {
		out, err := archiveGitRead(ctx, root, append(args, families...)...)
		if err != nil {
			return nil, fmt.Errorf("list removed history: %w", err)
		}
		for _, path := range strings.Split(string(out), "\x00") {
			path = strings.TrimLeft(path, "\n")
			if path == "" {
				continue
			}
			if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path))); errors.Is(err, os.ErrNotExist) {
				seen[path] = true
			} else if err != nil {
				return nil, fmt.Errorf("inspect removed history %s: %w", path, err)
			}
		}
	}
	paths := make([]string, 0, len(seen))
	for path := range seen {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

// ResolveMissingLinkByName never chooses among distinct same-named candidates.
// Duplicate occurrences of the same repository path are one candidate.
func ResolveMissingLinkByName(target string, candidates []string) (string, bool) {
	match := ""
	for _, candidate := range candidates {
		if filepath.Base(candidate) != filepath.Base(target) {
			continue
		}
		if match != "" && match != candidate {
			return "", false
		}
		match = candidate
	}
	return match, match != ""
}
