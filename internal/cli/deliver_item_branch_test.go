// Boundary: public deliver start, real local Git refs/worktrees and SQLite.
// Only machine readiness and detached owner launch are injected; origin is a
// temporary bare repository and every home/worktree is disposable.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func newItemBranchStartRepository(t *testing.T) (home, repo string) {
	t.Helper()
	root := t.TempDir()
	remote := filepath.Join(root, "origin.git")
	seed := filepath.Join(root, "seed")
	repo = filepath.Join(root, "clone")
	home = filepath.Join(root, "home")
	if err := os.MkdirAll(home, 0755); err != nil {
		t.Fatal(err)
	}
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.InitRepo(t, seed, "-b", "main")
	gittest.PersistIdentity(t, seed)
	writeImplementSpec(t, seed, implementTestSlug, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, seed, "implement", "commit", "push", "pull_request", "merge")
	mustWrite(t, filepath.Join(seed, ".roundfixrc.yml"), "worktree:\n  location: "+filepath.Join(root, "worktrees")+"\n")
	gittest.Run(t, seed, "add", ".")
	gittest.Run(t, seed, "commit", "-m", "seed")
	gittest.Run(t, seed, "remote", "add", "origin", remote)
	gittest.Run(t, seed, "push", "origin", "main")
	gittest.Run(t, root, "clone", remote, repo)
	gittest.Harden(t, repo)
	gittest.PersistIdentity(t, repo)
	resolved, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	repo = resolved
	setCommandEnvironmentForTest(t, home, repo)
	return home, repo
}

func itemBranchWithWork(t *testing.T, repo, suffix string, commits int) string {
	t.Helper()
	branch := deliveryBranchPrefix + implementTestSlug + "-" + suffix
	gittest.Run(t, repo, "checkout", "-b", branch, "origin/main")
	for i := 0; i < commits; i++ {
		mustWrite(t, filepath.Join(repo, "work.txt"), fmt.Sprintf("work %d\n", i))
		gittest.Run(t, repo, "add", "work.txt")
		gittest.Run(t, repo, "commit", "-m", "item work")
	}
	gittest.Run(t, repo, "checkout", "main")
	return branch
}

func TestDeliverStartContinuesTheItemBranchThatHoldsWork(t *testing.T) {
	for _, hasWorktree := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing worktree %v", hasWorktree), func(t *testing.T) {
			home, repo := newItemBranchStartRepository(t)
			branch := itemBranchWithWork(t, repo, "1111111111111111", 2)
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", branch))
			// Non-item and other-slug branches with work must not create ambiguity.
			gittest.Run(t, repo, "branch", deliveryBranchPrefix+implementTestSlug+"-legacy", branch)
			gittest.Run(t, repo, "branch", deliveryBranchPrefix+"other-2222222222222222", branch)
			var expectedWorktree string
			var out, diag bytes.Buffer
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.deliveryReadiness = readyDeliveryReadiness
				deps.startDeliveryOwner = func(ctx context.Context, loaded roundconfig.Loaded, _ commandEnvironment, _, _ io.Writer) int {
					if !strings.Contains(out.String(), "Continuing item branch "+branch+" for "+implementTestSlug) {
						t.Fatal("owner started before continuation output")
					}
					db, err := store.Open(ctx, home)
					if err != nil {
						t.Fatal(err)
					}
					defer db.Close()
					workflow := &commandDeliveryWorkflow{store: db, loaded: loaded, git: preflight.ExecGitRunner{}}
					ref, err := runworktree.ItemRefFor(repo, loaded.Config.Worktree.Location, branch)
					if err != nil {
						t.Fatal(err)
					}
					expectedWorktree = ref.Path
					if hasWorktree {
						gittest.Run(t, repo, "worktree", "add", ref.Path, branch)
						mustWrite(t, filepath.Join(ref.Path, "local-note.txt"), "keep me\n")
					}
					got, path, err := workflow.CreateItemBranch(ctx, repo, implementTestSlug)
					if err != nil || got != branch || path != expectedWorktree {
						t.Fatalf("branch=%q path=%q err=%v", got, path, err)
					}
					if strings.TrimSpace(gittest.Run(t, path, "rev-parse", "HEAD")) != head || mustRead(t, filepath.Join(path, "work.txt")) != "work 1\n" {
						t.Fatal("continued work was lost")
					}
					if hasWorktree && mustRead(t, filepath.Join(path, "local-note.txt")) != "keep me\n" {
						t.Fatal("worktree note was lost")
					}
					return exitOK
				}
			})
			if code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag); code != exitOK || diag.Len() != 0 {
				t.Fatalf("exit=%d stderr=%q", code, diag.String())
			}
			item := openDeliveryQueueForCLI(t, home, repo).Items[0]
			if item.Branch != branch || item.Worktree != expectedWorktree || item.Stage != store.DeliveryStageQueued {
				t.Fatalf("item=%+v", item)
			}
		})
	}
}

func TestDeliverStartRefusesTwoItemBranchesWithWork(t *testing.T) {
	home, repo := newItemBranchStartRepository(t)
	second := itemBranchWithWork(t, repo, "2222222222222222", 1)
	first := itemBranchWithWork(t, repo, "1111111111111111", 2)
	// Prove branch discovery does not fetch or reach readiness/owner launch.
	gittest.Run(t, repo, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing"))
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = func(context.Context, roundconfig.Loaded) []CheckResult { t.Fatal("readiness reached"); return nil }
		deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			t.Fatal("owner reached")
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag)
	reason := fmt.Sprintf("Spec %q has 2 item branches with commits origin/main lacks: %s, %s; delete every branch but the one to continue, then run roundfix deliver start again", implementTestSlug, first, second)
	if code != exitPreflight || out.Len() != 0 || !strings.Contains(diag.String(), reason) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diag.String())
	}
	if _, err := os.Stat(store.DatabasePath(home)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("database created: %v", err)
	}
}

func TestDeliverStartIgnoresAnItemBranchWithoutWork(t *testing.T) {
	home, repo := newItemBranchStartRepository(t)
	old := itemBranchWithWork(t, repo, "1111111111111111", 0)
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = readyDeliveryReadiness
		deps.startDeliveryOwner = func(ctx context.Context, loaded roundconfig.Loaded, _ commandEnvironment, _, _ io.Writer) int {
			db, err := store.Open(ctx, home)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			workflow := &commandDeliveryWorkflow{store: db, loaded: loaded, git: preflight.ExecGitRunner{}}
			branch, path, err := workflow.CreateItemBranch(ctx, repo, implementTestSlug)
			if err != nil || branch == old || !strings.HasPrefix(branch, deliveryBranchPrefix+implementTestSlug+"-") {
				t.Fatalf("branch=%q err=%v", branch, err)
			}
			if strings.TrimSpace(gittest.Run(t, path, "rev-parse", "HEAD")) != strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "origin/main")) {
				t.Fatal("new branch did not start at default")
			}
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	if code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag); code != exitOK || diag.Len() != 0 || strings.Contains(out.String(), "Continuing item branch") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diag.String())
	}
	if item := openDeliveryQueueForCLI(t, home, repo).Items[0]; item.Branch == old || item.Branch == "" {
		t.Fatalf("item=%+v", item)
	}
}
