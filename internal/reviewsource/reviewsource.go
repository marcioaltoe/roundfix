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

type WatchStatusRequest struct {
	Source         string
	PRNumber       string
	BaseRepository string
	HeadRepository string
	HeadBranch     string
	HeadSHA        string
}

type WatchStatus = watch.Status

type Source interface {
	FetchReviews(ctx context.Context, req FetchRequest) ([]ReviewItem, error)
}
