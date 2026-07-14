package reviewsource

import (
	"context"
	"time"

	"roundfix/internal/watch"
)

const SourceCodeRabbit = "coderabbit"

type FetchRequest struct {
	Source          string
	PRNumber        string
	BaseRepository  string
	HeadRepository  string
	HeadBranch      string
	HeadSHA         string
	IncludeNitpicks bool
}

type ReviewItem struct {
	Title                   string
	File                    string
	Line                    int
	Severity                string
	Author                  string
	Body                    string
	SourceRef               string
	ReviewHash              string
	SourceReviewID          string
	SourceReviewSubmittedAt time.Time
}

type ResolveRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	Issues         []ResolvedIssue
}

type ResolvedIssue struct {
	FilePath  string
	Status    string
	SourceRef string
}

type IssueResolveRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	SourceRef      string
}

type IssueCommentRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	SourceRef      string
	Marker         string
	Body           string
}

type IssueCommentResult struct {
	Posted  bool
	Skipped bool
}

type WatchStatusRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	HeadRepository string
	HeadBranch     string
	HeadSHA        string
}

type WatchStatus = watch.Status

type HeadCheckRequest struct {
	Source         string
	BaseRepository string
	HeadSHA        string
}

type Source interface {
	FetchReviews(ctx context.Context, req FetchRequest) ([]ReviewItem, error)
}
