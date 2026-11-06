// Suite: numbered-file Delivery Revalidation and public collision park status.
// Boundary IN: disposable Git clones, linked item worktrees, Spec plans and a real temporary Run Database.
// Boundary OUT: only an on-disk bare origin is fetched; no Run or GitHub action starts.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

const numberedAlphaBlocker = "numbered-file-collision: db/0002_beta.sql collides with db/0002_alpha.sql (alpha); next free number: db/0003"

func TestRevalidateReportsANumberedFileCollision(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"off", "ahead_branch", "ahead_plan", "default_branch", "ignored", "renumbered", "subject_branch", "recorded_plan", "missing_plan", "missing_branch", "absent_item", "base_file", "cleared_after_fetch"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home, source := newSpecCheckWorkspace(t, "clean")
			config := "delivery:\n  numbered_files: [db/*.sql]\n"
			if name == "off" {
				config = "{}\n"
			}
			mustWrite(t, filepath.Join(source, ".roundfixrc.yml"), config)
			if err := os.MkdirAll(filepath.Join(source, "db"), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWrite(t, filepath.Join(source, "db", "0001_init.sql"), "-- init\n")
			alphaCreates := []string(nil)
			if name == "ahead_plan" || name == "missing_branch" || name == "absent_item" {
				alphaCreates = []string{"db/0002_alpha.sql"}
			}
			betaCreates := []string{"db/0002_beta.sql"}
			if name == "renumbered" {
				betaCreates = []string{"db/0003_beta.sql"}
			}
			if name == "base_file" {
				betaCreates = []string{"db/0001_init.sql"}
			}
			if name == "subject_branch" {
				betaCreates = nil
			}
			if name == "ignored" {
				betaCreates = append(betaCreates, "db/meta/_journal.json")
			}
			writeNumberedSpec(t, source, "alpha", alphaCreates, nil)
			writeNumberedSpec(t, source, "beta", betaCreates, nil)
			if name == "ignored" {
				writeNumberedSpec(t, source, "beta", betaCreates, []string{"db/0002_done.sql"})
				writeNumberedSpec(t, source, "gamma", []string{"db/0002_gamma.sql"}, nil)
			}
			gittest.Run(t, source, "add", "-A")
			gittest.Run(t, source, "commit", "-m", "test: seed numbered-file declarations")
			gittest.Run(t, source, "switch", "main")
			gittest.Run(t, source, "merge", "--ff-only", "ma/spec-check")
			origin := filepath.Join(t.TempDir(), "origin.git")
			gittest.InitRepo(t, origin, "--bare", "--initial-branch=main")
			gittest.Run(t, source, "remote", "add", "origin", origin)
			gittest.Run(t, source, "push", "-u", "origin", "main")
			queueRepo := filepath.Join(t.TempDir(), "queue")
			gittest.Run(t, "", "clone", origin, queueRepo)
			gittest.Harden(t, queueRepo)
			workDir := filepath.Join(t.TempDir(), "beta")
			gittest.Run(t, queueRepo, "worktree", "add", "-b", "item-beta", workDir, "main")
			runStore, err := store.Open(t.Context(), home)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := runStore.Close(); err != nil {
					t.Error(err)
				}
			})
			slugs := []string{"alpha", "beta"}
			if name == "ignored" {
				slugs = append(slugs, "gamma")
			}
			queue, err := runStore.CreateDeliveryQueue(t.Context(), queueRepo, slugs)
			if err != nil {
				t.Fatal(err)
			}
			if name == "ahead_branch" || name == "renumbered" || name == "ignored" || name == "subject_branch" || name == "recorded_plan" || name == "missing_plan" {
				alphaDir := filepath.Join(t.TempDir(), "alpha")
				gittest.Run(t, queueRepo, "worktree", "add", "-b", "item-alpha", alphaDir, "main")
				if name == "recorded_plan" {
					writeNumberedSpec(t, alphaDir, "alpha", []string{"db/0002_alpha.sql"}, nil)
				} else {
					mustWrite(t, filepath.Join(alphaDir, "db", "0002_alpha.sql"), "-- alpha\n")
				}
				if name == "missing_plan" {
					if err := os.Remove(filepath.Join(alphaDir, "docs", "specs", "alpha", "_prd.md")); err != nil {
						t.Fatal(err)
					}
				}
				gittest.Run(t, alphaDir, "add", "-A")
				gittest.Run(t, alphaDir, "commit", "-m", "test: add earlier branch migration")
				item := queue.Items[0]
				item.Stage = store.DeliveryStageParked
				if name == "ignored" {
					item.Stage = store.DeliveryStageMerged
				}
				item.Branch, item.Worktree = "item-alpha", alphaDir
				if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), queueRepo, "alpha", item.Branch, item.Worktree); err != nil {
					t.Fatal(err)
				}
				if err := runStore.UpdateDeliveryQueueItem(t.Context(), queueRepo, item); err != nil {
					t.Fatal(err)
				}
			}

			if name == "missing_branch" {
				item := queue.Items[0]
				item.Branch = "missing-alpha"
				if err := runStore.UpdateDeliveryQueueItem(t.Context(), queueRepo, item); err != nil {
					t.Fatal(err)
				}
			}
			if name == "absent_item" {
				// Revalidating gamma must not treat alpha or beta as earlier
				// when this queue contains only those two items.
				writeNumberedSpec(t, workDir, "gamma", []string{"db/0002_gamma.sql"}, nil)
			}
			if name == "subject_branch" {
				mustWrite(t, filepath.Join(workDir, "db", "0002_beta.sql"), "-- beta\n")
				gittest.Run(t, workDir, "add", "-A")
				gittest.Run(t, workDir, "commit", "-m", "test: subject branch claims a migration")
			}
			if name == "ahead_branch" {
				// The item's config cannot disable the declaration on main.
				mustWrite(t, filepath.Join(workDir, ".roundfixrc.yml"), "{}\n")
				gittest.Run(t, workDir, "add", "-A")
				gittest.Run(t, workDir, "commit", "-m", "test: item config does not own declarations")
			}
			if name == "cleared_after_fetch" {
				mustWrite(t, filepath.Join(source, ".roundfixrc.yml"), "delivery:\n  numbered_files: []\n")
				gittest.Run(t, source, "add", "-A")
				gittest.Run(t, source, "commit", "-m", "test: clear declaration on remote")
				gittest.Run(t, source, "push", "origin", "main")
			}
			if name == "default_branch" {
				mustWrite(t, filepath.Join(source, "db", "0002_main.sql"), "-- main\n")
				gittest.Run(t, source, "add", "-A")
				gittest.Run(t, source, "commit", "-m", "test: advance remote after item worktree creation")
				gittest.Run(t, source, "push", "origin", "main")
			}
			workflow := &commandDeliveryWorkflow{
				loaded: roundconfig.Loaded{GitRoot: queueRepo, Config: roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}}},
				git:    preflight.ExecGitRunner{}, store: runStore,
			}
			runner := &numberedFileGitLog{}
			workflow.git = runner
			if name == "off" || name == "cleared_after_fetch" {
				workflow.store = nil
			}
			beforeRefs := gittest.Run(t, queueRepo, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads", "refs/tags")
			beforeQueue, _, err := runStore.DeliveryQueue(t.Context(), queueRepo)
			if err != nil {
				t.Fatal(err)
			}
			slug := "beta"
			if name == "absent_item" {
				slug = "gamma"
			}
			result, err := workflow.Revalidate(t.Context(), workDir, slug, nil)
			if err != nil {
				t.Fatalf("revalidate: %v", err)
			}
			want := ""
			switch name {
			case "ahead_branch", "ahead_plan", "subject_branch", "recorded_plan", "missing_plan", "missing_branch":
				want = numberedAlphaBlocker
			case "default_branch":
				want = "numbered-file-collision: db/0002_beta.sql collides with db/0002_main.sql (default branch); next free number: db/0003"
			}
			if got := result.NumberedFiles.Blocker(); got != want {
				t.Fatalf("blocker = %q, want %q", got, want)
			}
			if want == "" && !reflect.DeepEqual(result.NumberedFiles, delivery.NumberedFileReport{}) {
				t.Fatalf("report = %+v, want empty", result.NumberedFiles)
			}

			if name == "off" || name == "cleared_after_fetch" {
				last := runner.commands[len(runner.commands)-1]
				if last[0] != "show" || !strings.HasSuffix(last[1], ":.roundfixrc.yml") {
					t.Fatalf("Git continued after empty declaration: %v", runner.commands)
				}
			}
			if name == "absent_item" {
				// Remove only the intentionally untracked test plan before
				// checking that revalidation itself left the tree clean.
				if err := os.RemoveAll(filepath.Join(workDir, "docs", "specs", "gamma")); err != nil {
					t.Fatal(err)
				}
			}
			afterQueue, _, err := runStore.DeliveryQueue(t.Context(), queueRepo)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(beforeQueue, afterQueue) {
				t.Fatalf("revalidation changed queue: before=%+v after=%+v", beforeQueue, afterQueue)
			}
			if after := gittest.Run(t, queueRepo, "for-each-ref", "--format=%(refname) %(objectname)", "refs/heads", "refs/tags"); after != beforeRefs {
				t.Fatalf("revalidation changed local refs: %s", after)
			}
			if dirty := gittest.Run(t, workDir, "status", "--porcelain"); strings.TrimSpace(dirty) != "" {
				t.Fatalf("revalidation changed worktree: %s", dirty)
			}
			if name == "ignored" {
				// Isolate completed claims against a live default file. The combined
				// merged/later case above alone cannot expose counting completed Tasks.
				t.Run("completed", func(t *testing.T) {
					mustWrite(t, filepath.Join(source, "db", "0002_main.sql"), "-- main\n")
					gittest.Run(t, source, "add", "-A")
					gittest.Run(t, source, "commit", "-m", "test: expose completed claim")
					gittest.Run(t, source, "push", "origin", "main")
					writeNumberedSpec(t, workDir, "beta", []string{"db/0003_beta.sql"}, []string{"db/0002_done.sql"})
					result, err := workflow.Revalidate(t.Context(), workDir, "beta", nil)
					if err != nil {
						t.Fatal(err)
					}
					if got := result.NumberedFiles.Blocker(); got != "" {
						t.Fatalf("completed Task claims a number: %s", got)
					}
				})
			}
		})
	}
}

func TestDeliverStatusNamesTheNumberedFileCollisionPark(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{"alpha", "beta"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[1]
	item.Stage, item.Blocker = store.DeliveryStageParked, numberedAlphaBlocker
	item.Branch, item.Worktree = "item-beta", "/worktrees/beta"
	if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), repo, "beta", item.Branch, item.Worktree); err != nil {
		t.Fatal(err)
	}
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	next := "in /worktrees/beta, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry beta"
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, line := range []string{"Park: beta conflict: " + next, "Answer: " + next} {
		if !strings.Contains("\n"+stdout.String(), "\n"+line+"\n") {
			t.Errorf("stdout missing line %q: %s", line, stdout.String())
		}
	}
}

func writeNumberedSpec(t *testing.T, repo, slug string, pending, completed []string) {
	t.Helper()
	dir := filepath.Join(repo, "docs", "specs", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, "_prd.md"), fmt.Sprintf("---\nspec: %s\nstatus: active\n---\n\n# Numbered files\n", slug))
	manifest := fmt.Sprintf("---\nschema: spec-tasks/v1\nspec: %s\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n", slug)
	if len(completed) != 0 {
		manifest += "    - id: task_02\n      file: task_02.md\n      needs: []\n"
	}
	mustWrite(t, filepath.Join(dir, "_tasks.md"), manifest+"---\n\n# Task Graph\n")
	writeTask := func(id, status string, paths []string) {
		content := fmt.Sprintf("---\ntask: %s\nspec: %s\nstatus: %s\ntype: backend\ncomplexity: low\n---\n\n# Task: Migration\n\n## Context\n\n", id, slug, status)
		for _, path := range paths {
			content += "- creates: `" + path + "`\n"
		}
		content += "\n## Verification\n\n- `true` — expected: exit 0.\n"
		mustWrite(t, filepath.Join(dir, id+".md"), content)
	}
	writeTask("task_01", "pending", pending)
	if len(completed) != 0 {
		writeTask("task_02", "completed", completed)
	}
}

// Each instance belongs to one subtest; revalidation invokes Git serially.
type numberedFileGitLog struct{ commands [][]string }

func (runner *numberedFileGitLog) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	runner.commands = append(runner.commands, append([]string(nil), args...))
	return (preflight.ExecGitRunner{}).RunGit(ctx, workDir, args...)
}
