package delivery

// Suite: pull request boundary
// Invariant: retries observe GitHub state before creating or merging again.
// Boundary IN: GitHubCLI command construction, parsing, and idempotency.
// Boundary OUT: git, the GitHub CLI, and the network, replaced by a scripted runner.

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPushPublishesTheRecordedHeadNotTheCheckout(t *testing.T) {
	const headSHA = "0123456789abcdef"
	runner := newScriptedCommandRunner(t,
		commandStep{
			name: "git",
			args: []string{"log", "-p", "-U0", "--reverse", "--no-color", "--no-ext-diff", "--no-renames", "--format=commit %H", headSHA, "--not", "--remotes=origin"},
		},
		commandStep{
			name: "git",
			args: []string{"push", "origin", headSHA + ":refs/heads/feat/delivery"},
		},
		commandStep{
			name:   "git",
			args:   []string{"ls-remote", "--heads", "origin", "refs/heads/feat/delivery"},
			result: CommandResult{Stdout: headSHA + "\trefs/heads/feat/delivery\n"},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	remoteHead, err := boundary.PushBranch(t.Context(), "origin", "feat/delivery", headSHA)

	if err != nil {
		t.Fatalf("PushBranch returned error: %v", err)
	}
	want := RemoteHead{Remote: "origin", Branch: "feat/delivery", SHA: headSHA}
	if remoteHead != want {
		t.Fatalf("PushBranch() = %#v, want %#v", remoteHead, want)
	}
}

func TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead(t *testing.T) {
	runner := newScriptedCommandRunner(t,
		commandStep{
			name: "git",
			args: []string{"log", "-p", "-U0", "--reverse", "--no-color", "--no-ext-diff", "--no-renames", "--format=commit %H", "head-one", "--not", "--remotes=origin"},
		},
		commandStep{
			name: "git",
			args: []string{"push", "origin", "head-one:refs/heads/feat/delivery"},
		},
		commandStep{
			name: "git",
			args: []string{"ls-remote", "--heads", "origin", "refs/heads/feat/delivery"},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	_, err := boundary.PushBranch(t.Context(), "origin", "feat/delivery", "head-one")

	if err == nil || err.Error() != `read remote head after push: git did not report "refs/heads/feat/delivery"` {
		t.Fatalf("PushBranch error = %v, want missing remote-head error", err)
	}
}

func TestPullRequestBoundaryObservesRemoteHeadWithoutPushing(t *testing.T) {
	const headSHA = "0123456789abcdef"
	runner := newScriptedCommandRunner(t, commandStep{
		name:   "git",
		args:   []string{"ls-remote", "--heads", "origin", "refs/heads/feat/delivery"},
		result: CommandResult{Stdout: headSHA + "\trefs/heads/feat/delivery\n"},
	})
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	remoteHead, found, err := boundary.RemoteBranchHead(t.Context(), "origin", "feat/delivery")

	if err != nil {
		t.Fatalf("RemoteBranchHead returned error: %v", err)
	}
	if !found {
		t.Fatal("RemoteBranchHead did not find the scripted remote head")
	}
	want := RemoteHead{Remote: "origin", Branch: "feat/delivery", SHA: headSHA}
	if remoteHead != want {
		t.Fatalf("RemoteBranchHead() = %#v, want %#v", remoteHead, want)
	}
}

func TestPullRequestBoundaryReusesAnOpenPullRequest(t *testing.T) {
	runner := newScriptedCommandRunner(t, commandStep{
		name: "gh",
		args: []string{
			"pr", "list",
			"--head", "feat/delivery",
			"--base", "main",
			"--state", "open",
			"--limit", "2",
			"--json", pullRequestJSONFields,
		},
		result: CommandResult{Stdout: `[{"number":17,"url":"https://github.test/acme/repo/pull/17","state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-one","mergedAt":"","mergeCommit":null}]`},
	})
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	result, err := boundary.FindOrCreatePullRequest(t.Context(), PullRequestRequest{
		HeadBranch: "feat/delivery",
		BaseBranch: "main",
		Title:      "feat: durable delivery",
		Body:       "Delivery details.",
	})

	if err != nil {
		t.Fatalf("FindOrCreatePullRequest returned error: %v", err)
	}
	if result.Created {
		t.Fatal("FindOrCreatePullRequest reported creating an existing pull request")
	}
	if result.PullRequest.Number != "17" || result.PullRequest.HeadSHA != "head-one" {
		t.Fatalf("FindOrCreatePullRequest() = %#v, want existing pull request 17 at head-one", result)
	}
}

func TestPullRequestBoundaryCreatesOnlyWhenNoOpenPullRequestExists(t *testing.T) {
	runner := newScriptedCommandRunner(t,
		commandStep{
			name: "gh",
			args: []string{
				"pr", "list",
				"--head", "feat/delivery",
				"--base", "main",
				"--state", "open",
				"--limit", "2",
				"--json", pullRequestJSONFields,
			},
			result: CommandResult{Stdout: `[]`},
		},
		commandStep{
			name: "gh",
			args: []string{
				"pr", "create",
				"--head", "feat/delivery",
				"--base", "main",
				"--title", "feat: durable delivery",
				"--body", "Delivery details.",
			},
			result: CommandResult{Stdout: "https://github.test/acme/repo/pull/18\n"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "https://github.test/acme/repo/pull/18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: `{"number":18,"url":"https://github.test/acme/repo/pull/18","state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	result, err := boundary.FindOrCreatePullRequest(t.Context(), PullRequestRequest{
		HeadBranch: "feat/delivery",
		BaseBranch: "main",
		Title:      "feat: durable delivery",
		Body:       "Delivery details.",
	})

	if err != nil {
		t.Fatalf("FindOrCreatePullRequest returned error: %v", err)
	}
	if !result.Created {
		t.Fatal("FindOrCreatePullRequest did not report creating the missing pull request")
	}
	if result.PullRequest.Number != "18" || result.PullRequest.HeadSHA != "head-two" {
		t.Fatalf("FindOrCreatePullRequest() = %#v, want created pull request 18 at head-two", result)
	}
}

func TestPullRequestBoundaryReportsCurrentHeadChecks(t *testing.T) {
	const pullRequest = `{"number":18,"url":"https://github.test/acme/repo/pull/18","state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
		commandStep{
			name: "gh",
			args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"},
			result: CommandResult{
				Stdout:   `[{"bucket":"pending","link":"https://github.test/check/1","name":"verify","state":"IN_PROGRESS","workflow":"CI"}]`,
				ExitCode: 8,
			},
		},
		commandStep{
			name:   "gh",
			args:   []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"},
			result: CommandResult{Stdout: `{"total_count":0,"workflow_runs":[]}`},
		},
		commandStep{
			name: "git",
			args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	report, err := boundary.CurrentHeadChecks(t.Context(), "18")

	if err != nil {
		t.Fatalf("CurrentHeadChecks returned error: %v", err)
	}
	if report.HeadSHA != "head-two" {
		t.Fatalf("CurrentHeadChecks HeadSHA = %q, want head-two", report.HeadSHA)
	}
	wantChecks := []PullRequestCheck{{
		Name:     "verify",
		State:    "IN_PROGRESS",
		Bucket:   "pending",
		Link:     "https://github.test/check/1",
		Workflow: "CI",
	}}
	if !reflect.DeepEqual(report.Checks, wantChecks) {
		t.Fatalf("CurrentHeadChecks Checks = %#v, want %#v", report.Checks, wantChecks)
	}
}

func TestNoChecksReportedIsAnEmptyReport(t *testing.T) {
	const pullRequest = `{"number":18,"url":"https://github.test/acme/repo/pull/18","state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
		commandStep{
			name: "gh",
			args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"},
			result: CommandResult{
				Stderr:   "no checks reported on the 'feat/delivery' branch\n",
				ExitCode: 1,
			},
		},
		commandStep{
			name:   "gh",
			args:   []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"},
			result: CommandResult{Stdout: `{"total_count":0,"workflow_runs":[]}`},
		},
		commandStep{
			name: "git",
			args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	report, err := boundary.CurrentHeadChecks(t.Context(), "18")

	if err != nil {
		t.Fatalf("CurrentHeadChecks returned error: %v", err)
	}
	if report.HeadSHA != "head-two" {
		t.Fatalf("CurrentHeadChecks HeadSHA = %q, want head-two", report.HeadSHA)
	}
	if len(report.Checks) != 0 {
		t.Fatalf("CurrentHeadChecks Checks = %#v, want an empty report", report.Checks)
	}
}

func TestCheckCommandFailuresRemainErrors(t *testing.T) {
	const pullRequest = `{"number":18,"url":"https://github.test/acme/repo/pull/18","state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
		commandStep{
			name: "gh",
			args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"},
			result: CommandResult{
				Stderr:   "failed: no checks reported because GitHub was unavailable\n",
				ExitCode: 1,
			},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	_, err := boundary.CurrentHeadChecks(t.Context(), "18")

	if err == nil || err.Error() != "read pull request checks: failed: no checks reported because GitHub was unavailable" {
		t.Fatalf("CurrentHeadChecks error = %v, want the gh command failure", err)
	}
}

func TestPullRequestBoundaryReportsFailedCurrentHeadChecks(t *testing.T) {
	const pullRequest = `{"number":18,"state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
		commandStep{
			name: "gh",
			args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"},
			result: CommandResult{
				Stdout:   `[{"bucket":"fail","link":"https://github.test/check/1","name":"verify","state":"FAILURE","workflow":"CI"}]`,
				ExitCode: 1,
			},
		},
		commandStep{
			name:   "gh",
			args:   []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"},
			result: CommandResult{Stdout: `{"total_count":0,"workflow_runs":[]}`},
		},
		commandStep{
			name: "git",
			args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: pullRequest},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	report, err := boundary.CurrentHeadChecks(t.Context(), "18")

	if err != nil {
		t.Fatalf("CurrentHeadChecks returned error: %v", err)
	}
	if len(report.Checks) != 1 || report.Checks[0].Bucket != "fail" || report.Checks[0].State != "FAILURE" {
		t.Fatalf("CurrentHeadChecks Checks = %#v, want failed verify check", report.Checks)
	}
}

func TestPullRequestBoundaryRejectsChecksFromAMovedHead(t *testing.T) {
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: `{"number":18,"state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"},
			result: CommandResult{Stdout: `[]`},
		},
		commandStep{
			name:   "gh",
			args:   []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"},
			result: CommandResult{Stdout: `{"total_count":0,"workflow_runs":[]}`},
		},
		commandStep{
			name: "git",
			args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: `{"number":18,"state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-three","mergedAt":"","mergeCommit":null}`},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	_, err := boundary.CurrentHeadChecks(t.Context(), "18")

	if err == nil || err.Error() != `read pull request checks: PR Head Branch moved from "head-two" to "head-three" while checks were read` {
		t.Fatalf("CurrentHeadChecks error = %v, want moved-head error", err)
	}
}

func TestPullRequestBoundaryDoesNotMergeTwice(t *testing.T) {
	runner := newScriptedCommandRunner(t, commandStep{
		name:   "gh",
		args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
		result: CommandResult{Stdout: `{"number":18,"url":"https://github.test/acme/repo/pull/18","state":"MERGED","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"2026-09-24T13:00:00Z","mergeCommit":{"oid":"merge-one"}}`},
	})
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	result, err := boundary.MergePullRequest(t.Context(), "18", "head-two")

	if err != nil {
		t.Fatalf("MergePullRequest returned error: %v", err)
	}
	if !result.AlreadyMerged {
		t.Fatal("MergePullRequest did not report the existing merge")
	}
	if result.PullRequest.MergeCommit != "merge-one" {
		t.Fatalf("MergePullRequest merge commit = %q, want merge-one", result.PullRequest.MergeCommit)
	}
}

func TestPullRequestBoundaryMergesTheExpectedHead(t *testing.T) {
	runner := newScriptedCommandRunner(t,
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: `{"number":18,"state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"","mergeCommit":null}`},
		},
		commandStep{
			name: "gh",
			args: []string{"pr", "merge", "18", "--squash", "--match-head-commit", "head-two"},
		},
		commandStep{
			name:   "gh",
			args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
			result: CommandResult{Stdout: `{"number":18,"state":"MERGED","headRefName":"feat/delivery","headRefOid":"head-two","mergedAt":"2026-09-24T13:00:00Z","mergeCommit":{"oid":"merge-one"}}`},
		},
	)
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	result, err := boundary.MergePullRequest(t.Context(), "18", "head-two")

	if err != nil {
		t.Fatalf("MergePullRequest returned error: %v", err)
	}
	if result.AlreadyMerged {
		t.Fatal("MergePullRequest reported a new merge as pre-existing")
	}
	if result.PullRequest.MergeCommit != "merge-one" {
		t.Fatalf("MergePullRequest merge commit = %q, want merge-one", result.PullRequest.MergeCommit)
	}
}

func TestPullRequestBoundaryRefusesToMergeAnUnexpectedHead(t *testing.T) {
	runner := newScriptedCommandRunner(t, commandStep{
		name:   "gh",
		args:   []string{"pr", "view", "18", "--json", pullRequestJSONFields},
		result: CommandResult{Stdout: `{"number":18,"state":"OPEN","headRefName":"feat/delivery","headRefOid":"head-three","mergedAt":"","mergeCommit":null}`},
	})
	boundary := GitHubCLI{WorkDir: "/repo", Runner: runner}

	_, err := boundary.MergePullRequest(t.Context(), "18", "head-two")

	if err == nil || err.Error() != `merge pull request: PR Head Branch is at "head-three", expected "head-two"` {
		t.Fatalf("MergePullRequest error = %v, want unexpected-head error", err)
	}
}

type commandStep struct {
	name   string
	args   []string
	result CommandResult
	err    error
}

type scriptedCommandRunner struct {
	t     *testing.T
	steps []commandStep
	next  int
}

func newScriptedCommandRunner(t *testing.T, steps ...commandStep) *scriptedCommandRunner {
	t.Helper()
	runner := &scriptedCommandRunner{t: t, steps: steps}
	t.Cleanup(func() {
		if runner.next != len(runner.steps) {
			t.Errorf("command runner consumed %d of %d steps", runner.next, len(runner.steps))
		}
	})
	return runner
}

func (runner *scriptedCommandRunner) Run(_ context.Context, workDir, name string, args ...string) (CommandResult, error) {
	runner.t.Helper()
	if runner.next >= len(runner.steps) {
		runner.t.Fatalf("unexpected command: %s %#v", name, args)
	}
	step := runner.steps[runner.next]
	runner.next++
	if workDir != "/repo" {
		runner.t.Fatalf("command workdir = %q, want /repo", workDir)
	}
	if name != step.name || !reflect.DeepEqual(args, step.args) {
		runner.t.Fatalf("command = %s %#v, want %s %#v", name, args, step.name, step.args)
	}
	return step.result, step.err
}

type fakePullRequestBoundary struct {
	remoteBranchHead        func(context.Context, string, string) (RemoteHead, bool, error)
	pushBranch              func(context.Context, string, string, string) (RemoteHead, error)
	findOrCreatePullRequest func(context.Context, PullRequestRequest) (PullRequestResult, error)
	currentHeadChecks       func(context.Context, string) (CheckReport, error)
	mergePullRequest        func(context.Context, string, string) (MergeResult, error)
}

func (fake *fakePullRequestBoundary) WithWorkDir(string) PullRequestBoundary {
	return fake
}

var _ PullRequestBoundary = (*fakePullRequestBoundary)(nil)

func (fake *fakePullRequestBoundary) RemoteBranchHead(ctx context.Context, remote, branch string) (RemoteHead, bool, error) {
	if fake.remoteBranchHead == nil {
		return RemoteHead{}, false, errors.New("unexpected RemoteBranchHead call")
	}
	return fake.remoteBranchHead(ctx, remote, branch)
}

func (fake *fakePullRequestBoundary) PushBranch(ctx context.Context, remote, branch, head string) (RemoteHead, error) {
	if fake.pushBranch == nil {
		return RemoteHead{}, errors.New("unexpected PushBranch call")
	}
	return fake.pushBranch(ctx, remote, branch, head)
}

func (fake *fakePullRequestBoundary) FindOrCreatePullRequest(ctx context.Context, req PullRequestRequest) (PullRequestResult, error) {
	if fake.findOrCreatePullRequest == nil {
		return PullRequestResult{}, errors.New("unexpected FindOrCreatePullRequest call")
	}
	return fake.findOrCreatePullRequest(ctx, req)
}

func (fake *fakePullRequestBoundary) CurrentHeadChecks(ctx context.Context, number string) (CheckReport, error) {
	if fake.currentHeadChecks == nil {
		return CheckReport{}, errors.New("unexpected CurrentHeadChecks call")
	}
	return fake.currentHeadChecks(ctx, number)
}

func (fake *fakePullRequestBoundary) MergePullRequest(ctx context.Context, number, expectedHead string) (MergeResult, error) {
	if fake.mergePullRequest == nil {
		return MergeResult{}, errors.New("unexpected MergePullRequest call")
	}
	return fake.mergePullRequest(ctx, number, expectedHead)
}

func TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead(t *testing.T) {
	t.Parallel()
	const pr = `{"number":18,"baseRefName":"main","headRefOid":"head-two","mergeable":"MERGEABLE","mergeStateStatus":"CLEAN"}`
	const runs = `{"total_count":4,"workflow_runs":[
        {"id":9,"name":"CI | Backend","path":".github/workflows/ci-backend.yml","head_sha":"head-two","event":"pull_request","status":"completed","conclusion":"cancelled","run_number":9},
        {"id":10,"name":"CI | Backend","path":".github/workflows/ci-backend.yml","head_sha":"head-two","event":"pull_request","status":"in_progress","run_number":10,"run_attempt":1,"html_url":"https://github.test/acme/repo/actions/runs/10"},
        {"id":11,"name":"CI | Backend","path":".github/workflows/ci-backend.yml","head_sha":"head-one","event":"pull_request","status":"completed","conclusion":"success","run_number":11},
        {"id":1,"name":"CI | Conventions","path":".github/workflows/ci-conventions.yml","head_sha":"head-two","event":"pull_request","status":"completed","conclusion":"success","run_number":1,"html_url":"https://github.test/acme/repo/actions/runs/1"}
    ]}`
	tests := []struct {
		name      string
		step      int
		result    CommandResult
		runnerErr error
		wantError string
	}{
		{name: "reads head runs"},
		{name: "incomplete page", step: 2, result: CommandResult{Stdout: `{"total_count":101,"workflow_runs":[{"id":1}]}`}, wantError: "read workflow runs: 101 runs reported for head-two, more than one page of 1"},
		{name: "runs read fails", step: 2, result: CommandResult{ExitCode: 1, Stderr: "HTTP 403"}, wantError: "read workflow runs: HTTP 403"},
		{name: "ls-tree fails", step: 3, result: CommandResult{ExitCode: 1, Stderr: "missing head"}, wantError: "list workflow files: missing head"},
		{name: "git show fails", step: 4, result: CommandResult{ExitCode: 1, Stderr: "missing blob"}, wantError: "read workflow file: missing blob"},
		{name: "runs malformed", step: 2, result: CommandResult{Stdout: "{"}, wantError: "read workflow runs:"},
		{name: "runs runner fails", step: 2, runnerErr: context.Canceled, wantError: "read workflow runs: context canceled"},
		{name: "ls-tree runner fails", step: 3, runnerErr: context.Canceled, wantError: "list workflow files: context canceled"},
		{name: "git show runner fails", step: 4, runnerErr: context.Canceled, wantError: "read workflow file: context canceled"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			steps := []commandStep{
				{name: "gh", args: []string{"pr", "view", "18", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: pr}},
				{name: "gh", args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"}, result: CommandResult{Stdout: `[{"bucket":"pass","name":"Validate PR title","workflow":"CI | Conventions","state":"SUCCESS","link":"https://github.test/acme/repo/actions/runs/1/job/2"}]`}},
				{name: "gh", args: []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"}, result: CommandResult{Stdout: runs}},
				{name: "git", args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"}, result: CommandResult{Stdout: ".github/workflows/ci-backend.yml\x00.github/workflows/ci-conventions.yml\x00.github/workflows/README.md\x00"}},
				{name: "git", args: []string{"show", "head-two:.github/workflows/ci-backend.yml"}, result: CommandResult{Stdout: "on: {pull_request: {paths: ['internal/**']}}"}},
				{name: "git", args: []string{"show", "head-two:.github/workflows/ci-conventions.yml"}, result: CommandResult{Stdout: "on: {pull_request: {branches: [main], types: [opened, reopened, edited, synchronize]}}"}},
				{name: "gh", args: []string{"pr", "view", "18", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: pr}},
			}
			if tt.wantError != "" {
				steps[tt.step].result = tt.result
				steps[tt.step].err = tt.runnerErr
				steps = steps[:tt.step+1]
			}
			client := GitHubCLI{WorkDir: "/repo", Runner: newScriptedCommandRunner(t, steps...)}
			report, err := client.CurrentHeadChecks(t.Context(), "18")
			if tt.wantError != "" {
				if err == nil || !strings.HasPrefix(err.Error(), tt.wantError) {
					t.Fatalf("CurrentHeadChecks error = %v, want prefix %q", err, tt.wantError)
				}
				if tt.runnerErr != nil && !errors.Is(err, tt.runnerErr) {
					t.Fatalf("CurrentHeadChecks lost runner error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			wantRuns := []WorkflowRun{
				{ID: "10", Name: "CI | Backend", Path: ".github/workflows/ci-backend.yml", Event: "pull_request", Status: "in_progress", Attempt: 1, URL: "https://github.test/acme/repo/actions/runs/10"},
				{ID: "1", Name: "CI | Conventions", Path: ".github/workflows/ci-conventions.yml", Event: "pull_request", Status: "completed", Conclusion: "success", URL: "https://github.test/acme/repo/actions/runs/1"},
			}
			if !reflect.DeepEqual(report.WorkflowRuns, wantRuns) {
				t.Fatalf("WorkflowRuns = %+v, want %+v (run 9 is superseded by run 10)", report.WorkflowRuns, wantRuns)
			}
			wantExpected := []string{".github/workflows/ci-conventions.yml"}
			if !reflect.DeepEqual(report.ExpectedWorkflows, wantExpected) {
				t.Fatalf("ExpectedWorkflows = %v, want %v", report.ExpectedWorkflows, wantExpected)
			}
			if report.HeadSHA != "head-two" || report.Mergeable != "MERGEABLE" || report.MergeState != "CLEAN" || len(report.Checks) != 1 || report.Checks[0].Name != "Validate PR title" {
				t.Fatalf("CheckReport = %+v", report)
			}
		})
	}
}

func TestPullRequestReadTrimsTheBaseBranch(t *testing.T) {
	t.Parallel()
	pr, err := parsePullRequest([]byte(`{"number":18,"baseRefName":" main "}`))
	if err != nil || pr.BaseBranch != "main" {
		t.Fatalf("PullRequest = %+v, error = %v", pr, err)
	}
}

func TestCurrentHeadChecksKeepsWorkflowEventsSeparate(t *testing.T) {
	t.Parallel()
	const pr = `{"number":18,"baseRefName":" main ","headRefOid":"head-two"}`
	runner := newScriptedCommandRunner(t,
		commandStep{name: "gh", args: []string{"pr", "view", "18", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: pr}},
		commandStep{name: "gh", args: []string{"pr", "checks", "18", "--json", "bucket,link,name,state,workflow"}, result: CommandResult{Stdout: `[]`}},
		commandStep{name: "gh", args: []string{"api", "repos/{owner}/{repo}/actions/runs?head_sha=head-two&per_page=100"}, result: CommandResult{Stdout: `{"total_count":5,"workflow_runs":[
            {"id":30,"name":"CI","path":"ci.yml","event":"push","head_sha":"head-two","run_number":30,"status":"COMPLETED","conclusion":"SUCCESS"},
            {"id":20,"name":"CI","path":"ci.yml","event":"pull_request","head_sha":"head-two","run_number":20,"status":"IN_PROGRESS"},
            {"id":10,"name":"CI","path":"ci.yml","event":"pull_request","head_sha":"head-two","run_number":10,"status":"completed","conclusion":"failure"},
            {"id":40,"name":"CI","path":"ci.yml","event":"push","head_sha":"head-one","run_number":40},
            {"id":5,"name":"CI","path":"other.yml","event":"push","head_sha":"head-two","run_number":5}
        ]}`}},
		commandStep{name: "git", args: []string{"ls-tree", "-z", "--name-only", "head-two", "--", ".github/workflows/"}, result: CommandResult{Stdout: ".github/workflows/ci.yaml\x00"}},
		commandStep{name: "git", args: []string{"show", "head-two:.github/workflows/ci.yaml"}, result: CommandResult{Stdout: "on: {pull_request: {branches: main}}"}},
		commandStep{name: "gh", args: []string{"pr", "view", "18", "--json", pullRequestJSONFields}, result: CommandResult{Stdout: pr}},
	)
	report, err := (GitHubCLI{WorkDir: "/repo", Runner: runner}).CurrentHeadChecks(t.Context(), "18")
	if err != nil {
		t.Fatal(err)
	}
	want := []WorkflowRun{
		{ID: "30", Name: "CI", Path: "ci.yml", Event: "push", Status: "completed", Conclusion: "success"},
		{ID: "20", Name: "CI", Path: "ci.yml", Event: "pull_request", Status: "in_progress"},
		{ID: "5", Name: "CI", Path: "other.yml", Event: "push"},
	}
	if !reflect.DeepEqual(report.WorkflowRuns, want) {
		t.Fatalf("WorkflowRuns = %+v, want %+v", report.WorkflowRuns, want)
	}
	if !reflect.DeepEqual(report.ExpectedWorkflows, []string{".github/workflows/ci.yaml"}) {
		t.Fatalf("ExpectedWorkflows = %v", report.ExpectedWorkflows)
	}
}
