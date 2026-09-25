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

func TestClassifyBacklogEntryRetiresEveryClosedStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   string
		evidence string
		want     Retirement
	}{
		{name: "declined keeps retiring without additional evidence", status: "declined", want: Retirement{Retired: true, Reason: "declined"}},
		{name: "done with Spec", status: "done", evidence: "spec: 0164-knowledge-lifecycle-and-capture\n", want: Retirement{Retired: true, Reason: "done"}},
		{name: "deprecated with reason", status: "deprecated", evidence: "reason: replaced by current guidance\n", want: Retirement{Retired: true, Reason: "deprecated"}},
		{name: "superseded with Spec", status: "superseded", evidence: "spec: 0164-knowledge-lifecycle-and-capture\n", want: Retirement{Retired: true, Reason: "superseded"}},
		{name: "closed with reason", status: "closed", evidence: "reason: no implementation remains\n", want: Retirement{Retired: true, Reason: "closed"}},
		{name: "cancelled with Spec", status: "cancelled", evidence: "spec: 0164-knowledge-lifecycle-and-capture\n", want: Retirement{Retired: true, Reason: "cancelled"}},
		{name: "terminal without evidence", status: "done", evidence: "spec: null\nreason: null\n", want: Retirement{}},
		{name: "terminal with blank evidence", status: "closed", evidence: "spec: '  '\nreason: '  '\n", want: Retirement{}},
		{name: "open never retires", status: "open", evidence: "reason: deliberately retained\n", want: Retirement{}},
		{name: "promoted never retires", status: "promoted", evidence: "spec: 0164-knowledge-lifecycle-and-capture\n", want: Retirement{}},
		{name: "unknown never retires", status: "blocked", evidence: "reason: unknown is not terminal\n", want: Retirement{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			content := "---\nstatus: " + test.status + "\n" + test.evidence + "---\n\n# Backlog Entry\n"
			if got := ClassifyBacklogEntry([]byte(content)); got != test.want {
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
