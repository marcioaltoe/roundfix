package cli

// Suite: Archive Command pull request provenance and squash warnings.
// Boundary: full CLI transcripts against disposable repositories and a local bare origin.

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func archivePullRequestFixture(t *testing.T) (string, string, string) {
	t.Helper()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	gittest.Run(t, repo, "branch", "-m", "fix/archive-fixture")
	remote := t.TempDir()
	gittest.InitRepo(t, remote, "--bare", "--initial-branch=main")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "push", "origin", "main")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "docs: completed branch fixture")
	revision := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	return repo, filepath.Join(repo, "docs/history/specs", implementTestSlug+".md"), revision
}

func assertArchiveTranscript(t *testing.T, args []string, wantOut, wantErr string, wantCode int) {
	t.Helper()
	var out, errout bytes.Buffer
	code := runCLIContext(t, t.Context(), args, &out, &errout)
	if out.String() != wantOut || errout.String() != wantErr || code != wantCode {
		t.Fatalf("transcript: exit=%d stdout=%q stderr=%q\nwant: exit=%d stdout=%q stderr=%q", code, out.String(), errout.String(), wantCode, wantOut, wantErr)
	}
}

func archivePullRequestConfirmation(t *testing.T, repo, revision string, number bool) string {
	t.Helper()
	files := snapshotDirectoryFiles(t, filepath.Join(repo, "docs/specs", implementTestSlug))
	size := 0
	for _, content := range files {
		size += len(content)
	}
	suffix := ""
	if number {
		suffix = "; pull request 7"
	}
	return fmt.Sprintf("archived %s -> docs/history/specs/%s.md; removed %d file(s) (%d bytes) kept in Git at %.12s%s\n", implementTestSlug, implementTestSlug, len(files), size, revision, suffix)
}

func archivePullRequestPreflight(reason string) string {
	return "Preflight failed\n\nReason:\n  " + reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix archive --help' for usage.\n"
}

func archiveBranchRecord(t *testing.T) (string, string) {
	t.Helper()
	repo, path, revision := archivePullRequestFixture(t)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug}, archivePullRequestConfirmation(t, repo, revision, false), fmt.Sprintf("roundfix: warning: source_revision %.12s is not an ancestor of origin/main; a squash merge keeps it only in the pull request's head ref, so record the pull request with --pull-request once it opens\n", revision), 0)
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "docs: archive fixture")
	return repo, path
}

func TestArchiveOnABranchWarnsAboutASquashDroppedRevision(t *testing.T) {
	t.Parallel()
	repo, _, revision := archivePullRequestFixture(t)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug}, archivePullRequestConfirmation(t, repo, revision, false), fmt.Sprintf("roundfix: warning: source_revision %.12s is not an ancestor of origin/main; a squash merge keeps it only in the pull request's head ref, so record the pull request with --pull-request once it opens\n", revision), 0)
}

func TestArchiveWithAPullRequestNamesItsHeadRef(t *testing.T) {
	t.Parallel()
	repo, path, revision := archivePullRequestFixture(t)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--pull-request", "7"}, archivePullRequestConfirmation(t, repo, revision, true), fmt.Sprintf("roundfix: warning: source_revision %.12s is not an ancestor of origin/main; readers recover it from refs/pull/7/head\n", revision), 0)
	if got := archiveRecordAt(t, path).PullRequest; got != "7" {
		t.Fatalf("pull_request=%q", got)
	}
}

func TestArchiveRecordsAPullRequestInTheBranchRecord(t *testing.T) {
	t.Parallel()
	repo, path := archiveBranchRecord(t)
	before := snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--pull-request", "7"}, fmt.Sprintf("recorded pull request 7 in docs/history/specs/%s.md\n", implementTestSlug), "", 0)
	if got := archiveRecordAt(t, path).PullRequest; got != "7" {
		t.Fatalf("pull_request=%q", got)
	}
	after := snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))
	delete(before, "history/specs/"+implementTestSlug+".md")
	delete(after, "history/specs/"+implementTestSlug+".md")
	if !reflect.DeepEqual(before, after) {
		t.Fatal("record-only form changed other files")
	}
	content := mustRead(t, path)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--pull-request", "7"}, fmt.Sprintf("docs/history/specs/%s.md already records pull request 7\n", implementTestSlug), "", 0)
	if mustRead(t, path) != content {
		t.Fatal("idempotent record changed bytes")
	}
}

func TestArchiveRefusesAPullRequestForARecordOnTheDefaultBranch(t *testing.T) {
	t.Parallel()
	repo, path := archiveBranchRecord(t)
	gittest.Run(t, repo, "push", "origin", "HEAD:main")
	before := mustRead(t, path)
	reason := fmt.Sprintf("Archive Record docs/history/specs/%s.md is already on origin/main; an archived record stays byte-identical", implementTestSlug)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--pull-request", "7"}, "", archivePullRequestPreflight(reason), 2)
	if mustRead(t, path) != before {
		t.Fatal("default-branch record changed bytes")
	}
}

func TestArchiveHelpNamesThePullRequestOption(t *testing.T) {
	t.Parallel()
	const help = `Usage:
  roundfix archive <slug>
  roundfix archive <slug> --qa-override --approval <source> --reason <text>
  roundfix archive <slug> --plan
  roundfix archive <slug> --pull-request <number>

Archives a Spec after verifying either every Task is completed and the newest
QA Report has verdict: pass (or a partial verdict covered only by declared
Unreachable Acceptance), or a recorded supersession exists for a Spec without
a Task Graph. Writes one Archive Record and removes the committed Spec folder.
The removed bytes stay in Git at the record's source_revision. The destination
is docs/history/specs/<slug>.md for the built-in root, otherwise
<spec-root>/_archived/<slug>.md.
archive creates no Run and never pushes.
Refuses, writing nothing, while the Spec folder has uncommitted changes or Git ignores a file under it.

The QA override requires an approval source and reason, still requires every
non-QA Task completed, preserves the QA Task and Reports unchanged, and records
the observed QA outcome and archived HEAD. It is refused when QA already
qualifies for normal archive.

Options:
  --qa-override  Archive despite failed, missing, or unreadable QA evidence
  --approval     Source of the maintainer authority for the override
  --reason       Reason the unmet QA prerequisite is being waived
  --pull-request Store the delivery pull request number in the Archive Record

Exit codes:
  0  archived
  1  archive could not finish after writing; stderr says how to finish
  2  Preflight Validation failed
`
	assertArchiveTranscript(t, []string{"archive", "--help"}, help, "", 0)
}

func TestArchiveOnTheDefaultBranchDoesNotWarn(t *testing.T) {
	t.Parallel()
	for _, originHead := range []bool{true, false} {
		t.Run(fmt.Sprint(originHead), func(t *testing.T) {
			repo, _, revision := archivePullRequestFixture(t)
			gittest.Run(t, repo, "push", "origin", "HEAD:main")
			gittest.Run(t, repo, "checkout", "-B", "main")
			if !originHead {
				gittest.Run(t, repo, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			}
			assertArchiveTranscript(t, []string{"archive", implementTestSlug}, archivePullRequestConfirmation(t, repo, revision, false), "", 0)
		})
	}
}

func TestArchivePullRequestRefusals(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		record bool
		flags  []string
		reason string
	}{
		{"non-integer", false, []string{"--pull-request", "seven"}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"zero", false, []string{"--pull-request", "0"}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"negative", false, []string{"--pull-request", "-7"}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"empty", false, []string{"--pull-request", ""}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"leading zero", false, []string{"--pull-request", "07"}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"too long", false, []string{"--pull-request", "12345678901"}, "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"},
		{"plan", false, []string{"--pull-request", "7", "--plan"}, "--pull-request cannot be combined with --plan"},
		{"different stored number", true, []string{"--pull-request", "8"}, "Archive Record already records pull request 7; cannot replace it with 8"},
		{"unknown default branch", true, []string{"--pull-request", "7"}, "the default branch is unknown; cannot record a pull request"},
		{"active unknown default branch", false, []string{"--pull-request", "7"}, "the default branch is unknown; cannot record a pull request"},
		{"promote", true, []string{"--pull-request", "7", "--promote", "notes.md"}, "record-only --pull-request cannot be combined with --promote, --qa-override, --approval or --reason"},
		{"qa override", true, []string{"--pull-request", "7", "--qa-override", "--approval", "maintainer", "--reason", "approved"}, "record-only --pull-request cannot be combined with --promote, --qa-override, --approval or --reason"},
		{"approval", true, []string{"--pull-request", "7", "--approval", "maintainer"}, "--approval and --reason require --qa-override"},
		{"reason", true, []string{"--pull-request", "7", "--reason", "approved"}, "--approval and --reason require --qa-override"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var repo, path string
			if tc.record {
				repo, path = archiveBranchRecord(t)
			} else {
				repo, path, _ = archivePullRequestFixture(t)
			}
			if tc.name == "different stored number" {
				if _, err := spec.RecordArchivePullRequest(path, "7"); err != nil {
					t.Fatal(err)
				}
			}
			if strings.Contains(tc.name, "unknown default branch") {
				gittest.Run(t, repo, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			}
			before := snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))
			args := append([]string{"archive", implementTestSlug}, tc.flags...)
			assertArchiveTranscript(t, args, "", archivePullRequestPreflight(tc.reason), 2)
			if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, filepath.Join(repo, "docs"))) {
				t.Fatal("refusal changed files")
			}
		})
	}
}

func TestArchiveWithAnUnknownDefaultBranchWarns(t *testing.T) {
	t.Parallel()
	repo, _, revision := archivePullRequestFixture(t)
	gittest.Run(t, repo, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
	assertArchiveTranscript(t, []string{"archive", implementTestSlug}, archivePullRequestConfirmation(t, repo, revision, false), fmt.Sprintf("roundfix: warning: the default branch is unknown, so source_revision %.12s was not checked\n", revision), 0)
}

func TestArchivePullRequestRefusesAnExternalDefaultBranchRecord(t *testing.T) {
	t.Parallel()
	project, branchRecord := archiveBranchRecord(t)
	externalRepo, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	gittest.InitRepo(t, externalRepo, "--initial-branch=main")
	root := filepath.Join(externalRepo, "docs/specs")
	record := filepath.Join(root, "_archived", implementTestSlug+".md")
	mustMkdir(t, filepath.Dir(record))
	before := mustRead(t, branchRecord)
	mustWrite(t, record, before)
	gittest.Run(t, externalRepo, "add", "-A")
	gittest.Run(t, externalRepo, "commit", "-m", "docs: external default record")
	remote := t.TempDir()
	gittest.InitRepo(t, remote, "--bare", "--initial-branch=main")
	gittest.Run(t, externalRepo, "remote", "add", "origin", remote)
	gittest.Run(t, externalRepo, "push", "origin", "main")
	gittest.Run(t, externalRepo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, externalRepo, "checkout", "-b", "fix/external-record")
	mustWrite(t, filepath.Join(project, ".roundfixrc.yml"), fmt.Sprintf("specs:\n  root: %q\n", root))
	reason := fmt.Sprintf("Archive Record docs/specs/_archived/%s.md is already on origin/main; an archived record stays byte-identical", implementTestSlug)
	assertArchiveTranscript(t, []string{"archive", implementTestSlug, "--pull-request", "7"}, "", archivePullRequestPreflight(reason), 2)
	if mustRead(t, record) != before {
		t.Fatal("external default-branch record changed bytes")
	}
}
