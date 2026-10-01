package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

type commandDeliveryWorkflow struct {
	store  *store.Store
	loaded roundconfig.Loaded
	git    preflight.GitRunner
}

var _ delivery.ItemRecovery = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemWorkspace = (*commandDeliveryWorkflow)(nil)
var _ delivery.ItemRevalidator = (*commandDeliveryWorkflow)(nil)

const deliveryBranchPrefix = "roundfix/deliver-"

func newCommandDeliveryEngine(runStore *store.Store, loaded roundconfig.Loaded) deliveryEngine {
	workflow := &commandDeliveryWorkflow{
		store:  runStore,
		loaded: loaded,
		git:    preflight.ExecGitRunner{},
	}
	return delivery.NewEngine(runStore, delivery.EngineDependencies{
		Workspace:    workflow,
		Runner:       workflow,
		Reviewer:     workflow,
		Archiver:     workflow,
		Gate:         workflow,
		Authorizer:   workflow,
		Publication:  workflow,
		PullRequests: delivery.NewGitHubCLI(loaded.GitRoot),
		Checks:       delivery.NewGitHubCLI(loaded.GitRoot),
		Recovery:     workflow,
		Revalidator:  workflow,
		Log:          os.Stderr,
	})
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
		state.Archived = true
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
	branch, err := newDeliveryBranch(specSlug)
	if err != nil {
		return "", "", err
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
	if err := runworktree.CreateItem(ctx, ref, runworktree.ItemCreateOptions{
		HeadSHA:  remote + "/" + defaultBranch.Name,
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
			return "", "", fmt.Errorf(
				"candidate head %q is not represented by merge commit %q",
				candidateHead,
				mergeCommit,
			)
		}
	}

	return resolvedMergeCommit, resolvedHead, nil
}

func (workflow *commandDeliveryWorkflow) RunSpec(ctx context.Context, gitRoot, specSlug string) (delivery.RunResult, error) {
	beforeRun, beforeFound, err := workflow.latestImplementRun(ctx, workflow.loaded.GitRoot, specSlug)
	if err != nil {
		return delivery.RunResult{}, err
	}
	result, err := workflow.runRoundfix(ctx, gitRoot, "implement", "--spec", specSlug)
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
	return deliveryRunResult(result, candidateHead, before, after)
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

func (workflow *commandDeliveryWorkflow) Review(ctx context.Context, gitRoot, _ string, head string) (delivery.ReviewResult, error) {
	record, result, err := workflow.runReview(ctx, gitRoot)
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
	case reviewOutcomeReviewed, reviewOutcomeFindingsDismissed:
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

func (workflow *commandDeliveryWorkflow) RecordReviewOmission(ctx context.Context, gitRoot, _ string, head string) error {
	record, result, err := workflow.runReview(ctx, gitRoot)
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
	result, err := workflow.runRoundfix(ctx, gitRoot, "archive", specSlug)
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
	if _, err := workflow.git.RunGit(ctx, gitRoot, "add", "-A", "--", source, destination); err != nil {
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
	// The archive commit has moved the authorization record out of the active
	// Spec Root. Read the reviewed parent: the archive transition has already
	// proved that HEAD is its one exact Spec move.
	authorizationHead, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", "HEAD^")
	if err != nil {
		return delivery.Authorization{}, fmt.Errorf("read pre-archive authorization head: %w", err)
	}
	specsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, gitRoot)
	if err != nil {
		return delivery.Authorization{}, err
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

func (workflow *commandDeliveryWorkflow) runRoundfix(ctx context.Context, workDir string, args ...string) (roundfixCommandResult, error) {
	executable, err := os.Executable()
	if err != nil {
		return roundfixCommandResult{}, fmt.Errorf("resolve Roundfix executable: %w", err)
	}
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

func (workflow *commandDeliveryWorkflow) runReview(ctx context.Context, gitRoot string) (reviewRecord, roundfixCommandResult, error) {
	result, err := workflow.runRoundfix(ctx, gitRoot, "review")
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
	if !maps.Equal(sourceTree, destinationTree) {
		return false, nil
	}
	sourcePRD, err := workflow.git.RunGit(ctx, gitRoot, "show", parent+":"+source+"/_prd.md")
	if err != nil {
		return false, fmt.Errorf("read active Spec PRD: %w", err)
	}
	destinationPRD, err := workflow.git.RunGit(ctx, gitRoot, "show", head+":"+destination+"/_prd.md")
	if err != nil {
		return false, fmt.Errorf("read archived Spec PRD: %w", err)
	}
	if !archivePRDChangeIsExact([]byte(sourcePRD), []byte(destinationPRD), filepath.Base(source)) {
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

func archivePRDChangeIsExact(source, destination []byte, slug string) bool {
	sourceFrontmatter, sourceBody, ok := splitArchivePRD(source)
	if !ok {
		return false
	}
	destinationFrontmatter, destinationBody, ok := splitArchivePRD(destination)
	if !ok || !bytes.Equal(sourceBody, destinationBody) {
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
	if !statusOK || status != "archived" || !archivedOK || !validArchiveDate(archived) || !sourceSlugOK || sourceSlug != slug {
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
