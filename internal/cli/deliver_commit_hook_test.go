package cli

// Suite: queue commits through real Git and repository hooks.
// Invariant: bounded subjects, message-independent proof, and actionable refusals.
// Boundary IN: disposable repositories, homes, Git configuration and hook scripts.
// Boundary OUT: no forge, real user configuration, or network.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestPullRequestRecordSubjectNamesTheSpecNumber(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, slug, number, identifier string }{
		{"four digits", "0088-example", "280", "0088 "},
		{"five digits", "12345-example", "7", "12345 "},
		{"ten digits", "1234567890-example", "1234567890", "1234567890 "},
		{"eleven digits", "12345678901-example", "7", ""},
		{"digits alone", "0088", "7", ""},
		{"without digits", "example", "7", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := "docs: record pull request #" + tc.number + " in the " + tc.identifier + "Archive Record"
			if got := pullRequestRecordSubject(tc.slug, tc.number); got != want || len(got) > 70 {
				t.Fatalf("subject=%q (%d bytes), want %q, at most 70 bytes", got, len(got), want)
			}
		})
	}
}

// Reuse the archived item fixture, changing its record identity before the candidate.
func queueRecordHookFixture(t *testing.T, slug string) (*commandDeliveryWorkflow, string, string, string) {
	t.Helper()
	w, repo, head, path := pullRequestRecordFixture(t)
	if slug != implementTestSlug {
		record := archiveRecordAt(t, filepath.Join(repo, path))
		record.Spec = slug
		record.Source = strings.ReplaceAll(record.Source, implementTestSlug, slug)
		record.QAReport = strings.ReplaceAll(record.QAReport, implementTestSlug, slug)
		data, err := spec.RenderArchiveRecord(record)
		if err != nil {
			t.Fatal(err)
		}
		newPath := strings.ReplaceAll(path, implementTestSlug, slug)
		mustWrite(t, filepath.Join(repo, newPath), string(data))
		if err := os.Remove(filepath.Join(repo, path)); err != nil {
			t.Fatal(err)
		}
		gittest.Run(t, repo, "add", "-A")
		gittest.Run(t, repo, "commit", "-m", "docs: prepare record fixture")
		path = newPath
		head = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	}
	return w, repo, head, path
}

func TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	slug := "0088-" + strings.Repeat("x", 49)
	w, repo, candidate, _ := queueRecordHookFixture(t, slug)
	writeScriptFixture(t, filepath.Join(repo, ".git", "hooks", "commit-msg"), `IFS= read -r header < "$1"
if [ "${#header}" -gt 100 ]; then
 echo 'header must not be longer than 100 characters' >&2
 exit 1
fi
`)
	result, err := w.RecordPullRequest(commandContextForTest(t, t.Context()), repo, slug, candidate, "280")
	if err != nil || !result.Exact || result.Head == candidate || result.Head == "" {
		t.Fatalf("record=%+v err=%v", result, err)
	}
	want := "docs: record pull request #280 in the 0088 Archive Record"
	if got := strings.TrimSpace(gittest.Run(t, repo, "show", "-s", "--format=%s", result.Head)); got != want || len(got) != 57 {
		t.Fatalf("subject=%q", got)
	}
	if dirty := gittest.Run(t, repo, "status", "--porcelain"); dirty != "" {
		t.Fatalf("dirty=%s", dirty)
	}
}

func TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	slug := "0088-" + strings.Repeat("x", 49)
	for _, tc := range []struct{ name, subject string }{
		{"old long subject", "docs: record pull request #280 in the " + slug + " Archive Record"},
		{"short subject", "docs: record pull request #280 in the 0088 Archive Record"},
		{"wip", "wip"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, repo, candidate, path := queueRecordHookFixture(t, slug)
			if _, err := spec.RecordArchivePullRequest(filepath.Join(repo, path), "280"); err != nil {
				t.Fatal(err)
			}
			gittest.Run(t, repo, "add", "--", path)
			gittest.Run(t, repo, "commit", "-m", tc.subject)
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			result, err := w.RecordPullRequest(commandContextForTest(t, t.Context()), repo, slug, candidate, "280")
			if err != nil || !result.Exact || result.Head != head {
				t.Fatalf("adopt=%+v err=%v", result, err)
			}
			if got := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")); got != head {
				t.Fatalf("head=%s want=%s", got, head)
			}
			if count := strings.TrimSpace(gittest.Run(t, repo, "rev-list", "--count", candidate+"..HEAD")); count != "1" {
				t.Fatalf("commits=%s", count)
			}
		})
	}
}

func TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	for _, tc := range []struct {
		name, hook, output string
		both               bool
	}{
		{"record named hook", "commit-msg", "husky - commit-msg hook exited with code 1 (error)", false},
		// The runner repeats the exit error as its detail when the hook is silent.
		{"record silent hook", "commit-msg", "exit status 1", false},
		{"record two hooks", "pre-commit", "exit status 1", true},
		{"archive subject names another hook", "commit-msg", "exit status 1", false},
		{"derived merge", "commit-msg", "exit status 1", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			slug := implementTestSlug
			var w *commandDeliveryWorkflow
			var repo, candidate, path, subject, wrap string
			switch tc.name {
			case "derived merge":
				f := newConflictFixture(t, false, "printf 'regenerated\\n' > derived.txt", "")
				w, repo, candidate = f.workflow, f.repo, f.head
				subject = "chore: merge default branch and regenerate derived paths"
				wrap = "commit derived merge: "
			case "archive subject names another hook":
				home, r := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
				writeArchiveQAReport(t, r, spec.VerdictPass)
				slug = "0088-pre-commit-example"
				old := filepath.Join(r, "docs/specs", implementTestSlug)
				next := filepath.Join(r, "docs/specs", slug)
				if err := os.Rename(old, next); err != nil {
					t.Fatal(err)
				}
				if err := filepath.WalkDir(next, func(p string, d os.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() {
						mustWrite(t, p, strings.ReplaceAll(mustRead(t, p), implementTestSlug, slug))
					}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				commitArchiveFixture(t)
				gittest.PersistIdentity(t, r)
				setCommandEnvForTest(t, cliTestHelperEnv, "1")
				w, repo = retirementWorkflow(r, home), r
				candidate = strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
				subject = "docs: archive " + slug
				wrap = "commit archived Spec: "
			default:
				w, repo, candidate, path = queueRecordHookFixture(t, slug)
				subject = "docs: record pull request #7 in the 0001 Archive Record"
				wrap = "commit pull request record: "
			}
			script := "exit 1\n"
			if tc.name == "record named hook" {
				script = "echo '" + tc.output + "' >&2\nexit 1\n"
			}
			hookPath := filepath.Join(repo, ".git", "hooks", tc.hook)
			writeScriptFixture(t, hookPath, script)
			if tc.both {
				writeScriptFixture(t, filepath.Join(repo, ".git", "hooks", "commit-msg"), "exit 0\n")
			}
			ctx := commandContextForTest(t, t.Context())
			var err error
			switch tc.name {
			case "derived merge":
				_, err = w.ResolveConflict(ctx, repo, slug, candidate)
			case "archive subject names another hook":
				_, err = w.Archive(ctx, repo, slug, candidate)
			default:
				_, err = w.RecordPullRequest(ctx, repo, slug, candidate, "7")
			}
			hook := tc.hook
			if tc.both {
				hook = "pre-commit or commit-msg"
			}
			recovery := "the change stays staged in " + repo + "; commit it there with a subject the hook accepts, then run roundfix deliver retry " + slug
			if tc.name == "derived merge" {
				recovery = "the merge was aborted; merge the default branch into the item branch in " + repo + " by hand, commit it with a message the hook accepts, then run roundfix deliver retry " + slug
			}
			want := fmt.Sprintf("%sthe repository %s hook refused the commit %q (exit status 1): %s; %s", wrap, hook, subject, tc.output, recovery)
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v\nwant=%s", err, want)
			}
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
				t.Fatalf("exit error lost: %v", err)
			}
			item := store.DeliveryQueueItem{SpecSlug: slug, Blocker: delivery.BlockerDeliveryError + ": record pull request: " + err.Error()}
			if got := delivery.ClassifyPark(store.DeliveryQueue{}, item); got.Class != delivery.ParkClassEnvironment {
				t.Fatalf("park=%+v", got)
			}
			if tc.name == "derived merge" {
				assertConflictAborted(t, conflictFixture{workflow: w, repo: repo, head: candidate})
			} else if path != "" {
				if staged := strings.TrimSpace(gittest.Run(t, repo, "diff", "--cached", "--name-only")); staged != path {
					t.Fatalf("staged=%q", staged)
				}
				// Amend the fixture hook to accept the operator's conforming commit.
				mustWrite(t, hookPath+scriptFixtureSuffix, "exit 0\n")
				gittest.Run(t, repo, "commit", "-m", "docs: record pull request #7")
				hand := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
				adopted, e := w.RecordPullRequest(ctx, repo, slug, candidate, "7")
				if e != nil || !adopted.Exact || adopted.Head != hand {
					t.Fatalf("adopt=%+v err=%v", adopted, e)
				}
			}
		})
	}
}

func TestDeliveryQueueCommitWithoutAHookKeepsTheGitError(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	w, repo, candidate, path := queueRecordHookFixture(t, implementTestSlug)
	// Prepare an interrupted write so the stale lock reaches staging/commit directly.
	if _, err := spec.RecordArchivePullRequest(filepath.Join(repo, path), "7"); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, repo, "add", "--", path)
	// Install the lock only at commit: earlier Git operations also require the index.
	w.git = queueCommitLockRunner{GitRunner: w.git, repo: repo}
	_, err := w.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
	if err == nil || !strings.HasPrefix(err.Error(), "commit pull request record: git commit -m ") || strings.Contains(err.Error(), "hook refused") {
		t.Fatalf("err=%v", err)
	}
}

// Inject a real stale lock at the commit boundary, after staging has succeeded.
type queueCommitLockRunner struct {
	preflight.GitRunner
	repo string
}

func (runner queueCommitLockRunner) RunGit(ctx context.Context, dir string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "commit" {
		if err := os.WriteFile(filepath.Join(runner.repo, ".git", "index.lock"), nil, 0600); err != nil {
			return "", err
		}
	}
	return runner.GitRunner.RunGit(ctx, dir, args...)
}

// Run real-Git tests in a child with its own home and Git configuration, since
// the production runner and hook discovery intentionally inherit their environment.
func queueCommitTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_QUEUE_COMMIT_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	ctx, cancel := context.WithCancel(t.Context())
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	home := t.TempDir()
	config := filepath.Join(home, "gitconfig")
	mustWrite(t, config, "")
	cmd.Env = withEnvValue(gittest.IsolatedEnv(), "HOME", home)
	cmd.Env = withEnvValue(cmd.Env, "GIT_CONFIG_GLOBAL", config)
	cmd.Env = withEnvValue(cmd.Env, "GIT_CONFIG_SYSTEM", config)
	cmd.Env = withEnvValue(cmd.Env, key, t.Name())
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	waited := false
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil && !waited {
			_ = cmd.Wait()
		}
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("isolated queue commit test: %v\n%s", err, output.String())
	}
	return false
}

func TestDeliveryQueueCommitHookOutputIsBounded(t *testing.T) {
	t.Parallel()
	if !queueCommitTestHome(t) {
		return
	}
	for _, tc := range []struct{ name, output, want string }{
		{"whitespace", "  finding\n\t more  detail  ", "finding more detail"},
		{"exact byte limit", strings.Repeat("x", 510) + "é", strings.Repeat("x", 510) + "é"},
		{"cut inside rune", strings.Repeat("x", 511) + "é extra", strings.Repeat("x", 511) + " ..."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, repo, candidate, _ := queueRecordHookFixture(t, implementTestSlug)
			writeScriptFixture(t, filepath.Join(repo, ".git", "hooks", "commit-msg"), "printf '%s' '"+tc.output+"' >&2\nexit 1\n")
			_, err := w.RecordPullRequest(commandContextForTest(t, t.Context()), repo, implementTestSlug, candidate, "7")
			want := fmt.Sprintf("commit pull request record: the repository commit-msg hook refused the commit %q (exit status 1): %s; the change stays staged in %s; commit it there with a subject the hook accepts, then run roundfix deliver retry %s", "docs: record pull request #7 in the 0001 Archive Record", tc.want, repo, implementTestSlug)
			if err == nil || err.Error() != want || !utf8.ValidString(err.Error()) {
				t.Fatalf("error=%v\nwant=%s", err, want)
			}
		})
	}
}
