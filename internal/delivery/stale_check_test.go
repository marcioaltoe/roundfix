// Suite: stale check recovery; real local Git and scripted GitHub boundaries.
package delivery

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func staleFailure(attempt int, tip string) CheckFailure {
	return CheckFailure{RunID: "42", Attempt: attempt, TestedBase: strings.Repeat("a", 40), DefaultTip: tip, Stale: true}
}

func TestARetryAfterTheDefaultBranchMovedRerunsTheStaleCheckAndMerges(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{staleFailure(2, strings.Repeat("b", 40))}}
	item, log, polls := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("pending"), recoveryReport("pass")})
	want := "roundfix: check stale: Delivery Queue item example: verify tested " + strings.Repeat("a", 40) + ", default branch is at " + strings.Repeat("b", 40) + "; re-run (run 42)"
	if item.Stage != store.DeliveryStageMerging || item.Warning != "premise-changed: internal/existing" || recovery.reruns != 1 || polls != 3 || !strings.Contains(log, want) || strings.Contains(log, "check re-run:") {
		t.Fatalf("item=%+v reruns=%d polls=%d log=%q", item, recovery.reruns, polls, log)
	}
}

func TestAStaleCheckKeepsPollingWhileGitHubReportsTheAttemptItReran(t *testing.T) {
	current := outsideFailure(2)
	recovery := &fakeCheckRecovery{failures: []CheckFailure{staleFailure(2, "tip"), staleFailure(2, "tip"), current}}
	item, _, polls := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail"), recoveryReport("pass")})
	if item.Stage != store.DeliveryStageMerging || recovery.reruns != 1 || polls != 4 || item.Warning != "premise-changed: internal/existing" {
		t.Fatalf("item=%+v reruns=%d polls=%d", item, recovery.reruns, polls)
	}
}

func TestAStaleRerunThatFailsOnTheCurrentTipGetsTheOutsideChangeRerun(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{staleFailure(2, "tip"), outsideFailure(3), outsideFailure(3), outsideFailure(4)}}
	item, log, polls := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail")})
	if item.Blocker != BlockerFlakyCheck+": roundfix/internal/other" || ClassifyPark(store.DeliveryQueue{}, item).Class != ParkClassFlakyCheck || recovery.reruns != 2 || polls != 4 || strings.Count(log, "check re-run:") != 1 {
		t.Fatalf("item=%+v reruns=%d polls=%d log=%q", item, recovery.reruns, polls, log)
	}
}

func TestAStaleCheckWhoseRerunIsRefusedParksChecksFailed(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{staleFailure(2, "tip")}, rerunErr: errors.New("denied")}
	item, log, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail")})
	if item.Blocker != BlockerChecksFailed || recovery.reruns != 1 || !strings.Contains(log, "roundfix: check recovery: Delivery Queue item example: denied") {
		t.Fatalf("item=%+v reruns=%d log=%q", item, recovery.reruns, log)
	}
}

// Git failures must reach the adapter as exit codes, rather than failing the test helper.
type staleLocalGitRunner struct {
	script *scriptedCommandRunner
}

func (runner staleLocalGitRunner) Run(ctx context.Context, workDir, name string, args ...string) (CommandResult, error) {
	if name != "git" {
		return runner.script.Run(ctx, "/repo", name, args...)
	}
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = workDir
	output, err := command.CombinedOutput()
	result := CommandResult{Stdout: string(output)}
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return result, err
		}
		result.ExitCode = exit.ExitCode()
	}
	return result, nil
}

func TestGitHubCLIJudgesAFailedCheckByTheTipItTested(t *testing.T) {
	root := t.TempDir()
	remote, repo := filepath.Join(root, "remote.git"), filepath.Join(root, "item")
	gittest.InitRepo(t, remote, "--bare", "-b", "trunk")
	gittest.InitRepo(t, repo, "-b", "trunk")
	write := func(contents string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
		gittest.Run(t, repo, "add", ".")
		gittest.Run(t, repo, "commit", "-m", "fixture")
	}
	write("module example.test/repo\n")
	old := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "remote", "add", "upstream", remote)
	gittest.Run(t, repo, "push", "upstream", "trunk")
	gittest.Run(t, repo, "symbolic-ref", "refs/remotes/upstream/HEAD", "refs/remotes/upstream/trunk")
	write("module example.test/repo\n// current tip\n")
	tip := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.Run(t, repo, "push", "upstream", "trunk")
	gittest.Run(t, repo, "commit", "--allow-empty", "-m", "fixture descendant")
	descendant := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	for _, tc := range []struct {
		name, base string
		stale      bool
	}{
		{"older tip", old, true}, {"unknown commit", strings.Repeat("f", 40), true}, {"current tip", tip, false}, {"contains current tip", descendant, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Leave the local tracking ref behind; inspection must fetch before judging it.
			gittest.Run(t, repo, "update-ref", "refs/remotes/upstream/trunk", old)
			steps := []commandStep{
				{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":2}`}},
				{name: "gh", args: []string{"api", "repos/{owner}/{repo}/check-runs/7/annotations?per_page=100"}, result: CommandResult{Stdout: `[{"title":"tested-base","message":"` + tc.base + `"}]`}},
			}
			if !tc.stale {
				steps = append(steps, commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: "FAIL\texample.test/repo/internal/other\t1s"}})
			}
			client := GitHubCLI{Runner: staleLocalGitRunner{newScriptedCommandRunner(t, steps...)}}
			failure, err := client.InspectFailedCheck(t.Context(), repo, "upstream", tip, PullRequestCheck{Link: recoveryLink})
			if err != nil || failure.Stale != tc.stale || failure.TestedBase != tc.base || failure.DefaultTip != tip || (!tc.stale && !failure.OutsideChange) || (tc.stale && len(failure.Packages) != 0) {
				t.Fatalf("failure=%+v err=%v", failure, err)
			}
		})
	}
}

func TestGitHubCLIWithoutATestedBaseInspectsAsBefore(t *testing.T) {
	sha := strings.Repeat("a", 40)
	for _, tc := range []struct{ name, annotations string }{
		{"absent", `[]`},
		{"other annotation", `[{"title":"other","message":"` + sha + `"}]`},
		{"malformed", `[{"title":"tested-base","message":"bad"}]`},
		{"uppercase", `[{"title":"tested-base","message":"` + strings.Repeat("A", 40) + `"}]`},
		{"ambiguous", `[{"title":"tested-base","message":"` + sha + `"},{"title":"tested-base","message":"` + sha + `"}]`},
		{"valid and malformed", `[{"title":"tested-base","message":"` + sha + `"},{"title":"tested-base","message":"bad"}]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.test/repo\n"), 0600); err != nil {
				t.Fatal(err)
			}
			script := newScriptedCommandRunner(t,
				commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":1}`}},
				commandStep{name: "gh", args: []string{"api", "repos/{owner}/{repo}/check-runs/7/annotations?per_page=100"}, result: CommandResult{Stdout: tc.annotations}},
				commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: "FAIL\texample.test/repo/internal/other\t1s"}},
				commandStep{name: "git", args: []string{"symbolic-ref", "refs/remotes/origin/HEAD"}, result: CommandResult{Stdout: "refs/remotes/origin/trunk"}},
				commandStep{name: "git", args: []string{"fetch", "origin", "+refs/heads/trunk:refs/remotes/origin/trunk"}},
				commandStep{name: "git", args: []string{"merge-base", "refs/remotes/origin/trunk", "candidate"}, result: CommandResult{Stdout: "common"}},
				commandStep{name: "git", args: []string{"diff", "--name-only", "--no-renames", "-z", "common", "candidate"}, result: CommandResult{Stdout: "internal/change/file.go\x00"}},
			)
			failure, err := (GitHubCLI{Runner: recoveryCommandRunner{dir, script}}).InspectFailedCheck(t.Context(), dir, "origin", "candidate", PullRequestCheck{Link: recoveryLink})
			if err != nil || failure.Stale || failure.TestedBase != "" || failure.DefaultTip != "" || !failure.OutsideChange || failure.Attempt != 1 {
				t.Fatalf("failure=%+v err=%v", failure, err)
			}
		})
	}
}

func TestGitHubCLIRerunsAStaleFailurePastItsFirstAttempt(t *testing.T) {
	script := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"run", "rerun", "42", "--failed"}})
	if err := (GitHubCLI{Runner: script}).RerunFailedCheck(t.Context(), "/repo", staleFailure(2, "tip")); err != nil {
		t.Fatal(err)
	}
}

func TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{
		staleFailure(2, "tip1"), staleFailure(2, "tip1"), staleFailure(3, "tip2"),
		staleFailure(3, "tip2"), staleFailure(4, "tip3"), staleFailure(4, "tip3"),
	}}
	reports := make([]CheckReport, 6)
	for i := range reports {
		reports[i] = recoveryReport("fail")
	}
	item, _, polls := exerciseCheckRecovery(t, recovery, reports)
	if item.Blocker != BlockerChecksTimeout || recovery.reruns != 3 || polls != 6 {
		t.Fatalf("item=%+v reruns=%d polls=%d", item, recovery.reruns, polls)
	}
}

func TestStaleInspectionErrorsAreReportedBeforeReadingTheLog(t *testing.T) {
	sha := strings.Repeat("a", 40)
	steps := []commandStep{
		{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":2}`}},
		{name: "gh", args: []string{"api", "repos/{owner}/{repo}/check-runs/7/annotations?per_page=100"}, result: CommandResult{Stdout: `[{"title":"tested-base","message":"` + sha + `"}]`}},
		{name: "git", args: []string{"symbolic-ref", "refs/remotes/origin/HEAD"}, result: CommandResult{Stdout: "refs/remotes/origin/trunk"}},
		{name: "git", args: []string{"fetch", "origin", "+refs/heads/trunk:refs/remotes/origin/trunk"}},
		{name: "git", args: []string{"rev-parse", "refs/remotes/origin/trunk"}, result: CommandResult{Stdout: sha}},
		{name: "git", args: []string{"cat-file", "-e", sha + "^{commit}"}},
		{name: "git", args: []string{"merge-base", "--is-ancestor", "refs/remotes/origin/trunk", sha}},
	}
	for _, tc := range []struct {
		name  string
		index int
	}{
		{"annotation read", 1}, {"default resolution", 2}, {"fetch", 3}, {"tip resolution", 4}, {"ancestry", 6},
	} {
		t.Run(tc.name, func(t *testing.T) {
			commands := append([]commandStep(nil), steps[:tc.index+1]...)
			commands[tc.index].result = CommandResult{ExitCode: 2, Stderr: "inspection denied"}
			script := newScriptedCommandRunner(t, commands...)
			_, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", "origin", "candidate", PullRequestCheck{Link: recoveryLink})
			if err == nil || !strings.Contains(err.Error(), "inspection denied") {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestActionsLinkWithoutAJobKeepsThePreviousInspection(t *testing.T) {
	script := newScriptedCommandRunner(t,
		commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":2}`}},
		commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: "build [build failed]"}},
	)
	failure, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", "origin", "candidate", PullRequestCheck{Link: "https://github.test/acme/repo/actions/runs/42"})
	if err != nil || failure.RunID != "42" || failure.Attempt != 2 || failure.Stale || failure.OutsideChange || failure.TestedBase != "" {
		t.Fatalf("failure=%+v err=%v", failure, err)
	}
}

func TestOutsideChangePassAfterAStaleRerunKeepsTheExistingWarning(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{staleFailure(2, "tip"), outsideFailure(3)}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("fail"), recoveryReport("pass")})
	if item.Stage != store.DeliveryStageMerging || recovery.reruns != 2 || item.Warning != "premise-changed: internal/existing; flaky-check: verify passed on re-run" {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}
