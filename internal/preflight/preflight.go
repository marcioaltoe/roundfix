package preflight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CommandFetch   = "fetch"
	CommandResolve = "resolve"
	CommandWatch   = "watch"
)

type Request struct {
	Command                string
	WorkDir                string
	ArtifactDir            string
	PRNumber               string
	ExplicitBaseRepository string
	ExplicitHeadBranch     string
	ExplicitHeadRepository string
	AutoCommit             bool
	AutoPush               bool
	PushRemote             string
	PushBranch             string
	GitRunner              GitRunner
	PullRequestResolver    PullRequestResolver
}

type Result struct {
	Git         GitState
	PullRequest PullRequest
	PushPlan    PushPlan
}

type GitState struct {
	Root            string
	Branch          string
	HEAD            string
	UpstreamRemote  string
	UpstreamBranch  string
	UnpushedCommits int
	Dirty           []ChangedPath
}

type ChangedPath struct {
	Status string
	Path   string
}

type PullRequest struct {
	Number         string
	State          string
	BaseRepository string
	HeadBranch     string
	HeadRepository string
}

type PushPlan struct {
	Enabled bool
	Remote  string
	Branch  string
	Command []string
	Force   bool
}

type GitRunner interface {
	RunGit(ctx context.Context, workDir string, args ...string) (string, error)
}

type PullRequestResolver interface {
	ResolvePullRequest(ctx context.Context, workDir string, number string) (PullRequest, error)
}

type GHRunner interface {
	RunGH(ctx context.Context, workDir string, args ...string) (string, error)
}

type ExecGitRunner struct{}

func (runner ExecGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	if workDir == "" {
		workDir = "."
	}
	gitArgs := append([]string{"-C", workDir}, args...)
	cmd := exec.CommandContext(ctx, "git", gitArgs...)
	output, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(output), "\n")
	if err != nil {
		detail := strings.TrimSpace(text)
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return text, nil
}

type ExecGHRunner struct{}

func (runner ExecGHRunner) RunGH(ctx context.Context, workDir string, args ...string) (string, error) {
	if workDir == "" {
		workDir = "."
	}
	cmd := exec.CommandContext(ctx, "gh", args...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	text := strings.TrimRight(string(output), "\n")
	if err != nil {
		detail := strings.TrimSpace(text)
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("gh %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return text, nil
}

type GHPullRequestResolver struct {
	Runner GHRunner
}

func (resolver GHPullRequestResolver) ResolvePullRequest(ctx context.Context, workDir string, number string) (PullRequest, error) {
	if workDir == "" {
		workDir = "."
	}
	runner := resolver.Runner
	if runner == nil {
		runner = ExecGHRunner{}
	}
	output, err := runner.RunGH(ctx, workDir, "pr", "view", number, "--json", "number,state,url,headRefName,headRepository")
	if err != nil {
		return PullRequest{}, err
	}

	var raw struct {
		Number         int    `json:"number"`
		State          string `json:"state"`
		URL            string `json:"url"`
		HeadRefName    string `json:"headRefName"`
		HeadRepository struct {
			NameWithOwner string `json:"nameWithOwner"`
		} `json:"headRepository"`
	}
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		return PullRequest{}, fmt.Errorf("parse gh pull request metadata: %w", err)
	}
	baseRepository, err := baseRepositoryFromPullRequestURL(raw.URL)
	if err != nil {
		return PullRequest{}, fmt.Errorf("parse gh pull request URL: %w", err)
	}
	prNumber := number
	if raw.Number > 0 {
		prNumber = strconv.Itoa(raw.Number)
	}
	return PullRequest{
		Number:         prNumber,
		State:          raw.State,
		BaseRepository: baseRepository,
		HeadBranch:     raw.HeadRefName,
		HeadRepository: raw.HeadRepository.NameWithOwner,
	}, nil
}

func baseRepositoryFromPullRequestURL(rawURL string) (string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "pull" {
		return "", fmt.Errorf("expected GitHub pull request URL path /owner/repo/pull/<number>, got %q", parsed.Path)
	}
	if parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("expected owner and repository in pull request URL path %q", parsed.Path)
	}
	return parts[0] + "/" + parts[1], nil
}

type DirtyWorktreeError struct {
	ArtifactDir string
	Changes     []ChangedPath
}

func (err DirtyWorktreeError) Error() string {
	var builder strings.Builder
	builder.WriteString("dirty worktree outside the Artifact Directory")
	if err.ArtifactDir != "" {
		builder.WriteString(fmt.Sprintf(" %q", err.ArtifactDir))
	}
	builder.WriteString("; commit, stash, or remove these changes before running Roundfix again:")
	for _, change := range err.Changes {
		builder.WriteString(fmt.Sprintf("\n  %s %s", change.Status, change.Path))
	}
	return builder.String()
}

func Run(ctx context.Context, req Request) (Result, error) {
	gitRunner := req.GitRunner
	if gitRunner == nil {
		gitRunner = ExecGitRunner{}
	}
	prResolver := req.PullRequestResolver
	if prResolver == nil {
		prResolver = GHPullRequestResolver{}
	}

	gitState, err := InspectGit(ctx, req.WorkDir, gitRunner)
	if err != nil {
		return Result{}, err
	}

	blockingChanges := DirtyOutsideArtifact(gitState.Root, req.ArtifactDir, gitState.Dirty)
	if len(blockingChanges) > 0 {
		return Result{}, DirtyWorktreeError{
			ArtifactDir: req.ArtifactDir,
			Changes:     blockingChanges,
		}
	}

	pullRequest, err := ResolvePullRequest(ctx, ResolvePullRequestRequest{
		Number:                 req.PRNumber,
		ExplicitBaseRepository: req.ExplicitBaseRepository,
		ExplicitHeadBranch:     req.ExplicitHeadBranch,
		ExplicitHeadRepository: req.ExplicitHeadRepository,
		WorkDir:                gitState.Root,
		Resolver:               prResolver,
	})
	if err != nil {
		return Result{}, err
	}

	pushPlan, err := BuildPushPlan(PushPlanRequest{
		Command:    req.Command,
		AutoCommit: req.AutoCommit,
		AutoPush:   req.AutoPush,
		PushRemote: req.PushRemote,
		PushBranch: req.PushBranch,
		Git:        gitState,
	})
	if err != nil {
		return Result{}, err
	}

	return Result{
		Git:         gitState,
		PullRequest: pullRequest,
		PushPlan:    pushPlan,
	}, nil
}

func InspectGit(ctx context.Context, workDir string, runner GitRunner) (GitState, error) {
	if runner == nil {
		runner = ExecGitRunner{}
	}
	root, err := runner.RunGit(ctx, workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return GitState{}, fmt.Errorf("detect Git root: %w", err)
	}
	root = filepath.Clean(strings.TrimSpace(root))
	if root == "." || root == "" {
		return GitState{}, errors.New("detect Git root: git returned an empty root")
	}

	branch, err := runner.RunGit(ctx, root, "branch", "--show-current")
	if err != nil {
		return GitState{}, fmt.Errorf("detect current branch: %w", err)
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return GitState{}, errors.New("detect current branch: HEAD is detached")
	}

	head, err := runner.RunGit(ctx, root, "rev-parse", "HEAD")
	if err != nil {
		return GitState{}, fmt.Errorf("detect current HEAD: %w", err)
	}
	head = strings.TrimSpace(head)
	if head == "" {
		return GitState{}, errors.New("detect current HEAD: git returned an empty SHA")
	}

	state := GitState{
		Root:   root,
		Branch: branch,
		HEAD:   head,
	}

	upstream, err := runner.RunGit(ctx, root, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err == nil {
		state.UpstreamRemote, state.UpstreamBranch = ParseUpstream(strings.TrimSpace(upstream))
		if state.UpstreamRemote != "" && state.UpstreamBranch != "" {
			countText, err := runner.RunGit(ctx, root, "rev-list", "--count", "@{u}..HEAD")
			if err != nil {
				return GitState{}, fmt.Errorf("detect unpushed commit count: %w", err)
			}
			count, err := strconv.Atoi(strings.TrimSpace(countText))
			if err != nil {
				return GitState{}, fmt.Errorf("parse unpushed commit count %q: %w", countText, err)
			}
			state.UnpushedCommits = count
		}
	}

	status, err := runner.RunGit(ctx, root, "status", "--porcelain=v1", "-z")
	if err != nil {
		return GitState{}, fmt.Errorf("detect dirty worktree status: %w", err)
	}
	dirty, err := ParsePorcelainStatus(status)
	if err != nil {
		return GitState{}, err
	}
	state.Dirty = dirty
	return state, nil
}

func ParseUpstream(upstream string) (string, string) {
	remote, branch, ok := strings.Cut(upstream, "/")
	if !ok {
		return upstream, ""
	}
	return remote, branch
}

func ParsePorcelainStatus(raw string) ([]ChangedPath, error) {
	if raw == "" {
		return nil, nil
	}
	entries := strings.Split(raw, "\x00")
	changes := make([]ChangedPath, 0, len(entries))
	for index := 0; index < len(entries); index++ {
		entry := entries[index]
		if entry == "" {
			continue
		}
		if len(entry) < 4 {
			return nil, fmt.Errorf("parse git status entry %q: expected status and path", entry)
		}
		rawStatus := entry[:2]
		status := strings.TrimSpace(rawStatus)
		if status == "" {
			status = rawStatus
		}
		path := strings.TrimSpace(entry[3:])
		if path == "" {
			return nil, fmt.Errorf("parse git status entry %q: empty path", entry)
		}
		changes = append(changes, ChangedPath{
			Status: status,
			Path:   path,
		})
		if strings.ContainsAny(rawStatus, "RC") && index+1 < len(entries) {
			index++
		}
	}
	return changes, nil
}

func DirtyOutsideArtifact(gitRoot string, artifactDir string, dirty []ChangedPath) []ChangedPath {
	if len(dirty) == 0 {
		return nil
	}
	if artifactDir == "" {
		return dirty
	}
	root := filepath.Clean(gitRoot)
	artifact := filepath.Clean(artifactDir)
	blocking := make([]ChangedPath, 0, len(dirty))
	for _, change := range dirty {
		path := change.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		if isWithinPath(path, artifact) {
			continue
		}
		blocking = append(blocking, change)
	}
	return blocking
}

type ResolvePullRequestRequest struct {
	Number                 string
	ExplicitBaseRepository string
	ExplicitHeadBranch     string
	ExplicitHeadRepository string
	WorkDir                string
	Resolver               PullRequestResolver
}

func ResolvePullRequest(ctx context.Context, req ResolvePullRequestRequest) (PullRequest, error) {
	if req.Number == "" {
		return PullRequest{}, errors.New("Open Pull Request number is required")
	}
	if req.ExplicitHeadBranch != "" || req.ExplicitHeadRepository != "" || req.ExplicitBaseRepository != "" {
		if req.ExplicitHeadBranch == "" || req.ExplicitHeadRepository == "" {
			return PullRequest{}, errors.New("explicit pull request metadata requires both --head-branch and --head-repo")
		}
		baseRepository := req.ExplicitBaseRepository
		if baseRepository == "" {
			baseRepository = req.ExplicitHeadRepository
		}
		return PullRequest{
			Number:         req.Number,
			State:          "OPEN",
			BaseRepository: baseRepository,
			HeadBranch:     req.ExplicitHeadBranch,
			HeadRepository: req.ExplicitHeadRepository,
		}, nil
	}
	resolver := req.Resolver
	if resolver == nil {
		resolver = GHPullRequestResolver{}
	}
	pullRequest, err := resolver.ResolvePullRequest(ctx, req.WorkDir, req.Number)
	if err != nil {
		return PullRequest{}, fmt.Errorf("resolve Open Pull Request %s: %w", req.Number, err)
	}
	if !strings.EqualFold(pullRequest.State, "open") {
		return PullRequest{}, fmt.Errorf("Open Pull Request %s is not open; detected state %q", req.Number, pullRequest.State)
	}
	if strings.TrimSpace(pullRequest.HeadBranch) == "" {
		return PullRequest{}, fmt.Errorf("Open Pull Request %s is missing PR Head Branch metadata", req.Number)
	}
	if strings.TrimSpace(pullRequest.HeadRepository) == "" {
		return PullRequest{}, fmt.Errorf("Open Pull Request %s is missing Head Repository metadata", req.Number)
	}
	if strings.TrimSpace(pullRequest.BaseRepository) == "" {
		pullRequest.BaseRepository = pullRequest.HeadRepository
	}
	if pullRequest.Number == "" {
		pullRequest.Number = req.Number
	}
	return pullRequest, nil
}

type PushPlanRequest struct {
	Command    string
	AutoCommit bool
	AutoPush   bool
	PushRemote string
	PushBranch string
	Git        GitState
}

func BuildPushPlan(req PushPlanRequest) (PushPlan, error) {
	if !req.AutoPush || req.Command == CommandFetch {
		return PushPlan{}, nil
	}
	if !req.AutoCommit {
		return PushPlan{}, errors.New("auto-push requires auto-commit to be enabled")
	}
	remote := strings.TrimSpace(req.PushRemote)
	if remote == "" {
		remote = req.Git.UpstreamRemote
	}
	branch := strings.TrimSpace(req.PushBranch)
	if branch == "" {
		branch = req.Git.UpstreamBranch
	}
	if remote == "" || branch == "" {
		return PushPlan{}, errors.New("auto-push requires upstream remote and branch; set the local branch upstream or configure watch.push_remote and watch.push_branch before running Roundfix again")
	}
	return PushPlan{
		Enabled: true,
		Remote:  remote,
		Branch:  branch,
		Command: []string{"git", "push", remote, "HEAD:" + branch},
		Force:   false,
	}, nil
}

func isWithinPath(path string, dir string) bool {
	rel, err := filepath.Rel(dir, filepath.Clean(path))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
