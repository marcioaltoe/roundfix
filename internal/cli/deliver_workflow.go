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
	})
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
	branch = strings.TrimSpace(branch)
	itemWorktree = strings.TrimSpace(itemWorktree)
	if itemWorktree == "" {
		if branch == "" {
			return nil
		}
		if err := runworktree.CleanupItemBranch(ctx, strings.TrimSpace(gitRoot), branch); err != nil {
			return fmt.Errorf("remove item branch without recorded worktree: %w", err)
		}
		return nil
	}
	if err := runworktree.CleanupItem(ctx, runworktree.ItemRef{
		Path:     itemWorktree,
		Branch:   branch,
		UserRoot: strings.TrimSpace(gitRoot),
	}); err != nil {
		return fmt.Errorf("remove item worktree and branch: %w", err)
	}
	return nil
}

func (workflow *commandDeliveryWorkflow) RunSpec(ctx context.Context, gitRoot, specSlug string) (delivery.RunResult, error) {
	result, err := workflow.runRoundfix(ctx, gitRoot, "implement", "--spec", specSlug)
	if err != nil {
		return delivery.RunResult{}, fmt.Errorf("start Implement executor: %w", err)
	}
	run, found, err := workflow.latestImplementRun(ctx, workflow.loaded.GitRoot, specSlug)
	if err != nil {
		return delivery.RunResult{}, err
	}
	if result.exitCode != exitOK {
		if found && run.State == store.StateUnresolved {
			return delivery.RunResult{RunID: run.ID, Outcome: delivery.RunOutcomeUnresolved, Reason: strings.TrimSpace(result.stderr)}, nil
		}
		return delivery.RunResult{}, result.failure("roundfix implement")
	}
	head, err := workflow.git.RunGit(ctx, gitRoot, "rev-parse", "HEAD")
	if err != nil {
		return delivery.RunResult{}, fmt.Errorf("read Implement candidate head: %w", err)
	}
	runID := ""
	if found {
		runID = run.ID
	}
	return delivery.RunResult{
		RunID:            runID,
		Outcome:          delivery.RunOutcomeClean,
		CandidateCommits: []string{strings.TrimSpace(head)},
	}, nil
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
	if strings.TrimSpace(record.HeadCommit) != strings.TrimSpace(head) {
		return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeBlocked, Head: record.HeadCommit, Reason: "review record names a different head"}, nil
	}
	switch record.Outcome {
	case reviewOutcomeReviewed:
		return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeReviewed, Head: record.HeadCommit}, nil
	case reviewOutcomeFindings:
		return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeFindings, Head: record.HeadCommit, Reason: record.Findings}, nil
	case reviewOutcomeBlocked:
		return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeBlocked, Head: record.HeadCommit, Reason: record.Reason}, nil
	default:
		return delivery.ReviewResult{}, result.failure("roundfix review")
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
