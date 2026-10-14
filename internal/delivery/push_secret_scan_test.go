// Suite: push Secret Scan at the publication boundary.
// Boundary IN: real Git history and the Delivery Queue engine.
// Boundary OUT: a temporary bare remote and the existing GitHub boundary fake.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/secretscan"
	"roundfix/internal/store"
)

// Sequential: isolates process-wide Git configuration used by NewGitHubCLI.
func TestPushRefusesAHistoryThatAddsASecret(t *testing.T) {
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_SYSTEM", "/dev/null")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	t.Setenv("GIT_TEST_MAINT_AUTO_DETACH", "0")
	t.Setenv("GIT_OPTIONAL_LOCKS", "0")
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
	tip := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	client := NewGitHubCLI(repo)
	_, err := client.PushBranch(t.Context(), "origin", "feat/secret-history", tip)
	assertPushDiagnosticHasNoSample(t, err, sample)
	var refusal *secretscan.MatchesError
	if !errors.As(err, &refusal) {
		t.Fatal("push did not return a MatchesError")
	}
	want := secretscan.Match{Path: path, Line: 3, Kind: "github-token", Commit: added}
	if len(refusal.Matches) != 1 || refusal.Matches[0] != want {
		t.Fatal("refusal did not identify the adding commit, file, line 3 and github-token kind")
	}
	if err.Error() != "push branch: refuse push: "+want.String() {
		t.Fatal("refusal has an unexpected diagnostic")
	}
	if refs := gittest.Run(t, remote, "for-each-ref", "--format=%(refname)", "refs/heads/"); strings.TrimSpace(refs) != "" {
		t.Fatal("refused history published a remote branch")
	}
	// Scan the recorded clean head while the checkout still has tainted history.
	clean, err := client.PushBranch(t.Context(), "origin", "feat/clean-history", base)
	assertPushDiagnosticHasNoSample(t, err, sample)
	if err != nil || clean.SHA != base {
		t.Fatal("clean history did not push its recorded head")
	}
	if got := strings.TrimSpace(gittest.Run(t, remote, "rev-parse", "refs/heads/feat/clean-history")); got != base {
		t.Fatal("clean remote branch has the wrong head")
	}
}

func TestEnginePushParksSecretFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot, slug = "/repo", "spec-one"
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug}); err != nil {
		t.Fatal(err)
	}
	sample := "ghp_" + strings.Repeat("b2", 18)
	matches := secretscan.ScanContent("evidence.txt", []byte("first\nsecond\n"+sample))
	if len(matches) != 1 {
		t.Fatal("fixture did not produce one Secret Match")
	}
	matches[0].Commit = strings.Repeat("c", 40)
	second := matches[0]
	second.Path, second.Line = "other.txt", 7
	matches = append(matches, second)
	boundary := &secretRefusingBoundary{
		fakeDeliveryBoundary: newFakeDeliveryBoundary(),
		err:                  fmt.Errorf("push branch: refuse push: %w", &secretscan.MatchesError{Matches: matches}),
	}
	workflow := newFakeDeliveryWorkflow()
	engine := newTestDeliveryEngine(runStore, workflow, boundary.fakeDeliveryBoundary)
	engine.pullRequests = boundary
	_, err := engine.Run(ctx, gitRoot)
	assertPushDiagnosticHasNoSample(t, err, sample)
	if err != nil {
		t.Fatal("engine could not persist the secret refusal")
	}
	queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
	item := queue.Items[0]
	assertPushDiagnosticHasNoSample(t, errors.New(item.Blocker), sample)
	wantBlocker := "secret-found: " + matches[0].String() + ", " + matches[1].String()
	if item.Stage != store.DeliveryStageParked || item.Blocker != wantBlocker {
		t.Fatal("item did not park with the safe secret-found match list")
	}
	park := ClassifyPark(queue, item)
	wantNext := "remove the secret from every commit the blocker names on the item branch in " + item.Worktree + ", then run roundfix deliver retry " + slug
	if park.Class != ParkClassFinding || park.Next != wantNext {
		t.Fatal("secret-found has the wrong class or next action")
	}
	if boundary.pullRequestEffects != 0 || boundary.mergeEffects != 0 || len(boundary.pullRequests) != 0 {
		t.Fatal("secret refusal created a Pull Request or merged")
	}
	for _, event := range workflow.events[slug] {
		if event == "pull-request" {
			t.Fatal("secret refusal reached Pull Request creation")
		}
	}
}

type secretRefusingBoundary struct {
	*fakeDeliveryBoundary
	err error
}

func (boundary *secretRefusingBoundary) WithWorkDir(workDir string) PullRequestBoundary {
	boundary.fakeDeliveryBoundary.WithWorkDir(workDir)
	return boundary
}

func (boundary *secretRefusingBoundary) PushBranch(context.Context, string, string, string) (RemoteHead, error) {
	return RemoteHead{}, boundary.err
}

func assertPushDiagnosticHasNoSample(t *testing.T, err error, sample string) {
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

func TestPushRefusesAFailedHistoryRead(t *testing.T) {
	t.Parallel()
	readErr := errors.New("history runner unavailable")
	sample := "ghp_" + strings.Repeat("e5", 18)
	for _, tt := range []struct {
		name   string
		result CommandResult
		err    error
	}{
		{name: "runner error", err: readErr},
		{name: "nonzero exit", result: CommandResult{ExitCode: 128, Stderr: "history unavailable"}},
		{name: "partial patch", result: CommandResult{ExitCode: 128, Stdout: sample}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			runner := newScriptedCommandRunner(t, commandStep{
				name: "git", args: secretscan.PushedHistoryArgs("candidate", "origin"),
				result: tt.result, err: tt.err,
			})
			client := GitHubCLI{WorkDir: "/repo", Runner: runner}
			_, err := client.PushBranch(t.Context(), "origin", "feat/history", "candidate")
			assertPushDiagnosticHasNoSample(t, err, sample)
			if err == nil || !strings.HasPrefix(err.Error(), "push branch: read pushed history: ") {
				t.Fatal("failed history read did not refuse the push")
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Fatal("history runner error lost its cause")
			}
		})
	}
}
