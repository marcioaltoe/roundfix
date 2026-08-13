package spec

// Suite: document retirement classification
// Invariant: only terminal decision and intent statuses retire documents, and the retiring status is reported.
// Boundary IN: in-memory decision records and typed intent entries.
// Boundary OUT: filesystem discovery and relocation.

import "testing"

func TestClassifyADRRetirement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    Retirement
	}{
		{name: "proposed", content: retirementDocument("proposed"), want: Retirement{}},
		{name: "accepted", content: retirementDocument("accepted"), want: Retirement{}},
		{name: "rejected", content: retirementDocument("rejected"), want: Retirement{Retired: true, Reason: "rejected"}},
		{name: "deprecated", content: retirementDocument("deprecated"), want: Retirement{Retired: true, Reason: "deprecated"}},
		{name: "superseded", content: retirementDocument("superseded"), want: Retirement{Retired: true, Reason: "superseded"}},
		{name: "legacy-no-status", content: "# Legacy decision\n\nThis decision remains active.\n", want: Retirement{}},
		{name: "legacy-body-marked", content: "# Legacy decision\n\nStatus: Deprecated\n", want: Retirement{Retired: true, Reason: "deprecated"}},
		{name: "legacy-fenced-status-deep", content: legacyFencedStatusBeyondHeader, want: Retirement{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ClassifyADR([]byte(test.content)); got != test.want {
				t.Fatalf("ClassifyADR() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestClassifyBacklogEntryRetirement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
		want   Retirement
	}{
		{name: "open", status: "open", want: Retirement{}},
		{name: "promoted", status: "promoted", want: Retirement{}},
		{name: "declined", status: "declined", want: Retirement{Retired: true, Reason: "declined"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ClassifyBacklogEntry([]byte(retirementDocument(test.status))); got != test.want {
				t.Fatalf("ClassifyBacklogEntry() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func retirementDocument(status string) string {
	return "---\nstatus: " + status + "\n---\n\n# Document\n"
}

// legacyFencedStatusBeyondHeader carries a fenced `Status: Deprecated` line
// deep in the body, past the leading header lines that the legacy fallback
// scans. The marker must not retire the record.
var legacyFencedStatusBeyondHeader = "# Legacy decision\n\n" +
	"## Context\n" +
	"\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"Long narrative line.\n" +
	"```text\n" +
	"Status: Deprecated\n" +
	"```\n"
