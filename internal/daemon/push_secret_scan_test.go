// Suite: Final Push Secret Scan at the real Git boundary.
// Boundary IN: GitPusher over an added-then-removed runtime sample.
// Boundary OUT: a temporary local bare remote; no external publication.
package daemon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/secretscan"
	"roundfix/internal/testfixture"
)

func TestFinalPushRefusesAHistoryThatAddsASecret(t *testing.T) {
	t.Parallel()
	if runFinalPushInIsolatedProcess(t) {
		return
	}
	root := t.TempDir()
	repo, remote := filepath.Join(root, "item"), filepath.Join(root, "remote.git")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "test: clean base")
	base := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "checkout", "-b", "feat/secret-history")
	sample := "ghp_" + strings.Repeat("a1", 18)
	const path = "evidence.txt"
	if err := os.WriteFile(filepath.Join(repo, path), []byte("first\nsecond\n"+sample+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, repo, "add", path)
	gittest.Run(t, repo, "commit", "-m", "test: add sample")
	added := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	if err := os.WriteFile(filepath.Join(repo, path), []byte("first\nsecond\nremoved\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, repo, "add", path)
	gittest.Run(t, repo, "commit", "-m", "test: remove sample")
	pusher := GitPusher{}
	err := pusher.Push(t.Context(), PushRequest{WorkDir: repo, Remote: "origin", Branch: "feat/secret-history"})
	assertFinalPushDiagnosticHasNoSample(t, err, sample)
	var refusal *secretscan.MatchesError
	if !errors.As(err, &refusal) {
		t.Fatal("push did not return a MatchesError")
	}
	want := secretscan.Match{Path: path, Line: 3, Kind: "github-token", Commit: added}
	if len(refusal.Matches) != 1 || refusal.Matches[0] != want {
		t.Fatal("refusal did not identify the adding commit, file, line 3 and github-token kind")
	}
	if err.Error() != "run Final Push: refuse push: "+want.String() {
		t.Fatal("refusal has an unexpected diagnostic")
	}
	if refs := gittest.Run(t, remote, "for-each-ref", "--format=%(refname)", "refs/heads/"); strings.TrimSpace(refs) != "" {
		t.Fatal("refused history published a remote branch")
	}
	gittest.Run(t, repo, "checkout", "-b", "feat/clean-history", base)
	err = pusher.Push(t.Context(), PushRequest{WorkDir: repo, Remote: "origin", Branch: "feat/clean-history"})
	assertFinalPushDiagnosticHasNoSample(t, err, sample)
	if err != nil {
		t.Fatal("clean Final Push was refused")
	}
	if got := strings.TrimSpace(gittest.Run(t, remote, "rev-parse", "refs/heads/feat/clean-history")); got != base {
		t.Fatal("clean remote branch has the wrong head")
	}
}

func assertFinalPushDiagnosticHasNoSample(t *testing.T, err error, sample string) {
	t.Helper()
	if err == nil {
		return
	}
	for index := 0; index+8 <= len(sample); index++ {
		if strings.Contains(err.Error(), sample[index:index+8]) {
			t.Fatal("push diagnostic contains sample bytes")
		}
	}
}

func TestFinalPushRefusesAFailedHistoryRead(t *testing.T) {
	t.Parallel()
	if runFinalPushInIsolatedProcess(t) {
		return
	}
	root := t.TempDir()
	repo, remote := filepath.Join(root, "item"), filepath.Join(root, "remote.git")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	tracePath := os.Getenv("GIT_TRACE")
	err := (GitPusher{}).Push(t.Context(), PushRequest{WorkDir: repo, Remote: "origin", Branch: "feat/history"})
	if err == nil || !strings.HasPrefix(err.Error(), "run Final Push: read pushed history: ") {
		t.Fatal("unborn history did not refuse Final Push at the history read")
	}
	trace, readErr := os.ReadFile(tracePath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(trace), "git log ") || strings.Contains(string(trace), "git push ") {
		t.Fatal("failed history read must run git log and never git push")
	}
	if refs := gittest.Run(t, remote, "for-each-ref", "--format=%(refname)", "refs/heads/"); strings.TrimSpace(refs) != "" {
		t.Fatal("failed history read published a remote branch")
	}
}

// Re-execute only this test with isolated Git configuration. The parent remains
// parallel and never changes the process-wide environment used by other tests.
func runFinalPushInIsolatedProcess(t *testing.T) bool {
	t.Helper()
	const fixtureVariable = "ROUNDFIX_FINAL_PUSH_TEST_FIXTURE"
	if os.Getenv(fixtureVariable) == t.Name() {
		return false
	}
	command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+t.Name()+"$")
	command.Env = append(gittest.IsolatedEnv(),
		fixtureVariable+"="+t.Name(),
		"GIT_TRACE="+filepath.Join(t.TempDir(), "git.trace"),
	)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("Final Push isolated test: %v\n%s", err, output)
	}
	return true
}

func TestFinalPushHistoryReadFailureOmitsItsOutput(t *testing.T) {
	t.Parallel()
	sample := "ghp_" + strings.Repeat("c3", 18)
	const fixtureVariable = "ROUNDFIX_FINAL_PUSH_TEST_FIXTURE"
	if os.Getenv(fixtureVariable) != t.Name() {
		// The child's git prints a partial patch with the sample on stdout,
		// writes nothing to stderr and fails, only for git log.
		realGit, err := exec.LookPath("git")
		if err != nil {
			t.Fatal(err)
		}
		fakeGit := testfixture.FixtureBinary(t, "git", fmt.Sprintf(`package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "log" {
			fmt.Printf("commit 0\n+%%s\n", %q)
			os.Exit(1)
		}
	}
	command := exec.Command(%q, os.Args[1:]...)
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := command.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		os.Exit(1)
	}
}
`, sample, realGit))
		bin := filepath.Dir(fakeGit)
		command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+t.Name()+"$")
		command.Env = append(gittest.IsolatedEnv(),
			fixtureVariable+"="+t.Name(),
			"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("Final Push isolated test: %v\n%s", err, output)
		}
		return
	}
	root := t.TempDir()
	repo, remote := filepath.Join(root, "item"), filepath.Join(root, "remote.git")
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.InitRepo(t, remote, "--bare", "-b", "main")
	gittest.Run(t, repo, "remote", "add", "origin", remote)
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "test: base")
	err := (GitPusher{}).Push(t.Context(), PushRequest{WorkDir: repo, Remote: "origin", Branch: "feat/history"})
	if err == nil || !strings.HasPrefix(err.Error(), "run Final Push: read pushed history: ") {
		t.Fatal("failed history read did not refuse Final Push at the history read")
	}
	assertFinalPushDiagnosticHasNoSample(t, err, sample)
	if refs := gittest.Run(t, remote, "for-each-ref", "--format=%(refname)", "refs/heads/"); strings.TrimSpace(refs) != "" {
		t.Fatal("failed history read published a remote branch")
	}
}
