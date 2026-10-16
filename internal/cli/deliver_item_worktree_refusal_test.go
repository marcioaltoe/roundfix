// Suite: Delivery Queue refusals in recorded item worktrees.
// Invariant: queue commands name the owning checkout without acting on its queue.
// Boundary IN: public CLI dispatch, real linked Git worktrees, and the SQLite Run Database.
// Boundary OUT: Delivery Engine construction and detached owner launch are counted stubs.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const itemWorktreeRefusalSlug = "0300-example"

func TestDeliverRetryInAnItemWorktreeReproducesSurfaceTranscript(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, false)
	fixture.refuse(t, "retry", fixture.worktree, fixture.reason(fixture.worktree))
}

func TestDeliverStatusInAnItemWorktreeReproducesSurfaceTranscript(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, false)
	fixture.refuse(t, "status", fixture.worktree, fixture.reason(fixture.worktree))
}

func TestDeliverResumeAndStopInAnItemWorktreeNameTheOwningCheckout(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, false)
	for _, command := range []string{"resume", "stop"} {
		t.Run(command, func(t *testing.T) {
			fixture.refuse(t, command, fixture.worktree, fixture.reason(fixture.worktree))
		})
	}
}

func TestDeliverItemWorktreeRefusalComparesResolvedPaths(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, true)
	for _, command := range []string{"retry", "status", "resume", "stop"} {
		t.Run(command, func(t *testing.T) {
			fixture.refuse(t, command, fixture.worktree, fixture.reason(fixture.worktree))
		})
	}
}

func TestDeliverMissingQueueRefusalKeepsItsTextElsewhere(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, false)
	_, unrelated := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	for _, command := range []string{"retry", "status", "resume", "stop"} {
		t.Run(command, func(t *testing.T) {
			fixture.refuse(t, command, unrelated, fmt.Sprintf("Delivery Queue for repository %q does not exist", unrelated))
		})
	}
	t.Run("lookup unavailable", func(t *testing.T) {
		runStore, err := store.OpenReader(t.Context(), fixture.home)
		if err != nil {
			t.Fatal(err)
		}
		if err := runStore.Close(); err != nil {
			t.Fatal(err)
		}
		got := missingDeliveryQueue(t.Context(), runStore, fixture.worktree)
		want := fmt.Sprintf("Delivery Queue for repository %q does not exist", fixture.worktree)
		if got.Error() != want {
			t.Fatalf("unavailable lookup reason=%q, want %q", got, want)
		}
	})
}

func TestDeliverItemWorktreeRefusalUsesFirstRecordedMatch(t *testing.T) {
	t.Parallel()
	fixture := newDeliveryItemWorktreeRefusalFixture(t, false)
	// A parent path sorts before its child regardless of temporary directory names.
	firstRoot := filepath.Dir(fixture.repo)
	runStore, err := store.Open(t.Context(), fixture.home)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runStore.CreateDeliveryQueue(t.Context(), firstRoot, []string{"z-first", "a-second"}); err != nil {
		t.Fatal(err)
	}
	for _, slug := range []string{"a-second", "z-first"} {
		if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), firstRoot, slug, "test/"+slug, fixture.worktree); err != nil {
			t.Fatal(err)
		}
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	reason := fmt.Sprintf("Delivery Queue for repository %q does not exist; this checkout is the item worktree of z-first in the Delivery Queue of %q; run the command in %q", fixture.worktree, firstRoot, firstRoot)
	fixture.refuse(t, "status", fixture.worktree, reason)
}

type deliveryItemWorktreeRefusalFixture struct {
	home, repo, worktree string
	queue                store.DeliveryQueue
	engineBuilds, starts int
}

func newDeliveryItemWorktreeRefusalFixture(t *testing.T, linkedParent bool) *deliveryItemWorktreeRefusalFixture {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	worktree := filepath.Join(t.TempDir(), "item")
	gitImplement(t, repo, "worktree", "add", "-b", "test/item-refusal", worktree)
	resolved, err := filepath.EvalSymlinks(worktree)
	if err != nil {
		t.Fatal(err)
	}
	recorded := resolved
	if linkedParent {
		// Only the parent is a link; the item directory itself is real (ADR-0268).
		parentLink := filepath.Join(t.TempDir(), "worktree-root")
		if err := os.Symlink(filepath.Dir(resolved), parentLink); err != nil {
			t.Fatal(err)
		}
		recorded = filepath.Join(parentLink, filepath.Base(resolved))
	}
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{itemWorktreeRefusalSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage, item.Blocker = store.DeliveryStageParked, "run-unresolved"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), repo, itemWorktreeRefusalSlug, "test/item-refusal", recorded); err != nil {
		t.Fatal(err)
	}
	if err := runStore.ClaimDeliveryQueueOwner(t.Context(), repo, os.Getpid(), "recorded-owner"); err != nil {
		t.Fatal(err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	fixture := &deliveryItemWorktreeRefusalFixture{home: home, repo: repo, worktree: resolved, queue: openDeliveryQueueForCLI(t, home, repo)}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine {
			fixture.engineBuilds++
			return &retryCommandDeliveryEngine{}
		}
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			fixture.starts++
			return exitOK
		}
		dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{}
	})
	return fixture
}

func (fixture *deliveryItemWorktreeRefusalFixture) reason(worktree string) string {
	return fmt.Sprintf("Delivery Queue for repository %q does not exist; this checkout is the item worktree of %s in the Delivery Queue of %q; run the command in %q", worktree, itemWorktreeRefusalSlug, fixture.repo, fixture.repo)
}

func (fixture *deliveryItemWorktreeRefusalFixture) refuse(t *testing.T, command, worktree, reason string) {
	t.Helper()
	setCommandEnvironmentForTest(t, fixture.home, worktree)
	args := []string{"deliver", command}
	var want string
	if command == "retry" {
		args = append(args, itemWorktreeRefusalSlug)
		want = fmt.Sprintf("Retry refused\n\nReason:\n  retry Delivery Queue item %q: %s\n\nNo side effects:\n  Roundfix did not change the Delivery Queue item, start a queue owner, commit, or push.\n", itemWorktreeRefusalSlug, reason)
	} else {
		want = "Preflight failed\n\nReason:\n  " + reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix deliver " + command + " --help' for usage.\n"
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, args, &stdout, &stderr)
	if code != exitPreflight || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q, want exit=2 stdout=empty stderr=%q", code, stdout.String(), stderr.String(), want)
	}
	if fixture.starts != 0 || fixture.engineBuilds != 0 {
		t.Fatalf("refusal built %d engines and started %d owners", fixture.engineBuilds, fixture.starts)
	}
	if got := openDeliveryQueueForCLI(t, fixture.home, fixture.repo); !reflect.DeepEqual(got, fixture.queue) {
		t.Fatalf("owning queue changed: got=%+v want=%+v", got, fixture.queue)
	}
}
