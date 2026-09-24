package delivery

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"roundfix/internal/store"
)

const (
	BlockerRunUnresolved  = "run-unresolved"
	BlockerReviewFindings = "review-findings"
	BlockerReviewBlocked  = "review-blocked"
	BlockerReviewStale    = "review-stale"
	BlockerGateFailed     = "gate-failed"
	BlockerChecksFailed   = "checks-failed"
	BlockerChecksTimeout  = "checks-timeout"
	BlockerUnauthorized   = "unauthorized"
	BlockerDeliveryError  = "delivery-error"
)

const (
	defaultCheckTimeout  = 30 * time.Minute
	defaultCheckInterval = 15 * time.Second
)

type RunOutcome string

const (
	RunOutcomeClean      RunOutcome = "clean"
	RunOutcomeUnresolved RunOutcome = "unresolved"
)

type ReviewPolicy string

const (
	ReviewPolicyEnabled ReviewPolicy = "enabled"
	ReviewPolicyNone    ReviewPolicy = "none"
)

type ReviewOutcome string

const (
	ReviewOutcomeReviewed ReviewOutcome = "reviewed"
	ReviewOutcomeFindings ReviewOutcome = "findings"
	ReviewOutcomeBlocked  ReviewOutcome = "blocked"
)

type RunResult struct {
	RunID            string
	Outcome          RunOutcome
	CandidateCommits []string
	Reason           string
}

type ReviewResult struct {
	Outcome ReviewOutcome
	Head    string
	Reason  string
}

type ArchiveResult struct {
	Parent        string
	Head          string
	ExactSpecMove bool
}

type GateResult struct {
	Passed bool
	Reason string
}

type Authorization struct {
	Operations []string
}

type Publication struct {
	Remote     string
	HeadBranch string
	BaseBranch string
	Title      string
	Body       string
}

type CandidateRunner interface {
	RunSpec(ctx context.Context, gitRoot, specSlug string) (RunResult, error)
}

type ItemWorkspace interface {
	CreateItemBranch(ctx context.Context, gitRoot, specSlug string) (string, error)
	UseItemBranch(ctx context.Context, gitRoot, branch string) error
}

type PrePRReviewer interface {
	ReviewPolicy(ctx context.Context, gitRoot, specSlug string) (ReviewPolicy, error)
	Review(ctx context.Context, gitRoot, specSlug, head string) (ReviewResult, error)
	RecordReviewOmission(ctx context.Context, gitRoot, specSlug, head string) error
}

type CandidateArchiver interface {
	Archive(ctx context.Context, gitRoot, specSlug, reviewedHead string) (ArchiveResult, error)
}

type RepositoryGate interface {
	Gate(ctx context.Context, gitRoot, specSlug, head string) (GateResult, error)
}

type AuthorizationReader interface {
	Authorization(ctx context.Context, gitRoot, specSlug string) (Authorization, error)
}

type PublicationPlanner interface {
	Publication(ctx context.Context, gitRoot, specSlug, branch string) (Publication, error)
}

type Clock interface {
	Now() time.Time
}

type Sleeper interface {
	Sleep(context.Context, time.Duration) error
}

type EngineDependencies struct {
	Workspace     ItemWorkspace
	Runner        CandidateRunner
	Reviewer      PrePRReviewer
	Archiver      CandidateArchiver
	Gate          RepositoryGate
	Authorizer    AuthorizationReader
	Publication   PublicationPlanner
	PullRequests  PullRequestBoundary
	Clock         Clock
	Sleeper       Sleeper
	CheckTimeout  time.Duration
	CheckInterval time.Duration
}

type Engine struct {
	store         *store.Store
	workspace     ItemWorkspace
	runner        CandidateRunner
	reviewer      PrePRReviewer
	archiver      CandidateArchiver
	gate          RepositoryGate
	authorizer    AuthorizationReader
	publication   PublicationPlanner
	pullRequests  PullRequestBoundary
	clock         Clock
	sleeper       Sleeper
	checkTimeout  time.Duration
	checkInterval time.Duration
}

type EngineResult struct {
	Items []store.DeliveryQueueItem
}

func NewEngine(runStore *store.Store, dependencies EngineDependencies) *Engine {
	clock := dependencies.Clock
	if clock == nil {
		clock = realClock{}
	}
	sleeper := dependencies.Sleeper
	if sleeper == nil {
		sleeper = realSleeper{}
	}
	checkTimeout := dependencies.CheckTimeout
	if checkTimeout <= 0 {
		checkTimeout = defaultCheckTimeout
	}
	checkInterval := dependencies.CheckInterval
	if checkInterval <= 0 {
		checkInterval = defaultCheckInterval
	}
	return &Engine{
		store:         runStore,
		workspace:     dependencies.Workspace,
		runner:        dependencies.Runner,
		reviewer:      dependencies.Reviewer,
		archiver:      dependencies.Archiver,
		gate:          dependencies.Gate,
		authorizer:    dependencies.Authorizer,
		publication:   dependencies.Publication,
		pullRequests:  dependencies.PullRequests,
		clock:         clock,
		sleeper:       sleeper,
		checkTimeout:  checkTimeout,
		checkInterval: checkInterval,
	}
}

func (engine *Engine) Run(ctx context.Context, gitRoot string) (EngineResult, error) {
	if err := engine.validate(); err != nil {
		return EngineResult{}, err
	}
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return EngineResult{}, errors.New("run Delivery Engine: Git root is required")
	}
	queue, found, err := engine.store.DeliveryQueue(ctx, gitRoot)
	if err != nil {
		return EngineResult{}, fmt.Errorf("run Delivery Engine: read queue: %w", err)
	}
	if !found {
		return EngineResult{}, fmt.Errorf("run Delivery Engine: Delivery Queue for repository %q does not exist", gitRoot)
	}

	for index := range queue.Items {
		item := queue.Items[index]
		if item.Stage == store.DeliveryStageMerged || item.Stage == store.DeliveryStageParked {
			continue
		}
		if err := engine.advanceItem(ctx, gitRoot, &item); err != nil {
			if ctx.Err() != nil {
				return EngineResult{}, fmt.Errorf("deliver Spec %q: %w", item.SpecSlug, err)
			}
			reason := BlockerDeliveryError + ": " + err.Error()
			if parkErr := engine.park(ctx, gitRoot, &item, reason); parkErr != nil {
				return EngineResult{}, fmt.Errorf("deliver Spec %q: %w", item.SpecSlug, errors.Join(err, parkErr))
			}
		}
		queue.Items[index] = item
	}
	return EngineResult{Items: queue.Items}, nil
}

func (engine *Engine) validate() error {
	switch {
	case engine == nil:
		return errors.New("run Delivery Engine: engine is required")
	case engine.store == nil:
		return errors.New("run Delivery Engine: store is required")
	case engine.workspace == nil:
		return errors.New("run Delivery Engine: item workspace is required")
	case engine.runner == nil:
		return errors.New("run Delivery Engine: candidate runner is required")
	case engine.reviewer == nil:
		return errors.New("run Delivery Engine: reviewer is required")
	case engine.archiver == nil:
		return errors.New("run Delivery Engine: archiver is required")
	case engine.gate == nil:
		return errors.New("run Delivery Engine: repository gate is required")
	case engine.authorizer == nil:
		return errors.New("run Delivery Engine: authorization reader is required")
	case engine.publication == nil:
		return errors.New("run Delivery Engine: publication planner is required")
	case engine.pullRequests == nil:
		return errors.New("run Delivery Engine: pull request boundary is required")
	default:
		return nil
	}
}

func (engine *Engine) advanceItem(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	if item.Stage == store.DeliveryStageQueued {
		branch, err := engine.workspace.CreateItemBranch(ctx, gitRoot, item.SpecSlug)
		if err != nil {
			return fmt.Errorf("create item branch: %w", err)
		}
		item.Branch = strings.TrimSpace(branch)
		if item.Branch == "" {
			return errors.New("create item branch: branch is empty")
		}
		if err := engine.setStage(ctx, gitRoot, item, store.DeliveryStageRunning); err != nil {
			return err
		}
	} else {
		if strings.TrimSpace(item.Branch) == "" {
			return errors.New("recorded item branch is missing")
		}
		if err := engine.workspace.UseItemBranch(ctx, gitRoot, item.Branch); err != nil {
			return fmt.Errorf("use item branch %q: %w", item.Branch, err)
		}
	}
	for item.Stage != store.DeliveryStageMerged && item.Stage != store.DeliveryStageParked {
		switch item.Stage {
		case store.DeliveryStageRunning:
			if err := engine.runCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStageReviewing:
			if err := engine.reviewCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStageArchiving:
			if err := engine.archiveCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStageGating:
			if err := engine.gateCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStagePublishing:
			if err := engine.publishCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStageChecking:
			if err := engine.checkCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		case store.DeliveryStageMerging:
			if err := engine.mergeCandidate(ctx, gitRoot, item); err != nil {
				return err
			}
		default:
			return fmt.Errorf("stage %q is not resumable", item.Stage)
		}
	}
	return nil
}

func (engine *Engine) runCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	result, err := engine.runner.RunSpec(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		return fmt.Errorf("run Implement executor: %w", err)
	}
	switch result.Outcome {
	case RunOutcomeUnresolved:
		return engine.park(ctx, gitRoot, item, BlockerRunUnresolved)
	case RunOutcomeClean:
		if len(result.CandidateCommits) == 0 || strings.TrimSpace(result.CandidateCommits[len(result.CandidateCommits)-1]) == "" {
			return errors.New("clean Run returned no candidate head")
		}
	default:
		return fmt.Errorf("Implement executor returned invalid outcome %q", result.Outcome)
	}
	item.RunID = strings.TrimSpace(result.RunID)
	item.CandidateCommits = append([]string(nil), result.CandidateCommits...)
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageReviewing)
}

func (engine *Engine) reviewCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	policy, err := engine.reviewer.ReviewPolicy(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		return fmt.Errorf("resolve Pre-PR Review Policy: %w", err)
	}
	switch policy {
	case ReviewPolicyNone:
		if err := engine.reviewer.RecordReviewOmission(ctx, gitRoot, item.SpecSlug, head); err != nil {
			return fmt.Errorf("record configured review omission: %w", err)
		}
	case ReviewPolicyEnabled:
		result, err := engine.reviewer.Review(ctx, gitRoot, item.SpecSlug, head)
		if err != nil {
			return fmt.Errorf("run pre-PR review: %w", err)
		}
		switch result.Outcome {
		case ReviewOutcomeFindings:
			return engine.park(ctx, gitRoot, item, BlockerReviewFindings)
		case ReviewOutcomeBlocked:
			return engine.park(ctx, gitRoot, item, BlockerReviewBlocked)
		case ReviewOutcomeReviewed:
			if strings.TrimSpace(result.Head) != head {
				return engine.park(ctx, gitRoot, item, BlockerReviewStale)
			}
		default:
			return fmt.Errorf("pre-PR review returned invalid outcome %q", result.Outcome)
		}
	default:
		return fmt.Errorf("Pre-PR Review Policy %q is invalid", policy)
	}
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageArchiving)
}

func (engine *Engine) archiveCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	reviewedHead, err := candidateHead(*item)
	if err != nil {
		return err
	}
	result, err := engine.archiver.Archive(ctx, gitRoot, item.SpecSlug, reviewedHead)
	if err != nil {
		return fmt.Errorf("archive Spec: %w", err)
	}
	result.Parent = strings.TrimSpace(result.Parent)
	result.Head = strings.TrimSpace(result.Head)
	if result.Parent != reviewedHead || result.Head == "" || result.Head == reviewedHead || !result.ExactSpecMove {
		return engine.park(ctx, gitRoot, item, BlockerReviewStale)
	}
	item.CandidateCommits = append(item.CandidateCommits, result.Head)
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageGating)
}

func (engine *Engine) gateCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	result, err := engine.gate.Gate(ctx, gitRoot, item.SpecSlug, head)
	if err != nil {
		return fmt.Errorf("run repository gate: %w", err)
	}
	if !result.Passed {
		return engine.park(ctx, gitRoot, item, BlockerGateFailed)
	}
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStagePublishing)
}

func (engine *Engine) publishCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	authorized, err := engine.authorized(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		return err
	}
	if !authorized {
		return engine.park(ctx, gitRoot, item, BlockerUnauthorized)
	}
	publication, err := engine.publication.Publication(ctx, gitRoot, item.SpecSlug, item.Branch)
	if err != nil {
		return fmt.Errorf("plan publication: %w", err)
	}
	if err := validatePublication(publication); err != nil {
		return err
	}
	if publication.HeadBranch != item.Branch {
		return fmt.Errorf("plan publication: PR Head Branch %q does not match recorded item branch %q", publication.HeadBranch, item.Branch)
	}
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	if err := engine.push(ctx, gitRoot, item.SpecSlug, publication, head); err != nil {
		return err
	}
	pullRequest, err := engine.createPullRequest(ctx, gitRoot, item.SpecSlug, publication, head)
	if err != nil {
		return err
	}
	item.PullRequestNumber = pullRequest.Number
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageChecking)
}

func (engine *Engine) push(ctx context.Context, gitRoot, specSlug string, publication Publication, expectedHead string) error {
	intent, unmatched, err := engine.actionIntent(ctx, gitRoot, specSlug, store.DeliveryActionPush)
	if err != nil {
		return err
	}
	if unmatched {
		observed, found, err := engine.pullRequests.RemoteBranchHead(ctx, publication.Remote, publication.HeadBranch)
		if err != nil {
			return fmt.Errorf("reconcile push intent: %w", err)
		}
		if found {
			if observed.SHA == expectedHead {
				return engine.recordReceipt(ctx, intent.ID, expectedHead)
			}
		}
	} else {
		intent, err = engine.store.RecordDeliveryActionIntent(ctx, gitRoot, specSlug, store.DeliveryActionPush)
		if err != nil {
			return fmt.Errorf("record push intent: %w", err)
		}
	}

	remoteHead, err := engine.pullRequests.PushBranch(ctx, publication.Remote, publication.HeadBranch, expectedHead)
	if err != nil {
		return fmt.Errorf("push candidate: %w", err)
	}
	if remoteHead.SHA != expectedHead {
		return fmt.Errorf("push candidate: remote head is %q, expected %q", remoteHead.SHA, expectedHead)
	}
	return engine.recordReceipt(ctx, intent.ID, remoteHead.SHA)
}

func (engine *Engine) createPullRequest(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	publication Publication,
	expectedHead string,
) (PullRequest, error) {
	intent, unmatched, err := engine.actionIntent(ctx, gitRoot, specSlug, store.DeliveryActionCreatePullRequest)
	if err != nil {
		return PullRequest{}, err
	}
	if !unmatched {
		intent, err = engine.store.RecordDeliveryActionIntent(ctx, gitRoot, specSlug, store.DeliveryActionCreatePullRequest)
		if err != nil {
			return PullRequest{}, fmt.Errorf("record create-pull-request intent: %w", err)
		}
	}
	result, err := engine.pullRequests.FindOrCreatePullRequest(ctx, PullRequestRequest{
		HeadBranch: publication.HeadBranch,
		BaseBranch: publication.BaseBranch,
		Title:      publication.Title,
		Body:       publication.Body,
	})
	if err != nil {
		return PullRequest{}, fmt.Errorf("find or create pull request: %w", err)
	}
	if strings.TrimSpace(result.PullRequest.Number) == "" {
		return PullRequest{}, errors.New("find or create pull request: pull request number is empty")
	}
	if result.PullRequest.HeadBranch != publication.HeadBranch {
		return PullRequest{}, fmt.Errorf(
			"find or create pull request: PR Head Branch is %q, expected recorded item branch %q",
			result.PullRequest.HeadBranch,
			publication.HeadBranch,
		)
	}
	if result.PullRequest.HeadSHA != expectedHead {
		return PullRequest{}, fmt.Errorf("find or create pull request: PR Head Branch is at %q, expected %q", result.PullRequest.HeadSHA, expectedHead)
	}
	ownedByAnotherItem, err := engine.pullRequestOwnedByAnotherItem(ctx, gitRoot, specSlug, result.PullRequest.Number)
	if err != nil {
		return PullRequest{}, err
	}
	if ownedByAnotherItem {
		return PullRequest{}, fmt.Errorf("find or create pull request: pull request %s is recorded by another Delivery Queue item", result.PullRequest.Number)
	}
	if err := engine.recordReceipt(ctx, intent.ID, result.PullRequest.Number); err != nil {
		return PullRequest{}, err
	}
	return result.PullRequest, nil
}

func (engine *Engine) checkCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	deadline := engine.clock.Now().Add(engine.checkTimeout)
	for {
		report, err := engine.pullRequests.CurrentHeadChecks(ctx, item.PullRequestNumber)
		if err == nil {
			if report.HeadSHA != head {
				return engine.park(ctx, gitRoot, item, BlockerReviewStale)
			}
			pending := len(report.Checks) == 0
			for _, check := range report.Checks {
				switch strings.ToLower(strings.TrimSpace(check.Bucket)) {
				case "pass", "skipping":
				case "fail", "cancel", "cancelled":
					return engine.park(ctx, gitRoot, item, BlockerChecksFailed)
				default:
					pending = true
				}
			}
			if !pending {
				return engine.setStage(ctx, gitRoot, item, store.DeliveryStageMerging)
			}
		} else if ctx.Err() != nil {
			return fmt.Errorf("read current-head checks: %w", err)
		}
		remaining := deadline.Sub(engine.clock.Now())
		if remaining <= 0 {
			return engine.park(ctx, gitRoot, item, BlockerChecksTimeout)
		}
		wait := min(engine.checkInterval, remaining)
		if err := engine.sleeper.Sleep(ctx, wait); err != nil {
			return fmt.Errorf("wait for current-head checks: %w", err)
		}
	}
}

func (engine *Engine) mergeCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	authorized, err := engine.authorized(ctx, gitRoot, item.SpecSlug)
	if err != nil {
		return err
	}
	if !authorized {
		return engine.park(ctx, gitRoot, item, BlockerUnauthorized)
	}
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	intent, unmatched, err := engine.actionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionMerge)
	if err != nil {
		return err
	}
	if !unmatched {
		intent, err = engine.store.RecordDeliveryActionIntent(ctx, gitRoot, item.SpecSlug, store.DeliveryActionMerge)
		if err != nil {
			return fmt.Errorf("record merge intent: %w", err)
		}
	}
	result, err := engine.pullRequests.MergePullRequest(ctx, item.PullRequestNumber, head)
	if err != nil {
		var mismatch PullRequestHeadMismatchError
		if errors.As(err, &mismatch) {
			return engine.park(ctx, gitRoot, item, BlockerReviewStale)
		}
		return fmt.Errorf("merge pull request: %w", err)
	}
	if strings.TrimSpace(result.PullRequest.HeadSHA) != head {
		return engine.park(ctx, gitRoot, item, BlockerReviewStale)
	}
	mergeCommit := strings.TrimSpace(result.PullRequest.MergeCommit)
	if mergeCommit == "" {
		return errors.New("merge pull request: merge commit is empty")
	}
	if err := engine.recordReceipt(ctx, intent.ID, mergeCommit); err != nil {
		return err
	}
	item.MergeCommit = mergeCommit
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageMerged)
}

func (engine *Engine) pullRequestOwnedByAnotherItem(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	pullRequestNumber string,
) (bool, error) {
	queue, found, err := engine.store.DeliveryQueue(ctx, gitRoot)
	if err != nil {
		return false, fmt.Errorf("check pull request ownership: %w", err)
	}
	if !found {
		return false, errors.New("check pull request ownership: Delivery Queue does not exist")
	}
	for _, other := range queue.Items {
		if other.SpecSlug != specSlug && other.PullRequestNumber == pullRequestNumber {
			return true, nil
		}
	}
	return false, nil
}

func (engine *Engine) authorized(ctx context.Context, gitRoot, specSlug string) (bool, error) {
	authorization, err := engine.authorizer.Authorization(ctx, gitRoot, specSlug)
	if err != nil {
		return false, fmt.Errorf("read delivery authorization: %w", err)
	}
	for _, required := range []string{"push", "pull_request", "merge"} {
		if !slices.Contains(authorization.Operations, required) {
			return false, nil
		}
	}
	return true, nil
}

func (engine *Engine) actionIntent(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	action store.DeliveryAction,
) (store.DeliveryActionIntent, bool, error) {
	intents, err := engine.store.UnmatchedDeliveryActionIntents(ctx, gitRoot)
	if err != nil {
		return store.DeliveryActionIntent{}, false, fmt.Errorf("list unmatched Delivery Action intents: %w", err)
	}
	var match store.DeliveryActionIntent
	found := false
	for _, intent := range intents {
		if intent.SpecSlug != specSlug || intent.Action != action {
			continue
		}
		if found {
			return store.DeliveryActionIntent{}, false, fmt.Errorf("found more than one unmatched %q intent", action)
		}
		match = intent
		found = true
	}
	return match, found, nil
}

func (engine *Engine) recordReceipt(ctx context.Context, intentID int64, result string) error {
	if _, err := engine.store.RecordDeliveryActionReceipt(ctx, intentID, result); err != nil {
		return fmt.Errorf("record Delivery Action receipt: %w", err)
	}
	return nil
}

func (engine *Engine) setStage(
	ctx context.Context,
	gitRoot string,
	item *store.DeliveryQueueItem,
	stage store.DeliveryStage,
) error {
	item.Stage = stage
	item.Blocker = ""
	if err := engine.store.UpdateDeliveryQueueItem(ctx, gitRoot, *item); err != nil {
		return fmt.Errorf("persist stage %q: %w", stage, err)
	}
	return nil
}

func (engine *Engine) park(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem, blocker string) error {
	item.Stage = store.DeliveryStageParked
	item.Blocker = blocker
	if err := engine.store.UpdateDeliveryQueueItem(ctx, gitRoot, *item); err != nil {
		return fmt.Errorf("park item as %q: %w", blocker, err)
	}
	return nil
}

func candidateHead(item store.DeliveryQueueItem) (string, error) {
	if len(item.CandidateCommits) == 0 {
		return "", errors.New("candidate head is missing")
	}
	head := strings.TrimSpace(item.CandidateCommits[len(item.CandidateCommits)-1])
	if head == "" {
		return "", errors.New("candidate head is empty")
	}
	return head, nil
}

func validatePublication(publication Publication) error {
	switch {
	case strings.TrimSpace(publication.Remote) == "":
		return errors.New("plan publication: remote is required")
	case strings.TrimSpace(publication.HeadBranch) == "":
		return errors.New("plan publication: PR Head Branch is required")
	case strings.TrimSpace(publication.BaseBranch) == "":
		return errors.New("plan publication: base branch is required")
	case strings.TrimSpace(publication.HeadBranch) == strings.TrimSpace(publication.BaseBranch):
		return errors.New("plan publication: PR Head Branch cannot be the base branch")
	case strings.TrimSpace(publication.Title) == "":
		return errors.New("plan publication: pull request title is required")
	default:
		return nil
	}
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

type realSleeper struct{}

func (realSleeper) Sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
