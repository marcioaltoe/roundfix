package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const carryForwardReadyAction = "would carry forward with --carry-forward"

var carryForwardAcceptedStates = []string{store.StateStopped, store.StateUnresolved}

// specCarryForward is what one prior Run would hand back to the checkout.
type specCarryForward struct {
	Run        store.Run
	Candidates []spec.CarryForward
}

// carriable counts the candidates that passed every proof.
func (result specCarryForward) carriable() int {
	count := 0
	for _, candidate := range result.Candidates {
		if candidate.Action == carryForwardReadyAction {
			count++
		}
	}
	return count
}

// wouldCarry reports whether carry-forward would accept the complete set.
func (result specCarryForward) wouldCarry() bool {
	return len(result.Candidates) > 0 && result.carriable() == len(result.Candidates)
}

// inspectSpecCarryForwards reports, per prior terminal Run of one Spec in this
// repository, what carry-forward would do. Runs whose Run Worktree is no
// longer present are skipped rather than failing the inspection.
func inspectSpecCarryForwards(
	ctx context.Context,
	runStore *store.Store,
	repository string,
	resolvedSpecsRoot roundconfig.SpecsRoot,
	specSlug string,
) ([]specCarryForward, error) {
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{
		GitRoot: repository,
		States:  store.StatesTerminal,
	})
	if err != nil {
		return nil, fmt.Errorf("list prior Runs for Spec %q carry-forward: %w", specSlug, err)
	}
	selected := make([]store.Run, 0, len(runs))
	for _, run := range runs {
		if run.Kind != store.KindImplement ||
			strings.TrimSpace(run.SpecSlug) != strings.TrimSpace(specSlug) ||
			!slices.Contains(carryForwardAcceptedStates, run.State) {
			continue
		}
		present, err := carryForwardRunWorktreePresent(run)
		if err != nil {
			return nil, err
		}
		if !present {
			continue
		}
		selected = append(selected, run)
	}
	_, taskEvidence, err := loadReconcileTaskCoverage(ctx, runStore, selected)
	if err != nil {
		return nil, err
	}
	results := make([]specCarryForward, 0, len(selected))
	for _, run := range selected {
		candidates, err := inspectCarryForwards(
			ctx,
			repository,
			resolvedSpecsRoot,
			reconcileRunSelection{
				selected:     []store.Run{run},
				taskEvidence: taskEvidence,
			},
		)
		if err != nil {
			return nil, err
		}
		results = append(results, specCarryForward{Run: run, Candidates: candidates})
	}
	return results, nil
}

func carryForwardRunWorktreePresent(run store.Run) (bool, error) {
	if strings.TrimSpace(run.WorkDir) == "" {
		return false, nil
	}
	if _, err := os.Stat(run.WorkDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("inspect carry-forward Run %q Worktree %q: %w", run.ID, run.WorkDir, err)
	}
	return true, nil
}

func inspectCarryForwards(
	ctx context.Context,
	repository string,
	resolvedSpecsRoot roundconfig.SpecsRoot,
	runs reconcileRunSelection,
) ([]spec.CarryForward, error) {
	repoSpecsRoot, ok := repositoryRelativePath(repository, resolvedSpecsRoot.Path)
	if !ok {
		return nil, fmt.Errorf("carry-forward Specs Root %q is outside repository %q", resolvedSpecsRoot.Path, repository)
	}
	repoSpecsRoot = filepath.ToSlash(repoSpecsRoot)
	carried := make([]spec.CarryForward, 0)
	for _, run := range runs.selected {
		if !slices.Contains(carryForwardAcceptedStates, run.State) {
			continue
		}
		if strings.TrimSpace(run.WorkDir) == "" {
			return nil, fmt.Errorf("inspect carry-forward Run %q: recorded Run Worktree is missing", run.ID)
		}
		sourceSpecsRoot := specsRootForWorkDir(resolvedSpecsRoot, repository, run.WorkDir)
		graph, err := spec.Load(sourceSpecsRoot, run.SpecSlug)
		if err != nil {
			return nil, fmt.Errorf("load Run %q Spec %q: %w", run.ID, run.SpecSlug, err)
		}
		commitsByTask, err := carryForwardTaskCommits(ctx, run)
		if err != nil {
			return nil, err
		}
		evidence := runs.taskEvidence[run.ID]
		checkoutHead, err := reconcileGitText(ctx, repository, "rev-parse", "HEAD")
		if err != nil {
			return nil, fmt.Errorf("read checkout HEAD for Run %q carry-forward proof: %w", run.ID, err)
		}
		stagingWorktree, cleanup, err := createCarryForwardStaging(ctx, repository, checkoutHead)
		if err != nil {
			return nil, err
		}
		runCarried, inspectErr := inspectCarryForwardsForRun(
			ctx,
			stagingWorktree,
			repoSpecsRoot,
			run,
			graph,
			commitsByTask,
			evidence,
		)
		cleanupErr := cleanup()
		if inspectErr != nil || cleanupErr != nil {
			return nil, errors.Join(inspectErr, cleanupErr)
		}
		carried = append(carried, runCarried...)
	}
	return carried, nil
}

func inspectCarryForwardsForRun(
	ctx context.Context,
	stagingWorktree string,
	repoSpecsRoot string,
	run store.Run,
	graph *spec.Graph,
	commitsByTask map[string][]string,
	evidence map[string]reconcileTaskEvidence,
) ([]spec.CarryForward, error) {
	carried := make([]spec.CarryForward, 0)
	for index, task := range graph.Tasks {
		taskEvidence := evidence[task.ID]
		if !taskEvidence.settledCompleted {
			continue
		}
		candidate := spec.CarryForward{
			TaskID:      task.ID,
			RunID:       run.ID,
			TaskFile:    filepath.ToSlash(filepath.Join(repoSpecsRoot, task.File)),
			MovedInputs: []string{},
		}
		if !taskEvidence.verificationPassed {
			candidate.Action = "refuse"
			candidate.RefusalReason = fmt.Sprintf("Task %s has no passing Verification verdict in Run %s", task.ID, run.ID)
			carried = append(carried, candidate)
			continue
		}
		taskCommits := commitsByTask[task.ID]
		if len(taskCommits) != 1 {
			candidate.Action = "refuse"
			candidate.RefusalReason = fmt.Sprintf("Task %s has %d settlement commits in Run %s; expected exactly one", task.ID, len(taskCommits), run.ID)
			carried = append(carried, candidate)
			continue
		}
		candidate.Commit = taskCommits[0]
		if err := inspectCarryForwardCandidate(ctx, stagingWorktree, &candidate); err != nil {
			candidate.Action = "refuse"
			candidate.RefusalReason = err.Error()
			carried = append(carried, candidate)
			continue
		}
		if candidate.InputsMoved {
			candidate.Action = "refuse"
			candidate.RefusalReason = fmt.Sprintf(
				"Task %s declared input(s) moved: %s",
				task.ID,
				strings.Join(candidate.MovedInputs, ", "),
			)
		} else {
			candidate.Action = carryForwardReadyAction
			if stageErr := stageCarryForwardCandidate(ctx, stagingWorktree, candidate); stageErr != nil {
				// Only a commit that will not apply is this Task's own fault,
				// and only that is reported as its refusal. Provenance,
				// staging, and amend failures are Roundfix failing at its own
				// bookkeeping on work that applied cleanly; calling those a
				// refusal would tell the caller the Task is unfit when the
				// tool is.
				var conflict carryForwardConflictError
				if !errors.As(stageErr, &conflict) {
					return nil, fmt.Errorf("stage Task %s while proving serial carry-forward: %w", task.ID, stageErr)
				}
				candidate.Action = "refuse"
				candidate.RefusalReason = stageErr.Error()
				if abortErr := abortCarryForwardStaging(ctx, stagingWorktree); abortErr != nil {
					// Surfaced in the reason rather than returned: a failed
					// abort must not outrank the conflict that is the finding.
					candidate.RefusalReason += fmt.Sprintf("; aborting the staged cherry-pick also failed: %v", abortErr)
				}
				carried = append(carried, candidate)
				carried = append(carried, unstagedCarryForwards(repoSpecsRoot, run, graph.Tasks[index+1:], evidence, task.ID)...)
				return carried, nil
			}
		}
		carried = append(carried, candidate)
	}
	return carried, nil
}

// carryForwardConflictError marks a settlement commit that will not apply to
// the staged state. That is a property of the Task's own work, so it is
// reported as that Task's refusal. Every other staging failure is Roundfix
// failing at its own bookkeeping over work that applied cleanly, and stays an
// operational error.
type carryForwardConflictError struct {
	TaskID string
	Commit string
	Err    error
}

func (err carryForwardConflictError) Error() string {
	return fmt.Sprintf(
		"Task %s settlement commit %s does not apply to the staged carry-forward state: %v",
		err.TaskID,
		err.Commit,
		err.Err,
	)
}

func (err carryForwardConflictError) Unwrap() error { return err.Err }

// abortCarryForwardStaging clears the in-progress cherry-pick so nothing is
// left half-applied in the staging Worktree. Its failure is returned rather
// than swallowed, but the caller records it beside the conflict instead of
// returning it, so a failed abort never outranks the finding.
func abortCarryForwardStaging(ctx context.Context, stagingWorktree string) error {
	if _, err := reconcileGitRaw(ctx, stagingWorktree, "cherry-pick", "--abort"); err != nil {
		return fmt.Errorf("abort carry-forward staging cherry-pick: %w", err)
	}
	return nil
}

// unstagedCarryForwards records the Tasks whose proofs were never evaluated
// because an earlier Task in the same set could not be staged. Declared inputs
// are proved against the accumulating staged carries, so a set that stops
// staging has no baseline left to judge the rest against. Saying that is
// honest; inventing a verdict from a baseline that was never built is not.
func unstagedCarryForwards(
	repoSpecsRoot string,
	run store.Run,
	remaining []spec.Task,
	evidence map[string]reconcileTaskEvidence,
	blockingTaskID string,
) []spec.CarryForward {
	entries := make([]spec.CarryForward, 0, len(remaining))
	for _, task := range remaining {
		if !evidence[task.ID].settledCompleted {
			continue
		}
		entries = append(entries, spec.CarryForward{
			TaskID:      task.ID,
			RunID:       run.ID,
			TaskFile:    filepath.ToSlash(filepath.Join(repoSpecsRoot, task.File)),
			MovedInputs: []string{},
			Action:      "refuse",
			RefusalReason: fmt.Sprintf(
				"Task %s was not evaluated: Task %s could not be staged, so the set has no baseline to prove it against",
				task.ID,
				blockingTaskID,
			),
		})
	}
	return entries
}

func createCarryForwardStaging(
	ctx context.Context,
	repository string,
	checkoutHead string,
) (string, func() error, error) {
	tempRoot, err := os.MkdirTemp("", "roundfix-carry-forward-")
	if err != nil {
		return "", nil, fmt.Errorf("create carry-forward staging directory: %w", err)
	}
	stagingWorktree := filepath.Join(tempRoot, "worktree")
	if _, err := reconcileGitRaw(ctx, repository, "worktree", "add", "--detach", stagingWorktree, checkoutHead); err != nil {
		return "", nil, errors.Join(
			fmt.Errorf("create carry-forward staging Worktree: %w", err),
			os.RemoveAll(tempRoot),
		)
	}
	cleanup := func() error {
		_, worktreeErr := reconcileGitRaw(context.Background(), repository, "worktree", "remove", "--force", stagingWorktree)
		if worktreeErr != nil {
			worktreeErr = fmt.Errorf("remove carry-forward staging Worktree: %w", worktreeErr)
		}
		removeErr := os.RemoveAll(tempRoot)
		if removeErr != nil {
			removeErr = fmt.Errorf("remove carry-forward staging directory: %w", removeErr)
		}
		return errors.Join(worktreeErr, removeErr)
	}
	checkoutPatch, err := reconcileGitRaw(ctx, repository, "diff", "--binary", "--full-index", "HEAD", "--")
	if err != nil {
		return "", nil, errors.Join(fmt.Errorf("read checkout changes for carry-forward staging: %w", err), cleanup())
	}
	if len(checkoutPatch) > 0 {
		if _, err := reconcileGitRawInput(ctx, stagingWorktree, checkoutPatch, "apply", "--binary", "--whitespace=nowarn", "-"); err != nil {
			return "", nil, errors.Join(fmt.Errorf("apply checkout changes to carry-forward staging Worktree: %w", err), cleanup())
		}
	}
	return stagingWorktree, cleanup, nil
}

func stageCarryForwardCandidate(ctx context.Context, stagingWorktree string, candidate spec.CarryForward) error {
	if _, err := reconcileGitRaw(ctx, stagingWorktree, "cherry-pick", candidate.Commit); err != nil {
		return carryForwardConflictError{TaskID: candidate.TaskID, Commit: candidate.Commit, Err: err}
	}
	taskPath := filepath.Join(stagingWorktree, filepath.FromSlash(candidate.TaskFile))
	if err := spec.RecordCarryForward(taskPath, candidate.RunID, candidate.Commit); err != nil {
		return fmt.Errorf("record Task %s carry-forward provenance: %w", candidate.TaskID, err)
	}
	if _, err := reconcileGitRaw(ctx, stagingWorktree, "add", "--", candidate.TaskFile); err != nil {
		return fmt.Errorf("stage Task %s carry-forward provenance: %w", candidate.TaskID, err)
	}
	if _, err := reconcileGitRaw(ctx, stagingWorktree, "commit", "--amend", "--no-edit"); err != nil {
		return fmt.Errorf("amend Task %s carry-forward provenance: %w", candidate.TaskID, err)
	}
	return nil
}

func carryForwardTaskCommits(ctx context.Context, run store.Run) (map[string][]string, error) {
	if strings.TrimSpace(run.HeadSHA) == "" {
		return nil, fmt.Errorf("Run %q has no recorded starting commit", run.ID)
	}
	output, err := reconcileGitText(ctx, run.WorkDir, "rev-list", "--reverse", run.HeadSHA+"..HEAD")
	if err != nil {
		return nil, fmt.Errorf("list settlement commits for Run %q: %w", run.ID, err)
	}
	commits := make(map[string][]string)
	for _, commit := range strings.Fields(output) {
		message, err := reconcileGitText(ctx, run.WorkDir, "show", "-s", "--format=%B", commit)
		if err != nil {
			return nil, fmt.Errorf("read Run %q commit %s: %w", run.ID, commit, err)
		}
		if gitTrailerValue(message, "Roundfix-Spec") != run.SpecSlug {
			continue
		}
		taskID := gitTrailerValue(message, "Roundfix-Task")
		if taskID == "" {
			continue
		}
		commits[taskID] = append(commits[taskID], commit)
	}
	return commits, nil
}

func gitTrailerValue(message string, key string) string {
	prefix := strings.ToLower(strings.TrimSpace(key)) + ":"
	lines := strings.Split(message, "\n")
	for index := len(lines) - 1; index >= 0; index-- {
		line := strings.TrimSpace(lines[index])
		if strings.HasPrefix(strings.ToLower(line), prefix) {
			return strings.TrimSpace(line[len(prefix):])
		}
	}
	return ""
}

func inspectCarryForwardCandidate(ctx context.Context, stagingWorktree string, candidate *spec.CarryForward) error {
	if candidate == nil {
		return errors.New("carry-forward candidate is required")
	}
	changed, err := reconcileGitText(ctx, stagingWorktree, "diff-tree", "--no-commit-id", "--name-only", "-r", candidate.Commit)
	if err != nil {
		return fmt.Errorf("inspect Task %s settlement commit %s: %w", candidate.TaskID, candidate.Commit, err)
	}
	if !lineSet(changed)[candidate.TaskFile] {
		return fmt.Errorf("Task %s settlement commit %s does not contain %s", candidate.TaskID, candidate.Commit, candidate.TaskFile)
	}
	committedTask, found, err := reconcileGitBlob(ctx, stagingWorktree, candidate.Commit, candidate.TaskFile)
	if err != nil {
		return fmt.Errorf("read Task %s from settlement commit %s: %w", candidate.TaskID, candidate.Commit, err)
	}
	if !found {
		return fmt.Errorf("Task %s settlement commit %s does not contain %s", candidate.TaskID, candidate.Commit, candidate.TaskFile)
	}
	status, err := spec.CarryForwardStatus(candidate.TaskFile, committedTask)
	if err != nil {
		return fmt.Errorf("read Task %s status from settlement commit %s: %w", candidate.TaskID, candidate.Commit, err)
	}
	if status != spec.StatusCompleted {
		return fmt.Errorf("Task %s settlement commit %s records status %q, not completed", candidate.TaskID, candidate.Commit, status)
	}
	parent, err := reconcileGitText(ctx, stagingWorktree, "rev-parse", candidate.Commit+"^")
	if err != nil {
		return fmt.Errorf("resolve Task %s settlement parent: %w", candidate.TaskID, err)
	}
	parentTask, found, err := reconcileGitBlob(ctx, stagingWorktree, parent, candidate.TaskFile)
	if err != nil {
		return fmt.Errorf("read Task %s declared inputs from %s: %w", candidate.TaskID, parent, err)
	}
	if !found {
		return fmt.Errorf("Task %s input file %s is absent from settlement parent %s", candidate.TaskID, candidate.TaskFile, parent)
	}
	inputs, err := spec.CarryForwardInputs(filepath.Dir(candidate.TaskFile), candidate.TaskFile, parentTask)
	if err != nil {
		return fmt.Errorf("read Task %s declared inputs: %w", candidate.TaskID, err)
	}
	moved := make([]string, 0)
	for _, input := range inputs {
		_, existed, err := reconcileGitBlob(ctx, stagingWorktree, parent, input)
		if err != nil {
			return fmt.Errorf("read Task %s input %s from %s: %w", candidate.TaskID, input, parent, err)
		}
		if !existed {
			continue
		}
		if carryForwardPathCrossesSymlink(stagingWorktree, input) {
			moved = append(moved, input)
			continue
		}
		changed, err := reconcileGitInputChanged(ctx, stagingWorktree, parent, input)
		if err != nil {
			return fmt.Errorf("compare Task %s input %s against %s: %w", candidate.TaskID, input, parent, err)
		}
		if changed {
			moved = append(moved, input)
		}
	}
	slices.Sort(moved)
	candidate.MovedInputs = moved
	candidate.InputsMoved = len(moved) > 0
	return nil
}

func lineSet(output string) map[string]bool {
	set := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		line = filepath.ToSlash(strings.TrimSpace(line))
		if line != "" {
			set[line] = true
		}
	}
	return set
}

func carryForwardPathCrossesSymlink(repository string, relative string) bool {
	current := filepath.Clean(repository)
	for _, part := range strings.Split(filepath.FromSlash(relative), string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return false
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func carryForwardRefusalReason(candidates []spec.CarryForward) string {
	if len(candidates) == 0 {
		return "Run has no Tasks with completed settlement evidence"
	}
	refusals := make([]string, 0)
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.RefusalReason) != "" {
			refusals = append(refusals, candidate.RefusalReason)
		}
	}
	if len(refusals) == 0 {
		return ""
	}
	return "carry-forward refused the whole set: " + strings.Join(refusals, "; ")
}

func reconcileGitBlob(ctx context.Context, workDir string, revision string, path string) ([]byte, bool, error) {
	path = filepath.ToSlash(path)
	listing, err := reconcileGitRaw(ctx, workDir, "ls-tree", "-z", revision, "--", path)
	if err != nil {
		return nil, false, err
	}
	if len(listing) == 0 {
		return nil, false, nil
	}
	content, err := reconcileGitRaw(ctx, workDir, "show", revision+":"+path)
	if err != nil {
		return nil, false, err
	}
	return content, true, nil
}

// reconcileGitInputChanged reports whether the working-tree bytes for path
// differ from the bytes recorded at revision, applying Git's conversion rules
// (core.autocrlf, text/eol, and clean/smudge filters) rather than comparing raw
// bytes. A difference is reported as changed; nonexistent inputs and
// symlink-crossing paths are handled by the caller. Exit 1 from git diff means
// the path differs; any other Git error is returned.
func reconcileGitInputChanged(ctx context.Context, workDir string, revision string, path string) (bool, error) {
	output, err := reconcileGitRawDiff(ctx, workDir, "diff", "--quiet", revision, "--", filepath.ToSlash(path))
	if err != nil {
		return false, err
	}
	return output, nil
}

// reconcileGitRawDiff runs a git subcommand and maps a nonzero exit to a diff
// result without failing the whole helper: exit 1 signals a difference, other
// errors carry the diagnostic.
func reconcileGitRawDiff(ctx context.Context, workDir string, args ...string) (bool, error) {
	commandArgs := append([]string{"-c", "core.fsmonitor=false", "-c", "commit.gpgSign=false"}, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	command.Dir = workDir
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err == nil {
		return false, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return true, nil
	}
	diagnostic := strings.TrimSpace(string(output))
	if diagnostic == "" {
		return false, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return false, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, diagnostic)
}
