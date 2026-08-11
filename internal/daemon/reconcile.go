package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const BranchDispositionRecordSchemaVersion = "roundfix-branch-disposition/v1"

// RunTaskCoverage records the completed Tasks proven by one Run's journal.
type RunTaskCoverage struct {
	Run            store.Run
	CompletedTasks []string
}

// BranchDisposition records why one terminal Run Branch can or cannot be
// discarded. Commits and ChangedFiles describe the work created after the
// Run's recorded starting head.
type BranchDisposition struct {
	RunID         string   `json:"runId"`
	Branch        string   `json:"branch"`
	TargetBranch  string   `json:"targetBranch"`
	RunHead       string   `json:"runHead"`
	TargetHead    string   `json:"targetHead"`
	Worktree      string   `json:"worktree"`
	Superseded    bool     `json:"superseded"`
	Reachable     bool     `json:"reachable"`
	Commits       []string `json:"commits"`
	ChangedFiles  []string `json:"changedFiles"`
	Reason        string   `json:"reason"`
	RefusalReason string   `json:"refusalReason"`

	evidence *branchDispositionEvidence
}

type branchDispositionEvidence struct {
	run       store.Run
	coverages []RunTaskCoverage
	snapshot  branchDispositionSnapshot
}

type branchDispositionSnapshot struct {
	branch        string
	targetBranch  string
	runHead       string
	targetHead    string
	worktree      string
	superseded    bool
	reachable     bool
	commits       string
	changedFiles  string
	reason        string
	refusalReason string
}

type branchDispositionRecord struct {
	SchemaVersion string    `json:"schemaVersion"`
	RecordedAt    time.Time `json:"recordedAt"`
	BranchDisposition
}

// ClassifySupersededBranch proves whether a terminal Run Branch is superseded.
// It first proves every Run commit reachable from the target. When that fails,
// a later Clean Run for the same target and Spec must have completed every Task
// covered by this Run.
func ClassifySupersededBranch(
	ctx context.Context,
	run store.Run,
	coverages []RunTaskCoverage,
) (BranchDisposition, error) {
	result := BranchDisposition{
		RunID:        strings.TrimSpace(run.ID),
		Branch:       runworktree.BranchName(run.ID),
		TargetBranch: strings.TrimSpace(run.LocalBranch),
		Worktree:     strings.TrimSpace(run.WorkDir),
		Commits:      []string{},
		ChangedFiles: []string{},
	}
	if run.Kind != store.KindImplement {
		return result, fmt.Errorf("classify superseded Run Branch %q: Run kind %q is not %q", run.ID, run.Kind, store.KindImplement)
	}
	if !store.IsTerminalState(run.State) {
		return result, fmt.Errorf("classify superseded Run Branch %q: Run is Active", run.ID)
	}

	inspected, err := runworktree.InspectTerminalRun(ctx, run)
	if err != nil {
		return result, fmt.Errorf("classify superseded Run Branch %q: inspect terminal Run: %w", run.ID, err)
	}
	result.Branch = inspected.Branch
	result.TargetBranch = inspected.TargetBranch
	result.RunHead = inspected.RunHead
	result.TargetHead = inspected.TargetHead
	result.Worktree = inspected.Path
	if inspected.State == runworktree.ReconciliationReleased {
		result.RefusalReason = "branch condition failed: Run Worktree and Run Branch are already absent"
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	if inspected.State == runworktree.ReconciliationDirty || inspected.State == runworktree.ReconciliationUnknown {
		result.RefusalReason = fmt.Sprintf("worktree condition failed: %s", inspected.Reason)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	if result.RunHead == "" {
		result.RefusalReason = "branch condition failed: Run Branch head is unavailable"
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	if result.TargetHead == "" {
		result.RefusalReason = "target condition failed: target branch head is unavailable"
		return withBranchDispositionEvidence(result, run, coverages), nil
	}

	commits, err := runBranchCommits(ctx, run.GitRoot, run.HeadSHA, result.RunHead)
	if err != nil {
		result.RefusalReason = fmt.Sprintf("commit inventory condition failed: %v", err)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	result.Commits = commits
	changedFiles, err := runBranchChangedFiles(ctx, run.GitRoot, commits)
	if err != nil {
		result.RefusalReason = fmt.Sprintf("changed-file inventory condition failed: %v", err)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	result.ChangedFiles = changedFiles

	unreachable, err := unreachableCommits(ctx, run.GitRoot, commits, result.TargetHead)
	if err != nil {
		result.RefusalReason = fmt.Sprintf("reachability condition failed: %v", err)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	if len(unreachable) == 0 {
		result.Superseded = true
		result.Reachable = true
		result.Reason = fmt.Sprintf("every Run Branch commit is reachable from target branch %q", result.TargetBranch)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}

	coveredTasks, laterRun, covered := laterIntegratedRunCoverage(run, coverages)
	if covered {
		result.Superseded = true
		result.Reason = fmt.Sprintf(
			"later integrated Run %q covered Tasks [%s] for target branch %q",
			laterRun.ID,
			strings.Join(coveredTasks, ", "),
			result.TargetBranch,
		)
		return withBranchDispositionEvidence(result, run, coverages), nil
	}
	result.RefusalReason = fmt.Sprintf(
		"reachability condition failed: unreachable commit %s; later-Run Task coverage condition failed",
		unreachable[0],
	)
	return withBranchDispositionEvidence(result, run, coverages), nil
}

func withBranchDispositionEvidence(
	result BranchDisposition,
	run store.Run,
	coverages []RunTaskCoverage,
) BranchDisposition {
	result.Commits = slices.Clone(result.Commits)
	result.ChangedFiles = slices.Clone(result.ChangedFiles)
	result.evidence = &branchDispositionEvidence{
		run:       run,
		coverages: cloneRunTaskCoverages(coverages),
		snapshot:  snapshotBranchDisposition(result),
	}
	return result
}

func cloneRunTaskCoverages(coverages []RunTaskCoverage) []RunTaskCoverage {
	cloned := append([]RunTaskCoverage(nil), coverages...)
	for index := range cloned {
		cloned[index].CompletedTasks = slices.Clone(cloned[index].CompletedTasks)
	}
	return cloned
}

func snapshotBranchDisposition(result BranchDisposition) branchDispositionSnapshot {
	return branchDispositionSnapshot{
		branch:        result.Branch,
		targetBranch:  result.TargetBranch,
		runHead:       result.RunHead,
		targetHead:    result.TargetHead,
		worktree:      result.Worktree,
		superseded:    result.Superseded,
		reachable:     result.Reachable,
		commits:       strings.Join(result.Commits, "\x00"),
		changedFiles:  strings.Join(result.ChangedFiles, "\x00"),
		reason:        result.Reason,
		refusalReason: result.RefusalReason,
	}
}

func runBranchCommits(ctx context.Context, gitRoot, baseHead, runHead string) ([]string, error) {
	baseHead = strings.TrimSpace(baseHead)
	if baseHead == "" {
		return nil, errors.New("recorded starting head is missing")
	}
	ancestor, err := gitAncestor(ctx, gitRoot, baseHead, runHead)
	if err != nil {
		return nil, fmt.Errorf("inspect recorded starting head %q: %w", baseHead, err)
	}
	if !ancestor {
		return nil, fmt.Errorf("recorded starting head %q is not an ancestor of Run Branch head %q", baseHead, runHead)
	}
	output, err := runDispositionGit(ctx, gitRoot, "rev-list", "--reverse", baseHead+".."+runHead)
	if err != nil {
		return nil, fmt.Errorf("list Run Branch commits: %w", err)
	}
	commits := strings.Fields(output)
	return commits, nil
}

func runBranchChangedFiles(ctx context.Context, gitRoot string, commits []string) ([]string, error) {
	files := make(map[string]struct{})
	for _, commit := range commits {
		output, err := runDispositionGitBytes(
			ctx,
			gitRoot,
			"diff-tree",
			"--root",
			"-m",
			"--first-parent",
			"--no-commit-id",
			"--name-only",
			"-r",
			"-z",
			commit,
		)
		if err != nil {
			return nil, fmt.Errorf("list changed files for commit %s: %w", commit, err)
		}
		for _, raw := range strings.Split(string(output), "\x00") {
			path := strings.TrimSpace(raw)
			if path != "" {
				files[path] = struct{}{}
			}
		}
	}
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths, nil
}

func unreachableCommits(ctx context.Context, gitRoot string, commits []string, targetHead string) ([]string, error) {
	unreachable := make([]string, 0)
	for _, commit := range commits {
		reachable, err := gitAncestor(ctx, gitRoot, commit, targetHead)
		if err != nil {
			return nil, fmt.Errorf("inspect commit %s against target %s: %w", commit, targetHead, err)
		}
		if !reachable {
			unreachable = append(unreachable, commit)
		}
	}
	return unreachable, nil
}

func laterIntegratedRunCoverage(
	run store.Run,
	coverages []RunTaskCoverage,
) ([]string, store.Run, bool) {
	currentTasks := tasksForRun(run, coverages)
	if len(currentTasks) == 0 {
		return nil, store.Run{}, false
	}
	candidates := append([]RunTaskCoverage(nil), coverages...)
	sort.Slice(candidates, func(left, right int) bool {
		if candidates[left].Run.CreatedAt.Equal(candidates[right].Run.CreatedAt) {
			return candidates[left].Run.ID < candidates[right].Run.ID
		}
		return candidates[left].Run.CreatedAt.Before(candidates[right].Run.CreatedAt)
	})
	for _, candidate := range candidates {
		later := candidate.Run
		if later.ID == run.ID || later.State != store.StateClean || !later.CreatedAt.After(run.CreatedAt) ||
			!sameDispositionTarget(run, later) {
			continue
		}
		completed := taskSet(candidate.CompletedTasks)
		covered := true
		for _, task := range currentTasks {
			if !completed[task] {
				covered = false
				break
			}
		}
		if covered {
			return currentTasks, later, true
		}
	}
	return currentTasks, store.Run{}, false
}

func tasksForRun(run store.Run, coverages []RunTaskCoverage) []string {
	for _, coverage := range coverages {
		if coverage.Run.ID == run.ID {
			return sortedTasks(coverage.CompletedTasks)
		}
	}
	return nil
}

func sortedTasks(tasks []string) []string {
	set := taskSet(tasks)
	result := make([]string, 0, len(set))
	for task := range set {
		result = append(result, task)
	}
	sort.Strings(result)
	return result
}

func taskSet(tasks []string) map[string]bool {
	set := make(map[string]bool)
	for _, task := range tasks {
		if task = strings.TrimSpace(task); task != "" {
			set[task] = true
		}
	}
	return set
}

func sameDispositionTarget(left, right store.Run) bool {
	return strings.TrimSpace(left.GitRoot) == strings.TrimSpace(right.GitRoot) &&
		strings.TrimSpace(left.LocalBranch) == strings.TrimSpace(right.LocalBranch) &&
		strings.TrimSpace(left.SpecSlug) == strings.TrimSpace(right.SpecSlug)
}

// WriteBranchDispositionRecord persists the complete proof before Git state is
// changed. The record path must sit outside the Run Worktree it documents.
func WriteBranchDispositionRecord(path string, disposition BranchDisposition, recordedAt time.Time) error {
	if !disposition.Superseded || disposition.evidence == nil {
		return errors.New("write branch disposition record: branch was not classified superseded")
	}
	path = filepath.Clean(strings.TrimSpace(path))
	if path == "." || !filepath.IsAbs(path) {
		return errors.New("write branch disposition record: absolute record path is required")
	}
	if pathInside(path, disposition.Worktree) {
		return fmt.Errorf("write branch disposition record: path %q is inside Run Worktree %q", path, disposition.Worktree)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("write branch disposition record directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".branch-disposition-*.tmp")
	if err != nil {
		return fmt.Errorf("create branch disposition record: %w", err)
	}
	temporaryPath := temporary.Name()
	removeTemporary := true
	defer func() {
		if removeTemporary {
			_ = os.Remove(temporaryPath)
		}
	}()
	encoder := json.NewEncoder(temporary)
	encoder.SetIndent("", "  ")
	record := branchDispositionRecord{
		SchemaVersion:     BranchDispositionRecordSchemaVersion,
		RecordedAt:        recordedAt.UTC(),
		BranchDisposition: disposition,
	}
	if err := encoder.Encode(record); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("encode branch disposition record: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("sync branch disposition record: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close branch disposition record: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("publish branch disposition record: %w", err)
	}
	removeTemporary = false
	return nil
}

func pathInside(path, directory string) bool {
	directory = filepath.Clean(strings.TrimSpace(directory))
	if directory == "." || !filepath.IsAbs(directory) {
		return false
	}
	relative, err := filepath.Rel(directory, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// RecordAndDiscardSupersededBranch makes the durable record the mandatory
// first step of the destructive operation.
func RecordAndDiscardSupersededBranch(
	ctx context.Context,
	path string,
	disposition BranchDisposition,
	recordedAt time.Time,
) error {
	if err := WriteBranchDispositionRecord(path, disposition, recordedAt); err != nil {
		return err
	}
	return discardSupersededBranch(ctx, disposition)
}

func discardSupersededBranch(ctx context.Context, disposition BranchDisposition) error {
	evidence := disposition.evidence
	if evidence == nil || snapshotBranchDisposition(disposition) != evidence.snapshot {
		return errors.New("discard superseded Run Branch: disposition was not produced by inspection or was changed")
	}
	if !disposition.Superseded {
		return fmt.Errorf("discard superseded Run Branch %q: %s", disposition.Branch, disposition.RefusalReason)
	}
	fresh, err := ClassifySupersededBranch(ctx, evidence.run, evidence.coverages)
	if err != nil {
		return fmt.Errorf("revalidate superseded Run Branch %q: %w", disposition.Branch, err)
	}
	if snapshotBranchDisposition(fresh) != evidence.snapshot {
		return fmt.Errorf("discard superseded Run Branch %q: disposition proof is stale", disposition.Branch)
	}
	if _, err := os.Stat(fresh.Worktree); err == nil {
		if _, err := runDispositionGit(ctx, evidence.run.GitRoot, "worktree", "remove", fresh.Worktree); err != nil {
			return fmt.Errorf("discard superseded Run Branch %q: remove Run Worktree %q: %w", fresh.Branch, fresh.Worktree, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("discard superseded Run Branch %q: inspect Run Worktree %q: %w", fresh.Branch, fresh.Worktree, err)
	}
	if _, err := runDispositionGit(ctx, evidence.run.GitRoot, "branch", "-D", fresh.Branch); err != nil {
		return fmt.Errorf("discard superseded Run Branch %q: delete branch: %w", fresh.Branch, err)
	}
	return nil
}

func gitAncestor(ctx context.Context, gitRoot, ancestor, descendant string) (bool, error) {
	command := exec.CommandContext(ctx, "git", "-C", gitRoot, "merge-base", "--is-ancestor", ancestor, descendant)
	output, err := command.CombinedOutput()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("git merge-base --is-ancestor: %w: %s", err, strings.TrimSpace(string(output)))
}

func runDispositionGit(ctx context.Context, gitRoot string, args ...string) (string, error) {
	output, err := runDispositionGitBytes(ctx, gitRoot, args...)
	return strings.TrimSpace(string(output)), err
}

func runDispositionGitBytes(ctx context.Context, gitRoot string, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-C", gitRoot}, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}
