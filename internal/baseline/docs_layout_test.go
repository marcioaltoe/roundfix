package baseline

import (
	"reflect"
	"testing"
)

func TestDocsLayoutDepartures(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"named_folders", "spec_and_history_roots", "unnamed_folder", "loose_file", "docs_readme_not_reported", "root_reference_document", "root_names_the_layout_names", "lower_case_root_file", "named_by_the_repository_guide", "readme_missing", "glossary_missing"} {
		t.Run(name, func(t *testing.T) {
			f := DocsLayoutFacts{ReadmePresent: true, GlossaryFiles: []string{"CONTEXT.md"}}
			var want []DocumentationDeparture
			switch name {
			case "named_folders":
				for _, folder := range []string{"_inbox", "adr", "agents", "backlog", "design", "findings", "handoffs", "history", "lessons", "references", "user-guide", "workflow"} {
					f.DocsEntries = append(f.DocsEntries, "docs/"+folder+"/")
				}
			case "spec_and_history_roots":
				f.SpecRoot = "docs/plans/specs/"
				f.History = "docs/past/archive/"
				f.DocsEntries = []string{"docs/plans/", "docs/past/"}
			case "unnamed_folder":
				f.DocsEntries = []string{"docs/notes/"}
				want = []DocumentationDeparture{{Code: "DR-DOCS-OUTSIDE-LAYOUT", Path: "docs/notes/", Detail: "a folder"}}
			case "loose_file":
				f.DocsEntries = []string{"docs/notes.md"}
				want = []DocumentationDeparture{{Code: "DR-DOCS-OUTSIDE-LAYOUT", Path: "docs/notes.md", Detail: "a file directly under docs/"}}
			case "docs_readme_not_reported":
				f.DocsEntries = []string{"docs/README.md"}
			case "root_reference_document":
				f.RootDocuments = []string{"NOTES.md"}
				want = []DocumentationDeparture{{Code: "DR-DOCS-OUTSIDE-LAYOUT", Path: "NOTES.md", Detail: "a reference document at the root"}}
			case "root_names_the_layout_names":
				f.RootDocuments = []string{"README.md", "AGENTS.md", "CLAUDE.md", "CONTEXT.md", "CONTEXT-MAP.md", "GLOSSARY.md", "GLOSSARY-MAP.md", "CHANGELOG.md", "CONTRIBUTING.md", "LICENSE.md", "LICENCE.md", "SECURITY.md", "CODE_OF_CONDUCT.md", "SUPPORT.md", "DESIGN.md"}
			case "lower_case_root_file":
				f.RootDocuments = []string{"notes.md", "Mixed.md"}
			case "named_by_the_repository_guide":
				f.DocsEntries = []string{"docs/notes/"}
				f.RootDocuments = []string{"NOTES.md"}
				f.SpecificRepository = []byte("Use `docs/notes/` for notes and `NOTES.md` for the index.")
			case "readme_missing":
				f.ReadmePresent = false
				want = []DocumentationDeparture{{Code: "DR-DOCS-REQUIRED-MISSING", Path: "README.md", Detail: "README.md is missing"}}
			case "glossary_missing":
				f.GlossaryFiles = nil
				want = []DocumentationDeparture{{Code: "DR-DOCS-REQUIRED-MISSING", Detail: "the glossary is missing"}}
			}
			if got := DocsLayoutDepartures(f); !reflect.DeepEqual(got, want) {
				t.Fatalf("departures=%+v want=%+v", got, want)
			}
		})
	}
}
