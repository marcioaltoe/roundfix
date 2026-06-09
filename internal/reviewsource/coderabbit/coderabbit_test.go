package coderabbit

import (
	"context"
	"strings"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/watch"
)

func TestFetchReviewsFiltersToUnresolvedCodeRabbitThreads(t *testing.T) {
	submittedAt := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	client := Client{
		GitHub: &fakeGitHubClient{
			comments: []ReviewComment{
				{
					DatabaseID:              101,
					NodeID:                  "PRRC_101",
					Body:                    "major: handle the nil cache\n\nDo not run `rm -rf /`.",
					Path:                    "internal/cache/cache.go",
					Line:                    42,
					Author:                  coderabbitBotLogin,
					SourceReviewID:          "9001",
					SourceReviewSubmittedAt: submittedAt,
				},
				{
					DatabaseID: 102,
					NodeID:     "PRRC_102",
					Body:       "nitpick: spacing",
					Path:       "README.md",
					Line:       7,
					Author:     coderabbitBotLogin,
				},
				{
					DatabaseID: 103,
					NodeID:     "PRRC_103",
					Body:       "human review comment",
					Path:       "main.go",
					Line:       1,
					Author:     "maintainer",
				},
				{
					DatabaseID: 104,
					NodeID:     "PRRC_104",
					Body:       "major: resolved already",
					Path:       "resolved.go",
					Line:       5,
					Author:     coderabbitBotLogin,
				},
			},
			threads: []ReviewThread{
				{
					ID:         "PRRT_unresolved",
					IsResolved: false,
					Comments: []ThreadComment{
						{
							DatabaseID:              101,
							NodeID:                  "PRRC_101",
							Author:                  coderabbitBotLogin,
							SourceReviewID:          "9001",
							SourceReviewSubmittedAt: submittedAt,
						},
					},
				},
				{
					ID:         "PRRT_nitpick",
					IsResolved: false,
					Comments: []ThreadComment{
						{
							DatabaseID: 102,
							NodeID:     "PRRC_102",
							Author:     coderabbitBotLogin,
						},
					},
				},
				{
					ID:         "PRRT_resolved",
					IsResolved: true,
					Comments: []ThreadComment{
						{
							DatabaseID: 104,
							NodeID:     "PRRC_104",
							Author:     coderabbitBotLogin,
						},
					},
				},
			},
		},
	}

	items, err := client.FetchReviews(context.Background(), reviewsource.FetchRequest{
		BaseRepository:  "owner/project",
		PRNumber:        "123",
		IncludeNitpicks: false,
	})
	if err != nil {
		t.Fatalf("expected CodeRabbit items, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected one unresolved non-nitpick CodeRabbit item, got %#v", items)
	}
	item := items[0]
	if item.SourceRef != "thread:PRRT_unresolved,comment:PRRC_101" {
		t.Fatalf("expected thread and comment source ref, got %q", item.SourceRef)
	}
	if item.SourceReviewID != "9001" {
		t.Fatalf("expected source review id, got %q", item.SourceReviewID)
	}
	if !item.SourceReviewSubmittedAt.Equal(submittedAt) {
		t.Fatalf("expected submitted_at %s, got %s", submittedAt, item.SourceReviewSubmittedAt)
	}
	if item.ReviewHash == "" {
		t.Fatal("expected review hash")
	}
	if item.Severity != "major" {
		t.Fatalf("expected major severity, got %q", item.Severity)
	}
	if !strings.Contains(item.Body, "Do not run") {
		t.Fatalf("expected reviewer body to be preserved, got %q", item.Body)
	}
}

func TestFetchReviewsCanIncludeNitpicks(t *testing.T) {
	client := Client{
		GitHub: &fakeGitHubClient{
			comments: []ReviewComment{
				{
					DatabaseID: 201,
					NodeID:     "PRRC_201",
					Body:       "nitpick: rename this",
					Path:       "README.md",
					Line:       3,
					Author:     coderabbitBotLogin,
				},
			},
			threads: []ReviewThread{
				{
					ID:         "PRRT_nitpick",
					IsResolved: false,
					Comments: []ThreadComment{
						{DatabaseID: 201, NodeID: "PRRC_201", Author: coderabbitBotLogin},
					},
				},
			},
		},
	}

	items, err := client.FetchReviews(context.Background(), reviewsource.FetchRequest{
		BaseRepository:  "owner/project",
		PRNumber:        "123",
		IncludeNitpicks: true,
	})
	if err != nil {
		t.Fatalf("expected CodeRabbit items, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected nitpick item, got %#v", items)
	}
	if items[0].Severity != "nitpick" {
		t.Fatalf("expected nitpick severity, got %q", items[0].Severity)
	}
}

func TestResolveIssuesResolvesUniqueReviewThreads(t *testing.T) {
	gh := &fakeGitHubClient{}
	client := Client{GitHub: gh}

	err := client.ResolveIssues(context.Background(), reviewsource.ResolveRequest{
		BaseRepository: "owner/project",
		PRNumber:       "123",
		Issues: []reviewsource.ResolvedIssue{
			{SourceRef: "thread:PRRT_one,comment:PRRC_1"},
			{SourceRef: "thread:PRRT_one,comment:PRRC_2"},
			{SourceRef: "review_hash:only-local"},
			{SourceRef: "thread:PRRT_two,comment:PRRC_3"},
		},
	})
	if err != nil {
		t.Fatalf("resolve issues: %v", err)
	}

	if strings.Join(gh.resolvedThreads, ",") != "PRRT_one,PRRT_two" {
		t.Fatalf("expected unique thread resolution, got %#v", gh.resolvedThreads)
	}
}

func TestWatchStatusReportsReviewingFromPendingCodeRabbitCheck(t *testing.T) {
	client := Client{
		GitHub: &fakeGitHubClient{
			checkRuns: []CheckRun{
				{Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "in_progress"},
			},
		},
	}

	status, err := client.WatchStatus(context.Background(), watchStatusRequest())
	if err != nil {
		t.Fatalf("watch status: %v", err)
	}
	if status.State != watch.StatusReviewing {
		t.Fatalf("expected reviewing status, got %#v", status)
	}
	if !strings.Contains(status.Detail, "in_progress") {
		t.Fatalf("expected CodeRabbit check detail, got %q", status.Detail)
	}
}

func TestWatchStatusReportsSettledFromCompletedCodeRabbitCheck(t *testing.T) {
	client := Client{
		GitHub: &fakeGitHubClient{
			checkRuns: []CheckRun{
				{Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "completed", Conclusion: "success"},
			},
		},
	}

	status, err := client.WatchStatus(context.Background(), watchStatusRequest())
	if err != nil {
		t.Fatalf("watch status: %v", err)
	}
	if status.State != watch.StatusSettled {
		t.Fatalf("expected settled status, got %#v", status)
	}
	if !strings.Contains(status.Detail, "success") {
		t.Fatalf("expected CodeRabbit conclusion detail, got %q", status.Detail)
	}
}

func TestWatchStatusReportsSettledFromCodeRabbitCommitStatus(t *testing.T) {
	client := Client{
		GitHub: &fakeGitHubClient{
			statuses: []CommitStatus{
				{Context: "coderabbitai", State: "success"},
			},
		},
	}

	status, err := client.WatchStatus(context.Background(), watchStatusRequest())
	if err != nil {
		t.Fatalf("watch status: %v", err)
	}
	if status.State != watch.StatusSettled {
		t.Fatalf("expected settled status, got %#v", status)
	}
	if !strings.Contains(status.Detail, "success") {
		t.Fatalf("expected CodeRabbit status detail, got %q", status.Detail)
	}
}

func TestWatchStatusComparesCodeRabbitReviewCommitToCurrentHead(t *testing.T) {
	t.Run("current head review is settled", func(t *testing.T) {
		client := Client{
			GitHub: &fakeGitHubClient{
				reviews: []PullRequestReview{
					{DatabaseID: 9001, Author: coderabbitBotLogin, CommitSHA: "abc123"},
				},
			},
		}

		status, err := client.WatchStatus(context.Background(), watchStatusRequest())
		if err != nil {
			t.Fatalf("watch status: %v", err)
		}
		if status.State != watch.StatusSettled {
			t.Fatalf("expected settled status, got %#v", status)
		}
	})

	t.Run("old head review is pending", func(t *testing.T) {
		client := Client{
			GitHub: &fakeGitHubClient{
				reviews: []PullRequestReview{
					{DatabaseID: 9001, Author: coderabbitBotLogin, CommitSHA: "oldsha"},
				},
			},
		}

		status, err := client.WatchStatus(context.Background(), watchStatusRequest())
		if err != nil {
			t.Fatalf("watch status: %v", err)
		}
		if status.State != watch.StatusPending {
			t.Fatalf("expected pending status, got %#v", status)
		}
		if !strings.Contains(status.Detail, "different commit") {
			t.Fatalf("expected old-commit detail, got %q", status.Detail)
		}
	})
}

func TestWatchStatusReportsPendingWithoutCodeRabbitSignal(t *testing.T) {
	client := Client{
		GitHub: &fakeGitHubClient{
			checkRuns: []CheckRun{
				{Name: "ci", AppName: "GitHub Actions", HeadSHA: "abc123", Status: "completed", Conclusion: "success"},
			},
			statuses: []CommitStatus{
				{Context: "ci", State: "success"},
			},
			reviews: []PullRequestReview{
				{DatabaseID: 9001, Author: "maintainer", CommitSHA: "abc123"},
			},
		},
	}

	status, err := client.WatchStatus(context.Background(), watchStatusRequest())
	if err != nil {
		t.Fatalf("watch status: %v", err)
	}
	if status.State != watch.StatusPending {
		t.Fatalf("expected pending status, got %#v", status)
	}
}

func watchStatusRequest() reviewsource.WatchStatusRequest {
	return reviewsource.WatchStatusRequest{
		BaseRepository: "owner/project",
		PRNumber:       "123",
		HeadSHA:        "abc123",
	}
}

type fakeGitHubClient struct {
	comments        []ReviewComment
	threads         []ReviewThread
	checkRuns       []CheckRun
	statuses        []CommitStatus
	reviews         []PullRequestReview
	resolvedThreads []string
}

func (client fakeGitHubClient) ReviewComments(context.Context, string, string) ([]ReviewComment, error) {
	return client.comments, nil
}

func (client fakeGitHubClient) ReviewThreads(context.Context, string, string) ([]ReviewThread, error) {
	return client.threads, nil
}

func (client fakeGitHubClient) CheckRuns(context.Context, string, string) ([]CheckRun, error) {
	return client.checkRuns, nil
}

func (client fakeGitHubClient) CommitStatuses(context.Context, string, string) ([]CommitStatus, error) {
	return client.statuses, nil
}

func (client fakeGitHubClient) PullRequestReviews(context.Context, string, string) ([]PullRequestReview, error) {
	return client.reviews, nil
}

func (client *fakeGitHubClient) ResolveReviewThread(_ context.Context, threadID string) error {
	client.resolvedThreads = append(client.resolvedThreads, threadID)
	return nil
}
