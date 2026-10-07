package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// HistoryKindAction describes how a retired family leaves its full text in Git.
type HistoryKindAction string

// HistoryKindPlan contains only pending files and their measured replacement.
type HistoryKindPlan struct {
	Kind                    ArchiveKind
	Action                  HistoryKindAction
	Files                   []string
	BytesBefore, BytesAfter int64
	reduced                 map[string][]byte
}

var historyEntryProvenance = regexp.MustCompile("^Full text in Git at `([[:xdigit:]]{40}|[[:xdigit:]]{64})`: `[^`\\r\\n]+`\\.$")
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

// ReduceHistoryEntry preserves lifecycle metadata and a short, recoverable record.
func ReduceHistoryEntry(content []byte, revision, path string) ([]byte, error) {
	if IsReducedHistoryEntry(content) {
		return nil, fmt.Errorf("reduce history entry %q: already reduced", path)
	}
	lines := bytes.SplitAfter(content, []byte("\n"))
	frontEnd := 0
	if len(lines) > 0 && strings.TrimRight(string(lines[0]), "\r\n") == "---" {
		frontEnd = len(lines[0])
		closed := false
		for _, line := range lines[1:] {
			frontEnd += len(line)
			if strings.TrimRight(string(line), "\r\n") == "---" {
				closed = true
				break
			}
		}
		if !closed {
			frontEnd = 0
		}
	}
	if frontEnd == 0 {
		return nil, fmt.Errorf("reduce history entry %q: missing front matter", path)
	}
	body := strings.Split(string(content[frontEnd:]), "\n")
	titleIndex := -1
	for i, line := range body {
		if strings.HasPrefix(line, "# ") {
			titleIndex = i
			break
		}
	}
	if titleIndex < 0 {
		return nil, fmt.Errorf("reduce history entry %q: missing title", path)
	}
	var paragraph []string
	for _, line := range body[titleIndex+1:] {
		line = strings.TrimSpace(line)
		if line == "" || historyEntryHeading(line) {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		paragraph = append(paragraph, line)
	}
	provenance := fmt.Sprintf("Full text in Git at `%s`: `%s`.", revision, path)
	if !historyEntryProvenance.MatchString(provenance) {
		return nil, fmt.Errorf("reduce history entry %q: invalid Git provenance", path)
	}
	result := append([]byte(nil), content[:frontEnd]...)
	if len(result) > 0 && result[len(result)-1] != '\n' {
		result = append(result, '\n')
	}
	result = append(result, []byte("\n"+strings.TrimSuffix(body[titleIndex], "\r")+"\n\n")...)
	if len(paragraph) > 0 {
		result = append(result, []byte(strings.Join(strings.Fields(strings.Join(paragraph, " ")), " ")+"\n\n")...)
	}
	return append(result, []byte(provenance+"\n")...), nil
}

func historyEntryHeading(line string) bool {
	text := strings.TrimLeft(line, "#")
	count := len(line) - len(text)
	return count > 0 && count <= 6 && (text == "" || strings.HasPrefix(text, " ") || strings.HasPrefix(text, "\t"))
}

// PlanHistoryKinds measures pending families without reading retired ADRs.
func PlanHistoryKinds(repositoryRoot, revision string) ([]HistoryKindPlan, error) {
	var plans []HistoryKindPlan
	for _, kind := range []ArchiveKind{ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindReview, ArchiveKindHandoff} {
		plan := HistoryKindPlan{Kind: kind, Action: "remove", reduced: make(map[string][]byte)}
		if kind == ArchiveKindFinding || kind == ArchiveKindBacklog {
			plan.Action = "reduce"
		}
		root := filepath.Join(repositoryRoot, ArchiveDir(kind))
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if errors.Is(walkErr, os.ErrNotExist) && path == root {
				return nil
			}
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			if !entry.Type().IsRegular() {
				return fmt.Errorf("history entry %q is not a regular file", path)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if plan.Action == "reduce" && IsReducedHistoryEntry(content) {
				return nil
			}
			rel, err := filepath.Rel(repositoryRoot, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if plan.Action == "reduce" {
				reduced, err := ReduceHistoryEntry(content, revision, rel)
				if err != nil {
					return err
				}
				plan.reduced[rel] = reduced
				plan.BytesAfter += int64(len(reduced))
			}
			plan.Files = append(plan.Files, rel)
			plan.BytesBefore += int64(len(content))
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("plan history kind %s: %w", kind, err)
		}
		if len(plan.Files) > 0 {
			sort.Strings(plan.Files)
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// ApplyHistoryKind touches only a plan's files and directories they leave empty.
func ApplyHistoryKind(repositoryRoot string, plan HistoryKindPlan) error {
	action := HistoryKindAction("remove")
	switch plan.Kind {
	case ArchiveKindFinding, ArchiveKindBacklog:
		action = "reduce"
	case ArchiveKindReview, ArchiveKindHandoff:
	default:
		return fmt.Errorf("apply history kind %q: unsupported kind", plan.Kind)
	}
	if plan.Action != action {
		return fmt.Errorf("apply history kind %s: invalid action %q", plan.Kind, plan.Action)
	}
	root := filepath.Join(repositoryRoot, ArchiveDir(plan.Kind))
	// Validate the entire file list before the first mutation.
	for _, path := range plan.Files {
		if filepath.ToSlash(filepath.Clean(path)) != path || !strings.HasPrefix(path, ArchiveDir(plan.Kind)+"/") {
			return fmt.Errorf("apply history kind %s: path outside kind: %q", plan.Kind, path)
		}
		for current := filepath.Join(repositoryRoot, path); ; current = filepath.Dir(current) {
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
			if current == root {
				break
			}
		}
		if action == "reduce" && plan.reduced[path] == nil {
			return fmt.Errorf("apply history entry %q: missing planned bytes", path)
		}
	}
	for _, path := range plan.Files {
		absolute := filepath.Join(repositoryRoot, path)
		var err error
		if action == "reduce" {
			err = os.WriteFile(absolute, plan.reduced[path], 0o644)
		} else {
			err = os.Remove(absolute)
		}
		if err != nil {
			return fmt.Errorf("apply history entry %q: %w", path, err)
		}
		if action == "remove" {
			for dir := filepath.Dir(absolute); ; dir = filepath.Dir(dir) {
				entries, err := os.ReadDir(dir)
				if errors.Is(err, os.ErrNotExist) {
					break
				}
				if err != nil {
					return fmt.Errorf("read history directory %q: %w", dir, err)
				}
				if len(entries) != 0 {
					break
				}
				if err := os.Remove(dir); err != nil {
					return fmt.Errorf("remove empty history directory %q: %w", dir, err)
				}
				if dir == root {
					break
				}
			}
		}
	}
	return nil
}

// HistoryCitation locates a tracked Markdown line naming a removed path.
type HistoryCitation struct {
	Path   string
	Line   int
	Target string
}

// HistoryCitations reports citations as advice; a citation never refuses a plan.
func HistoryCitations(ctx context.Context, repositoryRoot string, removed []string) ([]HistoryCitation, error) {
	output, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "ls-files", "-z", "--", "*.md").Output()
	if err != nil {
		return nil, fmt.Errorf("list history citation sources: %w", err)
	}
	var citations []HistoryCitation
	targets := append([]string(nil), removed...)
	sort.Strings(targets)
	for _, path := range strings.Split(string(output), "\x00") {
		if path == "" || strings.HasPrefix(path, "docs/history/") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(repositoryRoot, path))
		if err != nil {
			return nil, fmt.Errorf("read history citation source %q: %w", path, err)
		}
		for i, line := range strings.Split(string(content), "\n") {
			for j, target := range targets {
				if target == "" || (j > 0 && targets[j-1] == target) {
					continue
				}
				if historyEntryCitesPath(line, path, strings.TrimSuffix(target, "/")) {
					citations = append(citations, HistoryCitation{Path: path, Line: i + 1, Target: target})
				}
			}
		}
	}
	return citations, nil
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
	for start := 0; start < len(line); {
		i := strings.Index(line[start:], target)
		if i < 0 {
			return false
		}
		i += start
		end := i + len(target)
		boundary := func(c byte) bool { return strings.ContainsRune(" \t\r`\"'()[]<>#.,;:", rune(c)) }
		if (i == 0 || boundary(line[i-1])) && (end == len(line) || line[end] == '/' || boundary(line[end])) {
			return true
		}
		start = end
	}
	return false
}
