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

type DeliveryQueue struct {
	GitRoot       string
	OwnerPID      int
	OwnerIdentity string
	Items         []DeliveryQueueItem
}

type DeliveryQueueItem struct {
	SpecSlug          string
	Position          int
	Stage             DeliveryStage
	Blocker           string
	Branch            string
	RunID             string
	CandidateCommits  []string
	PullRequestNumber string
	MergeCommit       string
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
	gitRoot = strings.TrimSpace(gitRoot)
	if gitRoot == "" {
		return DeliveryQueue{}, errors.New("create Delivery Queue: Git root is required")
	}
	if len(specSlugs) == 0 {
		return DeliveryQueue{}, errors.New("create Delivery Queue: at least one Spec slug is required")
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
INSERT INTO delivery_queues (git_root)
VALUES (?)`, gitRoot); err != nil {
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
	err := store.db.QueryRowContext(ctx, `
SELECT git_root, owner_pid, owner_identity
FROM delivery_queues
WHERE git_root = ?`, gitRoot).Scan(&queue.GitRoot, &ownerPID, &queue.OwnerIdentity)
	if errors.Is(err, sql.ErrNoRows) {
		return DeliveryQueue{}, false, nil
	}
	if err != nil {
		return DeliveryQueue{}, false, fmt.Errorf("read Delivery Queue for repository %q: %w", gitRoot, err)
	}
	if ownerPID.Valid {
		queue.OwnerPID = int(ownerPID.Int64)
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT spec_slug, position, stage, blocker, branch, run_id, candidate_commits,
       pull_request_number, merge_commit
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
    pull_request_number = ?, merge_commit = ?
WHERE git_root = ? AND spec_slug = ?`,
			item.Stage,
			item.Blocker,
			item.Branch,
			item.RunID,
			string(encodedCommits),
			item.PullRequestNumber,
			item.MergeCommit,
			gitRoot,
			item.SpecSlug,
		); err != nil {
			return fmt.Errorf("update Delivery Queue item %q: %w", item.SpecSlug, err)
		}
		return nil
	})
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
		&item.RunID,
		&candidateCommits,
		&item.PullRequestNumber,
		&item.MergeCommit,
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

func validDeliveryAction(action DeliveryAction) bool {
	switch action {
	case DeliveryActionPush, DeliveryActionCreatePullRequest, DeliveryActionMerge:
		return true
	default:
		return false
	}
}
