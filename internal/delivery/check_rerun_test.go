// Suite: bounded recovery of failed GitHub Actions checks.
// Boundary IN: real queue persistence, engine polling, and scripted CLI responses.
// Boundary OUT: GitHub and gh execution; every command uses an injected runner.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

const recoveryLink = "https://github.test/acme/repo/actions/runs/42/job/7"

func recoveryReport(bucket string) CheckReport {
	return CheckReport{HeadSHA: "candidate", Checks: []PullRequestCheck{{Name: "verify", Workflow: "CI", Link: recoveryLink, Bucket: bucket}}}
}

type fakeCheckRecovery struct {
	failures    []CheckFailure
	inspectErr  error
	rerunErr    error
	reruns      int
	inspections int
}

func (fake *fakeCheckRecovery) InspectFailedCheck(_ context.Context, workDir, _ string, head string, _ PullRequestCheck) (CheckFailure, error) {
	if workDir != "/worktrees/example" || head != "candidate" {
		return CheckFailure{}, errors.New("incorrect recovery worktree or head")
	}
	fake.inspections++
	if fake.inspectErr != nil {
		return CheckFailure{}, fake.inspectErr
	}
	if len(fake.failures) == 0 {
		return CheckFailure{}, errors.New("unexpected failure inspection")
	}
	failure := fake.failures[0]
	fake.failures = fake.failures[1:]
	return failure, nil
}
func (fake *fakeCheckRecovery) RerunFailedCheck(_ context.Context, workDir string, failure CheckFailure) error {
	if workDir != "/worktrees/example" || failure.RunID != "42" {
		return errors.New("incorrect re-run target")
	}
	fake.reruns++
	return fake.rerunErr
}

func outsideFailure(attempt int) CheckFailure {
	return CheckFailure{RunID: "42", Attempt: attempt, Packages: []string{"roundfix/internal/other"}, OutsideChange: true}
}

func exerciseCheckRecovery(t *testing.T, recovery CheckRecovery, reports []CheckReport) (store.DeliveryQueueItem, string, int) {
	t.Helper()
	ctx := t.Context()
	runStore := openDeliveryEngineStore(t, ctx)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/repo", []string{"example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Branch = "roundfix/deliver-example"
	item.Worktree = "/worktrees/example"
	item.Stage = store.DeliveryStageChecking
	item.PullRequestNumber = "1"
	item.CandidateCommits = []string{"candidate"}
	item.Warning = "premise-changed: internal/existing"
	recordDeliveryItemWorkspace(t, ctx, runStore, "/repo", &item)
	if err := runStore.UpdateDeliveryQueueItem(ctx, "/repo", item); err != nil {
		t.Fatal(err)
	}
	polls := 0
	boundary := &fakePullRequestBoundary{currentHeadChecks: func(_ context.Context, number string) (CheckReport, error) {
		if number != "1" {
			t.Fatalf("PR number=%q", number)
		}
		if polls >= len(reports) {
			t.Fatal("unexpected extra poll")
		}
		report := reports[polls]
		polls++
		return report, nil
	}}
	clock := &fakeDeliveryClock{now: time.Now()}
	var log bytes.Buffer
	engine := NewEngine(runStore, EngineDependencies{PullRequests: boundary, Checks: recovery, Clock: clock, Sleeper: clock, Log: &log, CheckInterval: time.Second, CheckTimeout: 5 * time.Second})
	if err := engine.checkCandidate(ctx, "/repo", &item); err != nil {
		t.Fatal(err)
	}
	return readDeliveryQueue(t, ctx, runStore, "/repo").Items[0], log.String(), polls
}

func TestAFailedCheckOutsideTheChangeIsRerunOnce(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1), outsideFailure(1)}}
	item, log, polls := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("fail"), recoveryReport("pending"), recoveryReport("pass")})
	if recovery.reruns != 1 || polls != 4 || item.Stage != store.DeliveryStageMerging || !strings.Contains(log, "roundfix: check re-run: Delivery Queue item example: verify (run 42)") {
		t.Fatalf("item=%+v reruns=%d polls=%d log=%q", item, recovery.reruns, polls, log)
	}
}
func TestARerunCheckThatPassesRecordsAWarning(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1)}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("pass")})
	if item.Stage != store.DeliveryStageMerging || item.Warning != "premise-changed: internal/existing; flaky-check: verify passed on re-run" || recovery.reruns != 1 {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}
func TestARerunCheckThatFailsAgainParksAsFlakyCheck(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1), outsideFailure(2)}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("pending"), recoveryReport("fail")})
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerFlakyCheck+": roundfix/internal/other" || recovery.reruns != 1 || ClassifyPark(store.DeliveryQueue{}, item).Class != ParkClassFlakyCheck {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}
func TestAFailedCheckInAChangedPackageParksWithoutARerun(t *testing.T) {
	failure := outsideFailure(1)
	failure.OutsideChange = false
	assertNoCheckRerun(t, &fakeCheckRecovery{failures: []CheckFailure{failure}}, "fail")
}
func TestAnUnattributableCheckFailureParksWithoutARerun(t *testing.T) {
	assertNoCheckRerun(t, &fakeCheckRecovery{failures: []CheckFailure{{RunID: "42", Attempt: 1}}}, "fail")
}
func TestACheckPastItsFirstAttemptIsNotRerun(t *testing.T) {
	assertNoCheckRerun(t, &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(2)}}, "fail")
}
func TestNilCheckRecoveryKeepsFailedCheckBehavior(t *testing.T) {
	item, _, _ := exerciseCheckRecovery(t, nil, []CheckReport{recoveryReport("fail")})
	if item.Blocker != BlockerChecksFailed {
		t.Fatalf("item=%+v", item)
	}
}
func TestCancelledChecksAreNotRerun(t *testing.T) {
	assertNoCheckRerun(t, &fakeCheckRecovery{}, "cancel")
}
func TestCheckInspectionErrorsDoNotRerun(t *testing.T) {
	assertNoCheckRerun(t, &fakeCheckRecovery{inspectErr: errors.New("log unavailable")}, "fail")
}
func TestCheckRerunErrorsParkTheItem(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1)}, rerunErr: errors.New("denied")}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail")})
	if item.Blocker != BlockerChecksFailed || recovery.reruns != 1 {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}
func TestSecondCheckFailureInChangedPackageIsNotCalledFlaky(t *testing.T) {
	second := outsideFailure(2)
	second.OutsideChange = false
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1), second}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("fail")})
	if item.Blocker != BlockerChecksFailed || recovery.reruns != 1 {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}
func assertNoCheckRerun(t *testing.T, recovery *fakeCheckRecovery, bucket string) {
	t.Helper()
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport(bucket)})
	if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerChecksFailed || recovery.reruns != 0 {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}

type recoveryCommandRunner struct {
	workDir string
	script  *scriptedCommandRunner
}

func (runner recoveryCommandRunner) Run(ctx context.Context, workDir, name string, args ...string) (CommandResult, error) {
	if workDir != runner.workDir {
		runner.script.t.Fatalf("workdir=%q want %q", workDir, runner.workDir)
	}
	return runner.script.Run(ctx, "/repo", name, args...)
}

func inspectRecoveryFixture(t *testing.T, log, changes string, attempt int, gitExpected bool) CheckFailure {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("// fixture\nmodule example.test/repo\n\ngo 1.26\n"), 0600); err != nil {
		t.Fatal(err)
	}
	steps := []commandStep{
		{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":` + strconv.Itoa(attempt) + `}`}},
		{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: log}},
	}
	if gitExpected {
		steps = append(steps,
			commandStep{name: "git", args: []string{"symbolic-ref", "refs/remotes/origin/HEAD"}, result: CommandResult{Stdout: "refs/remotes/origin/trunk"}},
			commandStep{name: "git", args: []string{"fetch", "origin", "+refs/heads/trunk:refs/remotes/origin/trunk"}},
			commandStep{name: "git", args: []string{"merge-base", "refs/remotes/origin/trunk", "candidate"}, result: CommandResult{Stdout: "common"}},
			commandStep{name: "git", args: []string{"diff", "--name-only", "--no-renames", "-z", "common", "candidate"}, result: CommandResult{Stdout: strings.ReplaceAll(changes, "\n", "\x00")}},
		)
	}
	script := newScriptedCommandRunner(t, steps...)
	client := GitHubCLI{WorkDir: "/wrong", Runner: recoveryCommandRunner{dir, script}}
	failure, err := client.InspectFailedCheck(t.Context(), dir, "origin", "candidate", PullRequestCheck{Link: recoveryLink})
	if err != nil {
		t.Fatal(err)
	}
	return failure
}

func TestGitHubCLIAttributesAFailedCheckToGoPackages(t *testing.T) {
	failure := inspectRecoveryFixture(t, "test\tstep\t2026-10-01T07:00:00Z FAIL\texample.test/repo/internal/other\t1.234s\nFAIL\texample.test/repo/internal/other\t1s\nFAIL\texample.test/repo/pkg/z\t2s", "internal/change/file.go\ninternal/otherish/file.go", 1, true)
	if !failure.OutsideChange || failure.Attempt != 1 || failure.RunID != "42" || !reflect.DeepEqual(failure.Packages, []string{"example.test/repo/internal/other", "example.test/repo/pkg/z"}) {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLIRejectsChangedPackageAttribution(t *testing.T) {
	failure := inspectRecoveryFixture(t, "FAIL\texample.test/repo/internal/other\t1s", "internal/other/file_test.go", 1, true)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLIBuildFailuresAreUnattributable(t *testing.T) {
	failure := inspectRecoveryFixture(t, "FAIL\texample.test/repo/internal/other\t1s\nFAIL\texample.test/repo/broken [build failed]", "", 1, false)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLISetupFailuresAreUnattributable(t *testing.T) {
	failure := inspectRecoveryFixture(t, "FAIL\texample.test/repo/internal/other\t1s\nFAIL\texample.test/repo/broken [setup failed]", "", 1, false)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLILogsWithoutGoPackagesAreUnattributable(t *testing.T) {
	failure := inspectRecoveryFixture(t, "a JavaScript test failed\nFAIL", "", 1, false)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLIExternalModuleFailuresAreUnattributable(t *testing.T) {
	failure := inspectRecoveryFixture(t, "FAIL\tanother.test/repo/pkg\t1s", "", 1, false)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLIRootPackageOverlapsEveryChangedPath(t *testing.T) {
	failure := inspectRecoveryFixture(t, "FAIL\texample.test/repo\t1s", "internal/change/file.go", 1, true)
	if failure.OutsideChange {
		t.Fatalf("failure=%+v", failure)
	}
}
func TestGitHubCLINonActionsCheckIsUnattributable(t *testing.T) {
	script := newScriptedCommandRunner(t)
	failure, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", "origin", "candidate", PullRequestCheck{Link: "https://ci.test/build/42"})
	if err != nil || failure.RunID != "" || failure.OutsideChange {
		t.Fatalf("failure=%+v err=%v", failure, err)
	}
}
func TestGitHubCLIRerunsOnlyFailedJobs(t *testing.T) {
	script := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"run", "rerun", "42", "--failed"}})
	if err := (GitHubCLI{Runner: script}).RerunFailedCheck(t.Context(), "/repo", outsideFailure(1)); err != nil {
		t.Fatal(err)
	}
}

func TestCheckRerunRestartsTheTimeout(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1)}}
	reports := []CheckReport{
		recoveryReport("pending"), recoveryReport("pending"), recoveryReport("pending"), recoveryReport("pending"),
		recoveryReport("fail"), recoveryReport("pending"), recoveryReport("pending"), recoveryReport("pending"), recoveryReport("pass"),
	}
	item, _, polls := exerciseCheckRecovery(t, recovery, reports)
	if item.Stage != store.DeliveryStageMerging || polls != len(reports) || recovery.reruns != 1 {
		t.Fatalf("item=%+v polls=%d reruns=%d", item, polls, recovery.reruns)
	}
}

func TestAnOldCheckAttemptCannotKeepResettingTheTimeout(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1), outsideFailure(1), outsideFailure(1), outsideFailure(1), outsideFailure(1), outsideFailure(1)}}
	reports := []CheckReport{recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail"), recoveryReport("fail")}
	item, _, polls := exerciseCheckRecovery(t, recovery, reports)
	if item.Blocker != BlockerChecksTimeout || polls != 6 || recovery.reruns != 1 {
		t.Fatalf("item=%+v polls=%d reruns=%d", item, polls, recovery.reruns)
	}
}

func TestMultipleFailedJobsInOneRunAreRerunTogetherOnce(t *testing.T) {
	failed := recoveryReport("fail")
	failed.Checks = append(failed.Checks, PullRequestCheck{Name: "lint", Workflow: "CI", Link: "https://github.test/acme/repo/actions/runs/42/job/8", Bucket: "fail"})
	passed := recoveryReport("pass")
	passed.Checks = append(passed.Checks, PullRequestCheck{Name: "lint", Workflow: "CI", Link: "https://github.test/acme/repo/actions/runs/42/job/80", Bucket: "pass"})
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1), outsideFailure(1)}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{failed, passed})
	if item.Stage != store.DeliveryStageMerging || recovery.reruns != 1 || item.Warning != "premise-changed: internal/existing; flaky-check: verify passed on re-run; flaky-check: lint passed on re-run" {
		t.Fatalf("item=%+v reruns=%d", item, recovery.reruns)
	}
}

func TestSkippedRerunDoesNotRecordAPassWarning(t *testing.T) {
	recovery := &fakeCheckRecovery{failures: []CheckFailure{outsideFailure(1)}}
	item, _, _ := exerciseCheckRecovery(t, recovery, []CheckReport{recoveryReport("fail"), recoveryReport("skipping")})
	if item.Stage != store.DeliveryStageMerging || item.Warning != "premise-changed: internal/existing" {
		t.Fatalf("item=%+v", item)
	}
}

func TestGitHubCLIRejectsMalformedAttemptResponses(t *testing.T) {
	script := newScriptedCommandRunner(t, commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: "invalid"}})
	_, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", "origin", "candidate", PullRequestCheck{Link: recoveryLink})
	if err == nil || !strings.Contains(err.Error(), "parse failed check attempt") {
		t.Fatalf("err=%v", err)
	}
}

func TestGitHubCLIReportsFailedLogReadErrors(t *testing.T) {
	script := newScriptedCommandRunner(t,
		commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":1}`}},
		commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{ExitCode: 1, Stderr: "log not available"}},
	)
	failure, err := (GitHubCLI{Runner: script}).InspectFailedCheck(t.Context(), "/repo", "origin", "candidate", PullRequestCheck{Link: recoveryLink})
	if err == nil || failure.OutsideChange || !strings.Contains(err.Error(), "log not available") {
		t.Fatalf("failure=%+v err=%v", failure, err)
	}
}

// This runner exercises actual local Git; gh remains a scripted boundary.
type localGitRecoveryRunner struct {
	t      *testing.T
	script *scriptedCommandRunner
}

func (runner localGitRecoveryRunner) Run(ctx context.Context, workDir, name string, args ...string) (CommandResult, error) {
	if name == "git" {
		return CommandResult{Stdout: gittest.Run(runner.t, workDir, args...)}, nil
	}
	return runner.script.Run(ctx, "/repo", name, args...)
}

func TestGitHubCLIAttributionUsesTheRefreshedMergeBaseWithRealGit(t *testing.T) {
	for _, renamed := range []bool{false, true} {
		name := "remote changes excluded"
		if renamed {
			name = "renamed package file included"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			remote := filepath.Join(root, "remote.git")
			repo := filepath.Join(root, "item")
			gittest.InitRepo(t, remote, "--bare", "-b", "trunk")
			gittest.InitRepo(t, repo, "-b", "trunk")
			write := func(relative, content string) {
				t.Helper()
				target := filepath.Join(repo, relative)
				if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(target, []byte(content), 0600); err != nil {
					t.Fatal(err)
				}
			}
			write("go.mod", "module example.test/repo\n\ngo 1.26\n")
			write("internal/other/file.go", "package other\n")
			write("internal/change/file.go", "package change\n")
			gittest.Run(t, repo, "add", ".")
			gittest.Run(t, repo, "commit", "-m", "fixture base")
			gittest.Run(t, repo, "remote", "add", "origin", remote)
			gittest.Run(t, repo, "push", "origin", "trunk")
			gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
			gittest.Run(t, repo, "checkout", "-b", "feat/item")
			if renamed {
				gittest.Run(t, repo, "mv", "internal/other/file.go", "internal/change/moved.go")
			} else {
				write("internal/change/file.go", "package change\n// item change\n")
			}
			gittest.Run(t, repo, "add", ".")
			gittest.Run(t, repo, "commit", "-m", "fixture item")
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			gittest.Run(t, repo, "checkout", "trunk")
			write("internal/other/remote.go", "package other\n")
			gittest.Run(t, repo, "add", ".")
			gittest.Run(t, repo, "commit", "-m", "fixture remote change")
			remoteHead := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			gittest.Run(t, repo, "push", "origin", "trunk")
			gittest.Run(t, repo, "update-ref", "refs/remotes/origin/trunk", remoteHead+"^")
			gittest.Run(t, repo, "checkout", "feat/item")
			script := newScriptedCommandRunner(t,
				commandStep{name: "gh", args: []string{"run", "view", "42", "--json", "attempt"}, result: CommandResult{Stdout: `{"attempt":1}`}},
				commandStep{name: "gh", args: []string{"run", "view", "42", "--log-failed"}, result: CommandResult{Stdout: "FAIL\texample.test/repo/internal/other\t1s"}},
			)
			client := GitHubCLI{Runner: localGitRecoveryRunner{t, script}}
			failure, err := client.InspectFailedCheck(t.Context(), repo, "origin", head, PullRequestCheck{Link: recoveryLink})
			if err != nil || failure.OutsideChange == renamed {
				t.Fatalf("failure=%+v err=%v renamed=%t", failure, err, renamed)
			}
			if refreshed := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "refs/remotes/origin/trunk")); refreshed != remoteHead {
				t.Fatalf("remote default head=%s want=%s", refreshed, remoteHead)
			}
		})
	}
}
