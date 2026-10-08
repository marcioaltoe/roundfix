package spec

// Suite: legacy unproven metadata.
// Invariant: legacy maps become stable text; active Specs remain string-only.
// Boundary IN: synthetic PRDs in temporary archive folders.
// Boundary OUT: archive records, refusal reasons, and preserved folder bytes.

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const legacyUnprovenMaps = `unproven:
  - satisfied-by: [task_04, task_05]
    claim: Catalog verification runs alone
    goal: G2
    row: 03
  - reason: No standalone evidence
    claim: Release is verified
    row: 04
  - "  Plain string item  "
  - claim: |
      A multi-line
      claim stays readable
    row: 05
  - zeta: last
    alpha: first
`

func legacyUnprovenFixture(t *testing.T, unproven string) LegacyConversionRequest {
	t.Helper()
	req := legacyProjectionFixture(t, "| task_01 | Implementation | backend | low | - |\n| task_qa | QA | qa | low | task_01 |\n")
	writeFile(t, filepath.Join(req.ArchiveRoot, req.Slug, "_prd.md"), "---\nspec: demo\nstatus: archived\ncreated: 2026-09-01\narchived: 2026-09-02\n"+unproven+"---\n\n# Demo\n\nPreserved outcome.\n")
	return req
}

func TestLegacyUnprovenMapsBecomeOneLineEach(t *testing.T) {
	req := legacyUnprovenFixture(t, legacyUnprovenMaps)
	want := []string{
		"row 03: Catalog verification runs alone (goal: G2; satisfied-by: task_04, task_05)",
		"row 04: Release is verified (reason: No standalone evidence)",
		"  Plain string item  ",
		"row 05: A multi-line claim stays readable",
		"alpha: first; zeta: last",
	}
	before := legacyTree(t, req)
	c := legacyPlan(t, req)
	legacyAssertWhole(t, req, before)
	if !reflect.DeepEqual(c.Record.Unproven, want) {
		t.Fatalf("unproven=%q want=%q", c.Record.Unproven, want)
	}
	archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil || !reflect.DeepEqual(archived.Record.Unproven, want) {
		t.Fatalf("legacy folder=%+v err=%v", archived, err)
	}
	if err := ApplyLegacyConversion(req.RepositoryRoot, c); err != nil {
		t.Fatal(err)
	}
	archived, err = ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil || archived.Form != ArchivedRecord || !reflect.DeepEqual(archived.Record.Unproven, want) {
		t.Fatalf("archive=%+v err=%v", archived, err)
	}
}

func TestLegacyUnprovenRefusesANestedMap(t *testing.T) {
	for _, tc := range []struct{ name, item string }{
		{"empty", "{}"},
		{"nested", "{claim: {text: unsupported}}"},
		{"sequence holds map", "{claim: [plain, {text: unsupported}]}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := legacyUnprovenFixture(t, "unproven:\n  - Preserved string\n  - "+tc.item+"\n")
			before := legacyTree(t, req)
			_, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: filepath.Join(req.ArchiveRoot, req.Slug), Legacy: true})
			if err == nil || err.Error() != "legacy unproven item 2 cannot be read as text" {
				t.Fatalf("builder error=%v", err)
			}
			_, err = PlanLegacyConversion(req)
			if err == nil || err.Error() != "build legacy record: legacy unproven item 2 cannot be read as text" {
				t.Fatalf("plan error=%v", err)
			}
			legacyAssertWhole(t, req, before)
		})
	}
}

func TestActiveSpecStillRefusesUnprovenMaps(t *testing.T) {
	req := legacyUnprovenFixture(t, legacyUnprovenMaps)
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	writeFile(t, filepath.Join(folder, "_prd.md"), "---\nspec: demo\nstatus: active\n"+legacyUnprovenMaps+"---\n\n# Demo\n")
	_, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder})
	if err == nil || !strings.HasPrefix(err.Error(), "parse archive PRD: yaml: unmarshal errors:") || !strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
		t.Fatalf("active error=%v", err)
	}
}
