// Suite: Doctor documentation rules
// Invariant: read-only advisory checks retain independent failures and exact diagnostics.
// Boundary IN: runCLI, temporary repositories and homes, existing readiness fakes.
// Boundary OUT: network, real home, machine and provider probes.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
)

func TestDoctorReportsTheDocumentationRules(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"export_list_missing", "export_list_stale", "export_marker", "export_process_path", "export_reason_missing", "export_adr_unmarked", "export_unread", "docs_required_missing", "docs_outside_layout", "docs_adr_retired", "docs_adr_amendment_unlisted", "docs_lesson_form", "both_ok", "skipped_without_manifest", "export_skipped_secondbrain_off", "twenty_lines_then_more", "exit_unchanged", "named_by_the_repository_guide", "order_after_claude_link", "proposed_adrs", "legacy_proposed_adr", "docs_required_glossary_missing", "root_reference_document", "marker_mentioned_in_comment"} {
		t.Run(name, func(t *testing.T) {
			repository, home := t.TempDir(), t.TempDir()
			setCommandEnvironmentForTest(t, home, repository)
			gittest.MarkRepository(t, repository)
			write := func(path, content string) { t.Helper(); writeBaselinePlanTestFile(t, repository, path, content) }
			remove := func(path string) {
				t.Helper()
				if err := os.Remove(filepath.Join(repository, path)); err != nil {
					t.Fatal(err)
				}
			}
			write(baselineSetupManifestPath, `{"modules":["core","context-workflow","secondbrain"]}`)
			write("README.md", "# Readme\n")
			write("CONTEXT.md", "# Context\n")
			facts, err := gatherDocumentationFacts(repository)
			if err != nil {
				t.Fatal(err)
			}
			standard := string(baseline.PlanExportList(nil, false, facts.Export).Content)
			write(baseline.ExportListPath, standard)
			exportOK := "baseline-export: ok (the Standard Export List is current; 0 line(s) of the repository's own)"
			docsOK := "baseline-docs: ok (the documentation follows the docs-layout guide)"
			wantExport, wantDocs := exportOK, docsOK
			var wantLines []string
			exportWarn := func(n int) {
				wantExport = fmt.Sprintf("baseline-export: warn (%d departure(s) from the export rules; next: run roundfix baseline update, then edit the repository's own lines)", n)
			}
			docsWarn := func(n int) {
				wantDocs = fmt.Sprintf("baseline-docs: warn (%d departure(s) from the documentation rules; next: see each line; roundfix baseline update fills an Amendment List and removes a retired ADR)", n)
			}
			adrPath := "docs/adr/0001-one.md"
			switch name {
			case "marker_mentioned_in_comment":
				write(baseline.ExportListPath, "# Restore "+baseline.ExportListBeginMarker+" here\n")
				exportWarn(1)
				wantLines = []string{"DR-EXPORT-LIST-MISSING: .secondbrain-export has no Standard Export List"}
			case "export_list_missing":
				remove(baseline.ExportListPath)
				exportWarn(1)
				wantLines = []string{"DR-EXPORT-LIST-MISSING: .secondbrain-export has no Standard Export List"}
			case "export_list_stale":
				write(baseline.ExportListPath, strings.Replace(standard, "docs/lessons/\n", "", 1))
				exportWarn(1)
				wantLines = []string{"DR-EXPORT-LIST-STALE: the Standard Export List differs from what roundfix baseline update would write: 1 line(s) to add, 0 to remove"}
			case "export_marker":
				write(baseline.ExportListPath, strings.Replace(standard, baseline.ExportListEndMarker, "", 1))
				exportWarn(1)
				wantLines = []string{"DR-EXPORT-MARKER: # roundfix:standard-export:end is missing; restore the pair of markers or delete both"}
			case "export_process_path":
				write(baseline.ExportListPath, standard+"\n# Local process\ndocs/agents/\n")
				exportWarn(1)
				wantLines = []string{fmt.Sprintf("DR-EXPORT-PROCESS-PATH: .secondbrain-export:%d: docs/agents/ publishes process; remove the line", strings.Count(standard, "\n")+3)}
			case "export_reason_missing":
				write(baseline.ExportListPath, standard+"\ncustom.md\n")
				exportWarn(1)
				wantLines = []string{fmt.Sprintf("DR-EXPORT-REASON-MISSING: .secondbrain-export:%d: custom.md has no reason; add a comment line directly above it", strings.Count(standard, "\n")+2)}
			case "export_adr_unmarked":
				write(adrPath, "---\nstatus: accepted\n---\n# One\n")
				exportWarn(1)
				wantLines = []string{"DR-EXPORT-ADR-UNMARKED: 1 of 1 ADR(s) carry no Knowledge Mark and are not exported"}
			case "export_unread":
				remove(baseline.ExportListPath)
				if err := os.Mkdir(filepath.Join(repository, baseline.ExportListPath), 0755); err != nil {
					t.Fatal(err)
				}
				anchored, err := os.OpenRoot(repository)
				if err != nil {
					t.Fatal(err)
				}
				_, readErr := anchored.ReadFile(baseline.ExportListPath)
				if err := anchored.Close(); err != nil {
					t.Fatal(err)
				}
				if readErr == nil {
					t.Fatal("expected directory read error")
				}
				wantExport = "baseline-export: warn (DR-EXPORT-UNREAD: read export list: " + readErr.Error() + "; next: re-run roundfix doctor after fixing the Setup Manifest)"
			case "docs_required_missing":
				remove("README.md")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-REQUIRED-MISSING: README.md is missing"}
			case "docs_required_glossary_missing":
				remove("CONTEXT.md")
				write(baseline.ExportListPath, strings.Replace(standard, "CONTEXT.md\n", "", 1))
				docsWarn(1)
				wantLines = []string{"DR-DOCS-REQUIRED-MISSING: the glossary is missing"}
			case "docs_outside_layout":
				write("docs/notes/note.md", "note\n")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-OUTSIDE-LAYOUT: docs/notes/: a folder the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md"}
			case "root_reference_document":
				write("NOTES.md", "notes\n")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-OUTSIDE-LAYOUT: NOTES.md: a reference document at the root the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md"}
			case "docs_adr_retired":
				write(adrPath, "---\nstatus: superseded\nknowledge: true\n---\n# Old\n")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-ADR-RETIRED: docs/adr/0001-one.md: a superseded ADR is still in the tree"}
			case "docs_adr_amendment_unlisted":
				write(adrPath, "---\nstatus: accepted\nknowledge: false\namended_by: [ADR-0002]\n---\n# One\n\nAmended by ADR-0002 and ADR-0003 (2026-11-19).\n")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-ADR-AMENDMENT-UNLISTED: docs/adr/0001-one.md: the body names ADR-0003 and amended_by does not"}
			case "docs_lesson_form":
				write("docs/lessons/2026-11-19-one.md", doctorValidLessonWithoutTitle)
				docsWarn(1)
				wantLines = []string{"DR-DOCS-LESSON-FORM: docs/lessons/2026-11-19-one.md: the title line is missing"}
			case "skipped_without_manifest":
				remove(baselineSetupManifestPath)
				wantExport = "baseline-export: skipped (no Setup Manifest at docs/agents/setup-context.json)"
				wantDocs = "baseline-docs: skipped (no Setup Manifest at docs/agents/setup-context.json)"
			case "export_skipped_secondbrain_off":
				write(baselineSetupManifestPath, `{"modules":["core","context-workflow"]}`)
				remove(baseline.ExportListPath)
				wantExport = "baseline-export: skipped (the Secondbrain is turned off in the Setup Manifest)"
			case "twenty_lines_then_more":
				var additions strings.Builder
				for i := 0; i < 23; i++ {
					fmt.Fprintf(&additions, "docs/agents/%02d.md\n", i)
				}
				write(baseline.ExportListPath, standard+"\n"+additions.String())
				exportWarn(46)
				// Two interleaved codes each retain their own limit and remainder.
				for i := 0; i < 20; i++ {
					wantLines = append(wantLines, fmt.Sprintf("DR-EXPORT-PROCESS-PATH: .secondbrain-export:%d: docs/agents/%02d.md publishes process; remove the line", strings.Count(standard, "\n")+2+i, i), fmt.Sprintf("DR-EXPORT-REASON-MISSING: .secondbrain-export:%d: docs/agents/%02d.md has no reason; add a comment line directly above it", strings.Count(standard, "\n")+2+i, i))
				}
				wantLines = append(wantLines, "DR-EXPORT-PROCESS-PATH: 3 more", "DR-EXPORT-REASON-MISSING: 3 more")
			case "named_by_the_repository_guide":
				write("docs/notes/note.md", "note\n")
				write("NOTES.md", "notes\n")
				write(baseline.SpecificRepositoryPath, "Use docs/notes/ for notes and NOTES.md for the index.\n")
			case "proposed_adrs":
				write(adrPath, "---\nstatus: proposed\n---\n# Proposed\n\nAmended by ADR-0002 (2026-11-19).\n")
			case "legacy_proposed_adr":
				write(adrPath, "# Proposed\n\nStatus: proposed\n\nAmended by ADR-0002 (2026-11-19).\n")
			case "both_ok":
				write(baseline.ExportListPath, standard+"\n# Additional knowledge\ncustom.md\n!docs/design/\n")
				wantExport = "baseline-export: ok (the Standard Export List is current; 2 line(s) of the repository's own)"
				write("docs/lessons/2026-11-19-one.md", strings.Replace(doctorValidLessonWithoutTitle, "---\n- Context:", "---\n# Lesson\n- Context:", 1))
			}
			node, wantCode := CheckStatusOK, exitOK
			if name == "exit_unchanged" {
				node, wantCode = CheckStatusFailed, exitRunFailed
				write("docs/notes/note.md", "note\n")
				docsWarn(1)
				wantLines = []string{"DR-DOCS-OUTSIDE-LAYOUT: docs/notes/: a folder the docs-layout guide does not name; move its content or name it in docs/agents/specific-repository.md"}
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: node}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: repository, HomeDir: home}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			before, homeBefore := snapshotDoctorPath(t, repository), snapshotDoctorPath(t, home)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			if code != wantCode || stderr.Len() != 0 {
				t.Fatalf("exit=%d want=%d stderr=%s stdout=%s", code, wantCode, stderr.String(), stdout.String())
			}
			lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			for _, want := range append([]string{wantExport, wantDocs}, wantLines...) {
				if strings.HasPrefix(want, "DR-") {
					want = "  " + want
				}
				if !slices.Contains(lines, want) {
					t.Fatalf("missing whole line %q:\n%s", want, stdout.String())
				}
			}
			var documentationLines []string
			start := slices.Index(lines, wantExport)
			end := slices.Index(lines, wantDocs)
			for _, line := range lines {
				if strings.HasPrefix(line, "  DR-EXPORT-") || strings.HasPrefix(line, "  DR-DOCS-") {
					documentationLines = append(documentationLines, strings.TrimPrefix(line, "  "))
				}
			}
			// Order of interleaved codes is not part of the surface; their full set is.
			slices.Sort(documentationLines)
			slices.Sort(wantLines)
			if !slices.Equal(documentationLines, wantLines) {
				t.Fatalf("departures=%q want=%q", documentationLines, wantLines)
			}
			if start < 1 || end <= start || !strings.HasPrefix(lines[start-1], "baseline-claude-link:") {
				t.Fatalf("documentation check order:\n%s", stdout.String())
			}
			for _, line := range lines[start+1 : end] {
				if !strings.HasPrefix(line, "  DR-EXPORT-") {
					t.Fatalf("unexpected line between documentation checks: %q", line)
				}
			}
			// The second check follows all details of the first and precedes residue.
			residue := slices.Index(lines, "residue: ok (no process residue found)")
			if residue < end+2 || !strings.HasPrefix(lines[residue-1], "flow-permissions:") || strings.Count(stdout.String(), "baseline-export:") != 1 || strings.Count(stdout.String(), "baseline-docs:") != 1 {
				t.Fatalf("duplicate or unordered checks:\n%s", stdout.String())
			}
			if !reflect.DeepEqual(snapshotDoctorPath(t, repository), before) || !reflect.DeepEqual(snapshotDoctorPath(t, home), homeBefore) {
				t.Fatal("Doctor changed repository or home")
			}
		})
	}
}

const doctorValidLessonWithoutTitle = "---\ncreated: 2026-11-19\n---\n- Context: work\n- Tried: attempt\n- Went wrong: symptom\n- Cause: proof\n- Solved by: solution\n- Applies when: here\n- Does not apply when: elsewhere\n- Links: evidence\n"
