package spec

// Suite: ADR reference forms in the shared history planner.
// Invariant: accepted references resolve exactly; existing lineage is preserved.
// Boundary IN: temporary ADR files and isolated Git histories.
// Boundary OUT: removal paths, successor bytes, lineage and per-path refusals.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestHistoryADRSuccessorResolvesEveryAcceptedForm(t *testing.T) {
	t.Parallel()
	const successor = "docs/adr/0064-new.md"
	const removed = "docs/adr/0002-old.md"
	cases := []struct {
		name, reference, status string
		accepted                bool
	}{
		{"ADR identifier", "ADR-0064", "accepted", true},
		{"number", "0064", "accepted", true},
		{"stem", "0064-new", "accepted", true},
		{"filename", "0064-new.md", "accepted", true},
		{"repository path", "docs/adr/0064-new.md", "accepted", true},
		{"relative path", "../adr/0064-new.md", "accepted", true},
		{"missing slug", "0064-missing", "accepted", false},
		{"missing filename", "0064-missing.md", "accepted", false},
		{"missing path", "docs/adr/0064-missing.md", "accepted", false},
		{"proposed number", "0064", "proposed", false},
		{"proposed filename", "0064-new.md", "proposed", false},
		{"unknown form", "new", "accepted", false},
		{"malformed identifier", "ADR-0064-new", "accepted", false},
		{"path without extension", "docs/adr/0064-new", "accepted", false},
		{"list instead of reference", "[ADR-0064]", "accepted", false},
		{"map instead of reference", "{adr: ADR-0064}", "accepted", false},
		{"null reference", "null", "accepted", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			old := "---\nstatus: superseded\nsuperseded_by: " + tc.reference + "\n---\n# Old\n"
			next := "---\nstatus: " + tc.status + "\n---\n# New\n\nSupersedes ADR-0002.\n"
			writeFile(t, filepath.Join(root, removed), old)
			writeFile(t, filepath.Join(root, successor), next)
			if tc.reference != "ADR-0064" && tc.reference != "0064" {
				// An explicit name must never fall back to another accepted ADR
				// with the same number, even when that file sorts first.
				writeFile(t, filepath.Join(root, "docs/adr/0064-a-decoy.md"), strings.ReplaceAll(next, "status: "+tc.status, "status: accepted"))
			}
			got, content, err := historyADRSuccessor(root, []byte(old))
			if err != nil {
				t.Fatal(err)
			}
			if tc.accepted {
				if got != successor || string(content) != next {
					t.Fatalf("successor=%q content=%q, want %q and exact bytes", got, content, successor)
				}
			} else if got != "" || len(content) != 0 {
				t.Fatalf("unexpected successor=%q content=%q", got, content)
			}
			plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil {
				t.Fatal(err)
			}
			if tc.accepted {
				if !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Refused) != 0 || len(plan.Lineage) != 0 {
					t.Fatalf("accepted reference plan=%+v", plan)
				}
			} else {
				want := []HistoryRefusal{{Path: removed, Reason: historyADRRefusal}}
				if len(plan.Files) != 0 || !reflect.DeepEqual(plan.Refused, want) || len(plan.Lineage) != 0 {
					t.Fatalf("refused reference plan=%+v, want refusal=%+v", plan, want)
				}
			}
		})
	}
}

func TestHistoryADRSuccessorThatSupersedesInFrontMatterGetsNoLineage(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	forms := []struct {
		name, reference string
		names           bool
	}{
		{"ADR identifier", "ADR-0002", true},
		{"number", "0002", true},
		{"stem", "0002-old", true},
		{"filename", "0002-old.md", true},
		{"repository path", "docs/adr/0002-old.md", true},
		{"relative path", "../adr/0002-old.md", true},
		{"another ADR", "ADR-0099", false},
		{"another slug", "0002-unrelated", false},
		{"number prefix", "00020", false},
	}
	for _, form := range forms {
		for _, list := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/list=%t", form.name, list), func(t *testing.T) {
				root, _ := gitLocatorRepo(t)
				const removed = "docs/adr/0002-old.md"
				const successor = "docs/adr/0064-new.md"
				value := fmt.Sprintf("%q", form.reference)
				if list {
					value = "\n  - \"ADR-0099\"\n  - " + value
				}
				original := "---\nstatus: accepted\nsupersedes: " + value + "\n---\n# New\n\nKeep these bytes.  \n"
				writeFile(t, filepath.Join(root, removed), "---\nstatus: superseded\nsuperseded_by: ADR-0064\n---\n# Old\n")
				writeFile(t, filepath.Join(root, successor), original)
				legacyCommit(t, root, "docs: predecessor reference fixture")
				plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
				if err != nil || !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Refused) != 0 {
					t.Fatalf("plan=%+v err=%v", plan, err)
				}
				wantLines := 1
				if form.names {
					wantLines = 0
				}
				if len(plan.Lineage) != wantLines {
					t.Fatalf("lineage=%+v, want %d line(s)", plan.Lineage, wantLines)
				}
				citations, err := plan.LineageCitations(root)
				if err != nil || len(citations) != wantLines {
					t.Fatalf("citations=%+v err=%v, want %d", citations, err, wantLines)
				}
				if err := ApplyHistoryKind(root, plan); err != nil {
					t.Fatal(err)
				}
				want := original
				if !form.names {
					want += "\n" + plan.Lineage[0].Line + "\n"
				}
				after, err := os.ReadFile(filepath.Join(root, successor))
				if err != nil || string(after) != want {
					t.Fatalf("successor=%q want=%q err=%v", after, want, err)
				}
			})
		}
	}
}

func TestHistoryPlanRemovesTheFluxusShapeWithoutALineageLine(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	const removed = "_archived/adr/0002-old.md"
	const successor = "docs/adr/0064-new.md"
	original := "---\nstatus: accepted\nsupersedes: 0002-old\n---\n# New\n\nConsolidated decision.\n"
	writeFile(t, filepath.Join(root, removed), "---\nstatus: superseded\nsuperseded_by: 0064-new\n---\n# Old\n")
	writeFile(t, filepath.Join(root, successor), original)
	plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
	if err != nil || !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Refused) != 0 || len(plan.Lineage) != 0 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	citations, err := plan.LineageCitations(root)
	if err != nil || len(citations) != 0 {
		t.Fatalf("citations=%+v err=%v", citations, err)
	}
	if err := ApplyHistoryKind(root, plan); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, removed)); !os.IsNotExist(err) {
		t.Fatalf("retired ADR remains: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(root, successor))
	if err != nil || string(after) != original {
		t.Fatalf("successor=%q want=%q err=%v", after, original, err)
	}
}
