package runevent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SelectionLifecycleRecord struct {
	Event               string
	ScopeKind           string
	ScopeID             string
	ScopeIdentity       string
	Category            string
	ProfileSource       string
	Attempt             int
	SelectionRole       string
	FallbackIndex       int
	Runtime             string
	Model               string
	ReasoningEffort     string
	Status              string
	ReasonCode          string
	Reason              string
	NextRuntime         string
	NextModel           string
	NextReasoningEffort string
}

func ProjectSelectionLifecycle(event RunEvent) (SelectionLifecycleRecord, bool, error) {
	if event.Source != SourceDaemon || !isSelectionLifecycleKind(event.Kind) {
		return SelectionLifecycleRecord{}, false, nil
	}
	fields, err := selectionPayloadFields(event)
	if err != nil {
		return SelectionLifecycleRecord{}, true, err
	}
	if _, hasAttempt := fields["attempt"]; hasAttempt {
		record, err := projectSelectionAttemptLifecycle(event, fields)
		return record, true, err
	}
	record, err := projectSelectionNotificationLifecycle(event, fields)
	return record, true, err
}

func SelectionLifecycleLine(record SelectionLifecycleRecord) string {
	summary := SelectionLifecycleSummary(record)
	if summary == "" {
		return ""
	}
	return BoundSummary("SELECTION " + summary)
}

func SelectionLifecycleSummary(record SelectionLifecycleRecord) string {
	scope := strings.TrimSpace(record.ScopeKind + " " + record.ScopeID)
	if strings.TrimSpace(record.Category) != "" {
		scope += " (" + strings.TrimSpace(record.Category) + ")"
	}
	parts := []string{strings.TrimSpace(scope)}
	if record.Attempt > 0 {
		parts = append(parts, fmt.Sprintf("attempt %d", record.Attempt))
	}
	if strings.TrimSpace(record.SelectionRole) != "" {
		parts = append(parts, strings.TrimSpace(record.SelectionRole))
	}
	if strings.TrimSpace(record.Status) != "" {
		parts = append(parts, strings.TrimSpace(record.Status))
	}
	source := strings.TrimSpace(record.ProfileSource)
	if source == "" {
		source = "unavailable"
	}
	parts = append(parts, "source "+source)
	if selection := selectionLifecycleTuple(record.Runtime, record.Model, record.ReasoningEffort); selection != "" {
		parts = append(parts, selection)
	}
	if next := selectionLifecycleTuple(record.NextRuntime, record.NextModel, record.NextReasoningEffort); next != "" {
		if record.FallbackIndex > 0 {
			parts = append(parts, fmt.Sprintf("-> fallback %d %s", record.FallbackIndex, next))
		} else {
			parts = append(parts, "-> "+next)
		}
	}
	if strings.TrimSpace(record.ReasonCode) != "" {
		reason := strings.TrimSpace(record.ReasonCode)
		if strings.TrimSpace(record.Reason) != "" {
			reason += ": " + strings.TrimSpace(record.Reason)
		}
		parts = append(parts, "reason "+reason)
	}
	return BoundSummary(strings.Join(nonEmptySelectionParts(parts), " "))
}

func SelectionLifecycleScopeKey(record SelectionLifecycleRecord) string {
	scopeIdentity := strings.TrimSpace(record.ScopeIdentity)
	if scopeIdentity != "" {
		return scopeIdentity
	}
	if strings.TrimSpace(record.ScopeKind) == "" || strings.TrimSpace(record.ScopeID) == "" {
		return ""
	}
	return strings.TrimSpace(record.ScopeKind) + ":" + strings.TrimSpace(record.ScopeID)
}

func isSelectionLifecycleKind(kind Kind) bool {
	switch kind {
	case KindDaemonAgentSelectionAttempt, KindDaemonAgentSelectionActive,
		KindDaemonAgentSelectionFallback, KindDaemonAgentSelectionExhausted,
		KindDaemonAgentSelectionClosed:
		return true
	default:
		return false
	}
}

func projectSelectionAttemptLifecycle(event RunEvent, fields map[string]json.RawMessage) (SelectionLifecycleRecord, error) {
	record := SelectionLifecycleRecord{}
	var err error
	record.Event, _ = selectionOptionalString(fields, "event")
	record.ScopeKind, err = selectionRequiredString(fields, event, "scope_kind")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ScopeID, err = selectionRequiredString(fields, event, "scope_id")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ScopeIdentity, _ = selectionOptionalString(fields, "scope_identity")
	if record.ScopeIdentity == "" {
		record.ScopeIdentity = record.ScopeKind + ":" + record.ScopeID
	}
	record.Category, err = selectionRequiredString(fields, event, "category")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ProfileSource, _ = selectionOptionalString(fields, "profile_source")
	record.Attempt, err = selectionOptionalInt(fields, "attempt")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "attempt", err)
	}
	if record.Attempt <= 0 {
		return SelectionLifecycleRecord{}, selectionMissingField(event, "attempt")
	}
	record.SelectionRole, err = selectionRequiredString(fields, event, "selection_role")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.FallbackIndex, err = selectionOptionalInt(fields, "fallback_index")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "fallback_index", err)
	}
	record.Runtime, err = selectionRequiredString(fields, event, "runtime")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.Model, err = selectionRequiredString(fields, event, "model")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ReasoningEffort, err = selectionOptionalString(fields, "reasoning_effort")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "reasoning_effort", err)
	}
	record.Status, err = selectionRequiredString(fields, event, "status")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ReasonCode, err = selectionOptionalString(fields, "reason_code")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "reason_code", err)
	}
	record.Reason, err = selectionOptionalString(fields, "reason")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "reason", err)
	}
	return record, nil
}

func projectSelectionNotificationLifecycle(event RunEvent, fields map[string]json.RawMessage) (SelectionLifecycleRecord, error) {
	record := SelectionLifecycleRecord{Status: selectionNotificationStatus(event.Kind)}
	var err error
	record.Event, _ = selectionOptionalString(fields, "event")
	record.ScopeKind, err = selectionRequiredString(fields, event, "scope_kind")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ScopeID, err = selectionRequiredString(fields, event, "scope_id")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.ScopeIdentity, _ = selectionOptionalString(fields, "scope_identity")
	if record.ScopeIdentity == "" {
		record.ScopeIdentity = record.ScopeKind + ":" + record.ScopeID
	}
	record.Category, err = selectionRequiredString(fields, event, "category")
	if err != nil {
		return SelectionLifecycleRecord{}, err
	}
	record.SelectionRole = "fallback"
	record.FallbackIndex, err = selectionOptionalInt(fields, "fallback_index")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "fallback_index", err)
	}
	record.ReasonCode, err = selectionOptionalString(fields, "reason_code")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "reason_code", err)
	}
	record.Reason, err = selectionOptionalString(fields, "reason")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, "reason", err)
	}
	record.Runtime, err = selectionNestedString(fields, "failed_selection", "runtime")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read failed selection runtime: %w", event.Kind, event.RunID, err)
	}
	record.Model, err = selectionNestedString(fields, "failed_selection", "model")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read failed selection model: %w", event.Kind, event.RunID, err)
	}
	record.ReasoningEffort, err = selectionNestedString(fields, "failed_selection", "reasoning_effort")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read failed selection reasoning: %w", event.Kind, event.RunID, err)
	}
	record.NextRuntime, err = selectionNestedString(fields, "next_selection", "runtime")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read next selection runtime: %w", event.Kind, event.RunID, err)
	}
	record.NextModel, err = selectionNestedString(fields, "next_selection", "model")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read next selection model: %w", event.Kind, event.RunID, err)
	}
	record.NextReasoningEffort, err = selectionNestedString(fields, "next_selection", "reasoning_effort")
	if err != nil {
		return SelectionLifecycleRecord{}, fmt.Errorf("project %s event for Run %q: read next selection reasoning: %w", event.Kind, event.RunID, err)
	}
	return record, nil
}

func selectionNotificationStatus(kind Kind) string {
	switch kind {
	case KindDaemonAgentSelectionFallback:
		return "failed"
	case KindDaemonAgentSelectionExhausted:
		return "exhausted"
	case KindDaemonAgentSelectionClosed:
		return "closed"
	case KindDaemonAgentSelectionActive:
		return "active"
	default:
		return "notified"
	}
}

func selectionPayloadFields(event RunEvent) (map[string]json.RawMessage, error) {
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

func selectionRequiredString(fields map[string]json.RawMessage, event RunEvent, key string) (string, error) {
	value, err := selectionOptionalString(fields, key)
	if err != nil {
		return "", fmt.Errorf("project %s event for Run %q: read payload field %q: %w", event.Kind, event.RunID, key, err)
	}
	if value == "" {
		return "", selectionMissingField(event, key)
	}
	return value, nil
}

func selectionOptionalString(fields map[string]json.RawMessage, key string) (string, error) {
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

func selectionOptionalInt(fields map[string]json.RawMessage, key string) (int, error) {
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

func selectionNestedString(fields map[string]json.RawMessage, objectKey string, key string) (string, error) {
	raw, ok := fields[objectKey]
	if !ok {
		return "", nil
	}
	var nested map[string]json.RawMessage
	if err := json.Unmarshal(raw, &nested); err != nil {
		return "", err
	}
	if nested == nil {
		return "", nil
	}
	return selectionOptionalString(nested, key)
}

func selectionMissingField(event RunEvent, key string) error {
	return fmt.Errorf("project %s event for Run %q: missing payload field %q", event.Kind, event.RunID, key)
}

func selectionLifecycleTuple(runtime string, model string, reasoningEffort string) string {
	runtime = strings.TrimSpace(runtime)
	model = strings.TrimSpace(model)
	if runtime == "" && model == "" {
		return ""
	}
	if runtime == "" {
		runtime = "-"
	}
	if model == "" {
		model = "-"
	}
	effort := strings.TrimSpace(reasoningEffort)
	if effort == "" {
		effort = "<model-managed>"
	}
	return runtime + "/" + model + "/" + effort
}

func nonEmptySelectionParts(parts []string) []string {
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			kept = append(kept, strings.TrimSpace(part))
		}
	}
	return kept
}
