package cli

// Boundary: deliver start/status, real Git archives and Run Database; only
// readiness, owner launch and the implementation/review/forge boundary are fake.
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type archivedResumeForge struct {
	t                *testing.T
	mergedHead       string
	specsRootMissing bool
	parkTestDeliveryFlow
	reviews []string
	created int
}

func (forge *archivedResumeForge) Review(_ context.Context, _, _, head string) (delivery.ReviewResult, error) {
	forge.reviews = append(forge.reviews, head)
	return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeReviewed, Head: head}, nil
}

func (forge *archivedResumeForge) FindOrCreatePullRequest(_ context.Context, request delivery.PullRequestRequest) (delivery.PullRequestResult, error) {
	for _, pr := range forge.pullRequests {
		if pr.State == "OPEN" && pr.HeadBranch == request.HeadBranch {
			return delivery.PullRequestResult{PullRequest: pr}, nil
		}
	}
	forge.created++
	return delivery.PullRequestResult{}, fmt.Errorf("unexpected creation for %s", request.HeadBranch)
}

func (forge *archivedResumeForge) MergePullRequest(ctx context.Context, number, expectedHead string) (delivery.MergeResult, error) {
	// Observe the unchanged item head before the real workspace retires the
	// merged branch and worktree. This also checks the absent-root shape.
	forge.mergedHead = strings.TrimSpace(gittest.Run(forge.t, forge.workDir, "rev-parse", "HEAD"))
	_, err := os.Stat(filepath.Join(forge.workDir, "docs/specs"))
	forge.specsRootMissing = errors.Is(err, os.ErrNotExist)
	return forge.parkTestDeliveryFlow.MergePullRequest(ctx, number, expectedHead)
}

func (forge *archivedResumeForge) WithWorkDir(workDir string) delivery.PullRequestBoundary {
	forge.parkTestDeliveryFlow.WithWorkDir(workDir)
	return forge
}

// A temporary home, bare origin and hardened clone use the same setup as the
// continuation tests. No real tool state or forge is consulted.
func archivedResumeBranch(t *testing.T, anotherSpec bool, mutation string) (home, repo, branch, head, archiveHead string) {
	t.Helper()
	home, repo = newItemBranchStartRepository(t)
	branch = deliveryBranchPrefix + implementTestSlug + "-1111111111111111"
	gittest.Run(t, repo, "checkout", "-b", branch)
	if anotherSpec {
		writeImplementSpec(t, repo, "0002-other", []implementSeed{{id: "task_01"}})
	}
	writeImplementSpec(t, repo, implementTestSlug, []implementSeed{{id: "task_01", status: "completed"}})
	mustWrite(t, filepath.Join(repo, "docs/specs", implementTestSlug, "_authorization.md"), implementFixtureAuthorization(implementTestSlug, "implement", "commit", "push", "pull_request", "merge"))
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "test: settle item")
	sourceHead := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	if mutation == "folder_removed_without_archive" {
		gittest.Run(t, repo, "rm", "-r", "docs/specs/"+implementTestSlug)
	} else {
		archived, err := spec.Archive(spec.ArchiveRequest{SpecsRoot: filepath.Join(repo, "docs/specs"), BuiltInRoot: true, Slug: implementTestSlug, SourceRevision: sourceHead, RepositoryRoot: repo})
		if err != nil {
			t.Fatal(err)
		}
		if mutation == "inexact_archive_commit" {
			mustWrite(t, filepath.Join(repo, "unrelated.txt"), "unrelated\n")
		}
		gittest.Run(t, repo, "add", "-A")
		// Deliberately no archive wording: the proof must judge the change.
		gittest.Run(t, repo, "commit", "-m", "docs: retire finished work")
		archiveHead = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
		changed, err := spec.RecordArchivePullRequest(archived.RecordPath, "619")
		if err != nil || !changed {
			t.Fatalf("record PR changed=%v err=%v", changed, err)
		}
		gittest.Run(t, repo, "add", "-A")
	}
	gittest.Run(t, repo, "commit", "-m", "docs: record disposition")
	if mutation == "no_removing_commit" {
		gittest.Run(t, repo, "checkout", "--orphan", "test/orphan")
		gittest.Run(t, repo, "commit", "-m", "test: archived tree without history")
		gittest.Run(t, repo, "branch", "-f", branch, "HEAD")
	}
	head = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "main")
	if mutation == "no_removing_commit" {
		gittest.Run(t, repo, "branch", "-D", "test/orphan")
	}
	return
}

func startArchivedResumeQueue(t *testing.T, home, repo string, forge *archivedResumeForge) store.DeliveryQueueItem {
	t.Helper()
	// Defense in depth: should a regression reach a roundfix step, its child
	// dispatches as the CLI helper instead of recursively running the suite.
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	forge.t = t
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = readyDeliveryReadiness
		deps.startDeliveryOwner = func(ctx context.Context, loaded roundconfig.Loaded, _ commandEnvironment, _, _ io.Writer) int {
			db, err := store.Open(ctx, home)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			workflow := &commandDeliveryWorkflow{store: db, loaded: loaded, git: preflight.ExecGitRunner{}}
			dependencies := commandDeliveryEngineDependencies(db, loaded)
			if dependencies.ArchivedItems == nil {
				t.Fatal("command engine lacks archive prover")
			}
			engine := delivery.NewEngine(db, delivery.EngineDependencies{
				Workspace: workflow, Recovery: workflow, History: workflow,
				Revalidator: workflow, ArchivedItems: workflow, Archiver: workflow, Recorder: workflow,
				Runner: forge, Reviewer: forge, Gate: forge, Authorizer: forge,
				Publication: forge, PullRequests: forge,
			})
			if _, err := engine.Run(ctx, repo); err != nil {
				t.Fatal(err)
			}
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	if code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag); code != exitOK || diag.Len() != 0 {
		t.Fatalf("start exit=%d stdout=%s stderr=%s", code, &out, &diag)
	}
	return openDeliveryQueueForCLI(t, home, repo).Items[0]
}

func TestDeliverStartResumesAnArchivedItemBranchAtReview(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"another_active_Spec", "only_active_Spec"} {
		t.Run(name, func(t *testing.T) {
			home, repo, branch, head, _ := archivedResumeBranch(t, name == "another_active_Spec", "")
			forge := &archivedResumeForge{parkTestDeliveryFlow: parkTestDeliveryFlow{
				remoteHeads:  map[string]string{branch: head},
				pullRequests: map[string]delivery.PullRequest{"619": {Number: "619", State: "OPEN", HeadBranch: branch, HeadSHA: head}},
			}}
			item := startArchivedResumeQueue(t, home, repo, forge)
			if item.Stage != store.DeliveryStageMerged || item.PullRequestNumber != "619" || item.Branch != branch {
				t.Fatalf("item=%+v", item)
			}
			if forge.runCalls != 0 || !slices.Equal(forge.reviews, []string{head}) || forge.created != 0 {
				t.Fatalf("runs=%d reviews=%v created=%d", forge.runCalls, forge.reviews, forge.created)
			}
			if forge.mergedHead != head {
				t.Fatalf("item head=%s want %s", forge.mergedHead, head)
			}
			if !slices.Equal(item.CandidateCommits, []string{head}) {
				t.Fatalf("candidates=%v", item.CandidateCommits)
			}
			if name == "only_active_Spec" {
				if !forge.specsRootMissing {
					t.Fatal("item still has a Specs Root")
				}
			}
		})
	}
}

func TestDeliverStartParksAnUnprovenArchivedItemBranch(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"inexact_archive_commit", "folder_removed_without_archive", "no_removing_commit"} {
		t.Run(name, func(t *testing.T) {
			home, repo, _, head, archiveHead := archivedResumeBranch(t, false, name)
			source := "docs/specs/" + implementTestSlug
			cause := ""
			switch name {
			case "inexact_archive_commit":
				cause = fmt.Sprintf("commit %s that removed %s is not an exact archive of the Spec", deliveryCommitName(archiveHead), source)
			case "folder_removed_without_archive":
				cause = fmt.Sprintf("item head %s holds neither %s nor an archive of the Spec", deliveryCommitName(head), source)
			case "no_removing_commit":
				cause = "no commit on the item branch removed " + source + "/_prd.md"
			}
			forge := &archivedResumeForge{}
			item := startArchivedResumeQueue(t, home, repo, forge)
			if item.Stage != store.DeliveryStageParked || item.Blocker != "archived-item-unproven: "+cause {
				t.Fatalf("item=%+v want cause %q", item, cause)
			}
			if forge.runCalls != 0 || len(forge.reviews) != 0 || forge.created != 0 {
				t.Fatalf("runs=%d reviews=%v created=%d", forge.runCalls, forge.reviews, forge.created)
			}
			var out, diag bytes.Buffer
			code := runCLI(t, []string{"deliver", "status"}, &out, &diag)
			want := "Park: " + implementTestSlug + " finding: in " + item.Worktree + ", make the item branch hold the Spec folder or one exact archive commit of the Spec, then run roundfix deliver retry " + implementTestSlug
			if code != exitOK || diag.Len() != 0 || !strings.Contains(out.String(), want) {
				t.Fatalf("status exit=%d stdout=%s stderr=%s want=%s", code, &out, &diag, want)
			}
		})
	}
}

func TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive(t *testing.T) {
	t.Parallel()
	home, repo := newItemBranchStartRepository(t)
	branch := itemBranchWithWork(t, repo, "1111111111111111", 1)
	gittest.Run(t, repo, "checkout", branch)
	taskPath := filepath.Join(repo, "docs/specs", implementTestSlug, "task_01.md")
	mustWrite(t, taskPath, mustRead(t, taskPath)+"\n## Consumers\n\nNone. This fixture changes no shared contract.\n")
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "test: declare control consumers")
	gittest.Run(t, repo, "checkout", "main")
	forge := &archivedResumeForge{parkTestDeliveryFlow: parkTestDeliveryFlow{parkSlug: implementTestSlug}}
	item := startArchivedResumeQueue(t, home, repo, forge)
	if item.Branch != branch || item.Stage != store.DeliveryStageParked || item.Blocker != delivery.BlockerRunUnresolved || forge.runCalls != 1 || len(forge.reviews) != 0 {
		t.Fatalf("item=%+v runs=%d reviews=%v", item, forge.runCalls, forge.reviews)
	}
}
