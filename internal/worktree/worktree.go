package worktree

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const (
	runBranchPrefix = "roundfix/run-"

	ModeFastForwardMerge = "ff-merge"
	ModeBranchMove       = "branch-move"
	ModePending          = "pending"

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
	UserRoot string
	Location string
	RunID    string
	HeadSHA  string
	CopyList []string
}

type IntegrationResult struct {
	Mode   string
	Reason string
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
	return ref, nil
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

func CleanupClean(ctx context.Context, ref Ref) error {
	if err := validateRef(ref); err != nil {
		return err
	}
	runner := execGitRunner{}
	if _, err := runner.Run(ctx, ref.UserRoot, "worktree", "remove", ref.Path); err != nil {
		return fmt.Errorf("remove Run Worktree %q: %w", ref.Path, err)
	}
	if err := deleteRunBranch(ctx, runner, ref.UserRoot, ref.Branch); err != nil {
		return err
	}
	return nil
}

func PruneTerminal(ctx context.Context, userRoot string, location string, isTerminalClean func(runID string) bool) error {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return errors.New("prune Run Worktrees: user root is required")
	}
	if isTerminalClean == nil {
		return errors.New("prune Run Worktrees: terminal Clean callback is required")
	}

	runner := execGitRunner{}
	if _, err := runner.Run(ctx, userRoot, "worktree", "prune"); err != nil {
		return fmt.Errorf("prune git worktrees: %w", err)
	}

	refs, err := terminalCandidates(ctx, runner, userRoot, location)
	if err != nil {
		return err
	}

	var errs []error
	for _, ref := range refs {
		if !isTerminalClean(ref.RunID) {
			continue
		}
		if _, err := os.Stat(ref.Path); err == nil {
			if _, err := runner.Run(ctx, userRoot, "worktree", "remove", ref.Path); err != nil {
				errs = append(errs, fmt.Errorf("remove terminal Run Worktree %q: %w", ref.Path, err))
				continue
			}
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("stat terminal Run Worktree %q: %w", ref.Path, err))
			continue
		}
		if err := deleteRunBranch(ctx, runner, userRoot, ref.Branch); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func BranchName(runID string) string {
	return runBranchPrefix + strings.TrimSpace(runID)
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

func terminalCandidates(ctx context.Context, runner gitRunner, userRoot string, location string) ([]Ref, error) {
	refsByRun := map[string]Ref{}
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
			runID := entry.Name()
			refsByRun[runID] = Ref{
				RunID:    runID,
				Path:     filepath.Join(worktreeDir, runID),
				Branch:   runBranchPrefix + runID,
				UserRoot: userRoot,
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
		runID := strings.TrimPrefix(branch, runBranchPrefix)
		if _, found := refsByRun[runID]; !found {
			refsByRun[runID] = Ref{
				RunID:    runID,
				Path:     filepath.Join(worktreeDir, runID),
				Branch:   branch,
				UserRoot: userRoot,
			}
		}
	}

	runIDs := make([]string, 0, len(refsByRun))
	for runID := range refsByRun {
		runIDs = append(runIDs, runID)
	}
	sort.Strings(runIDs)
	refs := make([]Ref, 0, len(runIDs))
	for _, runID := range runIDs {
		refs = append(refs, refsByRun[runID])
	}
	return refs, nil
}

func repoWorktreesDir(location, userRoot string) (string, error) {
	return deriveRootPath(location, userRoot)
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
