// Package specaudit reports the local Git survivors associated with a Spec.
// It reads the Run Database and local repository without changing either.
package specaudit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
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

type branchRef struct {
	Name         string
	Remote       string
	RemoteBranch string
}

type auditInputs struct {
	runsByWorktree  map[string]store.Run
	runsByRunBranch map[string]store.Run
	remoteRefs      map[string][]string
	associated      map[string]bool
	active          map[string]store.Run
	pullRequests    map[string]store.Run
	defaultRef      string
	defaultName     string
}

type deliveryTree struct {
	repoRoot     string
	artifactRoot string
	defaultRef   string
	defaultName  string
	branches     []string
}

type artifactClaim struct {
	path string
	tree deliveryTree
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
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{
		HomeDir: homeDir,
		WorkDir: repoRoot,
		Stderr:  io.Discard,
	})
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: load configuration: %w", err)
	}
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, repoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: %w", err)
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
	return audit(ctx, runner, repoRoot, resolvedSpecsRoot.Path, resolvedSpecsRoot.BuiltInRoot, slug, specRuns)
}

func audit(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	specsRoot string,
	builtInRoot bool,
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
	defaultBranches := defaultBranchAliases(ctx, runner, repoRoot, defaultRef, defaultName, branches)
	deliveryBranches, err := listDeliveryBranches(ctx, runner, repoRoot)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: enumerate branches: %w", err)
	}
	codeTree := deliveryTree{
		repoRoot:    cleanPath(repoRoot),
		defaultRef:  defaultRef,
		defaultName: defaultName,
		branches:    deliveryBranches,
	}
	specTree, err := resolveSpecDeliveryTree(ctx, runner, specsRoot, codeTree)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: resolve Spec Root tree: %w", err)
	}

	inputs := indexRuns(runs)
	inputs.defaultRef = defaultRef
	inputs.defaultName = defaultName
	for _, branch := range branches {
		if branch.RemoteBranch != "" {
			inputs.remoteRefs[branch.RemoteBranch] = append(inputs.remoteRefs[branch.RemoteBranch], branch.Name)
		}
		if defaultBranches[branch.Name] || branchAssociated(branch, inputs.associated) {
			continue
		}
		belongs, err := branchTipBelongsToSpec(ctx, runner, repoRoot, branch.Name, slug)
		if err != nil {
			return Result{}, fmt.Errorf("audit Spec survivors: inspect branch %q provenance: %w", branch.Name, err)
		}
		if belongs {
			inputs.associated[branch.Name] = true
		}
	}

	result := Result{
		Slug:        slug,
		Survivors:   []Survivor{},
		Undelivered: []Undelivered{},
	}
	localBranchSurvivors := make(map[string]int)
	for _, branch := range branches {
		if defaultBranches[branch.Name] || !branchAssociated(branch, inputs.associated) {
			continue
		}
		result.Survivors = append(result.Survivors, classifyBranch(ctx, runner, repoRoot, branch, inputs))
		if branch.Remote == "" {
			localBranchSurvivors[branch.Name] = len(result.Survivors) - 1
		}
	}
	scratchWorktreeSurvivors := make(map[string][]int)
	for _, candidate := range worktrees {
		matchingRun, hasMatchingRun := matchingWorktreeRun(candidate, inputs)
		if !hasMatchingRun && !inputs.associated[candidate.Branch] {
			continue
		}
		result.Survivors = append(
			result.Survivors,
			classifyWorktree(ctx, runner, repoRoot, candidate, matchingRun, hasMatchingRun, inputs),
		)
		if !hasMatchingRun && candidate.Branch != "" {
			scratchWorktreeSurvivors[candidate.Branch] = append(
				scratchWorktreeSurvivors[candidate.Branch],
				len(result.Survivors)-1,
			)
		}
	}
	coordinateScratchReclaims(result.Survivors, localBranchSurvivors, scratchWorktreeSurvivors)
	sort.Slice(result.Survivors, func(left, right int) bool {
		if result.Survivors[left].IsWorktree != result.Survivors[right].IsWorktree {
			return !result.Survivors[left].IsWorktree
		}
		return result.Survivors[left].Name < result.Survivors[right].Name
	})
	claimed, err := claimedArtifacts(
		ctx,
		runner,
		builtInRoot,
		slug,
		codeTree,
		specTree,
	)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: resolve claimed artifacts: %w", err)
	}
	result.Undelivered, err = findUndeliveredArtifacts(
		ctx,
		runner,
		claimed,
	)
	if err != nil {
		return Result{}, fmt.Errorf("audit Spec delivery: %w", err)
	}
	return result, nil
}

func claimedArtifacts(
	ctx context.Context,
	runner gitRunner,
	builtInRoot bool,
	slug string,
	codeTree deliveryTree,
	specTree deliveryTree,
) ([]artifactClaim, error) {
	activeSpec := filepath.ToSlash(filepath.Join(specTree.artifactRoot, slug))
	archivedSpec := archivedSpecArtifactPath(builtInRoot, specTree.artifactRoot, slug)

	archiveClaimed, err := artifactClaimed(ctx, runner, specTree, archivedSpec)
	if err != nil {
		return nil, err
	}
	specArtifact := activeSpec
	if archiveClaimed {
		specArtifact = archivedSpec
	}

	taskArtifacts, err := taskCommitArtifacts(ctx, runner, codeTree.repoRoot, slug)
	if err != nil {
		return nil, err
	}
	claimed := map[string]artifactClaim{
		specTree.repoRoot + "\x00" + specArtifact: {path: specArtifact, tree: specTree},
	}
	for _, artifact := range taskArtifacts {
		if archiveClaimed && codeTree.repoRoot == specTree.repoRoot &&
			(artifact == activeSpec || strings.HasPrefix(artifact, activeSpec+"/")) {
			artifact = archivedSpec + strings.TrimPrefix(artifact, activeSpec)
		}
		claimed[codeTree.repoRoot+"\x00"+artifact] = artifactClaim{path: artifact, tree: codeTree}
	}
	artifacts := make([]artifactClaim, 0, len(claimed))
	for _, artifact := range claimed {
		artifacts = append(artifacts, artifact)
	}
	sort.Slice(artifacts, func(left, right int) bool {
		if artifacts[left].path == artifacts[right].path {
			return artifacts[left].tree.repoRoot < artifacts[right].tree.repoRoot
		}
		return artifacts[left].path < artifacts[right].path
	})
	return artifacts, nil
}

func archivedSpecArtifactPath(builtInRoot bool, artifactRoot, slug string) string {
	// artifactRoot is a repo-relative path inside the owning git tree, so the
	// archive location is resolved within that tree using the Spec Root's
	// built-in classification.
	archiveDir := spec.ArchiveSpecRoot(filepath.FromSlash(artifactRoot), builtInRoot)
	return filepath.ToSlash(filepath.Join(archiveDir, slug))
}

func artifactClaimed(
	ctx context.Context,
	runner gitRunner,
	tree deliveryTree,
	artifact string,
) (bool, error) {
	if _, err := os.Stat(filepath.Join(tree.repoRoot, filepath.FromSlash(artifact))); err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("stat claimed artifact %q: %w", artifact, err)
	}
	if tree.defaultRef != "" {
		exists, err := treeHasArtifact(ctx, runner, tree.repoRoot, tree.defaultRef, artifact)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}
	for _, branch := range tree.branches {
		exists, err := treeHasArtifact(ctx, runner, tree.repoRoot, branch, artifact)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}
	return false, nil
}

func taskCommitArtifacts(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	slug string,
) ([]string, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"log",
		"--all",
		"--reverse",
		"--topo-order",
		"--format=%H",
		"--fixed-strings",
		"--grep=Roundfix-Spec: "+slug,
	)
	if err != nil {
		return nil, fmt.Errorf("list Task commits: %w", err)
	}
	artifacts := map[string]bool{}
	for _, commit := range nonEmptyLines(output) {
		isTask, err := isTaskCommit(ctx, runner, repoRoot, commit, slug)
		if err != nil {
			return nil, err
		}
		if !isTask {
			continue
		}
		changes, err := taskCommitChanges(ctx, runner, repoRoot, commit)
		if err != nil {
			return nil, err
		}
		for _, change := range changes {
			if change.deleted {
				delete(artifacts, change.path)
				continue
			}
			artifacts[change.path] = true
		}
	}
	paths := make([]string, 0, len(artifacts))
	for path := range artifacts {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func isTaskCommit(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	commit string,
	slug string,
) (bool, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"log",
		"-1",
		"--format=%(trailers:key=Roundfix-Spec,valueonly)%x00%(trailers:key=Roundfix-Task,valueonly)",
		commit,
		"--",
	)
	if err != nil {
		return false, fmt.Errorf("inspect Task commit %q: %w", commit, err)
	}
	specValues, taskValues, ok := strings.Cut(output, "\x00")
	if !ok {
		return false, fmt.Errorf("inspect Task commit %q: missing trailer delimiter", commit)
	}
	specMatches := false
	for _, candidate := range nonEmptyLines(specValues) {
		if candidate == slug {
			specMatches = true
			break
		}
	}
	return specMatches && len(nonEmptyLines(taskValues)) > 0, nil
}

type taskCommitChange struct {
	path    string
	deleted bool
}

func taskCommitChanges(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	commit string,
) ([]taskCommitChange, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"diff-tree",
		"--root",
		"--no-commit-id",
		"--name-status",
		"--no-renames",
		"-r",
		"-z",
		commit,
		"--",
	)
	if err != nil {
		return nil, fmt.Errorf("read Task commit %q artifacts: %w", commit, err)
	}
	fields := strings.Split(output, "\x00")
	changes := make([]taskCommitChange, 0, len(fields)/2)
	for index := 0; index+1 < len(fields); index += 2 {
		status := strings.TrimSpace(fields[index])
		path, err := cleanGitArtifactPath(fields[index+1])
		if err != nil {
			return nil, fmt.Errorf("read Task commit %q artifacts: %w", commit, err)
		}
		if status == "" || path == "" {
			continue
		}
		changes = append(changes, taskCommitChange{path: path, deleted: status == "D"})
	}
	return changes, nil
}

func cleanGitArtifactPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("Git returned absolute artifact path %q", path)
	}
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("Git returned artifact path outside the repository %q", path)
	}
	return filepath.ToSlash(clean), nil
}

func findUndeliveredArtifacts(
	ctx context.Context,
	runner gitRunner,
	claimed []artifactClaim,
) ([]Undelivered, error) {
	undelivered := []Undelivered{}
	for _, claim := range claimed {
		delivered := false
		if claim.tree.defaultRef != "" {
			var err error
			delivered, err = treeHasArtifact(
				ctx,
				runner,
				claim.tree.repoRoot,
				claim.tree.defaultRef,
				claim.path,
			)
			if err != nil {
				return nil, fmt.Errorf("inspect default branch artifact %q: %w", claim.path, err)
			}
		}
		if delivered {
			continue
		}
		heldBy, err := holdingBranch(
			ctx,
			runner,
			claim.tree.repoRoot,
			claim.path,
			claim.tree.defaultName,
			claim.tree.branches,
		)
		if err != nil {
			return nil, err
		}
		undelivered = append(undelivered, Undelivered{Artifact: claim.path, HeldBy: heldBy})
	}
	return undelivered, nil
}

func holdingBranch(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	artifact string,
	defaultName string,
	branches []string,
) (string, error) {
	for _, branch := range branches {
		if branch == defaultName {
			continue
		}
		exists, err := treeHasArtifact(ctx, runner, repoRoot, branch, artifact)
		if err != nil {
			return "", fmt.Errorf("inspect branch %q for artifact %q: %w", branch, artifact, err)
		}
		if exists {
			return branch, nil
		}
	}
	return "", nil
}

func treeHasArtifact(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	ref string,
	artifact string,
) (bool, error) {
	output, err := runner.Run(ctx, repoRoot, "ls-tree", "-z", "--name-only", ref, "--", artifact)
	if err != nil {
		return false, fmt.Errorf("read tree %q artifact %q: %w", ref, artifact, err)
	}
	for _, candidate := range strings.Split(output, "\x00") {
		if candidate == artifact {
			return true, nil
		}
	}
	return false, nil
}

func indexRuns(runs []store.Run) auditInputs {
	inputs := auditInputs{
		runsByWorktree:  make(map[string]store.Run),
		runsByRunBranch: make(map[string]store.Run),
		remoteRefs:      make(map[string][]string),
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
	branch branchRef,
	inputs auditInputs,
) Survivor {
	if run, exists := branchRun(branch, inputs.pullRequests); exists {
		return pullRequestSurvivor(branch.Name, false, run)
	}
	if run, exists := branchRun(branch, inputs.active); exists {
		return activeRunSurvivor(branch.Name, false, run)
	}
	classification, contentOnly := classifyGitRef(
		ctx,
		runner,
		repoRoot,
		branch.Name,
		inputs.defaultRef,
		inputs.defaultName,
	)
	classification.Name = branch.Name
	if classification.Kind == KindResidue {
		if branch.Remote != "" {
			classification.Reclaim = "git push --delete " + shellQuote(branch.Remote) + " " + shellQuote(branch.RemoteBranch)
		} else {
			deleteFlag := "-d"
			if contentOnly {
				deleteFlag = "-D"
			}
			classification.Reclaim = "git branch " + deleteFlag + " -- " + shellQuote(branch.Name)
		}
	}
	return classification
}

func coordinateScratchReclaims(
	survivors []Survivor,
	localBranchSurvivors map[string]int,
	scratchWorktreeSurvivors map[string][]int,
) {
	for branch, worktreeIndexes := range scratchWorktreeSurvivors {
		branchIndex, branchExists := localBranchSurvivors[branch]
		if !branchExists || survivors[branchIndex].Kind != KindResidue {
			continue
		}
		if len(worktreeIndexes) != 1 {
			survivors[branchIndex].Kind = KindPreserved
			survivors[branchIndex].Evidence = fmt.Sprintf(
				"branch %q is checked out in %d Run-less worktrees; it is preserved",
				branch,
				len(worktreeIndexes),
			)
			survivors[branchIndex].Reclaim = ""
			for _, worktreeIndex := range worktreeIndexes {
				if survivors[worktreeIndex].Kind != KindResidue {
					continue
				}
				survivors[worktreeIndex].Kind = KindPreserved
				survivors[worktreeIndex].Evidence = fmt.Sprintf(
					"worktree %q shares branch %q with another Run-less worktree; it is preserved",
					survivors[worktreeIndex].Name,
					branch,
				)
				survivors[worktreeIndex].Reclaim = ""
			}
			continue
		}
		worktreeIndex := worktreeIndexes[0]
		if survivors[worktreeIndex].Kind != KindResidue {
			survivors[branchIndex].Kind = KindPreserved
			survivors[branchIndex].Evidence = fmt.Sprintf(
				"branch %q is checked out in preserved Run-less worktree %q; it is preserved",
				branch,
				survivors[worktreeIndex].Name,
			)
			survivors[branchIndex].Reclaim = ""
			continue
		}
		reclaim := "git worktree remove -- " + shellQuote(survivors[worktreeIndex].Name) +
			" && git branch -D -- " + shellQuote(branch)
		survivors[branchIndex].Reclaim = reclaim
		survivors[worktreeIndex].Reclaim = reclaim
	}
}

func branchAssociated(branch branchRef, associated map[string]bool) bool {
	return associated[branch.Name] || (branch.RemoteBranch != "" && associated[branch.RemoteBranch])
}

func branchRun(branch branchRef, runs map[string]store.Run) (store.Run, bool) {
	if run, exists := runs[branch.Name]; exists {
		return run, true
	}
	if branch.RemoteBranch != "" {
		run, exists := runs[branch.RemoteBranch]
		return run, exists
	}
	return store.Run{}, false
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
	if activeRun, exists := inputs.active[candidate.Branch]; exists {
		return activeRunSurvivor(candidate.Path, true, activeRun)
	}
	if !hasMatchingRun {
		return classifyScratchWorktree(ctx, runner, repoRoot, candidate, inputs)
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
	classification, _ := classifyGitRef(ctx, runner, repoRoot, ref, inputs.defaultRef, inputs.defaultName)
	classification.Name = candidate.Path
	classification.IsWorktree = true
	if classification.Kind == KindResidue {
		classification.Reclaim = "git worktree remove -- " + shellQuote(candidate.Path)
	}
	return classification
}

func classifyScratchWorktree(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	candidate worktree,
	inputs auditInputs,
) Survivor {
	preserved := func(evidence string) Survivor {
		return Survivor{
			Name:       candidate.Path,
			IsWorktree: true,
			Kind:       KindPreserved,
			Evidence:   evidence,
		}
	}
	if candidate.Branch == "" {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run and no local branch; it is preserved",
			candidate.Path,
		))
	}
	dirty, evidence, err := worktreeChanges(ctx, runner, candidate.Path)
	if err != nil {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run, could not be inspected, and is preserved: %v",
			candidate.Path,
			err,
		))
	}
	if dirty {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run and is preserved: %s",
			candidate.Path,
			evidence,
		))
	}
	remoteRef, err := pushedRemoteRef(
		ctx,
		runner,
		repoRoot,
		candidate.Branch,
		inputs.remoteRefs[candidate.Branch],
	)
	if err != nil {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run, branch push state could not be determined, and it is preserved: %v",
			candidate.Path,
			err,
		))
	}
	if remoteRef == "" {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run and branch %q is unpushed; it is preserved",
			candidate.Path,
			candidate.Branch,
		))
	}
	integration, _ := classifyGitRef(
		ctx,
		runner,
		repoRoot,
		candidate.Branch,
		inputs.defaultRef,
		inputs.defaultName,
	)
	if integration.Kind != KindResidue {
		return preserved(fmt.Sprintf(
			"worktree %q has no matching Run and branch %q is pushed to %q, but %s; it is preserved",
			candidate.Path,
			candidate.Branch,
			remoteRef,
			integration.Evidence,
		))
	}
	return Survivor{
		Name:       candidate.Path,
		IsWorktree: true,
		Kind:       KindResidue,
		Evidence: fmt.Sprintf(
			"worktree %q has no matching Run; branch %q is pushed to %q and %s",
			candidate.Path,
			candidate.Branch,
			remoteRef,
			integration.Evidence,
		),
		Reclaim: "git worktree remove -- " + shellQuote(candidate.Path) +
			" && git branch -D -- " + shellQuote(candidate.Branch),
	}
}

func pushedRemoteRef(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	branch string,
	remoteRefs []string,
) (string, error) {
	var inspectErr error
	for _, remoteRef := range remoteRefs {
		uniqueCommits, err := countUniqueCommits(ctx, runner, repoRoot, branch, remoteRef)
		if err != nil {
			if inspectErr == nil {
				inspectErr = fmt.Errorf("compare branch %q with remote ref %q: %w", branch, remoteRef, err)
			}
			continue
		}
		if uniqueCommits == 0 {
			return remoteRef, nil
		}
	}
	return "", inspectErr
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
) (Survivor, bool) {
	if defaultRef == "" {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: "default branch could not be resolved; survivor is preserved",
		}, false
	}
	uniqueCommits, err := countUniqueCommits(ctx, runner, repoRoot, ref, defaultRef)
	if err != nil {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: fmt.Sprintf("integration could not be inspected against default branch %q: %v", defaultName, err),
		}, false
	}
	if uniqueCommits == 0 {
		return Survivor{
			Kind: KindResidue,
			Evidence: fmt.Sprintf(
				"survivor is reachable from default branch %q",
				defaultName,
			),
		}, false
	}
	runOnly, differingShared, err := compareContent(ctx, runner, repoRoot, ref, defaultRef)
	if err != nil {
		return Survivor{
			Kind:     KindPreserved,
			Evidence: fmt.Sprintf("content could not be compared with default branch %q: %v", defaultName, err),
		}, false
	}
	if runOnly == 0 && differingShared == 0 {
		return Survivor{
			Kind: KindResidue,
			Evidence: fmt.Sprintf(
				"survivor content is fully represented on default branch %q",
				defaultName,
			),
		}, true
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
	}, false
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

func listBranches(ctx context.Context, runner gitRunner, repoRoot string) ([]branchRef, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"for-each-ref",
		"--format=%(refname)",
		"refs/heads",
		"refs/remotes",
	)
	if err != nil {
		return nil, err
	}
	branches := make([]branchRef, 0)
	for _, ref := range nonEmptyLines(output) {
		switch {
		case strings.HasPrefix(ref, "refs/heads/"):
			branches = append(branches, branchRef{Name: strings.TrimPrefix(ref, "refs/heads/")})
		case strings.HasPrefix(ref, "refs/remotes/"):
			name := strings.TrimPrefix(ref, "refs/remotes/")
			remote, remoteName, found := strings.Cut(name, "/")
			if !found || remoteName == "HEAD" {
				continue
			}
			branches = append(branches, branchRef{
				Name:         name,
				Remote:       remote,
				RemoteBranch: remoteName,
			})
		}
	}
	sort.Slice(branches, func(left, right int) bool {
		return branches[left].Name < branches[right].Name
	})
	return branches, nil
}

func defaultBranchAliases(
	ctx context.Context,
	runner gitRunner,
	repoRoot string,
	defaultRef string,
	defaultName string,
	branches []branchRef,
) map[string]bool {
	aliases := map[string]bool{defaultName: true}
	if !strings.HasSuffix(defaultRef, "/HEAD") {
		return aliases
	}
	target, err := runner.Run(ctx, repoRoot, "symbolic-ref", "--quiet", "--short", defaultRef)
	if err != nil {
		return aliases
	}
	target = strings.TrimSpace(target)
	aliases[target] = true
	for _, branch := range branches {
		if branch.Name == target && branch.RemoteBranch != "" {
			aliases[branch.RemoteBranch] = true
		}
	}
	return aliases
}

func listDeliveryBranches(ctx context.Context, runner gitRunner, repoRoot string) ([]string, error) {
	output, err := runner.Run(
		ctx,
		repoRoot,
		"for-each-ref",
		"--format=%(refname:short)",
		"refs/heads",
		"refs/remotes",
	)
	if err != nil {
		return nil, err
	}
	branches := nonEmptyLines(output)
	sort.Strings(branches)
	return branches, nil
}

func resolveSpecDeliveryTree(
	ctx context.Context,
	runner gitRunner,
	specsRoot string,
	codeTree deliveryTree,
) (deliveryTree, error) {
	rootOutput, err := runner.Run(ctx, specsRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		return deliveryTree{}, fmt.Errorf("resolve Spec Root Git repository: %w", err)
	}
	specRepoRoot := cleanPath(rootOutput)
	resolvedSpecsRoot := cleanPath(specsRoot)
	relative, err := filepath.Rel(specRepoRoot, resolvedSpecsRoot)
	if err != nil {
		return deliveryTree{}, fmt.Errorf("make Spec Root repository-relative: %w", err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return deliveryTree{}, fmt.Errorf("Spec Root %q is outside its Git repository %q", specsRoot, specRepoRoot)
	}
	artifactRoot := ""
	if relative != "." {
		artifactRoot = filepath.ToSlash(relative)
	}
	if specRepoRoot == cleanPath(codeTree.repoRoot) {
		codeTree.repoRoot = specRepoRoot
		codeTree.artifactRoot = artifactRoot
		return codeTree, nil
	}
	defaultRef, defaultName, err := defaultBranch(ctx, runner, specRepoRoot)
	if err != nil {
		return deliveryTree{}, fmt.Errorf("resolve Spec Root default branch: %w", err)
	}
	branches, err := listDeliveryBranches(ctx, runner, specRepoRoot)
	if err != nil {
		return deliveryTree{}, fmt.Errorf("enumerate Spec Root branches: %w", err)
	}
	return deliveryTree{
		repoRoot:     specRepoRoot,
		artifactRoot: artifactRoot,
		defaultRef:   defaultRef,
		defaultName:  defaultName,
		branches:     branches,
	}, nil
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
