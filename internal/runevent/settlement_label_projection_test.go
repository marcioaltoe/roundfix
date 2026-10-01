package runevent

import (
	"encoding/json"
	"testing"
)

// Suite: Settlement Check labels in the public verification stream.
// Invariant: only the two exact Roundfix labels disclose a command in ordinary verification events.
// Boundary IN: daemon verification payloads and the verification category filter.
// Boundary OUT: journal persistence and CLI replay.
func TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot(t *testing.T) {
	t.Parallel()

	filter := StreamCategoryFilter{StreamCategoryVerification: {}}
	commands := []struct {
		command string
		want    string
	}{
		{"settlement check: spec consistency", "settlement check: spec consistency"},
		{"settlement check: authorization", "settlement check: authorization"},
		{"make verify", ""},
		{"settlement check: authorization; make verify", ""},
		{" settlement check: spec consistency", ""},
	}
	for _, command := range commands {
		for _, phase := range []VerificationPhase{VerificationPhaseStarted, VerificationPhaseCommandPassed, VerificationPhaseFailed} {
			t.Run(command.command+"/"+string(phase), func(t *testing.T) {
				payload, err := json.Marshal(map[string]any{
					"attempt": 1,
					"phase":   phase,
					"command": command.command,
					"task":    "task_07",
				})
				if err != nil {
					t.Fatalf("marshal verification payload: %v", err)
				}
				record, ok, err := ProjectStreamEvent(1, RunEvent{
					RunID:   "run_settlement",
					Source:  SourceDaemon,
					Kind:    KindDaemonVerification,
					Payload: payload,
				}, filter)
				if err != nil || !ok {
					t.Fatalf("project verification event: selected=%t, error=%v", ok, err)
				}
				if record.Category != StreamCategoryVerification || record.Phase != string(phase) || record.WorkItem != "task_07" {
					t.Fatalf("verification record = %#v", record)
				}
				encoded, err := json.Marshal(record)
				if err != nil {
					t.Fatalf("marshal stream record: %v", err)
				}
				var public map[string]json.RawMessage
				if err := json.Unmarshal(encoded, &public); err != nil {
					t.Fatalf("decode stream record: %v", err)
				}
				value, present := public["command"]
				if command.want == "" {
					if present {
						t.Fatalf("redacted command appears in public record: %s", encoded)
					}
					return
				}
				var got string
				if err := json.Unmarshal(value, &got); err != nil {
					t.Fatalf("decode projected command: %v", err)
				}
				if got != command.want {
					t.Fatalf("public command = %q, want %q", got, command.want)
				}
			})
		}
	}
}
