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
	"sort"
	"strings"
)

// IsLiveDocument classifies repository-relative Markdown by its live family.
// History and a configured archive root always take precedence over the Spec Root.
func IsLiveDocument(path, specRootRel, archiveRootRel string) bool {
	if !strings.HasSuffix(path, ".md") || archivePathUnder(path, "docs/history") || archivePathUnder(path, archiveRootRel) {
		return false
	}
	return archivePathUnder(path, specRootRel) || archivePathUnder(path, "docs/backlog") || archivePathUnder(path, "docs/findings") || archivePathUnder(path, "docs/references")
}

func archivePathUnder(path, root string) bool {
	if root == "." {
		return path != "" && !filepath.IsAbs(path) && path != ".." && !strings.HasPrefix(path, "../")
	}
	return root != "" && strings.HasPrefix(path, strings.TrimSuffix(root, "/")+"/")
}

// archiveCoreArtifact is the existing --promote refusal set, relative to a Spec.
func archiveCoreArtifact(path string) bool {
	base := filepath.Base(path)
	return (filepath.Dir(path) == "." && strings.HasPrefix(base, "_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md"))
}

// ArchiveCitationMoves maps removed repository paths to their surviving homes.
// folderFiles and promotion keys are repository-relative paths.
func ArchiveCitationMoves(folderRel, recordRel string, folderFiles []string, promotions map[string]string) map[string]string {
	moves := map[string]string{strings.TrimSuffix(folderRel, "/"): recordRel}
	for _, path := range folderFiles {
		if archiveCoreArtifact(strings.TrimPrefix(path, strings.TrimSuffix(folderRel, "/")+"/")) {
			moves[path] = recordRel
		}
	}
	for source, destination := range promotions {
		moves[source] = destination
	}
	return moves
}

// RelinkContent changes only scanned relative destinations and bounded path text.
// The scanner masks code; the path pass deliberately also covers spans and fences.
func RelinkContent(content []byte, citingPath string, moves map[string]string) ([]byte, int) {
	rewritten, count, _ := (CitationRewrite{Moves: moves}).Apply(content, citingPath, citingPath)
	return rewritten, count
}

// ArchiveRelinkPlan holds automatic and manual promotions and exact rewrites.
type ArchiveRelinkPlan struct {
	Promote    []string
	Relinked   map[string][]byte
	Unrelinked []HistoryCitation
}

// archiveGitRead preserves blob bytes and keeps preflight reads from refreshing
// the index or invoking a repository filesystem monitor.
func archiveGitRead(ctx context.Context, root string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", append([]string{"-C", root, "-c", "core.fsmonitor=false"}, args...)...)
	command.Env = authorizationGitEnvironment()
	var diagnostic bytes.Buffer
	command.Stderr = &diagnostic
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(diagnostic.String()), err)
	}
	return output, nil
}

// archiveGitWorkTree distinguishes a non-Git directory from failure in a worktree.
func archiveGitWorkTree(ctx context.Context, root string) (bool, error) {
	if root == "" {
		return false, nil
	}
	output, err := archiveGitRead(ctx, root, "rev-parse", "--is-inside-work-tree")
	if err == nil {
		return strings.TrimSpace(string(output)) == "true", nil
	}
	for dir := filepath.Clean(root); ; dir = filepath.Dir(dir) {
		_, statErr := os.Lstat(filepath.Join(dir, ".git"))
		if !errors.Is(statErr, os.ErrNotExist) {
			return false, fmt.Errorf("inspect archive Git work tree: %w", errors.Join(err, statErr))
		}
		if filepath.Dir(dir) == dir {
			break
		}
	}
	if ctx.Err() == nil && strings.Contains(err.Error(), "not a git repository") {
		return false, nil
	}
	return false, fmt.Errorf("inspect archive Git work tree: %w", err)
}

func archiveFolderFiles(root, folder string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(filepath.Join(root, filepath.FromSlash(folder)), func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, err
}

// archiveContentCitations uses the history boundary rule and supplements its
// inline-link detector with the archive scanner's angle and reference syntax.
func archiveContentCitations(content []byte, path string, targets []string, folder string) []HistoryCitation {
	byLine := map[int]map[string]bool{}
	add := func(line int, target string) {
		if byLine[line] == nil {
			byLine[line] = map[string]bool{}
		}
		byLine[line][target] = true
	}
	for i, line := range strings.Split(string(content), "\n") {
		for _, target := range targets {
			if historyEntryCitesPath(line, path, target) {
				add(i+1, target)
			}
		}
	}
	for _, span := range scanArchiveMarkdown(content) {
		target, _, relative := archiveLinkTarget(string(content[span.start:span.end]), filepath.Dir(path))
		if !relative {
			continue
		}
		target = filepath.ToSlash(target)
		for _, removed := range targets {
			if target == removed || strings.HasPrefix(target, removed+"/") {
				add(bytes.Count(content[:span.start], []byte("\n"))+1, removed)
			}
		}
	}
	var citations []HistoryCitation
	for line, named := range byLine {
		if len(named) > 1 {
			delete(named, folder)
		}
		for target := range named {
			citations = append(citations, HistoryCitation{Path: path, Line: line, Target: target})
		}
	}
	sort.Slice(citations, func(i, j int) bool {
		if citations[i].Line != citations[j].Line {
			return citations[i].Line < citations[j].Line
		}
		return citations[i].Target < citations[j].Target
	})
	return citations
}

// PlanArchiveRelink reads tracked Markdown without writing. Non-Git roots
// preserve the archive's library behavior by returning an empty plan.
func PlanArchiveRelink(ctx context.Context, repositoryRoot, specRootRel, archiveRootRel, folderRel, recordRel string, manual []string) (ArchiveRelinkPlan, error) {
	plan := ArchiveRelinkPlan{}
	inside, err := archiveGitWorkTree(ctx, repositoryRoot)
	if err != nil || !inside {
		return plan, err
	}
	files, err := archiveFolderFiles(repositoryRoot, folderRel)
	if err != nil {
		return plan, fmt.Errorf("list archive citation targets: %w", err)
	}
	targets := append(append([]string(nil), files...), folderRel)
	// Keep the shared detector as the source of ordinary history citations.
	citations, err := HistoryCitations(ctx, repositoryRoot, targets)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return plan, err
	}
	output, err := archiveGitRead(ctx, repositoryRoot, "ls-files", "-z", "--", "*.md")
	if err != nil {
		return plan, fmt.Errorf("list archive citation sources: %w", err)
	}
	// Include HEAD paths so an index deletion cannot hide a Live Document that
	// the archive would relink. New index entries need no HEAD blob until committed.
	headFiles, err := archiveGitRead(ctx, repositoryRoot, "ls-tree", "-r", "--name-only", "-z", "HEAD")
	if err != nil {
		return plan, fmt.Errorf("list committed archive citation sources: %w", err)
	}
	tracked := map[string]bool{}
	atHead := map[string]bool{}
	for _, path := range strings.Split(string(output), "\x00") {
		if path != "" {
			tracked[path] = true
		}
	}
	for _, path := range strings.Split(string(headFiles), "\x00") {
		if strings.HasSuffix(path, ".md") {
			tracked[path] = true
			atHead[path] = true
		}
	}
	contents := map[string][]byte{}
	committed := map[string][]byte{}
	for path := range tracked {
		if path == "" || archivePathUnder(path, folderRel) || archivePathUnder(path, archiveRootRel) || archivePathUnder(path, "docs/history") {
			continue
		}
		content, readErr := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			return plan, fmt.Errorf("read archive citation source %q: %w", path, readErr)
		}
		if atHead[path] && (IsLiveDocument(path, specRootRel, archiveRootRel) || readErr != nil) {
			head, err := archiveGitRead(ctx, repositoryRoot, "show", "HEAD:"+path)
			if err != nil {
				return plan, fmt.Errorf("read committed archive citation source %q: %w", path, err)
			}
			committed[path] = head
			if readErr != nil {
				content = head
			}
		}
		if readErr != nil && !atHead[path] {
			return plan, fmt.Errorf("read archive citation source %q: %w", path, readErr)
		}
		contents[path] = content
	}
	// Include reference definitions and decoded destinations, which the original
	// history detector's simple inline regex does not recognize.
	byCitation := map[HistoryCitation]bool{}
	for _, citation := range citations {
		if _, ok := contents[citation.Path]; ok {
			byCitation[citation] = true
		}
	}
	paths := make([]string, 0, len(contents))
	for path := range contents {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		for _, content := range [][]byte{contents[path], committed[path]} {
			for _, citation := range archiveContentCitations(content, path, targets, folderRel) {
				byCitation[citation] = true
			}
		}
	}
	citations = citations[:0]
	for citation := range byCitation {
		citations = append(citations, citation)
	}
	promotions := map[string]string{}
	for _, candidate := range manual {
		clean := filepath.ToSlash(filepath.Clean(candidate))
		promotions[folderRel+"/"+clean] = "docs/references/" + filepath.Base(clean)
	}
	for _, citation := range citations {
		if !IsLiveDocument(citation.Path, specRootRel, archiveRootRel) || citation.Target == folderRel {
			continue
		}
		rel := strings.TrimPrefix(citation.Target, folderRel+"/")
		if !archiveCoreArtifact(rel) {
			promotions[citation.Target] = "docs/references/" + filepath.Base(rel)
		}
	}
	for source := range promotions {
		plan.Promote = append(plan.Promote, strings.TrimPrefix(source, folderRel+"/"))
	}
	sort.Strings(plan.Promote)
	moves := ArchiveCitationMoves(folderRel, recordRel, files, promotions)
	plan.Relinked = map[string][]byte{}
	for _, path := range paths {
		content := contents[path]
		if IsLiveDocument(path, specRootRel, archiveRootRel) {
			rewritten, count := RelinkContent(committed[path], path, moves)
			if count == 0 {
				rewritten, count = RelinkContent(content, path, moves)
			}
			if count > 0 {
				plan.Relinked[path] = rewritten
				content = rewritten
			}
		}
		plan.Unrelinked = append(plan.Unrelinked, archiveContentCitations(content, path, targets, folderRel)...)
	}
	return plan, nil
}
