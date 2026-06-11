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
	KindAgentMessage     Kind = "agent.message"
	KindAgentThought     Kind = "agent.thought"
	KindAgentToolStarted Kind = "agent.tool_started"
	KindAgentToolUpdated Kind = "agent.tool_updated"
	KindAgentPlan        Kind = "agent.plan"
	KindAgentStatus      Kind = "agent.status"
	KindAgentRaw         Kind = "agent.raw"

	// Daemon kinds cover the orchestration loop: every user-meaningful
	// state transition appends one of these with a small payload of IDs
	// and counts. Large output stays in agent-source events.
	KindDaemonStatus           Kind = "daemon.status"
	KindDaemonReviewStatus     Kind = "daemon.review_status"
	KindDaemonQuietPeriod      Kind = "daemon.quiet_period"
	KindDaemonFetch            Kind = "daemon.fetch"
	KindDaemonSelection        Kind = "daemon.selection"
	KindDaemonBatch            Kind = "daemon.batch"
	KindDaemonVerification     Kind = "daemon.verification"
	KindDaemonCommit           Kind = "daemon.commit"
	KindDaemonPush             Kind = "daemon.push"
	KindDaemonSourceResolution Kind = "daemon.source_resolution"
	KindDaemonRetry            Kind = "daemon.retry"
	KindDaemonOutcome          Kind = "daemon.outcome"
)

// IsDaemonKind reports whether the kind belongs to the known daemon
// vocabulary. Readers render these from the bounded summary; unknown kinds
// stay skippable.
func IsDaemonKind(kind Kind) bool {
	switch kind {
	case KindDaemonStatus, KindDaemonReviewStatus, KindDaemonQuietPeriod,
		KindDaemonFetch, KindDaemonSelection, KindDaemonBatch,
		KindDaemonVerification, KindDaemonCommit, KindDaemonPush,
		KindDaemonSourceResolution, KindDaemonRetry, KindDaemonOutcome:
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

// Sink consumes published Run Events. Context comes first because durable
// adapters perform IO.
type Sink interface {
	Publish(ctx context.Context, event RunEvent) error
}

// Discard is a Sink that drops every event.
var Discard Sink = discardSink{}

type discardSink struct{}

func (discardSink) Publish(context.Context, RunEvent) error { return nil }
