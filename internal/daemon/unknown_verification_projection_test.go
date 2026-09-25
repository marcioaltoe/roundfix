package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

type capturedVerificationPublication struct {
	summary string
	payload map[string]any
}

func TestUnobservedVerificationProjectsAsUnknown(t *testing.T) {
	t.Parallel()
	const (
		command        = "make verify"
		reason         = "runner lost the command verdict"
		diagnosticPath = "/tmp/verification.log"
	)
	publications, publish := captureVerificationPublications()
	req := verificationAttemptRequest{
		BatchNumber: 4,
		WorkItem:    "task_04",
		Attempt:     2,
		Publish:     publish,
	}

	err := req.publishUnknownFailure(context.Background(), command, &VerificationUnknownError{
		Command:        command,
		DiagnosticPath: diagnosticPath,
		Err:            errors.New(reason),
	})

	if err != nil {
		t.Fatalf("publish unobserved Verification: %v", err)
	}
	records := projectVerificationPublications(t, publications)
	assertUnknownVerificationRecords(t, records, command, reason, diagnosticPath)
}

func TestUnobservedVerificationWithoutDiagnosticProjectsUnavailable(t *testing.T) {
	t.Parallel()
	const command = "make verify"
	publications, publish := captureVerificationPublications()
	req := verificationAttemptRequest{
		BatchNumber: 4,
		WorkItem:    "task_04",
		Attempt:     1,
		Publish:     publish,
	}

	err := req.publishUnknownFailure(context.Background(), command, &VerificationUnknownError{})

	if err != nil {
		t.Fatalf("publish unobserved Verification without diagnostics: %v", err)
	}
	records := projectVerificationPublications(t, publications)
	assertUnknownVerificationRecords(t, records, command, "reason unavailable", "unavailable")
}

func TestUnobservedPreconditionVerificationProjectsAsUnknown(t *testing.T) {
	t.Parallel()
	const (
		command        = "make verify"
		reason         = "runner stopped before observing the precondition verdict"
		diagnosticPath = "/tmp/precondition.log"
	)
	publications, publish := captureVerificationPublications()
	req := verificationAttemptRequest{
		BatchNumber:           4,
		WorkItem:              "task_04",
		Attempt:               1,
		FailureClassification: runevent.VerificationClassificationPrecondition,
		FailureReason:         runevent.VerificationReasonRepositoryNotGreenOnEntry,
		Publish:               publish,
	}

	err := req.publishUnknownFailure(context.Background(), command, &VerificationUnknownError{
		Command:        command,
		DiagnosticPath: diagnosticPath,
		Err:            errors.New(reason),
	})

	if err != nil {
		t.Fatalf("publish unobserved precondition Verification: %v", err)
	}
	records := projectVerificationPublications(t, publications)
	assertUnknownVerificationRecords(t, records, command, reason, diagnosticPath)
}

func TestDeterministicVerificationFailureProjectsUnclassified(t *testing.T) {
	t.Parallel()
	const (
		command        = "make verify"
		diagnosticPath = "/tmp/deterministic.log"
	)
	publications, publish := captureVerificationPublications()
	req := verificationAttemptRequest{
		BatchNumber: 4,
		WorkItem:    "task_04",
		Attempt:     1,
		Publish:     publish,
	}
	commandErr := &VerificationCommandError{
		Command:    command,
		OutputPath: diagnosticPath,
		Err:        errors.New("exit status 1"),
	}

	metadata, err := req.publishFailedCommand(context.Background(), command, commandErr, false)
	if err != nil {
		t.Fatalf("publish deterministic Verification failure: %v", err)
	}
	if err := req.publishVerdict(context.Background(), runevent.VerificationVerdictFailed, diagnosticPath, "", false, metadata); err != nil {
		t.Fatalf("publish deterministic Verification verdict: %v", err)
	}

	records := projectVerificationPublications(t, publications)
	if len(records) != 2 {
		t.Fatalf("projected records = %d, want 2", len(records))
	}
	for index, record := range records {
		if record.Classification != "" {
			t.Fatalf("record %d classification = %q, want empty", index, record.Classification)
		}
		if record.Reason != "" {
			t.Fatalf("record %d reason = %q, want empty", index, record.Reason)
		}
	}
	if records[0].Phase != string(runevent.VerificationPhaseFailed) {
		t.Fatalf("failed phase = %q, want %q", records[0].Phase, runevent.VerificationPhaseFailed)
	}
	if records[0].Verdict != "" {
		t.Fatalf("failed verdict = %q, want empty", records[0].Verdict)
	}
	if records[1].Phase != string(runevent.VerificationPhaseVerdict) {
		t.Fatalf("verdict phase = %q, want %q", records[1].Phase, runevent.VerificationPhaseVerdict)
	}
	if records[1].Verdict != string(runevent.VerificationVerdictFailed) {
		t.Fatalf("verdict = %q, want %q", records[1].Verdict, runevent.VerificationVerdictFailed)
	}
}

func captureVerificationPublications() (*[]capturedVerificationPublication, func(context.Context, string, map[string]any) error) {
	publications := []capturedVerificationPublication{}
	return &publications, func(_ context.Context, summary string, payload map[string]any) error {
		publications = append(publications, capturedVerificationPublication{summary: summary, payload: payload})
		return nil
	}
}

func projectVerificationPublications(t *testing.T, publications *[]capturedVerificationPublication) []runevent.StreamRecord {
	t.Helper()
	if publications == nil {
		t.Fatal("captured publications are nil")
	}
	records := make([]runevent.StreamRecord, 0, len(*publications))
	for index, publication := range *publications {
		payload, err := json.Marshal(publication.payload)
		if err != nil {
			t.Fatalf("encode publication %d: %v", index, err)
		}
		event := runevent.RunEvent{
			RunID:       "run_projection",
			Batch:       4,
			Source:      runevent.SourceDaemon,
			Kind:        runevent.KindDaemonVerification,
			ReviewIssue: "task_04",
			Summary:     publication.summary,
			Time:        time.Date(2026, 9, 25, 12, 0, index, 0, time.UTC),
			Payload:     payload,
		}
		record, ok, err := runevent.ProjectStreamEvent(int64(index+1), event, runevent.AllStreamCategories())
		if err != nil {
			t.Fatalf("project publication %d: %v", index, err)
		}
		if !ok {
			t.Fatalf("publication %d did not project", index)
		}
		records = append(records, record)
	}
	return records
}

func assertUnknownVerificationRecords(t *testing.T, records []runevent.StreamRecord, command string, reason string, diagnosticPath string) {
	t.Helper()
	if len(records) != 2 {
		t.Fatalf("projected records = %d, want 2", len(records))
	}
	for index, record := range records {
		if record.Schema != runevent.StreamSchema {
			t.Fatalf("record %d schema = %q, want %q", index, record.Schema, runevent.StreamSchema)
		}
		if record.RunID != "run_projection" {
			t.Fatalf("record %d run id = %q, want run_projection", index, record.RunID)
		}
		if record.Category != runevent.StreamCategoryVerification {
			t.Fatalf("record %d category = %q, want %q", index, record.Category, runevent.StreamCategoryVerification)
		}
		if record.Cursor != int64(index+1) {
			t.Fatalf("record %d cursor = %d, want %d", index, record.Cursor, index+1)
		}
		if record.Batch != 4 {
			t.Fatalf("record %d batch = %d, want 4", index, record.Batch)
		}
		if record.WorkItem != "task_04" {
			t.Fatalf("record %d work item = %q, want task_04", index, record.WorkItem)
		}
		if record.Classification != string(runevent.VerificationClassificationUnknown) {
			t.Fatalf("record %d classification = %q, want %q", index, record.Classification, runevent.VerificationClassificationUnknown)
		}
		if record.Command != command {
			t.Fatalf("record %d command = %q, want %q", index, record.Command, command)
		}
		if record.Reason != reason {
			t.Fatalf("record %d reason = %q, want %q", index, record.Reason, reason)
		}
		if record.DiagnosticPath != diagnosticPath {
			t.Fatalf("record %d diagnostic path = %q, want %q", index, record.DiagnosticPath, diagnosticPath)
		}
	}
	if records[0].Phase != string(runevent.VerificationPhaseFailed) {
		t.Fatalf("failed phase = %q, want %q", records[0].Phase, runevent.VerificationPhaseFailed)
	}
	if records[0].Verdict != "" {
		t.Fatalf("failed verdict = %q, want empty", records[0].Verdict)
	}
	if records[1].Phase != string(runevent.VerificationPhaseVerdict) {
		t.Fatalf("verdict phase = %q, want %q", records[1].Phase, runevent.VerificationPhaseVerdict)
	}
	if records[1].Verdict != string(runevent.VerificationVerdictFailed) {
		t.Fatalf("verdict = %q, want %q", records[1].Verdict, runevent.VerificationVerdictFailed)
	}
}
