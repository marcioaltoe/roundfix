package suiteguardcontract

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSanctionedRegenerationsReadArchiveRecords(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, archivedSpecRoot)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	record := []byte("---\nschema: " + archiveRecordSchema + "\nregeneration:\n  - command: make fixture\n    outputs: [generated/fixture.json]\n---\n# Fixture\n")
	path := filepath.Join(directory, "0001-fixture.md")
	if err := os.WriteFile(path, record, 0o644); err != nil {
		t.Fatal(err)
	}
	writeCleanupRegenerationFile(t, root, "docs/specs/current/_authorization.md", cleanupSpecGrantWithOutput("current", "make active", "generated/active.json"))
	writeCleanupRegenerationFile(t, root, "docs/history/specs/legacy/_authorization.md", cleanupSpecGrantWithOutput("legacy", "make legacy", "generated/legacy.json"))
	got, err := readSanctionedRegenerations(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []SanctionedRegeneration{
		{Command: "make active", Outputs: []string{"generated/active.json"}},
		{Command: "make fixture", Outputs: []string{"generated/fixture.json"}},
		{Command: "make legacy", Outputs: []string{"generated/legacy.json"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("regenerations = %+v, want %+v", got, want)
	}
	if err := os.WriteFile(path, []byte("---\nschema: unknown\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readSanctionedRegenerations(root); err == nil {
		t.Fatal("unknown schema accepted")
	}
}
