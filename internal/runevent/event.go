// Package runevent defines the Run Event seam: the product event type, the
// sink interface every consumer implements, and the fanout failure policy.
// It is a leaf module and must depend only on the standard library.
package runevent

import (
	"context"
	"encoding/json"
	"time"
)

// Source identifies which subsystem produced a Run Event.
type Source string

const (
	SourceAgent        Source = "agent"
	SourceDaemon       Source = "daemon"
	SourceVerification Source = "verification"
	SourceGit          Source = "git"
	SourceReviewSource Source = "review_source"
)

// Kind is a namespaced event kind. Readers must treat unknown kinds as
// skippable, never as errors.
type Kind string

const (
	KindAgentMessage          Kind = "agent.message"
	KindAgentThought          Kind = "agent.thought"
	KindAgentToolStarted      Kind = "agent.tool_started"
	KindAgentToolUpdated      Kind = "agent.tool_updated"
	KindAgentPlan             Kind = "agent.plan"
	KindAgentStatus           Kind = "agent.status"
	KindAgentRaw              Kind = "agent.raw"
	KindAgentSelectionReceipt Kind = "agent.selection_receipt"
	KindReviewSourceRequest   Kind = "review_source.request"

	// Daemon kinds cover the orchestration loop: every user-meaningful
	// state transition appends one of these with a small payload of IDs
	// and counts. Large output stays in agent-source events.
	KindDaemonStatus                  Kind = "daemon.status"
	KindDaemonReviewStatus            Kind = "daemon.review_status"
	KindDaemonQuietPeriod             Kind = "daemon.quiet_period"
	KindDaemonFetch                   Kind = "daemon.fetch"
	KindDaemonSelection               Kind = "daemon.selection"
	KindDaemonBatch                   Kind = "daemon.batch"
	KindDaemonVerification            Kind = "daemon.verification"
	KindDaemonCommit                  Kind = "daemon.commit"
	KindDaemonPush                    Kind = "daemon.push"
	KindDaemonSourceResolution        Kind = "daemon.source_resolution"
	KindDaemonRetry                   Kind = "daemon.retry"
	KindDaemonOutcome                 Kind = "daemon.outcome"
	KindDaemonAgentSelectionAttempt   Kind = "daemon.agent_selection_attempt"
	KindDaemonAgentSelectionActive    Kind = "daemon.agent_selection_active"
	KindDaemonAgentSelectionFallback  Kind = "daemon.agent_selection_fallback"
	KindDaemonAgentSelectionExhausted Kind = "daemon.agent_selection_exhausted"
	KindDaemonAgentSelectionClosed    Kind = "daemon.agent_selection_closed"

	// Spec Run kinds: daemon.task records one Task's phase transitions
	// (started, skipped, settled) and daemon.qa records the QA step. The
	// Task id rides in the existing ReviewIssue Work Item field.
	KindDaemonTask Kind = "daemon.task"
	KindDaemonQA   Kind = "daemon.qa"
)

type ReviewRequestOutcome string

const (
	ReviewRequestPublished    ReviewRequestOutcome = "published"
	ReviewRequestDeduplicated ReviewRequestOutcome = "deduplicated"
)

// VerificationPhase names the daemon.verification payload phase. The
// aggregate verdict is carried by the verdict phase only.
type VerificationPhase string

const (
	VerificationPhaseWaiting       VerificationPhase = "waiting"
	VerificationPhaseStarted       VerificationPhase = "started"
	VerificationPhaseCommandPassed VerificationPhase = "command-passed"
	VerificationPhaseFailed        VerificationPhase = "failed"
	VerificationPhaseVerdict       VerificationPhase = "verdict"
)

// VerificationVerdict names the aggregate outcome of one Verification
// attempt.
type VerificationVerdict string

const (
	VerificationVerdictPassed VerificationVerdict = "passed"
	VerificationVerdictFailed VerificationVerdict = "failed"
)

// RepeatedFailure names the retained Verification failure whose diagnostic
// signature matches the current failure for the same Work Item.
type RepeatedFailure struct {
	Signature string `json:"signature"`
	RunID     string `json:"run_id"`
	Attempt   int    `json:"attempt"`
}

// VerificationClassification names a bounded Verification failure class.
type VerificationClassification string

const (
	VerificationClassificationTemporary    VerificationClassification = "temporary"
	VerificationClassificationPrecondition VerificationClassification = "precondition"
	VerificationClassificationVacuous      VerificationClassification = "verification_vacuous"
	VerificationClassificationUnknown      VerificationClassification = "verification_unknown"
)

// VerificationReason names a bounded machine-readable Verification reason.
type VerificationReason string

const (
	VerificationReasonTemporaryFailure          VerificationReason = "temporary_verification_failure"
	VerificationReasonRepositoryNotGreenOnEntry VerificationReason = "repository_not_green_on_entry"
)

// IsDaemonKind reports whether the kind belongs to the known daemon
// vocabulary. Readers render these from the bounded summary; unknown kinds
// stay skippable.
func IsDaemonKind(kind Kind) bool {
	switch kind {
	case KindDaemonStatus, KindDaemonReviewStatus, KindDaemonQuietPeriod,
		KindDaemonFetch, KindDaemonSelection, KindDaemonBatch,
		KindDaemonVerification, KindDaemonCommit, KindDaemonPush,
		KindDaemonSourceResolution, KindDaemonRetry, KindDaemonOutcome,
		KindDaemonAgentSelectionAttempt, KindDaemonAgentSelectionActive,
		KindDaemonAgentSelectionFallback, KindDaemonAgentSelectionExhausted,
		KindDaemonAgentSelectionClosed,
		KindDaemonTask, KindDaemonQA:
		return true
	default:
		return false
	}
}

// RunEvent is one ordered product record of something meaningful that
// happened during a Run. For agent-source events the payload is the raw ACP
// session/update JSON exactly as the ACP Runtime sent it (ADR 0008);
// producers must never re-serialize or prune it.
type RunEvent struct {
	RunID       string
	Batch       int // 0 means outside any Batch
	Source      Source
	Kind        Kind
	ReviewIssue string
	ToolID      string
	ToolState   string
	Summary     string
	Time        time.Time
	Payload     json.RawMessage
}

// ReviewStatusPayload is the stable daemon.review_status payload. It carries
// the bounded Review Source Evidence observation without provider response
// bodies.
type ReviewStatusPayload struct {
	Phase           string    `json:"phase,omitempty"`
	StartedAt       time.Time `json:"started_at,omitzero"`
	Deadline        time.Time `json:"deadline,omitzero"`
	EvidenceState   string    `json:"evidence_state,omitempty"`
	EvidenceKind    string    `json:"evidence_kind,omitempty"`
	RetryStatus     string    `json:"retry_status,omitempty"`
	State           string    `json:"state"`
	Kind            string    `json:"kind"`
	Identity        string    `json:"identity"`
	ExpectedHeadSHA string    `json:"expected_head_sha"`
	ObservedHeadSHA string    `json:"observed_head_sha,omitempty"`
	ParentHeadSHA   string    `json:"parent_head_sha,omitempty"`
	Conclusion      string    `json:"conclusion,omitempty"`
	Detail          string    `json:"detail,omitempty"`
	Reason          string    `json:"reason,omitempty"`
}

// ReviewRequestPayload is the stable review_source.request payload.
type ReviewRequestPayload struct {
	HeadSHA string               `json:"head_sha"`
	Command string               `json:"command"`
	Outcome ReviewRequestOutcome `json:"outcome"`
}

const (
	SelectionReceiptEventApplied  = "agent_selection_receipt"
	SelectionReceiptStatusApplied = "applied"
)

// SelectionReceiptPayload records the effective reasoning effort observed
// after an Agent Session applies a runtime-deferred selection.
type SelectionReceiptPayload struct {
	Event                    string `json:"event"`
	Session                  string `json:"session"`
	Runtime                  string `json:"runtime"`
	Model                    string `json:"model"`
	RequestedReasoningEffort string `json:"requested_reasoning_effort"`
	ReasoningEffort          string `json:"reasoning_effort"`
	Status                   string `json:"status"`
}

// RetryPayload is the stable daemon.retry payload for one bounded Review
// Source retry episode.
type RetryPayload struct {
	Phase     string `json:"phase"`
	Operation string `json:"operation"`
	Reason    string `json:"reason"`
}

// NotificationReceiptPayload is the stable daemon.status payload for one
// completed best-effort terminal notification attempt.
type NotificationReceiptPayload struct {
	Event       string    `json:"event"`
	Route       string    `json:"route"`
	Status      string    `json:"status"`
	CompletedAt time.Time `json:"completed_at"`
	Reason      string    `json:"reason,omitempty"`
}

// OutcomePayload is the stable daemon.outcome payload. Optional evidence and
// recovery fields let non-Clean terminal outcomes remain actionable without
// changing the event kind or stream schema.
type OutcomePayload struct {
	State             string `json:"state"`
	Remaining         int    `json:"remaining"`
	Reason            string `json:"reason,omitempty"`
	NextAction        string `json:"next_action,omitempty"`
	ReviewIssuesKnown *bool  `json:"review_issues_known,omitempty"`
	ConsoleLog        string `json:"console_log,omitempty"`
	AttachCommand     string `json:"attach_command,omitempty"`
	EvidenceState     string `json:"evidence_state,omitempty"`
	EvidenceKind      string `json:"evidence_kind,omitempty"`
	EvidenceHeadSHA   string `json:"evidence_head_sha,omitempty"`
	VerifiedHeadSHA   string `json:"verified_head_sha,omitempty"`
}

// Sink consumes published Run Events. Context comes first because durable
// adapters perform IO.
type Sink interface {
	Publish(ctx context.Context, event RunEvent) error
}

// Discard is a Sink that drops every event.
var Discard Sink = discardSink{}

type discardSink struct{}

func (discardSink) Publish(context.Context, RunEvent) error { return nil }
