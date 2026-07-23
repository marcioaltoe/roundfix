package baseline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	transactionJournalSchema  = "roundfix/baseline-transaction/v1"
	transactionStateDirName   = "current"
	transactionJournalName    = "journal.json"
	transactionJournalNext    = "journal.next"
	transactionStagingDirName = "staged"
	maxTransactionJournalSize = 32 * 1024 * 1024
)

var ErrTransactionLocked = errors.New("another Baseline transaction owns this worktree")

// VerificationEvidence records the exact postimages verified by a completed
// file transaction.
type VerificationEvidence struct {
	VerifiedPostimages []Postimage
}

// Transaction is the recoverable mutation boundary for one confirmed
// Baseline Plan.
type Transaction interface {
	Apply(context.Context) (VerificationEvidence, error)
	Rollback(context.Context) error
	Close() error
}

// IncompleteRollbackError means repository mutation may remain and the
// Git-private journal must be recovered before another transaction can begin.
type IncompleteRollbackError struct {
	Cause    error
	Rollback error
}

func (err *IncompleteRollbackError) Error() string {
	return fmt.Sprintf("Baseline transaction rollback is incomplete: cause: %v; rollback: %v", err.Cause, err.Rollback)
}

func (err *IncompleteRollbackError) Unwrap() []error {
	return []error{err.Cause, err.Rollback}
}

type transactionPhase string

const (
	transactionPhaseJournaled          transactionPhase = "journaled"
	transactionPhaseStaging            transactionPhase = "staging"
	transactionPhaseStaged             transactionPhase = "staged"
	transactionPhasePreimagesValidated transactionPhase = "preimages_validated"
	transactionPhaseReplacing          transactionPhase = "replacing"
	transactionPhaseReplaced           transactionPhase = "replaced"
	transactionPhaseVerifying          transactionPhase = "verifying"
	transactionPhaseCommitting         transactionPhase = "committing"
	transactionPhaseRollingBack        transactionPhase = "rolling_back"
)

type transactionFaultPoint struct {
	Phase transactionPhase
	Path  string
}

func (point transactionFaultPoint) String() string {
	if point.Path == "" {
		return string(point.Phase)
	}
	return string(point.Phase) + ":" + point.Path
}

type transactionJournalPhase string

const (
	transactionJournalPrepared    transactionJournalPhase = "prepared"
	transactionJournalStaged      transactionJournalPhase = "staged"
	transactionJournalReplacing   transactionJournalPhase = "replacing"
	transactionJournalVerifying   transactionJournalPhase = "verifying"
	transactionJournalCommitted   transactionJournalPhase = "committed"
	transactionJournalRollingBack transactionJournalPhase = "rolling_back"
)

type transactionFileState struct {
	Exists          bool         `json:"exists"`
	Kind            PreimageKind `json:"kind"`
	Mode            uint32       `json:"mode"`
	LinkTarget      string       `json:"linkTarget,omitempty"`
	Bytes           int64        `json:"bytes"`
	Content         []byte       `json:"content,omitempty"`
	ContentIdentity string       `json:"contentIdentity,omitempty"`
}

type transactionJournalPreimage struct {
	Path  string               `json:"path"`
	State transactionFileState `json:"state"`
}

type transactionJournalEntry struct {
	Path   string               `json:"path"`
	Before transactionFileState `json:"before"`
	After  Postimage            `json:"after"`
}

type transactionJournal struct {
	SchemaVersion      string                       `json:"schemaVersion"`
	PlanDigest         string                       `json:"planDigest"`
	Phase              transactionJournalPhase      `json:"phase"`
	Preimages          []transactionJournalPreimage `json:"preimages"`
	Entries            []transactionJournalEntry    `json:"entries"`
	MutationOrder      []int                        `json:"mutationOrder"`
	CreatedDirectories []string                     `json:"createdDirectories"`
	TemporaryPaths     []string                     `json:"temporaryPaths"`
}

type fileTransaction struct {
	root           string
	anchored       *os.Root
	privateDir     string
	stateDir       string
	lock           *os.File
	journal        transactionJournal
	phaseHook      func(transactionFaultPoint) error
	applyAttempted bool
	applied        bool
	closed         bool
}

// BeginTransaction locks one Git worktree, recovers any interrupted mutation,
// validates the complete bounded preimage, and writes an exact recovery
// journal before returning.
func BeginTransaction(
	ctx context.Context,
	repository string,
	document PlanDocument,
) (Transaction, error) {
	return beginTransaction(ctx, repository, document, nil)
}

func beginTransaction(
	ctx context.Context,
	repository string,
	document PlanDocument,
	phaseHook func(transactionFaultPoint) error,
) (Transaction, error) {
	if ctx == nil {
		return nil, errors.New("begin Baseline transaction: context is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, _, err := inspectRepositoryIdentity(ctx, repository, nil)
	if err != nil {
		return nil, fmt.Errorf("begin Baseline transaction: %w", err)
	}
	privateDir, err := transactionPrivateDirectory(ctx, root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(privateDir, 0o700); err != nil {
		return nil, fmt.Errorf("create Git-private Baseline transaction directory: %w", err)
	}
	if err := syncDirectory(filepath.Dir(privateDir)); err != nil {
		return nil, fmt.Errorf("sync Git-private Baseline transaction parent: %w", err)
	}

	lock, err := os.OpenFile(filepath.Join(privateDir, "transaction.lock"), os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open Baseline transaction lock: %w", err)
	}
	if err := lockTransactionFile(lock); err != nil {
		_ = lock.Close()
		return nil, err
	}
	release := func() {
		_ = unlockTransactionFile(lock)
		_ = lock.Close()
	}

	anchored, err := os.OpenRoot(root)
	if err != nil {
		release()
		return nil, fmt.Errorf("open Baseline transaction repository root: %w", err)
	}
	stateDir := filepath.Join(privateDir, transactionStateDirName)
	if err := recoverInterruptedTransaction(context.WithoutCancel(ctx), root, anchored, privateDir, stateDir); err != nil {
		_ = anchored.Close()
		release()
		return nil, err
	}
	if err := ValidatePlanRepository(ctx, root, document); err != nil {
		_ = anchored.Close()
		release()
		return nil, fmt.Errorf("validate Baseline transaction plan: %w", err)
	}
	if err := validateTransactionPostimages(anchored, document.Postimages); err != nil {
		_ = anchored.Close()
		release()
		return nil, err
	}
	preimages, entries, err := captureTransactionEntries(anchored, document)
	if err != nil {
		_ = anchored.Close()
		release()
		return nil, err
	}
	if err := os.Mkdir(stateDir, 0o700); err != nil {
		_ = anchored.Close()
		release()
		return nil, fmt.Errorf("create Baseline transaction recovery state: %w", err)
	}
	transaction := &fileTransaction{
		root:       root,
		anchored:   anchored,
		privateDir: privateDir,
		stateDir:   stateDir,
		lock:       lock,
		phaseHook:  phaseHook,
		journal: transactionJournal{
			SchemaVersion:      transactionJournalSchema,
			PlanDigest:         document.PlanDigest,
			Phase:              transactionJournalPrepared,
			Preimages:          preimages,
			Entries:            entries,
			MutationOrder:      []int{},
			CreatedDirectories: []string{},
			TemporaryPaths:     []string{},
		},
	}
	if err := transaction.writeJournal(); err != nil {
		_ = os.RemoveAll(stateDir)
		_ = anchored.Close()
		release()
		return nil, err
	}
	if err := transaction.runPhaseHook(transactionFaultPoint{Phase: transactionPhaseJournaled}); err != nil {
		_ = transaction.rollback(context.WithoutCancel(ctx))
		_ = anchored.Close()
		release()
		return nil, err
	}
	return transaction, nil
}

func (transaction *fileTransaction) Apply(ctx context.Context) (VerificationEvidence, error) {
	if ctx == nil {
		return VerificationEvidence{}, errors.New("apply Baseline transaction: context is required")
	}
	if transaction.closed {
		return VerificationEvidence{}, errors.New("apply Baseline transaction: transaction is closed")
	}
	if transaction.applyAttempted {
		return VerificationEvidence{}, errors.New("apply Baseline transaction: transaction was already applied")
	}
	transaction.applyAttempted = true
	if err := ctx.Err(); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	if err := transaction.stagePostimages(ctx); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	if err := transaction.revalidatePreimages(ctx); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	for index := range transaction.journal.Entries {
		if err := transaction.replacePostimage(ctx, index); err != nil {
			return VerificationEvidence{}, transaction.failApply(ctx, err)
		}
	}
	transaction.journal.Phase = transactionJournalVerifying
	if err := transaction.writeJournal(); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	for _, entry := range transaction.journal.Entries {
		if err := ctx.Err(); err != nil {
			return VerificationEvidence{}, transaction.failApply(ctx, err)
		}
		point := transactionFaultPoint{Phase: transactionPhaseVerifying, Path: entry.Path}
		if err := transaction.runPhaseHook(point); err != nil {
			return VerificationEvidence{}, transaction.failApply(ctx, err)
		}
		if err := verifyTransactionPostimage(transaction.anchored, entry.After); err != nil {
			return VerificationEvidence{}, transaction.failApply(ctx, err)
		}
	}
	if err := transaction.runPhaseHook(transactionFaultPoint{Phase: transactionPhaseCommitting}); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	transaction.journal.Phase = transactionJournalCommitted
	if err := transaction.writeJournal(); err != nil {
		return VerificationEvidence{}, transaction.failApply(ctx, err)
	}
	if err := transaction.cleanupState(); err != nil {
		return VerificationEvidence{}, fmt.Errorf("complete Baseline transaction recovery cleanup: %w", err)
	}
	transaction.applied = true
	verified := make([]Postimage, len(transaction.journal.Entries))
	for index, entry := range transaction.journal.Entries {
		verified[index] = clonePostimage(entry.After)
	}
	return VerificationEvidence{VerifiedPostimages: verified}, nil
}

func (transaction *fileTransaction) Rollback(ctx context.Context) error {
	if ctx == nil {
		return errors.New("roll back Baseline transaction: context is required")
	}
	if transaction.closed {
		return errors.New("roll back Baseline transaction: transaction is closed")
	}
	if transaction.applied {
		return errors.New("roll back Baseline transaction: completed transaction is immutable")
	}
	if err := transaction.rollback(context.WithoutCancel(ctx)); err != nil {
		return &IncompleteRollbackError{Cause: errors.New("explicit rollback requested"), Rollback: err}
	}
	return nil
}

func (transaction *fileTransaction) Close() error {
	if transaction.closed {
		return nil
	}
	transaction.closed = true
	var result error
	if !transaction.applyAttempted && !transaction.applied {
		if err := transaction.rollback(context.Background()); err != nil {
			result = &IncompleteRollbackError{
				Cause:    errors.New("close unapplied Baseline transaction"),
				Rollback: err,
			}
		}
	}
	if err := transaction.anchored.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("close Baseline transaction repository root: %w", err))
	}
	if err := unlockTransactionFile(transaction.lock); err != nil {
		result = errors.Join(result, fmt.Errorf("unlock Baseline transaction: %w", err))
	}
	if err := transaction.lock.Close(); err != nil {
		result = errors.Join(result, fmt.Errorf("close Baseline transaction lock: %w", err))
	}
	return result
}

func (transaction *fileTransaction) failApply(ctx context.Context, cause error) error {
	rollbackErr := transaction.rollback(context.WithoutCancel(ctx))
	if rollbackErr != nil {
		return &IncompleteRollbackError{Cause: cause, Rollback: rollbackErr}
	}
	return cause
}

func (transaction *fileTransaction) stagePostimages(ctx context.Context) error {
	stagingDir := filepath.Join(transaction.stateDir, transactionStagingDirName)
	if err := os.Mkdir(stagingDir, 0o700); err != nil {
		return fmt.Errorf("create Baseline transaction staging directory: %w", err)
	}
	for index, entry := range transaction.journal.Entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		point := transactionFaultPoint{Phase: transactionPhaseStaging, Path: entry.Path}
		if err := transaction.runPhaseHook(point); err != nil {
			return err
		}
		if entry.After.Kind != PreimageRegular {
			continue
		}
		staged := transaction.stagedPath(index)
		if err := writeSyncedFile(staged, entry.After.Content, fs.FileMode(entry.After.Mode)); err != nil {
			return fmt.Errorf("stage Baseline postimage %q: %w", entry.Path, err)
		}
		if err := verifyStagedPostimage(staged, entry.After); err != nil {
			return err
		}
	}
	if err := syncDirectory(stagingDir); err != nil {
		return fmt.Errorf("sync Baseline transaction staging directory: %w", err)
	}
	transaction.journal.Phase = transactionJournalStaged
	if err := transaction.writeJournal(); err != nil {
		return err
	}
	return transaction.runPhaseHook(transactionFaultPoint{Phase: transactionPhaseStaged})
}

func (transaction *fileTransaction) revalidatePreimages(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for _, preimage := range transaction.journal.Preimages {
		if err := validatePathParents(transaction.anchored, preimage.Path); err != nil {
			return fmt.Errorf("revalidate Baseline transaction path %q: %w", preimage.Path, err)
		}
		current, err := captureTransactionState(transaction.anchored, preimage.Path)
		if err != nil {
			return fmt.Errorf("revalidate Baseline transaction preimage %q: %w", preimage.Path, err)
		}
		if !sameTransactionState(current, preimage.State) {
			return fmt.Errorf("Baseline transaction preimage is stale at %q", preimage.Path)
		}
	}
	if err := validateTransactionPostimages(transaction.anchored, transactionPostimages(transaction.journal)); err != nil {
		return err
	}
	if err := transaction.runPhaseHook(transactionFaultPoint{Phase: transactionPhasePreimagesValidated}); err != nil {
		return err
	}
	return nil
}

func (transaction *fileTransaction) replacePostimage(ctx context.Context, index int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	entry := transaction.journal.Entries[index]
	if postimageMatchesState(entry.After, entry.Before) {
		return nil
	}
	point := transactionFaultPoint{Phase: transactionPhaseReplacing, Path: entry.Path}
	if err := transaction.runPhaseHook(point); err != nil {
		return err
	}
	if err := transaction.ensureParentDirectories(entry.Path); err != nil {
		return err
	}
	current, err := captureTransactionState(transaction.anchored, entry.Path)
	if err != nil {
		return fmt.Errorf("revalidate Baseline destination %q: %w", entry.Path, err)
	}
	if !sameTransactionState(current, entry.Before) {
		return fmt.Errorf("Baseline transaction preimage is stale at %q", entry.Path)
	}
	if err := validateMutationDestination(transaction.anchored, entry.Path); err != nil {
		return err
	}
	transaction.journal.Phase = transactionJournalReplacing
	transaction.journal.MutationOrder = append(transaction.journal.MutationOrder, index)
	if err := transaction.writeJournal(); err != nil {
		return err
	}
	switch entry.After.Kind {
	case PreimageMissing:
		if err := transaction.anchored.Remove(entry.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove Baseline destination %q: %w", entry.Path, err)
		}
	case PreimageRegular:
		temporary, err := transaction.createDestinationTemporary(index, entry.After)
		if err != nil {
			return err
		}
		if err := transaction.anchored.Rename(temporary, entry.Path); err != nil {
			return fmt.Errorf("replace Baseline destination %q: %w", entry.Path, err)
		}
	default:
		return fmt.Errorf("replace Baseline destination %q: unsupported postimage kind %q", entry.Path, entry.After.Kind)
	}
	if err := transaction.syncRepositoryParent(entry.Path); err != nil {
		return err
	}
	return transaction.runPhaseHook(transactionFaultPoint{Phase: transactionPhaseReplaced, Path: entry.Path})
}

func (transaction *fileTransaction) createDestinationTemporary(
	index int,
	postimage Postimage,
) (string, error) {
	temporary := transaction.destinationTemporaryPath(index, postimage.Path, "apply")
	if err := transaction.recordTemporaryPath(temporary); err != nil {
		return "", err
	}
	file, err := transaction.anchored.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", fmt.Errorf("create Baseline destination temporary %q: %w", temporary, err)
	}
	staged, err := os.Open(transaction.stagedPath(index))
	if err != nil {
		_ = file.Close()
		return "", fmt.Errorf("open staged Baseline postimage %q: %w", postimage.Path, err)
	}
	_, copyErr := io.Copy(file, staged)
	stagedCloseErr := staged.Close()
	chmodErr := file.Chmod(fs.FileMode(postimage.Mode))
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(copyErr, stagedCloseErr, chmodErr, syncErr, closeErr); err != nil {
		return "", fmt.Errorf("write Baseline destination temporary %q: %w", temporary, err)
	}
	if err := transaction.syncRepositoryParent(temporary); err != nil {
		return "", err
	}
	return temporary, nil
}

func (transaction *fileTransaction) ensureParentDirectories(relative string) error {
	parent := path.Dir(relative)
	if parent == "." {
		return nil
	}
	current := ""
	for _, component := range strings.Split(parent, "/") {
		if current == "" {
			current = component
		} else {
			current += "/" + component
		}
		info, err := transaction.anchored.Lstat(current)
		switch {
		case err == nil:
			if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
				return fmt.Errorf("Baseline transaction parent %q must be a real directory", current)
			}
		case errors.Is(err, fs.ErrNotExist):
			if err := transaction.recordCreatedDirectory(current); err != nil {
				return err
			}
			if err := transaction.anchored.Mkdir(current, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
				return fmt.Errorf("create Baseline destination parent %q: %w", current, err)
			}
			if err := transaction.syncRepositoryParent(current); err != nil {
				return err
			}
		default:
			return fmt.Errorf("inspect Baseline destination parent %q: %w", current, err)
		}
	}
	return nil
}

func (transaction *fileTransaction) recordCreatedDirectory(relative string) error {
	for _, recorded := range transaction.journal.CreatedDirectories {
		if recorded == relative {
			return nil
		}
	}
	transaction.journal.CreatedDirectories = append(transaction.journal.CreatedDirectories, relative)
	return transaction.writeJournal()
}

func (transaction *fileTransaction) recordTemporaryPath(relative string) error {
	for _, recorded := range transaction.journal.TemporaryPaths {
		if recorded == relative {
			return nil
		}
	}
	transaction.journal.TemporaryPaths = append(transaction.journal.TemporaryPaths, relative)
	return transaction.writeJournal()
}

func (transaction *fileTransaction) rollback(ctx context.Context) error {
	transaction.journal.Phase = transactionJournalRollingBack
	if err := transaction.writeJournal(); err != nil {
		return err
	}
	for len(transaction.journal.MutationOrder) > 0 {
		position := len(transaction.journal.MutationOrder) - 1
		index := transaction.journal.MutationOrder[position]
		entry := transaction.journal.Entries[index]
		if err := transaction.runPhaseHook(transactionFaultPoint{
			Phase: transactionPhaseRollingBack,
			Path:  entry.Path,
		}); err != nil {
			return err
		}
		if err := restoreTransactionPreimage(transaction, index, entry); err != nil {
			return err
		}
		transaction.journal.MutationOrder = transaction.journal.MutationOrder[:position]
		if err := transaction.writeJournal(); err != nil {
			return err
		}
	}
	for len(transaction.journal.CreatedDirectories) > 0 {
		position := len(transaction.journal.CreatedDirectories) - 1
		relative := transaction.journal.CreatedDirectories[position]
		if err := transaction.runPhaseHook(transactionFaultPoint{
			Phase: transactionPhaseRollingBack,
			Path:  relative,
		}); err != nil {
			return err
		}
		if err := transaction.anchored.Remove(relative); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove rolled-back Baseline directory %q: %w", relative, err)
		}
		if err := transaction.syncRepositoryParent(relative); err != nil {
			return err
		}
		transaction.journal.CreatedDirectories = transaction.journal.CreatedDirectories[:position]
		if err := transaction.writeJournal(); err != nil {
			return err
		}
	}
	for _, temporary := range transaction.journal.TemporaryPaths {
		if err := removeTransactionTemporary(transaction.anchored, temporary); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("complete Baseline rollback after cancellation: %w", err)
	}
	return transaction.cleanupState()
}

func restoreTransactionPreimage(
	transaction *fileTransaction,
	index int,
	entry transactionJournalEntry,
) error {
	if err := transaction.ensureParentDirectories(entry.Path); err != nil {
		return err
	}
	if err := validateMutationDestination(transaction.anchored, entry.Path); err != nil {
		return fmt.Errorf("restore Baseline preimage %q: %w", entry.Path, err)
	}
	switch entry.Before.Kind {
	case PreimageMissing:
		if err := transaction.anchored.Remove(entry.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove rolled-back Baseline path %q: %w", entry.Path, err)
		}
	case PreimageRegular:
		temporary := transaction.destinationTemporaryPath(index, entry.Path, "restore")
		if err := transaction.recordTemporaryPath(temporary); err != nil {
			return err
		}
		if err := removeTransactionTemporary(transaction.anchored, temporary); err != nil {
			return err
		}
		file, err := transaction.anchored.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return fmt.Errorf("create Baseline rollback temporary %q: %w", temporary, err)
		}
		_, writeErr := file.Write(entry.Before.Content)
		chmodErr := file.Chmod(fs.FileMode(entry.Before.Mode))
		syncErr := file.Sync()
		closeErr := file.Close()
		if err := errors.Join(writeErr, chmodErr, syncErr, closeErr); err != nil {
			return fmt.Errorf("write Baseline rollback temporary %q: %w", temporary, err)
		}
		if err := transaction.anchored.Rename(temporary, entry.Path); err != nil {
			return fmt.Errorf("restore Baseline preimage %q: %w", entry.Path, err)
		}
	default:
		return fmt.Errorf("restore Baseline preimage %q: unsupported kind %q", entry.Path, entry.Before.Kind)
	}
	return transaction.syncRepositoryParent(entry.Path)
}

func (transaction *fileTransaction) syncRepositoryParent(relative string) error {
	parent := path.Dir(relative)
	if parent == "." {
		parent = ""
	}
	name := transaction.root
	if parent != "" {
		name = filepath.Join(transaction.root, filepath.FromSlash(parent))
	}
	if err := syncDirectory(name); err != nil {
		return fmt.Errorf("sync Baseline destination directory %q: %w", parent, err)
	}
	return nil
}

func (transaction *fileTransaction) runPhaseHook(point transactionFaultPoint) error {
	if transaction.phaseHook == nil {
		return nil
	}
	if err := transaction.phaseHook(point); err != nil {
		return fmt.Errorf("Baseline transaction phase %s: %w", point, err)
	}
	return nil
}

func (transaction *fileTransaction) stagedPath(index int) string {
	return filepath.Join(
		transaction.stateDir,
		transactionStagingDirName,
		fmt.Sprintf("%06d.postimage", index),
	)
}

func (transaction *fileTransaction) destinationTemporaryPath(
	index int,
	relative string,
	purpose string,
) string {
	token := strings.TrimPrefix(transaction.journal.PlanDigest, "sha256:")
	if len(token) > 12 {
		token = token[:12]
	}
	name := fmt.Sprintf(".roundfix-baseline-%s-%06d-%s.tmp", token, index, purpose)
	parent := path.Dir(relative)
	if parent == "." {
		return name
	}
	return path.Join(parent, name)
}

func (transaction *fileTransaction) writeJournal() error {
	if err := validateTransactionJournal(transaction.journal); err != nil {
		return fmt.Errorf("validate Baseline transaction journal: %w", err)
	}
	data, err := json.Marshal(transaction.journal)
	if err != nil {
		return fmt.Errorf("serialize Baseline transaction journal: %w", err)
	}
	next := filepath.Join(transaction.stateDir, transactionJournalNext)
	file, err := os.OpenFile(next, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open Baseline transaction journal: %w", err)
	}
	_, writeErr := file.Write(data)
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		return fmt.Errorf("write Baseline transaction journal: %w", err)
	}
	if err := os.Rename(next, filepath.Join(transaction.stateDir, transactionJournalName)); err != nil {
		return fmt.Errorf("replace Baseline transaction journal: %w", err)
	}
	if err := syncDirectory(transaction.stateDir); err != nil {
		return fmt.Errorf("sync Baseline transaction journal directory: %w", err)
	}
	return nil
}

func (transaction *fileTransaction) cleanupState() error {
	info, err := os.Lstat(transaction.stateDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect Baseline transaction recovery state: %w", err)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("Baseline transaction recovery state is not a real directory")
	}
	if err := os.RemoveAll(transaction.stateDir); err != nil {
		return fmt.Errorf("remove Baseline transaction recovery state: %w", err)
	}
	if err := syncDirectory(transaction.privateDir); err != nil {
		return fmt.Errorf("sync Baseline transaction private directory: %w", err)
	}
	return nil
}

func transactionPrivateDirectory(ctx context.Context, root string) (string, error) {
	output, err := (ExecGitRunner{}).RunGit(ctx, root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return "", fmt.Errorf("resolve Git-private Baseline transaction path: %w", err)
	}
	gitDir := filepath.Clean(strings.TrimSpace(output))
	if !filepath.IsAbs(gitDir) {
		return "", fmt.Errorf("resolve Git-private Baseline transaction path: Git returned non-absolute path %q", output)
	}
	info, err := os.Lstat(gitDir)
	if err != nil {
		return "", fmt.Errorf("inspect Git-private Baseline transaction path: %w", err)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return "", fmt.Errorf("inspect Git-private Baseline transaction path: %q must be a real directory", gitDir)
	}
	return filepath.Join(gitDir, "roundfix", "baseline-transaction"), nil
}

func recoverInterruptedTransaction(
	ctx context.Context,
	root string,
	anchored *os.Root,
	privateDir string,
	stateDir string,
) error {
	info, err := os.Lstat(stateDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect interrupted Baseline transaction: %w", err)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("interrupted Baseline transaction recovery state is unsafe")
	}
	journal, err := readTransactionJournal(stateDir)
	if errors.Is(err, fs.ErrNotExist) {
		entries, readErr := os.ReadDir(stateDir)
		if readErr != nil {
			return fmt.Errorf("inspect empty Baseline transaction recovery state: %w", readErr)
		}
		if len(entries) != 0 {
			return errors.New("interrupted Baseline transaction has recovery state without a journal")
		}
		if err := os.Remove(stateDir); err != nil {
			return fmt.Errorf("remove empty Baseline transaction recovery state: %w", err)
		}
		return syncDirectory(privateDir)
	}
	if err != nil {
		return err
	}
	recovery := &fileTransaction{
		root:       root,
		anchored:   anchored,
		privateDir: privateDir,
		stateDir:   stateDir,
		journal:    journal,
	}
	if journal.Phase == transactionJournalCommitted {
		for _, entry := range journal.Entries {
			if err := verifyTransactionPostimage(anchored, entry.After); err != nil {
				return fmt.Errorf("committed Baseline transaction cannot be conclusively recovered: %w", err)
			}
		}
		if err := recovery.cleanupState(); err != nil {
			return fmt.Errorf("clean committed Baseline transaction recovery state: %w", err)
		}
		return nil
	}
	if err := recovery.rollback(ctx); err != nil {
		return &IncompleteRollbackError{
			Cause:    errors.New("recover interrupted Baseline transaction"),
			Rollback: err,
		}
	}
	return nil
}

func readTransactionJournal(stateDir string) (transactionJournal, error) {
	journalPath := filepath.Join(stateDir, transactionJournalName)
	data, err := readBoundedTransactionFile(journalPath)
	if errors.Is(err, fs.ErrNotExist) {
		nextPath := filepath.Join(stateDir, transactionJournalNext)
		data, err = readBoundedTransactionFile(nextPath)
		if err != nil {
			return transactionJournal{}, err
		}
		if err := os.Rename(nextPath, journalPath); err != nil {
			return transactionJournal{}, fmt.Errorf("recover pending Baseline transaction journal: %w", err)
		}
	}
	if err != nil {
		return transactionJournal{}, fmt.Errorf("read Baseline transaction journal: %w", err)
	}
	var journal transactionJournal
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&journal); err != nil {
		return transactionJournal{}, fmt.Errorf("decode Baseline transaction journal: %w", err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return transactionJournal{}, errors.New("decode Baseline transaction journal: trailing JSON data")
	}
	if err := validateTransactionJournal(journal); err != nil {
		return transactionJournal{}, fmt.Errorf("validate Baseline transaction journal: %w", err)
	}
	return journal, nil
}

func readBoundedTransactionFile(name string) ([]byte, error) {
	info, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maxTransactionJournalSize {
		return nil, errors.New("Baseline transaction journal is unsafe or too large")
	}
	return os.ReadFile(name)
}

func validateTransactionJournal(journal transactionJournal) error {
	if journal.SchemaVersion != transactionJournalSchema {
		return fmt.Errorf("unsupported schema %q", journal.SchemaVersion)
	}
	if !strings.HasPrefix(journal.PlanDigest, "sha256:") ||
		len(strings.TrimPrefix(journal.PlanDigest, "sha256:")) != sha256.Size*2 {
		return errors.New("plan digest is invalid")
	}
	switch journal.Phase {
	case transactionJournalPrepared,
		transactionJournalStaged,
		transactionJournalReplacing,
		transactionJournalVerifying,
		transactionJournalCommitted,
		transactionJournalRollingBack:
	default:
		return fmt.Errorf("phase %q is invalid", journal.Phase)
	}
	if journal.Preimages == nil || journal.Entries == nil || journal.MutationOrder == nil ||
		journal.CreatedDirectories == nil || journal.TemporaryPaths == nil {
		return errors.New("collections must be JSON arrays")
	}
	preimagePaths := make(map[string]transactionFileState, len(journal.Preimages))
	for index, preimage := range journal.Preimages {
		if !repositoryPathIsSafe(preimage.Path) {
			return fmt.Errorf("preimage %d has unsafe path", index)
		}
		if index > 0 && journal.Preimages[index-1].Path >= preimage.Path {
			return errors.New("preimages must be in lexical path order")
		}
		if err := validateTransactionFileState(preimage.State); err != nil {
			return fmt.Errorf("preimage %q: %w", preimage.Path, err)
		}
		preimagePaths[preimage.Path] = preimage.State
	}
	for index, entry := range journal.Entries {
		if !repositoryPathIsSafe(entry.Path) || entry.After.Path != entry.Path {
			return fmt.Errorf("entry %d has unsafe or mismatched path", index)
		}
		if index > 0 && journal.Entries[index-1].Path >= entry.Path {
			return errors.New("entries must be in lexical path order")
		}
		if err := validateTransactionPostimage(entry.After); err != nil {
			return err
		}
		if err := validateTransactionFileState(entry.Before); err != nil {
			return fmt.Errorf("entry %q preimage: %w", entry.Path, err)
		}
		if preimage, ok := preimagePaths[entry.Path]; !ok || !sameTransactionState(preimage, entry.Before) {
			return fmt.Errorf("entry %q has no matching complete preimage", entry.Path)
		}
	}
	previous := -1
	for _, index := range journal.MutationOrder {
		if index < 0 || index >= len(journal.Entries) || index <= previous {
			return errors.New("mutation order is invalid")
		}
		previous = index
	}
	if err := validateSafeUniquePaths(journal.CreatedDirectories, false); err != nil {
		return fmt.Errorf("created directories: %w", err)
	}
	if err := validateSafeUniquePaths(journal.TemporaryPaths, true); err != nil {
		return fmt.Errorf("temporary paths: %w", err)
	}
	return nil
}

func validateSafeUniquePaths(paths []string, temporary bool) error {
	seen := make(map[string]struct{}, len(paths))
	for _, relative := range paths {
		if !repositoryPathIsSafe(relative) {
			return fmt.Errorf("unsafe path %q", relative)
		}
		if temporary && !strings.HasPrefix(path.Base(relative), ".roundfix-baseline-") {
			return fmt.Errorf("unrecognized temporary path %q", relative)
		}
		if _, duplicate := seen[relative]; duplicate {
			return fmt.Errorf("duplicate path %q", relative)
		}
		seen[relative] = struct{}{}
	}
	return nil
}

func validateTransactionFileState(state transactionFileState) error {
	switch state.Kind {
	case PreimageMissing:
		if state.Exists || state.Mode != 0 || state.LinkTarget != "" ||
			state.Bytes != 0 || len(state.Content) != 0 || state.ContentIdentity != "" {
			return errors.New("missing state carries file data")
		}
	case PreimageRegular:
		if !state.Exists || !validTransactionMode(state.Mode) ||
			state.Bytes != int64(len(state.Content)) ||
			state.ContentIdentity != transactionContentIdentity(state.Content) {
			return errors.New("regular state is invalid")
		}
	case PreimageSymlink:
		if !state.Exists || state.LinkTarget == "" || len(state.Content) != 0 {
			return errors.New("symlink state is invalid")
		}
	case PreimageDirectory, PreimageSpecial:
		if !state.Exists || len(state.Content) != 0 {
			return errors.New("non-file state is invalid")
		}
	default:
		return fmt.Errorf("unsupported kind %q", state.Kind)
	}
	return nil
}

func validateTransactionPostimages(root *os.Root, postimages []Postimage) error {
	for _, postimage := range postimages {
		if err := validateTransactionPostimage(postimage); err != nil {
			return err
		}
		if err := validatePathParents(root, postimage.Path); err != nil {
			return fmt.Errorf("validate Baseline destination %q: %w", postimage.Path, err)
		}
		if err := validateMutationDestination(root, postimage.Path); err != nil {
			return err
		}
	}
	return nil
}

func validateTransactionPostimage(postimage Postimage) error {
	if !repositoryPathIsSafe(postimage.Path) {
		return fmt.Errorf("unsafe Baseline postimage path %q", postimage.Path)
	}
	switch postimage.Kind {
	case PreimageMissing:
		if postimage.Mode != 0 || len(postimage.Content) != 0 || postimage.ContentIdentity != "" {
			return fmt.Errorf("missing Baseline postimage %q carries file data", postimage.Path)
		}
	case PreimageRegular:
		if !validTransactionMode(postimage.Mode) {
			return fmt.Errorf("Baseline postimage %q has unsafe mode %#o", postimage.Path, postimage.Mode)
		}
		if transactionContentIdentity(postimage.Content) != postimage.ContentIdentity {
			return fmt.Errorf("Baseline postimage %q content identity mismatch", postimage.Path)
		}
	default:
		return fmt.Errorf("Baseline postimage %q has unsupported kind %q", postimage.Path, postimage.Kind)
	}
	return nil
}

func validTransactionMode(mode uint32) bool {
	return mode != 0 && mode&^uint32(fs.ModePerm) == 0
}

func validatePathParents(root *os.Root, relative string) error {
	parent := path.Dir(relative)
	if parent == "." {
		return nil
	}
	current := ""
	for _, component := range strings.Split(parent, "/") {
		if current == "" {
			current = component
		} else {
			current += "/" + component
		}
		info, err := root.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("path component %q is a symlink", current)
		}
		if !info.IsDir() {
			return fmt.Errorf("path component %q is not a directory", current)
		}
	}
	return nil
}

func validateMutationDestination(root *os.Root, relative string) error {
	info, err := root.Lstat(relative)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect Baseline destination %q: %w", relative, err)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return fmt.Errorf("Baseline destination %q is a symlink", relative)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("Baseline destination %q is not a regular file", relative)
	}
	return nil
}

func captureTransactionEntries(
	root *os.Root,
	document PlanDocument,
) ([]transactionJournalPreimage, []transactionJournalEntry, error) {
	preimages := make([]transactionJournalPreimage, len(document.Preimages))
	exactByPath := make(map[string]transactionFileState, len(document.Preimages))
	for index, approved := range document.Preimages {
		if err := validatePathParents(root, approved.Path); err != nil {
			return nil, nil, fmt.Errorf("journal Baseline preimage %q: %w", approved.Path, err)
		}
		current, err := captureTransactionState(root, approved.Path)
		if err != nil {
			return nil, nil, fmt.Errorf("journal Baseline preimage %q: %w", approved.Path, err)
		}
		if !transactionStateMatchesPreimage(current, approved) {
			return nil, nil, fmt.Errorf("Baseline transaction preimage is stale at %q", approved.Path)
		}
		preimages[index] = transactionJournalPreimage{Path: approved.Path, State: current}
		exactByPath[approved.Path] = current
	}
	entries := make([]transactionJournalEntry, len(document.Postimages))
	for index, postimage := range document.Postimages {
		current := exactByPath[postimage.Path]
		if current.Kind != PreimageMissing && current.Kind != PreimageRegular {
			return nil, nil, fmt.Errorf("Baseline destination %q has unsafe preimage kind %q", postimage.Path, current.Kind)
		}
		entries[index] = transactionJournalEntry{
			Path:   postimage.Path,
			Before: current,
			After:  clonePostimage(postimage),
		}
	}
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Path < entries[right].Path
	})
	return preimages, entries, nil
}

func captureTransactionState(root *os.Root, relative string) (transactionFileState, error) {
	info, err := root.Lstat(relative)
	if errors.Is(err, fs.ErrNotExist) {
		return transactionFileState{Kind: PreimageMissing}, nil
	}
	if err != nil {
		return transactionFileState{}, err
	}
	state := transactionFileState{
		Exists: true,
		Mode:   uint32(info.Mode()),
		Bytes:  info.Size(),
	}
	switch {
	case info.Mode().IsRegular():
		state.Kind = PreimageRegular
		file, err := root.Open(relative)
		if err != nil {
			return transactionFileState{}, err
		}
		openedInfo, statErr := file.Stat()
		content, readErr := io.ReadAll(file)
		closeErr := file.Close()
		if err := errors.Join(statErr, readErr, closeErr); err != nil {
			return transactionFileState{}, err
		}
		if !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
			return transactionFileState{}, errors.New("path changed while its exact preimage was read")
		}
		state.Content = content
		state.ContentIdentity = transactionContentIdentity(content)
	case info.Mode()&fs.ModeSymlink != 0:
		state.Kind = PreimageSymlink
		state.LinkTarget, err = root.Readlink(relative)
		if err != nil {
			return transactionFileState{}, err
		}
	case info.IsDir():
		state.Kind = PreimageDirectory
	default:
		state.Kind = PreimageSpecial
	}
	return state, nil
}

func transactionStateMatchesPreimage(state transactionFileState, approved Preimage) bool {
	if state.Exists != approved.Exists ||
		state.Kind != approved.Kind ||
		state.Mode != approved.Mode ||
		state.LinkTarget != approved.LinkTarget ||
		state.Bytes != approved.Bytes {
		return false
	}
	if state.Kind == PreimageRegular {
		if int64(len(state.Content)) != approved.Bytes {
			return false
		}
		return approved.ContentIdentity == "" || state.ContentIdentity == approved.ContentIdentity
	}
	return true
}

func sameTransactionState(left, right transactionFileState) bool {
	return left.Exists == right.Exists &&
		left.Kind == right.Kind &&
		left.Mode == right.Mode &&
		left.LinkTarget == right.LinkTarget &&
		left.Bytes == right.Bytes &&
		left.ContentIdentity == right.ContentIdentity &&
		bytes.Equal(left.Content, right.Content)
}

func postimageMatchesState(postimage Postimage, state transactionFileState) bool {
	switch postimage.Kind {
	case PreimageMissing:
		return state.Kind == PreimageMissing
	case PreimageRegular:
		return state.Kind == PreimageRegular &&
			state.Mode == postimage.Mode &&
			state.ContentIdentity == postimage.ContentIdentity
	default:
		return false
	}
}

func verifyStagedPostimage(name string, postimage Postimage) error {
	info, err := os.Lstat(name)
	if err != nil {
		return fmt.Errorf("inspect staged Baseline postimage %q: %w", postimage.Path, err)
	}
	if !info.Mode().IsRegular() || uint32(info.Mode()) != postimage.Mode {
		return fmt.Errorf("staged Baseline postimage %q mode mismatch", postimage.Path)
	}
	content, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read staged Baseline postimage %q: %w", postimage.Path, err)
	}
	if transactionContentIdentity(content) != postimage.ContentIdentity {
		return fmt.Errorf("staged Baseline postimage %q content mismatch", postimage.Path)
	}
	return nil
}

func verifyTransactionPostimage(root *os.Root, postimage Postimage) error {
	if err := validatePathParents(root, postimage.Path); err != nil {
		return fmt.Errorf("verify Baseline postimage %q: %w", postimage.Path, err)
	}
	current, err := captureTransactionState(root, postimage.Path)
	if err != nil {
		return fmt.Errorf("verify Baseline postimage %q: %w", postimage.Path, err)
	}
	if !postimageMatchesState(postimage, current) {
		return fmt.Errorf("verify Baseline postimage %q: bytes, kind, or mode differ", postimage.Path)
	}
	return nil
}

func writeSyncedFile(name string, content []byte, mode fs.FileMode) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(content)
	chmodErr := file.Chmod(mode)
	syncErr := file.Sync()
	closeErr := file.Close()
	return errors.Join(writeErr, chmodErr, syncErr, closeErr)
}

func syncDirectory(name string) error {
	directory, err := os.Open(name)
	if err != nil {
		return err
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	return errors.Join(syncErr, closeErr)
}

func removeTransactionTemporary(root *os.Root, relative string) error {
	info, err := root.Lstat(relative)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect Baseline transaction temporary %q: %w", relative, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("Baseline transaction temporary %q is unsafe", relative)
	}
	if err := root.Remove(relative); err != nil {
		return fmt.Errorf("remove Baseline transaction temporary %q: %w", relative, err)
	}
	return nil
}

func transactionContentIdentity(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func transactionPostimages(journal transactionJournal) []Postimage {
	postimages := make([]Postimage, len(journal.Entries))
	for index, entry := range journal.Entries {
		postimages[index] = entry.After
	}
	return postimages
}

func clonePostimage(postimage Postimage) Postimage {
	postimage.Content = append([]byte(nil), postimage.Content...)
	return postimage
}
