package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const pullRequestJSONFields = "number,url,state,headRefName,headRefOid,mergedAt,mergeCommit"

// PullRequestBoundary is the external publication surface used by the
// delivery engine. Implementations must observe remote state before retrying
// create or merge operations.
type PullRequestBoundary interface {
	RemoteBranchHead(ctx context.Context, remote, branch string) (RemoteHead, bool, error)
	PushBranch(ctx context.Context, remote, branch, head string) (RemoteHead, error)
	FindOrCreatePullRequest(ctx context.Context, req PullRequestRequest) (PullRequestResult, error)
	CurrentHeadChecks(ctx context.Context, number string) (CheckReport, error)
	MergePullRequest(ctx context.Context, number, expectedHead string) (MergeResult, error)
}

type RemoteHead struct {
	Remote string
	Branch string
	SHA    string
}

type PullRequestRequest struct {
	HeadBranch string
	BaseBranch string
	Title      string
	Body       string
}

type PullRequest struct {
	Number      string
	URL         string
	State       string
	HeadBranch  string
	HeadSHA     string
	MergedAt    string
	MergeCommit string
}

type PullRequestResult struct {
	PullRequest PullRequest
	Created     bool
}

type CheckReport struct {
	HeadSHA string
	Checks  []PullRequestCheck
}

type PullRequestCheck struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Bucket   string `json:"bucket"`
	Link     string `json:"link"`
	Workflow string `json:"workflow"`
}

type MergeResult struct {
	PullRequest   PullRequest
	AlreadyMerged bool
}

type PullRequestHeadMismatchError struct {
	Operation string
	Found     string
	Expected  string
}

func (err PullRequestHeadMismatchError) Error() string {
	return fmt.Sprintf("%s: PR Head Branch is at %q, expected %q", err.Operation, err.Found, err.Expected)
}

// CommandResult preserves an invoked command's output and exit status. A
// non-zero exit is data rather than a runner error because gh uses exit 8 to
// report pending pull request checks.
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type CommandRunner interface {
	Run(ctx context.Context, workDir, name string, args ...string) (CommandResult, error)
}

// GitHubCLI implements PullRequestBoundary with git and the authenticated gh
// CLI in one repository working directory.
type GitHubCLI struct {
	WorkDir string
	Runner  CommandRunner
}

var _ PullRequestBoundary = GitHubCLI{}

func NewGitHubCLI(workDir string) GitHubCLI {
	return GitHubCLI{
		WorkDir: workDir,
		Runner:  execCommandRunner{},
	}
}

func (client GitHubCLI) RemoteBranchHead(ctx context.Context, remote, branch string) (RemoteHead, bool, error) {
	remote = strings.TrimSpace(remote)
	branch = strings.TrimSpace(branch)
	if remote == "" {
		return RemoteHead{}, false, errors.New("read remote branch head: remote is required")
	}
	if branch == "" {
		return RemoteHead{}, false, errors.New("read remote branch head: PR Head Branch is required")
	}
	if strings.HasPrefix(remote, "-") {
		return RemoteHead{}, false, fmt.Errorf("read remote branch head: remote %q cannot start with '-'", remote)
	}

	ref := "refs/heads/" + branch
	result, err := client.run(ctx, "git", "ls-remote", "--heads", remote, ref)
	if err != nil {
		return RemoteHead{}, false, fmt.Errorf("read remote branch head: %w", err)
	}
	if result.ExitCode != 0 {
		return RemoteHead{}, false, commandFailure("read remote branch head", result)
	}
	if strings.TrimSpace(result.Stdout) == "" {
		return RemoteHead{}, false, nil
	}
	sha, err := parseRemoteHead(result.Stdout, ref)
	if err != nil {
		return RemoteHead{}, false, err
	}
	return RemoteHead{Remote: remote, Branch: branch, SHA: sha}, true, nil
}

func (client GitHubCLI) PushBranch(ctx context.Context, remote, branch, head string) (RemoteHead, error) {
	remote = strings.TrimSpace(remote)
	branch = strings.TrimSpace(branch)
	head = strings.TrimSpace(head)
	if remote == "" {
		return RemoteHead{}, errors.New("push branch: remote is required")
	}
	if branch == "" {
		return RemoteHead{}, errors.New("push branch: PR Head Branch is required")
	}
	if head == "" {
		return RemoteHead{}, errors.New("push branch: reviewed head is required")
	}
	if strings.HasPrefix(remote, "-") {
		return RemoteHead{}, fmt.Errorf("push branch: remote %q cannot start with '-'", remote)
	}

	ref := "refs/heads/" + branch
	result, err := client.run(ctx, "git", "push", remote, head+":"+ref)
	if err != nil {
		return RemoteHead{}, fmt.Errorf("push branch: %w", err)
	}
	if result.ExitCode != 0 {
		return RemoteHead{}, commandFailure("push branch", result)
	}

	remoteHead, found, err := client.RemoteBranchHead(ctx, remote, branch)
	if err != nil {
		return RemoteHead{}, fmt.Errorf("read remote head after push: %w", err)
	}
	if !found {
		return RemoteHead{}, fmt.Errorf("read remote head after push: git did not report %q", ref)
	}
	return remoteHead, nil
}

func (client GitHubCLI) FindOrCreatePullRequest(ctx context.Context, req PullRequestRequest) (PullRequestResult, error) {
	req.HeadBranch = strings.TrimSpace(req.HeadBranch)
	req.BaseBranch = strings.TrimSpace(req.BaseBranch)
	req.Title = strings.TrimSpace(req.Title)
	if req.HeadBranch == "" {
		return PullRequestResult{}, errors.New("find or create pull request: PR Head Branch is required")
	}
	if req.BaseBranch == "" {
		return PullRequestResult{}, errors.New("find or create pull request: base branch is required")
	}

	args := []string{
		"pr", "list",
		"--head", req.HeadBranch,
		"--base", req.BaseBranch,
		"--state", "open",
		"--limit", "2",
		"--json", pullRequestJSONFields,
	}
	result, err := client.run(ctx, "gh", args...)
	if err != nil {
		return PullRequestResult{}, fmt.Errorf("find open pull request: %w", err)
	}
	if result.ExitCode != 0 {
		return PullRequestResult{}, commandFailure("find open pull request", result)
	}
	existing, err := parsePullRequestList(result.Stdout)
	if err != nil {
		return PullRequestResult{}, fmt.Errorf("parse open pull requests: %w", err)
	}
	switch len(existing) {
	case 1:
		return PullRequestResult{PullRequest: existing[0]}, nil
	case 0:
		// Continue to creation below.
	default:
		return PullRequestResult{}, fmt.Errorf("find open pull request: found %d pull requests for PR Head Branch %q", len(existing), req.HeadBranch)
	}

	if req.Title == "" {
		return PullRequestResult{}, errors.New("create pull request: title is required")
	}
	result, err = client.run(
		ctx,
		"gh",
		"pr", "create",
		"--head", req.HeadBranch,
		"--base", req.BaseBranch,
		"--title", req.Title,
		"--body", req.Body,
	)
	if err != nil {
		return PullRequestResult{}, fmt.Errorf("create pull request: %w", err)
	}
	if result.ExitCode != 0 {
		return PullRequestResult{}, commandFailure("create pull request", result)
	}
	target := strings.TrimSpace(result.Stdout)
	if target == "" {
		return PullRequestResult{}, errors.New("create pull request: gh returned an empty pull request URL")
	}

	created, err := client.pullRequest(ctx, target)
	if err != nil {
		return PullRequestResult{}, fmt.Errorf("read created pull request: %w", err)
	}
	return PullRequestResult{PullRequest: created, Created: true}, nil
}

func (client GitHubCLI) CurrentHeadChecks(ctx context.Context, number string) (CheckReport, error) {
	number = strings.TrimSpace(number)
	if number == "" {
		return CheckReport{}, errors.New("read pull request checks: pull request number is required")
	}

	before, err := client.pullRequest(ctx, number)
	if err != nil {
		return CheckReport{}, fmt.Errorf("read pull request head before checks: %w", err)
	}
	if before.HeadSHA == "" {
		return CheckReport{}, errors.New("read pull request checks: gh returned an empty PR Head Branch revision")
	}

	result, err := client.run(
		ctx,
		"gh",
		"pr", "checks", number,
		"--json", "bucket,link,name,state,workflow",
	)
	if err != nil {
		return CheckReport{}, fmt.Errorf("read pull request checks: %w", err)
	}
	if result.ExitCode != 0 && strings.TrimSpace(result.Stdout) == "" {
		return CheckReport{}, commandFailure("read pull request checks", result)
	}
	var checks []PullRequestCheck
	if err := json.Unmarshal([]byte(result.Stdout), &checks); err != nil {
		if result.ExitCode != 0 {
			return CheckReport{}, commandFailure("read pull request checks", result)
		}
		return CheckReport{}, fmt.Errorf("parse pull request checks: %w", err)
	}
	if checks == nil {
		checks = []PullRequestCheck{}
	}

	after, err := client.pullRequest(ctx, number)
	if err != nil {
		return CheckReport{}, fmt.Errorf("read pull request head after checks: %w", err)
	}
	if after.HeadSHA != before.HeadSHA {
		return CheckReport{}, fmt.Errorf(
			"read pull request checks: PR Head Branch moved from %q to %q while checks were read",
			before.HeadSHA,
			after.HeadSHA,
		)
	}
	return CheckReport{HeadSHA: before.HeadSHA, Checks: checks}, nil
}

func (client GitHubCLI) MergePullRequest(ctx context.Context, number, expectedHead string) (MergeResult, error) {
	number = strings.TrimSpace(number)
	expectedHead = strings.TrimSpace(expectedHead)
	if number == "" {
		return MergeResult{}, errors.New("merge pull request: pull request number is required")
	}
	if expectedHead == "" {
		return MergeResult{}, errors.New("merge pull request: expected PR Head Branch revision is required")
	}

	current, err := client.pullRequest(ctx, number)
	if err != nil {
		return MergeResult{}, fmt.Errorf("read pull request before merge: %w", err)
	}
	if current.HeadSHA != expectedHead {
		return MergeResult{}, PullRequestHeadMismatchError{
			Operation: "merge pull request",
			Found:     current.HeadSHA,
			Expected:  expectedHead,
		}
	}
	if current.isMerged() {
		if current.MergeCommit == "" {
			return MergeResult{}, errors.New("read existing merge: gh returned an empty merge commit")
		}
		return MergeResult{PullRequest: current, AlreadyMerged: true}, nil
	}

	result, err := client.run(
		ctx,
		"gh",
		"pr", "merge", number,
		"--squash",
		"--match-head-commit", expectedHead,
	)
	if err != nil {
		return MergeResult{}, fmt.Errorf("merge pull request: %w", err)
	}
	if result.ExitCode != 0 {
		return MergeResult{}, commandFailure("merge pull request", result)
	}

	merged, err := client.pullRequest(ctx, number)
	if err != nil {
		return MergeResult{}, fmt.Errorf("read pull request after merge: %w", err)
	}
	if !merged.isMerged() {
		return MergeResult{}, errors.New("merge pull request: GitHub did not report the pull request as merged")
	}
	if merged.HeadSHA != expectedHead {
		return MergeResult{}, PullRequestHeadMismatchError{
			Operation: "read merged pull request",
			Found:     merged.HeadSHA,
			Expected:  expectedHead,
		}
	}
	if merged.MergeCommit == "" {
		return MergeResult{}, errors.New("merge pull request: gh returned an empty merge commit")
	}
	return MergeResult{PullRequest: merged}, nil
}

func (client GitHubCLI) pullRequest(ctx context.Context, target string) (PullRequest, error) {
	result, err := client.run(
		ctx,
		"gh",
		"pr", "view", target,
		"--json", pullRequestJSONFields,
	)
	if err != nil {
		return PullRequest{}, err
	}
	if result.ExitCode != 0 {
		return PullRequest{}, commandFailure("read pull request", result)
	}
	return parsePullRequest([]byte(result.Stdout))
}

func (client GitHubCLI) run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	if client.Runner == nil {
		return CommandResult{}, errors.New("command runner is required")
	}
	workDir := strings.TrimSpace(client.WorkDir)
	if workDir == "" {
		workDir = "."
	}
	return client.Runner.Run(ctx, workDir, name, args...)
}

func parseRemoteHead(output, expectedRef string) (string, error) {
	var sha string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[1] != expectedRef {
			continue
		}
		if sha != "" {
			return "", fmt.Errorf("read remote head after push: git returned more than one %q", expectedRef)
		}
		sha = fields[0]
	}
	if sha == "" {
		return "", fmt.Errorf("read remote head after push: git did not report %q", expectedRef)
	}
	return sha, nil
}

type pullRequestPayload struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	State       string `json:"state"`
	HeadRefName string `json:"headRefName"`
	HeadRefOID  string `json:"headRefOid"`
	MergedAt    string `json:"mergedAt"`
	MergeCommit *struct {
		OID string `json:"oid"`
	} `json:"mergeCommit"`
}

func parsePullRequestList(output string) ([]PullRequest, error) {
	var payloads []pullRequestPayload
	if err := json.Unmarshal([]byte(output), &payloads); err != nil {
		return nil, err
	}
	pullRequests := make([]PullRequest, 0, len(payloads))
	for _, payload := range payloads {
		pullRequest, err := pullRequestFromPayload(payload)
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, pullRequest)
	}
	return pullRequests, nil
}

func parsePullRequest(output []byte) (PullRequest, error) {
	var payload pullRequestPayload
	if err := json.Unmarshal(output, &payload); err != nil {
		return PullRequest{}, fmt.Errorf("parse gh pull request metadata: %w", err)
	}
	return pullRequestFromPayload(payload)
}

func pullRequestFromPayload(payload pullRequestPayload) (PullRequest, error) {
	if payload.Number <= 0 {
		return PullRequest{}, errors.New("gh pull request metadata is missing number")
	}
	pullRequest := PullRequest{
		Number:     strconv.Itoa(payload.Number),
		URL:        strings.TrimSpace(payload.URL),
		State:      strings.TrimSpace(payload.State),
		HeadBranch: strings.TrimSpace(payload.HeadRefName),
		HeadSHA:    strings.TrimSpace(payload.HeadRefOID),
		MergedAt:   strings.TrimSpace(payload.MergedAt),
	}
	if payload.MergeCommit != nil {
		pullRequest.MergeCommit = strings.TrimSpace(payload.MergeCommit.OID)
	}
	return pullRequest, nil
}

func (pullRequest PullRequest) isMerged() bool {
	return strings.EqualFold(pullRequest.State, "merged") || pullRequest.MergedAt != "" || pullRequest.MergeCommit != ""
}

func commandFailure(operation string, result CommandResult) error {
	detail := strings.TrimSpace(result.Stderr)
	if detail == "" {
		detail = strings.TrimSpace(result.Stdout)
	}
	if detail == "" {
		detail = fmt.Sprintf("exit code %d", result.ExitCode)
	}
	if len(detail) > 1200 {
		detail = detail[len(detail)-1200:]
	}
	return fmt.Errorf("%s: %s", operation, detail)
}

type execCommandRunner struct{}

func (execCommandRunner) Run(ctx context.Context, workDir, name string, args ...string) (CommandResult, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = workDir
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	result := CommandResult{
		Stdout: strings.TrimRight(stdout.String(), "\n"),
		Stderr: strings.TrimRight(stderr.String(), "\n"),
	}
	if err == nil {
		return result, nil
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return result, ctxErr
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}
	return result, err
}
