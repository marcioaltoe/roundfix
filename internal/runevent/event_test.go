package runevent

import (
	"context"
	"errors"
	"testing"
)

func TestSpecRunKindsUseDaemonNamespace(t *testing.T) {
	if KindDaemonTask != "daemon.task" {
		t.Fatalf("expected daemon.task kind, got %q", KindDaemonTask)
	}
	if KindDaemonQA != "daemon.qa" {
		t.Fatalf("expected daemon.qa kind, got %q", KindDaemonQA)
	}
}

func TestVerificationEventVocabulary(t *testing.T) {
	if VerificationPhaseStarted != "started" {
		t.Fatalf("expected started phase, got %q", VerificationPhaseStarted)
	}
	if VerificationPhaseCommandPassed != "command-passed" {
		t.Fatalf("expected command-passed phase, got %q", VerificationPhaseCommandPassed)
	}
	if VerificationPhaseFailed != "failed" {
		t.Fatalf("expected failed phase, got %q", VerificationPhaseFailed)
	}
	if VerificationPhaseVerdict != "verdict" {
		t.Fatalf("expected verdict phase, got %q", VerificationPhaseVerdict)
	}
	if VerificationVerdictPassed != "passed" {
		t.Fatalf("expected passed verdict, got %q", VerificationVerdictPassed)
	}
	if VerificationVerdictFailed != "failed" {
		t.Fatalf("expected failed verdict, got %q", VerificationVerdictFailed)
	}
}

func TestIsDaemonKindCoversSpecRunKindsAndSkipsUnknown(t *testing.T) {
	tests := []struct {
		name     string
		kind     Kind
		expected bool
	}{
		{name: "task kind is daemon vocabulary", kind: KindDaemonTask, expected: true},
		{name: "qa kind is daemon vocabulary", kind: KindDaemonQA, expected: true},
		{name: "agent kind is not daemon vocabulary", kind: KindAgentMessage, expected: false},
		{name: "unknown daemon-prefixed kind stays skippable", kind: Kind("daemon.unknown"), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDaemonKind(tt.kind); got != tt.expected {
				t.Fatalf("IsDaemonKind(%q) = %v, want %v", tt.kind, got, tt.expected)
			}
		})
	}
}

func TestSourceFilterSinkDropsOnlyConfiguredSource(t *testing.T) {
	next := &recordingSink{}
	sink := NewSourceFilterSink(next, SourceAgent)
	events := []RunEvent{
		{Source: SourceDaemon, Kind: KindDaemonSelection, Summary: "selected"},
		{Source: SourceAgent, Kind: KindAgentRaw, Summary: "agent output"},
		{Source: SourceVerification, Kind: KindDaemonVerification, Summary: "verified"},
		{Source: SourceAgent, Kind: KindAgentStatus, Summary: "agent status"},
		{Source: SourceGit, Kind: KindDaemonCommit, Summary: "committed"},
	}

	for _, event := range events {
		if err := sink.Publish(context.Background(), event); err != nil {
			t.Fatalf("publish %s: %v", event.Summary, err)
		}
	}

	if len(next.events) != 3 {
		t.Fatalf("expected three forwarded events, got %#v", next.events)
	}
	for index, want := range []string{"selected", "verified", "committed"} {
		if next.events[index].Summary != want {
			t.Fatalf("expected forwarded event %d to be %q, got %#v", index, want, next.events)
		}
	}
}

func TestSourceFilterSinkPropagatesWrappedSinkErrors(t *testing.T) {
	next := &recordingSink{err: errors.New("writer failed")}
	sink := NewSourceFilterSink(next, SourceAgent)

	droppedErr := sink.Publish(context.Background(), RunEvent{Source: SourceAgent, Kind: KindAgentRaw})
	forwardedErr := sink.Publish(context.Background(), RunEvent{Source: SourceDaemon, Kind: KindDaemonOutcome})

	if droppedErr != nil {
		t.Fatalf("expected dropped event to skip wrapped sink error, got %v", droppedErr)
	}
	if !errors.Is(forwardedErr, next.err) {
		t.Fatalf("expected forwarded event to propagate wrapped sink error, got %v", forwardedErr)
	}
	if len(next.events) != 1 || next.events[0].Source != SourceDaemon {
		t.Fatalf("expected only the daemon event forwarded, got %#v", next.events)
	}
}
