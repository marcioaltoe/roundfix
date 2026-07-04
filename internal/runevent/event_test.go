package runevent

import "testing"

func TestSpecRunKindsUseDaemonNamespace(t *testing.T) {
	if KindDaemonTask != "daemon.task" {
		t.Fatalf("expected daemon.task kind, got %q", KindDaemonTask)
	}
	if KindDaemonQA != "daemon.qa" {
		t.Fatalf("expected daemon.qa kind, got %q", KindDaemonQA)
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
