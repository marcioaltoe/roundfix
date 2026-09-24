package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
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

func (workflow *commandDeliveryWorkflow) CreateItemBranch(ctx context.Context, gitRoot, specSlug string) (string, error) {
	branch, err := newDeliveryBranch(specSlug)
	if err != nil {
		return "", err
	}
	branch, err = workflow.store.RecordDeliveryQueueItemBranch(ctx, gitRoot, specSlug, branch)
	if err != nil {
		return "", fmt.Errorf("record item branch: %w", err)
	}
	state, err := preflight.InspectGit(ctx, gitRoot, workflow.git)
	if err != nil {
		return "", fmt.Errorf("inspect checkout before item branch creation: %w", err)
	}
	if len(state.Dirty) != 0 {
		return "", errors.New("create item branch: checkout has uncommitted changes")
	}
	if state.Branch == branch {
		return branch, nil
	}
	exists, err := localItemBranchExists(ctx, workflow.git, gitRoot, branch)
	if err != nil {
		return "", fmt.Errorf("inspect item branch %q: %w", branch, err)
	}
	if exists {
		if _, err := workflow.git.RunGit(ctx, gitRoot, "switch", branch); err != nil {
			return "", fmt.Errorf("reuse item branch %q: %w", branch, err)
		}
		return branch, nil
	}
	defaultBranch := preflight.DetectDefaultBranch(ctx, gitRoot, state.Branch, workflow.git)
	if defaultBranch.Source == preflight.DefaultBranchUndetermined {
		return "", errors.New("create item branch: repository default branch is unknown")
	}
	remote := strings.TrimSpace(workflow.loaded.Config.Watch.PushRemote)
	if remote == "" {
		remote = "origin"
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, "fetch", remote, defaultBranch.Name); err != nil {
		return "", fmt.Errorf("refresh default branch %q: %w", defaultBranch.Name, err)
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, "switch", "--no-track", "-c", branch, remote+"/"+defaultBranch.Name); err != nil {
		return "", fmt.Errorf("create item branch %q: %w", branch, err)
	}
	return branch, nil
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

func (workflow *commandDeliveryWorkflow) UseItemBranch(ctx context.Context, gitRoot, branch string) error {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return errors.New("use item branch: branch is required")
	}
	state, err := preflight.InspectGit(ctx, gitRoot, workflow.git)
	if err != nil {
		return fmt.Errorf("inspect checkout before selecting item branch: %w", err)
	}
	if len(state.Dirty) != 0 {
		return errors.New("use item branch: checkout has uncommitted changes")
	}
	if state.Branch == branch {
		return nil
	}
	if _, err := workflow.git.RunGit(ctx, gitRoot, "switch", branch); err != nil {
		return fmt.Errorf("switch to item branch %q: %w", branch, err)
	}
	return nil
}

func (workflow *commandDeliveryWorkflow) RunSpec(ctx context.Context, gitRoot, specSlug string) (delivery.RunResult, error) {
	result, err := workflow.runRoundfix(ctx, gitRoot, "implement", "--spec", specSlug)
	if err != nil {
		return delivery.RunResult{}, fmt.Errorf("start Implement executor: %w", err)
	}
	run, found, err := workflow.latestImplementRun(ctx, gitRoot, specSlug)
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
	if before.HEAD != strings.TrimSpace(reviewedHead) {
		return delivery.ArchiveResult{Parent: before.HEAD}, nil
	}
	if len(before.Dirty) != 0 {
		return delivery.ArchiveResult{Parent: before.HEAD}, nil
	}
	result, err := workflow.runRoundfix(ctx, gitRoot, "archive", specSlug)
	if err != nil {
		return delivery.ArchiveResult{}, fmt.Errorf("start Archive Command: %w", err)
	}
	if result.exitCode != exitOK {
		return delivery.ArchiveResult{}, result.failure("roundfix archive")
	}

	source, destination, err := workflow.archivePaths(specSlug)
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

func (workflow *commandDeliveryWorkflow) archivePaths(specSlug string) (string, string, error) {
	specsRoot, err := roundconfig.ResolveSpecsRoot(workflow.loaded, workflow.loaded.GitRoot)
	if err != nil {
		return "", "", err
	}
	source, err := filepath.Rel(workflow.loaded.GitRoot, filepath.Join(specsRoot.Path, specSlug))
	if err != nil {
		return "", "", fmt.Errorf("resolve active Spec path: %w", err)
	}
	destinationRoot := filepath.Join(specsRoot.Path, "_archived")
	if specsRoot.BuiltInRoot {
		destinationRoot = filepath.Join(workflow.loaded.GitRoot, filepath.FromSlash(spec.ArchiveDir(spec.ArchiveKindSpec)))
	}
	destination, err := filepath.Rel(workflow.loaded.GitRoot, filepath.Join(destinationRoot, specSlug))
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
