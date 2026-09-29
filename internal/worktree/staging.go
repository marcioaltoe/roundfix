package worktree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/store"
)

const (
	stagingRootPrefix             = "roundfix-carry-forward-"
	stagingWorktreeName           = "worktree"
	stagingOwnerFilename          = "owner.json"
	stagingInitializingLockReason = "initializing"
)

type stagingOwner struct {
	PID      int    `json:"pid"`
	Identity string `json:"identity"`
}

// CarryForwardStaging is a disposable detached worktree used to prove a
// carry-forward set before changing the user's checkout.
type CarryForwardStaging struct {
	Worktree string

	repository string
	root       string
}

// StagingCandidate describes one registered carry-forward staging worktree
// and the proof that permits or refuses its release.
type StagingCandidate struct {
	Worktree      string
	OwnerPID      int
	Proof         string
	Stale         bool
	RefusalReason string
}

type stagingRegistration struct {
	worktree   string
	lockReason string
}

// AddCarryForwardStaging creates an owned detached staging worktree. An empty
// parent uses the process temporary directory.
func AddCarryForwardStaging(
	ctx context.Context,
	repository string,
	head string,
	parent string,
) (CarryForwardStaging, error) {
	return addCarryForwardStaging(ctx, execGitRunner{}, repository, head, parent)
}

func addCarryForwardStaging(
	ctx context.Context,
	runner gitRunner,
	repository string,
	head string,
	parent string,
) (CarryForwardStaging, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(parent) == "" {
		parent = os.TempDir()
	}
	root, err := os.MkdirTemp(parent, stagingRootPrefix)
	if err != nil {
		return CarryForwardStaging{}, fmt.Errorf("create carry-forward staging directory: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(root); resolveErr == nil {
		root = resolved
	}
	staging := CarryForwardStaging{
		Worktree:   filepath.Join(root, stagingWorktreeName),
		repository: repository,
		root:       root,
	}
	cleanupRoot := func(cause error) (CarryForwardStaging, error) {
		removeErr := os.RemoveAll(root)
		if removeErr != nil {
			removeErr = fmt.Errorf("remove carry-forward staging directory: %w", removeErr)
		}
		return CarryForwardStaging{}, errors.Join(cause, removeErr)
	}

	identity, err := store.OwnerProcessIdentity(ctx, os.Getpid())
	if err != nil {
		return cleanupRoot(fmt.Errorf("read carry-forward staging owner identity: %w", err))
	}
	ownerData, err := json.Marshal(stagingOwner{PID: os.Getpid(), Identity: identity})
	if err != nil {
		return cleanupRoot(fmt.Errorf("encode carry-forward staging owner: %w", err))
	}
	ownerPath := filepath.Join(root, stagingOwnerFilename)
	if err := os.WriteFile(ownerPath, append(ownerData, '\n'), 0o600); err != nil {
		return cleanupRoot(fmt.Errorf("write carry-forward staging owner %q: %w", ownerPath, err))
	}
	if _, err := runWorktreeCommand(
		ctx,
		runner,
		repository,
		"worktree", "add", "--detach", staging.Worktree, head,
	); err != nil {
		return cleanupRoot(fmt.Errorf("create carry-forward staging Worktree: %w", err))
	}
	return staging, nil
}

// Remove unregisters a staging worktree through the administration lock and
// removes only its staging root.
func (staging CarryForwardStaging) Remove(ctx context.Context) error {
	_, worktreeErr := runWorktreeCommand(
		ctx,
		execGitRunner{},
		staging.repository,
		"worktree", "remove", "--force", staging.Worktree,
	)
	if worktreeErr != nil {
		worktreeErr = fmt.Errorf("remove carry-forward staging Worktree: %w", worktreeErr)
	}
	removeErr := os.RemoveAll(staging.root)
	if removeErr != nil {
		removeErr = fmt.Errorf("remove carry-forward staging directory: %w", removeErr)
	}
	return errors.Join(worktreeErr, removeErr)
}

// InspectCarryForwardStaging reports registered Roundfix staging worktrees.
// Staleness is based only on owner-death or legacy locked-initializing proof.
func InspectCarryForwardStaging(ctx context.Context, repository string) ([]StagingCandidate, error) {
	return inspectCarryForwardStaging(ctx, execGitRunner{}, repository, false)
}

func inspectCarryForwardStaging(
	ctx context.Context,
	runner gitRunner,
	repository string,
	lockHeld bool,
) ([]StagingCandidate, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var output string
	var err error
	if lockHeld {
		output, err = runner.Run(ctx, repository, "worktree", "list", "--porcelain", "-z")
	} else {
		output, err = runWorktreeCommand(ctx, runner, repository, "worktree", "list", "--porcelain", "-z")
	}
	if err != nil {
		return nil, fmt.Errorf("list carry-forward staging worktrees: %w", err)
	}
	registrations := parseStagingRegistrations(output)
	candidates := make([]StagingCandidate, 0, len(registrations))
	for _, registration := range registrations {
		if filepath.Base(registration.worktree) != stagingWorktreeName ||
			!strings.HasPrefix(filepath.Base(filepath.Dir(registration.worktree)), stagingRootPrefix) {
			continue
		}
		candidates = append(candidates, inspectStagingRegistration(ctx, registration))
	}
	sort.Slice(candidates, func(left, right int) bool {
		return candidates[left].Worktree < candidates[right].Worktree
	})
	return candidates, nil
}

func parseStagingRegistrations(output string) []stagingRegistration {
	registrations := make([]stagingRegistration, 0)
	current := stagingRegistration{}
	appendCurrent := func() {
		if current.worktree != "" {
			registrations = append(registrations, current)
		}
		current = stagingRegistration{}
	}
	for _, field := range strings.Split(output, "\x00") {
		if field == "" {
			appendCurrent()
			continue
		}
		if value, ok := strings.CutPrefix(field, "worktree "); ok {
			if current.worktree != "" {
				appendCurrent()
			}
			current.worktree = filepath.Clean(value)
			continue
		}
		if field == "locked" {
			current.lockReason = ""
			continue
		}
		if value, ok := strings.CutPrefix(field, "locked "); ok {
			current.lockReason = value
		}
	}
	appendCurrent()
	return registrations
}

func inspectStagingRegistration(ctx context.Context, registration stagingRegistration) StagingCandidate {
	candidate := StagingCandidate{Worktree: registration.worktree}
	ownerPath := filepath.Join(filepath.Dir(registration.worktree), stagingOwnerFilename)
	data, err := os.ReadFile(ownerPath)
	if errors.Is(err, os.ErrNotExist) {
		if registration.lockReason == stagingInitializingLockReason {
			candidate.Stale = true
			candidate.Proof = "legacy staging registration has no owner record and is locked initializing"
			return candidate
		}
		candidate.RefusalReason = "legacy staging registration has no owner record and is not locked initializing"
		if registration.lockReason != "" {
			candidate.RefusalReason = fmt.Sprintf(
				"legacy staging registration has no owner record and is locked %q, not locked initializing",
				registration.lockReason,
			)
		}
		return candidate
	}
	if err != nil {
		candidate.RefusalReason = fmt.Sprintf("owner record cannot be read: %v", err)
		return candidate
	}
	var owner stagingOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		candidate.RefusalReason = fmt.Sprintf("owner record cannot be decoded: %v", err)
		return candidate
	}
	candidate.OwnerPID = owner.PID
	if owner.PID <= 0 || owner.Identity == "" {
		candidate.RefusalReason = "owner record does not contain a comparable PID and process identity"
		return candidate
	}
	liveIdentity, err := store.OwnerProcessIdentity(ctx, owner.PID)
	if err != nil {
		candidate.Stale = true
		candidate.Proof = fmt.Sprintf("owner process %d identity cannot be read: %v", owner.PID, err)
		return candidate
	}
	if liveIdentity != owner.Identity {
		candidate.Stale = true
		candidate.Proof = fmt.Sprintf("owner process %d identity differs from the owner record", owner.PID)
		return candidate
	}
	candidate.RefusalReason = fmt.Sprintf("owner process %d is live with the recorded identity", owner.PID)
	return candidate
}

// ReleaseCarryForwardStaging re-proves one candidate and removes it while
// holding the repository administration lock for the complete act.
func ReleaseCarryForwardStaging(
	ctx context.Context,
	repository string,
	candidate StagingCandidate,
) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runner := execGitRunner{}
	return withWorktreeAdminLock(ctx, runner, repository, func() error {
		candidates, err := inspectCarryForwardStaging(ctx, runner, repository, true)
		if err != nil {
			return err
		}
		var current *StagingCandidate
		for index := range candidates {
			if candidates[index].Worktree == candidate.Worktree {
				current = &candidates[index]
				break
			}
		}
		if current == nil {
			return fmt.Errorf("release carry-forward staging %q: registration is absent", candidate.Worktree)
		}
		if !current.Stale {
			return fmt.Errorf("release carry-forward staging %q: %s", candidate.Worktree, current.RefusalReason)
		}
		if _, err := runner.Run(
			ctx,
			repository,
			"worktree", "remove", "--force", "--force", current.Worktree,
		); err != nil {
			return fmt.Errorf("release carry-forward staging Worktree %q: %w", current.Worktree, err)
		}
		root := filepath.Dir(current.Worktree)
		if err := os.RemoveAll(root); err != nil {
			return fmt.Errorf("remove carry-forward staging directory %q: %w", root, err)
		}
		return nil
	})
}
