package runevent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StreamSchema names the stable Supervisor Run Event Stream envelope.
const StreamSchema = "roundfix-events/v1"

// StreamCategory is the public event category vocabulary for
// roundfix-events/v1. It intentionally differs from internal Run Event kinds.
type StreamCategory string

const (
	StreamCategoryTaskStatus   StreamCategory = "task-status"
	StreamCategoryBatch        StreamCategory = "batch"
	StreamCategoryVerification StreamCategory = "verification"
	StreamCategoryOutcome      StreamCategory = "outcome"
	StreamCategorySelection    StreamCategory = "agent-selection"
)

// StreamRecord is one JSONL record in the Supervisor Run Event Stream.
type StreamRecord struct {
	Schema              string         `json:"schema"`
	RunID               string         `json:"run_id"`
	Category            StreamCategory `json:"category"`
	Time                string         `json:"time"`
	Cursor              int64          `json:"cursor"`
	Batch               int            `json:"batch,omitempty"`
	Attempt             int            `json:"attempt,omitempty"`
	Retry               int            `json:"retry,omitempty"`
	WorkItem            string         `json:"work_item,omitempty"`
	Phase               string         `json:"phase,omitempty"`
	Mode                string         `json:"mode,omitempty"`
	Classification      string         `json:"classification,omitempty"`
	RetryAvailable      *bool          `json:"retry_available,omitempty"`
	Status              string         `json:"status,omitempty"`
	Verdict             string         `json:"verdict,omitempty"`
	Outcome             string         `json:"outcome,omitempty"`
	Summary             string         `json:"summary,omitempty"`
	ScopeKind           string         `json:"scope_kind,omitempty"`
	ScopeID             string         `json:"scope_id,omitempty"`
	ScopeIdentity       string         `json:"scope_identity,omitempty"`
	WorkCategory        string         `json:"work_category,omitempty"`
	ProfileSource       string         `json:"profile_source,omitempty"`
	SelectionRole       string         `json:"selection_role,omitempty"`
	FallbackIndex       int            `json:"fallback_index,omitempty"`
	Runtime             string         `json:"runtime,omitempty"`
	Model               string         `json:"model,omitempty"`
	ReasoningEffort     string         `json:"reasoning_effort,omitempty"`
	NextRuntime         string         `json:"next_runtime,omitempty"`
	NextModel           string         `json:"next_model,omitempty"`
	NextReasoningEffort string         `json:"next_reasoning_effort,omitempty"`
	ReasonCode          string         `json:"reason_code,omitempty"`
	Reason              string         `json:"reason,omitempty"`
	NextAction          string         `json:"next_action,omitempty"`
	ReviewIssuesKnown   *bool          `json:"review_issues_known,omitempty"`
	ConsoleLog          string         `json:"console_log,omitempty"`
	AttachCommand       string         `json:"attach_command,omitempty"`
	EvidenceState       string         `json:"evidence_state,omitempty"`
	EvidenceKind        string         `json:"evidence_kind,omitempty"`
	EvidenceHeadSHA     string         `json:"evidence_head_sha,omitempty"`
	VerifiedHeadSHA     string         `json:"verified_head_sha,omitempty"`
}

// StreamCategoryFilter selects public Run Event Stream categories.
type StreamCategoryFilter map[StreamCategory]struct{}

// AllStreamCategories returns a filter that includes every stable category.
func AllStreamCategories() StreamCategoryFilter {
	return StreamCategoryFilter{
		StreamCategoryTaskStatus:   {},
		StreamCategoryBatch:        {},
		StreamCategoryVerification: {},
		StreamCategoryOutcome:      {},
		StreamCategorySelection:    {},
	}
}

// Includes reports whether the category is selected. A nil filter includes
// every stable category so callers can use nil as the default.
func (filter StreamCategoryFilter) Includes(category StreamCategory) bool {
	if filter == nil {
		return true
	}
	_, ok := filter[category]
	return ok
}

// ParseStreamCategoryFilter parses the comma-separated --filter value.
func ParseStreamCategoryFilter(raw string) (StreamCategoryFilter, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("event filter is required")
	}
	filter := StreamCategoryFilter{}
	for _, part := range strings.Split(raw, ",") {
		category := StreamCategory(strings.TrimSpace(part))
		if category == "" {
			return nil, fmt.Errorf("event filter contains an empty category")
		}
		if !isStreamCategory(category) {
			return nil, fmt.Errorf("unknown event category %q", category)
		}
		filter[category] = struct{}{}
	}
	return filter, nil
}

// ProjectStreamEvent maps one internal Run Event to the stable Supervisor
// stream. Unknown or unselected events are skipped. Malformed selected daemon
// payloads fail instead of falling back to human summaries.
func ProjectStreamEvent(cursor int64, event RunEvent, filter StreamCategoryFilter) (StreamRecord, bool, error) {
	category, ok := streamCategoryForEvent(event)
	if !ok || !filter.Includes(category) {
		return StreamRecord{}, false, nil
	}
	fields, err := streamPayloadFields(event)
	if err != nil {
		return StreamRecord{}, false, err
	}
	record := StreamRecord{
		Schema:   StreamSchema,
		RunID:    event.RunID,
		Category: category,
		Time:     streamTime(event.Time),
		Cursor:   cursor,
		Batch:    event.Batch,
		WorkItem: event.ReviewIssue,
	}
	switch category {
	case StreamCategoryTaskStatus:
		if record.WorkItem == "" {
			record.WorkItem, err = firstPayloadString(fields, "work_item", "task")
			if err != nil {
				return StreamRecord{}, false, err
			}
			if record.WorkItem == "" {
				return StreamRecord{}, false, streamMissingField(event, "task")
			}
		}
		record.Phase, err = requiredPayloadString(fields, event, "phase")
		if err != nil {
			return StreamRecord{}, false, err
		}
		record.Status, err = optionalPayloadString(fields, "status")
		if err != nil {
			return StreamRecord{}, false, err
		}
		record.Summary = taskStatusStreamSummary(record.WorkItem, record.Phase, record.Status)
	case StreamCategoryBatch:
		record.Phase, err = requiredPayloadString(fields, event, "phase")
		if err != nil {
			return StreamRecord{}, false, err
		}
		if record.Batch == 0 {
			record.Batch, err = optionalPayloadInt(fields, "batch")
			if err != nil {
				return StreamRecord{}, false, err
			}
		}
		record.Summary = batchStreamSummary(record.Batch, record.Phase)
	case StreamCategoryVerification:
		if err := projectVerificationRecord(&record, fields, event); err != nil {
			return StreamRecord{}, false, err
		}
	case StreamCategoryOutcome:
		if err := projectOutcomeRecord(&record, fields, event); err != nil {
			return StreamRecord{}, false, err
		}
	case StreamCategorySelection:
		if event.Kind == KindAgentSelectionReceipt {
			receipt, ok, err := ProjectSelectionReceipt(event)
			if err != nil {
				return StreamRecord{}, false, err
			}
			if !ok {
				return StreamRecord{}, false, nil
			}
			record.ScopeKind = "agent_session"
			record.ScopeID = receipt.Session
			record.ScopeIdentity = record.ScopeKind + ":" + record.ScopeID
			record.SelectionRole = "effective"
			record.Status = receipt.Status
			record.Runtime = receipt.Runtime
			record.Model = receipt.Model
			record.ReasoningEffort = receipt.ReasoningEffort
			record.Summary = SelectionReceiptLine(receipt)
			break
		}
		selection, ok, err := ProjectSelectionLifecycle(event)
		if err != nil {
			return StreamRecord{}, false, err
		}
		if !ok {
			return StreamRecord{}, false, nil
		}
		record.ScopeKind = selection.ScopeKind
		record.ScopeID = selection.ScopeID
		record.ScopeIdentity = selection.ScopeIdentity
		record.WorkCategory = selection.Category
		record.ProfileSource = selection.ProfileSource
		record.Attempt = selection.Attempt
		record.SelectionRole = selection.SelectionRole
		record.FallbackIndex = selection.FallbackIndex
		record.Status = selection.Status
		record.Runtime = selection.Runtime
		record.Model = selection.Model
		record.ReasoningEffort = selection.ReasoningEffort
		record.NextRuntime = selection.NextRuntime
		record.NextModel = selection.NextModel
		record.NextReasoningEffort = selection.NextReasoningEffort
		record.ReasonCode = selection.ReasonCode
		record.Reason = selection.Reason
		record.Summary = SelectionLifecycleLine(selection)
	}
	return record, true, nil
}

func projectVerificationRecord(record *StreamRecord, fields map[string]json.RawMessage, event RunEvent) error {
	var err error
	record.Phase, err = requiredPayloadString(fields, event, "phase")
	if err != nil {
		return err
	}
	record.Attempt, err = verificationPayloadAttempt(fields, event, record.Phase)
	if err != nil {
		return err
	}
	record.Retry, err = optionalPayloadInt(fields, "retry")
	if err != nil {
		return streamPayloadFieldError(event, "retry", err)
	}
	record.Mode, err = readOptionalString(fields, event, "mode")
	if err != nil {
		return err
	}
	record.Classification, err = readOptionalString(fields, event, "classification")
	if err != nil {
		return err
	}
	record.RetryAvailable, err = optionalPayloadBool(fields, "retry_available")
	if err != nil {
		return streamPayloadFieldError(event, "retry_available", err)
	}
	record.Reason, err = readOptionalString(fields, event, "reason")
	if err != nil {
		return err
	}
	if record.WorkItem == "" {
		record.WorkItem, err = firstPayloadString(fields, "work_item", "task")
		if err != nil {
			return err
		}
	}
	if record.Batch == 0 {
		record.Batch, err = optionalPayloadInt(fields, "batch")
		if err != nil {
			return err
		}
	}
	if record.Phase == string(VerificationPhaseVerdict) {
		record.Verdict, err = requiredPayloadString(fields, event, "verdict")
	} else {
		record.Verdict, err = optionalPayloadString(fields, "verdict")
	}
	if err != nil {
		return err
	}
	record.Summary = verificationStreamSummary(record.Attempt, record.Phase, record.Verdict)
	return nil
}

func projectOutcomeRecord(record *StreamRecord, fields map[string]json.RawMessage, event RunEvent) error {
	var err error
	record.Outcome, err = firstPayloadString(fields, "outcome", "state")
	if err != nil {
		return err
	}
	record.Reason, err = readOptionalString(fields, event, "reason")
	if err != nil {
		return err
	}
	record.NextAction, err = readOptionalString(fields, event, "next_action")
	if err != nil {
		return err
	}
	record.ReviewIssuesKnown, err = optionalPayloadBool(fields, "review_issues_known")
	if err != nil {
		return streamPayloadFieldError(event, "review_issues_known", err)
	}
	record.ConsoleLog, err = readOptionalString(fields, event, "console_log")
	if err != nil {
		return err
	}
	record.AttachCommand, err = readOptionalString(fields, event, "attach_command")
	if err != nil {
		return err
	}
	record.EvidenceState, err = readOptionalString(fields, event, "evidence_state")
	if err != nil {
		return err
	}
	record.EvidenceKind, err = readOptionalString(fields, event, "evidence_kind")
	if err != nil {
		return err
	}
	record.EvidenceHeadSHA, err = readOptionalString(fields, event, "evidence_head_sha")
	if err != nil {
		return err
	}
	record.VerifiedHeadSHA, err = readOptionalString(fields, event, "verified_head_sha")
	if err != nil {
		return err
	}
	record.Summary = outcomeStreamSummary(record.Outcome)
	return nil
}

func isStreamCategory(category StreamCategory) bool {
	switch category {
	case StreamCategoryTaskStatus, StreamCategoryBatch, StreamCategoryVerification, StreamCategoryOutcome, StreamCategorySelection:
		return true
	default:
		return false
	}
}

func streamCategoryForEvent(event RunEvent) (StreamCategory, bool) {
	if event.Source == SourceAgent && event.Kind == KindAgentSelectionReceipt {
		return StreamCategorySelection, true
	}
	if event.Source != SourceDaemon {
		return "", false
	}
	switch event.Kind {
	case KindDaemonTask:
		return StreamCategoryTaskStatus, true
	case KindDaemonBatch:
		return StreamCategoryBatch, true
	case KindDaemonVerification:
		return StreamCategoryVerification, true
	case KindDaemonOutcome:
		return StreamCategoryOutcome, true
	case KindDaemonAgentSelectionAttempt, KindDaemonAgentSelectionActive,
		KindDaemonAgentSelectionFallback, KindDaemonAgentSelectionExhausted,
		KindDaemonAgentSelectionClosed:
		return StreamCategorySelection, true
	default:
		return "", false
	}
}

func streamPayloadFields(event RunEvent) (map[string]json.RawMessage, error) {
	if len(event.Payload) == 0 {
		return nil, fmt.Errorf("project %s event for Run %q: payload is required", event.Kind, event.RunID)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(event.Payload, &fields); err != nil {
		return nil, fmt.Errorf("project %s event for Run %q: decode payload: %w", event.Kind, event.RunID, err)
	}
	if fields == nil {
		return nil, fmt.Errorf("project %s event for Run %q: payload object is required", event.Kind, event.RunID)
	}
	return fields, nil
}

func requiredPayloadString(fields map[string]json.RawMessage, event RunEvent, key string) (string, error) {
	value, err := optionalPayloadString(fields, key)
	if err != nil {
		return "", fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, key, err)
	}
	if value == "" {
		return "", streamMissingField(event, key)
	}
	return value, nil
}

func optionalPayloadString(fields map[string]json.RawMessage, key string) (string, error) {
	raw, ok := fields[key]
	if !ok {
		return "", nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", err
	}
	return value, nil
}

func readOptionalString(fields map[string]json.RawMessage, event RunEvent, key string) (string, error) {
	value, err := optionalPayloadString(fields, key)
	if err != nil {
		return "", streamPayloadFieldError(event, key, err)
	}
	return value, nil
}

func optionalPayloadBool(fields map[string]json.RawMessage, key string) (*bool, error) {
	raw, ok := fields[key]
	if !ok {
		return nil, nil
	}
	var value bool
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func firstPayloadString(fields map[string]json.RawMessage, keys ...string) (string, error) {
	for _, key := range keys {
		value, err := optionalPayloadString(fields, key)
		if err != nil {
			return "", fmt.Errorf("read payload field %q: %w", key, err)
		}
		if value != "" {
			return value, nil
		}
	}
	return "", nil
}

// verificationPayloadAttempt maps the three pre-repair Verification payload
// shapes to their single historical attempt without accepting malformed events
// emitted by the current attempt-aware producer.
func verificationPayloadAttempt(fields map[string]json.RawMessage, event RunEvent, phase string) (int, error) {
	attempt, err := optionalPayloadInt(fields, "attempt")
	if err != nil {
		return 0, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "attempt", err)
	}
	if attempt != 0 {
		return attempt, nil
	}
	legacy, err := isLegacyVerificationPayload(fields, event, phase)
	if err != nil {
		return 0, err
	}
	if legacy {
		return 1, nil
	}
	return 0, streamMissingField(event, "attempt")
}

func isLegacyVerificationPayload(fields map[string]json.RawMessage, event RunEvent, phase string) (bool, error) {
	var summaryPrefix string
	switch phase {
	case "started":
		summaryPrefix = "Verification started:"
	case "failed":
		summaryPrefix = "Verification failed:"
	case "passed":
		summaryPrefix = "Verification command passed:"
	default:
		return false, nil
	}
	if !strings.HasPrefix(event.Summary, summaryPrefix) {
		return false, nil
	}
	command, err := optionalPayloadString(fields, "command")
	if err != nil {
		return false, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "command", err)
	}
	if command == "" {
		return false, nil
	}
	if phase != "failed" {
		return true, nil
	}
	failure, err := optionalPayloadString(fields, "error")
	if err != nil {
		return false, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "error", err)
	}
	return failure != "", nil
}

func optionalPayloadInt(fields map[string]json.RawMessage, key string) (int, error) {
	raw, ok := fields[key]
	if !ok {
		return 0, nil
	}
	var value int
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func streamMissingField(event RunEvent, key string) error {
	return fmt.Errorf("project %s event for Run %q: missing payload field %q", event.Kind, event.RunID, key)
}

func streamPayloadFieldError(event RunEvent, key string, err error) error {
	return fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, key, err)
}

func streamTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func verificationStreamSummary(attempt int, phase string, verdict string) string {
	if phase == string(VerificationPhaseVerdict) && verdict != "" {
		return fmt.Sprintf("Verification attempt %d verdict: %s.", attempt, verdict)
	}
	return fmt.Sprintf("Verification attempt %d phase %s.", attempt, phase)
}

func taskStatusStreamSummary(workItem string, phase string, status string) string {
	if status != "" {
		return fmt.Sprintf("Task %s phase %s status %s.", workItem, phase, status)
	}
	return fmt.Sprintf("Task %s phase %s.", workItem, phase)
}

func batchStreamSummary(batch int, phase string) string {
	if batch > 0 {
		return fmt.Sprintf("Batch %03d phase %s.", batch, phase)
	}
	return fmt.Sprintf("Batch phase %s.", phase)
}

func outcomeStreamSummary(outcome string) string {
	if outcome != "" {
		return fmt.Sprintf("Run reached %s.", outcome)
	}
	return "Run outcome recorded."
}
