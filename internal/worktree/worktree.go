package worktree

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const (
	runBranchPrefix = store.RunBranchPrefix

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

type PendingRunWork struct {
	Branch       string
	WorktreePath string
	AheadCommits int
	FastForward  bool
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

type ReconciliationState string

const (
	ReconciliationSafe         ReconciliationState = "safe"
	ReconciliationSuperseded   ReconciliationState = "superseded"
	ReconciliationUnintegrated ReconciliationState = "unintegrated"
	ReconciliationDirty        ReconciliationState = "dirty"
	ReconciliationUnknown      ReconciliationState = "unknown"
	ReconciliationReleased     ReconciliationState = "released"
)

const (
	reconciliationReasonMaxBytes             = 160
	reconciliationReasonSafe                 = "Run Branch is integrated and Run Worktree is clean"
	reconciliationReasonUnintegrated         = "Run Branch is not integrated into the target branch"
	reconciliationReasonDirty                = "Run Worktree has tracked or untracked changes"
	reconciliationReasonReleased             = "Run Worktree and Run Branch are absent"
	reconciliationReasonRunMetadata          = "recorded Run metadata is invalid"
	reconciliationReasonWorktreePath         = "recorded Run Worktree path is unsafe"
	reconciliationReasonWorktreeUnregistered = "recorded Run Worktree is not registered in the Git root"
	reconciliationReasonWorktreeInspection   = "Run Worktree cleanliness could not be inspected"
	reconciliationReasonRunBranch            = "Run Branch could not be resolved unambiguously"
	reconciliationReasonTargetMetadata       = "recorded target branch is missing or invalid"
	reconciliationReasonTargetBranch         = "target branch could not be resolved unambiguously"
	reconciliationReasonAncestry             = "Run Branch ancestry could not be inspected"
	reconciliationReasonSupersededPrefix     = "Run QA report is superseded by "
)

const qaReportOnlyLogFormat = "%x00%x00%B%x00%x00"

type RunWorktreeReconciliation struct {
	RunID             string
	Outcome           string
	Path              string
	Branch            string
	TargetBranch      string
	RunHead           string
	TargetHead        string
	SupersedingReport string
	Reason            string
	State             ReconciliationState

	evidence *terminalRunReconciliationEvidence
}

// BranchSetClassification classifies the Run Branches attributed to one
// target. Current holds the branch with the newest QA evidence. Releasable
// holds only terminal branches proven superseded by the target or Current.
// Every Preserved branch has an entry in PreservedReasons.
type BranchSetClassification struct {
	Current          string
	CurrentReport    string
	Releasable       []string
	ReleasableProofs map[string]string
	Preserved        []string
	PreservedReasons map[string]string

	evidence *branchSetClassificationEvidence
}

type branchSetClassificationEvidence struct {
	gitRoot      string
	targetBranch string
	specSlug     string
	runs         []store.Run
	releasable   map[string]string
}

func InspectTerminalRun(ctx context.Context, run store.Run) (RunWorktreeReconciliation, error) {
	return inspectTerminalRun(ctx, execGitRunner{}, run)
}

// CountRetainedTerminalRuns reports terminal spec Runs with an existing
// recorded Run Worktree or Run Branch. Git inspection is batched per
// repository so listing cost does not grow with terminal Run history.
func CountRetainedTerminalRuns(ctx context.Context, runs []store.Run) (int, []error) {
	return countRetainedTerminalRuns(ctx, execGitRunner{}, runs)
}

func countRetainedTerminalRuns(ctx context.Context, runner gitRunner, runs []store.Run) (int, []error) {
	type repositoryRuns struct {
		root string
		runs []store.Run
	}

	groupsByRoot := make(map[string]*repositoryRuns)
	groups := make([]*repositoryRuns, 0)
	for _, run := range runs {
		if run.Kind != store.KindImplement || !store.IsTerminalState(run.State) {
			continue
		}
		group, found := groupsByRoot[run.GitRoot]
		if !found {
			group = &repositoryRuns{root: run.GitRoot}
			groupsByRoot[run.GitRoot] = group
			groups = append(groups, group)
		}
		group.runs = append(group.runs, run)
	}

	retained := 0
	var failures []error
	for _, group := range groups {
		branches := map[string]bool{}
		gitRoot, err := recordedGitRoot(ctx, runner, group.root)
		if err != nil {
			failures = append(
				failures,
				fmt.Errorf("inspect retained terminal Runs in repository %q: %w", group.root, err),
			)
		} else {
			branches, err = listRunBranches(ctx, runner, gitRoot)
			if err != nil {
				failures = append(
					failures,
					fmt.Errorf("inspect retained terminal Runs in repository %q: %w", group.root, err),
				)
			}
		}

		for _, run := range group.runs {
			pathPresent, pathErr := retainedRunWorktreePathExists(run.WorkDir)
			if pathErr != nil {
				failures = append(
					failures,
					fmt.Errorf("inspect retained terminal Run %q worktree %q: %w", run.ID, run.WorkDir, pathErr),
				)
			}
			if pathPresent || branches[BranchName(run.ID)] {
				retained++
			}
		}
	}
	return retained, failures
}

func retainedRunWorktreePathExists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	if strings.TrimSpace(path) != path || strings.ContainsAny(path, "\r\n\x00") ||
		!filepath.IsAbs(path) || filepath.Clean(path) != path {
		return false, errors.New("recorded Run Worktree path is invalid")
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, errors.New("recorded Run Worktree path is not a directory")
	}
	return true, nil
}

func listRunBranches(ctx context.Context, runner gitRunner, gitRoot string) (map[string]bool, error) {
	output, err := runner.Run(
		ctx,
		gitRoot,
		"for-each-ref",
		"--format=%(refname:short)",
		"refs/heads/"+runBranchPrefix+"*",
	)
	if err != nil {
		return nil, fmt.Errorf("list Run Branches: %w", err)
	}
	branches := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		branch := strings.TrimSpace(line)
		if strings.HasPrefix(branch, runBranchPrefix) {
			branches[branch] = true
		}
	}
	return branches, nil
}

// ClassifyRunBranchSet classifies existing Run Branches attributed by runs to
// one target. Run metadata is required because Git topology cannot identify a
// branch's target or prove that its Run is terminal.
func ClassifyRunBranchSet(
	ctx context.Context,
	gitRoot string,
	targetBranch string,
	specSlug string,
	runs []store.Run,
) (BranchSetClassification, error) {
	return classifyRunBranchSet(ctx, execGitRunner{}, gitRoot, targetBranch, specSlug, runs)
}

type runBranchSetCandidate struct {
	branch string
	head   string
	report string
	active bool
}

func classifyRunBranchSet(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	targetBranch string,
	specSlug string,
	runs []store.Run,
) (BranchSetClassification, error) {
	result := BranchSetClassification{}
	release := func(branch string, report string) {
		if result.ReleasableProofs == nil {
			result.ReleasableProofs = make(map[string]string)
		}
		if _, releasable := result.ReleasableProofs[branch]; releasable {
			return
		}
		result.Releasable = append(result.Releasable, branch)
		result.ReleasableProofs[branch] = report
	}
	preserve := func(branch string, reason string) {
		if result.PreservedReasons == nil {
			result.PreservedReasons = make(map[string]string)
		}
		if _, preserved := result.PreservedReasons[branch]; preserved {
			return
		}
		result.Preserved = append(result.Preserved, branch)
		result.PreservedReasons[branch] = boundedReconciliationReason(reason)
	}

	root, err := recordedGitRoot(ctx, runner, gitRoot)
	if err != nil {
		return result, fmt.Errorf("classify Run Branch set: %w", err)
	}
	targetBranch = strings.TrimSpace(targetBranch)
	if targetBranch == "" || !validLocalBranch(ctx, runner, root, targetBranch) {
		return result, errors.New("classify Run Branch set: target branch is missing or invalid")
	}
	specSlug = strings.TrimSpace(specSlug)
	cleanSlug, err := cleanPathSegment(specSlug)
	if err != nil || cleanSlug != specSlug {
		return result, errors.New("classify Run Branch set: Spec slug is missing or invalid")
	}
	targetHead, err := resolveUnambiguousLocalBranch(ctx, runner, root, targetBranch)
	if err != nil {
		return result, fmt.Errorf("classify Run Branch set: resolve target branch %q: %w", targetBranch, err)
	}
	if targetHead == "" {
		return result, fmt.Errorf("classify Run Branch set: resolve target branch %q: empty Git object ID", targetBranch)
	}
	branches, err := listRunBranches(ctx, runner, root)
	if err != nil {
		return result, fmt.Errorf("classify Run Branch set: %w", err)
	}

	targetReport, err := newestQAReportAtHead(ctx, runner, root, targetHead, specSlug)
	if err != nil && !errors.Is(err, spec.ErrNoQAReport) {
		return result, fmt.Errorf("classify Run Branch set: inspect target QA Reports: %w", err)
	}
	if errors.Is(err, spec.ErrNoQAReport) {
		targetReport = ""
	}

	seen := make(map[string]struct{})
	candidates := make([]runBranchSetCandidate, 0, len(runs))
	for _, run := range runs {
		if run.Kind != store.KindImplement ||
			strings.TrimSpace(run.LocalBranch) != targetBranch ||
			strings.TrimSpace(run.SpecSlug) != specSlug ||
			!samePath(run.GitRoot, root) {
			continue
		}
		branch := BranchName(run.ID)
		if !branches[branch] {
			continue
		}
		if _, duplicate := seen[branch]; duplicate {
			continue
		}
		seen[branch] = struct{}{}

		active := !store.IsTerminalState(run.State)
		head, resolveErr := resolveUnambiguousLocalBranch(ctx, runner, root, branch)
		if resolveErr != nil || head == "" {
			if active {
				preserve(branch, fmt.Sprintf("Run Branch belongs to Active Run %q", run.ID))
			} else {
				preserve(branch, reconciliationReasonRunBranch)
			}
			continue
		}
		if !active {
			if report, proven := supersedingQAReport(ctx, runner, root, targetHead, head, specSlug); proven {
				release(branch, report)
				continue
			}
		}
		report, reportErr := newestQAReportAtHead(ctx, runner, root, head, specSlug)
		if reportErr != nil {
			if active {
				preserve(branch, fmt.Sprintf("Run Branch belongs to Active Run %q", run.ID))
			} else if errors.Is(reportErr, spec.ErrNoQAReport) {
				preserve(branch, "Run Branch does not carry QA Report evidence")
			} else {
				preserve(branch, "Run Branch QA Report evidence could not be inspected")
			}
			continue
		}
		if targetReport != "" {
			newest, newestErr := spec.NewestQAReportFromPaths([]string{targetReport, report})
			if newestErr != nil || newest != report || report == targetReport {
				if active {
					preserve(branch, fmt.Sprintf("Run Branch belongs to Active Run %q", run.ID))
				} else {
					preserve(branch, "Run Branch does not carry newer QA Report evidence than the target branch")
				}
				continue
			}
		}
		candidates = append(candidates, runBranchSetCandidate{
			branch: branch,
			head:   head,
			report: report,
			active: active,
		})
	}

	if len(candidates) != 0 {
		reports := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			reports = append(reports, candidate.report)
		}
		currentReport, newestErr := spec.NewestQAReportFromPaths(reports)
		if newestErr != nil {
			for _, candidate := range candidates {
				preserve(candidate.branch, "Run Branch current QA Report evidence could not be determined")
			}
		} else {
			currentIndex := -1
			for index, candidate := range candidates {
				if candidate.report != currentReport {
					continue
				}
				if currentIndex != -1 {
					currentIndex = -2
					break
				}
				currentIndex = index
			}
			if currentIndex < 0 {
				for _, candidate := range candidates {
					preserve(candidate.branch, "Run Branch current QA Report evidence is not unique")
				}
			} else {
				current := candidates[currentIndex]
				result.Current = current.branch
				result.CurrentReport = current.report
				for index, candidate := range candidates {
					if index == currentIndex {
						continue
					}
					if candidate.active {
						preserve(candidate.branch, "Run Branch belongs to an Active Run")
						continue
					}
					if report, proven := supersedingQAReport(ctx, runner, root, current.head, candidate.head, specSlug); proven {
						release(candidate.branch, report)
						continue
					}
					preserve(candidate.branch, fmt.Sprintf(
						"Run Branch is not proven superseded by current evidence on %q",
						current.branch,
					))
				}
			}
		}
	}

	sort.Strings(result.Releasable)
	sort.Strings(result.Preserved)
	result.evidence = &branchSetClassificationEvidence{
		gitRoot:      root,
		targetBranch: targetBranch,
		specSlug:     specSlug,
		runs:         cloneRuns(runs),
		releasable:   maps.Clone(result.ReleasableProofs),
	}
	return result, nil
}

// ApplyRunBranchCandidate removes one Run Branch that a prior set
// classification proved releasable. It accepts only a candidate carried by
// that inspected classification, repeats the set proof, and re-inspects the
// candidate's registered worktree before removing either Git surface.
func ApplyRunBranchCandidate(ctx context.Context, inspected BranchSetClassification, branch string) error {
	evidence := inspected.evidence
	if evidence == nil {
		return errors.New("apply Run Branch candidate: classification was not produced by inspection")
	}
	proof, inspectedCandidate := evidence.releasable[branch]
	if !inspectedCandidate || proof == "" {
		return fmt.Errorf("apply Run Branch candidate %q: branch was not inspected as releasable", branch)
	}

	fresh, err := ClassifyRunBranchSet(
		ctx,
		evidence.gitRoot,
		evidence.targetBranch,
		evidence.specSlug,
		evidence.runs,
	)
	if err != nil {
		return fmt.Errorf("revalidate Run Branch candidate %q: %w", branch, err)
	}
	if fresh.ReleasableProofs[branch] != proof {
		return fmt.Errorf("apply Run Branch candidate %q: superseding proof is stale", branch)
	}

	run, found := runForBranch(evidence.runs, branch)
	if !found || !store.IsTerminalState(run.State) {
		return fmt.Errorf("apply Run Branch candidate %q: terminal Run metadata is unavailable", branch)
	}
	terminal, err := InspectTerminalRun(ctx, run)
	if err != nil {
		return fmt.Errorf("revalidate Run Branch candidate %q worktree: %w", branch, err)
	}
	if terminal.State == ReconciliationReleased {
		return nil
	}
	if terminal.State != ReconciliationUnintegrated {
		return fmt.Errorf(
			"apply Run Branch candidate %q: worktree classification changed to %q and must be preserved",
			branch,
			terminal.State,
		)
	}
	return cleanupTerminalRun(ctx, execGitRunner{}, terminal)
}

func runForBranch(runs []store.Run, branch string) (store.Run, bool) {
	for _, run := range runs {
		if BranchName(run.ID) == branch {
			return run, true
		}
	}
	return store.Run{}, false
}

func cloneRuns(runs []store.Run) []store.Run {
	cloned := append([]store.Run(nil), runs...)
	for index := range cloned {
		if cloned[index].OwnerPID != nil {
			pid := *cloned[index].OwnerPID
			cloned[index].OwnerPID = &pid
		}
		if cloned[index].CompletedAt != nil {
			completedAt := *cloned[index].CompletedAt
			cloned[index].CompletedAt = &completedAt
		}
	}
	return cloned
}

// QAReportOnlyBranch reports whether targetHead..runHead contains at least one
// commit and every commit matches the Daemon QA-commit contract for slug while
// changing only paths under that Spec's active or archived qa/ directory.
func QAReportOnlyBranch(
	ctx context.Context,
	gitRoot string,
	targetHead string,
	runHead string,
	slug string,
) (bool, error) {
	return qaReportOnlyBranch(ctx, execGitRunner{}, gitRoot, targetHead, runHead, slug)
}

func qaReportOnlyBranch(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	targetHead string,
	runHead string,
	slug string,
) (bool, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return false, nil
	}
	cleanSlug, err := cleanPathSegment(slug)
	if err != nil || cleanSlug != slug {
		return false, nil
	}
	targetHead = strings.TrimSpace(targetHead)
	runHead = strings.TrimSpace(runHead)
	if targetHead == "" || runHead == "" {
		return false, nil
	}

	output, err := runner.Run(
		ctx,
		gitRoot,
		"log",
		"--reverse",
		"--root",
		"--no-renames",
		"--format="+qaReportOnlyLogFormat,
		"--name-only",
		"-z",
		targetHead+".."+runHead,
	)
	if err != nil {
		return false, fmt.Errorf("read Run commits for QA-report-only proof: %w", err)
	}
	commits, parsed := parseQAReportOnlyLog(output)
	if !parsed || len(commits) == 0 {
		return false, nil
	}
	qaDirs := qaReportDirectories(slug)
	for _, commit := range commits {
		if !matchesQAReportCommitMessage(commit.message, slug) {
			return false, nil
		}
		if len(commit.paths) == 0 {
			return false, nil
		}
		for _, path := range commit.paths {
			if !pathUnderAnyGitDirectory(path, qaDirs) {
				return false, nil
			}
		}
	}
	return true, nil
}

type qaReportOnlyCommit struct {
	message string
	paths   []string
}

func parseQAReportOnlyLog(output string) ([]qaReportOnlyCommit, bool) {
	if output == "" {
		return nil, true
	}
	parts := strings.Split(output, "\x00\x00\x00")
	if len(parts)%2 != 0 || len(parts) == 0 || !strings.HasPrefix(parts[0], "\x00\x00") {
		return nil, false
	}
	parts[0] = strings.TrimPrefix(parts[0], "\x00\x00")

	commits := make([]qaReportOnlyCommit, 0, len(parts)/2)
	for index := 0; index < len(parts); index += 2 {
		pathOutput := parts[index+1]
		if !strings.HasPrefix(pathOutput, "\n") {
			return nil, false
		}
		pathOutput = strings.TrimPrefix(pathOutput, "\n")
		commits = append(commits, qaReportOnlyCommit{
			message: parts[index],
			paths:   nonEmptyNULTerms(pathOutput),
		})
	}
	return commits, true
}

func matchesQAReportCommitMessage(message string, slug string) bool {
	message = strings.TrimRight(message, "\r\n")
	for _, verdict := range []string{spec.VerdictPass, spec.VerdictFail, spec.VerdictPartial} {
		want := fmt.Sprintf(
			"docs: qa report for %s (%s)\n\nRoundfix-Spec: %s",
			slug,
			verdict,
			slug,
		)
		if message == want {
			return true
		}
	}
	return false
}

func qaReportDirectories(slug string) []string {
	return []string{
		filepath.ToSlash(filepath.Join("docs", "specs", slug, "qa")),
		filepath.ToSlash(filepath.Join("docs", "specs", "_archived", slug, "qa")),
	}
}

func pathUnderAnyGitDirectory(path string, directories []string) bool {
	if strings.TrimSpace(path) != path || path == "" {
		return false
	}
	for _, directory := range directories {
		if strings.HasPrefix(path, directory+"/") {
			return true
		}
	}
	return false
}

func nonEmptyNULTerms(output string) []string {
	terms := make([]string, 0)
	for _, term := range strings.Split(output, "\x00") {
		if term != "" {
			terms = append(terms, term)
		}
	}
	return terms
}

type terminalRunReconciliationEvidence struct {
	run              store.Run
	gitRoot          string
	snapshot         terminalRunReconciliationSnapshot
	worktreePresent  bool
	runBranchPresent bool
}

type terminalRunReconciliationSnapshot struct {
	runID             string
	outcome           string
	path              string
	branch            string
	targetBranch      string
	runHead           string
	targetHead        string
	supersedingReport string
	reason            string
	state             ReconciliationState
}

func inspectTerminalRun(ctx context.Context, runner gitRunner, run store.Run) (RunWorktreeReconciliation, error) {
	result := RunWorktreeReconciliation{
		RunID:        run.ID,
		Outcome:      run.State,
		Path:         run.WorkDir,
		Branch:       BranchName(run.ID),
		TargetBranch: run.LocalBranch,
		State:        ReconciliationUnknown,
		Reason:       reconciliationReasonRunMetadata,
	}

	gitRoot, err := recordedGitRoot(ctx, runner, run.GitRoot)
	if err != nil {
		return result, err
	}
	if run.Kind != store.KindImplement || !store.IsTerminalState(run.State) {
		return result, nil
	}
	if strings.TrimSpace(run.ID) != run.ID {
		return result, nil
	}
	if _, err := cleanPathSegment(run.ID); err != nil {
		return result, nil
	}
	if !validLocalBranch(ctx, runner, gitRoot, result.Branch) {
		return result, nil
	}

	worktrees, err := listRegisteredWorktrees(ctx, runner, gitRoot)
	if err != nil {
		result.Reason = reconciliationReasonWorktreeInspection
		return result, nil
	}
	worktree, worktreeState := recordedWorktree(result.Path, worktrees)
	if worktreeState == recordedWorktreeUnsafe {
		result.Reason = reconciliationReasonWorktreePath
		return result, nil
	}
	if worktreeState == recordedWorktreeUnregistered {
		result.Reason = reconciliationReasonWorktreeUnregistered
		return result, nil
	}
	worktreePresent := worktreeState == recordedWorktreePresent

	runBranchPresent, err := localBranchExists(ctx, runner, gitRoot, result.Branch)
	if err != nil {
		result.Reason = reconciliationReasonRunBranch
		return result, nil
	}
	if !worktreePresent && registeredBranchPath(worktrees, result.Branch) != "" {
		result.Reason = reconciliationReasonWorktreeUnregistered
		return result, nil
	}
	if !worktreePresent && !runBranchPresent {
		result.State = ReconciliationReleased
		result.Reason = reconciliationReasonReleased
		result.evidence = newTerminalRunReconciliationEvidence(run, gitRoot, result, false, false)
		return result, nil
	}

	var runHeadErr error
	if runBranchPresent {
		result.RunHead, runHeadErr = resolveUnambiguousLocalBranch(ctx, runner, gitRoot, result.Branch)
	}
	targetMetadataValid := validLocalBranch(ctx, runner, gitRoot, result.TargetBranch)
	targetBranchPresent := false
	var targetHeadErr error
	if targetMetadataValid {
		targetBranchPresent, targetHeadErr = localBranchExists(ctx, runner, gitRoot, result.TargetBranch)
		if targetHeadErr == nil && targetBranchPresent {
			result.TargetHead, targetHeadErr = resolveUnambiguousLocalBranch(ctx, runner, gitRoot, result.TargetBranch)
		}
	}

	if worktreePresent {
		if worktree.Branch != result.Branch {
			result.Reason = reconciliationReasonWorktreeUnregistered
			return result, nil
		}
		status, err := runner.Run(ctx, worktree.Path, "status", "--porcelain=v1", "--untracked-files=all")
		if err != nil {
			result.Reason = reconciliationReasonWorktreeInspection
			return result, nil
		}
		if strings.TrimSpace(status) != "" {
			result.State = ReconciliationDirty
			result.Reason = reconciliationReasonDirty
			return result, nil
		}
	}

	if runHeadErr != nil || result.RunHead == "" {
		result.Reason = reconciliationReasonRunBranch
		return result, nil
	}
	if !targetMetadataValid {
		result.Reason = reconciliationReasonTargetMetadata
		return result, nil
	}
	if targetHeadErr != nil {
		result.Reason = reconciliationReasonTargetBranch
		return result, nil
	}
	if !targetBranchPresent {
		return inspectDeletedTargetRunByContent(
			ctx,
			runner,
			run,
			gitRoot,
			result,
			worktreePresent,
			runBranchPresent,
		), nil
	}
	if result.TargetHead == "" {
		result.Reason = reconciliationReasonTargetBranch
		return result, nil
	}
	if _, err := runner.Run(ctx, gitRoot, "merge-base", "--is-ancestor", result.RunHead, result.TargetHead); err != nil {
		if isAncestryMiss(err) {
			if report, proven := supersedingQAReport(
				ctx,
				runner,
				gitRoot,
				result.TargetHead,
				result.RunHead,
				run.SpecSlug,
			); proven {
				result.State = ReconciliationSuperseded
				result.SupersedingReport = report
				result.Reason = supersededReconciliationReason(report)
				result.evidence = newTerminalRunReconciliationEvidence(
					run,
					gitRoot,
					result,
					worktreePresent,
					runBranchPresent,
				)
				return result, nil
			}
			result.State = ReconciliationUnintegrated
			result.Reason = reconciliationReasonUnintegrated
			result.evidence = newTerminalRunReconciliationEvidence(
				run,
				gitRoot,
				result,
				worktreePresent,
				runBranchPresent,
			)
			return result, nil
		}
		result.Reason = reconciliationReasonAncestry
		return result, nil
	}
	result.State = ReconciliationSafe
	result.Reason = reconciliationReasonSafe
	result.evidence = newTerminalRunReconciliationEvidence(run, gitRoot, result, worktreePresent, runBranchPresent)
	return result, nil
}

type preflightGitRunner struct {
	runner gitRunner
}

func (adapter preflightGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	return adapter.runner.Run(ctx, workDir, args...)
}

func inspectDeletedTargetRunByContent(
	ctx context.Context,
	runner gitRunner,
	run store.Run,
	gitRoot string,
	result RunWorktreeReconciliation,
	worktreePresent bool,
	runBranchPresent bool,
) RunWorktreeReconciliation {
	currentBranchOutput, currentBranchErr := runner.Run(ctx, gitRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
	currentBranch := ""
	if currentBranchErr == nil {
		currentBranch = strings.TrimSpace(currentBranchOutput)
	}
	defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, currentBranch, preflightGitRunner{runner: runner})
	if defaultBranch.Source == preflight.DefaultBranchUndetermined ||
		!validLocalBranch(ctx, runner, gitRoot, defaultBranch.Name) {
		result.Reason = reconciliationReasonTargetBranch
		return result
	}

	defaultHead, err := resolveUnambiguousLocalBranch(ctx, runner, gitRoot, defaultBranch.Name)
	if err != nil || defaultHead == "" {
		result.Reason = reconciliationReasonTargetBranch
		return result
	}
	result.TargetHead = defaultHead

	runOnly, differingShared, retainedRunDeletions, proven := compareRunContentToDefault(
		ctx,
		runner,
		gitRoot,
		result.RunHead,
		defaultHead,
	)
	if !proven {
		result.State = ReconciliationUnintegrated
		result.Reason = boundedReconciliationReason(fmt.Sprintf(
			"Run Branch content comparison could not prove integration against default branch %q",
			defaultBranch.Name,
		))
		return result
	}
	if runOnly != 0 || differingShared != 0 || retainedRunDeletions != 0 {
		var evidence []string
		if runOnly != 0 {
			evidence = append(evidence, fmt.Sprintf("%d Run-only file%s", runOnly, pluralSuffix(runOnly)))
		}
		if differingShared != 0 {
			evidence = append(evidence, fmt.Sprintf(
				"%d differing shared file%s",
				differingShared,
				pluralSuffix(differingShared),
			))
		}
		if retainedRunDeletions != 0 {
			evidence = append(evidence, fmt.Sprintf(
				"%d Run-deleted file%s retained by default",
				retainedRunDeletions,
				pluralSuffix(retainedRunDeletions),
			))
		}
		result.State = ReconciliationUnintegrated
		result.Reason = boundedReconciliationReason(fmt.Sprintf(
			"Run Branch content is not fully represented: %s against default branch %q",
			strings.Join(evidence, ", "),
			defaultBranch.Name,
		))
		return result
	}

	result.State = ReconciliationSafe
	result.Reason = boundedReconciliationReason(fmt.Sprintf(
		"Run Branch content is fully represented on default branch %q",
		defaultBranch.Name,
	))
	result.evidence = newTerminalRunReconciliationEvidence(
		run,
		gitRoot,
		result,
		worktreePresent,
		runBranchPresent,
	)
	return result
}

func compareRunContentToDefault(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	runHead string,
	defaultHead string,
) (runOnly int, differingShared int, retainedRunDeletions int, proven bool) {
	runOnlyOutput, err := runner.Run(
		ctx,
		gitRoot,
		"diff",
		"--name-only",
		"-z",
		"--no-renames",
		"--diff-filter=D",
		runHead,
		defaultHead,
		"--",
	)
	if err != nil {
		return 0, 0, 0, false
	}
	differingSharedOutput, err := runner.Run(
		ctx,
		gitRoot,
		"diff",
		"--name-only",
		"-z",
		"--no-renames",
		"--diff-filter=MT",
		runHead,
		defaultHead,
		"--",
	)
	if err != nil {
		return 0, 0, 0, false
	}
	mergeBaseOutput, err := runner.Run(ctx, gitRoot, "merge-base", runHead, defaultHead)
	if err != nil {
		return 0, 0, 0, false
	}
	mergeBase := strings.TrimSpace(mergeBaseOutput)
	if mergeBase == "" {
		return 0, 0, 0, false
	}
	runDeletedOutput, err := runner.Run(
		ctx,
		gitRoot,
		"diff",
		"--name-only",
		"-z",
		"--no-renames",
		"--diff-filter=D",
		mergeBase,
		runHead,
		"--",
	)
	if err != nil {
		return 0, 0, 0, false
	}
	runDeleted := nonEmptyNULTerms(runDeletedOutput)
	if len(runDeleted) != 0 {
		defaultTreeOutput, err := runner.Run(ctx, gitRoot, "ls-tree", "-r", "--name-only", "-z", defaultHead)
		if err != nil {
			return 0, 0, 0, false
		}
		defaultPaths := make(map[string]struct{})
		for _, path := range nonEmptyNULTerms(defaultTreeOutput) {
			defaultPaths[path] = struct{}{}
		}
		for _, path := range runDeleted {
			if _, retained := defaultPaths[path]; retained {
				retainedRunDeletions++
			}
		}
	}
	return len(nonEmptyNULTerms(runOnlyOutput)),
		len(nonEmptyNULTerms(differingSharedOutput)),
		retainedRunDeletions,
		true
}

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func supersededReconciliationReason(report string) string {
	reason := reconciliationReasonSupersededPrefix + strings.Map(func(char rune) rune {
		if char == '\r' || char == '\n' {
			return ' '
		}
		return char
	}, report)
	return boundedReconciliationReason(reason)
}

func boundedReconciliationReason(reason string) string {
	if len(reason) <= reconciliationReasonMaxBytes {
		return reason
	}

	const ellipsis = "..."
	budget := reconciliationReasonMaxBytes - len(ellipsis)
	cut := 0
	for index, char := range reason {
		next := index + utf8.RuneLen(char)
		if next > budget {
			break
		}
		cut = next
	}
	return reason[:cut] + ellipsis
}

func newTerminalRunReconciliationEvidence(run store.Run, gitRoot string, result RunWorktreeReconciliation, worktreePresent, runBranchPresent bool) *terminalRunReconciliationEvidence {
	return &terminalRunReconciliationEvidence{
		run:              run,
		gitRoot:          gitRoot,
		snapshot:         terminalRunSnapshot(result),
		worktreePresent:  worktreePresent,
		runBranchPresent: runBranchPresent,
	}
}

// SupersedingQAReport reports the target-side QA Report that supersedes
// targetHead..runHead, and whether supersession is proven. Both halves are
// required: QAReportOnlyBranch proves the branch holds nothing but QA reports,
// which is not proof that a newer report exists to supersede it. Callers that
// act on supersession — the reconcile classifier and Branch Integrity
// Preflight — must agree, or one offers a release the other refuses.
func SupersedingQAReport(
	ctx context.Context,
	gitRoot string,
	targetHead string,
	runHead string,
	slug string,
) (string, bool) {
	return supersedingQAReport(ctx, execGitRunner{}, gitRoot, targetHead, runHead, slug)
}

func supersedingQAReport(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	targetHead string,
	runHead string,
	slug string,
) (string, bool) {
	qaOnly, err := qaReportOnlyBranch(ctx, runner, gitRoot, targetHead, runHead, slug)
	if err != nil || !qaOnly {
		return "", false
	}
	runReport, err := newestQAReportAtHead(ctx, runner, gitRoot, runHead, slug)
	if err != nil {
		return "", false
	}
	targetReport, err := newestQAReportAtHead(ctx, runner, gitRoot, targetHead, slug)
	if err != nil || targetReport == runReport {
		return "", false
	}
	newest, err := spec.NewestQAReportFromPaths([]string{runReport, targetReport})
	if err != nil || newest != targetReport {
		return "", false
	}
	return targetReport, true
}

func newestQAReportAtHead(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	head string,
	slug string,
) (string, error) {
	directories := qaReportDirectories(slug)
	args := []string{"ls-tree", "-r", "--name-only", "-z", head, "--"}
	args = append(args, directories...)
	output, err := runner.Run(ctx, gitRoot, args...)
	if err != nil {
		return "", fmt.Errorf("list QA Reports at Git head %q: %w", head, err)
	}
	var reports []string
	for _, path := range nonEmptyNULTerms(output) {
		if !pathUnderAnyGitDirectory(path, directories) {
			continue
		}
		name := filepath.Base(filepath.FromSlash(path))
		if strings.HasPrefix(name, "qa-report-") && strings.HasSuffix(name, ".md") {
			reports = append(reports, path)
		}
	}
	return spec.NewestQAReportFromPaths(reports)
}

func terminalRunSnapshot(result RunWorktreeReconciliation) terminalRunReconciliationSnapshot {
	return terminalRunReconciliationSnapshot{
		runID:             result.RunID,
		outcome:           result.Outcome,
		path:              result.Path,
		branch:            result.Branch,
		targetBranch:      result.TargetBranch,
		runHead:           result.RunHead,
		targetHead:        result.TargetHead,
		supersedingReport: result.SupersedingReport,
		reason:            result.Reason,
		state:             result.State,
	}
}

type TerminalRunApplyError struct {
	RunID              string
	WorktreePath       string
	RunBranch          string
	WorktreeRemaining  bool
	RunBranchRemaining bool
	Err                error
}

func (err *TerminalRunApplyError) Error() string {
	if err == nil {
		return "apply terminal Run reconciliation failed"
	}
	var remaining []string
	if err.WorktreeRemaining {
		remaining = append(remaining, "worktree="+err.WorktreePath)
	}
	if err.RunBranchRemaining {
		remaining = append(remaining, "branch="+err.RunBranch)
	}
	if len(remaining) == 0 {
		remaining = append(remaining, "none")
	}
	return fmt.Sprintf("apply terminal Run %q: %v; remaining: %s", err.RunID, err.Err, strings.Join(remaining, " "))
}

func (err *TerminalRunApplyError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

type TerminalRunReconciliationStore interface {
	ReconcileIntegration(context.Context, store.IntegrationReconciliation) (store.Run, error)
}

func ApplyTerminalRun(ctx context.Context, runStore TerminalRunReconciliationStore, result RunWorktreeReconciliation) error {
	return applyTerminalRunWithStore(ctx, execGitRunner{}, runStore, result)
}

func applyTerminalRunWithStore(
	ctx context.Context,
	runner gitRunner,
	runStore TerminalRunReconciliationStore,
	result RunWorktreeReconciliation,
) error {
	fresh, released, err := revalidateTerminalRunApply(ctx, runner, result)
	if err != nil || released {
		return err
	}
	if runStore == nil {
		return errors.New("apply terminal Run reconciliation: Run Store is required before cleanup")
	}
	if _, err := runStore.ReconcileIntegration(ctx, store.IntegrationReconciliation{
		RunID:           fresh.RunID,
		PreviousOutcome: fresh.Outcome,
		Classification:  string(fresh.State),
		RunBranch:       fresh.Branch,
		RunHead:         fresh.RunHead,
		TargetBranch:    fresh.TargetBranch,
		TargetHead:      fresh.TargetHead,
		Worktree:        fresh.Path,
		Reason:          fresh.Reason,
		Action:          "cleanup",
		Time:            time.Now().UTC(),
	}); err != nil {
		return fmt.Errorf("persist terminal Run %q reconciliation before cleanup: %w", fresh.RunID, err)
	}
	return cleanupTerminalRun(ctx, runner, fresh)
}

func applyTerminalRun(ctx context.Context, runner gitRunner, result RunWorktreeReconciliation) error {
	fresh, released, err := revalidateTerminalRunApply(ctx, runner, result)
	if err != nil || released {
		return err
	}
	return cleanupTerminalRun(ctx, runner, fresh)
}

func revalidateTerminalRunApply(
	ctx context.Context,
	runner gitRunner,
	result RunWorktreeReconciliation,
) (RunWorktreeReconciliation, bool, error) {
	evidence := result.evidence
	if evidence == nil || terminalRunSnapshot(result) != evidence.snapshot {
		return RunWorktreeReconciliation{}, false, errors.New("apply terminal Run reconciliation: result was not produced by inspection or recorded metadata changed")
	}
	if result.State != ReconciliationSafe &&
		result.State != ReconciliationSuperseded &&
		result.State != ReconciliationReleased {
		return RunWorktreeReconciliation{}, false, fmt.Errorf("apply terminal Run reconciliation: classification %q is not safe", result.State)
	}

	fresh, err := inspectTerminalRun(ctx, runner, evidence.run)
	if err != nil {
		return RunWorktreeReconciliation{}, false, fmt.Errorf("revalidate terminal Run before cleanup: %w", err)
	}
	if fresh.State == ReconciliationReleased {
		return fresh, true, nil
	}
	if fresh.State != result.State ||
		fresh.RunHead != result.RunHead ||
		fresh.TargetHead != result.TargetHead {
		return RunWorktreeReconciliation{}, false, fmt.Errorf(
			"apply terminal Run reconciliation: evidence is stale: inspected state=%q Run head=%q target head=%q; current state=%q Run head=%q target head=%q",
			result.State,
			result.RunHead,
			result.TargetHead,
			fresh.State,
			fresh.RunHead,
			fresh.TargetHead,
		)
	}
	return fresh, false, nil
}

func cleanupTerminalRun(ctx context.Context, runner gitRunner, fresh RunWorktreeReconciliation) error {
	evidence := fresh.evidence
	if fresh.evidence.worktreePresent {
		if _, err := runner.Run(ctx, evidence.gitRoot, "worktree", "remove", fresh.Path); err != nil {
			return terminalRunApplyFailure(
				ctx,
				runner,
				evidence.gitRoot,
				fresh,
				fmt.Errorf("remove Run Worktree %q: %w", fresh.Path, err),
			)
		}
	}
	if fresh.evidence.runBranchPresent {
		if err := deleteRunBranch(ctx, runner, evidence.gitRoot, fresh.Branch); err != nil {
			return terminalRunApplyFailure(ctx, runner, evidence.gitRoot, fresh, err)
		}
	}
	return nil
}

func terminalRunApplyFailure(ctx context.Context, runner gitRunner, gitRoot string, result RunWorktreeReconciliation, operationErr error) error {
	worktreeRemaining, branchRemaining, inspectErr := remainingTerminalRunSurface(ctx, runner, gitRoot, result)
	if inspectErr != nil {
		operationErr = errors.Join(operationErr, fmt.Errorf("inspect remaining terminal Run surface: %w", inspectErr))
	}
	return &TerminalRunApplyError{
		RunID:              result.RunID,
		WorktreePath:       result.Path,
		RunBranch:          result.Branch,
		WorktreeRemaining:  worktreeRemaining,
		RunBranchRemaining: branchRemaining,
		Err:                operationErr,
	}
}

func remainingTerminalRunSurface(ctx context.Context, runner gitRunner, gitRoot string, result RunWorktreeReconciliation) (bool, bool, error) {
	worktrees, worktreeErr := listRegisteredWorktrees(ctx, runner, gitRoot)
	worktreeRemaining := samePath(registeredBranchPath(worktrees, result.Branch), result.Path)
	if !worktreeRemaining {
		if _, err := os.Stat(result.Path); err == nil {
			worktreeRemaining = true
		} else if !errors.Is(err, os.ErrNotExist) {
			worktreeErr = errors.Join(worktreeErr, fmt.Errorf("stat Run Worktree %q: %w", result.Path, err))
		}
	}
	branchRemaining, branchErr := localBranchExists(ctx, runner, gitRoot, result.Branch)
	return worktreeRemaining, branchRemaining, errors.Join(worktreeErr, branchErr)
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

// PriorChangedFiles returns repository-relative paths changed between the
// Run's persisted initial HEAD and the current worktree HEAD.
func PriorChangedFiles(ctx context.Context, workDir string, initialHead string) ([]string, error) {
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return nil, errors.New("resolve prior changed files: worktree root is required")
	}
	initialHead = strings.TrimSpace(initialHead)
	if initialHead == "" {
		return nil, errors.New("resolve prior changed files: initial HEAD is required")
	}
	runner := execGitRunner{}
	output, err := runner.Run(ctx, workDir, "diff", "--name-only", initialHead+"..HEAD", "--")
	if err != nil {
		return nil, fmt.Errorf("resolve prior changed files: %w", err)
	}
	seen := map[string]bool{}
	var paths []string
	for _, line := range strings.Split(output, "\n") {
		path := strings.TrimSpace(line)
		if path == "" {
			continue
		}
		clean := filepath.Clean(path)
		if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.IsAbs(clean) {
			return nil, fmt.Errorf("resolve prior changed files: git returned invalid path %q", path)
		}
		clean = filepath.ToSlash(clean)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		paths = append(paths, clean)
	}
	sort.Strings(paths)
	return paths, nil
}

func ListPendingRunWork(ctx context.Context, userRoot, baseBranch string) ([]PendingRunWork, error) {
	runner := execGitRunner{}
	return listPendingRunWork(ctx, runner, userRoot, baseBranch)
}

func IntegratePendingRunWork(ctx context.Context, userRoot, baseBranch, runBranch string) error {
	runner := execGitRunner{}
	return integratePendingRunWork(ctx, runner, userRoot, baseBranch, runBranch)
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

func listPendingRunWork(ctx context.Context, runner gitRunner, userRoot, baseBranch string) ([]PendingRunWork, error) {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return nil, errors.New("list pending Run Branch work: user root is required")
	}
	baseBranch = strings.TrimSpace(baseBranch)
	if baseBranch == "" {
		return nil, errors.New("list pending Run Branch work: base branch is required")
	}

	worktreePaths, err := worktreePathsByBranch(ctx, runner, userRoot)
	if err != nil {
		return nil, fmt.Errorf("list pending Run Branch work: %w", err)
	}
	output, err := runner.Run(ctx, userRoot, "for-each-ref", "--format=%(refname:short)", "refs/heads/roundfix/run-*")
	if err != nil {
		return nil, fmt.Errorf("list pending Run Branch work: list Run Branches: %w", err)
	}

	var pending []PendingRunWork
	for _, line := range strings.Split(output, "\n") {
		branch := strings.TrimSpace(line)
		if branch == "" || !strings.HasPrefix(branch, runBranchPrefix) || isTaskBranch(branch) {
			continue
		}
		ahead, err := branchAheadCommits(ctx, runner, userRoot, baseBranch, branch)
		if err != nil {
			return nil, fmt.Errorf("list pending Run Branch work: %w", err)
		}
		if ahead == 0 {
			continue
		}
		fastForward, err := branchCanFastForward(ctx, runner, userRoot, baseBranch, branch)
		if err != nil {
			return nil, fmt.Errorf("list pending Run Branch work: %w", err)
		}
		pending = append(pending, PendingRunWork{
			Branch:       branch,
			WorktreePath: worktreePaths[branch],
			AheadCommits: ahead,
			FastForward:  fastForward,
		})
	}
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Branch < pending[j].Branch
	})
	return pending, nil
}

func integratePendingRunWork(ctx context.Context, runner gitRunner, userRoot, baseBranch, runBranch string) error {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return errors.New("integrate pending Run Branch work: user root is required")
	}
	baseBranch = strings.TrimSpace(baseBranch)
	if baseBranch == "" {
		return errors.New("integrate pending Run Branch work: base branch is required")
	}
	runBranch = strings.TrimSpace(runBranch)
	if runBranch == "" {
		return errors.New("integrate pending Run Branch work: Run Branch is required")
	}
	if !strings.HasPrefix(runBranch, runBranchPrefix) || isTaskBranch(runBranch) {
		return fmt.Errorf("integrate pending Run Branch work: %q is not a Run Branch", runBranch)
	}

	fastForward, err := branchCanFastForward(ctx, runner, userRoot, baseBranch, runBranch)
	if err != nil {
		return fmt.Errorf("integrate pending Run Branch %q into %q: %w", runBranch, baseBranch, err)
	}
	if !fastForward {
		return fmt.Errorf("integrate pending Run Branch %q into %q: fast-forward is impossible", runBranch, baseBranch)
	}
	checkedOutPath, err := checkedOutBranchPath(ctx, runner, userRoot, baseBranch)
	if err != nil {
		return fmt.Errorf("integrate pending Run Branch %q into %q: %w", runBranch, baseBranch, err)
	}
	if samePath(checkedOutPath, userRoot) {
		if _, err := runner.Run(ctx, userRoot, "merge", "--ff-only", runBranch); err != nil {
			return fmt.Errorf("integrate pending Run Branch %q into %q: git merge --ff-only: %w", runBranch, baseBranch, err)
		}
		return nil
	}
	if checkedOutPath != "" {
		return fmt.Errorf("integrate pending Run Branch %q into %q: base branch is checked out at %s", runBranch, baseBranch, checkedOutPath)
	}

	runTip, err := gitRevision(ctx, runner, userRoot, runBranch)
	if err != nil {
		return fmt.Errorf("integrate pending Run Branch %q into %q: resolve Run Branch tip: %w", runBranch, baseBranch, err)
	}
	if _, err := runner.Run(ctx, userRoot, "branch", "-f", baseBranch, runTip); err != nil {
		return fmt.Errorf("integrate pending Run Branch %q into %q: git branch -f: %w", runBranch, baseBranch, err)
	}
	return nil
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

type TerminalRunLookup func(ctx context.Context, runID string) (store.Run, bool, error)

func PruneTerminal(
	ctx context.Context,
	userRoot string,
	location string,
	runStore TerminalRunReconciliationStore,
	loadTerminalRun TerminalRunLookup,
) error {
	_, err := PruneTerminalReport(ctx, userRoot, location, runStore, loadTerminalRun)
	return err
}

func PruneTerminalReport(
	ctx context.Context,
	userRoot string,
	location string,
	runStore TerminalRunReconciliationStore,
	loadTerminalRun TerminalRunLookup,
) ([]PrunedRef, error) {
	userRoot = filepath.Clean(strings.TrimSpace(userRoot))
	if userRoot == "." || userRoot == "" {
		return nil, errors.New("prune Run Worktrees: user root is required")
	}
	if runStore == nil {
		return nil, errors.New("prune Run Worktrees: Run Store is required before cleanup")
	}
	if loadTerminalRun == nil {
		return nil, errors.New("prune Run Worktrees: terminal Run lookup is required")
	}

	runner := execGitRunner{}
	refs, err := terminalCandidates(ctx, runner, userRoot, location)
	if err != nil {
		return nil, err
	}

	refsByRun := make(map[string][]terminalCandidate)
	var runIDs []string
	for _, ref := range refs {
		if _, found := refsByRun[ref.RunID]; !found {
			runIDs = append(runIDs, ref.RunID)
		}
		refsByRun[ref.RunID] = append(refsByRun[ref.RunID], ref)
	}
	sort.Strings(runIDs)

	var pruned []PrunedRef
	var errs []error
	for _, runID := range runIDs {
		run, found, err := loadTerminalRun(ctx, runID)
		if err != nil {
			errs = append(errs, fmt.Errorf("load terminal Run %q: %w", runID, err))
			continue
		}
		if !found || !store.IsTerminalState(run.State) || canonicalPath(run.GitRoot) != canonicalPath(userRoot) {
			continue
		}

		result, err := inspectTerminalRun(ctx, runner, run)
		if err != nil {
			errs = append(errs, fmt.Errorf("inspect terminal Run %q for pruning: %w", runID, err))
			continue
		}
		if result.State != ReconciliationSafe && result.State != ReconciliationReleased {
			continue
		}

		if result.State == ReconciliationSafe {
			if err := applyTerminalRunWithStore(ctx, runner, runStore, result); err != nil {
				errs = append(errs, err)
				continue
			}
			pruned = append(pruned, PrunedRef{
				RunID:  result.RunID,
				Path:   result.Path,
				Branch: result.Branch,
			})
		}

		prunedTasks, taskErrs := pruneReleasedRunTaskRefs(ctx, runner, userRoot, refsByRun[runID])
		pruned = append(pruned, prunedTasks...)
		errs = append(errs, taskErrs...)
	}
	return pruned, errors.Join(errs...)
}

func pruneReleasedRunTaskRefs(ctx context.Context, runner gitRunner, userRoot string, refs []terminalCandidate) ([]PrunedRef, []error) {
	var pruned []PrunedRef
	var errs []error
	worktrees, err := listRegisteredWorktrees(ctx, runner, userRoot)
	if err != nil {
		return nil, []error{err}
	}
	for _, ref := range refs {
		if ref.TaskID == "" {
			continue
		}
		if registeredPath := registeredBranchPath(worktrees, ref.Branch); registeredPath != "" {
			ref.Path = registeredPath
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
				errs = append(errs, fmt.Errorf("remove terminal Task Worktree %q: %w", ref.Path, err))
				continue
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("stat terminal Task Worktree %q: %w", ref.Path, err))
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
	return pruned, errs
}

func BranchName(runID string) string {
	return runBranchPrefix + strings.TrimSpace(runID)
}

// RunIDFromBranchName returns the Run ID encoded by the canonical Run Branch
// naming contract. An empty Run ID is not a valid Run Branch.
func RunIDFromBranchName(branch string) (string, bool) {
	runID, ok := strings.CutPrefix(strings.TrimSpace(branch), runBranchPrefix)
	if !ok {
		return "", false
	}
	runID = strings.TrimSpace(runID)
	return runID, runID != ""
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

type registeredWorktree struct {
	Path   string
	Branch string
}

type recordedWorktreeState uint8

const (
	recordedWorktreeAbsent recordedWorktreeState = iota
	recordedWorktreePresent
	recordedWorktreeUnregistered
	recordedWorktreeUnsafe
)

func recordedGitRoot(ctx context.Context, runner gitRunner, value string) (string, error) {
	root := strings.TrimSpace(value)
	if root == "" {
		return "", errors.New("inspect terminal Run: recorded Git root is required")
	}
	if root != value || strings.ContainsAny(root, "\r\n\x00") {
		return "", fmt.Errorf("inspect terminal Run: recorded Git root %q is invalid", value)
	}
	if !filepath.IsAbs(root) || filepath.Clean(root) != root {
		return "", fmt.Errorf("inspect terminal Run: recorded Git root %q must be a clean absolute path", value)
	}
	hasSymlink, err := pathContainsSymlink(root)
	if err != nil {
		return "", fmt.Errorf("inspect terminal Run: stat recorded Git root %q: %w", root, err)
	}
	if hasSymlink {
		return "", fmt.Errorf("inspect terminal Run: recorded Git root %q contains a symlink", root)
	}
	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("inspect terminal Run: stat recorded Git root %q: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("inspect terminal Run: recorded Git root %q is not a real directory", root)
	}
	output, err := runner.Run(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("inspect terminal Run: validate recorded Git root %q: %w", root, err)
	}
	topLevel := strings.TrimSpace(output)
	if topLevel == "" || canonicalPath(topLevel) != canonicalPath(root) {
		return "", fmt.Errorf("inspect terminal Run: recorded Git root %q is not the repository root", root)
	}
	return root, nil
}

func listRegisteredWorktrees(ctx context.Context, runner gitRunner, gitRoot string) ([]registeredWorktree, error) {
	output, err := runner.Run(ctx, gitRoot, "worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, fmt.Errorf("list registered Git worktrees: %w", err)
	}
	var worktrees []registeredWorktree
	var current registeredWorktree
	appendCurrent := func() {
		if current.Path != "" {
			worktrees = append(worktrees, current)
		}
		current = registeredWorktree{}
	}
	for _, field := range strings.Split(output, "\x00") {
		if field == "" {
			appendCurrent()
			continue
		}
		if value, ok := strings.CutPrefix(field, "worktree "); ok {
			current.Path = filepath.Clean(value)
			continue
		}
		if value, ok := strings.CutPrefix(field, "branch refs/heads/"); ok {
			current.Branch = value
		}
	}
	appendCurrent()
	return worktrees, nil
}

func recordedWorktree(path string, worktrees []registeredWorktree) (registeredWorktree, recordedWorktreeState) {
	if path == "" {
		return registeredWorktree{}, recordedWorktreeAbsent
	}
	if strings.ContainsAny(path, "\r\n\x00") {
		return registeredWorktree{}, recordedWorktreeUnsafe
	}
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return registeredWorktree{}, recordedWorktreeUnsafe
	}
	hasSymlink, err := pathContainsSymlink(path)
	if errors.Is(err, os.ErrNotExist) {
		return registeredWorktree{}, recordedWorktreeAbsent
	}
	if err != nil || hasSymlink {
		return registeredWorktree{}, recordedWorktreeUnsafe
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return registeredWorktree{}, recordedWorktreeUnsafe
	}
	canonical := canonicalPath(path)
	for _, worktree := range worktrees {
		if canonicalPath(worktree.Path) == canonical {
			return worktree, recordedWorktreePresent
		}
	}
	return registeredWorktree{}, recordedWorktreeUnregistered
}

func pathContainsSymlink(path string) (bool, error) {
	volume := filepath.VolumeName(path)
	current := volume + string(filepath.Separator)
	remainder := strings.TrimPrefix(path, current)
	for _, component := range strings.Split(remainder, string(filepath.Separator)) {
		if component == "" {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
	}
	return false, nil
}

func registeredBranchPath(worktrees []registeredWorktree, branch string) string {
	for _, worktree := range worktrees {
		if worktree.Branch == branch {
			return worktree.Path
		}
	}
	return ""
}

func validLocalBranch(ctx context.Context, runner gitRunner, gitRoot, branch string) bool {
	if branch == "" || strings.TrimSpace(branch) != branch || strings.HasPrefix(branch, "-") {
		return false
	}
	_, err := runner.Run(ctx, gitRoot, "check-ref-format", "refs/heads/"+branch)
	return err == nil
}

func localBranchExists(ctx context.Context, runner gitRunner, gitRoot, branch string) (bool, error) {
	_, err := runner.Run(ctx, gitRoot, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if err == nil {
		return true, nil
	}
	if isGitExitCode(err, 1) {
		return false, nil
	}
	return false, err
}

func resolveUnambiguousLocalBranch(ctx context.Context, runner gitRunner, gitRoot, branch string) (string, error) {
	ambiguous, err := localBranchIsAmbiguous(ctx, runner, gitRoot, branch)
	if err != nil {
		return "", err
	}
	if ambiguous {
		return "", fmt.Errorf("resolve local branch %q: short ref is ambiguous", branch)
	}
	output, err := runner.Run(ctx, gitRoot, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+branch+"^{commit}")
	if err != nil {
		return "", err
	}
	head := strings.TrimSpace(output)
	if head == "" || strings.ContainsAny(head, "\r\n") {
		return "", fmt.Errorf("resolve local branch %q: Git returned an invalid head", branch)
	}
	return head, nil
}

func localBranchIsAmbiguous(ctx context.Context, runner gitRunner, gitRoot, branch string) (bool, error) {
	candidates := map[string]bool{
		"refs/" + branch:                   true,
		"refs/tags/" + branch:              true,
		"refs/heads/" + branch:             true,
		"refs/remotes/" + branch:           true,
		"refs/remotes/" + branch + "/HEAD": true,
	}
	output, err := runner.Run(
		ctx,
		gitRoot,
		"for-each-ref",
		"--format=%(refname)",
		"refs/"+branch,
		"refs/tags/"+branch,
		"refs/heads/"+branch,
		"refs/remotes/"+branch,
	)
	if err != nil {
		return false, err
	}
	count := 0
	for _, line := range strings.Split(output, "\n") {
		if candidates[strings.TrimSpace(line)] {
			count++
		}
	}
	return count != 1, nil
}

func isGitExitCode(err error, code int) bool {
	var gitErr *gitCommandError
	if !errors.As(err, &gitErr) {
		return false
	}
	var exitErr *exec.ExitError
	return errors.As(gitErr.err, &exitErr) && exitErr.ExitCode() == code
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
	paths, err := worktreePathsByBranch(ctx, runner, userRoot)
	if err != nil {
		return "", err
	}
	return paths[targetBranch], nil
}

func worktreePathsByBranch(ctx context.Context, runner gitRunner, userRoot string) (map[string]string, error) {
	output, err := runner.Run(ctx, userRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("list git worktrees: %w", err)
	}
	paths := map[string]string{}
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
			if branch != "" && currentPath != "" {
				paths[branch] = currentPath
			}
		}
	}
	return paths, nil
}

func branchAheadCommits(ctx context.Context, runner gitRunner, userRoot, baseBranch, branch string) (int, error) {
	output, err := runner.Run(ctx, userRoot, "rev-list", "--count", baseBranch+".."+branch)
	if err != nil {
		return 0, fmt.Errorf("count commits on Run Branch %q ahead of %q: %w", branch, baseBranch, err)
	}
	countText := strings.TrimSpace(output)
	if countText == "" {
		return 0, fmt.Errorf("count commits on Run Branch %q ahead of %q: git returned an empty count", branch, baseBranch)
	}
	count, err := strconv.Atoi(countText)
	if err != nil {
		return 0, fmt.Errorf("count commits on Run Branch %q ahead of %q: parse count %q: %w", branch, baseBranch, countText, err)
	}
	return count, nil
}

func branchCanFastForward(ctx context.Context, runner gitRunner, userRoot, baseBranch, branch string) (bool, error) {
	if _, err := runner.Run(ctx, userRoot, "merge-base", "--is-ancestor", baseBranch, branch); err != nil {
		if isAncestryMiss(err) {
			return false, nil
		}
		return false, fmt.Errorf("check whether %q can fast-forward to %q: %w", baseBranch, branch, err)
	}
	return true, nil
}

func isTaskBranch(branch string) bool {
	suffix := strings.TrimPrefix(strings.TrimSpace(branch), runBranchPrefix)
	idx := strings.LastIndex(suffix, "-")
	return idx > 0 && idx < len(suffix)-1 && strings.HasPrefix(suffix[idx+1:], "task_")
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
