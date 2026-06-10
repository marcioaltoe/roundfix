package coderabbit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/watch"
)

const (
	coderabbitBotLogin = "coderabbitai[bot]"
	coderabbitAppLogin = "coderabbitai"
)

type Client struct {
	GitHub GitHubClient
}

type GitHubClient interface {
	ReviewComments(ctx context.Context, repo string, prNumber string) ([]ReviewComment, error)
	ReviewThreads(ctx context.Context, repo string, prNumber string) ([]ReviewThread, error)
	CheckRuns(ctx context.Context, repo string, headSHA string) ([]CheckRun, error)
	CommitStatuses(ctx context.Context, repo string, headSHA string) ([]CommitStatus, error)
	PullRequestReviews(ctx context.Context, repo string, prNumber string) ([]PullRequestReview, error)
	ResolveReviewThread(ctx context.Context, threadID string) error
}

type ReviewComment struct {
	DatabaseID              int64
	NodeID                  string
	Body                    string
	Path                    string
	Line                    int
	Author                  string
	SourceReviewID          string
	SourceReviewSubmittedAt time.Time
}

type ReviewThread struct {
	ID         string
	IsResolved bool
	Comments   []ThreadComment
}

type ThreadComment struct {
	DatabaseID              int64
	NodeID                  string
	Author                  string
	SourceReviewID          string
	SourceReviewSubmittedAt time.Time
}

type CheckRun struct {
	Name       string
	AppName    string
	HeadSHA    string
	Status     string
	Conclusion string
}

type CommitStatus struct {
	Context string
	State   string
}

type PullRequestReview struct {
	DatabaseID  int64
	Author      string
	State       string
	CommitSHA   string
	SubmittedAt time.Time
}

func (client Client) FetchReviews(ctx context.Context, req reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
	if strings.TrimSpace(req.BaseRepository) == "" {
		return nil, errors.New("CodeRabbit fetch requires base repository metadata")
	}
	gh := client.GitHub
	if gh == nil {
		gh = GHClient{}
	}

	comments, err := gh.ReviewComments(ctx, req.BaseRepository, req.PRNumber)
	if err != nil {
		return nil, fmt.Errorf("fetch CodeRabbit review comments: %w", err)
	}
	threads, err := gh.ReviewThreads(ctx, req.BaseRepository, req.PRNumber)
	if err != nil {
		return nil, fmt.Errorf("fetch CodeRabbit review threads: %w", err)
	}

	unresolved := unresolvedCommentThreads(threads)
	items := make([]reviewsource.ReviewItem, 0, len(comments))
	for _, comment := range comments {
		if !isCodeRabbitAuthor(comment.Author) {
			continue
		}
		thread, ok := unresolved[commentKey(comment.DatabaseID, comment.NodeID)]
		if !ok {
			continue
		}
		severity := detectSeverity(comment.Body)
		if severity == "nitpick" && !req.IncludeNitpicks {
			continue
		}
		submittedAt := comment.SourceReviewSubmittedAt
		sourceReviewID := comment.SourceReviewID
		if sourceReviewID == "" {
			sourceReviewID = thread.SourceReviewID
		}
		if submittedAt.IsZero() {
			submittedAt = thread.SourceReviewSubmittedAt
		}
		items = append(items, reviewsource.ReviewItem{
			Title:                   titleFromBody(comment.Body),
			File:                    comment.Path,
			Line:                    comment.Line,
			Severity:                severity,
			Author:                  comment.Author,
			Body:                    comment.Body,
			SourceRef:               fmt.Sprintf("thread:%s,comment:%s", thread.ThreadID, commentRef(comment)),
			ReviewHash:              reviewHash(comment.Body),
			SourceReviewID:          sourceReviewID,
			SourceReviewSubmittedAt: submittedAt,
		})
	}
	return items, nil
}

func (client Client) ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error {
	if strings.TrimSpace(req.BaseRepository) == "" {
		return errors.New("CodeRabbit source resolution requires base repository metadata")
	}
	gh := client.GitHub
	if gh == nil {
		gh = GHClient{}
	}

	threadIDs := uniqueThreadIDs(req.Issues)
	for _, threadID := range threadIDs {
		if err := gh.ResolveReviewThread(ctx, threadID); err != nil {
			return fmt.Errorf("resolve CodeRabbit review thread %s: %w", threadID, err)
		}
	}
	return nil
}

func (client Client) WatchStatus(ctx context.Context, req reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
	if err := ctx.Err(); err != nil {
		return reviewsource.WatchStatus{}, err
	}
	if strings.TrimSpace(req.BaseRepository) == "" {
		return reviewsource.WatchStatus{}, errors.New("CodeRabbit watch status requires base repository metadata")
	}
	if strings.TrimSpace(req.HeadSHA) == "" {
		return reviewsource.WatchStatus{}, errors.New("CodeRabbit watch status requires HEAD metadata")
	}
	gh := client.GitHub
	if gh == nil {
		gh = GHClient{}
	}

	checkRuns, err := gh.CheckRuns(ctx, req.BaseRepository, req.HeadSHA)
	if err != nil {
		return reviewsource.WatchStatus{}, fmt.Errorf("fetch CodeRabbit check runs: %w", err)
	}
	statuses, err := gh.CommitStatuses(ctx, req.BaseRepository, req.HeadSHA)
	if err != nil {
		return reviewsource.WatchStatus{}, fmt.Errorf("fetch CodeRabbit commit statuses: %w", err)
	}
	reviews, err := gh.PullRequestReviews(ctx, req.BaseRepository, req.PRNumber)
	if err != nil {
		return reviewsource.WatchStatus{}, fmt.Errorf("fetch CodeRabbit pull request reviews: %w", err)
	}
	return classifyWatchStatus(req.HeadSHA, checkRuns, statuses, reviews), nil
}

type GHClient struct{}

func (client GHClient) ReviewComments(ctx context.Context, repo string, prNumber string) ([]ReviewComment, error) {
	output, err := runGH(ctx, "api", "--paginate", fmt.Sprintf("repos/%s/pulls/%s/comments", repo, prNumber))
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID                  int64  `json:"id"`
		NodeID              string `json:"node_id"`
		Body                string `json:"body"`
		Path                string `json:"path"`
		Line                int    `json:"line"`
		OriginalLine        int    `json:"original_line"`
		PullRequestReviewID int64  `json:"pull_request_review_id"`
		CreatedAt           string `json:"created_at"`
		User                struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse pull request review comments: %w", err)
	}
	comments := make([]ReviewComment, 0, len(raw))
	for _, comment := range raw {
		line := comment.Line
		if line == 0 {
			line = comment.OriginalLine
		}
		createdAt, _ := parseGitHubTime(comment.CreatedAt)
		comments = append(comments, ReviewComment{
			DatabaseID:              comment.ID,
			NodeID:                  comment.NodeID,
			Body:                    comment.Body,
			Path:                    comment.Path,
			Line:                    line,
			Author:                  comment.User.Login,
			SourceReviewID:          strconv.FormatInt(comment.PullRequestReviewID, 10),
			SourceReviewSubmittedAt: createdAt,
		})
	}
	return comments, nil
}

func (client GHClient) ReviewThreads(ctx context.Context, repo string, prNumber string) ([]ReviewThread, error) {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" {
		return nil, fmt.Errorf("repository %q must be owner/name", repo)
	}
	number, err := strconv.Atoi(prNumber)
	if err != nil {
		return nil, fmt.Errorf("pull request %q must be numeric: %w", prNumber, err)
	}

	var all []ReviewThread
	var cursor string
	for {
		response, err := client.reviewThreadsPage(ctx, owner, name, number, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, response.Threads...)
		if !response.HasNextPage {
			return all, nil
		}
		cursor = response.EndCursor
	}
}

func (client GHClient) ResolveReviewThread(ctx context.Context, threadID string) error {
	query := `mutation($threadId:ID!) {
  resolveReviewThread(input: {threadId: $threadId}) {
    thread { id isResolved }
  }
}`
	_, err := runGH(ctx, "api", "graphql", "-f", "query="+query, "-f", "threadId="+threadID)
	return err
}

func (client GHClient) CheckRuns(ctx context.Context, repo string, headSHA string) ([]CheckRun, error) {
	output, err := runGH(ctx, "api", "-H", "Accept: application/vnd.github+json", fmt.Sprintf("repos/%s/commits/%s/check-runs?per_page=100", repo, headSHA))
	if err != nil {
		return nil, err
	}
	var raw struct {
		CheckRuns []struct {
			Name       string `json:"name"`
			HeadSHA    string `json:"head_sha"`
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
			App        *struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"app"`
		} `json:"check_runs"`
	}
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse check runs: %w", err)
	}
	checkRuns := make([]CheckRun, 0, len(raw.CheckRuns))
	for _, run := range raw.CheckRuns {
		appName := ""
		if run.App != nil {
			appName = firstNonEmpty(run.App.Name, run.App.Slug)
		}
		checkRuns = append(checkRuns, CheckRun{
			Name:       run.Name,
			AppName:    appName,
			HeadSHA:    run.HeadSHA,
			Status:     run.Status,
			Conclusion: run.Conclusion,
		})
	}
	return checkRuns, nil
}

func (client GHClient) CommitStatuses(ctx context.Context, repo string, headSHA string) ([]CommitStatus, error) {
	output, err := runGH(ctx, "api", fmt.Sprintf("repos/%s/commits/%s/status", repo, headSHA))
	if err != nil {
		return nil, err
	}
	var raw struct {
		Statuses []struct {
			Context string `json:"context"`
			State   string `json:"state"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse commit statuses: %w", err)
	}
	statuses := make([]CommitStatus, 0, len(raw.Statuses))
	for _, status := range raw.Statuses {
		statuses = append(statuses, CommitStatus{
			Context: status.Context,
			State:   status.State,
		})
	}
	return statuses, nil
}

func (client GHClient) PullRequestReviews(ctx context.Context, repo string, prNumber string) ([]PullRequestReview, error) {
	output, err := runGH(ctx, "api", "--paginate", fmt.Sprintf("repos/%s/pulls/%s/reviews", repo, prNumber))
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID          int64  `json:"id"`
		State       string `json:"state"`
		CommitID    string `json:"commit_id"`
		SubmittedAt string `json:"submitted_at"`
		User        struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	if err := json.Unmarshal(output, &raw); err != nil {
		return nil, fmt.Errorf("parse pull request reviews: %w", err)
	}
	reviews := make([]PullRequestReview, 0, len(raw))
	for _, review := range raw {
		submittedAt, _ := parseGitHubTime(review.SubmittedAt)
		reviews = append(reviews, PullRequestReview{
			DatabaseID:  review.ID,
			Author:      review.User.Login,
			State:       review.State,
			CommitSHA:   review.CommitID,
			SubmittedAt: submittedAt,
		})
	}
	return reviews, nil
}

type reviewThreadsPage struct {
	Threads     []ReviewThread
	HasNextPage bool
	EndCursor   string
}

func (client GHClient) reviewThreadsPage(ctx context.Context, owner string, name string, number int, cursor string) (reviewThreadsPage, error) {
	query := `query($owner:String!, $name:String!, $number:Int!, $cursor:String) {
  repository(owner: $owner, name: $name) {
    pullRequest(number: $number) {
      reviewThreads(first: 100, after: $cursor) {
        pageInfo { hasNextPage endCursor }
        nodes {
          id
          isResolved
          comments(first: 100) {
            nodes {
              databaseId
              id
              author { login }
              pullRequestReview { databaseId submittedAt }
            }
          }
        }
      }
    }
  }
}`
	args := []string{
		"api", "graphql",
		"-f", "query=" + query,
		"-f", "owner=" + owner,
		"-f", "name=" + name,
		"-F", fmt.Sprintf("number=%d", number),
	}
	if cursor != "" {
		args = append(args, "-f", "cursor="+cursor)
	}
	output, err := runGH(ctx, args...)
	if err != nil {
		return reviewThreadsPage{}, err
	}

	var raw struct {
		Data struct {
			Repository struct {
				PullRequest struct {
					ReviewThreads struct {
						PageInfo struct {
							HasNextPage bool   `json:"hasNextPage"`
							EndCursor   string `json:"endCursor"`
						} `json:"pageInfo"`
						Nodes []struct {
							ID         string `json:"id"`
							IsResolved bool   `json:"isResolved"`
							Comments   struct {
								Nodes []struct {
									DatabaseID int64  `json:"databaseId"`
									ID         string `json:"id"`
									Author     struct {
										Login string `json:"login"`
									} `json:"author"`
									PullRequestReview *struct {
										DatabaseID  int64  `json:"databaseId"`
										SubmittedAt string `json:"submittedAt"`
									} `json:"pullRequestReview"`
								} `json:"nodes"`
							} `json:"comments"`
						} `json:"nodes"`
					} `json:"reviewThreads"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
	}
	if err := json.Unmarshal(output, &raw); err != nil {
		return reviewThreadsPage{}, fmt.Errorf("parse pull request review threads: %w", err)
	}

	nodes := raw.Data.Repository.PullRequest.ReviewThreads.Nodes
	threads := make([]ReviewThread, 0, len(nodes))
	for _, node := range nodes {
		thread := ReviewThread{
			ID:         node.ID,
			IsResolved: node.IsResolved,
			Comments:   make([]ThreadComment, 0, len(node.Comments.Nodes)),
		}
		for _, comment := range node.Comments.Nodes {
			threadComment := ThreadComment{
				DatabaseID: comment.DatabaseID,
				NodeID:     comment.ID,
				Author:     comment.Author.Login,
			}
			if comment.PullRequestReview != nil {
				threadComment.SourceReviewID = strconv.FormatInt(comment.PullRequestReview.DatabaseID, 10)
				threadComment.SourceReviewSubmittedAt, _ = parseGitHubTime(comment.PullRequestReview.SubmittedAt)
			}
			thread.Comments = append(thread.Comments, threadComment)
		}
		threads = append(threads, thread)
	}
	pageInfo := raw.Data.Repository.PullRequest.ReviewThreads.PageInfo
	return reviewThreadsPage{
		Threads:     threads,
		HasNextPage: pageInfo.HasNextPage,
		EndCursor:   pageInfo.EndCursor,
	}, nil
}

func classifyWatchStatus(headSHA string, checkRuns []CheckRun, statuses []CommitStatus, reviews []PullRequestReview) reviewsource.WatchStatus {
	if status, ok := classifyCheckRuns(headSHA, checkRuns); ok {
		return status
	}
	if status, ok := classifyCommitStatuses(statuses); ok {
		return status
	}
	if status, ok := classifyPullRequestReviews(headSHA, reviews); ok {
		return status
	}
	return reviewsource.WatchStatus{
		State:  watch.StatusPending,
		Detail: "No CodeRabbit status or review found for the current HEAD",
	}
}

func classifyCheckRuns(headSHA string, checkRuns []CheckRun) (reviewsource.WatchStatus, bool) {
	for _, run := range checkRuns {
		if !isCodeRabbitSignal(run.Name, run.AppName) {
			continue
		}
		if run.HeadSHA != "" && run.HeadSHA != headSHA {
			continue
		}
		status := strings.ToLower(strings.TrimSpace(run.Status))
		conclusion := strings.ToLower(strings.TrimSpace(run.Conclusion))
		switch status {
		case "queued", "in_progress", "pending", "requested", "waiting":
			return reviewsource.WatchStatus{
				State:  watch.StatusReviewing,
				Detail: fmt.Sprintf("CodeRabbit check %q is %s for the current HEAD", run.Name, status),
			}, true
		case "completed":
			return reviewsource.WatchStatus{
				State:  watch.StatusSettled,
				Detail: fmt.Sprintf("CodeRabbit check %q completed for the current HEAD with conclusion %q", run.Name, firstNonEmpty(conclusion, "unknown")),
			}, true
		}
		if conclusion != "" {
			return reviewsource.WatchStatus{
				State:  watch.StatusSettled,
				Detail: fmt.Sprintf("CodeRabbit check %q has conclusion %q for the current HEAD", run.Name, conclusion),
			}, true
		}
		return reviewsource.WatchStatus{
			State:  watch.StatusPending,
			Detail: fmt.Sprintf("CodeRabbit check %q exists for the current HEAD without a terminal status", run.Name),
		}, true
	}
	return reviewsource.WatchStatus{}, false
}

func classifyCommitStatuses(statuses []CommitStatus) (reviewsource.WatchStatus, bool) {
	for _, status := range statuses {
		if !isCodeRabbitSignal(status.Context) {
			continue
		}
		state := strings.ToLower(strings.TrimSpace(status.State))
		switch state {
		case "pending":
			return reviewsource.WatchStatus{
				State:  watch.StatusReviewing,
				Detail: fmt.Sprintf("CodeRabbit commit status %q is pending for the current HEAD", status.Context),
			}, true
		case "success", "failure", "error":
			return reviewsource.WatchStatus{
				State:  watch.StatusSettled,
				Detail: fmt.Sprintf("CodeRabbit commit status %q reached %s for the current HEAD", status.Context, state),
			}, true
		default:
			return reviewsource.WatchStatus{
				State:  watch.StatusPending,
				Detail: fmt.Sprintf("CodeRabbit commit status %q is %q for the current HEAD", status.Context, state),
			}, true
		}
	}
	return reviewsource.WatchStatus{}, false
}

func classifyPullRequestReviews(headSHA string, reviews []PullRequestReview) (reviewsource.WatchStatus, bool) {
	latestCodeRabbitCommit := ""
	for _, review := range reviews {
		if !isCodeRabbitSignal(review.Author) {
			continue
		}
		if review.CommitSHA == headSHA {
			return reviewsource.WatchStatus{
				State:  watch.StatusSettled,
				Detail: fmt.Sprintf("CodeRabbit review %d is published for the current HEAD", review.DatabaseID),
			}, true
		}
		if review.CommitSHA != "" {
			latestCodeRabbitCommit = review.CommitSHA
		}
	}
	if latestCodeRabbitCommit != "" {
		return reviewsource.WatchStatus{
			State:  watch.StatusPending,
			Detail: "Latest CodeRabbit review is for a different commit",
		}, true
	}
	return reviewsource.WatchStatus{}, false
}

func isCodeRabbitSignal(values ...string) bool {
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if strings.Contains(normalized, "coderabbit") {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func runGH(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "gh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("gh %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return output, nil
}

type unresolvedThreadComment struct {
	ThreadID                string
	SourceReviewID          string
	SourceReviewSubmittedAt time.Time
}

func unresolvedCommentThreads(threads []ReviewThread) map[string]unresolvedThreadComment {
	index := make(map[string]unresolvedThreadComment)
	for _, thread := range threads {
		if thread.IsResolved {
			continue
		}
		for _, comment := range thread.Comments {
			if !isCodeRabbitAuthor(comment.Author) {
				continue
			}
			value := unresolvedThreadComment{
				ThreadID:                thread.ID,
				SourceReviewID:          comment.SourceReviewID,
				SourceReviewSubmittedAt: comment.SourceReviewSubmittedAt,
			}
			index[commentKey(comment.DatabaseID, comment.NodeID)] = value
		}
	}
	return index
}

func isCodeRabbitAuthor(author string) bool {
	normalized := strings.ToLower(strings.TrimSpace(author))
	return normalized == coderabbitBotLogin || normalized == coderabbitAppLogin
}

func uniqueThreadIDs(issues []reviewsource.ResolvedIssue) []string {
	seen := map[string]bool{}
	threadIDs := []string{}
	for _, issue := range issues {
		threadID, ok := threadIDFromSourceRef(issue.SourceRef)
		if !ok || seen[threadID] {
			continue
		}
		seen[threadID] = true
		threadIDs = append(threadIDs, threadID)
	}
	return threadIDs
}

func threadIDFromSourceRef(sourceRef string) (string, bool) {
	for _, part := range strings.Split(sourceRef, ",") {
		part = strings.TrimSpace(part)
		threadID, ok := strings.CutPrefix(part, "thread:")
		if ok && threadID != "" {
			return threadID, true
		}
	}
	return "", false
}

func commentKey(databaseID int64, nodeID string) string {
	if databaseID > 0 {
		return fmt.Sprintf("database:%d", databaseID)
	}
	return "node:" + nodeID
}

func commentRef(comment ReviewComment) string {
	if comment.NodeID != "" {
		return comment.NodeID
	}
	if comment.DatabaseID > 0 {
		return strconv.FormatInt(comment.DatabaseID, 10)
	}
	return reviewHash(comment.Body)
}

func titleFromBody(body string) string {
	for _, line := range strings.Split(body, "\n") {
		title := strings.TrimSpace(strings.TrimLeft(line, "# "))
		if title != "" {
			if len(title) > 80 {
				return title[:77] + "..."
			}
			return title
		}
	}
	return "CodeRabbit finding"
}

func detectSeverity(body string) string {
	lower := strings.ToLower(body)
	switch {
	case strings.Contains(lower, "nitpick"):
		return "nitpick"
	case strings.Contains(lower, "major"):
		return "major"
	case strings.Contains(lower, "minor"):
		return "minor"
	default:
		return "review"
	}
}

func reviewHash(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

func parseGitHubTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
