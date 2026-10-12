package cli

// Suite: Cancel Command transcripts and refusal atomicity.
// Boundary: CLI dispatch over committed temporary repositories and homes.

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func cancelFixture(t *testing.T, completed, pending int) (string, commandEnvironment) {
	t.Helper()
	seeds := make([]implementSeed, completed+pending)
	for i := range seeds {
		status := string(spec.StatusPending)
		if i < completed {
			status = string(spec.StatusCompleted)
		}
		seeds[i] = implementSeed{id: fmt.Sprintf("task_%02d", i+1), status: status}
	}
	home, repo := newImplementWorkspace(t, seeds)
	return repo, commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
}

func cancelRun(t *testing.T, env commandEnvironment, args ...string) (int, string, string) {
	t.Helper()
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), append([]string{"cancel"}, args...), &out, &errout, env)
	return code, out.String(), errout.String()
}

func cancelAssertTranscript(t *testing.T, code int, out, errout string, wantCode int, wantOut, wantErr string) {
	t.Helper()
	if code != wantCode || out != wantOut || errout != wantErr {
		t.Fatalf("exit=%d want=%d\nstdout=%q\nwant=%q\nstderr=%q\nwant=%q", code, wantCode, out, wantOut, errout, wantErr)
	}
}

func cancelRefusal(reason string) string {
	return "Preflight failed\n\nReason:\n  " + reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix cancel --help' for usage.\n"
}

func cancelSummary(t *testing.T, repo string, completed, pending int) string {
	t.Helper()
	files := historySnapshot(t, filepath.Join(repo, "docs/specs", implementTestSlug))
	size := 0
	for _, content := range files {
		size += len(content)
	}
	revision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	return fmt.Sprintf("cancelled %s -> docs/history/specs/%s.md; %d completed Task(s), %d pending; removed %d file(s) (%d bytes) kept in Git at %.12s\n", implementTestSlug, implementTestSlug, completed, pending, len(files), size, revision)
}

func cancelQueue(t *testing.T, repo string, env commandEnvironment, stage store.DeliveryStage) {
	t.Helper()
	db, err := store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = stage
	if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
}

func cancelUnchanged(t *testing.T, repo string, env commandEnvironment, wantErr string, args ...string) {
	t.Helper()
	before := historySnapshot(t, repo)
	status := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all")
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	hash := cancelFolderHash(t, folder)
	code, out, errout := cancelRun(t, env, args...)
	cancelAssertTranscript(t, code, out, errout, 2, "", wantErr)
	historyUnchanged(t, historyFixture{repo: repo}, before, status)
	if got := cancelFolderHash(t, folder); got != hash {
		t.Fatal("refusal changed Spec folder hash")
	}
}

func cancelFolderHash(t *testing.T, folder string) [32]byte {
	t.Helper()
	files := historySnapshot(t, folder)
	// Walk order is deterministic; include paths as well as bytes in the digest.
	var content bytes.Buffer
	err := filepath.WalkDir(folder, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(folder, path)
		if err != nil {
			return err
		}
		fmt.Fprintf(&content, "%s\x00%d\x00", filepath.ToSlash(rel), len(files[filepath.ToSlash(rel)]))
		content.Write(files[filepath.ToSlash(rel)])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return sha256.Sum256(content.Bytes())
}

func TestCancelTranscriptCancelsASpecWithSources(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 5, 3)
	root := filepath.Join(repo, "docs/specs", implementTestSlug, "references")
	historyWrite(t, filepath.Join(root, "_index.md"), "| path | type |\n| --- | --- |\n| entry.md | backlog |\n| finding.md | finding |\n")
	historyWrite(t, filepath.Join(root, "entry.md"), "---\nstatus: promoted\nspec: "+implementTestSlug+"\n---\n\n# Entry\n\nIntent.\n")
	historyWrite(t, filepath.Join(root, "finding.md"), "---\nstatus: done\nupdated_at: 2026-10-01\nabsorbed_by: "+implementTestSlug+"\n---\n\n# Finding\n\nEvidence.\n")
	historyCommit(t, repo, "docs: adopted sources")
	want := cancelSummary(t, repo, 5, 3) + "returned docs/backlog/entry.md as open\nreturned docs/findings/finding.md as pending\n"
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	cancelAssertTranscript(t, code, out, errout, 0, want, "")
	record := archiveRecordAt(t, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"))
	if record.Disposition != spec.ArchiveCancelled || record.CancelReason != "no consumer" || len(record.CompletedTasks) != 5 || len(record.PendingTasks) != 3 {
		t.Fatalf("record=%+v", record)
	}
	for path, status := range map[string]string{"docs/backlog/entry.md": "open", "docs/findings/finding.md": "pending"} {
		if !strings.Contains(mustRead(t, filepath.Join(repo, path)), "status: "+status+"\n") {
			t.Fatalf("source %s not returned", path)
		}
	}
	assertPathMissing(t, filepath.Join(repo, "docs/specs", implementTestSlug))
	assertPathMissing(t, store.DatabasePath(env.homeDir))
}

func TestCancelTranscriptWarnsAboutAParkedItemAndACitation(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	target := "docs/specs/" + implementTestSlug + "/"
	historyWrite(t, filepath.Join(repo, "docs/handoffs/notes.md"), "# Notes\n\nRead `"+target+"`.\n")
	historyCommit(t, repo, "docs: live citation")
	cancelQueue(t, repo, env, store.DeliveryStageParked)
	want := cancelSummary(t, repo, 0, 2)
	wantErr := "roundfix: warning: Delivery Queue item " + implementTestSlug + " stays parked; a Delivery Retry of a cancelled Spec is refused, so close its branch and pull request by hand\n" +
		"roundfix: warning: docs/handoffs/notes.md:3 cites " + target + ", which this cancellation removed\n"
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	cancelAssertTranscript(t, code, out, errout, 0, want, wantErr)
}

func TestCancelTranscriptRefusesAnUnfinishedQueueItem(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	cancelQueue(t, repo, env, store.DeliveryStageRunning)
	reason := fmt.Sprintf("Spec %q cannot be cancelled while its Delivery Queue item is at stage running; stop the queue with: roundfix deliver stop", implementTestSlug)
	cancelUnchanged(t, repo, env, cancelRefusal(reason), implementTestSlug, "--reason", "no consumer")
}

func TestCancelTranscriptRefusesAnActiveRun(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	db, err := store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	active, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, SpecSlug: implementTestSlug, LocalBranch: "roundfix/run-cancel-test", OwnerPID: os.Getpid()})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reason := fmt.Sprintf("Active Run %s holds Spec %q; stop it with: roundfix stop %s", active.ID, implementTestSlug, active.ID)
	cancelUnchanged(t, repo, env, cancelRefusal(reason), implementTestSlug, "--reason", "no consumer")
}

func TestCancelTranscriptRefusesARequiringSpec(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	writeImplementSpec(t, repo, "dependent", []implementSeed{{id: "task_01", status: string(spec.StatusPending)}})
	path := filepath.Join(repo, "docs/specs/dependent/_tasks.md")
	mustWrite(t, path, strings.Replace(mustRead(t, path), "schema: spec-tasks/v1\n", "schema: spec-tasks/v1\nrequires: ["+implementTestSlug+"]\n", 1))
	historyCommit(t, repo, "docs: requiring Spec")
	reason := fmt.Sprintf("Spec %q cannot be cancelled while active Spec dependent requires it; remove it from that Spec's requires first", implementTestSlug)
	cancelUnchanged(t, repo, env, cancelRefusal(reason), implementTestSlug, "--reason", "no consumer")
}

func TestCancelTranscriptPrintsHelp(t *testing.T) {
	t.Parallel()
	_, env := cancelFixture(t, 0, 2)
	want := `Usage:
  roundfix cancel <slug> --reason <text>

Cancels an active Spec that will not be delivered. Writes one Archive Record
with disposition cancelled, the reason, the date and the Tasks completed and
pending, returns each adopted Backlog Entry to docs/backlog/ (status: open)
and each adopted Finding to docs/findings/ (status: pending), and removes the
committed Spec folder. The removed bytes stay in Git at the record's
source_revision. cancel creates no Run, never changes a Task status, and
never commits or pushes.

Refuses while an Active Run or an unfinished Delivery Queue item holds the
Spec, while an active Spec requires it, and for a Spec that carries a
supersession.

Options:
  --reason  Why the Spec is cancelled, recorded in the Archive Record

Exit codes:
  0  cancelled
  1  cancel write failed
  2  Preflight Validation failed
`
	code, out, errout := cancelRun(t, env, "--help")
	cancelAssertTranscript(t, code, out, errout, 0, want, "")
}

func TestCancelRefusesEveryUnfinishedQueueStage(t *testing.T) {
	t.Parallel()
	for _, stage := range []store.DeliveryStage{store.DeliveryStageQueued, store.DeliveryStageRunning, store.DeliveryStageReviewing, store.DeliveryStageArchiving, store.DeliveryStageGating, store.DeliveryStagePublishing, store.DeliveryStageChecking, store.DeliveryStageMerging} {
		t.Run(string(stage), func(t *testing.T) {
			t.Parallel()
			repo, env := cancelFixture(t, 0, 2)
			cancelQueue(t, repo, env, stage)
			reason := fmt.Sprintf("Spec %q cannot be cancelled while its Delivery Queue item is at stage %s; stop the queue with: roundfix deliver stop", implementTestSlug, stage)
			cancelUnchanged(t, repo, env, cancelRefusal(reason), implementTestSlug, "--reason", "no consumer")
		})
	}
}

func TestCancelRefusalsLeaveTheTreeUnchanged(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"missing reason", "unknown Spec", "uncommitted", "supersession", "non-Markdown pin"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			repo, env := cancelFixture(t, 0, 2)
			args := []string{implementTestSlug, "--reason", "no consumer"}
			var reason string
			switch name {
			case "missing reason":
				args = []string{implementTestSlug}
				reason = "missing required --reason; pass --reason <text>"
			case "unknown Spec":
				args[0] = "unknown"
				reason = `unknown active Spec slug "unknown"`
			case "uncommitted":
				mustWrite(t, filepath.Join(repo, "docs/specs", implementTestSlug, "_prd.md"), "changed\n")
				reason = fmt.Sprintf("Spec %q has changes not committed at HEAD under docs/specs/%s; commit them before cancel so the Archive Record's source_revision holds the Spec", implementTestSlug, implementTestSlug)
			case "supersession":
				historyWrite(t, filepath.Join(repo, "docs/specs", implementTestSlug, "_supersession.md"), "supersession\n")
				historyCommit(t, repo, "docs: supersession")
				reason = fmt.Sprintf("cannot cancel a Spec with a supersession: cancellation destination %q already exists", filepath.Join(repo, "docs/specs", implementTestSlug, "_supersession.md"))
			case "non-Markdown pin":
				historyWrite(t, filepath.Join(repo, "pin.txt"), "docs/specs/"+implementTestSlug+"/\n")
				historyCommit(t, repo, "docs: pinned directory")
				reason = fmt.Sprintf("Spec %q cannot cancel while another file names its active directory docs/specs/%s/: pin.txt:1", implementTestSlug, implementTestSlug)
			}
			cancelUnchanged(t, repo, env, cancelRefusal(reason), args...)
		})
	}
}

func TestRootHelpNamesTheCancelCommand(t *testing.T) {
	t.Parallel()
	_, env := cancelFixture(t, 0, 2)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"--help"}, &out, &errout, env)
	if code != 0 || errout.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, &errout)
	}
	for _, line := range []string{"  roundfix archive <slug>\n  roundfix cancel <slug> --reason <text>\n", "  archive    Archive a completed Spec\n  cancel     Cancel an active Spec with an Archive Record\n"} {
		if !strings.Contains(out.String(), line) {
			t.Fatalf("help lacks %q", line)
		}
	}
}

func TestCancelAcceptsAMergedQueueItem(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	cancelQueue(t, repo, env, store.DeliveryStageMerged)
	want := cancelSummary(t, repo, 0, 2)
	before := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	cancelAssertTranscript(t, code, out, errout, 0, want, "")
	if got := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); got != before {
		t.Fatal("cancel committed")
	}
}

func TestCancelRefusesMalformedRequires(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	path := filepath.Join(repo, "docs/specs/dependent/_tasks.md")
	historyWrite(t, path, "---\nrequires: bad\n---\n")
	historyCommit(t, repo, "docs: invalid requires")
	cancelUnchanged(t, repo, env, cancelRefusal("read active Spec dependent prerequisites: _tasks.md: requires must be a list of strings"), implementTestSlug, "--reason", "no consumer")
}

func TestCancelReclaimsAnOrphanedRun(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	db, err := store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	active, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, SpecSlug: implementTestSlug, LocalBranch: "roundfix/run-cancel-orphan", OwnerPID: 2147483647})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	want := cancelSummary(t, repo, 0, 2)
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	if code != 0 || out != want || !strings.Contains(errout, "reclaimed orphaned Active Run "+active.ID) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, errout)
	}
	db, err = store.Open(t.Context(), env.homeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()
	run, found, err := db.Run(t.Context(), active.ID)
	if err != nil || !found || run.State != store.StateFailed {
		t.Fatalf("run=%+v found=%v err=%v", run, found, err)
	}
}

func TestCancelRefusesOccupiedSourceDestination(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	folder := filepath.Join(repo, "docs/specs", implementTestSlug, "references")
	historyWrite(t, filepath.Join(folder, "_index.md"), "| path | type |\n| --- | --- |\n| entry.md | backlog |\n")
	historyWrite(t, filepath.Join(folder, "entry.md"), "---\nstatus: promoted\n---\n\n# Entry\n")
	historyWrite(t, filepath.Join(repo, "docs/backlog/entry.md"), "existing\n")
	historyCommit(t, repo, "docs: occupied source")
	reason := fmt.Sprintf("cancellation destination %q already exists", filepath.Join(repo, "docs/backlog/entry.md"))
	cancelUnchanged(t, repo, env, cancelRefusal(reason), implementTestSlug, "--reason", "no consumer")
}

func TestCancelWriteFailureUsesExitOne(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	archiveRoot := filepath.Join(repo, "docs/history/specs")
	mustMkdir(t, filepath.Dir(archiveRoot))
	// A dangling parent passes destination validation but cannot be created.
	if err := os.Symlink("missing", archiveRoot); err != nil {
		t.Fatal(err)
	}
	historyCommit(t, repo, "docs: blocked archive parent")
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	before := historySnapshot(t, folder)
	status := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all")
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	if code != 1 || out != "" || !strings.HasPrefix(errout, "roundfix: cancel write failed: create cancellation destination directory:") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, errout)
	}
	if !reflect.DeepEqual(before, historySnapshot(t, folder)) {
		t.Fatal("write failure changed Spec bytes")
	}
	if got := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all"); got != status {
		t.Fatalf("write failure changed status: %q -> %q", status, got)
	}
	if target, err := os.Readlink(archiveRoot); err != nil || target != "missing" {
		t.Fatalf("archive parent changed: %q %v", target, err)
	}
}

func TestCancelArgumentRefusalsFollowContractOrder(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		args   []string
		reason string
	}{
		{"missing slug", nil, "missing required Spec slug; pass roundfix cancel <slug> --reason <text>"},
		{"flag without slug", []string{"--reason", "no consumer"}, "missing required Spec slug; pass roundfix cancel <slug> --reason <text>"},
		{"blank reason", []string{implementTestSlug, "--reason", " \t\n "}, "missing required --reason; pass --reason <text>"},
		{"reason before unexpected", []string{implementTestSlug, "extra"}, "missing required --reason; pass --reason <text>"},
		{"unexpected argument", []string{implementTestSlug, "--reason", "no consumer", "extra"}, `unexpected argument "extra"`},
		{"unsafe slug", []string{"../escape", "--reason", "no consumer"}, `unsafe Spec slug "../escape"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, env := cancelFixture(t, 0, 2)
			cancelUnchanged(t, repo, env, cancelRefusal(tc.reason), tc.args...)
			assertPathMissing(t, store.DatabasePath(env.homeDir))
		})
	}
}

func TestCancelRefusesOutsideGit(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	env.workDir = t.TempDir()
	cancelUnchanged(t, repo, env, cancelRefusal("cancel requires a git repository working tree"), implementTestSlug, "--reason", "no consumer")
}

func TestCancelWithoutATaskGraph(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	for _, name := range []string{"_tasks.md", "task_01.md", "task_02.md"} {
		if err := os.Remove(filepath.Join(folder, name)); err != nil {
			t.Fatal(err)
		}
	}
	historyCommit(t, repo, "docs: undecomposed Spec")
	want := cancelSummary(t, repo, 0, 0)
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "  no\n consumer  ")
	cancelAssertTranscript(t, code, out, errout, 0, want, "")
	record := archiveRecordAt(t, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"))
	if record.CancelReason != "no consumer" {
		t.Fatalf("reason=%q", record.CancelReason)
	}
}

func TestCancelWarnsOnlyOncePerRemainingTrackedCitationLine(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	folder := "docs/specs/" + implementTestSlug + "/"
	source := folder + "_prd.md"
	historyWrite(t, filepath.Join(repo, "docs/backlog/followup.md"), "# Follow-up\n\n[Spec](../specs/"+implementTestSlug+"/_prd.md)\n")
	historyWrite(t, filepath.Join(repo, "docs/handoffs/notes.md"), "# Notes\n\nRead `"+source+"`.\n")
	historyWrite(t, filepath.Join(repo, "docs/history/old.md"), "Read `"+source+"`.\n")
	historyCommit(t, repo, "docs: multiple citation surfaces")
	historyWrite(t, filepath.Join(repo, "untracked.md"), "Read `"+source+"`.\n")
	want := cancelSummary(t, repo, 0, 2)
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	cancelAssertTranscript(t, code, out, errout, 0, want, "roundfix: warning: docs/handoffs/notes.md:3 cites "+folder+", which this cancellation removed\n")
	if got := mustRead(t, filepath.Join(repo, "docs/backlog/followup.md")); got != "# Follow-up\n\n[Spec](../history/specs/"+implementTestSlug+".md)\n" {
		t.Fatalf("live document not relinked: %q", got)
	}
}

func TestCancelWarnsAboutACitationOnAPartlyRelinkedLine(t *testing.T) {
	t.Parallel()
	repo, env := cancelFixture(t, 0, 2)
	folder := "docs/specs/" + implementTestSlug + "/"
	relative := "../specs/" + implementTestSlug + "/_prd.md"
	// Relink preserves a code example's relative link while updating a normal link.
	historyWrite(t, filepath.Join(repo, "docs/backlog/followup.md"), "# Follow-up\n\n[Normal]("+relative+") and `[Code]("+relative+")`.\n")
	historyCommit(t, repo, "docs: mixed citation line")
	want := cancelSummary(t, repo, 0, 2)
	code, out, errout := cancelRun(t, env, implementTestSlug, "--reason", "no consumer")
	cancelAssertTranscript(t, code, out, errout, 0, want, "roundfix: warning: docs/backlog/followup.md:3 cites "+folder+", which this cancellation removed\n")
}
