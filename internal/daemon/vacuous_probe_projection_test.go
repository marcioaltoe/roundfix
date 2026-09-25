package daemon

import (
	"encoding/json"
	"reflect"
	"testing"

	"roundfix/internal/runevent"
)

// Suite: vacuous pre-work Verification publication
// Invariant: the publisher's payload projects the vacuous commands in probe order.
// Boundary IN: Daemon pre-work publisher and public Run Event Stream projection.
// Boundary OUT: Run Event Journal persistence and CLI replay.

func TestVacuousPreWorkEventProjectsItsCommands(t *testing.T) {
	t.Parallel()
	fixture := captureVacuousPreWorkEvent(t)
	var payload struct {
		Commands []string `json:"commands"`
	}
	if err := json.Unmarshal(fixture.event.Payload, &payload); err != nil {
		t.Fatalf("decode published vacuous event: %v", err)
	}
	if !reflect.DeepEqual(payload.Commands, fixture.offenders) {
		t.Fatalf("published commands = %q, want %q", payload.Commands, fixture.offenders)
	}

	record, ok, err := runevent.ProjectStreamEvent(1, fixture.event, runevent.AllStreamCategories())
	if err != nil {
		t.Fatalf("project published vacuous event: %v", err)
	}
	if !ok {
		t.Fatal("published vacuous event did not project")
	}
	if record.Classification != string(runevent.VerificationClassificationVacuous) {
		t.Fatalf("classification = %q, want %q", record.Classification, runevent.VerificationClassificationVacuous)
	}
	if !reflect.DeepEqual(record.Commands, fixture.offenders) {
		t.Fatalf("projected commands = %q, want %q", record.Commands, fixture.offenders)
	}
}
