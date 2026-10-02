// Suite: deterministic cause signatures.
// Invariant: each signature sees only its source, and the first matching entry wins.
package runcause

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoadRefusesABrokenTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*Table)
	}{
		{"class", func(t *Table) { t.Signatures[0].Class = "invented" }},
		{"duplicate signature", func(t *Table) { t.Signatures[1].ID = t.Signatures[0].ID }},
		{"duplicate trigger", func(t *Table) { t.Triggers[1].ID = t.Triggers[0].ID }},
		{"source", func(t *Table) { t.Signatures[0].Source = "invented" }},
		{"signature regexp", func(t *Table) { t.Signatures[0].Pattern = "[" }},
		{"trigger regexp", func(t *Table) { t.Triggers[0].Pattern = "[" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var table Table
			if err := json.Unmarshal(signatureBytes, &table); err != nil {
				t.Fatal(err)
			}
			test.mutate(&table)
			data, err := json.Marshal(table)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := load(data); err == nil {
				t.Fatal("broken table accepted")
			}
		})
	}
	table, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if table.SHA256 != fmt.Sprintf("%x", sha256.Sum256(signatureBytes)) {
		t.Fatal("digest differs from embedded bytes")
	}
}

func TestClassifyAppliesEachSignatureInOrder(t *testing.T) {
	t.Parallel()
	table, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct{ id, class, source, text string }{
		{"database-locked", "environment", "diagnostic", "SQLITE_BUSY"},
		{"disk-full", "environment", "diagnostic", "no space left on device"},
		{"network-unreachable", "environment", "diagnostic", "TLS handshake timeout"},
		{"killed-or-timed-out", "environment", "diagnostic", "panic: test timed out"},
		{"governed-path", "scope-or-authorization", "diagnostic", "undeclared path"},
		{"authorization-or-scope", "scope-or-authorization", "task", "outside its slice"},
		{"shared-section", "shared-section-contract", "command", "make skills-sync-check"},
		{"record-or-golden", "repository-convention", "diagnostic", "--- FAIL: TestCoverageEquivalence golden"},
		{"lint-or-format", "repository-convention", "diagnostic", "--- FAIL: TestLint: gofmt"},
		{"help-or-guide", "repository-convention", "task", "describe the command reference"},
		{"go-test-failure", "implementation-defect", "diagnostic", "--- FAIL: TestExample"},
		{"review-defect", "implementation-defect", "task", "review found two defects"},
	}
	if len(cases) != len(table.Signatures) {
		t.Fatal("signature coverage is not closed")
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			e := Evidence{}
			switch c.source {
			case "diagnostic":
				e.Diagnostic = c.text
			case "command":
				e.Command = c.text
			case "task":
				e.Task = c.text
			}
			class, id, source := table.Classify(e)
			if class != c.class || id != c.id || source != c.source {
				t.Fatalf("got %s %s %s", class, id, source)
			}
		})
	}
	for _, e := range []Evidence{{Command: "database is locked"}, {Task: "no space left on device"}, {Diagnostic: "outside its scope"}, {Command: "--- FAIL: TestExample"}, {Diagnostic: "help text"}, {Task: "unknown"}, {}} {
		class, id, source := table.Classify(e)
		if class != "unclassified" || id != "" || source != "" {
			t.Fatalf("cross-source/absent evidence classified: %+v -> %s %s %s", e, class, id, source)
		}
	}
	// any checks texts independently, diagnostic before command before Task.
	class, id, source := table.Classify(Evidence{Diagnostic: "golden", Command: "golden", Task: "golden"})
	if class != "repository-convention" || id != "record-or-golden" || source != "diagnostic" {
		t.Fatal("any source order changed")
	}
}

func TestTriggerNamesWhatAddedTheCorrectiveTask(t *testing.T) {
	t.Parallel()
	table, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ text, want string }{{"Pre-PR review found a gap", "pre-pr-review"}, {"QA finding F-001", "qa-gate"}, {"repository Verification failed", "verification"}, {"Pre-PR review and final QA and make verify", "pre-pr-review"}, {"unknown", "unknown"}} {
		if got := table.Trigger(c.text); got != c.want {
			t.Fatalf("%q -> %s, want %s", c.text, got, c.want)
		}
	}
}
