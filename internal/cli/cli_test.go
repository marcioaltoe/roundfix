package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	"roundfix/internal/watch"
)

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix watch") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.HasPrefix(stdout.String(), "roundfix ") {
		t.Fatalf("expected version output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunCommandHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix fetch --source coderabbit --pr <number>") {
		t.Fatalf("expected fetch help output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsCheck(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "check"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Roundfix skills check passed") {
		t.Fatalf("expected skills check output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix-watch") || !strings.Contains(stdout.String(), "roundfix-resolve-round") {
		t.Fatalf("expected both skill names, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCopiesArtifacts(t *testing.T) {
	targetDir := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "codex", "--dir", targetDir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Installed Roundfix skills for codex") {
		t.Fatalf("expected install output, got %q", stdout.String())
	}
	for _, path := range []string{
		"roundfix-watch/SKILL.md",
		"roundfix-watch/agents/openai.yaml",
		"roundfix-resolve-round/SKILL.md",
		"roundfix-resolve-round/agents/openai.yaml",
	} {
		if _, err := os.Stat(filepath.Join(targetDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallRejectsUnsupportedTarget(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "other"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported skill install target") {
		t.Fatalf("expected unsupported target error, got %q", stderr.String())
	}
}

func TestRunOperationalCommandAcceptsMVPFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedCode   int
		expectedOutput string
	}{
		{
			name:           "fetch",
			args:           []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"},
			expectedCode:   0,
			expectedOutput: "reached Fetched",
		},
		{
			name:           "resolve",
			args:           []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			expectedCode:   0,
			expectedOutput: "resolve selected 1 downloaded Unresolved Review Issue",
		},
		{
			name:           "watch",
			args:           []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			expectedCode:   0,
			expectedOutput: "Watch Run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != tt.expectedCode {
				t.Fatalf("expected exit code %d, got %d", tt.expectedCode, code)
			}
			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Fatalf("expected output to contain %q, got stdout=%q stderr=%q", tt.expectedOutput, stdout.String(), stderr.String())
			}
			if !strings.Contains(output, filepath.Join(repoDir, ".roundfix")) {
				t.Fatalf("expected artifact dir in output, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if strings.Contains(output, "not implemented yet") {
				t.Fatalf("did not expect scaffold message, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if tt.name == "fetch" {
				if !strings.Contains(stdout.String(), "did not start an Agent, commit, or push") {
					t.Fatalf("expected fetch safety confirmation, got %q", stdout.String())
				}
				if !strings.Contains(stdout.String(), "Review Issues: 1") {
					t.Fatalf("expected fetched issue count, got %q", stdout.String())
				}
				issuePath := filepath.Join(repoDir, ".roundfix", "reviews", "pr-123", "round-001", "issue_001.md")
				issueContent, err := os.ReadFile(issuePath)
				if err != nil {
					t.Fatalf("expected Review Issue artifact %s: %v", issuePath, err)
				}
				if !strings.Contains(string(issueContent), "source_ref: thread:PRRT_test,comment:PRRC_test") {
					t.Fatalf("expected source ref in issue artifact, got %s", string(issueContent))
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
			} else if tt.name == "resolve" {
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake verification output") {
					t.Fatalf("expected fake verification output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertAgentLogContains(t, repoDir, "fake agent output")
			} else if tt.name == "watch" {
				if !strings.Contains(stderr.String(), "Review Source status: settled") {
					t.Fatalf("expected fake Review Source status output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Fetched Round 001 with 1 Review Issue") {
					t.Fatalf("expected watch fetch output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake verification output") {
					t.Fatalf("expected fake verification output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "reached Clean") {
					t.Fatalf("expected watch Clean terminal outcome, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertAgentLogContains(t, repoDir, "fake agent output")
			}
		})
	}
}

func TestRunWatchTimeoutOffersManualReviewWithoutFetching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("watch timeout must not fetch Review Source issues")
		return nil, nil
	})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusPending},
			{State: watch.StatusPending},
			{State: watch.StatusPending},
		},
	}).Status)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1s
  review_timeout: 2s
  quiet_period: 1s
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected watch timeout exit 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Review Source status: pending") {
		t.Fatalf("expected pending Review Source status output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached TimedOut") {
		t.Fatalf("expected TimedOut terminal outcome, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "@coderabbitai review") {
		t.Fatalf("expected manual review trigger guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("timeout must not fetch or run Agent, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchStopRequestBeforeAgentMarksStopped(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	withSuccessfulPreflight(t, repoDir)
	withWatchStatus(t, func(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
		cancel()
		return reviewsource.WatchStatus{State: watch.StatusPending}, nil
	})
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("Stop Request before Agent must not fetch Review Source issues")
		return nil, nil
	})
	withChangedPaths(t, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Watch Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Watch Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Changed paths after Stop Request: none") {
		t.Fatalf("expected changed-path report, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("Stop Request before Agent must not fetch or run Agent, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveStopRequestDuringAgentPreservesWorkAndSkipsDaemonMutations(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withChangedPaths(t, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Resolve Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Resolve Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "partial agent output") {
		t.Fatalf("expected persisted Agent output in stream, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "M src/app.go") {
		t.Fatalf("expected changed path report, got %q", stderr.String())
	}
	if verifier.calls != 0 {
		t.Fatalf("Stop Request must not run verification, got %d calls", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("Stop Request must not create commits, got %d calls", committer.calls)
	}
	if sourceResolver.calls != 0 {
		t.Fatalf("Stop Request must not resolve Review Source threads, got %d calls", sourceResolver.calls)
	}
	if pusher.calls != 0 {
		t.Fatalf("Stop Request must not push, got %d calls", pusher.calls)
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusPending {
		t.Fatalf("Stop Request must preserve assigned issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")

	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, preflight.DirtyWorktreeError{
			ArtifactDir: filepath.Join(repoDir, ".roundfix"),
			Changes:     []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}},
		}
	})
	var retryStdout bytes.Buffer
	var retryStderr bytes.Buffer
	retryCode := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &retryStdout, &retryStderr)
	if retryCode != 2 {
		t.Fatalf("expected dirty preserved work to block next Run with exit 2, got %d", retryCode)
	}
	if !strings.Contains(retryStderr.String(), "M src/app.go") {
		t.Fatalf("expected dirty path on retry, got %q", retryStderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveSIGINTStopReturns130(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withChangedPaths(t, nil)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stoppedCode := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)
	code := exitForInterrupt(stoppedCode, true)

	if code != 130 {
		t.Fatalf("expected SIGINT exit 130, got %d", code)
	}
	if !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected Stopped output before SIGINT exit, got %q", stderr.String())
	}
}

func TestExitForWatchOutcome(t *testing.T) {
	tests := []struct {
		outcome string
		code    int
	}{
		{outcome: store.StateClean, code: 0},
		{outcome: store.StateMaxRoundsReached, code: 0},
		{outcome: store.StateStopped, code: 0},
		{outcome: store.StateBudgetExceeded, code: 1},
		{outcome: store.StateTimedOut, code: 1},
		{outcome: store.StateFailed, code: 1},
	}

	for _, tt := range tests {
		t.Run(tt.outcome, func(t *testing.T) {
			if code := exitForWatchOutcome(tt.outcome); code != tt.code {
				t.Fatalf("expected exit code %d for %s, got %d", tt.code, tt.outcome, code)
			}
		})
	}
}

func TestRunResolveRejectsMissingCompatibleArtifactsBeforeRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no downloaded Unresolved Review Issues in Compatible Artifacts") {
		t.Fatalf("expected missing Compatible Artifacts error, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix fetch --source coderabbit --pr 123") {
		t.Fatalf("expected fetch guidance, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveHonorsRoundSelector(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--round", "2", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit code 0, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 1 downloaded Unresolved Review Issue") {
		t.Fatalf("expected one selected issue, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Round scope: 002") {
		t.Fatalf("expected round scope in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push blocked: 1 Unresolved Review Issue") {
		t.Fatalf("expected Final Push to be blocked by unselected unresolved issue, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveDeduplicatesBeforeBatching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit code 0, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 2 downloaded Unresolved Review Issue(s)") {
		t.Fatalf("expected selected issue count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "assigned 1 newest occurrence(s) into 1 Batch(es)") {
		t.Fatalf("expected deduped assignment count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "associated 1 older duplicate occurrence(s)") {
		t.Fatalf("expected duplicate association count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Marked 1 older duplicate Review Issue occurrence") {
		t.Fatalf("expected older duplicate marker, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
		t.Fatalf("expected only newest source thread resolution, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push completed") {
		t.Fatalf("expected Final Push after all duplicates terminal, got %q", stderr.String())
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 1 {
		t.Fatalf("expected only newest issue to resolve source thread, got %d", got)
	}
	if sourceResolver.requests[0].Issues[0].FilePath != filepath.Join(repoDir, ".roundfix", "reviews", "pr-123", "round-002", "issue_001.md") {
		t.Fatalf("expected newest duplicate source resolution, got %#v", sourceResolver.requests[0].Issues[0])
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveVerificationFailureDoesNotCommit(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{err: errors.New("tests failed")}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call, got %d", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("expected no Batch commit after verification failure, got %d", committer.calls)
	}
	if sourceResolver.calls != 0 {
		t.Fatalf("expected no Review Source resolution after verification failure, got %d", sourceResolver.calls)
	}
	if pusher.calls != 0 {
		t.Fatalf("expected no Final Push after verification failure, got %d", pusher.calls)
	}
	if !strings.Contains(stderr.String(), "tests failed") {
		t.Fatalf("expected verification failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not commit, push, or resolve Review Source threads") {
		t.Fatalf("expected daemon-owned mutation boundary, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveCommitsSuccessfulBatchWithRemainingUnresolved(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	result := persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call, got %d", verifier.calls)
	}
	if committer.calls != 1 {
		t.Fatalf("expected one Batch commit, got %d", committer.calls)
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 1 {
		t.Fatalf("expected only assigned Batch issue to be source-resolved, got %d", got)
	}
	if pusher.calls != 0 {
		t.Fatalf("expected no Final Push while Unresolved Review Issues remain, got %d", pusher.calls)
	}
	if committer.messages[0] != daemon.BatchCommitMessage(1) {
		t.Fatalf("expected Batch commit message %q, got %q", daemon.BatchCommitMessage(1), committer.messages[0])
	}
	if !strings.Contains(stderr.String(), "assigned 2 newest occurrence(s) into 2 Batch(es)") {
		t.Fatalf("expected two planned Batches, got %q", stderr.String())
	}
	first, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(result.IssuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusResolved {
		t.Fatalf("expected first Batch issue resolved, got %q", first.Status)
	}
	if second.Status != rounds.StatusPending {
		t.Fatalf("expected second Batch issue to remain unresolved, got %q", second.Status)
	}
	if !strings.Contains(stderr.String(), "Final Push blocked: 1 Unresolved Review Issue") {
		t.Fatalf("expected Final Push blocked output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveFinalPushRunsOnceAfterAllUnresolvedTerminal(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	persistCLIReviewItems(t, repoDir, 2, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call for one Batch, got %d", verifier.calls)
	}
	if committer.calls != 1 {
		t.Fatalf("expected one Batch commit, got %d", committer.calls)
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 2 {
		t.Fatalf("expected both assigned terminal issues to be source-resolved together, got %d", got)
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one Final Push, got %d", pusher.calls)
	}
	if pusher.remotes[0] != "origin" || pusher.branches[0] != "feature/review" {
		t.Fatalf("expected push to origin feature/review, got %s %s", pusher.remotes[0], pusher.branches[0])
	}
	if strings.Contains(strings.Join(pusher.args, " "), "--force") {
		t.Fatalf("Final Push must not force-push, got args %#v", pusher.args)
	}
	if !strings.Contains(stderr.String(), "Final Push completed: git push origin HEAD:feature/review") {
		t.Fatalf("expected Final Push output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveResolvesInvalidAssignedIssueSourceThread(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withAgentRunner(t, &fakeAgentRunner{status: rounds.StatusInvalid})
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := sourceResolver.requests[0].Issues[0].Status; got != rounds.StatusInvalid {
		t.Fatalf("expected invalid issue to resolve source thread, got %q", got)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveProbeFailureDoesNotCreateRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("adapter missing")})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected probe preflight exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "adapter missing") {
		t.Fatalf("expected probe diagnostic, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunResolveAgentFailureMarksBatchFailed(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{runErr: errors.New("agent crashed")})
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "agent crashed") {
		t.Fatalf("expected Agent failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not commit, push, or resolve Review Source threads") {
		t.Fatalf("expected daemon-owned mutation boundary, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveRejectsIncompatibleArtifacts(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "other-branch")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Head Repository \"owner/project\", PR Head Branch \"feature/review\"") {
		t.Fatalf("expected incompatible artifact context, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunOperationalCommandRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "missing pull request with no input",
			args:     []string{"fetch", "--source", "coderabbit", "--no-input"},
			contains: "missing required --pr because --no-input disables Interactive Input",
		},
		{
			name:     "invalid pull request number",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "abc", "--no-input"},
			contains: "--pr must be a positive integer",
		},
		{
			name:     "unsupported source",
			args:     []string{"fetch", "--source", "other", "--pr", "123", "--no-input"},
			contains: "unsupported Review Source",
		},
		{
			name:     "unsupported agent",
			args:     []string{"resolve", "--pr", "123", "--agent", "other", "--no-input"},
			contains: "unsupported Agent",
		},
		{
			name:     "invalid max rounds",
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--max-rounds", "0", "--no-input"},
			contains: "--max-rounds must be greater than 0",
		},
		{
			name:     "unknown flag",
			args:     []string{"fetch", "--unknown"},
			contains: "flag provided but not defined",
		},
		{
			name:     "unexpected argument",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "123", "extra"},
			contains: "unexpected argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.contains) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.contains, stderr.String())
			}
			if !strings.Contains(stderr.String(), "did not create a Run") {
				t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunFetchCollectsMissingPullRequestWithInteractiveInput(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withCurrentPullRequestSuggestion(t, "321")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch exit 0, got %d", code)
	}
	if inputReq.Command != "fetch" {
		t.Fatalf("expected fetch Interactive Input request, got %#v", inputReq)
	}
	if inputReq.PRSuggestion.Value != "321" || inputReq.PRSuggestion.Source != "current" {
		t.Fatalf("expected current PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if !strings.Contains(stderr.String(), "Interactive Input collected command parameters.") {
		t.Fatalf("expected Interactive Input confirmation, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Roundfix") || !strings.Contains(stderr.String(), "state: FetchingIssues") {
		t.Fatalf("expected Live Run View output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "321" {
		t.Fatalf("expected remembered PR 321, got %#v", defaults)
	}
}

func TestRunFetchReportsBlankInteractivePullRequestWithoutSuggestingInteractiveAgain(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withCurrentPullRequestSuggestion(t, "")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		return req.Values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Interactive Input did not collect required --pr") {
		t.Fatalf("expected blank Interactive Input guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "use --interactive") {
		t.Fatalf("did not expect guidance to reopen Interactive Input, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveInteractiveInputSuggestsConfiguredAgentAndRememberedPullRequest(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if err := runStore.RememberInteractiveDefaults(context.Background(), store.InteractiveDefaults{PRNumber: "123", Agent: "claude"}); err != nil {
		t.Fatalf("remember defaults: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withCurrentPullRequestSuggestion(t, "")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--interactive"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected resolve exit 0, got %d", code)
	}
	if inputReq.PRSuggestion.Value != "123" || inputReq.PRSuggestion.Source != "remembered" {
		t.Fatalf("expected remembered PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if inputReq.AgentSuggestion.Value != "codex" || inputReq.AgentSuggestion.Source != "config" {
		t.Fatalf("expected configured Agent suggestion before remembered value, got %#v", inputReq.AgentSuggestion)
	}
	if !strings.Contains(stderr.String(), "state: ResolvingWithAgent") {
		t.Fatalf("expected resolving Live Run View, got %q", stderr.String())
	}
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "123" || defaults.Agent != "codex" {
		t.Fatalf("expected remembered PR and configured Agent after run, got %#v", defaults)
	}
}

func TestRunInteractiveInputRunsPreflightBeforeSideEffects(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("fetch must not run after post-input Preflight Validation failure")
		return nil, nil
	})
	withCurrentPullRequestSuggestion(t, "123")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("post-input preflight failed")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected post-input preflight exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "post-input preflight failed") {
		t.Fatalf("expected preflight failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
	_ = repoDir
}

func TestRunOperationalCommandAppliesConfigAndCLIArtifactDirPrecedence(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  artifact_dir: user-artifacts
`)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
defaults:
  artifact_dir: project-artifacts
`)

	var projectStdout bytes.Buffer
	var projectStderr bytes.Buffer
	projectCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &projectStdout, &projectStderr)

	if projectCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", projectCode)
	}
	if !strings.Contains(projectStdout.String(), filepath.Join(repoDir, "project-artifacts")) {
		t.Fatalf("expected project artifact dir to win over user config, got %q", projectStdout.String())
	}

	var cliStdout bytes.Buffer
	var cliStderr bytes.Buffer
	cliCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", "cli-artifacts", "--no-input"}, &cliStdout, &cliStderr)

	if cliCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", cliCode)
	}
	if !strings.Contains(cliStdout.String(), filepath.Join(repoDir, "cli-artifacts")) {
		t.Fatalf("expected CLI artifact dir to win over project config, got %q", cliStdout.String())
	}
}

func TestRunFetchRejectsDuplicateActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2 for duplicate Active Run, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "existing run_id="+active.ID) {
		t.Fatalf("expected existing run id in stderr, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunOperationalCommandReportsDirtyWorktreePreflight(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, preflight.DirtyWorktreeError{
			ArtifactDir: "/repo/.roundfix",
			Changes: []preflight.ChangedPath{
				{Status: "M", Path: "src/app.go"},
			},
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "M src/app.go") {
		t.Fatalf("expected changed path and status, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "commit, stash, or remove") {
		t.Fatalf("expected user action guidance, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunPreflightFailureLeavesBufferOutputPlainByDefault(t *testing.T) {
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("plain preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if strings.Contains(stderr.String(), "\x1b[") {
		t.Fatalf("did not expect ANSI color in buffer output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Preflight failed") || !strings.Contains(stderr.String(), "plain preflight failure") {
		t.Fatalf("expected formatted preflight failure, got %q", stderr.String())
	}
}

func TestRunPreflightFailureColorsOutputWhenForced(t *testing.T) {
	t.Setenv("ROUNDFIX_COLOR", "always")
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("colored preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "\x1b[31m") {
		t.Fatalf("expected forced ANSI color in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "colored preflight failure") {
		t.Fatalf("expected original error text, got %q", stderr.String())
	}
}

func TestRunOperationalCommandRejectsInvalidConfigAndArtifactDirectory(t *testing.T) {
	t.Run("invalid user config YAML", func(t *testing.T) {
		homeDir, _ := withCLIWorkspace(t)
		mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
		mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), "defaults:\n  agent: [")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "parse config") {
			t.Fatalf("expected config parse error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})

	t.Run("artifact path is a file", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		artifactFile := filepath.Join(repoDir, "artifact-file")
		mustWrite(t, artifactFile, "not a directory")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", artifactFile, "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "is not a directory") {
			t.Fatalf("expected artifact dir error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})
}

func withCLIWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	mustMkdir(t, filepath.Join(repoDir, ".git"))
	t.Setenv("HOME", homeDir)
	t.Chdir(repoDir)
	return homeDir, repoDir
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func withSuccessfulPreflight(t *testing.T, repoDir string) {
	t.Helper()
	withAgentRunner(t, &fakeAgentRunner{})
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	clock := &fakeWatchClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	withWatchTiming(t, clock, &fakeWatchSleeper{clock: clock})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{{State: watch.StatusSettled}},
	}).Status)
	withFetchReviewItems(t, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal: `$(rm -rf /)`.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		if req.pr == "" {
			return preflight.Result{}, errors.New("missing pr in test preflight")
		}
		return preflight.Result{
			Git: preflight.GitState{
				Root:            repoDir,
				Branch:          "feature/review",
				HEAD:            "abc123",
				UpstreamRemote:  "origin",
				UpstreamBranch:  "feature/review",
				UnpushedCommits: 1,
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{
				Enabled: req.name != "fetch",
				Remote:  "origin",
				Branch:  "feature/review",
				Command: []string{"git", "push", "origin", "HEAD:feature/review"},
			},
		}, nil
	})
}

func withFetchReviewItems(t *testing.T, items []reviewsource.ReviewItem) {
	t.Helper()
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		return items, nil
	})
}

func withFetchReviewItemsFunc(t *testing.T, fn func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error)) {
	t.Helper()
	old := fetchReviewItems
	fetchReviewItems = fn
	t.Cleanup(func() {
		fetchReviewItems = old
	})
}

func withPreflight(t *testing.T, fn func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error)) {
	t.Helper()
	old := runCommandPreflight
	runCommandPreflight = fn
	t.Cleanup(func() {
		runCommandPreflight = old
	})
}

func withAgentRunner(t *testing.T, runner agent.Runner) {
	t.Helper()
	old := runAgentRuntime
	runAgentRuntime = runner
	t.Cleanup(func() {
		runAgentRuntime = old
	})
}

func withVerifier(t *testing.T, verifier daemon.Verifier) {
	t.Helper()
	old := runVerificationGate
	runVerificationGate = verifier
	t.Cleanup(func() {
		runVerificationGate = old
	})
}

func withCommitter(t *testing.T, committer daemon.Committer) {
	t.Helper()
	old := createBatchCommit
	createBatchCommit = committer
	t.Cleanup(func() {
		createBatchCommit = old
	})
}

func withSourceResolver(t *testing.T, resolver *fakeSourceResolver) {
	t.Helper()
	old := resolveReviewSourceIssues
	resolveReviewSourceIssues = resolver.Resolve
	t.Cleanup(func() {
		resolveReviewSourceIssues = old
	})
}

func withPusher(t *testing.T, pusher daemon.Pusher) {
	t.Helper()
	old := runFinalPush
	runFinalPush = pusher
	t.Cleanup(func() {
		runFinalPush = old
	})
}

func withWatchStatus(t *testing.T, fn func(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error)) {
	t.Helper()
	old := watchReviewStatus
	watchReviewStatus = fn
	t.Cleanup(func() {
		watchReviewStatus = old
	})
}

func withWatchTiming(t *testing.T, clock watch.Clock, sleeper watch.Sleeper) {
	t.Helper()
	oldClock := watchClock
	oldSleeper := watchSleeper
	watchClock = clock
	watchSleeper = sleeper
	t.Cleanup(func() {
		watchClock = oldClock
		watchSleeper = oldSleeper
	})
}

func withInteractiveInput(t *testing.T, fn func(context.Context, roundtui.InputRequest) (roundtui.CommandValues, error)) {
	t.Helper()
	old := collectInteractiveInput
	collectInteractiveInput = fn
	t.Cleanup(func() {
		collectInteractiveInput = old
	})
}

func withCurrentPullRequestSuggestion(t *testing.T, value string) {
	t.Helper()
	old := suggestCurrentPullRequest
	suggestCurrentPullRequest = func(context.Context, string) (string, error) {
		return value, nil
	}
	t.Cleanup(func() {
		suggestCurrentPullRequest = old
	})
}

func withChangedPaths(t *testing.T, changes []preflight.ChangedPath) {
	t.Helper()
	old := inspectChangedPaths
	inspectChangedPaths = func(context.Context, string) ([]preflight.ChangedPath, error) {
		return changes, nil
	}
	t.Cleanup(func() {
		inspectChangedPaths = old
	})
}

func readInteractiveDefaults(t *testing.T, homeDir string) store.InteractiveDefaults {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for defaults: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for defaults: %v", err)
		}
	}()
	defaults, err := runStore.InteractiveDefaults(context.Background())
	if err != nil {
		t.Fatalf("read defaults: %v", err)
	}
	return defaults
}

func assertRunCount(t *testing.T, dbPath string, expected int) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, filepath.Dir(filepath.Dir(dbPath)))
	if err != nil {
		t.Fatalf("open store for count: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for count: %v", err)
		}
	}()
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != expected {
		t.Fatalf("expected %d Run record(s), got %d", expected, count)
	}
}

func assertNoRunDatabase(t *testing.T, homeDir string) {
	t.Helper()
	if _, err := os.Stat(store.DatabasePath(homeDir)); err == nil {
		t.Fatalf("expected no Run Database at %s", store.DatabasePath(homeDir))
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Database: %v", err)
	}
}

func assertNoActiveRun(t *testing.T, homeDir string, headRepository string, headBranch string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for active run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for active run: %v", err)
		}
	}()
	if active, found, err := runStore.ActiveRun(context.Background(), headRepository, headBranch); err != nil {
		t.Fatalf("lookup active run: %v", err)
	} else if found {
		t.Fatalf("expected no Active Run, got %#v", active)
	}
}

func assertAgentLogContains(t *testing.T, repoDir string, expected string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(repoDir, ".roundfix", "runs", "*", "agent", "batch-001.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one Agent log, got %#v", matches)
	}
	content, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read Agent log %s: %v", matches[0], err)
	}
	if !strings.Contains(string(content), expected) {
		t.Fatalf("expected Agent log to contain %q, got %s", expected, string(content))
	}
}

func persistCLIReviewIssue(t *testing.T, repoDir string, roundNumber int, headBranch string) rounds.PersistResult {
	t.Helper()
	return persistCLIReviewItems(t, repoDir, roundNumber, headBranch, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		},
	})
}

func persistCLIReviewItems(t *testing.T, repoDir string, roundNumber int, headBranch string, items []reviewsource.ReviewItem) rounds.PersistResult {
	t.Helper()
	result, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     headBranch,
		HeadSHA:        "abc123",
		Round:          roundNumber,
		CreatedAt:      time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		Items:          items,
	})
	if err != nil {
		t.Fatalf("persist CLI Review Issue artifact: %v", err)
	}
	return result
}

type fakeVerifier struct {
	err      error
	calls    int
	commands []string
}

func (verifier *fakeVerifier) Verify(_ context.Context, req daemon.VerifyRequest) error {
	verifier.calls++
	verifier.commands = append(verifier.commands, req.Command)
	if req.Stream != nil {
		if _, err := io.WriteString(req.Stream, "fake verification output\n"); err != nil {
			return err
		}
	}
	return verifier.err
}

type fakeCommitter struct {
	err      error
	calls    int
	workDirs []string
	messages []string
}

func (committer *fakeCommitter) Commit(_ context.Context, req daemon.CommitRequest) error {
	committer.calls++
	committer.workDirs = append(committer.workDirs, req.WorkDir)
	committer.messages = append(committer.messages, req.Message)
	return committer.err
}

type fakeSourceResolver struct {
	err      error
	calls    int
	requests []reviewsource.ResolveRequest
}

func (resolver *fakeSourceResolver) Resolve(_ context.Context, req reviewsource.ResolveRequest) error {
	resolver.calls++
	resolver.requests = append(resolver.requests, req)
	return resolver.err
}

type fakePusher struct {
	err      error
	calls    int
	workDirs []string
	remotes  []string
	branches []string
	args     []string
}

func (pusher *fakePusher) Push(_ context.Context, req daemon.PushRequest) error {
	pusher.calls++
	pusher.workDirs = append(pusher.workDirs, req.WorkDir)
	pusher.remotes = append(pusher.remotes, req.Remote)
	pusher.branches = append(pusher.branches, req.Branch)
	pusher.args = append(pusher.args, "push", req.Remote, "HEAD:"+req.Branch)
	return pusher.err
}

type fakeWatchStatus struct {
	err      error
	calls    int
	statuses []reviewsource.WatchStatus
}

func (source *fakeWatchStatus) Status(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
	source.calls++
	if source.err != nil {
		return reviewsource.WatchStatus{}, source.err
	}
	if len(source.statuses) == 0 {
		return reviewsource.WatchStatus{State: watch.StatusSettled}, nil
	}
	status := source.statuses[0]
	if len(source.statuses) > 1 {
		source.statuses = source.statuses[1:]
	}
	return status, nil
}

type fakeWatchClock struct {
	now time.Time
}

func (clock *fakeWatchClock) Now() time.Time {
	return clock.now
}

func (clock *fakeWatchClock) Advance(duration time.Duration) {
	clock.now = clock.now.Add(duration)
}

type fakeWatchSleeper struct {
	clock *fakeWatchClock
}

func (sleeper *fakeWatchSleeper) Sleep(_ context.Context, duration time.Duration) error {
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	return nil
}

type fakeAgentRunner struct {
	probeErr error
	runErr   error
	status   string
}

func (runner *fakeAgentRunner) Probe(context.Context, agent.RuntimeSpec) error {
	return runner.probeErr
}

func (runner *fakeAgentRunner) Run(_ context.Context, req agent.ExecuteRequest, stream io.Writer) (agent.ExecuteResult, error) {
	status := runner.status
	if status == "" {
		status = rounds.StatusResolved
	}
	for _, issue := range req.Batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, status, ""); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	output := "fake agent output\n"
	if _, err := io.WriteString(stream, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
		return agent.ExecuteResult{}, err
	}
	if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
		return agent.ExecuteResult{}, err
	}
	if runner.runErr != nil {
		return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, runner.runErr
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, nil
}

type fakeStoppingAgentRunner struct{}

func (runner *fakeStoppingAgentRunner) Probe(context.Context, agent.RuntimeSpec) error {
	return nil
}

func (runner *fakeStoppingAgentRunner) Run(_ context.Context, req agent.ExecuteRequest, stream io.Writer) (agent.ExecuteResult, error) {
	output := "partial agent output\n"
	if _, err := io.WriteString(stream, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
		return agent.ExecuteResult{}, err
	}
	if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
		return agent.ExecuteResult{}, err
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, agent.StopError{
		LogPath: req.LogPath,
		Output:  output,
		Err:     context.Canceled,
	}
}
