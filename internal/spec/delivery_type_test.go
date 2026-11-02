package spec

// Suite: PRD Delivery Types and their Archive Records.
// Boundary IN: authored frontmatter and disposable Spec folders.
// Boundary OUT: publication and release classification belong to the CLI tests.
// Invariant: an absent type preserves the original record bytes.

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPRDDeliveryTypeReadsTheFrontmatter(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, front, value string
		want               DeliveryType
		present, invalid   bool
	}{
		{name: "absent", front: "spec: 0001-example"},
		{name: "feat", front: "type: feat", want: DeliveryTypeFeat, present: true},
		{name: "fix", front: "type: fix", want: DeliveryTypeFix, present: true},
		{name: "perf", front: "type: perf", want: DeliveryTypePerf, present: true},
		{name: "refactor", front: "type: refactor", want: DeliveryTypeRefactor, present: true},
		{name: "docs", front: "type: docs", value: "docs", present: true, invalid: true},
		{name: "null", front: "type: null", value: "null", present: true, invalid: true},
		{name: "list", front: "type: [fix]", value: "[fix]", present: true, invalid: true},
		{name: "empty", front: "type: \"\"", value: "", present: true, invalid: true},
		{name: "number", front: "type: 1", value: "1", present: true, invalid: true},
		{name: "case", front: "type: Fix", value: "Fix", present: true, invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, present, err := ParsePRDDeliveryType([]byte("---\n" + tc.front + "\n---\n# Example\n\ntype: docs\n"))
			if got != tc.want || present != tc.present {
				t.Fatalf("type=%q present=%v, want %q %v", got, present, tc.want, tc.present)
			}
			if !tc.invalid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var unsupported *UnsupportedDeliveryTypeError
			if !errors.As(err, &unsupported) || unsupported.Value != tc.value || err.Error() != "PRD type \""+tc.value+"\" is not one of feat, fix, perf, refactor" {
				t.Fatalf("error=%v, want unsupported %q", err, tc.value)
			}
		})
	}
	t.Run("supported set", func(t *testing.T) {
		t.Parallel()
		want := []DeliveryType{DeliveryTypeFeat, DeliveryTypeFix, DeliveryTypePerf, DeliveryTypeRefactor}
		got := DeliveryTypes()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("types=%v, want %v", got, want)
		}
		got[0] = "changed"
		if !reflect.DeepEqual(DeliveryTypes(), want) {
			t.Fatal("supported set was mutated")
		}
	})
}

func TestArchiveRecordCarriesThePRDDeliveryType(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, front, want string
		invalid           bool
	}{
		{name: "fix", front: "type: fix\n", want: "fix"},
		{name: "absent"},
		{name: "unsupported", front: "type: docs\n", invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, "_prd.md"), "---\nspec: 0001-example\ncreated: 2026-11-01\n"+tc.front+"---\n\n# Example\n\nA useful outcome.\n")
			record, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: dir, Slug: "0001-example", Source: "docs/specs/0001-example", SourceRevision: strings.Repeat("a", 40), Archived: "2026-11-01"})
			if tc.invalid {
				var unsupported *UnsupportedDeliveryTypeError
				if !errors.As(err, &unsupported) || unsupported.Value != "docs" {
					t.Fatalf("error=%v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if record.Type != tc.want {
				t.Fatalf("type=%q, want %q", record.Type, tc.want)
			}
			record.DeliveryCommit = "abc123"
			content, err := RenderArchiveRecord(record)
			if err != nil {
				t.Fatal(err)
			}
			typeLine := ""
			if tc.want != "" {
				typeLine = "type: " + tc.want + "\n"
			}
			want := "---\nschema: roundfix/archive-record/v1\nspec: 0001-example\ntitle: Example\nstatus: archived\ncreated: \"2026-11-01\"\narchived: \"2026-11-01\"\ndisposition: no-qa\nsource: docs/specs/0001-example\nsource_revision: " + strings.Repeat("a", 40) + "\nqa_task: \"\"\nqa_report: \"\"\nqa_verdict: \"\"\nunproven: []\nadrs: []\nsources: []\nregeneration: []\npromoted: []\ndelivery_commit: abc123\n" + typeLine + "---\n\n# Example\n\nA useful outcome.\n"
			if string(content) != want {
				t.Fatalf("record bytes:\n%s\nwant:\n%s", content, want)
			}
			parsed, err := ParseArchiveRecord(content)
			if err != nil || parsed.Type != tc.want {
				t.Fatalf("parsed=%+v err=%v", parsed, err)
			}
		})
	}
}
