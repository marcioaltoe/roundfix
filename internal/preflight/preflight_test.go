package preflight

import (
	"context"
	"errors"
	"strings"
	"testing"
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

func TestRunRejectsDirtyWorktreeOutsideArtifactDirectory(t *testing.T) {
	runner := cleanGitRunner("origin/feature/review")
	runner.outputs[gitKey("status", "--porcelain=v1", "-z")] = " M .roundfix/reviews/pr-123/issue.md\x00?? src/app.go\x00"

	_, err := Run(context.Background(), Request{
		Command:                CommandFetch,
		WorkDir:                "/repo",
		ArtifactDir:            "/repo/.roundfix",
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		AutoCommit:             true,
		AutoPush:               true,
		GitRunner:              runner,
	})

	var dirtyErr DirtyWorktreeError
	if !errors.As(err, &dirtyErr) {
		t.Fatalf("expected DirtyWorktreeError, got %T %v", err, err)
	}
	if len(dirtyErr.Changes) != 1 {
		t.Fatalf("expected one blocking dirty path, got %#v", dirtyErr.Changes)
	}
	if dirtyErr.Changes[0] != (ChangedPath{Status: "??", Path: "src/app.go"}) {
		t.Fatalf("expected only outside-artifact change, got %#v", dirtyErr.Changes[0])
	}
	if !strings.Contains(err.Error(), "commit, stash, or remove") {
		t.Fatalf("expected user action guidance, got %q", err.Error())
	}
}

func TestFetchAllowsDirtyProjectConfig(t *testing.T) {
	runner := cleanGitRunner("origin/feature/review")
	runner.outputs[gitKey("status", "--porcelain=v1", "-z")] = "?? .roundfixrc.yml\x00"

	_, err := Run(context.Background(), Request{
		Command:                CommandFetch,
		WorkDir:                "/repo",
		ArtifactDir:            "/repo/.roundfix",
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		AutoCommit:             true,
		AutoPush:               true,
		GitRunner:              runner,
	})

	if err != nil {
		t.Fatalf("expected fetch to allow dirty Project Config, got %v", err)
	}
}

func TestResolveAllowsDirtyProjectConfig(t *testing.T) {
	runner := cleanGitRunner("origin/feature/review")
	runner.outputs[gitKey("status", "--porcelain=v1", "-z")] = "?? .roundfixrc.yml\x00"

	_, err := Run(context.Background(), Request{
		Command:                CommandResolve,
		WorkDir:                "/repo",
		ArtifactDir:            "/repo/.roundfix",
		PRNumber:               "123",
		ExplicitHeadBranch:     "feature/review",
		ExplicitHeadRepository: "owner/project",
		AutoCommit:             true,
		AutoPush:               true,
		GitRunner:              runner,
	})

	if err != nil {
		t.Fatalf("expected resolve to allow dirty Project Config, got %v", err)
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
		GitRunner:              cleanGitRunner(""),
	})

	if err == nil {
		t.Fatal("expected missing upstream to fail")
	}
	if !strings.Contains(err.Error(), "auto-push requires upstream remote and branch") {
		t.Fatalf("expected upstream guidance, got %q", err.Error())
	}
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
	runner := fakeGitRunner{
		outputs: map[string]string{
			gitKey("rev-parse", "--show-toplevel"):   "/repo",
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
