// Suite: Delivery Queue review lineage notices.
// Invariant: a restart is named before review, while continuing lineages keep their output.
// Boundary IN: public delivery commands, temporary Git repositories, review records, and the real queue store.
// Boundary OUT: reviewer execution and owner startup use the existing command seams.
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
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

func TestDeliverStatusNamesAReviewThatOpensANewLineage(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"rewritten head", "conflict-free merge of the default branch"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fixture, runner, prior := newDeliveryLineageNoticeFixture(t, name)
			head := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
			want := "0356-example\tparked\treview-findings\t" + fixture.repository + "\n"
			if name == "rewritten head" {
				want += fmt.Sprintf("Review: 0356-example head %s opens a new Reviewer Lineage at round 1, a full review of the candidate: the head does not descend from the reviewed head %s\n", head, prior.HeadCommit)
			}
			want += "Park: 0356-example finding: resolve the blocker, then run roundfix deliver retry 0356-example\n" +
				"Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\n" +
				"Usage: no Runs recorded\n" +
				"Pending question: 0356-example parked review-findings\n" +
				"Answer: resolve the blocker, then run roundfix deliver retry 0356-example\n"
			if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
				t.Fatalf("deliver status exit=%d stderr=%q stdout=%q, want %q", code, stderr.String(), stdout.String(), want)
			}
			if len(runner.prompts) != 1 {
				t.Fatalf("status ran the reviewer: %d prompts", len(runner.prompts))
			}
		})
	}
}

func TestDeliverRetryNamesAReviewThatOpensANewLineage(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, change string
		stage        store.DeliveryStage
		liveOwner    bool
	}{
		{"rewritten head", "rewritten head", store.DeliveryStageReviewing, false},
		{"conflict-free merge of the default branch", "conflict-free merge of the default branch", store.DeliveryStageReviewing, false},
		{"live owner", "rewritten head", store.DeliveryStageReviewing, true},
		{"returns to implementation", "rewritten head", store.DeliveryStageRunning, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture, runner, prior := newDeliveryLineageNoticeFixture(t, tc.change)
			head := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			engine := &retryCommandDeliveryEngine{retryResult: delivery.RetryResult{
				SpecSlug: "0356-example",
				Blocker:  "review-findings",
				Stage:    tc.stage,
			}}
			if tc.liveOwner {
				engine.retryResult.OwnerPID = os.Getpid()
				engine.retryResult.OwnerIdentity = "recorded-owner"
			}
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.ownerProcesses = &fakeDeliveryOwnerProcesses{}
				dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
				dependencies.startDeliveryOwner = func(_ context.Context, _ roundconfig.Loaded, _ commandEnvironment, stdout, _ io.Writer) int {
					if tc.liveOwner {
						t.Fatal("retry replaced a live owner")
					}
					if len(runner.prompts) != 1 {
						t.Fatalf("retry ran the reviewer before starting its owner: %d prompts", len(runner.prompts))
					}
					fmt.Fprintln(stdout, "Delivery Owner: test-owner")
					return exitOK
				}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"deliver", "retry", "0356-example"}, &stdout, &stderr)
			want := fmt.Sprintf("Retried 0356-example: review-findings -> %s\n", tc.stage)
			if tc.change == "rewritten head" && tc.stage == store.DeliveryStageReviewing {
				want += fmt.Sprintf("Review: 0356-example head %s opens a new Reviewer Lineage at round 1, a full review of the candidate: the head does not descend from the reviewed head %s\n", head, prior.HeadCommit)
			}
			if tc.liveOwner {
				want += fmt.Sprintf("Handed 0356-example to Delivery Queue owner PID %d.\n", os.Getpid())
			} else {
				want += "Delivery Owner: test-owner\n"
			}
			if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
				t.Fatalf("deliver retry exit=%d stderr=%q stdout=%q, want %q", code, stderr.String(), stdout.String(), want)
			}
			if engine.retryCalls != 1 || engine.gitRoot != fixture.repository || engine.specSlug != "0356-example" {
				t.Fatalf("retry calls=%d root=%q slug=%q", engine.retryCalls, engine.gitRoot, engine.specSlug)
			}
		})
	}
}

func TestReviewLineageNoticeIsAdvisory(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"codex", "claude", "none", "coderabbit", "same head", "no worktree", "not a repository", "missing artifacts", "missing record", "malformed record", "Git failure"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			change := "rewritten head"
			if name == "same head" {
				change = "same head"
			}
			fixture, runner, _ := newDeliveryLineageNoticeFixture(t, change)
			loaded, err := loadCommandConfig(commandEnvironmentForTest(t), io.Discard)
			if err != nil {
				t.Fatal(err)
			}
			worktree := fixture.repository
			var git preflight.GitRunner = preflight.ExecGitRunner{}
			recordPath := filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName)
			missingArtifacts := filepath.Join(t.TempDir(), "absent", "artifacts")
			switch name {
			case "claude", "none", "coderabbit":
				loaded.Config.PrePRReview.Provider = name
			case "no worktree":
				worktree = ""
			case "not a repository":
				worktree = t.TempDir()
			case "missing artifacts":
				loaded.Config.Defaults.ArtifactDir = missingArtifacts
			case "missing record":
				if err := os.Remove(recordPath); err != nil {
					t.Fatal(err)
				}
			case "malformed record":
				if err := os.WriteFile(recordPath, []byte("invalid JSON"), 0o600); err != nil {
					t.Fatal(err)
				}
			case "Git failure":
				head := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
				git = mainMergeLineageFailureGit{
					args: []string{"merge-base", "--is-ancestor", fixture.headCommit, head},
					err:  errors.New("Git ancestry unavailable"),
				}
			}
			before, beforeErr := os.ReadFile(recordPath)
			notice := reviewLineageNotice(t.Context(), git, loaded, "0356-example", worktree)
			wantNotice := name == "codex" || name == "claude"
			if (notice != "") != wantNotice {
				t.Fatalf("notice=%q, want notice=%t", notice, wantNotice)
			}
			after, afterErr := os.ReadFile(recordPath)
			if !bytes.Equal(before, after) || errors.Is(beforeErr, os.ErrNotExist) != errors.Is(afterErr, os.ErrNotExist) {
				t.Fatal("notice changed the review record")
			}
			if _, err := os.Stat(filepath.Dir(missingArtifacts)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("notice created an artifact directory: %v", err)
			}
			if len(runner.prompts) != 1 {
				t.Fatalf("notice ran the reviewer: %d prompts", len(runner.prompts))
			}
		})
	}
}

func newDeliveryLineageNoticeFixture(t *testing.T, change string) (reviewCommandFixture, *lineageReviewRunner, reviewRecord) {
	t.Helper()
	runner := newLineageRunner("Findings:\n- review.txt:2 Failure: candidate defect")
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, prior, stderr := runDefaultBranchReview(t)
	if code != exitRunFailed || stderr != "" || prior.Outcome != reviewOutcomeFindings || prior.Lineage == nil || prior.Lineage.Round != 1 {
		t.Fatalf("round one exit=%d record=%+v stderr=%q", code, prior, stderr)
	}
	if change == "rewritten head" {
		gittest.Run(t, fixture.repository, "commit", "--amend", "-m", "rewrite the reviewed head")
	} else if change == "conflict-free merge of the default branch" {
		mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
	}
	if err := os.MkdirAll(filepath.Join(fixture.repository, "docs", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	runStore, err := store.Open(t.Context(), fixture.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = runStore.Close() })
	queue, err := runStore.CreateDeliveryQueue(t.Context(), fixture.repository, []string{"0356-example"})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), fixture.repository, "0356-example", "feature/review", fixture.repository); err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Branch = "feature/review"
	item.Blocker = "review-findings"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), fixture.repository, item); err != nil {
		t.Fatal(err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	return fixture, runner, prior
}
