package spec

// Suite: source recovery through real Git and a local bare forge.
// Boundary IN: main-only clones and pull request refs on disk.
// Boundary OUT: network, credentials, pushes and settlement.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/testfixture"
)

func recoveryFixture(t *testing.T, number string) (string, string, ArchiveRecord, string) {
	t.Helper()
	source := t.TempDir()
	gittest.InitRepo(t, source, "-b", "main")
	gittest.Run(t, source, "commit", "--allow-empty", "-m", "docs: seed")
	gittest.Run(t, source, "checkout", "-b", "feat/source")
	dir := "docs/specs/0256-example"
	if err := os.MkdirAll(filepath.Join(source, dir), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, dir, "_prd.md"), []byte("# Historical source\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, source, "add", ".")
	gittest.Run(t, source, "commit", "-m", "docs: source")
	r := ArchiveRecord{Spec: "0256-example", Title: "Example", Source: dir, SourceRevision: strings.TrimSpace(gittest.Run(t, source, "rev-parse", "HEAD")), Disposition: ArchivePass, PullRequest: number}
	if err := os.RemoveAll(filepath.Join(source, dir)); err != nil {
		t.Fatal(err)
	}
	recordPath := ArchiveRecordPath("docs/history/specs", r.Spec)
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(source, recordPath)), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, recordPath), content, 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, source, "add", "-A")
	gittest.Run(t, source, "commit", "-m", "docs: archive")
	head := strings.TrimSpace(gittest.Run(t, source, "rev-parse", "HEAD"))
	gittest.Run(t, source, "checkout", "main")
	gittest.Run(t, source, "merge", "--squash", "feat/source")
	gittest.Run(t, source, "commit", "-m", "feat: delivery (#7)")
	repo, forge := gittest.PullRequestClone(t, source, head, r.SourceRevision)
	return repo, forge, r, recordPath
}

func requireRecovered(t *testing.T, repo, remote string, r ArchiveRecord, recordPath string) {
	t.Helper()
	result, err := EnsureSourceRevision(t.Context(), repo, remote, r, "HEAD", recordPath)
	if err != nil || result.Revision != r.SourceRevision || result.PullRequest != "7" || result.From != "pull-request" {
		t.Fatalf("recovery=%+v err=%v", result, err)
	}
	if got := gittest.Run(t, repo, "show", r.SourceRevision+":"+r.Source+"/_prd.md"); got != "# Historical source\n" {
		t.Fatalf("source=%q", got)
	}
	// The local path must work even after the forge becomes unreachable.
	result, err = EnsureSourceRevision(t.Context(), repo, filepath.Join(t.TempDir(), "absent"), r, "HEAD", recordPath)
	if err != nil || result.From != "local" {
		t.Fatalf("local recovery=%+v err=%v", result, err)
	}
}

func TestSourceRevisionRecoveryFetchesThePullRequestHead(t *testing.T) {
	t.Parallel()
	repo, _, r, p := recoveryFixture(t, "7")
	requireRecovered(t, repo, "origin", r, p)
}

func TestSourceRevisionRecoveryReadsTheNumberFromTheAddingCommit(t *testing.T) {
	t.Parallel()
	repo, _, r, p := recoveryFixture(t, "")
	// A later edit must not supply the number: it belongs to the adding commit.
	r.Outcome = "Later record detail."
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, p), content, 0600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, repo, "add", p)
	gittest.Run(t, repo, "commit", "-m", "docs: later change (#8)")
	requireRecovered(t, repo, "origin", r, p)
}

func TestSourceRevisionRecoveryWritesNoRef(t *testing.T) {
	t.Parallel()
	repo, _, r, p := recoveryFixture(t, "7")
	// Even a configured mapping for pull refs must not create a destination ref.
	gittest.Run(t, repo, "config", "--add", "remote.origin.fetch", "+refs/pull/*/head:refs/remotes/origin/pull/*")
	before := gittest.Run(t, repo, "for-each-ref")
	requireRecovered(t, repo, "origin", r, p)
	if after := gittest.Run(t, repo, "for-each-ref"); after != before {
		t.Fatalf("refs changed: before=%s after=%s", before, after)
	}
	if _, err := os.Stat(filepath.Join(repo, ".git", "FETCH_HEAD")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("FETCH_HEAD present: %v", err)
	}
}

func TestSourceRevisionRecoveryLeavesAnUnreachableRevisionUnavailable(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unreachable forge", "wrong pull request head"} {
		t.Run(name, func(t *testing.T) {
			repo, forge, r, p := recoveryFixture(t, "7")
			if name == "unreachable forge" {
				gittest.Run(t, repo, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
			} else {
				gittest.Run(t, forge, "fetch", "--no-write-fetch-head", repo, "main:refs/pull/7/head", "--force")
			}
			result, err := EnsureSourceRevision(t.Context(), repo, "origin", r, "HEAD", p)
			if !errors.Is(err, ErrSourceRevisionUnavailable) || result.From != "" {
				t.Fatalf("recovery=%+v err=%v", result, err)
			}
		})
	}
}

func TestSourceRevisionRecoveryRefusesAMalformedNumber(t *testing.T) {
	t.Parallel()
	repo, _, r, p := recoveryFixture(t, "7")
	before := gittest.Run(t, repo, "count-objects", "-v")
	for _, number := range []string{"0", "07", "-7", "7/head", "7:refs/heads/attack", "10000000000", "7\n"} {
		t.Run(number, func(t *testing.T) {
			r.PullRequest = number
			if _, err := EnsureSourceRevision(t.Context(), repo, "origin", r, "HEAD", p); !errors.Is(err, ErrSourceRevisionUnavailable) {
				t.Fatalf("number=%q err=%v", number, err)
			}
			if after := gittest.Run(t, repo, "count-objects", "-v"); after != before {
				t.Fatal("malformed number fetched objects")
			}
		})
	}
}

func TestPullRequestFromSubject(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ subject, number string }{
		{"feat: delivery (#7)", "7"}, {"feat: delivery (#9999999999)", "9999999999"},
		{"feat: (#0)", ""}, {"feat: (#07)", ""}, {"feat: (#10000000000)", ""}, {"feat: (#7) more", ""}, {"feat: #7", ""},
	} {
		t.Run(tc.subject, func(t *testing.T) {
			n, ok := PullRequestFromSubject(tc.subject)
			if n != tc.number || ok != (tc.number != "") {
				t.Fatalf("number=%q ok=%v", n, ok)
			}
		})
	}
}

func TestSourceRevisionRecoveryHonorsCancellation(t *testing.T) {
	t.Parallel()
	repo, _, r, p := recoveryFixture(t, "7")
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := EnsureSourceRevision(ctx, repo, "origin", r, "HEAD", p); !errors.Is(err, ErrSourceRevisionUnavailable) {
		t.Fatalf("canceled recovery=%v", err)
	}
}

func TestSourceRevisionRecoveryCancelsTheFetchProcess(t *testing.T) {
	t.Parallel()
	// Re-execution gives the fetch a private process environment while the test
	// remains parallel; no PATH or prompt setting is changed in the parent.
	if repo := os.Getenv("ROUNDFIX_TEST_RECOVERY_REPO"); repo != "" {
		r := ArchiveRecord{Spec: "0256-example", Source: "docs/specs/0256-example", SourceRevision: os.Getenv("ROUNDFIX_TEST_RECOVERY_REVISION"), PullRequest: "7"}
		marker := os.Getenv("ROUNDFIX_TEST_RECOVERY_MARKER")
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		_, err := EnsureSourceRevision(ctx, repo, "origin", r, "HEAD", ArchiveRecordPath("docs/history/specs", r.Spec))
		if !errors.Is(err, ErrSourceRevisionUnavailable) || !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Fatalf("fetch error=%v context=%v", err, ctx.Err())
		}
		data, err := os.ReadFile(marker)
		if err != nil {
			t.Fatal(err)
		}
		fields := strings.Fields(string(data))
		if len(fields) != 3 || fields[1] != "0" || fields[2] != "eof" {
			t.Fatalf("process boundary=%q", data)
		}
		if _, err := strconv.Atoi(fields[0]); err != nil {
			t.Fatal(err)
		}
		if err := exec.Command("kill", "-0", fields[0]).Run(); err == nil {
			t.Fatal("canceled fetch process remains alive")
		}

		return
	}
	repo, _, r, _ := recoveryFixture(t, "7")
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	binary := testfixture.FixtureBinary(t, "git", sourceRecoveryGitFixture)
	bin := filepath.Dir(binary)
	marker := filepath.Join(bin, "fetch-process")

	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestSourceRevisionRecoveryCancelsTheFetchProcess$")
	cmd.Env = append(gittest.IsolatedEnv(),
		"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"GIT_TERMINAL_PROMPT=1",
		"ROUNDFIX_TEST_RECOVERY_REPO="+repo,
		"ROUNDFIX_TEST_RECOVERY_REVISION="+r.SourceRevision,
		"ROUNDFIX_TEST_RECOVERY_MARKER="+marker,
		"ROUNDFIX_TEST_RECOVERY_GIT="+git,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("fetch process probe: %v\n%s", err, out)
	}
}

// The compiled fixture delegates local reads to real Git and blocks the fetch
// in its own process, so cancellation can prove that process was killed.
const sourceRecoveryGitFixture = `package main

import (
 "fmt"
 "io"
 "os"
 "os/exec"
 "time"
)

func main() {
 if len(os.Args) > 3 && os.Args[3] == "fetch" {
  input, err := io.ReadAll(os.Stdin)
  if err != nil { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
  state := "eof"
  if len(input) != 0 { state = "input" }
  marker := fmt.Sprintf("%d %s %s", os.Getpid(), os.Getenv("GIT_TERMINAL_PROMPT"), state)
  if err := os.WriteFile(os.Getenv("ROUNDFIX_TEST_RECOVERY_MARKER"), []byte(marker), 0600); err != nil {
   fmt.Fprintln(os.Stderr, err); os.Exit(1)
  }
  time.Sleep(30 * time.Second)
  return
 }
 cmd := exec.Command(os.Getenv("ROUNDFIX_TEST_RECOVERY_GIT"), os.Args[1:]...)
 cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
 if err := cmd.Run(); err != nil {
  if exit, ok := err.(*exec.ExitError); ok { os.Exit(exit.ExitCode()) }
  fmt.Fprintln(os.Stderr, err); os.Exit(1)
 }
}
`
