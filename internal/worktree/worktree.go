package worktree

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	runBranchPrefix = "roundfix/run-"

	ModeFastForwardMerge = "ff-merge"
	ModeBranchMove       = "branch-move"
	ModePending          = "pending"
	ModeTaskFastForward  = "ff"
	ModeTaskCherryPick   = "cherry-pick"
	ModeTaskConflict     = "conflict"

	ReasonOverlap             = "overlap"
	ReasonDiverged            = "diverged"
	ReasonNonAncestry         = "non-ancestry"
	ReasonCheckedOutElsewhere = "checked-out-elsewhere"
	ReasonDirtyTarget         = "dirty-target"
)

type Ref struct {
	RunID    string
	Path     string
	Branch   string
	UserRoot string
}

type CreateOptions struct {
	UserRoot        string
	Location        string
	RunID           string
	HeadSHA         string
	CopyList        []string
	Bootstrap       BootstrapSpec
	BootstrapOutput io.Writer
}

type TaskCreateOptions struct {
	CopyList        []string
	Bootstrap       BootstrapSpec
	BootstrapOutput io.Writer
}

type BootstrapSpec struct {
	Command string
	Timeout time.Duration
}

type BootstrapError struct {
	Command string
	Err     error
	Tail    string
}

func (err *BootstrapError) Error() string {
	reason := "unknown error"
	if err != nil && err.Err != nil {
		reason = err.Err.Error()
	}
	return fmt.Sprintf("worktree bootstrap failed: %s: %s", err.Command, reason)
}

func (err *BootstrapError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

type IntegrationResult struct {
	Mode   string
	Reason string
}

type TaskRef struct {
	RunID    string
	TaskID   string
	Path     string
	Branch   string
	UserRoot string
	BaseSHA  string
}

type TaskIntegration struct {
	Mode   string
	Reason string
}

type PrunedRef struct {
	RunID  string
	TaskID string
	Path   string
	Branch string
}

func Create(ctx context.Context, opts CreateOptions) (Ref, error) {
	ref, err := newRef(opts.UserRoot, opts.Location, opts.RunID)
	if err != nil {
		return Ref{}, err
	}
	headSHA := strings.TrimSpace(opts.HeadSHA)
	if headSHA == "" {
		return Ref{}, errors.New("create Run Worktree: HEAD is required")
	}

	if err := os.MkdirAll(filepath.Dir(ref.Path), 0o755); err != nil {
		return Ref{}, fmt.Errorf("create Run Worktree parent %q: %w", filepath.Dir(ref.Path), err)
	}

	runner := execGitRunner{}
	if _, err := runner.Run(ctx, ref.UserRoot, "worktree", "add", "-b", ref.Branch, ref.Path, headSHA); err != nil {
		return Ref{}, fmt.Errorf("create Run Worktree: %w", err)
	}
	if err := copyProvisionedFiles(ref.UserRoot, ref.Path, opts.CopyList); err != nil {
		return Ref{}, err
	}
	if err := runBootstrap(ctx, ref.Path, opts.Bootstrap, opts.BootstrapOutput); err != nil {
		return ref, err
	}
	return ref, nil
}

func CreateTask(ctx context.Context, run Ref, taskID string, copyList []string) (TaskRef, error) {
	return CreateTaskWithOptions(ctx, run, taskID, TaskCreateOptions{CopyList: copyList})
}

func CreateTaskWithOptions(ctx context.Context, run Ref, taskID string, opts TaskCreateOptions) (TaskRef, error) {
	if err := validateRef(run); err != nil {
		return TaskRef{}, err
	}
	taskID, err := cleanPathSegment(taskID)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create Task Worktree: %w", err)
	}

	runner := execGitRunner{}
	baseSHA, err := gitRevision(ctx, runner, run.UserRoot, run.Branch)
	if err != nil {
		return TaskRef{}, fmt.Errorf("resolve Run Branch tip for Task Worktree: %w", err)
	}
	ref := TaskRef{
		RunID:    run.RunID,
		TaskID:   taskID,
		Path:     filepath.Join(filepath.Dir(run.Path), taskPathSegment(run.RunID, taskID)),
		Branch:   TaskBranchName(run.RunID, taskID),
		UserRoot: run.UserRoot,
		BaseSHA:  baseSHA,
	}
	if err := os.MkdirAll(filepath.Dir(ref.Path), 0o755); err != nil {
		return TaskRef{}, fmt.Errorf("create Task Worktree parent %q: %w", filepath.Dir(ref.Path), err)
	}
	if _, err := runner.Run(ctx, ref.UserRoot, "worktree", "add", "-b", ref.Branch, ref.Path, baseSHA); err != nil {
		return TaskRef{}, fmt.Errorf("create Task Worktree: %w", err)
	}
	if err := copyProvisionedFiles(ref.UserRoot, ref.Path, opts.CopyList); err != nil {
		return TaskRef{}, err
	}
	if err := runBootstrap(ctx, ref.Path, opts.Bootstrap, opts.BootstrapOutput); err != nil {
		return ref, err
	}
	return ref, nil
}

func TaskRefFor(run Ref, taskID string) (TaskRef, error) {
	if err := validateRef(run); err != nil {
		return TaskRef{}, err
	}
	taskID, err := cleanPathSegment(taskID)
	if err != nil {
		return TaskRef{}, fmt.Errorf("derive Task Worktree ref: %w", err)
	}
	return TaskRef{
		RunID:    run.RunID,
		TaskID:   taskID,
		Path:     filepath.Join(filepath.Dir(run.Path), taskPathSegment(run.RunID, taskID)),
		Branch:   TaskBranchName(run.RunID, taskID),
		UserRoot: run.UserRoot,
	}, nil
}

func Integrate(ctx context.Context, ref Ref, targetBranch, runSHA string) (IntegrationResult, error) {
	if err := validateRef(ref); err != nil {
		return IntegrationResult{}, err
	}
	targetBranch = strings.TrimSpace(targetBranch)
	if targetBranch == "" {
		return IntegrationResult{}, errors.New("integrate Run Worktree: target branch is required")
	}
	runSHA = strings.TrimSpace(runSHA)
	if runSHA == "" {
		return IntegrationResult{}, errors.New("integrate Run Worktree: Run SHA is required")
	}

	runner := execGitRunner{}
	checkedOutPath, err := checkedOutBranchPath(ctx, runner, ref.UserRoot, targetBranch)
	if err != nil {
		return IntegrationResult{}, err
	}
	if samePath(checkedOutPath, ref.UserRoot) {
		if _, err := runner.Run(ctx, ref.UserRoot, "merge", "--ff-only", runSHA); err != nil {
			return IntegrationResult{Mode: ModePending, Reason: classifyMergeRefusal(err)}, nil
		}
		return IntegrationResult{Mode: ModeFastForwardMerge}, nil
	}
	if checkedOutPath != "" {
		return IntegrationResult{Mode: ModePending, Reason: ReasonCheckedOutElsewhere}, nil
	}

	if _, err := runner.Run(ctx, ref.UserRoot, "merge-base", "--is-ancestor", targetBranch, runSHA); err != nil {
		if isAncestryMiss(err) {
			return IntegrationResult{Mode: ModePending, Reason: ReasonNonAncestry}, nil
		}
		return IntegrationResult{}, fmt.Errorf("check target branch ancestry: %w", err)
	}
	if _, err := runner.Run(ctx, ref.UserRoot, "branch", "-f", targetBranch, runSHA); err != nil {
		return IntegrationResult{}, fmt.Errorf("move target branch: %w", err)
	}
	return IntegrationResult{Mode: ModeBranchMove}, nil
}

func IntegrateTask(ctx context.Context, run Ref, task TaskRef) (TaskIntegration, error) {
	if err := validateRef(run); err != nil {
		return TaskIntegration{}, err
	}
	if err := validateTaskRef(task); err != nil {
		return TaskIntegration{}, err
	}
	if run.RunID != task.RunID {
		return TaskIntegration{}, fmt.Errorf("integrate Task Worktree: Task Run ID %q does not match Run ID %q", task.RunID, run.RunID)
	}
	if filepath.Clean(run.UserRoot) != filepath.Clean(task.UserRoot) {
		return TaskIntegration{}, errors.New("integrate Task Worktree: Task and Run must share one user root")
	}

	runner := execGitRunner{}
	runTip, err := gitRevision(ctx, runner, run.UserRoot, run.Branch)
	if err != nil {
		return TaskIntegration{}, fmt.Errorf("resolve Run Branch tip: %w", err)
	}
	baseSHA, err := taskBaseSHA(ctx, runner, run, task)
	if err != nil {
		return TaskIntegration{}, err
	}
	if runTip == baseSHA {
		if _, err := runner.Run(ctx, run.Path, "merge", "--ff-only", task.Branch); err != nil {
			return TaskIntegration{}, fmt.Errorf("fast-forward Run Branch from Task Branch %q: %w", task.Branch, err)
		}
		return TaskIntegration{Mode: ModeTaskFastForward}, nil
	}

	commits, err := taskCommits(ctx, runner, run.UserRoot, baseSHA, task.Branch)
	if err != nil {
		return TaskIntegration{}, err
	}
	if len(commits) == 0 {
		return TaskIntegration{Mode: ModeTaskCherryPick}, nil
	}
	args := append([]string{"cherry-pick"}, commits...)
	if _, err := runner.Run(ctx, run.Path, args...); err != nil {
		paths, pathsErr := conflictingPaths(ctx, runner, run.Path)
		if pathsErr != nil {
			return TaskIntegration{}, errors.Join(
				fmt.Errorf("cherry-pick Task Branch %q: %w", task.Branch, err),
				fmt.Errorf("list conflicting paths: %w", pathsErr),
			)
		}
		if len(paths) == 0 {
			return TaskIntegration{}, fmt.Errorf("cherry-pick Task Branch %q: %w", task.Branch, err)
		}
		if _, abortErr := runner.Run(ctx, run.Path, "cherry-pick", "--abort"); abortErr != nil {
			return TaskIntegration{}, errors.Join(
				fmt.Errorf("cherry-pick Task Branch %q: %w", task.Branch, err),
				fmt.Errorf("abort conflicting Task cherry-pick: %w", abortErr),
			)
		}
		afterAbort, tipErr := gitRevision(ctx, runner, run.UserRoot, run.Branch)
		if tipErr != nil {
			return TaskIntegration{}, fmt.Errorf("verify Run Branch after Task conflict abort: %w", tipErr)
		}
		if afterAbort != runTip {
			return TaskIntegration{}, fmt.Errorf("Task conflict abort left Run Branch at %s, expected %s", afterAbort, runTip)
		}
		return TaskIntegration{Mode: ModeTaskConflict, Reason: strings.Join(paths, ", ")}, nil
	}
	return TaskIntegration{Mode: ModeTaskCherryPick}, nil
}

func CleanupClean(ctx context.Context, ref Ref) error {
	if err := validateRef(ref); err != nil {
		return err
	}
	runner := execGitRunner{}
	if _, err := runner.Run(ctx, ref.UserRoot, "worktree", "remove", "--force", ref.Path); err != nil {
		return fmt.Errorf("remove Run Worktree %q: %w", ref.Path, err)
	}
	if err := deleteRunBranch(ctx, runner, ref.UserRoot, ref.Branch); err != nil {
		return err
	}
	return nil
}

func CleanupTask(ctx context.Context, task TaskRef) error {
	if err := validateTaskRef(task); err != nil {
		return err
	}
	runner := execGitRunner{}
	if _, err := runner.Run(ctx, task.UserRoot, "worktree", "remove", "--force", task.Path); err != nil {
		return fmt.Errorf("remove Task Worktree %q: %w", task.Path, err)
	}
	if err := deleteRunBranch(ctx, runner, task.UserRoot, task.Branch); err != nil {
		return err
	}
	return nil
}

func PruneTerminal(ctx context.Context, userRoot string, location string, isTerminalRun func(runID string) bool) error {
	_, err := PruneTerminalReport(ctx, userRoot, location, isTerminalRun)
	return err
}

func PruneTerminalReport(ctx context.Context, userRoot string, location string, isTerminalRun func(runID string) bool) ([]PrunedRef, error) {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return nil, errors.New("prune Run Worktrees: user root is required")
	}
	if isTerminalRun == nil {
		return nil, errors.New("prune Run Worktrees: terminal callback is required")
	}

	runner := execGitRunner{}
	if _, err := runner.Run(ctx, userRoot, "worktree", "prune"); err != nil {
		return nil, fmt.Errorf("prune git worktrees: %w", err)
	}

	refs, err := terminalCandidates(ctx, runner, userRoot, location)
	if err != nil {
		return nil, err
	}

	var pruned []PrunedRef
	var errs []error
	for _, ref := range refs {
		if !isTerminalRun(ref.RunID) {
			continue
		}
		hasCommits, err := branchHasCommitsBeyondBase(ctx, runner, userRoot, ref.Branch)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if hasCommits {
			continue
		}
		if _, err := os.Stat(ref.Path); err == nil {
			if _, err := runner.Run(ctx, userRoot, "worktree", "remove", "--force", ref.Path); err != nil {
				errs = append(errs, fmt.Errorf("remove terminal Worktree %q: %w", ref.Path, err))
				continue
			}
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("stat terminal Worktree %q: %w", ref.Path, err))
			continue
		}
		if err := deleteRunBranch(ctx, runner, userRoot, ref.Branch); err != nil {
			errs = append(errs, err)
			continue
		}
		pruned = append(pruned, PrunedRef{
			RunID:  ref.RunID,
			TaskID: ref.TaskID,
			Path:   ref.Path,
			Branch: ref.Branch,
		})
	}
	return pruned, errors.Join(errs...)
}

func BranchName(runID string) string {
	return runBranchPrefix + strings.TrimSpace(runID)
}

func TaskBranchName(runID string, taskID string) string {
	return BranchName(runID) + "-" + strings.TrimSpace(taskID)
}

type gitRunner interface {
	Run(ctx context.Context, workDir string, args ...string) (string, error)
}

type execGitRunner struct{}

func (execGitRunner) Run(ctx context.Context, workDir string, args ...string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(workDir) == "" {
		workDir = "."
	}
	cmdArgs := append([]string{"-C", workDir, "-c", "core.fsmonitor=false"}, args...)
	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", &gitCommandError{
			args:   append([]string(nil), args...),
			stdout: stdout.String(),
			stderr: stderr.String(),
			err:    err,
		}
	}
	return stdout.String(), nil
}

type gitCommandError struct {
	args   []string
	stdout string
	stderr string
	err    error
}

func (err *gitCommandError) Error() string {
	detail := strings.TrimSpace(err.stderr)
	if detail == "" {
		detail = strings.TrimSpace(err.stdout)
	}
	if detail == "" {
		detail = err.err.Error()
	}
	return fmt.Sprintf("git %s failed: %s", strings.Join(err.args, " "), detail)
}

func (err *gitCommandError) Unwrap() error {
	return err.err
}

func newRef(userRoot, location, runID string) (Ref, error) {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return Ref{}, errors.New("create Run Worktree: user root is required")
	}
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return Ref{}, errors.New("create Run Worktree: Run ID is required")
	}
	if strings.ContainsAny(runID, `/\`) {
		return Ref{}, fmt.Errorf("create Run Worktree: Run ID %q must not contain path separators", runID)
	}
	path, err := deriveRootPath(location, userRoot, runID)
	if err != nil {
		return Ref{}, err
	}
	branch := BranchName(runID)
	return Ref{
		RunID:    runID,
		Path:     path,
		Branch:   branch,
		UserRoot: userRoot,
	}, nil
}

func validateRef(ref Ref) error {
	if strings.TrimSpace(ref.RunID) == "" {
		return errors.New("Run Worktree ref: Run ID is required")
	}
	if strings.TrimSpace(ref.Path) == "" {
		return errors.New("Run Worktree ref: path is required")
	}
	if strings.TrimSpace(ref.Branch) == "" {
		return errors.New("Run Worktree ref: Run Branch is required")
	}
	if strings.TrimSpace(ref.UserRoot) == "" {
		return errors.New("Run Worktree ref: user root is required")
	}
	return nil
}

func validateTaskRef(ref TaskRef) error {
	if strings.TrimSpace(ref.RunID) == "" {
		return errors.New("Task Worktree ref: Run ID is required")
	}
	if strings.TrimSpace(ref.TaskID) == "" {
		return errors.New("Task Worktree ref: Task ID is required")
	}
	if strings.TrimSpace(ref.Path) == "" {
		return errors.New("Task Worktree ref: path is required")
	}
	if strings.TrimSpace(ref.Branch) == "" {
		return errors.New("Task Worktree ref: Task Branch is required")
	}
	if strings.TrimSpace(ref.UserRoot) == "" {
		return errors.New("Task Worktree ref: user root is required")
	}
	return nil
}

func deriveRootPath(location, userRoot string, segments ...string) (string, error) {
	location = filepath.Clean(strings.TrimSpace(location))
	if location == "." || location == "" {
		return "", errors.New("derive Worktree path: location is required")
	}
	if !filepath.IsAbs(location) {
		return "", errors.New("derive Worktree path: location must be absolute")
	}
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return "", errors.New("derive Worktree path: user root is required")
	}
	parts := []string{location, repoSlug(userRoot)}
	for _, segment := range segments {
		cleaned, err := cleanPathSegment(segment)
		if err != nil {
			return "", err
		}
		parts = append(parts, cleaned)
	}
	return filepath.Join(parts...), nil
}

func cleanPathSegment(segment string) (string, error) {
	segment = strings.TrimSpace(segment)
	if segment == "" {
		return "", errors.New("derive Worktree path: path segment is required")
	}
	if strings.ContainsAny(segment, `/\`) {
		return "", fmt.Errorf("derive Worktree path: path segment %q must not contain path separators", segment)
	}
	if segment == "." || segment == ".." || filepath.Clean(segment) != segment {
		return "", fmt.Errorf("derive Worktree path: path segment %q must be clean", segment)
	}
	return segment, nil
}

func repoSlug(userRoot string) string {
	clean := filepath.Clean(userRoot)
	sum := sha256.Sum256([]byte(clean))
	hash := hex.EncodeToString(sum[:])[:8]
	return sanitizeSlugBase(filepath.Base(clean)) + "-" + hash
}

func sanitizeSlugBase(value string) string {
	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		switch {
		case unicode.IsLetter(char), unicode.IsDigit(char):
			builder.WriteRune(unicode.ToLower(char))
			lastDash = false
		case char == '-' || char == '_' || char == '.':
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		default:
			if !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
	}
	slug := strings.Trim(builder.String(), "-")
	if slug == "" {
		return "repo"
	}
	return slug
}

func copyProvisionedFiles(userRoot, worktreePath string, copyList []string) error {
	for _, entry := range copyList {
		rel, ok, err := cleanRelativeCopyPath(entry)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		source := filepath.Join(userRoot, rel)
		info, err := os.Stat(source)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Run Worktree copy skipped missing %s\n", filepath.ToSlash(rel))
			continue
		}
		if err != nil {
			return fmt.Errorf("stat Run Worktree copy source %q: %w", filepath.ToSlash(rel), err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("Run Worktree copy source %q is not a regular file", filepath.ToSlash(rel))
		}
		content, err := os.ReadFile(source)
		if err != nil {
			return fmt.Errorf("read Run Worktree copy source %q: %w", filepath.ToSlash(rel), err)
		}
		destination := filepath.Join(worktreePath, rel)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("create Run Worktree copy parent %q: %w", filepath.Dir(destination), err)
		}
		if err := os.WriteFile(destination, content, info.Mode().Perm()); err != nil {
			return fmt.Errorf("write Run Worktree copy destination %q: %w", filepath.ToSlash(rel), err)
		}
	}
	return nil
}

func cleanRelativeCopyPath(path string) (string, bool, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", false, nil
	}
	if filepath.IsAbs(path) {
		return "", false, fmt.Errorf("Run Worktree copy path %q must be relative", path)
	}
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false, fmt.Errorf("Run Worktree copy path %q must stay inside the repository", path)
	}
	return clean, true, nil
}

func runBootstrap(ctx context.Context, worktreeDir string, spec BootstrapSpec, out io.Writer) error {
	command := strings.TrimSpace(spec.Command)
	if command == "" {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	worktreeDir = strings.TrimSpace(worktreeDir)
	if worktreeDir == "" {
		return &BootstrapError{Command: command, Err: errors.New("worktree root is required")}
	}
	if spec.Timeout <= 0 {
		return &BootstrapError{Command: command, Err: errors.New("timeout must be greater than 0")}
	}
	if out == nil {
		out = io.Discard
	}

	runCtx, cancel := context.WithTimeout(ctx, spec.Timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, "sh", "-c", command)
	cmd.Dir = worktreeDir
	tail := &boundedTail{limit: 4096}
	stream := io.MultiWriter(out, tail)
	cmd.Stdout = stream
	cmd.Stderr = stream
	err := cmd.Run()
	if err == nil {
		return nil
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return &BootstrapError{
			Command: command,
			Err:     fmt.Errorf("timed out after %s", spec.Timeout),
			Tail:    tail.String(),
		}
	}
	if errors.Is(runCtx.Err(), context.Canceled) && ctx.Err() != nil {
		return ctx.Err()
	}
	return &BootstrapError{Command: command, Err: err, Tail: tail.String()}
}

type boundedTail struct {
	limit int
	data  []byte
}

func (tail *boundedTail) Write(p []byte) (int, error) {
	if tail.limit <= 0 {
		return len(p), nil
	}
	if len(p) >= tail.limit {
		tail.data = append(tail.data[:0], p[len(p)-tail.limit:]...)
		return len(p), nil
	}
	overflow := len(tail.data) + len(p) - tail.limit
	if overflow > 0 {
		tail.data = append(tail.data[:0], tail.data[overflow:]...)
	}
	tail.data = append(tail.data, p...)
	return len(p), nil
}

func (tail *boundedTail) String() string {
	if tail == nil {
		return ""
	}
	return string(tail.data)
}

func checkedOutBranchPath(ctx context.Context, runner gitRunner, userRoot, targetBranch string) (string, error) {
	output, err := runner.Run(ctx, userRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return "", fmt.Errorf("list git worktrees: %w", err)
	}
	var currentPath string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			currentPath = ""
			continue
		}
		if value, ok := strings.CutPrefix(line, "worktree "); ok {
			currentPath = filepath.Clean(value)
			continue
		}
		if value, ok := strings.CutPrefix(line, "branch "); ok {
			branch := strings.TrimPrefix(value, "refs/heads/")
			if branch == targetBranch {
				return currentPath, nil
			}
		}
	}
	return "", nil
}

func classifyMergeRefusal(err error) string {
	var gitErr *gitCommandError
	if !errors.As(err, &gitErr) {
		return ReasonDirtyTarget
	}
	text := strings.ToLower(gitErr.stderr + "\n" + gitErr.stdout)
	switch {
	case strings.Contains(text, "local changes") && strings.Contains(text, "would be overwritten"):
		return ReasonOverlap
	case strings.Contains(text, "not possible to fast-forward"), strings.Contains(text, "diverg"):
		return ReasonDiverged
	default:
		return ReasonDirtyTarget
	}
}

func isAncestryMiss(err error) bool {
	var gitErr *gitCommandError
	if !errors.As(err, &gitErr) {
		return false
	}
	var exitErr *exec.ExitError
	if !errors.As(gitErr.err, &exitErr) || exitErr.ExitCode() != 1 {
		return false
	}
	return strings.TrimSpace(gitErr.stderr) == ""
}

type terminalCandidate struct {
	RunID    string
	TaskID   string
	Path     string
	Branch   string
	UserRoot string
}

func terminalCandidates(ctx context.Context, runner gitRunner, userRoot string, location string) ([]terminalCandidate, error) {
	refsByBranch := map[string]terminalCandidate{}
	worktreeDir, err := repoWorktreesDir(location, userRoot)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(worktreeDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			ref, ok := terminalCandidateFromPath(worktreeDir, userRoot, entry.Name())
			if ok {
				refsByBranch[ref.Branch] = ref
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read Run Worktree directory %q: %w", worktreeDir, err)
	}

	output, err := runner.Run(ctx, userRoot, "for-each-ref", "--format=%(refname:short)", "refs/heads/roundfix/run-*")
	if err != nil {
		return nil, fmt.Errorf("list Run Branches: %w", err)
	}
	for _, line := range strings.Split(output, "\n") {
		branch := strings.TrimSpace(line)
		if branch == "" || !strings.HasPrefix(branch, runBranchPrefix) {
			continue
		}
		if _, found := refsByBranch[branch]; !found {
			ref, ok := terminalCandidateFromBranch(worktreeDir, userRoot, branch)
			if ok {
				refsByBranch[branch] = ref
			}
		}
	}

	branches := make([]string, 0, len(refsByBranch))
	for branch := range refsByBranch {
		branches = append(branches, branch)
	}
	sort.Strings(branches)
	refs := make([]terminalCandidate, 0, len(branches))
	for _, branch := range branches {
		refs = append(refs, refsByBranch[branch])
	}
	return refs, nil
}

func terminalCandidateFromPath(worktreeDir, userRoot, name string) (terminalCandidate, bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return terminalCandidate{}, false
	}
	runID := name
	taskID := ""
	if idx := strings.LastIndex(name, "."); idx > 0 && idx < len(name)-1 {
		runID = name[:idx]
		taskID = name[idx+1:]
	}
	branch := BranchName(runID)
	if taskID != "" {
		branch = TaskBranchName(runID, taskID)
	}
	return terminalCandidate{
		RunID:    runID,
		TaskID:   taskID,
		Path:     filepath.Join(worktreeDir, name),
		Branch:   branch,
		UserRoot: userRoot,
	}, true
}

func terminalCandidateFromBranch(worktreeDir, userRoot, branch string) (terminalCandidate, bool) {
	suffix := strings.TrimPrefix(branch, runBranchPrefix)
	if suffix == "" {
		return terminalCandidate{}, false
	}
	runID := suffix
	taskID := ""
	if idx := strings.LastIndex(suffix, "-"); idx > 0 && idx < len(suffix)-1 && strings.HasPrefix(suffix[idx+1:], "task_") {
		runID = suffix[:idx]
		taskID = suffix[idx+1:]
	}
	segment := runID
	if taskID != "" {
		segment = taskPathSegment(runID, taskID)
	}
	return terminalCandidate{
		RunID:    runID,
		TaskID:   taskID,
		Path:     filepath.Join(worktreeDir, segment),
		Branch:   branch,
		UserRoot: userRoot,
	}, true
}

func taskPathSegment(runID, taskID string) string {
	return strings.TrimSpace(runID) + "." + strings.TrimSpace(taskID)
}

func repoWorktreesDir(location, userRoot string) (string, error) {
	return deriveRootPath(location, userRoot)
}

func gitRevision(ctx context.Context, runner gitRunner, workDir, revision string) (string, error) {
	output, err := runner.Run(ctx, workDir, "rev-parse", strings.TrimSpace(revision))
	if err != nil {
		return "", err
	}
	sha := strings.TrimSpace(output)
	if sha == "" {
		return "", fmt.Errorf("git revision %q resolved to an empty SHA", revision)
	}
	return sha, nil
}

func taskBaseSHA(ctx context.Context, runner gitRunner, run Ref, task TaskRef) (string, error) {
	baseSHA := strings.TrimSpace(task.BaseSHA)
	if baseSHA != "" {
		return baseSHA, nil
	}
	output, err := runner.Run(ctx, run.UserRoot, "merge-base", run.Branch, task.Branch)
	if err != nil {
		return "", fmt.Errorf("resolve Task Branch base: %w", err)
	}
	baseSHA = strings.TrimSpace(output)
	if baseSHA == "" {
		return "", errors.New("resolve Task Branch base: git returned an empty SHA")
	}
	return baseSHA, nil
}

func taskCommits(ctx context.Context, runner gitRunner, userRoot, baseSHA, taskBranch string) ([]string, error) {
	output, err := runner.Run(ctx, userRoot, "rev-list", "--reverse", baseSHA+".."+taskBranch)
	if err != nil {
		return nil, fmt.Errorf("list Task Branch commits: %w", err)
	}
	var commits []string
	for _, line := range strings.Split(output, "\n") {
		sha := strings.TrimSpace(line)
		if sha != "" {
			commits = append(commits, sha)
		}
	}
	return commits, nil
}

func conflictingPaths(ctx context.Context, runner gitRunner, workDir string) ([]string, error) {
	output, err := runner.Run(ctx, workDir, "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, line := range strings.Split(output, "\n") {
		path := strings.TrimSpace(line)
		if path != "" {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func branchHasCommitsBeyondBase(ctx context.Context, runner gitRunner, userRoot, branch string) (bool, error) {
	tip, err := gitRevision(ctx, runner, userRoot, branch)
	if err != nil {
		return false, fmt.Errorf("resolve terminal Worktree Branch %q tip: %w", branch, err)
	}
	base, found, err := branchCreationBase(ctx, runner, userRoot, branch)
	if err != nil {
		return false, err
	}
	if !found {
		return true, nil
	}
	return tip != base, nil
}

func branchCreationBase(ctx context.Context, runner gitRunner, userRoot, branch string) (string, bool, error) {
	output, err := runner.Run(ctx, userRoot, "reflog", "show", "--format=%H", branch)
	if err != nil {
		return "", false, fmt.Errorf("read Worktree Branch %q reflog: %w", branch, err)
	}
	var base string
	for _, line := range strings.Split(output, "\n") {
		sha := strings.TrimSpace(line)
		if sha != "" {
			base = sha
		}
	}
	if base == "" {
		return "", false, nil
	}
	return base, true, nil
}

func deleteRunBranch(ctx context.Context, runner gitRunner, userRoot, branch string) error {
	if _, err := runner.Run(ctx, userRoot, "branch", "-D", branch); err != nil && !isMissingBranch(err) {
		return fmt.Errorf("delete Run Branch %q: %w", branch, err)
	}
	return nil
}

func isMissingBranch(err error) bool {
	var gitErr *gitCommandError
	if !errors.As(err, &gitErr) {
		return false
	}
	text := strings.ToLower(gitErr.stderr + "\n" + gitErr.stdout)
	return strings.Contains(text, "branch") && strings.Contains(text, "not found")
}

func samePath(left, right string) bool {
	if left == "" || right == "" {
		return false
	}
	return canonicalPath(left) == canonicalPath(right)
}

func canonicalPath(path string) string {
	clean := filepath.Clean(path)
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		return clean
	}
	return filepath.Clean(resolved)
}
