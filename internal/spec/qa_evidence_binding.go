package spec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	BorrowedEvidenceReused          = "reused"
	BorrowedEvidenceImported        = "imported"
	BorrowedEvidenceOutsideWorktree = "outside the worktree"
	BorrowedEvidenceMalformed       = "not a list of paths"
	BorrowedEvidenceForeignWorktree = "names another Run Worktree"
)

// BorrowedEvidence identifies evidence that the report's own Run did not produce.
type BorrowedEvidence struct {
	Line            int
	Subject, Reason string
}

// BorrowedEvidenceOptions supplies the Run and QA Report Import provenance.
type BorrowedEvidenceOptions struct {
	RunBranchPrefix string
	RunBranch       string
	AuditedHead     string
	ImportCommit    string
	ImportedPaths   []string
}

// FindBorrowedEvidence checks local evidence against the report's audited head.
// Outside Git it checks nothing. Read and Git failures remain errors.
func FindBorrowedEvidence(ctx context.Context, reportPath string, options BorrowedEvidenceOptions) ([]BorrowedEvidence, error) {
	if options.RunBranchPrefix == "" {
		return nil, errors.New("QA Report evidence requires RunBranchPrefix")
	}
	content, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, fmt.Errorf("read QA Report evidence: %w", err)
	}
	path, err := filepath.EvalSymlinks(reportPath)
	if err != nil {
		return nil, fmt.Errorf("resolve QA Report path: %w", err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve QA Report path: %w", err)
	}
	output, gitErr := archiveGitRead(ctx, filepath.Dir(path), "rev-parse", "--show-toplevel")
	if gitErr != nil {
		for dir := filepath.Dir(path); ; dir = filepath.Dir(dir) {
			_, statErr := os.Lstat(filepath.Join(dir, ".git"))
			if !errors.Is(statErr, os.ErrNotExist) {
				return nil, fmt.Errorf("resolve QA Report Git root: %w", errors.Join(gitErr, statErr))
			}
			if filepath.Dir(dir) == dir {
				break
			}
		}
		if ctx.Err() != nil || !strings.Contains(gitErr.Error(), "not a git repository") {
			return nil, fmt.Errorf("resolve QA Report Git root: %w", gitErr)
		}
		return nil, nil
	}
	root := strings.TrimSpace(string(output))
	head, err := borrowedEvidenceHead(ctx, root, path, options.AuditedHead)
	if err != nil {
		return nil, err
	}
	foreign, err := borrowedEvidenceWorktrees(ctx, root, options)
	if err != nil {
		return nil, err
	}
	exempt := borrowedEvidenceCarriedLines(content)
	items, err := borrowedEvidenceExecutables(root, content)
	if err != nil {
		return nil, err
	}
	evidenceDir := filepath.Join(filepath.Dir(path), "evidence")
	for _, link := range RelativeLinks(path, content) {
		if exempt[link.Line] || !borrowedEvidenceWithin(evidenceDir, filepath.FromSlash(link.Target)) {
			continue
		}
		err := filepath.WalkDir(filepath.FromSlash(link.Target), func(file string, entry fs.DirEntry, walkErr error) error {
			if errors.Is(walkErr, os.ErrNotExist) {
				return nil
			}
			if walkErr != nil {
				return fmt.Errorf("read QA evidence %s: %w", link.Destination, walkErr)
			}
			if entry.IsDir() {
				return nil
			}
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read QA evidence %s: %w", file, err)
			}
			subject := link.Destination
			if file != filepath.FromSlash(link.Target) {
				rel, err := filepath.Rel(filepath.FromSlash(link.Target), file)
				if err != nil {
					return fmt.Errorf("resolve QA evidence member: %w", err)
				}
				subject = strings.TrimRight(subject, "/") + "/" + filepath.ToSlash(rel)
			}
			rel, err := filepath.Rel(root, file)
			if err != nil {
				return fmt.Errorf("resolve QA evidence path: %w", err)
			}
			rel = filepath.ToSlash(rel)
			for _, comparison := range []struct{ commit, reason string }{{head, BorrowedEvidenceReused}, {options.ImportCommit, BorrowedEvidenceImported}} {
				if comparison.commit == "" || comparison.reason == BorrowedEvidenceImported && !slices.Contains(options.ImportedPaths, rel) {
					continue
				}
				same, err := borrowedEvidenceSameBlob(ctx, root, comparison.commit, rel, data)
				if err != nil {
					return err
				}
				if same {
					items = append(items, BorrowedEvidence{link.Line, subject, comparison.reason})
					return nil
				}
			}
			items = append(items, borrowedEvidenceForeignMatches(data, subject, link.Line, foreign, nil)...)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	items = append(items, borrowedEvidenceForeignMatches(content, filepath.Base(path), 0, foreign, exempt)...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Line < items[j].Line })
	return items, nil
}

func borrowedEvidenceWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

func borrowedEvidenceHead(ctx context.Context, root, report, explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	rel, err := filepath.Rel(root, report)
	if err != nil {
		return "", fmt.Errorf("resolve QA Report repository path: %w", err)
	}
	status, err := archiveGitRead(ctx, root, "status", "--porcelain", "--", rel)
	if err != nil {
		return "", fmt.Errorf("read QA Report status: %w", err)
	}
	// for-each-ref distinguishes an unborn branch without treating other Git
	// failures as an absent audited head.
	refs, err := archiveGitRead(ctx, root, "for-each-ref", "--format=%(refname)")
	if err != nil {
		return "", fmt.Errorf("read QA Report repository refs: %w", err)
	}
	head, headErr := archiveGitRead(ctx, root, "rev-parse", "--verify", "HEAD")
	if headErr != nil {
		branch, err := archiveGitRead(ctx, root, "symbolic-ref", "HEAD")
		if err == nil && !slices.Contains(strings.Fields(string(refs)), strings.TrimSpace(string(branch))) {
			return "", nil
		}
		return "", fmt.Errorf("read QA Report HEAD: %w", errors.Join(headErr, err))
	}
	if len(status) != 0 {
		return strings.TrimSpace(string(head)), nil
	}
	addition, err := archiveGitRead(ctx, root, "log", "--diff-filter=A", "--format=%H", "-1", "--", rel)
	if err != nil {
		return "", fmt.Errorf("read QA Report addition: %w", err)
	}
	if len(bytes.TrimSpace(addition)) == 0 {
		return "", nil
	}
	parents, err := archiveGitRead(ctx, root, "rev-list", "--parents", "-1", strings.TrimSpace(string(addition)))
	if err != nil {
		return "", fmt.Errorf("read QA Report audited parent: %w", err)
	}
	fields := strings.Fields(string(parents))
	if len(fields) < 2 {
		return "", nil
	}
	return fields[1], nil
}

func borrowedEvidenceSameBlob(ctx context.Context, root, commit, path string, data []byte) (bool, error) {
	listed, err := archiveGitRead(ctx, root, "ls-tree", "-z", commit, "--", path)
	if err != nil {
		return false, fmt.Errorf("read QA evidence tree: %w", err)
	}
	if len(listed) == 0 {
		return false, nil
	}
	fields := strings.Fields(strings.SplitN(string(listed), "\t", 2)[0])
	if len(fields) != 3 || fields[1] != "blob" {
		return false, nil
	}
	blob, err := archiveGitRead(ctx, root, "cat-file", "blob", fields[2])
	if err != nil {
		return false, fmt.Errorf("read QA evidence blob: %w", err)
	}
	return bytes.Equal(blob, data), nil
}

type borrowedEvidenceWorktree struct{ path, resolved string }

func borrowedEvidenceWorktrees(ctx context.Context, root string, options BorrowedEvidenceOptions) ([]borrowedEvidenceWorktree, error) {
	output, err := archiveGitRead(ctx, root, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("list QA evidence worktrees: %w", err)
	}
	type worktree struct{ path, branch string }
	var trees []worktree
	for _, block := range strings.Split(strings.TrimSpace(string(output)), "\n\n") {
		var tree worktree
		for _, line := range strings.Split(block, "\n") {
			if strings.HasPrefix(line, "worktree ") {
				tree.path = strings.TrimPrefix(line, "worktree ")
			} else if strings.HasPrefix(line, "branch refs/heads/") {
				tree.branch = strings.TrimPrefix(line, "branch refs/heads/")
			}
		}
		trees = append(trees, tree)
	}
	own := options.RunBranch
	for _, tree := range trees {
		if own == "" && tree.path == root && strings.HasPrefix(tree.branch, options.RunBranchPrefix) {
			own = tree.branch
		}
	}
	var foreign []borrowedEvidenceWorktree
	for _, tree := range trees {
		if tree.path == root || !strings.HasPrefix(tree.branch, options.RunBranchPrefix) || own != "" && (tree.branch == own || strings.HasPrefix(tree.branch, own+"-")) {
			continue
		}
		resolved, err := filepath.EvalSymlinks(tree.path)
		if err != nil {
			return nil, fmt.Errorf("resolve Run Worktree %s: %w", tree.path, err)
		}
		foreign = append(foreign, borrowedEvidenceWorktree{tree.path, resolved})
	}
	return foreign, nil
}

func borrowedEvidenceForeignMatches(data []byte, subject string, line int, trees []borrowedEvidenceWorktree, exempt map[int]bool) []BorrowedEvidence {
	var items []BorrowedEvidence
	for _, tree := range trees {
		for i, text := range strings.Split(string(data), "\n") {
			if exempt[i+1] {
				continue
			}
			if strings.Contains(text, tree.path+string(filepath.Separator)) || strings.Contains(text, tree.resolved+string(filepath.Separator)) {
				matchLine := line
				if matchLine == 0 {
					matchLine = i + 1
				}
				items = append(items, BorrowedEvidence{matchLine, subject, BorrowedEvidenceForeignWorktree + " " + tree.path})
				break
			}
		}
	}
	return items
}

func borrowedEvidenceExecutables(root string, content []byte) ([]BorrowedEvidence, error) {
	frontmatter, _, err := splitQAReportFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("read QA evidence front matter: %w", err)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(frontmatter, &doc); err != nil {
		return nil, fmt.Errorf("read QA evidence executables: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, nil
	}
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolve QA evidence root: %w", err)
	}
	fields := doc.Content[0].Content
	for i := 0; i+1 < len(fields); i += 2 {
		key, value := fields[i], fields[i+1]
		if key.Value != "user_flow_executables" {
			continue
		}
		malformed := []BorrowedEvidence{{key.Line + 1, key.Value, BorrowedEvidenceMalformed}}
		if value.Kind != yaml.SequenceNode {
			return malformed, nil
		}
		var items []BorrowedEvidence
		for _, entry := range value.Content {
			if entry.Kind != yaml.ScalarNode || entry.Tag != "!!str" {
				return malformed, nil
			}
			path := filepath.Clean(entry.Value)
			if !filepath.IsAbs(path) {
				path = filepath.Join(root, path)
			}
			if !borrowedEvidenceWithin(root, path) && !borrowedEvidenceWithin(resolved, path) {
				items = append(items, BorrowedEvidence{entry.Line + 1, entry.Value, BorrowedEvidenceOutsideWorktree})
			}
		}
		return items, nil
	}
	return nil, nil
}

func borrowedEvidenceCarriedLines(content []byte) map[int]bool {
	lines := strings.Split(string(content), "\n")
	exempt, carried := map[int]bool{}, map[string]bool{}
	inResults, statusColumn := false, -1
	fence := ""
	for i, line := range lines {
		text := strings.TrimSpace(line)
		if marker, ok := markdownFenceMarker(text); ok {
			if fence == "" {
				fence = marker
			} else if fence == marker {
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		if depth, heading, ok := markdownHeading(text); ok && depth <= 2 {
			inResults, statusColumn = depth == 2 && heading == "Results", -1
		}
		if !inResults {
			continue
		}
		cells := markdownTableCells(line)
		if i+1 < len(lines) && len(cells) > 0 && markdownTableSeparator(markdownTableCells(lines[i+1])) {
			statusColumn = -1
			for column, cell := range cells {
				if strings.EqualFold(strings.TrimSpace(cell), "Status") {
					statusColumn = column
				}
			}
		} else if statusColumn >= 0 && statusColumn < len(cells) && strings.HasPrefix(strings.TrimSpace(cells[statusColumn]), "carried (") {
			exempt[i+1], carried[strings.TrimSpace(cells[0])] = true, true
		}
	}
	block, fence := false, ""
	for i, line := range lines {
		text := strings.TrimSpace(line)
		if marker, ok := markdownFenceMarker(text); ok {
			if fence == "" {
				fence = marker
			} else if fence == marker {
				fence = ""
			}
		} else if fence == "" {
			if depth, heading, ok := markdownHeading(text); ok && depth <= 3 {
				block = depth == 3 && carried[heading]
			}
		}
		if block {
			exempt[i+1] = true
		}
	}
	return exempt
}

// BorrowedEvidenceError is the shared refusal of evidence from another Run.
type BorrowedEvidenceError struct{ Evidence []BorrowedEvidence }

func (err BorrowedEvidenceError) Error() string {
	parts := make([]string, 0, len(err.Evidence))
	for _, item := range err.Evidence {
		parts = append(parts, fmt.Sprintf("%s (line %d, %s)", item.Subject, item.Line, item.Reason))
	}
	return "QA Report evidence does not belong to this Run: " + strings.Join(parts, "; ")
}

// RequireOwnRunEvidence returns the shared refusal or an infrastructure error.
func RequireOwnRunEvidence(ctx context.Context, reportPath string, options BorrowedEvidenceOptions) error {
	items, err := FindBorrowedEvidence(ctx, reportPath, options)
	if err != nil {
		return err
	}
	if len(items) > 0 {
		return BorrowedEvidenceError{Evidence: items}
	}
	return nil
}
