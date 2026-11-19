package spec

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestReadADRHeader(t *testing.T) {
	t.Parallel()
	yes, no := true, false
	tests := []struct {
		name, content string
		front         bool
		status        string
		mark          *bool
		ids           []string
	}{
		{"mark_true", "---\nstatus: accepted\nknowledge: true\n---\n", true, "accepted", &yes, []string{}},
		{"mark_false", "---\nstatus: accepted\nknowledge: false\n---\n", true, "accepted", &no, []string{}},
		{"mark_absent", "---\nstatus: accepted\n---\n", true, "accepted", nil, []string{}},
		{"mark_not_boolean", "---\nknowledge: 'true'\n---\n", true, "", nil, []string{}},
		{"amended_by_listed", "---\namended_by: [ADR-0042, ADR-0002]\n---\n", true, "", nil, []string{"ADR-0042", "ADR-0002"}},
		{"no_frontmatter", "# Record\nStatus: accepted\n", false, "", nil, []string{}},
		{"fields_anywhere_in_header", "---\namended_by: [ADR-0003]\ncustom: kept\nknowledge: true\nstatus: ACCEPTED\n---\n", true, "accepted", &yes, []string{"ADR-0003"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadADRHeader([]byte(tt.content))
			want := ADRHeader{Frontmatter: tt.front, Status: tt.status, Knowledge: tt.mark, AmendedBy: tt.ids}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("header = %+v, want %+v", got, want)
			}
		})
	}
}

func TestADRAmendmentNotices(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, content string
		want          []string
	}{
		{"amended_by_one", "Amended by ADR-0003", []string{"ADR-0003"}},
		{"amended_by_two", "Amended by ADR-0003, ADR-0002", []string{"ADR-0002", "ADR-0003"}},
		{"narrowed_by", "Narrowed by ADR-0003", []string{"ADR-0003"}},
		{"extended_by", "Extended by ADR-0003", []string{"ADR-0003"}},
		{"superseded_in_part_by", "Superseded in part by ADR-0003", []string{"ADR-0003"}},
		{"linked_identifier", "**Amended by [ADR-0003](0003-ADR-0999.md).**", []string{"ADR-0003"}},
		{"joined_by_and", "Extended by [ADR-0003](three.md) and ADR-0002", []string{"ADR-0002", "ADR-0003"}},
		{"fenced_block_ignored", "```md\nAmended by ADR-0003\n```\n~~~\nNarrowed by ADR-0004\n~~~\n", []string{}},
		{"inline_code_ignored", "`Amended by ADR-0003` and ``Extended by ADR-0004``", []string{}},
		{"ascending_once", "Extended by ADR-10, ADR-2 and ADR-10.\nAmended by ADR-2.", []string{"ADR-2", "ADR-10"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ADRAmendmentNotices([]byte(tt.content)); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("notices = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithADRAmendedBy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, before, after string
		ids                 []string
	}{
		{"after_knowledge", "---\nstatus: accepted\nknowledge: true\ncustom: keep\n---\nBody\n", "---\nstatus: accepted\nknowledge: true\namended_by: [ADR-0003]\ncustom: keep\n---\nBody\n", []string{"ADR-0003"}},
		{"after_status", "---\nstatus: accepted\ncustom: keep\n---\n", "---\nstatus: accepted\namended_by: [ADR-0003]\ncustom: keep\n---\n", []string{"ADR-0003"}},
		{"replaces_the_list", "---\namended_by: [ADR-0003, ADR-0002]\nstatus: accepted\n---\n", "---\namended_by: [ADR-0002, ADR-0003]\nstatus: accepted\n---\n", []string{"ADR-0002"}},
		{"keeps_written_identifiers", "---\namended_by: [ADR-0007]\n---\n", "---\namended_by: [ADR-0003, ADR-0007]\n---\n", []string{"ADR-0003"}},
		{"ascending_by_number", "# Legacy\n", "---\namended_by: [ADR-2, ADR-10, ADR-100]\n---\n# Legacy\n", []string{"ADR-100", "ADR-10", "ADR-2"}},
		{"keeps_updated_at", "---\nstatus: accepted\nupdated_at: 2026-01-02T03:04:05Z # keep\n---\nBody without newline", "---\nstatus: accepted\namended_by: [ADR-0003]\nupdated_at: 2026-01-02T03:04:05Z # keep\n---\nBody without newline", []string{"ADR-0003"}},
		{"crlf", "---\r\nstatus: accepted\r\n---\r\nBody\r\n", "---\r\nstatus: accepted\r\namended_by: [ADR-0003]\r\n---\r\nBody\r\n", []string{"ADR-0003"}},
		{"no_frontmatter_gains_header", "# Legacy\nStatus: accepted\n", "---\namended_by: [ADR-0003]\n---\n# Legacy\nStatus: accepted\n", []string{"ADR-0003"}},
		{"same_bytes_twice", "---\nstatus: accepted\namended_by: [ADR-0003]\n---\n", "---\nstatus: accepted\namended_by: [ADR-0003]\n---\n", []string{"ADR-0003"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := []byte(tt.before)
			before := bytes.Clone(original)
			got := WithADRAmendedBy(original, tt.ids)
			if string(got) != tt.after {
				t.Fatalf("got %q, want %q", got, tt.after)
			}
			if !bytes.Equal(original, before) {
				t.Fatal("mutated input")
			}
			if twice := WithADRAmendedBy(got, tt.ids); !bytes.Equal(twice, got) {
				t.Fatalf("second write differs: %q", twice)
			}
			if strings.Contains(tt.before, "Status: accepted") && ClassifyADR(got).Retired {
				t.Fatal("legacy metadata changed lifecycle")
			}
		})
	}
}
