package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

type commandDeliveryWorkflow struct {
	store  *store.Store
	loaded roundconfig.Loaded
	git    preflight.GitRunner
	gh     delivery.GitHubCLI
	log    io.Writer
}

var _ delivery.PrerequisiteReader = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemHistory = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemRecovery = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemWorkspace = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemRevalidator = (*commandDeliveryWorkflow)(nil)
var _ delivery.MergeObserver = (*commandDeliveryWorkflow)(nil)

const deliveryBranchPrefix = "roundfix/deliver-"

func newCommandDeliveryEngine(runStore *store.Store, loaded roundconfig.Loaded) deliveryEngine {
	github := delivery.NewGitHubCLI(loaded.GitRoot)
	workflow := &commandDeliveryWorkflow{
		store:  runStore,
		loaded: loaded,
		git:    preflight.ExecGitRunner{},
		gh:     github,
	}
	return delivery.NewEngine(runStore, delivery.EngineDependencies{
		Merges:        workflow,
		Corrections:   workflow,
		Conflicts:     workflow,
		Workspace:     workflow,
		Runner:        workflow,
		Reviewer:      workflow,
		Archiver:      workflow,
		Gate:          workflow,
		Authorizer:    workflow,
		Publication:   workflow,
		PullRequests:  github,
		Checks:        github,
		Recovery:      workflow,
		History:       workflow,
		Revalidator:   workflow,
		Prerequisites: workflow,
		Log:           os.Stderr,
	})
}

func (workflow *commandDeliveryWorkflow) ObserveMerge(ctx context.Context, gitRoot string, item store.DeliveryQueueItem) (delivery.MergeObservation, error) {
	var pullRequest delivery.PullRequest
	if item.PullRequestNumber != "" {
		github := workflow.gh
		github.WorkDir = gitRoot
		var err error
		pullRequest, err = github.ViewPullRequest(ctx, item.PullRequestNumber)
		if err != nil {
			return delivery.MergeObservation{}, fmt.Errorf("read recorded pull request #%s: %w", item.PullRequestNumber, err)
		}
		if pullRequest.Merged() && pullRequest.HeadBranch == item.Branch && pullRequest.HeadSHA != "" && pullRequest.MergeCommit != "" {
			return delivery.MergeObservation{
				Merged: true, Head: pullRequest.HeadSHA, MergeCommit: pullRequest.MergeCommit,
				Evidence: "pull request #" + item.PullRequestNumber,
			}, nil
		}
	}

	head := ""
	if len(item.CandidateCommits) != 0 {
		head = strings.TrimSpace(item.CandidateCommits[len(item.CandidateCommits)-1])
	}
	if item.Branch != "" {
		runner := workflow.git
		if runner == nil {
			runner = preflight.ExecGitRunner{}
		}
		exists, err := localItemBranchExists(ctx, runner, gitRoot, item.Branch)
		if err != nil {
			return delivery.MergeObservation{}, fmt.Errorf("inspect item branch %q: %w", item.Branch, err)
		}
		if exists {
			output, err := runner.RunGit(ctx, gitRoot, "rev-parse", "--verify", "--end-of-options", "refs/heads/"+item.Branch+"^{commit}")
			if err != nil {
				return delivery.MergeObservation{}, fmt.Errorf("read item branch %q head: %w", item.Branch, err)
			}
			head = strings.TrimSpace(output)
		}
	}
	if head != "" {
		if proof, proven := runworktree.ProveDelivery(ctx, gitRoot, item.SpecSlug, head); proven {
			return delivery.MergeObservation{
				Merged: true, Head: head, MergeCommit: proof.DeliveryCommit,
				Evidence: fmt.Sprintf("Spec archived on default branch %q", proof.DefaultBranch),
			}, nil
		}
	}
	return delivery.MergeObservation{
		ClosedUnmerged: strings.EqualFold(pullRequest.State, "CLOSED") && !pullRequest.Merged(),
	}, nil
}

func (workflow *commandDeliveryWorkflow) InspectItem(
	ctx context.Context,
	workDir string,
	specSlug string,
) (delivery.ItemState, error) {
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, workDir)
	if err != nil {
		return delivery.ItemState{}, fmt.Errorf("resolve item Specs Root: %w", err)
	}
	head, err := workflow.git.RunGit(ctx, workDir, "rev-parse", "HEAD")
	if err != nil {
		return delivery.ItemState{}, fmt.Errorf("read item head: %w", err)
	}
	state := delivery.ItemState{Head: strings.TrimSpace(head)}

	specDir := filepath.Join(resolvedSpecsRoot.Path, specSlug)
	if _, err := os.Stat(specDir); err == nil {
		graph, err := spec.Load(resolvedSpecsRoot.Path, specSlug)
		if err != nil {
			return delivery.ItemState{}, fmt.Errorf("load item Spec %q: %w", specSlug, err)
		}
		for _, task := range graph.Tasks {
			if task.Status != spec.StatusCompleted {
				state.UnfinishedTasks = append(state.UnfinishedTasks, task.ID)
			}
		}
		return state, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return delivery.ItemState{}, fmt.Errorf("inspect item Spec %q: %w", specSlug, err)
	}

	_, archiveDestination, err := workflow.archivePaths(workDir, specSlug)
	if err != nil {
		return delivery.ItemState{}, fmt.Errorf("resolve item archive path: %w", err)
	}
	if _, err := os.Stat(filepath.Join(workDir, filepath.FromSlash(archiveDestination))); err == nil {
		content, err := os.ReadFile(filepath.Join(workDir, filepath.FromSlash(archiveDestination), "_prd.md"))
		if err != nil {
			return delivery.ItemState{}, fmt.Errorf("read archived item PRD: %w", err)
		}
		frontmatter, _, ok := splitArchivePRD(content)
		if !ok {
			return delivery.ItemState{}, errors.New("read archived item PRD: invalid frontmatter")
		}
		var metadata struct {
			QAOverride bool `yaml:"qa_override"`
		}
		if err := yaml.Unmarshal(frontmatter, &metadata); err != nil {
			return delivery.ItemState{}, fmt.Errorf("read archived item QA override: %w", err)
		}
		state.Archived = true
		state.QAOverride = metadata.QAOverride
		return state, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return delivery.ItemState{}, fmt.Errorf("inspect item Spec %q: active and archived Spec folders are missing", specSlug)
	} else {
		return delivery.ItemState{}, fmt.Errorf("inspect archived item Spec %q: %w", specSlug, err)
	}
}

type deliveryCarryForwardRefusal struct {
	runID       string
	workDir     string
	specSlug    string
	carriedRuns []string
	reason      string
	amendments  []string
}

func (err deliveryCarryForwardRefusal) Error() string {
	message := err.reason
	if len(err.carriedRuns) != 0 {
		message = fmt.Sprintf(
			"carried forward from Run(s) %s before Run %q refused: %s",
			strings.Join(err.carriedRuns, ", "),
			err.runID,
			err.reason,
		)
	}
	if len(err.amendments) != 0 {
		message += "; amended by " + strings.Join(err.amendments, ", ")
	}
	return message
}

func (err deliveryCarryForwardRefusal) NextAction() string {
	if len(err.amendments) != 0 {
		quotedAmendments := make([]string, 0, len(err.amendments))
		for _, amendment := range err.amendments {
			quotedAmendments = append(quotedAmendments, posixSingleQuote(amendment))
		}
		return strings.Join([]string{
			fmt.Sprintf(
				"git -C %s branch %s HEAD",
				posixSingleQuote(err.workDir),
				posixSingleQuote("roundfix-amended-"+err.runID),
			),
			fmt.Sprintf(
				"git -C %s reset --hard %s",
				posixSingleQuote(err.workDir),
				posixSingleQuote(err.amendments[0]+"^"),
			),
			fmt.Sprintf(
				"(cd %s && roundfix reconcile %s --carry-forward)",
				posixSingleQuote(err.workDir),
				posixSingleQuote(err.runID),
			),
			fmt.Sprintf(
				"git -C %s cherry-pick %s",
				posixSingleQuote(err.workDir),
				strings.Join(quotedAmendments, " "),
			),
			fmt.Sprintf("roundfix deliver retry %s", posixSingleQuote(err.specSlug)),
		}, "\n")
	}
	return fmt.Sprintf(
		"run `roundfix reconcile %s --carry-forward` in item worktree %q, then run `roundfix deliver retry %s`",
		err.runID,
		err.workDir,
		err.specSlug,
	)
}

func posixSingleQuote(argument string) string {
	return "'" + strings.ReplaceAll(argument, "'", "'\"'\"'") + "'"
}

func (workflow *commandDeliveryWorkflow) CarryForward(
	ctx context.Context,
	workDir string,
	specSlug string,
	branch string,
	runID string,
) (delivery.CarryForwardResult, error) {
	runs, err := workflow.deliveryCarryForwardRuns(ctx, specSlug, branch, runID)
	if err != nil {
		return delivery.CarryForwardResult{}, err
	}
	if len(runs) == 0 {
		return delivery.CarryForwardResult{}, nil
	}
	result := delivery.CarryForwardResult{RunID: runs[0].ID}
	var repository string
	var resolvedSpecsRoot roundconfig.SpecsRoot
	itemResolved := false
	resolveItem := func() error {
		if itemResolved {
			return nil
		}
		var err error
		repository, err = filepath.EvalSymlinks(workDir)
		if err != nil {
			return fmt.Errorf("resolve item worktree %q: %w", workDir, err)
		}
		resolvedSpecsRoot, err = roundconfig.ResolveSpecsRoot(workflow.loaded, repository)
		if err != nil {
			return fmt.Errorf("resolve item Specs Root: %w", err)
		}
		itemResolved = true
		return nil
	}

	for _, run := range runs {
		if !slices.Contains(runworktree.CarryForwardAcceptedOutcomes(), run.State) {
			continue
		}
		_, taskEvidence, err := loadReconcileTaskCoverage(ctx, workflow.store, []store.Run{run})
		if err != nil {
			return delivery.CarryForwardResult{}, fmt.Errorf("load Run %q Task coverage: %w", run.ID, err)
		}
		settled := taskEvidence[run.ID]
		settledCompleted := make(map[string]bool, len(settled))
		for taskID, evidence := range settled {
			if evidence.settledCompleted {
				settledCompleted[taskID] = true
			}
		}
		if len(settledCompleted) == 0 {
			continue
		}
		if err := resolveItem(); err != nil {
			return delivery.CarryForwardResult{}, err
		}
		graph, err := spec.Load(resolvedSpecsRoot.Path, specSlug)
		if err != nil {
			return delivery.CarryForwardResult{}, fmt.Errorf("load item Spec %q: %w", specSlug, err)
		}
		allCompleted := true
		matchedSettledTasks := 0
		for _, task := range graph.Tasks {
			if !settledCompleted[task.ID] {
				continue
			}
			matchedSettledTasks++
			if task.Status != spec.StatusCompleted {
				allCompleted = false
				break
			}
		}
		if allCompleted && matchedSettledTasks == len(settledCompleted) {
			continue
		}

		refuse := func(reason string, amendments []string) (delivery.CarryForwardResult, error) {
			carriedRuns := make([]string, 0, len(result.Runs))
			for _, carried := range result.Runs {
				carriedRuns = append(carriedRuns, carried.RunID)
			}
			return delivery.CarryForwardResult{}, deliveryCarryForwardRefusal{
				runID:       run.ID,
				workDir:     repository,
				specSlug:    specSlug,
				carriedRuns: carriedRuns,
				reason:      reason,
				amendments:  amendments,
			}
		}
		present, err := carryForwardRunWorktreePresent(run)
		if err != nil {
			return delivery.CarryForwardResult{}, err
		}
		if !present {
			return refuse(fmt.Sprintf("carry-forward Run %q Worktree is gone", run.ID), nil)
		}
		if resolvedSpecsRoot.External {
			return refuse(fmt.Sprintf("carry-forward Specs Root %q is external to item worktree %q", resolvedSpecsRoot.Path, repository), nil)
		}

		candidates, err := inspectCarryForwards(ctx, repository, resolvedSpecsRoot, reconcileRunSelection{
			selected:     []store.Run{run},
			taskEvidence: taskEvidence,
		})
		if err != nil {
			return delivery.CarryForwardResult{}, fmt.Errorf("inspect Run %q carry-forward: %w", run.ID, err)
		}
		if reason := carryForwardRefusalReason(candidates); reason != "" {
			amendments, ok, err := carryForwardAmendments(ctx, repository, run, candidates)
			if err != nil {
				return delivery.CarryForwardResult{}, fmt.Errorf("inspect Run %q carry-forward amendments: %w", run.ID, err)
			}
			if !ok {
				amendments = nil
			}
			return refuse(reason, amendments)
		}
		carriedTasks := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			if candidate.Action == carryForwardReadyAction {
				carriedTasks = append(carriedTasks, candidate.TaskID)
			}
		}
		if err := applyCarryForwards(ctx, repository, candidates); err != nil {
			return delivery.CarryForwardResult{}, fmt.Errorf("apply Run %q carry-forward: %w", run.ID, err)
		}
		if len(carriedTasks) > 0 {
			result.Runs = append(result.Runs, delivery.CarriedRun{RunID: run.ID, Carried: carriedTasks})
		}
	}
	return result, nil
}

func (workflow *commandDeliveryWorkflow) deliveryCarryForwardRuns(
	ctx context.Context,
	specSlug string,
	branch string,
	runID string,
) ([]store.Run, error) {
	var recorded store.Run
	hasRecorded := false
	runID = strings.TrimSpace(runID)
	if runID != "" {
		run, found, err := workflow.store.Run(ctx, runID)
		if err != nil {
			return nil, fmt.Errorf("read recorded Run %q: %w", runID, err)
		}
		if !found || run.Kind != store.KindImplement || strings.TrimSpace(run.SpecSlug) != strings.TrimSpace(specSlug) {
			return nil, fmt.Errorf("recorded Run %q is not an Implement Run of Spec %q", runID, specSlug)
		}
		recorded = run
		hasRecorded = true
	}

	runs, err := workflow.store.ListRuns(ctx, store.ListRunsQuery{
		GitRoot: workflow.loaded.GitRoot,
		States:  store.StatesTerminal,
	})
	if err != nil {
		return nil, fmt.Errorf("list Implement Runs for Spec %q: %w", specSlug, err)
	}
	branch = strings.TrimSpace(branch)
	selected := make([]store.Run, 0, len(runs)+1)
	seen := make(map[string]bool, len(runs)+1)
	if hasRecorded {
		selected = append(selected, recorded)
		seen[recorded.ID] = true
	}
	for _, run := range runs {
		if run.Kind == store.KindImplement &&
			strings.TrimSpace(run.SpecSlug) == strings.TrimSpace(specSlug) &&
			strings.TrimSpace(run.LocalBranch) == branch &&
			!seen[run.ID] {
			selected = append(selected, run)
			seen[run.ID] = true
		}
	}
	sort.Slice(selected, func(left, right int) bool {
		if selected[left].CreatedAt.Equal(selected[right].CreatedAt) {
			return selected[left].ID > selected[right].ID
		}
		return selected[left].CreatedAt.After(selected[right].CreatedAt)
	})
	return selected, nil
}

func (workflow *commandDeliveryWorkflow) CreateItemBranch(ctx context.Context, gitRoot, specSlug string) (string, string, error) {
	// A recorded workspace is durable recovery state, even when the remote
	// is unavailable. Only a new item selects among branches holding work.
	queue, _, err := workflow.store.DeliveryQueue(ctx, gitRoot)
	if err != nil {
		return "", "", fmt.Errorf("read item workspace: %w", err)
	}
	branch := ""
	for _, item := range queue.Items {
		if item.SpecSlug == specSlug {
			branch = item.Branch
			break
		}
	}
	base := ""
	if branch == "" {
		defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, "", workflow.git)
		if defaultBranch.Source == preflight.DefaultBranchUndetermined {
			return "", "", errors.New("create item branch: repository default branch is unknown")
		}
		remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
		if remote == "" {
			remote = "origin"
		}
		if _, err := workflow.git.RunGit(ctx, gitRoot, "fetch", remote, defaultBranch.Name); err != nil {
			return "", "", fmt.Errorf("refresh default branch %q: %w", defaultBranch.Name, err)
		}
		base = remote + "/" + defaultBranch.Name
		branches, err := existingItemBranches(ctx, workflow.git, gitRoot, base, specSlug)
		if err != nil {
			return "", "", err
		}
		if len(branches) > 1 {
			return "", "", ambiguousItemBranches(specSlug, base, branches)
		}
		if len(branches) == 1 {
			branch = branches[0]
		} else {
			branch, err = newDeliveryBranch(specSlug)
			if err != nil {
				return "", "", err
			}
		}
	}
	ref, err := runworktree.ItemRefFor(gitRoot, workflow.loaded.Config.Worktree.Location, branch)
	if err != nil {
		return "", "", err
	}
	branch, itemWorktree, provisioned, err := workflow.store.RecordDeliveryQueueItemWorktree(ctx, gitRoot, specSlug, branch, ref.Path)
	if err != nil {
		return "", "", fmt.Errorf("record item branch and worktree: %w", err)
	}
	ref, err = runworktree.ItemRefFor(gitRoot, workflow.loaded.Config.Worktree.Location, branch)
	if err != nil {
		return "", "", err
	}
	if ref.Path != itemWorktree {
		return "", "", fmt.Errorf("recorded item worktree %q does not match derived path %q", itemWorktree, ref.Path)
	}

	exists, err := localItemBranchExists(ctx, workflow.git, gitRoot, branch)
	if err != nil {
		return "", "", fmt.Errorf("inspect item branch %q: %w", branch, err)
	}
	if exists {
		if err := workflow.useAndProvisionItem(ctx, gitRoot, specSlug, ref, provisioned); err != nil {
			return "", "", fmt.Errorf("use recorded item worktree %q: %w", itemWorktree, err)
		}
		return branch, itemWorktree, nil
	}
	if provisioned {
		if err := workflow.store.SetDeliveryQueueItemWorktreeProvisioned(ctx, gitRoot, specSlug, false); err != nil {
			return "", "", err
		}
	}
	if base == "" {
		defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, "", workflow.git)
		if defaultBranch.Source == preflight.DefaultBranchUndetermined {
			return "", "", errors.New("create item branch: repository default branch is unknown")
		}
		remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
		if remote == "" {
			remote = "origin"
		}
		if _, err := workflow.git.RunGit(ctx, gitRoot, "fetch", remote, defaultBranch.Name); err != nil {
			return "", "", fmt.Errorf("refresh default branch %q: %w", defaultBranch.Name, err)
		}
		base = remote + "/" + defaultBranch.Name
	}
	if err := runworktree.CreateItem(ctx, ref, runworktree.ItemCreateOptions{
		HeadSHA:  base,
		CopyList: workflow.loaded.Config.Worktree.Copy,
		Bootstrap: runworktree.BootstrapSpec{
			Command: workflow.loaded.Config.Worktree.Bootstrap,
			Timeout: workflow.loaded.Config.Worktree.BootstrapTimeout,
		},
		BootstrapOutput: os.Stderr,
	}); err != nil {
		return "", "", fmt.Errorf("create item branch %q in worktree %q: %w", branch, itemWorktree, err)
	}
	if err := workflow.store.SetDeliveryQueueItemWorktreeProvisioned(ctx, gitRoot, specSlug, true); err != nil {
		return "", "", err
	}
	return branch, itemWorktree, nil
}

// existingItemBranches reads only local refs; callers choose when to refresh
// the default branch. A branch holds work when it has commits absent from base.
func existingItemBranches(ctx context.Context, git preflight.GitRunner, gitRoot, base, specSlug string) ([]string, error) {
	prefix := deliveryBranchPrefix + specSlug + "-"
	refs, err := git.RunGit(ctx, gitRoot, "for-each-ref", "--format=%(refname:strip=2)", "refs/heads/"+prefix+"*")
	if err != nil {
		return nil, fmt.Errorf("list item branches for %q: %w", specSlug, err)
	}
	pattern := regexp.MustCompile("^" + regexp.QuoteMeta(prefix) + "[0-9a-f]{16}$")
	var branches []string
	for _, branch := range strings.Fields(refs) {
		if !pattern.MatchString(branch) {
			continue
		}
		if base == "" {
			return nil, errors.New("inspect item branches: repository default branch is unknown")
		}
		output, err := git.RunGit(ctx, gitRoot, "rev-list", "--count", base+".."+branch)
		if err != nil {
			return nil, fmt.Errorf("count item branch %q commits beyond %q: %w", branch, base, err)
		}
		count, err := strconv.ParseUint(strings.TrimSpace(output), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("read item branch %q commit count: %w", branch, err)
		}
		if count > 0 {
			branches = append(branches, branch)
		}
	}
	sort.Strings(branches)
	return branches, nil
}

func ambiguousItemBranches(specSlug, base string, branches []string) error {
	return fmt.Errorf("Spec %q has %d item branches with commits %s lacks: %s; delete every branch but the one to continue, then run roundfix deliver start again", specSlug, len(branches), base, strings.Join(branches, ", "))
}

func newDeliveryBranch(specSlug string) (string, error) {
	specSlug = strings.TrimSpace(specSlug)
	if specSlug == "" {
		return "", errors.New("create item branch: Spec slug is required")
	}
	var suffix [8]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", fmt.Errorf("create item branch: generate branch suffix: %w", err)
	}
	return deliveryBranchPrefix + specSlug + "-" + hex.EncodeToString(suffix[:]), nil
}

func localItemBranchExists(ctx context.Context, runner preflight.GitRunner, gitRoot, branch string) (bool, error) {
	_, err := runner.RunGit(ctx, gitRoot, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}

func (workflow *commandDeliveryWorkflow) UseItemBranch(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	branch string,
	itemWorktree string,
	provisioned bool,
) (string, error) {
	specSlug = strings.TrimSpace(specSlug)
	if specSlug == "" {
		return "", errors.New("use item branch: Spec slug is required")
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return "", errors.New("use item branch: branch is required")
	}
	itemWorktree = strings.TrimSpace(itemWorktree)
	if itemWorktree == "" {
		return "", errors.New("use item branch: worktree is required")
	}
	err := workflow.useAndProvisionItem(ctx, gitRoot, specSlug, runworktree.ItemRef{
		Path:     itemWorktree,
		Branch:   branch,
		UserRoot: gitRoot,
	}, provisioned)
	if errors.Is(err, runworktree.ErrItemBranchMissing) {
		return "", delivery.ErrItemWorktreeMissing
	}
	if err != nil {
		return "", fmt.Errorf("use item branch: %w", err)
	}
	return itemWorktree, nil
}

func (workflow *commandDeliveryWorkflow) useAndProvisionItem(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	ref runworktree.ItemRef,
	provisioned bool,
) error {
	if _, err := os.Stat(ref.Path); errors.Is(err, os.ErrNotExist) {
		if provisioned {
			if err := workflow.store.SetDeliveryQueueItemWorktreeProvisioned(ctx, gitRoot, specSlug, false); err != nil {
				return err
			}
			provisioned = false
		}
	} else if err != nil {
		return fmt.Errorf("inspect recorded item worktree %q: %w", ref.Path, err)
	}
	if err := runworktree.UseItem(ctx, ref); err != nil {
		return err
	}
	if provisioned {
		return nil
	}
	if err := runworktree.ProvisionItem(ctx, ref, runworktree.ItemProvisionOptions{
		CopyList: workflow.loaded.Config.Worktree.Copy,
		Bootstrap: runworktree.BootstrapSpec{
			Command: workflow.loaded.Config.Worktree.Bootstrap,
			Timeout: workflow.loaded.Config.Worktree.BootstrapTimeout,
		},
		BootstrapOutput: os.Stderr,
	}); err != nil {
		return err
	}
	return workflow.store.SetDeliveryQueueItemWorktreeProvisioned(ctx, gitRoot, specSlug, true)
}

func (workflow *commandDeliveryWorkflow) RemoveItemBranch(ctx context.Context, gitRoot, branch, itemWorktree string) error {
	gitRoot = strings.TrimSpace(gitRoot)
	branch = strings.TrimSpace(branch)
	itemWorktree = strings.TrimSpace(itemWorktree)
	if itemWorktree == "" {
		if branch == "" {
			return nil
		}
		if err := runworktree.CleanupItemBranch(ctx, gitRoot, branch); err != nil {
			return fmt.Errorf("remove item branch without recorded worktree: %w", err)
		}
		return nil
	}
	ref, err := runworktree.ItemRefFor(gitRoot, workflow.loaded.Config.Worktree.Location, branch)
	if err != nil {
		return fmt.Errorf("derive item worktree for cleanup: %w", err)
	}
	if ref.Path != itemWorktree {
		return fmt.Errorf("recorded item worktree %q does not match derived path %q", itemWorktree, ref.Path)
	}
	if err := runworktree.CleanupItem(ctx, ref); err != nil {
		return fmt.Errorf("remove item worktree and branch: %w", err)
	}
	return nil
}

func (workflow *commandDeliveryWorkflow) ReleaseMergedRuns(
	ctx context.Context,
	gitRoot string,
	item store.DeliveryQueueItem,
) error {
	mergeCommit := strings.TrimSpace(item.MergeCommit)
	if mergeCommit == "" {
		return fmt.Errorf("release merged Spec %q Runs: merge commit is required", item.SpecSlug)
	}
	if len(item.CandidateCommits) == 0 {
		return fmt.Errorf("release merged Spec %q Runs: candidate head is required", item.SpecSlug)
	}
	head := strings.TrimSpace(item.CandidateCommits[len(item.CandidateCommits)-1])
	if head == "" {
		return fmt.Errorf("release merged Spec %q Runs: candidate head is required", item.SpecSlug)
	}
	resolvedMergeCommit, resolvedHead, err := workflow.resolveMergedReleaseEvidence(
		ctx,
		gitRoot,
		mergeCommit,
		head,
	)
	if err != nil {
		return fmt.Errorf("release merged Spec %q Runs: %w", item.SpecSlug, err)
	}
	merged := runworktree.MergedHead{
		SpecSlug:     strings.TrimSpace(item.SpecSlug),
		TargetBranch: strings.TrimSpace(item.Branch),
		Head:         resolvedHead,
		MergeCommit:  resolvedMergeCommit,
		PullRequest:  strings.TrimSpace(item.PullRequestNumber),
	}
	return releaseMergedSpecRuns(ctx, workflow.store, gitRoot, merged)
}

func (workflow *commandDeliveryWorkflow) resolveMergedReleaseEvidence(
	ctx context.Context,
	gitRoot string,
	mergeCommit string,
	candidateHead string,
) (string, string, error) {
	runner := workflow.git
	if runner == nil {
		runner = preflight.ExecGitRunner{}
	}
	resolveCommit := func(label string, value string) (string, error) {
		resolved, err := runner.RunGit(
			ctx,
			gitRoot,
			"rev-parse",
			"--verify",
			"--end-of-options",
			value+"^{commit}",
		)
		if err != nil {
			return "", fmt.Errorf("%s %q does not resolve to a commit", label, value)
		}
		return strings.TrimSpace(resolved), nil
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	remoteExists := false
	defaultBranch := preflight.DefaultBranch{}
	if _, err := runner.RunGit(ctx, gitRoot, "remote", "get-url", remote); err == nil {
		remoteExists = true
		currentBranch, _ := runner.RunGit(ctx, gitRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
		defaultBranch = preflight.DetectDefaultBranch(ctx, gitRoot, strings.TrimSpace(currentBranch), runner)
		if defaultBranch.Source == preflight.DefaultBranchUndetermined {
			for _, branch := range []string{"main", "master"} {
				if _, branchErr := runner.RunGit(
					ctx,
					gitRoot,
					"show-ref",
					"--verify",
					"--quiet",
					"refs/heads/"+branch,
				); branchErr == nil {
					defaultBranch = preflight.DefaultBranch{Name: branch, Source: preflight.DefaultBranchFromNameMatch}
					break
				}
			}
			if defaultBranch.Source == preflight.DefaultBranchUndetermined {
				return "", "", errors.New("default branch is unknown")
			}
		}
		if _, err := runner.RunGit(ctx, gitRoot, "fetch", remote, defaultBranch.Name); err != nil {
			return "", "", fmt.Errorf("refresh default branch %q: %w", defaultBranch.Name, err)
		}
	}

	resolvedMergeCommit, err := resolveCommit("merge commit", mergeCommit)
	if err != nil {
		return "", "", err
	}
	resolvedHead, err := resolveCommit("candidate head", candidateHead)
	if err != nil {
		return "", "", err
	}

	defaultHead := ""
	if remoteExists {
		defaultHead, err = resolveCommit(
			fmt.Sprintf("default branch %q", defaultBranch.Name),
			"refs/remotes/"+remote+"/"+defaultBranch.Name,
		)
		if err != nil {
			return "", "", err
		}
	} else {
		currentBranch, _ := runner.RunGit(ctx, gitRoot, "symbolic-ref", "--quiet", "--short", "HEAD")
		defaultBranch = preflight.DetectDefaultBranch(ctx, gitRoot, strings.TrimSpace(currentBranch), runner)
		if defaultBranch.Source == preflight.DefaultBranchUndetermined {
			for _, branch := range []string{"main", "master"} {
				candidate, candidateErr := resolveCommit(fmt.Sprintf("default branch %q", branch), "refs/heads/"+branch)
				if candidateErr == nil {
					defaultBranch = preflight.DefaultBranch{Name: branch, Source: preflight.DefaultBranchFromNameMatch}
					defaultHead = candidate
					break
				}
			}
			if defaultHead == "" {
				return "", "", errors.New("default branch is unknown")
			}
		}
		var defaultErr error
		if defaultHead == "" {
			defaultHead, defaultErr = resolveCommit(fmt.Sprintf("default branch %q", defaultBranch.Name), "refs/heads/"+defaultBranch.Name)
		}
		if defaultErr != nil && defaultBranch.Source == preflight.DefaultBranchFromOriginHead {
			defaultHead, defaultErr = resolveCommit(
				fmt.Sprintf("default branch %q", defaultBranch.Name),
				"refs/remotes/origin/"+defaultBranch.Name,
			)
		}
		if defaultErr != nil {
			return "", "", defaultErr
		}
	}
	if _, err := runner.RunGit(
		ctx,
		gitRoot,
		"merge-base",
		"--is-ancestor",
		resolvedMergeCommit,
		defaultHead,
	); err != nil {
		return "", "", fmt.Errorf(
			"merge commit %q is not on default branch %q",
			mergeCommit,
			defaultBranch.Name,
		)
	}

	if _, err := runner.RunGit(
		ctx,
		gitRoot,
		"merge-base",
		"--is-ancestor",
		resolvedHead,
		resolvedMergeCommit,
	); err != nil {
		candidateTree, candidateTreeErr := runner.RunGit(ctx, gitRoot, "rev-parse", resolvedHead+"^{tree}")
		mergeTree, mergeTreeErr := runner.RunGit(ctx, gitRoot, "rev-parse", resolvedMergeCommit+"^{tree}")
		if candidateTreeErr != nil || mergeTreeErr != nil || strings.TrimSpace(candidateTree) != strings.TrimSpace(mergeTree) {
			mergedCandidateTree, mergedCandidateTreeErr := runner.RunGit(
				ctx,
				gitRoot,
				"merge-tree",
				"--write-tree",
				resolvedMergeCommit+"^1",
				resolvedHead,
			)
			if mergeTreeErr != nil || mergedCandidateTreeErr != nil || strings.TrimSpace(mergedCandidateTree) != strings.TrimSpace(mergeTree) {
				return "", "", fmt.Errorf(
					"candidate head %q is not represented by merge commit %q",
					candidateHead,
					mergeCommit,
				)
			}
		}
	}

	return resolvedMergeCommit, resolvedHead, nil
}

func (workflow *commandDeliveryWorkflow) RunSpec(ctx context.Context, gitRoot, specSlug string) (delivery.RunResult, error) {
	beforeRun, beforeFound, err := workflow.latestImplementRun(ctx, workflow.loaded.GitRoot, specSlug)
	if err != nil {
		return delivery.RunResult{}, err
	}
	executable, err := workflow.stepExecutable(ctx, gitRoot, specSlug, "implement")
	if err != nil {
		return delivery.RunResult{}, err
	}
	result, err := workflow.runRoundfix(ctx, executable, gitRoot, "implement", "--spec", specSlug)
	if err != nil {
		return delivery.RunResult{}, fmt.Errorf("start Implement executor: %w", err)
	}
	afterRun, afterFound, err := workflow.latestImplementRun(ctx, workflow.loaded.GitRoot, specSlug)
	if err != nil {
		return delivery.RunResult{}, err
	}
	candidateHead := ""
	if result.exitCode == exitOK {
		head, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", "HEAD")
		if err != nil {
			return delivery.RunResult{}, fmt.Errorf("read Implement candidate head: %w", err)
		}
		candidateHead = head
	}
	var before, after *store.Run
	if beforeFound {
		before = &beforeRun
	}
	if afterFound {
		after = &afterRun
	}
	return workflow.runResult(ctx, gitRoot, specSlug, result, candidateHead, before, after)
}

// runResult reads QA evidence from the Run Branch, which need not have integrated
// into the item worktree after an unresolved Run.
func (workflow *commandDeliveryWorkflow) runResult(ctx context.Context, gitRoot, specSlug string, command roundfixCommandResult, candidateHead string, before, after *store.Run) (delivery.RunResult, error) {
	result, err := deliveryRunResult(command, candidateHead, before, after)
	if err == nil && result.Outcome == delivery.RunOutcomeUnresolved {
		if workflow.store != nil {
			events, readErr := workflow.store.RunEventsOfKinds(ctx, result.RunID, runevent.KindDaemonTask)
			if readErr != nil {
				return result, fmt.Errorf("read runtime infrastructure Run Events: %w", readErr)
			}
			for _, event := range events {
				var payload struct {
					Phase    string `json:"phase"`
					Recovery string `json:"recovery"`
					ScopeID  string `json:"scope_id"`
					Step     string `json:"step"`
				}
				if json.Unmarshal(event.Event.Payload, &payload) == nil && payload.Phase == "rollout_lost" && payload.Recovery == "exhausted" {
					result.RuntimeInfrastructure = payload.ScopeID + " lost its rollout at " + payload.Step
				}
			}
		}
		result.QAEnvironmentPartial = workflow.qaEnvironmentPartial(ctx, gitRoot, specSlug, store.RunBranchPrefix+result.RunID)
	}
	return result, err
}

func (workflow *commandDeliveryWorkflow) qaEnvironmentPartial(ctx context.Context, gitRoot, specSlug, branch string) bool {
	source, _, err := workflow.archivePaths(gitRoot, specSlug)
	if err != nil {
		return false
	}
	qaDir := source + "/qa/"
	paths, err := workflow.git.RunGit(ctx, gitRoot, "ls-tree", "-r", "--name-only", "-z", branch, "--", qaDir)
	if err != nil {
		return false
	}
	var reports []string
	for _, path := range strings.Split(paths, "\x00") {
		name := filepath.Base(path)
		if strings.HasPrefix(path, qaDir) && strings.HasPrefix(name, "qa-report-") && strings.HasSuffix(name, ".md") {
			reports = append(reports, path)
		}
	}
	newest, err := spec.NewestQAReportFromPaths(reports)
	if err != nil {
		return false
	}
	content, err := workflow.git.RunGit(ctx, gitRoot, "show", branch+":"+newest)
	if err != nil {
		return false
	}
	// Use the shared report reader so row counts, verdicts and pre-PR rows
	// retain exactly the QA Report contract.
	temporary, err := os.CreateTemp("", "roundfix-delivery-qa-*.md")
	if err != nil {
		return false
	}
	defer os.Remove(temporary.Name())
	_, writeErr := temporary.WriteString(content)
	closeErr := temporary.Close()
	if writeErr != nil || closeErr != nil {
		return false
	}
	report, err := spec.ReadQAReportFile(temporary.Name())
	return err == nil && report.Verdict == spec.VerdictPartial && report.RowsBlockedFinding == 0 && report.EnvironmentRowsNeedingOverride() > 0
}

func (workflow *commandDeliveryWorkflow) Descends(ctx context.Context, workDir, ancestor, head string) (bool, error) {
	_, err := workflow.git.RunGit(ctx, workDir, "merge-base", "--is-ancestor", ancestor, head)
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("inspect archived head ancestry: %w", err)
}

func (workflow *commandDeliveryWorkflow) RunStart(ctx context.Context, gitRoot, runID string) (string, error) {
	run, found, err := workflow.store.Run(ctx, runID)
	if err != nil {
		return "", fmt.Errorf("read Run %q: %w", runID, err)
	}
	if !found || run.Kind != store.KindImplement || strings.TrimSpace(run.HeadSHA) == "" {
		return "", fmt.Errorf("Run %q has no Implement start head for repository %q", runID, gitRoot)
	}
	repositoryRoot, err := roundconfig.RepositoryRoot(gitRoot)
	if err != nil {
		return "", fmt.Errorf("resolve queue repository for Run %q: %w", runID, err)
	}
	runRepositoryRoot := run.RepositoryRoot
	if runRepositoryRoot == "" {
		runRepositoryRoot, err = roundconfig.RepositoryRoot(run.GitRoot)
		if err != nil {
			return "", fmt.Errorf("resolve Run %q repository: %w", runID, err)
		}
	}
	if runRepositoryRoot != repositoryRoot {
		return "", fmt.Errorf("Run %q has no Implement start head for repository %q", runID, gitRoot)
	}
	return strings.TrimSpace(run.HeadSHA), nil
}

func deliveryRunResult(
	command roundfixCommandResult,
	candidateHead string,
	beforeRun *store.Run,
	afterRun *store.Run,
) (delivery.RunResult, error) {
	if command.exitCode == exitOK {
		runID := ""
		if afterRun != nil {
			runID = afterRun.ID
		}
		return delivery.RunResult{
			RunID:            runID,
			Outcome:          delivery.RunOutcomeClean,
			CandidateCommits: []string{strings.TrimSpace(candidateHead)},
		}, nil
	}

	createdRun := afterRun != nil && (beforeRun == nil || afterRun.ID != beforeRun.ID)
	if command.exitCode == exitRunFailed && createdRun {
		outcome := delivery.RunOutcome("")
		switch afterRun.State {
		case store.StateUnresolved:
			outcome = delivery.RunOutcomeUnresolved
		case store.StateBudgetExceeded:
			outcome = delivery.RunOutcomeBudgetExceeded
		}
		if outcome != "" {
			return delivery.RunResult{
				RunID:   afterRun.ID,
				Outcome: outcome,
				Reason:  strings.TrimSpace(command.stderr),
			}, nil
		}
	}
	return delivery.RunResult{}, command.failure("roundfix implement")
}

func (workflow *commandDeliveryWorkflow) ReviewPolicy(context.Context, string, string) (delivery.ReviewPolicy, error) {
	if workflow.loaded.Config.PrePRReview.Provider == "none" {
		return delivery.ReviewPolicyNone, nil
	}
	return delivery.ReviewPolicyEnabled, nil
}

func (workflow *commandDeliveryWorkflow) Review(ctx context.Context, gitRoot, specSlug string, head string) (delivery.ReviewResult, error) {
	record, result, err := workflow.runReview(ctx, gitRoot, specSlug)
	if err != nil {
		return delivery.ReviewResult{}, err
	}
	reviewResult, err := deliveryReviewResult(record, head)
	if err == nil {
		return reviewResult, nil
	}
	if result.exitCode != exitOK {
		return delivery.ReviewResult{}, result.failure("roundfix review")
	}
	return delivery.ReviewResult{}, err
}

func deliveryReviewResult(record reviewRecord, head string) (delivery.ReviewResult, error) {
	result := delivery.ReviewResult{
		Head:          record.HeadCommit,
		ArchivedSpecs: slices.Clone(record.ArchivedSpecs),
	}
	if strings.TrimSpace(record.HeadCommit) != strings.TrimSpace(head) {
		result.Outcome = delivery.ReviewOutcomeBlocked
		result.Reason = "review record names a different head"
		return result, nil
	}
	switch record.Outcome {
	case reviewOutcomeReviewed, reviewOutcomeFindingsDismissed, reviewOutcomeCeilingClosed:
		result.Outcome = delivery.ReviewOutcomeReviewed
		return result, nil
	case reviewOutcomeFindings:
		result.Outcome = delivery.ReviewOutcomeFindings
		result.Reason = record.Findings
		return result, nil
	case reviewOutcomeBlocked:
		result.Outcome = delivery.ReviewOutcomeBlocked
		result.Reason = record.Reason
		return result, nil
	default:
		return delivery.ReviewResult{}, fmt.Errorf("review record outcome %q cannot drive delivery", record.Outcome)
	}
}

func (workflow *commandDeliveryWorkflow) RecordReviewOmission(ctx context.Context, gitRoot, specSlug string, head string) error {
	record, result, err := workflow.runReview(ctx, gitRoot, specSlug)
	if err != nil {
		return err
	}
	if result.exitCode != exitOK {
		return result.failure("roundfix review")
	}
	if record.Outcome != reviewOutcomeOmitted || strings.TrimSpace(record.HeadCommit) != strings.TrimSpace(head) {
		return fmt.Errorf("roundfix review recorded outcome %q at head %q, want omitted at %q", record.Outcome, record.HeadCommit, head)
	}
	return nil
}

func (workflow *commandDeliveryWorkflow) Archive(ctx context.Context, gitRoot, specSlug, reviewedHead string) (delivery.ArchiveResult, error) {
	before, err := preflight.InspectGit(ctx, gitRoot, workflow.git)
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("inspect candidate before archive: %w", err)
	}
	if len(before.Dirty) != 0 {
		return delivery.ArchiveResult{Parent: before.HEAD}, nil
	}
	if before.HEAD != strings.TrimSpace(reviewedHead) {
		return workflow.reconcileArchiveCommit(ctx, gitRoot, specSlug, strings.TrimSpace(reviewedHead), before.HEAD)
	}
	sourcePath, destinationPath, err := workflow.archivePaths(gitRoot, specSlug)
	if err != nil {
		return delivery.ArchiveResult{}, err
	}
	active, err := workflow.git.RunGit(ctx, gitRoot, "ls-tree", "--name-only", reviewedHead, "--", sourcePath+"/_prd.md")
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("inspect active Spec at reviewed head: %w", err)
	}
	archived, err := workflow.git.RunGit(ctx, gitRoot, "ls-tree", "--name-only", reviewedHead, "--", destinationPath+"/_prd.md", destinationPath+".md")
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("inspect archived Spec at reviewed head: %w", err)
	}
	if strings.TrimSpace(active) == "" && strings.TrimSpace(archived) != "" {
		return delivery.ArchiveResult{Head: before.HEAD, AlreadyArchived: true}, nil
	}
	executable, err := workflow.stepExecutable(ctx, gitRoot, specSlug, "archive")
	if err != nil {
		return delivery.ArchiveResult{}, err
	}
	result, err := workflow.runRoundfix(ctx, executable, gitRoot, "archive", specSlug)
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("start Archive Command: %w", err)
	}
	if result.exitCode != exitOK {
		return delivery.ArchiveResult{}, result.failure("roundfix archive")
	}

	source, destination, err := workflow.archivePaths(gitRoot, specSlug)
	if err != nil {
		return delivery.ArchiveResult{}, err
	}
	exactMove, err := workflow.archiveDiffIsExact(ctx, gitRoot, source, destination)
	if err != nil || !exactMove {
		return delivery.ArchiveResult{Parent: before.HEAD, ExactSpecMove: false}, err
	}
	stagePaths := []string{source}
	if _, err := os.Stat(filepath.Join(gitRoot, destination+".md")); err == nil {
		stagePaths = append(stagePaths, destination+".md")
		archived, err := spec.ReadArchivedSpec(filepath.Join(gitRoot, filepath.Dir(destination)), specSlug)
		if err != nil {
			return delivery.ArchiveResult{}, err
		}
		stagePaths = append(stagePaths, archived.Record.Promoted...)
	} else {
		stagePaths = append(stagePaths, destination)
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, append([]string{"add", "-A", "--"}, stagePaths...)...); err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("stage archived Spec: %w", err)
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, "commit", "-m", "docs: archive "+specSlug); err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("commit archived Spec: %w", err)
	}
	head, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", "HEAD")
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("read archive commit: %w", err)
	}
	return delivery.ArchiveResult{Parent: before.HEAD, Head: strings.TrimSpace(head), ExactSpecMove: true}, nil
}

func (workflow *commandDeliveryWorkflow) reconcileArchiveCommit(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	reviewedHead string,
	archiveHead string,
) (delivery.ArchiveResult, error) {
	parent, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", archiveHead+"^")
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("read archive commit parent: %w", err)
	}
	parent = strings.TrimSpace(parent)
	result := delivery.ArchiveResult{Parent: parent, Head: archiveHead}
	if parent != reviewedHead {
		return result, nil
	}
	source, destination, err := workflow.archivePaths(gitRoot, specSlug)
	if err != nil {
		return delivery.ArchiveResult{}, err
	}
	result.ExactSpecMove, err = workflow.archiveCommitIsExact(ctx, gitRoot, parent, archiveHead, source, destination)
	return result, err
}

func (workflow *commandDeliveryWorkflow) Gate(ctx context.Context, gitRoot, specSlug, _ string) (delivery.GateResult, error) {
	artifactDir, err := roundconfig.ValidateArtifactDirectory(
		workflow.loaded.Config.Defaults.ArtifactDir,
		workflow.loaded.GitRoot,
		workflow.loaded.HomeDir,
	)
	if err != nil {
		return delivery.GateResult{}, err
	}
	outputPath := filepath.Join(artifactDir, "delivery", specSlug, "repository-gate.log")
	_, err = (daemon.ExecVerifier{}).Verify(ctx, daemon.VerifyRequest{
		WorkDir:    gitRoot,
		Command:    workflow.loaded.Config.Defaults.Verification,
		OutputPath: outputPath,
	})
	if err != nil {
		var commandErr *daemon.VerificationCommandError
		if errors.As(err, &commandErr) {
			return delivery.GateResult{Passed: false, Reason: err.Error()}, nil
		}
		return delivery.GateResult{}, err
	}
	return delivery.GateResult{Passed: true}, nil
}

func (workflow *commandDeliveryWorkflow) Authorization(ctx context.Context, gitRoot, specSlug string) (delivery.Authorization, error) {
	specsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, gitRoot)
	if err != nil {
		return delivery.Authorization{}, err
	}
	source, _, err := workflow.archivePaths(gitRoot, specSlug)
	if err != nil {
		return delivery.Authorization{}, err
	}
	// Locate the actual archive transition even when later commits follow it.
	archiveHead, err := workflow.git.RunGit(ctx, gitRoot, "log", "--first-parent", "--diff-filter=D", "--format=%H", "-1", "HEAD", "--", source+"/_prd.md")
	if err != nil {
		return delivery.Authorization{}, fmt.Errorf("find archive authorization transition: %w", err)
	}
	archiveHead = strings.TrimSpace(archiveHead)
	if archiveHead == "" {
		return delivery.Authorization{}, errors.New("read pre-archive authorization head: active PRD deletion not found")
	}
	authorizationHead, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", archiveHead+"^")
	if err != nil {
		return delivery.Authorization{}, fmt.Errorf("read pre-archive authorization head: %w", err)
	}

	resolution := spec.ReadSpecAuthorization(ctx, gitRoot, specsRoot.Path, specSlug, strings.TrimSpace(authorizationHead))
	operations := make([]string, 0, 3)
	for _, operation := range []spec.AuthorizationOperation{
		spec.AuthorizationOperationPush,
		spec.AuthorizationOperationPullRequest,
		spec.AuthorizationOperationMerge,
	} {
		if resolution.Permits(operation) {
			operations = append(operations, string(operation))
		}
	}
	return delivery.Authorization{Operations: operations}, nil
}

func (workflow *commandDeliveryWorkflow) Publication(ctx context.Context, gitRoot, specSlug, branch string) (delivery.Publication, error) {
	state, err := preflight.InspectGit(ctx, gitRoot, workflow.git)
	if err != nil {
		return delivery.Publication{}, fmt.Errorf("inspect publication candidate: %w", err)
	}
	defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, state.Branch, workflow.git)
	if defaultBranch.Source == preflight.DefaultBranchUndetermined {
		return delivery.Publication{}, errors.New("plan publication: repository default branch is unknown")
	}
	branch = strings.TrimSpace(branch)
	if state.Branch != branch {
		return delivery.Publication{}, fmt.Errorf("plan publication: checkout branch is %q, recorded item branch is %q", state.Branch, branch)
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	return delivery.Publication{
		Remote:     remote,
		HeadBranch: branch,
		BaseBranch: defaultBranch.Name,
		Title:      "feat: deliver " + specSlug,
		Body:       "Delivers Spec " + specSlug + ".",
	}, nil
}

type roundfixCommandResult struct {
	stdout   string
	stderr   string
	exitCode int
}

func (result roundfixCommandResult) failure(operation string) error {
	detail := strings.TrimSpace(result.stderr)
	if detail == "" {
		detail = strings.TrimSpace(result.stdout)
	}
	if detail == "" {
		detail = fmt.Sprintf("exit code %d", result.exitCode)
	}
	return fmt.Errorf("%s failed with exit code %d: %s", operation, result.exitCode, detail)
}

// stepExecutable selects a binary using only the owner's loaded declaration.
func (workflow *commandDeliveryWorkflow) stepExecutable(ctx context.Context, workDir, specSlug, step string) (string, error) {
	declaration := workflow.loaded.Config.Delivery.ItemBinary
	if !declaration.Declared() {
		return os.Executable()
	}
	if _, err := workflow.git.RunGit(ctx, workDir, "check-ignore", "-q", "--", declaration.Path); err != nil {
		return "", fmt.Errorf("item binary path %q is not ignored by Git", declaration.Path)
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(workflow.loaded.Config.Defaults.ArtifactDir, workflow.loaded.GitRoot, workflow.loaded.HomeDir)
	if err != nil {
		return "", err
	}
	_, err = (daemon.ExecVerifier{}).Verify(ctx, daemon.VerifyRequest{
		WorkDir:    workDir,
		Command:    declaration.Build,
		OutputPath: filepath.Join(artifactDir, "delivery", specSlug, "item-binary-build.log"),
	})
	if err != nil {
		return "", fmt.Errorf("build item binary: %w", err)
	}
	executable, err := filepath.Abs(filepath.Join(workDir, filepath.FromSlash(declaration.Path)))
	if err != nil {
		return "", fmt.Errorf("resolve item binary: %w", err)
	}
	probe, err := workflow.runRoundfix(ctx, executable, workDir, "migrate", "--check")
	if err != nil {
		return "", fmt.Errorf("run item binary %q: %w", executable, err)
	}
	log := workflow.log
	if log == nil {
		log = os.Stderr
	}
	if probe.exitCode == exitOK {
		if _, err := fmt.Fprintf(log, "roundfix: Delivery Queue item %s: %s runs the item binary %s\n", specSlug, step, executable); err != nil {
			return "", fmt.Errorf("write item binary selection: %w", err)
		}
		return executable, nil
	}
	detail := ""
	for _, output := range []string{probe.stderr, probe.stdout} {
		for _, line := range strings.Split(output, "\n") {
			if strings.TrimSpace(line) != "" {
				detail = strings.TrimSpace(line)
				break
			}
		}
		if detail != "" {
			break
		}
	}
	if _, err := fmt.Fprintf(log, "roundfix: notice: Delivery Queue item %s: %s runs the owner's binary; the item binary's migrate --check exited %d: %s\n", specSlug, step, probe.exitCode, detail); err != nil {
		return "", fmt.Errorf("write item binary fallback: %w", err)
	}
	return os.Executable()
}

func (workflow *commandDeliveryWorkflow) runRoundfix(ctx context.Context, executable, workDir string, args ...string) (roundfixCommandResult, error) {
	command := exec.CommandContext(ctx, executable, args...)
	command.Dir = workDir
	command.Env = deliveryCommandEnvironment(os.Environ(), workflow.loaded.HomeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	runErr := command.Run()
	result := roundfixCommandResult{stdout: stdout.String(), stderr: stderr.String(), exitCode: exitOK}
	if runErr == nil {
		return result, nil
	}
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result, nil
	}
	return roundfixCommandResult{}, runErr
}

func (workflow *commandDeliveryWorkflow) latestImplementRun(ctx context.Context, gitRoot, specSlug string) (store.Run, bool, error) {
	runs, err := workflow.store.ListRuns(ctx, store.ListRunsQuery{GitRoot: gitRoot, States: store.StatesAll})
	if err != nil {
		return store.Run{}, false, fmt.Errorf("list Implement Runs: %w", err)
	}
	for _, run := range runs {
		if run.Kind == store.KindImplement && run.SpecSlug == specSlug {
			return run, true, nil
		}
	}
	return store.Run{}, false, nil
}

func (workflow *commandDeliveryWorkflow) runReview(ctx context.Context, gitRoot, specSlug string) (reviewRecord, roundfixCommandResult, error) {
	executable, err := workflow.stepExecutable(ctx, gitRoot, specSlug, "review")
	if err != nil {
		return reviewRecord{}, roundfixCommandResult{}, err
	}
	result, err := workflow.runRoundfix(ctx, executable, gitRoot, "review")
	if err != nil {
		return reviewRecord{}, roundfixCommandResult{}, fmt.Errorf("start Pre-PR Review Command: %w", err)
	}
	var record reviewRecord
	if err := json.Unmarshal([]byte(result.stdout), &record); err != nil {
		if result.exitCode != exitOK {
			return reviewRecord{}, result, result.failure("roundfix review")
		}
		return reviewRecord{}, result, fmt.Errorf("decode review record: %w", err)
	}
	return record, result, nil
}

func (workflow *commandDeliveryWorkflow) archivePaths(gitRoot, specSlug string) (string, string, error) {
	specsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, gitRoot)
	if err != nil {
		return "", "", err
	}
	source, err := filepath.Rel(gitRoot, filepath.Join(specsRoot.Path, specSlug))
	if err != nil {
		return "", "", fmt.Errorf("resolve active Spec path: %w", err)
	}
	destinationRoot := filepath.Join(specsRoot.Path, "_archived")
	if specsRoot.BuiltInRoot {
		destinationRoot = filepath.Join(gitRoot, filepath.FromSlash(spec.ArchiveDir(spec.ArchiveKindSpec)))
	}
	destination, err := filepath.Rel(gitRoot, filepath.Join(destinationRoot, specSlug))
	if err != nil {
		return "", "", fmt.Errorf("resolve archived Spec path: %w", err)
	}
	return filepath.ToSlash(source), filepath.ToSlash(destination), nil
}

func (workflow *commandDeliveryWorkflow) archiveDiffIsExact(ctx context.Context, gitRoot, source, destination string) (bool, error) {
	if _, err := os.Lstat(filepath.Join(gitRoot, destination+".md")); err == nil {
		parent, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", "HEAD")
		if err != nil {
			return false, err
		}
		return workflow.archiveRetirementIsExact(ctx, gitRoot, strings.TrimSpace(parent), "", source, destination)
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	state, err := preflight.InspectGit(ctx, gitRoot, workflow.git)
	if err != nil {
		return false, fmt.Errorf("inspect archive diff: %w", err)
	}
	sawSource := false
	sawDestination := false
	for _, change := range state.Dirty {
		changed := filepath.ToSlash(strings.TrimSuffix(strings.TrimSpace(change.Path), "/"))
		switch {
		case changed == source || strings.HasPrefix(changed, source+"/"):
			sawSource = true
		case changed == destination || strings.HasPrefix(changed, destination+"/"):
			sawDestination = true
		default:
			return false, nil
		}
	}
	return sawSource && sawDestination, nil
}

func (workflow *commandDeliveryWorkflow) archiveCommitIsExact(
	ctx context.Context,
	gitRoot string,
	parent string,
	head string,
	source string,
	destination string,
) (bool, error) {
	exists, err := workflow.gitObjectExists(ctx, gitRoot, head+":"+destination+".md")
	if err != nil {
		return false, err
	}
	if exists {
		return workflow.archiveRetirementIsExact(ctx, gitRoot, parent, head, source, destination)
	}
	changed, err := workflow.git.RunGit(ctx, gitRoot, "diff", "--name-only", "--no-renames", "-z", parent, head, "--")
	if err != nil {
		return false, fmt.Errorf("inspect archive commit diff: %w", err)
	}
	sawSource := false
	sawDestination := false
	for _, rawPath := range strings.Split(strings.TrimSuffix(changed, "\x00"), "\x00") {
		changedPath := filepath.ToSlash(rawPath)
		if changedPath == "" {
			continue
		}
		switch {
		case changedPath == source || strings.HasPrefix(changedPath, source+"/"):
			sawSource = true
		case changedPath == destination || strings.HasPrefix(changedPath, destination+"/"):
			sawDestination = true
		default:
			return false, nil
		}
	}
	if !sawSource || !sawDestination {
		return false, nil
	}
	sourceTree, err := workflow.gitTreeEntries(ctx, gitRoot, parent+":"+source)
	if err != nil {
		return false, fmt.Errorf("read active Spec tree: %w", err)
	}
	destinationTree, err := workflow.gitTreeEntries(ctx, gitRoot, head+":"+destination)
	if err != nil {
		return false, fmt.Errorf("read archived Spec tree: %w", err)
	}
	sourcePRDEntry, ok := sourceTree["_prd.md"]
	if !ok {
		return false, nil
	}
	destinationPRDEntry, ok := destinationTree["_prd.md"]
	sourcePRDKind := gitTreeEntryKind(sourcePRDEntry)
	if !ok || sourcePRDKind == "" || sourcePRDKind != gitTreeEntryKind(destinationPRDEntry) {
		return false, nil
	}
	delete(sourceTree, "_prd.md")
	delete(destinationTree, "_prd.md")
	if len(sourceTree) != len(destinationTree) {
		return false, nil
	}
	for name, sourceEntry := range sourceTree {
		destinationEntry, exists := destinationTree[name]
		if !exists {
			return false, nil
		}
		if sourceEntry == destinationEntry {
			continue
		}
		kind := gitTreeEntryKind(sourceEntry)
		if filepath.Ext(name) != ".md" || (kind != "100644 blob" && kind != "100755 blob") || kind != gitTreeEntryKind(destinationEntry) {
			return false, nil
		}
		sourceContent, err := workflow.git.RunGit(ctx, gitRoot, "show", parent+":"+source+"/"+name)
		if err != nil {
			return false, fmt.Errorf("read active Spec file %q: %w", name, err)
		}
		destinationContent, err := workflow.git.RunGit(ctx, gitRoot, "show", head+":"+destination+"/"+name)
		if err != nil {
			return false, fmt.Errorf("read archived Spec file %q: %w", name, err)
		}
		if !spec.ArchiveLinksMatch([]byte(sourceContent), []byte(destinationContent),
			filepath.Dir(filepath.Join(source, name)), filepath.Dir(filepath.Join(destination, name)), source) {
			return false, nil
		}
	}
	sourcePRD, err := workflow.git.RunGit(ctx, gitRoot, "show", parent+":"+source+"/_prd.md")
	if err != nil {
		return false, fmt.Errorf("read active Spec PRD: %w", err)
	}
	destinationPRD, err := workflow.git.RunGit(ctx, gitRoot, "show", head+":"+destination+"/_prd.md")
	if err != nil {
		return false, fmt.Errorf("read archived Spec PRD: %w", err)
	}
	if !archivePRDChangeIsExact([]byte(sourcePRD), []byte(destinationPRD), source, destination) {
		return false, nil
	}
	sourceAtHead, err := workflow.gitObjectExists(ctx, gitRoot, head+":"+source)
	if err != nil {
		return false, fmt.Errorf("inspect active Spec after archive: %w", err)
	}
	destinationAtParent, err := workflow.gitObjectExists(ctx, gitRoot, parent+":"+destination)
	if err != nil {
		return false, fmt.Errorf("inspect archive destination before archive: %w", err)
	}
	return !sourceAtHead && !destinationAtParent, nil
}

// archiveRetirementIsExact proves both committed and unstaged record retirements.
// An empty head reads the working tree while keeping provenance anchored at parent.
func (workflow *commandDeliveryWorkflow) archiveRetirementIsExact(ctx context.Context, root, parent, head, source, destination string) (bool, error) {
	read := func(path string) ([]byte, error) {
		if head == "" {
			return os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		}
		text, err := workflow.git.RunGit(ctx, root, "show", head+":"+path)
		return []byte(text), err
	}
	content, err := read(destination + ".md")
	if err != nil {
		return false, err
	}
	record, err := spec.ParseArchiveRecord(content)
	if err != nil {
		return false, nil
	}
	if record.Spec != filepath.Base(source) || record.Source != source || record.SourceRevision != parent {
		return false, nil
	}
	sourceTree, err := workflow.gitTreeEntries(ctx, root, parent+":"+source)
	if err != nil {
		return false, err
	}
	if _, ok := sourceTree["_prd.md"]; !ok {
		return false, nil
	}
	prd, err := workflow.git.RunGit(ctx, root, "show", parent+":"+source+"/_prd.md")
	if err != nil {
		return false, err
	}
	fm, body, ok := splitArchivePRD([]byte(prd))
	if !ok {
		return false, nil
	}
	var meta struct{ Created string }
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return false, nil
	}
	title := ""
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}
	if record.Title != title || record.Created != meta.Created {
		return false, nil
	}
	allowed := map[string]bool{destination + ".md": true}
	for name := range sourceTree {
		allowed[source+"/"+name] = true
	}
	for _, path := range record.Promoted {
		if !strings.HasPrefix(path, "docs/references/") || filepath.ToSlash(filepath.Clean(path)) != path || allowed[path] {
			return false, nil
		}
		allowed[path] = true
		var blob string
		if head == "" {
			info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
			if err != nil {
				return false, err
			}
			if !info.Mode().IsRegular() {
				return false, nil
			}
			blob, err = workflow.git.RunGit(ctx, root, "hash-object", "--", path)
		} else {
			entries, treeErr := workflow.gitTreeEntries(ctx, root, head)
			if treeErr != nil {
				return false, treeErr
			}
			entry, ok := entries[path]
			if !ok || (gitTreeEntryKind(entry) != "100644 blob" && gitTreeEntryKind(entry) != "100755 blob") {
				return false, nil
			}
			blob = strings.Fields(entry)[2]
		}
		if err != nil {
			return false, err
		}
		matched := false
		for _, entry := range sourceTree {
			fields := strings.Fields(entry)
			if len(fields) == 3 && (gitTreeEntryKind(entry) == "100644 blob" || gitTreeEntryKind(entry) == "100755 blob") && fields[2] == strings.TrimSpace(blob) {
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}
	for path := range allowed {
		if !strings.HasPrefix(path, source+"/") {
			exists, err := workflow.gitObjectExists(ctx, root, parent+":"+path)
			if err != nil {
				return false, err
			}
			if exists {
				return false, nil
			}
		}
	}
	var changes []string
	if head == "" {
		state, err := preflight.InspectGit(ctx, root, workflow.git)
		if err != nil {
			return false, err
		}
		for _, change := range state.Dirty {
			changes = append(changes, change.Path)
		}
		if _, err := os.Lstat(filepath.Join(root, source)); !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
	} else {
		changed, err := workflow.git.RunGit(ctx, root, "diff", "--name-only", "--no-renames", "-z", parent, head, "--")
		if err != nil {
			return false, err
		}
		changes = strings.Split(strings.TrimSuffix(changed, "\x00"), "\x00")
		exists, err := workflow.gitObjectExists(ctx, root, head+":"+source)
		if err != nil || exists {
			return false, err
		}
	}
	if len(changes) != len(allowed) {
		return false, nil
	}
	seen := map[string]bool{}
	for _, path := range changes {
		if !allowed[path] || seen[path] {
			return false, nil
		}
		seen[path] = true
	}
	return true, nil
}

func (workflow *commandDeliveryWorkflow) gitTreeEntries(ctx context.Context, gitRoot, object string) (map[string]string, error) {
	listing, err := workflow.git.RunGit(ctx, gitRoot, "ls-tree", "-r", "-z", object)
	if err != nil {
		return nil, err
	}
	entries := make(map[string]string)
	for _, record := range strings.Split(strings.TrimSuffix(listing, "\x00"), "\x00") {
		if record == "" {
			continue
		}
		identity, name, ok := strings.Cut(record, "\t")
		if !ok || identity == "" || name == "" {
			return nil, fmt.Errorf("malformed ls-tree entry %q", record)
		}
		entries[filepath.ToSlash(name)] = identity
	}
	return entries, nil
}

func gitTreeEntryKind(entry string) string {
	fields := strings.Fields(entry)
	if len(fields) != 3 {
		return ""
	}
	return fields[0] + " " + fields[1]
}

func archivePRDChangeIsExact(source, destination []byte, sourceDir, destinationDir string) bool {
	sourceFrontmatter, sourceBody, ok := splitArchivePRD(source)
	if !ok {
		return false
	}
	destinationFrontmatter, destinationBody, ok := splitArchivePRD(destination)
	if !ok || (!bytes.Equal(sourceBody, destinationBody) && !spec.ArchiveLinksMatch(sourceBody, destinationBody, sourceDir, destinationDir, sourceDir)) {
		return false
	}
	var sourceValues map[string]any
	if err := yaml.Unmarshal(sourceFrontmatter, &sourceValues); err != nil {
		return false
	}
	var destinationValues map[string]any
	if err := yaml.Unmarshal(destinationFrontmatter, &destinationValues); err != nil {
		return false
	}
	status, statusOK := destinationValues["status"].(string)
	archived, archivedOK := destinationValues["archived"].(string)
	sourceSlug, sourceSlugOK := destinationValues["source_slug"].(string)
	if !statusOK || status != "archived" || !archivedOK || !validArchiveDate(archived) || !sourceSlugOK || sourceSlug != filepath.Base(sourceDir) {
		return false
	}
	for _, key := range []string{"status", "archived", "source_slug", "unproven"} {
		delete(sourceValues, key)
		delete(destinationValues, key)
	}
	return reflect.DeepEqual(sourceValues, destinationValues)
}

func splitArchivePRD(content []byte) ([]byte, []byte, bool) {
	const opening = "---\n"
	if !bytes.HasPrefix(content, []byte(opening)) {
		return nil, nil, false
	}
	rest := content[len(opening):]
	end := bytes.Index(rest, []byte("\n---"))
	if end < 0 {
		return nil, nil, false
	}
	bodyStart := end + len("\n---")
	for bodyStart < len(rest) && (rest[bodyStart] == '\n' || rest[bodyStart] == '\r') {
		bodyStart++
	}
	return rest[:end], rest[bodyStart:], true
}

func validArchiveDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

func (workflow *commandDeliveryWorkflow) gitObjectExists(ctx context.Context, gitRoot, object string) (bool, error) {
	_, err := workflow.git.RunGit(ctx, gitRoot, "cat-file", "-e", object)
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}
	return false, err
}

func deliveryCommandEnvironment(environment []string, homeDir string) []string {
	result := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if key == "HOME" {
			continue
		}
		result = append(result, entry)
	}
	return append(result, "HOME="+homeDir)
}

// UnmetPrerequisites reads merge evidence from the refreshed delivery default,
// never from the owner's checkout or an unmerged item branch.
func (workflow *commandDeliveryWorkflow) UnmetPrerequisites(ctx context.Context, gitRoot, specSlug string) ([]string, error) {
	root, err := roundconfig.ResolveSpecsRoot(workflow.loaded, gitRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve prerequisite Specs Root: %w", err)
	}
	if root.External {
		return nil, errors.New("read prerequisites: Specs Root is outside the delivery repository")
	}
	defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, "", workflow.git)
	if defaultBranch.Source == preflight.DefaultBranchUndetermined {
		return nil, errors.New("read prerequisites: repository default branch is unknown")
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, "fetch", remote, defaultBranch.Name); err != nil {
		return nil, fmt.Errorf("refresh prerequisite default branch %q: %w", defaultBranch.Name, err)
	}
	ref := "refs/remotes/" + remote + "/" + defaultBranch.Name
	manifestPath, err := filepath.Rel(gitRoot, filepath.Join(root.Path, specSlug, "_tasks.md"))
	if err != nil {
		return nil, fmt.Errorf("resolve prerequisite manifest path: %w", err)
	}
	content, err := workflow.git.RunGit(ctx, gitRoot, "show", ref+":"+filepath.ToSlash(manifestPath))
	if err != nil {
		return nil, fmt.Errorf("read prerequisite _tasks.md at %s: %w", ref, err)
	}
	requires, err := spec.ParseRequiredSpecs([]byte(content), specSlug)
	if err != nil {
		return nil, err
	}
	if len(requires) == 0 {
		return nil, nil
	}
	archiveRoot, err := filepath.Rel(gitRoot, spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot))
	if err != nil {
		return nil, fmt.Errorf("resolve prerequisite archive root: %w", err)
	}
	// One immutable-tree read checks every prerequisite, without a subprocess
	// per slug or treating a failed Git read as proof of absence.
	paths, err := workflow.git.RunGit(ctx, gitRoot, "ls-tree", "-r", "--name-only", ref, "--", filepath.ToSlash(archiveRoot))
	if err != nil {
		return nil, fmt.Errorf("read prerequisite archives at %s: %w", ref, err)
	}
	present := make(map[string]bool)
	for _, path := range strings.Split(paths, "\n") {
		present[path] = true
	}
	var unmet []string
	for _, slug := range requires {
		if !present[filepath.ToSlash(filepath.Join(archiveRoot, slug, "_prd.md"))] {
			unmet = append(unmet, slug)
		}
	}
	return unmet, nil
}

var _ delivery.ConflictResolver = (*commandDeliveryWorkflow)(nil)

// ResolveConflict trusts only declarations committed on the refreshed default.
func (workflow *commandDeliveryWorkflow) ResolveConflict(ctx context.Context, workDir, specSlug, head string) (resolution delivery.ConflictResolution, resultErr error) {
	state, err := preflight.InspectGit(ctx, workDir, workflow.git)
	if err != nil {
		return resolution, fmt.Errorf("inspect conflict worktree: %w", err)
	}
	if len(state.Dirty) != 0 || state.HEAD != head {
		return resolution, errors.New("resolve conflict requires a clean worktree at the candidate head")
	}
	branch := preflight.DetectDefaultBranch(ctx, workDir, state.Branch, workflow.git)
	if branch.Source == preflight.DefaultBranchUndetermined {
		return resolution, errors.New("resolve conflict: default branch is unknown")
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	ref := "refs/remotes/" + remote + "/" + branch.Name
	if _, err := workflow.git.RunGit(ctx, workDir, "fetch", remote, "+refs/heads/"+branch.Name+":"+ref); err != nil {
		return resolution, fmt.Errorf("fetch conflict default: %w", err)
	}
	defaultHead, err := workflow.git.RunGit(ctx, workDir, "rev-parse", ref+"^{commit}")
	if err != nil {
		return resolution, fmt.Errorf("resolve conflict default commit: %w", err)
	}
	defaultHead = strings.TrimSpace(defaultHead)
	config, err := roundconfig.DeliveryConfigAtCommit(ctx, workflow.git, workDir, workflow.loaded.UserConfigPath, defaultHead)
	if err != nil {
		return resolution, err
	}
	declarations := config.Delivery.DerivedPaths
	matches := func(name string) bool {
		for _, declaration := range declarations {
			if declaration.Matches(name) {
				return true
			}
		}
		return false
	}
	initialUntracked, err := workflow.git.RunGit(ctx, workDir, "ls-files", "--others", "-z")
	if err != nil {
		return resolution, fmt.Errorf("inspect initial untracked paths: %w", err)
	}
	existingUntracked := make(map[string]bool)
	for _, name := range nulPaths(initialUntracked) {
		existingUntracked[name] = true
	}
	_, mergeErr := workflow.git.RunGit(ctx, workDir, "-c", "merge.conflictStyle=merge", "merge", "--no-ff", "--no-commit", defaultHead)
	merging, err := workflow.gitObjectExists(ctx, workDir, "MERGE_HEAD")
	if err != nil {
		return resolution, fmt.Errorf("inspect conflict merge: %w", err)
	}
	if !merging {
		if mergeErr != nil {
			return resolution, fmt.Errorf("merge conflict default: %w", mergeErr)
		}
		return resolution, nil
	}
	committed := false
	// The worktree was clean. Restore tracked files and remove only newly
	// created untracked files before aborting an unsuccessful regeneration.
	defer func() {
		if committed {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()
		untracked, readErr := workflow.git.RunGit(cleanupCtx, workDir, "ls-files", "--others", "-z")
		var cleanupErr error
		if readErr != nil {
			cleanupErr = readErr
		} else {
			for _, name := range nulPaths(untracked) {
				if existingUntracked[name] {
					continue
				}
				cleanupErr = errors.Join(cleanupErr, os.Remove(filepath.Join(workDir, filepath.FromSlash(name))))
			}
		}
		_, restoreErr := workflow.git.RunGit(cleanupCtx, workDir, "restore", "--source=HEAD", "--staged", "--worktree", "--", ".")
		_, abortErr := workflow.git.RunGit(cleanupCtx, workDir, "merge", "--abort")
		resultErr = errors.Join(resultErr, cleanupErr, restoreErr, abortErr)
	}()
	conflicts, err := workflow.git.RunGit(ctx, workDir, "diff", "--name-only", "--diff-filter=U", "-z")
	if err != nil {
		return resolution, fmt.Errorf("list conflict paths: %w", err)
	}
	paths := nulPaths(conflicts)
	matched := make([]bool, len(declarations))
	lineResolutions := make(map[string][]byte)
	for _, name := range paths {
		if matches(name) {
			for index, declaration := range declarations {
				if declaration.Matches(name) {
					matched[index] = true
				}
			}
			continue
		}
		lineScoped := false
		for _, declaration := range declarations {
			lineScoped = lineScoped || declaration.MatchesLines(name)
		}
		if !lineScoped {
			resolution.SourcePaths = append(resolution.SourcePaths, name)
			continue
		}
		content, mode, err := readDerivedLineFile(workDir, name)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return resolution, fmt.Errorf("read line-scoped conflict %s: %w", name, err)
		}
		for index, declaration := range declarations {
			if !declaration.MatchesLines(name) || !mode.IsRegular() {
				continue
			}
			resolved, ok := resolveDerivedLineHunks(content, declaration.Lines.Match)
			if ok {
				lineResolutions[name] = resolved
				matched[index] = true
			}
		}
		if _, ok := lineResolutions[name]; !ok {
			resolution.SourcePaths = append(resolution.SourcePaths, name)
		}
	}
	if len(resolution.SourcePaths) != 0 {
		return resolution, nil
	}
	if mergeErr != nil && len(paths) == 0 {
		return resolution, fmt.Errorf("merge conflict default: %w", mergeErr)
	}
	for _, name := range paths {
		if content, ok := lineResolutions[name]; ok {
			if err := os.WriteFile(filepath.Join(workDir, filepath.FromSlash(name)), content, 0600); err != nil {
				return resolution, fmt.Errorf("take default derived lines: %w", err)
			}
		} else if _, err := workflow.git.RunGit(ctx, workDir, "checkout", "--theirs", "--", name); err != nil {
			return resolution, fmt.Errorf("take default derived path: %w", err)
		}
		if _, err := workflow.git.RunGit(ctx, workDir, "add", "--", name); err != nil {
			return resolution, fmt.Errorf("stage default derived path: %w", err)
		}
	}
	baseline, err := workflow.git.RunGit(ctx, workDir, "write-tree")
	if err != nil {
		return resolution, fmt.Errorf("snapshot merge before regeneration: %w", err)
	}
	// Snapshot every tracked line-scoped file, including cleanly merged files
	// that a matched command may rewrite. Preserve bytes and file mode exactly.
	lineBaseline := make(map[string]derivedLineSnapshot)
	tracked, err := workflow.git.RunGit(ctx, workDir, "ls-files", "-z")
	if err != nil {
		return resolution, fmt.Errorf("list line-scoped merge paths: %w", err)
	}
	for _, name := range nulPaths(tracked) {
		if matches(name) {
			continue
		}
		for _, declaration := range declarations {
			if declaration.MatchesLines(name) {
				content, mode, err := readDerivedLineFile(workDir, name)
				if err != nil {
					return resolution, fmt.Errorf("snapshot derived lines %s: %w", name, err)
				}
				lineBaseline[name] = derivedLineSnapshot{content: content, mode: mode}
				break
			}
		}
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(workflow.loaded.Config.Defaults.ArtifactDir, workflow.loaded.GitRoot, workflow.loaded.HomeDir)
	if err != nil {
		return resolution, err
	}
	for index, declaration := range declarations {
		if !matched[index] {
			continue
		}
		_, err := (daemon.ExecVerifier{}).Verify(ctx, daemon.VerifyRequest{WorkDir: workDir, Command: declaration.Regenerate, OutputPath: filepath.Join(artifactDir, "delivery", specSlug, fmt.Sprintf("derived-regeneration-%d.log", index+1))})
		if err != nil {
			return resolution, fmt.Errorf("regenerate derived paths: %w", err)
		}
		resolution.Regenerated = append(resolution.Regenerated, declaration.Regenerate)
	}
	changed, err := workflow.git.RunGit(ctx, workDir, "diff", "--name-only", "-z", strings.TrimSpace(baseline), "--")
	if err != nil {
		return resolution, fmt.Errorf("inspect regeneration changes: %w", err)
	}
	untracked, err := workflow.git.RunGit(ctx, workDir, "ls-files", "--others", "-z")
	if err != nil {
		return resolution, fmt.Errorf("inspect regeneration new paths: %w", err)
	}
	changedPaths := nulPaths(changed)
	for _, name := range nulPaths(untracked) {
		if !existingUntracked[name] {
			changedPaths = append(changedPaths, name)
		}
	}
	for _, name := range changedPaths {
		allowed := matches(name)
		if !allowed {
			before, exists := lineBaseline[name]
			after, mode, err := readDerivedLineFile(workDir, name)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return resolution, fmt.Errorf("inspect regenerated lines %s: %w", name, err)
			}
			if exists && err == nil && mode.IsRegular() && mode == before.mode {
				for _, declaration := range declarations {
					if declaration.MatchesLines(name) && derivedLineChangesAllowed(before.content, after, declaration.Lines.Match) {
						allowed = true
						break
					}
				}
			}
		}
		if !allowed {
			resolution.SourcePaths = append(resolution.SourcePaths, "regenerated "+name+" outside delivery.derived_paths")
		}
	}
	if len(resolution.SourcePaths) > 0 {
		return resolution, nil
	}
	for _, name := range changedPaths {
		if _, err := workflow.git.RunGit(ctx, workDir, "add", "-A", "--", name); err != nil {
			return resolution, fmt.Errorf("stage regenerated path: %w", err)
		}
	}
	if _, err := workflow.git.RunGit(ctx, workDir, "commit", "-m", "chore: merge default branch and regenerate derived paths\n\nRoundfix-Delivery: derived-merge"); err != nil {
		return resolution, fmt.Errorf("commit derived merge: %w", err)
	}
	committed = true
	newHead, err := workflow.git.RunGit(ctx, workDir, "rev-parse", "HEAD")
	if err != nil {
		return resolution, fmt.Errorf("read derived merge head: %w", err)
	}
	resolution.Head = strings.TrimSpace(newHead)
	return resolution, nil
}

type derivedLineSnapshot struct {
	content []byte
	mode    os.FileMode
}

func readDerivedLineFile(workDir, name string) ([]byte, os.FileMode, error) {
	name = filepath.Join(workDir, filepath.FromSlash(name))
	info, err := os.Lstat(name)
	if err != nil {
		return nil, 0, err
	}
	if !info.Mode().IsRegular() {
		return nil, info.Mode(), nil
	}
	content, err := os.ReadFile(name)
	return content, info.Mode(), err
}

// Resolve only well-formed, two-sided text hunks. Keep all bytes outside
// hunks, including cleanly merged source edits, and take the incoming side.
func resolveDerivedLineHunks(content []byte, pattern string) ([]byte, bool) {
	if bytes.ContainsRune(content, '\x00') {
		return nil, false
	}
	match := regexp.MustCompile(pattern) // Project Config has validated it.
	lines := strings.SplitAfter(string(content), "\n")
	var output strings.Builder
	var ours, theirs []string
	state, width, hunks := 0, 0, 0
	for _, line := range lines {
		text := strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(text, "<<<<<<<") {
			if state != 0 {
				return nil, false
			}
			width = len(text) - len(strings.TrimLeft(text, "<"))
			if !strings.HasPrefix(text[width:], " ") {
				return nil, false
			}
			state = 1
			ours, theirs = nil, nil
			hunks++
		} else if text == strings.Repeat("=", width) && state == 1 {
			state = 2
		} else if strings.HasPrefix(text, ">>>>>>>") {
			if state != 2 || !strings.HasPrefix(text, strings.Repeat(">", width)+" ") {
				return nil, false
			}
			if len(ours) != len(theirs) {
				return nil, false
			}
			for index, oldLine := range ours {
				newLine := theirs[index]
				if oldLine != newLine && (!match.MatchString(strings.TrimSuffix(oldLine, "\n")) || !match.MatchString(strings.TrimSuffix(newLine, "\n"))) {
					return nil, false
				}
				output.WriteString(newLine)
			}
			state = 0
		} else if state != 0 {
			if state == 1 {
				ours = append(ours, line)
			} else {
				theirs = append(theirs, line)
			}
		} else {
			output.WriteString(line)
		}
	}
	return []byte(output.String()), state == 0 && hunks > 0
}

func derivedLineChangesAllowed(before, after []byte, pattern string) bool {
	oldLines := strings.SplitAfter(string(before), "\n")
	newLines := strings.SplitAfter(string(after), "\n")
	if len(oldLines) != len(newLines) {
		return false
	}
	match := regexp.MustCompile(pattern) // Project Config has validated it.
	for index, oldLine := range oldLines {
		newLine := newLines[index]
		if oldLine != newLine && (!match.MatchString(strings.TrimSuffix(oldLine, "\n")) || !match.MatchString(strings.TrimSuffix(newLine, "\n"))) {
			return false
		}
	}
	return true
}

func nulPaths(output string) []string {
	if output == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(output, "\x00"), "\x00")
}
