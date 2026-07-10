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
)

// StreamRecord is one JSONL record in the Supervisor Run Event Stream.
type StreamRecord struct {
	Schema   string         `json:"schema"`
	RunID    string         `json:"run_id"`
	Category StreamCategory `json:"category"`
	Time     string         `json:"time"`
	Cursor   int64          `json:"cursor"`
	Batch    int            `json:"batch,omitempty"`
	Attempt  int            `json:"attempt,omitempty"`
	WorkItem string         `json:"work_item,omitempty"`
	Phase    string         `json:"phase,omitempty"`
	Status   string         `json:"status,omitempty"`
	Verdict  string         `json:"verdict,omitempty"`
	Outcome  string         `json:"outcome,omitempty"`
	Summary  string         `json:"summary,omitempty"`
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
		record.Attempt, err = requiredPayloadInt(fields, event, "attempt")
		if err != nil {
			return StreamRecord{}, false, err
		}
		record.Phase, err = requiredPayloadString(fields, event, "phase")
		if err != nil {
			return StreamRecord{}, false, err
		}
		if record.WorkItem == "" {
			record.WorkItem, err = firstPayloadString(fields, "work_item", "task")
			if err != nil {
				return StreamRecord{}, false, err
			}
		}
		if record.Batch == 0 {
			record.Batch, err = optionalPayloadInt(fields, "batch")
			if err != nil {
				return StreamRecord{}, false, err
			}
		}
		if record.Phase == string(VerificationPhaseVerdict) {
			record.Verdict, err = requiredPayloadString(fields, event, "verdict")
			if err != nil {
				return StreamRecord{}, false, err
			}
		} else {
			record.Verdict, err = optionalPayloadString(fields, "verdict")
			if err != nil {
				return StreamRecord{}, false, err
			}
		}
		record.Summary = verificationStreamSummary(record.Attempt, record.Phase, record.Verdict)
	case StreamCategoryOutcome:
		record.Outcome, err = firstPayloadString(fields, "outcome", "state")
		if err != nil {
			return StreamRecord{}, false, err
		}
		record.Summary = outcomeStreamSummary(record.Outcome)
	}
	return record, true, nil
}

func isStreamCategory(category StreamCategory) bool {
	switch category {
	case StreamCategoryTaskStatus, StreamCategoryBatch, StreamCategoryVerification, StreamCategoryOutcome:
		return true
	default:
		return false
	}
}

func streamCategoryForEvent(event RunEvent) (StreamCategory, bool) {
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

func requiredPayloadInt(fields map[string]json.RawMessage, event RunEvent, key string) (int, error) {
	value, err := optionalPayloadInt(fields, key)
	if err != nil {
		return 0, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, key, err)
	}
	if value == 0 {
		return 0, streamMissingField(event, key)
	}
	return value, nil
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
