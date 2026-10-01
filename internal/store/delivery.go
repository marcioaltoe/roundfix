package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type DeliveryStage string

const (
	DeliveryStageQueued     DeliveryStage = "queued"
	DeliveryStageRunning    DeliveryStage = "running"
	DeliveryStageReviewing  DeliveryStage = "reviewing"
	DeliveryStageArchiving  DeliveryStage = "archiving"
	DeliveryStageGating     DeliveryStage = "gating"
	DeliveryStagePublishing DeliveryStage = "publishing"
	DeliveryStageChecking   DeliveryStage = "checking"
	DeliveryStageMerging    DeliveryStage = "merging"
	DeliveryStageMerged     DeliveryStage = "merged"
	DeliveryStageParked     DeliveryStage = "parked"
)

type DeliveryAction string

const (
	DeliveryActionPush              DeliveryAction = "push"
	DeliveryActionCreatePullRequest DeliveryAction = "create-pull-request"
	DeliveryActionMerge             DeliveryAction = "merge"
)

// ErrDeliveryRetryLimit reports that a Delivery Queue item has used every
// retry permitted by its queue.
var ErrDeliveryRetryLimit = errors.New("Delivery Queue item retry limit reached")

type DeliveryQueueLimits struct {
	Deadline   time.Time
	MaxRetries int
	MaxTokens  int64
}

type DeliveryQueue struct {
	GitRoot       string
	OwnerPID      int
	OwnerIdentity string
	Limits        DeliveryQueueLimits
	Items         []DeliveryQueueItem
}

type DeliveryQueueItem struct {
	SpecSlug            string
	Position            int
	Stage               DeliveryStage
	Blocker             string
	Branch              string
	Worktree            string
	WorktreeProvisioned bool
	RunID               string
	CandidateCommits    []string
	PullRequestNumber   string
	MergeCommit         string
	RetryCount          int
	Warning             string
}

type DeliveryActionIntent struct {
	ID        int64
	GitRoot   string
	SpecSlug  string
	Action    DeliveryAction
	CreatedAt time.Time
}

type DeliveryActionReceipt struct {
	IntentID  int64
	Result    string
	CreatedAt time.Time
}

// CreateDeliveryQueue records one ordered Delivery Queue for a repository.
// The queue is durable until a later command explicitly removes it.
func (store *Store) CreateDeliveryQueue(ctx context.Context, gitRoot string, specSlugs []string) (DeliveryQueue, error) {
	return store.CreateDeliveryQueueWithLimits(ctx, gitRoot, specSlugs, DeliveryQueueLimits{})
}

// CreateDeliveryQueueWithLimits records one ordered Delivery Queue and the
// limits every process advancing it must enforce.
func (store *Store) CreateDeliveryQueueWithLimits(
	ctx context.Context,
	gitRoot string,
	specSlugs []string,
	limits DeliveryQueueLimits,
) (DeliveryQueue, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return DeliveryQueue{}, errors.New("create Delivery Queue: Git root is required")
	}
	if len(specSlugs) == 0 {
		return DeliveryQueue{}, errors.New("create Delivery Queue: at least one Spec slug is required")
	}
	if limits.MaxRetries < 0 {
		return DeliveryQueue{}, errors.New("create Delivery Queue: Max retries must not be negative")
	}
	if limits.MaxTokens < 0 {
		return DeliveryQueue{}, errors.New("create Delivery Queue: Max tokens must not be negative")
	}
	deadlineUnix := int64(0)
	if !limits.Deadline.IsZero() {
		deadlineUnix = limits.Deadline.UTC().Unix()
		limits.Deadline = time.Unix(deadlineUnix, 0).UTC()
	}

	items := make([]DeliveryQueueItem, len(specSlugs))
	seen := make(map[string]struct{}, len(specSlugs))
	for position, specSlug := range specSlugs {
		specSlug = strings.TrimSpace(specSlug)
		if specSlug == "" {
			return DeliveryQueue{}, fmt.Errorf("create Delivery Queue: Spec slug at position %d is required", position)
		}
		if _, exists := seen[specSlug]; exists {
			return DeliveryQueue{}, fmt.Errorf("create Delivery Queue: Spec slug %q appears more than once", specSlug)
		}
		seen[specSlug] = struct{}{}
		items[position] = DeliveryQueueItem{
			SpecSlug:         specSlug,
			Position:         position,
			Stage:            DeliveryStageQueued,
			CandidateCommits: []string{},
		}
	}

	queue := DeliveryQueue{
		GitRoot: gitRoot,
		Limits:  limits,
		Items:   items,
	}
	err := store.withWriteTx(ctx, "Delivery Queue creation", func(tx *sql.Tx) error {
		var ownerPID sql.NullInt64
		err := tx.QueryRowContext(ctx, `
SELECT owner_pid
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(&ownerPID)
		switch {
		case err == nil:
			if ownerPID.Valid && ownerPID.Int64 > 0 {
				return fmt.Errorf("create Delivery Queue: repository %q already has owner PID %d", gitRoot, ownerPID.Int64)
			}
			var unfinished int
			if err := tx.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM delivery_queue_items
WHERE git_root = ? AND stage NOT IN (?, ?)`,
				gitRoot,
				DeliveryStageMerged,
				DeliveryStageParked,
			).Scan(&unfinished); err != nil {
				return fmt.Errorf("inspect existing Delivery Queue: %w", err)
			}
			if unfinished > 0 {
				return fmt.Errorf("create Delivery Queue: repository %q already has %d unfinished item(s)", gitRoot, unfinished)
			}
			if _, err := tx.ExecContext(ctx, `DELETE FROM delivery_queues WHERE git_root = ?`, gitRoot); err != nil {
				return fmt.Errorf("replace terminal Delivery Queue: %w", err)
			}
		case errors.Is(err, sql.ErrNoRows):
			// No queue exists yet.
		case err != nil:
			return fmt.Errorf("inspect existing Delivery Queue: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO delivery_queues (git_root, deadline_unix, max_retries, max_tokens)
VALUES (?, ?, ?, ?)`, gitRoot, deadlineUnix, limits.MaxRetries, limits.MaxTokens); err != nil {
			return fmt.Errorf("insert Delivery Queue for repository %q: %w", gitRoot, err)
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `
INSERT INTO delivery_queue_items (
	git_root, spec_slug, position, stage, blocker, branch, run_id,
	candidate_commits, pull_request_number, merge_commit
) VALUES (?, ?, ?, ?, '', '', '', '[]', '', '')`,
				gitRoot,
				item.SpecSlug,
				item.Position,
				item.Stage,
			); err != nil {
				return fmt.Errorf("insert Delivery Queue item %q: %w", item.SpecSlug, err)
			}
		}
		return nil
	})
	if err != nil {
		return DeliveryQueue{}, err
	}
	return queue, nil
}

// DeliveryQueue returns one repository's persisted queue, ordered by the
// positions recorded when it was created.
func (store *Store) DeliveryQueue(ctx context.Context, gitRoot string) (DeliveryQueue, bool, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return DeliveryQueue{}, false, errors.New("read Delivery Queue: Git root is required")
	}

	var queue DeliveryQueue
	var ownerPID sql.NullInt64
	var deadlineUnix int64
	err := store.db.QueryRowContext(ctx, `
SELECT git_root, owner_pid, owner_identity, deadline_unix, max_retries, max_tokens
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(
		&queue.GitRoot,
		&ownerPID,
		&queue.OwnerIdentity,
		&deadlineUnix,
		&queue.Limits.MaxRetries,
		&queue.Limits.MaxTokens,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return DeliveryQueue{}, false, nil
	}
	if err != nil {
		return DeliveryQueue{}, false, fmt.Errorf("read Delivery Queue for repository %q: %w", gitRoot, err)
	}
	if ownerPID.Valid {
		queue.OwnerPID = int(ownerPID.Int64)
	}
	if deadlineUnix != 0 {
		queue.Limits.Deadline = time.Unix(deadlineUnix, 0).UTC()
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT spec_slug, position, stage, blocker, branch, worktree, worktree_provisioned, run_id, candidate_commits,
	   pull_request_number, merge_commit, retry_count, warning
FROM delivery_queue_items
WHERE git_root = ?
ORDER BY position`, gitRoot)
	if err != nil {
		return DeliveryQueue{}, false, fmt.Errorf("list Delivery Queue items: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	queue.Items = []DeliveryQueueItem{}
	for rows.Next() {
		item, err := scanDeliveryQueueItem(rows)
		if err != nil {
			return DeliveryQueue{}, false, fmt.Errorf("scan Delivery Queue item: %w", err)
		}
		queue.Items = append(queue.Items, item)
	}
	if err := rows.Err(); err != nil {
		return DeliveryQueue{}, false, fmt.Errorf("iterate Delivery Queue items: %w", err)
	}
	return queue, true, nil
}

// ClaimDeliveryQueueOwner records the process that exclusively advances one
// repository's Delivery Queue. Repeating the same claim is idempotent; a
// different owner is refused until stop releases the recorded process.
func (store *Store) ClaimDeliveryQueueOwner(ctx context.Context, gitRoot string, pid int, identity string) error {
	gitRoot = strings.TrimSpace(gitRoot)
	identity = strings.TrimSpace(identity)
	if gitRoot == "" {
		return errors.New("claim Delivery Queue owner: Git root is required")
	}
	if pid < 1 {
		return errors.New("claim Delivery Queue owner: PID is required")
	}
	if identity == "" {
		return errors.New("claim Delivery Queue owner: process identity is required")
	}

	return store.withWriteTx(ctx, "Delivery Queue owner claim", func(tx *sql.Tx) error {
		var storedPID sql.NullInt64
		var storedIdentity string
		if err := tx.QueryRowContext(ctx, `
SELECT owner_pid, owner_identity
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(&storedPID, &storedIdentity); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("claim Delivery Queue owner: queue for repository %q does not exist", gitRoot)
			}
			return fmt.Errorf("read Delivery Queue owner: %w", err)
		}
		if storedPID.Valid && storedPID.Int64 > 0 {
			if int(storedPID.Int64) == pid && storedIdentity == identity {
				return nil
			}
			return fmt.Errorf(
				"claim Delivery Queue owner: repository %q already has owner PID %d",
				gitRoot,
				storedPID.Int64,
			)
		}
		if _, err := tx.ExecContext(ctx, `
UPDATE delivery_queues
SET owner_pid = ?, owner_identity = ?
WHERE git_root = ?`, pid, identity, gitRoot); err != nil {
			return fmt.Errorf("record Delivery Queue owner: %w", err)
		}
		return nil
	})
}

// ReleaseDeliveryQueueOwner clears only the owner identity the caller proved.
// A stale process can never clear a newer owner's claim.
func (store *Store) ReleaseDeliveryQueueOwner(ctx context.Context, gitRoot string, pid int, identity string) (bool, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	identity = strings.TrimSpace(identity)
	if gitRoot == "" {
		return false, errors.New("release Delivery Queue owner: Git root is required")
	}
	if pid < 1 {
		return false, errors.New("release Delivery Queue owner: PID is required")
	}
	if identity == "" {
		return false, errors.New("release Delivery Queue owner: process identity is required")
	}

	released := false
	err := store.withWriteTx(ctx, "Delivery Queue owner release", func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
UPDATE delivery_queues
SET owner_pid = NULL, owner_identity = ''
WHERE git_root = ? AND owner_pid = ? AND owner_identity = ?`, gitRoot, pid, identity)
		if err != nil {
			return fmt.Errorf("clear Delivery Queue owner: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read Delivery Queue owner release result: %w", err)
		}
		released = affected > 0
		return nil
	})
	return released, err
}

// ReleaseIdleDeliveryQueueOwner clears the caller's owner claim only when no
// item can advance. The owner check, idle check, and release share one write
// transaction so a concurrent retry cannot be stranded after the owner exits.
func (store *Store) ReleaseIdleDeliveryQueueOwner(
	ctx context.Context,
	gitRoot string,
	pid int,
	identity string,
) (bool, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	identity = strings.TrimSpace(identity)
	if gitRoot == "" {
		return false, errors.New("release idle Delivery Queue owner: Git root is required")
	}
	if pid < 1 {
		return false, errors.New("release idle Delivery Queue owner: PID is required")
	}
	if identity == "" {
		return false, errors.New("release idle Delivery Queue owner: process identity is required")
	}

	released := false
	err := store.withWriteTx(ctx, "idle Delivery Queue owner release", func(tx *sql.Tx) error {
		var storedPID sql.NullInt64
		var storedIdentity string
		if err := tx.QueryRowContext(ctx, `
SELECT owner_pid, owner_identity
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(&storedPID, &storedIdentity); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("release idle Delivery Queue owner: queue for repository %q does not exist", gitRoot)
			}
			return fmt.Errorf("read Delivery Queue owner: %w", err)
		}
		recordedPID := 0
		if storedPID.Valid {
			recordedPID = int(storedPID.Int64)
		}
		if recordedPID != pid || storedIdentity != identity {
			return fmt.Errorf(
				"release idle Delivery Queue owner: recorded owner is PID %d with identity %q",
				recordedPID,
				storedIdentity,
			)
		}

		var advanceable int
		if err := tx.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM delivery_queue_items
WHERE git_root = ? AND stage NOT IN (?, ?)`,
			gitRoot,
			DeliveryStageMerged,
			DeliveryStageParked,
		).Scan(&advanceable); err != nil {
			return fmt.Errorf("inspect Delivery Queue for advanceable items: %w", err)
		}
		if advanceable > 0 {
			return nil
		}

		result, err := tx.ExecContext(ctx, `
UPDATE delivery_queues
SET owner_pid = NULL, owner_identity = ''
WHERE git_root = ? AND owner_pid = ? AND owner_identity = ?`, gitRoot, pid, identity)
		if err != nil {
			return fmt.Errorf("clear idle Delivery Queue owner: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read idle Delivery Queue owner release result: %w", err)
		}
		if affected != 1 {
			return errors.New("clear idle Delivery Queue owner: owner changed during transaction")
		}
		released = true
		return nil
	})
	return released, err
}

// RecordDeliveryQueueItemWorktree records the first item branch and worktree
// together and returns the durable pair on retries.
func (store *Store) RecordDeliveryQueueItemWorktree(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	branch string,
	worktree string,
) (string, string, bool, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	specSlug = strings.TrimSpace(specSlug)
	branch = strings.TrimSpace(branch)
	worktree = strings.TrimSpace(worktree)
	if gitRoot == "" {
		return "", "", false, errors.New("record Delivery Queue item worktree: Git root is required")
	}
	if specSlug == "" {
		return "", "", false, errors.New("record Delivery Queue item worktree: Spec slug is required")
	}
	if branch == "" {
		return "", "", false, errors.New("record Delivery Queue item worktree: branch is required")
	}
	if worktree == "" {
		return "", "", false, errors.New("record Delivery Queue item worktree: worktree is required")
	}

	recordedBranch := ""
	recordedWorktree := ""
	recordedProvisioned := false
	err := store.withWriteTx(ctx, fmt.Sprintf("Delivery Queue item %q worktree recording", specSlug), func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
UPDATE delivery_queue_items
SET branch = ?, worktree = ?
WHERE git_root = ? AND spec_slug = ? AND branch = '' AND worktree = ''`,
			branch, worktree, gitRoot, specSlug); err != nil {
			return fmt.Errorf("record Delivery Queue item %q worktree: %w", specSlug, err)
		}
		if err := tx.QueryRowContext(ctx, `
SELECT branch, worktree, worktree_provisioned
FROM delivery_queue_items
WHERE git_root = ? AND spec_slug = ?`, gitRoot, specSlug).Scan(
			&recordedBranch,
			&recordedWorktree,
			&recordedProvisioned,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("record Delivery Queue item %q worktree: item does not exist", specSlug)
			}
			return fmt.Errorf("read Delivery Queue item %q worktree: %w", specSlug, err)
		}
		if strings.TrimSpace(recordedBranch) == "" {
			return fmt.Errorf("record Delivery Queue item %q worktree: recorded branch is empty", specSlug)
		}
		if strings.TrimSpace(recordedWorktree) == "" {
			return fmt.Errorf("record Delivery Queue item %q worktree: recorded worktree is empty", specSlug)
		}
		return nil
	})
	if err != nil {
		return "", "", false, err
	}
	return recordedBranch, recordedWorktree, recordedProvisioned, nil
}

// SetDeliveryQueueItemWorktreeProvisioned records whether the current item
// worktree finished its copy and bootstrap steps.
func (store *Store) SetDeliveryQueueItemWorktreeProvisioned(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	provisioned bool,
) error {
	gitRoot = strings.TrimSpace(gitRoot)
	specSlug = strings.TrimSpace(specSlug)
	if gitRoot == "" {
		return errors.New("record Delivery Queue item provisioning: Git root is required")
	}
	if specSlug == "" {
		return errors.New("record Delivery Queue item provisioning: Spec slug is required")
	}

	return store.withWriteTx(ctx, fmt.Sprintf("Delivery Queue item %q provisioning update", specSlug), func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
UPDATE delivery_queue_items
SET worktree_provisioned = ?
WHERE git_root = ? AND spec_slug = ?`, provisioned, gitRoot, specSlug)
		if err != nil {
			return fmt.Errorf("record Delivery Queue item %q provisioning: %w", specSlug, err)
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read Delivery Queue item %q provisioning update: %w", specSlug, err)
		}
		if updated != 1 {
			return fmt.Errorf("record Delivery Queue item %q provisioning: item does not exist", specSlug)
		}
		return nil
	})
}

// UpdateDeliveryQueueItem persists one item's current delivery state without
// changing its original queue position.
func (store *Store) UpdateDeliveryQueueItem(ctx context.Context, gitRoot string, item DeliveryQueueItem) error {
	gitRoot = strings.TrimSpace(gitRoot)
	item.SpecSlug = strings.TrimSpace(item.SpecSlug)
	item.Blocker = strings.TrimSpace(item.Blocker)
	item.Branch = strings.TrimSpace(item.Branch)
	item.RunID = strings.TrimSpace(item.RunID)
	item.PullRequestNumber = strings.TrimSpace(item.PullRequestNumber)
	item.MergeCommit = strings.TrimSpace(item.MergeCommit)
	item.Warning = strings.TrimSpace(item.Warning)
	if gitRoot == "" {
		return errors.New("update Delivery Queue item: Git root is required")
	}
	if item.SpecSlug == "" {
		return errors.New("update Delivery Queue item: Spec slug is required")
	}
	if !validDeliveryStage(item.Stage) {
		return fmt.Errorf("update Delivery Queue item %q: stage %q is invalid", item.SpecSlug, item.Stage)
	}
	commits := item.CandidateCommits
	if commits == nil {
		commits = []string{}
	}
	encodedCommits, err := json.Marshal(commits)
	if err != nil {
		return fmt.Errorf("encode candidate commits for Delivery Queue item %q: %w", item.SpecSlug, err)
	}

	return store.withWriteTx(ctx, fmt.Sprintf("Delivery Queue item %q update", item.SpecSlug), func(tx *sql.Tx) error {
		var storedPosition int
		if err := tx.QueryRowContext(ctx, `
SELECT position
FROM delivery_queue_items
WHERE git_root = ? AND spec_slug = ?`, gitRoot, item.SpecSlug).Scan(&storedPosition); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("update Delivery Queue item %q: item does not exist", item.SpecSlug)
			}
			return fmt.Errorf("read Delivery Queue item %q before update: %w", item.SpecSlug, err)
		}
		if item.Position != storedPosition {
			return fmt.Errorf(
				"update Delivery Queue item %q: position %d does not match stored position %d",
				item.SpecSlug,
				item.Position,
				storedPosition,
			)
		}

		if _, err := tx.ExecContext(ctx, `
UPDATE delivery_queue_items
SET stage = ?, blocker = ?, branch = ?, run_id = ?, candidate_commits = ?,
	pull_request_number = ?, merge_commit = ?, warning = ?
WHERE git_root = ? AND spec_slug = ?`,
			item.Stage,
			item.Blocker,
			item.Branch,
			item.RunID,
			string(encodedCommits),
			item.PullRequestNumber,
			item.MergeCommit,
			item.Warning,
			gitRoot,
			item.SpecSlug,
		); err != nil {
			return fmt.Errorf("update Delivery Queue item %q: %w", item.SpecSlug, err)
		}
		return linkDeliveryQueueRun(ctx, tx, gitRoot, item)
	})
}

// RetryDeliveryQueueItem moves one parked item back to an advanceable stage
// only when its persisted blocker still matches the operator's observation.
func (store *Store) RetryDeliveryQueueItem(
	ctx context.Context,
	gitRoot string,
	item DeliveryQueueItem,
	parkedBlocker string,
) (int, string, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	item.SpecSlug = strings.TrimSpace(item.SpecSlug)
	item.RunID = strings.TrimSpace(item.RunID)
	parkedBlocker = strings.TrimSpace(parkedBlocker)
	if gitRoot == "" {
		return 0, "", errors.New("retry Delivery Queue item: Git root is required")
	}
	if item.SpecSlug == "" {
		return 0, "", errors.New("retry Delivery Queue item: Spec slug is required")
	}
	if item.Position < 0 {
		return 0, "", fmt.Errorf("retry Delivery Queue item %q: position must not be negative", item.SpecSlug)
	}
	if !retryDeliveryStage(item.Stage) {
		return 0, "", fmt.Errorf(
			"retry Delivery Queue item %q: target stage %q is not retryable",
			item.SpecSlug,
			item.Stage,
		)
	}
	commits := item.CandidateCommits
	if commits == nil {
		commits = []string{}
	}
	encodedCommits, err := json.Marshal(commits)
	if err != nil {
		return 0, "", fmt.Errorf("encode candidate commits for Delivery Queue item %q: %w", item.SpecSlug, err)
	}

	ownerPID := 0
	ownerIdentity := ""
	err = store.withWriteTx(ctx, fmt.Sprintf("Delivery Queue item %q retry", item.SpecSlug), func(tx *sql.Tx) error {
		var storedSpecSlug string
		var storedStage DeliveryStage
		var storedBlocker string
		var retryCount int
		var maxRetries int
		if err := tx.QueryRowContext(ctx, `
SELECT item.spec_slug, item.stage, item.blocker, item.retry_count, queue.max_retries
FROM delivery_queue_items item
JOIN delivery_queues queue ON queue.git_root = item.git_root
WHERE item.git_root = ? AND item.position = ?`, gitRoot, item.Position).Scan(
			&storedSpecSlug,
			&storedStage,
			&storedBlocker,
			&retryCount,
			&maxRetries,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf(
					"retry Delivery Queue item %q: item at position %d does not exist",
					item.SpecSlug,
					item.Position,
				)
			}
			return fmt.Errorf("read Delivery Queue item %q before retry: %w", item.SpecSlug, err)
		}
		if storedSpecSlug != item.SpecSlug {
			return fmt.Errorf(
				"retry Delivery Queue item %q: stored item at position %d is %q",
				item.SpecSlug,
				item.Position,
				storedSpecSlug,
			)
		}
		if storedStage != DeliveryStageParked || storedBlocker != parkedBlocker {
			return fmt.Errorf(
				"retry Delivery Queue item %q: stored item has stage %q and blocker %q; want stage %q and blocker %q",
				item.SpecSlug,
				storedStage,
				storedBlocker,
				DeliveryStageParked,
				parkedBlocker,
			)
		}
		if maxRetries != 0 && retryCount >= maxRetries {
			return fmt.Errorf(
				"retry Delivery Queue item %q: retry count %d reached queue limit %d: %w",
				item.SpecSlug,
				retryCount,
				maxRetries,
				ErrDeliveryRetryLimit,
			)
		}

		if _, err := tx.ExecContext(ctx, `
UPDATE delivery_queue_items
SET stage = ?, blocker = '', run_id = ?, candidate_commits = ?, retry_count = retry_count + 1
WHERE git_root = ? AND position = ?`,
			item.Stage,
			item.RunID,
			string(encodedCommits),
			gitRoot,
			item.Position,
		); err != nil {
			return fmt.Errorf("retry Delivery Queue item %q: %w", item.SpecSlug, err)
		}

		if err := linkDeliveryQueueRun(ctx, tx, gitRoot, item); err != nil {
			return err
		}
		var storedOwnerPID sql.NullInt64
		if err := tx.QueryRowContext(ctx, `
SELECT owner_pid, owner_identity
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(&storedOwnerPID, &ownerIdentity); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("retry Delivery Queue item %q: queue does not exist", item.SpecSlug)
			}
			return fmt.Errorf("read Delivery Queue owner after retrying item %q: %w", item.SpecSlug, err)
		}
		if storedOwnerPID.Valid {
			ownerPID = int(storedOwnerPID.Int64)
		}
		return nil
	})
	if err != nil {
		return 0, "", err
	}
	return ownerPID, ownerIdentity, nil
}

func (store *Store) RecordDeliveryActionIntent(
	ctx context.Context,
	gitRoot string,
	specSlug string,
	action DeliveryAction,
) (DeliveryActionIntent, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	specSlug = strings.TrimSpace(specSlug)
	if gitRoot == "" {
		return DeliveryActionIntent{}, errors.New("record Delivery Action intent: Git root is required")
	}
	if specSlug == "" {
		return DeliveryActionIntent{}, errors.New("record Delivery Action intent: Spec slug is required")
	}
	if !validDeliveryAction(action) {
		return DeliveryActionIntent{}, fmt.Errorf("record Delivery Action intent: action %q is invalid", action)
	}

	intent := DeliveryActionIntent{
		GitRoot:   gitRoot,
		SpecSlug:  specSlug,
		Action:    action,
		CreatedAt: store.now(),
	}
	err := store.withWriteTx(ctx, "Delivery Action intent recording", func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, `
INSERT INTO delivery_action_intents (git_root, spec_slug, action, created_at)
VALUES (?, ?, ?, ?)`, gitRoot, specSlug, action, formatTime(intent.CreatedAt))
		if err != nil {
			return fmt.Errorf("insert Delivery Action intent for Spec %q: %w", specSlug, err)
		}
		intent.ID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("read Delivery Action intent identity: %w", err)
		}
		return nil
	})
	if err != nil {
		return DeliveryActionIntent{}, err
	}
	return intent, nil
}

func (store *Store) RecordDeliveryActionReceipt(
	ctx context.Context,
	intentID int64,
	result string,
) (DeliveryActionReceipt, error) {
	if intentID < 1 {
		return DeliveryActionReceipt{}, errors.New("record Delivery Action receipt: intent ID is required")
	}
	result = strings.TrimSpace(result)
	if result == "" {
		return DeliveryActionReceipt{}, errors.New("record Delivery Action receipt: result is required")
	}
	receipt := DeliveryActionReceipt{
		IntentID:  intentID,
		Result:    result,
		CreatedAt: store.now(),
	}
	err := store.withWriteTx(ctx, "Delivery Action receipt recording", func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
INSERT INTO delivery_action_receipts (intent_id, result, created_at)
VALUES (?, ?, ?)`, intentID, result, formatTime(receipt.CreatedAt)); err != nil {
			return fmt.Errorf("insert receipt for Delivery Action intent %d: %w", intentID, err)
		}
		return nil
	})
	if err != nil {
		return DeliveryActionReceipt{}, err
	}
	return receipt, nil
}

func (store *Store) UnmatchedDeliveryActionIntents(ctx context.Context, gitRoot string) ([]DeliveryActionIntent, error) {
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return nil, errors.New("list unmatched Delivery Action intents: Git root is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT intent.id, intent.git_root, intent.spec_slug, intent.action, intent.created_at
FROM delivery_action_intents intent
LEFT JOIN delivery_action_receipts receipt ON receipt.intent_id = intent.id
WHERE intent.git_root = ? AND receipt.intent_id IS NULL
ORDER BY intent.id`, gitRoot)
	if err != nil {
		return nil, fmt.Errorf("list unmatched Delivery Action intents: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	intents := []DeliveryActionIntent{}
	for rows.Next() {
		intent, err := scanDeliveryActionIntent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan unmatched Delivery Action intent: %w", err)
		}
		intents = append(intents, intent)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate unmatched Delivery Action intents: %w", err)
	}
	return intents, nil
}

func (store *Store) DeliveryActionReceipt(
	ctx context.Context,
	intentID int64,
) (DeliveryActionReceipt, bool, error) {
	if intentID < 1 {
		return DeliveryActionReceipt{}, false, errors.New("read Delivery Action receipt: intent ID is required")
	}
	var receipt DeliveryActionReceipt
	var createdAt string
	err := store.db.QueryRowContext(ctx, `
SELECT intent_id, result, created_at
FROM delivery_action_receipts
WHERE intent_id = ?`, intentID).Scan(&receipt.IntentID, &receipt.Result, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return DeliveryActionReceipt{}, false, nil
	}
	if err != nil {
		return DeliveryActionReceipt{}, false, fmt.Errorf("read Delivery Action receipt: %w", err)
	}
	receipt.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return DeliveryActionReceipt{}, false, fmt.Errorf("read Delivery Action receipt time: %w", err)
	}
	return receipt, true, nil
}

type deliveryQueueItemScanner interface {
	Scan(dest ...any) error
}

func scanDeliveryQueueItem(row deliveryQueueItemScanner) (DeliveryQueueItem, error) {
	var item DeliveryQueueItem
	var candidateCommits string
	if err := row.Scan(
		&item.SpecSlug,
		&item.Position,
		&item.Stage,
		&item.Blocker,
		&item.Branch,
		&item.Worktree,
		&item.WorktreeProvisioned,
		&item.RunID,
		&candidateCommits,
		&item.PullRequestNumber,
		&item.MergeCommit,
		&item.RetryCount,
		&item.Warning,
	); err != nil {
		return DeliveryQueueItem{}, err
	}
	if err := json.Unmarshal([]byte(candidateCommits), &item.CandidateCommits); err != nil {
		return DeliveryQueueItem{}, fmt.Errorf("decode candidate commits for Spec %q: %w", item.SpecSlug, err)
	}
	if item.CandidateCommits == nil {
		item.CandidateCommits = []string{}
	}
	return item, nil
}

func scanDeliveryActionIntent(row deliveryQueueItemScanner) (DeliveryActionIntent, error) {
	var intent DeliveryActionIntent
	var createdAt string
	if err := row.Scan(&intent.ID, &intent.GitRoot, &intent.SpecSlug, &intent.Action, &createdAt); err != nil {
		return DeliveryActionIntent{}, err
	}
	parsed, err := parseTime(createdAt)
	if err != nil {
		return DeliveryActionIntent{}, err
	}
	intent.CreatedAt = parsed
	return intent, nil
}

func validDeliveryStage(stage DeliveryStage) bool {
	switch stage {
	case DeliveryStageQueued,
		DeliveryStageRunning,
		DeliveryStageReviewing,
		DeliveryStageArchiving,
		DeliveryStageGating,
		DeliveryStagePublishing,
		DeliveryStageChecking,
		DeliveryStageMerging,
		DeliveryStageMerged,
		DeliveryStageParked:
		return true
	default:
		return false
	}
}

func retryDeliveryStage(stage DeliveryStage) bool {
	switch stage {
	case DeliveryStageQueued, DeliveryStageRunning, DeliveryStageReviewing, DeliveryStageGating, DeliveryStageChecking:
		return true
	default:
		return false
	}
}

func validDeliveryAction(action DeliveryAction) bool {
	switch action {
	case DeliveryActionPush, DeliveryActionCreatePullRequest, DeliveryActionMerge:
		return true
	default:
		return false
	}
}
