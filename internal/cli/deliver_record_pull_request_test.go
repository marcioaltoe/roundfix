package cli

// Suite: delivery Pull Request Record Commit using real Git.
// Invariant: recording and resuming change only the number and renderer outcome.
// Boundary IN: workflow, Archive Command, temporary Git repositories and queue.
// Boundary OUT: the end-to-end forge is local and never reaches the network.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func pullRequestRecordFixture(t *testing.T) (*commandDeliveryWorkflow, string, string, string) {
	t.Helper()
	workflow, repo, _, _, destination := retirementFixture(t, false)
	origin := t.TempDir()
	gittest.InitRepo(t, origin, "--bare", "--initial-branch=main")
	gittest.Run(t, repo, "remote", "add", "origin", origin)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "remote", "set-head", "origin", "main")
	gittest.Run(t, repo, "switch", "-c", "roundfix/deliver-record")
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "docs: archive Spec")
	gittest.PersistIdentity(t, repo)
	setCommandEnvForTest(t, cliTestHelperEnv, "1")
	head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	return workflow, repo, head, destination + ".md"
}

func TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord(t *testing.T) {
	t.Parallel()
	workflow, repo, candidate, path := pullRequestRecordFixture(t)
	before := archiveRecordAt(t, filepath.Join(repo, path))
	result, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
	if err != nil || !result.Exact || result.Head == "" || result.Head == candidate {
		t.Fatalf("record=%+v err=%v", result, err)
	}
	after := archiveRecordAt(t, filepath.Join(repo, path))
	before.PullRequest = "7"
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("record changed other fields: before=%+v after=%+v", before, after)
	}
	if changed := gittest.Run(t, repo, "diff", "--name-only", candidate, result.Head); strings.TrimSpace(changed) != path {
		t.Fatalf("changed=%q", changed)
	}
	wantSubject := "docs: record pull request #7 in the 0001 Archive Record"
	if subject := strings.TrimSpace(gittest.Run(t, repo, "show", "-s", "--format=%s", result.Head)); subject != wantSubject {
		t.Fatalf("subject=%q", subject)
	}
	// Simulate a crash after commit but before the queue saved the new head.
	adopted, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
	if err != nil || adopted != result {
		t.Fatalf("adopt=%+v err=%v", adopted, err)
	}
	again, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, result.Head, "7")
	if err != nil || again != (delivery.PullRequestRecord{}) {
		t.Fatalf("idempotence=%+v err=%v", again, err)
	}
	if count := strings.TrimSpace(gittest.Run(t, repo, "rev-list", "--count", candidate+"..HEAD")); count != "1" {
		t.Fatalf("record commits=%s", count)
	}
	if dirty := gittest.Run(t, repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%s", dirty)
	}
}

func TestDeliveryWorkflowRefusesARecordCommitThatChangesAnotherPath(t *testing.T) {
	t.Parallel()
	for _, mutation := range []string{"extra path", "title", "source revision", "outcome", "number", "mode", "grandchild", "merge", "unknown field", "extra prose", "body title"} {
		t.Run(mutation, func(t *testing.T) {
			workflow, repo, candidate, path := pullRequestRecordFixture(t)
			if _, err := spec.RecordArchivePullRequest(filepath.Join(repo, path), "7"); err != nil {
				t.Fatal(err)
			}
			switch mutation {
			case "extra path":
				mustWrite(t, filepath.Join(repo, "unrelated.md"), "Unrelated\n")
			case "mode":
				if err := os.Chmod(filepath.Join(repo, path), 0755); err != nil {
					t.Fatal(err)
				}
			case "unknown field":
				data := mustRead(t, filepath.Join(repo, path))
				mustWrite(t, filepath.Join(repo, path), strings.Replace(data, "schema:", "unreviewed: true\nschema:", 1))
			case "extra prose":
				mustWrite(t, filepath.Join(repo, path), mustRead(t, filepath.Join(repo, path))+"\nAdditional unreviewed paragraph\n")
			case "body title":
				data := mustRead(t, filepath.Join(repo, path))
				mustWrite(t, filepath.Join(repo, path), strings.Replace(data, "# ", "# Different ", 1))
			case "title", "source revision", "outcome", "number":
				record := archiveRecordAt(t, filepath.Join(repo, path))
				switch mutation {
				case "title":
					record.Title = "Other title"
				case "source revision":
					record.SourceRevision = strings.Repeat("b", 40)
				case "outcome":
					record.Outcome = "Invented outcome"
				case "number":
					record.PullRequest = "8"
				}
				data, err := spec.RenderArchiveRecord(record)
				if err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(repo, path), string(data))
			}
			gittest.Run(t, repo, "add", "-A")
			if mutation == "merge" {
				// Same tree change, two parents: it is not the candidate's one child.
				gittest.Run(t, repo, "branch", "other-parent", candidate+"^")
				gittest.Run(t, repo, "commit", "-m", "docs: record number")
				tree := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD^{tree}"))
				merged := strings.TrimSpace(gittest.Run(t, repo, "commit-tree", tree, "-p", candidate, "-p", "other-parent", "-m", "docs: merge number"))
				gittest.Run(t, repo, "reset", "--hard", merged)
			} else {
				gittest.Run(t, repo, "commit", "-m", "docs: inexact record")
				if mutation == "grandchild" {
					gittest.Run(t, repo, "commit", "--allow-empty", "-m", "extra child")
				}
			}
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			result, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
			if err != nil || result.Exact || result.Head != head {
				t.Fatalf("inexact=%+v err=%v", result, err)
			}
			if got := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); got != head {
				t.Fatalf("refusal wrote commit=%s", got)
			}
		})
	}
}

func TestDeliveryWorkflowAdoptsAnUncommittedRecord(t *testing.T) {
	t.Parallel()
	for _, staged := range []bool{false, true} {
		t.Run(fmt.Sprint("staged=", staged), func(t *testing.T) {
			workflow, repo, candidate, path := pullRequestRecordFixture(t)
			if _, err := spec.RecordArchivePullRequest(filepath.Join(repo, path), "7"); err != nil {
				t.Fatal(err)
			}
			// A valid alternate YAML layout must survive adoption byte for byte.
			original := strings.Replace(mustRead(t, filepath.Join(repo, path)), "pull_request: \"7\"", "pull_request: '7'", 1)
			mustWrite(t, filepath.Join(repo, path), original)
			if staged {
				gittest.Run(t, repo, "add", "--", path)
			}
			// No Archive Command can run successfully with this helper disabled.
			setCommandEnvForTest(t, cliTestHelperEnv, "0")
			result, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
			if err != nil || !result.Exact || result.Head == candidate {
				t.Fatalf("adopt=%+v err=%v", result, err)
			}
			if got := mustRead(t, filepath.Join(repo, path)); got != original {
				t.Fatal("adoption rewrote record")
			}
			if count := strings.TrimSpace(gittest.Run(t, repo, "rev-list", "--count", candidate+"..HEAD")); count != "1" {
				t.Fatalf("commits=%s", count)
			}
		})
	}
}

func TestDeliveryWorkflowRefusesAnUncommittedInexactRecord(t *testing.T) {
	t.Parallel()
	for _, mutation := range []string{"extra path", "staged extra path", "unchanged number", "other number", "other field", "mode"} {
		t.Run(mutation, func(t *testing.T) {
			workflow, repo, candidate, path := pullRequestRecordFixture(t)
			if _, err := spec.RecordArchivePullRequest(filepath.Join(repo, path), "7"); err != nil {
				t.Fatal(err)
			}
			switch mutation {
			case "extra path", "staged extra path":
				mustWrite(t, filepath.Join(repo, "unrelated.md"), "Unrelated\n")
				if mutation == "staged extra path" {
					gittest.Run(t, repo, "add", "unrelated.md")
				}
			case "mode":
				if err := os.Chmod(filepath.Join(repo, path), 0755); err != nil {
					t.Fatal(err)
				}
			default:
				record := archiveRecordAt(t, filepath.Join(repo, path))
				switch mutation {
				case "unchanged number":
					record.PullRequest = ""
					record.Outcome = "Changed"
				case "other number":
					record.PullRequest = "8"
				case "other field":
					record.Relinked = []string{"docs/changed.md"}
				}
				data, err := spec.RenderArchiveRecord(record)
				if err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(repo, path), string(data))
			}
			result, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
			if err != nil || result.Exact || result.Head != candidate {
				t.Fatalf("inexact=%+v err=%v", result, err)
			}
			if got := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); got != candidate {
				t.Fatal("refusal committed")
			}
		})
	}
}

func TestDeliveryWorkflowSkipsALegacyArchiveFolder(t *testing.T) {
	t.Parallel()
	repo, _, head := commitLinkRewritingArchive(t, []string{"_prd.md", "task_01.md"}, nil)
	result, err := retirementWorkflow(repo, "").RecordPullRequest(t.Context(), repo, implementTestSlug, head, "7")
	if err != nil || result != (delivery.PullRequestRecord{}) {
		t.Fatalf("legacy=%+v err=%v", result, err)
	}
}

func TestDeliveryRecordProofAllowsOnlyTheRenderedOutcomeShortening(t *testing.T) {
	t.Parallel()
	workflow, repo, candidate, path := pullRequestRecordFixture(t)
	record := archiveRecordAt(t, filepath.Join(repo, path))
	record.Outcome = strings.Repeat("A long outcome without punctuation ", 100)
	before, err := spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	old, err := spec.ParseArchiveRecord(before)
	if err != nil {
		t.Fatal(err)
	}
	old.PullRequest = "7"
	after, err := spec.RenderArchiveRecord(old)
	if err != nil {
		t.Fatal(err)
	}
	newRecord, err := spec.ParseArchiveRecord(after)
	if err != nil {
		t.Fatal(err)
	}
	if len(newRecord.Outcome) >= len(old.Outcome) {
		t.Fatal("fixture did not shorten outcome")
	}
	mustWrite(t, filepath.Join(repo, path), string(before))
	gittest.Run(t, repo, "add", path)
	gittest.Run(t, repo, "commit", "-m", "docs: long outcome")
	candidate = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	result, err := workflow.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
	if err != nil || !result.Exact {
		t.Fatalf("shortened record=%+v err=%v", result, err)
	}
	if got := mustRead(t, filepath.Join(repo, path)); got != string(after) {
		t.Fatal("unexpected shortening")
	}
}

func TestCommandDeliveryEngineRecordsThePullRequest(t *testing.T) {
	t.Parallel()
	// The real constructor resolves gh through the process PATH. Run this
	// scenario in a child with its own PATH, without changing the suite's env.
	const forgeBinEnv = "ROUNDFIX_PULL_REQUEST_RECORD_FORGE_BIN"
	bin := os.Getenv(forgeBinEnv)
	if bin == "" {
		bin = t.TempDir()
		command := exec.CommandContext(t.Context(), scriptFixtureBinaryPath,
			"-test.run=^TestCommandDeliveryEngineRecordsThePullRequest$", "-test.count=1", "-test.v")
		command.Env = withEnvValue(isolatedGitEnvForTest(), forgeBinEnv, bin)
		command.Env = withEnvValue(command.Env, "PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
		command.Env = withEnvValue(command.Env, "HOME", t.TempDir())
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("command engine with isolated forge: %v\n%s", err, output)
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
	workflow := &commandDeliveryWorkflow{store: runStore, loaded: cfg, git: retirementWorkflow(repo, home).git}
	branch, workDir, err := workflow.CreateItemBranch(commandContextForTest(t, t.Context()), repo, implementTestSlug)
	if err != nil {
		t.Fatal(err)
	}
	candidate := strings.TrimSpace(gittest.Run(t, workDir, "rev-parse", "HEAD"))
	archived, err := spec.Archive(spec.ArchiveRequest{SpecsRoot: filepath.Join(workDir, "docs/specs"), BuiltInRoot: true, Slug: implementTestSlug, SourceRevision: candidate, RepositoryRoot: workDir})
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
 'pr create') touch "$state/open"; printf 'https://local.invalid/pull/7\n'; exit 0 ;;
 'pr checks') printf '[{"name":"verify","bucket":"pass","state":"SUCCESS"}]'; exit 0 ;;
 'pr merge')
  test "$3" = 7
  test "$4" = --squash
  test "$5" = --match-head-commit
  test "$6" = "$head"
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
status=OPEN
merge='null'
if test -f "$state/merged"; then status=MERGED; merge="{\"oid\":\"$(cat "$state/merged")\"}"; fi
payload=$(printf '{"number":7,"state":"%%s","headRefName":"%%s","headRefOid":"%%s","mergeable":"MERGEABLE","mergeStateStatus":"CLEAN","mergeCommit":%%s}' "$status" "$branch" "$head" "$merge")
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
	path, err := filepath.Rel(workDir, archived.RecordPath)
	if err != nil {
		t.Fatal(err)
	}
	mergedContent := gittest.Run(t, origin, "show", item.MergeCommit+":"+filepath.ToSlash(path))
	record, err := spec.ParseArchiveRecord([]byte(mergedContent))
	if err != nil || record.PullRequest != "7" {
		t.Fatalf("merged record=%+v err=%v", record, err)
	}
	if len(item.CandidateCommits) != 3 {
		t.Fatalf("candidates=%v", item.CandidateCommits)
	}
	if count := strings.TrimSpace(gittest.Run(t, origin, "rev-list", "--count", archiveHead+".."+item.CandidateCommits[2])); count != "1" {
		t.Fatalf("record commits=%s", count)
	}
	refs := strings.Fields(gittest.Run(t, origin, "for-each-ref", "--format=%(refname)"))
	if !reflect.DeepEqual(refs, []string{"refs/heads/main", "refs/heads/" + branch}) {
		t.Fatalf("forge refs=%v", refs)
	}
	if _, err := engine.Run(commandContextForTest(t, t.Context()), repo); err != nil {
		t.Fatal(err)
	}
	if got := readDeliveryItemForCLI(t, runStore, repo); !reflect.DeepEqual(item.CandidateCommits, got.CandidateCommits) {
		t.Fatal("resume appended another commit")
	}
}
