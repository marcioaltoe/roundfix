// Suite: delivery Base Update through the command engine and real Git.
// Invariant: a strict local forge never receives a merge at a behind head.
// Boundary IN: real queue, item worktree, repository gate, pushes and merge commit.
// Boundary OUT: gh is a local script and never reaches a network forge.
package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestCommandDeliveryEngineBringsABehindPullRequestUpToDate(t *testing.T) {
	t.Parallel()
	// The real constructor resolves gh through the process PATH. Run this
	// scenario in a child with its own PATH, without changing the suite's env.
	const forgeBinEnv = "ROUNDFIX_BASE_UPDATE_FORGE_BIN"
	bin := os.Getenv(forgeBinEnv)
	if bin == "" {
		bin = t.TempDir()
		command := exec.CommandContext(t.Context(), scriptFixtureBinaryPath,
			"-test.run=^TestCommandDeliveryEngineBringsABehindPullRequestUpToDate$", "-test.count=1", "-test.v")
		command.Env = withEnvValue(isolatedGitEnvForTest(), forgeBinEnv, bin)
		command.Env = withEnvValue(command.Env, "PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
		command.Env = withEnvValue(command.Env, "HOME", t.TempDir())
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("command engine with isolated forge: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "roundfix: base update: Delivery Queue item ") {
			t.Fatalf("missing Base Update log in child output: %s", output)
		}
		t.Logf("isolated forge scenario: %s", output)
		return
	}
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: "completed"}})
	setImplementFixtureAuthorizationOperations(t, repo, "implement", "commit", "push", "pull_request", "merge")
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	mustWrite(t, filepath.Join(repo, "docs/specs/.gitkeep"), "")
	commitArchiveFixture(t)
	gittest.PersistIdentity(t, repo)
	origin := t.TempDir()
	gittest.InitRepo(t, origin, "--bare", "--initial-branch=main")
	gittest.Run(t, repo, "remote", "add", "origin", origin)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "remote", "set-head", "origin", "main")
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Error(err)
		}
	})
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
	if err != nil {
		t.Fatal(err)
	}
	cfg := retirementWorkflow(repo, home).loaded
	cfg.Config.Worktree.Location = filepath.Join(t.TempDir(), "worktrees")
	cfg.Config.Defaults.Verification = "true"
	cfg.Config.Defaults.ArtifactDir = t.TempDir()
	workflow := &commandDeliveryWorkflow{store: runStore, loaded: cfg, git: retirementWorkflow(repo, home).git}
	branch, workDir, err := workflow.CreateItemBranch(commandContextForTest(t, t.Context()), repo, implementTestSlug)
	if err != nil {
		t.Fatal(err)
	}
	candidate := strings.TrimSpace(gittest.Run(t, workDir, "rev-parse", "HEAD"))
	_, err = spec.Archive(spec.ArchiveRequest{SpecsRoot: filepath.Join(workDir, "docs/specs"), BuiltInRoot: true, Slug: implementTestSlug, SourceRevision: candidate, RepositoryRoot: workDir})
	if err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, workDir, "add", "-A")
	gittest.Run(t, workDir, "commit", "-m", "docs: archive Spec")
	archiveHead := strings.TrimSpace(gittest.Run(t, workDir, "rev-parse", "HEAD"))
	item := queue.Items[0]
	item.Branch, item.Worktree, item.WorktreeProvisioned = branch, workDir, true
	item.Stage = store.DeliveryStagePublishing
	item.CandidateCommits = []string{candidate, archiveHead}
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	// Move main only after the item branch was cut, using a separate clone.
	movedClone := t.TempDir()
	gittest.Run(t, movedClone, "clone", origin, ".")
	mustWrite(t, filepath.Join(movedClone, "moved.txt"), "The base branch moved.\n")
	gittest.Run(t, movedClone, "add", "moved.txt")
	gittest.Run(t, movedClone, "commit", "-m", "docs: move the base branch")
	gittest.Run(t, movedClone, "push", "origin", "main")
	movedMain := strings.TrimSpace(gittest.Run(t, movedClone, "rev-parse", "HEAD"))
	// gh is a strict local forge: publication uses real git pushes, merge builds
	// a squash commit in the bare origin without pushing any additional ref.
	stateDir := t.TempDir()
	script := fmt.Sprintf(`#!/bin/sh
set -eu
state=%s
origin=%s
branch=%s
head=$(git --git-dir="$origin" rev-parse "refs/heads/$branch")
case "$1 $2" in
 'api '*) printf '{"total_count":0,"workflow_runs":[]}'; exit 0 ;;
 'pr create') touch "$state/open"; printf 'https://local.invalid/pull/7\n'; exit 0 ;;
 'pr checks') printf '[{"name":"verify","bucket":"pass","state":"SUCCESS"}]'; exit 0 ;;
 'pr merge')
  test "$3" = 7
  test "$4" = --squash
  test "$5" = --match-head-commit
  test "$6" = "$head"
  if ! git --git-dir="$origin" merge-base --is-ancestor main "$head"; then
   touch "$state/refused"
   printf 'X Pull request #7 is not mergeable: the base branch policy prohibits the merge.\n' >&2
   exit 1
  fi
  tree=$(git --git-dir="$origin" rev-parse "$head^{tree}")
  base=$(git --git-dir="$origin" rev-parse main)
  merge=$(git --git-dir="$origin" -c user.name=Forge -c user.email=forge@example.com commit-tree "$tree" -p "$base" -m 'feat: deliver Spec (#7)')
  git --git-dir="$origin" update-ref refs/heads/main "$merge"
  printf '%%s' "$merge" > "$state/merged"
  exit 0 ;;
 'pr list'|'pr view') ;;
 *) printf 'unexpected forge command: %%s\n' "$*" >&2; exit 2 ;;
esac
if test "$2" = list && ! test -f "$state/open"; then printf '[]'; exit 0; fi
merge_state=CLEAN
if ! git --git-dir="$origin" merge-base --is-ancestor main "$head"; then merge_state=BEHIND; fi
status=OPEN
merge='null'
if test -f "$state/merged"; then status=MERGED; merge="{\"oid\":\"$(cat "$state/merged")\"}"; fi
payload=$(printf '{"number":7,"state":"%%s","headRefName":"%%s","headRefOid":"%%s","mergeable":"MERGEABLE","mergeStateStatus":"%%s","mergeCommit":%%s}' "$status" "$branch" "$head" "$merge_state" "$merge")
if test "$2" = list; then printf '[%%s]' "$payload"; else printf '%%s' "$payload"; fi
`, shellQuoteImplement(stateDir), shellQuoteImplement(origin), shellQuoteImplement(branch))
	gh := filepath.Join(bin, "gh")
	writeScriptFixture(t, gh, script)
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	engine := newCommandDeliveryEngine(runStore, cfg)
	if _, err := engine.Run(commandContextForTest(t, t.Context()), repo); err != nil {
		t.Fatal(err)
	}
	item = readDeliveryItemForCLI(t, runStore, repo)
	if item.Stage != store.DeliveryStageMerged || item.PullRequestNumber != "7" {
		t.Fatalf("item=%+v", item)
	}
	if len(item.CandidateCommits) != 4 {
		t.Fatalf("candidates=%v", item.CandidateCommits)
	}
	updated := item.CandidateCommits[len(item.CandidateCommits)-1]
	parents := strings.Fields(gittest.Run(t, origin, "show", "-s", "--format=%P", updated))
	if len(parents) != 2 || parents[0] != item.CandidateCommits[len(item.CandidateCommits)-2] || parents[1] != movedMain {
		t.Fatalf("Base Update parents=%v, moved main=%s", parents, movedMain)
	}
	message := gittest.Run(t, origin, "show", "-s", "--format=%B", updated)
	if !strings.Contains(message, "Roundfix-Delivery: derived-merge") {
		t.Fatalf("Base Update message=%q", message)
	}
	if content := gittest.Run(t, origin, "show", item.MergeCommit+":moved.txt"); content != "The base branch moved.\n" {
		t.Fatalf("merged base file=%q", content)
	}
	if _, err := os.Stat(filepath.Join(stateDir, "refused")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected policy refusal marker: %v", err)
	}
}
