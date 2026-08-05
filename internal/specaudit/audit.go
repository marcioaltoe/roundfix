// Package specaudit reports the local Git survivors associated with a Spec.
// It reads the Run Database and local repository without changing either.
package specaudit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"roundfix/internal/store"
)

// Kind is a survivor classification.
type Kind string

const (
	KindPullRequest Kind = "pull-request"
	KindPending     Kind = "pending"
	KindResidue     Kind = "residue"
	KindPreserved   Kind = "preserved"
)

// Survivor is one branch or worktree that remains after a Spec cycle.
type Survivor struct {
	Name       string `json:"name"`
	IsWorktree bool   `json:"is_worktree"`
	Kind       Kind   `json:"kind"`
	Evidence   string `json:"evidence"`
	Reclaim    string `json:"reclaim,omitempty"`
}

// Undelivered identifies a claimed artifact absent from the default branch.
type Undelivered struct {
	Artifact string `json:"artifact"`
	HeldBy   string `json:"held_by"`
}

// Result is the read-only audit result for one Spec.
type Result struct {
	Slug        string        `json:"slug"`
	Survivors   []Survivor    `json:"survivors"`
	Undelivered []Undelivered `json:"undelivered"`
}

type gitRunner interface {
	Run(ctx context.Context, workDir string, args ...string) (string, error)
}

type execGitRunner struct{}

func (execGitRunner) Run(ctx context.Context, workDir string, args ...string) (string, error) {
	gitArgs := []string{"--no-optional-locks"}
	if strings.TrimSpace(workDir) != "" {
		gitArgs = append(gitArgs, "-C", workDir)
	}
	gitArgs = append(gitArgs, "-c", "core.fsmonitor=false")
	gitArgs = append(gitArgs, args...)
	command := exec.CommandContext(ctx, "git", gitArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

type worktree struct {
	Path   string
	Head   string
	Branch string
}

type auditInputs struct {
	runsByWorktree  map[string]store.Run
	runsByRunBranch map[string]store.Run
	associated      map[string]bool
	active          map[string]store.Run
	pullRequests    map[string]store.Run
	defaultRef      string
	defaultName     string
}

// Audit reads local Git state and the Run Database. It mutates nothing.
func Audit(ctx context.Context, repoRoot, homeDir, slug string) (result Result, err error) {
	repoRoot = strings.TrimSpace(repoRoot)
	homeDir = strings.TrimSpace(homeDir)
	slug = strings.TrimSpace(slug)
	if repoRoot == "" {
		return Result{}, errors.New("audit Spec survivors: repository root is required")
	}
	if homeDir == "" {
		return Result{}, errors.New("audit Spec survivors: home directory is required")
	}
	if slug == "" {
		return Result{}, errors.New("audit Spec survivors: Spec slug is required")
	}

	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec survivors: %w", err)
	}
	defer func() {
		if closeErr := reader.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("close Run Database reader: %w", closeErr)
		}
	}()

	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{
		GitRoot: repoRoot,
		States:  store.StatesAll,
	})
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec survivors: list Runs: %w", err)
	}
	specRuns := make([]store.Run, 0, len(runs))
	for _, run := range runs {
		if strings.TrimSpace(run.SpecSlug) == slug {
			specRuns = append(specRuns, run)
		}
	}

	runner := execGitRunner{}
	return audit(ctx, runner, repoRoot, slug, specRuns)
}

func audit(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	slug string,
	runs []store.Run,
) (Result, error) {
	branches, err := listBranches(ctx, runner, repoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec survivors: enumerate branches: %w", err)
	}
	worktrees, err := listWorktrees(ctx, runner, repoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec survivors: enumerate worktrees: %w", err)
	}
	defaultRef, defaultName, err := defaultBranch(ctx, runner, repoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec survivors: resolve default branch: %w", err)
	}

	inputs := indexRuns(runs)
	inputs.defaultRef = defaultRef
	inputs.defaultName = defaultName
	for _, branch := range branches {
		if branch == defaultName || inputs.associated[branch] {
			continue
		}
		belongs, err := branchTipBelongsToSpec(ctx, runner, repoRoot, branch, slug)
		if err != nil {
			return Result{}, fmt.Errorf("audit Spec survivors: inspect branch %q provenance: %w", branch, err)
		}
		if belongs {
			inputs.associated[branch] = true
		}
	}

	result := Result{
		Slug:        slug,
		Survivors:   []Survivor{},
		Undelivered: []Undelivered{},
	}
	for _, branch := range branches {
		if branch == defaultName || !inputs.associated[branch] {
			continue
		}
		result.Survivors = append(result.Survivors, classifyBranch(ctx, runner, repoRoot, branch, inputs))
	}
	for _, candidate := range worktrees {
		matchingRun, hasMatchingRun := matchingWorktreeRun(candidate, inputs)
		if !hasMatchingRun && !inputs.associated[candidate.Branch] {
			continue
		}
		result.Survivors = append(
			result.Survivors,
			classifyWorktree(ctx, runner, repoRoot, candidate, matchingRun, hasMatchingRun, inputs),
		)
	}
	sort.Slice(result.Survivors, func(left, right int) bool {
		if result.Survivors[left].IsWorktree != result.Survivors[right].IsWorktree {
			return !result.Survivors[left].IsWorktree
		}
		return result.Survivors[left].Name < result.Survivors[right].Name
	})
	return result, nil
}

func indexRuns(runs []store.Run) auditInputs {
	inputs := auditInputs{
		runsByWorktree:  make(map[string]store.Run),
		runsByRunBranch: make(map[string]store.Run),
		associated:      make(map[string]bool),
		active:          make(map[string]store.Run),
		pullRequests:    make(map[string]store.Run),
	}
	for _, run := range runs {
		branches := runBranches(run)
		for _, branch := range branches {
			inputs.associated[branch] = true
			if !store.IsTerminalState(run.State) {
				inputs.active[branch] = run
			}
			if strings.TrimSpace(run.PRNumber) != "" {
				if _, exists := inputs.pullRequests[branch]; !exists {
					inputs.pullRequests[branch] = run
				}
			}
		}
		if run.Kind == store.KindImplement {
			runBranch := store.RunBranchPrefix + run.ID
			inputs.runsByRunBranch[runBranch] = run
		}
		if workDir := cleanPath(run.WorkDir); workDir != "" {
			existing, exists := inputs.runsByWorktree[workDir]
			if !exists || (store.IsTerminalState(existing.State) && !store.IsTerminalState(run.State)) {
				inputs.runsByWorktree[workDir] = run
			}
		}
	}
	return inputs
}

func runBranches(run store.Run) []string {
	seen := make(map[string]bool)
	branches := make([]string, 0, 3)
	for _, branch := range []string{run.HeadBranch, run.LocalBranch} {
		branch = strings.TrimSpace(branch)
		if branch != "" && !seen[branch] {
			seen[branch] = true
			branches = append(branches, branch)
		}
	}
	if run.Kind == store.KindImplement {
		branch := store.RunBranchPrefix + run.ID
		if !seen[branch] {
			branches = append(branches, branch)
		}
	}
	return branches
}

func classifyBranch(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	branch string,
	inputs auditInputs,
) Survivor {
	if run, exists := inputs.pullRequests[branch]; exists {
		return pullRequestSurvivor(branch, false, run)
	}
	if run, exists := inputs.active[branch]; exists {
		return activeRunSurvivor(branch, false, run)
	}
	classification := classifyGitRef(ctx, runner, repoRoot, branch, inputs.defaultRef, inputs.defaultName)
	classification.Name = branch
	if classification.Kind == KindResidue {
		classification.Reclaim = "git branch -d -- " + shellQuote(branch)
	}
	return classification
}

func classifyWorktree(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	candidate worktree,
	run store.Run,
	hasMatchingRun bool,
	inputs auditInputs,
) Survivor {
	if !hasMatchingRun {
		return Survivor{
			Name:       candidate.Path,
			IsWorktree: true,
			Kind:       KindPreserved,
			Evidence:   fmt.Sprintf("worktree %q has no matching Run in the Run Database", candidate.Path),
		}
	}
	if pullRequestRun, exists := inputs.pullRequests[candidate.Branch]; exists {
		return pullRequestSurvivor(candidate.Path, true, pullRequestRun)
	}
	if !store.IsTerminalState(run.State) {
		return activeRunSurvivor(candidate.Path, true, run)
	}
	dirty, evidence, err := worktreeChanges(ctx, runner, candidate.Path)
	if err != nil {
		return Survivor{
			Name:       candidate.Path,
			IsWorktree: true,
			Kind:       KindPreserved,
			Evidence:   fmt.Sprintf("worktree %q could not be inspected and is preserved: %v", candidate.Path, err),
		}
	}
	if dirty {
		return Survivor{
			Name:       candidate.Path,
			IsWorktree: true,
			Kind:       KindPending,
			Evidence:   evidence,
		}
	}
	ref := candidate.Branch
	if ref == "" {
		ref = candidate.Head
	}
	classification := classifyGitRef(ctx, runner, repoRoot, ref, inputs.defaultRef, inputs.defaultName)
	classification.Name = candidate.Path
	classification.IsWorktree = true
	if classification.Kind == KindResidue {
		classification.Reclaim = "git worktree remove -- " + shellQuote(candidate.Path)
	}
	return classification
}

func pullRequestSurvivor(name string, isWorktree bool, run store.Run) Survivor {
	number := strings.TrimPrefix(strings.TrimSpace(run.PRNumber), "#")
	return Survivor{
		Name:       name,
		IsWorktree: isWorktree,
		Kind:       KindPullRequest,
		Evidence: fmt.Sprintf(
			"survivor backs Pull Request #%s recorded by Run %s",
			number,
			run.ID,
		),
	}
}

func activeRunSurvivor(name string, isWorktree bool, run store.Run) Survivor {
	return Survivor{
		Name:       name,
		IsWorktree: isWorktree,
		Kind:       KindPreserved,
		Evidence: fmt.Sprintf(
			"Active Run %s owns this survivor in state %s; it is preserved",
			run.ID,
			run.State,
		),
	}
}

func classifyGitRef(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	ref string,
	defaultRef string,
	defaultName string,
) Survivor {
	if defaultRef == "" {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: "default branch could not be resolved; survivor is preserved",
		}
	}
	uniqueCommits, err := countUniqueCommits(ctx, runner, repoRoot, ref, defaultRef)
	if err != nil {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: fmt.Sprintf("integration could not be inspected against default branch %q: %v", defaultName, err),
		}
	}
	if uniqueCommits == 0 {
		return Survivor{
			Kind: KindResidue,
			Evidence: fmt.Sprintf(
				"survivor is reachable from default branch %q",
				defaultName,
			),
		}
	}
	runOnly, differingShared, err := compareContent(ctx, runner, repoRoot, ref, defaultRef)
	if err != nil {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: fmt.Sprintf("content could not be compared with default branch %q: %v", defaultName, err),
		}
	}
	if runOnly == 0 && differingShared == 0 {
		return Survivor{
			Kind: KindResidue,
			Evidence: fmt.Sprintf(
				"survivor content is fully represented on default branch %q",
				defaultName,
			),
		}
	}
	return Survivor{
		Kind: KindPending,
		Evidence: fmt.Sprintf(
			"%d commit%s not represented on default branch %q: %d branch-only file%s, %d differing shared file%s",
			uniqueCommits,
			pluralSuffix(uniqueCommits),
			defaultName,
			runOnly,
			pluralSuffix(runOnly),
			differingShared,
			pluralSuffix(differingShared),
		),
	}
}

func countUniqueCommits(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	ref string,
	defaultRef string,
) (int, error) {
	output, err := runner.Run(ctx, repoRoot, "rev-list", "--count", ref, "--not", defaultRef)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return 0, fmt.Errorf("parse unique commit count %q: %w", output, err)
	}
	return count, nil
}

func compareContent(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	ref string,
	defaultRef string,
) (runOnly int, differingShared int, err error) {
	runOnlyOutput, err := runner.Run(
		ctx,
		repoRoot,
		"diff",
		"--name-only",
		"-z",
		"--no-renames",
		"--diff-filter=D",
		ref,
		defaultRef,
		"--",
	)
	if err != nil {
		return 0, 0, err
	}
	differingOutput, err := runner.Run(
		ctx,
		repoRoot,
		"diff",
		"--name-only",
		"-z",
		"--no-renames",
		"--diff-filter=MT",
		ref,
		defaultRef,
		"--",
	)
	if err != nil {
		return 0, 0, err
	}
	return countNULTerms(runOnlyOutput), countNULTerms(differingOutput), nil
}

func listBranches(ctx context.Context, runner gitRunner, repoRoot string) ([]string, error) {
	output, err := runner.Run(ctx, repoRoot, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return nil, err
	}
	branches := nonEmptyLines(output)
	sort.Strings(branches)
	return branches, nil
}

func listWorktrees(ctx context.Context, runner gitRunner, repoRoot string) ([]worktree, error) {
	output, err := runner.Run(ctx, repoRoot, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}
	return parseWorktrees(output), nil
}

func parseWorktrees(output string) []worktree {
	worktrees := []worktree{}
	current := worktree{}
	appendCurrent := func() {
		if current.Path != "" {
			worktrees = append(worktrees, current)
		}
		current = worktree{}
	}
	for _, field := range strings.Split(output, "\x00") {
		if field == "" {
			appendCurrent()
			continue
		}
		key, value, _ := strings.Cut(field, " ")
		switch key {
		case "worktree":
			if current.Path != "" {
				appendCurrent()
			}
			current.Path = cleanPath(value)
		case "HEAD":
			current.Head = strings.TrimSpace(value)
		case "branch":
			current.Branch = strings.TrimPrefix(strings.TrimSpace(value), "refs/heads/")
		}
	}
	appendCurrent()
	return worktrees
}

func defaultBranch(ctx context.Context, runner gitRunner, repoRoot string) (ref string, name string, err error) {
	candidates := []struct {
		ref  string
		name string
	}{
		{ref: "refs/remotes/origin/HEAD", name: "origin/HEAD"},
		{ref: "refs/heads/main", name: "main"},
		{ref: "refs/heads/master", name: "master"},
	}
	for _, candidate := range candidates {
		output, err := runner.Run(ctx, repoRoot, "for-each-ref", "--format=%(refname)", candidate.ref)
		if err != nil {
			return "", "", err
		}
		if strings.TrimSpace(output) != "" {
			return candidate.ref, candidate.name, nil
		}
	}
	return "", "", nil
}

func branchTipBelongsToSpec(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	branch string,
	slug string,
) (bool, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"log",
		"-1",
		"--format=%(trailers:key=Roundfix-Spec,valueonly)",
		branch,
		"--",
	)
	if err != nil {
		return false, err
	}
	for _, candidate := range nonEmptyLines(output) {
		if strings.TrimSpace(candidate) == slug {
			return true, nil
		}
	}
	return false, nil
}

func matchingWorktreeRun(candidate worktree, inputs auditInputs) (store.Run, bool) {
	if run, exists := inputs.runsByWorktree[cleanPath(candidate.Path)]; exists {
		return run, true
	}
	if run, exists := inputs.runsByRunBranch[candidate.Branch]; exists {
		return run, true
	}
	return store.Run{}, false
}

func worktreeChanges(ctx context.Context, runner gitRunner, worktreePath string) (bool, string, error) {
	output, err := runner.Run(ctx, worktreePath, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return false, "", err
	}
	changed := nonEmptyLines(output)
	if len(changed) == 0 {
		return false, "", nil
	}
	return true, fmt.Sprintf(
		"worktree %q holds %d unintegrated changed path%s",
		worktreePath,
		len(changed),
		pluralSuffix(len(changed)),
	), nil
}

func countNULTerms(output string) int {
	count := 0
	for _, term := range strings.Split(output, "\x00") {
		if term != "" {
			count++
		}
	}
	return count
}

func nonEmptyLines(output string) []string {
	lines := []string{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved
	}
	return path
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
