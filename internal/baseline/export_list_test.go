package baseline

import (
	"bytes"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestPlanExportList(t *testing.T) {
	t.Parallel()
	facts := ExportListFacts{GlossaryFiles: []string{"CONTEXT.md"}, KnowledgeADRs: []string{"docs/adr/0007-rule.md"}, ProcessPaths: []string{"AGENTS.md", "CLAUDE.md", "docs/agents/", "docs/specs/", "docs/history/", "docs/findings/", "docs/backlog/", "docs/handoffs/"}, DocsEntries: []string{"docs/adr/", "docs/agents/", "docs/design/", "docs/history/", "docs/private/", "docs/specs/", "docs/user-guide/"}}
	standard := []string{"README.md", "CONTEXT.md", "docs/user-guide/", "docs/references/", "docs/design/", "docs/lessons/", "docs/adr/0007-rule.md"}
	block := ExportListBeginMarker + "\n# Kept by the Baseline: roundfix baseline update rewrites the lines between these two markers.\n" + strings.Join(standard, "\n") + "\n" + ExportListEndMarker + "\n"
	t.Run("create", func(t *testing.T) {
		p := PlanExportList(nil, false, facts)
		if p.Action != "create" || string(p.Content) != block || !reflect.DeepEqual(p.Standard, standard) || !reflect.DeepEqual(p.Added, standard) {
			t.Fatalf("%+v", p)
		}
	})
	t.Run("convert_the_template", func(t *testing.T) {
		before := "# template\nREADME.md\ndocs/\nAGENTS.md\n!docs/history/**\n# Own\nnotes.md\n"
		p := PlanExportList([]byte(before), true, facts)
		if p.Action != "update" || string(p.Content) != block+"\n# template\n# Own\nnotes.md\n" || !reflect.DeepEqual(p.Converted, []string{"README.md", "docs/", "AGENTS.md", "!docs/history/**"}) {
			t.Fatalf("%+v\n%s", p, p.Content)
		}
	})
	t.Run("convert_reports_losses", func(t *testing.T) {
		p := PlanExportList([]byte("docs/\nAGENTS.md\n# Own\ndocs/private/\n"), true, facts)
		want := []ExportLoss{{Path: "AGENTS.md", Kind: "process"}, {Path: "docs/adr/", Kind: "other"}, {Path: "docs/agents/", Kind: "process"}, {Path: "docs/history/", Kind: "process"}, {Path: "docs/specs/", Kind: "process"}}
		if !reflect.DeepEqual(p.StopsExporting, want) {
			t.Fatalf("losses = %+v, want %+v", p.StopsExporting, want)
		}
	})
	t.Run("replace_between_markers", func(t *testing.T) {
		p := PlanExportList([]byte("# prefix\r\n"+ExportListBeginMarker+"\nold.md\n"+ExportListEndMarker+"\n# own\r\n!docs/design/\r\n"), true, facts)
		if string(p.Content) != "# prefix\r\n"+block+"# own\r\n!docs/design/\r\n" || !reflect.DeepEqual(p.Removed, []string{"old.md"}) {
			t.Fatalf("%+v\n%q", p, p.Content)
		}
	})
	t.Run("unchanged", func(t *testing.T) {
		p := PlanExportList([]byte(block), true, facts)
		if p.Action != "unchanged" || len(p.Added)+len(p.Removed) != 0 {
			t.Fatalf("%+v", p)
		}
	})
	cases := []struct{ name, input, refusal string }{
		{"refuse_begin_missing", ExportListEndMarker + "\n", ExportListBeginMarker + " is missing"},
		{"refuse_end_missing", ExportListBeginMarker + "\n", ExportListEndMarker + " is missing"},
		{"refuse_begin_repeated", ExportListBeginMarker + "\n" + block, ExportListBeginMarker + " appears 2 times"},
		{"refuse_end_repeated", block + ExportListEndMarker + "\n", ExportListEndMarker + " appears 2 times"},
		{"refuse_end_before_begin", ExportListEndMarker + "\n" + ExportListBeginMarker + "\n", ExportListEndMarker + " comes before " + ExportListBeginMarker},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p := PlanExportList([]byte(tt.input), true, facts)
			if p.Action != "refused" || p.Refusal != tt.refusal || string(p.Content) != tt.input {
				t.Fatalf("%+v", p)
			}
		})
	}
	t.Run("same_bytes_same_plan", func(t *testing.T) {
		data := []byte("docs/\nAGENTS.md\n")
		original := bytes.Clone(data)
		a, b := PlanExportList(data, true, facts), PlanExportList(data, true, facts)
		if !reflect.DeepEqual(a, b) || !bytes.Equal(data, original) {
			t.Fatal("not pure")
		}
	})
	t.Run("never_a_process_path", func(t *testing.T) {
		f := facts
		f.GlossaryFiles = append(append([]string{}, facts.GlossaryFiles...), "AGENTS.md", "docs/specs/GLOSSARY.md")
		f.KnowledgeADRs = append(append([]string{}, facts.KnowledgeADRs...), "docs/history/old.md")
		p := PlanExportList(nil, false, f)
		if !reflect.DeepEqual(p.Standard, standard) {
			t.Fatalf("%+v", p.Standard)
		}
	})
	t.Run("own_lines_kept", func(t *testing.T) {
		own := "# own bytes\r\n!docs/design/\r\nnotes.md"
		p := PlanExportList([]byte(block+own), true, facts)
		if string(p.Content) != block+own {
			t.Fatalf("own bytes changed: %q", p.Content)
		}
	})
	t.Run("apply_then_plan_unchanged", func(t *testing.T) {
		p := PlanExportList([]byte("docs/\n# Own\nnotes.md\n"), true, facts)
		next := PlanExportList(p.Content, true, facts)
		if next.Action != "unchanged" || !bytes.Equal(next.Content, p.Content) || !slices.Equal(next.Standard, p.Standard) {
			t.Fatalf("%+v", next)
		}
	})
}
