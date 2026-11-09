package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUnreachableUnreadReadsTheSectionAsUnreachableDoes(t *testing.T) {
	t.Parallel()

	fixtures, err := filepath.Glob(filepath.Join("testdata", "unreachable", "*", "_prd.md"))
	if err != nil || len(fixtures) == 0 {
		t.Fatalf("fixtures = %v, %v", fixtures, err)
	}
	for _, path := range fixtures {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			t.Parallel()
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want, wantErr := Unreachable(filepath.Dir(path))
			got, gotErr := parseUnreachableDeclarations(content, path)
			unread, unreadErr := UnreachableUnread(content, path)
			if !reflect.DeepEqual(got, want) || !reflect.DeepEqual(gotErr, wantErr) || !reflect.DeepEqual(unreadErr, wantErr) {
				t.Fatalf("parse = %+v, %v; unread error = %v; Unreachable = %+v, %v", got, gotErr, unreadErr, want, wantErr)
			}
			if len(unread) != 0 {
				t.Fatalf("unread = %+v, want none", unread)
			}
		})
	}
	t.Run("numbered and prose", func(t *testing.T) {
		t.Parallel()
		content := []byte("# PRD\n\n## Unreachable Acceptance\n\n1. A release rerun\n   needs the registry\n\n- criterion: release\n  reason: network\n  satisfied-by: operator\nUnindented prose\n")
		got, err := UnreachableUnread(content, "_prd.md")
		want := []UnreachableUnreadLine{{Line: 5, Text: "1. A release rerun"}, {Line: 6, Text: "needs the registry"}, {Line: 11, Text: "Unindented prose"}}
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("unread = %+v, %v; want %+v", got, err, want)
		}
	})
}

func TestUnreachableUnreadSectionBoundaries(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		section string
		unread  []UnreachableUnreadLine
	}{
		{"none after comment", "<!-- comment\nignored\n-->\n\nNone. All criteria\nare reachable.\n", nil},
		{"none paragraph ends at blank", "None. Reachable.\n\nlater prose\n", []UnreachableUnreadLine{{Line: 4, Text: "later prose"}}},
		{"blank in comment ends none paragraph", "None. Reachable.\n<!-- comment\n\n-->\nlater prose\n", []UnreachableUnreadLine{{Line: 6, Text: "later prose"}}},
		{"none must be first", "intro\n\nNone.\n", []UnreachableUnreadLine{{Line: 2, Text: "intro"}, {Line: 4, Text: "None."}}},
		{"field without declaration", "  reason: network\n", []UnreachableUnreadLine{{Line: 2, Text: "reason: network"}}},
		{"space and tab continuations", "- criterion: release\n wrapped\n\ttabbed\n  reason: network\n  satisfied-by: operator\n", nil},
		{"comment fields do not open declaration", "<!-- comment\n- criterion: example\n  reason: network\n  satisfied-by: operator\n-->\n  orphan\n", []UnreachableUnreadLine{{Line: 7, Text: "orphan"}}},
		{"comment field cannot enable continuation", "- plain bullet\n<!-- comment\n  reason: example\n-->\n orphan\n", []UnreachableUnreadLine{{Line: 6, Text: "orphan"}}},
		{"unread after malformed declaration", "- plain bullet\n  orphan\n- criterion: release\n  reason: network\n  satisfied-by: operator\nlast prose\n", []UnreachableUnreadLine{{Line: 3, Text: "orphan"}, {Line: 7, Text: "last prose"}}},
		{"level two ends section", "unread\n## Next\noutside\n", []UnreachableUnreadLine{{Line: 2, Text: "unread"}}},
		{"level one ends section", "unread\n# Next\noutside\n", []UnreachableUnreadLine{{Line: 2, Text: "unread"}}},
		{"level three stays inside", "### Details\n", []UnreachableUnreadLine{{Line: 2, Text: "### Details"}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			content := []byte("## Unreachable Acceptance\n" + tt.section)
			got, gotErr := UnreachableUnread(content, "_prd.md")
			_, wantErr := parseUnreachableDeclarations(content, "_prd.md")
			if !reflect.DeepEqual(gotErr, wantErr) || len(got) != len(tt.unread) {
				t.Fatalf("unread = %+v, %v; want %+v, %v", got, gotErr, tt.unread, wantErr)
			}
			for index := range got {
				if got[index] != tt.unread[index] {
					t.Errorf("unread[%d] = %+v, want %+v", index, got[index], tt.unread[index])
				}
			}
		})
	}
}
