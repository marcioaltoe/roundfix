package spec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// TaskTouchSet is the repository files one Task is known to touch and where
// each path was learned. The set is evidence from written artifacts, not an
// inference about what the Task intends to edit.
type TaskTouchSet struct {
	TaskID string
	Paths  map[string]TouchSource
}

// TouchSource names one artifact that made a Task's repository path known.
type TouchSource string

const (
	TouchFromVerification TouchSource = "verification command"
	TouchFromContext      TouchSource = "declared context"
	TouchFromPriorRun     TouchSource = "prior Run settlement commit"
)

// WaveCollision is two Tasks the Task Graph permits in one Wave that are
// known to touch at least one common repository file.
type WaveCollision struct {
	First  string
	Second string
	Paths  map[string]TouchSource
}

// Collisions reports every pair the Task Graph permits in one Wave that is
// known to touch a common repository file. It reads files and Git objects
// directly; it does not execute Verification, Git, or any other command.
func Collisions(repoRoot string, graph *Graph) ([]WaveCollision, error) {
	if graph == nil {
		return nil, errors.New("find Task Graph collisions: graph is required")
	}
	root, err := collisionRepoRoot(repoRoot)
	if err != nil {
		return nil, err
	}

	sets := make([]TaskTouchSet, len(graph.Tasks))
	for index, task := range graph.Tasks {
		paths, err := declaredTaskTouches(root, task)
		if err != nil {
			return nil, fmt.Errorf("find Task Graph collisions for Task %q: %w", task.ID, err)
		}
		sets[index] = TaskTouchSet{TaskID: task.ID, Paths: paths}
	}
	prior, err := priorRunTaskTouches(root, graph)
	if err != nil {
		return nil, fmt.Errorf("find Task Graph collisions from prior Run: %w", err)
	}
	for index := range sets {
		for path := range prior[sets[index].TaskID] {
			addTouchSource(sets[index].Paths, path, TouchFromPriorRun)
		}
	}

	ordered := taskDependencyClosure(graph.Tasks)
	var collisions []WaveCollision
	for first := 0; first < len(sets); first++ {
		for second := first + 1; second < len(sets); second++ {
			firstID, secondID := sets[first].TaskID, sets[second].TaskID
			if ordered[firstID][secondID] || ordered[secondID][firstID] {
				continue
			}
			shared := sharedTaskTouches(sets[first].Paths, sets[second].Paths)
			if len(shared) == 0 {
				continue
			}
			collisions = append(collisions, WaveCollision{
				First:  firstID,
				Second: secondID,
				Paths:  shared,
			})
		}
	}
	return collisions, nil
}

func collisionRepoRoot(repoRoot string) (string, error) {
	if strings.TrimSpace(repoRoot) == "" {
		return "", errors.New("find Task Graph collisions: repository root is required")
	}
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve repository root %q: %w", repoRoot, err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("stat repository root %q: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("repository root %q is not a directory", root)
	}
	return root, nil
}

func declaredTaskTouches(repoRoot string, task Task) (map[string]TouchSource, error) {
	paths := make(map[string]TouchSource)
	for _, command := range task.Verification {
		for _, candidate := range shellWords(command) {
			path, exists, err := repositoryFile(repoRoot, candidate)
			if err != nil {
				return nil, fmt.Errorf("inspect Verification path %q: %w", candidate, err)
			}
			if exists {
				addTouchSource(paths, path, TouchFromVerification)
			}
		}
	}
	for _, ref := range task.Context {
		path, exists, err := repositoryFile(repoRoot, ref.Path)
		if err != nil {
			return nil, fmt.Errorf("inspect Context path %q: %w", ref.Path, err)
		}
		if exists {
			addTouchSource(paths, path, TouchFromContext)
		}
	}
	return paths, nil
}

// repositoryFile reads a candidate lifted out of a Verification command or a
// declared Context reference. Such a candidate is a token from a shell line,
// not a path: it can carry surrounding whitespace, a leading dash that is
// really a flag, or a glob that names no single file. Those are filtered here
// and nowhere else.
func repositoryFile(repoRoot, candidate string) (string, bool, error) {
	candidate = strings.TrimSpace(candidate)
	if candidate == "" || strings.HasPrefix(candidate, "-") ||
		strings.ContainsAny(candidate, "$`*?[]{}") {
		return "", false, nil
	}
	return repositoryFileExact(repoRoot, candidate)
}

// repositoryFileExact reads a path exactly as given. Git reports a path
// verbatim, so trimming it or rejecting it for shell metacharacters would
// answer about a different file than the one history names — or discard a real
// one, since `foo[1].go` is an ordinary filename.
func repositoryFileExact(repoRoot, candidate string) (string, bool, error) {
	if candidate == "" {
		return "", false, nil
	}

	path := candidate
	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, filepath.FromSlash(path))
	}
	path = filepath.Clean(path)
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false, nil
	}
	relative = filepath.ToSlash(relative)
	if relative == ".git" || strings.HasPrefix(relative, ".git/") {
		return "", false, nil
	}

	resolved, err := filepath.EvalSymlinks(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	resolvedRelative, err := filepath.Rel(repoRoot, resolved)
	if err != nil || resolvedRelative == ".." || strings.HasPrefix(resolvedRelative, ".."+string(filepath.Separator)) {
		return "", false, nil
	}
	info, err := os.Stat(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if !info.Mode().IsRegular() {
		return "", false, nil
	}
	return relative, true, nil
}

// shellWords is deliberately smaller than a shell parser. It separates the
// literal words a Verification command writes down while refusing expansion;
// repositoryFile then accepts only words that already resolve to files.
func shellWords(command string) []string {
	var words []string
	var word strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if word.Len() > 0 {
			words = append(words, word.String())
			word.Reset()
		}
	}
	for _, char := range command {
		if escaped {
			word.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' && quote != '\'' {
			escaped = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				word.WriteRune(char)
			}
			continue
		}
		switch char {
		case '\'', '"':
			quote = char
		case ' ', '\t', '\r', '\n', ';', '|', '&', '(', ')', '<', '>':
			flush()
		default:
			word.WriteRune(char)
		}
	}
	if escaped {
		word.WriteRune('\\')
	}
	flush()
	return words
}

func addTouchSource(paths map[string]TouchSource, path string, source TouchSource) {
	if _, exists := paths[path]; !exists {
		paths[path] = source
	}
}

func sharedTaskTouches(first, second map[string]TouchSource) map[string]TouchSource {
	shared := make(map[string]TouchSource)
	for path, firstSource := range first {
		secondSource, exists := second[path]
		if !exists {
			continue
		}
		shared[path] = mergedTouchSource(firstSource, secondSource)
	}
	return shared
}

func mergedTouchSource(first, second TouchSource) TouchSource {
	if first == second {
		return first
	}
	sources := []string{string(first), string(second)}
	sort.Strings(sources)
	return TouchSource(strings.Join(sources, "; "))
}

func taskDependencyClosure(tasks []Task) map[string]map[string]bool {
	needs := make(map[string][]string, len(tasks))
	for _, task := range tasks {
		needs[task.ID] = append([]string(nil), task.Needs...)
	}
	closure := make(map[string]map[string]bool, len(tasks))
	for _, task := range tasks {
		seen := make(map[string]bool)
		stack := append([]string(nil), needs[task.ID]...)
		for len(stack) > 0 {
			last := len(stack) - 1
			need := stack[last]
			stack = stack[:last]
			if seen[need] {
				continue
			}
			seen[need] = true
			stack = append(stack, needs[need]...)
		}
		closure[task.ID] = seen
	}
	return closure
}

// priorRunTaskTouches reports the paths the newest prior settlement commit
// changed for each Task in the graph.
//
// It invokes git rather than reading object files, because that is what this
// repository does for exactly these questions in internal/worktree,
// internal/specaudit, internal/speccheck, and internal/cli. Trailers are read
// through git's own interpolation for the same reason: a hand-rolled scan of
// every message line would read a Roundfix-Spec mentioned in a body as though
// it were a trailer, and git already knows where the terminal block starts.
//
// Only the newest settlement per Task counts. A Task that settled in more than
// one Run — the ordinary shape when a Run ends Unresolved and the next one
// re-executes it — would otherwise union stale paths into current evidence and
// refuse a Wave that is safe.
func priorRunTaskTouches(repoRoot string, graph *Graph) (map[string]map[string]bool, error) {
	touches := make(map[string]map[string]bool)
	slug := strings.TrimSpace(graph.Spec.Slug)
	if slug == "" {
		return touches, nil
	}
	wanted := make(map[string]bool, len(graph.Tasks))
	for _, task := range graph.Tasks {
		wanted[task.ID] = true
	}

	// A repository with no commits, or none carrying the trailer, is not an
	// error: the prior-Run source is one of three, and absent evidence is not a
	// failure.
	const format = "--format=%H%x1f%(trailers:key=Roundfix-Spec,valueonly)%x1f" +
		"%(trailers:key=Roundfix-Task,valueonly)%x1e"
	// --all, not HEAD: a settlement commit lives on its Run Branch until
	// integration moves it, and a Run that ended Unresolved leaves it there.
	// Searching HEAD alone would miss the prior Run this source exists to read.
	available, err := priorRunHistoryReadable(repoRoot)
	if err != nil {
		return nil, err
	}
	if !available {
		return touches, nil
	}
	output, err := collisionGit(repoRoot, "log", "--all", "--grep=Roundfix-Task:", format)
	if err != nil {
		return nil, err
	}
	// git log reports newest first, so the first record for a Task is its
	// newest settlement and every later one is superseded.
	for _, record := range strings.Split(output, "\x1e") {
		fields := strings.Split(strings.TrimSpace(record), "\x1f")
		if len(fields) != 3 {
			continue
		}
		commit := strings.TrimSpace(fields[0])
		if strings.TrimSpace(fields[1]) != slug {
			continue
		}
		taskID := strings.TrimSpace(fields[2])
		if taskID == "" || !wanted[taskID] {
			continue
		}
		if _, seen := touches[taskID]; seen {
			continue
		}
		// The newest settlement is the evidence, readable or not: falling back
		// to an older one would answer with paths this Task has since stopped
		// touching.
		touches[taskID] = make(map[string]bool)
		// -z with quotePath off: git otherwise quotes and escapes any path
		// outside plain ASCII, and a quoted path is not the path.
		changed, err := collisionGit(repoRoot, "-c", "core.quotePath=false",
			"diff-tree", "--no-commit-id", "--name-only", "-r", "-z", "--root", commit)
		if err != nil {
			// The Task keeps its claim above, so no older settlement answers
			// in its place; only this Task's prior-Run paths are omitted, and
			// the other two sources still speak for it.
			continue
		}
		for _, line := range strings.Split(changed, "\x00") {
			// History names paths the tree no longer carries. A file two Tasks
			// both deleted is not a file they now share, so the same
			// repository-boundary and regular-file check the declared sources
			// use applies here; without it this source could report a collision
			// on a path that does not exist.
			relative, ok, err := repositoryFileExact(repoRoot, line)
			if err != nil || !ok {
				continue
			}
			touches[taskID][relative] = true
		}
	}
	return touches, nil
}

// priorRunHistoryReadable answers whether this repository can be read for
// settlement history at all.
//
// Two failures are absent evidence rather than a fault: a path outside any
// repository, because a Spec can be checked outside a working tree, and a
// repository with no commits. `rev-parse --verify --quiet HEAD` separates them
// by exit status — 1 and silent when the ref does not resolve, 128 with a
// fatal when the repository itself is the problem — so only the one fatal that
// means "no repository" is tolerated. A permission failure or a damaged object
// store reaches neither branch: reporting those as absent evidence would let
// an unsafe Wave through, which is the failure this rule exists to remove.
func priorRunHistoryReadable(repoRoot string) (bool, error) {
	if _, err := collisionGit(repoRoot, "rev-parse", "--verify", "--quiet", "HEAD"); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return false, err
		}
		stderr := strings.TrimSpace(string(exit.Stderr))
		if exit.ExitCode() == 1 && stderr == "" {
			return false, nil
		}
		if strings.Contains(stderr, "not a git repository") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func collisionGit(repoRoot string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", repoRoot}, args...)...)
	output, err := command.Output()
	if err != nil {
		// Carry git's own words: "read Git history" alone leaves a reader with
		// nothing to act on.
		var exit *exec.ExitError
		if errors.As(err, &exit) && len(exit.Stderr) > 0 {
			return "", fmt.Errorf("read Git history for Task Graph collisions: %w: %s", err, strings.TrimSpace(string(exit.Stderr)))
		}
		return "", fmt.Errorf("read Git history for Task Graph collisions: %w", err)
	}
	return string(output), nil
}
