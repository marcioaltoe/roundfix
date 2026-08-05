// Suite: operational Preflight Validation
// Invariant: review Runs refuse unsafe state before any Run-side mutation can begin.
// Boundary IN: Git state, Open Pull Request metadata, push planning, and review-request coherence.
// Boundary OUT: CLI exit-code rendering and Run execution, owned by internal/cli.
package preflight

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestInspectGitDetectsRepositoryState(t *testing.T) {
	runner := fakeGitRunner{
		outputs: map[string]string{
			gitKey("rev-parse", "--show-toplevel"):                              "/repo",
			gitKey("branch", "--show-current"):                                  "feature/review",
			gitKey("rev-parse", "HEAD"):                                         "abc123",
			gitKey("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"): "origin/feature/review",
			gitKey("rev-list", "--count", "@{u}..HEAD"):                         "2",
			gitKey("status", "--porcelain=v1", "-z"):                            " M src/app.go\x00?? docs/note.md\x00",
		},
	}

	state, err := InspectGit(context.Background(), "/repo/subdir", runner)
	if err != nil {
		t.Fatalf("expected git state, got %v", err)
	}

	if state.Root != "/repo" {
		t.Fatalf("expected root /repo, got %q", state.Root)
	}
	if state.Branch != "feature/review" {
		t.Fatalf("expected branch feature/review, got %q", state.Branch)
	}
	if state.HEAD != "abc123" {
		t.Fatalf("expected HEAD abc123, got %q", state.HEAD)
	}
	if state.UpstreamRemote != "origin" || state.UpstreamBranch != "feature/review" {
		t.Fatalf("expected upstream origin/feature/review, got %q/%q", state.UpstreamRemote, state.UpstreamBranch)
	}
	if state.UnpushedCommits != 2 {
		t.Fatalf("expected 2 unpushed commits, got %d", state.UnpushedCommits)
	}
	if len(state.Dirty) != 2 {
		t.Fatalf("expected 2 dirty paths, got %#v", state.Dirty)
	}
	if state.Dirty[0] != (ChangedPath{Status: "M", Path: "src/app.go"}) {
		t.Fatalf("expected modified path, got %#v", state.Dirty[0])
	}
	if state.Dirty[1] != (ChangedPath{Status: "??", Path: "docs/note.md"}) {
		t.Fatalf("expected untracked path, got %#v", state.Dirty[1])
	}
}

func TestRunResolvesExplicitPullRequestAndBuildsNonForcePushPlan(t *testing.T) {
	result, err := Run(context.Background(), Request{
		Command:                CommandResolve,
		WorkDir:                "/repo",
		ArtifactDir:            "/repo/.roundfix",
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		AutoCommit:             true,
		AutoPush:               true,
		RequestReview:          true,
		GitRunner:              cleanGitRunner("origin/feature/review"),
		PullRequestResolver:    failingPullRequestResolver{},
	})
	if err != nil {
		t.Fatalf("expected preflight success, got %v", err)
	}

	if result.PullRequest.HeadBranch != "feature/review" {
		t.Fatalf("expected explicit PR Head Branch, got %q", result.PullRequest.HeadBranch)
	}
	if result.PullRequest.HeadRepository != "owner/project" {
		t.Fatalf("expected explicit Head Repository, got %q", result.PullRequest.HeadRepository)
	}
	if !result.PushPlan.Enabled {
		t.Fatal("expected auto-push plan")
	}
	if result.PushPlan.Force {
		t.Fatal("expected push plan to never force-push")
	}
	gotCommand := strings.Join(result.PushPlan.Command, " ")
	if gotCommand != "git push origin HEAD:feature/review" {
		t.Fatalf("expected boring push command, got %q", gotCommand)
	}
}

func TestRunResolvesPullRequestThroughInjectedResolver(t *testing.T) {
	resolver := &fakePullRequestResolver{
		pr: PullRequest{
			Number:         "123",
			State:          "OPEN",
			HeadBranch:     "feature/review",
			HeadRepository: "owner/project",
		},
	}

	result, err := Run(context.Background(), Request{
		Command:             CommandFetch,
		WorkDir:             "/repo",
		ArtifactDir:         "/repo/.roundfix",
		PRNumber:            "123",
		AutoCommit:          true,
		AutoPush:            true,
		GitRunner:           cleanGitRunner(""),
		PullRequestResolver: resolver,
	})
	if err != nil {
		t.Fatalf("expected preflight success, got %v", err)
	}

	if !resolver.called {
		t.Fatal("expected injected pull request resolver to be used")
	}
	if result.PullRequest.HeadRepository != "owner/project" {
		t.Fatalf("expected resolver Head Repository, got %q", result.PullRequest.HeadRepository)
	}
	if result.PushPlan.Enabled {
		t.Fatal("fetch must not plan a Final Push")
	}
}

func TestGHPullRequestResolverUsesSupportedFieldsAndParsesBaseRepositoryFromURL(t *testing.T) {
	runner := &fakeGHRunner{
		outputs: map[string]string{
			ghKey("pr", "view", "20", "--json", "number,state,url,headRefName,headRepository"): `{
				"number": 20,
				"state": "OPEN",
				"url": "https://github.com/owner/base/pull/20",
				"headRefName": "feature/review",
				"headRepository": {"nameWithOwner": "contributor/fork"}
			}`,
		},
	}

	pr, err := (GHPullRequestResolver{Runner: runner}).ResolvePullRequest(context.Background(), "/repo", "20")

	if err != nil {
		t.Fatalf("expected pull request metadata, got %v", err)
	}
	if pr.BaseRepository != "owner/base" {
		t.Fatalf("expected base repository from PR URL, got %q", pr.BaseRepository)
	}
	if pr.HeadRepository != "contributor/fork" {
		t.Fatalf("expected Head Repository from gh JSON, got %q", pr.HeadRepository)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one gh call, got %#v", runner.calls)
	}
	fields := strings.Join(runner.calls[0], " ")
	if strings.Contains(fields, "baseRepository") {
		t.Fatalf("must not request unsupported baseRepository field, got %q", fields)
	}
	if !strings.Contains(fields, "url") {
		t.Fatalf("expected resolver to request supported url field, got %q", fields)
	}
}

func TestRunRejectsMissingUpstreamWhenAutoPushCanRun(t *testing.T) {
	_, err := Run(context.Background(), Request{
		Command:                CommandWatch,
		WorkDir:                "/repo",
		ArtifactDir:            "/repo/.roundfix",
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		AutoCommit:             true,
		AutoPush:               true,
		RequestReview:          true,
		GitRunner:              cleanGitRunner(""),
	})

	if err == nil {
		t.Fatal("expected missing upstream to fail")
	}
	if !strings.Contains(err.Error(), "auto-push requires upstream remote and branch") {
		t.Fatalf("expected upstream guidance, got %q", err.Error())
	}
}

func TestRunEnforcesReviewRequestCoherence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                          string
		command                       string
		requestReview                 bool
		autoReviewEnabled             bool
		autoIncrementalReview         bool
		autoPauseAfterReviewedCommits *int
		writeCodeRabbitConfig         bool
		wantPushTriggersReview        bool
		wantError                     string
		wantNextAction                string
	}{
		{
			name:                  "resolve refuses finite Review Source defaults without asking",
			command:               CommandResolve,
			autoReviewEnabled:     true,
			autoIncrementalReview: true,
			wantError:             "would strand the Run waiting for a review nobody requests",
			wantNextAction:        "set review_source.request_review to true in Project Config",
		},
		{
			name:                          "resolve runs when zero pause keeps automatic reviews available",
			command:                       CommandResolve,
			autoReviewEnabled:             true,
			autoIncrementalReview:         true,
			autoPauseAfterReviewedCommits: intPointer(0),
			writeCodeRabbitConfig:         true,
			wantPushTriggersReview:        true,
		},
		{
			name:                          "resolve refuses an explicit finite pause without asking",
			command:                       CommandResolve,
			autoReviewEnabled:             true,
			autoIncrementalReview:         true,
			autoPauseAfterReviewedCommits: intPointer(2),
			writeCodeRabbitConfig:         true,
			wantError:                     "would strand the Run waiting for a review nobody requests",
			wantNextAction:                "set review_source.request_review to true in Project Config",
		},
		{
			name:                  "resolve runs when pushes do not trigger reviews and asking is enabled",
			command:               CommandResolve,
			requestReview:         true,
			autoReviewEnabled:     false,
			writeCodeRabbitConfig: true,
		},
		{
			name:                  "resolve refuses a stranded Run",
			command:               CommandResolve,
			autoReviewEnabled:     false,
			writeCodeRabbitConfig: true,
			wantError:             "would strand the Run waiting for a review nobody requests",
			wantNextAction:        "set review_source.request_review to true in Project Config",
		},
		{
			name:                          "resolve refuses duplicate reviews",
			command:                       CommandResolve,
			requestReview:                 true,
			autoReviewEnabled:             true,
			autoIncrementalReview:         true,
			autoPauseAfterReviewedCommits: intPointer(0),
			writeCodeRabbitConfig:         true,
			wantPushTriggersReview:        true,
			wantError:                     "would request a duplicate review after every push",
			wantNextAction:                "set review_source.request_review to false in Project Config",
		},
		{
			name:                  "watch refuses a stranded Run",
			command:               CommandWatch,
			autoReviewEnabled:     false,
			writeCodeRabbitConfig: true,
			wantError:             "would strand the Run waiting for a review nobody requests",
			wantNextAction:        "set review_source.request_review to true in Project Config",
		},
		{
			name:                          "watch refuses duplicate reviews",
			command:                       CommandWatch,
			requestReview:                 true,
			autoReviewEnabled:             true,
			autoIncrementalReview:         true,
			autoPauseAfterReviewedCommits: intPointer(0),
			writeCodeRabbitConfig:         true,
			wantPushTriggersReview:        true,
			wantError:                     "would request a duplicate review after every push",
			wantNextAction:                "set review_source.request_review to false in Project Config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			if tt.writeCodeRabbitConfig {
				content := fmt.Sprintf("reviews:\n  auto_review:\n    enabled: %t\n    auto_incremental_review: %t\n", tt.autoReviewEnabled, tt.autoIncrementalReview)
				if tt.autoPauseAfterReviewedCommits != nil {
					content += fmt.Sprintf("    auto_pause_after_reviewed_commits: %d\n", *tt.autoPauseAfterReviewedCommits)
				}
				if err := os.WriteFile(filepath.Join(repo, codeRabbitConfigName), []byte(content), 0o644); err != nil {
					t.Fatalf("write CodeRabbit config: %v", err)
				}
			}

			_, err := Run(context.Background(), Request{
				Command:                tt.command,
				WorkDir:                repo,
				PRNumber:               "123",
				ExplicitHeadBranch:     "feature/review",
				ExplicitHeadRepository: "owner/project",
				RequestReview:          tt.requestReview,
				GitRunner:              cleanGitRunnerAt(repo, ""),
			})
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected coherent review settings, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected incoherent review settings to fail")
			}
			for _, want := range []string{
				filepath.Join(repo, codeRabbitConfigName),
				fmt.Sprintf("auto_review.enabled=%t", tt.autoReviewEnabled),
				fmt.Sprintf("auto_review.auto_incremental_review=%t", tt.autoIncrementalReview),
				fmt.Sprintf("auto_review.auto_pause_after_reviewed_commits=%d", valueOrDefault(tt.autoPauseAfterReviewedCommits, defaultCodeRabbitAutoPauseAfterReviewedCommits)),
				fmt.Sprintf("pushTriggersReview=%t", tt.wantPushTriggersReview),
				fmt.Sprintf("review_source.request_review=%t", tt.requestReview),
				tt.wantError,
				tt.wantNextAction,
			} {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("expected refusal to contain %q, got %q", want, err.Error())
				}
			}
		})
	}
}

func TestRunExemptsFetchFromReviewRequestCoherence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		requestReview         bool
		autoReviewEnabled     bool
		autoIncrementalReview bool
	}{
		{name: "push review on and asking off", autoReviewEnabled: true, autoIncrementalReview: true},
		{name: "push review off and asking on", requestReview: true},
		{name: "push review off and asking off"},
		{name: "push review on and asking on", requestReview: true, autoReviewEnabled: true, autoIncrementalReview: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := t.TempDir()
			content := fmt.Sprintf("reviews:\n  auto_review:\n    enabled: %t\n    auto_incremental_review: %t\n", tt.autoReviewEnabled, tt.autoIncrementalReview)
			if err := os.WriteFile(filepath.Join(repo, codeRabbitConfigName), []byte(content), 0o644); err != nil {
				t.Fatalf("write CodeRabbit config: %v", err)
			}

			_, err := Run(context.Background(), Request{
				Command:                CommandFetch,
				WorkDir:                repo,
				PRNumber:               "123",
				ExplicitHeadBranch:     "feature/review",
				ExplicitHeadRepository: "owner/project",
				RequestReview:          tt.requestReview,
				GitRunner:              cleanGitRunnerAt(repo, ""),
			})
			if err != nil {
				t.Fatalf("fetch must ignore review request coherence: %v", err)
			}
		})
	}
}

func TestRunTreatsUnreadableCodeRabbitConfigAsDefaults(t *testing.T) {
	t.Parallel()
	repo := t.TempDir()
	configPath := filepath.Join(repo, codeRabbitConfigName)
	if err := os.Mkdir(configPath, 0o755); err != nil {
		t.Fatalf("create unreadable CodeRabbit config path: %v", err)
	}

	_, err := Run(context.Background(), Request{
		Command:                CommandResolve,
		WorkDir:                repo,
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		RequestReview:          false,
		GitRunner:              cleanGitRunnerAt(repo, ""),
	})

	if err == nil {
		t.Fatal("expected finite Review Source defaults without asking to refuse a stranded Run")
	}
	for _, want := range []string{
		configPath,
		"auto_review.enabled=true",
		"auto_review.auto_incremental_review=true",
		"auto_review.auto_pause_after_reviewed_commits=5",
		"pushTriggersReview=false",
		"review_source.request_review=false",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected unreadable-file refusal to contain %q, got %q", want, err.Error())
		}
	}
}

func intPointer(value int) *int {
	return &value
}

func valueOrDefault(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func TestResolvePullRequestRejectsClosedMetadata(t *testing.T) {
	_, err := ResolvePullRequest(context.Background(), ResolvePullRequestRequest{
		Number:  "123",
		WorkDir: "/repo",
		Resolver: &fakePullRequestResolver{
			pr: PullRequest{
				Number:         "123",
				State:          "CLOSED",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
		},
	})

	if err == nil {
		t.Fatal("expected closed pull request to fail")
	}
	if !strings.Contains(err.Error(), "is not open") {
		t.Fatalf("expected open PR error, got %q", err.Error())
	}
}

type fakeGitRunner struct {
	outputs map[string]string
	errors  map[string]error
}

func (runner fakeGitRunner) RunGit(_ context.Context, _ string, args ...string) (string, error) {
	key := gitKey(args...)
	if err := runner.errors[key]; err != nil {
		return "", err
	}
	if value, ok := runner.outputs[key]; ok {
		return value, nil
	}
	return "", errors.New("missing fake git output for " + strings.Join(args, " "))
}

type fakeGHRunner struct {
	outputs map[string]string
	errors  map[string]error
	calls   [][]string
}

func (runner *fakeGHRunner) RunGH(_ context.Context, _ string, args ...string) (string, error) {
	copied := append([]string(nil), args...)
	runner.calls = append(runner.calls, copied)
	key := ghKey(args...)
	if err := runner.errors[key]; err != nil {
		return "", err
	}
	if value, ok := runner.outputs[key]; ok {
		return value, nil
	}
	return "", errors.New("missing fake gh output for " + strings.Join(args, " "))
}

type fakePullRequestResolver struct {
	pr     PullRequest
	err    error
	called bool
}

func (resolver *fakePullRequestResolver) ResolvePullRequest(_ context.Context, _ string, _ string) (PullRequest, error) {
	resolver.called = true
	if resolver.err != nil {
		return PullRequest{}, resolver.err
	}
	return resolver.pr, nil
}

type failingPullRequestResolver struct{}

func (resolver failingPullRequestResolver) ResolvePullRequest(context.Context, string, string) (PullRequest, error) {
	return PullRequest{}, errors.New("resolver should not be called")
}

func cleanGitRunner(upstream string) fakeGitRunner {
	return cleanGitRunnerAt("/repo", upstream)
}

func cleanGitRunnerAt(root string, upstream string) fakeGitRunner {
	runner := fakeGitRunner{
		outputs: map[string]string{
			gitKey("rev-parse", "--show-toplevel"):   root,
			gitKey("branch", "--show-current"):       "feature/review",
			gitKey("rev-parse", "HEAD"):              "abc123",
			gitKey("status", "--porcelain=v1", "-z"): "",
		},
		errors: map[string]error{},
	}
	upstreamKey := gitKey("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if upstream == "" {
		runner.errors[upstreamKey] = errors.New("no upstream")
		return runner
	}
	runner.outputs[upstreamKey] = upstream
	runner.outputs[gitKey("rev-list", "--count", "@{u}..HEAD")] = "0"
	return runner
}

func gitKey(args ...string) string {
	return strings.Join(args, "\x00")
}

func ghKey(args ...string) string {
	return strings.Join(args, "\x00")
}

func TestExecGitRunnerParsesStdoutSeparatelyFromStderr(t *testing.T) {
	repo := initGitRepoForTest(t)
	runner := ExecGitRunner{}

	if _, err := runner.RunGit(context.Background(), repo, "checkout", "-b", "feature/noise"); err != nil {
		t.Fatalf("create branch: %v", err)
	}
	// git checkout reports the branch switch on stderr; parsed output must stay empty.
	output, err := runner.RunGit(context.Background(), repo, "checkout", "main")
	if err != nil {
		t.Fatalf("checkout main: %v", err)
	}
	if output != "" {
		t.Fatalf("expected stderr kept out of parsed output, got %q", output)
	}
}

func TestExecGitRunnerReportsStderrDetailOnFailure(t *testing.T) {
	repo := initGitRepoForTest(t)

	_, err := ExecGitRunner{}.RunGit(context.Background(), repo, "checkout", "missing-branch")

	if err == nil {
		t.Fatal("expected checkout of missing branch to fail")
	}
	if !strings.Contains(err.Error(), "missing-branch") {
		t.Fatalf("expected stderr detail in error, got %v", err)
	}
}

func TestExecGitRunnerDisablesFSMonitorPerInvocation(t *testing.T) {
	repo := initGitRepoForTest(t)
	gittest.AppendConfig(t, repo, "[core]\n\tfsmonitor = true\n")

	value, err := ExecGitRunner{}.RunGit(context.Background(), repo, "config", "core.fsmonitor")

	if err != nil {
		t.Fatalf("read effective fsmonitor config: %v", err)
	}
	if value != "false" {
		t.Fatalf("expected per-invocation fsmonitor override to win, got %q", value)
	}
}

func TestExecGitRunnerStatusDetectionSurvivesFSMonitorNoise(t *testing.T) {
	repo := initGitRepoForTest(t)
	hook := filepath.Join(repo, ".git", "noisy-fsmonitor.sh")
	script := "#!/bin/sh\necho 'fsmonitor_ipc__send_query: noise' >&2\nexit 1\n"
	if err := os.WriteFile(hook, []byte(script), 0o755); err != nil {
		t.Fatalf("write noisy fsmonitor hook: %v", err)
	}
	runGitForSetup(t, repo, "config", "core.fsmonitor", hook)
	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatalf("dirty worktree: %v", err)
	}

	status, err := ExecGitRunner{}.RunGit(context.Background(), repo, "status", "--porcelain=v1", "-z")

	if err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != " M a.txt\x00" {
		t.Fatalf("expected exact porcelain record despite fsmonitor noise, got %q", status)
	}
}

func initGitRepoForTest(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.AppendConfig(t, repo, "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n")
	if err := os.WriteFile(filepath.Join(repo, "a.txt"), []byte("initial\n"), 0o644); err != nil {
		t.Fatalf("write seed file: %v", err)
	}
	runGitForSetup(t, repo, "add", "a.txt")
	runGitForSetup(t, repo, "commit", "-m", "initial")
	return repo
}

func runGitForSetup(t *testing.T, workDir string, args ...string) {
	t.Helper()
	cmdArgs := append([]string{"-C", workDir}, gitConfigArgsForTest()...)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = isolatedGitEnvForTest()
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
}

func gitConfigArgsForTest() []string {
	return gittest.ConfigArgs()
}

func isolatedGitEnvForTest() []string {
	return gittest.IsolatedEnv()
}
