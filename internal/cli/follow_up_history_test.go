package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/queueforecast"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type followUpReadGit func(context.Context, string, ...string) (string, error)

func (run followUpReadGit) RunGit(ctx context.Context, dir string, args ...string) (string, error) {
	return run(ctx, dir, args...)
}

func followUpCommitAt(t *testing.T, repo string, at time.Time, message string) {
	t.Helper()
	gitImplement(t, repo, "add", ".")
	args := append(gitConfigArgsForTest(), "commit", "--allow-empty", "-m", message)
	cmd := exec.CommandContext(t.Context(), "git", args...)
	cmd.Dir = repo
	cmd.Env = append(isolatedGitEnvForTest(), "GIT_COMMITTER_DATE="+at.UTC().Format(time.RFC3339), "GIT_AUTHOR_DATE="+at.UTC().Format(time.RFC3339))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("commit: %v\n%s", err, out)
	}
}

func newFollowUpReadInput(t *testing.T, now time.Time) followUpQueueInput {
	t.Helper()
	repo, home := t.TempDir(), t.TempDir()
	gittest.InitRepo(t, repo, "--initial-branch=main")
	var err error
	repo, err = filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	home, err = filepath.EvalSymlinks(home)
	if err != nil {
		t.Fatal(err)
	}
	followUpWrite(t, repo, "README.md", "# Widget\n")
	mustMkdir(t, filepath.Join(repo, "docs/specs"))
	followUpCommitAt(t, repo, now.AddDate(0, 0, -60), "feat: seed widget")
	gitImplement(t, repo, "remote", "add", "origin", "git@github.com:owner/widget.git")
	return followUpQueueInput{Context: t.Context(), Loaded: roundconfig.Loaded{GitRoot: repo, HomeDir: home}, SpecRoot: filepath.Join(repo, "docs/specs"), ArchiveRoot: filepath.Join(repo, "docs/history/specs"), Now: now, Zone: time.FixedZone("BRT", -3*60*60), Git: preflight.ExecGitRunner{}}
}

func followUpArchiveFixture(t *testing.T, in followUpQueueInput, slug string, disposition spec.ArchiveDisposition) {
	t.Helper()
	record := spec.ArchiveRecord{Spec: slug, Title: "Widget", Created: in.Now.AddDate(0, 0, -2).Format("2006-01-02"), Archived: in.Now.Format("2006-01-02"), Disposition: disposition, Source: "docs/specs/" + slug, SourceRevision: strings.Repeat("a", 40), Outcome: "Widget delivered."}
	if disposition == spec.ArchiveCancelled {
		record.CancelReason = "scope changed"
	}
	data, err := spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := spec.ParseArchiveRecord(data); err != nil {
		t.Fatal(err)
	}
	followUpWrite(t, in.ArchiveRoot, slug+".md", string(data))
}

func TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 17, 0, 0, 0, time.UTC)
	names := []string{"start_is_the_first_implement_run", "end_is_the_archive_records_commit_time", "plan_is_the_prds_commit_time", "a_cancelled_record_is_not_delivered", "a_spec_without_a_run_is_left_out", "before_the_window_is_not_read", "one_git_log", "no_default_branch_is_unavailable", "a_failed_log_is_unavailable", "an_unparseable_git_log", "a_malformed_commit_time", "an_unparseable_record_is_not_delivered", "the_first_commit_for_a_path_counts", "absolute_roots", "an_edited_working_tree_record_is_not_read", "a_removed_working_tree_record_is_not_read", "a_committed_cancellation_wins_over_the_working_tree"}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			in := newFollowUpReadInput(t, now)
			plan, start, end := now.Add(-8*time.Hour), now.Add(-5*time.Hour), now.Add(-time.Hour)
			if name == "before_the_window_is_not_read" {
				plan = now.AddDate(0, 0, -45)
			}
			followUpWrite(t, in.SpecRoot, "0001-widget/_prd.md", "---\nstatus: active\n---\n# Widget\n")
			followUpCommitAt(t, in.Loaded.GitRoot, plan, "docs: plan widget")
			disposition := spec.ArchiveNoQA
			if name == "a_cancelled_record_is_not_delivered" || name == "a_committed_cancellation_wins_over_the_working_tree" {
				disposition = spec.ArchiveCancelled
			}
			followUpArchiveFixture(t, in, "0001-widget", disposition)
			followUpArchiveFixture(t, in, "0002-review", spec.ArchiveNoQA)
			followUpArchiveFixture(t, in, "0004-no-run", spec.ArchiveNoQA)
			if name == "an_unparseable_record_is_not_delivered" {
				followUpWrite(t, in.ArchiveRoot, "0001-widget.md", "invalid record")
			}
			followUpCommitAt(t, in.Loaded.GitRoot, end, "docs: archive widget")
			recordPath := filepath.Join(in.ArchiveRoot, "0001-widget.md")
			switch name {
			case "an_edited_working_tree_record_is_not_read":
				followUpWrite(t, in.ArchiveRoot, "0001-widget.md", "invalid record")
			case "a_removed_working_tree_record_is_not_read":
				if err := os.Remove(recordPath); err != nil {
					t.Fatal(err)
				}
			case "a_committed_cancellation_wins_over_the_working_tree":
				followUpArchiveFixture(t, in, "0001-widget", spec.ArchiveNoQA)
			}
			runs := []store.Run{
				{Kind: store.KindImplement, SpecSlug: "0001-widget", CreatedAt: start.Add(time.Hour)},
				{Kind: store.KindImplement, SpecSlug: "0001-widget", CreatedAt: start},
				{Kind: store.KindWatch, SpecSlug: "0001-widget", CreatedAt: start.Add(-time.Hour)},
				{Kind: store.KindWatch, SpecSlug: "0002-review", CreatedAt: start},
				{Kind: store.KindImplement, SpecSlug: "0003-unfinished", CreatedAt: start.Add(2 * time.Hour)},
			}
			input := deliveryHistoryInput{Runs: runs, Checkout: in.Loaded.GitRoot, Branch: "main", SpecRoot: "docs/specs", ArchiveRoot: "docs/history/specs", RetentionDays: 30, Now: now}
			if name == "absolute_roots" {
				input.SpecRoot = in.SpecRoot
				input.ArchiveRoot = in.ArchiveRoot
			}
			count := 0
			input.Git = followUpReadGit(func(ctx context.Context, dir string, args ...string) (string, error) {
				if name == "no_default_branch_is_unavailable" {
					return "", errors.New("no refs")
				}
				if args[0] == "log" {
					count++
					want := []string{"log", "refs/heads/main", "--since=" + now.AddDate(0, 0, -31).Format(time.RFC3339), "--no-renames", "--diff-filter=A", "--name-only", "--format=%x00%cI", "--", "docs/history/specs/*.md", "docs/specs/*/_prd.md"}
					if dir != in.Loaded.GitRoot || !reflect.DeepEqual(args, want) {
						t.Fatalf("dir=%s args=%q want=%q", dir, args, want)
					}
					switch name {
					case "a_failed_log_is_unavailable":
						return "", errors.New("history unavailable\nprivate second line")
					case "an_unparseable_git_log":
						return "broken output", nil
					case "a_malformed_commit_time":
						return "\x00invalid\n\ndocs/history/specs/0001-widget.md\n", nil
					case "the_first_commit_for_a_path_counts":
						return "\x00" + end.Format(time.RFC3339) + "\n\ndocs/history/specs/0001-widget.md\n\x00" + end.Add(-time.Hour).Format(time.RFC3339) + "\n\ndocs/history/specs/0001-widget.md\n\x00" + plan.Format(time.RFC3339) + "\n\ndocs/specs/0001-widget/_prd.md\n", nil
					}
				}
				return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
			})
			before := followUpSnapshot(t, in.Loaded.GitRoot, in.Loaded.HomeDir)
			history, source := readDeliveryHistory(t.Context(), input)
			if name == "no_default_branch_is_unavailable" || name == "a_failed_log_is_unavailable" || (name == "an_unparseable_git_log" || name == "a_malformed_commit_time") {
				if len(history) != 0 || source.Name != "delivery-history" || source.State != "unavailable" || source.Count != 0 || source.Detail == "" || strings.Contains(source.Detail, "\n") {
					t.Fatalf("history=%+v source=%+v", history, source)
				}
				if name == "a_failed_log_is_unavailable" && source.Detail != "read delivery history: history unavailable" {
					t.Fatalf("detail=%q", source.Detail)
				}
			} else {
				wantEnd := end
				if disposition == spec.ArchiveCancelled || name == "an_unparseable_record_is_not_delivered" {
					wantEnd = time.Time{}
				}
				wantPlan := plan
				if name == "before_the_window_is_not_read" {
					wantPlan = time.Time{}
				}
				want := []queueforecast.SpecTimes{{Spec: "0001-widget", Plan: wantPlan, Start: start, End: wantEnd}, {Spec: "0003-unfinished", Start: start.Add(2 * time.Hour)}}
				if !reflect.DeepEqual(history, want) {
					t.Fatalf("history=%+v want=%+v", history, want)
				}
				delivered := 1
				if wantEnd.IsZero() {
					delivered = 0
				}
				if source != followUpSourceResult("delivery-history", delivered, nil) {
					t.Fatalf("source=%+v", source)
				}
			}
			wantCalls := 1
			if name == "no_default_branch_is_unavailable" {
				wantCalls = 0
			}
			if count != wantCalls {
				t.Fatalf("log calls=%d want=%d", count, wantCalls)
			}
			if after := followUpSnapshot(t, in.Loaded.GitRoot, in.Loaded.HomeDir); !reflect.DeepEqual(before, after) {
				t.Fatal("history read changed its inputs")
			}
		})
	}
}

func TestTheArchiveRecordsAreReadWithAFixedNumberOfGitProcesses(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 16, 17, 0, 0, 0, time.UTC)
	calls := map[int]int{}
	for _, records := range []int{1, 6} {
		in := newFollowUpReadInput(t, now)
		runs := []store.Run{}
		for i := range records {
			slug := fmt.Sprintf("%04d-widget", i+1)
			followUpArchiveFixture(t, in, slug, spec.ArchiveNoQA)
			runs = append(runs, store.Run{Kind: store.KindImplement, SpecSlug: slug, CreatedAt: now.Add(-5 * time.Hour)})
		}
		followUpCommitAt(t, in.Loaded.GitRoot, now.Add(-time.Hour), "docs: archive widgets")
		count := 0
		input := deliveryHistoryInput{Runs: runs, Checkout: in.Loaded.GitRoot, Branch: "main", SpecRoot: "docs/specs", ArchiveRoot: "docs/history/specs", RetentionDays: 30, Now: now, Git: followUpReadGit(func(ctx context.Context, dir string, args ...string) (string, error) {
			count++
			return (preflight.ExecGitRunner{}).RunGit(ctx, dir, args...)
		})}
		history, source := readDeliveryHistory(t.Context(), input)
		if len(history) != records || source != followUpSourceResult("delivery-history", records, nil) {
			t.Fatalf("records=%d history=%+v source=%+v", records, history, source)
		}
		calls[records] = count
	}
	if calls[1] != calls[6] {
		t.Fatalf("Git processes grow with the archive records: %v", calls)
	}
}

func TestAnAddedFilesDiffIsReadByItsHunkCounts(t *testing.T) {
	t.Parallel()
	raw := "diff --git a/r/one.md b/r/one.md\nnew file mode 100644\nindex 0000000..1111111\n--- /dev/null\n+++ b/r/one.md\n@@ -0,0 +1,3 @@\n+diff --git a/r/x.md b/r/x.md\n++++ b/r/x.md\n+@@ -0,0 +1 @@\ndiff --git a/r/two.md b/r/two.md\nnew file mode 100644\n--- /dev/null\n+++ b/r/two.md\n@@ -0,0 +1 @@\n+no newline\n\\ No newline at end of file\ndiff --git a/r/short.md b/r/short.md\n--- /dev/null\n+++ b/r/short.md\n@@ -0,0 +1,2 @@\n+only one line"
	want := map[string][]byte{"r/one.md": []byte("diff --git a/r/x.md b/r/x.md\n+++ b/r/x.md\n@@ -0,0 +1 @@\n"), "r/two.md": []byte("no newline")}
	if got := parseAddedFilesDiff(raw); !reflect.DeepEqual(got, want) {
		t.Fatalf("records=%q", got)
	}
}
