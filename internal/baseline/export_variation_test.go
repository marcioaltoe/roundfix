package baseline

import (
	"reflect"
	"testing"
)

func TestExportVariationDepartures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, content string
		want          []DocumentationDeparture
	}{
		{"process_path", "# reason\ndocs/agents/\ndocs/\n!docs/agents/*\n", []DocumentationDeparture{{Code: "DR-EXPORT-PROCESS-PATH", Path: "docs/agents/", Line: 2}, {Code: "DR-EXPORT-PROCESS-PATH", Path: "docs/", Line: 3}}},
		{"reason_missing", "custom.md\n", []DocumentationDeparture{{Code: "DR-EXPORT-REASON-MISSING", Path: "custom.md", Line: 1}}},
		{"comment_stands_for_following_lines", "# reason\ncustom.md\n!docs/design/\n", nil},
		{"marker_not_a_reason", ExportListBeginMarker + "\nREADME.md\n" + ExportListEndMarker + "\ncustom.md\n", []DocumentationDeparture{{Code: "DR-EXPORT-REASON-MISSING", Path: "custom.md", Line: 4}}},
		{"blank_line_breaks_the_reason", "# reason\n\ncustom.md\n", []DocumentationDeparture{{Code: "DR-EXPORT-REASON-MISSING", Path: "custom.md", Line: 3}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := ExportListFacts{ProcessPaths: []string{"AGENTS.md", "docs/agents/", "plans/"}}
			if got := ExportVariationDepartures([]byte(tc.content), f); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("departures=%+v want=%+v", got, tc.want)
			}
		})
	}
}
