package delivery

// Suite: pull request boundary
// Invariant: retries observe GitHub state before creating or merging again.
// Boundary IN: GitHubCLI command construction, parsing, and idempotency.
// Boundary OUT: git, the GitHub CLI, and the network, replaced by a scripted runner.

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestPushPublishesTheRecordedHeadNotTheCheckout(t *testing.T) {
	const headSHA = "0123456789abcdef"
	runner := newScriptedCommandRunner(t,
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
