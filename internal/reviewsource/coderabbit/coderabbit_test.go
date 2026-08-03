package coderabbit

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/watch"
)

func TestFetchReviewsFiltersToUnresolvedCodeRabbitThreads(t *testing.T) {
	t.Parallel()
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

func TestFetchReviewsAcceptsCodeRabbitAppLoginFromGraphQLThreads(t *testing.T) {
	t.Parallel()
	client := Client{
		GitHub: &fakeGitHubClient{
			comments: []ReviewComment{
				{
					DatabaseID: 301,
					NodeID:     "PRRC_301",
					Body:       "major: update the smoke URL",
					Path:       "apps/api/tests/smoke/production-smoke.test.ts",
					Line:       7,
					Author:     coderabbitBotLogin,
				},
			},
			threads: []ReviewThread{
				{
					ID:         "PRRT_graphql_login",
					IsResolved: false,
					Comments: []ThreadComment{
						{
							DatabaseID: 301,
							NodeID:     "PRRC_301",
							Author:     coderabbitAppLogin,
						},
					},
				},
			},
		},
	}

	items, err := client.FetchReviews(context.Background(), reviewsource.FetchRequest{
		BaseRepository: "owner/project",
		PRNumber:       "123",
	})
	if err != nil {
		t.Fatalf("expected CodeRabbit items, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected CodeRabbit item with app login thread author, got %#v", items)
	}
	if items[0].SourceRef != "thread:PRRT_graphql_login,comment:PRRC_301" {
		t.Fatalf("expected stable thread/comment source ref, got %q", items[0].SourceRef)
	}
}

func TestFetchReviewsCanIncludeNitpicks(t *testing.T) {
	t.Parallel()
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

func TestTitleFromBodyStripsCodeRabbitMarkup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "table fragment heading with emoji",
			body: "| ## 🧹 Clean up API docs |\n\nThe generated docs need trimming.",
			want: "Clean up API docs",
		},
		{
			name: "heading with shortcode emoji and markdown emphasis",
			body: "### :warning: **Potential issue** | avoid panic in setup\n\nDetails.",
			want: "Potential issue avoid panic in setup",
		},
		{
			name: "surrounding whitespace and emoji marker",
			body: "  ## 🛠️ major: handle nil cache  \n\nDetails.",
			want: "major: handle nil cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := titleFromBody(tt.body)
			if got != tt.want {
				t.Fatalf("expected title %q, got %q", tt.want, got)
			}
			for _, disallowed := range []string{"|", "#", "🧹", "🛠", ":warning:", "**"} {
				if strings.Contains(got, disallowed) {
					t.Fatalf("expected title %q to omit %q", got, disallowed)
				}
			}
		})
	}
}

func TestResolveIssuesResolvesUniqueReviewThreads(t *testing.T) {
	t.Parallel()
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

func TestResolveIssueResolvesOneReviewThread(t *testing.T) {
	t.Parallel()
	gh := &fakeGitHubClient{}
	client := Client{GitHub: gh}

	err := client.ResolveIssue(context.Background(), reviewsource.IssueResolveRequest{
		BaseRepository: "owner/project",
		PRNumber:       "123",
		SourceRef:      "thread:PRRT_thread,comment:PRRC_1",
	})
	if err != nil {
		t.Fatalf("resolve issue: %v", err)
	}

	if strings.Join(gh.resolvedThreads, ",") != "PRRT_thread" {
		t.Fatalf("expected single thread resolution, got %#v", gh.resolvedThreads)
	}
}

func TestReplyToIssueUsesMarkerForIdempotency(t *testing.T) {
	t.Parallel()
	const marker = "<!-- roundfix:outcome run=run_1 issue=abc action=invalid -->"
	tests := []struct {
		name        string
		threads     []ReviewThread
		wantPosted  bool
		wantSkipped bool
		wantReplies int
	}{
		{
			name: "posts when marker absent",
			threads: []ReviewThread{{
				ID: "PRRT_thread",
				Comments: []ThreadComment{
					{Body: "previous human reply"},
				},
			}},
			wantPosted:  true,
			wantReplies: 1,
		},
		{
			name: "skips when marker already present",
			threads: []ReviewThread{{
				ID: "PRRT_thread",
				Comments: []ThreadComment{
					{Body: "already handled\n" + marker},
				},
			}},
			wantSkipped: true,
		},
		{
			name: "posts when marker is embedded in quoted text",
			threads: []ReviewThread{{
				ID: "PRRT_thread",
				Comments: []ThreadComment{
					{Body: "quoted previous reply: " + marker},
				},
			}},
			wantPosted:  true,
			wantReplies: 1,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gh := &fakeGitHubClient{threads: testCase.threads}
			client := Client{GitHub: gh}

			result, err := client.ReplyToIssue(context.Background(), reviewsource.IssueCommentRequest{
				BaseRepository: "owner/project",
				PRNumber:       "123",
				SourceRef:      "thread:PRRT_thread,comment:PRRC_1",
				Marker:         marker,
				Body:           marker + "\n\nRoundfix outcome: invalid.",
			})
			if err != nil {
				t.Fatalf("reply to issue: %v", err)
			}

			if result.Posted != testCase.wantPosted || result.Skipped != testCase.wantSkipped {
				t.Fatalf("unexpected result: %+v", result)
			}
			if len(gh.replies) != testCase.wantReplies {
				t.Fatalf("expected %d replies, got %+v", testCase.wantReplies, gh.replies)
			}
		})
	}
}

func TestGHClientWriteMutationsInvokeGHOnce(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		act      func(context.Context, GHClient) error
		wantArgs []string
	}{
		{
			name: "reply to review thread",
			act: func(ctx context.Context, client GHClient) error {
				return client.ReplyToReviewThread(ctx, "PRRT_thread", "failed verification details")
			},
			wantArgs: []string{
				"api",
				"graphql",
				"-f",
				"query=" + replyToReviewThreadMutation,
				"-f",
				"threadId=PRRT_thread",
				"-f",
				"body=failed verification details",
			},
		},
		{
			name: "comment on pull request",
			act: func(ctx context.Context, client GHClient) error {
				return client.CommentOnPullRequest(ctx, "octo/base", 123, "bypass audit details")
			},
			wantArgs: []string{
				"api",
				"-X",
				"POST",
				"repos/octo/base/issues/123/comments",
				"-f",
				"body=bypass audit details",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var calls [][]string
			withRunGH(t, func(_ context.Context, args ...string) ([]byte, error) {
				calls = append(calls, append([]string(nil), args...))
				return []byte(`{}`), nil
			})

			if err := testCase.act(context.Background(), GHClient{}); err != nil {
				t.Fatalf("write mutation failed: %v", err)
			}

			if len(calls) != 1 {
				t.Fatalf("expected one gh invocation, got %#v", calls)
			}
			assertStringSlicesEqual(t, calls[0], testCase.wantArgs)
			if testCase.name == "reply to review thread" && !strings.Contains(calls[0][3], "addPullRequestReviewThreadReply") {
				t.Fatalf("expected reply mutation name in args, got %#v", calls[0])
			}
		})
	}
}

func TestGHClientWriteMutationFailuresWrapCause(t *testing.T) {
	t.Parallel()
	cause := errors.New("gh failed")
	tests := []struct {
		name          string
		act           func(context.Context, GHClient) error
		wantOperation string
	}{
		{
			name: "reply failure",
			act: func(ctx context.Context, client GHClient) error {
				return client.ReplyToReviewThread(ctx, "PRRT_thread", "body")
			},
			wantOperation: "reply to review thread",
		},
		{
			name: "comment failure",
			act: func(ctx context.Context, client GHClient) error {
				return client.CommentOnPullRequest(ctx, "octo/base", 123, "body")
			},
			wantOperation: "comment on pull request",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			withRunGH(t, func(context.Context, ...string) ([]byte, error) {
				return nil, cause
			})

			err := testCase.act(context.Background(), GHClient{})

			if err == nil {
				t.Fatal("expected write mutation failure")
			}
			if !errors.Is(err, cause) {
				t.Fatalf("expected error to wrap cause, got %v", err)
			}
			if !strings.Contains(err.Error(), testCase.wantOperation) {
				t.Fatalf("expected error to name operation %q, got %v", testCase.wantOperation, err)
			}
		})
	}
}

func TestRoundfixCommentMarkerHelpers(t *testing.T) {
	t.Parallel()
	marker := RoundfixCommentMarker("run:run_123", "issue:abc")
	if marker != "<!-- roundfix: run:run_123 issue:abc -->" {
		t.Fatalf("unexpected marker %q", marker)
	}
	body := RoundfixCommentBody("failed verification details", marker)
	if body != "failed verification details\n\n<!-- roundfix: run:run_123 issue:abc -->" {
		t.Fatalf("unexpected marker body %q", body)
	}
	if !HasRoundfixCommentMarker(marker, "unrelated\n  "+marker+"  \nmore") {
		t.Fatal("expected marker helper to detect marker line")
	}
	if HasRoundfixCommentMarker(marker, "unrelated comment", "prefix "+marker+" suffix") {
		t.Fatal("expected marker helper to require a standalone marker line")
	}
}

func TestWatchStatusReportsReviewingFromPendingCodeRabbitCheck(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestEvidenceHierarchyPrecedence(t *testing.T) {
	t.Parallel()
	unresolvedThread := ReviewThread{
		ID:         "thread-1",
		IsResolved: false,
		Comments: []ThreadComment{
			{Author: coderabbitBotLogin},
		},
	}
	tests := []struct {
		name       string
		checkRuns  []CheckRun
		statuses   []CommitStatus
		reviews    []PullRequestReview
		threads    []ReviewThread
		wantState  reviewsource.EvidenceState
		wantKind   reviewsource.EvidenceKind
		wantReason string
	}{
		{
			name: "explicit skip wins over pending and success",
			checkRuns: []CheckRun{
				{
					DatabaseID:    41,
					Name:          "CodeRabbit",
					AppName:       "CodeRabbit",
					HeadSHA:       "abc123",
					Status:        "completed",
					Conclusion:    "success",
					OutputTitle:   "Review skipped",
					OutputSummary: "Review skipped because this pull request contains too many files.",
				},
				{DatabaseID: 42, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "in_progress"},
			},
			statuses:   []CommitStatus{{Context: "coderabbitai", State: "success"}},
			wantState:  reviewsource.EvidenceSkipped,
			wantKind:   reviewsource.EvidenceKindCheckRun,
			wantReason: "Review skipped because this pull request contains too many files.",
		},
		{
			name: "pending check wins over successful status",
			checkRuns: []CheckRun{
				{DatabaseID: 43, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "in_progress"},
			},
			statuses:  []CommitStatus{{Context: "coderabbitai", State: "success"}},
			wantState: reviewsource.EvidenceReviewing,
			wantKind:  reviewsource.EvidenceKindCheckRun,
		},
		{
			name: "successful check verifies with no unresolved threads",
			checkRuns: []CheckRun{
				{DatabaseID: 44, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "completed", Conclusion: "success"},
			},
			wantState: reviewsource.EvidenceVerified,
			wantKind:  reviewsource.EvidenceKindCheckRun,
		},
		{
			name: "successful check is reviewed with unresolved thread",
			checkRuns: []CheckRun{
				{DatabaseID: 45, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "completed", Conclusion: "success"},
			},
			threads:   []ReviewThread{unresolvedThread},
			wantState: reviewsource.EvidenceReviewed,
			wantKind:  reviewsource.EvidenceKindCheckRun,
		},
		{
			name: "current approval verifies with no unresolved threads",
			reviews: []PullRequestReview{
				{DatabaseID: 9001, Author: coderabbitBotLogin, State: "APPROVED", CommitSHA: "abc123"},
			},
			wantState: reviewsource.EvidenceVerified,
			wantKind:  reviewsource.EvidenceKindReviewApproval,
		},
		{
			name: "current approval is reviewed with unresolved thread",
			reviews: []PullRequestReview{
				{DatabaseID: 9002, Author: coderabbitBotLogin, State: "APPROVED", CommitSHA: "abc123"},
			},
			threads:   []ReviewThread{unresolvedThread},
			wantState: reviewsource.EvidenceReviewed,
			wantKind:  reviewsource.EvidenceKindReviewApproval,
		},
		{
			name: "commented current review is never approval",
			reviews: []PullRequestReview{
				{DatabaseID: 9003, Author: coderabbitBotLogin, State: "COMMENTED", CommitSHA: "abc123"},
			},
			wantState: reviewsource.EvidenceReviewed,
			wantKind:  reviewsource.EvidenceKindNone,
		},
		{
			name:      "no usable signal stays pending",
			wantState: reviewsource.EvidencePending,
			wantKind:  reviewsource.EvidenceKindNone,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			client := Client{GitHub: &fakeGitHubClient{
				checkRuns: testCase.checkRuns,
				statuses:  testCase.statuses,
				reviews:   testCase.reviews,
				threads:   testCase.threads,
			}}

			evidence, err := client.Evidence(context.Background(), evidenceRequest())
			if err != nil {
				t.Fatalf("classify evidence: %v", err)
			}
			if evidence.State != testCase.wantState || evidence.Kind != testCase.wantKind {
				t.Fatalf("evidence = %#v, want state %q kind %q", evidence, testCase.wantState, testCase.wantKind)
			}
			if evidence.ExpectedHeadSHA != "abc123" {
				t.Fatalf("expected head = %q, want abc123", evidence.ExpectedHeadSHA)
			}
			if evidence.Reason != testCase.wantReason {
				t.Fatalf("reason = %q, want %q", evidence.Reason, testCase.wantReason)
			}
		})
	}
}

func TestEvidenceReviewingSkipsReviewAndThreadRequests(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		checkRuns []CheckRun
		statuses  []CommitStatus
		wantKind  reviewsource.EvidenceKind
	}{
		{
			name: "pending check run",
			checkRuns: []CheckRun{
				{DatabaseID: 43, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "abc123", Status: "in_progress"},
			},
			wantKind: reviewsource.EvidenceKindCheckRun,
		},
		{
			name: "pending commit status",
			statuses: []CommitStatus{
				{Context: "coderabbitai", State: "pending"},
			},
			wantKind: reviewsource.EvidenceKindCommitStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewCalls := 0
			threadCalls := 0
			client := Client{GitHub: &fakeGitHubClient{
				checkRuns:   tt.checkRuns,
				statuses:    tt.statuses,
				reviewCalls: &reviewCalls,
				threadCalls: &threadCalls,
			}}

			evidence, err := client.Evidence(context.Background(), evidenceRequest())
			if err != nil {
				t.Fatalf("classify reviewing Evidence: %v", err)
			}
			if evidence.State != reviewsource.EvidenceReviewing || evidence.Kind != tt.wantKind {
				t.Fatalf("reviewing Evidence = %#v, want kind %q", evidence, tt.wantKind)
			}
			if reviewCalls != 0 || threadCalls != 0 {
				t.Fatalf("reviewing Evidence fetched reviews=%d threads=%d, want 0 calls", reviewCalls, threadCalls)
			}
		})
	}
}

func TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals(t *testing.T) {
	t.Parallel()
	client := Client{GitHub: &fakeGitHubClient{
		checkRuns: []CheckRun{
			{DatabaseID: 51, Name: "CodeRabbit", AppName: "CodeRabbit", Status: "completed", Conclusion: "success"},
			{DatabaseID: 52, Name: "CodeRabbit", AppName: "CodeRabbit", HeadSHA: "oldsha", Status: "completed", Conclusion: "success"},
		},
		reviews: []PullRequestReview{
			{DatabaseID: 9004, Author: coderabbitBotLogin, State: "APPROVED", CommitSHA: "oldsha"},
		},
	}}

	evidence, err := client.Evidence(context.Background(), evidenceRequest())
	if err != nil {
		t.Fatalf("classify evidence: %v", err)
	}
	if evidence.State != reviewsource.EvidencePending || evidence.Kind != reviewsource.EvidenceKindNone {
		t.Fatalf("stale or unbound signal verified expected head: %#v", evidence)
	}
	if evidence.ExpectedHeadSHA != "abc123" || evidence.ObservedHeadSHA != "oldsha" {
		t.Fatalf("expected stale-head detail, got %#v", evidence)
	}
}

func TestHeadCheckMapsGitHubCheckRunJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		fixture string
		want    watch.HeadCheckState
	}{
		{
			name: "success",
			fixture: `{
				"total_count": 1,
				"check_runs": [{
					"name": "CodeRabbit",
					"head_sha": "abc123",
					"status": "completed",
					"conclusion": "success",
					"app": {"name": "CodeRabbit", "slug": "coderabbitai"}
				}]
			}`,
			want: watch.CheckSuccess,
		},
		{
			name: "failure",
			fixture: `{
				"total_count": 1,
				"check_runs": [{
					"name": "CodeRabbit",
					"head_sha": "abc123",
					"status": "completed",
					"conclusion": "failure",
					"app": {"name": "CodeRabbit", "slug": "coderabbitai"}
				}]
			}`,
			want: watch.CheckFailure,
		},
		{
			name: "in progress",
			fixture: `{
				"total_count": 1,
				"check_runs": [{
					"name": "CodeRabbit",
					"head_sha": "abc123",
					"status": "in_progress",
					"app": {"name": "CodeRabbit", "slug": "coderabbitai"}
				}]
			}`,
			want: watch.CheckPending,
		},
		{
			name: "absent",
			fixture: `{
				"total_count": 1,
				"check_runs": [{
					"name": "ci",
					"head_sha": "abc123",
					"status": "completed",
					"conclusion": "success",
					"app": {"name": "GitHub Actions", "slug": "github-actions"}
				}]
			}`,
			want: watch.CheckMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkRuns, err := parseCheckRuns([]byte(tt.fixture))
			if err != nil {
				t.Fatalf("parse check runs fixture: %v", err)
			}
			client := Client{GitHub: &fakeGitHubClient{checkRuns: checkRuns}}

			got, err := client.HeadCheck(context.Background(), reviewsource.HeadCheckRequest{
				PRNumber:       "123",
				BaseRepository: "owner/project",
				HeadSHA:        "abc123",
			})

			if err != nil {
				t.Fatalf("head check: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCheckRunOutputJSONMapping(t *testing.T) {
	t.Parallel()
	const sensitiveText = "authorization: Bearer secret-that-must-not-escape"
	fixture := `{
		"total_count": 1,
		"check_runs": [{
			"id": 42,
			"name": "CodeRabbit",
			"head_sha": "abc123",
			"status": "completed",
			"conclusion": "success",
			"app": {"name": "CodeRabbit", "slug": "coderabbitai"},
			"output": {
				"title": "Review skipped",
				"summary": "CodeRabbit skipped this review because the change set is too large.",
				"text": "` + sensitiveText + `"
			}
		}]
	}`

	checkRuns, err := parseCheckRuns([]byte(fixture))
	if err != nil {
		t.Fatalf("parse check runs fixture: %v", err)
	}
	if len(checkRuns) != 1 {
		t.Fatalf("check run count = %d, want 1", len(checkRuns))
	}
	got := checkRuns[0]
	if got.DatabaseID != 42 {
		t.Fatalf("database ID = %d, want 42", got.DatabaseID)
	}
	if got.OutputTitle != "Review skipped" {
		t.Fatalf("output title = %q", got.OutputTitle)
	}
	if got.OutputSummary != "CodeRabbit skipped this review because the change set is too large." {
		t.Fatalf("output summary = %q", got.OutputSummary)
	}
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal check run: %v", err)
	}
	if strings.Contains(string(encoded), sensitiveText) {
		t.Fatalf("check run retained raw output text: %s", encoded)
	}
}

func TestSkipSignalStructuredOutputRemainsAvailable(t *testing.T) {
	t.Parallel()
	fixture := `{
		"check_runs": [{
			"name": "CodeRabbit",
			"head_sha": "abc123",
			"status": "completed",
			"conclusion": "success",
			"output": {
				"title": "Review skipped",
				"summary": "Review skipped because this pull request contains too many files."
			}
		}]
	}`
	checkRuns, err := parseCheckRuns([]byte(fixture))
	if err != nil {
		t.Fatalf("parse check runs fixture: %v", err)
	}
	if checkRuns[0].OutputTitle == "" || checkRuns[0].OutputSummary == "" {
		t.Fatalf("structured skip fields were discarded: %+v", checkRuns[0])
	}
}

func TestSkipSignalDoesNotInferFromArbitrarySuccessfulText(t *testing.T) {
	t.Parallel()
	fixture := `{
		"check_runs": [{
			"name": "CodeRabbit",
			"head_sha": "abc123",
			"status": "completed",
			"conclusion": "success",
			"output": {
				"title": "Review completed",
				"summary": "All checks passed successfully."
			}
		}]
	}`
	checkRuns, err := parseCheckRuns([]byte(fixture))
	if err != nil {
		t.Fatalf("parse check runs fixture: %v", err)
	}
	if checkRuns[0].OutputTitle != "Review completed" || checkRuns[0].OutputSummary != "All checks passed successfully." {
		t.Fatalf("successful output was rewritten as another signal: %+v", checkRuns[0])
	}
}

func TestTransientClassificationMatrix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
	}{
		{name: "temporary DNS", err: &net.DNSError{Err: "temporary failure", Name: "api.github.com", IsTemporary: true}},
		{name: "connection reset", err: syscall.ECONNRESET},
		{name: "HTTP 429", err: errors.New("HTTP 429: rate limited")},
		{name: "GitHub HTTP 500", err: errors.New("HTTP 500: internal server error")},
		{name: "GitHub HTTP 503", err: errors.New("HTTP 503: service unavailable")},
		{name: "non-parent timeout", err: context.DeadlineExceeded},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			client := Client{GitHub: &fakeGitHubClient{checkRunsErr: testCase.err}}
			_, err := client.WatchStatus(context.Background(), watchStatusRequest())
			if err == nil {
				t.Fatal("expected WatchStatus failure")
			}
			if !reviewsource.IsTransient(err) {
				t.Fatalf("error = %T %v, want transient", err, err)
			}
			var transient *reviewsource.TransientError
			if !errors.As(err, &transient) {
				t.Fatalf("error = %T, want *reviewsource.TransientError", err)
			}
			if transient.Operation != "fetch CodeRabbit check runs" {
				t.Fatalf("operation = %q", transient.Operation)
			}
			if !errors.Is(err, testCase.err) {
				t.Fatalf("transient error did not wrap cause %v", testCase.err)
			}
		})
	}
}

func TestTransientPermanentFailureMatrix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
	}{
		{name: "authentication", err: errors.New("HTTP 401: bad credentials")},
		{name: "authorization", err: errors.New("HTTP 403: forbidden")},
		{name: "invalid request", err: errors.New("HTTP 422: validation failed")},
		{name: "malformed response", err: &json.SyntaxError{Offset: 1}},
		{name: "cancellation", err: context.Canceled},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			client := Client{GitHub: &fakeGitHubClient{checkRunsErr: testCase.err}}
			_, err := client.WatchStatus(context.Background(), watchStatusRequest())
			if err == nil {
				t.Fatal("expected WatchStatus failure")
			}
			if reviewsource.IsTransient(err) {
				t.Fatalf("error = %T %v, want permanent", err, err)
			}
		})
	}
}

func TestTransientParentCancellationIsPermanent(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := Client{GitHub: &fakeGitHubClient{
		checkRunsErr: &net.DNSError{Err: "temporary failure", Name: "api.github.com", IsTemporary: true},
	}}

	_, err := client.WatchStatus(ctx, watchStatusRequest())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want parent cancellation", err)
	}
	if reviewsource.IsTransient(err) {
		t.Fatalf("parent cancellation classified transient: %v", err)
	}
}

func watchStatusRequest() reviewsource.WatchStatusRequest {
	return reviewsource.WatchStatusRequest{
		BaseRepository: "owner/project",
		PRNumber:       "123",
		HeadSHA:        "abc123",
	}
}

func evidenceRequest() reviewsource.EvidenceRequest {
	return reviewsource.EvidenceRequest{
		BaseRepository:  "owner/project",
		PRNumber:        "123",
		ExpectedHeadSHA: "abc123",
	}
}

type fakeGitHubClient struct {
	comments          []ReviewComment
	threads           []ReviewThread
	checkRuns         []CheckRun
	checkRunsErr      error
	statuses          []CommitStatus
	reviews           []PullRequestReview
	resolvedThreads   []string
	replies           []reviewThreadReplyCall
	prComments        []pullRequestCommentCall
	commentRepository string
	reviewCalls       *int
	threadCalls       *int
}

func (client fakeGitHubClient) ReviewComments(context.Context, string, string) ([]ReviewComment, error) {
	return client.comments, nil
}

func (client fakeGitHubClient) ReviewThreads(context.Context, string, string) ([]ReviewThread, error) {
	if client.threadCalls != nil {
		(*client.threadCalls)++
	}
	return client.threads, nil
}

func (client fakeGitHubClient) CheckRuns(context.Context, string, string) ([]CheckRun, error) {
	return client.checkRuns, client.checkRunsErr
}

func (client fakeGitHubClient) CommitStatuses(context.Context, string, string) ([]CommitStatus, error) {
	return client.statuses, nil
}

func (client fakeGitHubClient) PullRequestReviews(context.Context, string, string) ([]PullRequestReview, error) {
	if client.reviewCalls != nil {
		(*client.reviewCalls)++
	}
	return client.reviews, nil
}

func (client *fakeGitHubClient) ResolveReviewThread(_ context.Context, threadID string) error {
	client.resolvedThreads = append(client.resolvedThreads, threadID)
	return nil
}

type reviewThreadReplyCall struct {
	ThreadID string
	Body     string
}

type pullRequestCommentCall struct {
	PRNumber int
	Body     string
}

func (client *fakeGitHubClient) ReplyToReviewThread(_ context.Context, threadID string, body string) error {
	client.replies = append(client.replies, reviewThreadReplyCall{ThreadID: threadID, Body: body})
	return nil
}

func (client *fakeGitHubClient) CommentOnPullRequest(_ context.Context, repository string, prNumber int, body string) error {
	client.commentRepository = repository
	client.prComments = append(client.prComments, pullRequestCommentCall{PRNumber: prNumber, Body: body})
	return nil
}

func withRunGH(t *testing.T, runner func(context.Context, ...string) ([]byte, error)) {
	t.Helper()
	previous := runGH
	runGH = runner
	t.Cleanup(func() {
		runGH = previous
	})
}

func assertStringSlicesEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected args %#v, got %#v", want, got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("expected args %#v, got %#v", want, got)
		}
	}
}
