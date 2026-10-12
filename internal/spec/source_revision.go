package spec

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const SourceRecoveryFetchTimeout = 60 * time.Second

var ErrSourceRevisionUnavailable = errors.New("source_revision is unavailable")

// SourceRecovery names the recovered commit and its source.
type SourceRecovery struct{ Revision, PullRequest, From string }

var sourceRecoveryRevision = regexp.MustCompile(`^(?:[0-9a-fA-F]{40}|[0-9a-fA-F]{64})$`)

var sourceRecoveryNumber = regexp.MustCompile(`^[1-9][0-9]{0,9}$`)
var sourceRecoverySubject = regexp.MustCompile(`\(#([1-9][0-9]{0,9})\)$`)

// PullRequestFromSubject reads only a trailing delivery pull request number.
func PullRequestFromSubject(subject string) (string, bool) {
	match := sourceRecoverySubject.FindStringSubmatch(subject)
	if match == nil {
		return "", false
	}
	return match[1], true
}

// EnsureSourceRevision fetches missing objects from the delivery pull request.
// It never updates refs or FETCH_HEAD, and cancellation kills and waits for Git.
func EnsureSourceRevision(ctx context.Context, repo, remote string, record ArchiveRecord, head, recordPath string) (SourceRecovery, error) {
	result := SourceRecovery{Revision: record.SourceRevision, PullRequest: record.PullRequest}
	if !sourceRecoveryRevision.MatchString(record.SourceRevision) {
		return result, fmt.Errorf("validate source revision: %w", ErrSourceRevisionUnavailable)
	}
	available := func() error {
		return exec.CommandContext(ctx, "git", "-C", repo, "cat-file", "-e", record.SourceRevision+"^{commit}").Run()
	}
	if available() == nil {
		result.From = "local"
		return result, nil
	}
	if result.PullRequest == "" {
		subject, err := exec.CommandContext(ctx, "git", "-C", repo, "log", "-1", "--diff-filter=A", "--format=%s", head, "--", recordPath).Output()
		if err != nil {
			return result, fmt.Errorf("read record adding commit: %w: %w", ErrSourceRevisionUnavailable, err)
		}
		result.PullRequest, _ = PullRequestFromSubject(strings.TrimSpace(string(subject)))
	}
	if !sourceRecoveryNumber.MatchString(result.PullRequest) {
		return result, fmt.Errorf("resolve pull request number: %w", ErrSourceRevisionUnavailable)
	}
	if remote == "" {
		remote = "origin"
	}
	fetchCtx, cancel := context.WithTimeout(ctx, SourceRecoveryFetchTimeout)
	defer cancel()
	// An empty refmap also prevents configured mappings from creating a ref.
	cmd := exec.CommandContext(fetchCtx, "git", "-C", repo, "fetch", "--no-tags", "--no-write-fetch-head", "--refmap=", "--no-auto-maintenance", "--", remote, "refs/pull/"+result.PullRequest+"/head")
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if err := cmd.Run(); err != nil {
		return result, fmt.Errorf("fetch pull request head: %w: %w", ErrSourceRevisionUnavailable, err)
	}
	if err := available(); err != nil {
		return result, fmt.Errorf("check fetched source revision: %w: %w", ErrSourceRevisionUnavailable, err)
	}
	result.From = "pull-request"
	return result, nil
}
