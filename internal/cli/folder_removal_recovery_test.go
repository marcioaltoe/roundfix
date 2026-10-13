package cli

// Suite: forward recovery transcripts for archive and cancel.
// Boundary: CLI dispatch, real filesystem permissions, and a POSIX shell.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func TestArchiveFolderRemovalFailureRecoversForward(t *testing.T) {
	t.Parallel()
	repo, env := folderRemovalFixture(t, "archive")
	assertCommandRecoversForward(t, repo, env, "archive")
}

func TestCancelFolderRemovalFailureRecoversForward(t *testing.T) {
	t.Parallel()
	repo, env := folderRemovalFixture(t, "cancel")
	assertCommandRecoversForward(t, repo, env, "cancel")
}

func TestArchiveFolderRemovalFailureRecordsThePullRequest(t *testing.T) {
	t.Parallel()
	repo, env := folderRemovalFixture(t, "archive")
	gittest.Run(t, repo, "update-ref", "refs/remotes/origin/main", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	assertCommandRecoversForward(t, repo, env, "archive", "--pull-request", "263")
	record := archiveRecordAt(t, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"))
	if record.PullRequest != "263" {
		t.Fatalf("pull_request=%q, want 263", record.PullRequest)
	}
}

func TestFolderRemovalFinishCommandQuotesTheFolder(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	folder := filepath.Join(root, "a space and 'quote")
	references := filepath.Join(folder, "references")
	mustMkdir(t, references)
	mustWrite(t, filepath.Join(references, "note.md"), "remove only this folder\n")
	sibling := filepath.Join(root, "keep.md")
	mustWrite(t, sibling, "keep these bytes\n")
	makeCLIReferencesReadOnly(t, folder)
	quoted := "'" + strings.ReplaceAll(folder, "'", "'\\''") + "'"
	want := "chmod -R u+w " + quoted + " && rm -rf " + quoted
	command := folderRemovalFinishCommand(folder)
	if command != want {
		t.Fatalf("command=%q, want %q", command, want)
	}
	runFolderRemovalFinish(t, command)
	assertPathMissing(t, folder)
	if got := mustRead(t, sibling); got != "keep these bytes\n" {
		t.Fatalf("sibling changed: %q", got)
	}
}

func folderRemovalFixture(t *testing.T, command string) (string, commandEnvironment) {
	t.Helper()
	pending := 0
	if command == "cancel" {
		pending = 1
	}
	repo, env := cancelFixture(t, 1, pending)
	folder := filepath.Join(repo, "docs/specs", implementTestSlug)
	if command == "archive" {
		writeArchiveQAReport(t, repo, spec.VerdictPass, "rows_passed: 1")
	}
	mustWrite(t, filepath.Join(repo, ".gitignore"), "*.log\n")
	historyWrite(t, filepath.Join(folder, "references/note.md"), "Promoted note.\n")
	historyWrite(t, filepath.Join(repo, "docs/backlog/followup.md"), "[Note](../specs/"+implementTestSlug+"/references/note.md)\n[Spec](../specs/"+implementTestSlug+"/_prd.md)\n")
	if command == "cancel" {
		historyWrite(t, filepath.Join(folder, "references/_index.md"), "| path | type |\n| --- | --- |\n| entry.md | backlog |\n| finding.md | finding |\n")
		historyWrite(t, filepath.Join(folder, "references/entry.md"), "---\nstatus: promoted\nspec: "+implementTestSlug+"\n---\n\n# Entry\n\nIntent.\n")
		historyWrite(t, filepath.Join(folder, "references/finding.md"), "---\nstatus: done\nupdated_at: 2026-10-01\nabsorbed_by: "+implementTestSlug+"\n---\n\n# Finding\n\nEvidence.\n")
	}
	historyCommit(t, repo, "docs: folder removal fixture")
	makeCLIReferencesReadOnly(t, folder)
	return repo, env
}

func makeCLIReferencesReadOnly(t *testing.T, folder string) {
	t.Helper()
	references := filepath.Join(folder, "references")
	t.Cleanup(func() {
		if err := os.Chmod(references, 0o755); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Error(err)
		}
	})
	if err := os.Chmod(references, 0o555); err != nil {
		t.Fatal(err)
	}
}

func assertCommandRecoversForward(t *testing.T, repo string, env commandEnvironment, command string, extra ...string) {
	t.Helper()
	before := historySnapshot(t, repo)
	revision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	folderRel := "docs/specs/" + implementTestSlug
	folder := filepath.Join(repo, filepath.FromSlash(folderRel))
	recordRel := "docs/history/specs/" + implementTestSlug + ".md"
	args := []string{command, implementTestSlug}
	if command == "cancel" {
		args = append(args, "--reason", "no consumer")
	}
	args = append(args, extra...)
	var stdout, stderr bytes.Buffer
	code := runWithContext(t.Context(), args, &stdout, &stderr, env)
	lines := strings.Split(strings.TrimSuffix(stderr.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected exactly two stderr lines: %q", stderr.String())
	}
	prefix := fmt.Sprintf("roundfix: %s wrote %s but could not remove the Spec folder %s: ", command, recordRel, folderRel)
	cause, ok := strings.CutPrefix(lines[0], prefix)
	if !ok || !strings.Contains(cause, filepath.Join(folder, "references")) || !strings.Contains(cause, "permission denied") {
		t.Fatalf("expected real removal permission failure: %q", lines[0])
	}
	quoted := "'" + strings.ReplaceAll(folder, "'", "'\\''") + "'"
	finish := "chmod -R u+w " + quoted + " && rm -rf " + quoted
	wantErr := fmt.Sprintf("roundfix: %s wrote %s but could not remove the Spec folder %s: %s\nroundfix: finish the removal with: %s\n", command, recordRel, folderRel, cause, finish)
	cancelAssertTranscript(t, code, stdout.String(), stderr.String(), 1, "", wantErr)
	recordPath := filepath.Join(repo, filepath.FromSlash(recordRel))
	record := archiveRecordAt(t, recordPath)
	disposition := spec.ArchivePass
	if command == "cancel" {
		disposition = spec.ArchiveCancelled
	}
	if record.Disposition != disposition || record.SourceRevision != revision || record.Source != folderRel || !reflect.DeepEqual(record.Promoted, []string{"docs/references/note.md"}) || !reflect.DeepEqual(record.Relinked, []string{"docs/backlog/followup.md"}) {
		t.Fatalf("record=%+v", record)
	}
	if len(extra) > 0 && record.PullRequest != "263" {
		t.Fatalf("pull request not recorded before recovery: %+v", record)
	}
	expected := make(map[string][]byte)
	for path, data := range before {
		if !strings.HasPrefix(path, folderRel+"/") {
			expected[path] = data
		}
	}
	expected[recordRel] = []byte(mustRead(t, recordPath))
	expected["docs/references/note.md"] = []byte("Promoted note.\n")
	expected["docs/backlog/followup.md"] = []byte("[Note](../references/note.md)\n[Spec](../history/specs/" + implementTestSlug + ".md)\n")
	if command == "cancel" {
		if record.CancelReason != "no consumer" || !reflect.DeepEqual(record.CompletedTasks, []string{"task_01"}) || !reflect.DeepEqual(record.PendingTasks, []string{"task_02"}) || !reflect.DeepEqual(record.ReturnedSources, []string{"docs/backlog/entry.md", "docs/findings/finding.md"}) {
			t.Fatalf("cancel record=%+v", record)
		}
		expected["docs/backlog/entry.md"] = []byte("---\nstatus: open\nspec: null\n---\n\n# Entry\n\nIntent.\n")
		expected["docs/findings/finding.md"] = []byte("---\nstatus: pending\nupdated_at: " + record.Archived + "\n---\n\n# Finding\n\nEvidence.\n\n## Addendum — " + record.Archived + " — Spec " + implementTestSlug + " cancelled\n\nSpec " + implementTestSlug + " was cancelled on " + record.Archived + " before it implemented this finding: no consumer\n")
	}
	for path, want := range expected {
		if got := mustRead(t, filepath.Join(repo, filepath.FromSlash(path))); got != string(want) {
			t.Fatalf("write %s not preserved: got %q, want %q", path, got, want)
		}
	}
	assertPathExists(t, filepath.Join(folder, "references/note.md"))
	assertPathMissing(t, filepath.Join(folder, "_prd.md"))
	assertNoRunDatabase(t, env.homeDir)
	printed, ok := strings.CutPrefix(lines[1], "roundfix: finish the removal with: ")
	if !ok {
		t.Fatalf("missing finishing command: %q", lines[1])
	}
	runFolderRemovalFinish(t, printed)
	assertPathMissing(t, folder)
	if got := historySnapshot(t, repo); !reflect.DeepEqual(got, expected) {
		t.Fatalf("finished tree differs from successful %s tree:\ngot=%q\nwant=%q", command, got, expected)
	}
	if got := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); got != revision {
		t.Fatalf("command committed: %s -> %s", revision, got)
	}
}

func runFolderRemovalFinish(t *testing.T, command string) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "sh", "-c", command)
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + t.TempDir()}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("finishing command failed: %v\n%s", err, output)
	}
}
