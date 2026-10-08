// Suite: carry-forward staging ownership and reclamation.
// Invariant: Roundfix removes a registered staging worktree only after proving its owner is stale.
// Boundary IN: owner records, process identity, Git registration, and the worktree administration lock.
// Boundary OUT: reconcile reporting, owned by internal/cli/reconcile_staging_test.go.
package worktree

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestCarryForwardStagingRecordsItsOwner(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	staging, err := AddCarryForwardStaging(t.Context(), repoDir, head, t.TempDir())
	if err != nil {
		t.Fatalf("add carry-forward staging: %v", err)
	}
	t.Cleanup(func() { removeStagingForTest(t, staging) })

	record := readStagingOwnerForTest(t, filepath.Join(filepath.Dir(staging.Worktree), stagingOwnerFilename))
	wantIdentity, err := store.OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("read current process identity: %v", err)
	}
	if record.PID != os.Getpid() || record.Identity != wantIdentity {
		t.Fatalf("staging owner = pid:%d identity:%q, want pid:%d identity:%q", record.PID, record.Identity, os.Getpid(), wantIdentity)
	}
	if registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain"); !strings.Contains(registered, "worktree "+staging.Worktree) {
		t.Fatalf("staging worktree %q is not registered:\n%s", staging.Worktree, registered)
	}
}

func TestCarryForwardStagingWaitsForTheAdminLock(t *testing.T) {
	t.Parallel()
	commonDir := t.TempDir()
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	holder := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(ctx context.Context, _ string, _ ...string) (string, error) {
			started <- struct{}{}
			select {
			case <-release:
				return "", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
	}
	holderResult := make(chan error, 1)
	go func() {
		_, err := runWorktreeCommand(context.Background(), holder, commonDir, "worktree", "prune")
		holderResult <- err
	}()
	waitForAdminCommandStart(t, started)

	waiter := &stagingAdminTestRunner{
		commonDir: commonDir,
		resolved:  make(chan struct{}, 1),
	}
	waitCtx, cancelWait := context.WithCancel(context.Background())
	waitResult := make(chan error, 1)
	parent := t.TempDir()
	go func() {
		_, err := addCarryForwardStaging(waitCtx, waiter, commonDir, "HEAD", parent)
		waitResult <- err
	}()
	waitForAdminCommandStart(t, waiter.resolved)
	cancelWait()
	err := <-waitResult
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("add while administration lock is held: %v, want context canceled", err)
	}
	if waiter.commandRan {
		t.Fatal("carry-forward worktree add ran while another administration lock holder was active")
	}

	close(release)
	if err := <-holderResult; err != nil {
		t.Fatalf("release administration lock holder: %v", err)
	}
}

func TestCarryForwardStagingSweepReleasesADeadOwnersWorktree(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	staging, err := AddCarryForwardStaging(t.Context(), repoDir, head, t.TempDir())
	if err != nil {
		t.Fatalf("add carry-forward staging: %v", err)
	}
	t.Cleanup(func() { removeStagingForTest(t, staging) })
	deadPID, deadIdentity := exitedOwnerForTest(t)
	writeStagingOwnerForTest(t, filepath.Join(filepath.Dir(staging.Worktree), stagingOwnerFilename), stagingOwner{
		PID:      deadPID,
		Identity: deadIdentity,
	})

	candidate := singleStagingCandidateForTest(t, repoDir)
	if !candidate.Stale || candidate.OwnerPID != deadPID || !strings.Contains(candidate.Proof, "owner process") {
		t.Fatalf("dead-owner candidate = %+v, want stale proof for PID %d", candidate, deadPID)
	}
	if err := ReleaseCarryForwardStaging(t.Context(), repoDir, candidate); err != nil {
		t.Fatalf("release dead-owner staging: %v", err)
	}
	assertStagingReleasedForTest(t, repoDir, staging.Worktree)
}

func TestCarryForwardStagingSweepKeepsALiveOwnersWorktree(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	staging, err := AddCarryForwardStaging(t.Context(), repoDir, head, t.TempDir())
	if err != nil {
		t.Fatalf("add carry-forward staging: %v", err)
	}
	t.Cleanup(func() { removeStagingForTest(t, staging) })

	candidate := singleStagingCandidateForTest(t, repoDir)
	if candidate.Stale || candidate.OwnerPID != os.Getpid() {
		t.Fatalf("live-owner candidate = %+v, want preserved PID %d", candidate, os.Getpid())
	}
	if !strings.Contains(candidate.RefusalReason, "live") {
		t.Fatalf("live-owner refusal = %q, want live owner reason", candidate.RefusalReason)
	}
}

func TestCarryForwardStagingSweepReleasesALockedInitializingLegacyWorktree(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := addLegacyStagingForTest(t, repoDir, head, true)

	candidate := singleStagingCandidateForTest(t, repoDir)
	if !candidate.Stale || candidate.OwnerPID != 0 || !strings.Contains(candidate.Proof, "locked initializing") {
		t.Fatalf("locked legacy candidate = %+v, want stale locked initializing proof", candidate)
	}
	if err := ReleaseCarryForwardStaging(t.Context(), repoDir, candidate); err != nil {
		t.Fatalf("release locked legacy staging: %v", err)
	}
	assertStagingReleasedForTest(t, repoDir, worktreePath)
}

func TestCarryForwardStagingSweepKeepsAnUnlockedLegacyWorktree(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := addLegacyStagingForTest(t, repoDir, head, false)
	t.Cleanup(func() {
		removeStagingForTest(t, CarryForwardStaging{repository: repoDir, root: filepath.Dir(worktreePath), Worktree: worktreePath})
	})

	candidate := singleStagingCandidateForTest(t, repoDir)
	if candidate.Stale || candidate.OwnerPID != 0 {
		t.Fatalf("unlocked legacy candidate = %+v, want preserved legacy staging", candidate)
	}
	if !strings.Contains(candidate.RefusalReason, "not locked initializing") {
		t.Fatalf("unlocked legacy refusal = %q, want missing locked initializing proof", candidate.RefusalReason)
	}
}

func TestCarryForwardStagingOwnerHelperProcess(t *testing.T) {
	t.Parallel()
	if os.Getenv("ROUNDFIX_STAGING_OWNER_HELPER") != "1" {
		return
	}
	if _, err := io.Copy(io.Discard, os.Stdin); err != nil {
		t.Fatalf("wait for owner helper release: %v", err)
	}
}

type stagingAdminTestRunner struct {
	commonDir  string
	resolved   chan struct{}
	commandRan bool
}

func (runner *stagingAdminTestRunner) Run(_ context.Context, _ string, args ...string) (string, error) {
	if len(args) == 3 && args[0] == "rev-parse" && args[1] == "--path-format=absolute" && args[2] == "--git-common-dir" {
		runner.resolved <- struct{}{}
		return runner.commonDir + "\n", nil
	}
	if len(args) >= 2 && args[0] == "worktree" {
		runner.commandRan = true
		return "", nil
	}
	return "", errors.New("unexpected Git command")
}

func newStagingRepository(t *testing.T) (string, string) {
	t.Helper()
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "tracked\n")
	gitWorktreeTest(t, repoDir, "add", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	return repoDir, strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
}

func addLegacyStagingForTest(t *testing.T, repoDir string, head string, locked bool) string {
	t.Helper()
	root, err := os.MkdirTemp(t.TempDir(), stagingRootPrefix)
	if err != nil {
		t.Fatalf("create legacy staging root: %v", err)
	}
	worktreePath := filepath.Join(root, stagingWorktreeName)
	gitWorktreeTest(t, repoDir, "worktree", "add", "--detach", worktreePath, head)
	if locked {
		gitWorktreeTest(t, repoDir, "worktree", "lock", "--reason", stagingInitializingLockReason, worktreePath)
	}
	return worktreePath
}

func exitedOwnerForTest(t *testing.T) (int, string) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("open helper pipe: %v", err)
	}
	command := exec.Command(os.Args[0], "-test.run=^TestCarryForwardStagingOwnerHelperProcess$")
	command.Env = append(os.Environ(), "ROUNDFIX_STAGING_OWNER_HELPER=1")
	command.Stdin = reader
	if err := command.Start(); err != nil {
		_ = reader.Close()
		_ = writer.Close()
		t.Fatalf("start owner helper: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close parent helper reader: %v", err)
	}
	pid := command.Process.Pid
	identity, err := store.OwnerProcessIdentity(t.Context(), pid)
	if err != nil {
		_ = writer.Close()
		_ = command.Wait()
		t.Fatalf("read helper process identity: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("release owner helper: %v", err)
	}
	if err := command.Wait(); err != nil {
		t.Fatalf("wait for owner helper exit: %v", err)
	}
	return pid, identity
}

func singleStagingCandidateForTest(t *testing.T, repoDir string) StagingCandidate {
	t.Helper()
	candidates, err := InspectCarryForwardStaging(t.Context(), repoDir)
	if err != nil {
		t.Fatalf("inspect carry-forward staging: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("staging candidates = %+v, want one", candidates)
	}
	return candidates[0]
}

func readStagingOwnerForTest(t *testing.T, path string) stagingOwner {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read staging owner %q: %v", path, err)
	}
	var owner stagingOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		t.Fatalf("decode staging owner %q: %v", path, err)
	}
	return owner
}

func writeStagingOwnerForTest(t *testing.T, path string, owner stagingOwner) {
	t.Helper()
	data, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("encode staging owner: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatalf("write staging owner %q: %v", path, err)
	}
}

func removeStagingForTest(t *testing.T, staging CarryForwardStaging) {
	t.Helper()
	registered := gitWorktreeTest(t, staging.repository, "worktree", "list", "--porcelain")
	if strings.Contains(registered, "worktree "+staging.Worktree) {
		gitWorktreeTest(t, staging.repository, "worktree", "remove", "--force", "--force", staging.Worktree)
	}
	if err := os.RemoveAll(staging.root); err != nil {
		t.Errorf("remove staging root %q: %v", staging.root, err)
	}
}

func assertStagingReleasedForTest(t *testing.T, repoDir string, worktreePath string) {
	t.Helper()
	if registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain"); strings.Contains(registered, "worktree "+worktreePath) {
		t.Fatalf("staging worktree remains registered:\n%s", registered)
	}
	if _, err := os.Stat(filepath.Dir(worktreePath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("staging root stat error = %v, want not exist", err)
	}
}
