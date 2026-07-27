package reviewsource

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"roundfix/internal/watch"
)

const SourceCodeRabbit = "coderabbit"

type EvidenceState string

const (
	EvidencePending   EvidenceState = "pending"
	EvidenceReviewing EvidenceState = "reviewing"
	EvidenceReviewed  EvidenceState = "reviewed"
	EvidenceVerified  EvidenceState = "verified"
	EvidenceSkipped   EvidenceState = "skipped"
	EvidenceFailed    EvidenceState = "failed"
)

type EvidenceKind string

const (
	EvidenceKindNone                   EvidenceKind = "none"
	EvidenceKindCheckRun               EvidenceKind = "check_run"
	EvidenceKindCommitStatus           EvidenceKind = "commit_status"
	EvidenceKindReviewApproval         EvidenceKind = "review_approval"
	EvidenceKindArtifactOnlyDescendant EvidenceKind = "artifact_only_descendant"
)

// MaxEvidenceDetailLength bounds Review Source-authored detail in bytes.
const MaxEvidenceDetailLength = 2048

type Evidence struct {
	State           EvidenceState `json:"state"`
	Kind            EvidenceKind  `json:"kind"`
	Identity        string        `json:"identity"`
	ExpectedHeadSHA string        `json:"expected_head_sha"`
	ObservedHeadSHA string        `json:"observed_head_sha,omitempty"`
	Conclusion      string        `json:"conclusion,omitempty"`
	Detail          string        `json:"detail,omitempty"`
	Reason          string        `json:"reason,omitempty"`
}

// BoundEvidenceDetail truncates Review Source-authored text on a rune boundary.
func BoundEvidenceDetail(detail string) string {
	if len(detail) <= MaxEvidenceDetailLength {
		return detail
	}
	cut := MaxEvidenceDetailLength
	for cut > 0 && !utf8.RuneStart(detail[cut]) {
		cut--
	}
	return detail[:cut] + "…"
}

// TransientError identifies a failed Review Source operation that may be
// retried within the caller's existing timeout and Run Budget.
type TransientError struct {
	Operation string `json:"operation"`
	Err       error  `json:"-"`
}

func (err *TransientError) Error() string {
	operation := BoundEvidenceDetail(err.Operation)
	if operation == "" {
		return "temporary Review Source failure"
	}
	return fmt.Sprintf("%s: temporary Review Source failure", operation)
}

func (err *TransientError) Unwrap() error {
	return err.Err
}

// IsTransient reports whether err wraps a typed transient Review Source
// failure.
func IsTransient(err error) bool {
	var transient *TransientError
	return errors.As(err, &transient)
}

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
