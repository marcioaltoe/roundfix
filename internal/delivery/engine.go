package delivery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"roundfix/internal/store"
)

const (
	BlockerRunUnresolved          = "run-unresolved"
	BlockerReviewFindings         = "review-findings"
	BlockerReviewBlocked          = "review-blocked"
	BlockerReviewStale            = "review-stale"
	BlockerCorrectiveSpecRequired = "corrective-spec-required"
	BlockerGateFailed             = "gate-failed"
	BlockerChecksFailed           = "checks-failed"
	BlockerChecksTimeout          = "checks-timeout"
	BlockerUnauthorized           = "unauthorized"
	BlockerDeliveryError          = "delivery-error"
	BlockerItemWorktreeMissing    = "item-worktree-missing"
	BlockerRunBudgetExceeded      = "run-budget-exceeded"
	BlockerQueueDeadline          = "queue-deadline"
	BlockerRevalidationFailed     = "revalidation-failed"
)

const WarningPremiseChanged = "premise-changed"

const cleanupWarningPrefix = "warning: cleanup failed"

// ErrItemWorktreeMissing reports that neither an item's recorded worktree nor
// its recorded branch remains available for resume.
var ErrItemWorktreeMissing = errors.New("item worktree and branch are missing")

const (
	defaultCheckTimeout  = 30 * time.Minute
	defaultCheckInterval = 15 * time.Second
)

type RunOutcome string

const (
	RunOutcomeClean          RunOutcome = "clean"
	RunOutcomeUnresolved     RunOutcome = "unresolved"
	RunOutcomeBudgetExceeded RunOutcome = "budget-exceeded"
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
	Outcome       ReviewOutcome
	Head          string
	Reason        string
	ArchivedSpecs []string
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

type ItemState struct {
	Archived        bool
	UnfinishedTasks []string
	Head            string
}

type CarryForwardResult struct {
	RunID   string
	Carried []string
}

type Revalidation struct {
	Findings        []string // unique, sorted finding codes
	ChangedPremises []string // unique, sorted repository-relative paths
	ChangedBy       []string // the prior merge commits that changed at least one of them, in queue order
}

type ItemRevalidator interface {
	Revalidate(ctx context.Context, workDir, specSlug string, priorMerges []string) (Revalidation, error)
}

type ItemRecovery interface {
	InspectItem(ctx context.Context, workDir, specSlug string) (ItemState, error)
	CarryForward(ctx context.Context, workDir, specSlug, branch, runID string) (CarryForwardResult, error)
}

type RetryResult struct {
	SpecSlug      string
	Blocker       string
	Stage         store.DeliveryStage
	CarriedFrom   CarryForwardResult
	OwnerPID      int
	OwnerIdentity string
}

type CandidateRunner interface {
	RunSpec(ctx context.Context, gitRoot, specSlug string) (RunResult, error)
}

type ItemWorkspace interface {
	CreateItemBranch(ctx context.Context, gitRoot, specSlug string) (branch, worktree string, err error)
	UseItemBranch(ctx context.Context, gitRoot, specSlug, branch, worktree string, provisioned bool) (string, error)
	ReleaseMergedRuns(ctx context.Context, gitRoot string, item store.DeliveryQueueItem) error
	RemoveItemBranch(ctx context.Context, gitRoot, branch, worktree string) error
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
	Recovery      ItemRecovery
	Revalidator   ItemRevalidator
	Log           io.Writer
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
	recovery      ItemRecovery
	revalidator   ItemRevalidator
	log           io.Writer
	clock         Clock
	sleeper       Sleeper
	checkTimeout  time.Duration
	checkInterval time.Duration
}

type EngineResult struct {
	Items []store.DeliveryQueueItem
}

func NewEngine(runStore *store.Store, dependencies EngineDependencies) *Engine {
	log := dependencies.Log
	if log == nil {
		log = io.Discard
	}
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
		recovery:      dependencies.Recovery,
		revalidator:   dependencies.Revalidator,
		log:           log,
		clock:         clock,
		sleeper:       sleeper,
		checkTimeout:  checkTimeout,
		checkInterval: checkInterval,
	}
}

func (engine *Engine) Retry(ctx context.Context, gitRoot, specSlug string) (RetryResult, error) {
	if err := engine.validateRetry(); err != nil {
		return RetryResult{}, err
	}
	gitRoot = strings.TrimSpace(gitRoot)
	specSlug = strings.TrimSpace(specSlug)
	if gitRoot == "" {
		return RetryResult{}, errors.New("retry Delivery Queue item: Git root is required")
	}
	if specSlug == "" {
		return RetryResult{}, errors.New("retry Delivery Queue item: Spec slug is required")
	}

	queue, found, err := engine.store.DeliveryQueue(ctx, gitRoot)
	if err != nil {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: read queue: %w", specSlug, err)
	}
	if !found {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: Delivery Queue for repository %q does not exist", specSlug, gitRoot)
	}
	var item store.DeliveryQueueItem
	itemFound := false
	for _, candidate := range queue.Items {
		if candidate.SpecSlug == specSlug {
			item = candidate
			itemFound = true
			break
		}
	}
	if !itemFound {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: Delivery Queue does not contain the item", specSlug)
	}
	if item.Stage != store.DeliveryStageParked {
		return RetryResult{}, fmt.Errorf(
			"retry Delivery Queue item %q: item has stage %q and blocker %q; want stage %q",
			specSlug,
			item.Stage,
			item.Blocker,
			store.DeliveryStageParked,
		)
	}
	if item.Blocker == BlockerQueueDeadline {
		return RetryResult{}, fmt.Errorf(
			"retry Delivery Queue item %q: queue deadline %s was reached; start a new queue with roundfix deliver start",
			specSlug,
			queue.Limits.Deadline.UTC().Format(time.RFC3339),
		)
	}
	if queue.Limits.MaxRetries != 0 && item.RetryCount >= queue.Limits.MaxRetries {
		return RetryResult{}, fmt.Errorf(
			"retry Delivery Queue item %q: retry count %d reached queue limit %d: %w",
			specSlug,
			item.RetryCount,
			queue.Limits.MaxRetries,
			store.ErrDeliveryRetryLimit,
		)
	}

	workDir, err := engine.workspace.UseItemBranch(
		ctx,
		gitRoot,
		item.SpecSlug,
		item.Branch,
		item.Worktree,
		item.WorktreeProvisioned,
	)
	if err != nil {
		if errors.Is(err, ErrItemWorktreeMissing) {
			return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: item branch %q is missing: %w", specSlug, item.Branch, err)
		}
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: use item branch %q: %w", specSlug, item.Branch, err)
	}
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: item worktree is empty", specSlug)
	}

	state, err := engine.recovery.InspectItem(ctx, workDir, item.SpecSlug)
	if err != nil {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: inspect item: %w", specSlug, err)
	}
	result := RetryResult{SpecSlug: item.SpecSlug, Blocker: item.Blocker}
	correctivePrefix := BlockerCorrectiveSpecRequired + ":"
	if strings.HasPrefix(item.Blocker, correctivePrefix) {
		archivedSpecs := strings.TrimSpace(strings.TrimPrefix(item.Blocker, correctivePrefix))
		candidate, err := candidateHead(item)
		if err != nil {
			return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: %w", specSlug, err)
		}
		head := strings.TrimSpace(state.Head)
		if head != candidate {
			return RetryResult{}, fmt.Errorf(
				"retry Delivery Queue item %q: archived Specs %s were reviewed at parked candidate head %q, but the item head is %q; author a corrective Spec with its own authorization and QA gate",
				specSlug,
				archivedSpecs,
				candidate,
				head,
			)
		}
		item.Stage = store.DeliveryStageReviewing
	} else if state.Archived {
		candidate, err := candidateHead(item)
		if err != nil {
			return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: %w", specSlug, err)
		}
		head := strings.TrimSpace(state.Head)
		if head != candidate {
			return RetryResult{}, fmt.Errorf(
				"retry Delivery Queue item %q: archived item head %q differs from candidate head %q",
				specSlug,
				head,
				candidate,
			)
		}
		if strings.TrimSpace(item.PullRequestNumber) == "" {
			item.Stage = store.DeliveryStageGating
		} else {
			item.Stage = store.DeliveryStageChecking
		}
	} else {
		if strings.TrimSpace(item.RunID) == "" {
			revalidation, err := engine.revalidator.Revalidate(ctx, workDir, item.SpecSlug, nil)
			if err != nil {
				return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: revalidate item: %w", specSlug, err)
			}
			if len(revalidation.Findings) > 0 {
				return RetryResult{}, fmt.Errorf(
					"retry Delivery Queue item %q: revalidation failed: %s",
					specSlug,
					strings.Join(revalidation.Findings, ", "),
				)
			}
		}
		carried, err := engine.recovery.CarryForward(
			ctx,
			workDir,
			item.SpecSlug,
			item.Branch,
			item.RunID,
		)
		if err != nil {
			return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: carry forward: %w", specSlug, err)
		}
		carried.RunID = strings.TrimSpace(carried.RunID)
		item.RunID = carried.RunID
		result.CarriedFrom = carried

		state, err = engine.recovery.InspectItem(ctx, workDir, item.SpecSlug)
		if err != nil {
			return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: inspect item after carry-forward: %w", specSlug, err)
		}
		if len(state.UnfinishedTasks) > 0 {
			item.Stage = store.DeliveryStageRunning
		} else {
			head := strings.TrimSpace(state.Head)
			if head == "" {
				return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: current item head is empty", specSlug)
			}
			if len(item.CandidateCommits) == 0 || strings.TrimSpace(item.CandidateCommits[len(item.CandidateCommits)-1]) != head {
				item.CandidateCommits = append(item.CandidateCommits, head)
			}
			item.Stage = store.DeliveryStageReviewing
		}
	}

	ownerPID, ownerIdentity, err := engine.store.RetryDeliveryQueueItem(ctx, gitRoot, item, result.Blocker)
	if err != nil {
		return RetryResult{}, fmt.Errorf("retry Delivery Queue item %q: %w", specSlug, err)
	}
	result.Stage = item.Stage
	result.OwnerPID = ownerPID
	result.OwnerIdentity = ownerIdentity
	return result, nil
}

func (engine *Engine) validateRetry() error {
	switch {
	case engine == nil:
		return errors.New("retry Delivery Queue item: engine is required")
	case engine.store == nil:
		return errors.New("retry Delivery Queue item: store is required")
	case engine.workspace == nil:
		return errors.New("retry Delivery Queue item: item workspace is required")
	case engine.recovery == nil:
		return errors.New("retry Delivery Queue item: item recovery is required")
	case engine.revalidator == nil:
		return errors.New("retry Delivery Queue item: item revalidator is required")
	default:
		return nil
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
		if item.Stage == store.DeliveryStageParked {
			continue
		}
		if item.Stage == store.DeliveryStageQueued &&
			!queue.Limits.Deadline.IsZero() &&
			!engine.clock.Now().Before(queue.Limits.Deadline) {
			if err := engine.park(ctx, gitRoot, &item, BlockerQueueDeadline); err != nil {
				return EngineResult{}, fmt.Errorf("park Spec %q at queue deadline: %w", item.SpecSlug, err)
			}
			queue.Items[index] = item
			continue
		}
		if item.Stage != store.DeliveryStageMerged {
			priorMerges := mergedCommitsBefore(queue.Items, index)
			if err := engine.advanceItem(ctx, gitRoot, &item, priorMerges); err != nil {
				if ctx.Err() != nil {
					return EngineResult{}, fmt.Errorf("deliver Spec %q: %w", item.SpecSlug, err)
				}
				reason := BlockerDeliveryError + ": " + err.Error()
				if parkErr := engine.park(ctx, gitRoot, &item, reason); parkErr != nil {
					return EngineResult{}, fmt.Errorf("deliver Spec %q: %w", item.SpecSlug, errors.Join(err, parkErr))
				}
			}
		}
		if item.Stage == store.DeliveryStageMerged {
			releaseErr := engine.workspace.ReleaseMergedRuns(ctx, gitRoot, item)
			removeErr := engine.workspace.RemoveItemBranch(ctx, gitRoot, item.Branch, item.Worktree)
			if cleanupErr := errors.Join(releaseErr, removeErr); cleanupErr != nil {
				item.Blocker = cleanupWarningPrefix + ": " + cleanupErr.Error()
				if persistErr := engine.store.UpdateDeliveryQueueItem(ctx, gitRoot, item); persistErr != nil {
					return EngineResult{}, fmt.Errorf(
						"record cleanup warning for merged Spec %q after %v: %w",
						item.SpecSlug,
						cleanupErr,
						persistErr,
					)
				}
			} else if strings.HasPrefix(item.Blocker, cleanupWarningPrefix) {
				item.Blocker = ""
				if err := engine.store.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
					return EngineResult{}, fmt.Errorf("clear cleanup warning for merged Spec %q: %w", item.SpecSlug, err)
				}
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
	case engine.revalidator == nil:
		return errors.New("run Delivery Engine: item revalidator is required")
	default:
		return nil
	}
}

func (engine *Engine) advanceItem(
	ctx context.Context,
	gitRoot string,
	item *store.DeliveryQueueItem,
	priorMerges []string,
) error {
	if item.Stage == store.DeliveryStageQueued {
		branch, itemWorktree, err := engine.workspace.CreateItemBranch(ctx, gitRoot, item.SpecSlug)
		if err != nil {
			return fmt.Errorf("create item worktree: %w", err)
		}
		item.Branch = strings.TrimSpace(branch)
		if item.Branch == "" {
			return errors.New("create item worktree: branch is empty")
		}
		item.Worktree = strings.TrimSpace(itemWorktree)
		if item.Worktree == "" {
			return errors.New("create item worktree: path is empty")
		}
		item.WorktreeProvisioned = true
		revalidation, err := engine.revalidator.Revalidate(ctx, item.Worktree, item.SpecSlug, priorMerges)
		if len(revalidation.ChangedPremises) > 0 {
			item.Warning = fmt.Sprintf(
				"%s: %s (merge %s)",
				WarningPremiseChanged,
				strings.Join(revalidation.ChangedPremises, ", "),
				strings.Join(revalidation.ChangedBy, ", "),
			)
			fmt.Fprintf(engine.log, "roundfix: warning: Delivery Queue item %s: %s\n", item.SpecSlug, item.Warning)
		}
		if err != nil {
			return fmt.Errorf("revalidate item: %w", err)
		}
		if len(revalidation.Findings) > 0 {
			return engine.park(
				ctx,
				gitRoot,
				item,
				BlockerRevalidationFailed+": "+strings.Join(revalidation.Findings, ", "),
			)
		}
		if err := engine.setStage(ctx, gitRoot, item, store.DeliveryStageRunning); err != nil {
			return err
		}
	} else {
		if strings.TrimSpace(item.Branch) == "" {
			return errors.New("recorded item branch is missing")
		}
		itemWorktree, err := engine.workspace.UseItemBranch(
			ctx,
			gitRoot,
			item.SpecSlug,
			item.Branch,
			item.Worktree,
			item.WorktreeProvisioned,
		)
		if err != nil {
			if errors.Is(err, ErrItemWorktreeMissing) {
				return engine.park(ctx, gitRoot, item, BlockerItemWorktreeMissing)
			}
			return fmt.Errorf("use item branch %q: %w", item.Branch, err)
		}
		item.Worktree = strings.TrimSpace(itemWorktree)
		if item.Worktree == "" {
			return errors.New("use item branch: worktree is empty")
		}
		item.WorktreeProvisioned = true
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

func mergedCommitsBefore(items []store.DeliveryQueueItem, index int) []string {
	commits := make([]string, 0, index)
	for _, item := range items[:index] {
		mergeCommit := strings.TrimSpace(item.MergeCommit)
		if item.Stage == store.DeliveryStageMerged && mergeCommit != "" {
			commits = append(commits, mergeCommit)
		}
	}
	return commits
}

func (engine *Engine) runCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	result, err := engine.runner.RunSpec(ctx, workDir, item.SpecSlug)
	if err != nil {
		return fmt.Errorf("run Implement executor: %w", err)
	}
	item.RunID = strings.TrimSpace(result.RunID)
	switch result.Outcome {
	case RunOutcomeUnresolved:
		return engine.park(ctx, gitRoot, item, BlockerRunUnresolved)
	case RunOutcomeBudgetExceeded:
		return engine.park(ctx, gitRoot, item, BlockerRunBudgetExceeded)
	case RunOutcomeClean:
		if len(result.CandidateCommits) == 0 || strings.TrimSpace(result.CandidateCommits[len(result.CandidateCommits)-1]) == "" {
			return errors.New("clean Run returned no candidate head")
		}
	default:
		return fmt.Errorf("Implement executor returned invalid outcome %q", result.Outcome)
	}
	item.CandidateCommits = append([]string(nil), result.CandidateCommits...)
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageReviewing)
}

func (engine *Engine) reviewCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	policy, err := engine.reviewer.ReviewPolicy(ctx, workDir, item.SpecSlug)
	if err != nil {
		return fmt.Errorf("resolve Pre-PR Review Policy: %w", err)
	}
	switch policy {
	case ReviewPolicyNone:
		if err := engine.reviewer.RecordReviewOmission(ctx, workDir, item.SpecSlug, head); err != nil {
			return fmt.Errorf("record configured review omission: %w", err)
		}
	case ReviewPolicyEnabled:
		result, err := engine.reviewer.Review(ctx, workDir, item.SpecSlug, head)
		if err != nil {
			return fmt.Errorf("run pre-PR review: %w", err)
		}
		switch result.Outcome {
		case ReviewOutcomeFindings:
			if len(result.ArchivedSpecs) > 0 {
				archivedSpecs := slices.Clone(result.ArchivedSpecs)
				slices.Sort(archivedSpecs)
				return engine.park(
					ctx,
					gitRoot,
					item,
					BlockerCorrectiveSpecRequired+": "+strings.Join(archivedSpecs, ", "),
				)
			}
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
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	reviewedHead, err := candidateHead(*item)
	if err != nil {
		return err
	}
	result, err := engine.archiver.Archive(ctx, workDir, item.SpecSlug, reviewedHead)
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
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	result, err := engine.gate.Gate(ctx, workDir, item.SpecSlug, head)
	if err != nil {
		return fmt.Errorf("run repository gate: %w", err)
	}
	if !result.Passed {
		return engine.park(ctx, gitRoot, item, BlockerGateFailed)
	}
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStagePublishing)
}

func (engine *Engine) publishCandidate(ctx context.Context, gitRoot string, item *store.DeliveryQueueItem) error {
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	authorized, err := engine.authorized(ctx, workDir, item.SpecSlug)
	if err != nil {
		return err
	}
	if !authorized {
		return engine.park(ctx, gitRoot, item, BlockerUnauthorized)
	}
	publication, err := engine.publication.Publication(ctx, workDir, item.SpecSlug, item.Branch)
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
	if err := engine.push(ctx, gitRoot, workDir, item.SpecSlug, publication, head); err != nil {
		return err
	}
	pullRequest, err := engine.createPullRequest(ctx, gitRoot, workDir, item.SpecSlug, publication, head)
	if err != nil {
		return err
	}
	item.PullRequestNumber = pullRequest.Number
	return engine.setStage(ctx, gitRoot, item, store.DeliveryStageChecking)
}

func (engine *Engine) push(ctx context.Context, gitRoot, workDir, specSlug string, publication Publication, expectedHead string) error {
	pullRequests := engine.pullRequests.WithWorkDir(workDir)
	intent, unmatched, err := engine.actionIntent(ctx, gitRoot, specSlug, store.DeliveryActionPush)
	if err != nil {
		return err
	}
	if unmatched {
		observed, found, err := pullRequests.RemoteBranchHead(ctx, publication.Remote, publication.HeadBranch)
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

	remoteHead, err := pullRequests.PushBranch(ctx, publication.Remote, publication.HeadBranch, expectedHead)
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
	workDir string,
	specSlug string,
	publication Publication,
	expectedHead string,
) (PullRequest, error) {
	pullRequests := engine.pullRequests.WithWorkDir(workDir)
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
	result, err := pullRequests.FindOrCreatePullRequest(ctx, PullRequestRequest{
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
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	pullRequests := engine.pullRequests.WithWorkDir(workDir)
	head, err := candidateHead(*item)
	if err != nil {
		return err
	}
	deadline := engine.clock.Now().Add(engine.checkTimeout)
	for {
		report, err := pullRequests.CurrentHeadChecks(ctx, item.PullRequestNumber)
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
	workDir, err := itemWorkDir(*item)
	if err != nil {
		return err
	}
	authorized, err := engine.authorized(ctx, workDir, item.SpecSlug)
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
	result, err := engine.pullRequests.WithWorkDir(workDir).MergePullRequest(ctx, item.PullRequestNumber, head)
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

func (engine *Engine) authorized(ctx context.Context, workDir, specSlug string) (bool, error) {
	authorization, err := engine.authorizer.Authorization(ctx, workDir, specSlug)
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

func itemWorkDir(item store.DeliveryQueueItem) (string, error) {
	workDir := strings.TrimSpace(item.Worktree)
	if workDir == "" {
		return "", errors.New("recorded item worktree is missing")
	}
	return workDir, nil
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
