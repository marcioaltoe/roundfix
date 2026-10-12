package cli

// Suite: archived review context and Run causes through their public commands.
// Boundary IN: main-only clones, a local bare forge and a disposable Run store.
// Boundary OUT: network, real providers and daemon settlement.

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func cloneCLIRecovery(t *testing.T, repo, branch string, r spec.ArchiveRecord, number string) string {
	t.Helper()
	r.PullRequest = number
	if number == "unreachable" {
		r.PullRequest = "7"
	}
	p := writeReaderRecord(t, repo, r)
	gittest.Run(t, repo, "add", p)
	gittest.Run(t, repo, "commit", "--amend", "-m", "feat: deliver archived Spec (#7)")
	clone, _ := gittest.PullRequestClone(t, repo, branch, r.SourceRevision)
	// Both command readers hold Project Config and must honor this remote.
	gittest.Run(t, clone, "remote", "rename", "origin", "delivery")
	if number == "unreachable" {
		gittest.Run(t, clone, "remote", "set-url", "delivery", filepath.Join(t.TempDir(), "absent.git"))
	}
	gittest.Run(t, clone, "remote", "add", "origin", filepath.Join(t.TempDir(), "unreachable.git"))
	canonical, err := filepath.EvalSymlinks(clone)
	if err != nil {
		t.Fatal(err)
	}
	return canonical
}

func TestReviewContextRecoversASquashDroppedRevision(t *testing.T) {
	t.Parallel()
	for _, number := range []string{"7", "", "unreachable"} {
		t.Run(map[string]string{"7": "record", "": "adding subject", "unreachable": "unreachable forge"}[number], func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
			f := newReviewCommandFixture(t, "codex", runner)
			r := squashReaderRecord(t, f.repository, readerRecordSlug, false, false)
			clone := cloneCLIRecovery(t, f.repository, "feature/review", r, number)
			f.repository = clone
			f.headCommit = strings.TrimSpace(gittest.Run(t, clone, "rev-parse", "HEAD"))
			writeReviewCommandConfig(t, clone, "codex", f.artifactDir)
			file := filepath.Join(clone, ".roundfixrc.yml")
			mustWrite(t, file, "defaults:\n  artifact_dir: "+f.artifactDir+"\npre_pr_review:\n  provider: codex\nwatch:\n  push_remote: delivery\n")
			setCommandEnvironmentForTest(t, f.homeDir, clone)
			code, record, diag := f.run(t)
			if number == "unreachable" {
				if code != exitOK || runner.preparedCalls != 1 || !slices.Contains(record.SkippedSpecs, readerRecordSlug) || slices.Contains(record.Specs, readerRecordSlug) || strings.Contains(runner.request.Prompt, "Source technical design") {
					t.Fatalf("unavailable review exit=%d record=%+v stderr=%s", code, record, diag)
				}
				return
			}
			if code != exitOK || runner.preparedCalls != 1 || !slices.Contains(record.Specs, readerRecordSlug) || slices.Contains(record.SkippedSpecs, readerRecordSlug) || !strings.Contains(runner.request.Prompt, "Keep archive context from Git") || !strings.Contains(runner.request.Prompt, "Source technical design") {
				t.Fatalf("exit=%d record=%+v stderr=%s prompt=%s", code, record, diag, runner.request.Prompt)
			}
			gittest.Run(t, clone, "cat-file", "-e", r.SourceRevision+"^{commit}")
		})
	}
}

func TestRunCausesRecoverASquashDroppedRevision(t *testing.T) {
	t.Parallel()
	for _, number := range []string{"7", "", "unreachable"} {
		t.Run(map[string]string{"7": "record", "": "adding subject", "unreachable": "unreachable forge"}[number], func(t *testing.T) {
			env, r := archivedCausesCLI(t, false)
			clone := cloneCLIRecovery(t, env.workDir, "feat/archive-record-source", r, number)
			db, err := sql.Open("sqlite", store.DatabasePath(env.homeDir))
			if err != nil {
				t.Fatal(err)
			}
			canonical, err := filepath.EvalSymlinks(clone)
			if err != nil {
				t.Fatal(err)
			}
			_, updateErr := db.Exec("UPDATE runs SET git_root = ?, repository_root = ?", canonical, canonical)
			closeErr := db.Close()
			if updateErr != nil {
				t.Fatal(updateErr)
			}
			if closeErr != nil {
				t.Fatal(closeErr)
			}
			env.workDir = clone
			mustWrite(t, filepath.Join(clone, ".roundfixrc.yml"), "watch:\n  push_remote: delivery\n")
			code, out, diag := runCausesCLI(t, env, "--format", "json")
			var report runcause.Report
			if err := json.Unmarshal([]byte(out), &report); err != nil {
				t.Fatalf("JSON=%s stderr=%s err=%v", out, diag, err)
			}
			if number == "unreachable" {
				if code != exitOK || diag != "" || len(report.ArchivedSpecs) != 1 || report.ArchivedSpecs[0].SourceAvailable || len(report.Items) != 3 || len(report.Tasks) != 4 || report.Tasks[3].Corrective != nil {
					t.Fatalf("unavailable causes exit=%d report=%+v stderr=%s", code, report, diag)
				}
				return
			}
			if code != exitOK || diag != "" || len(report.ArchivedSpecs) != 1 || !report.ArchivedSpecs[0].SourceAvailable || len(report.Items) != 4 || len(report.Tasks) != 4 || report.Tasks[3].Corrective == nil || !*report.Tasks[3].Corrective {
				t.Fatalf("exit=%d report=%+v stderr=%s", code, report, diag)
			}
			gittest.Run(t, clone, "cat-file", "-e", r.SourceRevision+"^{commit}")
		})
	}
}
